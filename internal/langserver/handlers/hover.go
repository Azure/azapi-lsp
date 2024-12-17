package handlers

import (
	"context"

	lsctx "github.com/Azure/azapi-lsp/internal/context"
	"github.com/Azure/azapi-lsp/internal/langserver/handlers/hover"
	ilsp "github.com/Azure/azapi-lsp/internal/lsp"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

func (svc *service) TextDocumentHover(ctx context.Context, params lsp.TextDocumentPositionParams) (*lsp.Hover, error) {
	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return nil, err
	}

	_, err = ilsp.ClientCapabilities(ctx)
	if err != nil {
		return nil, err
	}

	doc, err := fs.GetDocument(ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI))
	if err != nil {
		return nil, err
	}

	fPos, err := ilsp.FilePositionFromDocumentPosition(params, doc)
	if err != nil {
		return nil, err
	}

	data, err := doc.Text()
	if err != nil {
		return nil, err
	}

	telemetrySender, err := lsctx.Telemetry(ctx)
	if err != nil {
		return nil, err
	}

	svc.logger.Printf("Looking for hover data at %q -> %#v", doc.Filename(), fPos.Position())
	hoverData := hover.HoverAtPos(ctx, data, doc.Filename(), fPos.Position(), svc.logger, telemetrySender)
	svc.logger.Printf("received hover data: %#v", hoverData)

	if hoverData != nil {
		telemetrySender.SendEvent(ctx, "textDocument/hover", nil)
	}

	return hoverData, nil
}
