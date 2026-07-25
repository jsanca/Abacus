# ABA-003 — Backend Foundation: Executable Runtime

| Field | Value |
|---|---|
| Task | ABA-003 — Backend Foundation & Executable Runtime |
| Status | Complete |
| Author | Clio |
| Date | 2026-07-24 |
| Report type | Implementation |

## Summary

Delivered the executable Go backend foundation for Abacus. The server builds, starts, and responds to `GET /health`. SIGINT and SIGTERM trigger bounded graceful shutdown. All quality gates pass.

## Files Created or Changed

| File | Change |
|---|---|
| `server/go.mod` | Go module `github.com/jsanca/abacus/server`, chi/v5 and cors dependencies |
| `server/go.sum` | Resolved dependency checksums |
| `server/cmd/api/main.go` | Entry point; signal-aware context; delegates to `run()` |
| `server/internal/config/config.go` | Environment-driven config; `time.Duration` values; fail-fast validation |
| `server/internal/config/config_test.go` | Table-driven tests for defaults, overrides, and invalid durations |
| `server/internal/observability/observability.go` | Concrete no-op `Observer`; telemetry extension boundary |
| `server/internal/transport/json.go` | `writeJSONResponse` helper; JSON Content-Type and encoding |
| `server/internal/transport/handler.go` | `newHealthHandler` — `GET /health` → `{"status":"ok"}` |
| `server/internal/transport/router.go` | Chi router; RequestID, Recoverer, Timeout, CORS, request logging middleware |
| `server/internal/transport/handler_test.go` | `httptest` tests: 200 health, 405 method-not-allowed, 404 unknown routes |
| `server/internal/app/app.go` | `Runtime.Run`: goroutine lifecycle, signal context, graceful `Shutdown` |
| `server/internal/app/app_test.go` | Context cancellation triggers clean shutdown; pre-cancelled context exits cleanly |
| `server/.golangci.yml` | golangci-lint v2 config: errcheck, govet, staticcheck, ineffassign, unused, misspell, revive, goimports |
| `server/Dockerfile` | Multi-stage build; non-root `appuser`; alpine runtime; port 8080 |
| `server/bin/api` | Compiled binary (local build artifact, not committed) |
| `docker-compose.yml` | `server` service with all environment variables |
| `Makefile` | `fmt`, `lint`, `test`, `coverage`, `build`, `run-server`, `docker-up`, `docker-down`, `verify` |
| `docs/engineering/adr/0001-defer-production-telemetry.md` | ADR: accepted decision to defer production telemetry |
| `docs/engineering/ROADMAP.md` | ABA-003 marked Complete |
| `docs/engineering/ENGINEERING_LOG.md` | ABA-003 artifacts and telemetry deferral logged |
| `docs/engineering/TECHNICAL_DEBT.md` | TD-001: production telemetry deferred |

## Architecture Implemented

```
main
  └─ signal.NotifyContext (SIGINT/SIGTERM)
  └─ config.Load()
  └─ observability.Observer (no-op)
  └─ transport.NewRouter(logger, observer, RouterConfig)
       middleware: RequestID → Recoverer → Timeout → CORS → requestLogging
       GET /health → newHealthHandler
  └─ http.Server{Addr, Handler, timeouts}
  └─ app.NewRuntime(httpServer, observer, logger, shutdownTimeout)
       └─ Runtime.Run(ctx)
            goroutine: ListenAndServe → serverErrors chan
            select: serverErrors | ctx.Done
            on Done: httpServer.Shutdown(shutdownCtx with timeout)
```

**Lifecycle ownership:** The goroutine that calls `ListenAndServe` sends to a buffered channel (`cap=1`) and exits when the server stops. `Runtime.Run` owns the shutdown sequence. No unmanaged goroutines.

## Commands Executed and Results

```
make build
  → cd server && go build -o bin/api ./cmd/api    ✓

make test
  → ok  internal/app          1.173s
  → ok  internal/config       2.191s
  → ok  internal/transport    3.247s              ✓

make coverage
  → internal/app        77.8%
  → internal/config     84.1%
  → internal/transport  95.5%
  → total               73.5%                     ✓

make verify
  → fmt check: no issues
  → golangci-lint: 0 issues
  → go test -race: all pass
  → build: ok                                     ✓

docker compose config   → valid                   ✓
docker compose build server → image built         ✓

Manual: GET /health
  → HTTP 200  {"status":"ok"}
  → structured slog JSON logged                   ✓
```

## Docker Verification

- Multi-stage build: `golang:1.26-alpine` builder → `alpine:3.22` runtime.
- Non-root: `appuser` created with `adduser -D`.
- Image size: ~10 MB (stripped binary on alpine).
- `docker compose build server` completed without errors.

## Assumptions

- Module path `github.com/jsanca/abacus/server` follows the project's GitHub identity. No remote is configured yet; this is consistent with the bootstrap state.
- `gofmt` is used for `make fmt` (recursive directory mode). `goimports` is enforced through golangci-lint's formatter pipeline.
- golangci-lint v2 (`2.12.2`) installed via Homebrew; `goimports` is a v2 formatter (not a linter).

## Accepted Debt

- **TD-001:** Production telemetry deferred. `observability.Observer` no-op retained as the extension boundary. See ADR 0001.

## Remaining Work

None for this task. Next task is the Operation Registry (ABA-005 or equivalent).

## Non-goals Confirmed Not Implemented

- No operation registry, arithmetic functions, or calculation endpoints.
- No `GET /operations` or manifest endpoint.
- No database, authentication, or frontend integration.
- No OpenAPI, Kubernetes, or Terraform.

## No-commit Confirmation

No commits were created.
