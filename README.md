# Abacus

## Project Overview

Abacus is an operation-based calculator service created as a software engineering take-home exercise.

## Objectives

- Deliver a clear, maintainable calculator experience across desktop and mobile.
- Keep the operation catalogue, validation, and execution model consistent.
- Demonstrate proportionate engineering decisions and a verifiable delivery process.

## Architecture Overview

The project will use a React + TypeScript client and a Go HTTP service. The backend operation registry is the single source of truth: it produces the operation manifest consumed by the client, validates requests through a shared declarative model, and executes operations with functions. Local services will run with Docker Compose and the backend will support graceful shutdown.

## Technology Stack

- Frontend: React and TypeScript
- Backend: Go with the Chi HTTP router
- Local execution: Docker Compose
- Logging: structured logging with `log/slog`

## Repository Structure

```text
client/             React + TypeScript frontend
server/             Go backend
docs/engineering/   Process, prompts, debt register, and ADRs
docs/knowledge/     Long-lived engineering knowledge
.claude/skills/     Claude-oriented project skills
.codex/skills/      Codex-oriented project skills
.opencode/skills/   OpenCode-oriented project skills
```

## Running the Project

Start the complete local stack with:

```sh
docker compose up --build
```

The calculator is available at <http://localhost:3000>. The backend API is available at <http://localhost:8080>.

| URL | Description |
|---|---|
| `http://localhost:3000` | React calculator frontend |
| `http://localhost:8080/health` | Backend health check |
| `http://localhost:8080/api/v1/operations` | Operation manifest |
| `http://localhost:8080/api/v1/calculations` | Execute a calculation |

Stop the stack with `docker compose down` (or `make docker-down`).

For local development with live reload, run the backend (`make run-server`) and the frontend (`npm run dev` in `client/`) separately. The Vite dev server proxies `/api` to `localhost:8080`.

## Testing

**Backend:**
```sh
make test          # unit tests
make verify        # fmt + lint + race tests + build
```

**Frontend:**
```sh
cd client && npm test
cd client && npm run verify   # lint + build + coverage
```

## API

See [API Contract](docs/knowledge/api/API_CONTRACT.md) for endpoint documentation, request/response shapes, validation error codes, and curl examples.

## Docker

Docker Compose runs the production React static build and the Go backend together. The client is served by unprivileged Nginx with SPA fallback and API proxy; the backend serves the operation manifest and executes calculations via the immutable registry.

## Engineering Process

The agreed workflow is documented in [Engineering Process](docs/engineering/ENGINEERING_PROCESS.md).

## AI-assisted Development

AI-assisted work remains subject to human review and project conventions. When prompt usage materially influences an engineering artifact or decision, record it in the [Engineering Log](docs/engineering/ENGINEERING_LOG.md).

## Design Decisions

The initial architecture decisions are:

- React + TypeScript frontend and Go backend using Chi.
- Docker Compose for local execution and a runtime with graceful shutdown.
- An immutable operation registry as the single source of truth.
- A backend-generated operation manifest with shared declarative validation.
- Function-based operation execution.
- A responsive, mobile-capable, operation-based calculator UI rather than a keypad-based UI.
- Architecture kept proportional to the assignment scope.

## Future Improvements

Future enhancements and deferred decisions will be tracked in the [technical debt register](docs/engineering/TECHNICAL_DEBT.md) and ADRs.
