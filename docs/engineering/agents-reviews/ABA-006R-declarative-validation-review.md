# ABA-006R — Declarative Validation Architecture Review

**Date:** 2026-07-25
**Status:** Complete
**Reviewer:** Deep Pro
**Role:** Independent Go & Domain Architecture Reviewer
**Target:** 30–40 minutes
**Commit:** None
**Repository:** /Users/jonathan/code/Abacus

---

## Objective

Review the final declarative validation architecture delivered through ABA-006 and refined by ABA-006.1. Evaluate the cohesive design for correctness, idiomaticity, safety, immutability, contract compatibility, proportionality, and REST readiness.

---

## Scope Reviewed

- **Go:** `validation.go`, `operation.go`, `operations.go`, `registry.go`, `manifest.go`
- **Tests:** `validation_test.go`, `registry_test.go`, `operations_test.go`
- **Frontend:** `validation.ts`, `api/types.ts`
- **Reports:** ABA-006 implementation, ABA-006.1 implementation
- **Verification:** `make verify`, `go test -race ./...`, `docker compose build server`

---

## Verification Results

```
make verify         → All checks passed (fmt, lint 0 issues, race tests, build)
go test -race ./...  → calculator: ok (coverage: 95.3% of statements)
docker compose build server → Image built successfully
```

### Coverage Detail (key functions)

| Function | Coverage |
|---|---|
| `Validate` | 100.0% |
| `copyOperation` | 100.0% |
| `NewRegistry` | 100.0% |
| `validateRegistryOperation` | 90.0% |
| `evaluate` (ComparisonExpression) | 81.8% |
| `evaluate` (AllOfExpression) | 100.0% |
| `resolveOperand` | 100.0% |
| `copyExpression` | 77.8% |
| `validateExpressionDefinition` | 94.1% |
| `Manifest` / `projectOperation` | 100.0% |

Uncovered branches are unreachable under the sealed interface (`expressionNode()` no-ops, default cases in `copyExpression` and `validateExpressionDefinition`, and the `copyExpression` nil-input guard).

---

## Findings

### F1: `copyExpression` default branch silently aliases (MINOR)

**Location:** `server/internal/calculator/validation.go:185–186`

```go
default:
    return expression
```

**Evidence:** Under the current sealed interface (only `ComparisonExpression` and `AllOfExpression`), this default branch is unreachable. If a third expression type were added but the `copyExpression` switch were not updated, the new type would be returned by reference — silently breaking the registry's deep immutability guarantee.

**Why it matters:** The `copyExpression` contract is that it returns an independent copy. The default branch violates that contract. Since deep immutability is a documented invariant of the registry (proven by dedicated tests), any path that could break it silently is a risk. The sealed interface currently prevents this at the type level, but the function itself does not enforce its own contract.

**Recommended change:** Replace the default branch with a panic with a descriptive message:

```go
default:
    panic(fmt.Sprintf("copyExpression: unsupported expression type %T", expression))
```

This follows the `template.Must`-style init-time pattern — if this branch is ever reached, something has broken the sealed interface and it must fail loudly rather than silently.

---

### N1: `validateExpressionDefinition` default branch is unreachable (NOTE)

**Location:** `server/internal/calculator/validation.go:226–228`

Same sealed-interface dead-code pattern as F1. The default branch `"unsupported expression type %T"` cannot be reached because no other types implement `Expression`. It exists as defensive code. Consider a comment documenting this intent.

---

### N2: Manifest `version` and `shortcut` fields not in TypeScript contract (NOTE)

**Location:** `server/internal/calculator/manifest.go:13–14`, `28`

Previously recorded as finding F1 in ABA-005R. The manifest includes `version` and `shortcut` fields that are absent from the TypeScript `OperationManifest` type. These are forward-compatible extras — JSON deserialization ignores unknown fields. Should be aligned during frontend-backend integration. Not a regression.

---

## Review Dimensions

### 1. Ownership and Cohesion

`Operation.Validate(operands)` is the sole validation entry point. Each operation declaration in `operations.go` carries metadata, operand definitions, validation rules, and its executor together — a single point of truth when adding an operation. The centralized `Evaluate` and `Validate` free functions from ABA-006 were removed. There is no operation-ID-based validation dispatcher elsewhere.

**Assessment:** PASS. Cohesive single entry point.

### 2. Expression Interface

The `Expression` interface carries two unexported methods: `expressionNode()` (sealing marker) and `evaluate(operands []float64) (bool, error)`. Only `ComparisonExpression` and `AllOfExpression` implement it. The interface cannot be extended from outside the package. Both implementations use value receivers — appropriate for immutable value types.

Nil-interface safety is addressed at three layers: `copyExpression` guards input, `validateExpressionDefinition` rejects nil at registry construction, and `Operation.Validate` guards the expression before calling `evaluate`.

**Assessment:** PASS. Well-designed sealed interface.

### 3. Runtime Evaluation

