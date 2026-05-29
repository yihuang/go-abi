//go:build !uint256

package tests

import (
	"testing"

	"github.com/yihuang/go-abi"
)

// Lazy view benchmarks - compare view creation and field access patterns
// Run with: go test -bench='LazyView' -benchmem ./tests/...

func mustEncode(b *testing.B, v abi.Encode) []byte {
	encoded, err := v.Encode()
	if err != nil {
		b.Fatal(err)
	}
	return encoded
}

// must avoids b.Helper(), which adds ~68 ns/op when called in a b.N loop and
// would dominate these benchmarks.
func must(b *testing.B, err error) {
	if err != nil {
		b.Fatal(err)
	}
}

// SmallIntegers - all static fields (10 fields: uint8-64, int8-64)
// Best case for lazy views: no dynamic fields to parse
func BenchmarkGoABI_Decode_SmallIntegers_LazyView(b *testing.B) {
	args := createSmallIntegersData()
	encoded := mustEncode(b, args)

	b.Run("CreateView", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := DecodeTestSmallIntegersCallView(encoded)
			must(b, err)
		}
	})

	view, _, err := DecodeTestSmallIntegersCallView(encoded)
	must(b, err)

	b.Run("AccessOneField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := view.U64()
			must(b, err)
		}
	})

	b.Run("AccessAllFields", func(b *testing.B) {
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
	})
}

// MixedTypes - static + dynamic fields (bytes32, bytes, bool, uint8, Item[])
func BenchmarkGoABI_Decode_MixedTypes_LazyView(b *testing.B) {
	args := createMixedTypesData()
	encoded := mustEncode(b, args)

	b.Run("CreateView", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := DecodeTestMixedTypesCallView(encoded)
			must(b, err)
		}
	})

	view, _, err := DecodeTestMixedTypesCallView(encoded)
	must(b, err)

	b.Run("AccessStaticField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := view.Flag()
			must(b, err)
		}
	})

	b.Run("AccessDynamicField", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := view.DynamicData()
			must(b, err)
		}
	})

	b.Run("AccessItems", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			itemsView, err := view.Items()
			must(b, err)
			for j := 0; j < itemsView.Len(); j++ {
				_, err = itemsView.Get(j)
				must(b, err)
			}
		}
	})

	b.Run("Materialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, _, err := DecodeTestMixedTypesCallView(encoded)
			must(b, err)
			_, err = v.Materialize()
			must(b, err)
		}
	})
}

// ComplexDynamicTuples - deeply nested tuples (User2[] with UserProfile, UserMetadata2)
func BenchmarkGoABI_Decode_ComplexDynamicTuples_LazyView(b *testing.B) {
	args := createComplexDynamicTuplesData()
	encoded := mustEncode(b, args)

	b.Run("CreateView", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := DecodeTestComplexDynamicTuplesCallView(encoded)
			must(b, err)
		}
	})

	view, _, err := DecodeTestComplexDynamicTuplesCallView(encoded)
	must(b, err)

	b.Run("AccessFirstUserName", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			usersView, err := view.Users()
			must(b, err)
			userView, err := usersView.Get(0)
			must(b, err)
			profileView, err := userView.Profile()
			must(b, err)
			_, err = profileView.Name()
			must(b, err)
		}
	})

	b.Run("Materialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, _, err := DecodeTestComplexDynamicTuplesCallView(encoded)
			must(b, err)
			_, err = v.Materialize()
			must(b, err)
		}
	})
}

// NestedDynamicArrays - multi-dimensional arrays (uint256[][], address[][3][], string[][])
func BenchmarkGoABI_Decode_NestedDynamicArrays_LazyView(b *testing.B) {
	args := createNestedDynamicArraysData()
	encoded := mustEncode(b, args)

	b.Run("CreateView", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := DecodeTestNestedDynamicArraysCallView(encoded)
			must(b, err)
		}
	})

	b.Run("Materialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, _, err := DecodeTestNestedDynamicArraysCallView(encoded)
			must(b, err)
			_, err = v.Materialize()
			must(b, err)
		}
	})
}

// DeeplyNested - 4 levels of nesting (Level1 -> Level2 -> Level3 -> Level4)
func BenchmarkGoABI_Decode_DeeplyNested_LazyView(b *testing.B) {
	args := createDeeplyNestedData()
	encoded := mustEncode(b, args)

	b.Run("CreateView", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := DecodeTestDeeplyNestedCallView(encoded)
			must(b, err)
		}
	})

	view, _, err := DecodeTestDeeplyNestedCallView(encoded)
	must(b, err)

	b.Run("AccessInnerValue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dataView, err := view.Data()
			must(b, err)
			l2View, err := dataView.Level1()
			must(b, err)
			l3View, err := l2View.Level2()
			must(b, err)
			l4View, err := l3View.Level3()
			must(b, err)
			_, err = l4View.Value()
			must(b, err)
		}
	})
}

// LeafAccess compares "read one deeply-nested leaf" against full Decode +
// access of the same leaf. This is the workload the lazy option is meant to
// win at — pulling out a single field from a payload without paying for the
// rest. Run alongside BenchmarkGoABI_Decode_DeeplyNested for a Materialize
// baseline.
func BenchmarkGoABI_Decode_DeeplyNested_LeafAccess(b *testing.B) {
	args := createDeeplyNestedData()
	encoded := mustEncode(b, args)

	b.Run("FullDecodeThenLeaf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var decoded TestDeeplyNestedCall
			_, err := decoded.Decode(encoded)
			must(b, err)
			_ = decoded.Data.Level1.Level2.Level3.Value
		}
	})

	b.Run("LazyViewLeaf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			view, _, err := DecodeTestDeeplyNestedCallView(encoded)
			must(b, err)
			dataView, err := view.Data()
			must(b, err)
			l2View, err := dataView.Level1()
			must(b, err)
			l3View, err := l2View.Level2()
			must(b, err)
			l4View, err := l3View.Level3()
			must(b, err)
			_, err = l4View.Value()
			must(b, err)
		}
	})
}

// FixedTypes - static arrays and fixed bytes
func BenchmarkGoABI_Decode_FixedTypes_LazyView(b *testing.B) {
	b.Run("FixedArrays", func(b *testing.B) {
		encoded := mustEncode(b, createFixedArraysData())
		for i := 0; i < b.N; i++ {
			_, _, err := DecodeTestFixedArraysCallView(encoded)
			must(b, err)
		}
	})

	b.Run("FixedBytes", func(b *testing.B) {
		encoded := mustEncode(b, createFixedBytesData())
		for i := 0; i < b.N; i++ {
			_, _, err := DecodeTestFixedBytesCallView(encoded)
			must(b, err)
		}
	})
}
