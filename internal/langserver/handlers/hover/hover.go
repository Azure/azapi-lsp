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
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/schema"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/parser"
)

func HoverAtPos(data []byte, filename string, pos hcl.Pos, logger *log.Logger) *lang.HoverData {
	file, _ := hclsyntax.ParseConfig(data, filename, hcl.InitialPos)
	block, err := parser.BlockAtPos(file, pos)
	if err == nil && block != nil && len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
		if attribute := parser.AttributeAtPos(block, pos); attribute != nil {
			switch attribute.Name {
			case "type":
				if parser.ContainsPos(attribute.NameRange, pos) {
					return Hover("type", "required", "string <resource-type>@<api-version>",
						"Azure Resource Manager type.", attribute.NameRange)
				}
			case "parent_id":
				if parser.ContainsPos(attribute.NameRange, pos) {
					typeAttribute := parser.AttributeWithName(block, "type")
					if typeAttribute != nil {
						if typeValue := parser.ToLiteral(typeAttribute.Expr); typeValue != nil && len(*typeValue) != 0 {
							parentType := GetParentType(*typeValue)
							return Hover("parent_id", "required", "string",
								fmt.Sprintf("The ID of `%s` which is the parent resource in which this resource is created.", parentType), attribute.NameRange)
						}
					}
				}
			case "body":
				if parser.ContainsPos(attribute.NameRange, pos) {
					return Hover("body", "optional", "string",
						"A JSON object that contains the request body used to create and update azure resource.", attribute.NameRange)
				}
				typeValue := parser.ExtractAzureResourceType(block)
				if typeValue == nil {
					break
				}
				def, _ := azure.GetResourceDefinitionByResourceType(*typeValue)
				if def == nil {
					break
				}

				hclNode := parser.JsonEncodeExpressionToHclNode(data, attribute.Expr)
				if hclNode == nil {
					break
				}
				hclNodes := parser.HclNodeArraysOfPos(hclNode, pos)
				if len(hclNodes) == 0 {
					break
				}
				lastHclNode := hclNodes[len(hclNodes)-1]

				if parser.ContainsPos(lastHclNode.KeyRange, pos) {
					scopes := schema.GetDef(def.AsTypeBase(), hclNodes, 0)
					defs := schema.GetDef(def.AsTypeBase(), hclNodes[0:len(hclNodes)-1], 0)
					props := make([]schema.Property, 0)
					for _, def := range defs {
						props = append(props, schema.GetAllowedProperties(def, scopes)...)
					}
					logger.Printf("received allowed keys: %#v", props)
					if len(props) != 0 {
						return Hover(props[0].Name, string(props[0].Modifier), props[0].Type, props[0].Description, lastHclNode.KeyRange)
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
