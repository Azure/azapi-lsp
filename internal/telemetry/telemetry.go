package telemetry

import (
	"context"
	"fmt"
	"time"

	lsp "github.com/Azure/azapi-lsp/internal/protocol"
)

var _ Sender = &Telemetry{}

type Telemetry struct {
	version  int
	notifier Notifier
}

type Notifier interface {
	Notify(ctx context.Context, method string, params interface{}) error
}

func NewSender(version int, notifier Notifier) (*Telemetry, error) {
	if version != lsp.TelemetryFormatVersion {
		return nil, fmt.Errorf("unsupported telemetry format version: %d", version)
	}

	return &Telemetry{
		version:  version,
		notifier: notifier,
	}, nil
}

func (t *Telemetry) SendEvent(ctx context.Context, name string, properties map[string]interface{}) {
	for _, limiter := range telemetryRateLimiters {
		if !limiter.Allow(name, properties) {
			return
		}
	}

	_ = t.notifier.Notify(ctx, "telemetry/event", lsp.TelemetryEvent{
		Version:    t.version,
		Name:       name,
		Properties: properties,
	})
}

var telemetryRateLimiters []TelemetryRateLimiter

func init() {
	telemetryRateLimiters = make([]TelemetryRateLimiter, 0)
	telemetryRateLimiters = append(telemetryRateLimiters, &HoverTelemetryRateLimiter{})
	telemetryRateLimiters = append(telemetryRateLimiters, &CompletionTelemetryRateLimiter{})
}

type TelemetryRateLimiter interface {
	Allow(name string, properties map[string]interface{}) bool
}

var _ TelemetryRateLimiter = &HoverTelemetryRateLimiter{}

type HoverTelemetryRateLimiter struct {
	lastSentTimeStamp int64
}

func (h *HoverTelemetryRateLimiter) Allow(name string, properties map[string]interface{}) bool {
	if name != "textDocument/hover" {
		return true
	}
	// always allow the resource-definition hover
	if properties != nil && properties["kind"] != nil {
		if kind, ok := properties["kind"].(string); ok && kind == "resource-definition" {
			return true
		}
	}

	timeStampNow := time.Now().Unix()
	if h.lastSentTimeStamp == 0 {
		h.lastSentTimeStamp = timeStampNow
		return true
	}

	// minimum interval between two hover telemetry events is 1 hour
	if timeStampNow-h.lastSentTimeStamp < 3600 {
		return false
	}
	h.lastSentTimeStamp = timeStampNow
	return true
}

var _ TelemetryRateLimiter = &CompletionTelemetryRateLimiter{}

type CompletionTelemetryRateLimiter struct {
	lastSentTimeStamp int64
}

func (c *CompletionTelemetryRateLimiter) Allow(name string, properties map[string]interface{}) bool {
	if name != "textDocument/completion" {
		return true
	}

	// always allow code-sample and template completions
	if properties != nil && properties["kind"] != nil {
		if kind, ok := properties["kind"].(string); ok {
			if kind == "code-sample" || kind == "template" {
				return true
			}
		}
	}

	timeStampNow := time.Now().Unix()
	if c.lastSentTimeStamp == 0 {
		c.lastSentTimeStamp = timeStampNow
		return true
	}

	// minimum interval between two completion telemetry events is 1 hour
	if timeStampNow-c.lastSentTimeStamp < 3600 {
		return false
	}
	c.lastSentTimeStamp = timeStampNow
	return true
}
