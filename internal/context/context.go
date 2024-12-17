package context

import (
	"context"
	"github.com/Azure/azapi-lsp/internal/filesystem"
	"github.com/Azure/azapi-lsp/internal/langserver/diagnostics"
	"github.com/Azure/azapi-lsp/internal/langserver/session"
	"github.com/Azure/azapi-lsp/internal/telemetry"
)

type contextKey struct {
	Name string
}

func (k *contextKey) String() string {
	return k.Name
}

var (
	ctxDs            = &contextKey{"document storage"}
	ctxDiagsNotifier = &contextKey{"diagnostics notifier"}
	ctxLsVersion     = &contextKey{"language server version"}
	ctxClientCaller  = &contextKey{Name: "client caller"}
	ctxTelemetry     = &contextKey{"telemetry"}
)

func missingContextErr(ctxKey *contextKey) *MissingContextErr {
	return &MissingContextErr{ctxKey}
}

func WithDocumentStorage(ctx context.Context, fs filesystem.DocumentStorage) context.Context {
	return context.WithValue(ctx, ctxDs, fs)
}

func DocumentStorage(ctx context.Context) (filesystem.DocumentStorage, error) {
	fs, ok := ctx.Value(ctxDs).(filesystem.DocumentStorage)
	if !ok {
		return nil, missingContextErr(ctxDs)
	}

	return fs, nil
}

func WithDiagnosticsNotifier(ctx context.Context, diags *diagnostics.Notifier) context.Context {
	return context.WithValue(ctx, ctxDiagsNotifier, diags)
}

func DiagnosticsNotifier(ctx context.Context) (*diagnostics.Notifier, error) {
	diags, ok := ctx.Value(ctxDiagsNotifier).(*diagnostics.Notifier)
	if !ok {
		return nil, missingContextErr(ctxDiagsNotifier)
	}

	return diags, nil
}

func WithLanguageServerVersion(ctx context.Context, version string) context.Context {
	return context.WithValue(ctx, ctxLsVersion, version)
}

func LanguageServerVersion(ctx context.Context) (string, bool) {
	version, ok := ctx.Value(ctxLsVersion).(string)
	if !ok {
		return "", false
	}
	return version, true
}

func WithClientCaller(ctx context.Context, clientCaller session.ClientCaller) context.Context {
	return context.WithValue(ctx, ctxClientCaller, clientCaller)
}

func ClientCaller(ctx context.Context) (session.ClientCaller, error) {
	clientCaller, ok := ctx.Value(ctxClientCaller).(session.ClientCaller)
	if !ok {
		return nil, missingContextErr(ctxClientCaller)
	}

	return clientCaller, nil
}

func WithClientNotifier(ctx context.Context, clientNotifier session.ClientNotifier) context.Context {
	return context.WithValue(ctx, ctxClientCaller, clientNotifier)
}

func ClientNotifier(ctx context.Context) (session.ClientNotifier, error) {
	clientNotifier, ok := ctx.Value(ctxClientCaller).(session.ClientNotifier)
	if !ok {
		return nil, missingContextErr(ctxClientCaller)
	}

	return clientNotifier, nil
}

func WithTelemetry(ctx context.Context, telemetry telemetry.Sender) context.Context {
	return context.WithValue(ctx, ctxTelemetry, telemetry)
}

func Telemetry(ctx context.Context) (telemetry.Sender, error) {
	tel, ok := ctx.Value(ctxTelemetry).(telemetry.Sender)
	if !ok {
		return nil, missingContextErr(ctxTelemetry)
	}

	return tel, nil
}
