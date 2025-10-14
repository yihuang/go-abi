# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go library for generating Go structs and encoding methods from Ethereum contract ABI definitions. The library automatically generates type-safe Go code for interacting with smart contracts without runtime type reflection.

## Key Architecture

### Core Components

- **Generator** (`generator.go`): Main code generation engine that converts ABI definitions to Go code
- **Type System** (`types.go`): Helper functions for ABI type analysis (dynamic type detection, size calculation, etc.)
- **Command Line Tool** (`cmd/main.go`): CLI interface for generating code from ABI JSON files

### Code Generation Process

1. **ABI Parsing**: Uses `github.com/ethereum/go-ethereum/accounts/abi` to parse ABI JSON
2. **Type Mapping**: Converts Solidity types to Go types with optimizations for common integer sizes
3. **Struct Generation**: Creates Go structs for function arguments with proper field names
4. **Encoding Methods**: Generates `Encode()`, `EncodeTo()`, `EncodeWithSelector()` methods for ABI encoding
5. **Size Calculation**: Generates `TotalSize()` methods for pre-calculating encoded data size

### Type Mappings

- `address` → `common.Address`
- `uintN/intN` (8,16,32,64) → native Go types (`uint8`, `int16`, etc.)
- `uintN/intN` (other sizes) → `*big.Int`
- `bool` → `bool`
- `string` → `string`
- `bytes` → `[]byte`
- `bytesN` → `[N]byte`
- `type[]` → `[]GoType`
- `type[N]` → `[N]GoType`

## Development Commands

### Building and Testing

```bash
# Generate abi code
go generate ./...

# Run all tests
go test ./...
```

## Usage Patterns

### Generated Code Structure

For each function in the ABI, the generator creates:
- `{FunctionName}Args` struct with typed fields
- `Encode()` method for ABI encoding
- `EncodeWithSelector()` method including function selector
- `EncodeTo(buf []byte)` method for direct buffer encoding
- `TotalSize()` method for size calculation
- `Selector()` method for function selector

### Performance Optimizations

- Uses native Go types for 8, 16, 32, 64-bit integers to avoid `big.Int` allocations
- Direct buffer encoding with `EncodeTo()` method to avoid temporary allocations
- Pre-calculated size methods for efficient buffer allocation
- Optimized offset/length encoding using `encoding/binary.BigEndian`

## Testing Strategy

Tests are located in the `tests/` directory and verify:
- Generated code compiles correctly
- ABI encoding matches go-ethereum's implementation
- Complex type support (arrays, tuples, nested structures)
- Performance characteristics

Tests use `go:generate` directives to regenerate test code from ABI JSON files.

## Key Files

- `generator.go`: Main code generation logic
- `types.go`: ABI type analysis utilities
- `cmd/main.go`: CLI interface
- `tests/`: Test suite with example ABIs
- `doc.go`: Package documentation
- `abi-spec.rst`: ABI Specification

## Dependencies

- `github.com/ethereum/go-ethereum`: For ABI parsing and reference encoding
- `golang.org/x/text`: For proper case conversion in field names
