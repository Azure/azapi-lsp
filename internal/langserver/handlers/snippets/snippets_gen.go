package snippets

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"golang.org/x/exp/slices"
)

func parseSnippet(filepath string) (*Snippet, error) {
	// #nosec G304
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	file, diags := hclsyntax.ParseConfig(data, filepath, hcl.InitialPos)
	if diags.HasErrors() {
		return nil, diags
	}

	blocks := file.Body.(*hclsyntax.Body).Blocks
	if len(blocks) == 0 {
		return nil, fmt.Errorf("no block found in %s", filepath)
	}

	addrTypeMap := make(map[string]string)
	for _, block := range blocks {
		if block.Type == "resource" || block.Type == "data" {
			typeValue := ""
			for _, attr := range block.Body.Attributes {
				if attr.Name == "type" {
					typeValue = strings.Trim(stringValue(data, attr.Expr.Range()), `"`)
					break
				}
			}
			addr := strings.Join(block.Labels, ".")
			if block.Type == "data" {
				addr = "data." + addr
			}
			addrTypeMap[addr] = typeValue
		}
	}

	lastBlock := blocks[len(blocks)-1]
	if lastBlock.Labels[0] != "azapi_resource" {
		// skip non azapi_resource block
		return nil, nil
	}
	typeValue := ""
	fields := make([]Field, 0)
	index := 1
	for _, attr := range lastBlock.Body.Attributes {
		if attr.Name == "type" {
			typeValue = strings.Trim(stringValue(data, attr.Expr.Range()), `"`)
		}
		value := stringValue(data, attr.Range())
		vars := attr.Expr.Variables()
		for _, variable := range vars {
			rawContent := stringValue(data, variable.SourceRange())
			placeholder := placeholderContent(rawContent, addrTypeMap)
			format := `${%d:"%s"}`
			if attr.Name == "name" || attr.Name == "location" {
				format = `"${%d:%s}"`
			}
			value = strings.ReplaceAll(value, rawContent, fmt.Sprintf(format, index, placeholder))
			index++
		}
		fields = append(fields, Field{
			Name:  attr.Name,
			Value: value,
		})
	}
	for _, block := range lastBlock.Body.Blocks {
		value := stringValue(data, block.Range())
		for _, attr := range block.Body.Attributes {
			vars := attr.Expr.Variables()
			for _, variable := range vars {
				rawContent := stringValue(data, variable.SourceRange())
				placeholder := placeholderContent(rawContent, addrTypeMap)
				value = strings.ReplaceAll(value, rawContent, fmt.Sprintf("${%d:%s}", index, placeholder))
				index++
			}
		}
		fields = append(fields, Field{
			Name:  block.Type,
			Value: value,
		})
	}

	slices.SortFunc(fields, func(i, j Field) int {
		return i.Order() - j.Order()
	})

	return &Snippet{
		AzureResourceType: parseResourceType(typeValue),
		Fields:            fields,
	}, nil
}

func placeholderContent(content string, typeMap map[string]string) string {
	out := "TODO"
	parts := strings.Split(content, ".")
	if len(parts) == 0 {
		return out
	}
	lastPart := parts[len(parts)-1]
	if lastPart == "subscription_id" {
		return "subscription id"
	}
	if lastPart == "tenant_id" {
		return "tenant id"
	}
	if lastPart == "output" {
		return out
	}
	switch parts[0] {
	case "var":
		if lastPart == "resource_name" {
			return "The name of the resource"
		}
		return fmt.Sprintf("%s", lastPart)
	case "local":
		return fmt.Sprintf("%s", lastPart)
	case "data":
		addr := strings.Join(parts[:len(parts)-1], ".")
		if typeValue, ok := typeMap[addr]; ok {
			return fmt.Sprintf("The %s of the %s resource", lastPart, typeValue)
		}
	case "azapi_resource", "azapi_resource_action", "azapi_update_resource":
		addr := strings.Join(parts[:len(parts)-1], ".")
		if len(parts) == 2 {
			addr = content
		}
		if typeValue, ok := typeMap[addr]; ok {
			return fmt.Sprintf("The %s of the %s resource", lastPart, typeValue)
		}
	}
	return out
}

func stringValue(data []byte, h hcl.Range) string {
	return string(data[h.Start.Byte:h.End.Byte])
}
