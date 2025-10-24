package tests

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// Benchmark data setup functions
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
		AddressMatrix: [][]common.Address{
			{
				common.HexToAddress("0x1111111111111111111111111111111111111111"),
				common.HexToAddress("0x2222222222222222222222222222222222222222"),
				common.HexToAddress("0x3333333333333333333333333333333333333333"),
			},
			{
				common.HexToAddress("0x4444444444444444444444444444444444444444"),
				common.HexToAddress("0x5555555555555555555555555555555555555555"),
			},
			{
				common.HexToAddress("0x6666666666666666666666666666666666666666"),
			},
		},
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

// Benchmark functions for go-abi generated code
func BenchmarkGoABI_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := args.EncodeWithSelector()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_NestedDynamicArrays(b *testing.B) {
	args := createNestedDynamicArraysData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := args.EncodeWithSelector()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_MixedTypes(b *testing.B) {
	args := createMixedTypesData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := args.EncodeWithSelector()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark functions for go-ethereum ABI encoding
func BenchmarkGoEthereum_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ComprehensiveTestABIDef.Pack("testComplexDynamicTuples", args.Users)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoEthereum_NestedDynamicArrays(b *testing.B) {
	args := createNestedDynamicArraysData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ComprehensiveTestABIDef.Pack("testNestedDynamicArrays",
			args.Matrix, args.AddressMatrix)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoEthereum_MixedTypes(b *testing.B) {
	args := createMixedTypesData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ComprehensiveTestABIDef.Pack("testMixedTypes",
			args.FixedData, args.DynamicData, args.Flag, args.Count, args.Items)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Additional benchmarks for different encoding methods
func BenchmarkGoABI_EncodeOnly_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := args.Encode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_EncodeTo_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	size := args.EncodedSize()
	buf := make([]byte, size)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := args.EncodeTo(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Memory allocation benchmarks
func BenchmarkGoABI_MemoryAllocations_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := args.EncodeWithSelector()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoEthereum_MemoryAllocations_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ComprehensiveTestABIDef.Pack("testComplexDynamicTuples", args.Users)
		if err != nil {
			b.Fatal(err)
		}
	}
}
