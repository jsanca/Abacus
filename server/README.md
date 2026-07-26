# Server

The Abacus server is a Go + Chi HTTP API. At startup it creates the immutable calculator operation registry, serves its manifest at `GET /api/v1/operations`, and validates, executes, and formats calculations at `POST /api/v1/calculations`.

From the repository root:

```sh
make run-server  # start the API on :8080
make test        # run server tests
make verify      # format check, lint, race tests, and build
```

Use `GET /health` for a readiness probe. The full request/response schema is in the [API contract](../docs/knowledge/api/API_CONTRACT.md).
