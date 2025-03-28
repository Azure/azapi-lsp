package tfschema

import lsp "github.com/Azure/azapi-lsp/internal/protocol"

func init() {
	Resources = make([]Resource, 0)

	retryProperty := Property{
		Name:              "retry",
		Modifier:          "Optional",
		Type:              "block",
		Description:       "Configuration block for custom retry policy.",
		CompletionNewText: "retry = {\n\terror_message_regex = [$0]\n}",
		NestedProperties: []Property{
			{
				Name:              "error_message_regex",
				Modifier:          "Required",
				Type:              "list<string>",
				Description:       "A list of regular expressions to match against error messages. If any of the regular expressions match, the request will be retried.",
				CompletionNewText: "error_message_regex = [$0]",
			},

			{
				Name:              "interval_seconds",
				Modifier:          "Optional",
				Type:              "number",
				Description:       "The base number of seconds to wait between retries. Default is `10`.",
				CompletionNewText: "interval_seconds = $0",
			},
			{
				Name:              "max_interval_seconds",
				Modifier:          "Optional",
				Type:              "number",
				Description:       "The maximum number of seconds to wait between retries. Default is `180`.",
				CompletionNewText: "max_interval_seconds = $0",
			},
			{
				Name:              "multiplier",
				Modifier:          "Optional",
				Type:              "number",
				Description:       "The multiplier to apply to the interval between retries. Default is `1.5`.",
				CompletionNewText: "multiplier = $0",
			},
			{
				Name:              "randomization_factor",
				Modifier:          "Optional",
				Type:              "number",
				Description:       "The randomization factor to apply to the interval between retries. The formula for the randomized interval is: `RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])`. Therefore set to zero `0.0` for no randomization. Default is `0.5`.",
				CompletionNewText: "randomization_factor = $0",
			},
		},
	}

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
					ValueCandidatesFunc: typeCandidates,
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
					Type:                  "dynamic",
					Description:           "An HCL object that contains the request body used to create and update azure resource.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{dynamicPlaceholderCandidate()}),
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
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				{
					Name:              "replace_triggers_external_values",
					Modifier:          "Optional",
					Type:              "dynamic",
					Description:       "Will trigger a replace of the resource when the value changes and is not `null`.",
					CompletionNewText: `replace_triggers_external_values = $0`,
				},

				retryProperty,

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

				{
					Name:              "create_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the create request.",
					CompletionNewText: `create_headers = $0`,
				},

				{
					Name:              "update_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the update request.",
					CompletionNewText: `update_headers = $0`,
				},

				{
					Name:              "read_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the read request.",
					CompletionNewText: `read_headers = $0`,
				},

				{
					Name:              "delete_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the delete request.",
					CompletionNewText: `delete_headers = $0`,
				},

				{
					Name:              "create_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the create request.",
					CompletionNewText: `create_query_parameters = $0`,
				},

				{
					Name:              "update_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the update request.",
					CompletionNewText: `update_query_parameters = $0`,
				},

				{
					Name:              "read_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the read request.",
					CompletionNewText: `read_query_parameters = $0`,
				},

				{
					Name:              "delete_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the delete request.",
					CompletionNewText: `delete_query_parameters = $0`,
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
					ValueCandidatesFunc: typeCandidates,
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
					Type:                  "dynamic",
					Description:           "An HCL object that contains the request body used to create and update azure resource.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{dynamicPlaceholderCandidate()}),
					GenericCandidatesFunc: bodyCandidates,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				retryProperty,

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

				{
					Name:              "update_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the update request.",
					CompletionNewText: `update_headers = $0`,
				},

				{
					Name:              "read_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the read request.",
					CompletionNewText: `read_headers = $0`,
				},

				{
					Name:              "update_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the update request.",
					CompletionNewText: `update_query_parameters = $0`,
				},

				{
					Name:              "read_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the read request.",
					CompletionNewText: `read_query_parameters = $0`,
				},
			},
		},
		Resource{
			Name: "resource.azapi_data_plane_resource",
			Properties: []Property{
				{
					Name:              "type",
					Modifier:          "Required",
					Type:              "string <resource-type>@<api-version>",
					Description:       "Azure Resource Manager type.",
					CompletionNewText: `type = "$0"`,
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
					Name:                "body",
					Modifier:            "Optional",
					Type:                "dynamic",
					Description:         "An HCL object that contains the request body used to create and update azure resource.",
					CompletionNewText:   `body = $0`,
					ValueCandidatesFunc: FixedValueCandidatesFunc([]lsp.CompletionItem{dynamicPlaceholderCandidate()}),
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				retryProperty,

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

				{
					Name:              "create_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the create request.",
					CompletionNewText: `create_headers = $0`,
				},

				{
					Name:              "update_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the update request.",
					CompletionNewText: `update_headers = $0`,
				},

				{
					Name:              "read_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the read request.",
					CompletionNewText: `read_headers = $0`,
				},

				{
					Name:              "delete_headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the delete request.",
					CompletionNewText: `delete_headers = $0`,
				},

				{
					Name:              "create_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the create request.",
					CompletionNewText: `create_query_parameters = $0`,
				},

				{
					Name:              "update_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the update request.",
					CompletionNewText: `update_query_parameters = $0`,
				},

				{
					Name:              "read_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the read request.",
					CompletionNewText: `read_query_parameters = $0`,
				},

				{
					Name:              "delete_query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the delete request.",
					CompletionNewText: `delete_query_parameters = $0`,
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
					ValueCandidatesFunc: typeCandidates,
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
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				retryProperty,

				{
					Name:              "headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the request.",
					CompletionNewText: `headers = $0`,
				},

				{
					Name:              "query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the request.",
					CompletionNewText: `query_parameters = $0`,
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
					ValueCandidatesFunc: typeCandidates,
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
					Type:                  "dynamic",
					Description:           "An HCL object that contains the request body.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{dynamicPlaceholderCandidate()}),
					GenericCandidatesFunc: bodyCandidates,
				},

				{
					Name:              "locks",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of ARM resource IDs which are used to avoid create/modify/delete azapi resources at the same time.",
					CompletionNewText: `locks = [$0]`,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				{
					Name:              "sensitive_response_export_values",
					Modifier:          "Optional",
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body as sensitive.",
					CompletionNewText: `sensitive_response_export_values = [$0]`,
				},

				retryProperty,

				{
					Name:              "headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the request.",
					CompletionNewText: `headers = $0`,
				},

				{
					Name:              "query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the request.",
					CompletionNewText: `query_parameters = $0`,
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
					ValueCandidatesFunc: typeCandidates,
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
					Type:                  "dynamic",
					Description:           "An HCL object that contains the request body.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{dynamicPlaceholderCandidate()}),
					GenericCandidatesFunc: bodyCandidates,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				{
					Name:              "sensitive_response_export_values",
					Modifier:          "Optional",
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body as sensitive.",
					CompletionNewText: `sensitive_response_export_values = [$0]`,
				},

				retryProperty,

				{
					Name:              "headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the request.",
					CompletionNewText: `headers = $0`,
				},

				{
					Name:              "query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the request.",
					CompletionNewText: `query_parameters = $0`,
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
					ValueCandidatesFunc: typeCandidates,
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
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				retryProperty,

				{
					Name:              "headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the request.",
					CompletionNewText: `headers = $0`,
				},

				{
					Name:              "query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the request.",
					CompletionNewText: `query_parameters = $0`,
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
					ValueCandidatesFunc: typeCandidates,
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
		Resource{
			Name: "ephemeral.azapi_resource_action",
			Properties: []Property{
				{
					Name:                "type",
					Modifier:            "Required",
					Type:                "string <resource-type>@<api-version>",
					Description:         "Azure Resource Manager type.",
					CompletionNewText:   `type = "$0"`,
					ValueCandidatesFunc: typeCandidates,
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
					Type:                  "dynamic",
					Description:           "An HCL object that contains the request body.",
					CompletionNewText:     `body = $0`,
					ValueCandidatesFunc:   FixedValueCandidatesFunc([]lsp.CompletionItem{dynamicPlaceholderCandidate()}),
					GenericCandidatesFunc: bodyCandidates,
				},

				{
					Name:              "locks",
					Modifier:          "Optional",
					Type:              "list<string>",
					Description:       "A list of ARM resource IDs which are used to avoid create/modify/delete azapi resources at the same time.",
					CompletionNewText: `locks = [$0]`,
				},

				{
					Name:              "response_export_values",
					Modifier:          "Optional",
					Type:              "list<string> or map<string, string>",
					Description:       "The attribute can accept either a list or a map of path that needs to be exported from response body.",
					CompletionNewText: `response_export_values = [$0]`,
				},

				retryProperty,

				{
					Name:              "headers",
					Modifier:          "Optional",
					Type:              "map<string, string>",
					Description:       "A mapping of headers which should be sent with the request.",
					CompletionNewText: `headers = $0`,
				},

				{
					Name:              "query_parameters",
					Modifier:          "Optional",
					Type:              "map<string, list<string>>",
					Description:       "A mapping of query parameters which should be sent with the request.",
					CompletionNewText: `query_parameters = $0`,
				},
			},
		},
	)
}
