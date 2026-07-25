ABA-005 — Immutable Operation Registry & Manifest Projection
Field	Value
Status	Ready
Owner	Clio
Role	Main Developer — Go Domain Foundation
Target	30–45 minutes
Hard stop	60 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Implement the immutable operation registry that serves as the backend source of truth for Abacus calculator capabilities.

The registry must:

define the seven supported operations;
associate each operation with an executable Go function;
preserve deterministic presentation order;
support lookup by operation ID;
expose a safe manifest projection for the frontend;
validate its own configuration during application startup;
remain immutable after construction and safe for concurrent reads.

This task does not expose HTTP calculator endpoints and does not implement the declarative validation engine.

Required skills

Apply:

.codex/skills/process/osk-execution-timebox/SKILL.md
.codex/skills/process/osk-engineering-reporting/SKILL.md

Also apply the mirrored Go and architecture guidance available to Codex. At minimum, use the canonical project guidance represented by:

.claude/skills/osk-go-guide/SKILL.md
.claude/skills/architecture/osk-architecture-review/SKILL.md

Follow:

CLAUDE.md
AGENTS.md
Required knowledge inputs

Before implementation, inspect:

client/src/features/calculator/api/types.ts
client/src/features/calculator/api/mockCalculatorApi.ts
docs/knowledge/use-cases/
docs/knowledge/acceptance-tests/
docs/knowledge/glossary/

The current TypeScript contract and mock manifest are an important integration input.

Do not copy TypeScript mechanically into Go. Reconcile the models deliberately so that the future HTTP projection can replace the frontend mock without redesigning the client.

Document any contract mismatch discovered.

Architectural intent

The core model should remain close to:

type ExecuteFunc func(operands []float64) (float64, error)

type Operation struct {
    Definition OperationDefinition
    Execute    ExecuteFunc
}

type Registry struct {
    operationsByID map[string]Operation
    operations     []Operation
    defaultID      string
}

Exact names may change when a clearer idiomatic model emerges.

The important properties are:

Operation
├── frontend-safe definition and metadata
└── trusted backend execution function

Registry
├── ordered operations
├── lookup by ID
├── default operation
└── manifest projection

Do not introduce one concrete strategy type per arithmetic operation.

Supported operations

Register exactly these seven operations in this deterministic order:

Order	ID	Name	Symbol	Shortcut	Arity
1	addition	Addition	+	+	2
2	subtraction	Subtraction	−	-	2
3	multiplication	Multiplication	×	*	2
4	division	Division	÷	/	2
5	exponentiation	Exponentiation	xʸ	^	2
6	square-root	Square Root	√	r	1
7	percentage	Percentage	%	%	2

The default operation is:

addition
Operand semantics
Addition
Operand 1: First number
Operand 2: Second number
Subtraction
Operand 1: Minuend or First number
Operand 2: Subtrahend or Second number

Prefer user-friendly labels consistent with the current frontend unless domain terminology materially improves clarity.

Multiplication
Operand 1: First number
Operand 2: Second number
Division
Operand 1: Dividend
Operand 2: Divisor
Exponentiation
Operand 1: Base
Operand 2: Exponent
Square Root
Operand 1: Number or Radicand

Use prefix notation.

Percentage
Operand 1: Base value
Operand 2: Percentage

Percentage semantics are:

percentage of base = base × percentage / 100

For example:

15% of 200 = 30

The second operand includes a frontend-safe, non-editable suffix:

%

The range validation from 0 through 100 belongs to ABA-006, not this task.

Scope
1. Calculator package

Create a proportional calculator/domain package, likely under:

server/internal/calculator/

A reasonable structure could be:

server/internal/calculator/
├── operation.go
├── operations.go
├── registry.go
├── manifest.go
└── *_test.go

This is guidance, not mandatory ceremony. Prefer fewer files when clearer.

Do not place calculator-domain behavior in:

transport
app
config
observability
2. Operation definition

Define a frontend-safe operation model containing the metadata needed by the React client.

It should support at least:

stable operation ID;
human-readable name;
mathematical symbol;
keyboard shortcut;
arity;
operand definitions;
notation or presentation metadata where genuinely required.

An operand definition should support at least:

stable operand ID or key;
user-facing label;
placeholder where the existing client contract expects one;
optional suffix;
required position/order.

