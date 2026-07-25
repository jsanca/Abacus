ABA-004 — Responsive React Frontend with Mocked Calculator API
Field	Value
Status	Ready
Owner	Elito
Role	Frontend Developer · UX Implementer · Knowledge Curator
Target	30–45 minutes
Hard stop	60 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Build the initial Abacus frontend using React, TypeScript, and Vite against a mocked calculator API.

The frontend must establish:

the responsive operation-based calculator experience;
visual direction inspired by the exported Google Stitch artifacts;
component boundaries;
accessibility;
state transitions;
backend contract types;
mocked operation discovery and calculation behavior.

The frontend must not depend on the real Go backend during this task.

Required inputs

Review all exported Stitch artifacts:

docs/knowledge/design/stitch_abacus_calculator_design_app/
├── DESIGN.md
├── code.html
└── screen.png
docs/knowledge/design/stitch_abacus_calculator_design_mobile/
├── DESIGN.md
├── code.html
└── screen.png

These artifacts are engineering-reviewed visual references, not pixel-perfect implementation specifications.

The React implementation is the source of truth and may diverge to improve:

usability;
accessibility;
responsiveness;
semantic HTML;
backend-driven behavior;
maintainability;
correctness.
Required skills and instructions

Apply:

repository CLAUDE.md;
repository AGENTS.md;
.codex/skills/process/osk-execution-timebox/SKILL.md;
.codex/skills/process/osk-engineering-reporting/SKILL.md;
relevant architecture and knowledge-curation skills.

Use descriptive component, property, function, hook, and state names.

Avoid compressed names such as op, val, res, cfg, or data where a more precise name is available.

Technical foundation

Initialize the client under:

client/

Use:

React;
TypeScript;
Vite;
Vitest;
React Testing Library.

Avoid unnecessary state-management or data-fetching frameworks.

Do not introduce Redux, Zustand, React Query, a UI framework, or a form framework unless a demonstrated need emerges.

Prefer:

native React state;
typed domain models;
a small API boundary;
plain CSS, CSS Modules, or another lightweight styling method already justified by the repository.
Functional experience

Implement an operation-based calculator, not a physical numeric keypad.

Supported operations in the mocked manifest:

Addition
Subtraction
Multiplication
Division
Exponentiation
Square Root
Percentage

All seven operator controls must be visible and selectable with one click.

No dropdown is allowed for primary operator selection.

Layout
Desktop

Use the Stitch desktop proposal as visual inspiration:

centered calculator card;
clean neutral canvas;
restrained indigo/purple accent;
prominent result area;
clear input hierarchy;
generous whitespace;
professional fintech tone.
Mobile

Implement a genuinely responsive mobile layout:

no horizontal scrolling;
comfortable touch targets;
stacked content where appropriate;
all seven operators visible;
operators may wrap into two balanced rows;
primary action easy to reach;
result remains visually prominent.

Do not reproduce Stitch's omission of operators in the mobile design.

Required design refinements

Relative to the exported Stitch reference:

Preserve
overall clean fintech aesthetic;
neutral background;
indigo/purple accent direction;
centered card concept;
rounded but restrained controls;
clear result hierarchy;
simple operation-based interaction.
Change
display all seven operations on desktop and mobile;
use unambiguous mathematical symbols:
+
−
×
÷
xʸ
√
%
make inputs look clearly editable rather than disabled;
add visible hover and keyboard focus states;
strengthen result prominence;
add the required UI states;
support operation-specific layouts;
preserve accessibility and semantic structure.
Remove

Do not implement:

settings;
application navigation;
privacy policy links;
terms links;
help center links;
copyright footer;
history management;
decorative features outside assignment scope.
Interaction model

The UI must support these states:

loading capabilities
idle
editing
calculating
resolved
validation error
service error
Initial loading

Since the real frontend will depend on backend-discovered operations, display a capabilities-loading state while the mock manifest loads.

Do not silently hardcode the visible UI independently of the mocked manifest.

Operation selection

Selecting an operator:

changes the active state using aria-pressed;
updates operand labels and presentation from manifest metadata;
uses a single click;
supports keyboard activation;
moves focus appropriately.
Unary operation

For square root:

only one operand is shown;
prefix notation is used;
the second operand is removed from tab order and layout.
Percentage

For percentage:

Operand 1 represents the base value;
Operand 2 represents a percentage from 0 through 100;
show a fixed, non-editable % suffix in Operand 2;
the mock operation should calculate the percentage of Operand 1.
Calculation continuity

After resolving:

144 ÷ 12 = 12

keep the completed expression and result visible.

When the user selects a new operator:

move the previous result into Operand 1;
clear Operand 2 for binary operations;
clear the previous resolved state as appropriate;
focus Operand 2;
for square root, apply the unary layout and focus Operand 1 as appropriate.

Also provide the secondary explicit continuity methods already agreed where they fit naturally, without cluttering the primary flow.

Clear

Provide a restrained secondary Clear action.

Clear must reset:

operands;
selected or default operation as defined by the UX;
result;
field errors;
service errors;
focus to the first operand.
Keyboard behavior

At minimum:

Enter calculates when valid;
Escape clears;
standard operator shortcuts select the associated operation where they do not interfere with text input;
focus remains visible and predictable.

Document shortcuts in an accessible, non-intrusive way.

Numeric input handling

