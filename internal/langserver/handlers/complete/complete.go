package complete

import (
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/common"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func CandidatesAtPos(data []byte, filename string, pos hcl.Pos, logger *log.Logger) []lsp.CompletionItem {
	file, _ := hclsyntax.ParseConfig(data, filename, hcl.InitialPos)
	block, err := common.BlockAtPos(file, pos)
	candidateList := make([]lsp.CompletionItem, 0)
	if err == nil && block != nil && len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
		if attribute := common.AttributeAtPos(block, pos); attribute != nil {
			switch attribute.Name {
			case "type":
				prefix := common.ToLiteral(attribute.Expr)
				candidateList = typeCandidates(prefix, editRangeFromExprRange(attribute.Expr, pos))
			case "body":
				if attribute.Expr != nil {
					if _, ok := attribute.Expr.(*hclsyntax.LiteralValueExpr); ok && common.ToLiteral(attribute.Expr) == nil {
						candidateList = append(candidateList, bodyFuncCandidate(editRangeFromExprRange(attribute.Expr, pos)))
					}
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
				rangeMaps := common.RangeMapArraysOfPos(rangeMap, pos)
				if len(rangeMaps) == 0 {
					break
				}
				lastRangeMap := rangeMaps[len(rangeMaps)-1]

				switch {
				case common.ContainsPos(lastRangeMap.KeyRange, pos):
					// input a property with a prefix
					rangeMaps := rangeMaps[0 : len(rangeMaps)-1]
					keys := azure.ListProperties(def, rangeMaps)
					if len(rangeMaps) == 1 {
						keys = ignorePulledOutProperties(keys)
					}
					logger.Printf("received allowed keys: %#v", keys)
					editRange := ilsp.HCLRangeToLSP(lastRangeMap.KeyRange)
					candidateList = keyCandidates(keys, editRange)
				case !lastRangeMap.KeyRange.Empty() && !lastRangeMap.EqualRange.Empty() && lastRangeMap.Children == nil:
					// input property =
					values := azure.ListValues(def, rangeMaps)
					editRange := lastRangeMap.ValueRange
					if lastRangeMap.Value == nil {
						editRange.End = pos
					}
					candidateList = valueCandidates(values, ilsp.HCLRangeToLSP(editRange))
					logger.Printf("received allowed keys: %#v", values)
				case common.ContainsPos(lastRangeMap.ValueRange, pos):
					// input a property
					keys := azure.ListProperties(def, rangeMaps)
					if len(rangeMaps) == 1 {
						keys = ignorePulledOutProperties(keys)
					}
					logger.Printf("received allowed keys: %#v", keys)
					editRange := ilsp.HCLRangeToLSP(hcl.Range{Start: pos, End: pos, Filename: filename})
					candidateList = keyCandidates(keys, editRange)
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

func ignorePulledOutProperties(input []azure.Property) []azure.Property {
	res := make([]azure.Property, 0)
	// ignore properties pulled out from body
	for _, p := range input {
		if p.Name != "name" && p.Name != "location" && p.Name != "identity" && p.Name != "tags" {
			res = append(res, p)
		}
	}
	return res
}
