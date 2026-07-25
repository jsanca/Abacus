# 0001 — Defer Production Telemetry

| Field | Value |
|---|---|
| Status | Accepted |
| Date | 2026-07-24 |
| Author | Clio |

## Context

Abacus is a self-contained calculator exercise. Adding a real metrics or tracing backend (e.g., OpenTelemetry, Prometheus) would introduce external dependencies, additional configuration surface, and operational complexity disproportionate to the project scope.

## Decision

Introduce a concrete `observability.Observer` type with no-op lifecycle methods as an explicit telemetry extension boundary. The boundary exists in the codebase and is wired into the runtime; the no-op implementation is used for the duration of development.

## Consequences

- No external telemetry dependency is introduced.
- Telemetry events (server start, server stop, request handled) are identified and documented but not forwarded anywhere in the current implementation.
- A production observer can be wired in by replacing the `*observability.Observer` at the composition root without restructuring the codebase.

## Alternatives Considered

- **Omit the boundary entirely:** Simpler, but makes telemetry a structural change rather than a configuration change when production requirements emerge.
- **Introduce OpenTelemetry immediately:** Proportionate for a long-lived production service but excessive for a scoped exercise with no telemetry consumer.
