package main_test

import (
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/hover"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/complete"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/validate"
)

func TestValidate(t *testing.T) {
	_, _ = validate.ValidateFile([]byte(config), filename)
}

func TestHover(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	pos := hcl.Pos{Line: 89, Column: 10, Byte: 1910}
	hover.HoverAtPos([]byte(config), filename, pos, logger)
}

func TestCompletion(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	pos := hcl.Pos{Line: 75, Column: 6, Byte: 1910}
	complete.CandidatesAtPos([]byte(config), filename, pos, logger)
}

const filename = "main.tf"
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

provider "azurerm-restapi" {
  schema_validation_enabled = true
  default_tags              = local.default_tags
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
  type = "Microsoft.Containerregistry/Registries@2021-06-01-preview"
  body = jsonencode({
    sku = {
      p1 = 
		= v2
      name = "Classic"
    }
    properties = {
      adminUserEnabled = ""
      
    }
  })
}

resource "azurerm-restapi_resource" "test4" {
  type = "Microsoft.Machinelearningservices/Workspaces/Computes@2021-07-01"
  body = jsonencode({
    properties = {
      properties = {
         address = "2.2.2.2"
         administratorAccount = {
           password = ""
         }
      }
      computeType = "ComputeInstance"
    }
  })
}
`

/*



 */
