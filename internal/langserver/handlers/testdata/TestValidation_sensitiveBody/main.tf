resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "henglu415"
  location = "westus"
  body     = {} # Resource group does not require additional properties in the body
  lifecycle {
    ignore_changes = [
      tags,
    ]
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "henglu415storage"
  location  = "westus"
  body = {
    kind = "StorageV2"
    sku = {
      name = "Premium_LRS"
    }
  }
}

ephemeral "azapi_resource_action" "listKeys" {
  type        = "Microsoft.Storage/storageAccounts@2024-01-01"
  resource_id = azapi_resource.storageAccount.id
  action      = "listKeys"
  response_export_values = {
    primary_access_key = "keys[0].value"
  }
}


resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2023-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "henglu415workspace"
  location  = "westus"
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }

}

resource "azapi_resource" "storageInsightConfig" {
  type      = "Microsoft.OperationalInsights/workspaces/storageInsightConfigs@2023-09-01"
  parent_id = azapi_resource.workspace.id
  name      = "storageInsightConfig"
  location  = "westus"
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }

  body = {
    properties = {
      storageAccount = {
        id = azapi_resource.storageAccount.id
      }
    }
  }

  sensitive_body = {
    properties = {
      storageAccount = {
        key = ephemeral.azapi_resource_action.listKeys.output.primary_access_key
      }
    }
  }
}
