package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	context2 "github.com/Azure/azapi-lsp/internal/context"
	lsctx "github.com/Azure/azapi-lsp/internal/context"
	ilsp "github.com/Azure/azapi-lsp/internal/lsp"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/Azure/aztfmigrate/azurerm"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/Azure/aztfmigrate/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const tempFolderName = "aztfmigrate_temp"
const importFileName = "imports.tf"
const planFileName = "planfile"

type AztfMigrateCommand struct {
}

var _ CommandHandler = &AztfMigrateCommand{}

func (c AztfMigrateCommand) Handle(ctx context.Context, arguments []json.RawMessage) (interface{}, error) {
	var params lsp.CodeActionParams
	if len(arguments) != 0 {
		err := json.Unmarshal(arguments[0], &params)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
		}
	}

	clientCaller, err := context2.ClientCaller(ctx)
	if err != nil {
		return nil, err
	}

	clientNotifier, err := context2.ClientNotifier(ctx)
	if err != nil {
		return nil, err
	}

	reportProgress(ctx, "Parsing Terraform configurations...", 0)
	defer reportProgress(ctx, "Migration completed.", 100)

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := fs.GetDocument(ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI))
	if err != nil {
		return nil, err
	}

	// creating temp workspace
	workingDirectory := getWorkingDirectory(string(params.TextDocument.URI), runtime.GOOS)
	tempDir := filepath.Join(workingDirectory, tempFolderName)
	if err := os.MkdirAll(tempDir, 0750); err != nil {
		return nil, fmt.Errorf("creating temp workspace %q: %+v", tempDir, err)
	}
	defer func() {
		err := os.RemoveAll(path.Join(tempDir, "terraform.tfstate"))
		if err != nil {
			log.Printf("[ERROR] removing temp workspace %q: %+v", tempDir, err)
		}
	}()

	startDocPos := lsp.TextDocumentPositionParams{
		TextDocument: params.TextDocument,
		Position:     params.Range.Start,
	}
	startPos, err := ilsp.FilePositionFromDocumentPosition(startDocPos, doc)
	if err != nil {
		return nil, err
	}

	endDocPos := lsp.TextDocumentPositionParams{
		TextDocument: params.TextDocument,
		Position:     params.Range.End,
	}
	endPos, err := ilsp.FilePositionFromDocumentPosition(endDocPos, doc)
	if err != nil {
		return nil, err
	}

	data, err := doc.Text()
	if err != nil {
		return nil, err
	}

	// parsing the document
	syntaxDoc, diags := hclsyntax.ParseConfig(data, "", hcl.InitialPos)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing the HCL file: %s", diags.Error())
	}
	writeDoc, diags := hclwrite.ParseConfig(data, "", hcl.InitialPos)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parsing the HCL file: %s", diags.Error())
	}
	syntaxBlockMap := map[string]*hclsyntax.Block{}
	writeBlockMap := map[string]*hclwrite.Block{}

	body, ok := syntaxDoc.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("failed to parse HCL syntax")
	}
	addresses := make([]string, 0)
	for _, block := range body.Blocks {
		if startPos.Position().Byte <= block.Range().Start.Byte && block.Range().End.Byte <= endPos.Position().Byte {
			if block.Type != "resource" {
				continue
			}
			address := strings.Join(block.Labels, ".")
			addresses = append(addresses, address)
			syntaxBlockMap[address] = block
		}
	}
	for _, block := range writeDoc.Body().Blocks() {
		address := strings.Join(block.Labels(), ".")
		if _, ok := syntaxBlockMap[address]; ok {
			writeBlockMap[address] = block
		}
	}

	reportProgress(ctx, "Running Terraform plan...", 10)
	// running terraform plan
	terraform, err := tf.NewTerraform(workingDirectory, false)
	if err != nil {
		return nil, err
	}

	// AzureRM provider will honor env.var "AZURE_HTTP_USER_AGENT" when constructing for HTTP "User-Agent" header.
	// #nosec G104
	_ = os.Setenv("AZURE_HTTP_USER_AGENT", "aztfmigrate-vscode")
	// The following env.vars are used to disable enhanced validation and skip provider registration, to speed up the process.
	// #nosec G104
	_ = os.Setenv("ARM_PROVIDER_ENHANCED_VALIDATION", "false")
	// #nosec G104
	_ = os.Setenv("ARM_SKIP_PROVIDER_REGISTRATION", "true")

	options := make([]tfexec.PlanOption, 0)
	for address := range syntaxBlockMap {
		options = append(options, tfexec.Target(address))
	}
	planfile := path.Join(tempDir, planFileName)
	options = append(options, tfexec.Out(planfile))
	_, err = terraform.GetExec().Plan(ctx, options...)
	if err != nil {
		return nil, err
	}
	plan, err := terraform.GetExec().ShowPlanFile(ctx, planfile)
	if err != nil {
		return nil, err
	}

	allResources := types.ListResourcesFromPlan(plan)
	resources := make([]types.AzureResource, 0)
	for _, r := range allResources {
		if syntaxBlockMap[r.OldAddress(nil)] == nil {
			continue
		}

		switch resource := r.(type) {
		case *types.AzapiResource:
			if len(resource.Instances) == 0 {
				continue
			}
			resourceId := resource.Instances[0].ResourceId
			resourceTypes, exact, err := azurerm.GetAzureRMResourceType(resourceId)
			if err != nil {
				log.Printf("failed to get resource type for %s: %v", resourceId, err)
				continue
			}
			if len(resourceTypes) == 0 {
				log.Printf("failed to get resource type for %s", resourceId)
				continue
			}
			if !exact {
				log.Printf("multiple resource types found for %s: %v", resourceId, resourceTypes)
			}
			resource.ResourceType = resourceTypes[0]
			resources = append(resources, resource)

		case *types.AzapiUpdateResource:
			resourceTypes, exact, err := azurerm.GetAzureRMResourceType(resource.Id)
			if err != nil {
				log.Printf("failed to get resource type for %s: %v", resource.Id, err)
				continue
			}
			if len(resourceTypes) == 0 {
				log.Printf("failed to get resource type for %s", resource.Id)
				continue
			}
			if !exact {
				log.Printf("multiple resource types found for %s: %v", resource.Id, resourceTypes)
			}
			resource.ResourceType = resourceTypes[0]
			resources = append(resources, resource)

		case *types.AzurermResource:
			if len(resource.Instances) == 0 {
				continue
			}
			resources = append(resources, resource)
		}
	}

	tempTerraform, err := tf.NewTerraform(tempDir, false)
	if err != nil {
		return nil, err
	}

	if err = os.WriteFile(filepath.Join(tempDir, importFileName), []byte(importConfig(resources)), 0600); err != nil {
		return nil, err
	}

	for index, r := range resources {
		// #nosec G115
		reportProgress(ctx, fmt.Sprintf("Migrating resource %d/%d...", index+1, len(resources)), 40+uint32(50.0*index/len(resources)))
		if err := r.GenerateNewConfig(tempTerraform); err != nil {
			log.Printf("[ERROR] %+v", err)
			_ = clientNotifier.Notify(ctx, "window/showMessage", lsp.ShowMessageParams{
				Type:    lsp.Error,
				Message: "Failed to generate new config for " + r.OldAddress(nil),
			})
		}
	}

	reportProgress(ctx, "Updating Terraform configurations...", 90)
	// migrate depends_on, lifecycle, provisioner
	for _, r := range resources {
		existingBlock := writeBlockMap[r.OldAddress(nil)]
		if existingBlock == nil {
			continue
		}
		migratedBlock := r.MigratedBlock()
		if attr := existingBlock.Body().GetAttribute("depends_on"); attr != nil {
			migratedBlock.Body().SetAttributeRaw("depends_on", attr.Expr().BuildTokens(nil))
		}
		for _, block := range existingBlock.Body().Blocks() {
			if block.Type() == "lifecycle" || block.Type() == "provisioner" {
				migratedBlock.Body().AppendBlock(block)
			}
		}
	}

	// update config
	emptyFile := hclwrite.NewEmptyFile()
	outputs := make([]types.Output, 0)
	resourcesMap := make(map[string]types.AzureResource)
	for _, r := range resources {
		if r.IsMigrated() {
			outputs = append(outputs, r.Outputs()...)
		}
		resourcesMap[r.OldAddress(nil)] = r
	}

	for _, addr := range addresses {
		r := resourcesMap[addr]
		if r == nil {
			continue
		}
		if !r.IsMigrated() {
			emptyFile.Body().AppendBlock(writeBlockMap[r.OldAddress(nil)])
			emptyFile.Body().AppendNewline()
			continue
		}
		if writeBlockMap[r.OldAddress(nil)] == nil {
			continue
		}

		emptyFile.Body().AppendUnstructuredTokens(types.CommentOutBlock(writeBlockMap[r.OldAddress(nil)]))
		emptyFile.Body().AppendNewline()
		if removedBlock := r.RemovedBlock(); removedBlock != nil {
			emptyFile.Body().AppendBlock(removedBlock)
			emptyFile.Body().AppendNewline()
		}
		if importBlock := r.ImportBlock(); importBlock != nil {
			emptyFile.Body().AppendBlock(importBlock)
			emptyFile.Body().AppendNewline()
		}
		if migratedBlock := r.MigratedBlock(); migratedBlock != nil {
			types.ReplaceOutputs(migratedBlock, outputs)
			emptyFile.Body().AppendBlock(migratedBlock)
			emptyFile.Body().AppendNewline()
		}
	}

	_, _ = clientCaller.Callback(ctx, "workspace/applyEdit", lsp.ApplyWorkspaceEditParams{
		Label: "Update config",
		Edit: lsp.WorkspaceEdit{
			Changes: map[string][]lsp.TextEdit{
				string(params.TextDocument.URI): {
					{
						Range:   params.Range,
						NewText: string(hclwrite.Format(emptyFile.Bytes())),
					},
				},
			},
		},
	})

	return nil, nil
}

