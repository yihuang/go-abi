//go:build !uint256

package tests

import (
	"math/big"
	"testing"
)

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
			args.Matrix, args.AddressMatrix, args.DymMatrix)
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

func BenchmarkGoABI_Encode_SmallIntegers(b *testing.B) {
	args := createSmallIntegersData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := args.EncodeWithSelector()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoEthereum_Encode_SmallIntegers(b *testing.B) {
	args := createSmallIntegersData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ComprehensiveTestABIDef.Pack("testSmallIntegers",
			args.U8, args.U16, big.NewInt(int64(args.U24)), args.U32, args.U64,
			args.I8, args.I16, big.NewInt(int64(args.I24)), args.I32, args.I64)
		if err != nil {
			b.Fatal(err)
		}
	}
}
