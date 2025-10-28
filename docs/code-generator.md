# Code Generator Design

ABI spec distinguish static and dynamic types. Static types are encoded in-place and dynamic types are encoded at a separately allocated location after the current block.

The size of static buffer is known at code generation time based on the types, while the dynamic part depends on the runtime value.

## Panic

The generated code is free of panic by design:

1. For an encode call, we pre-calculate the total size based on the value, and allocate the whole buffer at once, the intermdiiate encoding calls don't need to worry about bound checks.

2. The decoding logic must do neccerary bound checks to avoid panic.

## Generation

## Composible Generation

We derive a deterministic identifier for each ABI type, which is used to name unique functions.

### Size Functions

Size functions are only generated for dynamic types, because static type's size is known at generation time.

```golang
func size_string(v string) int {
  return 32 + pad32(len(v))
}

func size_array3_string(s [3]string) int {
  size := 32 * 3
  size += size_string(s[0])
  size += size_string(s[1])
  size += size_string(s[2])
  return size
}

func size_slice_string(s []string) int {
  size := 32 + 32 * len(s)
  for _, item := range s {
    size += size_string(item)
  }
  return size
}

func size_dynamic_tuple(a, b, c) int {
  return 32 * 3 + size_$t(a) + size_$t(b) + size_$t(c)
}

func size_static_tuple(a, b, c) int{
  return size_$t(a), size_$t(b), size_$t(c)
}
```

### Encode Functions

### Decode Function