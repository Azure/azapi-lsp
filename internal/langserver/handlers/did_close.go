package handlers

import (
	"context"

	lsctx "github.com/Azure/azapi-lsp/internal/context"
	ilsp "github.com/Azure/azapi-lsp/internal/lsp"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

func TextDocumentDidClose(ctx context.Context, params lsp.DidCloseTextDocumentParams) error {
	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return err
	}

	fh := ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI)
	err = fs.CloseAndRemoveDocument(fh)
	if err != nil {
		return err
	}

	return nil
}
