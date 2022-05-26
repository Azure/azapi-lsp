resource "azapi_resource" "dataflow" {
  type = "Microsoft.DataFactory/factories/dataflows@2018-06-01"
  name = "hengludf0509"
  parent_id = azurerm_data_factory.test.id
  body = jsonencode({
    properties = {
      type = "Flowlet"
      typeProperties = {

      }

    }
  })
}