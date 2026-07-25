# ABA-005 — Immutable Operation Registry & Manifest Projection

| Field | Value |
|---|---|
| Task | ABA-005 — Immutable Operation Registry & Manifest Projection |
| Status | Complete |
| Author | Clio |
| Date | 2026-07-25 |
| Report type | Implementation |

## Summary

Delivered the immutable operation registry as the backend source of truth for Abacus calculator capabilities. The registry holds seven operations in deterministic order, validates its own invariants at construction time, is safe for concurrent reads, and projects a frontend-safe manifest. All quality gates pass.

## Package Structure

```
server/internal/calculator/
├── operation.go        — ExecuteFunc, Arity, OperandDefinition, OperationDefinition, Operation, copyOperation
├── operations.go       — DefaultOperations, NewDefaultRegistry, seven operation builders, seven executor functions
├── registry.go         — Registry struct, NewRegistry, Find, Operations, DefaultOperationID, validateRegistryOperation
├── manifest.go         — Manifest, ManifestOperation, ManifestOperand types, Registry.Manifest, projectOperation
├── registry_test.go    — construction, behavior, manifest, and defensive-copy tests
└── operations_test.go  — execution correctness and wrong-arity tests
```

## Registry Invariants

Validated at construction time, before the HTTP server starts:

- At least one operation is registered.
- Default operation ID is non-empty.
- Default operation ID exists in the registry.
- Each operation ID is non-empty and unique.
- Each operation has a non-empty name, symbol, and shortcut.
- Arity is `1` or `2`; no other values accepted.
- Operand count equals arity.
- All operand IDs are non-empty and unique within the operation.
- Executor function is not nil.

## Immutability Strategy

The constructor (`NewRegistry`) allocates new slices for `operationsByID` and `operations`. Each `Operation` value is deep-copied via `copyOperation`, which allocates an independent backing array for the `Operands` slice. After construction no write path exists: there is no `Register`, `Remove`, or mutation method. The struct fields are unexported.

`Find` and `Operations` return deep copies for the same reason: callers receive values whose `Operands` slices do not alias registry-owned storage.

## Concurrency Reasoning

No writes occur after successful construction. Go's memory model guarantees that a value fully initialized before a goroutine is started is safely readable by that goroutine and any subsequent readers without synchronization. No mutex, `sync.Map`, or atomic values are used or needed.

## Lookup Method: `Find`

`Find(operationID string) (Operation, bool)` was chosen over `Get` and `Lookup`:
- `Get` reads as a map accessor and implies only value retrieval.
- `Lookup` is common in networking/DNS contexts and feels slightly off for a domain registry.
- `Find` reads naturally at the call site (`registry.Find("addition")`) and signals intentional search with an existence result, consistent with the comma-ok idiom.

## Manifest Projection

`Registry.Manifest()` returns a `Manifest` value containing:
- `version: "1"` — stable schema version.
- `defaultOperationId: "addition"` — using the exact JSON key from the TypeScript contract.
- `operations: []ManifestOperation` — ordered to match the registry's declared product order.

Each `ManifestOperation` is allocated independently from registry-owned storage. `ManifestOperand` fields are value copies. The `Execute` function is excluded from the projection.

`validations` is always `[]` (empty JSON array); see Contract Mismatches below.

## Executor Behavior

All seven executors guard against wrong operand count (`ErrWrongOperandCount`). Domain validations (divisor ≠ 0, radicand ≥ 0, percentage in 0–100) are deliberately absent; they belong to the declarative validation engine in ABA-006.

For mathematically invalid inputs that pass the arity check (e.g. `division(5, 0)`, `sqrt(-1)`), Go's `float64` semantics apply (`+Inf`, `NaN`). This is documented in each executor and is the accepted interim behavior.

Percentage semantics: `base × percentage / 100` — `executePercentage(200, 15) = 30`. ✓

## Application Wiring

