---
name: writing-go-code
description: >-
  Use when writing, modifying, or reviewing .go files, implementing new Go
  functions or packages, or when Go code style and conventions are relevant.
---

# Writing Go Code

## Style Priorities (in order)

1. **Clarity** -- purpose and rationale clear to the reader
2. **Simplicity** -- simplest way to accomplish the goal
3. **Concision** -- high signal-to-noise ratio
4. **Maintainability** -- easy to modify correctly
5. **Consistency** -- consistent with surrounding codebase

## Core Rules

- Follow [Google Go Style Guide](https://google.github.io/styleguide/go/guide)
- Format all code with `golangci-lint fmt`
- MixedCaps/mixedCaps only -- never snake_case (even constants: `MaxLength` not `MAX_LENGTH`)
- No fixed line length -- refactor long lines instead of splitting
- Shorter names in Go than other languages; context reduces need for verbosity
- Comments explain **why**, not what
- Allow code to speak for itself with self-describing symbol names rather than redundant comments
- Use least-powerful mechanism: language primitive > stdlib > external dependency

## Formatting

- Imports grouped: stdlib, external, internal (blank line between groups)
- Avoid magic numbers -- use named constants
- No unnecessary levels of abstraction

## Naming

- Exported: `PascalCase`
- Unexported: `camelCase`
- Short receiver names (1-2 chars matching type initial)
- Acronyms keep case: `HTTPClient`, `xmlParser`
- Package names: short, lowercase, no underscores, no `util`/`common`/`base`
- Names should not feel repetitive when used: `queue.New()` not `queue.NewQueue()`
- Predictable names -- a user should be able to predict the name in a given context

## Function Design

- Keep functions small and focused
- Prefer composition over embedding
- Use functional options pattern for flexible constructors
- Return concrete types, accept interfaces
- Avoid `init()` unless absolutely necessary
- Avoid variable shadowing
- Do not over-nest control flow (flatten with early returns)

## Testing

- Table-driven tests where appropriate
- Name tests consistently: `TestXxx`, `BenchmarkXxx`
- Use `t.Helper()` in helper functions
- Use `t.Context()` rather than `context.Background()` or `context.TODO()`
- Avoid global state in tests
- Tests should provide clear, actionable diagnostics on failure

## Common Pitfalls

See these references for detailed guidance on frequent Go mistakes:

- **[Interfaces and type design](interfaces-and-types.md)** -- when designing interfaces, choosing receivers, or structuring types
- **[Error handling patterns](error-handling.md)** -- when handling, wrapping, or comparing errors
- **[Concurrency safety](concurrency-safety.md)** -- when writing goroutines, channels, or shared-state code
- **[Collections and numeric types](collections-and-numerics.md)** -- when working with slices, maps, or numeric conversions
