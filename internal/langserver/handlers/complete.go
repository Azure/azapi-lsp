package handlers

import (
	"context"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func (svc *service) TextDocumentComplete(ctx context.Context, params lsp.CompletionParams) (lsp.CompletionList, error) {
	var list lsp.CompletionList

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return list, err
	}

	cc, err := ilsp.ClientCapabilities(ctx)
	if err != nil {
		return list, err
	}

	doc, err := fs.GetDocument(ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI))
	if err != nil {
		return list, err
	}

	d, err := svc.decoderForDocument(ctx, doc)
	if err != nil {
		return list, err
	}

	expFeatures, err := lsctx.ExperimentalFeatures(ctx)
	if err != nil {
		return list, err
	}

	d.PrefillRequiredFields = expFeatures.PrefillRequiredFields

	fPos, err := ilsp.FilePositionFromDocumentPosition(params.TextDocumentPositionParams, doc)

	if err != nil {
		return list, err
	}

	svc.logger.Printf("Looking for candidates at %q -> %#v", doc.Filename(), fPos.Position())

	_, _ = d.CandidatesAtPos(doc.Filename(), fPos.Position())
	candidates := lang.Candidates{
		List:       []lang.Candidate{
			{
				Label: "henglu",
				Description: lang.Markdown("##title\nthis is description"),
				Detail: "block",
				IsDeprecated: false,
				TextEdit: lang.TextEdit{
					Range:   hcl.Range{
						Filename: doc.Filename(),
						End: fPos.Position(),
						Start: hcl.Pos{
							Line: fPos.Position().Line,
							Column: fPos.Position().Column - 1,
						},
					},
					NewText: "henglu = ",
					Snippet: "henglu = ",
				},
				Kind: lang.AttributeCandidateKind,
				TriggerSuggest: true,
			},
		},
		IsComplete: true,
	}
	svc.logger.Printf("received candidates: %#v", candidates)
	return ilsp.ToCompletionList(candidates, cc.TextDocument), err
}
