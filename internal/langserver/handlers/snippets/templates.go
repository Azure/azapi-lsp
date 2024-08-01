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

func TemplateCandidates(editRange lsp.Range) []lsp.CompletionItem {
	templates := make([]CompletionModel, 0)
	data, err := templateJSON.ReadFile("templates.json")
	if err != nil {
		return nil
	}
	err = json.Unmarshal(data, &templates)
	if err != nil {
		return nil
	}

	candidates := make([]lsp.CompletionItem, 0)
	for _, template := range templates {
		candidates = append(candidates, lsp.CompletionItem{
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
	return candidates
}
