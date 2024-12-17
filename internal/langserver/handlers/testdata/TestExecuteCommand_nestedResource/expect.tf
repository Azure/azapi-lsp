/*
Note: This is a generated HCL content from the JSON input which is based on the latest API version available.
To import the resource, please run the following command:
terraform import azapi_resource.subnet /subscriptions/0000000/resourceGroups/acctestRG-databricks-240530173443928447/providers/Microsoft.Network/virtualNetworks/acctest-vnet-240530173443928447/subnets/acctest-sn-private-240530173443928447?api-version=2024-05-01

Or add the below config:
import {
  id = "/subscriptions/0000000/resourceGroups/acctestRG-databricks-240530173443928447/providers/Microsoft.Network/virtualNetworks/acctest-vnet-240530173443928447/subnets/acctest-sn-private-240530173443928447?api-version=2024-05-01"
  to = azapi_resource.subnet
}
*/

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = "/subscriptions/0000000/resourceGroups/acctestRG-databricks-240530173443928447/providers/Microsoft.Network/virtualNetworks/acctest-vnet-240530173443928447"
  name      = "acctest-sn-private-240530173443928447"
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
      delegations = [{
        id   = "/subscriptions/0000000/resourceGroups/acctestRG-databricks-240530173443928447/providers/Microsoft.Network/virtualNetworks/acctest-vnet-240530173443928447/subnets/acctest-sn-private-240530173443928447/delegations/acctest"
        name = "acctest"
        properties = {
          serviceName = "Microsoft.Databricks/workspaces"
        }
        type = "Microsoft.Network/virtualNetworks/subnets/delegations"
      }]
      networkSecurityGroup = {
        id = "/subscriptions/0000000/resourceGroups/acctestRG-databricks-240530173443928447/providers/Microsoft.Network/networkSecurityGroups/acctest-nsg-240530173443928447"
      }
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpoints                  = []
    }
  }
}
