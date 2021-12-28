package main_test

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"log"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	file, dialog := hclsyntax.ParseConfig([]byte(config), "", hcl.InitialPos)


	pos := hcl.Pos{Line: 91, Column: 23, Byte: 1910}

	block, err := blockAtPos(file, pos)

	filename := "main.tf"
	if err == nil && block != nil {
		if len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
			attribute, err := attributeAtPos(block, pos)
			if err == nil && attribute != nil {
				switch attribute.Name {
				case "type":
					value := toLiteral(attribute.Expr)
					if value == nil {
						fmt.Println("nil")
					} else if strings.Contains(*value, "@") {
						fmt.Println(*value)
					}
					candidates := typeCandidates()
					for i := range candidates {
						candidates[i].TextEdit.Range = hcl.Range{
							Filename: filename,
							End:      pos,
							Start:    attribute.Expr.Range().Start,
						}
					}
					break
				case "body":
					typeAttr := attributeWithName(block, "type")
					_ = typeAttr.Name
					if r, err := rangeOfJsonEncodeBody(attribute.Expr); err == nil {
						temp := ([]byte(config))[r.Start.Byte:r.End.Byte]
						tokens, dialog1 := hclsyntax.LexExpression(temp, "", r.Start)
						path := buildRangeMap(tokens)
						_, _ = path, dialog1
					}

					fmt.Println("")
					break
				}
			}
		}
	}
	fmt.Println(file, dialog, pos)
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

