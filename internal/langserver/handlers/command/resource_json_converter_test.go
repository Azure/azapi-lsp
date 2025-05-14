package command

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func Test_ParseResourceJson(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple resource json",
			input: `
{
    "properties": {
        "customerId": "547614eb-a6da-43a7-9dd0-f004940a6865",
        "provisioningState": "Succeeded",
        "sku": {
            "name": "PerGB2018",
            "lastSkuUpdate": "2024-07-19T05:47:34.5001244Z"
        },
        "retentionInDays": 30,
        "features": {
            "legacy": 0,
            "searchVersion": 1,
            "enableLogAccessUsingOnlyResourcePermissions": true,
            "disableLocalAuth": false
        },
        "workspaceCapping": {
            "dailyQuotaGb": -1,
            "quotaNextResetTime": "2024-07-19T22:00:00Z",
            "dataIngestionStatus": "RespectQuota"
        },
        "publicNetworkAccessForIngestion": "Enabled",
        "publicNetworkAccessForQuery": "Enabled",
        "createdDate": "2024-07-19T05:47:34.5001244Z",
        "modifiedDate": "2024-07-19T05:47:58.6330522Z"
    },
    "location": "eastus",
    "tags": {},
    "id": "/subscriptions/000000/resourceGroups/acctestRG-240719134659555592/providers/Microsoft.OperationalInsights/workspaces/acctestLAW-240719134659555592",
    "name": "acctestLAW-240719134659555592",
    "type": "Microsoft.OperationalInsights/workspaces",
    "etag": "\"1b004b62-0000-0100-0000-6699fe0e0000\""
}`,
			expected: `resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2025-02-01"
  parent_id = "/subscriptions/000000/resourceGroups/acctestRG-240719134659555592"
  name      = "acctestLAW-240719134659555592"
  location  = "eastus"
  body = {
    etag = "\"1b004b62-0000-0100-0000-6699fe0e0000\""
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
        legacy                                      = 0
        searchVersion                               = 1
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
`,
		},
		{
			name: "resource from ARM template",
			input: `{
            "type": "Microsoft.OperationalInsights/workspaces",
            "apiVersion": "2023-09-01",
            "name": "${var.foo}",
            "location": "eastus",
            "properties": {
                "sku": {
                    "name": "PerGB2018"
                },
                "retentionInDays": 30,
                "features": {
                    "legacy": 0,
                    "searchVersion": 1,
                    "enableLogAccessUsingOnlyResourcePermissions": true,
                    "disableLocalAuth": false
                },
                "workspaceCapping": {
                    "dailyQuotaGb": -1
                },
                "publicNetworkAccessForIngestion": "Enabled",
                "publicNetworkAccessForQuery": "Enabled"
            }
        }`,
			expected: `resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2023-09-01"
  parent_id = "/subscriptions/$${var.subscriptionId}/resourceGroups/$${var.resourceGroupName}"
  name      = "$${var.foo}"
  location  = "eastus"
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
        legacy                                      = 0
        searchVersion                               = 1
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
`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			block, err := ParseResourceJson(tc.input)
			if err != nil {
				t.Errorf("ParseResourceJson() error = %v", err)
				return
			}

			result := string(hclwrite.Format(block.BuildTokens(nil).Bytes()))
			if result != tc.expected {
				t.Errorf("ParseResourceJson() = %v, want %v", result, tc.expected)
			}
		})
	}
}
