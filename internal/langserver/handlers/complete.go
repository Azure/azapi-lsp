package handlers

import (
	"context"
	"fmt"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
	"log"
	"strings"

	"github.com/ms-henglu/azurerm-restapi-lsp/internal/utils"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
	"github.com/zclconf/go-cty/cty"
)

func (svc *service) TextDocumentComplete(ctx context.Context, params lsp.CompletionParams) (lsp.CompletionList, error) {
	var list lsp.CompletionList

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return list, err
	}

	cc, err := ilsp.ClientCapabilities(ctx)
	if err != nil {
		return list, err
	}

	doc, err := fs.GetDocument(ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI))
	if err != nil {
		return list, err
	}

	d, err := svc.decoderForDocument(ctx, doc)
	if err != nil {
		return list, err
	}

	expFeatures, err := lsctx.ExperimentalFeatures(ctx)
	if err != nil {
		return list, err
	}

	d.PrefillRequiredFields = expFeatures.PrefillRequiredFields

	fPos, err := ilsp.FilePositionFromDocumentPosition(params.TextDocumentPositionParams, doc)
	if err != nil {
		return list, err
	}

	svc.logger.Printf("Looking for candidates at %q -> %#v", doc.Filename(), fPos.Position())

	_, _ = d.CandidatesAtPos(doc.Filename(), fPos.Position())

	data, err := doc.Text()
	if err != nil {
		return list, err
	}
	file, _ := hclsyntax.ParseConfig(data, doc.Filename(), hcl.InitialPos)
	block, err := blockAtPos(file, fPos.Position())

	candidateList := make([]lsp.CompletionItem, 0)
	if err == nil && block != nil && len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
		if attribute := attributeAtPos(block, fPos.Position()); attribute != nil {
			switch attribute.Name {
			case "type":
				expRange := attribute.Expr.Range()
				if expRange.Start.Line != expRange.End.Line && expRange.End.Column == 1 && expRange.End.Line-1 == fPos.Position().Line {
					expRange.End = fPos.Position()
				}
				r := ilsp.HCLRangeToLSP(expRange)
				prefix := toLiteral(attribute.Expr)
				candidateList = typeCandidates(prefix, r)
				break
			case "body":
				typeAttr := attributeWithName(block, "type")
				if typeAttr == nil {
					break
				}

				typeValue := toLiteral(typeAttr.Expr)
				if typeValue == nil {
					break
				}

				parts := strings.Split(*typeValue, "@")
				if len(parts) != 2 {
					break
				}

				def, _ := azure.GetResourceDefinition(parts[0], parts[1])
				if def == nil {
					break
				}

				r, err := rangeOfJsonEncodeBody(attribute.Expr)
				if err != nil {
					break
				}

				tokens, _ := hclsyntax.LexExpression(data[r.Start.Byte:r.End.Byte], "", r.Start)
				rangeMap := buildRangeMap(tokens)

				pos := fPos.Position()
				targetRangeMap := rangeMapOfPos(rangeMap, pos)
				lastRangeMap := targetRangeMap[len(targetRangeMap)-1]

				key := ""
				for _, rangeMap := range targetRangeMap {
					key = fmt.Sprintf("%s.%v", key, rangeMap.Key)
				}
				start := len(".<nil>.dummy.")
				end := strings.LastIndex(key, ".")
				if start < end && start < len(key) && end > 0{
					key = key[start : end]
				} else {
					break
				}

				svc.logger.Printf("key: %v\n", key)
				switch {
				case ContainsPos(lastRangeMap.KeyRange, pos):
					// input a property
					editRange := ilsp.HCLRangeToLSP(lastRangeMap.KeyRange)
					keys := azure.ListProperties(def, key)
					svc.logger.Printf("received allowed keys: %#v", keys)
					candidateList = keyCandidates(keys, editRange)
					break
				case !lastRangeMap.KeyRange.Empty() && !lastRangeMap.EqualRange.Empty():
					// input property =
					editRange := ilsp.HCLRangeToLSP(lastRangeMap.ValueRange)
					values := azure.ListValues(def, key)
					candidateList = valueCandidates(values, editRange)
				}

				break
			}
		}
	}

	candidates := lsp.CompletionList{
		IsIncomplete: false,
		Items:        candidateList,
	}
	svc.logger.Printf("received candidates: %#v", candidates)
	_ = cc
	return candidates, err
}

