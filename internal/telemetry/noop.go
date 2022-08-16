package telemetry

import (
	"context"
	"io"
	"log"
)

type NoopSender struct {
	Logger *log.Logger
}

func (t *NoopSender) log() *log.Logger { //nolint
	if t.Logger != nil {
		return t.Logger
	}
	return log.New(io.Discard, "", 0)
}

func (t *NoopSender) SendEvent(ctx context.Context, name string, properties map[string]interface{}) {}
