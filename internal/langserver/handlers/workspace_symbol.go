package handlers

import (
	"context"

	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func (svc *service) WorkspaceSymbol(ctx context.Context, params lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	cc, err := ilsp.ClientCapabilities(ctx)
	if err != nil {
		return nil, err
	}

	symbols, err := svc.decoder.Symbols(ctx, params.Query)
	if err != nil {
		return nil, err
	}

	return ilsp.WorkspaceSymbols(symbols, cc.Workspace.Symbol), nil
}
