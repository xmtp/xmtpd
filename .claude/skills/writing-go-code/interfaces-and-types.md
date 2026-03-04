# Interfaces and Type Design

## Interfaces

- Define interfaces on the **consumer** side, not the producer
- Don't define interfaces until you need them
- Don't return interfaces from constructors or public APIs
- Keep interfaces small and focused: 1-2 methods
- Prefer `any` over `interface{}`, but avoid either unless necessary

## Structs and Methods

- Avoid embedding pointer types unless necessary
- Don't overuse getters/setters -- prefer public fields when it makes sense
- **Value receivers** when method doesn't mutate state or require pointer semantics
- **Pointer receivers** when method mutates state or struct is large

## Generics (Go 1.18+)

- Only use when they simplify code or add real flexibility
- Avoid over-engineering with type parameters
- Keep constraint complexity low and readable