func keyCandidates(keys []string, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	for _, key := range keys {
		content := key
		candidates = append(candidates, lsp.CompletionItem{
			Label: content,
			Kind:  lsp.ValueCompletion,
			Documentation: lsp.MarkupContent{
				Kind:  lsp.MarkupKind("markdown"),
				Value: fmt.Sprintf("Property: `%s`  \n", content),
			},
			SortText:         content,
			InsertTextFormat: lsp.SnippetTextFormat,
			InsertTextMode:   lsp.AdjustIndentation,
			TextEdit: &lsp.TextEdit{
				Range:   r,
				NewText: fmt.Sprintf(`%s = $0`, content),
			},
			Command: &lsp.Command{
				Command: "editor.action.triggerSuggest",
				Title:   "Suggest",
			},
		})
	}
	return candidates
}

func valueCandidates(values []string, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	for _, value := range values {
		content := value
		candidates = append(candidates, lsp.CompletionItem{
			Label: fmt.Sprintf(`"%s"`, content),
			Kind:  lsp.ValueCompletion,
			Documentation: lsp.MarkupContent{
				Kind:  lsp.MarkupKind("markdown"),
				Value: fmt.Sprintf("Value: `%s`  \n", content),
			},
			SortText:         content,
			InsertTextFormat: lsp.PlainTextTextFormat,
			InsertTextMode:   lsp.AdjustIndentation,
			TextEdit: &lsp.TextEdit{
				Range:   r,
				NewText: fmt.Sprintf(`"%s"`, content),
			},
		})
	}
	return candidates
}

func typeCandidates(prefix *string, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	if prefix == nil || !strings.Contains(*prefix, "@") {
		for resourceType := range azure.GetAzureSchema().Resources {
			resourceType = strings.Title(strings.ToLower(resourceType))
			candidates = append(candidates, lsp.CompletionItem{
				Label: fmt.Sprintf(`"%s"`, resourceType),
				Kind:  lsp.ValueCompletion,
				Documentation: lsp.MarkupContent{
					Kind:  lsp.MarkupKind("markdown"),
					Value: fmt.Sprintf("Type: `%s`  \n", resourceType),
				},
				SortText:         resourceType,
				InsertTextFormat: lsp.SnippetTextFormat,
				InsertTextMode:   lsp.AdjustIndentation,
				TextEdit: &lsp.TextEdit{
					Range:   r,
					NewText: fmt.Sprintf(`"%s@$0"`, resourceType),
				},
				Command: &lsp.Command{
					Command: "editor.action.triggerSuggest",
					Title:   "Suggest",
				},
			})
		}
	} else {
		resourceType := (*prefix)[0:strings.Index(*prefix, "@")]
		if resource, ok := azure.GetAzureSchema().Resources[strings.ToUpper(resourceType)]; ok {
			for _, resourceDef := range resource.Definitions {
				apiVersion := resourceDef.ApiVersion
				candidates = append(candidates, lsp.CompletionItem{
					Label: fmt.Sprintf(`"%s@%s"`, resourceType, apiVersion),
					Kind:  lsp.ValueCompletion,
					Documentation: lsp.MarkupContent{
						Kind:  lsp.MarkupKind("markdown"),
						Value: fmt.Sprintf("Type: `%s`  \nAPI Version: `%s`", resourceType, apiVersion),
					},
					SortText:         apiVersion,
					InsertTextFormat: lsp.PlainTextTextFormat,
					InsertTextMode:   lsp.AdjustIndentation,
					TextEdit: &lsp.TextEdit{
						Range:   r,
						NewText: fmt.Sprintf(`"%s@%s"`, resourceType, apiVersion),
					},
				})
			}
		}
	}

	return candidates
}

func toLiteral(expression hclsyntax.Expression) *string {
	value, dialog := expression.Value(&hcl.EvalContext{})
	if dialog != nil && dialog.HasErrors() {
		return nil
	}
	if value.Type() == cty.String && !value.IsNull() && value.IsKnown() {
		v := value.AsString()
		return &v
	}
	return nil
}

func blockAtPos(file *hcl.File, pos hcl.Pos) (*hclsyntax.Block, error) {
	body, isHcl := file.Body.(*hclsyntax.Body)
	if !isHcl {
		return nil, fmt.Errorf("file is not hcl")
	}

	for _, b := range body.Blocks {
		if ContainsPos(b.Range(), pos) {
			return b, nil
		}
	}
	return nil, nil
}

func attributeAtPos(block *hclsyntax.Block, pos hcl.Pos) *hclsyntax.Attribute {
	if block == nil {
		return nil
	}

	for _, attr := range block.Body.Attributes {
		if ContainsPos(attr.SrcRange, pos) {
			return attr
		}
		if ContainsPos(attr.Expr.Range(), pos) {
			return attr
		}
	}

	return nil
}

