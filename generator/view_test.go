package generator

import (
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// mustParseABI parses an ABI JSON string and returns the ABI definition.
// It calls t.Fatal if parsing fails.
func mustParseABI(t *testing.T, abiJSON string) abi.ABI {
	t.Helper()
	abiDef, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		t.Fatalf("Failed to parse ABI: %v", err)
	}
	return abiDef
}

// mustGenerate generates code from an ABI with the given options.
// It calls t.Fatal if generation fails.
func mustGenerate(t *testing.T, abiDef abi.ABI, opts ...Option) string {
	t.Helper()
	g := NewGenerator(opts...)
	code, err := g.GenerateFromABI(abiDef)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}
	return code
}

// assertContains checks that code contains all the given substrings.
func assertContains(t *testing.T, code string, substrs ...string) {
	t.Helper()
	for _, s := range substrs {
		if !strings.Contains(code, s) {
			t.Errorf("Expected code to contain %q", s)
		}
	}
}

// assertNotContains checks that code does not contain any of the given substrings.
func assertNotContains(t *testing.T, code string, substrs ...string) {
	t.Helper()
	for _, s := range substrs {
		if strings.Contains(code, s) {
			t.Errorf("Expected code to NOT contain %q", s)
		}
	}
}

func TestZeroValue(t *testing.T) {
	tests := []struct {
		goType   string
		expected string
	}{
		{"*big.Int", "nil"},
		{"*uint256.Int", "nil"},
		{"[]byte", "nil"},
		{"[]string", "nil"},
		{"[5]byte", "[5]byte{}"},
		{"[32]byte", "[32]byte{}"},
		{"bool", "false"},
		{"string", `""`},
		{"uint8", "0"},
		{"uint64", "0"},
		{"uint256", "0"},
		{"int8", "0"},
		{"int64", "0"},
		{"common.Address", "common.Address{}"},
		{"MyStruct", "MyStruct{}"},
	}

	for _, tt := range tests {
		t.Run(tt.goType, func(t *testing.T) {
			result := zeroValue(tt.goType)
			if result != tt.expected {
				t.Errorf("zeroValue(%q) = %q, want %q", tt.goType, result, tt.expected)
			}
		})
	}
}

func TestCountDynamicFields(t *testing.T) {
	t.Run("MixedFields", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "testMixed",
			"inputs": [
				{"name": "staticField", "type": "uint256"},
				{"name": "dynamicField", "type": "string"},
				{"name": "staticArray", "type": "uint256[3]"},
				{"name": "dynamicArray", "type": "uint256[]"}
			],
			"outputs": []
		}]`)

		method := abiDef.Methods["testMixed"]
		s := StructFromArguments("TestMixedCall", method.Inputs)

		count := countDynamicFields(s)
		if count != 2 {
			t.Errorf("countDynamicFields = %d, want 2 (string and dynamic array)", count)
		}
	})

	t.Run("AllStatic", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "testStatic",
			"inputs": [
				{"name": "a", "type": "uint256"},
				{"name": "b", "type": "address"},
				{"name": "c", "type": "bool"}
			],
			"outputs": []
		}]`)

		method := abiDef.Methods["testStatic"]
		s := StructFromArguments("TestStaticCall", method.Inputs)

		count := countDynamicFields(s)
		if count != 0 {
			t.Errorf("countDynamicFields = %d, want 0", count)
		}
	})
}

func TestCollectSliceTypes(t *testing.T) {
	abiDef := mustParseABI(t, `[{
		"type": "function",
		"name": "testSlices",
		"inputs": [
			{"name": "strings", "type": "string[]"},
			{"name": "uints", "type": "uint256[]"},
			{"name": "nested", "type": "uint256[][]"}
		],
		"outputs": [
			{"name": "addresses", "type": "address[]"}
		]
	}]`)

	sliceTypes := collectSliceTypes(abiDef)

	// Should find: string[], uint256[], uint256[][], address[]
	if len(sliceTypes) < 4 {
		t.Errorf("collectSliceTypes found %d types, want at least 4", len(sliceTypes))
	}

	// Verify all are slice types
	for _, st := range sliceTypes {
		if st.T != abi.SliceTy {
			t.Errorf("collectSliceTypes returned non-slice type: %v", st.T)
		}
	}
}

