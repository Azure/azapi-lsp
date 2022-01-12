package lsp

import (
	"testing"

	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func TestFileHandler_invalid(t *testing.T) {
	path := "invalidpath"
	fh := FileHandlerFromDocumentURI(lsp.DocumentURI(path))
	if fh.Valid() {
		t.Fatalf("Expected %q to be invalid", path)
	}
}
