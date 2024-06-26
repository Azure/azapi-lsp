package schema

import (
	"fmt"
	"strings"

	"github.com/Azure/azapi-lsp/internal/azure/types"
	"github.com/Azure/azapi-lsp/internal/parser"
)

func GetDef(resourceType *types.TypeBase, hclNodes []*parser.HclNode, index int) []*types.TypeBase {
	if resourceType == nil {
		return nil
	}
	if len(hclNodes) == index {
		if t, ok := (*resourceType).(*types.DiscriminatedObjectType); ok {
			if discriminator, ok := hclNodes[index-1].Children[t.Discriminator]; ok && discriminator != nil && discriminator.Value != nil {
				if discriminatorValue := strings.Trim(*discriminator.Value, `"`); len(discriminatorValue) > 0 {
					if t.Elements[discriminatorValue] != nil && t.Elements[discriminatorValue].Type != nil {
						selectedDiscriminatedObjectType := &types.DiscriminatedObjectType{
							Name:           t.Name,
							Discriminator:  t.Discriminator,
							BaseProperties: t.BaseProperties,
							Elements: map[string]*types.TypeReference{
								discriminatorValue: t.Elements[discriminatorValue],
							},
						}
						return []*types.TypeBase{selectedDiscriminatedObjectType.AsTypeBase()}
					}
				}
			}
		}
		return []*types.TypeBase{resourceType}
	}
	key := hclNodes[index].Key
	switch t := (*resourceType).(type) {
	case *types.ArrayType:
		if t.ItemType != nil {
			if strings.Contains(key, ".") {
				return GetDef(t.ItemType.Type, hclNodes, index+1)
			}
			return GetDef(t.ItemType.Type, hclNodes, index)
		}
		return nil
	case *types.DiscriminatedObjectType:
		if value, ok := t.BaseProperties[key]; ok {
			if value.Type != nil {
				return GetDef(value.Type.Type, hclNodes, index+1)
			}
		}
		if index != 0 {
			if discriminator, ok := hclNodes[index-1].Children[t.Discriminator]; ok && discriminator != nil && discriminator.Value != nil {
				if discriminatorValue := strings.Trim(*discriminator.Value, `"`); len(discriminatorValue) > 0 {
					if t.Elements[discriminatorValue] != nil && t.Elements[discriminatorValue].Type != nil {
						return GetDef(t.Elements[discriminatorValue].Type, hclNodes, index)
					}
				}
			}
		}
		res := make([]*types.TypeBase, 0)
		for _, discriminator := range t.Elements {
			if resourceType := GetDef(discriminator.Type, hclNodes, index); resourceType != nil {
				res = append(res, resourceType...)
			}
		}
		return res
	case *types.ObjectType:
		if value, ok := t.Properties[key]; ok {
			if value.Type != nil {
				return GetDef(value.Type.Type, hclNodes, index+1)
			}
		}
		if t.AdditionalProperties != nil {
			return GetDef(t.AdditionalProperties.Type, hclNodes, index+1)
		}
	case *types.ResourceType:
		if t.Body != nil {
			return GetDef(t.Body.Type, hclNodes, index+1)
		}
	case *types.ResourceFunctionType:
		if t.Input != nil {
			return GetDef(t.Input.Type, hclNodes, index+1)
		}
	case *types.AnyType:
		return []*types.TypeBase{resourceType}
	case *types.BooleanType:
		return []*types.TypeBase{resourceType}
	case *types.IntegerType:
		return []*types.TypeBase{resourceType}
	case *types.StringType:
		return []*types.TypeBase{resourceType}
	case *types.StringLiteralType:
		return []*types.TypeBase{resourceType}
	case *types.UnionType:
		res := make([]*types.TypeBase, 0)
		for _, element := range t.Elements {
			res = append(res, GetDef(element.Type, hclNodes, index)...)
		}
		return res
	}
	return nil
}

func GetAllowedProperties(resourceType *types.TypeBase) []Property {
	if resourceType == nil {
		return []Property{}
	}
	props := make([]Property, 0)
	switch t := (*resourceType).(type) {
	case *types.ArrayType:
		return props
	case *types.DiscriminatedObjectType:
		for key, value := range t.BaseProperties {
			if prop := PropertyFromObjectProperty(key, value); prop != nil {
				props = append(props, *prop)
			}
		}
		for _, discriminator := range t.Elements {
			props = append(props, GetAllowedProperties(discriminator.Type)...)
		}
	case *types.ObjectType:
		for key, value := range t.Properties {
			if prop := PropertyFromObjectProperty(key, value); prop != nil {
				props = append(props, *prop)
			}
		}
		if t.AdditionalProperties != nil {
			props = append(props, GetAllowedProperties(t.AdditionalProperties.Type)...)
		}
	case *types.ResourceType:
		if t.Body != nil {
			return GetAllowedProperties(t.Body.Type)
		}
	case *types.AnyType:
	case *types.BooleanType:
	case *types.IntegerType:
	case *types.StringType:
	case *types.StringLiteralType:
	case *types.UnionType:
	}
	return props
}

func GetAllowedValues(resourceType *types.TypeBase) []string {
	if resourceType == nil {
		return nil
	}
	values := make([]string, 0)
	switch t := (*resourceType).(type) {
	case *types.ResourceType:
		if t.Body != nil {
			return GetAllowedValues(t.Body.Type)
		}
	case *types.StringLiteralType:
		return []string{fmt.Sprintf(`"%s"`, t.Value)}
	case *types.UnionType:
		for _, element := range t.Elements {
			values = append(values, GetAllowedValues(element.Type)...)
		}
		return values
	case *types.DiscriminatedObjectType:
	case *types.ObjectType:
	case *types.ArrayType:
	case *types.AnyType:
	case *types.BooleanType:
		values = append(values, "true", "false")
	case *types.IntegerType:
	case *types.StringType:
	}
	return values
}

func PropertyFromObjectProperty(propertyName string, property types.ObjectProperty) *Property {
	if property.IsReadOnly() {
		return nil
	}
	description := ""
	if property.Description != nil {
		description = *property.Description
	}
	modifier := Optional
	if property.IsRequired() {
		modifier = Required
	}
	propertyType := ""
	if property.Type != nil {
		propertyType = GetTypeName(property.Type.Type)
	}
	return &Property{
		Name:        propertyName,
		Description: description,
		Modifier:    modifier,
		Type:        propertyType,
	}
}

func GetTypeName(typeBase *types.TypeBase) string {
	if typeBase == nil {
		return ""
	}
	switch t := (*typeBase).(type) {
	case *types.ArrayType:
		return "array"
	case *types.DiscriminatedObjectType:
		return "object"
	case *types.ObjectType:
		return "object"
	case *types.ResourceType:
		return "object"
	case *types.AnyType:
		return "any"
	case *types.BooleanType:
		return "boolean"
	case *types.IntegerType:
		return "int"
	case *types.StringType:
		return "string"
	case *types.StringLiteralType:
		return "string"
	case *types.UnionType:
		for _, element := range t.Elements {
			return GetTypeName(element.Type)
		}
	}
	return ""
}
