package handlers

import (
	"context"

	lsctx "github.com/Azure/azapi-lsp/internal/context"
	"github.com/Azure/azapi-lsp/internal/langserver/handlers/complete"
	ilsp "github.com/Azure/azapi-lsp/internal/lsp"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

func (svc *service) TextDocumentComplete(ctx context.Context, params lsp.CompletionParams) (lsp.CompletionList, error) {
	var list lsp.CompletionList

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return list, err
	}

	_, err = ilsp.ClientCapabilities(ctx)
	if err != nil {
		return list, err
	}

	doc, err := fs.GetDocument(ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI))
	if err != nil {
		return list, err
	}

	fPos, err := ilsp.FilePositionFromDocumentPosition(params.TextDocumentPositionParams, doc)
	if err != nil {
		return list, err
	}

	svc.logger.Printf("Looking for candidates at %q -> %#v", doc.Filename(), fPos.Position())

	data, err := doc.Text()
	if err != nil {
		return list, err
	}

	candidates := lsp.CompletionList{
		IsIncomplete: false,
		Items:        complete.CandidatesAtPos(data, doc.Filename(), fPos.Position(), svc.logger),
	}

	if len(candidates.Items) > 0 {
		telemetrySender, err := lsctx.Telemetry(ctx)
		if err != nil {
			return candidates, err
		}
		telemetrySender.SendEvent(ctx, "textDocument/completion", nil)
	}
	return candidates, err
}
