package parser

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func ToLiteral(expression hclsyntax.Expression) *string {
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

func ToLiteralBoolean(expression hclsyntax.Expression) *bool {
	value, dialog := expression.Value(&hcl.EvalContext{})
	if dialog != nil && dialog.HasErrors() {
		return nil
	}
	if value.Type() == cty.Bool && !value.IsNull() && value.IsKnown() {
		v := value.True()
		return &v
	}
	return nil
}

func BlockAtPos(body *hclsyntax.Body, pos hcl.Pos) *hclsyntax.Block {
	for _, b := range body.Blocks {
		if ContainsPos(b.Range(), pos) {
			return b
		}
	}
	return nil
}

func AttributeAtPos(block *hclsyntax.Block, pos hcl.Pos) *hclsyntax.Attribute {
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

func AttributeWithName(block *hclsyntax.Block, name string) *hclsyntax.Attribute {
	if block == nil {
		return nil
	}
	for _, attr := range block.Body.Attributes {
		if attr.Name == name {
			return attr
		}
	}
	return nil
}
