package transport

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/jsanca/abacus/server/internal/calculator"
	"github.com/jsanca/abacus/server/internal/observability"
)

// RouterConfig holds the values needed to configure the HTTP router.
type RouterConfig struct {
	// AllowedOrigin is the CORS origin permitted for local frontend development.
	AllowedOrigin string
	// RequestTimeout is the maximum duration for a single HTTP request to complete.
	RequestTimeout time.Duration
}

// NewRouter constructs and returns a configured Chi router with middleware and routes registered.
func NewRouter(logger *slog.Logger, observer *observability.Observer, registry *calculator.Registry, routerConfig RouterConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(routerConfig.RequestTimeout))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{routerConfig.AllowedOrigin},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		MaxAge:         300,
	}))
	router.Use(requestLoggingMiddleware(logger, observer))

	router.Get("/health", newHealthHandler(logger))

	router.Route("/api/v1", func(apiRouter chi.Router) {
		apiRouter.Get("/operations", newOperationsHandler(logger, registry))
		apiRouter.Post("/calculations", newCalculateHandler(logger, registry))
	})

	return router
}

// requestLoggingMiddleware returns middleware that logs each completed request
// and notifies the observer.
func requestLoggingMiddleware(logger *slog.Logger, observer *observability.Observer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()
			wrapped := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

			next.ServeHTTP(wrapped, request)

			duration := time.Since(start)
			logger.InfoContext(request.Context(), "request handled",
				"method", request.Method,
				"path", request.URL.Path,
				"status", wrapped.Status(),
				"duration_ms", duration.Milliseconds(),
				"request_id", middleware.GetReqID(request.Context()),
			)
			observer.OnRequestHandled(request.Method, request.URL.Path, wrapped.Status(), duration)
		})
	}
}
