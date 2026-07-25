AT-001 — Calculator Experience Playwright Recipe
Purpose

Define the browser-level acceptance scenarios that Elito will automate using Playwright after the real frontend and backend are integrated.

These tests verify product behavior across the deployed client/server boundary. They are not replacements for React unit tests or Go API tests.

Test level
Browser
  ↓
React client
  ↓
HTTP API
  ↓
Go calculator service

Except where a failure must be simulated, Playwright should exercise the real backend rather than mock calculation responses.

Required Playwright setup

Recommended location:

client/
├── playwright.config.ts
└── tests/
    └── e2e/
        ├── calculator-happy-path.spec.ts
        ├── calculator-validation.spec.ts
        ├── calculator-continuity.spec.ts
        ├── calculator-failures.spec.ts
        ├── calculator-accessibility.spec.ts
        └── calculator-responsive.spec.ts

Playwright should:

start the backend and frontend through webServer configuration or a repository command;
use stable accessible locators;
avoid CSS implementation selectors;
run at least Chromium desktop and a mobile viewport;
capture trace and screenshot on failure;
run serially only where a scenario truly shares state;
keep every test independently executable.

Preferred locators:

page.getByRole("button", { name: "Division" })
page.getByLabel("Dividend")
page.getByLabel("Divisor")
page.getByRole("button", { name: "Calculate" })
page.getByRole("status")
page.getByRole("alert")

Add data-testid only where no stable semantic locator exists.

Preconditions
Backend health check succeeds.
Frontend is available.
Operation manifest includes all seven expected operations.
Browser begins with a fresh page state.
Acceptance scenarios
AT-001.01 — Load calculator capabilities

Given

the backend is available;

When

the user opens Abacus;

Then

the capabilities-loading state is eventually removed;
all seven operation controls are visible;
one operation has an understandable selected state;
the calculator is ready for input.
AT-001.02 — Add two positive numbers

Given

Addition is selected;

When

Operand 1 is 10;
Operand 2 is 5;
the user activates Calculate;

Then

the result is 15;
the completed expression communicates 10 + 5 = 15;
no error is visible.
AT-001.03 — Subtract to a negative result

When

the user calculates 5 − 10;

Then

the result is -5.
AT-001.04 — Multiply decimal values

When

the user calculates 2.5 × 4;

Then

the result is 10.
AT-001.05 — Divide successfully

When

the user calculates 144 ÷ 12;

Then

the result is 12.
AT-001.06 — Reject division by zero

Given

Division is selected;

When

Operand 1 is 10;
Operand 2 is 0;
the user calculates;

Then

no calculation request should be sent when client-side projected validation already rejects it, where observable;
an error identifies the divisor;
no result is displayed;
Operand 1 remains 10;
Operand 2 remains 0.

The test must assert behavior, not hardcode the exact English message unless the message is part of the public contract.

AT-001.07 — Exponentiation

When

the user calculates 2 raised to 8;

Then

the result is 256.
AT-001.08 — Square root unary layout

When

the user selects Square Root;

Then

only one operand input is visible;
a second operand is absent from the tab order;
prefix notation is communicated.

When

the user calculates the square root of 144;

Then

the result is 12.
AT-001.09 — Reject invalid real square root

When

the user attempts to calculate the square root of -1;

Then

the operation is rejected;
the affected operand is identified;
no invalid or NaN result is shown.
AT-001.10 — Percentage lower boundary

When

the user calculates 0% of 200;

Then

the result is 0;
Operand 2 displays a fixed % suffix.
AT-001.11 — Percentage upper boundary

When

the user calculates 100% of 200;

Then

the result is 200.
AT-001.12 — Reject percentage below range

When

the percentage operand is -1;

Then

validation prevents a successful calculation;
the percentage field is identified.
AT-001.13 — Reject percentage above range

When

the percentage operand is 101;

Then

validation prevents a successful calculation;
the percentage field is identified.
AT-001.14 — Continue from previous result

Given

the user has calculated 144 ÷ 12 = 12;

When

the user selects Multiplication;

Then

Operand 1 becomes 12;
Operand 2 is empty;
Multiplication is active;
focus is on Operand 2.

When

the user enters 3 and calculates;

Then

the result is 36.
AT-001.15 — Continue from unary result

Given

the user has calculated √144 = 12;

When

the user selects Addition;

Then

Operand 1 becomes 12;
Operand 2 is empty;
focus moves to Operand 2.
AT-001.16 — Clear with button

Given

the calculator contains operands, a result, and possibly an error;

When

the user activates Clear;

Then

operand inputs are reset;
result and expression are removed;
errors are removed;
focus returns to Operand 1.
AT-001.17 — Calculate using Enter

Given

valid operands are entered;

When

the user presses Enter;

Then

the calculation is submitted once;
the result is displayed.
AT-001.18 — Clear using Escape

Given

the calculator has entered values;

When

the user presses Escape;

Then

the calculator returns to its cleared state.
AT-001.19 — Operator keyboard shortcuts

