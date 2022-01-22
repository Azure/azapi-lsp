resource "azurerm-restapi_resource" "test" {
  name = "acctest1774"
  parent_id = azurerm_batch_account.test.id
  type = "Microsoft.DataFactory/factories@2018-06-01"
  schema_validation_enabled
  body =  jsonencode({
    identity = {
      type = "SystemAssigned"
    }
    properties = {
      encryption = {
        identity = {

        }
      }
    }
  })
}
