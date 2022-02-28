resource "azapi_resource" "test" {
  name = "acctest1774"
  parent_id = azurerm_batch_account.test.id
  type = "Microsoft.DataFactory/factories@2018-06-01"
  body =
}
