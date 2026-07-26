# Architecture Knowledge

Abacus keeps the calculator catalogue and rules in the Go domain. The API projects that source of truth as a manifest; the React client renders it without embedding a duplicate list of operations or validation rules. Consequential trade-offs are recorded in the [engineering ADRs](../../engineering/adr/README.md).

## Component view

```mermaid
flowchart LR
  Browser[Browser] --> React[React client]
  React --> Nginx[Nginx]
  Nginx --> API[Go API]
  API --> Registry[Immutable operation registry]
```

## Calculation sequence

```mermaid
sequenceDiagram
  actor User
  participant React as React client
  participant API as Go API
  participant Registry as Operation registry

  User->>React: Select operation / enter values
  React->>API: GET /operations (on load)
  API->>Registry: Project manifest
  Registry-->>API: Manifest
  API-->>React: Manifest
  User->>React: Calculate
  React->>React: Projected validation
  React->>API: POST /calculations
  API->>Registry: Lookup, validate, execute, format
  Registry-->>API: Expression and result
  API-->>React: Calculation response
  React-->>User: Result
```
