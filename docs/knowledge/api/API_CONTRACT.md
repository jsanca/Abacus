# API Contract

Base path: `/api/v1`

---

## GET /api/v1/operations

Returns the operation manifest. The frontend uses this on load to render the calculator UI and interpret validation rules client-side.

### Response — 200 OK

```json
{
  "version": "1",
  "defaultOperationId": "addition",
  "operations": [
    {
      "id": "addition",
      "name": "Addition",
      "symbol": "+",
      "shortcut": "+",
      "arity": 2,
      "operands": [
        { "id": "first", "label": "First number", "placeholder": "0" },
        { "id": "second", "label": "Second number", "placeholder": "0" }
      ],
      "validations": []
    },
    {
      "id": "division",
      "name": "Division",
      "symbol": "÷",
      "shortcut": "/",
      "arity": 2,
      "operands": [
        { "id": "first", "label": "Dividend", "placeholder": "0" },
        { "id": "second", "label": "Divisor", "placeholder": "0" }
      ],
      "validations": [
        {
          "message": "The divisor must not be zero.",
          "expression": { "kind": "comparison", "operand": "second", "operator": "notEqual", "value": 0 }
        }
      ]
    },
    {
      "id": "square-root",
      "name": "Square Root",
      "symbol": "√",
      "shortcut": "r",
      "arity": 1,
      "operands": [
        { "id": "first", "label": "Number", "placeholder": "0" }
      ],
      "validations": [
        {
          "message": "The number must be zero or greater.",
          "expression": { "kind": "comparison", "operand": "first", "operator": "greaterThanOrEqual", "value": 0 }
        }
      ]
    },
    {
      "id": "percentage",
      "name": "Percentage",
      "symbol": "%",
      "shortcut": "%",
      "arity": 2,
      "operands": [
        { "id": "first", "label": "Base value", "placeholder": "0" },
        { "id": "second", "label": "Percentage", "placeholder": "0", "suffix": "%" }
      ],
      "validations": [
        {
          "message": "Percentage must be between 0 and 100.",
          "expression": {
            "kind": "allOf",
            "expressions": [
              { "kind": "comparison", "operand": "second", "operator": "greaterThanOrEqual", "value": 0 },
              { "kind": "comparison", "operand": "second", "operator": "lessThanOrEqual", "value": 100 }
            ]
          }
        }
      ]
    }
  ]
}
```

### curl example

```sh
curl http://localhost:8080/api/v1/operations
```

---

## POST /api/v1/calculations

Validates and executes a calculation. The backend owns expression formatting.

### Request body

```json
{
  "operationId": "division",
  "operands": [144, 12]
}
```

| Field | Type | Description |
|---|---|---|
| `operationId` | string | Stable operation ID from the manifest |
| `operands` | number[] | Operand values in declaration order |

### Response — 200 OK

```json
{
  "expression": "144 ÷ 12",
  "result": 12
}
```

Unary example (`√9`):

```json
{
  "expression": "√9",
  "result": 3
}
```

Percentage example (`15% of 200`):

```json
{
  "expression": "15% of 200",
  "result": 30
}
```

### Error responses

All errors follow the shape:

```json
{ "code": "...", "message": "..." }
```

| HTTP status | `code` | Condition |
|---|---|---|
| 400 | `malformed_request` | Request body is not valid JSON |
| 404 | `unknown_operation` | `operationId` is not registered |
| 422 | `validation_failed` | An operation constraint was violated |
| 500 | `calculation_failed` | Internal evaluation or execution error |

### Validation failure example

```sh
curl -X POST http://localhost:8080/api/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operationId":"division","operands":[10,0]}'
```

```json
{
  "code": "validation_failed",
  "message": "The divisor must not be zero."
}
```

### curl examples

```sh
# Addition
curl -X POST http://localhost:8080/api/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operationId":"addition","operands":[3,4]}'

# Square root
curl -X POST http://localhost:8080/api/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operationId":"square-root","operands":[9]}'

# Percentage (returns "15% of 200")
curl -X POST http://localhost:8080/api/v1/calculations \
  -H 'Content-Type: application/json' \
  -d '{"operationId":"percentage","operands":[200,15]}'
```

---

## Validation expression schema

The `expression` field in each validation uses a discriminated union by `kind`.

### `comparison`

```json
{
  "kind": "comparison",
  "operand": "first" | "second",
  "operator": "equal" | "notEqual" | "greaterThan" | "greaterThanOrEqual" | "lessThan" | "lessThanOrEqual",
  "value": <number>
}
```

### `allOf`

```json
{
  "kind": "allOf",
  "expressions": [ <expression>, ... ]
}
```

An empty `expressions` array is vacuously true. Children are evaluated in declaration order; evaluation short-circuits on the first false child.
