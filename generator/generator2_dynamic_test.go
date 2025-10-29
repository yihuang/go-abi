package generator

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestDynamicTypes(t *testing.T) {
	// Test ABI with dynamic types
	abiJSON := `[
		{
			"name": "setMessage",
			"type": "function",
			"inputs": [
				{
					"name": "message",
					"type": "string"
				},
				{
					"name": "data",
					"type": "bytes"
				}
			],
			"outputs": [
				{
					"name": "success",
					"type": "bool"
				}
			]
		},
		{
			"name": "addItems",
			"type": "function",
			"inputs": [
				{
					"name": "items",
					"type": "string[]"
				}
			],
			"outputs": [
				{
					"name": "count",
					"type": "uint256"
				}
			]
		}
	]`

	abiDef, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		t.Fatalf("Failed to parse ABI: %v", err)
	}

	// Test that we can generate code for dynamic types
	generator := NewGenerator2()
	generatedCode, err := generator.GenerateFromABI(abiDef)
	if err != nil {
		t.Fatalf("Failed to generate code for dynamic types: %v", err)
	}

	// Check that the generated code contains expected dynamic type functions
	expectedFunctions := []string{
		"encode_string",
		"encode_bytes",
		"encode_string_array",
	}

	for _, expectedFunc := range expectedFunctions {
		if !strings.Contains(generatedCode, expectedFunc) {
			t.Errorf("Generated code missing expected function: %s", expectedFunc)
		}
	}

	// Check that dynamic types are properly handled in tuple encoding
	if !strings.Contains(generatedCode, "dynamicSize") {
		t.Error("Generated code should handle dynamic size calculations")
	}

	// Check that size functions are generated for dynamic types
	expectedSizeFunctions := []string{
		"size_string",
		"size_bytes",
		"size_string_array",
	}

	for _, expectedFunc := range expectedSizeFunctions {
		if !strings.Contains(generatedCode, expectedFunc+"(") {
			t.Errorf("Generated code missing expected size function: %s", expectedFunc)
		}
	}
}