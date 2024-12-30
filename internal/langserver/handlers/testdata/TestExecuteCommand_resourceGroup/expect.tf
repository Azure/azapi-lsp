/*
Note: This is a generated HCL content from the JSON input which is based on the latest API version available.
To import the resource, please run the following command:
terraform import azapi_resource.resourceGroup /subscriptions/0000000/resourceGroups/app-rg5wg5lo7sa2hxrh6bs6?api-version=2024-11-01

Or add the below config:
import {
  id = "/subscriptions/0000000/resourceGroups/app-rg5wg5lo7sa2hxrh6bs6?api-version=2024-11-01"
  to = azapi_resource.resourceGroup
}
*/

resource "azapi_resource" "resourceGroup" {
  type      = "Microsoft.Resources/resourceGroups@2024-11-01"
  parent_id = "/subscriptions/0000000"
  name      = "app-rg5wg5lo7sa2hxrh6bs6"
  location  = "eastus"
  body = {
    properties = {}
  }
  tags = {
    Creator     = "jiasli@microsoft.com"
    DateCreated = "2024-05-07T12:24:37Z"
  }
}
