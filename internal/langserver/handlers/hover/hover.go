package hover

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
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
