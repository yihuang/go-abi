package testdata

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

func TestDirectTupleNestingInHumanReadableABI(t *testing.T) {
	// Test direct tuple nesting in human readable ABI (without struct definitions)
	humanABI := []string{
		"function testSimpleTuple((uint256 a, uint256 b) pair)",
		"function testNestedTuple((uint256 a, (address b, string c) b) complex)",
		"function testTupleArray((uint256 a, uint256 b)[] pairs)",
		"function testMixed((uint256 a, address b) tuple1, (string a, bytes b) tuple2)",
		"function testDeepNested(((uint256 a, address b) a, (string a, bytes b) b) deep)",
	}

	// Parse human-readable ABI
	abiJSON, err := abi.ParseHumanReadableABI(humanABI)
	require.NoError(t, err, "Failed to parse human-readable ABI with direct tuple nesting")

	// Parse JSON ABI
	abiDef, err := ethabi.JSON(bytes.NewReader(abiJSON))
	require.NoError(t, err)

	// Verify we have the expected functions
	expectedFunctions := []string{"testSimpleTuple", "testNestedTuple", "testTupleArray", "testMixed", "testDeepNested"}
	for _, expectedFunc := range expectedFunctions {
		_, found := abiDef.Methods[expectedFunc]
		require.True(t, found, "Expected function %s not found in ABI", expectedFunc)
	}

	// Generate Go code from the ABI
	generator := abi.NewGenerator("testdata")
	generatedCode, err := generator.GenerateFromABI(abiDef)
	require.NoError(t, err)

	// Basic sanity check on generated code
	require.Contains(t, generatedCode, "type TestSimpleTupleCall struct")
	require.Contains(t, generatedCode, "type TestNestedTupleCall struct")
	require.Contains(t, generatedCode, "type TestTupleArrayCall struct")

	t.Log("✅ Direct tuple nesting in human readable ABI is working correctly!")
}
