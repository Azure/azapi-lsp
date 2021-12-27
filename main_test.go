package main_test

import (
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
	Range hcl.Range
	Value interface{}
	Key interface{}
}

func buildRangeMap(tokens hclsyntax.Tokens) *RangeMap {
	stack := make([]*RangeMap, 0)
	dummy := &RangeMap{
		Children: make(map[string]*RangeMap),
	}
	stack = append(stack, dummy)

	keyStack := make([]string, 0)
	indexStack := make([]int, 0)
	keyStack = append(keyStack, "dummy")
	indexStack = append(indexStack, -1)

	isKey := false
	var value *string
	for _, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenOBrace:
			key := "key_placeholder"
			if isKey {
				log.Printf("[WARN] key is empty")
			} else {
				key = getKey(keyStack, indexStack)
			}
			rangeMap := &RangeMap{
				Key: key,
				Children: make(map[string]*RangeMap),
			}
			stack[len(stack) - 1].Children[key] = rangeMap
			stack = append(stack, rangeMap)
			isKey = true
			break
		case hclsyntax.TokenCBrace:
			// waiting for value, but got none
			if !isKey {
				key := getKey(keyStack, indexStack)
				stack[len(stack) - 1].Children[key] = &RangeMap{
					Key: key,
					Value: nil,
				}
				keyStack = keyStack[0: len(keyStack) - 1]
				indexStack = indexStack[0: len(indexStack) - 1]
				isKey = true
			}
			stack = stack[0: len(stack) - 1]
			keyStack = keyStack[0: len(keyStack) - 1]
			indexStack = indexStack[0: len(indexStack) - 1]
			break
		case hclsyntax.TokenOBrack:
			if len(keyStack) == 0 {
				log.Printf("[WARN] key is empty")
				keyStack = append(keyStack, "key_placeholder")
				indexStack = append(indexStack, -1)
			}
			key := keyStack[len(keyStack) - 1]
			keyStack = append(keyStack, key)
			indexStack = append(indexStack, 0)
			break
		case hclsyntax.TokenCBrack:
			if !isKey && value != nil {
				key := getKey(keyStack, indexStack)
				stack[len(stack) - 1].Children[key] = &RangeMap{
					Key: key,
					Value: value,
				}
				value = nil
				isKey = true
			}
			keyStack = keyStack[0: len(keyStack) - 2]
			indexStack = indexStack[0: len(indexStack) - 2]
			break
		case hclsyntax.TokenIdent:
			if !isKey {
				key := getKey(keyStack, indexStack)
				stack[len(stack) - 1].Children[key] = &RangeMap{
					Key: key,
					Value: value,
				}
				keyStack = keyStack[0: len(keyStack) - 1]
				indexStack = indexStack[0: len(indexStack) - 1]
				isKey = true
				value = nil
			}
			if isKey {
				keyStack = append(keyStack, string(token.Bytes))
				indexStack = append(indexStack, -1)
				isKey = false
			} else {
				log.Println("[WARN] expect value but got key")
			}
			break
		case hclsyntax.TokenEqual:
			if isKey {
				log.Printf("[WARN] key is empty")
			}
			break
		case hclsyntax.TokenNewline:
			break
		case hclsyntax.TokenComma:
			if !isKey {
				key := getKey(keyStack, indexStack)
				stack[len(stack) - 1].Children[key] = &RangeMap{
					Key: key,
					Value: value,
				}
				indexStack[len(indexStack) - 1] ++
				value = nil
			}
			break
		default:
			if !isKey {
				if value == nil {
					value = pointerString(string(token.Bytes))
				} else {
					*value = *value + string(token.Bytes)
				}
			}
			break
		}
	}

	if len(stack) == 0 {
		return nil
	}
	return stack[0]
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

func getKey(keyStack []string, indexStack []int) string {
	key := "missing_key_placeholder"
	if len(keyStack) == 0 {
		log.Printf("[WARN] key is empty")
	} else {
		key = keyStack[len(keyStack) - 1]
	}
	if len(indexStack) != 0 && indexStack[len(indexStack) - 1] != -1 {
		key = fmt.Sprintf("%s.%d", key, indexStack[len(indexStack) - 1])
	}
	return key
}

func ContainsPos(r hcl.Range, pos hcl.Pos) bool {
	afterStart := pos.Line > r.Start.Line || pos.Line == r.Start.Line && pos.Column >= r.Start.Column
	beforeEnd := pos.Line < r.End.Line || pos.Line == r.End.Line && pos.Column <= r.End.Column
	return afterStart && beforeEnd
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
    p1 = ["p1.0", "p1.1"]
    p2 = [
      {
        p3 = "p2.0.p3"
      },
      {
        p3 = "p2.0.p3"
      }
    ]
  })
}

`
