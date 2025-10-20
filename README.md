# Go ABI Code Generator

Generating statically typed Go structs and encode/decode methods from Ethereum contract ABI definitions.

[ABI Specification](https://github.com/argotorg/solidity/blob/v0.8.30/docs/abi-spec.rst)

## Why

- Deterministic
- Static typing
- Efficiency

## Features

- **Single Allocation**,  allocate buffer once during the encoding.
- **Static Typing**, natural type mapping from ABI to golang.

## Quick Start

Embed the generator command in code comments:

```bash
//go:generate go run github.com/yihuang/go-abi/cmd -input test_abi.json -output test_abi.abi.go
```

Run `go generate` to sync the generated files.

```bash
$ go generate ./...
```

## Usage Examples

### Basic Example

```go
// Use the generated code
args := &TransferArgs{
	To:     common.HexToAddress("0x742d35Cc6634C0532925a3b8DfBc6A8c4f4C9F7F"),
	Amount: big.NewInt(1000),
}

// Encode the arguments
encoded, err := args.Encode()

// Encode with selector
encoded, err := args.EncodeWithSelector()
```

## Example

See `tests/` directory for example ABIs and generated code.

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
