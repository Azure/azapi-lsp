package common

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/utils"
)

type RangeMap struct {
	Children                         map[string]*RangeMap
	KeyRange, ValueRange, EqualRange hcl.Range
	Value                            *string
	Key                              string
}

func (rm RangeMap) GetRange() hcl.Range {
	return RangeOver(RangeOver(rm.KeyRange, rm.ValueRange), rm.EqualRange)
}

func (rm RangeMap) IsValueArray() bool {
	if rm.Value != nil {
		return false
	}
	prefix := rm.Key + "."
	for _, child := range rm.Children {
		if !strings.HasPrefix(child.Key, prefix) {
			return false
		}
	}
	return true
}

func (rm RangeMap) IsValueMap() bool {
	if rm.Value != nil {
		return false
	}
	prefix := rm.Key + "."
	for _, child := range rm.Children {
		if strings.HasPrefix(child.Key, prefix) {
			return false
		}
	}
	return true
}

func RangeMapArraysOfPos(rangeMap *RangeMap, pos hcl.Pos) []*RangeMap {
	if rangeMap == nil {
		return nil
	}
	if ContainsPos(rangeMap.GetRange(), pos) {
		res := make([]*RangeMap, 0)
		if rangeMap.Key != "" {
			res = append(res, rangeMap)
		}
		for _, child := range rangeMap.Children {
			if arr := RangeMapArraysOfPos(child, pos); len(arr) != 0 {
				return append(res, arr...)
			}
		}
		return res
	}
	return nil
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

func BuildRangeMap(tokens hclsyntax.Tokens) *RangeMap {
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
		case hclsyntax.TokenEqual: // =
			top := len(stack) - 1
			state := stack[top]
			if state.ExpectKey {
				log.Printf("[WARN] expect key but got =")
				stack[top].ExpectKey = false
			}
			stack[top].EqualRange = token.Range
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
		case hclsyntax.TokenComma:
			top := len(stack) - 1
			state := stack[top]
			if !state.ExpectKey {
				if state.Index == nil {
					log.Printf("[WARN] unexpected symbol: ,")
					break
				}
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
		}
	}

	if len(stack) == 0 {
		return nil
	}

	root := stack[0].CurrentRangeMap
	updateValueRange(root)
	fixEmptyValueRange(root)
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

func fixEmptyValueRange(rangeMap *RangeMap) {
	if rangeMap == nil {
		return
	}
	if rangeMap.Children == nil {
		if !rangeMap.KeyRange.Empty() && !rangeMap.EqualRange.Empty() && rangeMap.ValueRange.Empty() {
			line := rangeMap.EqualRange.End.Line
			column := rangeMap.EqualRange.End.Column
			rangeMap.ValueRange = hcl.Range{
				Start: hcl.Pos{Line: line, Column: column + 1, Byte: rangeMap.EqualRange.End.Byte},
				End:   hcl.Pos{Line: line + 1, Byte: rangeMap.EqualRange.End.Byte + 1},
			}
		}
	} else {
		for _, child := range rangeMap.Children {
			fixEmptyValueRange(child)
		}
	}
}
