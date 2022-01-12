package handlers

import (
	"fmt"
	"testing"

	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/session"
)

func TestLangServer_codeActionWithoutInitialization(t *testing.T) {
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

func TestLangServer_codeAction_basic(t *testing.T) {
	// code action is not implemented yet
	t.SkipNow()
	tmpDir := TempDir(t)

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

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
		Method: "textDocument/didOpen",
		ReqParams: fmt.Sprintf(`{
		"textDocument": {
			"version": 0,
			"languageId": "terraform",
			"text": "provider  \"test\"   {\n\n      }\n",
			"uri": "%s/main.tf"
		}
	}`, tmpDir.URI()),
	})
	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method: "textDocument/codeAction",
		ReqParams: fmt.Sprintf(`{
			"textDocument": { "uri": "%s/main.tf" },
			"range": {
				"start": { "line": 0, "character": 0 },
				"end": { "line": 1, "character": 0 }
			},
			"context": { "diagnostics": [], "only": ["source.formatAll.terraform"] }
		}`, tmpDir.URI()),
	}, fmt.Sprintf(`{
			"jsonrpc": "2.0",
			"id": 3,
			"result": [
				{
					"title": "Format Document",
					"kind": "source.formatAll.terraform",
					"edit":{
						"changes":{
							"%s/main.tf": [
								{
									"range": {
										"start": {
											"line": 0,
											"character": 0
										},
										"end": {
											"line": 1,
											"character": 0
										}
									},
									"newText": "provider \"test\" {\n"
								},
								{
									"range": {
										"start": {
											"line": 2,
											"character": 0
										},
										"end": {
											"line": 3,
											"character": 0
										}
									},
									"newText": "}\n"
								}
							]
						}
					}
				}
			]
		}`, tmpDir.URI()))
}

func TestLangServer_codeAction_no_code_action_requested(t *testing.T) {
	tmpDir := TempDir(t)

	tests := []struct {
		name    string
		request *langserver.CallRequest
		want    string
	}{
		{
			name: "no code action requested",
			request: &langserver.CallRequest{
				Method: "textDocument/codeAction",
				ReqParams: fmt.Sprintf(`{
						"textDocument": { "uri": "%s/main.tf" },
						"range": {
							"start": { "line": 0, "character": 0 },
							"end": { "line": 1, "character": 0 }
						},
						"context": { "diagnostics": [], "only": [""] }
					}`, tmpDir.URI()),
			},
			want: `{
				"jsonrpc": "2.0",
				"id": 3,
				"result": null
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
			stop := ls.Start(t)
			defer stop()

			ls.Call(t, &langserver.CallRequest{
				Method: "initialize",
				ReqParams: fmt.Sprintf(`{
					"capabilities": {},
					"rootUri": %q,
					"processId": 123456
			}`, tmpDir.URI()),
			})
			ls.Notify(t, &langserver.CallRequest{
				Method:    "initialized",
				ReqParams: "{}",
			})
			ls.Call(t, &langserver.CallRequest{
				Method: "textDocument/didOpen",
				ReqParams: fmt.Sprintf(`{
				"textDocument": {
					"version": 0,
					"languageId": "terraform",
					"text": "provider  \"test\"   {\n\n      }\n",
					"uri": "%s/main.tf"
				}
			}`, tmpDir.URI()),
			})

			ls.CallAndExpectResponse(t, tt.request, tt.want)
		})
	}
}
