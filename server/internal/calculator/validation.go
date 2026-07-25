package calculator

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Sentinel errors for expression evaluation failures.
var (
	// ErrUnsupportedOperandReference is returned when an expression references an
	// operand position that is not a known constant in this package.
	ErrUnsupportedOperandReference = errors.New("unsupported operand reference")
	// ErrOperandOutOfBounds is returned when the referenced operand position is valid
	// but the provided operand slice is too short to supply a value at that position.
	ErrOperandOutOfBounds = errors.New("operand out of bounds")
	// ErrUnsupportedComparisonOperator is returned when a ComparisonExpression carries
	// an operator value that is not one of the six supported operators.
	ErrUnsupportedComparisonOperator = errors.New("unsupported comparison operator")
)

// Expression is the sealed base type for all validation expressions.
// Only types defined in this package can implement Expression.
type Expression interface {
	expressionNode()
	evaluate(operands []float64) (bool, error)
}

// OperandReference identifies which operand value to read during evaluation.
type OperandReference string

// Supported OperandReference values.
const (
	// FirstOperand references the first positional operand (operands[0]).
	FirstOperand OperandReference = "first"
	// SecondOperand references the second positional operand (operands[1]).
	SecondOperand OperandReference = "second"
)

// ComparisonOperator defines the set of supported comparison operations.
type ComparisonOperator string

// Supported ComparisonOperator values. These match the frontend TypeScript contract.
const (
	Equal              ComparisonOperator = "equal"
	NotEqual           ComparisonOperator = "notEqual"
	GreaterThan        ComparisonOperator = "greaterThan"
	GreaterThanOrEqual ComparisonOperator = "greaterThanOrEqual"
	LessThan           ComparisonOperator = "lessThan"
	LessThanOrEqual    ComparisonOperator = "lessThanOrEqual"
)

// ComparisonExpression evaluates to true when the referenced operand satisfies
// the comparison against the literal value.
type ComparisonExpression struct {
	// Operand identifies which operand value to compare.
	Operand OperandReference
	// Operator is the comparison relationship to apply.
	Operator ComparisonOperator
	// Value is the literal constant on the right-hand side of the comparison.
	Value float64
}

func (ComparisonExpression) expressionNode() {}

func (expr ComparisonExpression) evaluate(operands []float64) (bool, error) {
	operandValue, err := resolveOperand(expr.Operand, operands)
	if err != nil {
		return false, err
	}
	switch expr.Operator {
	case Equal:
		return operandValue == expr.Value, nil
	case NotEqual:
		return operandValue != expr.Value, nil
	case GreaterThan:
		return operandValue > expr.Value, nil
	case GreaterThanOrEqual:
		return operandValue >= expr.Value, nil
	case LessThan:
		return operandValue < expr.Value, nil
	case LessThanOrEqual:
		return operandValue <= expr.Value, nil
	default:
		return false, fmt.Errorf("%w: %q", ErrUnsupportedComparisonOperator, expr.Operator)
	}
}

// MarshalJSON produces the discriminated-union JSON expected by the frontend:
// { "kind": "comparison", "operand": "...", "operator": "...", "value": N }
func (expr ComparisonExpression) MarshalJSON() ([]byte, error) {
	type wire struct {
		Kind     string             `json:"kind"`
		Operand  OperandReference   `json:"operand"`
		Operator ComparisonOperator `json:"operator"`
		Value    float64            `json:"value"`
	}
	return json.Marshal(wire{
		Kind:     "comparison",
		Operand:  expr.Operand,
		Operator: expr.Operator,
		Value:    expr.Value,
	})
}

// AllOfExpression evaluates to true when every child expression evaluates to true.
// An empty Expressions slice evaluates to true (vacuous truth).
type AllOfExpression struct {
	// Expressions is the ordered list of conditions that must all be satisfied.
	Expressions []Expression
}

func (AllOfExpression) expressionNode() {}

