package calculator_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jsanca/abacus/server/internal/calculator"
)

// --- Construction tests ---

func TestNewDefaultRegistry_succeeds(t *testing.T) {
	registry, err := calculator.NewDefaultRegistry()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestNewRegistry_invalidInputs(t *testing.T) {
	validExecute := func(operands []float64) (float64, error) { return 0, nil }

	validOp := calculator.Operation{
		Definition: calculator.OperationDefinition{
			ID:       "addition",
			Name:     "Addition",
			Symbol:   "+",
			Shortcut: "+",
			Arity:    calculator.Binary,
			Operands: []calculator.OperandDefinition{
				{ID: calculator.OperandFirst, Label: "First", Placeholder: "0"},
				{ID: calculator.OperandSecond, Label: "Second", Placeholder: "0"},
			},
		},
		Execute: validExecute,
	}

	tests := []struct {
		name               string
		defaultOperationID string
		operations         []calculator.Operation
		wantErrFragment    string
	}{
		{
			name:               "empty operations slice",
			defaultOperationID: "addition",
			operations:         []calculator.Operation{},
			wantErrFragment:    "at least one operation",
		},
		{
			name:               "empty default operation ID",
			defaultOperationID: "",
			operations:         []calculator.Operation{validOp},
			wantErrFragment:    "non-empty default",
		},
		{
			name:               "unknown default operation",
			defaultOperationID: "nonexistent",
			operations:         []calculator.Operation{validOp},
			wantErrFragment:    "not registered",
		},
		{
			name:               "duplicate operation IDs",
			defaultOperationID: "addition",
			operations:         []calculator.Operation{validOp, validOp},
			wantErrFragment:    "duplicate operation ID",
		},
		{
			name:               "empty operation ID",
			defaultOperationID: "addition",
			operations: []calculator.Operation{
				{
					Definition: calculator.OperationDefinition{ID: "", Name: "X", Symbol: "+", Shortcut: "+", Arity: calculator.Binary,
						Operands: []calculator.OperandDefinition{
							{ID: calculator.OperandFirst, Label: "A", Placeholder: "0"},
							{ID: calculator.OperandSecond, Label: "B", Placeholder: "0"},
						}},
					Execute: validExecute,
				},
			},
			wantErrFragment: "ID must not be empty",
		},
		{
			name:               "empty operation name",
			defaultOperationID: "op1",
			operations: []calculator.Operation{
				{
					Definition: calculator.OperationDefinition{ID: "op1", Name: "", Symbol: "+", Shortcut: "+", Arity: calculator.Binary,
						Operands: []calculator.OperandDefinition{
							{ID: calculator.OperandFirst, Label: "A", Placeholder: "0"},
							{ID: calculator.OperandSecond, Label: "B", Placeholder: "0"},
						}},
					Execute: validExecute,
				},
			},
			wantErrFragment: "name must not be empty",
		},
		{
			name:               "unsupported arity",
			defaultOperationID: "op1",
			operations: []calculator.Operation{
				{
					Definition: calculator.OperationDefinition{ID: "op1", Name: "Op", Symbol: "+", Shortcut: "+", Arity: 3,
						Operands: []calculator.OperandDefinition{
							{ID: "a", Label: "A", Placeholder: "0"},
							{ID: "b", Label: "B", Placeholder: "0"},
							{ID: "c", Label: "C", Placeholder: "0"},
						}},
					Execute: validExecute,
				},
			},
			wantErrFragment: "unsupported arity",
		},
		{
			name:               "operand count mismatch",
			defaultOperationID: "op1",
			operations: []calculator.Operation{
				{
					Definition: calculator.OperationDefinition{ID: "op1", Name: "Op", Symbol: "+", Shortcut: "+", Arity: calculator.Binary,
						Operands: []calculator.OperandDefinition{
							{ID: calculator.OperandFirst, Label: "Only one", Placeholder: "0"},
						}},
					Execute: validExecute,
				},
			},
			wantErrFragment: "arity 2 but 1 operand",
		},
		{
			name:               "duplicate operand IDs",
			defaultOperationID: "op1",
			operations: []calculator.Operation{
				{
					Definition: calculator.OperationDefinition{ID: "op1", Name: "Op", Symbol: "+", Shortcut: "+", Arity: calculator.Binary,
						Operands: []calculator.OperandDefinition{
							{ID: calculator.OperandFirst, Label: "A", Placeholder: "0"},
							{ID: calculator.OperandFirst, Label: "B", Placeholder: "0"},
						}},
					Execute: validExecute,
				},
			},
			wantErrFragment: "duplicate operand ID",
		},
		{
			name:               "nil executor",
			defaultOperationID: "op1",
			operations: []calculator.Operation{
				{
					Definition: calculator.OperationDefinition{ID: "op1", Name: "Op", Symbol: "+", Shortcut: "+", Arity: calculator.Unary,
						Operands: []calculator.OperandDefinition{
							{ID: calculator.OperandFirst, Label: "A", Placeholder: "0"},
						}},
					Execute: nil,
				},
			},
			wantErrFragment: "executor must not be nil",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := calculator.NewRegistry(testCase.defaultOperationID, testCase.operations)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), testCase.wantErrFragment) {
				t.Errorf("expected error containing %q, got: %q", testCase.wantErrFragment, err.Error())
			}
		})
	}
}

