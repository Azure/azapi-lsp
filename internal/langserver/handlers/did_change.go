package handlers

import (
	"context"
	"fmt"

	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/validate"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func TextDocumentDidChange(ctx context.Context, params lsp.DidChangeTextDocumentParams) error {
	p := lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{
				URI: params.TextDocument.URI,
			},
			Version: params.TextDocument.Version,
		},
		ContentChanges: params.ContentChanges,
	}

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return err
	}

	fh := ilsp.VersionedFileHandler(p.TextDocument)
	f, err := fs.GetDocument(fh)
	if err != nil {
		return err
	}

	// Versions don't have to be consecutive, but they must be increasing
	if int(p.TextDocument.Version) <= f.Version() {
		_ = fs.CloseAndRemoveDocument(fh)
		return fmt.Errorf("old version (%d) received, current version is %d. "+
			"Unable to update %s. This is likely a bug, please report it",
			int(p.TextDocument.Version), f.Version(), p.TextDocument.URI)
	}

	changes, err := ilsp.DocumentChanges(params.ContentChanges, f)
	if err != nil {
		return err
	}
	err = fs.ChangeDocument(fh, changes)
	if err != nil {
		return err
	}

	notifier, err := lsctx.DiagnosticsNotifier(ctx)
	if err != nil {
		return err
	}

	src, err := f.Text()
	if err != nil {
		return err
	}

	diags := validate.NewDiagnostics(src, f.Filename())
	notifier.PublishHCLDiags(ctx, f.Dir(), diags)

	return nil
}
