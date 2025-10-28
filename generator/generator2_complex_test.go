package generator

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestGenerator2ComplexFunction(t *testing.T) {
	// Test with a complex ABI including arrays and tuples
	abiJSON := `[
		{
			"type": "function",
			"name": "transferBatch",
			"inputs": [
				{"name": "recipients", "type": "address[]"},
				{"name": "amounts", "type": "uint256[]"}
			],
			"outputs": [],
			"stateMutability": "nonpayable"
		},
		{
			"type": "function",
			"name": "processUserData",
			"inputs": [
				{
					"name": "user",
					"type": "tuple",
					"components": [
						{"name": "address", "type": "address"},
						{"name": "name", "type": "string"},
						{"name": "age", "type": "uint256"}
					]
				}
			],
			"outputs": [],
			"stateMutability": "nonpayable"
		}
	]`

	abiDef, err := abi.JSON(bytes.NewReader([]byte(abiJSON)))
	if err != nil {
		t.Fatalf("Failed to parse ABI: %v", err)
	}

	generator := NewGenerator2(PackageName("test"))
	code, err := generator.GenerateFromABI(abiDef)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Basic validation of generated code
	if len(code) == 0 {
		t.Fatal("Generated code is empty")
	}

	// Check that all expected encoding functions are generated
	expectedFunctions := []string{
		"encode_address",
		"encode_uint256",
		"encode_address_array",
		"encode_uint256_array",
	}

	for _, funcName := range expectedFunctions {
		if !strings.Contains(code, funcName) {
			t.Errorf("Generated code missing function: %s", funcName)
		}
	}

	// Check that structs are generated for tuples
	if !strings.Contains(code, "type ") || !strings.Contains(code, "struct {") {
		t.Error("Generated code missing struct definitions")
	}

	// Check that function selectors are generated
	if !strings.Contains(code, "TransferBatchSelector") || !strings.Contains(code, "ProcessUserDataSelector") {
		t.Error("Generated code missing function selectors")
	}

	// Check that EncodeWithSelector methods are generated
	if !strings.Contains(code, "EncodeWithSelector") {
		t.Error("Generated code missing EncodeWithSelector methods")
	}

	t.Logf("Generated code length: %d bytes", len(code))
	t.Logf("Generated functions found: %d", strings.Count(code, "func encode_"))
	t.Logf("Generated structs found: %d", strings.Count(code, "type ")-1) // -1 for package declaration

	// Print a sample of the generated code
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if strings.Contains(line, "func encode_") {
			t.Logf("Sample encoding function: %s", strings.TrimSpace(line))
			if i+2 < len(lines) {
				t.Logf("Next line: %s", strings.TrimSpace(lines[i+1]))
			}
			break
		}
	}
}