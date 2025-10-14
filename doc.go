/*
Package abi provides a Go library for generating Go structs and encoding methods from Ethereum contract ABI definitions.

Overview

This library provides an ABI code generator that automatically generates:
- Go structs for function arguments
- Encode() methods for ABI encoding
- Function selectors for contract calls
- Proper type mappings between Solidity and Go types

Quick Start

Generate code from ABI JSON using the command-line tool:

	go run cmd/generate/main.go -input contract.abi.json -package mycontract

Or use programmatically:

	generator := abi.NewGenerator("mypackage")
	generatedCode, err := generator.GenerateFromABI(abiJSON)

Example

Input ABI:

	[
		{
			"name": "transfer",
			"type": "function",
			"inputs": [
				{"name": "to", "type": "address"},
				{"name": "amount", "type": "uint256"}
			],
			"outputs": [{"name": "", "type": "bool"}]
		}
	]

Generated Code:

	package mypackage

	import (
		"github.com/ethereum/go-ethereum/common"
		"math/big"
	)

	// TransferArgs represents the arguments for transfer function
	type TransferArgs struct {
		To     common.Address `json:"to"`
		Amount *big.Int       `json:"amount"`
	}

	// Encode encodes transfer arguments to ABI bytes
	func (args *TransferArgs) Encode() ([]byte, error) {
		values := []interface{}{}
		values = append(values, args.To)
		values = append(values, args.Amount)
		return Encode(values)
	}

	// Selector returns the function selector for transfer
	func (*TransferArgs) Selector() [4]byte {
		return transferSelector
	}

	// transferSelector is the function selector for transfer(address,uint256)
	var transferSelector = [4]byte{0xa9, 0x05, 0x9c, 0xbb}

Features

- Automatic code generation from ABI JSON
- Type-safe Go structs for function arguments
- ABI encoding using go-ethereum library
- Function selector generation
- Support for arrays and complex types
- Command-line interface and programmatic API

Type Mappings

Solidity types are mapped to Go types as follows:

	address     -> common.Address
	uintN/intN  -> *big.Int
	bool        -> bool
	string      -> string
	bytes       -> []byte
	bytesN      -> [N]byte
	type[]      -> []GoType
	type[N]     -> [N]GoType

See the examples directory for complete usage examples.
*/
package abi