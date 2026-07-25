// Package app owns the Abacus HTTP server lifecycle: startup, signal-driven
// graceful shutdown, and drain.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jsanca/abacus/server/internal/observability"
)

// Runtime owns the HTTP server and is responsible for its full lifecycle:
// startup, signal-driven graceful shutdown, and in-flight request drain.
type Runtime struct {
	httpServer      *http.Server
	observer        *observability.Observer
	logger          *slog.Logger
	shutdownTimeout time.Duration
}

// NewRuntime constructs a Runtime that owns the given HTTP server.
func NewRuntime(
	httpServer *http.Server,
	observer *observability.Observer,
	logger *slog.Logger,
	shutdownTimeout time.Duration,
) *Runtime {
	return &Runtime{
		httpServer:      httpServer,
		observer:        observer,
		logger:          logger,
		shutdownTimeout: shutdownTimeout,
	}
}

// Run starts the HTTP server and blocks until ctx is cancelled or an unexpected
// server failure occurs. On cancellation, Run initiates a bounded graceful shutdown
// and waits for in-flight requests to drain before returning.
func (runtime *Runtime) Run(ctx context.Context) error {
	serverErrors := make(chan error, 1)

	go func() {
		runtime.logger.Info("server starting", "address", runtime.httpServer.Addr)
		runtime.observer.OnServerStart(runtime.httpServer.Addr)
		serverErrors <- runtime.httpServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("unexpected server failure: %w", err)

	case <-ctx.Done():
		runtime.logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), runtime.shutdownTimeout)
	defer cancel()

	if err := runtime.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown did not complete within %v: %w", runtime.shutdownTimeout, err)
	}

	runtime.observer.OnServerStop()
	runtime.logger.Info("graceful shutdown complete")
	return nil
}
