//go:build uint256

package tests

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/test-go/testify/require"
)

// TestLazyViewUint256Field exercises the lazy view path for a uint256 field
// (mapped to *uint256.Int under the uint256 build tag). Closes the
// (uint256, lazy) gap left by lazy_test.go (gated !uint256) and
// comprehensive_uint256.abi.go (compiled but not exercised at runtime).
func TestLazyViewUint256Field(t *testing.T) {
	cases := map[string]*uint256.Int{
		"small":      uint256.NewInt(123456789),
		"max-uint64": uint256.NewInt(^uint64(0)),
		"large":      new(uint256.Int).Lsh(uint256.NewInt(1), 200),
	}
	for name, val := range cases {
		t.Run(name, func(t *testing.T) {
			ret := &GetRecordsReturn{Field1: val}
			encoded, err := ret.Encode()
			require.NoError(t, err)

			view, n, err := DecodeGetRecordsReturnView(encoded)
			require.NoError(t, err)
			require.Equal(t, len(encoded), n)

			got, err := view.Field1()
			require.NoError(t, err)
			require.True(t, val.Eq(got), "Field1 got %s, want %s", got, val)

			mat, err := view.Materialize()
			require.NoError(t, err)
			require.True(t, val.Eq(mat.Field1), "Materialize Field1 got %s, want %s", mat.Field1, val)
		})
	}
}
