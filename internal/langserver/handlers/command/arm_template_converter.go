package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Context struct {
	File         *hclwrite.File
	typeLabelMap map[string]bool
	idRefMap     map[string]string
}

func NewContext() *Context {
	file, _ := hclwrite.ParseConfig([]byte(hclTemplate), "main.tf", hcl.InitialPos)
	return &Context{
		File:         file,
		typeLabelMap: make(map[string]bool),
		idRefMap:     make(map[string]string),
	}
}

func (c *Context) AppendBlock(block *hclwrite.Block) {
	if block.Type() == "resource" {
		typeValue := string(block.Body().GetAttribute("type").Expr().BuildTokens(nil).Bytes())
		typeValue = strings.Trim(typeValue, " \"")
		tfLabel := block.Labels()[1]
		if typeLabel := fmt.Sprintf("%s-%s", typeValue, tfLabel); c.typeLabelMap[typeLabel] {
			for i := 1; ; i++ {
				tfLabel = fmt.Sprintf("%s%d", block.Labels()[1], i)
				newTypeLabel := fmt.Sprintf("%s-%s", typeValue, tfLabel)
				if !c.typeLabelMap[newTypeLabel] {
					newLabels := []string{block.Labels()[0], tfLabel}
					block.SetLabels(newLabels)
					c.typeLabelMap[newTypeLabel] = true
					break
				}
			}
		} else {
			c.typeLabelMap[typeLabel] = true
		}
		nameValue := string(block.Body().GetAttribute("name").Expr().BuildTokens(nil).Bytes())
		nameValue = strings.Trim(nameValue, " \"")
		parentIdValue := string(block.Body().GetAttribute("parent_id").Expr().BuildTokens(nil).Bytes())
		parentIdValue = strings.Trim(parentIdValue, " \"")
		c.idRefMap[buildResourceId(nameValue, parentIdValue, typeValue)] = fmt.Sprintf("azapi_resource.%s.id", tfLabel)
	}

	c.File.Body().AppendBlock(block)
	c.File.Body().AppendNewline()
}

func (c *Context) String() string {
	result := string(c.File.Bytes())

	for id, ref := range c.idRefMap {
		result = strings.ReplaceAll(result, fmt.Sprintf(`"%s"`, id), ref)
	}

	result = string(hclwrite.Format([]byte(result)))
	// TODO: improve it
	result = strings.ReplaceAll(result, "$${", "${")
	return result
}

func convertARMTemplate(input string) (string, error) {
	var model ARMTemplateModel
	err := json.Unmarshal([]byte(input), &model)
	if err != nil {
		return "", err
	}

	c := NewContext()

	for key, parameter := range model.Parameters {
		varBlock := hclwrite.NewBlock("variable", []string{key})
		switch strings.ToLower(parameter.Type) {
		case "string":
			varBlock.Body().SetAttributeTraversal("type", hcl.Traversal{hcl.TraverseRoot{Name: "string"}})
		case "securestring":
			varBlock.Body().SetAttributeTraversal("type", hcl.Traversal{hcl.TraverseRoot{Name: "string"}})
			varBlock.Body().SetAttributeValue("sensitive", cty.True)
		default:
			// Todo: support other types
			varBlock.Body().SetAttributeTraversal("type", hcl.Traversal{hcl.TraverseRoot{Name: strings.ToLower(parameter.Type)}})
		}
		varBlock.Body().SetAttributeValue("default", cty.StringVal(parameter.DefaultValue))
		c.AppendBlock(varBlock)
	}

	for _, resource := range model.Resources {
		res := flattenARMExpression(resource)
		data, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return "", fmt.Errorf("unable to marshal JSON content: %v", err)
		}
		resourceJson := string(data)
		resourceBlock, err := ParseResourceJson(resourceJson)
		if err != nil {
			return "", fmt.Errorf("unable to parse resource JSON content: %v", err)
		}
		if resourceBlock == nil {
			return "", fmt.Errorf("resource block is nil")
		}
		c.AppendBlock(resourceBlock)
	}

	return c.String(), nil
}

const hclTemplate = `
variable "subscriptionId" {
  type    = string
  description = "The subscription id"
}

variable "resourceGroupName" {
  type    = string
  description = "The resource group name"
}

`

type ARMTemplateParameterModel struct {
	DefaultValue string `json:"defaultValue"`
	Type         string `json:"type"`
}

type ARMTemplateModel struct {
	Schema         string                               `json:"$schema"`
	ContentVersion string                               `json:"contentVersion"`
	Parameters     map[string]ARMTemplateParameterModel `json:"parameters"`
	Variables      interface{}                          `json:"variables"`
	Resources      []map[string]interface{}             `json:"resources"`
}

func buildResourceId(name, parentId, resourceType string) string {
	azureResourceType := resourceType[:strings.Index(resourceType, "@")]
	azureResourceId := ""
	switch {
	case strings.Count(azureResourceType, "/") == 1:
		// build azure resource id
		switch azureResourceType {
		case arm.ResourceGroupResourceType.String():
			azureResourceId = fmt.Sprintf("%s/resourceGroups/%s", parentId, name)
		case arm.SubscriptionResourceType.String():
			azureResourceId = fmt.Sprintf("/subscriptions/%s", name)
		case arm.TenantResourceType.String():
			azureResourceId = "/"
		case arm.ProviderResourceType.String():
			// avoid duplicated `/` if parent_id is tenant scope
			scopeId := parentId
			if parentId == "/" {
				scopeId = ""
			}
			azureResourceId = fmt.Sprintf("%s/providers/%s", scopeId, name)
		default:
			// avoid duplicated `/` if parent_id is tenant scope
			scopeId := parentId
			if parentId == "/" {
				scopeId = ""
			}
			azureResourceId = fmt.Sprintf("%s/providers/%s/%s", scopeId, azureResourceType, name)
		}
	default:
		// build azure resource id
		lastType := azureResourceType[strings.LastIndex(azureResourceType, "/")+1:]
		azureResourceId = fmt.Sprintf("%s/%s/%s", parentId, lastType, name)
	}

	return azureResourceId
}
