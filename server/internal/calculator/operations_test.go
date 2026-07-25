package calculator_test

import (
	"errors"
	"testing"

	"github.com/jsanca/abacus/server/internal/calculator"
)

func TestOperation_execution_correctResults(t *testing.T) {
	tests := []struct {
		operationID string
		operands    []float64
		want        float64
	}{
		{operationID: "addition", operands: []float64{10, 5}, want: 15},
		{operationID: "subtraction", operands: []float64{5, 10}, want: -5},
		{operationID: "multiplication", operands: []float64{2.5, 4}, want: 10},
		{operationID: "division", operands: []float64{144, 12}, want: 12},
		{operationID: "exponentiation", operands: []float64{2, 8}, want: 256},
		{operationID: "square-root", operands: []float64{144}, want: 12},
		{operationID: "percentage", operands: []float64{200, 15}, want: 30},
	}

	registry, err := calculator.NewDefaultRegistry()
	if err != nil {
		t.Fatalf("registry setup failed: %v", err)
	}

	for _, testCase := range tests {
		t.Run(testCase.operationID, func(t *testing.T) {
			operation, found := registry.Find(testCase.operationID)
			if !found {
				t.Fatalf("operation %q not found", testCase.operationID)
			}
			result, err := operation.Execute(testCase.operands)
			if err != nil {
				t.Fatalf("Execute returned error: %v", err)
			}
			if result != testCase.want {
				t.Errorf("expected %v, got %v", testCase.want, result)
			}
		})
	}
}

func TestOperation_wrongArity_returnsErrWrongOperandCount(t *testing.T) {
	tests := []struct {
		operationID string
		badOperands []float64
	}{
		{operationID: "addition", badOperands: []float64{1}},
		{operationID: "subtraction", badOperands: []float64{1, 2, 3}},
		{operationID: "multiplication", badOperands: []float64{}},
		{operationID: "division", badOperands: []float64{10}},
		{operationID: "exponentiation", badOperands: []float64{2}},
		{operationID: "square-root", badOperands: []float64{4, 9}},
		{operationID: "percentage", badOperands: []float64{100}},
	}

	registry, _ := calculator.NewDefaultRegistry()

	for _, testCase := range tests {
		t.Run(testCase.operationID, func(t *testing.T) {
			t.Parallel()
			operation, _ := registry.Find(testCase.operationID)

			// Must not panic.
			result, err := operation.Execute(testCase.badOperands)
			if err == nil {
				t.Errorf("expected error for wrong arity, got result %v", result)
			}
			if !errors.Is(err, calculator.ErrWrongOperandCount) {
				t.Errorf("expected errors.Is(err, ErrWrongOperandCount), got: %v", err)
			}
		})
	}
}