func TestExternalTupleReferences(t *testing.T) {
	t.Run("TypeReferences", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "testTuple",
			"inputs": [{
				"name": "data",
				"type": "tuple",
				"components": [
					{"name": "name", "type": "string"},
					{"name": "value", "type": "uint256"}
				]
			}],
			"outputs": []
		}]`)

		method := abiDef.Methods["testTuple"]
		tupleType := method.Inputs[0].Type

		// Without external tuples
		if typeReferencesExternalTuple(tupleType, nil) {
			t.Error("Expected false for empty external tuples map")
		}

		if typeReferencesExternalTuple(tupleType, map[string]string{}) {
			t.Error("Expected false for empty external tuples map")
		}

		// With matching external tuple
		s := StructFromTuple(tupleType)
		externalTuples := map[string]string{s.Name: "ExternalType"}

		if !typeReferencesExternalTuple(tupleType, externalTuples) {
			t.Errorf("Expected true when tuple %s is in external tuples", s.Name)
		}
	})

	t.Run("StructReferences", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "test",
			"inputs": [{
				"name": "container",
				"type": "tuple",
				"components": [
					{"name": "inner", "type": "tuple", "components": [{"name": "x", "type": "uint256"}]},
					{"name": "value", "type": "uint256"}
				]
			}],
			"outputs": []
		}]`)

		method := abiDef.Methods["test"]
		containerType := method.Inputs[0].Type
		containerStruct := StructFromTuple(containerType)

		// Get inner tuple name
		innerType := containerType.TupleElems[0]
		innerStruct := StructFromTuple(*innerType)

		// Without external tuples
		if structReferencesExternalTuple(containerStruct, nil) {
			t.Error("Expected false for nil external tuples")
		}

		// With inner tuple as external
		externalTuples := map[string]string{innerStruct.Name: "ExternalInner"}
		if !structReferencesExternalTuple(containerStruct, externalTuples) {
			t.Error("Expected true when nested tuple is external")
		}
	})
}

func TestSliceViewTypeName(t *testing.T) {
	abiDef := mustParseABI(t, `[{
		"type": "function",
		"name": "test",
		"inputs": [{"name": "items", "type": "uint256[]"}],
		"outputs": []
	}]`)

	method := abiDef.Methods["test"]
	sliceType := method.Inputs[0].Type

	name := sliceViewTypeName(sliceType)
	if !strings.HasSuffix(name, "View") {
		t.Errorf("sliceViewTypeName = %q, want suffix 'View'", name)
	}
}

func TestGenerateLazyViews(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "getProfile",
			"inputs": [{
				"name": "profile",
				"type": "tuple",
				"components": [
					{"name": "name", "type": "string"},
					{"name": "age", "type": "uint64"}
				]
			}],
			"outputs": [{"name": "success", "type": "bool"}]
		}]`)

		code := mustGenerate(t, abiDef, GenerateLazy(true))
		assertContains(t, code,
			"View struct {",
			"View(data []byte)",
			"Materialize()",
			"Raw() []byte",
		)
	})

	t.Run("Disabled", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "test",
			"inputs": [{"name": "value", "type": "uint256"}],
			"outputs": []
		}]`)

		code := mustGenerate(t, abiDef) // Without GenerateLazy option
		assertNotContains(t, code, "DecodeTestCallView")
	})

	t.Run("WithSlices", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "getItems",
			"inputs": [{
				"name": "items",
				"type": "tuple[]",
				"components": [
					{"name": "id", "type": "uint32"},
					{"name": "data", "type": "bytes"}
				]
			}],
			"outputs": []
		}]`)

		code := mustGenerate(t, abiDef, GenerateLazy(true))
		assertContains(t, code, "Len() int", "Get(i int)")
	})

	t.Run("WithExternalTuples", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "process",
			"inputs": [{
				"name": "data",
				"type": "tuple",
				"components": [{"name": "value", "type": "uint256"}]
			}],
			"outputs": []
		}]`)

		method := abiDef.Methods["process"]
		s := StructFromTuple(method.Inputs[0].Type)

		code := mustGenerate(t, abiDef,
			GenerateLazy(true),
			ExternalTuples(map[string]string{s.Name: "ExternalData"}),
		)
		assertNotContains(t, code, s.Name+"View struct")
	})

	t.Run("AllStaticTuple", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "getPoint",
			"inputs": [{
				"name": "point",
				"type": "tuple",
				"components": [
					{"name": "x", "type": "uint256"},
					{"name": "y", "type": "uint256"},
					{"name": "z", "type": "uint256"}
				]
			}],
			"outputs": []
		}]`)

		code := mustGenerate(t, abiDef, GenerateLazy(true))
		assertContains(t, code, "data []byte")
		assertNotContains(t, code, "offsets [0]") // All-static tuples don't need offsets
	})

	t.Run("NestedDynamic", func(t *testing.T) {
		abiDef := mustParseABI(t, `[{
			"type": "function",
			"name": "getUser",
			"inputs": [{
				"name": "user",
				"type": "tuple",
				"components": [
					{"name": "name", "type": "string"},
					{"name": "emails", "type": "string[]"},
					{
						"name": "profile",
						"type": "tuple",
						"components": [
							{"name": "bio", "type": "string"},
							{"name": "age", "type": "uint64"}
						]
					}
				]
			}],
			"outputs": []
		}]`)

		code := mustGenerate(t, abiDef, GenerateLazy(true))

		viewCount := strings.Count(code, "View struct {")
		if viewCount < 2 {
			t.Errorf("Expected at least 2 View structs for nested types, got %d", viewCount)
		}

		assertContains(t, code, "Profile()") // Getter for nested tuple field
	})
}
