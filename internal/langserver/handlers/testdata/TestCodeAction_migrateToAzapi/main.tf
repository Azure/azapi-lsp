resource "azurerm_resource_group" "test" {
  name     = "acctesthenglu109"
  location = "westus"
}

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2019-06-01"
  parent_id = azurerm_resource_group.test.id
  name      = "acctesthenglu109"
  location  = azurerm_resource_group.test.location
  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
  response_export_values = ["properties.name", "properties.sku.name"]
}