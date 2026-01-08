//go:build !uint256

package tests

import (
	"testing"
)

// Benchmark functions for go-abi generated code
func BenchmarkGoABI_Decode_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var decoded TestComplexDynamicTuplesCall
		_, err := decoded.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_Decode_NestedDynamicArrays(b *testing.B) {
	args := createNestedDynamicArraysData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var decoded TestNestedDynamicArraysCall
		_, err := decoded.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_Decode_MixedTypes(b *testing.B) {
	args := createMixedTypesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var decoded TestMixedTypesCall
		_, err := decoded.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark functions for go-ethereum ABI decoding
func BenchmarkGoEthereum_Decode_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	encoded, err := ComprehensiveTestABIDef.Pack("testComplexDynamicTuples", args.Users)
	if err != nil {
		b.Fatal(err)
	}
	arguments := ComprehensiveTestABIDef.Methods["testComplexDynamicTuples"].Inputs

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := arguments.Unpack(encoded[4:])
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoEthereum_Decode_NestedDynamicArrays(b *testing.B) {
	args := createNestedDynamicArraysData()
	encoded, err := ComprehensiveTestABIDef.Pack("testNestedDynamicArrays",
		args.Matrix, args.AddressMatrix, args.DymMatrix)
	if err != nil {
		b.Fatal(err)
	}
	arguments := ComprehensiveTestABIDef.Methods["testNestedDynamicArrays"].Inputs
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := arguments.Unpack(encoded[4:])
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoEthereum_Decode_MixedTypes(b *testing.B) {
	args := createMixedTypesData()
	encoded, err := ComprehensiveTestABIDef.Pack("testMixedTypes",
		args.FixedData, args.DynamicData, args.Flag, args.Count, args.Items)
	if err != nil {
		b.Fatal(err)
	}
	arguments := ComprehensiveTestABIDef.Methods["testMixedTypes"].Inputs

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := arguments.Unpack(encoded[4:])
		if err != nil {
			b.Fatal(err)
		}
	}
}
