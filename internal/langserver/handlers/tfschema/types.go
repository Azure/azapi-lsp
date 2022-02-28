package tfschema

import (
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var Resources []Resource

type Resource struct {
	Name       string
	Properties []Property
}

type Property struct {
	Name              string
	Modifier          string
	Type              string
	Description       string
	CompletionNewText string

	GenericCandidatesFunc GenericCandidatesFunc
	ValueCandidatesFunc   ValueCandidatesFunc
	NestedProperties      []Property
}

func (r *Resource) GetProperty(propertyName string) *Property {
	if r == nil {
		return nil
	}
	for _, prop := range r.Properties {
		if prop.Name == propertyName {
			return &prop
		}
	}
	return nil
}

type GenericCandidatesFunc func(data []byte, filename string, block *hclsyntax.Block, attribute *hclsyntax.Attribute, pos hcl.Pos, property *Property) []lsp.CompletionItem
type ValueCandidatesFunc func(prefix *string, r lsp.Range) []lsp.CompletionItem

var _ ValueCandidatesFunc = FixedValueCandidatesFunc(nil)

func FixedValueCandidatesFunc(fixedItems []lsp.CompletionItem) ValueCandidatesFunc {
	return func(prefix *string, r lsp.Range) []lsp.CompletionItem {
		for i := range fixedItems {
			if fixedItems[i].TextEdit != nil {
				fixedItems[i].TextEdit.Range = r
			}
		}
		return fixedItems
	}
}
