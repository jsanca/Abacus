# Project Overview

React + TypeScript frontend (Vite) and Go (Chi) backend calculator service. The backend operation registry is the single source of truth: it produces the operation manifest consumed by the client, validates requests, and executes operations.

Also see `CLAUDE.md` for additional conventions.

# Commands

**All client commands run from `client/`:**

| Command | What it does |
|---------|-------------|
| `npm run dev` | Vite dev server |
| `npm run lint` | ESLint |
| `npm run test` | Vitest (jsdom, single run) |
| `npm run coverage` | Vitest with coverage |
| `npm run build` | `tsc -b && vite build` |
| `npm run verify` | lint → build → coverage (run this before committing) |

**All server commands run from repo root:**

| Command | What it does |
|---------|-------------|
| `make fmt` | `gofmt -w .` in `server/` |
| `make lint` | `golangci-lint run ./...` in `server/` |
| `make test` | `go test ./...` in `server/` |
| `make coverage` | `go test -coverprofile` + HTML coverage in `server/` |
| `make build` | `go build -o bin/api ./cmd/api` in `server/` |
| `make run-server` | `go run ./cmd/api` in `server/` |
| `make verify` | gofmt check → lint → `go test -race` → build (run this before committing) |
| `make docker-up` | `docker compose up -d` |
| `make docker-down` | `docker compose down` |

## Running a single Go test

```bash
cd server && go test -run TestName ./internal/calculator/
```

## Gotchas

- `npm run coverage` uses `--pool=forks --no-file-parallelism` — these flags are required for vitest coverage in jsdom.
- `make fmt` only runs `gofmt` (not `goimports`). The linter checks import ordering via `goimports`, so `make fmt` won't fix import issues — the `make verify` step will catch them.
- There is no CI or pre-commit hooks. Always run `make verify` and `npm run verify` before committing.

# Architecture

- **Frontend entry:** `client/src/main.tsx` → `App.tsx` → `features/calculator/components/Calculator.tsx`
- **Frontend features:** `client/src/features/<name>/` with `api/`, `model/`, `components/`, `test/` subdirectories.
- **Frontend API:** `client/src/features/calculator/api/calculatorApi.ts` — the calculator consumes the live `/api/v1` manifest and calculation endpoints.
- **Go module:** `github.com/jsanca/abacus/server`
- **Backend entry:** `server/cmd/api/main.go` → `internal/app/` (lifecycle) → `internal/transport/` (Chi router + middleware) → `internal/calculator/` (registry + operations)
- **Observability:** `internal/observability/` is a no-op boundary — zero-value `Observer{}` is safe to use. No real telemetry yet.
- **Docker Compose:** frontend on `:3000`, backend on `:8080`. Uses `docker compose up --build` (or `make docker-up`).

# Conventions

## TypeScript
- `verbatimModuleSyntax: true` — use `import type` for type-only imports.
- Full strict mode; `noUnusedLocals` and `noUnusedParameters` enabled.

## Go
- Use structured logging (`log/slog`); never `fmt.Println`.
- Function-based operation execution (`type ExecuteFunc func(...)`) over interface hierarchies.
- Registry is concrete and immutable after startup; no interface unless a consumer truly needs it.
- Graceful shutdown is a blocker: handle OS signals, stop accepting new requests, drain in-flight requests, exit.
- Exported symbols require doc comments (enforced by `revive` in `.golangci.yml`).
- Linter config at `server/.golangci.yml`: `goimports` formatter, `revive` with exported-name checking.

## Testing
- Go: table-driven tests, `net/http/httptest` for HTTP, `-race` flag required.
- Client: Vitest + React Testing Library + jsdom. Setup file at `client/src/test/setup.ts`.

## Errors
- Use explicit domain errors, not anonymous string comparisons.
- Wrap errors with `%w` for chain inspection.

# OpenCode Skills

Project skills in `.opencode/skills/`:

| Skill | Purpose |
|-------|---------|
| `osk-go-guide` | Go implementation and review guidance |
| `osk-architecture-review` | System structure and boundary validation |
| `osk-knowledge-curator` | Documentation integrity and governance |
| `osk-engineering-reporting` | Durable implementation and review records |
| `osk-execution-timebox` | Execution bounding with recovery checkpoints |
