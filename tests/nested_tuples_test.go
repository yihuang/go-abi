package testdata

import (
	"bytes"
	"strings"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

func TestNestedTuplesInHumanReadableABI(t *testing.T) {
	// Test nested tuples in human readable ABI
	humanABI := []string{
		"struct Address { string street; string city; }",
		"struct User { string name; Address addr; uint256 balance; }",
		"function createUser(User user)",
		"function batchCreateUsers(User[] users)",
		"function complexNested(User[2] fixedUsers, User[] dynamicUsers)",
	}

	// Parse human-readable ABI
	abiJSON, err := abi.ParseHumanReadableABI(humanABI)
	if err != nil {
		t.Fatalf("Failed to parse human-readable ABI: %v", err)
	}

	// Parse JSON ABI
	abiDef, err := ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		t.Fatalf("Failed to parse JSON ABI: %v", err)
	}

	// Verify we have the expected functions
	expectedFunctions := []string{"createUser", "batchCreateUsers", "complexNested"}
	for _, expectedFunc := range expectedFunctions {
		found := false
		for name := range abiDef.Methods {
			if strings.HasPrefix(name, expectedFunc) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected function %s not found in ABI", expectedFunc)
		}
	}

	// Verify structs are properly parsed
	// The human readable ABI parser should have converted structs to tuples
	// We can verify this by checking that the generated code compiles
	generator := abi.NewGenerator("testdata")
	generatedCode, err := generator.GenerateFromABI(abiDef)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Basic sanity check on generated code
	if !strings.Contains(generatedCode, "type Address struct") {
		t.Error("Generated code should contain Address struct")
	}
	if !strings.Contains(generatedCode, "type User struct") {
		t.Error("Generated code should contain User struct")
	}
	if !strings.Contains(generatedCode, "type CreateUserCall struct") {
		t.Error("Generated code should contain CreateUserCall struct")
	}

	t.Log("✅ Nested tuples in human readable ABI are working correctly!")
}