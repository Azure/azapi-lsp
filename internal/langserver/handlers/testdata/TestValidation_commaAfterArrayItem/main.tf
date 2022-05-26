resource "azapi_resource" "diagnostics" {
  type      = "microsoft.insights/diagnosticSettings@2017-05-01-preview"
  name      = local.diagnosticsName
  parent_id = azapi_resource.automationAccount.id

  count = (var.enableDiagnostics) ? 1 : 0
  body = jsonencode({
    properties = {
      workspaceId      = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.logAnalyticsResourceGroup}/providers/Microsoft.OperationalInsights/workspaces/${var.logAnalyticsWorkspaceName}"
      storageAccountId = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/resourceGroups/${var.diagnosticStorageAccountResourceGroup}/providers/Microsoft.Storage/storageAccounts/${var.diagnosticStorageAccountName}"
      logs = [
        {
          category = "JobLogs"
          enabled  = true
        },
        {
          category = "JobStreams"
          enabled  = true
        },
        {
          category = "DscNodeStatus"
          enabled  = true
        },
      ]
      metrics = [
        {
          category = "AllMetrics"
          enabled  = true
        },
      ]
    }
  })
}