# go-abi

A Go implementation of Ethereum ABI encoding/decoding using code generation to avoid reliance on reflection.

## Features

- **Code Generation**: Generate type-safe ABI bindings at compile time
- **No Reflection**: Avoid runtime reflection for better performance
- **Type Safety**: Compile-time validation of ABI types
- **Ethereum Compatible**: Full compatibility with Ethereum ABI specification

## Installation

```bash
go get github.com/huangyi/go-abi
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/huangyi/go-abi"
)

func main() {
	// Your ABI usage here
}
```

## Code Generation

Use the provided code generator to create type-safe bindings from ABI JSON:

```bash
go run cmd/generator/main.go -abi contract.abi -out bindings.go
```

## License

MIT