Avoid free-form frontend templates or a generic rendering language.

Prefer small explicit metadata such as:

prefix
infix
suffix

only when the current client genuinely needs it.

3. Function-based execution

Define a function type for executable operations.

Preferred direction:

type ExecuteFunc func(operands []float64) (float64, error)

The operation already owns its function, so callers should not pass the operator separately.

Implement trusted executors for:

addition;
subtraction;
multiplication;
division;
exponentiation;
square root;
percentage.

At this stage:

executors may assume that arity and domain rules were validated by callers;
they must not panic for malformed operand slices;
they should return a clear internal error when invoked with the wrong number of operands;
do not duplicate the future declarative validation engine inside each executor.

Do not yet implement operation-specific domain validations such as:

divisor not zero;
square-root operand non-negative;
percentage range;
exponent domain restrictions.

Those belong to ABA-006.

It is acceptable for direct executor invocation to return Go floating-point behavior for mathematically invalid inputs until the validation layer is introduced, provided the limitation is explicit in tests or the report.

4. Immutable registry

Implement construction similar to:

func NewRegistry(
    defaultOperationID string,
    operations []Operation,
) (*Registry, error)

or a clearer equivalent.

The constructor must validate:

registry contains at least one operation;
default operation ID is not empty;
default operation exists;
every operation ID is non-empty;
operation IDs are unique;
operation name is non-empty;
symbol is non-empty;
shortcut is coherent;
arity is supported;
operand count equals arity;
operand identifiers are present and unique within the operation;
executor function is not nil.

Once construction succeeds:

operations cannot be registered, replaced or removed;
no public Register method exists;
no mutex is required;
concurrent reads are safe because no writes occur.

Do not use:

sync.Map
sync.Mutex
atomic values
package-level mutable registry
5. Lookup

Provide an idiomatic lookup method.

Preferred forms include:

func (registry *Registry) Get(operationID string) (Operation, bool)

or:

func (registry *Registry) Lookup(operationID string) (Operation, bool)

Choose one and explain briefly in the report.

Do not use Java-style names mechanically. GetByID is acceptable only when it materially improves clarity over Get or Lookup.

Consider whether returning Operation by value is safe and clear. The returned operation must not expose mutable registry internals.

6. Deterministic ordering

Go map iteration is intentionally nondeterministic.

The registry must preserve operation order independently of its lookup map.

A suitable internal model is:

type Registry struct {
    operationsByID map[string]Operation
    operations     []Operation
}

Manifest order must always match the declared product order.

7. Manifest projection

Define a frontend-safe manifest projection.

Conceptually:

type Manifest struct {
    Version            string                `json:"version"`
    DefaultOperationID string                `json:"defaultOperationId"`
    Operations         []OperationDefinition `json:"operations"`
}

Use the exact JSON field names required to align with the client contract.

The manifest must:

include the stable version;
include the default operation ID;
preserve deterministic operation order;
include only safe definition metadata;
exclude executor functions;
exclude internal runtime details;
return defensive copies where slices could otherwise expose registry internals.

The registry should generate the manifest from the same registered operations used for backend execution.

Do not maintain a second independently authored manifest catalogue.

8. Default registry construction

Provide an explicit startup construction function, for example:

func NewDefaultRegistry() (*Registry, error)

or:

func DefaultOperations() []Operation

combined with NewRegistry.

Prefer explicit, inspectable composition.

Avoid:

init();
reflection;
plugin loading;
YAML;
dynamic scripting;
expression-based executor definitions;
package-level mutable initialization.
9. Application composition

Wire the registry during application startup so invalid operation configuration fails fast before the HTTP server starts.

The runtime and router do not need to consume it yet if no endpoint exists.

A temporary explicit construction in run() is acceptable, but do not create an unused dependency merely to satisfy compilation.

If wiring the registry without a consumer would create ornamental code, keep the constructor tested and document that application wiring begins with the REST or calculation-service task.

Use engineering judgment and report the choice.

Validation model boundary

ABA-005 may prepare types for validations only when required to keep the manifest contract coherent with the existing frontend.

Do not implement:

expression evaluation;
comparison operators;
logical conjunction;
violation generation;
domain validation execution.

Avoid adding empty placeholder abstractions.

A validation field may be omitted from the manifest for now or represented as an empty collection if required by the frozen client contract. Choose based on compatibility with the existing React client and document the decision.

