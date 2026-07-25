# ABA-005R — Operation Registry Architecture & Domain Review

**Verdict:** PASS WITH OBSERVATIONS

**Reviewer:** Deep Pro
**Date:** 2026-07-25
**Repository:** /Users/jonathan/code/Abacus

---

## 1. Domain Model

**Assessment: PASS**

The domain model is clean and proportional. `Operation` pairs a frontend-safe `OperationDefinition` with a trusted `ExecuteFunc` — no transport, HTTP, or configuration concerns leaked.

`Arity` is a named int type with constants `Unary` (1) and `Binary` (2). `OperandDefinition` carries presentation metadata (Label, Placeholder, Suffix) — this is intentional: the registry IS the manifest source of truth, so presentation metadata belongs here.

`ExecuteFunc` is the simplest idiomatic representation: a function type, not an interface hierarchy. Correct for stateless arithmetic.

`OperandFirst` and `OperandSecond` are stable string constants forming the client contract. No discovery needed.

## 2. Registry Design

**Assessment: PASS**

`NewRegistry` validates all invariants at construction: non-empty operations, non-empty default ID, no duplicate IDs, valid arity, operand count matches arity, no duplicate operand IDs, non-nil executor. Returns `(*Registry, error)` — there is no path to a partially-constructed registry.

**Immutability**: Struct fields are unexported. No method mutates the registry. All accessors return deep copies:

- Constructor deep-copies every `Operation` and its operand slice via `copyOperation`.
- `Find()` returns a deep copy.
- `Operations()` returns a new slice of deep copies.
- `Manifest()` allocates new `ManifestOperation` and `ManifestOperand` slices.

**No synchronization**: No mutex or atomic. Concurrent reads are safe because no writes occur after construction returns.

**Defensive copies verified by tests**: `TestRegistry_Operations_returnsDefensiveCopy`, `TestRegistry_immutableAfterSourceSliceMutation`, `TestManifest_callerCannotMutateRegistryOperands`.

**Ordering**: A map for O(1) lookup + an ordered `[]Operation` slice for deterministic manifest projection. Standard Go pattern. The order is product-defined in `DefaultOperations()` and preserved by index-mirroring in `NewRegistry`.

## 3. ExecuteFunc

**Assessment: PASS**

```go
type ExecuteFunc func(operands []float64) (float64, error)
```

Every executor guards against wrong operand count and returns `ErrWrongOperandCount` wrapped with `%w`. No panics possible. Edge cases (division by zero → ±Inf, sqrt of negative → NaN) are deferred to ABA-006 validation, which is the documented plan.

No unnecessary interfaces. No speculation.

## 4. Manifest Projection

**Assessment: PASS**

`Registry.Manifest()` produces a `Manifest` struct with new allocations at every level — the caller cannot mutate registry internals through the manifest.

`json.RawMessage("[]")` for the `validations` field: this serializes as `"validations":[]` in JSON, matching the TypeScript `ValidationDefinition[]` contract. The frontend can iterate the empty array without special-casing. This is an acceptable transitional decision and should be replaced with real validations in ABA-006.

Executor is excluded — `ManifestOperation` has no `Execute` field. JSON tags use camelCase consistently.

Deterministic order verified: operations appear in the same order as `DefaultOperations()`.

## 5. Frontend Compatibility

**Assessment: PASS WITH OBSERVATIONS**

| Field | Go (JSON key) | TypeScript | Match |
|-------|--------------|------------|-------|
| `DefaultOperationID` | `defaultOperationId` | `defaultOperationId` | ✓ |
| `Version` | `version` | **missing** | Extra, ignored |
| `Operations[].ID` | `id` | `id` | ✓ |
| `Operations[].Name` | `name` | `name` | ✓ |
| `Operations[].Symbol` | `symbol` | `symbol` | ✓ |
| `Operations[].Shortcut` | `shortcut` | **missing** | Extra, ignored |
| `Operations[].Arity` | `arity` | `arity` | ✓ |
| `Operations[].Operands` | `operands` | `operands` | ✓ |
| `Operations[].Validations` | `validations` | `validations` | ✓ |
| `Operands[].ID` | `id` | `id` | ✓ |
| `Operands[].Label` | `label` | `label` | ✓ |
| `Operands[].Placeholder` | `placeholder` | `placeholder` | ✓ |
| `Operands[].Suffix` | `suffix` | `suffix?` | ✓ |

