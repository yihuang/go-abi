package abi

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
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
		name     string
		data     string
		signed   bool
		expected *big.Int
		expErr   bool
		err      error
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
			name:   "insufficient data",
			data:   "80000000000000000000000000000000000000000000000000000000000000",
			signed: true,
			expErr: true,
			err:    io.ErrUnexpectedEOF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data []byte
			var err error

			data, err = hex.DecodeString(tt.data)
			require.NoError(t, err)

			result, err := DecodeBigInt(data, tt.signed)

			if tt.expErr {
				require.Error(t, err)
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)

				// roundtrip
				buf := make([]byte, 32)
				err = EncodeBigInt(result, buf, tt.signed)
				require.NoError(t, err)
				require.Equal(t, tt.data, hex.EncodeToString(buf))
			}
		})
	}
}

// TestFixedBytesEncodePadding verifies that fixed bytes encode functions
// properly zero out the padding region of the 32-byte ABI slot.
// Without this, reusing a buffer via EncodeTo would leave garbage in padding bytes.
func TestFixedBytesEncodePadding(t *testing.T) {
	tests := []struct {
		name   string
		size   int
		encode func(value []byte, buf []byte) (int, error)
	}{
		{"bytes1", 1, func(value []byte, buf []byte) (int, error) {
			v := [1]byte(value)
			return EncodeBytes1(v, buf)
		}},
		{"bytes3", 3, func(value []byte, buf []byte) (int, error) {
			v := [3]byte(value)
			return EncodeBytes3(v, buf)
		}},
		{"bytes7", 7, func(value []byte, buf []byte) (int, error) {
			v := [7]byte(value)
			return EncodeBytes7(v, buf)
		}},
		{"bytes15", 15, func(value []byte, buf []byte) (int, error) {
			v := [15]byte(value)
			return EncodeBytes15(v, buf)
		}},
		{"bytes31", 31, func(value []byte, buf []byte) (int, error) {
			v := [31]byte(value)
			return EncodeBytes31(v, buf)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill buffer with non-zero poison bytes to simulate a reused buffer
			buf := bytes.Repeat([]byte{0xFF}, 32)

			value := make([]byte, tt.size)
			for i := range value {
				value[i] = byte(i + 1)
			}

			n, err := tt.encode(value, buf)
			require.NoError(t, err)
			require.Equal(t, 32, n,
				"fixed-bytes encoder should report writing a full 32-byte word, got %d", n)

			// Verify: bytes from tt.size to 32 should be zero (padding)
			for i := tt.size; i < 32; i++ {
				require.Equal(t, byte(0x00), buf[i],
					"padding byte at position %d should be zero, got 0x%02x", i, buf[i])
			}

			// Verify: first tt.size bytes should be the value
			for i := 0; i < tt.size; i++ {
				require.Equal(t, byte(i+1), buf[i],
					"value byte at position %d mismatch", i)
			}
		})
	}
}

// TestFixedBytesDecodePadding verifies that fixed-bytes decode functions
// correctly reject dirty padding bytes, accept clean padding, and consume
// the full 32-byte ABI word.
func TestFixedBytesDecodePadding(t *testing.T) {
	tests := []struct {
		name   string
		decode func(data []byte) (interface{}, int, error)
		size   int
	}{
		{"bytes1", func(data []byte) (interface{}, int, error) { return DecodeBytes1(data) }, 1},
		{"bytes3", func(data []byte) (interface{}, int, error) { return DecodeBytes3(data) }, 3},
		{"bytes7", func(data []byte) (interface{}, int, error) { return DecodeBytes7(data) }, 7},
		{"bytes31", func(data []byte) (interface{}, int, error) { return DecodeBytes31(data) }, 31},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/clean", func(t *testing.T) {
			// Clean 32-byte input with proper zero padding
			// Only set bytes within the value region to avoid false dirty padding
			data := make([]byte, 32)
			for i := 0; i < tt.size && i < 5; i++ {
				data[i] = byte(i + 1)
			}
			_, n, err := tt.decode(data)
			require.NoError(t, err)
			// Fixed-size byte decoders should consume the full 32-byte ABI word
			require.Equal(t, 32, n,
				"Decode%s should return 32 bytes consumed (ABI word), got %d", tt.name, n)
		})

		t.Run(tt.name+"/dirty", func(t *testing.T) {
			// Dirty padding: non-zero byte just after the value
			data := make([]byte, 32)
			data[tt.size] = 0x01 // dirty padding
			_, _, err := tt.decode(data)
			require.True(t, errors.Is(err, ErrDirtyPadding))
		})
	}
}

