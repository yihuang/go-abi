package tests

import (
	"errors"
	"io"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"github.com/yihuang/go-abi"
)

var TestAddress = common.HexToAddress("0x1234567890123456789012345678901234567890")

func BenchEncode(b *testing.B, call abi.Tuple) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := call.Encode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchEncodeTo(b *testing.B, call abi.Tuple) {
	buf := make([]byte, call.EncodedSize())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := call.EncodeTo(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchDecode(b *testing.B, encoded []byte, newCall func() abi.Decode) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		call := newCall()
		_, err := call.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

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
