package generator

import (
	"bytes"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestGenerator2SimpleFunction(t *testing.T) {
	// Test with a simple ABI: transfer(address,uint256)
	abiJSON := `[
		{
			"type": "function",
			"name": "transfer",
			"inputs": [
				{"name": "to", "type": "address"},
				{"name": "amount", "type": "uint256"}
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

	// Check that package declaration is present
	if !bytes.Contains([]byte(code), []byte("package test")) {
		t.Error("Generated code missing package declaration")
	}

	// Check that function selectors are generated
	if !bytes.Contains([]byte(code), []byte("TransferSelector")) {
		t.Error("Generated code missing function selector")
	}

	// Check that encoding functions are generated
	if !bytes.Contains([]byte(code), []byte("encode_address")) {
		t.Error("Generated code missing address encoding function")
	}

	if !bytes.Contains([]byte(code), []byte("encode_uint256")) {
		t.Error("Generated code missing uint256 encoding function")
	}

	t.Logf("Generated code length: %d bytes", len(code))
	t.Logf("Generated code preview:\n%s", code[:500])
}