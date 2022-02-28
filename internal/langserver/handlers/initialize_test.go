package handlers

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	ilsp "github.com/Azure/azapi-lsp/internal/lsp"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"

	"github.com/Azure/azapi-lsp/internal/langserver"
	"github.com/creachadair/jrpc2/code"
)

func TestInitialize_twice(t *testing.T) {
	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {},
	    "rootUri": %q,
	    "processId": 12345
	}`, TempDir(t).URI()),
	})
	ls.CallAndExpectError(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {},
	    "rootUri": %q,
	    "processId": 12345
	}`, TempDir(t).URI()),
	}, code.SystemError.Err())
}

func TestInitialize_withIncompatibleTerraformVersion(t *testing.T) {
	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {},
	    "processId": 12345,
	    "rootUri": %q
	}`, TempDir(t).URI()),
	})
}

func TestInitialize_multipleFolders(t *testing.T) {
	rootDir := TempDir(t)
	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {},
	    "rootUri": %q,
	    "processId": 12345,
	    "workspaceFolders": [
	    	{
	    		"uri": %q,
	    		"name": "root"
	    	}
	    ]
	}`, rootDir.URI(), rootDir.URI()),
	})
}

func TestInitialize_ignoreDirectoryNames(t *testing.T) {
	tmpDir := TempDir(t, "plugin", "ignore")
	pluginDir := filepath.Join(tmpDir.Dir(), "plugin")
	emptyDir := filepath.Join(tmpDir.Dir(), "ignore")

	InitPluginCache(t, pluginDir)
	InitPluginCache(t, emptyDir)

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
			"capabilities": {},
			"rootUri": %q,
			"processId": 12345,
			"initializationOptions": {
				"ignoreDirectoryNames": [%q]
			}
	}`, tmpDir.URI(), "ignore"),
	})
}

func TestName(t *testing.T) {
	serverCaps := lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: lsp.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    lsp.Incremental,
			},
			CompletionProvider: lsp.CompletionOptions{
				ResolveProvider:   false,
				TriggerCharacters: []string{".", "["},
			},
			CodeActionProvider: lsp.CodeActionOptions{
				CodeActionKinds: ilsp.SupportedCodeActions.AsSlice(),
				ResolveProvider: false,
			},
			HoverProvider: true,

			DeclarationProvider:        false,
			DefinitionProvider:         false,
			CodeLensProvider:           nil,
			ReferencesProvider:         false,
			DocumentFormattingProvider: false,
			DocumentSymbolProvider:     false,
			WorkspaceSymbolProvider:    false,
			Workspace:                  nil,
		},
	}
	data, _ := json.Marshal(serverCaps)
	fmt.Println(string(data))
}