`resolveOperand` returns `ErrOperandOutOfBounds` for missing operands and `ErrUnsupportedOperandReference` for unknown references — never silently resolves to zero. `ComparisonExpression.evaluate` returns `ErrUnsupportedComparisonOperator` for unknown operators. `AllOfExpression.evaluate` short-circuits on first failure and propagates child errors with index context. All errors use `%w` wrapping, preserving the chain for `errors.Is` inspection.

Error context is structured: `operation "id": validation "rule": allOf child at index N: operand out of bounds: ...`

No panics were found in evaluation paths. The nil-expression guard in `Operation.Validate` prevents a nil-interface panic before evaluation.

**Assessment:** PASS. Never silently resolves to zero. Errors are wrapped and inspectable.

### 4. Operation Validation

`Operation.Validate` returns `(*ValidationDefinition, error)` — two distinct paths: violation (nil error, non-nil pointer) and evaluation failure (non-nil error). The first violated validation is returned; no error collection. The happy path returns `nil, nil` with no heap allocation.

The `&violated` expression takes the address of a local variable, which Go correctly escapes to the heap. The pointer is to a copy of the validation definition (`violated := validation`), so the registry's stored definition is not exposed.

The two-return design is natural for the future application service:

```go
violation, err := operation.Validate(operands)
if err != nil { /* evaluation failure */ }
if violation != nil { /* user-facing constraint */ }
```

**Assessment:** PASS. Clean dual-purpose return, minimal allocations.

### 5. Definition Validation (Centralized vs. Expression-Owned)

`validateExpressionDefinition` performs a type switch on the sealed expression set, validating structural integrity against operation arity. Both expression types already own `evaluate()` (runtime) and `MarshalJSON()` (wire). Adding `validateDefinition(arity)` would follow the same ownership pattern.

However, `evaluate` is runtime behavior (needs operand values) while `validateDefinition` is startup behavior (needs operation arity metadata). These are different concerns evaluated at different lifecycle phases. For two expression types, a single centralized validator is proportional and clearer than distributing the concern across types.

**Explicit trade-off:** The central switch is the right choice. It's 35 lines, covers all structural invariants, and is easier to audit than tracing through two separate `validateDefinition` method implementations. Concentration on nodes without clear value is avoided per the skill guidance.

**Assessment:** PASS. Centralized validator is proportional and auditable.

### 6. Deep Immutability

| Path | How | Verified |
|---|---|---|
| Source validation slices → registry | `copyOperation` deep-copies each validation via `copyExpression` | `TestRegistry_immutableAfterSourceSliceMutation` |
| `AllOfExpression.Expressions` slices | `copyExpression` allocates new slice, recursively copies children | `TestRegistry_DeepImmutability_mutatingSourceAllOf` |
| `Find()` results → registry | `Find` returns `copyOperation` | `TestRegistry_DeepImmutability_mutatingReturnedOperation` |
| Manifest projections → registry | `projectOperation` calls `copyExpression` per validation | `TestManifest_callerCannotMutateRegistryOperands` |
| Interface values | All concrete types are value types (`ComparisonExpression`, `AllOfExpression`); no interior pointers | Structural |

The one gap is the `copyExpression` default branch (see F1).

**Assessment:** PASS WITH OBSERVATIONS. Proven by tests for every access path. One minor note on the default branch.

### 7. Registry Fail-Fast Behavior

`validateRegistryOperation` validates every `ValidationDefinition` at construction:

- Validation ID non-empty
- Message non-empty
- Expression non-nil
- Operand reference is a known constant
- Operand reference is compatible with operation arity
- Comparison operator is supported
- AllOf children are non-nil and recursively valid

All checks complete before `NewRegistry` returns. No malformed validation reaches a user request.

Unreachable-at-startup cases: the `default` branch for unknown expression types (N1) and the `copyExpression` default (F1).

**Assessment:** PASS. Comprehensive startup validation.

### 8. Domain and Wire-Model Coupling

Expression types own `MarshalJSON`. The manifest is a primary domain projection — the JSON is the contract for the React client, not an incidental transport layer. Only one JSON representation exists (the TypeScript discriminated union). A separate transport mapper would add indirection without value: the domain types are already stable, and the marshaling logic is co-located with the type definition.

The `Expression` field in `ManifestValidation` uses the `Expression` interface type. When serialized, Go's `json.Marshal` dispatches to the concrete type's `MarshalJSON` method. This is idiomatic and requires no intermediate DTO.

**Assessment:** PASS. Co-located serialization is proportional.

### 9. Frontend Contract Compatibility

**Wire format comparison:**

| Concept | Go JSON | TypeScript |
|---|---|---|
| Comparison | `{"kind":"comparison","operand":"first"\|"second","operator":"equal"\|...,"value":N}` | `{ kind: 'comparison'; operand: 'first' \| 'second'; operator: ComparisonOperator; value: number }` |
| AllOf | `{"kind":"allOf","expressions":[...]}` | `{ kind: 'allOf'; expressions: ValidationExpression[] }` |
| Operators | `equal, notEqual, greaterThan, greaterThanOrEqual, lessThan, lessThanOrEqual` | Identical |
| Validation wrapper | `{"message":"...","expression":{...}}` | `{ message: string; expression: ValidationExpression }` |

