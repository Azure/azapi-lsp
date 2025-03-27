resource "azapi_resource" "diagnostics" {
  type      = "microsoft.insights/diagnosticSettings@2021-05-01-preview"
  name      = "somename"
  parent_id = azapi_resource.automationAccount.id
  body = jsonencode({
    properties = {
      metrics = [
        {
          category = "AllMetrics"
          enabled  = true
        },
      ]
    }
  })
}