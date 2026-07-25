package transport_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jsanca/abacus/server/internal/calculator"
	"github.com/jsanca/abacus/server/internal/observability"
	"github.com/jsanca/abacus/server/internal/transport"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
	observer := &observability.Observer{}
	registry, err := calculator.NewDefaultRegistry()
	if err != nil {
		t.Fatalf("setup registry: %v", err)
	}
	return transport.NewRouter(logger, observer, registry, transport.RouterConfig{
		AllowedOrigin:  "http://localhost:5173",
		RequestTimeout: 5 * time.Second,
	})
}

// --- /health ---

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

// --- GET /api/v1/operations ---

func TestOperationsEndpoint_returnsManifest(t *testing.T) {
	router := newTestRouter(t)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/operations", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if ct := recorder.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type %q, got %q", "application/json", ct)
	}

	var body map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode response body: %v", err)
	}
	if body["defaultOperationId"] != "addition" {
		t.Errorf("expected defaultOperationId %q, got %v", "addition", body["defaultOperationId"])
	}
	operations, ok := body["operations"].([]any)
	if !ok {
		t.Fatal("operations field missing or not an array")
	}
	if len(operations) != 7 {
		t.Errorf("expected 7 operations, got %d", len(operations))
	}
	if body["version"] == "" {
		t.Error("expected non-empty version field")
	}
}

// --- POST /api/v1/calculations ---

func postCalculation(t *testing.T, router http.Handler, body any) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculations", bytes.NewReader(data))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestCalculateEndpoint_binaryOperation(t *testing.T) {
	router := newTestRouter(t)
	recorder := postCalculation(t, router, map[string]any{
		"operationId": "addition",
		"operands":    []float64{3, 4},
	})

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["result"] != float64(7) {
		t.Errorf("expected result 7, got %v", body["result"])
	}
	if body["expression"] != "3 + 4" {
		t.Errorf("expected expression %q, got %q", "3 + 4", body["expression"])
	}
}

func TestCalculateEndpoint_unaryOperation(t *testing.T) {
	router := newTestRouter(t)
	recorder := postCalculation(t, router, map[string]any{
		"operationId": "square-root",
		"operands":    []float64{9},
	})

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["result"] != float64(3) {
		t.Errorf("expected result 3, got %v", body["result"])
	}
	if body["expression"] != "√9" {
		t.Errorf("expected expression %q, got %q", "√9", body["expression"])
	}
}

func TestCalculateEndpoint_unknownOperation(t *testing.T) {
	router := newTestRouter(t)
	recorder := postCalculation(t, router, map[string]any{
		"operationId": "nonexistent",
		"operands":    []float64{1, 2},
	})

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["code"] != "unknown_operation" {
		t.Errorf("expected code %q, got %q", "unknown_operation", body["code"])
	}
}

func TestCalculateEndpoint_validationFailure_divisionByZero(t *testing.T) {
	router := newTestRouter(t)
	recorder := postCalculation(t, router, map[string]any{
		"operationId": "division",
		"operands":    []float64{10, 0},
	})

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, recorder.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["code"] != "validation_failed" {
		t.Errorf("expected code %q, got %q", "validation_failed", body["code"])
	}
	if body["message"] == "" {
		t.Error("expected non-empty validation message")
	}
}

func TestCalculateEndpoint_expressionFormats(t *testing.T) {
	router := newTestRouter(t)
	tests := []struct {
		operationID string
		operands    []float64
		expression  string
	}{
		{"addition", []float64{10, 5}, "10 + 5"},
		{"subtraction", []float64{5, 10}, "5 − 10"},
		{"multiplication", []float64{2.5, 4}, "2.5 × 4"},
		{"division", []float64{144, 12}, "144 ÷ 12"},
		{"exponentiation", []float64{2, 8}, "2 xʸ 8"},
		{"square-root", []float64{144}, "√144"},
		{"percentage", []float64{200, 15}, "15% of 200"},
	}
	for _, test := range tests {
		t.Run(test.operationID, func(t *testing.T) {
			recorder := postCalculation(t, router, map[string]any{
				"operationId": test.operationID,
				"operands":    test.operands,
			})
			if recorder.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
			}
			var body map[string]any
			if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if body["expression"] != test.expression {
				t.Errorf("expected expression %q, got %q", test.expression, body["expression"])
			}
		})
	}
}

func TestCalculateEndpoint_wrongOperandCountReturnsError(t *testing.T) {
	router := newTestRouter(t)
	tests := []struct {
		operationID string
		operands    []float64
	}{
		{"addition", []float64{1}},
		{"square-root", []float64{9, 4}},
		{"percentage", []float64{200}},
	}
	for _, test := range tests {
		t.Run(test.operationID, func(t *testing.T) {
			recorder := postCalculation(t, router, map[string]any{
				"operationId": test.operationID,
				"operands":    test.operands,
			})
			if recorder.Code != http.StatusInternalServerError {
				t.Errorf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
			}
			var body map[string]string
			if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if body["code"] != "calculation_failed" {
				t.Errorf("expected code %q, got %q", "calculation_failed", body["code"])
			}
		})
	}
}

func TestCalculateEndpoint_malformedJSON(t *testing.T) {
	router := newTestRouter(t)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculations",
		bytes.NewBufferString("not-valid-json"))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["code"] != "malformed_request" {
		t.Errorf("expected code %q, got %q", "malformed_request", body["code"])
	}
}
