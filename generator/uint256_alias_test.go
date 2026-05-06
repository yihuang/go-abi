package generator

import (
	"bytes"
	"strings"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

// TestStdlibUint256NoConcreteRef ensures stdlib uint256 output uses the
// *Uint256 alias and never references holiman/uint256 or uint256.Int directly.
func TestStdlibUint256NoConcreteRef(t *testing.T) {
	abiJSON := `[
		{"name":"u256","type":"function","inputs":[{"name":"v","type":"uint256"}],"outputs":[]},
		{"name":"u72","type":"function","inputs":[{"name":"v","type":"uint72"}],"outputs":[]}
	]`
	abiDef, err := ethabi.JSON(bytes.NewReader([]byte(abiJSON)))
	if err != nil {
		t.Fatal(err)
	}
	code, err := NewGenerator(Stdlib(true), UseUint256(true)).GenerateFromABI(abiDef)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(code, `"github.com/holiman/uint256"`) {
		t.Error("stdlib uint256 output unexpectedly imports github.com/holiman/uint256")
	}
	if strings.Contains(code, "uint256.Int") {
		t.Error("stdlib uint256 output unexpectedly names uint256.Int (should use Uint256 alias)")
	}
	if !strings.Contains(code, "*Uint256") {
		t.Error("stdlib uint256 output does not reference the Uint256 alias")
	}
}
