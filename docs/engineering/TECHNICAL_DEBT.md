# Technical Debt Register

| ID | Description | Decision | Status |
| --- | --- | --- | --- |
| TD-001 | Production telemetry deferred. The `observability.Observer` no-op boundary is in place; a real metrics or tracing backend is not wired in. | Accepted for project scope; see ADR 0001. | Open |
| TD-002 | Manifest `validations` field is always `[]`. Client-side validation (divisor ≠ 0, radicand ≥ 0, percentage 0–100) is disabled until ABA-006 delivers the declarative validation engine. | Accepted; ABA-006 will populate the field. | Open |
