# osk-go-guide — Go Engineering Skill

Derived from and inspired by the [Uber Go Style Guide](https://github.com/uber-go/guide).

## Purpose

This skill guides Go implementation and code review for the Abacus project and future Go services. It helps agents:

- Write idiomatic, maintainable, production-quality Go.
- Review Go code for correctness, clarity, and idiomaticity.
- Identify Java-style patterns accidentally transferred into Go.
- Keep architecture proportional to the problem.
- Produce actionable, ranked findings rather than stylistic noise.

## Engineering Philosophy

The OSK approach is practical engineering judgment, not rule recital:

- Keep the happy path obvious.
- Prefer explicit ownership for resources and lifecycle.
- Favor small, composable abstractions.
- Prefer immutable configuration after startup.
- Keep architecture proportional to the problem size.
- Optimize for maintainability before cleverness.

This philosophy complements the curated Uber guidance below. It helps agents apply the guidance to the project rather than treating it as a mechanical checklist.

## Activation Conditions

Apply this skill when:

- Implementing Go services, packages, APIs, CLIs, or runtimes.
- Reviewing Go architecture or code.
- Reviewing goroutine lifecycles or concurrency safety.
- Remediating lint findings.
- Reviewing tests or test patterns.
- Evaluating interfaces, registries, factories, or package structure.

## Core Principles

1. **Clarity over cleverness.** Write code a new reader can understand in one pass.
2. **Explicit dependency construction.** Construct and wire dependencies at startup; avoid package-level mutable state.
3. **Composition over hierarchy.** Prefer small, focused types composed together over deep inheritance-like type trees.
4. **Functions as first-class strategies.** Use function types (`type ExecuteFunc func(...) (...)`) before introducing interface hierarchies for stateless behavior selection.
5. **Interfaces defined by consumers.** Declare interfaces in the package that uses them, not the package that implements them.
6. **No speculative interfaces.** Do not create an interface until at least two concrete implementations exist or a consumer has a demonstrated need.
7. **Immutable configuration where possible.** Once initialized, configuration and registries must not change.
8. **Package-level mutable state avoided.** Use dependency injection; do not expose mutable global variables or function pointers.
9. **Lifecycle ownership must be explicit.** Every goroutine, connection, and resource must have a clear owner responsible for its shutdown.
10. **Goroutines must have a defined termination path.** Every goroutine must be stoppable and waitable.
11. **Errors must preserve context.** Wrap errors with `%w` so callers can inspect the chain. Never discard error context without reason.
12. **Names must describe intent.** Avoid abbreviations, one-letter names (except trivial loop indices), and compressed identifiers.
13. **Exported symbols require Go documentation comments.** Every exported function, type, constant, and variable must have a doc comment.
14. **Abstractions must remain proportional to the project.** Do not build for imagined future scale.

## Decision Heuristics

- **Introduce an interface** when a consumer needs to depend on a stable capability across multiple real implementations, or when a boundary genuinely needs substitution. Define it at the consuming boundary.
- **Prefer a concrete type** when there is one implementation and direct use makes ownership and behavior clearer. Do not add an interface solely for a mock.
- **Use composition** when independently useful behavior can be assembled from focused types. Do not turn simple delegation into a hierarchy.
- **Use a function type** for stateless, interchangeable behavior such as an arithmetic operation. Use an interface only when the behavior needs multiple cohesive methods or stateful lifecycle.
- **Choose immutability** for configuration, registries, and metadata built at startup and read thereafter. It reduces synchronization and makes ownership clear.
- **Add abstraction** only when it removes present complexity, protects a real boundary, or enables a demonstrated product requirement. Explain the benefit before adding the layer.

## When to Deviate

Guidance serves the product and maintainers; it is not a substitute for judgment. Deviate deliberately when API consistency outweighs an internal preference, a simpler design is better than theoretical extensibility, or a product requirement justifies an exception. Keep the deviation small, document the reason when it is durable, and do not use an exception to conceal avoidable complexity.

## Ownership Model

The creator of a resource must make its ownership explicit: identify who owns it, who may use it, and who shuts it down. The runtime owns HTTP server startup and graceful shutdown; a component that starts goroutines owns their cancellation and wait; startup owns registry construction before it becomes immutable; and the composition root owns telemetry lifecycle. Clear ownership prevents leaks, concurrent mutation, and shutdown ambiguity.

## BAD / GOOD Examples

### Avoid speculative interfaces

```go
// BAD: one implementation and no consumer boundary require this interface.
type Registry interface { Find(string) (Operation, bool) }

// GOOD: depend on the concrete registry until a real boundary requires otherwise.
type Registry struct { operations map[string]Operation }
```

### Make registry ownership explicit

```go
// BAD: mutable package-level state with unclear initialization and ownership.
var operations = map[string]Operation{}

// GOOD: construct at startup, then expose read-only behavior.
registry := NewRegistry(operationDefinitions)
```

### Use functions for stateless operations

```go
// BAD: hierarchy for a single stateless calculation.
type AdditionStrategy interface { Execute([]float64) (float64, error) }

// GOOD: a function describes the behavior directly.
type ExecuteFunc func(operands []float64) (float64, error)
```

### Prefer descriptive naming

```go
// BAD
func (svc *Service) Do(ctx context.Context, p Payload) error

// GOOD
func (service *CalculatorService) ExecuteOperation(ctx context.Context, request OperationRequest) error
```

## Review Checklist

### Naming

- [ ] Package names are lower-case, singular, and not `util`, `common`, `helpers`, or `lib`.
- [ ] Variable names describe intent; no `p`, `svc`, `mgr`, `ctx2`, or similar compressed identifiers.
- [ ] Acronyms and initialisms follow Go convention (`HTTPServer`, `parseURL`, not `HttpServer`, `parseUrl`).
- [ ] Error variables use `Err` (exported) or `err` (unexported) prefix; error types use `Error` suffix.
- [ ] Function names use MixedCaps; test function names may use underscores for grouping (`TestFoo_WhatIsBeingTested`).
- [ ] No shadowing of built-in names (`error`, `string`, `len`, etc.).

### Packages and Layout

- [ ] Packages group related functionality; no catch-all `/pkg`, `/util`, `/common`, or `/helpers` directories.
- [ ] Imports are organized: stdlib group first, then external, separated by a blank line.
- [ ] Import aliases used only when the last element of the import path does not match the package name.
- [ ] No `init()` functions except for deterministic, order-independent, I/O-free initialization (prefer explicit calls from `main()`).

### Functions and Methods

- [ ] Functions are short and do one thing.
- [ ] Early returns reduce nesting; no unnecessary `else` branches.
- [ ] Parameters use descriptive types rather than bare booleans where clarity demands it.
- [ ] `os.Exit` / `log.Fatal` appear only in `main()`; all other functions return errors.
- [ ] At most one `os.Exit` call in `main()`; prefer a `run() error` function.
- [ ] Constructors use `&T{Field: value}` not `new(T)`.
- [ ] Struct initialization uses named fields.
- [ ] Zero-valued structs declared with `var v T` not `v := T{}`.

### Structs and Ownership

- [ ] No embedded mutexes; mutexes use named fields (`mu sync.Mutex`).
- [ ] Structs containing `sync.Mutex`, `sync.WaitGroup`, or similar synchronization primitives are not copied (pass by pointer).
- [ ] Zero-value mutexes used directly (`var mu sync.Mutex`), not `new(sync.Mutex)`.
- [ ] Slices and maps defensive-copied when received from or returned to external callers.
- [ ] Embedded types in public structs avoided unless the inner type's API genuinely belongs as part of the outer type's public API.

### Interfaces

- [ ] Interfaces defined in the consuming package, not the implementing package.
- [ ] No interfaces with only one concrete implementation (unless the consumer genuinely needs the abstraction).
- [ ] Interfaces not created solely for testing; prefer concrete types or function injection.
- [ ] Compile-time compliance checks used for exported types: `var _ SomeInterface = (*MyType)(nil)`.
- [ ] Never use pointer-to-interface; interfaces are passed by value.

### Errors

- [ ] Error sentinel variables use `errors.New` for static messages.
- [ ] Dynamic error messages use `fmt.Errorf` with `%w` to wrap underlying errors.
- [ ] Custom error types used only when callers need `errors.As` matching with dynamic context.
- [ ] Errors handled exactly once: wrap-and-return up the call stack, or log-and-degrade; do not both log and return the same error.
- [ ] Context added to errors without repeating `"failed to"` chains.
- [ ] No panics in production code (except `template.Must`-style init-time patterns).

### Context

- [ ] `context.Context` is always the first parameter of functions.
- [ ] Context is never stored in struct fields; it is passed explicitly through the call chain.
- [ ] Context-derived values are propagated to all downstream calls that accept context.

### Concurrency and Lifecycle

- [ ] Every goroutine has a defined stop mechanism (cancellation via context, `done` channel, or `sync.WaitGroup`).
- [ ] The code that starts a goroutine is responsible for waiting for it to exit (or the code that owns the component).
- [ ] No goroutines spawned in `init()` functions.
- [ ] `sync.WaitGroup` used for multiple goroutines; `done chan struct{}` for a single goroutine.
- [ ] No writes to shared data without synchronization; no reads of data that may be concurrently written.
- [ ] Concurrent reads of immutable data (e.g., a registry populated at startup) require no mutex.
- [ ] Graceful shutdown is implemented: signal handling, drain/wait, then exit.

### Collections

- [ ] Empty slices returned as `nil`, not `[]T{}`.
- [ ] Slice emptiness checked with `len(s) == 0`, never `s == nil`.
- [ ] Maps initialized with `make()` for empty maps; literals for fixed content.
- [ ] Capacity hints provided to `make()` when the expected size is known.
- [ ] Maps with value-type values: code accounts for non-addressable values (pointer-receiver methods won't work on `map[int]S` — store `*S` instead).

### Logging and Telemetry

- [ ] Structured logging via `log/slog`; no `fmt.Println` for application output.
- [ ] Log levels used appropriately: debug for detailed flow, info for key events, warn for recoverable issues, error for failures.
- [ ] Telemetry abstraction respected; no-op implementations allowed but must not become ornamental architecture.

### HTTP Services

- [ ] Handlers accept `http.ResponseWriter` and `*http.Request`; no custom handler signatures without clear justification.
- [ ] Request-scoped context propagated from `r.Context()`.
- [ ] Response encoding/decoding uses standard `encoding/json` or a thin wrapper; no heavy serialization frameworks without justification.
- [ ] Tests use `net/http/httptest` for full-stack HTTP testing.

### Tests

- [ ] Table-driven tests used for multiple inputs on the same logic: `[]struct{ give, want }` pattern.
- [ ] Sub-tests use `t.Run(tt.name, func(t *testing.T) { ... })`.
- [ ] `t.Parallel()` used with loop variable capture (`tt := tt`).
- [ ] Test tables avoid conditional mock setup, branching logic, or functions in table rows; split into separate test functions when tests diverge.
- [ ] Test helpers use `t.Helper()`.
- [ ] Test files use `_test.go` suffix and live in the same package (`package foo_test` for black-box; `package foo` for white-box).

### Linting and Formatting

- [ ] Code passes `goimports` / `gofmt` formatting.
- [ ] `go vet` reports no issues.
- [ ] `golangci-lint` configured and run in CI (at minimum: `errcheck`, `govet`, `staticcheck`, `goimports`).
- [ ] No commented-out code in committed files.

### Documentation

- [ ] Every exported symbol has a Go doc comment starting with the symbol name.
- [ ] Package doc comments exist (one `// Package foo ...` comment in any file).
- [ ] Comments explain *why*, not *what* — the code already says what it does.

## Severity Model

| Severity | Meaning |
|----------|---------|
| **BLOCKER** | Correctness error, data race, goroutine leak, lifecycle leak, unsafe behavior, or missing error handling that could cause silent failure. Must be fixed before merge. |
| **MAJOR** | Non-idiomatic design or significant maintainability risk (e.g., global mutable state, speculative interface, interface in wrong package, context stored in struct). Should be fixed before merge unless explicitly justified. |
| **MINOR** | Localized clarity, naming, or consistency issue. Fix at author's discretion; do not block merge over MINOR findings alone. |
| **NOTE** | Optional improvement, educational observation, or recommendation that does not indicate a problem. No action required. |

Multiple findings at the same severity may collectively warrant escalation (e.g., five MINOR naming issues in one file may indicate a MAJOR consistency problem).

## Expected Review Output

Each finding uses this format:

```
Finding:  <one-line description of the issue>
Severity: BLOCKER | MAJOR | MINOR | NOTE
Location: <file>:<line> — <symbol or code snippet>
Why it matters: <impact on correctness, maintainability, or clarity>
Recommended change: <the smallest reasonable fix>
Reference: <link to relevant section in this skill or external source>
```

Review findings must:
- Identify concrete code, not general impressions.
- Explain impact, not just state a rule.
- Propose the smallest reasonable correction.
- Avoid demanding abstractions without demonstrated need.
- Distinguish correctness issues (BLOCKER) from preferences (MINOR/NOTE).

## Abacus-Specific Guidance

The following applies specifically to the Abacus project:

- **Registry:** Use a concrete `Registry` instance, not a package-global mutable map. The registry must be immutable after startup.
- **Concurrent reads:** When the registry is populated once at startup and never written again, concurrent reads from multiple goroutines require no mutex.
- **Operation execution:** Prefer function-based operation strategies (`type ExecuteFunc func(operands []float64) (float64, error)`) over class-like hierarchies with `AdditionStrategy` / `DivisionStrategy` interfaces for stateless arithmetic.
- **Registry interface:** Do not create an interface for the registry unless a consumer has a demonstrated need. The concrete type suffices until an actual consumer requires abstraction.
- **Graceful shutdown:** The runtime must handle OS signals, stop accepting new requests, drain in-flight requests, and then exit. This is a BLOCKER requirement.
- **Logging:** `log/slog` for structured application logging. No `fmt.Println` for application output.
- **Naming:** Descriptive names throughout. No one-letter or compressed identifiers except trivial loop indices.
- **Documentation:** All exported Go symbols must have Go documentation comments.
- **HTTP testing:** Use `net/http/httptest` for API-level tests.
- **Linting:** `golangci-lint` integrated into `make lint` and `make verify`.
- **Telemetry:** The telemetry boundary may use a no-op implementation during early development but must not grow into ornamental architecture (empty interfaces, unused hooks, parameterized-but-unused backends).

## Explicit Anti-Patterns

The following patterns should be flagged on sight:

- Java-style `AdditionStrategy`, `DivisionStrategy`, and deep interface hierarchies for stateless arithmetic functions.
- Global mutable registries (`var ops = map[string]Operation{}` at package level).
- Interfaces created only to enable mocking in tests.
- `context.Context` stored in struct fields.
- Goroutines started without cancellation, ownership, or a defined stop path.
- `fmt.Println` as application logging.
- Meaningless names: `p`, `svc`, `mgr`, `ctx2`, `data`, `info`, `tmp`, `ret`.
- Getters and setters copied mechanically from Java (`GetX()` / `SetX()` for simple struct fields).
- Excessive `/pkg`, `/util`, `/common`, or `/helpers` packages.
- Wrapping every small behavior in a constructor or factory.
- Mutexes guarding immutable data.
- Comments that merely restate the code (`// Add adds two numbers` above `func Add(a, b int) int`).
- Panics in production code paths.
- Unparameterized format strings used as logging templates.
- Empty interface (`interface{}` / `any`) used where a concrete type would suffice.
- `go test` without `-race` in CI.

## Source Attribution

This skill is derived from and inspired by the [Uber Go Style Guide](https://github.com/uber-go/guide) (Apache 2.0 licensed). The guidance has been curated, organized, and adapted for practical implementation and review use. The skill does not reproduce the source guide verbatim. Detailed notes on selection, omission, and adaptation are in `references/uber-go-guide-notes.md`.

The `golang-standards/project-layout` repository is not recognized as an official Go project layout standard and is not referenced by this skill.