func attributeWithName(block *hclsyntax.Block, name string) *hclsyntax.Attribute {
	if block == nil {
		return nil
	}
	for _, attr := range block.Body.Attributes {
		if attr.Name == name {
			return attr
		}
	}
	return nil
}

func rangeOfJsonEncodeBody(expression hclsyntax.Expression) (*hcl.Range, error) {
	if expression == nil {
		return nil, nil
	}
	if funcCallExpr, ok := expression.(*hclsyntax.FunctionCallExpr); ok {
		if funcCallExpr.Name != "jsonencode" {
			return nil, fmt.Errorf("expression is not funcation jsonencode")
		}
		if len(funcCallExpr.Args) != 1 {
			return nil, fmt.Errorf("invalid args length for jsonencode")
		}
		r := funcCallExpr.Args[0].Range()
		return &r, nil
	} else {
		return nil, fmt.Errorf("expression is not funcation call expression")
	}
}

type RangeMap struct {
	Children                         map[string]*RangeMap
	KeyRange, ValueRange, EqualRange hcl.Range
	Value                            interface{}
	Key                              interface{}
}

func (rm RangeMap) GetRange() hcl.Range {
	return RangeOver(RangeOver(rm.KeyRange, rm.ValueRange), rm.EqualRange)
}

type State struct {
	CurrentRangeMap                  *RangeMap
	Key                              *string
	Value                            *string
	KeyRange, ValueRange, EqualRange hcl.Range
	Index                            *int
	ExpectKey                        bool
}

func (s State) GetCurrentKey() string {
	key := "key_placeholder"
	if s.Key != nil {
		key = *s.Key
	}
	if s.Index != nil {
		return fmt.Sprintf("%s.%d", key, *s.Index)
	}
	return key
}

func buildRangeMap(tokens hclsyntax.Tokens) *RangeMap {
	stack := make([]State, 0)
	stack = append(stack, State{
		CurrentRangeMap: &RangeMap{
			Children: make(map[string]*RangeMap),
		},
		Key:       utils.String("dummy"),
		Index:     nil,
		Value:     nil,
		ExpectKey: false,
	})

	for _, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenOBrace: // {
			top := len(stack) - 1
			state := stack[top]
			key := state.GetCurrentKey()
			if state.ExpectKey {
				log.Printf("[WARN] expect key but got {")
			}
			rangeMap := &RangeMap{
				Key:        key,
				KeyRange:   state.KeyRange,
				ValueRange: token.Range,
				EqualRange: state.EqualRange,
				Children:   make(map[string]*RangeMap),
			}
			stack[top].CurrentRangeMap.Children[key] = rangeMap
			stack[top].ExpectKey = stack[top].Index == nil // if there's an index, then expect value because p = [{}]
			stack = append(stack, State{
				ExpectKey:       true,
				CurrentRangeMap: rangeMap,
			})
			break
		case hclsyntax.TokenCBrace: // }
			top := len(stack) - 1
			state := stack[top]
			if !state.ExpectKey {
				log.Printf("[WARN] expect value but got }")
				key := state.GetCurrentKey()
				stack[top].CurrentRangeMap.Children[key] = &RangeMap{
					Key:        key,
					Value:      state.Value,
					KeyRange:   state.KeyRange,
					ValueRange: state.ValueRange,
					EqualRange: state.EqualRange,
				}
			}
			stack[top].CurrentRangeMap.ValueRange = RangeOver(stack[top].CurrentRangeMap.ValueRange, token.Range)
			stack = stack[0:top]
			break

		case hclsyntax.TokenOBrack: // [
			top := len(stack) - 1
			state := stack[top]
			if state.ExpectKey {
				log.Printf("[WARN] expect key but got [")
			}
			key := state.GetCurrentKey()
			rangeMap := &RangeMap{
				Key:        key,
				KeyRange:   state.KeyRange,
				EqualRange: state.EqualRange,
				ValueRange: token.Range,
				Children:   make(map[string]*RangeMap),
			}
			stack[top].CurrentRangeMap.Children[key] = rangeMap
			stack[top].ExpectKey = true
			stack = append(stack, State{
				ExpectKey:       false,
				Key:             state.Key,
				Index:           utils.Int(0),
				CurrentRangeMap: rangeMap,
			})
			break
		case hclsyntax.TokenCBrack: // ]
			top := len(stack) - 1
			state := stack[top]
			if _, ok := state.CurrentRangeMap.Children[state.GetCurrentKey()]; !ok {
				log.Printf("[WARN] expect value but got }")
				key := state.GetCurrentKey()
				stack[top].CurrentRangeMap.Children[key] = &RangeMap{
					Key:        key,
					Value:      state.Value,
					ValueRange: state.ValueRange,
				}
			}
			stack[top].CurrentRangeMap.ValueRange = RangeOver(stack[top].CurrentRangeMap.ValueRange, token.Range)
			stack = stack[0:top]
			top = len(stack) - 1
			break
		case hclsyntax.TokenIdent:
			top := len(stack) - 1
			state := stack[top]
			if state.ExpectKey {
				stack[top].Key = utils.String(string(token.Bytes))
				stack[top].KeyRange = token.Range
				stack[top].Index = nil
				stack[top].Value = nil
				stack[top].ValueRange = hcl.Range{}
				stack[top].EqualRange = hcl.Range{}
				stack[top].ExpectKey = false
			} else {
				if stack[top].Value == nil {
					stack[top].Value = utils.String(string(token.Bytes))
				} else {
					*stack[top].Value = *stack[top].Value + string(token.Bytes)
				}
				stack[top].ValueRange = RangeOver(stack[top].ValueRange, token.Range)
			}
			break
		case hclsyntax.TokenEqual: // =
			top := len(stack) - 1
			state := stack[top]
			if state.ExpectKey {
				log.Printf("[WARN] expect key but got =")
				stack[top].ExpectKey = false
			}
			stack[top].EqualRange = token.Range
			break
		case hclsyntax.TokenNewline:
			top := len(stack) - 1
			state := stack[top]
			if !state.ExpectKey {
				if stack[top].Index == nil {
					key := state.GetCurrentKey()
					stack[top].CurrentRangeMap.Children[key] = &RangeMap{
						Key:        key,
						Value:      state.Value,
						KeyRange:   state.KeyRange,
						ValueRange: state.ValueRange,
						EqualRange: state.EqualRange,
					}
					stack[top].Key = nil
					stack[top].Value = nil
					stack[top].ExpectKey = true
				}
			}
			break
		case hclsyntax.TokenComma:
			top := len(stack) - 1
			state := stack[top]
			if !state.ExpectKey {
				if state.Value != nil {
					key := state.GetCurrentKey()
					stack[top].CurrentRangeMap.Children[key] = &RangeMap{
						Key:        key,
						Value:      state.Value,
						ValueRange: state.ValueRange,
					}
				}
				*stack[top].Index++
				stack[top].Value = nil
				stack[top].ValueRange = hcl.Range{}
			} else {
				log.Printf("[WARN] expect key but got ,")
			}
			break
		default:
			top := len(stack) - 1
			state := stack[top]
			if !state.ExpectKey {
				if stack[top].Value == nil {
					stack[top].Value = utils.String(string(token.Bytes))
				} else {
					*stack[top].Value = *stack[top].Value + string(token.Bytes)
				}
				stack[top].ValueRange = RangeOver(stack[top].ValueRange, token.Range)
			}
			break
		}
	}

	if len(stack) == 0 {
		return nil
	}

	root := stack[0].CurrentRangeMap
	updateValueRange(root)
	return root
}

