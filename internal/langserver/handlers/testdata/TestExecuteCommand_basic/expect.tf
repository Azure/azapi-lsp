/*
Note: This is a generated HCL content from the JSON input which is based on the latest API version available.
To import the resource, please run the following command:
terraform import azapi_resource.managedCluster /subscriptions/0000000/resourcegroups/acctestRG-aks-240626160102776968/providers/Microsoft.ContainerService/managedClusters/acctestaks240626160102776968?api-version=2025-02-02-preview

Or add the below config:
import {
  id = "/subscriptions/0000000/resourcegroups/acctestRG-aks-240626160102776968/providers/Microsoft.ContainerService/managedClusters/acctestaks240626160102776968?api-version=2025-02-02-preview"
  to = azapi_resource.managedCluster
}
*/

resource "azapi_resource" "managedCluster" {
  type      = "Microsoft.ContainerService/managedClusters@2025-02-02-preview"
  parent_id = "/subscriptions/0000000/resourceGroups/acctestRG-aks-240626160102776968"
  name      = "acctestaks240626160102776968"
  location  = "eastus"
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      agentPoolProfiles = [{
        count                  = 3
        enableAutoScaling      = false
        enableEncryptionAtHost = false
        enableFIPS             = false
        enableNodePublicIP     = false
        enableUltraSSD         = false
        kubeletDiskType        = "OS"
        maxPods                = 30
        mode                   = "System"
        name                   = "np"
        orchestratorVersion    = "1.28"
        osDiskSizeGB           = 128
        osDiskType             = "Managed"
        osSKU                  = "Ubuntu"
        osType                 = "Linux"
        powerState = {
          code = "Running"
        }
        scaleDownMode = "Delete"
        type          = "VirtualMachineScaleSets"
        upgradeSettings = {
          maxSurge = "10%"
        }
        vmSize = "Standard_DS2_v2"
      }]
      autoUpgradeProfile = {
        nodeOSUpgradeChannel = "NodeImage"
        upgradeChannel       = "none"
      }
      azureMonitorProfile = {
        metrics = {
          enabled          = false
          kubeStateMetrics = {}
        }
      }
      disableLocalAccounts = false
      dnsPrefix            = "acctestaks240626160102776968"
      enableRBAC           = true
      kubernetesVersion    = "1.28"
      linuxProfile = {
        adminUsername = "acctestuser240626160102776968"
        ssh = {
          publicKeys = [{
            keyData = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
          }]
        }
      }
      networkProfile = {
        dnsServiceIP = "10.10.0.10"
        ipFamilies   = ["IPv4"]
        loadBalancerProfile = {
          backendPoolType = "nodeIPConfiguration"
          managedOutboundIPs = {
            count = 1
          }
        }
        loadBalancerSku  = "standard"
        networkDataplane = "azure"
        networkPlugin    = "azure"
        networkPolicy    = "azure"
        outboundType     = "loadBalancer"
        serviceCidr      = "10.10.0.0/16"
        serviceCidrs     = ["10.10.0.0/16"]
      }
      nodeResourceGroup = "MC_acctestRG-aks-240626160102776968_acctestaks240626160102776968_eastus"
      oidcIssuerProfile = {
        enabled = false
      }
      securityProfile = {
        imageCleaner = {
          enabled       = false
          intervalHours = 48
        }
      }
      servicePrincipalProfile = {
        clientId = "msi"
      }
      storageProfile = {
        diskCSIDriver = {
          enabled = true
        }
        fileCSIDriver = {
          enabled = true
        }
        snapshotController = {
          enabled = true
        }
      }
      supportPlan = "KubernetesOfficial"
      windowsProfile = {
        adminUsername  = "azureuser"
        enableCSIProxy = true
        gmsaProfile = {
          dnsServer      = ""
          enabled        = true
          rootDomainName = ""
        }
        licenseType = "None"
      }
      workloadAutoScalerProfile = {}
    }
    sku = {
      name = "Base"
      tier = "Free"
    }
  }
}
