package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azapi-lsp/internal/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type KeyValueFormat int

const (
	KeyEqualValue KeyValueFormat = iota
	QuotedKeyEqualValue
	QuotedKeyColonValue
)

type HclNode struct {
	Children                         map[string]*HclNode
	KeyRange, ValueRange, EqualRange hcl.Range
	Value                            *string
	Key                              string
	KeyValueFormat                   KeyValueFormat
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
	ExpectEqual                      bool
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
			KeyValueFormat: KeyEqualValue,
			Children:       make(map[string]*HclNode),
		},
		Key:         utils.String("dummy"),
		Index:       nil,
		Value:       nil,
		ExpectKey:   false,
		ExpectEqual: false,
	})

	for _, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenOBrace: // {
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
			state := stack[top]
			key := state.GetCurrentKey()
			if state.ExpectKey {
				log.Printf("[WARN] expect key but got {")
			}
			hclNode := &HclNode{
				Key:            key,
				KeyRange:       state.KeyRange,
				ValueRange:     token.Range,
				EqualRange:     state.EqualRange,
				Children:       make(map[string]*HclNode),
				KeyValueFormat: stack[top].CurrentHclNode.KeyValueFormat,
			}
			stack[top].CurrentHclNode.Children[key] = hclNode
			if stack[top].Index == nil {
				stack[top].ExpectKey = true
				stack[top].ExpectEqual = true
			} else {
				// if there's an index, then expect value because p = [{}]
				stack[top].ExpectKey = false
				stack[top].ExpectEqual = false
				// the hclNode is an object in an array, so use the range of "{" as the key range
				hclNode.KeyRange = token.Range
			}
			stack = append(stack, State{
				ExpectKey:      true,
				ExpectEqual:    true,
				CurrentHclNode: hclNode,
			})
		case hclsyntax.TokenCBrace: // }
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
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
			if top < 0 {
				return nil
			}
			state := stack[top]
			if state.ExpectKey {
				log.Printf("[WARN] expect key but got [")
			}
			if state.Value != nil && *state.Value != "" {
				updateStateValue(&stack[top], token)
				break
			}
			key := state.GetCurrentKey()
			hclNode := &HclNode{
				Key:            key,
				KeyRange:       state.KeyRange,
				EqualRange:     state.EqualRange,
				ValueRange:     token.Range,
				KeyValueFormat: state.CurrentHclNode.KeyValueFormat,
				Children:       make(map[string]*HclNode),
			}
			stack[top].CurrentHclNode.Children[key] = hclNode
			stack[top].ExpectKey = true
			stack[top].ExpectEqual = stack[top].ExpectKey
			stack = append(stack, State{
				ExpectKey:      false,
				ExpectEqual:    false,
				Key:            state.Key,
				Index:          utils.Int(0),
				CurrentHclNode: hclNode,
			})
		case hclsyntax.TokenCBrack: // ]
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
			state := stack[top]
			if state.Value != nil && *state.Value != "" && strings.Contains(*state.Value, "[") {
				updateStateValue(&stack[top], token)
				break
			}
			if _, ok := state.CurrentHclNode.Children[state.GetCurrentKey()]; !ok {
				log.Printf("[WARN] expect value but got }")
				// avoid adding an empty element
				if !state.ValueRange.Empty() {
					key := state.GetCurrentKey()
					stack[top].CurrentHclNode.Children[key] = &HclNode{
						Key:        key,
						Value:      state.Value,
						ValueRange: state.ValueRange,
					}
				}
			}
			stack[top].CurrentHclNode.ValueRange = RangeOver(stack[top].CurrentHclNode.ValueRange, token.Range)
			stack = stack[0:top]
		case hclsyntax.TokenQuotedLit:
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
			if stack[top].ExpectKey {
				foundKey(&stack[top], token)
				stack[top].CurrentHclNode.KeyValueFormat = QuotedKeyEqualValue
			} else {
				updateStateValue(&stack[top], token)
			}
		case hclsyntax.TokenIdent:
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
			if stack[top].ExpectKey {
				foundKey(&stack[top], token)
			} else {
				updateStateValue(&stack[top], token)
			}
		case hclsyntax.TokenColon: // :
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
			if stack[top].ExpectEqual {
				if stack[top].ExpectKey {
					log.Printf("[WARN] expect key but got =")
					stack[top].ExpectKey = false
				}

				stack[top].CurrentHclNode.KeyValueFormat = QuotedKeyColonValue
				stack[top].EqualRange = token.Range
				stack[top].ExpectEqual = false
			} else {
				updateStateValue(&stack[top], token)
			}
		case hclsyntax.TokenEqual: // =
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
			if stack[top].ExpectEqual {
				if stack[top].ExpectKey {
					log.Printf("[WARN] expect key but got =")
					stack[top].ExpectKey = false
				}
				stack[top].EqualRange = token.Range
				stack[top].ExpectEqual = false
			} else {
				updateStateValue(&stack[top], token)
			}
		case hclsyntax.TokenNewline:
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
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
					stack[top].ExpectEqual = stack[top].ExpectKey
				}
			}
		case hclsyntax.TokenComma:
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
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
				if top < 0 {
					return nil
				}
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
						stack[top].ExpectEqual = stack[top].ExpectKey
					}
				}
			}
		default:
			top := len(stack) - 1
			if top < 0 {
				return nil
			}
			if !stack[top].ExpectEqual {
				updateStateValue(&stack[top], token)
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

func foundKey(st *State, token hclsyntax.Token) {
	st.Key = utils.String(string(token.Bytes))
	st.KeyRange = token.Range
	st.Index = nil
	st.Value = nil
	st.ValueRange = hcl.Range{}
	st.EqualRange = hcl.Range{}
	st.ExpectKey = false
}

func updateStateValue(st *State, token hclsyntax.Token) {
	if st.Value == nil {
		st.Value = utils.String(string(token.Bytes))
	} else {
		*st.Value = *st.Value + string(token.Bytes)
	}

	st.ValueRange = RangeOver(st.ValueRange, token.Range)
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
