package azure

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure/types"
)

var schema *Schema

//go:embed generated
var StaticFiles embed.FS

var mutex sync.Mutex

func GetAzureSchema() *Schema {
	mutex.Lock()
	defer mutex.Unlock()
	if schema == nil {
		data, err := StaticFiles.ReadFile("generated/index.json")
		if err != nil {
			log.Printf("[ERROR] failed to load schema index: %+v", err)
			return nil
		}
		err = json.Unmarshal(data, &schema)
		if err != nil {
			log.Printf("[ERROR] failed to unmarshal schema index: %+v", err)
			return nil
		}
	}
	return schema
}

func GetApiVersions(resourceType string) []string {
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return []string{}
	}
	res := make([]string, 0)
	resourceType = strings.ToUpper(resourceType)
	if azureSchema.Resources[resourceType] != nil {
		for _, v := range azureSchema.Resources[resourceType].Definitions {
			res = append(res, v.ApiVersion)
		}
	}
	sort.Strings(res)
	return res
}

func GetResourceDefinitionByResourceType(azureResourceType string) (*types.ResourceType, error) {
	parts := strings.Split(azureResourceType, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("input %s is invalid", azureResourceType)
	}

	return GetResourceDefinition(parts[0], parts[1])
}

func GetResourceDefinition(resourceType, apiVersion string) (*types.ResourceType, error) {
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return nil, fmt.Errorf("failed to load azure schema index")
	}
	resourceType = strings.ToUpper(resourceType)
	if azureSchema.Resources[resourceType] != nil {
		for _, v := range azureSchema.Resources[resourceType].Definitions {
			if v.ApiVersion == apiVersion {
				return v.GetDefinition()
			}
		}
	}
	return nil, fmt.Errorf("failed to find resource type %s api-version %s in azure schema index", resourceType, apiVersion)
}

func ListProperties(resourceType *types.ResourceType, path string) []Property {
	if resourceType == nil {
		return []Property{}
	}
	defs := GetDef(resourceType.AsTypeBase(), path)
	props := make([]Property, 0)
	for _, def := range defs {
		props = append(props, GetAllowedProperties(def)...)
	}
	return props
}

func ListValues(resourceType *types.ResourceType, path string) []string {
	if resourceType == nil {
		return []string{}
	}
	defs := GetDef(resourceType.AsTypeBase(), path)
	values := make([]string, 0)
	for _, def := range defs {
		values = append(values, GetAllowedValues(def)...)
	}
	return values
}

func GetDef(resourceType *types.TypeBase, path string) []*types.TypeBase {
	if resourceType == nil {
		return nil
	}
	if len(path) == 0 {
		return []*types.TypeBase{resourceType}
	}
	key := path
	rest := ""
	if strings.Contains(key, ".") {
		key = path[0:strings.Index(path, ".")]
		rest = path[len(key)+1:]
	}
	switch t := (*resourceType).(type) {
	case *types.ArrayType:
		if t.ItemType != nil {
			if key == "#" {
				return GetDef(t.ItemType.Type, rest)
			}
			return GetDef(t.ItemType.Type, path)
		}
		return nil
	case *types.DiscriminatedObjectType:
		if value, ok := t.BaseProperties[key]; ok {
			if value.Type != nil {
				return GetDef(value.Type.Type, rest)
			}
		}
		res := make([]*types.TypeBase, 0)
		for _, discriminator := range t.Elements {
			if resourceType := GetDef(discriminator.Type, path); resourceType != nil {
				res = append(res, resourceType...)
			}
		}
		return res
	case *types.ObjectType:
		if value, ok := t.Properties[key]; ok {
			if value.Type != nil {
				return GetDef(value.Type.Type, rest)
			}
		}
		if t.AdditionalProperties != nil {
			return GetDef(t.AdditionalProperties.Type, rest)
		}
		break
	case *types.ResourceType:
		if t.Body != nil {
			return GetDef(t.Body.Type, path)
		}
		break
	case *types.BuiltInType:
		return []*types.TypeBase{resourceType}
	case *types.StringLiteralType:
		return []*types.TypeBase{resourceType}
	case *types.UnionType:
		res := make([]*types.TypeBase, 0)
		for _, element := range t.Elements {
			res = append(res, GetDef(element.Type, path)...)
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
		break
	case *types.ObjectType:
		for key, value := range t.Properties {
			if prop := PropertyFromObjectProperty(key, value); prop != nil {
				props = append(props, *prop)
			}
		}
		if t.AdditionalProperties != nil {
			props = append(props, GetAllowedProperties(t.AdditionalProperties.Type)...)
		}
		break
	case *types.ResourceType:
		if t.Body != nil {
			return GetAllowedProperties(t.Body.Type)
		}
		break
	case *types.BuiltInType:
	case *types.StringLiteralType:
	case *types.UnionType:
		return props
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
		break
	case *types.StringLiteralType:
		return []string{t.Value}
	case *types.UnionType:
		for _, element := range t.Elements {
			values = append(values, GetAllowedValues(element.Type)...)
		}
		return values
	case *types.DiscriminatedObjectType:
	case *types.ObjectType:
	case *types.ArrayType:
	case *types.BuiltInType:
	}
	return values
}

type Modifier string

const (
	Optional Modifier = "Optional"
	Required Modifier = "Required"
)

type Property struct {
	Name        string
	Description string
	Modifier    Modifier
	Type        string
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
	case *types.BuiltInType:
		return t.Kind.String()
	case *types.StringLiteralType:
		return "string"
	case *types.UnionType:
		for _, element := range t.Elements {
			return GetTypeName(element.Type)
		}
	}
	return ""
}
