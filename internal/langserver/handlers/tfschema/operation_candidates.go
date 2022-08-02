package tfschema

import (
	"fmt"
	"strings"

	"github.com/Azure/azapi-lsp/internal/azure"
	"github.com/Azure/azapi-lsp/internal/parser"
	lsp "github.com/Azure/azapi-lsp/internal/protocol"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func actionCandidates(_ []byte, _ string, block *hclsyntax.Block, attribute *hclsyntax.Attribute, pos hcl.Pos, _ *Property) []lsp.CompletionItem {
	typeValue := parser.ExtractAzureResourceType(block)
	if typeValue == nil {
		return nil
	}

	parts := strings.Split(*typeValue, "@")
	if len(parts) != 2 {
		return nil
	}

	functions, err := azure.ListResourceFunctions(parts[0], parts[1])
	if err != nil {
		return nil
	}

	values := make([]string, 0)
	for _, function := range functions {
		def, err := function.GetDefinition()
		if err == nil && def != nil {
			values = append(values, fmt.Sprintf(`"%s"`, def.Name))
		}
	}

	return valueCandidates(values, editRangeFromExprRange(attribute.Expr, pos), false)
}
