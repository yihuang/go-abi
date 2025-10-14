package abi

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/test-go/testify/require"
)

func TestEncodeBigInt(t *testing.T) {
	t.Run("signed", func(t *testing.T) {
		buf := make([]byte, 32)
		err := EncodeBigInt(big.NewInt(-100), buf, true)
		require.NoError(t, err)
		require.Equal(t, "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c", hex.EncodeToString(buf))
	})

	t.Run("unsigned", func(t *testing.T) {
		buf := make([]byte, 32)
		err := EncodeBigInt(big.NewInt(100), buf, false)
		require.NoError(t, err)
		require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000064", hex.EncodeToString(buf))
	})
}
