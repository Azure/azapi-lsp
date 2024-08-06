package command

import (
	"reflect"
	"testing"
)

func Test_flattenARMExpression(t *testing.T) {
	testcases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "map",
			input: map[string]interface{}{
				"foo": "bar",
				"baz": "[parameters('myParam')]",
			},
			expected: map[string]interface{}{
				"foo": "bar",
				"baz": "${var.myParam}",
			},
		},
		{
			name: "array",
			input: []interface{}{
				"foo",
				"[parameters('myParam')]",
			},
			expected: []interface{}{
				"foo",
				"${var.myParam}",
			},
		},
		{
			name:     "string",
			input:    "[parameters('myParam')]",
			expected: "${var.myParam}",
		},
		{
			name:     "string without brackets",
			input:    "parameters('myParam')",
			expected: "parameters('myParam')",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			actual := flattenARMExpression(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("flattenARMExpression(%v) => %v, want %v", tt.input, actual, tt.expected)
			}
		})
	}
}

func Test_evaluateARMTemplateExpression(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic resourceId",
			input:    "resourceId('Microsoft.Network/virtualNetworks', 'myVnet')",
			expected: "/subscriptions/${var.subscriptionId}/resourceGroups/${var.resourceGroupName}/providers/Microsoft.Network/virtualNetworks/myVnet",
		},
		{
			name:     "resourceId with resource group name",
			input:    "resourceId('myResourceGroup', 'Microsoft.Network/virtualNetworks', 'myVnet')",
			expected: "/subscriptions/${var.subscriptionId}/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet",
		},
		{
			name:     "resourceId with subscription id and resource group name",
			input:    "resourceId('mySubscription', 'myResourceGroup', 'Microsoft.Network/virtualNetworks', 'myVnet')",
			expected: "/subscriptions/mySubscription/resourceGroups/myResourceGroup/providers/Microsoft.Network/virtualNetworks/myVnet",
		},
		{
			name:     "resourceId with multiple segments",
			input:    "resourceId('Microsoft.Network/virtualNetworks/subnets', 'myVnet', 'mySubnet')",
			expected: "/subscriptions/${var.subscriptionId}/resourceGroups/${var.resourceGroupName}/providers/Microsoft.Network/virtualNetworks/subnets/myVnet/mySubnet",
		},
		{
			name:     "parameters",
			input:    "parameters('myParam')",
			expected: "${var.myParam}",
		},
		{
			name:     "concat",
			input:    "concat('foo', 'bar')",
			expected: "foobar",
		},
		{
			name:     "complex expression",
			input:    `concat('foo', resourceId('Microsoft.Network/publicIPAddresses', parameters('myParam')), 'bar')`,
			expected: "foo/subscriptions/${var.subscriptionId}/resourceGroups/${var.resourceGroupName}/providers/Microsoft.Network/publicIPAddresses/${var.myParam}bar",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := evaluateARMTemplateExpression(tt.input)
			if err != nil {
				t.Errorf("evaluateARMTemplateExpression(%v) => %v, want %v", tt.input, actual, tt.expected)
			}
		})
	}
}