All seven operation IDs, operand labels, placeholders, suffixes, symbols, and the default operation match between Go and the mock.

**Integration risk**: The Go manifest includes `version` and `shortcut` fields that the TypeScript `OperationManifest` / `OperationDefinition` interfaces do not define. These are forward-compatible (extra JSON keys are silently ignored). However, `shortcut` should be consumed by the frontend during integration to eliminate the hardcoded `shortcutOperationIds` map in `Calculator.tsx:12`. This is a known integration item, not a contract break.

## 6. Public API Surface

**Assessment: PASS**

Every exported symbol has a present or near-future consumer:

| Symbol | Consumer |
|--------|----------|
| `Registry` | `main.go`, tests |
| `NewRegistry` | composition root, tests |
| `NewDefaultRegistry` | `main.go`, tests |
| `DefaultOperations` | composition root, tests |
| `Registry.Find` | future API handler, tests |
| `Registry.Operations` | tests |
| `Registry.Manifest` | future API handler, tests |
| `Registry.DefaultOperationID` | tests |
| `Operation`, `OperationDefinition`, `OperandDefinition` | callers constructing operations, tests |
| `ExecuteFunc` | callers assigning executors |
| `Arity`, `Unary`, `Binary` | callers constructing operations, tests |
| `OperandFirst`, `OperandSecond` | callers constructing operand definitions |
| `Manifest`, `ManifestOperation`, `ManifestOperand` | transport layer serialization |
| `ManifestVersion` | manifest generation |
| `ErrWrongOperandCount` | callers inspecting errors, tests |

No speculative APIs. Nothing package-private that shouldn't be.

## 7. Go Idioms

**Assessment: PASS**

- Single `calculator` package — no over-fragmentation.
- Constructor returns `(*Registry, error)` — no zero-value trickery.
- Receiver named `registry` — matches type, not abbreviated.
- Pointer receiver on `Registry` — correct since struct contains map + slice.
- `%w` for error wrapping; `errors.Is` used in tests.
- `make(T, len)` for pre-allocated maps and slices.
- Table-driven tests with `t.Run`.
- Every exported symbol has a doc comment (verified by `revive` linter).
- No interfaces without consumers, no getter/setter methods.
- No Java influence detected.

## 8. Testing

**Assessment: PASS**

**Coverage**: `internal/calculator` at 96.7%.

**Construction tests** (`TestNewRegistry_invalidInputs`): 10 sub-cases covering every invariant — empty operations, empty default, unknown default, duplicates, empty ID, empty name, unsupported arity, operand count mismatch, duplicate operand IDs, nil executor.

**Behavior tests**: Finds all 7 operations by ID, unknown ID returns false, default operation ID verified, deterministic ordering verified.

**Defensive-copy tests**: Three tests verify that:
1. Mutating `Operations()` return value doesn't affect registry.
2. Mutating source slice after construction doesn't affect registry.
3. Mutating manifest does not affect registry.

**Executor tests**: All 7 operations with correct results. Each operation tested with wrong arity (including empty slice). Uses `errors.Is` for sentinel error checking.

**Manifest tests**: Version, default operation, count (7), deterministic order, JSON field names, unary square-root (single operand, correct ID), percentage suffix, validations as empty array.

Tests verify behavior and invariants, not implementation details. No mocking — tests use the real `NewDefaultRegistry()`.

## 9. Future Readiness

**Assessment: PASS**

| Planned work | Impact on registry |
|-------------|-------------------|
| ABA-006 Validation Model | Populate `Validations` field — no structural change |
| ABA-007 REST API | `Manifest()` → GET /manifest; `Find()` → lookup for calculation — no new registry methods needed |
| Frontend integration | JSON contract is compatible — no schema changes |

