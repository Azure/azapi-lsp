package command

import (
	"context"
	"encoding/json"
	"fmt"

	lsctx "github.com/Azure/azapi-lsp/internal/context"
	pluralize "github.com/gertd/go-pluralize"
)

var pluralizeClient = pluralize.NewClient()

type ConvertJsonCommand struct {
}

var _ CommandHandler = &ConvertJsonCommand{}

type ConvertJsonResponse struct {
	HCLContent string `json:"hclcontent"`
}

func (c ConvertJsonCommand) Handle(ctx context.Context, arguments []json.RawMessage) (interface{}, error) {
	params := ParseCommandArgs(arguments)
	content, ok := params.GetString("jsoncontent")
	if !ok {
		return nil, nil
	}

	telemetrySender, err := lsctx.Telemetry(ctx)
	if err != nil {
		return nil, err
	}

	var model map[string]interface{}
	err = json.Unmarshal([]byte(content), &model)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON content: %w", err)
	}

	result := ""
	if model["$schema"] != nil {
		telemetrySender.SendEvent(ctx, "ConvertJsonToAzapi", map[string]interface{}{
			"status": "started",
			"kind":   "arm-template",
		})
		result, err = convertARMTemplate(ctx, content, telemetrySender)
		if err != nil {
			telemetrySender.SendEvent(ctx, "ConvertJsonToAzapi", map[string]interface{}{
				"status": "failed",
				"kind":   "arm-template",
				"error":  err.Error(),
			})
			return nil, err
		}
	} else {
		telemetrySender.SendEvent(ctx, "ConvertJsonToAzapi", map[string]interface{}{
			"status": "started",
			"kind":   "resource-json",
		})
		result, err = convertResourceJson(ctx, content, telemetrySender)
		if err != nil {
			telemetrySender.SendEvent(ctx, "ConvertJsonToAzapi", map[string]interface{}{
				"status": "failed",
				"kind":   "resource-json",
				"error":  err.Error(),
			})
			return nil, err
		}
	}

	response := ConvertJsonResponse{
		HCLContent: result,
	}

	return response, nil
}
