package validate

import (
	"fmt"
	"strings"

	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/diagnostics"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure/types"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/handlers/common"
	"github.com/ms-henglu/azurerm-restapi-lsp/internal/utils"
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
		if block.Type == "resource" && len(block.Labels) > 0 && strings.HasPrefix(block.Labels[0], "azurerm-restapi") {
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
	schemaValidationAttr := common.AttributeWithName(block, "schema_validation_enabled")
	if schemaValidationAttr != nil {
		if enabled := common.ToLiteralBoolean(schemaValidationAttr.Expr); enabled != nil && !*enabled {
			return nil
		}
	}
	attribute := common.AttributeWithName(block, "body")
	if attribute == nil {
		return nil
	}

	typeValue := common.ExtractAzureResourceType(block)
	if typeValue == nil {
		return nil
	}
	def, _ := azure.GetResourceDefinitionByResourceType(*typeValue)
	if def == nil {
		return nil
	}

	rangeMap := common.JsonEncodeExpressionToRangeMap(src, attribute.Expr)
	if rangeMap == nil {
		return nil
	}
	if dummy, ok := rangeMap.Children["dummy"]; ok {
		dummy.KeyRange = attribute.NameRange
		diags := Validate(dummy, def.AsTypeBase())
		// patch resource doesn't need to check on required properties
		if block.Labels[0] == "azurerm-restapi_patch_resource" {
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

func Validate(rangeMap *common.RangeMap, typeBase *types.TypeBase) hcl.Diagnostics {
	if typeBase == nil || rangeMap == nil {
		return nil
	}
	diags := make([]*hcl.Diagnostic, 0)
	switch t := (*typeBase).(type) {
	case *types.ArrayType:
		if !rangeMap.IsValueArray() {
			summary := fmt.Sprintf("`%s`'s value `%s` is invalid, expect an array", rangeMap.Key, *rangeMap.Value)
			if rangeMap.Value != nil {
				summary = fmt.Sprintf("`%s`'s value is invalid, expect an array", rangeMap.Key)
			}
			diags = append(diags, newDiagnostic(summary, rangeMap.ValueRange))
			break
		}
		if t.ItemType == nil {
			break
		}
		for _, child := range rangeMap.Children {
			diags = append(diags, Validate(child, t.ItemType.Type)...)
		}
	case *types.DiscriminatedObjectType:
		if !rangeMap.IsValueMap() {
			summary := fmt.Sprintf("`%s`'s value `%s` is invalid, expect an object", rangeMap.Key, *rangeMap.Value)
			if rangeMap.Value != nil {
				summary = fmt.Sprintf("`%s`'s value is invalid, expect an object", rangeMap.Key)
			}
			diags = append(diags, newDiagnostic(summary, rangeMap.ValueRange))
			break
		}

		// check base properties
		otherProperties := make(map[string]*common.RangeMap)
		for key, value := range rangeMap.Children {
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
			if value.IsRequired() && rangeMap.Children[key] == nil {
				diags = append(diags, newDiagnostic(ErrorShouldDefine(key), rangeMap.KeyRange))
			}
		}

		// check other properties which should be defined in discriminated objects
		if _, ok := otherProperties[t.Discriminator]; !ok {
			diags = append(diags, newDiagnostic(ErrorShouldDefine(t.Discriminator), rangeMap.KeyRange))
			break
		}

		discriminator := ""

		discriminatorRange := rangeMap.KeyRange
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
				other := &common.RangeMap{
					Key:        rangeMap.Key,
					KeyRange:   rangeMap.KeyRange,
					Children:   otherProperties,
					EqualRange: rangeMap.EqualRange,
					ValueRange: rangeMap.ValueRange,
				}
				diags = append(diags, Validate(other, t.Elements[discriminator].Type)...)
			}
		} else {
			diags = append(diags, newDiagnostic(ErrorMismatch(t.Discriminator, "string", fmt.Sprintf("%T", otherProperties[t.Discriminator])), discriminatorRange))
		}
	case *types.ObjectType:
		if !rangeMap.IsValueMap() {
			summary := fmt.Sprintf("`%s`'s value `%s` is invalid, expect an object", rangeMap.Key, *rangeMap.Value)
			if rangeMap.Value != nil {
				summary = fmt.Sprintf("`%s`'s value is invalid, expect an object", rangeMap.Key)
			}
			diags = append(diags, newDiagnostic(summary, rangeMap.ValueRange))
			break
		}
		// check properties defined in body, but not in schema
		for key, value := range rangeMap.Children {
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
			if value.IsRequired() && rangeMap.Children[key] == nil {
				// skip name in body
				if rangeMap.Key == "dummy" && (key == "name" || key == "location") {
					continue
				}
				diags = append(diags, newDiagnostic(ErrorShouldDefine(key), rangeMap.KeyRange))
			}
		}
	case *types.ResourceType:
		if t.Body != nil {
			return Validate(rangeMap, t.Body.Type)
		}
	case *types.BuiltInType:
	case *types.StringLiteralType:
		if rangeMap.Value != nil {
			value := strings.TrimSpace(*rangeMap.Value)
			if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
				value = strings.TrimPrefix(strings.TrimSuffix(value, `"`), `"`)
				if value != t.Value {
					diags = append(diags, newDiagnostic(ErrorMismatch(rangeMap.Key, t.Value, value), rangeMap.ValueRange))
				}
			}
		}
	case *types.UnionType:
		valid := false
		for _, element := range t.Elements {
			if element.Type == nil {
				continue
			}
			temp := Validate(rangeMap, element.Type)
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
				diags = append(diags, newDiagnostic(ErrorNotMatchAny(rangeMap.Key), rangeMap.GetRange()))
			} else {
				value := ""
				if rangeMap.Value != nil {
					value = *rangeMap.Value
				}
				diags = append(diags, newDiagnostic(ErrorNotMatchAnyValues(rangeMap.Key, value, options), rangeMap.ValueRange))
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
