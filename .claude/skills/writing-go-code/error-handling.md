# Error Handling Patterns

## Core Rules

- Always check errors -- never ignore them
- Wrap errors with context: `fmt.Errorf("doing X: %w", err)`
- Use `errors.Is` and `errors.As` for comparison (not `==`)
- Avoid panics except in truly unrecoverable cases

## Wrapping Pattern

```go
if err := doSomething(); err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

## Sentinel Errors

Define at package level, check with `errors.Is`:

```go
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) { ... }
```

## Custom Error Types

Use `errors.As` to extract typed errors:

```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println(pathErr.Path)
}
```
