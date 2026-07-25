ABA-007 — Calculator REST API & Frontend Integration
Field	Value
Status	Ready
Owner	Clio
Role	Main Developer — Transport & Application Integration
Target	45–60 minutes
Hard stop	75 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Expose the immutable calculator domain through a REST API.

This slice publishes the backend-generated manifest and calculator execution while preserving the architecture established in previous slices.

The React frontend must be able to replace the mock adapter with the HTTP implementation without requiring UI redesign.

Precondition

Before implementing the REST endpoints, resolve ABA-006R finding F1.

Current implementation:

default:
    return expression

Replace it with a fail-fast implementation:

default:
    panic(fmt.Sprintf(
        "copyExpression: unsupported expression type %T",
        expression,
    ))

This preserves the registry's deep immutability guarantee.

No additional task is required; include this remediation as part of ABA-007.

Required skills

Apply:

.codex/skills/process/osk-execution-timebox/SKILL.md
.codex/skills/process/osk-engineering-reporting/SKILL.md

.claude/skills/osk-go-guide/SKILL.md
.claude/skills/architecture/osk-architecture-review/SKILL.md

Follow:

CLAUDE.md
AGENTS.md
Scope

Implement only:

GET
GET /api/v1/operations

Returns:

Manifest

generated directly from:

Registry.Manifest()

No duplicated DTO construction.

POST
POST /api/v1/calculations

Request:

{
  "operationId": "...",
  "operands": [...]
}

Response:

{
  "expression": "144 ÷ 12",
  "result": 12
}

Expression formatting belongs to the backend.

The frontend must simply display it.

Application flow

The endpoint should read naturally:

Lookup operation

↓

Validate

↓

Execute

↓

Build response

Not:

Switch(operationID)

No operation-specific dispatch outside the registry.

Errors

Implement a small error contract.

Suggested shape:

{
  "code": "...",
  "message": "..."
}

Cover:

unknown operation
validation failure
malformed request
internal evaluation failure

Use appropriate HTTP status codes.

Frontend integration

Replace:

mockCalculatorApi

with:

calculatorApi

The UI should require minimal changes.

Specifically:

remove hardcoded shortcut map;
consume shortcut from the manifest;
consume version from the manifest;
keep existing validation interpreter;
keep existing UX unchanged.
Documentation

Update:

README
API Contract
Engineering Log
Roadmap

Add:

docs/knowledge/api/API_CONTRACT.md

Include:

endpoints
request examples
response examples
validation errors
curl examples
Tests

Backend

manifest endpoint
calculate endpoint
validation failures
unknown operation
malformed JSON

Frontend

manifest fetch
calculation request
error rendering

Keep Playwright out of scope.

Validation

Run:

make fmt
make lint
make test
make coverage
make build
make verify

go test -race ./...

docker compose build

docker compose up

Verify manually:

GET /api/v1/operations

POST /api/v1/calculations
Deliverables
REST endpoints
HTTP client adapter
mock removal
API Contract
updated README
updated screenshots if necessary
Constraints

Do not:

redesign the frontend
introduce OpenAPI generation
introduce authentication
introduce persistence
introduce middleware beyond what already exists
implement Playwright
commit
Definition of Done
F1 remediation completed.
/api/v1/operations returns the registry manifest.
/api/v1/calculations executes through the immutable registry.
No operation-specific switches outside the registry.
Frontend consumes the live backend.
Mock adapter removed.
Existing UX preserved.
Error contract implemented.
API Contract documented.
Quality gates pass.
Docker stack works end-to-end.
No commits created.
