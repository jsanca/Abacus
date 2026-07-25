// Package calculator implements the Abacus operation registry, manifest projection,
// and trusted arithmetic executors. It is the domain layer and must not import
// transport, app, config, or observability packages.
package calculator

import "errors"

// Arity represents the number of operands a calculator operation requires.
type Arity int

// Supported arity values. Every registered operation must declare exactly one of these.
const (
	Unary  Arity = 1
	Binary Arity = 2
)

// Operand position identifiers used by the calculator UI.
// These are stable string values referenced by the client contract.
const (
	OperandFirst  = "first"
	OperandSecond = "second"
)

// ErrWrongOperandCount is returned by an executor when the provided operand
// slice length does not match the operation's declared arity.
var ErrWrongOperandCount = errors.New("wrong number of operands")

// ExecuteFunc performs a calculator operation on the given operands and returns
// the numeric result or an error. Callers are expected to validate arity and
// domain constraints before invoking an executor; executors guard only against
// wrong operand count. Domain validation (e.g. divisor not zero) is introduced
// in ABA-006.
type ExecuteFunc func(operands []float64) (float64, error)

// OperandDefinition describes a single input slot for a calculator operation.
// All fields are safe to expose to the frontend.
type OperandDefinition struct {
	// ID is a stable identifier used by the frontend to reference this operand.
	// Supported values are OperandFirst and OperandSecond.
	ID string
	// Label is the user-facing name displayed beside the input field.
	Label string
	// Placeholder is shown inside the input field before a value is entered.
	Placeholder string
	// Suffix is an optional non-editable annotation displayed after the input.
	// For example, the percentage operation's second operand carries the suffix "%".
	Suffix string
}

// OperationDefinition contains the frontend-safe metadata for a calculator operation.
// It does not include the executor function.
type OperationDefinition struct {
	// ID is the stable operation identifier used in registry lookups and the client contract.
	ID string
	// Name is the human-readable operation name displayed in the UI.
	Name string
	// Symbol is the mathematical symbol used in expressions and the operation selector.
	Symbol string
	// Shortcut is the keyboard character that activates this operation.
	Shortcut string
	// Arity is the number of operands the operation requires.
	Arity Arity
	// Operands describes each input slot in presentation order.
	Operands []OperandDefinition
}

// Operation pairs a frontend-safe definition with a trusted backend executor.
// The Execute function must not be exposed to the frontend or serialized.
type Operation struct {
	// Definition holds the metadata used in the manifest projection.
	Definition OperationDefinition
	// Execute performs the arithmetic for this operation.
	Execute ExecuteFunc
}

// copyOperation returns a deep copy of an Operation with its own Operands backing array.
// This prevents callers from aliasing the registry's internal operand storage.
func copyOperation(operation Operation) Operation {
	operandsCopy := make([]OperandDefinition, len(operation.Definition.Operands))
	copy(operandsCopy, operation.Definition.Operands)
	operation.Definition.Operands = operandsCopy
	return operation
}
