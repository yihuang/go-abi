package generator

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestExternalTuples(t *testing.T) {
	// Define a simple ABI with a tuple
	abiJSON := `[
		{
			"type": "function",
			"name": "processUserData",
			"inputs": [
				{
					"name": "data",
					"type": "tuple",
					"components": [
						{"name": "address", "type": "address"},
						{"name": "name", "type": "string"},
						{"name": "amount", "type": "uint256"}
					]
				}
			],
			"outputs": []
		}
	]`

	abiDef, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		t.Fatalf("Failed to parse ABI: %v", err)
	}

	// Test without ExternalTuples option (should generate tuple)
	generator := NewGenerator()
	code, err := generator.GenerateFromABI(abiDef)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Should contain the generated tuple struct
	if !contains(code, "type ProcessUserDataCall struct") {
		t.Error("Expected generated code to contain ProcessUserDataCall struct")
	}

	// Test with ExternalTuples option (should not generate tuple)
	// The tuple name is generated from the tuple structure, not the function name
	extTuples := map[string]string{
		"Tuple_b53c1574": "MyCustomUserData",
	}

	generatorWithExternal := NewGenerator(ExternalTuples(extTuples))
	codeWithExternal, err := generatorWithExternal.GenerateFromABI(abiDef)
	if err != nil {
		t.Fatalf("Failed to generate code with external tuples: %v", err)
	}

	// Should NOT contain the generated nested tuple struct
	if contains(codeWithExternal, "type Tuple_b53c1574 struct") {
		t.Error("Expected generated code to NOT contain Tuple_b53c1574 struct when using external tuple")
	}

	// Should use the external type name in the function signature
	if !contains(codeWithExternal, "MyCustomUserData") {
		t.Error("Expected generated code to use external tuple type name")
	}

	// The function input struct should still be generated but use the external type
	if !contains(codeWithExternal, "Data MyCustomUserData") {
		t.Error("Expected function input struct to use external tuple type")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}
