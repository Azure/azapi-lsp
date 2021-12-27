package handlers

import (
	"context"
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
	"github.com/zclconf/go-cty/cty"
	"strings"
)

func (svc *service) TextDocumentComplete(ctx context.Context, params lsp.CompletionParams) (lsp.CompletionList, error) {
	var list lsp.CompletionList

	fs, err := lsctx.DocumentStorage(ctx)
	if err != nil {
		return list, err
	}

	cc, err := ilsp.ClientCapabilities(ctx)
	if err != nil {
		return list, err
	}

	doc, err := fs.GetDocument(ilsp.FileHandlerFromDocumentURI(params.TextDocument.URI))
	if err != nil {
		return list, err
	}

	d, err := svc.decoderForDocument(ctx, doc)
	if err != nil {
		return list, err
	}

	expFeatures, err := lsctx.ExperimentalFeatures(ctx)
	if err != nil {
		return list, err
	}

	d.PrefillRequiredFields = expFeatures.PrefillRequiredFields

	fPos, err := ilsp.FilePositionFromDocumentPosition(params.TextDocumentPositionParams, doc)

	if err != nil {
		return list, err
	}

	svc.logger.Printf("Looking for candidates at %q -> %#v", doc.Filename(), fPos.Position())

	_, _ = d.CandidatesAtPos(doc.Filename(), fPos.Position())

	data, err := doc.Text()
	if err != nil {
		return list, err
	}
	file, _ := hclsyntax.ParseConfig(data, doc.Filename(), hcl.InitialPos)
	block, err := blockAtPos(file, fPos.Position())

	candidateList := make([]lsp.CompletionItem, 0)
	if err == nil && block != nil && len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
		if attribute := attributeAtPos(block, fPos.Position()); attribute != nil {
			switch attribute.Name {
			case "type":
				expRange := attribute.Expr.Range()
				if expRange.Start.Line != expRange.End.Line && expRange.End.Column == 1 && expRange.End.Line-1 == fPos.Position().Line {
					expRange.End = fPos.Position()
				}
				r := ilsp.HCLRangeToLSP(expRange)
				prefix := toLiteral(attribute.Expr)
				candidateList = typeCandidates(prefix, r)
				break
			case "body":
				break
			}
		}
	}

	candidates := lsp.CompletionList{
		IsIncomplete: false,
		Items:        candidateList,
	}
	svc.logger.Printf("received candidates: %#v", candidates)
	_ = cc
	return candidates, err
}

func typeCandidates(prefix *string, r lsp.Range) []lsp.CompletionItem {
	types := map[string][]string{
		"Microsoft.MachineLearningServices/workspaces": {
			"2021-05-01",
			"2021-07-01",
		},
		"Microsoft.MachineLearningServices/workspaces/computes": {
			"2021-05-01",
			"2021-07-01",
		},
	}

	candidates := make([]lsp.CompletionItem, 0)
	if prefix == nil || !strings.Contains(*prefix, "@") {
		for resourceType := range types {
			candidates = append(candidates, lsp.CompletionItem{
				Label: fmt.Sprintf(`"%s"`, resourceType),
				Kind:  lsp.ValueCompletion,
				Documentation: lsp.MarkupContent{
					Kind:  lsp.MarkupKind("markdown"),
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
		for _, apiVersion := range types[resourceType] {
			candidates = append(candidates, lsp.CompletionItem{
				Label: fmt.Sprintf(`"%s@%s"`, resourceType, apiVersion),
				Kind:  lsp.ValueCompletion,
				Documentation: lsp.MarkupContent{
					Kind:  lsp.MarkupKind("markdown"),
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

	return candidates
}

func toLiteral(expression hclsyntax.Expression) *string {
	value, dialog := expression.Value(&hcl.EvalContext{})
	if dialog != nil && dialog.HasErrors() {
		return nil
	}
	if value.Type() == cty.String && !value.IsNull() && value.IsKnown() {
		v := value.AsString()
		return &v
	}
	return nil
}

func blockAtPos(file *hcl.File, pos hcl.Pos) (*hclsyntax.Block, error) {
	body, isHcl := file.Body.(*hclsyntax.Body)
	if !isHcl {
		return nil, fmt.Errorf("file is not hcl")
	}

	for _, b := range body.Blocks {
		if ContainsPos(b.Range(), pos) {
			return b, nil
		}
	}
	return nil, nil
}

func attributeAtPos(block *hclsyntax.Block, pos hcl.Pos) *hclsyntax.Attribute {
	if block == nil {
		return nil
	}

	for _, attr := range block.Body.Attributes {
		if ContainsPos(attr.SrcRange, pos) {
			return attr
		}
		if ContainsPos(attr.Expr.Range(), pos) {
			return attr
		}
	}

	return nil
}

func ContainsPos(r hcl.Range, pos hcl.Pos) bool {
	afterStart := pos.Line > r.Start.Line || pos.Line == r.Start.Line && pos.Column >= r.Start.Column
	beforeEnd := pos.Line < r.End.Line || pos.Line == r.End.Line && pos.Column <= r.End.Column
	return afterStart && beforeEnd
}
