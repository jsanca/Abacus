// Package transport provides the HTTP router, handlers, and JSON response utilities
// for the Abacus backend.
package transport

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// writeJSONResponse encodes value as JSON and writes it with the given status code.
// Encoding failures are logged; the status header is already sent at that point.
func writeJSONResponse(writer http.ResponseWriter, logger *slog.Logger, statusCode int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		logger.Error("failed to encode JSON response", "error", err)
	}
}
