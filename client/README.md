# Client

The Abacus client is a responsive React + TypeScript calculator built with Vite. It renders the live Go API operation manifest, applies its projected validation generically, and calls the proxied `/api/v1` endpoints through `calculatorApi.ts`.

## Commands

From this directory:

- `npm run dev` — start the Vite development server.
- `npm run lint` — run ESLint.
- `npm run test` — run Vitest and React Testing Library tests.
- `npm run coverage` — produce test coverage.
- `npm run build` — type-check and produce a production build.
- `npm run verify` — run all meaningful checks.
- `npm run test:e2e` — run Playwright against an already-running stack.
- `npm run test:e2e:headed` — run the same browser tests with a visible browser.

For a clean, stack-owning acceptance run, use `make e2e` from the repository root. Start the Docker stack first (`docker compose up -d`) before running `npm run test:e2e` directly, or use `make e2e-local` from the repository root.
