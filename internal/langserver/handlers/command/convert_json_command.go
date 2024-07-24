package command

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azapi-lsp/internal/azure"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var pluralizeClient = pluralize.NewClient()

type ConvertJsonCommand struct {
}

var _ CommandHandler = &ConvertJsonCommand{}

type ConvertJsonResponse struct {
	HCLContent string `json:"hclcontent"`
}

type resourceModel struct {
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	Location string            `json:"location"`
	ID       string            `json:"id"`
	Tags     map[string]string `json:"tags"`
	Identity *identityModel    `json:"identity"`
}

type identityModel struct {
	Type                   string                 `json:"type"`
	UserAssignedIdentities map[string]interface{} `json:"userAssignedIdentities"`
}

func (c ConvertJsonCommand) Handle(ctx context.Context, params CommandArgs) (interface{}, error) {
	content, ok := params.GetString("jsoncontent")
	if !ok {
		return nil, nil
	}

	var model resourceModel
	err := json.Unmarshal([]byte(content), &model)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON content: %w", err)
	}

	armId, err := arm.ParseResourceID(model.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse resource ID: %w", err)
	}

	model.Type = armId.ResourceType.String()

	apiVersions := azure.GetApiVersions(model.Type)
	apiVersion := "TODO"
	if len(apiVersions) > 0 {
		apiVersion = apiVersions[len(apiVersions)-1]
	}

	label := pluralizeClient.Singular(LastSegment(model.Type))
	block := hclwrite.NewBlock("resource", []string{"azapi_resource", label})
	block.Body().SetAttributeValue("type", cty.StringVal(fmt.Sprintf("%s@%s", model.Type, apiVersion)))
	block.Body().SetAttributeValue("parent_id", cty.StringVal(armId.Parent.String()))
	block.Body().SetAttributeValue("name", cty.StringVal(model.Name))

	if model.Location != "" {
		block.Body().SetAttributeValue("location", cty.StringVal(model.Location))
	}

	if model.Identity != nil && !strings.EqualFold(model.Identity.Type, "None") {
		identityBlock := hclwrite.NewBlock("identity", nil)
		identityBlock.Body().SetAttributeValue("type", cty.StringVal(model.Identity.Type))
		if len(model.Identity.UserAssignedIdentities) > 0 {
			identityIds := make([]cty.Value, 0)
			for k := range model.Identity.UserAssignedIdentities {
				identityIds = append(identityIds, cty.StringVal(k))
			}
			identityBlock.Body().SetAttributeValue("identity_ids", cty.ListVal(identityIds))
		} else {
			identityBlock.Body().SetAttributeValue("identity_ids", cty.ListValEmpty(cty.String))
		}
		block.Body().AppendBlock(identityBlock)
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal([]byte(content), &bodyMap)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON content: %w", err)
	}
	delete(bodyMap, "type")
	delete(bodyMap, "name")
	delete(bodyMap, "location")
	delete(bodyMap, "id")
	delete(bodyMap, "tags")
	delete(bodyMap, "identity")
	def, err := azure.GetResourceDefinition(model.Type, apiVersion)
	if err == nil && def != nil {
		writeOnlyBody := def.GetWriteOnly(NormalizeObject(bodyMap))
		bodyMap, ok = writeOnlyBody.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unable to get write only body")
		}
	}
	block.Body().SetAttributeValue("body", toCtyValue(bodyMap))

	if len(model.Tags) > 0 {
		tags := map[string]cty.Value{}
		for k, v := range model.Tags {
			tags[k] = cty.StringVal(v)
		}
		block.Body().SetAttributeValue("tags", cty.MapVal(tags))
	}

	importBlock := hclwrite.NewBlock("import", nil)
	importBlock.Body().SetAttributeValue("id", cty.StringVal(fmt.Sprintf("%s?api-version=%s", model.ID, apiVersion)))
	importBlock.Body().SetAttributeTraversal("to", hcl.Traversal{hcl.TraverseRoot{Name: "azapi_resource"}, hcl.TraverseAttr{Name: label}})

	result := fmt.Sprintf(`/*
Note: This is a generated HCL content from the JSON input which is based on the latest API version available.
To import the resource, please run the following command:
terraform import azapi_resource.%s %s?api-version=%s

Or add the below config:
%s*/

%s`, label, model.ID, apiVersion, string(hclwrite.Format(importBlock.BuildTokens(nil).Bytes())), string(hclwrite.Format(block.BuildTokens(nil).Bytes())))

	response := ConvertJsonResponse{
		HCLContent: result,
	}

	return response, nil
}

func toCtyValue(input interface{}) cty.Value {
	if input == nil {
		return cty.NullVal(cty.DynamicPseudoType)
	}
	switch v := input.(type) {
	case map[string]interface{}:
		m := map[string]cty.Value{}
		for k, v := range v {
			m[k] = toCtyValue(v)
		}
		return cty.ObjectVal(m)
	case []interface{}:
		l := make([]cty.Value, len(v))
		for i, e := range v {
			l[i] = toCtyValue(e)
		}
		return cty.TupleVal(l)
	case string:
		return cty.StringVal(v)
	case bool:
		return cty.BoolVal(v)
	case float64:
		return cty.NumberFloatVal(v)
	case float32:
		return cty.NumberFloatVal(float64(v))
	case int:
		return cty.NumberIntVal(int64(v))
	case int64:
		return cty.NumberIntVal(v)
	case int32:
		return cty.NumberIntVal(int64(v))
	default:
		return cty.NullVal(cty.DynamicPseudoType)
	}
}

func LastSegment(input string) string {
	id := strings.Trim(input, "/")
	components := strings.Split(id, "/")
	if len(components) == 0 {
		return ""
	}
	return components[len(components)-1]
}

func NormalizeObject(input interface{}) interface{} {
	jsonString, _ := json.Marshal(input)
	var output interface{}
	_ = json.Unmarshal(jsonString, &output)
	return output
}
