# Go ABI Code Generator

Generating statically typed Go structs and encode/decode methods from Ethereum contract ABI definitions.

[![Go Reference](https://pkg.go.dev/badge/github.com/yihuang/go-abi.svg)](https://pkg.go.dev/github.com/yihuang/go-abi)
[![Go Report Card](https://goreportcard.com/badge/github.com/yihuang/go-abi)](https://goreportcard.com/report/github.com/yihuang/go-abi)
[![Tests](https://github.com/yihuang/go-abi/actions/workflows/go.yml/badge.svg)](https://github.com/yihuang/go-abi/actions/workflows/go.yml)

## Why go-abi?

### 🚀 **High Performance**
- **Single Allocation**: Allocate buffer once during encoding for maximum efficiency
- **Zero-Copy Decoding**: Direct memory operations with proper bounds checking

### 🛡️ **Type Safety**
- **Static Typing**: Natural type mapping from ABI to Go
- **No Reflection**: Generated code is fast and type-safe

### 📝 **Developer Experience**
- **Human Readable ABI**: Write ABIs in readable format directly in Go
- **Auto-Generation**: Keep code in sync with `go generate`

## Quick Start

### From Human Readable ABI

Define ABI directly in your Go source file:

```go
package mycontract

//go:generate go run github.com/yihuang/go-abi/cmd -var ContractABI -output mycontract.abi.go

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
go run github.com/yihuang/go-abi/cmd -input contract.abi.json -output mycontract.abi.go
```

## Usage Examples

### Call Functions

```go
// Create transfer call
call := erc20.TransferCall{
    To:     common.HexToAddress("0x1000000000000000000000000000000000000000"),
    Amount: big.NewInt(1000000000000000000), // 1 ETH
}

// Encode with function selector
data, err := call.EncodeWithSelector()

// Execute EVM function with calldata
ret, err := CallEVM(data)

// Decode return value
var result erc20.BalanceOfReturn
_, err := result.Decode(ret)

fmt.Printf("Balance: %s\n", result.Balance)
```

### Working with Events

```go
// Create event struct
var transfer erc20.TransferEvent

// Decode topics and data
err := abi.DecodeEvent(&transfer, log.Topics, log.Data)

// Encode events to topics and data
topics, data, err := abi.EncodeEvent(&transfer)
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

## Performance

See [benchmarks](tests/encode_benchmark_test.go) for detailed performance comparisons with go-ethereum.

Benchmark results:

```
$ go test -bench=. -benchmem -benchtime=5s ./tests/...
goos: darwin
goarch: arm64
pkg: github.com/yihuang/go-abi/tests
cpu: Apple M3 Max
BenchmarkGoABI_ComplexDynamicTuples-16                           	14185347	      440.3 ns/op	   3072 B/op	      1 allocs/op
BenchmarkGoABI_NestedDynamicArrays-16                            	20589320	      287.8 ns/op	   1280 B/op	      1 allocs/op
BenchmarkGoABI_MixedTypes-16                                     	44322256	      133.8 ns/op	    896 B/op	      1 allocs/op
BenchmarkGoEthereum_ComplexDynamicTuples-16                      	 415578	    14214 ns/op	  49385 B/op	    372 allocs/op
BenchmarkGoEthereum_NestedDynamicArrays-16                       	1378192	     4327 ns/op	  13504 B/op	    151 allocs/op
BenchmarkGoEthereum_MixedTypes-16                                	1391475	     4309 ns/op	   9296 B/op	    116 allocs/op
BenchmarkGoABI_EncodeOnly_ComplexDynamicTuples-16                	14153460	      443.3 ns/op	   3072 B/op	      1 allocs/op
BenchmarkGoABI_EncodeTo_ComplexDynamicTuples-16                  	41144373	      146.1 ns/op	      0 B/op	      0 allocs/op
BenchmarkGoABI_MemoryAllocations_ComplexDynamicTuples-16         	13971866	      465.7 ns/op	   3072 B/op	      1 allocs/op
BenchmarkGoEthereum_MemoryAllocations_ComplexDynamicTuples-16    	 383490	    14732 ns/op	  49385 B/op	    372 allocs/op
PASS
ok  	github.com/yihuang/go-abi/tests	71.688s
```

## References

- [Ethereum ABI Specification](https://docs.soliditylang.org/en/latest/abi-spec.html)
- [Human Readable ABI](https://abitype.dev/api/human)
- [go-ethereum](https://github.com/ethereum/go-ethereum)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

[MIT License](LICENSE)
