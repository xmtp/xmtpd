# Concurrency Safety

## Context

- Always `defer cancel()` when using `context.WithCancel` or `context.WithTimeout`
- Propagate context through call chains -- don't store in structs

## Goroutines

- Never leak goroutines -- ensure they exit cleanly
- Every goroutine needs a clear shutdown signal (context, done channel, or WaitGroup)
- Avoid data races -- use mutexes or channels, not both for same data
- Don't rely on goroutine scheduling order

## Channels

- Never close a channel from the receiving side
- Understand blocking semantics before choosing buffered vs unbuffered
- Use `select` with `context.Done()` to avoid goroutine leaks in channel operations

## Mutexes

- Keep critical sections small
- Prefer `sync.RWMutex` when reads vastly outnumber writes
- Always unlock in defer: `mu.Lock(); defer mu.Unlock()`
