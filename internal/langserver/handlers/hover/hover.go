package hover

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure/types"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/common"
)

func HoverAtPos(data []byte, filename string, pos hcl.Pos, logger *log.Logger) *lang.HoverData {
	file, _ := hclsyntax.ParseConfig(data, filename, hcl.InitialPos)
	block, err := common.BlockAtPos(file, pos)
	if err == nil && block != nil && len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
		if attribute := common.AttributeAtPos(block, pos); attribute != nil {
			switch attribute.Name {
			case "type":
				if common.ContainsPos(attribute.NameRange, pos) {
					return Hover("type", "required", "string <resource-type>@<api-version>",
						"Azure Resource Manager type.", attribute.NameRange)
				}
			case "parent_id":
				if common.ContainsPos(attribute.NameRange, pos) {
					typeAttribute := common.AttributeWithName(block, "type")
					if typeAttribute != nil {
						if typeValue := common.ToLiteral(typeAttribute.Expr); typeValue != nil && len(*typeValue) != 0 {
							parentType := GetParentType(*typeValue)
							return Hover("parent_id", "required", "string",
								fmt.Sprintf("The ID of `%s` which is the parent resource in which this resource is created.", parentType), attribute.NameRange)
						}
					}
				}
			case "body":
				if common.ContainsPos(attribute.NameRange, pos) {
					return Hover("body", "optional", "string",
						"A JSON object that contains the request body used to create and update azure resource.", attribute.NameRange)
				}
				typeValue := common.ExtractAzureResourceType(block)
				if typeValue == nil {
					break
				}
				def, _ := azure.GetResourceDefinitionByResourceType(*typeValue)
				if def == nil {
					break
				}

				rangeMap := common.JsonEncodeExpressionToRangeMap(data, attribute.Expr)
				if rangeMap == nil {
					break
				}
				targetRangeMap := common.RangeMapOfPos(rangeMap, pos)
				if len(targetRangeMap) == 0 {
					break
				}
				lastRangeMap := targetRangeMap[len(targetRangeMap)-1]

				key := common.BuildKeyFromRangeMaps(targetRangeMap)
				logger.Printf("key: %v\n", key)

				if common.ContainsPos(lastRangeMap.KeyRange, pos) {
					start := len("..dummy.")
					end := strings.LastIndex(key, ".")
					propName := ""
					if start < end && start < len(key) && end > 0 {
						propName = key[end+1:]
						key = key[start:end]
					} else {
						propName = key[end+1:]
						key = ""
					}

					// input a property
					props := azure.ListProperties(def, key)
					logger.Printf("state: input key, key: %s", props)
					logger.Printf("received allowed keys: %#v", props)
					for _, prop := range props {
						if prop.Name == propName {
							return Hover(prop.Name, string(prop.Modifier), prop.Type, prop.Description, lastRangeMap.KeyRange)
						}
					}
				}
			}
		}
	}
	return nil
}

func Hover(name string, modifier string, propType string, description string, r hcl.Range) *lang.HoverData {
	return &lang.HoverData{
		Range:   r,
		Content: lang.Markdown(fmt.Sprintf("```\n%s: %s(%s)\n```\n%s", name, modifier, propType, description)),
	}
}

func GetParentType(resourceType string) string {
	parts := strings.Split(resourceType, "/")
	if len(parts) <= 2 {
		def, err := azure.GetResourceDefinitionByResourceType(resourceType)
		if err != nil || def == nil || len(def.ScopeTypes) == 0 {
			return "Azure Resource"
		}
		res := make([]string, 0)
		for _, scope := range def.ScopeTypes {
			switch scope {
			case types.Unknown:
				return "Azure Resource(Unknown scope)"
			case types.Tenant:
				res = append(res, "Tenant")
			case types.Subscription:
				res = append(res, "Subscription")
			case types.ManagementGroup:
				res = append(res, "Microsoft.Management/managementGroups")
			case types.ResourceGroup:
				res = append(res, "Microsoft.Resources/resourceGroups")
			case types.Extension:
				res = append(res, "Azure Resource(Extension scope)")
			}
		}
		return strings.Join(res, ", ")
	}
	return strings.Join(parts[0:len(parts)-1], "/")
}
