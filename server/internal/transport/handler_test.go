package transport_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jsanca/abacus/server/internal/observability"
	"github.com/jsanca/abacus/server/internal/transport"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	observer := &observability.Observer{}
	return transport.NewRouter(logger, observer, transport.RouterConfig{
		AllowedOrigin:  "http://localhost:5173",
		RequestTimeout: 5 * time.Second,
	})
}

func TestHealthEndpoint_returnsOKWithJSONBody(t *testing.T) {
	router := newTestRouter(t)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", contentType)
	}

	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status field %q, got %q", "ok", body["status"])
	}
}

func TestHealthEndpoint_methodNotAllowed(t *testing.T) {
	router := newTestRouter(t)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			request := httptest.NewRequest(method, "/health", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d for %s, got %d", http.StatusMethodNotAllowed, method, recorder.Code)
			}
		})
	}
}

func TestUnknownRoute_returnsNotFound(t *testing.T) {
	router := newTestRouter(t)

	request := httptest.NewRequest(http.MethodGet, "/not-a-route", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}
