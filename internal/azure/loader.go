package azure

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/Azure/azapi-lsp/internal/azure/types"
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
	for key, value := range azureSchema.Resources {
		if strings.EqualFold(key, resourceType) {
			for _, v := range value.Definitions {
				res = append(res, v.ApiVersion)
			}
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
	for key, value := range azureSchema.Resources {
		if strings.EqualFold(key, resourceType) {
			for _, v := range value.Definitions {
				if v.ApiVersion == apiVersion {
					return v.GetDefinition()
				}
			}
		}
	}
	return nil, fmt.Errorf("failed to find resource type %s api-version %s in azure schema index", resourceType, apiVersion)
}

func ListResourceFunctions(resourceType, apiVersion string) ([]FunctionDefinition, error) {
	res := make([]FunctionDefinition, 0)
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return nil, fmt.Errorf("failed to load azure schema index")
	}
	for key, value := range azureSchema.Functions {
		if strings.EqualFold(key, resourceType) {
			for _, v := range value.Definitions {
				if v.ApiVersion == apiVersion {
					res = append(res, v)
				}
			}
		}
	}

	return res, nil
}

func GetResourceFunction(resourceType, apiVersion, name string) (*types.ResourceFunctionType, error) {
	azureSchema := GetAzureSchema()
	if azureSchema == nil {
		return nil, fmt.Errorf("failed to load azure schema index")
	}
	for key, value := range azureSchema.Functions {
		if strings.EqualFold(key, resourceType) {
			for _, v := range value.Definitions {
				if v.ApiVersion == apiVersion {
					def, err := v.GetDefinition()
					if err == nil && def.Name == name {
						return def, nil
					}
				}
			}
		}
	}
	return nil, nil
}
