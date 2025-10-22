package testdata

import (
	"bytes"
	"strings"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

//go:generate go run ../cmd/main.go -var OverloadABI -output overload.abi.go

var OverloadABI = []string{
	"function overloaded1(address to, uint256 amount) returns (bool)",
	"function overloaded1(address from, address to, uint256 amount) returns (bool)",
	"function overloaded1(address from, address to, uint256 amount, bytes data) returns (bool)",
	"function overloaded2(address account) view returns (uint256)",
	"function overloaded2() view returns (uint256)",
}

func TestOverloadedFunctions(t *testing.T) {
	// Parse human-readable ABI
	abiJSON, err := abi.ParseHumanReadableABI(OverloadABI)
	if err != nil {
		t.Fatalf("Failed to parse human-readable ABI: %v", err)
	}

	// Parse JSON ABI
	abiDef, err := ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		t.Fatalf("Failed to parse JSON ABI: %v", err)
	}

	// Check that we have multiple overloaded1 functions
	// Go-ethereum renames overloaded functions with suffixes
	overloaded1Count := 0
	for name := range abiDef.Methods {
		if strings.HasPrefix(name, "overloaded1") {
			overloaded1Count++
		}
	}

	if overloaded1Count != 3 {
		t.Errorf("Expected 3 overloaded1 functions, got %d", overloaded1Count)
	}

	// Check that we have multiple overloaded2 functions
	overloaded2Count := 0
	for name := range abiDef.Methods {
		if strings.HasPrefix(name, "overloaded2") {
			overloaded2Count++
		}
	}

	if overloaded2Count != 2 {
		t.Errorf("Expected 2 overloaded2 functions, got %d", overloaded2Count)
	}

	// Verify function signatures are different
	overloaded1Signatures := make(map[string]bool)
	for name, method := range abiDef.Methods {
		if strings.HasPrefix(name, "overloaded1") {
			overloaded1Signatures[method.Sig] = true
		}
	}

	if len(overloaded1Signatures) != 3 {
		t.Errorf("Expected 3 unique overloaded1 signatures, got %d", len(overloaded1Signatures))
	}

	overloaded2Signatures := make(map[string]bool)
	for name, method := range abiDef.Methods {
		if strings.HasPrefix(name, "overloaded2") {
			overloaded2Signatures[method.Sig] = true
		}
	}

	if len(overloaded2Signatures) != 2 {
		t.Errorf("Expected 2 unique overloaded2 signatures, got %d", len(overloaded2Signatures))
	}
}