**Exact match.** No discrepancies.

**Edge case behavior:**

| Case | Go | TypeScript | Compatible? |
|---|---|---|---|
| Missing operand | `ErrOperandOutOfBounds` error | Rejected by `operands.length !== operation.arity` | Yes — different layers |
| NaN as operand | Comparison evaluates (NaN >= 0 → false) | `Number.isFinite(NaN)` → false, rejected early | Different layers, both correct |
| ±Inf as operand | Correct IEEE-754 comparison | `Number.isFinite(Infinity)` → false, rejected early | Different layers |
| Negative zero | `-0 >= 0` → true, `-0 < 0` → false | Same IEEE-754 behavior | Compatible |
| Empty allOf | Vacuously true | `[].every(callback)` → true | Compatible |

The behavioral differences (NaN/Infinity rejection) occur at different layers: TypeScript at UI validation entry, Go at domain rule evaluation. Both approaches are correct. The Go validation evaluates domain constraints (divisor ≠ 0, radicand ≥ 0); the TypeScript also adds input sanity checks.

**Assessment:** PASS. Exact wire compatibility. Edge case behavior is the intended design.

### 10. Test Quality

**Covered:**

- All seven domain operations validated through `Operation.Validate`
- All three constrained operations: division (valid + violation), square root (valid + violation), percentage (5 boundary cases)
- All four unconstrained operations: always pass
- First violation in declaration order
- All evaluation error cases: missing first/second operand, unsupported reference, unsupported operator, nil expression
- AllOf: all-pass, short-circuit, error propagation, vacuous truth
- Registry startup rejection: 8 distinct malformed-definition cases
- Deep immutability: source mutation, Find copy mutation, manifest mutation
- Manifest: field-by-field JSON structure, deterministic ordering, no executor leakage
- Sentinel error matching with `errors.Is` in all error tests

**Not covered (low value):**

- Individual operator branches in `ComparisonExpression.evaluate` are exercised implicitly through domain operation tests rather than isolated unit tests
- `copyExpression` default branch and `validateExpressionDefinition` default branch are unreachable under the sealed interface
- NaN/Infinity as operands (not a domain concern — handled at the application layer)

**Assessment:** PASS. Comprehensive, well-structured. Tests exercise the primary domain path (`Operation.Validate`) and all failure modes.

### 11. Scope Proportionality

The solution contains exactly:
- 2 expression types (`ComparisonExpression`, `AllOfExpression`)
- 1 unexported evaluate contract
- 1 recursive evaluator per expression type
- 1 centralized structural validator
- 1 deep-copy function
- 6 comparison operators
- 2 operand references

It is not a DSL, parser, rules engine, reflection-based evaluator, plugin system, or unnecessary type hierarchy. It sits in the intended sweet spot between hardcoded if-statements per operation and a generic expression framework.

**Assessment:** PASS. Proportional to the calculator domain.

### 12. REST Readiness

| Concern | Status |
|---|---|
| Operation manifest endpoint | `Registry.Manifest()` returns `Manifest` struct — ready to serialize |
| Calculation service | `operation.Validate()` → `operation.Execute()` — natural call sequence |
| Validation-error response | `*ValidationDefinition` with `Message` for user, `ID` for tracing |
| Evaluation-error response | Error wraps sentinel + context (operation ID, validation ID) |
| Frontend HTTP integration | JSON contract matches TypeScript discriminated union exactly |

The `ManifestValidation` struct excludes the backend `ID` — the manifest is consumer-facing and the ID is only needed for tracing violations server-side. The frontend `ValidationDefinition` type has no `id` field. Compatible.

**Assessment:** PASS. Ready for REST contract publication.

---

## Architectural Confidence

| Dimension | Score (1–5) |
|---|---|
| Domain cohesion | 5 |
| Go idiomaticity | 5 |
| Runtime correctness | 5 |
| Startup validation | 5 |
| Immutability | 4 |
| Frontend contract compatibility | 5 |
| Scope proportionality | 5 |
| REST readiness | 5 |

---

## Verdict: PASS WITH OBSERVATIONS

**1 MINOR finding (F1)** — `copyExpression` default branch silently returns shared reference rather than failing loudly. Recommended change: panic with a descriptive message in the unreachable default branch.

**2 NOTES (N1, N2)** — unreachable default branch in `validateExpressionDefinition`, and previously-documented manifest field extras.

---

## Closing Assessment

If I inherited this validation model as a Staff Engineer, I would publish it as the contract used by the REST API and React client. The design is cohesive: `Operation.Validate` is the single authoritative entry point, expression types own their evaluation through the sealed interface, and the manifest projection is an exact match for the TypeScript discriminated union. Deep immutability is proven for every access path. Startup validation catches malformed definitions before they reach a user request. The implementation is proportional — two expression types, no generic engine, no reflection. The one concern (F1) is a dead-code safety gap in `copyExpression` that a one-line change would close. I would apply that change and ship.
