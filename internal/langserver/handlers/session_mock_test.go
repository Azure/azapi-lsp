package handlers

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/Azure/azapi-lsp/internal/filesystem"
	"github.com/Azure/azapi-lsp/internal/langserver/session"
	"github.com/creachadair/jrpc2/handler"
)

type MockSessionInput struct {
	Filesystem         filesystem.Filesystem
	AdditionalHandlers map[string]handler.Func
}

type mockSession struct {
	mockInput *MockSessionInput

	stopFunc     func()
	stopCalled   bool
	stopCalledMu *sync.RWMutex
}

func (ms *mockSession) new(srvCtx context.Context) session.Session {
	sessCtx, stopSession := context.WithCancel(srvCtx)
	ms.stopFunc = stopSession

	var fs filesystem.Filesystem
	fs = filesystem.NewFilesystem()
	var handlers map[string]handler.Func
	if ms.mockInput != nil {
		if ms.mockInput.Filesystem != nil {
			fs = ms.mockInput.Filesystem
		}
		handlers = ms.mockInput.AdditionalHandlers
	}

	svc := &service{
		logger:             testLogger(),
		srvCtx:             srvCtx,
		sessCtx:            sessCtx,
		stopSession:        ms.stop,
		fs:                 fs,
		additionalHandlers: handlers,
	}

	return svc
}

func testLogger() *log.Logger {
	if testing.Verbose() {
		return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}

	return log.New(io.Discard, "", 0)
}

func (ms *mockSession) stop() {
	ms.stopCalledMu.Lock()
	defer ms.stopCalledMu.Unlock()

	ms.stopFunc()
	ms.stopCalled = true
}

func (ms *mockSession) StopFuncCalled() bool {
	ms.stopCalledMu.RLock()
	defer ms.stopCalledMu.RUnlock()

	return ms.stopCalled
}

func newMockSession(input *MockSessionInput) *mockSession {
	return &mockSession{
		mockInput:    input,
		stopCalledMu: &sync.RWMutex{},
	}
}

func NewMockSession(input *MockSessionInput) session.SessionFactory {
	return newMockSession(input).new
}
