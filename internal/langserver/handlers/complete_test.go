package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azapi-lsp/internal/langserver"
	"github.com/Azure/azapi-lsp/internal/langserver/session"
)

func TestCompletion_withoutInitialization(t *testing.T) {
	ls := langserver.NewLangServerMock(t, NewMockSession(nil))
	stop := ls.Start(t)
	defer stop()

	ls.CallAndExpectError(t, &langserver.CallRequest{
		Method: "textDocument/completion",
		ReqParams: fmt.Sprintf(`{
			"textDocument": {
				"uri": "%s/main.tf"
			},
			"position": {
				"character": 0,
				"line": 1
			}
		}`, TempDir(t).URI()),
	}, session.SessionNotInitialized.Err())
}

func TestCompletion_apiVersion(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(4, 43, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_jsonencode(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(5, 13, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_value(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(7, 15, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_payload_value(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(7, 15, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_propWrappedByQuote(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(9, 11, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_payload_propWrappedByQuote(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(9, 11, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_action(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(3, 12, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_actionBody(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(5, 5, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_payload_actionBody(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(5, 5, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_discriminated(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(11, 7, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_payload_discriminated(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(11, 7, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_prop(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(12, 11, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_payload_prop(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(12, 11, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_codeSample(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expectRaw, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.json", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	expectRaw = []byte(strings.ReplaceAll(string(expectRaw), "<", "\\u003c"))
	expectRaw = []byte(strings.ReplaceAll(string(expectRaw), ">", "\\u003e"))

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

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(3, 3, tmpDir.URI()),
	}, string(expectRaw))
}

func TestCompletion_template(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
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
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(1, 14, tmpDir.URI()),
	})

	var result map[string]interface{}
	err = json.Unmarshal(rawResponse.Result, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	items := result["items"].([]interface{})
	if len(items) < 100 {
		t.Fatalf("expected at least 100 items, got %d", len(items))
	}
}

func TestCompletion_templateInsideOtherResource(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	config, err := os.ReadFile(fmt.Sprintf("./testdata/%s/main.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
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
		Method:    "textDocument/completion",
		ReqParams: buildReqParamsCompletion(2, 3, tmpDir.URI()),
	})

	var result map[string]interface{}
	err = json.Unmarshal(rawResponse.Result, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	items := result["items"].([]interface{})
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func buildReqParamsCompletion(line int, character int, uri string) string {
	param := make(map[string]interface{})
	textDocument := make(map[string]interface{})
	textDocument["uri"] = fmt.Sprintf("%s/main.tf", uri)
	param["textDocument"] = textDocument
	position := make(map[string]interface{})
	position["character"] = character - 1
	position["line"] = line - 1
	param["position"] = position
	paramJson, _ := json.Marshal(param)
	return string(paramJson)
}

func buildReqParamsTextDocument(text string, uri string) string {
	param := make(map[string]interface{})
	textDocument := make(map[string]interface{})
	textDocument["version"] = 0
	textDocument["languageId"] = "terraform"
	textDocument["text"] = text
	textDocument["uri"] = fmt.Sprintf("%s/main.tf", uri)
	param["textDocument"] = textDocument
	paramJson, _ := json.Marshal(param)
	return string(paramJson)
}
