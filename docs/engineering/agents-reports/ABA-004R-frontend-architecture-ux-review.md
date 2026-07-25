# ABA-004R — Frontend Architecture & UX Quality Review

**Status:** Complete
**Date:** 2026-07-24
**Reviewer:** Deep Pro (Architecture Reviewer)
**Reviewed artifact:** ABA-004 — Responsive React Frontend against mocked calculator API
**References:** `ABA-004` task spec, `UC-001` use case, `AGENTS.md`, implementation source under `client/src/`

## Verdict: CHANGES REQUIRED (3 MAJOR findings)

Three MAJOR findings must be resolved before the frontend architecture is approved:
- No Retry action for manifest loading failure (use-case violation).
- Clear does not reset to the default operation (use-case violation).
- Keyboard effect re-attaches on every render (correctness risk).

## 1. Scope

This review covers the frontend slice delivered by ABA-004:

| Layer | What | Scope |
|-------|------|-------|
| Entry / wiring | `main.tsx`, `App.tsx` | Included |
| UI component | `Calculator.tsx`, `Calculator.css` | Included |
| Mock API boundary | `mockCalculatorApi.ts`, `types.ts` | Included |
| Validation model | `validation.ts` | Included |
| Tests | `Calculator.test.tsx` | Included |
| Deferred | Go backend, Docker, financial precision | Excluded |

## 2. Architecture Assessment

### 2.1 Execution-Path Trace

```
main.tsx (StrictMode, createRoot)
  └─ App.tsx
       └─ Calculator.tsx
            ├─ useEffect → getOperations() → mockCalculatorApi.ts → OperationManifest
            │                                                              ↑
            │                                                    types.ts (contract)
            ├─ selectOperation() → state transition (selectedOperationId, calculatorState)
            ├─ controlled inputs (firstOperandText, secondOperandText)
            ├─ submitCalculation()
            │    ├─ validateOperation() → validation.ts (generic expression evaluator)
            │    └─ calculate() → mockCalculatorApi.ts → execute() → CalculationResponse
            └─ render tree
                 ├─ operator buttons (manifest-driven)
                 ├─ operand inputs (arity-driven)
                 ├─ status/error message (role=alert / aria-live=polite)
                 └─ Clear + Calculate actions
```

### 2.2 Positive Findings

**P1 — Clean feature-based structure.** The calculator lives under `features/calculator/` with `api/`, `model/`, `components/`, `test/` subdirectories. Each directory has a single clear responsibility. This pattern scales to additional features without cross-contamination.

**P2 — Well-defined API contract boundary.** `types.ts` defines `OperationManifest`, `OperationDefinition`, `OperandDefinition`, `ValidationDefinition`, `CalculationRequest`, `CalculationResponse`, and `ApiError` as a coherent contract. The mock API in `mockCalculatorApi.ts` implements the same type signature the future Go endpoints will expose. When the backend is ready, replacing the adapter is a one-file drop-in.

**P3 — Truly generic validation interpreter.** `validation.ts:3-21` evaluates manifest-provided comparison expressions (`equal`, `notEqual`, `greaterThan`, `lessThan`, `greaterThanOrEqual`, `lessThanOrEqual`) and `allOf` conjunctions. The code contains zero operation-specific branches; division-by-zero is enforced via the manifest's `{ operand: 'second', operator: 'notEqual', value: 0 }` expression. This satisfies the critical requirement that the frontend must not contain domain-specific validation rules.

**P4 — Manifest-driven rendering.** All seven operations are rendered from `OperationManifest.operations`. Operands are rendered from `OperationDefinition.operands`. Arity controls input count, and suffix metadata (percentage `%`) is applied generically. No operation is hardcoded in the render tree.

**P5 — Proper adapter isolation.** The mock API exposes `getOperations()` and `calculate()` as async functions. The component consumes them through a standard `try/catch` boundary. The `configureMockScenario()` escape hatch is exported but only used in tests — the production UI has no debug toggles.

**P6 — Accessible semantics present.** The implementation includes:
- `fieldset`/`legend` grouping for operators
- `aria-pressed` on operator buttons
- `aria-live="polite"` on the result panel
- `role="alert"` on error messages
- `aria-label` on operator buttons (human-readable names, not just symbols)
- `aria-describedby` linking inputs to the status message
- Visible focus styles (`:focus-visible` with 3px outline)

