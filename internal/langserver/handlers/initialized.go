package handlers

import (
	"context"

	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

func Initialized(ctx context.Context, params lsp.InitializedParams) error {
	return nil
}
