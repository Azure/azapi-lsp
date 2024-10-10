package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/Azure/azapi-lsp/internal/langserver"
	"github.com/Azure/azapi-lsp/internal/langserver/session"
	"github.com/Azure/azapi-lsp/internal/protocol"
)

func TestCodeAction_withoutInitialization(t *testing.T) {
	ls := langserver.NewLangServerMock(t, NewMockSession(nil))
	stop := ls.Start(t)
	defer stop()

	ls.CallAndExpectError(t, &langserver.CallRequest{
		Method: "textDocument/codeAction",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "terraform",
			"text": "provider \"github\" {}",
			"uri": "%s/main.tf"
		}
	}`, TempDir(t).URI()),
	}, session.SessionNotInitialized.Err())
}

func TestCodeAction_migrateToAzureRM(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	reqParams := buildReqParamsCodeAction(6, 1, 19, 2, tmpDir.URI())

	expected := []protocol.CodeAction{
		{
			Title: "Migrate to AzureRM Provider",
			Kind:  "refactor.rewrite",
			Edit:  protocol.WorkspaceEdit{},
			Command: &protocol.Command{
				Title:     "Migrate to AzureRM Provider",
				Command:   "azapi.aztfmigrate",
				Arguments: []json.RawMessage{[]byte(reqParams)},
			},
		},
	}

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
		"capabilities": {},
		"rootUri": %q,
		"processId": 12345
	}`, tmpDir.URI()),
	})
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Call(t, &langserver.CallRequest{
		Method:    "textDocument/didOpen",
		ReqParams: buildReqParamsTextDocument(string(config), tmpDir.URI()),
	})

	rawResponse := ls.Call(t, &langserver.CallRequest{
		Method:    "textDocument/codeAction",
		ReqParams: reqParams,
	})

	var result []protocol.CodeAction
	err = json.Unmarshal(rawResponse.Result, &result)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestCodeAction_migrateToAzapi(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	reqParams := buildReqParamsCodeAction(1, 1, 4, 2, tmpDir.URI())

	expected := []protocol.CodeAction{
		{
			Title: "Migrate to Azapi Provider",
			Kind:  "refactor.rewrite",
			Edit:  protocol.WorkspaceEdit{},
			Command: &protocol.Command{
				Title:     "Migrate to Azapi Provider",
				Command:   "azapi.aztfmigrate",
				Arguments: []json.RawMessage{[]byte(reqParams)},
			},
		},
	}

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
		"capabilities": {},
		"rootUri": %q,
		"processId": 12345
	}`, tmpDir.URI()),
	})
	ls.Notify(t, &langserver.CallRequest{
		Method:    "initialized",
		ReqParams: "{}",
	})
	ls.Call(t, &langserver.CallRequest{
		Method:    "textDocument/didOpen",
		ReqParams: buildReqParamsTextDocument(string(config), tmpDir.URI()),
	})

	rawResponse := ls.Call(t, &langserver.CallRequest{
		Method:    "textDocument/codeAction",
		ReqParams: reqParams,
	})

	var result []protocol.CodeAction
	err = json.Unmarshal(rawResponse.Result, &result)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func buildReqParamsCodeAction(startLine, startCharacter, endLine, endCharacter int, uri string) string {
	param := protocol.CodeActionParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.DocumentURI(fmt.Sprintf("%s/main.tf", uri)),
		},
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      uint32(startLine - 1),
				Character: uint32(startCharacter - 1),
			},
			End: protocol.Position{
				Line:      uint32(endLine - 1),
				Character: uint32(endCharacter - 1),
			},
		},
		Context: protocol.CodeActionContext{},
	}
	paramJson, _ := json.Marshal(param)
	return string(paramJson)
}
