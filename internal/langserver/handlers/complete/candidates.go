package complete

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/schema"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
)

func keyCandidates(props []schema.Property, r lsp.Range) []lsp.CompletionItem {
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
		if prop.Modifier == schema.Required {
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
			Command: constTriggerSuggestCommand(),
		})
	}
	return candidates
}

func requiredPropertiesCandidates(propertySets []schema.PropertySet, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	for _, ps := range propertySets {
		if len(ps.Properties) == 0 {
			continue
		}
		props := make([]schema.Property, 0)
		for _, prop := range ps.Properties {
			props = append(props, prop)
		}
		for range props {
			for i := 0; i < len(props)-1; i++ {
				if props[i].Name > props[i+1].Name {
					props[i], props[i+1] = props[i+1], props[i]
				}
			}
		}
		newText := ""
		index := 1
		for _, prop := range props {
			if len(prop.Value) != 0 {
				newText += fmt.Sprintf("%s = \"%s\"\n", prop.Name, prop.Value)
			} else {
				switch prop.Type {
				case "string":
					newText += fmt.Sprintf(`%s = "$%d"`, prop.Name, index)
				case "array":
					newText += fmt.Sprintf(`%s = [$%d]`, prop.Name, index)
				case "object":
					newText += fmt.Sprintf("%s = {\n\t$%d\n}", prop.Name, index)
				default:
					newText += fmt.Sprintf(`%s = $%d`, prop.Name, index)
				}
				newText += "\n"
				index++
			}
		}

		label := "required-properties"
		if len(ps.Name) != 0 {
			label = fmt.Sprintf("required-properties-%s", ps.Name)
		}
		detail := "Required properties"
		if len(ps.Name) != 0 {
			detail = fmt.Sprintf("Required properties - %s", ps.Name)
		}
		candidates = append(candidates, lsp.CompletionItem{
			Label:  label,
			Kind:   lsp.SnippetCompletion,
			Detail: detail,
			Documentation: lsp.MarkupContent{
				Kind:  "markdown",
				Value: fmt.Sprintf("Type: `%s`  \n```\n%s\n```\n", ps.Name, newText),
			},
			SortText:         "0",
			InsertTextFormat: lsp.SnippetTextFormat,
			InsertTextMode:   lsp.AdjustIndentation,
			TextEdit: &lsp.TextEdit{
				Range:   r,
				NewText: newText,
			},
			Command: constTriggerSuggestCommand(),
		})
	}
	return candidates
}

func valueCandidates(values []string, r lsp.Range) []lsp.CompletionItem {
	candidates := make([]lsp.CompletionItem, 0)
	for _, value := range values {
		literal := strings.Trim(value, `"`)
		candidates = append(candidates, lsp.CompletionItem{
			Label: value,
			Kind:  lsp.ValueCompletion,
			Documentation: lsp.MarkupContent{
				Kind:  "markdown",
				Value: fmt.Sprintf("Value: `%s`  \n", literal),
			},
			SortText:         "0" + literal,
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
		Command: constTriggerSuggestCommand(),
	}
}

func constTriggerSuggestCommand() *lsp.Command {
	return &lsp.Command{
		Command: "editor.action.triggerSuggest",
		Title:   "Suggest",
	}
}
