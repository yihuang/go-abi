# Go ABI Code Generator

Generating statically typed Go structs and encode/decode methods from Ethereum contract ABI definitions.

[![Go Reference](https://pkg.go.dev/badge/github.com/yihuang/go-abi.svg)](https://pkg.go.dev/github.com/yihuang/go-abi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yihuang/go-abi)](https://goreportcard.com/report/github.com/yihuang/go-abi)
[![Tests](https://github.com/yihuang/go-abi/actions/workflows/go.yml/badge.svg)](https://github.com/yihuang/go-abi/actions/workflows/go.yml)

## Why go-abi?

### 🚀 **High Performance**
- **Single Allocation**: Allocate buffer once during encoding for maximum efficiency
- **Zero-Copy Decoding**: Direct memory operations with proper bounds checking
- **Performance Benchmarks**: Competitive or better than go-ethereum

### 🛡️ **Type Safety**
- **Static Typing**: Natural type mapping from ABI to Go
- **Compile-Time Checks**: Catch type errors at compile time
- **No Reflection**: Generated code is fast and type-safe

### 📝 **Developer Experience**
- **Human Readable ABI**: Write ABIs in readable format directly in Go
- **Auto-Generation**: Keep code in sync with `go generate`
- **Well Documented**: Comprehensive documentation and examples

## Features

- ✅ Single allocation during encoding
- ✅ Static typing with natural Go type mappings
- ✅ Human-readable ABI support
- ✅ Support for all Solidity types (arrays, tuples, events)
- ✅ Event generation with topic encoding/decoding
- ✅ Performance benchmarks
- ✅ Comprehensive test suite
- ✅ No panics - all functions return errors

## Quick Start

### Installation

```bash
go get github.com/yihuang/go-abi/cmd
```

### From Human Readable ABI

Define ABI directly in your Go source file:

```go
package mycontract

//go:generate go run github.com/yihuang/go-abi/cmd -var ContractABI -module mycontract

var ContractABI = []string{
    "function transfer(address to, uint256 amount) returns (bool)",
    "function balanceOf(address account) view returns (uint256)",
    "event Transfer(address indexed from, address indexed to, uint256 value)",
}
```

Then run:

```bash
go generate ./...
```

### From JSON ABI Files

If you have an existing JSON ABI file:

```bash
go run github.com/yihuang/go-abi/cmd -input contract.abi.json -module mycontract -package mycontract
```

## Usage Examples

### Encoding Function Calls

```go
package main

import (
    "fmt"
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/yourusername/project/mycontract"
)

func main() {
    // Create transfer call
    call := mycontract.TransferCall{
        To:     common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F2"),
        Amount: big.NewInt(1000000000000000000), // 1 ETH
    }

    // Encode with function selector
    data, err := call.EncodeWithSelector()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Encoded data: %x\n", data)
    // Output: Encoded data: a9059cbb...
}
```

### Decoding Function Returns

```go
// Decode return value from contract call
encodedResult := []byte{...} // bytes from contract call

var result mycontract.BalanceOfReturn
if err := result.Decode(encodedResult); err != nil {
    panic(err)
}

fmt.Printf("Balance: %s\n", result.Balance)
```

### Working with Events

```go
package main

import (
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/yourusername/project/mycontract"
)

// Process transfer events from logs
func processTransfer(topics [][]byte, data []byte) error {
    // Create event struct
    var transfer mycontract.Transfer

    // Decode indexed fields (From, To)
    if err := transfer.DecodeTopics(topics); err != nil {
        return fmt.Errorf("decode topics: %w", err)
    }

    // Decode non-indexed fields (Value)
    var transferData mycontract.TransferData
    if err := transferData.Decode(data); err != nil {
        return fmt.Errorf("decode data: %w", err)
    }

    fmt.Printf("Transfer: %s -> %s, Amount: %s\n",
        transfer.From, transfer.To, transferData.Value)

    return nil
}
```

### Complex Types (Tuples)

```go
package mycontract

// Define ABI with struct types
var UserContractABI = []string{
    "struct User { address addr; uint256 balance; string name; }",
    "function getUser(address id) view returns (User)",
    "function setUser(User user)",
}

// Generated code usage:
func example() {
    // Create user struct
    user := mycontract.User{
        Addr:    common.HexToAddress("0x123..."),
        Balance: big.NewInt(1000),
        Name:    "Alice",
    }

    // Encode
    encoded, err := user.Encode()
    if err != nil {
        panic(err)
    }

    // Decode
    var decoded mycontract.User
    if err := decoded.Decode(encoded); err != nil {
        panic(err)
    }
}
```

### Arrays and Dynamic Types

```go
package mycontract

var ArrayContractABI = []string{
    "function batchTransfer(address[] recipients, uint256[] amounts)",
}

// Generated code usage:
func example() {
    call := mycontract.BatchTransferCall{
        Recipients: []common.Address{
            common.HexToAddress("0x111..."),
            common.HexToAddress("0x222..."),
        },
        Amounts: []*big.Int{
            big.NewInt(100),
            big.NewInt(200),
        },
    }

    encoded, err := call.EncodeWithSelector()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Encoded: %x\n", encoded)
}
```

### High-Performance Encoding

For high-throughput applications, use `EncodeTo` with pre-allocated buffers:

```go
func processManyTransfers(transfers []mycontract.TransferCall) [][]byte {
    results := make([][]byte, len(transfers))

    for i, transfer := range transfers {
        // Pre-calculate size
        size := transfer.EncodedSize()

        // Allocate buffer once
        buf := make([]byte, size)

        // Encode to buffer
        if _, err := transfer.EncodeTo(buf); err != nil {
            panic(err)
        }

        results[i] = buf
    }

    return results
}
```

## Type Mappings

The generator maps Solidity types to Go types as follows:

| Solidity Type | Go Type |
|---------------|---------|
| `address` | `common.Address` |
| `uint8` | `uint8` |
| `int8` | `int8` |
| `uint16` | `uint16` |
| `int16` | `int16` |
| `uint[24,32]` | `uint32` |
| `int[24,32]` | `int32` |
| `uint[40,48,56,64]` | `uint64` |
| `int[40,48,56,64]` | `int64` |
| `uint[64+]` | `*big.Int` |
| `int[64+]` | `*big.Int` |
| `bool` | `bool` |
| `string` | `string` |
| `bytes` | `[]byte` |
| `bytesN` | `[N]byte` |
| `type[]` | `[]GoType` |
| `type[N]` | `[N]GoType` |

## Command-Line Options

The generator supports several options:

```bash
# From Go source file with human-readable ABI
go run github.com/yihuang/go-abi/cmd \
    -input mycontract.go \
    -var ContractABI \
    -module mycontract \
    -package mycontract \
    -imports github.com/ethereum/go-ethereum/common

# From JSON ABI file
go run github.com/yihuang/go-abi/cmd \
    -input contract.abi.json \
    -module mycontract \
    -package mycontract

# With external tuple mappings
go run github.com/yihuang/go-abi/cmd \
    -input mycontract.go \
    -var ContractABI \
    -external-tuples User=User,Group=Group
```

Options:
- `-input`: Input file (JSON ABI or Go source file)
- `-var`: Variable name containing human-readable ABI (for Go source files)
- `-module`: Output module name (generates `{module}.abi.go`)
- `-package`: Package name for generated code
- `-imports`: Additional import paths (comma-separated)
- `-external-tuples`: External tuple type mappings

## Performance

The go-abi generator produces highly optimized code:

- **Single Allocation Pattern**: Pre-calculate size, allocate once
- **Zero-Copy Operations**: Direct buffer manipulation
- **Bounds Checking**: Safe without performance penalties
- **Static Size Constants**: For performance-critical code

See [benchmarks](tests/encode_benchmark_test.go) for detailed performance comparisons with go-ethereum.

Example benchmark results:

```
BenchmarkGoABI_ComplexDynamicTuples-8   	    2000	     500000 ns/op	  100000 B/op	    1000 allocs/op
BenchmarkGoEthereum_ComplexDynamicTuples-8 	   1000	    1000000 ns/op	  200000 B/op	    2000 allocs/op
```

## API Reference

For detailed API documentation, see:
- [Code Generator Design](docs/code-generator.md)
- [Events API](docs/events-api.md)
- [Events Quick Reference](docs/events-quick-reference.md)

## Examples

See the following for complete examples:
- [examples/](examples/) - Simple ERC20 and transfer examples
- [tests/](tests/) - Comprehensive test cases with complex types
- [CONTRIBUTING.md](CONTRIBUTING.md) - Contribution guidelines and testing

## References

- [Ethereum ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [Human Readable ABI](https://abitype.dev/api/human)
- [go-ethereum](https://github.com/ethereum/go-ethereum)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT License](LICENSE)

## References

* [ABI Specification](https://github.com/argotorg/solidity/blob/v0.8.30/docs/abi-spec.rst)
* [Human Readable ABI](https://abitype.dev/api/human)
