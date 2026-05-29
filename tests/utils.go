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

var testAddressMatrix = [][3][]common.Address{{
	{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
	},
	{
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
		common.HexToAddress("0x4444444444444444444444444444444444444444"),
		common.HexToAddress("0x5555555555555555555555555555555555555555"),
	},
	{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
	},
}}

var testAddresses5 = [5]common.Address{
	common.HexToAddress("0x1111111111111111111111111111111111111111"),
	common.HexToAddress("0x2222222222222222222222222222222222222222"),
	common.HexToAddress("0x3333333333333333333333333333333333333333"),
	common.HexToAddress("0x4444444444444444444444444444444444444444"),
	common.HexToAddress("0x5555555555555555555555555555555555555555"),
}

var testBytes32s2 = [2][32]byte{
	{0x01, 0x02, 0x03},
	{0x04, 0x05, 0x06},
}

func createTestMatrix[T any](newInt func(int64) T) [][]T {
	return [][]T{
		{newInt(1), newInt(2), newInt(3), newInt(4), newInt(5)},
		{newInt(6), newInt(7), newInt(8)},
		{newInt(9), newInt(10)},
		{newInt(11)},
	}
}

func createTestUints3[T any](newInt func(int64) T) [3]T {
	return [3]T{newInt(100), newInt(200), newInt(300)}
}

var testUserData = []struct {
	Id        int64
	Name      string
	Emails    []string
	CreatedAt int64
	Tags      []string
}{
	{1, "User 1", []string{"user1@example.com", "user1@gmail.com", "user1@test.org"}, 1234567890, []string{"tag1", "tag2", "tag3", "tag4", "tag5"}},
	{2, "User 2 with a longer name for testing", []string{"user2@example.com", "user2@work.com"}, 9876543210, []string{"tag6", "tag7"}},
	{3, "User 3", []string{"user3@example.com"}, 5555555555, []string{"tag8", "tag9", "tag10", "tag11"}},
}

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

func createMixedTypesData() TestMixedTypesCall {
	return TestMixedTypesCall{
		FixedData:   [32]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
		DynamicData: []byte{0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12},
		Flag:        true,
		Count:       uint8(42),
		Items: []Item{
			{
				Id:     uint32(1),
				Data:   []byte{0x13, 0x14, 0x15},
				Active: true,
			},
			{
				Id:     uint32(2),
				Data:   []byte{0x16, 0x17, 0x18, 0x19, 0x1a},
				Active: false,
			},
			{
				Id:     uint32(3),
				Data:   []byte{0x1b, 0x1c},
				Active: true,
			},
		},
	}
}

func createSmallIntegersData() TestSmallIntegersCall {
	return TestSmallIntegersCall{
		U8:  255,
		U16: 65535,
		U24: 16777215,
		U32: 4294967295,
		U64: 18446744073709551615,
		I8:  -128,
		I16: -32768,
		I24: -8388608,
		I32: -2147483648,
		I64: -9223372036854775808,
	}
}

func createFixedBytesData() TestFixedBytesCall {
	return TestFixedBytesCall{
		Data3:  [3]byte{0x01, 0x02, 0x03},
		Data7:  [7]byte{0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a},
		Data15: [15]byte{0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19},
	}
}
