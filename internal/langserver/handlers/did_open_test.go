package handlers

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/filesystem"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/session"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/terraform/exec"
	"github.com/stretchr/testify/mock"
)

func TestLangServer_didOpenWithoutInitialization(t *testing.T) {
	ls := langserver.NewLangServerMock(t, NewMockSession(nil))
	stop := ls.Start(t)
	defer stop()

	ls.CallAndExpectError(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "terraform",
			"text": "provider \"github\" {}",
			"uri": "%s/main.tf"
		}
	}`, TempDir(t).URI())}, session.SessionNotInitialized.Err())
}

func TestLangServer_didOpenLanguageIdStored(t *testing.T) {
	tmpDir := TempDir(t)
	fs := filesystem.NewFilesystem()

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		TerraformCalls: &exec.TerraformMockCalls{
			PerWorkDir: map[string][]*mock.Call{
				tmpDir.Dir(): validTfMockCalls(),
			},
		},
		Filesystem: fs,
	}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {},
	    "rootUri": %q,
	    "processId": 12345
	}`, tmpDir.URI())})
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})

	originalText := `variable "service_host" {
  default = "blah"
}
`
	ls.Call(t, &langserver.CallRequest{
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
    "textDocument": {
        "languageId": "terraform",
        "version": 0,
        "uri": "%s/main.tf",
        "text": %q
    }
}`, TempDir(t).URI(), originalText)})
	path := filepath.Join(TempDir(t).Dir(), "main.tf")
	doc, err := fs.GetDocument(lsp.FileHandlerFromPath(path))
	if err != nil {
		t.Fatal(err)
	}
	languageID := doc.LanguageID()
	if diff := cmp.Diff(languageID, string("terraform")); diff != "" {
		t.Fatalf("unexpected languageID: %s", diff)
	}
	fullPath := doc.FullPath()
	if diff := cmp.Diff(fullPath, string(path)); diff != "" {
		t.Fatalf("unexpected fullPath: %s", diff)
	}
	version := doc.Version()
	if diff := cmp.Diff(version, int(0)); diff != "" {
		t.Fatalf("unexpected version: %s", diff)
	}
}
