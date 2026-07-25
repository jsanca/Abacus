ABA-006R — Declarative Validation Architecture Review
Field	Value
Status	Ready
Owner	Deep Pro
Role	Independent Go & Domain Architecture Reviewer
Target	30–40 minutes
Hard stop	50 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Review the final declarative validation architecture delivered through ABA-006 and refined by ABA-006.1.

Evaluate only the current cohesive design. Do not review or recommend restoring the superseded centralized Evaluate and free Validate functions.

The review must determine whether the validation model is:

cohesive;
idiomatic in Go;
safe against malformed definitions;
deeply immutable;
compatible with the frontend contract;
proportional to the calculator domain;
ready to become a public REST contract.

No code modifications are allowed.

Required skills

Apply the OpenCode copies of:

.opencode/skills/process/osk-execution-timebox/SKILL.md
.opencode/skills/process/osk-engineering-reporting/SKILL.md
.opencode/skills/architecture/osk-architecture-review/SKILL.md
.opencode/skills/osk-go-guide/SKILL.md

Use the actual mirrored paths available in the repository if their organization differs.

Follow:

CLAUDE.md
AGENTS.md
Required inputs

Review:

ABA-006 task specification;
ABA-006 implementation report;
ABA-006.1 refinement task;
ABA-006.1 implementation report;
server/internal/calculator/validation.go;
operation.go;
operations.go;
registry.go;
manifest.go;
all calculator tests;
frontend TypeScript validation types and interpreter;
relevant use cases and acceptance tests.
Final architecture under review
Operation.Validate(operands)
        │
        ▼
ValidationDefinition
        │
        ▼
Expression.evaluate(operands)
        │
        ├── ComparisonExpression
        └── AllOfExpression

Registry construction separately performs fail-fast structural validation:

NewRegistry
    │
    ▼
validateRegistryOperation
    │
    ▼
validateExpressionDefinition

Manifest projection serializes the same expression definitions consumed by backend evaluation.

Review dimensions
1. Ownership and cohesion

Evaluate whether:

Operation.Validate is the correct primary domain entry point;
operation declarations keep metadata, validations and execution together;
each expression owning evaluate improves cohesion;
arithmetic execution remains appropriately function-based;
no competing validation entry point exists.

Determine whether the current model provides a practical single point of truth when adding an operation.

2. Expression interface

Review:

type Expression interface {
    expressionNode()
    evaluate(operands []float64) (bool, error)
}

Assess:

sealed-interface technique;
unexported methods;
package boundary;
exhaustiveness;
pointer versus value implementations;
nil-interface edge cases;
whether external consumers need capabilities they cannot obtain.

Confirm the interface remains domain-specific rather than becoming a generic rule engine.

3. Runtime evaluation

Review:

operand resolution;
missing operand behavior;
unsupported operand references;
unsupported comparison operators;
allOf short-circuiting;
error wrapping and sentinel errors;
vacuous truth for empty allOf;
panic safety.

Verify that malformed or incomplete operand input cannot silently resolve to zero.

4. Operation validation

Review Operation.Validate for:

validation declaration order;
first-violation behavior;
contextual errors;
returned pointer safety;
allocations;
nil expressions;
absence of duplicated rules.

Confirm that the future application service can use it naturally:

violation, err := operation.Validate(operands)

before invoking:

operation.Execute(operands)
5. Definition validation

Give special attention to:

validateExpressionDefinition(expression Expression, arity Arity) error

Determine whether the centralized type switch is:

a reasonable closed-union validator;
unnecessarily duplicative;
better located on the expression types;
missing any important structural invariant.

Do not recommend moving it merely to eliminate a switch. Recommend change only if it materially improves correctness, cohesion or maintainability.

Assess whether expression types should own a private method such as:

validateDefinition(arity Arity) error

and explicitly compare the trade-offs:

central validator
vs.
expression-owned definition validation

Consider that expression types already own:

evaluation;
JSON serialization.

Avoid concentrating responsibilities on the nodes without clear value.

6. Deep immutability

Verify recursively that:

