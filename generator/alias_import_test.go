package generator

import (
	"bytes"
	"strings"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

func TestAliasImports(t *testing.T) {
	// Test with alias imports
	g := NewGenerator(
		ExtraImports([]ImportSpec{
			{
				Path:  "github.com/ethereum/go-ethereum/common",
				Alias: "cmn",
			},
			{
				Path:  "math",
				Alias: "math",
			},
			{
				Path: "time",
			},
		}),
	)

	// Create a simple ABI for testing
	abiJSON := `[
		{
			"name": "test",
			"type": "function",
			"inputs": [{"name": "value", "type": "uint256"}],
			"outputs": []
		}
	]`

	// Parse ABI and generate code
	abiDef, err := ethabi.JSON(bytes.NewReader([]byte(abiJSON)))
	if err != nil {
		t.Fatal(err)
	}

	code, err := g.GenerateFromABI(abiDef)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that alias imports are properly formatted
	if !strings.Contains(code, `cmn "github.com/ethereum/go-ethereum/common"`) {
		t.Error("Expected alias import 'cmn \"github.com/ethereum/go-ethereum/common\"' not found")
	}
	if !strings.Contains(code, `math "math"`) {
		t.Error("Expected alias import 'math \"math\"' not found")
	}
	if !strings.Contains(code, `"time"`) {
		t.Error("Expected regular import '\"time\"' not found")
	}

	// Verify that the import block contains the alias imports
	if !strings.Contains(code, `cmn "github.com/ethereum/go-ethereum/common"`) {
		t.Error("Expected alias import 'cmn \"github.com/ethereum/go-ethereum/common\"' not found")
	}
	if !strings.Contains(code, `math "math"`) {
		t.Error("Expected alias import 'math \"math\"' not found")
	}
	if !strings.Contains(code, `"time"`) {
		t.Error("Expected regular import '\"time\"' not found")
	}
}
