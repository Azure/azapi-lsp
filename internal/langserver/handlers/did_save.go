package handlers

import (
	"context"

	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/cmd"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/command"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func (lh *logHandler) TextDocumentDidSave(ctx context.Context, params lsp.DidSaveTextDocumentParams) error {
	expFeatures, err := lsctx.ExperimentalFeatures(ctx)
	if err != nil {
		return err
	}
	if !expFeatures.ValidateOnSave {
		return nil
	}

	fh := ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI)
	dh := ilsp.FileHandlerFromDirPath(fh.Dir())

	_, err = command.TerraformValidateHandler(ctx, cmd.CommandArgs{
		"uri": dh.URI(),
	})

	return err
}
