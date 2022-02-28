package lsp

import (
	"testing"

	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

var validUnixPath = "file:///valid/path/to/file.tf" //nolint

func TestFileHandler_invalid(t *testing.T) {
	path := "invalidpath"
	fh := FileHandlerFromDocumentURI(lsp.DocumentURI(path))
	if fh.Valid() {
		t.Fatalf("Expected %q to be invalid", path)
	}
}
