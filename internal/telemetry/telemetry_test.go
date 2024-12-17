package telemetry

import (
	"testing"

	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

func Test_HoverTelemetryRateLimiter(t *testing.T) {
	rateLimiter := HoverTelemetryRateLimiter{}

	inputs := []lsp.TelemetryEvent{
		{
			Name: "textDocument/hover",
		},
		{
			Name: "textDocument/completion",
		},
		{
			Name:       "textDocument/hover",
			Properties: make(map[string]interface{}),
		},
		{
			Name: "textDocument/hover",
			Properties: map[string]interface{}{
				"kind": "resource-definition",
			},
		},
	}

	expected := []bool{
		true,
		true,
		false,
		true,
	}

	for i, input := range inputs {
		if got := rateLimiter.Allow(input.Name, input.Properties); got != expected[i] {
			t.Errorf("for input no. %d, got %v, want %v", i, got, expected[i])
		}
	}
}

func Test_CompletionTelemetryRateLimiter(t *testing.T) {
	rateLimiter := CompletionTelemetryRateLimiter{}

	inputs := []lsp.TelemetryEvent{
		{
			Name: "textDocument/completion",
		},
		{
			Name: "textDocument/hover",
		},
		{
			Name:       "textDocument/completion",
			Properties: make(map[string]interface{}),
		},
		{
			Name: "textDocument/completion",
			Properties: map[string]interface{}{
				"kind": "code-sample",
			},
		},
		{
			Name: "textDocument/completion",
			Properties: map[string]interface{}{
				"kind": "template",
			},
		},
	}

	expected := []bool{
		true,
		true,
		false,
		true,
		true,
	}

	for i, input := range inputs {
		if got := rateLimiter.Allow(input.Name, input.Properties); got != expected[i] {
			t.Errorf("for input no. %d, got %v, want %v", i, got, expected[i])
		}
	}
}
