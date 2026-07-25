// Command api is the entry point for the Abacus HTTP server.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jsanca/abacus/server/internal/app"
	"github.com/jsanca/abacus/server/internal/config"
	"github.com/jsanca/abacus/server/internal/observability"
	"github.com/jsanca/abacus/server/internal/transport"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	serverConfig, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	observer := &observability.Observer{}

	router := transport.NewRouter(logger, observer, transport.RouterConfig{
		AllowedOrigin:  serverConfig.AllowedOrigin,
		RequestTimeout: serverConfig.RequestTimeout,
	})

	httpServer := &http.Server{
		Addr:         serverConfig.ListenAddress,
		Handler:      router,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}

	runtime := app.NewRuntime(httpServer, observer, logger, serverConfig.ShutdownTimeout)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return runtime.Run(ctx)
}
