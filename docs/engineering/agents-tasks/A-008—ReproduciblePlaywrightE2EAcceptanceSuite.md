ABA-008 — Reproducible Playwright E2E Acceptance Suite
Field	Value
Status	Ready
Owner	Elito
Role	Frontend Developer · QA Automation Engineer · Knowledge Curator
Target	45–60 minutes
Hard stop	75 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Implement a reproducible Playwright end-to-end acceptance suite for Abacus.

The suite must exercise the real integrated application:

Browser
  ↓
React client
  ↓
Nginx proxy
  ↓
Go REST API
  ↓
Operation Registry
  ↓
Validation
  ↓
Execution

The canonical E2E command must own the complete Docker lifecycle so the suite can run from a clean machine state without requiring manually started services.

The suite must also remain convenient for local development by optionally supporting execution against an already-running stack.

Required skills

Apply:

.claude/skills/process/osk-execution-timebox/SKILL.md
.claude/skills/process/osk-engineering-reporting/SKILL.md
.claude/skills/architecture/osk-architecture-review/SKILL.md

Follow:

CLAUDE.md
AGENTS.md
Required inputs

Review:

docs/knowledge/use-cases/UC-001-perform-calculation.md
docs/knowledge/acceptance-tests/AT-001-calculator-experience.md
docs/knowledge/api/API_CONTRACT.md

Also inspect:

client/
server/
docker-compose.yml
Makefile
README.md

The acceptance-test document is the functional recipe. Do not invent different product behavior from the implementation.

Environment strategy
Canonical execution mode

The official command must:

remove any existing Abacus containers and orphans;
build the current client and server images;
start both services;
wait until the backend and frontend are ready;
run Playwright against the real stack;
retain failure evidence;
stop and remove the stack even when tests fail;
return Playwright’s original exit status.

Conceptually:

docker compose down --remove-orphans
docker compose up --build -d
wait-for-stack
playwright test
docker compose down --remove-orphans

Do not implement readiness using an arbitrary long sleep.

Poll:

http://localhost:8080/health
http://localhost:3000/

with:

bounded retries;
short intervals;
meaningful failure output;
an overall timeout.
Development mode

Provide a second command that runs Playwright against an already-running stack without rebuilding or shutting it down.

For example:

npm run test:e2e:local

This mode is for rapid iteration only.

The canonical evidence command remains the clean-stack mode.

Suggested commands

Expose under client/package.json:

npm run test:e2e
npm run test:e2e:local
npm run test:e2e:headed
npm run test:e2e:report

At repository root expose:

make e2e
make e2e-local

Recommended semantics:

make e2e

owns Docker lifecycle.

make e2e-local

assumes services already run.

Do not make make verify automatically run full E2E unless the repository deliberately wants every verification to build Docker and launch browsers. Keep unit verification fast and E2E explicit.

Playwright setup

Install and configure:

@playwright/test;
Chromium at minimum;
one desktop project;
one mobile project.

Suggested files:

client/
├── playwright.config.ts
└── tests/
    └── e2e/
        ├── calculator.page.ts
        ├── calculator-happy-path.spec.ts
        ├── calculator-validation.spec.ts
        ├── calculator-continuity.spec.ts
        ├── calculator-keyboard.spec.ts
        ├── calculator-failures.spec.ts
        ├── calculator-responsive.spec.ts
        └── calculator-accessibility.spec.ts

Keep the number of files proportional. Consolidate when fewer files make the suite clearer.

Playwright configuration

Configure:

baseURL: http://localhost:3000
trace: retain-on-failure
screenshot: only-on-failure
video: retain-on-failure

Use retries only where justified by CI-like environmental variability. Do not mask flaky tests with excessive retries.

Recommended projects:

Desktop Chromium
Mobile Chromium — approximately 375 × 667 or a standard mobile device profile

No cross-browser matrix is required for this take-home unless time remains after the primary suite is stable.

Locator strategy

Prefer accessible locators:

page.getByRole('button', { name: 'Division' })
page.getByLabel('Dividend')
page.getByLabel('Divisor')
page.getByRole('button', { name: 'Calculate' })
page.getByRole('status')
page.getByRole('alert')

Avoid:

CSS implementation selectors;
generated classes;
DOM-position selectors;
brittle full-page snapshots.

Add data-testid only when semantic locators cannot express the user-facing concept reliably.

Page object

A small page object is appropriate.

Suggested shape:

class CalculatorPage {
  constructor(page: Page)

