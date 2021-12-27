package main_test

import (
	"fmt"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	file, dialog := hclsyntax.ParseConfig([]byte(config), "", hcl.InitialPos)
	pos := hcl.Pos{Line: 80, Column: 19, Byte: 1910}

	block, err := blockAtPos(file, pos)


	filename := "main.tf"
	if err == nil && block != nil {
		if len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
			attribute, err := attributeAtPos(block, pos)
			if err == nil && attribute != nil {
				switch attribute.Name {
				case "type":
					candidates := typeCandidates()
					for i, _ := range candidates {
						candidates[i].TextEdit.Range = hcl.Range{
							Filename: filename,
							End:      pos,
							Start: attribute.Expr.Range().Start,
						}
					}
					break
				case "body":
					break
				}
			}
		}
	}
	fmt.Println(file, dialog, pos)
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

func attributeAtPos(block *hclsyntax.Block, pos hcl.Pos) (*hclsyntax.Attribute, error) {
	if block == nil {
		return nil, nil
	}

	for _, attr := range block.Body.Attributes {
		if ContainsPos(attr.SrcRange, pos) {
			return attr, nil
		}
	}

	return nil, nil
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
  type = "Microsof"
  body = jsonencode({
    p1 = "value1"
    p2 = {
        p3 = {
            p4 = d
        }
    }
  })
}

`
