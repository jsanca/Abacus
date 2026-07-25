# ABA-003R — Backend Foundation Architecture & Go Review

**Status:** Complete
**Date:** 2026-07-24
**Reviewer:** Deep Pro
**Verdict:** PASS WITH OBSERVATIONS

## Scope

Reviewed the Go backend foundation delivered by ABA-003. Evaluated all ten review dimensions specified by the task: package organization, runtime ownership, configuration, HTTP transport, logging, observability boundary, Go idioms, Docker/tooling, testing, and future readiness.

| Reference | Source |
|-----------|--------|
| Task spec | `docs/engineering/agents-tasks/ABA-003—BackendFoundation-ExecutableRuntime.md` |
| Implementation report | `docs/engineering/agents-reports/ABA-003-backend-foundation.md` |
| Architecture decision | `docs/engineering/adr/0001-defer-production-telemetry.md` |
| Go guidance | `.opencode/skills/osk-go-guide/SKILL.md` |
| Conventions | `CLAUDE.md`, `AGENTS.md` |

## Validation Evidence

| Command | Result |
|---------|--------|
| `go test ./...` | All pass (app, config, transport) |
| `go test -race ./...` | All pass, no races |
| `golangci-lint run ./...` | 0 issues |
| `go build ./cmd/api` | Success |
| `make verify` | All checks passed |
| `docker compose config` | Valid |
| `docker compose build server` | Built successfully |
| `GET /health` | `{"status":"ok"}` (HTTP 200) |
| SIGTERM shutdown | Graceful: "shutdown signal received" → "graceful shutdown complete" |

## Architecture Assessment

### Execution Path

```
main.go
  └─ signal.NotifyContext (SIGINT, SIGTERM)
  └─ config.Load() → ServerConfig (env-driven, validated)
  └─ &observability.Observer{} (no-op)
  └─ transport.NewRouter(logger, observer, RouterConfig)
       ├─ middleware: RequestID → Recoverer → Timeout → CORS → requestLogging
       └─ GET /health → newHealthHandler → writeJSONResponse
  └─ &http.Server{Addr, Handler, timeouts}
  └─ app.NewRuntime(httpServer, observer, logger, shutdownTimeout)
       └─ Runtime.Run(ctx)
            ├─ goroutine: ListenAndServe → serverErrors chan (cap=1)
            ├─ select: serverErrors | ctx.Done
            ├─ Done: bounded Shutdown(shutdownCtx) via context.WithTimeout
            └─ observer.OnServerStop()
```

### Summary by Dimension

| Dimension | Verdict | Key Observations |
|-----------|---------|------------------|
| 1. Package organization | PASS | Clean, minimal, each package has a single responsibility. No `util`/`common`/`helpers`. |
| 2. Runtime ownership | PASS | Goroutine has explicit owner and stop path. Buffered channel prevents leaks. Shutdown sequence is correct. |
| 3. Configuration | PASS | `time.Duration` throughout. Fail-fast validation. Immutable after loading. |
| 4. HTTP transport | PASS | Chi middleware order is correct. JSON helper is simple and adequate. No business logic in transport. |
| 5. Logging | PASS | Structured `slog` throughout. No `fmt.Println`. Appropriate levels (INFO/ERROR). Request ID propagated. |
| 6. Observability boundary | PASS | Proportionate no-op Observer. Three focused hooks. No ornamental architecture. ADR 0001 documents the deferral. |
| 7. Go idioms | PASS (2 MINOR) | All exported symbols documented. `%w` error wrapping. No interfaces. Minor receiver naming and Observer receiver style observations. |
| 8. Docker & tooling | PASS (1 MINOR) | Multi-stage build, non-root user, small image. `make fmt` uses `gofmt` but lint expects `goimports` — minor inconsistency. |
| 9. Testing | PASS (1 MINOR) | Table-driven config tests. `httptest` for router. Runtime lifecycle tests. `time.Sleep` flake risk. |
| 10. Future readiness | PASS | Composition root (`run()`) is the natural wiring point for the operation registry, validation model, and calculator routes. No architectural changes needed. |

## Findings

### MINOR — AAR-003-001: `make fmt` uses `gofmt` but lint checks `goimports`

**Severity:** MINOR
**Location:** `Makefile:3-4`, `server/.golangci.yml:21-24`
**Evidence:** `make fmt` runs `gofmt -w .`. The golangci-lint config enables `goimports` as a formatter. A developer running `make fmt` may still fail `make lint` if import ordering is wrong.
**Impact:** Minor developer friction — running `make fmt` should produce code that passes `make lint`.
**Recommendation:** Change `make fmt` to use `goimports -w .` (requires `goimports` to be installed). Alternatively, document that `golangci-lint run --fix` is the authoritative formatter.

### MINOR — AAR-003-002: `time.Sleep` in runtime test creates flake risk

