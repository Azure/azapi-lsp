package command

import (
	"testing"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func Test_Context(t *testing.T) {
	res1 := hclwrite.NewBlock("resource", []string{"azapi_resource", "managedCluster"})
	res1.Body().SetAttributeValue("type", cty.StringVal("Microsoft.ContainerService/managedClusters@2021-03-01"))
	res1.Body().SetAttributeValue("name", cty.StringVal("myCluster"))
	res1.Body().SetAttributeValue("parent_id", cty.StringVal("myParentId"))

	res2 := hclwrite.NewBlock("resource", []string{"azapi_resource", "managedCluster"})
	res2.Body().SetAttributeValue("type", cty.StringVal("Microsoft.ContainerService/managedClusters@2021-03-01"))
	res2.Body().SetAttributeValue("name", cty.StringVal("myCluster"))
	res2.Body().SetAttributeValue("parent_id", cty.StringVal("myParentId"))

	var1 := hclwrite.NewBlock("variable", []string{"foo"})

	tests := []struct {
		name   string
		input  []*hclwrite.Block
		expect string
	}{
		{
			name:  "empty",
			input: []*hclwrite.Block{},
			expect: `
variable "subscriptionId" {
  type        = string
  description = "The subscription id"
}

variable "resourceGroupName" {
  type        = string
  description = "The resource group name"
}

`,
		},
		{
			name: "with variables",
			input: []*hclwrite.Block{
				var1,
			},
			expect: `
variable "subscriptionId" {
  type        = string
  description = "The subscription id"
}

variable "resourceGroupName" {
  type        = string
  description = "The resource group name"
}

variable "foo" {
}

`,
		},
		{
			name: "with resources",
			input: []*hclwrite.Block{
				res1,
				res2,
			},
			expect: `
variable "subscriptionId" {
  type        = string
  description = "The subscription id"
}

variable "resourceGroupName" {
  type        = string
  description = "The resource group name"
}

resource "azapi_resource" "managedCluster" {
  type      = "Microsoft.ContainerService/managedClusters@2021-03-01"
  name      = "myCluster"
  parent_id = "myParentId"
}

resource "azapi_resource" "managedCluster1" {
  type      = "Microsoft.ContainerService/managedClusters@2021-03-01"
  name      = "myCluster"
  parent_id = "myParentId"
}

`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewContext()
			for _, b := range tt.input {
				ctx.AppendBlock(b)
			}
			if got := ctx.String(); got != tt.expect {
				t.Errorf("unexpected result, got: %s, expect: %s", got, tt.expect)
			}
		})
	}
}
