# Go ABI Code Generator

Generating statically typed Go structs and encode/decode methods from Ethereum contract ABI definitions.

## Why

- Deterministic
- Static typing
- Efficiency

## Features

- **Single Allocation**,  allocate buffer once during the encoding.
- **Static Typing**, natural type mapping from ABI to golang.
- **Human Readable ABI Support**, generate code directly from human-readable ABI definitions.

## Quick Start

### From Human Readable ABI

Define ABI directly in Go source files:

```go
//go:generate go run github.com/yihuang/go-abi/cmd -var ContractABI -output contract.abi.go

var ContractABI = []string{
    "function transfer(address to, uint256 amount) returns (bool)",
    "function balanceOf(address account) view returns (uint256)",
    "event Transfer(address indexed from, address indexed to, uint256 value)",
}
```

### From JSON ABI Files

Embed the generator command in code comments:

```bash
//go:generate go run github.com/yihuang/go-abi/cmd -input test_abi.json -output test_abi.abi.go
```

### Sync

Run `go generate` to sync all the files.

```bash
$ go generate ./...
```

## Example

See `tests/` and `examples/` directories for example ABIs and generated code.

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

## References

* [ABI Specification](https://github.com/argotorg/solidity/blob/v0.8.30/docs/abi-spec.rst)
* [Human Readable ABI](https://abitype.dev/api/human)
