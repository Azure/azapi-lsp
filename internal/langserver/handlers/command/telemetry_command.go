package command

import (
	"context"
	"encoding/json"
	"fmt"
	lsctx "github.com/Azure/azapi-lsp/internal/context"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

type TelemetryCommand struct {
}

func (t TelemetryCommand) Handle(ctx context.Context, arguments []json.RawMessage) (interface{}, error) {
	telemetrySender, err := lsctx.Telemetry(ctx)
	if err != nil {
		return nil, err
	}

	var params lsp.TelemetryEvent
	if len(arguments) != 0 {
		err := json.Unmarshal(arguments[0], &params)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
		}
	}

	telemetrySender.SendEvent(ctx, params.Name, params.Properties)
	return nil, nil
}