Tests

Add table-driven tests covering the operation model and registry.

Registry construction

Test:

successful default registry construction;
empty operation list;
empty default operation ID;
unknown default operation;
duplicate operation IDs;
empty operation ID;
empty operation name;
unsupported arity;
operand count mismatch;
duplicate operand IDs;
nil executor.
Registry behavior

Test:

lookup of each supported operation;
unknown operation lookup;
default operation;
deterministic order;
registry remains unchanged when caller mutates the source slice after construction;
manifest caller cannot mutate registry-owned definitions through returned slices.
Operation execution

For valid normal inputs, test:

Operation	Operands	Result
Addition	10, 5	15
Subtraction	5, 10	-5
Multiplication	2.5, 4	10
Division	144, 12	12
Exponentiation	2, 8	256
Square Root	144	12
Percentage	200, 15	30

Also test wrong arity does not panic and returns an error.

Do not test future validation behavior as if it already existed.

Manifest projection

Test:

version;
default operation;
seven operations;
stable order;
JSON field names;
square-root arity and notation;
percentage operand suffix;
executor functions are not serializable or exposed;
empty validations behavior matches the agreed temporary contract.
Naming and documentation
Use descriptive names.
Preserve standard Go initialisms such as ID, HTTP, and JSON.
Avoid compressed names except conventional locals such as err.
Every exported Go symbol requires a Go documentation comment beginning with the symbol name.
Comments explain intent, ownership or trade-offs rather than repeating implementation.
Explicit non-goals

Do not implement:

HTTP operation endpoint;
calculation endpoint;
transport DTO mapping;
declarative validation interpreter;
validation error contract;
calculator application service;
frontend HTTP integration;
Playwright;
OpenAPI;
runtime registration of operations;
configuration-driven YAML operations;
dynamic plugin loading;
database persistence.
Documentation

Update:

docs/engineering/ROADMAP.md
docs/engineering/ENGINEERING_LOG.md

Add an ADR only if a new durable decision appears beyond the already agreed architecture.

Potential decision to document, if not already covered:

The operation registry is immutable after startup and uses an ordered slice plus an ID lookup map.

Do not create an ADR merely to restate code when the decision is already captured sufficiently in durable architecture documentation.

Update glossary or architecture knowledge when terms such as Operation, Registry, Manifest, and ExecuteFunc are not yet defined.

Do not add technical debt for work intentionally assigned to ABA-006 or ABA-007.

Required verification

Run and record:

make fmt
make lint
make test
make coverage
make build
make verify

Also run:

cd server
go test -race ./...

If Docker build context includes the changed backend source, verify:

docker compose build server

No calculator endpoint needs manual HTTP verification in this task.

Deliverables

Expected artifacts include:

server/internal/calculator/operation.go
server/internal/calculator/operations.go
server/internal/calculator/registry.go
server/internal/calculator/manifest.go
server/internal/calculator/*_test.go

File names may differ if a simpler organization is clearer.

Also deliver:

updated roadmap;
updated Engineering Log;
architecture/glossary updates where useful;
implementation report.
Report and checkpoint

Create:

docs/agents/reports/ABA-005-operation-registry.md

Use the repository’s canonical report path if it differs.

The report must include:

final package structure;
registry invariants;
immutability strategy;
concurrency reasoning;
lookup method chosen;
manifest projection;
client-contract reconciliation;
executor behavior;
tests and command results;
contract mismatches found;
deferred validation behavior;
files changed;
no-commit confirmation.

If the hard stop is reached, stop and create:

docs/agents/checkpoints/CHECKPOINT-ABA-005.md
Definition of Done
All seven operations are registered.
Every operation has metadata and an executable Go function.
Percentage semantics are correct.
Square Root is unary.
Registry construction validates its invariants.
Registry is immutable after construction.
Concurrent reads require no synchronization.
Lookup by stable operation ID works.
Manifest order is deterministic.
Manifest is generated from registered operations.
Manifest excludes backend executors.
Manifest aligns with the existing client contract or explicitly documents any mismatch.
Tests cover construction, lookup, order, defensive copies, execution, arity and projection.
Formatting, lint, tests, race detector, coverage, build and verification pass.
No validation engine, HTTP calculator endpoint or frontend integration was introduced.
Documentation is current.
No commits were created.
