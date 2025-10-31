package abi

import (
	"encoding/hex"
	"io"
	"math/big"
	"testing"

	"github.com/test-go/testify/require"
)

func TestEncodeBigInt(t *testing.T) {
	tests := []struct {
		name     string
		value    *big.Int
		signed   bool
		expected string
	}{
		{
			name:     "signed negative",
			value:    big.NewInt(-100),
			signed:   true,
			expected: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c",
		},
		{
			name:     "unsigned positive",
			value:    big.NewInt(100),
			signed:   false,
			expected: "0000000000000000000000000000000000000000000000000000000000000064",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 32)
			err := EncodeBigInt(tt.value, buf, tt.signed)
			require.NoError(t, err)
			require.Equal(t, tt.expected, hex.EncodeToString(buf))
		})
	}
}

func TestDecodeBigInt(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		signed      bool
		expected    *big.Int
		expectError bool
		error       error
	}{
		{
			name:     "signed positive",
			data:     "0000000000000000000000000000000000000000000000000000000000000064",
			signed:   true,
			expected: big.NewInt(100),
		},
		{
			name:     "signed negative",
			data:     "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff9c",
			signed:   true,
			expected: big.NewInt(-100),
		},
		{
			name:     "unsigned positive",
			data:     "0000000000000000000000000000000000000000000000000000000000000064",
			signed:   false,
			expected: big.NewInt(100),
		},
		{
			name:   "unsigned with high bit set",
			data:   "8000000000000000000000000000000000000000000000000000000000000000",
			signed: false,
			expected: func() *big.Int {
				expected := new(big.Int)
				expected.SetString("8000000000000000000000000000000000000000000000000000000000000000", 16)
				return expected
			}(),
		},
		{
			name:     "zero signed",
			data:     "0000000000000000000000000000000000000000000000000000000000000000",
			signed:   true,
			expected: new(big.Int).SetBytes([]byte{0x00}),
		},
		{
			name:     "zero unsigned",
			data:     "0000000000000000000000000000000000000000000000000000000000000000",
			signed:   false,
			expected: new(big.Int).SetBytes([]byte{0x00}),
		},
		{
			name:   "max int256",
			data:   "7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			signed: true,
			expected: func() *big.Int {
				expected := new(big.Int)
				expected.SetString("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
				return expected
			}(),
		},
		{
			name:   "min int256",
			data:   "8000000000000000000000000000000000000000000000000000000000000000",
			signed: true,
			expected: func() *big.Int {
				expected := new(big.Int)
				expected.SetString("-8000000000000000000000000000000000000000000000000000000000000000", 16)
				return expected
			}(),
		},
		{
			name:        "insufficient data",
			data:        "80000000000000000000000000000000000000000000000000000000000000",
			signed:      true,
			expectError: true,
			error:       io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			var err error

			data, err = hex.DecodeString(tt.data)
			require.NoError(t, err)

			result, err := DecodeBigInt(data, tt.signed)

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, tt.error, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