The registry requires no architectural changes for any planned work.

---

## Validation Results

| Command | Result |
|---------|--------|
| `make fmt` | PASS |
| `make lint` | PASS (0 issues) |
| `make test` | PASS (all packages) |
| `make coverage` | PASS (83.9% total, 96.7% calculator) |
| `make build` | PASS |
| `make verify` | PASS (fmt → lint → race → build) |
| `go test -race ./...` | PASS (no races) |
| `docker compose build server` | PASS |

---

## Findings

### MINOR

**F1: Manifest includes `version` and `shortcut` fields not yet in frontend TypeScript contract**

- **Severity:** MINOR
- **Location:** `server/internal/calculator/manifest.go:17,30`
- **Evidence:** Go `Manifest` has `Version string \`json:"version"\`` and `ManifestOperation` has `Shortcut string \`json:"shortcut"\``. The TypeScript `OperationManifest` and `OperationDefinition` interfaces do not define these fields.
- **Why it matters:** Extra JSON keys are forward-compatible at runtime but represent an incomplete contract alignment. The `shortcut` field in particular duplicates the hardcoded `shortcutOperationIds` map in `Calculator.tsx:12`. During frontend integration, the client should consume `shortcut` from the manifest instead of maintaining its own mapping.
- **Recommendation:** Defer to the frontend-backend integration task. Align the TypeScript interfaces with the Go manifest (add `version` and `shortcut` fields) and remove the hardcoded shortcut map from Calculator.tsx.

### NOTE

**N2: `Operations()` and `Find()` return `Operation` with an exposed `Execute` field**

- **Severity:** NOTE
- **Location:** `server/internal/calculator/registry.go:60,76`
- **Evidence:** `Operations()` returns `[]Operation` and `Find()` returns `(Operation, bool)`. The `Operation` type includes `Execute ExecuteFunc`, meaning callers could invoke execution directly without any mediation layer.
- **Why it matters:** Currently harmless — tests are the only callers, and the future REST handler (ABA-007) will be the natural mediation layer. No action needed now.
- **Recommendation:** Ensure the ABA-007 REST handler mediates all execution. Do not expose `Execute` through the transport boundary.

**N3: `DefaultOperations()` does not deep-copy returned elements**

- **Severity:** NOTE
- **Location:** `server/internal/calculator/operations.go:10`
- **Evidence:** `DefaultOperations()` returns a fresh `[]Operation` each call (via composite literal), but does not deep-copy individual `Operation` structs' operand slices. A caller could mutate the returned operand slice before passing it to `NewRegistry`.
- **Why it matters:** Harmless — `NewRegistry` deep-copies via `copyOperation` on construction, so pre-construction mutations cannot affect the registry. No action needed.
- **Recommendation:** No change required. The defensive layer at `NewRegistry` is the correct single point of protection.

---

## Architectural Confidence

| Aspect | Score | Notes |
|--------|-------|-------|
| Domain model | 5/5 | Clean, no leaks, proportional to scope |
| Go idiomaticity | 5/5 | Simple functions, concrete types, proper error wrapping, documented |
| Public API design | 5/5 | Every export has a consumer, nothing speculative |
| Immutability | 5/5 | Provably immutable after construction, deep-copy on all access paths |
| Test quality | 5/5 | Behavior-focused, covers invariants, 96.7% coverage |
| Frontend integration readiness | 4/5 | JSON contract compatible; `version` and `shortcut` fields need frontend alignment during integration |
| Long-term maintainability | 5/5 | Simple abstractions, supports all planned extensions without rearchitecture |

---

## Definition of Done

- [x] Registry architecture reviewed
- [x] Immutability strategy validated
- [x] Defensive-copy implementation reviewed
- [x] Manifest projection reviewed
- [x] Frontend compatibility evaluated
- [x] Public API surface evaluated
- [x] Go idioms reviewed
- [x] Tests reviewed
- [x] Validation commands executed
- [x] Future readiness assessed
- [x] Final verdict issued
- [x] No code changed
- [x] No commit created
