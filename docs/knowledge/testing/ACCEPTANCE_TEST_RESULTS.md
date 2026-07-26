# Playwright acceptance test results

## Latest verified run

- Date: 2026-07-25
- Command: `make e2e-local`
- Environment: Docker Compose production client (Nginx) and Go API, Chromium
- Result: **40 passed** — 20 desktop Chromium tests and 20 mobile Chromium tests in 22.7 seconds.

The run covered the full operation manifest, all real calculations, validation and API rejection paths, result continuity/focus, keyboard and reset behavior, failure recovery, delayed submission, responsive accessibility checks, and curated screenshots.

## Reproducible command

Use `make e2e` for the clean, stack-owning canonical execution. It returns the Playwright result and always removes the Compose resources it started.

## Evidence

![Desktop calculator](evidence/desktop-calculator.png)

![Mobile calculator](evidence/mobile-calculator.png)
