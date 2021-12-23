package handlers

import (
	"fmt"
	"testing"

	"github.com/creachadair/jrpc2/code"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/terraform/exec"
	"github.com/stretchr/testify/mock"
)

func TestShutdown_twice(t *testing.T) {
	tmpDir := TempDir(t)
	ls := langserver.NewLangServerMock(t, NewMockSession(&MockSessionInput{
		TerraformCalls: &exec.TerraformMockCalls{
			PerWorkDir: map[string][]*mock.Call{
				tmpDir.Dir(): validTfMockCalls(),
			},
		},
	}))
	stop := ls.Start(t)
	defer stop()

	ls.Call(t, &langserver.CallRequest{
		Method: "initialize",
		ReqParams: fmt.Sprintf(`{
	    "capabilities": {},
	    "rootUri": %q,
	    "processId": 12345
	}`, TempDir(t).URI())})
	ls.Call(t, &langserver.CallRequest{
		Method: "shutdown", ReqParams: `{}`})

	ls.CallAndExpectError(t, &langserver.CallRequest{
		Method: "shutdown", ReqParams: `{}`},
		code.InvalidRequest.Err())
}
