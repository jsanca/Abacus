ABA-003 — Backend Foundation & Executable Runtime
Field	Value
Status	Ready
Owner	Clio
Role	Main Developer — Go Backend Implementation
Target	20–30 minutes
Hard stop	45 minutes
Commit	No
Repository	/Users/jonathan/code/Abacus
Objective

Establish the executable Go backend foundation for Abacus, including application composition, configuration, Chi routing, HTTP runtime lifecycle, structured logging, a minimal observability boundary, development tooling, and initial Docker execution.

This task creates the platform on which the operation registry, declarative validation model, and calculator API will be implemented.

It must not implement the operation catalogue, calculation behavior, or operation manifest yet.

Required skills

Apply:

.claude/skills/osk-go-guide/SKILL.md
.claude/skills/process/osk-execution-timebox/SKILL.md
.claude/skills/process/osk-engineering-reporting/SKILL.md
.claude/skills/architecture/osk-architecture-review/SKILL.md where applicable

Follow repository instructions in CLAUDE.md and AGENTS.md.

Architectural intent

The backend must demonstrate disciplined ownership of application resources:

main
  ↓
configuration
  ↓
dependency composition
  ↓
HTTP router
  ↓
application runtime
  ↓
Run
  ↓
SIGINT / SIGTERM
  ↓
graceful shutdown

The component that creates a resource must make its ownership and shutdown responsibility explicit.

Scope
1. Initialize the Go service

Create one Go module under:

server/

Use a single go.mod.

Do not create multiple Go modules or speculative packages.

Suggested proportional structure:

server/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── app/
│   ├── config/
│   ├── observability/
│   └── transport/
│       └── http/
├── go.mod
└── go.sum

The exact internal file distribution may evolve if a simpler idiomatic structure is clearer.

2. Configuration

Implement explicit environment-driven configuration for the minimum required settings:

HTTP port or listen address;
graceful shutdown timeout;
HTTP server timeouts;
allowed frontend origin if CORS is needed for local development.

Configuration must:

use descriptive names;
provide sensible local defaults;
validate invalid values at startup;
fail fast with contextual errors;
use time.Duration, not untyped integers, for time values.

Do not introduce a configuration framework unless clearly justified.

3. Application composition

Keep dependency construction explicit.

Prefer a composition function such as:

func run() error

or an application builder with equivalent clarity.

main() should remain small and own process exit behavior.

Avoid:

package-level mutable state;
service locators;
reflection-based dependency injection;
speculative interfaces.
4. HTTP router

Use Chi.

Create an executable router with at least:

GET /health

Expected response:

{
  "status": "ok"
}

Use a consistent JSON writer and error response mechanism that can later support the calculator API.

Do not add calculator routes yet.

5. Middleware

Configure only proportional middleware:

request ID;
panic recovery;
structured request logging;
request timeout;
CORS where local frontend/backend separation requires it.

Consider Chi's maintained middleware before writing custom alternatives.

Do not add authentication, rate limiting, tracing infrastructure, or unrelated middleware.

6. Runtime and graceful shutdown

Implement an application runtime that owns the HTTP server.

Required lifecycle behavior:

start the HTTP server;
distinguish expected http.ErrServerClosed from unexpected runtime failures;
listen for SIGINT and SIGTERM;
stop accepting new requests;
allow in-flight requests to drain through http.Server.Shutdown;
enforce the configured shutdown timeout;
return contextual errors;
emit structured lifecycle logs.

Do not start unmanaged goroutines.

Every goroutine introduced must have:

an owner;
a stop mechanism;
a wait path.
7. Structured logging

Use log/slog.

Do not use fmt.Println, log.Printf, or anonymous application logging.

Log at least:

application startup;
listen address;
shutdown signal received;
graceful shutdown completion;
unexpected runtime failure.

Do not log sensitive request content.

8. Observability boundary

Introduce a deliberately small observability or telemetry port with a no-op implementation.

