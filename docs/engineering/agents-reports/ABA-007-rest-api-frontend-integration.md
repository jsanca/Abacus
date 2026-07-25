# ABA-007 — REST API & Frontend Integration: Agent Report

**Date:** 2026-07-25
**Status:** Complete
**Agent:** Clio (claude-sonnet-4-6)

---

## Precondition: F1 remediation

`copyExpression` default branch changed from silent alias (`return expression`) to:

```go
default:
    panic(fmt.Sprintf("copyExpression: unsupported expression type %T", expression))
```

This preserves the registry's deep immutability guarantee and makes misuse fail loudly at development time.

---

## Backend: REST endpoints

### `server/internal/transport/calculator.go` (new)

Two handlers and supporting types:

**`GET /api/v1/operations`** — `newOperationsHandler`
- Returns `registry.Manifest()` as JSON.
- No DTO translation; the manifest type serializes directly.

**`POST /api/v1/calculations`** — `newCalculateHandler`

Application flow, exactly:
```
Lookup operation → Validate → Execute → Build response
```

No operation-specific dispatch. The registry is the only place that knows which operation to run.

**Error contract:**

| HTTP | `code` | Condition |
|---|---|---|
| 400 | `malformed_request` | Invalid JSON body |
| 404 | `unknown_operation` | Unknown `operationId` |
| 422 | `validation_failed` | Operation constraint violated |
| 500 | `calculation_failed` | Internal evaluation or execution error |

**Expression formatting:**

Backend-owned per task requirements:
- Unary: `√9`
- Binary: `144 ÷ 12`

`formatFloat` uses `strconv.FormatFloat(value, 'f', -1, 64)`: integer-valued floats render without decimal point (`12.0` → `"12"`).

### `server/internal/transport/router.go` (modified)

- `NewRouter` accepts `*calculator.Registry` as a new parameter.
- Routes registered under `/api/v1/`:
  - `GET /operations`
  - `POST /calculations`

### `server/cmd/api/main.go` (modified)

Registry is now a named dependency wired to the router:

```go
registry, err := calculator.NewDefaultRegistry()
// ...
router := transport.NewRouter(logger, observer, registry, transport.RouterConfig{...})
```

The old discard pattern `if _, err := ...` is replaced.

### Backend tests added (`handler_test.go`)

- `TestOperationsEndpoint_returnsManifest` — 200, JSON, 7 operations, version field
- `TestCalculateEndpoint_binaryOperation` — addition, expression `"3 + 4"`, result `7`
- `TestCalculateEndpoint_unaryOperation` — square-root, expression `"√9"`, result `3`
- `TestCalculateEndpoint_unknownOperation` — 404, code `unknown_operation`
- `TestCalculateEndpoint_validationFailure_divisionByZero` — 422, code `validation_failed`
- `TestCalculateEndpoint_malformedJSON` — 400, code `malformed_request`

---

## Frontend: HTTP adapter and mock removal

### `client/src/features/calculator/api/calculatorApi.ts` (new)

Real `fetch`-based adapter:
- `getOperations()` → `GET /api/v1/operations`
- `calculate(request)` → `POST /api/v1/calculations`

Network failures are caught and rethrown as typed `ApiError` objects with appropriate codes. Non-OK responses parse the error body or fall back to `service_unavailable`.

### `client/src/features/calculator/api/mockCalculatorApi.ts` (deleted)

Mock adapter removed entirely. No code imports from it.

### `client/src/features/calculator/api/types.ts` (modified)

- `OperationDefinition`: added `shortcut: string`
- `OperationManifest`: added `version: string`
- `ApiError.code`: added `'calculation_failed'`

### `client/src/features/calculator/components/Calculator.tsx` (modified)

- Import changed from `mockCalculatorApi` to `calculatorApi`.
- Hardcoded `shortcutOperationIds` map removed.
- `useMemo` computes `shortcutMap` from manifest operations: `{ [op.shortcut]: op.id }`.
- Keyboard handler uses `shortcutMap[event.key]` instead of the static constant.
- `shortcutMap` added to keyboard effect dependencies.

### `client/src/features/calculator/test/Calculator.test.tsx` (rewritten)

- `vi.mock('../api/calculatorApi')` replaces `configureMockScenario`.
- `beforeEach`: `vi.resetAllMocks()` + default resolved manifest and calculation mocks.
- Tests that trigger calculations set specific `mockResolvedValue` per test.
- Service error test uses `mockRejectedValue`.
- Retry test uses `mockRejectedValueOnce` + `mockResolvedValueOnce` for the two-call sequence.
- Mock manifest includes `shortcut` and `version` fields to match the live backend contract.

### `client/vite.config.ts` (modified)

Dev server proxy added:
```ts
server: { proxy: { '/api': 'http://localhost:8080' } }
```

### `client/nginx.conf` (modified)

API proxy block added before the SPA fallback:
```nginx
location /api/ {
    proxy_pass http://server:8080;
    proxy_set_header Host $host;
    ...
}
```

### `docker-compose.yml` (modified)

`client` service gains `depends_on: server`.

---

## API contract

Documented at `docs/knowledge/api/API_CONTRACT.md`: endpoints, request/response examples, validation expression schema, error codes, and curl examples.

---

## Verification

```
make verify                        → All checks passed
go test -race ./...                → All packages green
cd client && npm test              → 13/13 tests passed
docker compose build               → Both images built
docker compose up + manual curl    → All endpoints verified
```

Manual verification:
- `GET /api/v1/operations` — 200, manifest with 7 operations
- `POST /api/v1/calculations` `{"operationId":"division","operands":[144,12]}` → `{"expression":"144 ÷ 12","result":12}`
- `POST /api/v1/calculations` `{"operationId":"division","operands":[10,0]}` → `{"code":"validation_failed","message":"The divisor must not be zero."}`
- `GET http://localhost:3000/api/v1/operations` (via nginx proxy) — 200 ✓
- `POST http://localhost:3000/api/v1/calculations` (via nginx proxy) — 200 ✓

---

## No commit was created.
