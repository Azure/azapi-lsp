resource "azapi_resource_action" "test" {
  type = "Microsoft.AppPlatform/Spring@2022-05-01-preview"
  operation = "regenerateTestKey"
  body = jsonencode({

  })
}