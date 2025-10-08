# Go-ABI Usage Guide

## Overview

go-abi is a Go implementation of Ethereum ABI encoding/decoding that uses code generation to avoid runtime reflection, providing better performance and type safety.

## Quick Start

### Basic Usage

```go
package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/huangyi/go-abi/pkg/abi"
)

func main() {
	// Create an ABI instance
	abiInstance := abi.ExampleUint256()

	// Encode a function call
	to := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	value := big.NewInt(1000000)

	encoded, err := abiInstance.EncodeFunctionCall("transfer", to, value)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Encoded: %s\n", hex.EncodeToString(encoded))
}
```

### Individual Type Encoding

```go
// Encode a uint256
uint256Type := &abi.Uint256{}
encoded, err := uint256Type.Encode(big.NewInt(42))

// Encode a boolean
boolType := &abi.Bool{}
encoded, err := boolType.Encode(true)

// Encode an address
addrType := &abi.Address{}
addr := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
encoded, err := addrType.Encode(addr)
```

### Decoding

```go
// Decode a uint256
var decoded big.Int
bytesRead, err := uint256Type.Decode(encodedData, &decoded)

// Decode a boolean
var result bool
bytesRead, err := boolType.Decode(encodedData, &result)
```

## Code Generation

### Using the ABI Generator

1. Create an ABI JSON file for your contract:

```json
[
  {
    "type": "function",
    "name": "transfer",
    "inputs": [
      {"name": "to", "type": "address"},
      {"name": "value", "type": "uint256"}
    ],
    "outputs": [
      {"name": "success", "type": "bool"}
    ]
  }
]
```

2. Generate Go bindings:

```bash
go run cmd/generator/main.go -abi contract.abi -out bindings.go -pkg main
```

3. Use the generated code:

```go
package main

import "math/big"

func main() {
	// Use generated types
	transfer := Transfer{
		To: [20]byte{...},
		Value: big.NewInt(1000),
	}

	// Encode using generated methods
	encoded, err := transfer.Encode()
	if err != nil {
		panic(err)
	}

	// Send to Ethereum network...
}
```

### Source Code Generation

You can also generate bindings from Go source code with ABI annotations:

```go
// abi:generate
// Transfer represents a token transfer
type Transfer struct {
	From  [20]byte
	To    [20]byte
	Value *big.Int
}
```

Run the generator:

```go
generator := abi.NewGenerator()
generated, err := generator.GenerateFromSource(sourceCode, "main")
```

## Supported Types

- **Unsigned Integers**: `uint8`, `uint16`, `uint32`, `uint64`, `uint128`, `uint256`
- **Signed Integers**: `int8`, `int16`, `int32`, `int64`, `int128`, `int256`
- **Boolean**: `bool`
- **Address**: `address` (20 bytes)
- **Bytes**: `bytes` (dynamic)
- **String**: `string` (dynamic)
- **Arrays**: Fixed and dynamic arrays
- **Tuples**: Struct types (via code generation)

## Performance Benefits

- **No Reflection**: All type information is known at compile time
- **Type Safety**: Compile-time validation of ABI types
- **Efficient**: Direct memory operations without runtime type checks
- **Predictable**: Consistent performance characteristics

## Testing

Run the test suite:

```bash
go test ./pkg/abi/...
```

## Examples

See the `examples/` directory for complete working examples:

- `examples/simple/` - Basic encoding/decoding examples
- `examples/generator/` - Code generation examples

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request