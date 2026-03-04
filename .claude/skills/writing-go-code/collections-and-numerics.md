# Collections and Numeric Types

## Slices

- Distinguish between nil and empty slices (`var s []int` vs `s := []int{}`)
- Preallocate capacity when size is known: `make([]T, 0, n)`
- Avoid memory leaks from slicing large arrays (copy instead)
- Check capacity when copying or appending

## Maps

- Always initialize before use: `m := make(map[K]V)`
- Check existence with two-value: `val, ok := m[key]`
- Iteration order is random -- never depend on it
- Don't compare maps with `==` -- use `maps.Equal` or `reflect.DeepEqual`

## Numeric Types

- Avoid default `int` unless the size is appropriate for the platform
- Watch for integer overflow and underflow
- Avoid float precision issues -- use `math/big` for exact arithmetic when needed
- Be explicit about type conversions