**P7 — Responsive layout.** Desktop (7-column operator grid) and mobile (4-column wrap, 560px breakpoint) are both handled. No horizontal scrolling at mobile widths. Touch targets meet minimum sizing (48px+).

## 3. Findings

### AAR-001 — No Retry action for manifest loading failure

**Severity:** MAJOR
**Location:** `Calculator.tsx:31-35`, `Calculator.tsx:104`
**Finding:** When `getOperations()` rejects (manifest failure), the component transitions to `service-error` and renders only a message in a card: `<section className="calculator-card loading-card">...</section>`. There is no button or mechanism to retry loading the manifest.

**Why it matters:** UC-001 F1 explicitly requires: "A Retry action is available." The user is left with no recovery path — they must refresh the page. This violates the use-case contract.

**Recommendation:** Render a Retry button alongside the error message that calls `getOperations()` again. Consider a small `retryCount` to prevent infinite retry loops, or keep it simple and let the user control retries.

### AAR-002 — Clear does not reset to the default operation

**Severity:** MAJOR
**Location:** `Calculator.tsx:64-72`
**Finding:** The `clearCalculator()` function resets operand texts, expression, result, and calculator state, but it does not reset `selectedOperationId` to the manifest's `defaultOperationId`. The previously selected operation persists.

**Why it matters:** UC-001 A4 specifies: "The initial or default operation is selected." After Clear, the calculator is logically restarting, and keeping the previous operation selected is inconsistent with the accepted interaction model. If the user was on Square Root (unary), cleared, and the default is Addition (binary), the UI would show stale unary operand layout.

**Recommendation:** Add `setSelectedOperationId(manifest?.defaultOperationId ?? '')` after the existing `setCalculatorState('idle')` call, guarded by a `manifest` check (which is always defined when `clearCalculator` can be called).

### AAR-003 — Keyboard shortcut effect has no dependency array

**Severity:** MAJOR
**Location:** `Calculator.tsx:38-47`
**Finding:** The `useEffect` that registers the global `keydown` listener has no dependency array (not even `[]`). This means the effect runs after every render: it removes the previous listener and re-adds an identical one. In React 18 Strict Mode with dev double-rendering, this is wasteful; in concurrent rendering, repeated attach/detach could cause subtle event-ordering issues.

```ts
useEffect(() => {
  const handleShortcut = (event: globalThis.KeyboardEvent) => { ... };
  window.addEventListener('keydown', handleShortcut);
  return () => window.removeEventListener('keydown', handleShortcut);
});  // <-- missing dependency array
```

**Why it matters:** While the cleanup prevents duplicate listeners, the constant re-registration is a correctness risk. If a future change adds state-dependent logic inside `handleShortcut`, stale closures would appear. The React exhaustivity lint rule correctly flags this.

**Recommendation:** Add `[]` as the dependency array since the effect depends on nothing that changes between renders. Wrap `selectOperation` and `clearCalculator` in `useCallback` if they are needed inside the effect (or use refs for stable references).

### AAR-004 — Fragile error-state detection via string inclusion

**Severity:** MINOR
**Location:** `Calculator.tsx:121`
**Finding:** The error class and `role="alert"` are gated on `calculatorState.includes('error')`. This works for the current states (`validation-error`, `service-error`) but would silently match any future state whose name contains "error" (e.g., `timeout-error`, `rejected-by-backend`).

**Why it matters:** Low immediate risk, but string-matching on state names couples rendering to naming convention rather than explicit semantics. If a non-error state were later added with "error" in its name, it would incorrectly trigger alert semantics.

**Recommendation:** Use explicit comparison: `calculatorState === 'validation-error' || calculatorState === 'service-error'`, or define a helper set of error states.

### AAR-005 — No `prefers-reduced-motion` media query

**Severity:** MINOR
**Location:** `Calculator.css:14` (button transitions)
**Finding:** The CSS applies `transition: background .15s, color .15s, transform .15s` to operator buttons and `transform: translateY(1px)` on `:active` for operator, calculate, and clear buttons. There is no `@media (prefers-reduced-motion: reduce)` to disable these animations.

