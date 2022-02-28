package handlers

import (
	"fmt"
	"testing"

	"github.com/Azure/azapi-lsp/internal/langserver"
	"github.com/creachadair/jrpc2/code"
)

func TestShutdown_twice(t *testing.T) {
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
	ls.Call(t, &langserver.CallRequest{
		Method: "shutdown", ReqParams: `{}`,
	})

	ls.CallAndExpectError(t, &langserver.CallRequest{
		Method: "shutdown", ReqParams: `{}`,
	},
		code.InvalidRequest.Err())
}
