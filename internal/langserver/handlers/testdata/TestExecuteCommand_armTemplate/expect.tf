
variable "subscriptionId" {
  type        = string
  description = "The subscription id"
}

variable "resourceGroupName" {
  type        = string
  description = "The resource group name"
}

variable "Spring_acctest_sc_123_name" {
  type    = string
  default = "acctest-sc-123"
}

resource "azapi_resource" "Spring" {
  type      = "Microsoft.AppPlatform/Spring@2024-05-01-preview"
  parent_id = "/subscriptions/${var.subscriptionId}/resourceGroups/${var.resourceGroupName}"
  name      = "${var.Spring_acctest_sc_123_name}"
  location  = "westeurope"
  body = {
    properties = {
      networkProfile = {
        outboundType = "loadBalancer"
      }
      zoneRedundant = false
    }
    sku = {
      name = "S0"
      tier = "Standard"
    }
  }
}

resource "azapi_resource" "app" {
  type      = "Microsoft.AppPlatform/Spring/apps@2024-05-01-preview"
  parent_id = azapi_resource.Spring.id
  name      = "acctest-sca-123"
  location  = "westeurope"
  body = {
    properties = {
      addonConfigs = {
        applicationConfigurationService = {}
        configServer                    = {}
        serviceRegistry                 = {}
      }
      customPersistentDisks = []
      enableEndToEndTLS     = false
      httpsOnly             = false
      ingressSettings = {
        backendProtocol      = "Default"
        readTimeoutInSeconds = 300
        sendTimeoutInSeconds = 60
        sessionAffinity      = "None"
        sessionCookieMaxAge  = 0
      }
      persistentDisk = {
        mountPath = "/persistent"
        sizeInGB  = 0
      }
      public = false
      temporaryDisk = {
        mountPath = "/tmp"
        sizeInGB  = 5
      }
      testEndpointAuthState = "Enabled"
    }
  }
}

resource "azapi_resource" "configServer" {
  type      = "Microsoft.AppPlatform/Spring/configServers@2024-05-01-preview"
  parent_id = azapi_resource.Spring.id
  name      = "default"
  body = {
    properties = {
      configServer = {}
    }
  }
}

