package tests

import (
	"bytes"
	"strings"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/test-go/testify/require"
	"github.com/yihuang/go-abi"
)

//go:generate go run ../cmd -var OverloadABI -output overload.abi.go

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
	require.NoError(t, err, "Failed to parse human-readable ABI with overloaded functions")

	// Parse JSON ABI
	abiDef, err := ethabi.JSON(bytes.NewReader(abiJSON))
	require.NoError(t, err)

	// Check that we have multiple overloaded1 functions
	// Go-ethereum renames overloaded functions with suffixes
	overloaded1Count := 0
	for name := range abiDef.Methods {
		if strings.HasPrefix(name, "overloaded1") {
			overloaded1Count++
		}
	}

	require.Equal(t, 3, overloaded1Count, "Expected 3 overloaded1 functions")

	// Check that we have multiple overloaded2 functions
	overloaded2Count := 0
	for name := range abiDef.Methods {
		if strings.HasPrefix(name, "overloaded2") {
			overloaded2Count++
		}
	}

	require.Equal(t, 2, overloaded2Count, "Expected 2 overloaded2 functions")

	// Verify function signatures are different
	overloaded1Signatures := make(map[string]bool)
	for name, method := range abiDef.Methods {
		if strings.HasPrefix(name, "overloaded1") {
			overloaded1Signatures[method.Sig] = true
		}
	}

	require.Equal(t, 3, len(overloaded1Signatures), "Expected 3 unique overloaded1 signatures")

	overloaded2Signatures := make(map[string]bool)
	for name, method := range abiDef.Methods {
		if strings.HasPrefix(name, "overloaded2") {
			overloaded2Signatures[method.Sig] = true
		}
	}

	require.Equal(t, 2, len(overloaded2Signatures), "Expected 2 unique overloaded2 signatures")
}
