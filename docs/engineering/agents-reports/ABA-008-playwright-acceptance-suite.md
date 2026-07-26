# A-008 — Playwright acceptance suite report

## Outcome

Implemented a reproducible, real-stack Playwright acceptance suite for the Abacus calculator. The suite uses Chromium desktop and mobile projects and reaches the live Go API through the production React build and Nginx proxy.

## Delivered

- `@playwright/test`, Chromium project configuration, E2E scripts, and ignored transient reports/artifacts.
- Reusable calculator page object and independent acceptance, interaction, resilience, accessibility, and evidence tests.
- `make e2e` Compose lifecycle ownership with clean start, build, readiness checks, guaranteed teardown, and preserved Playwright exit status.
- `make e2e-local` for an existing stack.
- Explicit result landmark semantics to make the live result announcement testable and accessible.
- Test strategy, verified results, and curated desktop/mobile screenshots.

## Verification

`npm run lint && npm run build` passed in `client/`.

`make e2e-local` passed with 40 tests: 20 desktop Chromium and 20 mobile Chromium (26.8 seconds). The canonical `make e2e` script owns Docker lifecycle and performs the same suite after readiness checks.

## Notes

No external services or test-only API endpoints were introduced. Network delay and service-failure scenarios use Playwright request routing inside the browser context; ordinary operation and validation paths use the real backend.
