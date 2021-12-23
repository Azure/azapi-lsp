package complete

import (
	"fmt"
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
		switch prop.Type {
		case "string":
			newText = fmt.Sprintf(`%s = "$0"`, content)
			break
		case "array":
			newText = fmt.Sprintf(`%s = [$0]`, content)
			break
		case "object":
			newText = fmt.Sprintf("%s = {\n\t$0\n}", content)
			break
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
			SortText:         content,
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
			resourceType = strings.Title(strings.ToLower(resourceType))
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
		if resource, ok := azure.GetAzureSchema().Resources[strings.ToUpper(resourceType)]; ok {
			for _, resourceDef := range resource.Definitions {
				apiVersion := resourceDef.ApiVersion
				candidates = append(candidates, lsp.CompletionItem{
					Label: fmt.Sprintf(`"%s@%s"`, resourceType, apiVersion),
					Kind:  lsp.ValueCompletion,
					Documentation: lsp.MarkupContent{
						Kind:  "markdown",
						Value: fmt.Sprintf("Type: `%s`  \nAPI Version: `%s`", resourceType, apiVersion),
					},
					SortText:         apiVersion,
					InsertTextFormat: lsp.PlainTextTextFormat,
					InsertTextMode:   lsp.AdjustIndentation,
					TextEdit: &lsp.TextEdit{
						Range:   r,
						NewText: fmt.Sprintf(`"%s@%s"`, resourceType, apiVersion),
					},
				})
			}
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
