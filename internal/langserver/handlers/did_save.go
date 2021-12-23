package handlers

import (
	"context"

	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func (lh *logHandler) TextDocumentDidSave(ctx context.Context, params lsp.DidSaveTextDocumentParams) error {
	return nil
}
