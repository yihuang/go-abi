package tests

import (
	"errors"
	"io"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/test-go/testify/require"
	"github.com/yihuang/go-abi"
)

func DecodeRoundTrip[T any, PT interface {
	abi.Tuple
	*T
}](t *testing.T, orig PT) {
	data, err := orig.Encode()
	require.NoError(t, err)

	var decoded T
	_, err = PT(&decoded).Decode(data)
	require.NoError(t, err)

	require.Equal(t, orig, &decoded)

	// test ErrUnexpectedEOF
	for i := 0; i < len(data); i++ {
		_, err = PT(&decoded).Decode(data[:i])
		require.Error(t, err)
		require.True(t, errors.Is(err, io.ErrUnexpectedEOF))
	}

	// test validation with bit flipping
	// Test a selection of bit positions to verify validation is working
	if len(data) > 0 {
		// Test diverse positions across the entire data
		for pos := 0; pos < len(data)*8; pos++ {
			flipped := make([]byte, len(data))
			copy(flipped, data)
			bitIndex := pos / 8
			bitOffset := pos % 8
			flipped[bitIndex] ^= 1 << bitOffset

			var flippedDecoded T
			_, err := PT(&flippedDecoded).Decode(flipped)

			// it either cause error or unequal result
			if err == nil {
				require.NotEqual(t, orig, &flippedDecoded, "orig: %v, flipped at bit %d", orig, pos)
			}
		}
	}
}

// CompareDecoder compares decode error conditions between go-abi and go-ethereum.
// It tests that both libraries validate inputs the same way by toggling each bit
// in the input and verifying that both decoders either succeed or fail together.
func CompareDecoder(t *testing.T, tuple abi.Tuple, abiSpec ethabi.ABI, method string, input []byte) {
	// For each bit position, toggle that bit and test both decoders
	bits := len(input) * 8
	for i := 0; i < bits; i++ {
		// Clone input and toggle bit i
		In := make([]byte, len(input))
		copy(In, input)
		inIndex := i / 8
		bitIndex := i % 8
		In[inIndex] ^= 1 << bitIndex

		// Test go-abi decode
		_, err1 := tuple.Decode(In)

		// Test go-ethereum unpack
		_, err2 := abiSpec.Unpack(method, In)

		// Assert that both either succeed or fail together
		// err1 and err2 should both be nil or both be non-nil
		hasErr1 := err1 != nil
		hasErr2 := err2 != nil

		if hasErr1 != hasErr2 {
			t.Errorf("Bit %d: error mismatch - go-abi: %v, go-ethereum: %v", i, err1, err2)
		}
	}
}

func EventDecodeRoundTrip[T any, PT interface {
	abi.Event
	*T
}](t *testing.T, orig PT) {
	topics, data, err := abi.EncodeEvent(orig)
	require.NoError(t, err)

	var decoded T
	err = abi.DecodeEvent(PT(&decoded), topics, data)
	require.NoError(t, err)

	require.Equal(t, orig, &decoded)

	// test ErrUnexpectedEOF for data
	for i := 0; i < len(data); i++ {
		err = abi.DecodeEvent(PT(&decoded), topics, data[:i])
		require.Error(t, err)
		require.True(t, errors.Is(err, io.ErrUnexpectedEOF))
	}
}
