ABA-006 — Declarative Validation Model
Field	Value
Status	Ready
Owner	Clio
Role	Main Developer — Domain Modeling
Target	35–50 minutes
Hard stop	60 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Implement the immutable declarative validation model used by the Abacus calculator.

The backend must define validation rules as typed domain objects.

Those same rules will later be projected through the manifest so both:

backend
frontend

validate calculator operations from the same source of truth.

This task implements:

validation model;
evaluator;
manifest projection;
integration with the operation registry.

This task does not expose REST endpoints.

Required skills

Apply:

.codex/skills/process/osk-execution-timebox/SKILL.md
.codex/skills/process/osk-engineering-reporting/SKILL.md

.claude/skills/osk-go-guide/SKILL.md
.claude/skills/architecture/osk-architecture-review/SKILL.md

Follow:

CLAUDE.md
AGENTS.md
Required inputs

Review:

ABA-005 implementation
current frontend validation interpreter
frontend manifest contract
use cases
acceptance tests

The Go model should naturally replace the mocked frontend validation definitions.

Architectural intent

The validation model is not a scripting language.

It is not an expression parser.

It is not a rule engine.

It is a small closed domain model.

Only implement the minimum concepts required by the calculator.

Supported validation expressions

Implement only these expression types.

Operand reference

Represents:

operand[0]

operand[1]
Literal

Represents:

0

100
Comparison

Supported operators:

equal

notEqual

greaterThan

greaterThanOrEqual

lessThan

lessThanOrEqual

No arithmetic.

No function invocation.

No nesting inside comparisons.

Logical composition

Implement only:

allOf

No:

anyOf

not

xor

implication
Validation definition

A validation should resemble:

type ValidationDefinition struct {
    ID         string
    Message    string
    Expression Expression
}

Exact naming may change if a clearer Go model emerges.

Expression hierarchy

Keep the hierarchy intentionally small.

A suitable direction:

Expression

├── ComparisonExpression

├── AllOfExpression

├── OperandReference

└── Literal

Avoid reflection.

Avoid dynamic dispatch through strings.

Avoid generic map-based ASTs.

Evaluator

Implement a recursive evaluator.

Conceptually:

Evaluate(expression, operands)

returns:

bool

No error collection yet.

No partial evaluation.

No optimization.

No expression rewriting.

Validation execution

Implement:

Validate(operation, operands)

Returns:

nil when all validations succeed

or

the first violated validation.

The first failure is sufficient.

Do not accumulate multiple errors.

Operation integration

Every operation now owns:

Validations

Examples:

Division

operand[1] != 0

Square Root

operand[0] >= 0

Percentage

operand[1] >= 0

AND

operand[1] <= 100

Addition

No validations.

Subtraction

No validations.

Multiplication

No validations.

Exponentiation

No validations.

Manifest projection

Remove the temporary:

json.RawMessage("[]")

Replace with the strongly typed validation projection.

The frontend should receive the same logical tree represented by Go.

No executable code.

No serialized functions.

No scripts.

JSON contract

Maintain compatibility with the existing frontend interpreter whenever possible.

If a contract improvement is required:

document it;
keep it minimal.

Do not redesign the frontend contract.

Tests

Cover:

Expressions
every comparison operator
operand references
literals
allOf
Validation execution

Division

12 / 3

PASS
12 / 0

FAIL

Square Root

sqrt(9)

PASS
sqrt(-1)

FAIL

Percentage

50%

PASS
101%

FAIL
-5%

FAIL

Operations without validations

Always pass.

Manifest

Verify:

validations serialized
deterministic ordering
no executor leakage
Documentation

Update:

ROADMAP

ENGINEERING_LOG

Glossary

Explain:

ValidationDefinition
Expression
Comparison
OperandReference
Literal
AllOf

No ADR unless a genuinely new architectural decision emerges.

Validation

Run:

make fmt
make lint
make test
make coverage
make build
make verify

go test -race ./...
Deliverables

Expected additions:

server/internal/calculator/validation.go

comparison.go

evaluator.go

validation_test.go

Exact file organization is free provided it remains proportional.

Constraints

Do not implement:

parser
DSL
scripting
reflection
runtime expression compilation
YAML
plugins
REST endpoints
frontend integration
OpenAPI

Do not introduce a generic rules engine.

Do not commit.

Definition of Done
Typed validation model implemented.
Recursive evaluator implemented.
Operations own declarative validations.
Registry projects typed validations.
Temporary json.RawMessage removed.
Frontend contract remains compatible or differences documented.
Division, Square Root and Percentage validations implemented.
Tests cover expressions and operation validation.
Formatting, lint, tests, coverage, race detector and verification pass.
Documentation updated.
No REST endpoint introduced.
No commit created.

Important:
Prefer explicit domain types over generic AST nodes.

If at any point the implementation starts resembling a scripting language or a generic expression engine, stop and simplify. The goal is to model exactly the calculator's validation needs—not to create a reusable rules engine.
