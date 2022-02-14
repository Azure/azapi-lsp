package complete

import "github.com/ms-henglu/azurerm-restapi-lsp/internal/langserver/schema"

func ignorePulledOutProperties(input []schema.Property) []schema.Property {
	res := make([]schema.Property, 0)
	// ignore properties pulled out from body
	for _, p := range input {
		if !isPropertyPulledOut(p) {
			res = append(res, p)
		}
	}
	return res
}

func ignorePulledOutPropertiesFromPropertySet(properties map[string]schema.Property) map[string]schema.Property {
	res := make(map[string]schema.Property, 0)
	// ignore properties pulled out from body
	for _, p := range properties {
		if !isPropertyPulledOut(p) {
			res[p.Name] = p
		}
	}
	return res
}

func isPropertyPulledOut(p schema.Property) bool {
	return p.Name == "name" || p.Name == "location" || p.Name == "identity" || p.Name == "tags"
}
