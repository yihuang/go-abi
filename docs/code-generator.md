# Code Generator Design

The ABI specification distinguishes between static and dynamic types. Static types are encoded in-place and dynamic types are encoded at a separately allocated location after the current block.

The size of the static buffer is known at code generation time based on the types, while the dynamic part depends on the runtime value.

## Panic-Free Design

The generated code is free of panic by design:

1. For an encode call, we pre-calculate the total size based on the value, and allocate the whole buffer at once. The intermediate encoding calls don't need to worry about bound checks.

2. The decoding logic must do necessary bound checks to avoid panic.

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

Encode functions generate ABI-encoded bytes for each type. They follow a pattern of pre-calculating size, allocating buffer, then encoding:

```golang
// Static types (known size at compile time)
func encode_uint256(v *big.Int, buf []byte) (int, error) {
    return abi.EncodeBigInt(v, buf[:32], false), nil
}

func encode_address(v common.Address, buf []byte) (int, error) {
    copy(buf[12:32], v[:])
    return 32, nil
}

// Dynamic types (size calculated at runtime)
func encode_string(s string, buf []byte) (int, error) {
    size := len(s)
    // Encode length (32 bytes)
    binary.BigEndian.PutUint64(buf[24:32], uint64(size))

    // Encode data
    copy(buf[32:], []byte(s))

    return 32 + abi.Pad32(size), nil
}

// Slice encoding (dynamic array)
func encode_slice_uint256(arr []*big.Int, buf []byte) (int, error) {
    offset := 32 + 32 * len(arr) // offset for array data
    binary.BigEndian.PutUint64(buf[24:32], uint64(len(arr)))

    for i, v := range arr {
        // Write offset for element i
        binary.BigEndian.PutUint64(buf[24+32*i:32+32*i], uint64(offset))
        // Encode element
        abi.EncodeBigInt(v, buf[offset:offset+32], false)
        offset += 32
    }

    return offset, nil
}

// Complex types compose simpler types
func encode_tuple_person(p Person, buf []byte) (int, error) {
    // Encode each field sequentially
    offset := 64 // static part size

    // Field 1: name (string)
    copy(buf[0:], buf[32:64]) // offset placeholder
    nameSize := encode_string(p.Name, buf[offset:])
    offset += nameSize

    // Field 2: age (uint256)
    encode_uint256(big.NewInt(int64(p.Age)), buf[32:])

    // Field 3: address (address)
    encode_address(p.Address, buf[64:])

    return offset, nil
}
```

#### EncodeTo Pattern

All generated types support `EncodeTo(buf []byte)` for efficient buffer reuse:

```golang
// Generated code
func (p Person) EncodeTo(buf []byte) (int, error) {
    size := p.EncodedSize()
    if len(buf) < size {
        return 0, fmt.Errorf("buffer too small: got %d, need %d", len(buf), size)
    }

    // Encoding logic
    return encodedSize, nil
}

// Convenience method
func (p Person) Encode() ([]byte, error) {
    buf := make([]byte, p.EncodedSize())
    _, err := p.EncodeTo(buf)
    return buf, err
}
```

### Decode Functions

Decode functions read ABI-encoded bytes back into Go types. They include necessary bounds checking:

```golang
// Static type decoding
func decode_uint256(data []byte) (*big.Int, int, error) {
    if len(data) < 32 {
        return nil, 0, io.ErrUnexpectedEOF
    }

    result := new(big.Int)
    result.SetBytes(data[:32])
    return result, 32, nil
}

func decode_address(data []byte) (common.Address, int, error) {
    if len(data) < 32 {
        return common.Address{}, 0, io.ErrUnexpectedEOF
    }

    var result common.Address
    copy(result[:], data[12:32])
    return result, 32, nil
}

// Dynamic type decoding with bounds checking
func decode_string(data []byte) (string, int, error) {
    if len(data) < 32 {
        return "", 0, io.ErrUnexpectedEOF
    }

    // Read length
    length := int(binary.BigEndian.Uint64(data[24:32]))

    // Calculate total size (32 + padded data)
    totalSize := 32 + abi.Pad32(length)
    if len(data) < totalSize {
        return "", 0, io.ErrUnexpectedEOF
    }

    // Decode data
    result := string(data[32 : 32+length])
    return result, totalSize, nil
}

// Slice decoding (dynamic array)
func decode_slice_uint256(data []byte) ([]*big.Int, int, error) {
    if len(data) < 32 {
        return nil, 0, io.ErrUnexpectedEOF
    }

    // Read array length
    length := int(binary.BigEndian.Uint64(data[24:32]))

    // Check we have enough offsets
    neededSize := 32 + 32*length
    if len(data) < neededSize {
        return nil, 0, io.ErrUnexpectedEOF
    }

    result := make([]*big.Int, length)

    // Decode each element
    for i := 0; i < length; i++ {
        offset := int(binary.BigEndian.Uint64(data[24+32*i:32+32*i]))
        elem, _, err := decode_uint256(data[offset:])
        if err != nil {
            return nil, 0, err
        }
        result[i] = elem
    }

    return result, neededSize, nil
}

// Tuple decoding
func (t *Person) Decode(data []byte) (int, error) {
    if len(data) < 64 {
        return 0, io.ErrUnexpectedEOF
    }

    // Decode static fields
    name, nameSize, err := decode_string(data[0:])
    if err != nil {
        return 0, err
    }

    age, _, err := decode_uint256(data[32:])
    if err != nil {
        return 0, err
    }

    addr, _, err := decode_address(data[64:])
    if err != nil {
        return 0, err
    }

    t.Name = name
    t.Age = int(age.Int64())
    t.Address = addr

    return 64 + nameSize, nil
}
```

