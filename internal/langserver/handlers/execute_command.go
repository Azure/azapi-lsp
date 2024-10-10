package handlers

import (
	"context"
	"fmt"

	"github.com/Azure/azapi-lsp/internal/langserver/handlers/command"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

var handlerMap = map[string]command.CommandHandler{}

const CommandConvertJsonToAzapi = "azapi.convertJsonToAzapi"
const CommandAztfMigrate = "azapi.aztfmigrate"

func availableCommands() []string {
	return []string{CommandConvertJsonToAzapi, CommandAztfMigrate}
}

func init() {
	handlerMap = make(map[string]command.CommandHandler)
	handlerMap[CommandConvertJsonToAzapi] = command.ConvertJsonCommand{}
	handlerMap[CommandAztfMigrate] = command.AztfMigrateCommand{}
}

func (svc *service) WorkspaceExecuteCommand(ctx context.Context, params lsp.ExecuteCommandParams) (interface{}, error) {
	if params.Command == "editor.action.triggerSuggest" {
		// If this was actually received by the server, it means the client
		// does not support explicit suggest triggering, so we fail silently
		// TODO: Revisit once https://github.com/microsoft/language-server-protocol/issues/1117 is addressed
		return nil, nil
	}

	handler, ok := handlerMap[params.Command]
	if !ok {
		return nil, fmt.Errorf("command %q not found", params.Command)
	}
	return handler.Handle(ctx, params.Arguments)
}
