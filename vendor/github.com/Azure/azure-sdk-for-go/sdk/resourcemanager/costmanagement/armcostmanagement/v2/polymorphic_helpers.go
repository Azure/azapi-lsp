//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcostmanagement

import "encoding/json"

func unmarshalBenefitRecommendationPropertiesClassification(rawMsg json.RawMessage) (BenefitRecommendationPropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b BenefitRecommendationPropertiesClassification
	switch m["scope"] {
	case string(ScopeShared):
		b = &SharedScopeBenefitRecommendationProperties{}
	case string(ScopeSingle):
		b = &SingleScopeBenefitRecommendationProperties{}
	default:
		b = &BenefitRecommendationProperties{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalBenefitUtilizationSummaryClassification(rawMsg json.RawMessage) (BenefitUtilizationSummaryClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b BenefitUtilizationSummaryClassification
	switch m["kind"] {
	case string(BenefitKindIncludedQuantity):
		b = &IncludedQuantityUtilizationSummary{}
	case string(BenefitKindSavingsPlan):
		b = &SavingsPlanUtilizationSummary{}
	default:
		b = &BenefitUtilizationSummary{}
	}
	if err := json.Unmarshal(rawMsg, b); err != nil {
		return nil, err
	}
	return b, nil
}

func unmarshalBenefitUtilizationSummaryClassificationArray(rawMsg json.RawMessage) ([]BenefitUtilizationSummaryClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]BenefitUtilizationSummaryClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalBenefitUtilizationSummaryClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}