/*
Note: This is a generated HCL content from the JSON input which is based on the latest API version available.
To import the resource, please run the following command:
terraform import azapi_resource.automationAccount /subscriptions/0b1f6471-1bf0-4dda-aec3-cb9272f09590/resourceGroups/acctest1660/providers/Microsoft.Automation/automationAccounts/hengluacctest5030?api-version=2023-11-01

Or add the below config:
import {
  id = "/subscriptions/0b1f6471-1bf0-4dda-aec3-cb9272f09590/resourceGroups/acctest1660/providers/Microsoft.Automation/automationAccounts/hengluacctest5030?api-version=2023-11-01"
  to = azapi_resource.automationAccount
}
*/

resource "azapi_resource" "automationAccount" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  parent_id = "/subscriptions/0b1f6471-1bf0-4dda-aec3-cb9272f09590/resourceGroups/acctest1660"
  name      = "hengluacctest5030"
  location  = "westeurope"
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      encryption = {
        identity = {
          userAssignedIdentity = null
        }
        keySource = "Microsoft.Automation"
      }
      sku = {
        capacity = null
        family   = null
        name     = "Basic"
      }
    }
  }
}
