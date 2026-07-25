package calculator_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/jsanca/abacus/server/internal/calculator"
)

// makeUnaryOp constructs a minimal unary operation with the given validation expression
// for use in tests that need to bypass registry construction validation.
func makeUnaryOp(expr calculator.Expression) calculator.Operation {
	return calculator.Operation{
		Definition: calculator.OperationDefinition{
			ID:    "test-op",
			Arity: calculator.Unary,
			Operands: []calculator.OperandDefinition{
				{ID: calculator.OperandFirst, Label: "Number", Placeholder: "0"},
			},
			Validations: []calculator.ValidationDefinition{
				{ID: "test-rule", Message: "test constraint", Expression: expr},
			},
		},
		Execute: func(operands []float64) (float64, error) { return 0, nil },
	}
}

// makeBinaryOp constructs a minimal binary operation with the given validation expression.
func makeBinaryOp(expr calculator.Expression) calculator.Operation {
	return calculator.Operation{
		Definition: calculator.OperationDefinition{
			ID:    "test-op",
			Arity: calculator.Binary,
			Operands: []calculator.OperandDefinition{
				{ID: calculator.OperandFirst, Label: "First", Placeholder: "0"},
				{ID: calculator.OperandSecond, Label: "Second", Placeholder: "0"},
			},
			Validations: []calculator.ValidationDefinition{
				{ID: "test-rule", Message: "test constraint", Expression: expr},
			},
		},
		Execute: func(operands []float64) (float64, error) { return 0, nil },
	}
}

// --- Operation.Validate: domain operations ---

func TestOperation_Validate_division(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	op, _ := registry.Find("division")

	if violation, err := op.Validate([]float64{12, 3}); violation != nil || err != nil {
		t.Errorf("12 / 3: expected no violation, got violation=%v err=%v", violation, err)
	}

	violation, err := op.Validate([]float64{12, 0})
	if err != nil {
		t.Fatalf("12 / 0: unexpected error: %v", err)
	}
	if violation == nil {
		t.Fatal("12 / 0: expected violation, got nil")
	}
	if violation.ID != "division-by-zero" {
		t.Errorf("expected violation ID %q, got %q", "division-by-zero", violation.ID)
	}
}

func TestOperation_Validate_squareRoot(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	op, _ := registry.Find("square-root")

	if violation, err := op.Validate([]float64{9}); violation != nil || err != nil {
		t.Errorf("sqrt(9): expected no violation, got violation=%v err=%v", violation, err)
	}

	violation, err := op.Validate([]float64{-1})
	if err != nil {
		t.Fatalf("sqrt(-1): unexpected error: %v", err)
	}
	if violation == nil {
		t.Fatal("sqrt(-1): expected violation, got nil")
	}
	if violation.ID != "square-root-negative" {
		t.Errorf("expected violation ID %q, got %q", "square-root-negative", violation.ID)
	}
}

func TestOperation_Validate_percentage(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	op, _ := registry.Find("percentage")

	tests := []struct {
		name     string
		operands []float64
		wantFail bool
	}{
		{"50% passes", []float64{200, 50}, false},
		{"0% passes", []float64{200, 0}, false},
		{"100% passes", []float64{200, 100}, false},
		{"101% fails", []float64{200, 101}, true},
		{"-5% fails", []float64{200, -5}, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			violation, err := op.Validate(testCase.operands)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if testCase.wantFail && violation == nil {
				t.Error("expected validation failure, got nil")
			}
			if !testCase.wantFail && violation != nil {
				t.Errorf("expected no violation, got: %v", violation.Message)
			}
		})
	}
}

func TestOperation_Validate_unconstrained(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	ids := []string{"addition", "subtraction", "multiplication", "exponentiation"}

	for _, operationID := range ids {
		t.Run(operationID, func(t *testing.T) {
			op, _ := registry.Find(operationID)
			if violation, err := op.Validate([]float64{1, 2}); violation != nil || err != nil {
				t.Errorf("%s: expected no violation, got violation=%v err=%v", operationID, violation, err)
			}
		})
	}
}

