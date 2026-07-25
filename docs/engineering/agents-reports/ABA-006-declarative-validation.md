# ABA-006 — Declarative Validation Model: Agent Report

**Date:** 2026-07-25
**Status:** Complete
**Agent:** Clio (claude-sonnet-4-6)

---

## Objective

Implement the declarative validation layer for the Abacus operation registry: a typed expression hierarchy, recursive evaluator, manifest projection of validation trees, and domain constraints for the three operations that require them.

---

## Deliverables

### `server/internal/calculator/validation.go` (new)

- **Sealed `Expression` interface** — unexported marker method `expressionNode()` closes the type set to this package, enabling exhaustive type-switch evaluation without reflection.
- **`ComparisonExpression`** — leaf node comparing one operand (by `OperandReference`) to a literal `float64` value. Implements `MarshalJSON` to produce `{"kind":"comparison","operand":"...","operator":"...","value":N}`.
- **`AllOfExpression`** — composite node; evaluates to true when all children pass (vacuously true for empty slice). Implements `MarshalJSON` to produce `{"kind":"allOf","expressions":[...]}`.
- **`ValidationDefinition`** — named rule pairing a backend `ID` with a user-facing `Message` and an `Expression` tree.
- **`Evaluate(expression, operands)`** — recursive evaluator; type-switches on `ComparisonExpression` (six operators) and `AllOfExpression`. Short-circuits on first failing child.
- **`Validate(definition, operands)`** — iterates `definition.Validations` in declaration order; returns a pointer to the first violated `ValidationDefinition`, or nil.

### `server/internal/calculator/operation.go` (modified)

- Added `Validations []ValidationDefinition` to `OperationDefinition`.
- Updated `copyOperation` to deep-copy the `Validations` slice, preserving the aliasing-prevention invariant established in ABA-005.

### `server/internal/calculator/operations.go` (modified)

Validation rules added to three operations:

| Operation | Rule ID | Expression |
|---|---|---|
| Division | `division-by-zero` | `second != 0` |
| Square Root | `square-root-negative` | `first >= 0` |
| Percentage | `percentage-range` | `second >= 0 AND second <= 100` |

### `server/internal/calculator/manifest.go` (modified)

- Removed `encoding/json` import (no more `json.RawMessage` placeholder).
- Added `ManifestValidation` struct projecting `Message` and `Expression`; the backend-only `ID` is excluded.
- Changed `ManifestOperation.Validations` from `json.RawMessage` to `[]ManifestValidation`.
- `projectOperation` now builds the validations projection from `definition.Validations`.

### `server/internal/calculator/registry_test.go` (modified)

- Removed `TestManifest_validationsAlwaysEmptyArray` — the ABA-005 placeholder test that asserted all operations serialize an empty validations array. Superseded by `TestManifest_validationsSerializedCorrectly` in the new validation test file.

### `server/internal/calculator/validation_test.go` (new)

Comprehensive test coverage across four concerns:

| Test | Cases |
|---|---|
| `TestEvaluate_comparisonOperators` | 16 table-driven cases covering all six operators |
| `TestEvaluate_operandReferences` | first and second operand index resolution |
| `TestEvaluate_allOf` | all-pass, one-fail, empty vacuous truth |
| `TestValidate_division` | 12/3 pass, 12/0 fail with ID assertion |
| `TestValidate_squareRoot` | sqrt(9) pass, sqrt(-1) fail with ID assertion |
| `TestValidate_percentage` | 5 cases: 50%, 0%, 100% pass; 101%, -5% fail |
| `TestValidate_operationsWithoutValidations_alwaysPass` | addition, subtraction, multiplication, exponentiation |
| `TestManifest_validationsSerializedCorrectly` | JSON structure for addition (empty), division (comparison kind/operand/operator), percentage (allOf with 2 inner expressions), square-root |
| `TestManifest_validationsDeterministicOrdering` | two Manifest() calls produce identical JSON |
| `TestManifest_executorNotExposedInValidation` | serialized manifest contains no "func" substring |

---

## Key Design Decisions

**Sealed interface over open interface**: `expressionNode()` being unexported means external packages cannot implement `Expression`, keeping the evaluator's type switch exhaustive and the expression set intentionally closed. The frontend contract defines exactly two expression shapes; open extensibility would be premature.

**`MarshalJSON` on concrete types**: Each expression type carries its own serialization logic, producing the discriminated-union `{"kind":"..."}` envelope the TypeScript client expects. This keeps serialization co-located with the type definition rather than scattered across the manifest projection.

**`Validate` returns pointer to first violation**: Returning `*ValidationDefinition` gives callers both the user-facing message and the stable backend ID for tracing, without allocating in the happy path (nil return).

**`copyOperation` deep-copies validations**: Preserves the aliasing-prevention invariant from ABA-005. The `ValidationDefinition` slice is copied; the `Expression` interface values inside it are not (interface values are immutable from the caller's perspective since `Expression` has no mutating methods).

---

## Debt Resolved

**TD-002** — The `validations: []` manifest placeholder introduced in ABA-005 is now resolved. All three constrained operations emit live validation trees; the client can evaluate them without a server round-trip.

---

## Verification

```
make verify        # fmt check, lint (0 issues), race tests (all pass), build
go test -race ./... # All packages green
```
