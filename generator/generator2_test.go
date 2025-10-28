package generator

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestGenTypeIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    abi.Type
		expected string
	}{
		{
			name:     "uint256",
			input:    mustParseType("uint256"),
			expected: "uint256",
		},
		{
			name:     "int256",
			input:    mustParseType("int256"),
			expected: "int256",
		},
		{
			name:     "address",
			input:    mustParseType("address"),
			expected: "address",
		},
		{
			name:     "bool",
			input:    mustParseType("bool"),
			expected: "bool",
		},
		{
			name:     "string",
			input:    mustParseType("string"),
			expected: "string",
		},
		{
			name:     "bytes",
			input:    mustParseType("bytes"),
			expected: "bytes",
		},
		{
			name:     "bytes32",
			input:    mustParseType("bytes32"),
			expected: "bytes32",
		},
		{
			name:     "uint8",
			input:    mustParseType("uint8"),
			expected: "uint8",
		},
		{
			name:     "int64",
			input:    mustParseType("int64"),
			expected: "int64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenTypeIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("GenTypeIdentifier() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewGenerator2(t *testing.T) {
	generator := NewGenerator2()
	if generator == nil {
		t.Fatal("NewGenerator2() returned nil")
	}
	if generator.Options.PackageName != "abi" {
		t.Errorf("Expected default package name 'abi', got %s", generator.Options.PackageName)
	}
}

func mustParseType(typeStr string) abi.Type {
	t, err := abi.NewType(typeStr, "", nil)
	if err != nil {
		panic(err)
	}
	return t
}