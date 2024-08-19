// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azapi-lsp/internal/langserver"
)

func TestExecuteCommand_basic(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	input, err := os.ReadFile(fmt.Sprintf("./testdata/%s/input.json", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expect, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	expectMap := make(map[string]interface{})
	expectMap["hclcontent"] = string(expect)

	expectJson, _ := json.Marshal(expectMap)

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

	reqParams := make(map[string]interface{})
	reqParams["command"] = "azapi.convertJsonToAzapi"
	reqParams["arguments"] = []string{fmt.Sprintf("jsoncontent=%s", string(input))}

	reqParamsString, _ := json.Marshal(reqParams)

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "workspace/executeCommand",
		ReqParams: string(reqParamsString)},
		fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 2,
		"result": %s
	}`, string(expectJson)))
}

func TestExecuteCommand_resourceGroup(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	input, err := os.ReadFile(fmt.Sprintf("./testdata/%s/input.json", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expect, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	expectMap := make(map[string]interface{})
	expectMap["hclcontent"] = string(expect)

	expectJson, _ := json.Marshal(expectMap)

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

	reqParams := make(map[string]interface{})
	reqParams["command"] = "azapi.convertJsonToAzapi"
	reqParams["arguments"] = []string{fmt.Sprintf("jsoncontent=%s", string(input))}

	reqParamsString, _ := json.Marshal(reqParams)

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "workspace/executeCommand",
		ReqParams: string(reqParamsString)},
		fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 2,
		"result": %s
	}`, string(expectJson)))
}

func TestExecuteCommand_nestedResource(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	input, err := os.ReadFile(fmt.Sprintf("./testdata/%s/input.json", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expect, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	expectMap := make(map[string]interface{})
	expectMap["hclcontent"] = string(expect)

	expectJson, _ := json.Marshal(expectMap)

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

	reqParams := make(map[string]interface{})
	reqParams["command"] = "azapi.convertJsonToAzapi"
	reqParams["arguments"] = []string{fmt.Sprintf("jsoncontent=%s", string(input))}

	reqParamsString, _ := json.Marshal(reqParams)

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "workspace/executeCommand",
		ReqParams: string(reqParamsString)},
		fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 2,
		"result": %s
	}`, string(expectJson)))
}

func TestExecuteCommand_armTemplate(t *testing.T) {
	tmpDir := TempDir(t)
	InitPluginCache(t, tmpDir.Dir())

	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{}))
	stop := ls.Start(t)
	defer stop()

	input, err := os.ReadFile(fmt.Sprintf("./testdata/%s/input.json", t.Name()))
	if err != nil {
		t.Fatal(err)
	}

	expect, err := os.ReadFile(fmt.Sprintf("./testdata/%s/expect.tf", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	expectMap := make(map[string]interface{})
	expectMap["hclcontent"] = string(expect)

	expectJson, _ := json.Marshal(expectMap)

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

	reqParams := make(map[string]interface{})
	reqParams["command"] = "azapi.convertJsonToAzapi"
	reqParams["arguments"] = []string{fmt.Sprintf("jsoncontent=%s", string(input))}

	reqParamsString, _ := json.Marshal(reqParams)

	ls.CallAndExpectResponse(t, &langserver.CallRequest{
		Method:    "workspace/executeCommand",
		ReqParams: string(reqParamsString)},
		fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 2,
		"result": %s
	}`, string(expectJson)))
}
