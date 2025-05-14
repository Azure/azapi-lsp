package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azapi-lsp/internal/azure"
	"github.com/Azure/azapi-lsp/internal/azure/types"
	"github.com/Azure/azapi-lsp/internal/telemetry"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type resourceModel struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	APIVersion string            `json:"apiVersion"`
	Location   string            `json:"location"`
	ID         string            `json:"id"`
	Tags       map[string]string `json:"tags"`
	Identity   *identityModel    `json:"identity"`
	DependsOn  []string          `json:"dependsOn"`
}

type identityModel struct {
	Type                   string                 `json:"type"`
	UserAssignedIdentities map[string]interface{} `json:"userAssignedIdentities"`
}

func convertResourceJson(ctx context.Context, input string, telemetrySender telemetry.Sender) (string, error) {
	var model resourceModel
	err := json.Unmarshal([]byte(input), &model)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal JSON content: %w", err)
	}

	block, err := ParseResourceJson(input)
	if err != nil {
		return "", err
	}

	label := ""
	if len(block.Labels()) == 2 {
		label = block.Labels()[1]
	}
	apiVersion := ""
	typeValue := string(block.Body().GetAttribute("type").Expr().BuildTokens(nil).Bytes())
	typeValue = strings.Trim(typeValue, " \"")
	if parts := strings.Split(typeValue, "@"); len(parts) == 2 {
		apiVersion = parts[1]
	}

	telemetrySender.SendEvent(ctx, "ConvertJsonToAzapi", map[string]interface{}{
		"status": "completed",
		"kind":   "resource-json",
		"type":   typeValue,
	})

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

	return result, nil
}

func ParseResourceJson(content string) (*hclwrite.Block, error) {
	var model resourceModel
	err := json.Unmarshal([]byte(content), &model)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal JSON content: %w", err)
	}

	block := hclwrite.NewBlock("resource", []string{})
	parentId := ""
	if armId, err := arm.ParseResourceID(model.ID); err == nil {
		// use the resource type from the ID
		model.Type = armId.ResourceType.String()
		parentId = armId.Parent.String()
	}
	if parentId == "" {
		for _, v := range model.DependsOn {
			armId, err := arm.ParseResourceID(v)
			if err != nil {
				continue
			}
			if armId.ResourceType.String() == GetParentType(model.Type) {
				parentId = armId.String()
				break
			}
		}
	}
	if parentId == "" {
		parentId = "/subscriptions/${var.subscriptionId}/resourceGroups/${var.resourceGroupName}"
	}

	if model.APIVersion == "" {
		apiVersions := azure.GetApiVersions(model.Type)
		apiVersion := "TODO"
		if len(apiVersions) > 0 {
			apiVersion = apiVersions[len(apiVersions)-1]
		}
		model.APIVersion = apiVersion
	}

	label := pluralizeClient.Singular(LastSegment(model.Type))
	block.SetLabels([]string{"azapi_resource", label})
	block.Body().SetAttributeValue("type", cty.StringVal(fmt.Sprintf("%s@%s", model.Type, model.APIVersion)))
	block.Body().SetAttributeValue("parent_id", cty.StringVal(parentId))

	nameValue := model.Name[strings.LastIndex(model.Name, "/")+1:]
	block.Body().SetAttributeValue("name", cty.StringVal(nameValue))

	def, _ := azure.GetResourceDefinition(model.Type, model.APIVersion)
	if model.Location != "" && (def == nil || canResourceHaveProperty(def, "location")) {
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
	delete(bodyMap, "location")
	delete(bodyMap, "id")
	delete(bodyMap, "tags")
	delete(bodyMap, "identity")
	bodyMap["name"] = nameValue
	if label == "basicPublishingCredentialsPolicy" {
		fmt.Println(bodyMap)
	}
	if def != nil {
		writeOnlyBody := def.GetWriteOnly(NormalizeObject(bodyMap))
		ok := false
		bodyMap, ok = writeOnlyBody.(map[string]interface{})
		if !ok {
			log.Printf("[ERROR] unable to get write only body, result: %v", writeOnlyBody)
		}
	}
	delete(bodyMap, "name")
	block.Body().SetAttributeValue("body", toCtyValue(bodyMap))

	if len(model.Tags) > 0 {
		tags := map[string]cty.Value{}
		for k, v := range model.Tags {
			tags[k] = cty.StringVal(v)
		}
		block.Body().SetAttributeValue("tags", cty.MapVal(tags))
	}

	return block, nil
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

func GetParentType(resourceType string) string {
	parts := strings.Split(resourceType, "/")
	if len(parts) <= 2 {
		return ""
	}
	return strings.Join(parts[0:len(parts)-1], "/")
}

func NormalizeObject(input interface{}) interface{} {
	jsonString, _ := json.Marshal(input)
	var output interface{}
	_ = json.Unmarshal(jsonString, &output)
	return output
}

func canResourceHaveProperty(resourceDef *types.ResourceType, property string) bool {
	if resourceDef == nil || resourceDef.Body == nil || resourceDef.Body.Type == nil {
		return false
	}
	objectType, ok := (*resourceDef.Body.Type).(*types.ObjectType)
	if !ok {
		return false
	}
	if prop, ok := objectType.Properties[property]; ok {
		if !prop.IsReadOnly() {
			return true
		}
	}
	return false
}
