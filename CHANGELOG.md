## v0.4.0(unreleased)
### ENHANCEMENTS:
- Update `azapi` resource schemas.
- Update bicep types to https://github.com/Azure/bicep-types-az/commit/6bc73cdafff6e71eebb75994f5683a9dfba096df

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