resource "azapi_resource" "route" {
  type      = "Microsoft.AppPlatform/Spring/gateways/routeConfigs@2022-01-01-preview"
  name      = "henglu38"
  parent_id = azapi_resource.gateway.id

  body = jsonencode({
    properties = {
      appResourceId = azapi_resource.app.id
      routes = [
        {
          description = "this description"  // add some comment // add some comment
          filters     = ["StripPrefix=2", "RateLimit=1,1s"] // add some comment
          order       = 1
          predicates  = ["Path=/api5/customer/**"]
          ssoEnabled  = false
          tags        = ["tag1", "tag2"]
          title       = "myApp route config"
          tokenRelay  = true
          uri         = "testuri"
        }
      ]
    }
  })

}