package transport

import (
	"log/slog"
	"net/http"
)

type healthBody struct {
	Status string `json:"status"`
}

// newHealthHandler returns an HTTP handler for GET /health.
func newHealthHandler(logger *slog.Logger) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeJSONResponse(writer, logger, http.StatusOK, healthBody{Status: "ok"})
	}
}
