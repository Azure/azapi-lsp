package validate

import (
	"fmt"
	"strings"

	"github.com/Azure/azapi-lsp/internal/azure"
	"github.com/Azure/azapi-lsp/internal/azure/types"
	"github.com/Azure/azapi-lsp/internal/langserver/diagnostics"
	"github.com/Azure/azapi-lsp/internal/parser"
	"github.com/Azure/azapi-lsp/internal/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func NewDiagnostics(src []byte, filename string) diagnostics.Diagnostics {
	diags := diagnostics.NewDiagnostics()
	_, schemaDiags := ValidateFile(src, filename)
	diags.EmptyRootDiagnostic()
	validateDiags := make(map[string]hcl.Diagnostics)
	validateDiags[filename] = schemaDiags
	diags.Append("schema validate", validateDiags)
	return diags
}

func ValidateFile(src []byte, filename string) (*hcl.File, hcl.Diagnostics) {
	file, _ := hclsyntax.ParseConfig(src, filename, hcl.InitialPos)
	if file == nil {
		return nil, nil
	}
	body, isHcl := file.Body.(*hclsyntax.Body)
	if !isHcl {
		return nil, nil
	}

	diags := make([]*hcl.Diagnostic, 0)
	for _, block := range body.Blocks {
		if block.Type == "resource" && len(block.Labels) > 0 && strings.HasPrefix(block.Labels[0], "azapi_") {
			if diag := ValidateBlock(src, block); diag != nil {
				diags = append(diags, diag...)
			}
		}
	}
	return file, diags
}

