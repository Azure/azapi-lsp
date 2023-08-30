## v1.9.0 (unreleased)
FEATURES:
- **New Data Source**: azapi_resource_list
- **New Data Source**: azapi_resource_id

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
