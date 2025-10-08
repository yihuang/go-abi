package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/huangyi/go-abi/pkg/abi"
)

func main() {
	// Example: Encoding a simple transfer
	to := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	value := big.NewInt(1000000)

	encoded, err := abi.EncodeSimpleTransfer(to, value)
	if err != nil {
		log.Fatal("Encoding failed:", err)
	}

	fmt.Printf("Encoded transfer call: %s\n", hex.EncodeToString(encoded))

	// Example: Using the ABI directly
	abiInstance := abi.ExampleUint256()

	// Encode using ABI
	encoded2, err := abiInstance.EncodeFunctionCall("transfer", to, value)
	if err != nil {
		log.Fatal("Encoding failed:", err)
	}

	fmt.Printf("Encoded via ABI: %s\n", hex.EncodeToString(encoded2))

	// Example of basic type encoding
	fmt.Println("\nBasic type encoding examples:")

	// Encode a uint256
	uint256Type := &abi.Uint256{}
	encodedUint, err := uint256Type.Encode(big.NewInt(42))
	if err != nil {
		log.Fatal("Uint256 encoding failed:", err)
	}
	fmt.Printf("Uint256(42): %s\n", hex.EncodeToString(encodedUint))

	// Encode a bool
	boolType := &abi.Bool{}
	encodedBool, err := boolType.Encode(true)
	if err != nil {
		log.Fatal("Bool encoding failed:", err)
	}
	fmt.Printf("Bool(true): %s\n", hex.EncodeToString(encodedBool))

	// Encode an address
	addrType := &abi.Address{}
	addr := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	encodedAddr, err := addrType.Encode(addr)
	if err != nil {
		log.Fatal("Address encoding failed:", err)
	}
	fmt.Printf("Address: %s\n", hex.EncodeToString(encodedAddr))

	fmt.Println("\nLibrary initialized successfully!")
}