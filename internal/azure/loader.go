package azure

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/ms-henglu/azurerm-restapi-lsp/internal/azure/types"
)

var schema *Schema

//go:embed generated
var StaticFiles embed.FS

func GetAzureSchema() *Schema {
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

func ListProperties(resourceType *types.ResourceType, path string) []string {
	if resourceType == nil {
		return []string{}
	}
	def := GetDef(resourceType.AsTypeBase(), path)
	return GetAllowedProperties(def)
}

func ListValues(resourceType *types.ResourceType, path string) []string {
	if resourceType == nil {
		return []string{}
	}
	def := GetDef(resourceType.AsTypeBase(), path)
	return GetAllowedValues(def)
}

func GetDef(resourceType *types.TypeBase, path string) *types.TypeBase{
	if len(path) == 0 || resourceType == nil {
		return resourceType
	}
	key := path
	rest := ""
	if strings.Contains(key, ".") {
		key = path[0: strings.Index(path, ".")]
		rest = path[len(key) + 1:]
	}
	switch t := (*resourceType).(type) {
	case *types.ArrayType:
		if t.ItemType != nil {
			return GetDef(t.ItemType.Type, path)
		}
		return nil
	case *types.DiscriminatedObjectType:
		if value, ok := t.BaseProperties[key]; ok {
			if value.Type != nil {
				return GetDef(value.Type.Type, rest)
			}
		}
		for _, discriminator := range t.Elements {
			if resourceType := GetDef(discriminator.Type, path); resourceType != nil {
				return resourceType
			}
		}
		break
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
	case *types.StringLiteralType:
	case *types.UnionType:
		return resourceType
	}
	return nil
}

func GetAllowedProperties(resourceType *types.TypeBase) []string{
	if resourceType == nil {
		return []string{}
	}
	keys := make([]string, 0)
	switch t := (*resourceType).(type) {
	case *types.ArrayType:
		return keys
	case *types.DiscriminatedObjectType:
		for key := range t.BaseProperties {
			keys = append(keys, key)
		}
		for _, discriminator := range t.Elements {
			keys = append(keys, GetAllowedProperties(discriminator.Type)...)
		}
		break
	case *types.ObjectType:
		for key := range t.Properties {
			keys = append(keys, key)
		}
		if t.AdditionalProperties != nil {
			keys = append(keys, GetAllowedProperties(t.AdditionalProperties.Type)...)
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
		return keys
	}
	return keys
}


func GetAllowedValues(resourceType *types.TypeBase) []string{
	if resourceType == nil {
		return []string{}
	}
	values := make([]string, 0)

	return values
}