// --- Behavior tests ---

func TestRegistry_Find_allOperations(t *testing.T) {
	registry, err := calculator.NewDefaultRegistry()
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	ids := []string{"addition", "subtraction", "multiplication", "division", "exponentiation", "square-root", "percentage"}
	for _, operationID := range ids {
		t.Run(operationID, func(t *testing.T) {
			operation, found := registry.Find(operationID)
			if !found {
				t.Fatalf("expected operation %q to be found", operationID)
			}
			if operation.Definition.ID != operationID {
				t.Errorf("expected ID %q, got %q", operationID, operation.Definition.ID)
			}
		})
	}
}

func TestRegistry_Find_unknownOperationReturnsFalse(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	_, found := registry.Find("nonexistent")
	if found {
		t.Error("expected found=false for unknown operation ID")
	}
}

func TestRegistry_DefaultOperationID(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	if registry.DefaultOperationID() != "addition" {
		t.Errorf("expected default operation %q, got %q", "addition", registry.DefaultOperationID())
	}
}

func TestRegistry_Operations_deterministicOrder(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	operations := registry.Operations()

	expectedOrder := []string{
		"addition", "subtraction", "multiplication", "division",
		"exponentiation", "square-root", "percentage",
	}

	if len(operations) != len(expectedOrder) {
		t.Fatalf("expected %d operations, got %d", len(expectedOrder), len(operations))
	}
	for index, expectedID := range expectedOrder {
		if operations[index].Definition.ID != expectedID {
			t.Errorf("position %d: expected %q, got %q", index, expectedID, operations[index].Definition.ID)
		}
	}
}

func TestRegistry_Operations_returnsDefensiveCopy(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()

	first := registry.Operations()
	first[0].Definition.Operands[0].Label = "mutated"

	second := registry.Operations()
	if second[0].Definition.Operands[0].Label == "mutated" {
		t.Error("mutating the returned slice propagated into the registry")
	}
}

func TestRegistry_immutableAfterSourceSliceMutation(t *testing.T) {
	operations := calculator.DefaultOperations()
	registry, err := calculator.NewRegistry("addition", operations)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Mutate the source slice's first element.
	operations[0].Definition.Name = "Mutated Name"

	// The registry should retain the original name.
	operation, _ := registry.Find("addition")
	if operation.Definition.Name != "Addition" {
		t.Errorf("registry affected by post-construction mutation of input slice: got name %q", operation.Definition.Name)
	}
}

// --- Manifest projection tests ---

func TestManifest_version(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	manifest := registry.Manifest()
	if manifest.Version != "1" {
		t.Errorf("expected version %q, got %q", "1", manifest.Version)
	}
}

func TestManifest_defaultOperation(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	manifest := registry.Manifest()
	if manifest.DefaultOperationID != "addition" {
		t.Errorf("expected defaultOperationId %q, got %q", "addition", manifest.DefaultOperationID)
	}
}

