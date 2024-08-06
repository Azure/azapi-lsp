package command

import (
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
)

func flattenARMExpression(input interface{}) interface{} {
	if input == nil {
		return nil
	}
	switch v := input.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{})
		for key, value := range v {
			res[key] = flattenARMExpression(value)
		}
		return res
	case []interface{}:
		res := make([]interface{}, 0)
		for _, value := range v {
			res = append(res, flattenARMExpression(value))
		}
		return res
	case string:
		if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
			if output, err := evaluateARMTemplateExpression(v[1 : len(v)-1]); err == nil {
				return output
			}
		}
		return v
	default:
		return v
	}
}

func evaluateARMTemplateExpression(input string) (string, error) {
	env := map[string]interface{}{
		"resourceId": func(input ...string) string {
			resourceType := ""
			subscriptionId := "var.subscriptionId"
			resourceGroupName := "var.resourceGroupName"
			index := -1
			for i, v := range input {
				if strings.Contains(v, "/") {
					resourceType = v
					index = i
					break
				}
			}

			switch index {
			case -1:
				return "resourceId"
			case 1:
				resourceGroupName = input[0]
			case 2:
				subscriptionId = input[0]
				resourceGroupName = input[1]
			}
			scopeId := fmt.Sprintf("/subscriptions/${%s}/resourceGroups/${%s}", subscriptionId, resourceGroupName)
			parts := strings.Split(resourceType, "/")
			resourceId := fmt.Sprintf("%s/providers/%s", scopeId, parts[0])
			for i := 1; i < len(parts); i++ {
				name := ""
				if inputIndex := index + i; inputIndex < len(input) {
					name = input[inputIndex]
				}
				resourceId = fmt.Sprintf("%s/%s/%s", resourceId, parts[i], name)
			}
			return resourceId
		},
		"parameters": func(name string) string {
			return fmt.Sprintf("${var.%s}", name)
		},
		"concat": func(input ...string) string {
			return strings.Join(input, "")
		},
	}

	program, err := expr.Compile(input, expr.Env(env))
	if err != nil {
		return "", err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return "", err
	}

	if output == nil {
		return "", fmt.Errorf("output is nil")
	}

	v, ok := output.(string)
	if !ok {
		return "", fmt.Errorf("output is not string")
	}
	return v, nil
}
