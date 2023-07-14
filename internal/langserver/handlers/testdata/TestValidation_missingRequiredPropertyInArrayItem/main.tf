
resource "azapi_resource" "test" {
  type = "Microsoft.ContainerService/managedClusters@2023-05-02-preview"
  body = jsonencode({
    extendedLocation = {
    
    }
    properties = {
      agentPoolProfiles = [
        {
          mode = "System"
        }, 
        {

        }
      ]
    }
  })
}