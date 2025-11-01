# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive benchmarks comparing go-abi performance with go-ethereum
- Detailed documentation for events generation API
- Quick reference guide for common event patterns
- Stdlib implementation for common primitive types to reduce code duplication

### Changed
- Improved code generator to be more composable and maintainable
- Enhanced error handling with proper error wrapping

## [0.1.0] - 2024-XX-XX

### Added
- Support for human-readable ABI format
- Code generation from both JSON and human-readable ABIs
- Support for all standard Solidity types:
  - Primitive types (uint, int, bool, address, string, bytes)
  - Fixed and dynamic arrays
  - Tuples and nested structs
  - Events with indexed and non-indexed parameters
- Function selector generation (both byte array and uint32 constants)
- Encode/Decode methods with buffer reuse support
- Size calculation methods for dynamic types
- Support for external tuple types
- Support for function overloading
- Support for small integers (non-standard bit sizes)
- Comprehensive test suite with 61 test files
- Benchmark tests for performance comparison
- Command-line tool for code generation
- Programmatic API for integration

### Features
- **Single Allocation**: Allocate buffer once during encoding for efficiency
- **Static Typing**: Natural type mapping from ABI to Go
- **Human Readable ABI**: Generate code directly from human-readable definitions
- **Error-Free**: All functions return errors instead of panicking
- **Composable**: Generated code is modular and composable
- **Performance**: Optimized for minimal allocations and maximum throughput

### Generator Capabilities
- Parse JSON ABI files
- Parse human-readable ABI from Go source variables
- Generate encode/decode functions for all types
- Generate tuple structs with methods
- Generate function and event selectors
- Support for custom imports and external tuples
- Type-safe code generation with proper Go types

### Supported Type Mappings
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

### Internal Changes
- Refactored code generator for better composability
- Implemented stdlib for primitive types to reduce duplication
- Fixed import handling to properly clone import slices
- Made function prefix optional in generated code
- Fixed handling of empty functions
- Improved nested array parsing and encoding
- Fixed event generation to include top-level event structs
- Separated indexed fields handling for events
- Optimized event generated code for compactness
- Implemented nested struct handling
- Added alias import support for extra imports
- Added command flag for extra imports
- Fixed name conflicts with multiple dynamic tuples
- Implemented external tuples functionality
- Made generator more customizable

### Bug Fixes
- Fixed nested fixed size array missing bound check
- Fixed bug in nested fixed size array handling
- Fixed human readable parser not handling nested arrays correctly
- Fixed event indexed fields not being handled separately
- Fixed empty functions not being generated
- Fixed imports slice not being cloned
- Fixed function prefix not being optional

### Documentation
- Added README.md with quick start guide
- Added comprehensive code generator design documentation
- Added events API documentation
- Added events quick reference guide
- Added package documentation in doc.go
- Fixed typos and improved clarity in documentation

### Testing
- Added comprehensive test suite covering all features
- Added benchmark tests for performance comparison
- Added tests for small integers and non-standard types
- Added tests for nested arrays and tuples
- Added tests for events encoding/decoding
- Added tests for function overloading
- Added tests for external tuples

### Performance
- Optimized encoding for single allocation pattern
- Reduced code duplication through stdlib functions
- Improved buffer reuse with EncodeTo methods
- Static size constants for performance-critical code
- Benchmarks show competitive or better performance vs go-ethereum

[Unreleased]: https://github.com/yihuang/go-abi/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/yihuang/go-abi/releases/tag/v0.1.0
