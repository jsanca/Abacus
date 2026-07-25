ABA-004-FIX-001 — Frontend Review Remediation
Field	Value
Status	Ready
Owner	Elito
Role	Frontend Developer · Knowledge Curator
Target	20–30 minutes
Hard stop	45 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Resolve the findings raised in ABA-004R — Frontend Architecture & UX Quality Review while preserving the current proportional architecture.

The task must correct functional and lifecycle issues without decomposing the calculator into unnecessary components or introducing speculative abstractions.

Required skills

Apply:

.claude/skills/process/osk-execution-timebox/SKILL.md
.claude/skills/process/osk-engineering-reporting/SKILL.md
.claude/skills/architecture/osk-architecture-review/SKILL.md

Follow:

CLAUDE.md
AGENTS.md
Review input

Use the review findings from:

docs/agents/reviews/ABA-004R-frontend-architecture-ux-review.md

or the repository’s actual review path if different.

Scope

Resolve:

AAR-001
AAR-002
AAR-003
AAR-004
AAR-005
AAR-006

Do not modify unrelated frontend behavior.

Required changes
1. Manifest loading Retry

When loading calculator capabilities fails:

display the service error;
render an accessible Retry action;
invoke the manifest-loading function again;
display the loading state during retry;
clear stale service errors when retry begins;
initialize the default operation when retry succeeds.

Avoid page reload as the recovery mechanism.

Prefer extracting a small function such as:

async function loadManifest(): Promise<void>

when it reduces duplication between the initial effect and Retry action.

Do not introduce automatic retry loops or retry counters unless required.

2. Clear resets the default operation

Update Clear behavior so it resets:

first operand;
second operand;
expression;
result;
validation and service messages;
calculator state;
selected operation to manifest.defaultOperationId;
focus to Operand 1.

The final behavior must match UC-001 A4.

3. Stabilize global keyboard listener

Correct the global keyboard-shortcut effect.

Do not merely add an empty dependency array if doing so captures stale state or handlers.

Choose the smallest correct design, for example:

stable callbacks with useCallback and complete effect dependencies;
a stable event handler using refs;
another React-idiomatic approach with explicit lifecycle.

Requirements:

one effective listener;
cleanup on unmount;
no stale manifest, result, selected operation, or calculator state;
no repeated unnecessary attach/detach on each render;
existing input-conflict protections remain intact.

Document the chosen dependency strategy briefly in the report.

4. Explicit error-state semantics

Replace:

calculatorState.includes('error')

with explicit semantic logic.

A small helper is acceptable:

function isErrorState(state: CalculatorState): boolean

Keep the state union exhaustive and readable.

5. Reduced-motion support

Add support for:

@media (prefers-reduced-motion: reduce)

Disable nonessential transitions and press transforms for interactive calculator controls.

Do not remove ordinary visible focus feedback.

6. Declarative focus management

Replace window.setTimeout(..., 0) focus scheduling with a more robust React lifecycle mechanism.

The required focus behaviors remain:

selecting a binary operation after a resolved result focuses Operand 2;
selecting a unary operation focuses Operand 1;
Clear focuses Operand 1;
continuation moves the prior result into Operand 1 before focusing the next appropriate input.

A focused useEffect driven by explicit focus intent is preferred over timing assumptions.

Avoid a generalized focus-management framework.

Component decomposition guidance

Do not split Calculator.tsx merely because it contains substantial JSX.

The current component remains acceptable as the state and orchestration owner for this scope.

Natural future extraction candidates include:

OperatorSelector
OperandEditor
ResultPanel
useKeyboardShortcuts
useOperandFocus

However, extract one only when it materially improves the implementation of the requested fixes.

Do not create:

one component per button;
generic design-system primitives;
prop-heavy presentational fragmentation;
a reducer or state-machine library;
a custom hook solely to move a few lines out of the file.

Prefer cohesive responsibility boundaries over file-size reduction.

Tests

Add or update behavior-focused tests for:

Manifest Retry
initial manifest request fails;
Retry action is visible;
Retry triggers another request;
successful retry renders all operations and the default operation.
Clear default operation
select a non-default operation, preferably Square Root;
enter values or obtain a result;
activate Clear;
verify default operation is restored;
verify binary layout is restored when the default is binary;
verify focus returns to Operand 1.
Keyboard listener

Test user-observable behavior rather than listener internals:

shortcuts continue working after multiple state transitions and rerenders;
a shortcut triggers one operation change;
no duplicate submission or duplicated state transition occurs.
Focus behavior
continuing from a binary result focuses Operand 2;
selecting Square Root focuses Operand 1;
Clear focuses Operand 1.
Error semantics
validation and service errors use alert semantics;
normal status messages do not.

The reduced-motion rule may be validated by static stylesheet inspection if browser emulation is disproportionate for this slice.

Verification

Run and report:

cd client
npm run lint
npm run test
npm run coverage
npm run build
npm run verify

npm run verify must complete and its final exit status must be recorded explicitly.

Also inspect the calculator manually at:

desktop width;
mobile width;
manifest failure and retry;
Clear from Square Root;
result continuation;
keyboard shortcuts.
Documentation

Update:

docs/engineering/ROADMAP.md
docs/engineering/ENGINEERING_LOG.md

Add technical debt only if a review finding is consciously deferred.

Do not create an ADR for localized React remediation.

Deliverables

Expected changes include:

client/src/features/calculator/components/Calculator.tsx
client/src/features/calculator/components/Calculator.css
client/src/features/calculator/test/Calculator.test.tsx

Additional files are allowed only when they create a clear cohesive boundary.

Report

Create:

docs/agents/reports/ABA-004-FIX-001-frontend-review-remediation.md

Use the repository’s actual standard reports path if different.

The report must include:

each AAR finding and its resolution;
keyboard-effect dependency strategy;
focus-management strategy;
tests added or changed;
exact command results;
explicit npm run verify result;
manual checks;
files modified;
remaining observations;
no-commit confirmation.

If the hard stop is reached, stop and create:

docs/agents/checkpoints/CHECKPOINT-ABA-004-FIX-001.md
Constraints
Do not integrate the Go backend.
Do not add Playwright.
Do not rewrite the visual design.
Do not introduce Redux, Zustand, React Query, or state-machine libraries.
Do not broadly refactor the calculator.
Do not commit.
Do not turn the review remediation into component atomization.
Definition of Done
Manifest failure has a functional Retry path.
Clear restores the default operation.
Keyboard listener lifecycle is stable and free of stale state.
Error states are identified explicitly.
Reduced-motion preferences are respected.
Focus management no longer depends on setTimeout.
New behavior tests pass.
Lint, test, coverage, build, and npm run verify complete successfully.
Documentation and the implementation report are current.
The frontend remains proportional and ready for re-review.
