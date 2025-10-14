# Go ABI Code Generator

A comprehensive Go library for generating Go structs and encoding methods from Ethereum contract ABI definitions.

## Overview

This library provides an ABI code generator that automatically generates:
- Go structs for function arguments.
- Encode/Decode methods without runtime type reflection.
- Proper type mappings between Solidity and Go types

## Features

- **Type Safety**: Work with statically typed structures.
- **Single Allocation**: Allocate buffer once during the encoding.

## Installation

```bash
go get github.com/yihuang/go-abi
```

## Quick Start

### Go Generate

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
package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"github.com/yihuang/go-abi"
)

func Test() {
	// Use the generated code
	args := &TransferArgs{
		To:     common.HexToAddress("0x742d35Cc6634C0532925a3b8DfBc6A8c4f4C9F7F"),
		Amount: big.NewInt(1000),
	}

	// Encode the arguments
	encoded, err := args.Encode()
	if err != nil {
		panic(err)
	}

	// Get the function selector
	selector := args.Selector()

	fmt.Printf("Encoded data: %x\n", encoded)
	fmt.Printf("Function selector: %x\n", selector)
}
```

## Example

See `tests/` directory for example ABIs and generated code.

## Type Mappings

The generator maps Solidity types to Go types as follows:

| Solidity Type | Go Type | Example |
|---------------|---------|---------|
| `address` | `common.Address` | `common.HexToAddress("0x...")` |
| `uintN` | `*big.Int` | `big.NewInt(1000)` |
| `intN` | `*big.Int` | `big.NewInt(-100)` |
| `bool` | `bool` | `true` |
| `string` | `string` | `"hello"` |
| `bytes` | `[]byte` | `[]byte{0x01, 0x02}` |
| `bytesN` | `[N]byte` | `[32]byte{...}` |
| `type[]` | `[]GoType` | `[]common.Address` |
| `type[N]` | `[N]GoType` | `[10]*big.Int` |
