## v0.5.0 (Unreleased)
FEATURES:
- **New Data Source**: azapi_operation
- **New Resource**: azapi_operation

ENHANCEMENTS:
- Update bicep types to https://github.com/ms-henglu/bicep-types-az/commit/08fe00486aec8ee0e35a7776f81da090c32476de

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