package lsp

import (
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/hashicorp/hcl-lang/lang"
)

func Command(cmd lang.Command) (lsp.Command, error) {
	lspCmd := lsp.Command{
		Title:   cmd.Title,
		Command: cmd.ID,
	}

	for _, arg := range cmd.Arguments {
		jsonArg, err := arg.MarshalJSON()
		if err != nil {
			return lspCmd, err
		}
		lspCmd.Arguments = append(lspCmd.Arguments, jsonArg)
	}

	return lspCmd, nil
}