func TestManifest_sevenOperations(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	manifest := registry.Manifest()
	if len(manifest.Operations) != 7 {
		t.Errorf("expected 7 operations, got %d", len(manifest.Operations))
	}
}

func TestManifest_deterministicOrder(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	manifest := registry.Manifest()

	expectedOrder := []string{
		"addition", "subtraction", "multiplication", "division",
		"exponentiation", "square-root", "percentage",
	}
	for index, expectedID := range expectedOrder {
		if manifest.Operations[index].ID != expectedID {
			t.Errorf("position %d: expected %q, got %q", index, expectedID, manifest.Operations[index].ID)
		}
	}
}

func TestManifest_jsonFieldNames(t *testing.T) {
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

	requiredTopLevel := []string{"version", "defaultOperationId", "operations"}
	for _, field := range requiredTopLevel {
		if _, exists := raw[field]; !exists {
			t.Errorf("missing top-level JSON field %q", field)
		}
	}

	ops, _ := raw["operations"].([]any)
	if len(ops) == 0 {
		t.Fatal("operations array is empty in JSON")
	}
	firstOp, _ := ops[0].(map[string]any)
	requiredOpFields := []string{"id", "name", "symbol", "shortcut", "arity", "operands", "validations"}
	for _, field := range requiredOpFields {
		if _, exists := firstOp[field]; !exists {
			t.Errorf("missing operation JSON field %q", field)
		}
	}

	operands, _ := firstOp["operands"].([]any)
	if len(operands) == 0 {
		t.Fatal("operands array is empty in JSON")
	}
	firstOperand, _ := operands[0].(map[string]any)
	requiredOperandFields := []string{"id", "label", "placeholder"}
	for _, field := range requiredOperandFields {
		if _, exists := firstOperand[field]; !exists {
			t.Errorf("missing operand JSON field %q", field)
		}
	}
}

func TestManifest_squareRoot_unaryWithNoSecondOperand(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	manifest := registry.Manifest()

	var squareRoot *calculator.ManifestOperation
	for index := range manifest.Operations {
		if manifest.Operations[index].ID == "square-root" {
			squareRoot = &manifest.Operations[index]
			break
		}
	}
	if squareRoot == nil {
		t.Fatal("square-root not found in manifest")
	}
	if squareRoot.Arity != 1 {
		t.Errorf("expected arity 1, got %d", squareRoot.Arity)
	}
	if len(squareRoot.Operands) != 1 {
		t.Errorf("expected 1 operand, got %d", len(squareRoot.Operands))
	}
	if squareRoot.Operands[0].ID != calculator.OperandFirst {
		t.Errorf("expected operand ID %q, got %q", calculator.OperandFirst, squareRoot.Operands[0].ID)
	}
}

func TestManifest_percentage_secondOperandHasSuffix(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()
	manifest := registry.Manifest()

	var percentage *calculator.ManifestOperation
	for index := range manifest.Operations {
		if manifest.Operations[index].ID == "percentage" {
			percentage = &manifest.Operations[index]
			break
		}
	}
	if percentage == nil {
		t.Fatal("percentage not found in manifest")
	}
	if len(percentage.Operands) < 2 {
		t.Fatalf("expected 2 operands, got %d", len(percentage.Operands))
	}
	if percentage.Operands[1].Suffix != "%" {
		t.Errorf("expected second operand suffix %q, got %q", "%", percentage.Operands[1].Suffix)
	}
}

func TestManifest_callerCannotMutateRegistryOperands(t *testing.T) {
	registry, _ := calculator.NewDefaultRegistry()

	manifest := registry.Manifest()
	manifest.Operations[0].Operands[0].Label = "mutated"

	secondManifest := registry.Manifest()
	if secondManifest.Operations[0].Operands[0].Label == "mutated" {
		t.Error("mutating the manifest propagated into the registry")
	}
}

// --- Registry construction: validation definition validation ---