source validation slices cannot mutate registry state;
nested AllOfExpression.Expressions slices are copied;
Find() results cannot mutate registry state;
manifest projections cannot mutate registry state;
interface values do not conceal mutable pointers or slices;
copyExpression handles every supported expression type correctly.

Review the default branch in copyExpression:

default:
    return expression

Determine whether this is safe under a sealed interface or whether an unsupported internal expression should fail more explicitly.

7. Registry fail-fast behavior

Review structural validation for:

validation ID;
message;
nil expression;
operand compatibility with arity;
known comparison operator;
nested allOf;
nil children;
unknown expression type.

Determine whether malformed backend-authored rules are rejected before the server starts.

Identify any configuration errors that remain discoverable only during a request.

8. Domain and wire-model coupling

Expression types currently own MarshalJSON.

Evaluate whether this coupling is proportional:

domain expression
        +
wire serialization

Consider:

the manifest is a primary projection of the domain;
only one JSON representation exists;
TypeScript requires a discriminated union;
transport concerns should not generally leak into domain behavior.

Do not recommend transport mappers unless they reduce real coupling or improve testability enough to justify the extra model.

9. Frontend contract compatibility

Compare Go JSON with TypeScript:

comparison;
allOf;
operand references;
operators;
literal values;
nested expressions;
validations array;
validation message;
validation ID inclusion or exclusion.

Verify exact runtime compatibility, not merely conceptual similarity.

Determine whether frontend and backend evaluate all supported rules equivalently, including edge cases such as:

missing operands;
NaN;
infinities;
negative zero;
empty allOf.

Flag only material divergences relevant to the application.

10. Test quality

Review tests for:

operation validation;
every comparison operator;
boundaries;
malformed expressions;
startup rejection;
short-circuit behavior;
sentinel error matching;
deep immutability;
deterministic JSON;
contract stability.

Identify valuable missing tests without optimizing for percentage alone.

11. Scope proportionality

Confirm the solution did not become:

a DSL;
a parser;
a generic rules engine;
a reflection-based evaluator;
a plugin system;
an unnecessary type hierarchy.

Assess whether the implementation sits in the intended sweet spot between hardcoded rules and overengineering.

12. REST readiness

Determine whether the model is ready for:

operation-manifest endpoint;
calculation service;
validation-error response;
frontend HTTP integration.

Identify any change that should happen before publishing the JSON contract.

Required verification

Run and report:

make fmt
make lint
make test
make coverage
make build
make verify

cd server
go test -race ./...

Also run:

docker compose build server

Inspect serialized manifest examples for:

division;
square root;
percentage;
operation without validations.
Severity model

Use:

BLOCKER — panic, incorrect validation, mutation leak or unsafe public contract;
MAJOR — architectural or contract problem that should be fixed before REST;
MINOR — localized maintainability, clarity or test issue;
NOTE — optional improvement.

Do not elevate stylistic preference into a finding.

Expected output

Create:

docs/engineering/agents-reviews/ABA-006R-declarative-validation-review.md

Use the repository’s canonical review path if different.

Each finding must include:

Finding
Severity
Location
Evidence
Why it matters
Smallest recommended change

Finish with:

PASS
PASS WITH OBSERVATIONS
CHANGES REQUIRED

Also include an Architectural Confidence table from 1–5 for:

domain cohesion;
Go idiomaticity;
runtime correctness;
startup validation;
immutability;
frontend contract compatibility;
scope proportionality;
REST readiness.

End with a concise answer to:

If you inherited this validation model as a Staff Engineer, would you publish it as the contract used by the REST API and React client? Why or why not?

Constraints

Do not:

modify code;
restore the old free functions;
add expression types;
redesign the product rules;
recommend a generic evaluator framework;
add REST endpoints;
commit.
Definition of Done
Final cohesive architecture reviewed.
Operation.Validate ownership assessed.
Expression evaluation reviewed.
validateExpressionDefinition trade-off explicitly assessed.
Deep-copy implementation verified.
Fail-fast registry behavior reviewed.
JSON and TypeScript contracts compared.
Tests and quality gates executed.
REST readiness determined.
Final verdict issued.
No code or commits created.
