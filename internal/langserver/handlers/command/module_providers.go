package command

import (
	"context"
	"fmt"

	"github.com/creachadair/jrpc2/code"
	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/cmd"
	op "github.com/ms-henglu/azurerm-restapi-lsp/internal/terraform/module/operation"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/uri"
	tfaddr "github.com/hashicorp/terraform-registry-address"
)

const moduleProvidersVersion = 0

type moduleProvidersResponse struct {
	FormatVersion        int                            `json:"v"`
	ProviderRequirements map[string]providerRequirement `json:"provider_requirements"`
	InstalledProviders   map[string]string              `json:"installed_providers"`
}

type providerRequirement struct {
	DisplayName       string `json:"display_name"`
	VersionConstraint string `json:"version_constraint,omitempty"`
	DocsLink          string `json:"docs_link,omitempty"`
}

func ModuleProvidersHandler(ctx context.Context, args cmd.CommandArgs) (interface{}, error) {
	response := moduleProvidersResponse{
		FormatVersion:        moduleProvidersVersion,
		ProviderRequirements: make(map[string]providerRequirement),
		InstalledProviders:   make(map[string]string),
	}

	modUri, ok := args.GetString("uri")
	if !ok || modUri == "" {
		return response, fmt.Errorf("%w: expected module uri argument to be set", code.InvalidParams.Err())
	}

	if !uri.IsURIValid(modUri) {
		return response, fmt.Errorf("URI %q is not valid", modUri)
	}

	modPath, err := uri.PathFromURI(modUri)
	if err != nil {
		return response, err
	}

	mm, err := lsctx.ModuleFinder(ctx)
	if err != nil {
		return response, err
	}

	mod, _ := mm.ModuleByPath(modPath)
	if mod == nil {
		return response, nil
	}

	if mod.MetaState == op.OpStateUnknown || mod.MetaErr != nil {
		return response, nil
	}

	for provider, version := range mod.Meta.ProviderRequirements {
		response.ProviderRequirements[provider.String()] = providerRequirement{
			DisplayName:       provider.ForDisplay(),
			VersionConstraint: version.String(),
			DocsLink:          getProviderDocumentationLink(provider),
		}
	}

	for provider, version := range mod.InstalledProviders {
		response.InstalledProviders[provider.String()] = version.String()
	}

	return response, nil
}

func getProviderDocumentationLink(provider tfaddr.Provider) string {
	if provider.IsLegacy() || provider.IsBuiltIn() || provider.Hostname != "registry.terraform.io" {
		return ""
	}

	return fmt.Sprintf(`https://registry.terraform.io/providers/%s/latest`, provider.ForDisplay())
}
