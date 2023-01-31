package complete

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/Azure/azapi-lsp/internal/langserver/handlers/tfschema"
	ilsp "github.com/Azure/azapi-lsp/internal/lsp"
	"github.com/Azure/azapi-lsp/internal/parser"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func CandidatesAtPos(data []byte, filename string, pos hcl.Pos, logger *log.Logger) []lsp.CompletionItem {
	file, _ := hclsyntax.ParseConfig(data, filename, hcl.InitialPos)

	body, isHcl := file.Body.(*hclsyntax.Body)
	if !isHcl {
		logger.Printf("file is not hcl")
		return nil
	}
	block := parser.BlockAtPos(body, pos)
	candidateList := make([]lsp.CompletionItem, 0)
	if block != nil && len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azapi") {
		resourceName := fmt.Sprintf("%s.%s", block.Type, block.Labels[0])
		resource := tfschema.GetResourceSchema(resourceName)
		if resource == nil {
			return candidateList
		}
		if attribute := parser.AttributeAtPos(block, pos); attribute != nil {
			property := resource.GetProperty(attribute.Name)
			if property != nil {
				if property.GenericCandidatesFunc != nil {
					candidateList = append(candidateList, property.GenericCandidatesFunc(data, filename, block, attribute, pos, property)...)
				} else if property.ValueCandidatesFunc != nil {
					prefix := parser.ToLiteral(attribute.Expr)
					candidateList = append(candidateList, property.ValueCandidatesFunc(prefix, editRangeFromExprRange(attribute.Expr, pos))...)
				}
			}
		} else {
			if subBody := parser.BlockAtPos(block.Body, pos); subBody != nil {
				property := resource.GetProperty(subBody.Type)
				if property == nil {
					return nil
				}
				if attribute := parser.AttributeAtPos(subBody, pos); attribute != nil {
					for _, p := range property.NestedProperties {
						if p.Name == attribute.Name && p.ValueCandidatesFunc != nil {
							candidateList = append(candidateList, p.ValueCandidatesFunc(nil, editRangeFromExprRange(attribute.Expr, pos))...)
							break
						}
					}
				}
			} else {
				//input a top level property
				editRange := lsp.Range{
					Start: ilsp.HCLPosToLSP(pos),
					End:   ilsp.HCLPosToLSP(pos),
				}
				editRange.Start.Character = 2
				candidateList = append(candidateList, tfschema.PropertiesCandidates(resource.Properties, editRange)...)
			}
		}
	}

	sort.Slice(candidateList, func(i, j int) bool { return candidateList[i].SortText < candidateList[j].SortText })
	return candidateList
}

func editRangeFromExprRange(expression hclsyntax.Expression, pos hcl.Pos) lsp.Range {
	expRange := expression.Range()
	if expRange.Start.Line != expRange.End.Line && expRange.End.Column == 1 && expRange.End.Line-1 == pos.Line {
		expRange.End = pos
	}
	return ilsp.HCLRangeToLSP(expRange)
}