  async open(): Promise<void>
  async waitUntilReady(): Promise<void>
  async selectOperation(name: string): Promise<void>
  async enterFirstOperand(value: string): Promise<void>
  async enterSecondOperand(value: string): Promise<void>
  async calculate(): Promise<void>
  async clear(): Promise<void>

  firstOperand(): Locator
  secondOperand(): Locator
  result(): Locator
  status(): Locator
  error(): Locator
}

Do not hide entire acceptance scenarios behind methods such as:

verifyDivisionWorks()

The test should still tell the business story explicitly.

Required acceptance coverage
1. Capability loading

Verify:

the real GET /api/v1/operations request succeeds;
all seven operation controls render;
Addition is selected by default;
no hardcoded fallback silently replaces a failed manifest.
2. Seven real calculations

Use the real backend for:

Operation	Inputs	Expected result	Expected expression
Addition	10, 5	15	10 + 5 = 15
Subtraction	5, 10	-5	5 − 10 = -5
Multiplication	2.5, 4	10	2.5 × 4 = 10
Division	144, 12	12	144 ÷ 12 = 12
Exponentiation	2, 8	256	2 xʸ 8 = 256
Square Root	144	12	√144 = 12
Percentage	200, 15	30	15% of 200 = 30

The tests must confirm a real calculation request reached:

POST /api/v1/calculations

Do not mock these happy paths.

3. Declarative validation

Verify:

Division
10 ÷ 0

Expected:

client validation error;
no incorrect result;
operand values preserved;
calculate request is not sent when client-side projected validation prevents submission.
Square Root
√-1

Expected rejection.

Percentage

Verify:

0%
100%
-1%
101%

Expected:

0 and 100 accepted;
-1 and 101 rejected;
fixed % suffix remains visible.
4. Backend authority

Client-side validation prevents ordinary invalid requests, so add at least one API-level integration assertion outside the browser UI or through Playwright request context:

POST /api/v1/calculations
division [10, 0]

Expected:

HTTP 422
code = validation_failed

This proves the backend independently enforces the same rule.

Use Playwright’s request fixture or a small E2E API spec. Do not bypass the stack with internal Go calls.

5. Continuity
Binary to binary

Given:

144 ÷ 12 = 12

When Multiplication is selected:

Operand 1 becomes 12;
Operand 2 is empty;
Operand 2 has focus.

Then calculate:

12 × 3 = 36
Unary to binary

Given:

√144 = 12

When Addition is selected:

Operand 1 becomes 12;
Operand 2 is empty;
Operand 2 has focus.
Binary to unary

Given a completed binary calculation, when Square Root is selected:

result becomes Operand 1;
second operand is absent;
first operand has focus.

Use Playwright focus assertions:

await expect(locator).toBeFocused()
6. Clear

Verify button and Escape behavior:

operands cleared;
result cleared;
errors cleared;
default operation restored;
focus returns to Operand 1.
7. Keyboard

Verify:

Enter → calculate
Escape → clear
+ → Addition
- → Subtraction
* → Multiplication
/ → Division
^ → Exponentiation
r → Square Root
% → Percentage

Also verify operator shortcuts do not fire while the user is typing inside an operand field.

8. Loading and duplicate submission

Delay the calculation route using Playwright interception.

Verify:

calculating state is visible;
Calculate becomes disabled or otherwise protected;
repeated activation does not create duplicate requests;
exactly one calculation result appears.

Interception is allowed here because the purpose is to make latency observable, not replace backend behavior.

9. Calculation service failure

Intercept:

/api/v1/calculations

or temporarily stop the server only if deterministic orchestration is practical.

Verify:

service error appears;
operands remain;
operation remains;
no incorrect result appears;
user can retry after recovery.

Prefer route interception for this isolated infrastructure scenario because it keeps the suite deterministic.

10. Manifest failure and retry

Intercept the first manifest request to fail.

Verify:

service-unavailable state;
no hardcoded operation controls appear;
Retry is visible;
next request succeeds;
calculator becomes usable with all seven operations.
11. Responsive behavior
Desktop

Verify:

all seven operations visible;
no horizontal overflow;
calculator card fully usable.
Mobile

At approximately:

375 × 667

verify:

all seven operators visible;
operator grid wraps;
no horizontal overflow;
inputs and actions remain operable;
result remains visible;
square-root unary layout works;
percentage suffix remains visible.

Avoid exact pixel assertions.

12. Accessibility-observable behavior

Verify:

selected operator exposes aria-pressed;
result uses live/status semantics;
errors use alert semantics;
invalid inputs expose appropriate accessible state where implemented;
operator group has an accessible name;
all primary actions are keyboard reachable;
focus transitions match the expected flow.

Do not claim full WCAG compliance from Playwright alone.

