package app_test

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jsanca/abacus/server/internal/app"
	"github.com/jsanca/abacus/server/internal/observability"
)

func TestRuntime_contextCancellationTriggersGracefulShutdown(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	observer := &observability.Observer{}

	httpServer := &http.Server{
		Addr:    ":0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	}

	runtime := app.NewRuntime(httpServer, observer, logger, 5*time.Second)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() {
		done <- runtime.Run(ctx)
	}()

	// Allow the server goroutine time to call ListenAndServe before signalling shutdown.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("expected nil error on graceful shutdown, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return within timeout after context cancellation")
	}
}

func TestRuntime_immediateContextCancellationExitsCleanly(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	observer := &observability.Observer{}

	httpServer := &http.Server{
		Addr:    ":0",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	}

	runtime := app.NewRuntime(httpServer, observer, logger, 5*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before Run is called

	done := make(chan error, 1)
	go func() {
		done <- runtime.Run(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("expected nil error when context already cancelled, got: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return within timeout")
	}
}