func importConfig(resources []types.AzureResource) string {
	const providerConfig = `
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
  subscription_id = "%s"
}

provider "azapi" {
}
`

	config := ""
	for _, r := range resources {
		config += r.EmptyImportConfig()
	}
	subscriptionId := ""
	for _, r := range resources {
		switch resource := r.(type) {
		case *types.AzapiResource:
			for _, instance := range resource.Instances {
				if strings.HasPrefix(instance.ResourceId, "/subscriptions/") {
					subscriptionId = strings.Split(instance.ResourceId, "/")[2]
					break
				}
			}
		case *types.AzapiUpdateResource:
			if strings.HasPrefix(resource.Id, "/subscriptions/") {
				subscriptionId = strings.Split(resource.Id, "/")[2]
			}
		}
		if subscriptionId != "" {
			break
		}
	}
	config = fmt.Sprintf(providerConfig, subscriptionId) + config

	return config
}

func reportProgress(ctx context.Context, message string, percentage uint32) {
	clientCaller, err := context2.ClientCaller(ctx)
	if err != nil {
		log.Printf("[ERROR] failed to get client caller: %+v", err)
		return
	}

	clientNotifier, err := context2.ClientNotifier(ctx)
	if err != nil {
		log.Printf("[ERROR] failed to get client notifier: %+v", err)
		return
	}
	switch percentage {
	case 0:
		_, _ = clientCaller.Callback(ctx, "window/workDoneProgress/create", lsp.WorkDoneProgressCreateParams{
			Token: "aztfmigrate",
		})
		_ = clientNotifier.Notify(ctx, "$/progress", map[string]interface{}{
			"token": "aztfmigrate",
			"value": lsp.WorkDoneProgressBegin{
				Kind:        "begin",
				Title:       "Azure providers migration",
				Cancellable: false,
				Message:     message,
				Percentage:  0,
			},
		})
	case 100:
		_ = clientNotifier.Notify(ctx, "$/progress", map[string]interface{}{
			"token": "aztfmigrate",
			"value": lsp.WorkDoneProgressEnd{
				Kind:    "end",
				Message: message,
			},
		})
	default:
		_ = clientNotifier.Notify(ctx, "$/progress", map[string]interface{}{
			"token": "aztfmigrate",
			"value": lsp.WorkDoneProgressReport{
				Kind:        "report",
				Cancellable: false,
				Message:     message,
				Percentage:  percentage,
			},
		})
	}
}

func getWorkingDirectory(uri string, os string) string {
	workingDirectory := path.Dir(strings.TrimPrefix(uri, "file://"))
	if os == "windows" {
		workingDirectory, _ = url.QueryUnescape(workingDirectory)
		workingDirectory = strings.ReplaceAll(workingDirectory, "/", "\\")
		workingDirectory = strings.TrimPrefix(workingDirectory, "\\")
	}
	return workingDirectory
}
