package tests

import (
	"errors"
	"io"
	"math/big"
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

// Benchmark data setup functions - shared across all benchmark files
func createComplexDynamicTuplesData() TestComplexDynamicTuplesCall {
	return TestComplexDynamicTuplesCall{
		Users: []User2{
			{
				Id: big.NewInt(1),
				Profile: UserProfile{
					Name:   "User 1",
					Emails: []string{"user1@example.com", "user1@gmail.com", "user1@test.org"},
					Metadata: UserMetadata2{
						CreatedAt: big.NewInt(1234567890),
						Tags:      []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
					},
				},
			},
			{
				Id: big.NewInt(2),
				Profile: UserProfile{
					Name:   "User 2 with a longer name for testing",
					Emails: []string{"user2@example.com", "user2@work.com"},
					Metadata: UserMetadata2{
						CreatedAt: big.NewInt(9876543210),
						Tags:      []string{"tag6", "tag7"},
					},
				},
			},
			{
				Id: big.NewInt(3),
				Profile: UserProfile{
					Name:   "User 3",
					Emails: []string{"user3@example.com"},
					Metadata: UserMetadata2{
						CreatedAt: big.NewInt(5555555555),
						Tags:      []string{"tag8", "tag9", "tag10", "tag11"},
					},
				},
			},
		},
	}
}

func createNestedDynamicArraysData() TestNestedDynamicArraysCall {
	return TestNestedDynamicArraysCall{
		Matrix: [][]*big.Int{
			{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5)},
			{big.NewInt(6), big.NewInt(7), big.NewInt(8)},
			{big.NewInt(9), big.NewInt(10)},
			{big.NewInt(11)},
		},
		AddressMatrix: [][3][]common.Address{{
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
		}},
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

func createDeeplyNestedData() TestDeeplyNestedCall {
	return TestDeeplyNestedCall{
		Data: Level1{
			Level1: Level2{
				Level2: Level3{
					Level3: Level4{
						Value:       big.NewInt(999),
						Description: "Deeply nested value",
					},
				},
			},
		},
	}
}

func createFixedArraysData() TestFixedArraysCall {
	return TestFixedArraysCall{
		Addresses: [5]common.Address{
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
			common.HexToAddress("0x3333333333333333333333333333333333333333"),
			common.HexToAddress("0x4444444444444444444444444444444444444444"),
			common.HexToAddress("0x5555555555555555555555555555555555555555"),
		},
		Uints: [3]*big.Int{
			big.NewInt(100),
			big.NewInt(200),
			big.NewInt(300),
		},
		Bytes32s: [2][32]byte{
			{0x01, 0x02, 0x03},
			{0x04, 0x05, 0x06},
		},
	}
}

func createFixedBytesData() TestFixedBytesCall {
	return TestFixedBytesCall{
		Data3:  [3]byte{0x01, 0x02, 0x03},
		Data7:  [7]byte{0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a},
		Data15: [15]byte{0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19},
	}
}