The registry is constructed and discarded in `run()` for fail-fast startup validation:
```go
if _, err := calculator.NewDefaultRegistry(); err != nil {
    return fmt.Errorf("invalid operation registry: %w", err)
}
```

The registry is not threaded as a named dependency because no HTTP endpoint consumes it yet. Application wiring begins in ABA-007 when the operations and calculation endpoints are introduced.

## Client-Contract Reconciliation

| Field | TypeScript OperationManifest | Go Manifest | Status |
|---|---|---|---|
| `operations` | `OperationDefinition[]` | `[]ManifestOperation` | ✓ Compatible |
| `defaultOperationId` | `string` | `string` | ✓ Compatible |
| `version` | not present | `"1"` | Extra field; TS client ignores unknown fields |
| `shortcut` per operation | not present | `string` | Extra field; TS client ignores unknown fields |
| `validations` per operation | `ValidationDefinition[]` | `[]` (always empty) | See mismatch below |

## Contract Mismatches Found

**validations always empty:** The TypeScript client uses `validations` from the manifest to run client-side validation (e.g. division divisor ≠ 0, square-root radicand ≥ 0). The Go manifest returns `[]` for all operations. As a result:
- Client-side validation is disabled until ABA-006 delivers the declarative validation engine.
- The happy path (valid inputs) is unaffected.
- This is documented in `manifest.go` and registered below in Technical Debt.

## Tests and Command Results

```
make fmt          → no format issues                              ✓
make lint         → 0 issues (golangci-lint)                     ✓
make test         → calculator 100%, all packages pass           ✓
make coverage     → calculator 96.7%                             ✓
make build        → ok                                           ✓
make verify       → all checks passed                            ✓
go test -race     → all pass, no data races                      ✓
docker compose build server → image rebuilt                      ✓
```

Coverage detail for `server/internal/calculator`:
- `validateRegistryOperation`: 87% (empty-operand-ID path; covered by operand count mismatch cases)
- All other functions: 100%

## Deferred Validation Behavior

Domain validations (divisor ≠ 0, square-root non-negative, percentage 0–100) are deferred to ABA-006. Executors return Go float64 semantics for invalid inputs when arity is satisfied.

The `Validations json.RawMessage` field in `ManifestOperation` is set to `json.RawMessage("[]")` — a correct empty JSON array that satisfies the TypeScript `ValidationDefinition[]` field without introducing a premature type.

## Files Created or Changed

| File | Change |
|---|---|
| `server/internal/calculator/operation.go` | Created: types, constants, `copyOperation` |
| `server/internal/calculator/operations.go` | Created: `DefaultOperations`, `NewDefaultRegistry`, 7 builders, 7 executors |
| `server/internal/calculator/registry.go` | Created: `Registry`, `NewRegistry`, `Find`, `Operations`, `DefaultOperationID`, `validateRegistryOperation` |
| `server/internal/calculator/manifest.go` | Created: `Manifest`, `ManifestOperation`, `ManifestOperand`, `Registry.Manifest`, `projectOperation` |
| `server/internal/calculator/registry_test.go` | Created: 14 test functions covering construction, behavior, manifest, defensive copies |
| `server/internal/calculator/operations_test.go` | Created: execution correctness (7 ops) and wrong-arity (7 ops) |
| `server/cmd/api/main.go` | Updated: fail-fast registry validation at startup |
| `docs/engineering/ROADMAP.md` | ABA-005 marked Complete |
| `docs/engineering/ENGINEERING_LOG.md` | ABA-005 entry added |
| `docs/knowledge/glossary/README.md` | Operation, Registry, Manifest, ExecuteFunc, Arity, ValidationDefinition defined |
| `docs/engineering/agents-reports/ABA-005-operation-registry.md` | This report |

## Accepted Technical Debt

TD-002: `validations` in the manifest is always `[]`. Client-side validation (divisor ≠ 0, etc.) is disabled until ABA-006 populates the declarative validation engine.

## No-commit Confirmation

No commits were created.
