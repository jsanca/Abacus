# Engineering Team

## Clio — Main Developer

Responsibilities:

- Backend implementation
- Frontend implementation
- Unit tests
- Docker

## Deep Pro — Architecture Reviewer

Responsibilities:

- Go idiomatic review
- React review
- Design review
- Trade-off review
- Final implementation review

## Elito — Knowledge Curator

Responsibilities:

- Documentation
- ADRs
- README
- Engineering Process
- Prompt curation
- Technical debt tracking
- Knowledge consistency

# Project Overview

React + TypeScript frontend (Vite) and Go (Chi) backend calculator service. The backend operation registry is the single source of truth: it produces the operation manifest consumed by the client, validates requests, and executes operations.

**Current state:** The frontend is implemented against a mock API boundary. The Go backend is planned but not yet built (see `docs/engineering/ROADMAP.md`).

# Commands

All client commands run from `client/`:

| Command | What it does |
|---------|-------------|
| `npm run dev` | Vite dev server |
| `npm run lint` | ESLint |
| `npm run test` | Vitest (jsdom) |
| `npm run coverage` | Vitest with coverage |
| `npm run build` | `tsc -b && vite build` |
| `npm run verify` | lint → test → coverage → build (run this before committing) |

Server commands (`make lint`, `make test`, `make verify`) will be added when the Go backend is implemented.

# Architecture

- **Frontend entry:** `client/src/main.tsx` → `App.tsx` → `features/calculator/components/Calculator.tsx`
- **Frontend feature structure:** Each feature under `client/src/features/<name>/` has `api/`, `model/`, `components/`, and `test/` directories.
- **Mock API boundary:** `client/src/features/calculator/api/mockCalculatorApi.ts` — models the future backend endpoints. The manifest and calculation types are in `api/types.ts`.
- **Validation model:** `client/src/features/calculator/model/validation.ts` — shared declarative validation consuming the operation manifest.
- **Test setup:** Vitest with jsdom environment, setup file at `client/src/test/setup.ts` (imports `@testing-library/jest-dom/vitest`).
- **Backend (planned):** Go service under `server/`, Chi router, immutable operation registry, function-based execution, graceful shutdown.

# Conventions

## Naming
- Descriptive names; no abbreviations or one-letter variables except trivial loop indices.

## Go (when backend is built)
- Apply the Go engineering skill at `.opencode/skills/osk-go-guide/SKILL.md` — it covers naming, packages, ownership, error handling, concurrency, testing, and anti-patterns.
- Registry must be concrete and immutable after startup; no interface unless a consumer truly needs it.
- Function-based operation execution (`type ExecuteFunc func(...)`) over interface hierarchies for stateless arithmetic.
- Structured logging via `log/slog`; no `fmt.Println`.
- Graceful shutdown is a BLOCKER requirement: handle OS signals, stop accepting new requests, drain in-flight requests, exit.

## React / TypeScript
- `verbatimModuleSyntax: true` — use `import type` for type-only imports.
- Strict TypeScript mode; `noUnusedLocals` and `noUnusedParameters` enabled.

## Documentation
- Exported Go symbols require Go doc comments.
- Architecture decisions in `docs/engineering/adr/`.
- Long-lived knowledge in `docs/knowledge/`.
- Engineering process, roadmap, and debt register in `docs/engineering/`.

## Errors
- Prefer explicit domain errors over anonymous string comparisons.
- Wrap errors with `%w` so callers can inspect the chain.

# OpenCode Skills

Project skills live in `.opencode/skills/`. Agents should consult them when their descriptions match:

| Skill | Purpose |
|-------|---------|
| `osk-go-guide` | Go implementation and review guidance |
| `osk-architecture-review` | System structure and boundary validation |
| `osk-knowledge-curator` | Documentation integrity and governance |
| `osk-engineering-reporting` | Durable implementation and review records |
| `osk-execution-timebox` | Execution bounding with recovery checkpoints |

# Testing

- Client: Vitest + React Testing Library + jsdom. Use `npm run test` (single run) or `npm run coverage`.
- Go (planned): Table-driven tests, `net/http/httptest` for HTTP, `-race` flag required in CI.
