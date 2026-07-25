# ABA-003R — Backend Foundation Architecture & Go Review

**Status:** Complete | **Date:** 2026-07-24 | **Reviewer:** Deep Pro | **Verdict:** PASS WITH OBSERVATIONS

## Scope

Architecture and Go-idiom review of the ABA-003 backend foundation. Evaluated package organization, runtime ownership, configuration, HTTP transport, logging, observability boundary, Go idioms, Docker/tooling, testing, and future readiness. No calculator-domain or operation-registry review.

## Validation Evidence

| Command | Result |
|---------|--------|
| `go test ./...` | All pass (3 packages) |
| `go test -race ./...` | All pass, no data races |
| `golangci-lint run ./...` | 0 issues |
| `go build ./cmd/api` | Success |
| `make verify` | All checks passed |
| `docker compose config` | Valid |
| `docker compose build server` | Built successfully |
| `GET /health` | `{"status":"ok"}` (HTTP 200) |
| Graceful shutdown (SIGTERM) | "shutdown signal received" → "graceful shutdown complete" |

## Dimension Summary

| # | Dimension | Verdict |
|---|-----------|---------|
| 1 | Package organization | PASS |
| 2 | Runtime ownership | PASS |
| 3 | Configuration | PASS |
| 4 | HTTP transport | PASS |
| 5 | Logging | PASS |
| 6 | Observability boundary | PASS |
| 7 | Go idioms | PASS (2 MINOR) |
| 8 | Docker & tooling | PASS (1 MINOR) |
| 9 | Testing | PASS (1 MINOR) |
| 10 | Future readiness | PASS |

## Findings

### MINOR — AAR-003-001: `make fmt` uses `gofmt` but lint checks `goimports`

**Location:** `Makefile:3-4`, `server/.golangci.yml:21-24`
**Finding:** `make fmt` runs `gofmt -w .` while golangci-lint enforces `goimports`. A developer running `make fmt` may still fail lint if import ordering is incorrect.
**Recommendation:** Change `make fmt` to `goimports -w .`, or document that `golangci-lint run --fix` is the authoritative formatter.

### MINOR — AAR-003-002: `time.Sleep(50ms)` in runtime test

**Location:** `server/internal/app/app_test.go:34`
**Finding:** The test sleeps 50ms to wait for `ListenAndServe` before signalling shutdown. On a slow CI runner this could race.
**Recommendation:** Poll for server readiness (HTTP GET to ephemeral port) instead of a fixed sleep. The generous 5s test timeout provides a safety net in practice.

### MINOR — AAR-003-003: `Observer` uses pointer receivers on a stateless struct

**Location:** `server/internal/observability/observability.go:17-24`
**Finding:** All Observer methods use pointer receivers but the struct has no fields. This forces callers to use `&observability.Observer{}`. The package doc says "The zero value is a ready-to-use no-op implementation" — value receivers would match this contract.
**Recommendation:** Switch to value receivers until state justifies pointer receivers.

## Positive Findings

1. **Flawless lifecycle ownership** — one goroutine for `ListenAndServe`, buffered error channel (no leak), `http.ErrServerClosed` distinguished from real failures, bounded drain with shutdown timeout derived from `context.Background()` (not the cancelled signal context).
2. **Config is `time.Duration` throughout** — no untyped integers, `time.ParseDuration` for parsing, fail-fast `validate()`, immutable after loading, every error wrapped with `%w`.
3. **Correct middleware ordering** — RequestID (identity) → Recoverer (safety) → Timeout (boundary) → CORS (policy) → request logging (observation, uses `WrapResponseWriter` for final status).
4. **Proportionate observability** — three concrete no-op hooks on a struct, wired at real call sites, zero interfaces. ADR 0001 documents the deferral.
5. **Clean composition root** — `main()` is 5 lines (error/slog/os.Exit). `run()` is the composition function. Exactly matches the Go idiom.
6. **Intentional linter config** — 7 linters + goimports formatter. `revive` enforces exported doc comments. No blanket linter creep.
7. **Production-shaped Docker image** — multi-stage, non-root user, stripped binary, `CGO_ENABLED=0`, Alpine. Compose uses explicit env vars.

## Conclusion

The backend foundation is well-crafted, idiomatic Go. No BLOCKER or MAJOR findings. All quality gates pass. The three MINOR observations do not affect correctness or architecture.

**Ready for the Operation Registry (next roadmap initiative).** The `run()` composition root is the natural wiring point, and `NewRouter()` is ready for additional routes.
