package handlers

import (
	"context"
	"fmt"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	lsctx "github.com/ms-henglu/azurerm-restapi-lsp/internal/context"
	ilsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/lsp"
	lsp "github.com/ms-henglu/azurerm-restapi-lsp/internal/protocol"
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

	candidateList := make([]lang.Candidate, 0)
	if err == nil && block != nil {
		if len(block.Labels) != 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
			attribute, err := attributeAtPos(block, fPos.Position())
			if err == nil && attribute != nil {
				switch attribute.Name {
				case "type":
					candidateList = typeCandidates()
					for i, _ := range candidateList {
						candidateList[i].TextEdit.Range = hcl.Range{
							Filename: doc.Filename(),
							End: fPos.Position(),
							Start: attribute.Expr.Range().Start,
						}
					}
					break
				case "body":
					break
				}
			}
		}
	}


	candidates := lang.Candidates{
		List:candidateList,
		IsComplete: true,
	}
	svc.logger.Printf("received candidates: %#v", candidates)
	return ilsp.ToCompletionList(candidates, cc.TextDocument), err
}


func typeCandidates() []lang.Candidate {
	types := []string{
		"Microsoft.MachineLearningServices/workspaces@2021-07-01",
		"Microsoft.MachineLearningServices/workspaces/computes@2021-07-1",
	}

	candidates := make([]lang.Candidate, 0)
	for _, resourceType := range types {
		candidates = append(candidates, lang.Candidate{
			Label:        resourceType,
			Description:  lang.PlainText("this is description for " + resourceType),
			Detail:       "string",
			IsDeprecated: false,
			TextEdit: lang.TextEdit{
				NewText: fmt.Sprintf(`"%s"`, resourceType),
				Snippet: fmt.Sprintf(`"%s"`, resourceType),
			},
			Kind:           lang.AttributeCandidateKind,
			TriggerSuggest: true,
		})
	}
	return candidates
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

func attributeAtPos(block *hclsyntax.Block, pos hcl.Pos) (*hclsyntax.Attribute, error) {
	if block == nil {
		return nil, nil
	}

	for _, attr := range block.Body.Attributes {
		if ContainsPos(attr.SrcRange, pos) {
			return attr, nil
		}
	}

	return nil, nil
}

func ContainsPos(r hcl.Range, pos hcl.Pos) bool {
	afterStart := pos.Line > r.Start.Line || pos.Line == r.Start.Line && pos.Column >= r.Start.Column
	beforeEnd := pos.Line < r.End.Line || pos.Line == r.End.Line && pos.Column <= r.End.Column
	return afterStart && beforeEnd
}