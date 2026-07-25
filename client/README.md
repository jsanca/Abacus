# Client

The Abacus client is a responsive React + TypeScript calculator built with Vite. It renders an operation manifest, applies its projected validation generically, and calls a small mocked API boundary until the Go backend is available.

## Commands

From this directory:

- `npm run dev` — start the Vite development server.
- `npm run lint` — run ESLint.
- `npm run test` — run Vitest and React Testing Library tests.
- `npm run coverage` — produce test coverage.
- `npm run build` — type-check and produce a production build.
- `npm run verify` — run all meaningful checks.

The mock API is intentionally internal to the feature and models the future operation-discovery and calculation endpoints. It is not a substitute for the backend contract.