// buildValidOp returns a minimal valid binary operation.
func buildValidOp(id string) calculator.Operation {
	return calculator.Operation{
		Definition: calculator.OperationDefinition{
			ID:       id,
			Name:     "Op",
			Symbol:   "+",
			Shortcut: "+",
			Arity:    calculator.Binary,
			Operands: []calculator.OperandDefinition{
				{ID: calculator.OperandFirst, Label: "First", Placeholder: "0"},
				{ID: calculator.OperandSecond, Label: "Second", Placeholder: "0"},
			},
		},
		Execute: func(operands []float64) (float64, error) { return 0, nil },
	}
}

// buildOpWithValidation returns a binary operation carrying the given validation.
func buildOpWithValidation(id string, validation calculator.ValidationDefinition) calculator.Operation {
	op := buildValidOp(id)
	op.Definition.Validations = []calculator.ValidationDefinition{validation}
	return op
}

func assertRegistryRejectsOp(t *testing.T, op calculator.Operation, wantFragment string) {
	t.Helper()
	_, err := calculator.NewRegistry(op.Definition.ID, []calculator.Operation{op})
	if err == nil {
		t.Fatal("expected registry construction error, got nil")
	}
	if !strings.Contains(err.Error(), wantFragment) {
		t.Errorf("expected error containing %q, got: %q", wantFragment, err.Error())
	}
}

func TestRegistry_ValidationDefinition_emptyID(t *testing.T) {
	op := buildOpWithValidation("op", calculator.ValidationDefinition{
		ID:      "",
		Message: "msg",
		Expression: calculator.ComparisonExpression{
			Operand: calculator.FirstOperand, Operator: calculator.Equal, Value: 0,
		},
	})
	assertRegistryRejectsOp(t, op, "ID must not be empty")
}

func TestRegistry_ValidationDefinition_emptyMessage(t *testing.T) {
	op := buildOpWithValidation("op", calculator.ValidationDefinition{
		ID:      "rule",
		Message: "",
		Expression: calculator.ComparisonExpression{
			Operand: calculator.FirstOperand, Operator: calculator.Equal, Value: 0,
		},
	})
	assertRegistryRejectsOp(t, op, "message must not be empty")
}

func TestRegistry_ValidationDefinition_nilExpression(t *testing.T) {
	op := buildOpWithValidation("op", calculator.ValidationDefinition{
		ID:         "rule",
		Message:    "msg",
		Expression: nil,
	})
	assertRegistryRejectsOp(t, op, "expression must not be nil")
}

func TestRegistry_ValidationDefinition_unsupportedOperandReference(t *testing.T) {
	op := buildOpWithValidation("op", calculator.ValidationDefinition{
		ID:      "rule",
		Message: "msg",
		Expression: calculator.ComparisonExpression{
			Operand: calculator.OperandReference("bogus"), Operator: calculator.Equal, Value: 0,
		},
	})
	assertRegistryRejectsOp(t, op, "unsupported operand reference")
}

func TestRegistry_ValidationDefinition_secondOperandOnUnaryOperation(t *testing.T) {
	op := calculator.Operation{
		Definition: calculator.OperationDefinition{
			ID:       "unary-op",
			Name:     "Op",
			Symbol:   "√",
			Shortcut: "r",
			Arity:    calculator.Unary,
			Operands: []calculator.OperandDefinition{
				{ID: calculator.OperandFirst, Label: "Number", Placeholder: "0"},
			},
			Validations: []calculator.ValidationDefinition{
				{
					ID:      "rule",
					Message: "msg",
					Expression: calculator.ComparisonExpression{
						Operand:  calculator.SecondOperand, // invalid for unary
						Operator: calculator.Equal,
						Value:    0,
					},
				},
			},
		},
		Execute: func(operands []float64) (float64, error) { return 0, nil },
	}
	assertRegistryRejectsOp(t, op, "requires arity")
}

func TestRegistry_ValidationDefinition_unsupportedOperator(t *testing.T) {
	op := buildOpWithValidation("op", calculator.ValidationDefinition{
		ID:      "rule",
		Message: "msg",
		Expression: calculator.ComparisonExpression{
			Operand: calculator.FirstOperand, Operator: calculator.ComparisonOperator("bogus"), Value: 0,
		},
	})
	assertRegistryRejectsOp(t, op, "unsupported comparison operator")
}