Use controlled text state with numeric input semantics, rather than relying blindly on browser number-input behavior.

A suitable approach is:

type="text"
inputMode="decimal"

The implementation must support legitimate editing states such as:

""
"-"
"12."
"-0."

before converting to a number for submission.

The frontend may know only structural requirements:

operation exists;
one or two operands are required;
entered values can be represented as numbers.

Operation-specific validation must come from the mocked backend manifest and be interpreted generically.

Do not write logic such as:

if (operation.id === "division" && secondOperand === 0)
Mocked contract

Create typed fixtures representing the agreed backend projection.

At minimum model:

type OperationManifest
type OperationDefinition
type OperandDefinition
type ValidationDefinition
type CalculationRequest
type CalculationResponse
type ApiError

The exact field names should remain coherent and implementation-ready.

The mocked API must expose conceptual equivalents of:

GET /api/v1/operations
POST /api/v1/calculations

Use an asynchronous boundary with a small artificial delay only if needed to make states testable. Do not add a noticeable artificial product delay.

Declarative validation

The frontend must use a generic interpreter for the limited validation expressions supplied by the mock manifest.

The frontend must not understand domain rules such as “division by zero.”

It should understand only the portable validation vocabulary agreed for the contract, for example:

operand reference;
literal value;
equality and comparison operations;
logical conjunction where required.

Keep the expression model minimal.

Do not create a general-purpose language, dynamic JavaScript execution, eval, or new Function.

The backend remains authoritative even when the client later performs the same projected rules for immediate feedback.

Mock behavior

The mock service should support:

successful operation manifest loading;
manifest loading failure;
successful calculations;
calculation loading;
declarative validation failures;
unknown operation response;
simulated backend/service failure.

Keep failure controls internal to tests or development helpers; do not clutter the production UI with debug toggles.

Accessibility

At minimum:

semantic labels for operands;
fieldset and legend or another appropriate grouping for operators;
aria-pressed for selected operator buttons;
aria-live="polite" for successful results;
role="alert" or equivalent for errors;
visible focus styles;
sufficient contrast;
minimum practical touch target sizing;
no state communicated by color alone;
logical tab order.
Suggested proportional structure
client/
├── src/
│   ├── app/
│   ├── features/
│   │   └── calculator/
│   │       ├── api/
│   │       ├── components/
│   │       ├── model/
│   │       ├── hooks/
│   │       └── test/
│   ├── shared/
│   └── main.tsx
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts

This is guidance, not a mandatory folder ceremony. Prefer fewer packages if clearer.

Testing

Add meaningful tests for:

manifest-driven rendering of all seven operations;
unary square-root layout;
percentage suffix and range validation from manifest;
selecting an operator in one click;
successful mocked calculation;
division validation without frontend operation-specific logic;
calculation continuity from a previous result;
clear behavior;
loading state;
service error state;
keyboard calculation;
basic responsive class or structure where practical.

Use behavior-focused tests. Avoid testing implementation details or exact CSS values.

Quality commands

Provide and run:

npm run lint
npm run test
npm run coverage
npm run build

Add a client verification command that runs the meaningful checks together:

npm run verify

Integrate it into the root Makefile if the backend task has already established one, without disrupting parallel work. When concurrent file edits would conflict, document the required integration instead of overwriting another agent's changes.

Documentation

Update durable artifacts:

mark ABA-002 complete if not already updated;
add ABA-004 to ROADMAP.md;
add the task and artifacts to ENGINEERING_LOG.md;
record material deviations from Stitch in the implementation report;
update technical debt only for accepted limitations;
update docs/knowledge/design/ only when a concise implementation note adds durable value.

Do not rewrite the Stitch exports.

Explicit non-goals

Do not:

call the real backend;
implement Docker;
implement authentication;
implement persistent history;
use a physical keypad;
copy code.html directly into React;
reproduce Stitch pixel-for-pixel;
add legal links or settings;
add financial precision behavior not yet agreed;
invent additional calculator operations;
introduce a large design system or component library.
Deliverables
Vite React TypeScript application;
responsive desktop and mobile calculator;
seven manifest-driven operations;
mocked capability and calculation APIs;
generic projected-validation interpreter;
continuous calculation flow;
loading, success, validation, and service-error states;
keyboard and accessibility behavior;
tests, linting, coverage command, build, and verification command;
implementation report and updated engineering artifacts.
Report and checkpoint

Create:

docs/engineering/agents-reports/ABA-004-frontend-foundation.md

If the hard stop is reached, stop and create:

docs/engineering/agents-reports/CHECKPOINT-ABA-004.md

Include:

status;
files created and changed;
Stitch artifacts reviewed;
design elements preserved;
intentional deviations;
mock contract;
UX state model;
accessibility work;
tests and command results;
integration assumptions;
remaining work;
no-commit confirmation.
Definition of Done
The client runs independently using mocked APIs.
Desktop and mobile layouts are responsive.
All seven operators are visible and selectable in one click.
Stitch is reflected as inspiration without being treated as a fixed specification.
Square root and percentage have their required specialized presentation.
Operation-specific validation is manifest-driven and generically interpreted.
Calculation continuity works.
Required loading and error states exist.
Keyboard navigation and accessible semantics are implemented.
Tests, lint, coverage generation, build, and npm run verify succeed.
Durable documentation is updated.
No real backend integration or commits were introduced.
