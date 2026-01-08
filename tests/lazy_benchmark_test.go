//go:build !uint256

package tests

import (
	"testing"
)

// Lazy view benchmarks - compare view creation and field access patterns
// Run with: go test -bench='Lazy' -benchmem ./tests/...

// SmallIntegers - all static fields (10 fields: uint8-64, int8-64)
// Best case for lazy views: no dynamic fields to parse

func BenchmarkGoABI_LazyView_SmallIntegers(b *testing.B) {
	args := createSmallIntegersData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := DecodeTestSmallIntegersCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_SmallIntegers_AccessOneField(b *testing.B) {
	args := createSmallIntegersData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	view, _, err := DecodeTestSmallIntegersCallView(encoded)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := view.U64()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_SmallIntegers_AccessAllFields(b *testing.B) {
	args := createSmallIntegersData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	view, _, err := DecodeTestSmallIntegersCallView(encoded)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = view.U8()
		_, _ = view.U16()
		_, _ = view.U24()
		_, _ = view.U32()
		_, _ = view.U64()
		_, _ = view.I8()
		_, _ = view.I16()
		_, _ = view.I24()
		_, _ = view.I32()
		_, _ = view.I64()
	}
}

// MixedTypes - static + dynamic fields (bytes32, bytes, bool, uint8, Item[])

func BenchmarkGoABI_LazyView_MixedTypes(b *testing.B) {
	args := createMixedTypesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := DecodeTestMixedTypesCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_MixedTypes_AccessStaticField(b *testing.B) {
	args := createMixedTypesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	view, _, err := DecodeTestMixedTypesCallView(encoded)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := view.Flag()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_MixedTypes_AccessDynamicField(b *testing.B) {
	args := createMixedTypesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	view, _, err := DecodeTestMixedTypesCallView(encoded)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := view.DynamicData()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_MixedTypes_AccessItems(b *testing.B) {
	args := createMixedTypesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	view, _, err := DecodeTestMixedTypesCallView(encoded)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		itemsView, err := view.Items()
		if err != nil {
			b.Fatal(err)
		}
		for j := 0; j < itemsView.Len(); j++ {
			_, err = itemsView.Get(j)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkGoABI_LazyView_MixedTypes_Materialize(b *testing.B) {
	args := createMixedTypesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		view, _, err := DecodeTestMixedTypesCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
		_, err = view.Materialize()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ComplexDynamicTuples - deeply nested tuples (User2[] with UserProfile, UserMetadata2)

func BenchmarkGoABI_LazyView_ComplexDynamicTuples(b *testing.B) {
	args := createComplexDynamicTuplesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := DecodeTestComplexDynamicTuplesCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_ComplexDynamicTuples_AccessFirstUser(b *testing.B) {
	args := createComplexDynamicTuplesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	view, _, err := DecodeTestComplexDynamicTuplesCallView(encoded)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		usersView, err := view.Users()
		if err != nil {
			b.Fatal(err)
		}
		_, err = usersView.Get(0)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_ComplexDynamicTuples_AccessFirstUserName(b *testing.B) {
	args := createComplexDynamicTuplesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	view, _, err := DecodeTestComplexDynamicTuplesCallView(encoded)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		usersView, err := view.Users()
		if err != nil {
			b.Fatal(err)
		}
		userView, err := usersView.Get(0)
		if err != nil {
			b.Fatal(err)
		}
		profileView, err := userView.Profile()
		if err != nil {
			b.Fatal(err)
		}
		_, err = profileView.Name()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_ComplexDynamicTuples_Materialize(b *testing.B) {
	args := createComplexDynamicTuplesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		view, _, err := DecodeTestComplexDynamicTuplesCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
		_, err = view.Materialize()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// NestedDynamicArrays - multi-dimensional arrays (uint256[][], address[][3][], string[][])

func BenchmarkGoABI_LazyView_NestedDynamicArrays(b *testing.B) {
	args := createNestedDynamicArraysData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := DecodeTestNestedDynamicArraysCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_NestedDynamicArrays_Materialize(b *testing.B) {
	args := createNestedDynamicArraysData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		view, _, err := DecodeTestNestedDynamicArraysCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
		_, err = view.Materialize()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// DeeplyNested - 4 levels of nesting (Level1 -> Level2 -> Level3 -> Level4)

func BenchmarkGoABI_LazyView_DeeplyNested(b *testing.B) {
	args := createDeeplyNestedData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := DecodeTestDeeplyNestedCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGoABI_LazyView_DeeplyNested_AccessInnerValue(b *testing.B) {
	args := createDeeplyNestedData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	view, _, err := DecodeTestDeeplyNestedCallView(encoded)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dataView, err := view.Data()
		if err != nil {
			b.Fatal(err)
		}
		l2View, err := dataView.Level1()
		if err != nil {
			b.Fatal(err)
		}
		l3View, err := l2View.Level2()
		if err != nil {
			b.Fatal(err)
		}
		l4View, err := l3View.Level3()
		if err != nil {
			b.Fatal(err)
		}
		_, err = l4View.Value()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// FixedArrays - static arrays (address[5], uint256[3], bytes32[2])

func BenchmarkGoABI_LazyView_FixedArrays(b *testing.B) {
	args := createFixedArraysData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := DecodeTestFixedArraysCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// FixedBytes - fixed size bytes (bytes3, bytes7, bytes15)

func BenchmarkGoABI_LazyView_FixedBytes(b *testing.B) {
	args := createFixedBytesData()
	encoded, err := args.Encode()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := DecodeTestFixedBytesCallView(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}
