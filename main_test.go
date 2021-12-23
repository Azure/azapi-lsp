package main_test

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"testing"
)



func TestName(t *testing.T) {
	file, dialog := hclsyntax.ParseConfig([]byte(config), "", hcl.InitialPos)
	pos := hcl.Pos{Line: 80, Column: 4}

	if body, ok := file.Body.(*hclsyntax.Body); ok {
		blocks := body.Blocks
		block := blocks[len(blocks) - 1]
		body := block.Body
		fmt.Println(body)
	}
	fmt.Println(file, dialog, pos)
}


const config =
`terraform {
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
  body1 = <<BODY
  {
    "p1": "value1",
    "p2": "${var.input}"
  }
  BODY
  body = jsonencode({
    p1 = "value1"
    p2 = 
  })
}

`