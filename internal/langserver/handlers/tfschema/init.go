package tfschema

import lsp "github.com/Azure/azapi-lsp/internal/protocol"

func init() {
	Resources = make([]Resource, 0)

	Resources = append(Resources,
		Resource{
			Name: "resource.azapi_resource",
			Properties: []Property{
				{
					Name:                "type",
					Modifier:            "Required",
					Type:                "string <resource-type>@<api-version>",
					Description:         "Azure Resource Manager type.",
					CompletionNewText:   `type = "$0"`,
					ValueCandidatesFunc: resourceTypeCandidates,
				},

				{
					Name:              "name",
					Modifier:          "Required",
					Type:              "string",
					Description:       "Specifies the name of the azure resource. Changing this forces a new resource to be created.",
					CompletionNewText: `name = "$0"`,
				},

				{
					Name:              "parent_id",
					Modifier:          "Required",
					Type:              "string",
					Description:       "The ID of the azure resource in which this resource is created. Changing this forces a new resource to be created.",
					CompletionNewText: `parent_id = $0`,
				},

				{
					Name:                "location",
					Modifier:            "Optional",
					Type:                "string",
					Description:         "The Azure Region where the azure resource should exist.",
					CompletionNewText:   `location = "$0"`,
					ValueCandidatesFunc: locationCandidates,
				},

				{
					Name:              "identity",
					Modifier:          "Optional",
					Type:              "block",
					Description:       "Managed identities which should be assigned to the azure resource.",
					CompletionNewText: "identity {\n\ttype = \"$1\"\n\tidentity_ids = [$2]\n}\n",
					NestedProperties: []Property{
						{
							Name:                "type",
							Modifier:            "Required",
							Type:                "string",
							Description:         "The Type of Identity which should be used for this azure resource. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned, UserAssigned`.",
							CompletionNewText:   `type = "$0"`,
							ValueCandidatesFunc: identityTypesCandidates,
						},

						{
							Name:              "identity_ids",
							Modifier:          "Optional",
							Type:              "list<string>",
							Description:       "A list of User Managed Identity ID's which should be assigned to the azure resource.",
							CompletionNewText: "identity_ids = [$0]",
						},
					},
				},

				{
					Name:                  "body",
					Modifier:              "Optional",
					Type:                  "string <JSON>",
					Description:           "A JSON object that contains the request body used to create and update azure resource.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{bodyJsonencodeFuncCandidate()}),
					GenericCandidatesFunc: bodyCandidates,
				},

				{
					Name:              "tags",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of tags which should be assigned to the azure resource.",
					CompletionNewText: `tags = $0`,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				{
					Name:                "schema_validation_enabled",
					Modifier:            "Optional",
					Type:                "bool",
					Description:         "Whether enabled the validation on `type` and `body` with embedded schema. Defaults to `true`.",
					CompletionNewText:   `schema_validation_enabled = $0`,
					ValueCandidatesFunc: boolCandidates,
				},

				{
					Name:              "locks",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of ARM resource IDs which are used to avoid create/modify/delete azapi resources at the same time.",
					CompletionNewText: `locks = [$0]`,
				},

				{
					Name:              "ignore_body_changes",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of properties that should be ignored when comparing the `body` with its current state..",
					CompletionNewText: `ignore_body_changes = [$0]`,
				},

				{
					Name:                "ignore_casing",
					Modifier:            "Optional",
					Type:                "bool",
					Description:         "Whether ignore incorrect casing returned in `body` to suppress plan-diff. Defaults to `false`.",
					CompletionNewText:   `ignore_casing = $0`,
					ValueCandidatesFunc: boolCandidates,
				},

				{
					Name:                "ignore_missing_property",
					Modifier:            "Optional",
					Type:                "bool",
					Description:         "Whether ignore not returned properties like credentials in `body` to suppress plan-diff. Defaults to `false`.",
					CompletionNewText:   `ignore_missing_property = $0`,
					ValueCandidatesFunc: boolCandidates,
				},
			},
		},
		Resource{
			Name: "resource.azapi_update_resource",
			Properties: []Property{
				{
					Name:                "type",
					Modifier:            "Required",
					Type:                "string <resource-type>@<api-version>",
					Description:         "Azure Resource Manager type.",
					CompletionNewText:   `type = "$0"`,
					ValueCandidatesFunc: resourceTypeCandidates,
				},

				{
					Name:              "name",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "Specifies the name of the azure resource. Changing this forces a new resource to be created.\n\nConfiguring `name` and `parent_id` is an alternative way to configure `resource_id`.",
					CompletionNewText: `name = "$0"`,
				},

				{
					Name:              "parent_id",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "The ID of the azure resource in which this resource is created. Changing this forces a new resource to be created.\n\nConfiguring `name` and `parent_id` is an alternative way to configure `resource_id`.",
					CompletionNewText: `parent_id = $0`,
				},

				{
					Name:              "resource_id",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "The ID of an existing azure source. Changing this forces a new resource to be created.\n\nConfiguring `name` and `parent_id` is an alternative way to configure `resource_id`.",
					CompletionNewText: `resource_id = $0`,
				},

				{
					Name:                  "body",
					Modifier:              "Optional",
					Type:                  "string <JSON>",
					Description:           "A JSON object that contains the request body used to create and update azure resource.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{bodyJsonencodeFuncCandidate()}),
					GenericCandidatesFunc: bodyCandidates,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				{
					Name:              "locks",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of ARM resource IDs which are used to avoid create/modify/delete azapi resources at the same time.",
					CompletionNewText: `locks = [$0]`,
				},

				{
					Name:                "ignore_casing",
					Modifier:            "Optional",
					Type:                "bool",
					Description:         "Whether ignore incorrect casing returned in `body` to suppress plan-diff. Defaults to `false`.",
					CompletionNewText:   `ignore_casing = $0`,
					ValueCandidatesFunc: boolCandidates,
				},

				{
					Name:                "ignore_missing_property",
					Modifier:            "Optional",
					Type:                "bool",
					Description:         "Whether ignore not returned properties like credentials in `body` to suppress plan-diff. Defaults to `false`.",
					CompletionNewText:   `ignore_missing_property = $0`,
					ValueCandidatesFunc: boolCandidates,
				},
			},
		},
		Resource{
			Name: "data.azapi_resource",
			Properties: []Property{
				{
					Name:                "type",
					Modifier:            "Required",
					Type:                "string <resource-type>@<api-version>",
					Description:         "Azure Resource Manager type.",
					CompletionNewText:   `type = "$0"`,
					ValueCandidatesFunc: dataSourceTypeCandidates,
				},

				{
					Name:              "name",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "Specifies the name of the azure resource.",
					CompletionNewText: `name = "$0"`,
				},

				{
					Name:              "parent_id",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "The ID of the azure resource in which this resource is created.",
					CompletionNewText: `parent_id = $0`,
				},

				{
					Name:              "resource_id",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "The ID of an existing azure source.\n\nConfiguring `name` and `parent_id` is an alternative way to configure `resource_id`.",
					CompletionNewText: `resource_id = $0`,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},
			},
		},
		Resource{
			Name: "resource.azapi_resource_action",
			Properties: []Property{
				{
					Name:                "type",
					Modifier:            "Required",
					Type:                "string <resource-type>@<api-version>",
					Description:         "Azure Resource Manager type.",
					CompletionNewText:   `type = "$0"`,
					ValueCandidatesFunc: resourceTypeCandidates,
				},

				{
					Name:              "resource_id",
					Modifier:          "Required",
					Type:              "string",
					Description:       "The ID of an existing azure source.",
					CompletionNewText: `resource_id = $0`,
				},

				{
					Name:                  "action",
					Modifier:              "Optional",
					Type:                  "string",
					Description:           "Specifies the name of the azure resource action.",
					CompletionNewText:     `action = "$0"`,
					GenericCandidatesFunc: actionCandidates,
				},

				{
					Name:                "method",
					Modifier:            "Optional",
					Type:                "string",
					Description:         "Specifies the Http method of the azure resource action. Defaults to `POST`",
					CompletionNewText:   `method = $0`,
					ValueCandidatesFunc: resourceHttpMethodCandidates,
				},

				{
					Name:                  "body",
					Modifier:              "Optional",
					Type:                  "string <JSON>",
					Description:           "A JSON object that contains the request body.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{bodyJsonencodeFuncCandidate()}),
					GenericCandidatesFunc: bodyCandidates,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},
			},
		},
		Resource{
			Name: "data.azapi_resource_action",
			Properties: []Property{
				{
					Name:                "type",
					Modifier:            "Required",
					Type:                "string <resource-type>@<api-version>",
					Description:         "Azure Resource Manager type.",
					CompletionNewText:   `type = "$0"`,
					ValueCandidatesFunc: resourceTypeCandidates,
				},

				{
					Name:              "resource_id",
					Modifier:          "Required",
					Type:              "string",
					Description:       "The ID of an existing azure source.",
					CompletionNewText: `resource_id = $0`,
				},

				{
					Name:                  "action",
					Modifier:              "Optional",
					Type:                  "string",
					Description:           "Specifies the name of the azure resource action.",
					CompletionNewText:     `action = "$0"`,
					GenericCandidatesFunc: actionCandidates,
				},

				{
					Name:                "method",
					Modifier:            "Optional",
					Type:                "string",
					Description:         "Specifies the Http method of the azure resource action. Defaults to `POST`",
					CompletionNewText:   `method = $0`,
					ValueCandidatesFunc: dataSourceHttpMethodCandidates,
				},

				{
					Name:                  "body",
					Modifier:              "Optional",
					Type:                  "string <JSON>",
					Description:           "A JSON object that contains the request body.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{bodyJsonencodeFuncCandidate()}),
					GenericCandidatesFunc: bodyCandidates,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},
			},
		},
		Resource{
			Name: "data.azapi_resource_list",
			Properties: []Property{
				{
					Name:                "type",
					Modifier:            "Required",
					Type:                "string <resource-type>@<api-version>",
					Description:         "Azure Resource Manager type.",
					CompletionNewText:   `type = "$0"`,
					ValueCandidatesFunc: dataSourceTypeCandidates,
				},

				{
					Name:              "parent_id",
					Modifier:          "Required",
					Type:              "string",
					Description:       "The parent resource ID to list resources under.",
					CompletionNewText: `parent_id = $0`,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},
			},
		},
		Resource{
			Name: "data.azapi_resource_id",
			Properties: []Property{
				{
					Name:                "type",
					Modifier:            "Required",
					Type:                "string <resource-type>@<api-version>",
					Description:         "Azure Resource Manager type.",
					CompletionNewText:   `type = "$0"`,
					ValueCandidatesFunc: dataSourceTypeCandidates,
				},

				{
					Name:              "name",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "Specifies the name of the azure resource.\n\nConfiguring `name` and `parent_id` is an alternative way to configure `resource_id`.",
					CompletionNewText: `name = "$0"`,
				},

				{
					Name:              "parent_id",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "The ID of the azure resource in which this resource is created.\n\nConfiguring `name` and `parent_id` is an alternative way to configure `resource_id`.",
					CompletionNewText: `parent_id = $0`,
				},

				{
					Name:              "resource_id",
					Modifier:          "Optional",
					Type:              "string",
					Description:       "The ID of an existing azure source.\n\nConfiguring `name` and `parent_id` is an alternative way to configure `resource_id`.",
					CompletionNewText: `resource_id = $0`,
				},
			},
		},
	)
}