#### Decoding with Struct Methods

Generated types include convenient `Decode` methods:

```golang
// Decode directly into struct
func (p *Person) Decode(data []byte) error {
    _, err := p.decodeFromBuffer(data)
    return err
}

// Decode with offset tracking (for embedded decoding)
func (p *Person) decodeFromBuffer(data []byte) (int, error) {
    // ... decoding logic
    return offset, nil
}
```

## Stdlib Functions

To reduce code duplication, common primitive types use stdlib functions instead of generating duplicate code:

- **Primitive types**: address, bool, uint8, uint256, int256, string, bytes, bytes32
- **Slice types**: address[], bool[], uint8[], uint256[], int256[], string[], bytes[], bytes32[]

Generated code calls these stdlib functions:

```golang
// Instead of generating encode_uint256 locally:
func encode_myUint256(v *big.Int, buf []byte) (int, error) {
    return abi.EncodeUint256(v, buf) // Uses stdlib
}

// This reduces generated code size and improves maintainability
```

## Error Handling

All generated functions follow Go best practices:

1. **Return errors, never panic**: Even in debug builds
2. **Bounds checking**: Verify buffer/data sizes before access
3. **Error wrapping**: Use `%w` verb for wrapped errors
4. **Clear error messages**: Describe what went wrong and where

```golang
// Example error handling patterns
if len(data) < requiredSize {
    return 0, fmt.Errorf("data too short: got %d bytes, need %d: %w",
        len(data), requiredSize, io.ErrUnexpectedEOF)
}

if err != nil {
    return 0, fmt.Errorf("failed to encode field %s: %w", fieldName, err)
}
```

## Performance Optimizations

### Single Allocation Pattern

The generator follows a "pre-calculate, allocate once" pattern:

```golang
func (p Person) Encode() ([]byte, error) {
    // Pre-calculate total size
    size := p.EncodedSize()

    // Allocate buffer once
    buf := make([]byte, size)

    // Encode to pre-allocated buffer
    _, err := p.EncodeTo(buf)
    return buf, err
}
```

### Buffer Reuse

Users can pre-allocate buffers for high-throughput scenarios:

```golang
// Allocate once
buf := make([]byte, expectedSize)

// Reuse buffer
for _, person := range people {
    n, err := person.EncodeTo(buf)
    process(buf[:n], err)
}
```

### Static Size Constants

For performance-critical code, static size constants are generated:

```golang
const PersonStaticSize = 96 // bytes

func (p Person) EncodedSize() int {
    return PersonStaticSize + size_string(p.Name)
}
```

## Type-Specific Patterns

### Integers

- **uint8-uint64**: Native Go types for efficiency
- **uint256/int256**: `*big.Int` with helper functions
- **Small integers (<64 bits)**: Native types where possible

### Arrays

- **Fixed arrays**: Statically sized, encoded inline
- **Dynamic arrays**: Length prefix + dynamic data
- **Nested arrays**: Recursive encoding

### Tuples

- **Static tuples**: All fields static, encoded inline
- **Dynamic tuples**: Contains dynamic fields, uses offset references
- **Nested tuples**: Recursive structure support

## Code Generation Pipeline

1. **Parse ABI**: JSON or human-readable format
2. **Collect types**: Gather all unique types needed
3. **Generate functions**: Create encode/decode/size functions
4. **Generate structs**: Create tuple/struct types
5. **Generate methods**: Add Encode/Decode/EncodedSize methods
6. **Generate selectors**: Add function/event selector constants
7. **Format code**: Apply goimports and format code

Each step is deterministic, ensuring reproducible builds.