For each supported shortcut:

When

the shortcut is pressed outside conflicting text-entry behavior;

Then

the associated operation becomes active;
the selected state is exposed accessibly.

Use a table-driven Playwright loop only when failures remain independently identifiable.

AT-001.20 — Preserve values after service failure

Given

the calculation endpoint is unavailable or returns a server failure;

When

the user calculates a valid expression;

Then

a service error is displayed;
operands remain intact;
the selected operation remains intact;
no incorrect result is shown;
the user can retry.

This scenario may use Playwright route interception because it intentionally simulates infrastructure failure.

AT-001.21 — Manifest loading failure and retry

Given

the operation-manifest request initially fails;

When

the user opens Abacus;

Then

a calculator-unavailable state is shown;
hardcoded operation controls are not silently substituted;
a Retry action is visible.

When

the next manifest request succeeds and Retry is activated;

Then

the calculator becomes usable;
all operations appear.
AT-001.22 — Loading state prevents duplicate submission

Given

a calculation request is pending;

When

the user activates Calculate repeatedly;

Then

only one request is active;
the primary action communicates progress or becomes disabled;
duplicate results are not produced.

A short route delay may be used to make the pending state observable.

AT-001.23 — Desktop responsive layout

Run at a desktop viewport such as 1440 × 900.

Assert:

calculator content is fully visible;
all seven operations are visible;
no horizontal scrollbar exists;
result and actions remain logically ordered.

Avoid asserting exact pixels.

AT-001.24 — Mobile responsive layout

Run using a mobile Playwright project or viewport near 390 × 844.

Assert:

no horizontal overflow exists;
all seven operation controls remain visible;
operator controls wrap appropriately;
controls have usable dimensions;
inputs, result, and actions remain visible and operable.
AT-001.25 — Keyboard-only primary flow

Using only keyboard actions:

focus Operand 1;
enter a value;
navigate to an operation;
select it;
enter Operand 2;
calculate;
verify the result.

Assert that focus is always visible and the flow does not trap the user.

AT-001.26 — Accessible result announcement

When

a calculation succeeds;

Then

the result is exposed through an aria-live region or equivalent status semantics.

Playwright should validate the semantic attribute or role; it cannot fully emulate every screen reader announcement.

AT-001.27 — Accessible validation error

When

an operand fails validation;

Then

an alert or equivalent accessible error is present;
the affected input exposes its invalid state;
the input is programmatically associated with the error when applicable.
Playwright implementation conventions
Test data

Use named test cases:

const arithmeticScenarios = [
  {
    name: "adds positive integers",
    operationName: "Addition",
    firstOperand: "10",
    secondOperand: "5",
    expectedResult: "15",
  },
];

Avoid opaque data like:

["add", "10", "5", "15"]
Page object

A small page object is appropriate because the same interaction vocabulary repeats:

class CalculatorPage {
  async open(): Promise<void>
  async selectOperation(name: string): Promise<void>
  async enterFirstOperand(value: string): Promise<void>
  async enterSecondOperand(value: string): Promise<void>
  async calculate(): Promise<void>
  async clear(): Promise<void>
  async expectResult(value: string): Promise<void>
}

Do not hide assertions or entire business flows behind oversized page-object methods.

Good:

await calculatorPage.selectOperation("Division");
await calculatorPage.enterFirstOperand("144");
await calculatorPage.enterSecondOperand("12");
await calculatorPage.calculate();
await calculatorPage.expectResult("12");

Avoid:

await calculatorPage.verifyDivisionWorks();
Network policy

Use the real backend for:

supported operation discovery;
arithmetic happy paths;
domain validation;
error contract;
integration behavior.

Use route interception only for controlled infrastructure scenarios:

manifest unavailable;
calculation service failure;
delayed response;
malformed response if that scenario becomes required.
Assertions

Assert user-observable and accessible behavior:

role;
accessible name;
visible value;
focus;
error association;
result;
request count;
absence of horizontal overflow.

Avoid:

internal React state;
component names;
exact class names;
exact CSS colors;
snapshots of the full page as the primary assertion;
arbitrary timeouts.

Use Playwright auto-waiting and web-first assertions.

Required commands

Eventually expose:

npm run test:e2e
npm run test:e2e:headed
npm run test:e2e:report

And at root:

make e2e

The E2E command must start or require a documented deterministic stack.

Evidence on failure

Configure:

trace: retain-on-failure
screenshot: only-on-failure
video: retain-on-failure

No generated Playwright evidence should be committed unless intentionally included as a review artifact.

Definition of acceptance-suite readiness
Every automated scenario links to UC-001 and its relevant rule.
Tests use semantic locators.
Happy-path tests use the real client/server stack.
Infrastructure errors are simulated explicitly and sparingly.
Desktop and mobile projects run.
Tests are independent and deterministic.
Retry, loading, validation, continuity, keyboard, and accessibility paths are represented.
The suite can be executed through a documented single command.