func typeCandidates() []lang.Candidate {
	types := []string{
		"Microsoft.MachineLearningServices/workspaces@2021-07-01",
		"Microsoft.MachineLearningServices/workspaces/computes@2021-07-1",
	}

	candidates := make([]lang.Candidate, 0)
	for _, resourceType := range types {
		candidates = append(candidates, lang.Candidate{
			Label:        fmt.Sprintf(`"%s"`, resourceType),
			Description:  lang.PlainText("this is description for " + resourceType),
			Detail:       "string",
			IsDeprecated: false,
			TextEdit: lang.TextEdit{
				NewText: fmt.Sprintf(`"%s"`, resourceType),
				Snippet: fmt.Sprintf(`"%s"`, resourceType),
			},
			Kind:           lang.AttributeCandidateKind,
			TriggerSuggest: true,
		})
	}
	return candidates
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

func attributeAtPos(block *hclsyntax.Block, pos hcl.Pos) (*hclsyntax.Attribute, error) {
	if block == nil {
		return nil, nil
	}

	for _, attr := range block.Body.Attributes {
		if ContainsPos(attr.SrcRange, pos) {
			return attr, nil
		}

		if ContainsPos(attr.Expr.Range(), pos) {
			return attr, nil
		}
	}

	return nil, nil
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
	Children map[string]*RangeMap
	KeyRange, ValueRange, EqualRange hcl.Range
	Value interface{}
	Key interface{}
}

func (rm RangeMap) GetRange() hcl.Range  {
	return RangeOver(rm.KeyRange, rm.ValueRange)
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
	if a.Start.Line > b.Start.Line || a.Start.Line == b.Start.Line && a.Start.Column > b.Start.Column {
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

type State struct {
	CurrentRangeMap *RangeMap
	Key *string
	Value *string
	KeyRange, ValueRange, EqualRange hcl.Range
	Index *int
	ExpectKey bool
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
		Key: String("dummy"),
		Index: nil,
		Value: nil,
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
				Key: key,
				KeyRange: state.KeyRange,
				EqualRange: state.EqualRange,
				Children: make(map[string]*RangeMap),
			}
			stack[top].CurrentRangeMap.Children[key] = rangeMap
			stack[top].CurrentRangeMap.ValueRange = RangeOver(stack[top].CurrentRangeMap.ValueRange, token.Range)
			stack[top].ExpectKey = stack[top].Index == nil // if there's an index, then expect value because p = [{}]
			stack = append(stack, State{
				ExpectKey: true,
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
					Key: key,
					Value: state.Value,
					KeyRange: state.KeyRange,
					ValueRange: state.ValueRange,
					EqualRange: state.EqualRange,
				}
			}
			stack[top].CurrentRangeMap.ValueRange = RangeOver(stack[top].CurrentRangeMap.ValueRange, token.Range)
			stack = stack[0: top]
			break

		case hclsyntax.TokenOBrack: // [
			top := len(stack) - 1
			state := stack[top]
			if state.ExpectKey {
				log.Printf("[WARN] expect key but got [")
			}
			key := state.GetCurrentKey()
			rangeMap := &RangeMap{
				Key: key,
				KeyRange: state.KeyRange,
				EqualRange: state.EqualRange,
				ValueRange: token.Range,
				Children: make(map[string]*RangeMap),
			}
			stack[top].CurrentRangeMap.Children[key] = rangeMap
			stack[top].ExpectKey = true
			stack = append(stack, State{
				ExpectKey: false,
				Key: state.Key,
				Index: Int(0),
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
					Key: key,
					Value: state.Value,
					ValueRange: state.ValueRange,
				}
			}
			stack[top].CurrentRangeMap.ValueRange = RangeOver(stack[top].CurrentRangeMap.ValueRange, token.Range)
			stack = stack[0: top]
			top = len(stack) - 1
			break
		case hclsyntax.TokenIdent:
			top := len(stack) - 1
			state := stack[top]
			if state.ExpectKey {
				stack[top].Key = String(string(token.Bytes))
				stack[top].KeyRange = token.Range
				stack[top].Index = nil
				stack[top].Value = nil
				stack[top].ExpectKey = false
			} else {
				if stack[top].Value == nil {
					stack[top].Value = pointerString(string(token.Bytes))
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
						Key: key,
						Value: state.Value,
						KeyRange: state.KeyRange,
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
						Key: key,
						Value: state.Value,
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
					stack[top].Value = pointerString(string(token.Bytes))
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

	data, _ := json.Marshal(root)
	fmt.Println(string(data))
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

func pointerString(input string) *string{
	return &input
}

func getNextIndex(key string, rangeMap *RangeMap) int {
	index := 0
	for _ = range rangeMap.Children {
		if _, ok := rangeMap.Children[fmt.Sprintf("%s.%d", key, index)]; ok {
			return index
		}
		index++
	}
	return index
}

func ContainsPos(r hcl.Range, pos hcl.Pos) bool {
	afterStart := pos.Line > r.Start.Line || pos.Line == r.Start.Line && pos.Column >= r.Start.Column
	beforeEnd := pos.Line < r.End.Line || pos.Line == r.End.Line && pos.Column <= r.End.Column
	return afterStart && beforeEnd
}

func Bool(input bool) *bool {
	return &input
}

func Int(input int) *int {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func Int64(input int64) *int64 {
	return &input
}

func Float(input float64) *float64 {
	return &input
}

func String(input string) *string {
	return &input
}

const config = `terraform {
  required_providers {
    azurerm-restapi = {
      source = "Azure/azurerm-restapi"
    }
  }
}

provider "azurerm" {
  features {}
}


locals {
  default_tags = {
    P1 = "provider"
    P2 = "provider"
  }
}
provider "azurerm-restapi" {
  schema_validation_enabled = true
  default_tags = local.default_tags
}

resource "azurerm_resource_group" "test" {
  name     = "henglu1220-resource-group"
  location = "west europe"
}



resource "azurerm_container_registry" "test1" {
  name                = "acctesthenglu1220"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false
}

resource "azurerm-restapi_resource" "test2" {
  resource_id = "${azurerm_container_registry.test1.id}/scopeMaps/acctesthenglu1220"
  type        = "Microsoft.ContainerRegistry/registries/scopeMaps@2020-11-01-preview"
  body        = <<BODY
   {
      "properties": {
        "description": "Developer Scopes",
        "actions": [
          "repositories/testrepo/content/read"
        ]
      }
    }
  BODY
  
}

resource "azurerm-restapi_resource" "test" {
  resource_id = "${azurerm_resource_group.test.id}/providers/Microsoft.Automation/automationAccounts/henglu1220"
  type        = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

  tags = merge(local.default_tags, {
    P2 = "resource"
    P3 = "resource"
  })

}
resource "azurerm-restapi_resource" "test3" {
  type  = "Microsoft.MachineLearningServices/workspaces@2021-07-01"
  body1 = <<BODY
  {
    "p1": "value1",
    "p2": "${var.input}"
  }
  BODY
  body = jsonencode({
    network_rule_set = [
    {
      default_action = "Deny"
      ip_rule = [
        {
             = var.action
          ip_range = 
        },
        {
          action   = var.action
          ip_range = "2.2.2.2/32"
        }
      ]
      virtual_network = []
    }
  ]
  })
}

`
/*



 */