**Severity:** MINOR
**Location:** `server/internal/app/app_test.go:34`
**Evidence:** `time.Sleep(50 * time.Millisecond)` waits for the server goroutine to reach `ListenAndServe` before signalling cancellation. On a heavily loaded CI runner, 50ms may be insufficient.
**Impact:** Low — the test also has a 5-second timeout that catches the hang. If the sleep is too short, the "immediate context cancellation" test path may be exercised instead of the intended shutdown path.
**Recommendation:** Poll for server readiness (e.g., loop HTTP GET to the ephemeral port until it responds, with a short deadline) instead of sleeping. Or accept the pattern and rely on the generous test timeout.

### MINOR — AAR-003-003: `Observer` uses pointer receivers on a stateless struct

**Severity:** MINOR
**Location:** `server/internal/observability/observability.go:17-24`
**Evidence:** All three `Observer` methods use pointer receivers (`*Observer`) but the struct contains no fields. There is no state to share or mutate. Using pointer receivers forces callers to use `&observability.Observer{}` even though the zero value `observability.Observer{}` is equally valid.
**Impact:** Low — the code works correctly. Value receivers would be more idiomatic for a stateless struct and would eliminate the need for the `&` in the composition root. If future telemetry state is added (e.g., a metrics client), pointer receivers become correct; but the decision should be deferred until state exists.
**Recommendation:** Change receivers to value (`func (o Observer)`) until state justifies pointer receivers. This is consistent with the package doc comment: "The zero value is a ready-to-use no-op implementation."

### NOTE — N-003-001: Receiver name `runtime` is descriptive but unconventional

**Location:** `server/internal/app/app.go:43`
**Observation:** The `Run` method uses `runtime` as its receiver name. In idiomatic Go, short receiver names are conventional — `r` or `rt` would be more typical. The osk-go-guide says no abbreviations but does not prescribe receiver length — single-letter receivers are standard Go convention. The current name is clear and unambiguous. No action required unless the team standardizes receiver naming.

### NOTE — N-003-002: `NewRouter` parameter list will grow with calculator routes

**Location:** `server/internal/transport/router.go:24`
**Observation:** `NewRouter` currently accepts `(logger, observer, RouterConfig)`. When calculator routes are added (requiring an operation registry reference), the parameter count will grow. A `RouterDependencies` value type or a continuation of parameterized construction are both valid. No action required now — the current signature is proportionate.

## Positive Findings

1. **Flawless lifecycle ownership.** `Runtime.Run()` starts one goroutine for `ListenAndServe`, sends errors through a buffered channel (no leak), distingushes `http.ErrServerClosed` from real failures, and drains gracefully with bounded timeout. The shutdown context is correctly derived from `context.Background()`, not the cancelled signal context.

2. **Configuration is bulletproof.** `time.Duration` throughout (no untyped integers). `time.ParseDuration` for parsing. Fail-fast validation of all values. Immutable after `Load()`. Every error wraps with `%w`.

3. **Disciplined middleware ordering.** Chi middleware follows the correct stack: `RequestID` first (identity), `Recoverer` (safety net), `Timeout` (boundary), `CORS` (cross-origin policy), request logging last (observes the final response status from `WrapResponseWriter`).

4. **Proportionate observability boundary.** Three concrete no-op methods on a struct. Wired into real call sites. Required zero interfaces. The ADR documents the deliberate deferral. No ornamental architecture.

5. **`writeJSONResponse` handles the hard HTTP truth correctly.** Sets Content-Type before WriteHeader. Acknowledges the Go HTTP limitation that encoding failures after header-write can't change the status code. Logs the error without panicking.

6. **Clean composition root.** `main()` is 5 lines: error check, log fatal, os.Exit. `run()` is the composition function — load config, construct dependencies, wire, run. The Go idiom is respected precisely.

7. **Linter config is intentional.** Exactly 7 linters enabled with rationale. `goimports` as formatter. Exported-symbol doc comments enforced via `revive`. No blanket linter set that generates noise.

8. **Docker image is production-shaped.** Multi-stage build, non-root user, stripped binary (`-ldflags="-s -w"`), `CGO_ENABLED=0`, Alpine runtime. Compose service uses explicit environment variables (no .env file leakage).

## Conclusion

The backend foundation is well-crafted Go. Lifecycle ownership is explicit, configuration is robust, the observability boundary is proportionate, and all quality gates pass. The three MINOR findings (Makefile formatter inconsistency, test `time.Sleep` fragility, Observer receiver style) and two observations (receiver naming, future parameter growth) do not indicate correctness or architecture issues.

**The foundation is ready for the Operation Registry (next roadmap initiative).** No structural changes are required — the `run()` composition root is the natural point to wire new dependencies, and the Chi router is already configured for additional routes.
