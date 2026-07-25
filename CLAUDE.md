# Abacus Project Instructions

## Naming

- Prefer descriptive names.
- Avoid abbreviations.
- Avoid one-letter variables except for trivial loop indices.

## Documentation

- Exported Go symbols must include Go documentation comments.
- Public React components should include concise documentation when appropriate.

## Logging

- Use structured logging with `log/slog`.
- Do not use `fmt.Println` for application logging.

## Telemetry

- Respect the observability abstraction.
- Use the no-op implementation unless otherwise specified.

## Errors

- Prefer explicit domain errors.
- Avoid anonymous string comparisons.

## Architecture

- Keep abstractions proportional.
- Favor composition.
- Avoid speculative interfaces.