**Why it matters:** WCAG SC 2.3.3 (Animation from Interactions) requires respecting system-level motion preferences. Users who have enabled reduced motion will still experience the transition and press-down translate effects.

**Recommendation:** Add:
```css
@media (prefers-reduced-motion: reduce) {
  .operator-button, .calculate-button, .clear-button { transition: none; }
  .operator-button:active, .calculate-button:active, .clear-button:active { transform: none; }
}
```

### AAR-006 — `setTimeout(fn, 0)` focus management is fragile

**Severity:** MINOR
**Location:** `Calculator.tsx:61,71`
**Finding:** Post-state-update focus is managed via `window.setTimeout(..., 0)`. This relies on React's synchronous flush timing and can break if React changes its scheduling. React 18's automatic batching already ensures state updates are flushed before the timeout fires, but the pattern is a race-condition surface.

**Why it matters:** Low immediate risk (works correctly today), but the pattern is a known anti-pattern. A more robust approach would use `useEffect` that watches the state triggering the focus change, or the `ref`-based autoFocus pattern.

**Recommendation:** Replace with a `useEffect` that watches `calculatorState` and `selectedOperationId` and focuses the appropriate input when conditions are met. This makes the intent declarative and eliminates the timing dependency.

## 4. Architecture Properties Checklist

| Property | Status | Notes |
|----------|--------|-------|
| Dependency direction | PASS | Component → API, Component → Validation. No inversions. |
| Adapter isolation | PASS | Mock API is a clean swappable boundary behind async functions. |
| API contract | PASS | Types form a complete, self-describing contract. |
| Validation genericity | PASS | Zero operation-specific logic in validation or UI. |
| Statelessness (API) | PASS | Mock API functions are stateless aside from the test-scenario flag. |
| Error propagation | PASS | Errors flow through a typed `ApiError` envelope; UI distinguishes service vs validation. |
| Test concern coverage | PASS (with gaps) | 10 tests cover primary and alternative flows. No test for manifest-failure retry (the feature doesn't exist). No test for default-operation reset on Clear. |
| Accessibility | PASS (MINOR gap) | Semantic HTML, ARIA roles, focus styles. Missing `prefers-reduced-motion`. |
| Responsiveness | PASS | Desktop and mobile layouts at 560px breakpoint; 7 operators visible at both widths. |

## 5. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Use-case non-compliance (retry, clear default) | Medium — users hit these flows | Medium — blocked recovery | Implement AAR-001 and AAR-002 |
| Keyboard listener leak in future changes | Low — code works correctly today | High — subtle stale-closure bugs | Implement AAR-003 |
| State-transition race during calculating | Low — button disabled during calculating | Low — no data corruption | Accept; resolve if observed |
| JavaScript float precision | High (inherent) | Low — out of scope per task spec | Accepted limitation |

## 6. Recommendation

Resolve the three MAJOR findings (AAR-001, AAR-002, AAR-003), then re-review. The three MINOR findings (AAR-004, AAR-005, AAR-006) should be addressed but do not block approval.

The underlying architecture — feature-based structure, typed contract boundary, generic validation interpreter, mock adapter isolation, and manifest-driven rendering — is sound and well-aligned with the intended design.

## 7. Positive Findings (Architecture)

1. **Feature directory convention** (`api/`, `model/`, `components/`, `test/`) establishes a repeatable pattern.
2. **Contract-first types** (`types.ts`) define the full API surface before implementation.
3. **Generic validation interpreter** (`validation.ts`) correctly abstracts operation-specific rules into declarative expressions.
4. **Async API boundary** (`mockCalculatorApi.ts`) mirrors the future HTTP contract; replacing it requires no component changes.
5. **No framework sprawl** — no Redux, Zustand, React Query, or UI library. Pure React state with controlled inputs.

## 8. Deferred Findings

- **Financial/decimal precision** — explicitly out of scope (task paragraph 432, report paragraph 42).
- **Manifest persistence across page refresh** — out of scope; the backend will own session concerns.
- **CSS methodology (plain CSS vs CSS Modules)** — the task permitted plain CSS; the choice is sufficient for current scope.