func (expr AllOfExpression) evaluate(operands []float64) (bool, error) {
	for index, child := range expr.Expressions {
		passed, err := child.evaluate(operands)
		if err != nil {
			return false, fmt.Errorf("allOf child at index %d: %w", index, err)
		}
		if !passed {
			return false, nil
		}
	}
	return true, nil
}

// MarshalJSON produces the discriminated-union JSON expected by the frontend:
// { "kind": "allOf", "expressions": [ ... ] }
func (expr AllOfExpression) MarshalJSON() ([]byte, error) {
	type wire struct {
		Kind        string       `json:"kind"`
		Expressions []Expression `json:"expressions"`
	}
	return json.Marshal(wire{Kind: "allOf", Expressions: expr.Expressions})
}

// ValidationDefinition is a named, declarative constraint on an operation's operands.
// ID is a stable backend identifier; Message is the user-facing error text.
type ValidationDefinition struct {
	// ID is a stable backend identifier used to trace which rule was violated.
	ID string
	// Message is the user-facing error text returned when the validation fails.
	Message string
	// Expression is the declarative predicate; the validation fails when it evaluates to false.
	Expression Expression
}

// resolveOperand returns the operand value at the given reference position.
// It returns ErrOperandOutOfBounds when the slice is too short, or
// ErrUnsupportedOperandReference when the reference is not a known constant.
func resolveOperand(reference OperandReference, operands []float64) (float64, error) {
	switch reference {
	case FirstOperand:
		if len(operands) < 1 {
			return 0, fmt.Errorf("%w: %q requires at least 1 operand, got %d", ErrOperandOutOfBounds, reference, len(operands))
		}
		return operands[0], nil
	case SecondOperand:
		if len(operands) < 2 {
			return 0, fmt.Errorf("%w: %q requires at least 2 operands, got %d", ErrOperandOutOfBounds, reference, len(operands))
		}
		return operands[1], nil
	default:
		return 0, fmt.Errorf("%w: %q", ErrUnsupportedOperandReference, reference)
	}
}

// copyExpression returns a deep copy of an Expression.
// ComparisonExpression is a value type and is copied by value.
// AllOfExpression.Expressions is recursively copied to prevent slice aliasing.
func copyExpression(expression Expression) Expression {
	if expression == nil {
		return nil
	}
	switch expr := expression.(type) {
	case ComparisonExpression:
		return expr
	case AllOfExpression:
		childrenCopy := make([]Expression, len(expr.Expressions))
		for index, child := range expr.Expressions {
			childrenCopy[index] = copyExpression(child)
		}
		return AllOfExpression{Expressions: childrenCopy}
	default:
		panic(fmt.Sprintf("copyExpression: unsupported expression type %T", expression))
	}
}

// validateExpressionDefinition checks that an expression is structurally valid
// for an operation of the given arity. It recurses into AllOfExpression children.
// Returns a descriptive error on the first structural violation found.
func validateExpressionDefinition(expression Expression, arity Arity) error {
	if expression == nil {
		return errors.New("expression must not be nil")
	}
	switch expr := expression.(type) {
	case ComparisonExpression:
		switch expr.Operand {
		case FirstOperand:
			// valid for any arity
		case SecondOperand:
			if arity < Binary {
				return fmt.Errorf("operand reference %q requires arity %d or greater, got %d", expr.Operand, Binary, arity)
			}
		default:
			return fmt.Errorf("%w: %q", ErrUnsupportedOperandReference, expr.Operand)
		}
		switch expr.Operator {
		case Equal, NotEqual, GreaterThan, GreaterThanOrEqual, LessThan, LessThanOrEqual:
			// valid
		default:
			return fmt.Errorf("%w: %q", ErrUnsupportedComparisonOperator, expr.Operator)
		}
		return nil
	case AllOfExpression:
		for index, child := range expr.Expressions {
			if child == nil {
				return fmt.Errorf("allOf expression at index %d must not be nil", index)
			}
			if err := validateExpressionDefinition(child, arity); err != nil {
				return fmt.Errorf("allOf expression at index %d: %w", index, err)
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported expression type %T", expression)
	}
}
