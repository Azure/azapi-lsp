package handlers

import (
	"context"

	"github.com/hashicorp/hcl-lang/lang"
	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func (svc *service) References(ctx context.Context, params lsp.ReferenceParams) ([]lsp.Location, error) {
	list := make([]lsp.Location, 0)

	fs, err := lsctx.DocumentStorage(ctx)
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

	path := lang.Path{
		Path:       doc.Dir(),
		LanguageID: doc.LanguageID(),
	}

	origins := svc.decoder.ReferenceOriginsTargetingPos(path, doc.Filename(), fPos.Position())

	return ilsp.RefOriginsToLocations(origins), nil
}
