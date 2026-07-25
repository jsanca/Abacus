package calculator

import (
	"fmt"
	"math"
)

// DefaultOperations returns the seven standard Abacus calculator operations in
// deterministic product-defined order. Each call returns a new slice.
func DefaultOperations() []Operation {
	return []Operation{
		additionOperation(),
		subtractionOperation(),
		multiplicationOperation(),
		divisionOperation(),
		exponentiationOperation(),
		squareRootOperation(),
		percentageOperation(),
	}
}

// NewDefaultRegistry is the standard composition root entry point. It constructs
// an immutable Registry from the canonical operation set with "addition" as default.
func NewDefaultRegistry() (*Registry, error) {
	return NewRegistry("addition", DefaultOperations())
}

func additionOperation() Operation {
	return Operation{
		Definition: OperationDefinition{
			ID:       "addition",
			Name:     "Addition",
			Symbol:   "+",
			Shortcut: "+",
			Arity:    Binary,
			Operands: []OperandDefinition{
				{ID: OperandFirst, Label: "First number", Placeholder: "0"},
				{ID: OperandSecond, Label: "Second number", Placeholder: "0"},
			},
		},
		Execute: executeAddition,
	}
}

func subtractionOperation() Operation {
	return Operation{
		Definition: OperationDefinition{
			ID:       "subtraction",
			Name:     "Subtraction",
			Symbol:   "−",
			Shortcut: "-",
			Arity:    Binary,
			Operands: []OperandDefinition{
				{ID: OperandFirst, Label: "First number", Placeholder: "0"},
				{ID: OperandSecond, Label: "Second number", Placeholder: "0"},
			},
		},
		Execute: executeSubtraction,
	}
}

func multiplicationOperation() Operation {
	return Operation{
		Definition: OperationDefinition{
			ID:       "multiplication",
			Name:     "Multiplication",
			Symbol:   "×",
			Shortcut: "*",
			Arity:    Binary,
			Operands: []OperandDefinition{
				{ID: OperandFirst, Label: "First number", Placeholder: "0"},
				{ID: OperandSecond, Label: "Second number", Placeholder: "0"},
			},
		},
		Execute: executeMultiplication,
	}
}

func divisionOperation() Operation {
	return Operation{
		Definition: OperationDefinition{
			ID:       "division",
			Name:     "Division",
			Symbol:   "÷",
			Shortcut: "/",
			Arity:    Binary,
			Operands: []OperandDefinition{
				{ID: OperandFirst, Label: "Dividend", Placeholder: "0"},
				{ID: OperandSecond, Label: "Divisor", Placeholder: "0"},
			},
			Validations: []ValidationDefinition{
				{
					ID:      "division-by-zero",
					Message: "The divisor must not be zero.",
					Expression: ComparisonExpression{
						Operand:  SecondOperand,
						Operator: NotEqual,
						Value:    0,
					},
				},
			},
		},
		Execute: executeDivision,
	}
}

func exponentiationOperation() Operation {
	return Operation{
		Definition: OperationDefinition{
			ID:       "exponentiation",
			Name:     "Exponentiation",
			Symbol:   "xʸ",
			Shortcut: "^",
			Arity:    Binary,
			Operands: []OperandDefinition{
				{ID: OperandFirst, Label: "Base", Placeholder: "0"},
				{ID: OperandSecond, Label: "Exponent", Placeholder: "0"},
			},
		},
		Execute: executeExponentiation,
	}
}

func squareRootOperation() Operation {
	return Operation{
		Definition: OperationDefinition{
			ID:       "square-root",
			Name:     "Square Root",
			Symbol:   "√",
			Shortcut: "r",
			Arity:    Unary,
			Operands: []OperandDefinition{
				// Prefix notation: the symbol precedes the operand (√144, not 144√).
				{ID: OperandFirst, Label: "Number", Placeholder: "0"},
			},
			Validations: []ValidationDefinition{
				{
					ID:      "square-root-negative",
					Message: "The number must be zero or greater.",
					Expression: ComparisonExpression{
						Operand:  FirstOperand,
						Operator: GreaterThanOrEqual,
						Value:    0,
					},
				},
			},
		},
		Execute: executeSquareRoot,
	}
}

func percentageOperation() Operation {
	return Operation{
		Definition: OperationDefinition{
			ID:       "percentage",
			Name:     "Percentage",
			Symbol:   "%",
			Shortcut: "%",
			Arity:    Binary,
			Operands: []OperandDefinition{
				{ID: OperandFirst, Label: "Base value", Placeholder: "0"},
				// The non-editable suffix "%" is displayed alongside this field.
				{ID: OperandSecond, Label: "Percentage", Placeholder: "0", Suffix: "%"},
			},
			Validations: []ValidationDefinition{
				{
					ID:      "percentage-range",
					Message: "Percentage must be between 0 and 100.",
					Expression: AllOfExpression{
						Expressions: []Expression{
							ComparisonExpression{Operand: SecondOperand, Operator: GreaterThanOrEqual, Value: 0},
							ComparisonExpression{Operand: SecondOperand, Operator: LessThanOrEqual, Value: 100},
						},
					},
				},
			},
		},
		Execute: executePercentage,
	}
}

func executeAddition(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, fmt.Errorf("%w: addition requires 2, got %d", ErrWrongOperandCount, len(operands))
	}
	return operands[0] + operands[1], nil
}

func executeSubtraction(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, fmt.Errorf("%w: subtraction requires 2, got %d", ErrWrongOperandCount, len(operands))
	}
	return operands[0] - operands[1], nil
}

func executeMultiplication(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, fmt.Errorf("%w: multiplication requires 2, got %d", ErrWrongOperandCount, len(operands))
	}
	return operands[0] * operands[1], nil
}

func executeDivision(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, fmt.Errorf("%w: division requires 2, got %d", ErrWrongOperandCount, len(operands))
	}
	// Division by zero returns ±Inf or NaN per Go float64 semantics.
	// The domain validation constraint (divisor ≠ 0) is enforced in ABA-006.
	return operands[0] / operands[1], nil
}

func executeExponentiation(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, fmt.Errorf("%w: exponentiation requires 2, got %d", ErrWrongOperandCount, len(operands))
	}
	return math.Pow(operands[0], operands[1]), nil
}

func executeSquareRoot(operands []float64) (float64, error) {
	if len(operands) != 1 {
		return 0, fmt.Errorf("%w: square-root requires 1, got %d", ErrWrongOperandCount, len(operands))
	}
	// math.Sqrt of a negative number returns NaN per Go float64 semantics.
	// The domain validation constraint (radicand ≥ 0) is enforced in ABA-006.
	return math.Sqrt(operands[0]), nil
}

func executePercentage(operands []float64) (float64, error) {
	if len(operands) != 2 {
		return 0, fmt.Errorf("%w: percentage requires 2, got %d", ErrWrongOperandCount, len(operands))
	}
	// Percentage semantics: what percentage of the base equals the result.
	// e.g. 15% of 200 = 200 × 15 / 100 = 30.
	return operands[0] * operands[1] / 100, nil
}
