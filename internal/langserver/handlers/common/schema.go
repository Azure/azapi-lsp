package common

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func ExtractAzureResourceType(block *hclsyntax.Block) *string {
	typeAttr := AttributeWithName(block, "type")
	if typeAttr == nil {
		return nil
	}

	return ToLiteral(typeAttr.Expr)
}

func JsonEncodeExpressionToRangeMap(data []byte, expression hclsyntax.Expression) *RangeMap {
	r, err := rangeOfJsonEncodeBody(expression)
	if err != nil {
		return nil
	}
	tokens, _ := hclsyntax.LexExpression(data[r.Start.Byte:r.End.Byte], "", r.Start)
	return BuildRangeMap(tokens)
}

func rangeOfJsonEncodeBody(expression hclsyntax.Expression) (*hcl.Range, error) {
	if expression == nil {
		return nil, nil
	}
	if funcCallExpr, ok := expression.(*hclsyntax.FunctionCallExpr); ok {
		if funcCallExpr.Name != "jsonencode" {
			return nil, fmt.Errorf("expression is not funcation jsonencode")
		}
		if len(funcCallExpr.Args) != 1 {
			return nil, fmt.Errorf("invalid args length for jsonencode")
		}
		r := funcCallExpr.Args[0].Range()
		return &r, nil
	} else {
		return nil, fmt.Errorf("expression is not funcation call expression")
	}
}
