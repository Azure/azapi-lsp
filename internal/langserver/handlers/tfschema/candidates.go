package tfschema

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func PropertiesCandidates(props []Property, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	for index, prop := range props {
		candidates = append(candidates, lsp.CompletionItem{
			Label:  prop.Name,
			Kind:   lsp.PropertyCompletion,
			Detail: fmt.Sprintf("%s (%s)", prop.Name, prop.Modifier),
			Documentation: lsp.MarkupContent{
				Kind:  "markdown",
				Value: fmt.Sprintf("Type: `%s`  \n%s\n", prop.Type, prop.Description),
			},
			SortText:         fmt.Sprintf("%04d", index),
			InsertTextFormat: lsp.SnippetTextFormat,
			InsertTextMode:   lsp.AdjustIndentation,
			TextEdit: &lsp.TextEdit{
				Range:   r,
				NewText: prop.CompletionNewText,
			},
			Command: constTriggerSuggestCommand(),
		})
	}
	return candidates
}

func valueCandidates(values []string, r lsp.Range, isOrdered bool) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	for index, value := range values {
		literal := strings.Trim(value, `"`)
		sortText := "0" + literal
		if isOrdered {
			sortText = fmt.Sprintf("%04d", index)
		}
		candidates = append(candidates, lsp.CompletionItem{
			Label: value,
			Kind:  lsp.ValueCompletion,
			Documentation: lsp.MarkupContent{
				Kind:  "markdown",
				Value: fmt.Sprintf("Value: `%s`  \n", literal),
			},
			SortText:         sortText,
			InsertTextFormat: lsp.PlainTextTextFormat,
			InsertTextMode:   lsp.AdjustIndentation,
			TextEdit: &lsp.TextEdit{
				Range:   r,
				NewText: value,
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
				Command: constTriggerSuggestCommand(),
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

func locationCandidates(_ *string, r lsp.Range) []lsp.CompletionItem {
	values := make([]string, 0)
	for _, location := range supportedLocations() {
		values = append(values, fmt.Sprintf(`"%s"`, location))
	}
	return valueCandidates(values, r, true)
}

func identityTypesCandidates(_ *string, r lsp.Range) []lsp.CompletionItem {
	values := []string{
		`"SystemAssigned"`,
		`"UserAssigned"`,
		`"SystemAssigned, UserAssigned"`,
	}
	return valueCandidates(values, r, false)
}

func boolCandidates(_ *string, r lsp.Range) []lsp.CompletionItem {
	return valueCandidates([]string{"true", "false"}, r, false)
}

func bodyJsonencodeFuncCandidate() lsp.CompletionItem {
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
			NewText: "jsonencode({\n\t$0\n})",
		},
		Command: constTriggerSuggestCommand(),
	}
}

func constTriggerSuggestCommand() *lsp.Command {
	return &lsp.Command{
		Command: "editor.action.triggerSuggest",
		Title:   "Suggest",
	}
}
