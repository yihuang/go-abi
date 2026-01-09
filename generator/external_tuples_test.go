package generator

import (
	"testing"
)

func TestExternalTuples(t *testing.T) {
	abiDef := mustParseABI(t, `[{
		"type": "function",
		"name": "processUserData",
		"inputs": [{
			"name": "data",
			"type": "tuple",
			"components": [
				{"name": "address", "type": "address"},
				{"name": "name", "type": "string"},
				{"name": "amount", "type": "uint256"}
			]
		}],
		"outputs": []
	}]`)

	// Test without ExternalTuples option (should generate tuple)
	code := mustGenerate(t, abiDef)
	assertContains(t, code, "type ProcessUserDataCall struct")

	// Test with ExternalTuples option (should not generate tuple)
	// The tuple name is generated from the tuple structure, not the function name
	extTuples := map[string]string{
		"Tupleb53c1574": "MyCustomUserData",
	}

	codeWithExternal := mustGenerate(t, abiDef, ExternalTuples(extTuples))
	assertNotContains(t, codeWithExternal, "type Tuple_b53c1574 struct")
	assertContains(t, codeWithExternal, "MyCustomUserData", "Data MyCustomUserData")
}
