# Playwright acceptance test strategy

## Purpose

The acceptance suite verifies the deployed local path: Chromium browser → production React bundle → Nginx `/api` proxy → Go API → operation registry, validation, execution, and response formatting. It does not replace focused unit tests; it confirms that the independently tested parts integrate correctly.

## Execution modes

- `make e2e` is the canonical reproducible command. It removes prior Compose resources, builds and starts the stack, probes the backend health endpoint and browser entry point, runs Playwright, preserves its exit status, and tears the stack down.
- `make e2e-local` is the fast feedback command for an already-running local stack.

Both modes run the same test files under Chromium desktop and iPhone-sized mobile projects. The suite is parallel-safe: each test creates its own browser context and only routes its own requests.

## Coverage

- Real manifest loading and all seven registry operations: addition, subtraction, multiplication, division, exponentiation, square root, and percentage.
- Projected client validation, authoritative API validation, boundary values, and the no-request division-by-zero rule.
- Result continuation and focus for binary-to-binary, binary-to-unary, and unary-to-binary transitions; Clear, Escape, Enter, and manifest-sourced operator shortcuts.
- Delayed execution/duplicate-submit protection, recoverable calculation failures, and manifest failure/retry without a hardcoded operator fallback.
- Semantic controls, live results and status, percent suffix rendering, responsive overflow checks, and desktop/mobile screenshot evidence.

## Evidence and diagnostics

Curated passing screenshots are retained in [evidence](evidence/). On a failure, Playwright retains a screenshot, video, and trace under `client/test-results/`; its HTML report is written to `client/playwright-report/`. Those volatile diagnostics are intentionally ignored by Git.
