package schema

import "github.com/Azure/azapi-lsp/internal/azure/types"

func GetRequiredPropertySet(typeBase *types.TypeBase) []PropertySet {
	if typeBase == nil {
		return nil
	}
	switch t := (*typeBase).(type) {
	case *types.DiscriminatedObjectType:
		res := make([]PropertySet, 0)
		for name, element := range t.Elements {
			if element == nil {
				continue
			}
			propertySet := GetRequiredPropertySet(element.Type)
			if len(propertySet) == 1 {
				requiredProps := propertySet[0].Properties
				requiredProps[t.Discriminator] = Property{
					Name:  t.Discriminator,
					Value: name,
				}

				res = append(res, PropertySet{
					Name:       name,
					Properties: requiredProps,
				})
			}
		}
		return res
	case *types.ObjectType:
		requiredProps := make(map[string]Property)
		for propName, prop := range t.Properties {
			if prop.IsRequired() {
				if value := PropertyFromObjectProperty(propName, prop); value != nil {
					requiredProps[value.Name] = *value
				}
			}
		}
		return []PropertySet{{
			Name:       t.Name,
			Properties: requiredProps,
		}}
	case *types.ResourceType:
		if t.Body != nil {
			return GetRequiredPropertySet(t.Body.Type)
		}
	case *types.UnionType:
		res := make([]PropertySet, 0)
		for _, element := range t.Elements {
			if element.Type == nil {
				continue
			}
			res = append(res, GetRequiredPropertySet(element.Type)...)
		}
		return res
	}
	return nil
}
