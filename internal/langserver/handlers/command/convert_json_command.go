package command

import (
	"context"
	"encoding/json"
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
)

var pluralizeClient = pluralize.NewClient()

type ConvertJsonCommand struct {
}

var _ CommandHandler = &ConvertJsonCommand{}

type ConvertJsonResponse struct {
	HCLContent string `json:"hclcontent"`
}

func (c ConvertJsonCommand) Handle(ctx context.Context, params CommandArgs) (interface{}, error) {
	content, ok := params.GetString("jsoncontent")
	if !ok {
		return nil, nil
	}

	var model map[string]interface{}
	err := json.Unmarshal([]byte(content), &model)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON content: %w", err)
	}

	result := ""
	if model["$schema"] != nil {
		result, err = convertARMTemplate(content)
		if err != nil {
			return nil, err
		}
	} else {
		result, err = convertResourceJson(content)
		if err != nil {
			return nil, err
		}
	}

	response := ConvertJsonResponse{
		HCLContent: result,
	}

	return response, nil
}