func updateValueRange(rangeMap *RangeMap) {
	if rangeMap == nil {
		return
	}
	for _, v := range rangeMap.Children {
		updateValueRange(v)
		rangeMap.ValueRange = RangeOver(rangeMap.ValueRange, v.GetRange())
	}
}

func rangeMapOfPos(rangeMap *RangeMap, pos hcl.Pos) []*RangeMap {
	if rangeMap == nil {
		return nil
	}
	res := make([]*RangeMap, 0)
	if ContainsPos(rangeMap.GetRange(), pos) {
		res = append(res, rangeMap)
		for _, child := range rangeMap.Children {
			if arr := rangeMapOfPos(child, pos); len(arr) != 0 {
				res = append(res, arr...)
				break
			}
		}
		return res
	}
	return nil
}

func ContainsPos(r hcl.Range, pos hcl.Pos) bool {
	afterStart := pos.Line > r.Start.Line || pos.Line == r.Start.Line && pos.Column >= r.Start.Column
	beforeEnd := pos.Line < r.End.Line || pos.Line == r.End.Line && pos.Column <= r.End.Column
	return afterStart && beforeEnd
}

func RangeOver(a hcl.Range, b hcl.Range) hcl.Range {
	if a.Empty() {
		return b
	}
	if b.Empty() {
		return a
	}

	var start, end hcl.Pos
	if a.Start.Line < b.Start.Line || a.Start.Line == b.Start.Line && a.Start.Column < b.Start.Column {
		start = a.Start
	} else {
		start = b.Start
	}
	if a.End.Line > b.End.Line || a.End.Line == b.End.Line && a.End.Column > b.End.Column {
		end = a.End
	} else {
		end = b.End
	}
	return hcl.Range{
		Filename: a.Filename,
		Start:    start,
		End:      end,
	}
}