func TestOperation_Validate_firstViolationInDeclarationOrder(t *testing.T) {
	op := calculator.Operation{
		Definition: calculator.OperationDefinition{
			ID:    "multi-rule",
			Arity: calculator.Binary,
			Operands: []calculator.OperandDefinition{
				{ID: calculator.OperandFirst, Label: "First", Placeholder: "0"},
				{ID: calculator.OperandSecond, Label: "Second", Placeholder: "0"},
			},
			Validations: []calculator.ValidationDefinition{
				{
					ID:      "first-rule",
					Message: "first fails",
					Expression: calculator.ComparisonExpression{
						Operand:  calculator.FirstOperand,
						Operator: calculator.GreaterThan,
						Value:    10,
					},
				},
				{
					ID:      "second-rule",
					Message: "second fails",
					Expression: calculator.ComparisonExpression{
						Operand:  calculator.SecondOperand,
						Operator: calculator.GreaterThan,
						Value:    10,
					},
				},
			},
		},
		Execute: func(operands []float64) (float64, error) { return 0, nil },
	}

	// Both rules fail; only the first should be returned.
	violation, err := op.Validate([]float64{1, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if violation == nil {
		t.Fatal("expected violation, got nil")
	}
	if violation.ID != "first-rule" {
		t.Errorf("expected first-rule violation, got %q", violation.ID)
	}
}

// --- Operation.Validate: evaluation error cases ---

func TestOperation_Validate_missingFirstOperand(t *testing.T) {
	op := makeUnaryOp(calculator.ComparisonExpression{
		Operand:  calculator.FirstOperand,
		Operator: calculator.Equal,
		Value:    0,
	})

	_, err := op.Validate([]float64{})
	if err == nil {
		t.Fatal("expected error for missing first operand, got nil")
	}
	if !errors.Is(err, calculator.ErrOperandOutOfBounds) {
		t.Errorf("expected ErrOperandOutOfBounds, got: %v", err)
	}
}

func TestOperation_Validate_missingSecondOperand(t *testing.T) {
	op := makeBinaryOp(calculator.ComparisonExpression{
		Operand:  calculator.SecondOperand,
		Operator: calculator.Equal,
		Value:    0,
	})

	_, err := op.Validate([]float64{1})
	if err == nil {
		t.Fatal("expected error for missing second operand, got nil")
	}
	if !errors.Is(err, calculator.ErrOperandOutOfBounds) {
		t.Errorf("expected ErrOperandOutOfBounds, got: %v", err)
	}
}

func TestOperation_Validate_unsupportedOperandReference(t *testing.T) {
	op := makeBinaryOp(calculator.ComparisonExpression{
		Operand:  calculator.OperandReference("bogus"),
		Operator: calculator.Equal,
		Value:    0,
	})

	_, err := op.Validate([]float64{1, 2})
	if err == nil {
		t.Fatal("expected error for unsupported operand reference, got nil")
	}
	if !errors.Is(err, calculator.ErrUnsupportedOperandReference) {
		t.Errorf("expected ErrUnsupportedOperandReference, got: %v", err)
	}
}

func TestOperation_Validate_unsupportedComparisonOperator(t *testing.T) {
	op := makeUnaryOp(calculator.ComparisonExpression{
		Operand:  calculator.FirstOperand,
		Operator: calculator.ComparisonOperator("bogus"),
		Value:    0,
	})

	_, err := op.Validate([]float64{1})
	if err == nil {
		t.Fatal("expected error for unsupported comparison operator, got nil")
	}
	if !errors.Is(err, calculator.ErrUnsupportedComparisonOperator) {
		t.Errorf("expected ErrUnsupportedComparisonOperator, got: %v", err)
	}
}

func TestOperation_Validate_nilExpression_doesNotPanic(t *testing.T) {
	op := calculator.Operation{
		Definition: calculator.OperationDefinition{
			Validations: []calculator.ValidationDefinition{
				{ID: "rule", Message: "msg", Expression: nil},
			},
		},
		Execute: func(operands []float64) (float64, error) { return 0, nil },
	}

	_, err := op.Validate([]float64{})
	if err == nil {
		t.Error("expected error for nil expression, got nil")
	}
}

// --- AllOf: evaluation behavior ---

func TestOperation_Validate_allOf_allChildrenPass(t *testing.T) {
	op := makeBinaryOp(calculator.AllOfExpression{
		Expressions: []calculator.Expression{
			calculator.ComparisonExpression{Operand: calculator.SecondOperand, Operator: calculator.GreaterThanOrEqual, Value: 0},
			calculator.ComparisonExpression{Operand: calculator.SecondOperand, Operator: calculator.LessThanOrEqual, Value: 100},
		},
	})

	violation, err := op.Validate([]float64{0, 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if violation != nil {
		t.Errorf("expected no violation, got: %v", violation.Message)
	}
}

func TestOperation_Validate_allOf_firstChildFailsShortCircuits(t *testing.T) {
	// Child 1 always fails; child 2 would return ErrOperandOutOfBounds on a unary call.
	// Short-circuit means child 2 is never evaluated, so we get a violation, not an error.
	op := makeUnaryOp(calculator.AllOfExpression{
		Expressions: []calculator.Expression{
			calculator.ComparisonExpression{
				Operand:  calculator.FirstOperand,
				Operator: calculator.GreaterThan,
				Value:    1000,
			},
			calculator.ComparisonExpression{
				Operand:  calculator.SecondOperand, // would error: no second operand
				Operator: calculator.Equal,
				Value:    0,
			},
		},
	})

	violation, err := op.Validate([]float64{50})
	if err != nil {
		t.Errorf("expected no error (short-circuit before second child), got: %v", err)
	}
	if violation == nil {
		t.Error("expected violation, got nil")
	}
}

func TestOperation_Validate_allOf_childErrorPropagates(t *testing.T) {
	op := makeBinaryOp(calculator.AllOfExpression{
		Expressions: []calculator.Expression{
			calculator.ComparisonExpression{
				Operand:  calculator.SecondOperand,
				Operator: calculator.Equal,
				Value:    0,
			},
		},
	})

	// Calling with only 1 operand causes the child to error.
	_, err := op.Validate([]float64{1})
	if err == nil {
		t.Fatal("expected error from AllOf child propagation, got nil")
	}
	if !errors.Is(err, calculator.ErrOperandOutOfBounds) {
		t.Errorf("expected ErrOperandOutOfBounds, got: %v", err)
	}
}

func TestOperation_Validate_allOf_emptyIsVacuouslyTrue(t *testing.T) {
	op := makeUnaryOp(calculator.AllOfExpression{Expressions: []calculator.Expression{}})

	violation, err := op.Validate([]float64{1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if violation != nil {
		t.Errorf("expected no violation for empty AllOf, got: %v", violation.Message)
	}
}

// --- Manifest: validation serialization ---

func TestManifest_validationsSerializedCorrectly(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	manifest := registry.Manifest()

	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	ops := raw["operations"].([]any)
	byID := make(map[string]map[string]any, len(ops))
	for _, op := range ops {
		opMap := op.(map[string]any)
		byID[opMap["id"].(string)] = opMap
	}

	// Addition has no validations.
	additionValidations := byID["addition"]["validations"].([]any)
	if len(additionValidations) != 0 {
		t.Errorf("addition: expected 0 validations, got %d", len(additionValidations))
	}

	// Division has one comparison validation.
	divisionValidations := byID["division"]["validations"].([]any)
	if len(divisionValidations) != 1 {
		t.Fatalf("division: expected 1 validation, got %d", len(divisionValidations))
	}
	divVal := divisionValidations[0].(map[string]any)
	if divVal["message"] != "The divisor must not be zero." {
		t.Errorf("division validation message mismatch: %v", divVal["message"])
	}
	divExpr := divVal["expression"].(map[string]any)
	if divExpr["kind"] != "comparison" {
		t.Errorf("division expression kind: expected %q, got %q", "comparison", divExpr["kind"])
	}
	if divExpr["operand"] != "second" {
		t.Errorf("division expression operand: expected %q, got %q", "second", divExpr["operand"])
	}
	if divExpr["operator"] != "notEqual" {
		t.Errorf("division expression operator: expected %q, got %q", "notEqual", divExpr["operator"])
	}

	// Percentage has one allOf validation.
	percentageValidations := byID["percentage"]["validations"].([]any)
	if len(percentageValidations) != 1 {
		t.Fatalf("percentage: expected 1 validation, got %d", len(percentageValidations))
	}
	pctVal := percentageValidations[0].(map[string]any)
	pctExpr := pctVal["expression"].(map[string]any)
	if pctExpr["kind"] != "allOf" {
		t.Errorf("percentage expression kind: expected %q, got %q", "allOf", pctExpr["kind"])
	}
	pctInner := pctExpr["expressions"].([]any)
	if len(pctInner) != 2 {
		t.Errorf("percentage allOf: expected 2 inner expressions, got %d", len(pctInner))
	}

	// Square root has one comparison validation.
	sqrtValidations := byID["square-root"]["validations"].([]any)
	if len(sqrtValidations) != 1 {
		t.Fatalf("square-root: expected 1 validation, got %d", len(sqrtValidations))
	}
}

func TestManifest_validationsDeterministicOrdering(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()

	first := registry.Manifest()
	second := registry.Manifest()

	firstData, _ := json.Marshal(first)
	secondData, _ := json.Marshal(second)

	if string(firstData) != string(secondData) {
		t.Error("manifest JSON is not deterministic across calls")
	}
}

func TestManifest_executorNotExposedInValidation(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	manifest := registry.Manifest()

	data, _ := json.Marshal(manifest)
	raw := string(data)

	if containsFunc(raw) {
		t.Error("serialized manifest appears to contain a function reference")
	}
}

func containsFunc(s string) bool {
	return len(s) > 4 && findSubstring(s, "func")
}

func findSubstring(s, sub string) bool {
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
