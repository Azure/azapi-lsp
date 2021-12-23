package command

import (
	"context"
	"fmt"

	"github.com/creachadair/jrpc2/code"
	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/cmd"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/diagnostics"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/errors"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/progress"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/terraform/module"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/uri"
)

func TerraformValidateHandler(ctx context.Context, args cmd.CommandArgs) (interface{}, error) {
	dirUri, ok := args.GetString("uri")
	if !ok || dirUri == "" {
		return nil, fmt.Errorf("%w: expected module uri argument to be set", code.InvalidParams.Err())
	}

	if !uri.IsURIValid(dirUri) {
		return nil, fmt.Errorf("URI %q is not valid", dirUri)
	}

	dh := ilsp.FileHandlerFromDirURI(lsp.DocumentURI(dirUri))

	modMgr, err := lsctx.ModuleManager(ctx)
	if err != nil {
		return nil, err
	}

	mod, err := modMgr.ModuleByPath(dh.Dir())
	if err != nil {
		if module.IsModuleNotFound(err) {
			mod, err = modMgr.AddModule(dh.Dir())
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	tfExec, err := module.TerraformExecutorForModule(ctx, mod.Path)
	if err != nil {
		return nil, errors.EnrichTfExecError(err)
	}

	notifier, err := lsctx.DiagnosticsNotifier(ctx)
	if err != nil {
		return nil, err
	}

	progress.Begin(ctx, "Validating")
	defer func() {
		progress.End(ctx, "Finished")
	}()
	progress.Report(ctx, "Running terraform validate ...")
	jsonDiags, err := tfExec.Validate(ctx)
	if err != nil {
		return nil, err
	}

	diags := diagnostics.NewDiagnostics()
	validateDiags := diagnostics.HCLDiagsFromJSON(jsonDiags)
	diags.EmptyRootDiagnostic()
	diags.Append("terraform validate", validateDiags)
	diags.Append("HCL", mod.ModuleDiagnostics.AsMap())
	diags.Append("HCL", mod.VarsDiagnostics.AutoloadedOnly().AsMap())

	notifier.PublishHCLDiags(ctx, mod.Path, diags)

	return nil, nil
}
