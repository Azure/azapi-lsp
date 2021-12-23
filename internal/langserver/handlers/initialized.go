package handlers

import (
	"context"

	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func Initialized(ctx context.Context, params lsp.InitializedParams) error {
	return nil
}
