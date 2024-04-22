package hover

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azapi-lsp/internal/azure"
	"github.com/Azure/azapi-lsp/internal/azure/types"
	"github.com/Azure/azapi-lsp/internal/langserver/handlers/tfschema"
	"github.com/Azure/azapi-lsp/internal/langserver/schema"
	ilsp "github.com/Azure/azapi-lsp/internal/lsp"
	"github.com/Azure/azapi-lsp/internal/parser"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func HoverAtPos(data []byte, filename string, pos hcl.Pos, logger *log.Logger) *lsp.Hover {
	file, _ := hclsyntax.ParseConfig(data, filename, hcl.InitialPos)
	body, isHcl := file.Body.(*hclsyntax.Body)
	if !isHcl {
		logger.Printf("file is not hcl")
		return nil
	}
	block := parser.BlockAtPos(body, pos)
	if block != nil && len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azapi") {
		resourceName := fmt.Sprintf("%s.%s", block.Type, block.Labels[0])
		resource := tfschema.GetResourceSchema(resourceName)
		if resource == nil {
			return nil
		}
		if attribute := parser.AttributeAtPos(block, pos); attribute != nil {
			property := resource.GetProperty(attribute.Name)
			if property == nil {
				return nil
			}
			switch attribute.Name {
			case "parent_id":
				if parser.ContainsPos(attribute.NameRange, pos) {
					typeAttribute := parser.AttributeWithName(block, "type")
					if typeAttribute != nil {
						if typeValue := parser.ToLiteral(typeAttribute.Expr); typeValue != nil && len(*typeValue) != 0 {
							parentType := GetParentType(*typeValue)
							return Hover(property.Name, property.Modifier, property.Type,
								fmt.Sprintf("The ID of `%s` which is the parent resource in which this resource is created.", parentType), attribute.NameRange)
						}
					}
				}
			case "body":
				if parser.ContainsPos(attribute.NameRange, pos) {
					return Hover(property.Name, property.Modifier, property.Type, property.Description, attribute.NameRange)
				}

				bodyDef := tfschema.BodyDefinitionFromBlock(block)
				if bodyDef == nil {
					return nil
				}

				hclNode := parser.JsonEncodeExpressionToHclNode(data, attribute.Expr)
				if hclNode == nil {
					tokens, _ := hclsyntax.LexExpression(data[attribute.Expr.Range().Start.Byte:attribute.Expr.Range().End.Byte], filename, attribute.Expr.Range().Start)
					hclNode = parser.BuildHclNode(tokens)
				}
				if hclNode == nil {
					break
				}

				return hoverOnBody(hclNode, pos, bodyDef)
			default:
				if !parser.ContainsPos(attribute.NameRange, pos) {
					return nil
				}
				return Hover(property.Name, property.Modifier, property.Type, property.Description, attribute.NameRange)
			}
		} else {
			if subBody := parser.BlockAtPos(block.Body, pos); subBody != nil {
				property := resource.GetProperty(subBody.Type)
				if property == nil {
					return nil
				}
				if attribute := parser.AttributeAtPos(subBody, pos); attribute != nil {
					for _, p := range property.NestedProperties {
						if p.Name == attribute.Name {
							return Hover(p.Name, p.Modifier, p.Type, p.Description, attribute.NameRange)
						}
					}
				} else {
					return Hover(property.Name, property.Modifier, property.Type, property.Description, subBody.TypeRange)
				}
			}
		}
	}
	return nil
}

func hoverOnBody(hclNode *parser.HclNode, pos hcl.Pos, bodyDef types.TypeBase) *lsp.Hover {
	hclNodes := parser.HclNodeArraysOfPos(hclNode, pos)
	if len(hclNodes) == 0 {
		return nil
	}
	lastHclNode := hclNodes[len(hclNodes)-1]

	if parser.ContainsPos(lastHclNode.KeyRange, pos) {
		defs := schema.GetDef(bodyDef.AsTypeBase(), hclNodes[0:len(hclNodes)-1], 0)
		props := make([]schema.Property, 0)
		for _, def := range defs {
			props = append(props, schema.GetAllowedProperties(def)...)
		}
		if len(props) != 0 {
			index := -1
			for i := range props {
				if props[i].Name == lastHclNode.Key {
					index = i
					break
				}
			}
			if index == -1 {
				return nil
			}
			return Hover(props[index].Name, string(props[index].Modifier), props[index].Type, props[index].Description, lastHclNode.KeyRange)
		}
	}
	return nil
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

func Hover(name string, modifier string, propType string, description string, r hcl.Range) *lsp.Hover {
	return &lsp.Hover{
		Range: ilsp.HCLRangeToLSP(r),
		Contents: lsp.MarkupContent{
			Kind:  lsp.Markdown,
			Value: fmt.Sprintf("```\n%s: %s(%s)\n```\n%s", name, modifier, propType, description),
		},
	}
}
