package complete

import (
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure/types"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/schema"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/parser"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func CandidatesAtPos(data []byte, filename string, pos hcl.Pos, logger *log.Logger) []lsp.CompletionItem {
	file, _ := hclsyntax.ParseConfig(data, filename, hcl.InitialPos)
	block, err := parser.BlockAtPos(file, pos)
	candidateList := make([]lsp.CompletionItem, 0)
	if err == nil && block != nil && len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
		if attribute := parser.AttributeAtPos(block, pos); attribute != nil {
			switch attribute.Name {
			case "type":
				prefix := parser.ToLiteral(attribute.Expr)
				candidateList = typeCandidates(prefix, editRangeFromExprRange(attribute.Expr, pos))
			case "body":
				if attribute.Expr != nil {
					if _, ok := attribute.Expr.(*hclsyntax.LiteralValueExpr); ok && parser.ToLiteral(attribute.Expr) == nil {
						candidateList = append(candidateList, bodyFuncCandidate(editRangeFromExprRange(attribute.Expr, pos)))
					}
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

				switch {
				case parser.ContainsPos(lastHclNode.KeyRange, pos):
					// input a property with a prefix
					hclNodes := hclNodes[0 : len(hclNodes)-1]
					defs := schema.GetDef(def.AsTypeBase(), hclNodes, 0)
					keys := make([]schema.Property, 0)
					for _, def := range defs {
						keys = append(keys, schema.GetAllowedProperties(def, []*types.TypeBase{})...)
					}
					if len(hclNodes) == 1 {
						keys = ignorePulledOutProperties(keys)
					}
					logger.Printf("received allowed keys: %#v", keys)
					editRange := ilsp.HCLRangeToLSP(lastHclNode.KeyRange)
					candidateList = keyCandidates(keys, editRange)
				case !lastHclNode.KeyRange.Empty() && !lastHclNode.EqualRange.Empty() && lastHclNode.Children == nil:
					// input property =
					defs := schema.GetDef(def.AsTypeBase(), hclNodes, 0)
					values := make([]string, 0)
					for _, def := range defs {
						values = append(values, schema.GetAllowedValues(def)...)
					}
					editRange := lastHclNode.ValueRange
					if lastHclNode.Value == nil {
						editRange.End = pos
					}
					candidateList = valueCandidates(values, ilsp.HCLRangeToLSP(editRange))
					logger.Printf("received allowed keys: %#v", values)
				case parser.ContainsPos(lastHclNode.ValueRange, pos):
					// input a property
					defs := schema.GetDef(def.AsTypeBase(), hclNodes, 0)
					keys := make([]schema.Property, 0)
					for _, def := range defs {
						keys = append(keys, schema.GetAllowedProperties(def, []*types.TypeBase{})...)
					}
					if len(hclNodes) == 1 {
						keys = ignorePulledOutProperties(keys)
					}
					logger.Printf("received allowed keys: %#v", keys)
					editRange := ilsp.HCLRangeToLSP(hcl.Range{Start: pos, End: pos, Filename: filename})
					candidateList = keyCandidates(keys, editRange)

					if len(lastHclNode.Children) == 0 {
						propertySets := make([]schema.PropertySet, 0)
						for _, def := range defs {
							propertySets = append(propertySets, schema.GetRequiredPropertySet(def)...)
						}
						if len(hclNodes) == 1 {
							for i, ps := range propertySets {
								propertySets[i].Name = ""
								propertySets[i].Properties = ignorePulledOutPropertiesFromPropertySet(ps.Properties)
							}
						}
						candidateList = append(candidateList, requiredPropertiesCandidates(propertySets, editRange)...)
					}
				}
			}
		}
	}
	return candidateList
}

func editRangeFromExprRange(expression hclsyntax.Expression, pos hcl.Pos) lsp.Range {
	expRange := expression.Range()
	if expRange.Start.Line != expRange.End.Line && expRange.End.Column == 1 && expRange.End.Line-1 == pos.Line {
		expRange.End = pos
	}
	return ilsp.HCLRangeToLSP(expRange)
}
