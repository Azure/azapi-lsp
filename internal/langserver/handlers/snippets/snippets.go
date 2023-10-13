package snippets

import (
	"fmt"
	"log"
	"strings"

	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Snippet struct {
	AzureResourceType string
	Fields            []Field
}

type Field struct {
	Name  string
	Value string
}

func (field Field) Order() int {
	switch field.Name {
	case "type":
		return -1
	case "parent_id":
		return 0
	case "name":
		return 1
	case "resource_id":
		return 2
	case "location":
		return 3
	case "identity":
		return 4
	case "body":
		return 5
	case "tags":
		return 6
	case "schema_validation_enabled":
		return 7
	case "response_export_values":
		return 8
	default:
		return 9
	}
}

func CodeSampleCandidates(block *hclsyntax.Block, editRange lsp.Range) []lsp.CompletionItem {
	if block == nil {
		return nil
	}

	typeValue := ""
	for _, attr := range block.Body.Attributes {
		if attr.Name == "type" {
			v, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				log.Printf("[ERROR] %s", diags)
				return nil
			}
			typeValue = v.AsString()
			break
		}
	}

	azureResourceType := parseResourceType(typeValue)
	if snippet, ok := snippetMap[strings.ToLower(azureResourceType)]; ok {
		newText := ""
		for _, field := range snippet.Fields {
			if _, ok := block.Body.Attributes[field.Name]; ok {
				continue
			}
			newText += field.Value + "\n"
		}
		if newText == "" {
			return nil
		}
		return []lsp.CompletionItem{
			{
				Label:  "code sample",
				Kind:   lsp.SnippetCompletion,
				Detail: "Code Sample",
				Documentation: lsp.MarkupContent{
					Kind:  "markdown",
					Value: fmt.Sprintf("```\n%s\n```\n", newText),
				},
				SortText:         "0",
				InsertTextFormat: lsp.SnippetTextFormat,
				InsertTextMode:   lsp.AdjustIndentation,
				TextEdit: &lsp.TextEdit{
					Range:   editRange,
					NewText: newText,
				},
			},
		}
	}

	return nil
}

func parseResourceType(typeValue string) string {
	parts := strings.Split(typeValue, "@")
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}
