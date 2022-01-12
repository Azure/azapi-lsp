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
				targetRangeMap := common.RangeMapOfPos(rangeMap, pos)
				if len(targetRangeMap) == 0 {
					break
				}
				lastRangeMap := targetRangeMap[len(targetRangeMap)-1]

				key := common.BuildKeyFromRangeMaps(targetRangeMap)
				logger.Printf("key: %v\n", key)

				switch {
				case common.ContainsPos(lastRangeMap.KeyRange, pos):
					start := len("..dummy.")
					end := strings.LastIndex(key, ".")
					if start < end && start < len(key) && end > 0 {
						key = key[start:end]
					} else {
						key = ""
					}
					// input a property
					editRange := ilsp.HCLRangeToLSP(lastRangeMap.KeyRange)
					keys := azure.ListProperties(def, key)
					logger.Printf("state: input key, key: %s", key)
					logger.Printf("received allowed keys: %#v", keys)
					candidateList = keyCandidates(keys, editRange)
				case !lastRangeMap.KeyRange.Empty() && !lastRangeMap.EqualRange.Empty() && lastRangeMap.Children == nil:
					start := len("..dummy.")
					if start < len(key) {
						key = key[start:]
					}
					// input property =
					editRange := lastRangeMap.ValueRange
					if lastRangeMap.Value == nil {
						editRange.End = pos
					}
					values := azure.ListValues(def, key)
					candidateList = valueCandidates(values, ilsp.HCLRangeToLSP(editRange))
					logger.Printf("state: input value, key: %s", key)
					logger.Printf("received allowed keys: %#v", values)
				case common.ContainsPos(lastRangeMap.ValueRange, pos):
					start := len("..dummy.")
					if start < len(key) {
						key = key[start:]
					} else {
						key = ""
					}
					// input a property
					editRange := ilsp.HCLRangeToLSP(hcl.Range{Start: pos, End: pos, Filename: filename})
					keys := azure.ListProperties(def, key)
					logger.Printf("state: input key without prefix, key: %s", key)
					logger.Printf("received allowed keys: %#v", keys)
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