It should represent a real extension boundary, such as runtime or HTTP lifecycle observations, without creating ornamental architecture.

Requirements:

concrete no-op implementation;
no external telemetry dependency;
no unused generic hooks;
no elaborate metrics vocabulary;
document why production telemetry is deferred.

Prefer consumer-owned interfaces only when the consumer actually needs one. A small concrete observer with no-op functions is acceptable if more idiomatic.

9. Docker foundation

Create a multi-stage backend Dockerfile.

Requirements:

build the Go binary in a builder stage;
run as a non-root user where practical;
produce a small runtime image;
expose or document the backend port;
no development tooling in the final image.

Create or update root Compose configuration with the backend service executable.

The frontend service may remain a placeholder until ABA-004 produces a buildable client.

Use the repository’s existing Compose filename rather than creating a second convention.

10. Makefile and verification

Create or update a root Makefile with at least:

make fmt
make lint
make test
make coverage
make build
make run-server
make docker-up
make docker-down
make verify

make verify should run the meaningful backend quality gates in deterministic order:

format check
lint
tests
race detector where appropriate
build

Use golangci-lint and provide a small, intentional configuration.

Do not enable a huge arbitrary linter set that generates noise.

11. Tests

Add tests for:

configuration defaults and invalid configuration;
health endpoint using net/http/httptest;
runtime handling of expected server closure where reasonably testable;
any custom middleware or JSON helper introduced.

Prefer table-driven tests when several equivalent cases exist.

Do not chase an arbitrary coverage percentage in this task.

Naming and documentation rules
Use descriptive names.
Avoid p, svc, mgr, cfg, tmp, ret, or similar compressed identifiers.
Preserve conventional short identifiers only where genuinely idiomatic, such as ctx, err, r, and w in standard HTTP signatures.
Every exported Go symbol requires a Go documentation comment beginning with the symbol name.
Comments must explain intent or trade-offs, not restate code.
Explicit non-goals

Do not implement:

operation registry;
arithmetic functions;
declarative validation;
GET /operations;
calculation endpoints;
frontend integration;
OpenAPI;
database;
authentication;
production telemetry backend;
Kubernetes or Terraform.
Documentation

Update:

docs/engineering/ROADMAP.md
docs/engineering/ENGINEERING_LOG.md
docs/engineering/TECHNICAL_DEBT.md only for accepted debt
ADRs only when a durable material trade-off emerges

The engineering log must link durable artifacts and must not reconstruct conversation history.

Record the accepted decision that production telemetry is deferred while an observability boundary is retained, unless already documented.

Validation commands

Run and report the exact results of:

make fmt
make lint
make test
make coverage
make build
make verify
docker compose config
docker compose build server

Also run the service and manually verify:

GET /health
Deliverables

Expected deliverables include:

executable Go service;
Chi router;
health endpoint;
configuration;
graceful lifecycle management;
structured logging;
minimal no-op observability boundary;
backend Dockerfile;
Compose backend service;
Makefile verification commands;
tests;
updated durable documentation.
Checkpoint and report

Create:

docs/engineering/agents-reports/ABA-003-backend-foundation.md

If the hard stop is reached before completion, stop and create:

docs/engineering/agents-reports/CHECKPOINT-ABA-003.md

The report or checkpoint must include:

status;
files created or changed;
architecture implemented;
lifecycle ownership;
commands executed and results;
Docker verification;
assumptions;
accepted debt;
remaining work;
no-commit confirmation.
Definition of Done
The server builds and starts.
/health returns JSON successfully.
SIGINT/SIGTERM trigger bounded graceful shutdown.
Logs are structured through slog.
No unmanaged goroutines exist.
Configuration is validated at startup.
The observability boundary has a meaningful no-op implementation.
Docker can build and run the server.
Lint, tests, coverage generation, and build are available through Make.
make verify succeeds.
Documentation and engineering log are current.
No operation or calculator-domain behavior was prematurely implemented.
No commits were created.