// TestDecodeBigIntNegative tests decoding and encoding of negative int256 values (roundtrip)
func TestDecodeBigIntNegative(t *testing.T) {
	negativeValues := []*big.Int{
		big.NewInt(-1),
		big.NewInt(-128),
		big.NewInt(-32768),
		big.NewInt(-2147483648),
		big.NewInt(-9223372036854775808),
		func() *big.Int {
			v := new(big.Int).Lsh(common.Big1, 255)
			v.Neg(v)
			return v
		}(), // -2^255 (min int256)
	}

	for _, val := range negativeValues {
		t.Run(val.String(), func(t *testing.T) {
			buf := make([]byte, 32)
			err := EncodeBigInt(val, buf, true)
			require.NoError(t, err)

			result, err := DecodeBigInt(buf, true)
			require.NoError(t, err)
			require.Equal(t, val, result)
		})
	}

	t.Run("negative_unsigned_rejected", func(t *testing.T) {
		// Encoding negative value as unsigned should fail
		buf := make([]byte, 32)
		err := EncodeBigInt(big.NewInt(-1), buf, false)
		require.True(t, errors.Is(err, ErrNegativeValue))
	})
}

// TestDecodeDynamicMalformedLength ensures that DecodeString and DecodeBytes
// both return a clean error (not a runtime panic) when the encoded length overflows Pad32.
// TestDecodeUintPadding verifies that DecodeUint rejects dirty padding in high bytes.
func TestDecodeUintPadding(t *testing.T) {
	// uint8 max = 255, a valid encoding is 0x0000...00FF
	// If the high bytes have non-zero bits for an unsigned value, reject.
	data := make([]byte, 32)
	data[31] = 0xAB // valid uint8 value 0xAB
	data[30] = 0x01 // dirty padding byte

	_, err := DecodeUint[uint8](data, MaxUint8)
	require.True(t, errors.Is(err, ErrDirtyPadding))

	// Also test that a clean value passes
	data2 := make([]byte, 32)
	data2[31] = 0xAB
	v, err := DecodeUint[uint8](data2, MaxUint8)
	require.NoError(t, err)
	require.Equal(t, uint8(0xAB), v)
}

func TestDecodeIntSignExtension(t *testing.T) {
	// int8: valid negative value -1 = 0xFF...FF
	// Upper bytes should be all 1s for negative values
	t.Run("negative_with_partial_extension_rejected", func(t *testing.T) {
		// Only lower 32 bits are 0xFFFFFFFF, upper 96 bits are zero
		// This is NOT proper sign extension for int8
		data := make([]byte, 32)
		data[28] = 0xFF
		data[29] = 0xFF
		data[30] = 0xFF
		data[31] = 0xFF // 32-bit sign extension for -1, but not 256-bit

		_, err := DecodeInt[int8](data, MinInt8, MaxInt8)
		require.True(t, errors.Is(err, ErrDirtyPadding))
	})

	t.Run("positive_with_garbage_rejected", func(t *testing.T) {
		// Positive value with garbage in upper bytes
		data := make([]byte, 32)
		data[31] = 0x05 // value = 5
		data[30] = 0x01 // garbage in padding

		_, err := DecodeInt[int8](data, MinInt8, MaxInt8)
		require.True(t, errors.Is(err, ErrDirtyPadding))
	})

	t.Run("proper_sign_extension_accepted", func(t *testing.T) {
		data := make([]byte, 32)
		for i := 0; i < 31; i++ {
			data[i] = 0xFF
		}
		data[31] = 0x80 // int8 value -128 with proper 256-bit sign extension

		v, err := DecodeInt[int8](data, MinInt8, MaxInt8)
		require.NoError(t, err)
		require.Equal(t, int8(-128), v)
	})
}

