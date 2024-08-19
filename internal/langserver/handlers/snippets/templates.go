package snippets

import (
	"embed"
	"encoding/json"

	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

//go:embed templates.json
var templateJSON embed.FS

type CompletionModel struct {
	Label         string             `json:"label"`
	Documentation DocumentationModel `json:"documentation"`
	SortText      string             `json:"sortText"`
	TextEdit      TextEditModel      `json:"textEdit"`
}

type TextEditModel struct {
	NewText string `json:"newText"`
}

type DocumentationModel struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

var templateCandidates []lsp.CompletionItem

func TemplateCandidates(editRange lsp.Range) []lsp.CompletionItem {
	if len(templateCandidates) != 0 {
		for i := range templateCandidates {
			templateCandidates[i].TextEdit.Range = editRange
		}
		return templateCandidates
	}
	templates := make([]CompletionModel, 0)
	data, err := templateJSON.ReadFile("templates.json")
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data, &templates)
	if err != nil {
		return nil
	}

	for _, template := range templates {
		templateCandidates = append(templateCandidates, lsp.CompletionItem{
			Label:  template.Label,
			Kind:   lsp.SnippetCompletion,
			Detail: "Code Sample",
			Documentation: lsp.MarkupContent{
				Kind:  "markdown",
				Value: template.Documentation.Value,
			},
			SortText:         template.SortText,
			InsertTextFormat: lsp.SnippetTextFormat,
			InsertTextMode:   lsp.AdjustIndentation,
			TextEdit: &lsp.TextEdit{
				Range:   editRange,
				NewText: template.TextEdit.NewText,
			},
		})
	}
	return templateCandidates
}
