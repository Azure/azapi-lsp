// Code generated by generate-provider-schema/run.sh; DO NOT EDIT.
package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

var ProviderVersion = "tags/v3.114.0"

var ProviderSchemaInfo ProviderSchema

func init() {
	if err := json.Unmarshal(b, &ProviderSchemaInfo); err != nil {
		fmt.Fprintf(os.Stderr, "unmarshalling the provider schema: %s", err)
		os.Exit(1)
	}
}