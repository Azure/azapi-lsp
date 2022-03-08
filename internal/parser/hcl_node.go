package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azapi-lsp/internal/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type HclNode struct {
	Children                         map[string]*HclNode
	KeyRange, ValueRange, EqualRange hcl.Range
	Value                            *string
	Key                              string
}

func (rm HclNode) GetRange() hcl.Range {
	return RangeOver(RangeOver(rm.KeyRange, rm.ValueRange), rm.EqualRange)
}

func (rm HclNode) IsValueArray() bool {
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

func (rm HclNode) IsValueMap() bool {
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

func HclNodeArraysOfPos(hclNode *HclNode, pos hcl.Pos) []*HclNode {
	if hclNode == nil {
		return nil
	}
	if ContainsPos(hclNode.GetRange(), pos) {
		res := make([]*HclNode, 0)
		if hclNode.Key != "" {
			res = append(res, hclNode)
		}
		for _, child := range hclNode.Children {
			if arr := HclNodeArraysOfPos(child, pos); len(arr) != 0 {
				return append(res, arr...)
			}
		}
		return res
	}
	return nil
}

type State struct {
	CurrentHclNode                   *HclNode
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

func BuildHclNode(tokens hclsyntax.Tokens) *HclNode {
	stack := make([]State, 0)
	stack = append(stack, State{
		CurrentHclNode: &HclNode{
			Children: make(map[string]*HclNode),
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
			hclNode := &HclNode{
				Key:        key,
				KeyRange:   state.KeyRange,
				ValueRange: token.Range,
				EqualRange: state.EqualRange,
				Children:   make(map[string]*HclNode),
			}
			stack[top].CurrentHclNode.Children[key] = hclNode
			stack[top].ExpectKey = stack[top].Index == nil // if there's an index, then expect value because p = [{}]
			stack = append(stack, State{
				ExpectKey:      true,
				CurrentHclNode: hclNode,
			})
		case hclsyntax.TokenCBrace: // }
			top := len(stack) - 1
			state := stack[top]
			if !state.ExpectKey {
				log.Printf("[WARN] expect value but got }")
				key := state.GetCurrentKey()
				stack[top].CurrentHclNode.Children[key] = &HclNode{
					Key:        key,
					Value:      state.Value,
					KeyRange:   state.KeyRange,
					ValueRange: state.ValueRange,
					EqualRange: state.EqualRange,
				}
			}
			stack[top].CurrentHclNode.ValueRange = RangeOver(stack[top].CurrentHclNode.ValueRange, token.Range)
			stack = stack[0:top]
		case hclsyntax.TokenOBrack: // [
			top := len(stack) - 1
			state := stack[top]
			if state.ExpectKey {
				log.Printf("[WARN] expect key but got [")
			}
			key := state.GetCurrentKey()
			hclNode := &HclNode{
				Key:        key,
				KeyRange:   state.KeyRange,
				EqualRange: state.EqualRange,
				ValueRange: token.Range,
				Children:   make(map[string]*HclNode),
			}
			stack[top].CurrentHclNode.Children[key] = hclNode
			stack[top].ExpectKey = true
			stack = append(stack, State{
				ExpectKey:      false,
				Key:            state.Key,
				Index:          utils.Int(0),
				CurrentHclNode: hclNode,
			})
		case hclsyntax.TokenCBrack: // ]
			top := len(stack) - 1
			state := stack[top]
			if _, ok := state.CurrentHclNode.Children[state.GetCurrentKey()]; !ok {
				log.Printf("[WARN] expect value but got }")
				key := state.GetCurrentKey()
				stack[top].CurrentHclNode.Children[key] = &HclNode{
					Key:        key,
					Value:      state.Value,
					ValueRange: state.ValueRange,
				}
			}
			stack[top].CurrentHclNode.ValueRange = RangeOver(stack[top].CurrentHclNode.ValueRange, token.Range)
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
					stack[top].CurrentHclNode.Children[key] = &HclNode{
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
					stack[top].CurrentHclNode.Children[key] = &HclNode{
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
		case hclsyntax.TokenComment:
			comment := string(token.Bytes)
			if strings.HasSuffix(comment, "\n") {
				top := len(stack) - 1
				state := stack[top]
				if !state.ExpectKey {
					if stack[top].Index == nil {
						key := state.GetCurrentKey()
						stack[top].CurrentHclNode.Children[key] = &HclNode{
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

	root := stack[0].CurrentHclNode
	updateValueRange(root)
	fixEmptyValueRange(root)
	return root
}

func updateValueRange(hclNode *HclNode) {
	if hclNode == nil {
		return
	}
	for _, v := range hclNode.Children {
		updateValueRange(v)
		hclNode.ValueRange = RangeOver(hclNode.ValueRange, v.GetRange())
	}
}

func fixEmptyValueRange(hclNode *HclNode) {
	if hclNode == nil {
		return
	}
	if hclNode.Children == nil {
		if !hclNode.KeyRange.Empty() && !hclNode.EqualRange.Empty() && hclNode.ValueRange.Empty() {
			line := hclNode.EqualRange.End.Line
			column := hclNode.EqualRange.End.Column
			hclNode.ValueRange = hcl.Range{
				Start: hcl.Pos{Line: line, Column: column + 1, Byte: hclNode.EqualRange.End.Byte},
				End:   hcl.Pos{Line: line + 1, Byte: hclNode.EqualRange.End.Byte + 1},
			}
		}
	} else {
		for _, child := range hclNode.Children {
			fixEmptyValueRange(child)
		}
	}
}
