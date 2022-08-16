package tfschema

import (
	"fmt"
	"strings"

	"github.com/Azure/azapi-lsp/internal/azure"
	"github.com/Azure/azapi-lsp/internal/azure/types"
	"github.com/Azure/azapi-lsp/internal/langserver/schema"
	ilsp "github.com/Azure/azapi-lsp/internal/lsp"
	"github.com/Azure/azapi-lsp/internal/parser"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func bodyCandidates(data []byte, filename string, block *hclsyntax.Block, attribute *hclsyntax.Attribute, pos hcl.Pos, property *Property) []lsp.CompletionItem {
	if attribute.Expr != nil {
		if _, ok := attribute.Expr.(*hclsyntax.LiteralValueExpr); ok && parser.ToLiteral(attribute.Expr) == nil {
			if property != nil {
				return property.ValueCandidatesFunc(nil, editRangeFromExprRange(attribute.Expr, pos))
			}
		}
	}

	typeValue := parser.ExtractAzureResourceType(block)
	if typeValue == nil {
		return nil
	}
	var bodyDef types.TypeBase
	def, err := azure.GetResourceDefinitionByResourceType(*typeValue)
	if err != nil || def == nil {
		return nil
	}
	bodyDef = def
	if len(block.Labels) >= 2 && block.Labels[0] == "azapi_resource_action" {
		parts := strings.Split(*typeValue, "@")
		if len(parts) != 2 {
			return nil
		}
		actionName := parser.ExtractAction(block)
		if actionName != nil && len(*actionName) != 0 {
			resourceFuncDef, err := azure.GetResourceFunction(parts[0], parts[1], *actionName)
			if err != nil || resourceFuncDef == nil {
				return nil
			}
			bodyDef = resourceFuncDef
		}
	}

	return buildCandidates(data, filename, attribute, pos, bodyDef)
}

func buildCandidates(data []byte, filename string, attribute *hclsyntax.Attribute, pos hcl.Pos, def types.TypeBase) []lsp.CompletionItem {
	candidateList := make([]lsp.CompletionItem, 0)
	hclNode := parser.JsonEncodeExpressionToHclNode(data, attribute.Expr)
	if hclNode == nil {
		return nil
	}
	hclNodes := parser.HclNodeArraysOfPos(hclNode, pos)
	if len(hclNodes) == 0 {
		return nil
	}
	lastHclNode := hclNodes[len(hclNodes)-1]

	switch {
	case parser.ContainsPos(lastHclNode.KeyRange, pos):
		// input a property with a prefix
		hclNodes := hclNodes[0 : len(hclNodes)-1]
		defs := schema.GetDef(def.AsTypeBase(), hclNodes, 0)
		keys := make([]schema.Property, 0)
		for _, def := range defs {
			keys = append(keys, schema.GetAllowedProperties(def)...)
		}
		if len(hclNodes) == 1 {
			keys = ignorePulledOutProperties(keys)
		}
		editRange := ilsp.HCLRangeToLSP(lastHclNode.KeyRange)
		candidateList = keyCandidates(keys, editRange)
	case !lastHclNode.KeyRange.Empty() && !lastHclNode.EqualRange.Empty() && lastHclNode.Children == nil:
		// input property =
		defs := schema.GetDef(def.AsTypeBase(), hclNodes, 0)
		values := make([]string, 0)
		for _, def := range defs {
			values = append(values, schema.GetAllowedValues(def)...)
		}
		editRange := lastHclNode.ValueRange
		if lastHclNode.Value == nil {
			editRange.End = pos
		}
		candidateList = valueCandidates(values, ilsp.HCLRangeToLSP(editRange), false)
	case parser.ContainsPos(lastHclNode.ValueRange, pos):
		// input a property
		defs := schema.GetDef(def.AsTypeBase(), hclNodes, 0)
		keys := make([]schema.Property, 0)
		for _, def := range defs {
			keys = append(keys, schema.GetAllowedProperties(def)...)
		}
		if len(hclNodes) == 1 {
			keys = ignorePulledOutProperties(keys)
		}
		editRange := ilsp.HCLRangeToLSP(hcl.Range{Start: pos, End: pos, Filename: filename})
		candidateList = keyCandidates(keys, editRange)

		if len(lastHclNode.Children) == 0 {
			propertySets := make([]schema.PropertySet, 0)
			for _, def := range defs {
				propertySets = append(propertySets, schema.GetRequiredPropertySet(def)...)
			}
			if len(hclNodes) == 1 {
				for i, ps := range propertySets {
					propertySets[i].Name = ""
					propertySets[i].Properties = ignorePulledOutPropertiesFromPropertySet(ps.Properties)
				}
			}
			candidateList = append(candidateList, requiredPropertiesCandidates(propertySets, editRange)...)
		}
	}
	return candidateList
}

func editRangeFromExprRange(expression hclsyntax.Expression, pos hcl.Pos) lsp.Range {
	expRange := expression.Range()
	if expRange.Start.Line != expRange.End.Line && expRange.End.Column == 1 && expRange.End.Line-1 == pos.Line {
		expRange.End = pos
	}
	return ilsp.HCLRangeToLSP(expRange)
}

func ignorePulledOutProperties(input []schema.Property) []schema.Property {
	res := make([]schema.Property, 0)
	// ignore properties pulled out from body
	for _, p := range input {
		if !isPropertyPulledOut(p) {
			res = append(res, p)
		}
	}
	return res
}

func ignorePulledOutPropertiesFromPropertySet(properties map[string]schema.Property) map[string]schema.Property {
	res := make(map[string]schema.Property)
	// ignore properties pulled out from body
	for _, p := range properties {
		if !isPropertyPulledOut(p) {
			res[p.Name] = p
		}
	}
	return res
}

func isPropertyPulledOut(p schema.Property) bool {
	return p.Name == "name" || p.Name == "location" || p.Name == "identity" || p.Name == "tags"
}

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
