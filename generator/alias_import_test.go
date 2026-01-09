package generator

import (
	"testing"
)

func TestAliasImports(t *testing.T) {
	abiDef := mustParseABI(t, `[{
		"name": "test",
		"type": "function",
		"inputs": [{"name": "value", "type": "uint256"}],
		"outputs": []
	}]`)

	code := mustGenerate(t, abiDef,
		ExtraImports([]ImportSpec{
			{Path: "github.com/ethereum/go-ethereum/common", Alias: "cmn"},
			{Path: "math", Alias: "math"},
			{Path: "time"},
		}),
	)

	// Verify that alias imports are properly formatted
	assertContains(t, code,
		`cmn "github.com/ethereum/go-ethereum/common"`,
		`math "math"`,
		`"time"`,
	)

	// Verify that the import block contains the alias imports
	assertContains(t, code,
		`cmn "github.com/ethereum/go-ethereum/common"`,
		`math "math"`,
		`"time"`,
	)
}
