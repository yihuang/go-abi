//go:build !uint256

package tests

import (
	"errors"
	"io"
	"slices"
	"testing"

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
	if len(data) > 0 {
		// Test diverse positions across the entire data
		for pos := 0; pos < len(data)*8; pos++ {
			flipped := slices.Clone(data)
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

func DecodePackedRoundTrip[T any, PT interface {
	abi.PackedTuple
	*T
}](t *testing.T, orig PT) {
	data, err := orig.PackedEncode()
	require.NoError(t, err)

	var decoded T
	_, err = PT(&decoded).PackedDecode(data)
	require.NoError(t, err)

	require.Equal(t, orig, &decoded)

	// test ErrUnexpectedEOF
	for i := range len(data) {
		_, err = PT(&decoded).PackedDecode(data[:i])
		require.Error(t, err)
		require.True(t, errors.Is(err, io.ErrUnexpectedEOF))
	}

	// test validation with bit flipping
	if len(data) > 0 {
		// Test diverse positions across the entire data
		for pos := 0; pos < len(data)*8; pos++ {
			flipped := slices.Clone(data)
			bitIndex := pos / 8
			bitOffset := pos % 8
			flipped[bitIndex] ^= 1 << bitOffset

			var flippedDecoded T
			_, err := PT(&flippedDecoded).PackedDecode(flipped)

			// it either cause error or unequal result
			if err == nil {
				require.NotEqual(t, orig, &flippedDecoded, "orig: %v, flipped at bit %d", orig, pos)
			}
		}
	}
}