func ValidateBlock(src []byte, block *hclsyntax.Block) hcl.Diagnostics {
	if block == nil {
		return nil
	}
	schemaValidationAttr := parser.AttributeWithName(block, "schema_validation_enabled")
	if schemaValidationAttr != nil {
		if enabled := parser.ToLiteralBoolean(schemaValidationAttr.Expr); enabled != nil && !*enabled {
			return nil
		}
	}
	attribute := parser.AttributeWithName(block, "body")
	if attribute == nil {
		return nil
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
	if len(block.Labels) >= 2 && block.Labels[0] == "azapi_operation" {
		parts := strings.Split(*typeValue, "@")
		if len(parts) != 2 {
			return nil
		}
		operationName := parser.ExtractOperation(block)
		if operationName != nil && len(*operationName) != 0 {
			resourceFuncDef, err := azure.GetResourceFunction(parts[0], parts[1], *operationName)
			if err != nil || resourceFuncDef == nil {
				return nil
			}
			bodyDef = resourceFuncDef
		}
	}

	hclNode := parser.JsonEncodeExpressionToHclNode(src, attribute.Expr)
	if hclNode == nil {
		return nil
	}
	if dummy, ok := hclNode.Children["dummy"]; ok {
		dummy.KeyRange = attribute.NameRange
		diags := Validate(dummy, bodyDef.AsTypeBase())
		// update resource doesn't need to check on required properties
		if block.Labels[0] == "azapi_update_resource" {
			res := hcl.Diagnostics{}
			for _, diag := range diags {
				// TODO: don't hardcode here
				if !strings.HasSuffix(diag.Summary, " is required, but no definition was found") {
					res = append(res, diag)
				}
			}
			return res
		} else {
			return diags
		}
	}
	return nil
}

func Validate(hclNode *parser.HclNode, typeBase *types.TypeBase) hcl.Diagnostics {
	if typeBase == nil || hclNode == nil {
		return nil
	}
	diags := make([]*hcl.Diagnostic, 0)
	switch t := (*typeBase).(type) {
	case *types.ArrayType:
		if !hclNode.IsValueArray() {
			summary := fmt.Sprintf("`%s`'s value `%s` is invalid, expect an array", hclNode.Key, *hclNode.Value)
			if hclNode.Value != nil {
				summary = fmt.Sprintf("`%s`'s value is invalid, expect an array", hclNode.Key)
			}
			diags = append(diags, newDiagnostic(summary, hclNode.ValueRange))
			break
		}
		if t.ItemType == nil {
			break
		}
		for _, child := range hclNode.Children {
			diags = append(diags, Validate(child, t.ItemType.Type)...)
		}
	case *types.DiscriminatedObjectType:
		if !hclNode.IsValueMap() {
			summary := fmt.Sprintf("`%s`'s value `%s` is invalid, expect an object", hclNode.Key, *hclNode.Value)
			if hclNode.Value != nil {
				summary = fmt.Sprintf("`%s`'s value is invalid, expect an object", hclNode.Key)
			}
			diags = append(diags, newDiagnostic(summary, hclNode.ValueRange))
			break
		}

		// check base properties
		otherProperties := make(map[string]*parser.HclNode)
		for key, value := range hclNode.Children {
			if def, ok := t.BaseProperties[key]; ok {
				if def.IsReadOnly() {
					diags = append(diags, newDiagnostic(ErrorShouldNotDefineReadOnly(key), value.KeyRange))
					continue
				}
				if def.Type != nil {
					diags = append(diags, Validate(value, def.Type.Type)...)
				}
			} else {
				otherProperties[key] = value
			}
		}

		// check required base properties
		for key, value := range t.BaseProperties {
			if value.IsRequired() && hclNode.Children[key] == nil {
				diags = append(diags, newDiagnostic(ErrorShouldDefine(key), hclNode.KeyRange))
			}
		}

		// check other properties which should be defined in discriminated objects
		if _, ok := otherProperties[t.Discriminator]; !ok {
			diags = append(diags, newDiagnostic(ErrorShouldDefine(t.Discriminator), hclNode.KeyRange))
			break
		}

		discriminator := ""

		discriminatorRange := hclNode.KeyRange
		discriminatorProp := otherProperties[t.Discriminator]
		if discriminatorProp != nil {
			if discriminatorProp.Value != nil {
				discriminator = strings.TrimPrefix(strings.TrimSuffix(strings.TrimSpace(*discriminatorProp.Value), `"`), `"`)
			}
			discriminatorRange = discriminatorProp.KeyRange
		}

		if len(discriminator) != 0 {
			switch {
			case t.Elements[discriminator] == nil:
				options := make([]string, 0)
				for key := range t.Elements {
					options = append(options, key)
				}
				diags = append(diags, newDiagnostic(ErrorNotMatchAnyValues(t.Discriminator, discriminator, options), discriminatorProp.ValueRange))
			case t.Elements[discriminator].Type != nil:
				other := &parser.HclNode{
					Key:        hclNode.Key,
					KeyRange:   hclNode.KeyRange,
					Children:   otherProperties,
					EqualRange: hclNode.EqualRange,
					ValueRange: hclNode.ValueRange,
				}
				diags = append(diags, Validate(other, t.Elements[discriminator].Type)...)
			}
		} else {
			diags = append(diags, newDiagnostic(ErrorMismatch(t.Discriminator, "string", fmt.Sprintf("%T", otherProperties[t.Discriminator])), discriminatorRange))
		}
	case *types.ObjectType:
		if !hclNode.IsValueMap() {
			summary := fmt.Sprintf("`%s`'s value `%s` is invalid, expect an object", hclNode.Key, *hclNode.Value)
			if hclNode.Value != nil {
				summary = fmt.Sprintf("`%s`'s value is invalid, expect an object", hclNode.Key)
			}
			diags = append(diags, newDiagnostic(summary, hclNode.ValueRange))
			break
		}
		// check properties defined in body, but not in schema
		for key, value := range hclNode.Children {
			if def, ok := t.Properties[key]; ok {
				if def.IsReadOnly() {
					diags = append(diags, newDiagnostic(ErrorShouldNotDefineReadOnly(key), value.KeyRange))
					continue
				}
				if def.Type != nil {
					diags = append(diags, Validate(value, def.Type.Type)...)
				}
				continue
			}
			if t.AdditionalProperties != nil {
				diags = append(diags, Validate(value, t.AdditionalProperties.Type)...)
			} else {
				options := make([]string, 0)
				for key := range t.Properties {
					options = append(options, key)
				}
				diags = append(diags, newDiagnostic(ErrorShouldNotDefine(key, options), value.KeyRange))
			}
		}

		// check properties required in schema, but not in body
		for key, value := range t.Properties {
			if value.IsRequired() && hclNode.Children[key] == nil {
				// skip name in body
				if hclNode.Key == "dummy" && (key == "name" || key == "location") {
					continue
				}
				diags = append(diags, newDiagnostic(ErrorShouldDefine(key), hclNode.KeyRange))
			}
		}
	case *types.ResourceType:
		if t.Body != nil {
			return Validate(hclNode, t.Body.Type)
		}
	case *types.ResourceFunctionType:
		if t.Input != nil {
			return Validate(hclNode, t.Input.Type)
		}
	case *types.BuiltInType:
	case *types.StringLiteralType:
		if hclNode.Value != nil {
			value := strings.TrimSpace(*hclNode.Value)
			if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
				value = strings.TrimPrefix(strings.TrimSuffix(value, `"`), `"`)
				if value != t.Value {
					diags = append(diags, newDiagnostic(ErrorMismatch(hclNode.Key, t.Value, value), hclNode.ValueRange))
				}
			}
		}
	case *types.UnionType:
		valid := false
		for _, element := range t.Elements {
			if element.Type == nil {
				continue
			}
			temp := Validate(hclNode, element.Type)
			if len(temp) == 0 {
				valid = true
				break
			}
		}
		if !valid {
			options := make([]string, 0)
			for _, element := range t.Elements {
				if element.Type != nil {
					if stringLiteralType, ok := (*element.Type).(*types.StringLiteralType); ok {
						options = append(options, stringLiteralType.Value)
					}
				}
			}
			if len(options) == 0 {
				diags = append(diags, newDiagnostic(ErrorNotMatchAny(hclNode.Key), hclNode.GetRange()))
			} else {
				value := ""
				if hclNode.Value != nil {
					value = *hclNode.Value
				}
				diags = append(diags, newDiagnostic(ErrorNotMatchAnyValues(hclNode.Key, value, options), hclNode.ValueRange))
			}
		}
	}
	return diags
}

