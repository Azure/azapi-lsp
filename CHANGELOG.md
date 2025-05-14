## v2.4.0

ENHANCEMENTS:
- `azapi` resources/data sources: Update embedded schema to the latest version.
- Support specifying the provider version used in the migration in the `terraform` block.
- `azapi_resource`, `azapi_update_resource` resources: Support `sensitive_body` field, which is used to specify the write-only properties in the request body.
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/0ce6ee9ce836e6847eaa92a6ac4ecd7ef4b89d0b

BUG FIXES:
- Fix a bug that resource group's api-version `2024-07-01` is disabled in the provider.

## v2.3.0

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/c4c1c04cee8c5362b705f1519cf0cd701ef65f6b

## v2.2.0

FEATURES:
- Support collecting telemetry data.

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/4da2e194de989ed72552add82b9a5ead5223695b

## v2.1.0

FEATURES:
- Support `refactor/rewrite` code action which can trigger the command to convert resources between `azapi` and `azurerm` providers.
- Support `aztfmigrate` command which can convert resources between `azapi` and `azurerm` providers.

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/401bed53e5495fb79f6c3357d9befb9fea158b1f

## v2.0.1
 
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/c3ff45dfffe7f229447639b5982a1e2deadc1b71

## v2.0.0-beta
FEATURES:
- Support `workspace/executeCommand` protocol which can convert ARMTemplate and resource JSON content to azapi configuration.
- `azapi_resource` resource: Support `replace_triggers_external_values` field which is used to trigger a replacement of the resource.
- `azapi` resources and data sources: Support `retry` field, which is used to specify the retry configuration.
- `azapi` resources and data sources: Support `headers` and `query_parameters` fields, which are used to specify the headers and query parameters.
- `azapi` resources and data sources: The `response_export_values` field supports JMESPath expressions.

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/7492c6d0a12a07f97b955661bf6df83d51bbb14d

## v1.15.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/5ccee7fe1b353e40ed86bfc530ee185faa43a288

## v1.14.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/37dcb1890e3661255614961f470676b486272ff2

## v1.13.0
ENHANCEMENTS:
- Support the new bicep types.
- `azapi_resource` resource: The `body` field supports dynamic schema.
- `azapi_update_resource` resource: The `body` field supports dynamic schema.
- `azapi_resource_action` resource: The `body` field supports dynamic schema.
- `azapi_resource_action` data source: The `body` field supports dynamic schema.
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/48ce933ad85391b60ee02cf83e17c9b28d31a7b1

## v1.12.0
ENHANCEMENTS:
- `azapi_resource` resource: Support `ignore_body_changes` property.
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/4abd79ba2baa05ba3c8364919b670ab43a9bf69c

BUG FIXES:
- Fix a bug that auto-completion with verified examples showed for wrong resource types.

## v1.11.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/fcfe2a66a04575f204767182fc575612c82eebc1

## v1.10.0
FEATURES:
- Support auto-completion with verified examples.

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/505b813ce50368156e3da1b86f07977b5a913be9

## v1.9.0
FEATURES:
- **New Data Source**: azapi_resource_list
- **New Data Source**: azapi_resource_id

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/1d8fec8184258cdf967b1288b156e01f7cbc8ca9

## v1.8.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/c616eb1ad4980f63c0d6b436a63701e175a62224

BUG FIXES:
- Fix a bug that field `name`'s value is not parsed correctly when validate schema.
- Fix a bug that missing required properties in array item is not reported properly.
- Fix a bug that using brackets in value is not parsed correctly.

## v1.7.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/0536b68e779fba100b9fbe32737c38d75396e2cf

## v1.6.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/da15d0376faa02a6e891dee315910535cef2a13f

## v1.5.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/b8626aecc5f47b70086580956adfcd1e267a49e6

## v1.4.0
ENHANCEMENTS:
- Support syntax which keys are wrapped by quotes.
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/a915acab5788d890aed774ec818535b44311d16d

BUG FIXES:
- Fix a bug that schema validation requires redundant `name` fields both in resource and in body.

## v1.3.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/78ec1b99699a4bf44869bd13f1b0ed7d92a99c27

## v1.2.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/019b2d62fe84508582b8c54ce3d91c2b4840e624

## v1.1.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/e641570bedc5004498d3e374adb60fdfd3521b09

## v1.0.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/a6dabb0cd645c17a1accf3ec1be4d7930e982b23

## v0.6.0
ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/5451fcd5e1bf4d8313d475d8e3dc28efc7a77e2a

## v0.5.0
FEATURES:
- **New Data Source**: azapi_resource_action
- **New Resource**: azapi_resource_action

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/813d8bbc9ecf432a2a0ff2769627592fae34369f

## v0.4.0
### ENHANCEMENTS:
- Update `azapi` resource schemas.
- Data source `azapi_resource`: Supports read only types.
- Update bicep types to https://github.com/Azure/bicep-types-az/commit/ea703e2aba0d1c024f33124ee2cd34bc0c6084b5

## v0.3.0
### ENHANCEMENTS:
- Update bicep types to https://github.com/Azure/bicep-types-az/commit/644ff521c92ce8d493f6da977af12377f32abffc

### Bugfix
- Base properties of a discriminated object are not showed in completion list
- When there's comma after last array item, HclNode parser falsely think there's one more item.

## v0.2.0
### Enhancements
- Update bicep types to https://github.com/Azure/bicep-types-az/commit/57f3ecc750648562cf170ef456ef39533872b101

## v0.1.0
### Features
- Completion of `azapi` resources and data sources
- Completion of allowed azure resource types when input `type` in `azapi` resources
- Completion of allowed azure resource properties when input `body` in `azapi` resources, limitation: it only works when use `jsonencode` function to build the JSON
- Better completion for discriminated object
- Completion for all required properties
- Show hint when hover on `azapi` resources
- Show diagnostics for properties defined inside `body`
