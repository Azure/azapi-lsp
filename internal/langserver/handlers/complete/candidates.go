package complete

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func keyCandidates(props []azure.Property, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	propSet := make(map[string]bool)
	for _, prop := range props {
		if propSet[prop.Name] {
			continue
		}
		propSet[prop.Name] = true
		content := prop.Name
		newText := ""
		sortText := fmt.Sprintf("1%s", content)
		if prop.Modifier == azure.Required {
			sortText = fmt.Sprintf("0%s", content)
		}
		switch prop.Type {
		case "string":
			newText = fmt.Sprintf(`%s = "$0"`, content)
		case "array":
			newText = fmt.Sprintf(`%s = [$0]`, content)
		case "object":
			newText = fmt.Sprintf("%s = {\n\t$0\n}", content)
		default:
			newText = fmt.Sprintf(`%s = $0`, content)
		}
		candidates = append(candidates, lsp.CompletionItem{
			Label:  content,
			Kind:   lsp.PropertyCompletion,
			Detail: fmt.Sprintf("%s (%s)", prop.Name, prop.Modifier),
			Documentation: lsp.MarkupContent{
				Kind:  "markdown",
				Value: fmt.Sprintf("Type: `%s`  \n%s\n", prop.Type, prop.Description),
			},
			SortText:         sortText,
			InsertTextFormat: lsp.SnippetTextFormat,
			InsertTextMode:   lsp.AdjustIndentation,
			TextEdit: &lsp.TextEdit{
				Range:   r,
				NewText: newText,
			},
			Command: &lsp.Command{
				Command: "editor.action.triggerSuggest",
				Title:   "Suggest",
			},
		})
	}
	return candidates
}

func valueCandidates(values []string, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	for _, value := range values {
		content := value
		candidates = append(candidates, lsp.CompletionItem{
			Label: fmt.Sprintf(`"%s"`, content),
			Kind:  lsp.ValueCompletion,
			Documentation: lsp.MarkupContent{
				Kind:  "markdown",
				Value: fmt.Sprintf("Value: `%s`  \n", content),
			},
			SortText:         content,
			InsertTextFormat: lsp.PlainTextTextFormat,
			InsertTextMode:   lsp.AdjustIndentation,
			TextEdit: &lsp.TextEdit{
				Range:   r,
				NewText: fmt.Sprintf(`"%s"`, content),
			},
		})
	}
	return candidates
}

func typeCandidates(prefix *string, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	if prefix == nil || !strings.Contains(*prefix, "@") {
		for resourceType := range azure.GetAzureSchema().Resources {
			candidates = append(candidates, lsp.CompletionItem{
				Label: fmt.Sprintf(`"%s"`, resourceType),
				Kind:  lsp.ValueCompletion,
				Documentation: lsp.MarkupContent{
					Kind:  "markdown",
					Value: fmt.Sprintf("Type: `%s`  \n", resourceType),
				},
				SortText:         resourceType,
				InsertTextFormat: lsp.SnippetTextFormat,
				InsertTextMode:   lsp.AdjustIndentation,
				TextEdit: &lsp.TextEdit{
					Range:   r,
					NewText: fmt.Sprintf(`"%s@$0"`, resourceType),
				},
				Command: &lsp.Command{
					Command: "editor.action.triggerSuggest",
					Title:   "Suggest",
				},
			})
		}
	} else {
		resourceType := (*prefix)[0:strings.Index(*prefix, "@")]
		apiVersions := azure.GetApiVersions(resourceType)
		sort.Strings(apiVersions)
		length := len(apiVersions)
		for index, apiVersion := range apiVersions {
			candidates = append(candidates, lsp.CompletionItem{
				Label: fmt.Sprintf(`"%s@%s"`, resourceType, apiVersion),
				Kind:  lsp.ValueCompletion,
				Documentation: lsp.MarkupContent{
					Kind:  "markdown",
					Value: fmt.Sprintf("Type: `%s`  \nAPI Version: `%s`", resourceType, apiVersion),
				},
				SortText:         fmt.Sprintf("%d", length-index),
				InsertTextFormat: lsp.PlainTextTextFormat,
				InsertTextMode:   lsp.AdjustIndentation,
				TextEdit: &lsp.TextEdit{
					Range:   r,
					NewText: fmt.Sprintf(`"%s@%s"`, resourceType, apiVersion),
				},
			})
		}
	}

	return candidates
}

func bodyFuncCandidate(r lsp.Range) lsp.CompletionItem {
	return lsp.CompletionItem{
		Label: `jsonencode({})`,
		Kind:  lsp.ValueCompletion,
		Documentation: lsp.MarkupContent{
			Kind:  "markdown",
			Value: "`jsonencode` encodes a given value to a string using JSON syntax.",
		},
		SortText:         `jsonencode`,
		InsertTextFormat: lsp.SnippetTextFormat,
		InsertTextMode:   lsp.AdjustIndentation,
		TextEdit: &lsp.TextEdit{
			Range:   r,
			NewText: "jsonencode({\n\t$0\n})",
		},
		Command: &lsp.Command{
			Command: "editor.action.triggerSuggest",
			Title:   "Suggest",
		},
	}
}
