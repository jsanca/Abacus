package transport

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jsanca/abacus/server/internal/calculator"
)

// errorBody is the standard error response shape for all API error conditions.
type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type calculationRequest struct {
	OperationID string    `json:"operationId"`
	Operands    []float64 `json:"operands"`
}

type calculationResponse struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

// newOperationsHandler returns a handler for GET /api/v1/operations.
// It responds with the registry manifest serialized as JSON.
func newOperationsHandler(logger *slog.Logger, registry *calculator.Registry) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writeJSONResponse(writer, logger, http.StatusOK, registry.Manifest())
	}
}

// newCalculateHandler returns a handler for POST /api/v1/calculations.
// Application flow: lookup → validate → execute → build response.
func newCalculateHandler(logger *slog.Logger, registry *calculator.Registry) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req calculationRequest
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			writeJSONResponse(writer, logger, http.StatusBadRequest, errorBody{
				Code:    "malformed_request",
				Message: "Request body is not valid JSON.",
			})
			return
		}

		operation, found := registry.Find(req.OperationID)
		if !found {
			writeJSONResponse(writer, logger, http.StatusNotFound, errorBody{
				Code:    "unknown_operation",
				Message: "The requested operation is not available.",
			})
			return
		}

		violation, err := operation.Validate(req.Operands)
		if err != nil {
			logger.ErrorContext(request.Context(), "validation error",
				"operation", req.OperationID, "error", err)
			writeJSONResponse(writer, logger, http.StatusInternalServerError, errorBody{
				Code:    "calculation_failed",
				Message: "An internal error occurred while validating the request.",
			})
			return
		}
		if violation != nil {
			writeJSONResponse(writer, logger, http.StatusUnprocessableEntity, errorBody{
				Code:    "validation_failed",
				Message: violation.Message,
			})
			return
		}

		result, err := operation.Execute(req.Operands)
		if err != nil {
			logger.ErrorContext(request.Context(), "execution error",
				"operation", req.OperationID, "error", err)
			writeJSONResponse(writer, logger, http.StatusInternalServerError, errorBody{
				Code:    "calculation_failed",
				Message: "An internal error occurred while executing the calculation.",
			})
			return
		}

		expressionStr, err := formatExpression(operation, req.Operands)
		if err != nil {
			logger.ErrorContext(request.Context(), "formatting error",
				"operation", req.OperationID, "error", err)
			writeJSONResponse(writer, logger, http.StatusInternalServerError, errorBody{
				Code:    "calculation_failed",
				Message: "An internal error occurred while formatting the result.",
			})
			return
		}

		writeJSONResponse(writer, logger, http.StatusOK, calculationResponse{
			Expression: expressionStr,
			Result:     result,
		})
	}
}

// formatExpression formats the expression string for a resolved calculation.
// If the operation provides a FormatExpression function it is called; otherwise
// defaultFormatExpression applies standard unary prefix or binary infix notation.
func formatExpression(operation calculator.Operation, operands []float64) (string, error) {
	if operation.FormatExpression != nil {
		return operation.FormatExpression(operation.Definition, operands)
	}
	return defaultFormatExpression(operation.Definition, operands)
}

// defaultFormatExpression produces standard notation: prefix for unary (e.g. √9),
// infix for binary (e.g. 144 ÷ 12). Returns a descriptive error on wrong operand count.
func defaultFormatExpression(definition calculator.OperationDefinition, operands []float64) (string, error) {
	if definition.Arity == calculator.Unary {
		if len(operands) < 1 {
			return "", fmt.Errorf("format %q expression: need 1 operand, got %d", definition.ID, len(operands))
		}
		return definition.Symbol + formatFloat(operands[0]), nil
	}
	if len(operands) < 2 {
		return "", fmt.Errorf("format %q expression: need 2 operands, got %d", definition.ID, len(operands))
	}
	return formatFloat(operands[0]) + " " + definition.Symbol + " " + formatFloat(operands[1]), nil
}

// formatFloat returns the shortest decimal string representation of value.
// Integer-valued floats are rendered without a decimal point (12.0 → "12").
func formatFloat(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
