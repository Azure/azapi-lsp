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