func newDiagnostic(summary string, r hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Summary:  summary,
		Subject:  utils.Range(r),
		Severity: hcl.DiagError,
	}
}

func ErrorMismatch(key, expected, actual string) string {
	return fmt.Sprintf("`%s` is invalid, expect `%s` but got `%s`", strings.TrimPrefix(key, "."), expected, actual)
}

func ErrorNotMatchAny(key string) string {
	return fmt.Sprintf("`%s` doesn't match any accepted values", strings.TrimPrefix(key, "."))
}

func ErrorNotMatchAnyValues(key string, value string, options []string) string {
	suggestion := getSuggestion(value, options)
	return fmt.Sprintf("`%s`'s value `%s` is invalid. The supported values are [%s]. Do you mean `%s`? ",
		strings.TrimPrefix(key, "."),
		value,
		strings.Join(options, ", "),
		suggestion)
}

func ErrorShouldNotDefineReadOnly(key string) string {
	return fmt.Sprintf("`%s` is not expected here, it's read only", strings.TrimPrefix(key, "."))
}

func ErrorShouldNotDefine(key string, options []string) string {
	suggestion := getSuggestion(key, options)
	return fmt.Sprintf("`%s` is not expected here. Do you mean `%s`? ", strings.TrimPrefix(key, "."), strings.TrimPrefix(suggestion, "."))
}

func ErrorShouldDefine(key string) string {
	return fmt.Sprintf("`%s` is required, but no definition was found", strings.TrimPrefix(key, "."))
}

func getSuggestion(value string, options []string) string {
	suggestion := ""
	distance := 1 << 16
	for _, option := range options {
		if dist := editDistance(value, option); dist < distance {
			distance = dist
			suggestion = option
		}
	}
	return suggestion
}

func editDistance(a, b string) int {
	n, m := len(a), len(b)
	f := make([][]int, n+1)
	for i := range f {
		f[i] = make([]int, m+1, 1<<16)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			f[i][j] = 1 << 16
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if a[i-1] == b[j-1] {
				f[i][j] = f[i-1][j-1]
			}
			if f[i][j] > f[i-1][j]+1 {
				f[i][j] = f[i-1][j] + 1
			}
			if f[i][j] > f[i][j-1]+1 {
				f[i][j] = f[i][j-1] + 1
			}
			if f[i][j] > f[i-1][j-1]+1 {
				f[i][j] = f[i-1][j-1] + 1
			}
		}
	}
	return f[n][m]
}