Test independence

Every test must:

begin from a fresh page;
establish its own selected operation and input state;
not depend on execution order;
not share mutable browser state;
leave network interception cleaned up.

Do not run tests serially unless an explicit reason is documented.

Test evidence

Create:

docs/knowledge/testing/
├── TEST_STRATEGY.md
├── ACCEPTANCE_TEST_RESULTS.md
└── evidence/
    ├── desktop-calculator.png
    ├── mobile-calculator.png
    └── playwright-summary.png
TEST_STRATEGY.md

Explain:

Go unit/API tests;
React component tests;
Playwright E2E;
what each layer proves;
why happy-path E2E uses the real backend;
why route interception is limited to infrastructure failures and latency.
ACCEPTANCE_TEST_RESULTS.md

Include:

execution date;
command used;
environment;
browser projects;
passed/failed/skipped totals;
duration;
acceptance scenarios covered;
links to the evidence images;
known limitations;
final verdict.

Do not paste enormous console logs.

Screenshots

Produce clean, user-facing screenshots from Playwright:

desktop successful calculation;
mobile successful calculation.

Do not include browser DevTools, private local paths, or debugging clutter.

playwright-summary.png may show the generated report summary or terminal result, provided it is readable and contains no sensitive information.

Artifact policy

Do not commit ordinary generated artifacts such as:

playwright-report/
test-results/
videos/
traces/

Add them to .gitignore when needed.

Commit only deliberately curated evidence under:

docs/knowledge/testing/evidence/

Retain traces, screenshots and videos automatically on failure for local inspection, but keep them out of Git.

Docker orchestration

Implement a script under a proportional path such as:

scripts/run-e2e.sh

or:

client/scripts/run-e2e.sh

Responsibilities:

clean stack
build and start
wait for readiness
run Playwright
capture exit code
always tear down
return original exit code

Use a shell trap:

trap cleanup EXIT

or an equivalent robust mechanism.

The script must not leave containers running after success or failure.

Print useful diagnostics if readiness fails:

docker compose ps
docker compose logs server
docker compose logs client

Do not dump full logs on every successful run.

CI-readiness

No hosted CI pipeline is required in this task.

However, make e2e must be non-interactive and suitable for later invocation by CI.

Avoid:

headed browser requirement;
prompts;
machine-specific paths;
dependence on globally installed Playwright;
dependence on an already-running Docker stack.

Use project-local npm tooling.

Documentation updates

Update:

README.md
docs/engineering/ROADMAP.md
docs/engineering/ENGINEERING_LOG.md

README should include:

make e2e

and explain that it:

rebuilds the stack;
executes Playwright against the real application;
tears the stack down.

Do not write the final complete README redesign yet; only add accurate test execution instructions.

Verification

Run:

cd client
npm run lint
npm run test
npm run coverage
npm run build
npm run verify
npx playwright install chromium

Then from root:

make e2e

Run it at least twice from a clean state to detect lifecycle or port-reuse issues.

Also run the local mode once against an already-running stack:

docker compose up --build -d
make e2e-local
docker compose down

Confirm after canonical E2E:

docker compose ps

shows no remaining Abacus containers.

Report

Create:

docs/engineering/agents-reports/ABA-008-playwright-acceptance-suite.md

Include:

final test structure;
Docker lifecycle strategy;
readiness mechanism;
browser projects;
scenarios automated;
real-backend versus intercepted tests;
commands executed;
results from two clean runs;
evidence created;
flaky tests encountered and resolution;
files changed;
known limitations;
no-commit confirmation.

If the hard stop is reached, create:

docs/engineering/agents-reports/CHECKPOINT-ABA-008.md
Constraints

Do not:

mock happy-path calculations;
use arbitrary sleeps for readiness;
leave Docker containers running;
add authentication or persistence;
redesign the UI;
modify operation semantics;
add a hosted CI pipeline;
commit generated traces or videos;
commit.
Definition of Done
Playwright is configured and executable locally.
make e2e owns clean Docker startup and teardown.
make e2e-local supports an existing stack.
All seven operations are tested against the real backend.
Client-side and backend validation are tested.
Continuity and focus are tested.
Clear, Enter, Escape and all shortcuts are tested.
Loading, duplicate submission and failure recovery are tested.
Desktop and mobile layouts are tested.
Accessibility-observable behavior is tested.
Tests are independent and deterministic.
Curated evidence is stored under documentation.
Generated Playwright artifacts remain ignored.
The canonical suite passes twice from a clean state.
Unit verification remains green.
Documentation and the engineering report are current.
No commits were created.