func TestDecodeBoolValidation(t *testing.T) {
	// Valid values: 0 and 1
	t.Run("valid_true", func(t *testing.T) {
		data := make([]byte, 32)
		data[31] = 0x01
		v, n, err := DecodeBool(data)
		require.NoError(t, err)
		require.True(t, v)
		require.Equal(t, 32, n)
	})

	t.Run("valid_false", func(t *testing.T) {
		data := make([]byte, 32)
		v, n, err := DecodeBool(data)
		require.NoError(t, err)
		require.False(t, v)
		require.Equal(t, 32, n)
	})

	t.Run("invalid_2", func(t *testing.T) {
		data := make([]byte, 32)
		data[31] = 0x02
		_, _, err := DecodeBool(data)
		require.True(t, errors.Is(err, ErrDirtyPadding))
	})

	t.Run("invalid_0xFF", func(t *testing.T) {
		data := make([]byte, 32)
		data[31] = 0xFF
		_, _, err := DecodeBool(data)
		require.True(t, errors.Is(err, ErrDirtyPadding))
	})

	t.Run("dirty_padding_high_bytes", func(t *testing.T) {
		data := make([]byte, 32)
		data[31] = 0x01
		data[30] = 0x01 // dirty padding
		_, _, err := DecodeBool(data)
		require.True(t, errors.Is(err, ErrDirtyPadding))
	})
}

func TestDecodeAddressValidation(t *testing.T) {
	t.Run("valid_address", func(t *testing.T) {
		data := make([]byte, 32)
		data[12] = 0x12
		data[31] = 0x34
		addr, n, err := DecodeAddress(data)
		require.NoError(t, err)
		require.Equal(t, 32, n)
		require.Equal(t, common.HexToAddress("0x1200000000000000000000000000000000000034"), addr)
	})

	t.Run("dirty_padding_first_12_bytes", func(t *testing.T) {
		data := make([]byte, 32)
		data[11] = 0x01 // dirty padding in the first 12 bytes
		data[12] = 0x12
		data[31] = 0x34
		_, _, err := DecodeAddress(data)
		require.True(t, errors.Is(err, ErrDirtyPadding))
	})

	t.Run("insufficient_data", func(t *testing.T) {
		// Must not panic on short data
		_, _, err := DecodeAddress([]byte{0x01, 0x02})
		require.True(t, errors.Is(err, io.ErrUnexpectedEOF),
			"DecodeAddress should return io.ErrUnexpectedEOF for short data, not panic")
	})

	t.Run("insufficient_data_31_bytes", func(t *testing.T) {
		// Exactly 31 bytes should also fail (needs 32)
		_, _, err := DecodeAddress(make([]byte, 31))
		require.True(t, errors.Is(err, io.ErrUnexpectedEOF))
	})
}

func TestDecodeDynamicMalformedLength(t *testing.T) {
	// 9223372036854775787 = 0x7FFFFFFFFFFFFFEB
	poisonLen := uint64(9223372036854775787)
	data := make([]byte, 32)
	for i := 0; i < 8; i++ {
		data[31-i] = byte(poisonLen >> (i * 8))
	}

	tests := []struct {
		name   string
		decode func([]byte) (int, error)
	}{
		{
			name: "DecodeString",
			decode: func(d []byte) (int, error) {
				_, n, err := DecodeString(d)
				return n, err
			},
		},
		{
			name: "DecodeBytes",
			decode: func(d []byte) (int, error) {
				_, n, err := DecodeBytes(d)
				return n, err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.decode(data)
			require.Error(t, err)
			require.True(t, errors.Is(err, io.ErrUnexpectedEOF))
		})
	}
}

func TestWriteUint256ShortBuffer(t *testing.T) {
	val := new(Uint256)
	short := make([]byte, 16)
	err := WriteUint256(val, short)
	require.Error(t, err)
	require.True(t, errors.Is(err, io.ErrShortBuffer))

	buf := make([]byte, 32)
	require.NoError(t, WriteUint256(val, buf))
}
