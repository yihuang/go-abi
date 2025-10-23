package testdata

import (
	"bytes"
	"strings"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

func TestDirectTupleNestingInHumanReadableABI(t *testing.T) {
	// Test direct tuple nesting in human readable ABI (without struct definitions)
	humanABI := []string{
		"function testSimpleTuple((uint256, uint256) pair)",
		"function testNestedTuple((uint256, (address, string)) complex)",
		"function testTupleArray((uint256, uint256)[] pairs)",
		"function testMixed((uint256, address) tuple1, (string, bytes) tuple2)",
		"function testDeepNested(((uint256, address), (string, bytes)) deep)",
	}

	// Parse human-readable ABI
	abiJSON, err := abi.ParseHumanReadableABI(humanABI)
	if err != nil {
		t.Fatalf("Failed to parse human-readable ABI with direct tuple nesting: %v", err)
	}

	// Parse JSON ABI
	abiDef, err := ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		t.Fatalf("Failed to parse JSON ABI: %v", err)
	}

	// Verify we have the expected functions
	expectedFunctions := []string{"testSimpleTuple", "testNestedTuple", "testTupleArray", "testMixed", "testDeepNested"}
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

	// Generate Go code from the ABI
	generator := abi.NewGenerator("testdata")
	generatedCode, err := generator.GenerateFromABI(abiDef)
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Basic sanity check on generated code
	if !strings.Contains(generatedCode, "type TestSimpleTupleCall struct") {
		t.Error("Generated code should contain TestSimpleTupleCall struct")
	}
	if !strings.Contains(generatedCode, "type TestNestedTupleCall struct") {
		t.Error("Generated code should contain TestNestedTupleCall struct")
	}
	if !strings.Contains(generatedCode, "type TestTupleArrayCall struct") {
		t.Error("Generated code should contain TestTupleArrayCall struct")
	}

	t.Log("Generated code for direct tuple nesting:")
	// Print relevant parts of generated code
	lines := strings.Split(generatedCode, "\n")
	for i, line := range lines {
		if strings.Contains(line, "type") && strings.Contains(line, "struct") {
			t.Logf("  %s", line)
			// Print the next few lines to see struct fields
			for j := i + 1; j < len(lines) && j < i+10; j++ {
				if strings.Contains(lines[j], "}") {
					break
				}
				t.Logf("    %s", lines[j])
			}
		}
	}

	t.Log("✅ Direct tuple nesting in human readable ABI is working correctly!")
}