func TestRegistry_ValidationDefinition_nilAllOfChild(t *testing.T) {
	op := buildOpWithValidation("op", calculator.ValidationDefinition{
		ID:      "rule",
		Message: "msg",
		Expression: calculator.AllOfExpression{
			Expressions: []calculator.Expression{nil},
		},
	})
	assertRegistryRejectsOp(t, op, "must not be nil")
}

func TestRegistry_ValidationDefinition_malformedAllOfChild(t *testing.T) {
	op := buildOpWithValidation("op", calculator.ValidationDefinition{
		ID:      "rule",
		Message: "msg",
		Expression: calculator.AllOfExpression{
			Expressions: []calculator.Expression{
				calculator.ComparisonExpression{
					Operand:  calculator.OperandReference("bogus"),
					Operator: calculator.Equal,
					Value:    0,
				},
			},
		},
	})
	assertRegistryRejectsOp(t, op, "unsupported operand reference")
}

// --- Deep immutability ---

func TestRegistry_DeepImmutability_mutatingSourceAllOf(t *testing.T) {
	// Build a source AllOfExpression with two children.
	source := []calculator.Expression{
		calculator.ComparisonExpression{Operand: calculator.SecondOperand, Operator: calculator.GreaterThanOrEqual, Value: 0},
		calculator.ComparisonExpression{Operand: calculator.SecondOperand, Operator: calculator.LessThanOrEqual, Value: 100},
	}
	op := buildOpWithValidation("pct", calculator.ValidationDefinition{
		ID:      "range",
		Message: "must be 0–100",
		Expression: calculator.AllOfExpression{
			Expressions: source,
		},
	})
	registry, err := calculator.NewRegistry("pct", []calculator.Operation{op})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Mutate the source slice after registry construction.
	source[0] = calculator.ComparisonExpression{
		Operand:  calculator.SecondOperand,
		Operator: calculator.Equal,
		Value:    999,
	}

	found, _ := registry.Find("pct")

	// Registry should still enforce the original rules (>= 0 AND <= 100).
	if violation, err := found.Validate([]float64{0, 50}); violation != nil || err != nil {
		t.Errorf("50 should pass original rules; violation=%v err=%v", violation, err)
	}
	if violation, _ := found.Validate([]float64{0, -1}); violation == nil {
		t.Error("-1 should fail original rules; registry appears affected by source mutation")
	}

	// Manifest should still reflect the original expression.
	manifestData, _ := json.Marshal(registry.Manifest())
	raw := string(manifestData)
	if !strings.Contains(raw, "greaterThanOrEqual") {
		t.Error("manifest does not contain original operator after source mutation")
	}
	if strings.Contains(raw, "999") {
		t.Error("manifest was affected by source mutation")
	}
}

func TestRegistry_DeepImmutability_mutatingReturnedOperation(t *testing.T) {
	source := []calculator.Expression{
		calculator.ComparisonExpression{Operand: calculator.SecondOperand, Operator: calculator.GreaterThanOrEqual, Value: 0},
		calculator.ComparisonExpression{Operand: calculator.SecondOperand, Operator: calculator.LessThanOrEqual, Value: 100},
	}
	op := buildOpWithValidation("pct", calculator.ValidationDefinition{
		ID:      "range",
		Message: "must be 0–100",
		Expression: calculator.AllOfExpression{
			Expressions: source,
		},
	})
	registry, err := calculator.NewRegistry("pct", []calculator.Operation{op})
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Get a copy from Find and mutate its expression tree.
	returned, _ := registry.Find("pct")
	allOf := returned.Definition.Validations[0].Expression.(calculator.AllOfExpression)
	allOf.Expressions[0] = calculator.ComparisonExpression{
		Operand:  calculator.SecondOperand,
		Operator: calculator.Equal,
		Value:    999,
	}

	// The registry's stored copy must be unaffected.
	canonical, _ := registry.Find("pct")
	if violation, _ := canonical.Validate([]float64{0, -1}); violation == nil {
		t.Error("-1 should fail registry's canonical rules; Find copy mutation propagated into registry")
	}
}
