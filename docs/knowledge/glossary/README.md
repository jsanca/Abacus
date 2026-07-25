# Glossary

Shared product and engineering terms used consistently across requirements, implementation, and review.

## Calculator Domain

**Operation**
A named, identifiable calculator capability pairing a frontend-safe `OperationDefinition` with a trusted backend `ExecuteFunc`. The definition is safe to serialize in the manifest; the function is never exposed to the client.

**OperationDefinition**
The frontend-safe metadata for a single operation: stable ID, human-readable name, mathematical symbol, keyboard shortcut, arity, and ordered operand definitions.

**OperandDefinition**
The metadata for a single input slot within an operation: stable ID (`first` or `second`), user-facing label, placeholder text, and an optional non-editable suffix.

**ExecuteFunc**
A Go function type `func(operands []float64) (float64, error)` that performs arithmetic for one operation. Callers must validate arity and domain constraints before invoking; executors guard only against wrong operand count.

**Registry**
The immutable, ordered, server-side collection of all supported operations. Constructed once at startup and never modified thereafter. Concurrent reads require no synchronization. Provides ordered enumeration (for manifest projection) and O(1) lookup by operation ID.

**Manifest**
The frontend-safe JSON projection of the Registry. Contains `version`, `defaultOperationId`, and an ordered `operations` array of `ManifestOperation` records. Excludes executor functions and backend-only runtime details.

**Arity**
The number of operands an operation requires. Supported values: `1` (unary, e.g. square root) and `2` (binary, e.g. addition).

**Default operation**
The operation selected when the calculator loads for the first time. Currently: `addition`.

## Validation

**ValidationDefinition**
A declarative rule associated with an operation: a stable backend `ID`, a user-facing `Message`, and an `Expression` tree that must evaluate to `true` for the operands to be valid. The first violated rule is returned by `Validate`; the `ID` is backend-only and excluded from the manifest.

**Expression**
A sealed interface with two unexported methods (`expressionNode()` marker and `evaluate(operands []float64) (bool, error)`). Only `ComparisonExpression` and `AllOfExpression` are valid implementations. Each type owns the semantics of its own evaluation; the interface is not exposed as a general-purpose evaluator.

**ComparisonExpression**
A leaf `Expression` that compares one operand (by `OperandReference`) against a literal `float64` value using a `ComparisonOperator`. Serializes to `{"kind":"comparison","operand":"...","operator":"...","value":N}`.

**AllOfExpression**
A composite `Expression` that evaluates to true when every child expression in its `Expressions` slice evaluates to true. An empty slice is vacuously true. Serializes to `{"kind":"allOf","expressions":[...]}`.

**OperandReference**
A typed string identifying which positional operand to read during expression evaluation. Supported values: `"first"` (operands[0]) and `"second"` (operands[1]).

**ComparisonOperator**
A typed string identifying the comparison relationship in a `ComparisonExpression`. Supported values: `equal`, `notEqual`, `greaterThan`, `greaterThanOrEqual`, `lessThan`, `lessThanOrEqual`. These match the TypeScript frontend contract.
