//go:build !uint256

package tests

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/test-go/testify/require"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/yihuang/go-abi"
)

//go:generate go run ../cmd -var ComplexReturnsABI -output complex_returns.abi.go -prefix complexreturns

// ComplexReturnsABI contains human-readable ABI definitions for testing complex return types
// including nested tuples, fixed-length arrays of tuples, dynamic arrays of tuples, deeply nested structs, etc.
var ComplexReturnsABI = []string{
	// Basic structs for nesting
	"struct Simple { uint256 value; string name }",
	"struct NestedLevel1 { Simple inner; uint256 extra }",
	"struct NestedLevel2 { NestedLevel1 level1; bool active }",

	// Tuples with arrays
	"struct TupleWithArray { uint256 id; uint256[] values }",
	"struct TupleWithFixedArray { uint256 id; uint256[3] values }",
	"struct NestedTupleArray { string label; TupleWithArray[] items }",

	// Mixed types
	"struct MixedReturn { uint256 a; string b; bool c; address d; bytes e }",

	// Fixed array of tuples
	"struct FixedTupleArray { Simple[2] pairs }",

	// Deeply nested (3+ levels)
	"struct DeepNested { uint256 id; NestedLevel2 nested }",

	// Multi-dimensional array of tuples
	"struct MatrixStruct { uint256 id; Simple[][] rows }",

	// Fixed-size arrays inside tuples
	"struct FixedBytesInTuple { bytes32 hash; bytes4 prefix }",

	// Struct with all basic types
	"struct AllTypes { uint256 u; int256 i; address a; bool b; bytes32 h; bytes d; string s }",

	// Functions returning nested tuples
	"function getNestedLevel1() returns (NestedLevel1)",
	"function getNestedLevel2() returns (NestedLevel2)",
	"function getDeepNested() returns (DeepNested)",

	// Functions returning tuples with arrays
	"function getTupleWithArray() returns (TupleWithArray)",
	"function getTupleWithFixedArray() returns (TupleWithFixedArray)",
	"function getNestedTupleArrayReturn() returns (NestedTupleArray)",

	// Functions returning mixed types
	"function getMixedReturn() returns (MixedReturn)",

	// Functions returning arrays of tuples
	"function getSimpleArray() returns (Simple[])",
	"function getNestedLevel1Array() returns (NestedLevel1[])",
	"function getSimpleFixedArray() returns (Simple[2])",
	"function getNestedLevel1FixedArray() returns (NestedLevel1[3])",

	// Multiple return values with complex types
	"function getMultipleComplex() returns (NestedLevel1, Simple, bool)",

	// Inline tuple return
	"function getInlineTupleReturn() returns ((uint256 a, string b) pair, (address c, bool d) other)",

	// Fixed array of tuples
	"function getFixedTupleArray() returns (FixedTupleArray)",

	// Fixed bytes inside tuple
	"function getFixedBytesInTuple() returns (FixedBytesInTuple)",

	// All basic types
	"function getAllTypesReturn() returns (AllTypes)",

	// Array of tuples with nested dynamic elements
	"function getMatrixStruct() returns (MatrixStruct)",

	// Function returning empty tuple
	"function getEmptyReturn() returns ()",
}

var ComplexReturnsABIDef ethabi.ABI

func init() {
	var err error
	abiJSON, err := abi.ParseHumanReadableABI(ComplexReturnsABI)
	if err != nil {
		panic(err)
	}
	ComplexReturnsABIDef, err = ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		panic(err)
	}
}

// =============================================================================
// Nested Tuple Returns Tests
// =============================================================================

func TestComplexReturnsNestedLevel1(t *testing.T) {
	args := &GetNestedLevel1Return{
		Field1: NestedLevel1{
			Inner: Simple{
				Value: big.NewInt(42),
				Name:  "hello",
			},
			Extra: big.NewInt(100),
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getNestedLevel1"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetNestedLevel1Return
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

func TestComplexReturnsNestedLevel2(t *testing.T) {
	args := &GetNestedLevel2Return{
		Field1: NestedLevel2{
			Level1: NestedLevel1{
				Inner: Simple{
					Value: big.NewInt(999),
					Name:  "deeply nested",
				},
				Extra: big.NewInt(888),
			},
			Active: true,
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getNestedLevel2"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetNestedLevel2Return
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

func TestComplexReturnsDeepNested(t *testing.T) {
	args := &GetDeepNestedReturn{
		Field1: DeepNested{
			Id: big.NewInt(1),
			Nested: NestedLevel2{
				Level1: NestedLevel1{
					Inner: Simple{
						Value: big.NewInt(777),
						Name:  "deep",
					},
					Extra: big.NewInt(666),
				},
				Active: false,
			},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getDeepNested"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetDeepNestedReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// Tuples with Arrays Returns Tests
// =============================================================================

func TestComplexReturnsTupleWithArray(t *testing.T) {
	args := &GetTupleWithArrayReturn{
		Field1: TupleWithArray{
			Id:     big.NewInt(123),
			Values: []*big.Int{big.NewInt(10), big.NewInt(20), big.NewInt(30)},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getTupleWithArray"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetTupleWithArrayReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

func TestComplexReturnsTupleWithFixedArray(t *testing.T) {
	args := &GetTupleWithFixedArrayReturn{
		Field1: TupleWithFixedArray{
			Id: big.NewInt(456),
			Values: [3]*big.Int{
				big.NewInt(100),
				big.NewInt(200),
				big.NewInt(300),
			},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getTupleWithFixedArray"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetTupleWithFixedArrayReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

func TestComplexReturnsNestedTupleArray(t *testing.T) {
	args := &GetNestedTupleArrayReturnReturn{
		Field1: NestedTupleArray{
			Label: "test items",
			Items: []TupleWithArray{
				{
					Id:     big.NewInt(1),
					Values: []*big.Int{big.NewInt(10), big.NewInt(20)},
				},
				{
					Id:     big.NewInt(2),
					Values: []*big.Int{big.NewInt(30)},
				},
			},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getNestedTupleArrayReturn"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetNestedTupleArrayReturnReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// Mixed Types Return Test
// =============================================================================

func TestComplexReturnsMixedReturn(t *testing.T) {
	args := &GetMixedReturnReturn{
		Field1: MixedReturn{
			A: big.NewInt(42),
			B: "mixed string",
			C: true,
			D: common.HexToAddress("0x1111111111111111111111111111111111111111"),
			E: []byte{0x01, 0x02, 0x03, 0x04},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getMixedReturn"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetMixedReturnReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// Arrays of Tuples Returns Tests
// =============================================================================

func TestComplexReturnsSimpleArray(t *testing.T) {
	args := &GetSimpleArrayReturn{
		Field1: []Simple{
			{Value: big.NewInt(1), Name: "first"},
			{Value: big.NewInt(2), Name: "second"},
			{Value: big.NewInt(3), Name: "third"},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getSimpleArray"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetSimpleArrayReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Len(t, decoded.Field1, 3)
	require.Equal(t, args, &decoded)
}

func TestComplexReturnsNestedLevel1Array(t *testing.T) {
	args := &GetNestedLevel1ArrayReturn{
		Field1: []NestedLevel1{
			{
				Inner: Simple{Value: big.NewInt(10), Name: "a"},
				Extra: big.NewInt(100),
			},
			{
				Inner: Simple{Value: big.NewInt(20), Name: "b"},
				Extra: big.NewInt(200),
			},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getNestedLevel1Array"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetNestedLevel1ArrayReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Len(t, decoded.Field1, 2)
	require.Equal(t, args, &decoded)
}

func TestComplexReturnsSimpleFixedArray(t *testing.T) {
	args := &GetSimpleFixedArrayReturn{
		Field1: [2]Simple{
			{Value: big.NewInt(100), Name: "fixed a"},
			{Value: big.NewInt(200), Name: "fixed b"},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getSimpleFixedArray"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetSimpleFixedArrayReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args.Field1, decoded.Field1)
	require.Equal(t, args, &decoded)
}

func TestComplexReturnsNestedLevel1FixedArray(t *testing.T) {
	args := &GetNestedLevel1FixedArrayReturn{
		Field1: [3]NestedLevel1{
			{
				Inner: Simple{Value: big.NewInt(1), Name: "x"},
				Extra: big.NewInt(10),
			},
			{
				Inner: Simple{Value: big.NewInt(2), Name: "y"},
				Extra: big.NewInt(20),
			},
			{
				Inner: Simple{Value: big.NewInt(3), Name: "z"},
				Extra: big.NewInt(30),
			},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getNestedLevel1FixedArray"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetNestedLevel1FixedArrayReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args.Field1, decoded.Field1)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// Multiple Complex Returns Test
// =============================================================================

func TestComplexReturnsMultipleComplex(t *testing.T) {
	args := &GetMultipleComplexReturn{
		Field1: NestedLevel1{
			Inner: Simple{
				Value: big.NewInt(55),
				Name:  "multi inner",
			},
			Extra: big.NewInt(66),
		},
		Field2: Simple{
			Value: big.NewInt(77),
			Name:  "multi simple",
		},
		Field3: true,
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getMultipleComplex"].Outputs.Pack(args.Field1, args.Field2, args.Field3)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetMultipleComplexReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// Inline Tuple Return Test
// =============================================================================

func TestComplexReturnsInlineTupleReturn(t *testing.T) {
	// The inline tuple return:
	//   returns ((uint256 a, string b) pair, (address c, bool d) other)
	// Generates struct with named fields Pair and Other.

	// Build the data using go-ethereum
	pairValue := struct {
		A *big.Int
		B string
	}{A: big.NewInt(12345), B: "inline pair"}

	otherValue := struct {
		C common.Address
		D bool
	}{C: common.HexToAddress("0x2222222222222222222222222222222222222222"), D: false}

	method, ok := ComplexReturnsABIDef.Methods["getInlineTupleReturn"]
	require.True(t, ok)

	goEthEncoded, err := method.Outputs.Pack(pairValue, otherValue)
	require.NoError(t, err)

	// Now decode using our generated types
	var decoded GetInlineTupleReturnReturn
	_, err = decoded.Decode(goEthEncoded)
	require.NoError(t, err)

	// Verify the decoded values match
	require.Equal(t, pairValue.A, decoded.Pair.A)
	require.Equal(t, pairValue.B, decoded.Pair.B)
	require.Equal(t, otherValue.C, decoded.Other.C)
	require.Equal(t, otherValue.D, decoded.Other.D)

	// Now encode using our types
	encoded, err := decoded.Encode()
	require.NoError(t, err)
	require.Equal(t, goEthEncoded, encoded)
}

// =============================================================================
// Fixed Tuple Array Return Test
// =============================================================================

func TestComplexReturnsFixedTupleArray(t *testing.T) {
	args := &GetFixedTupleArrayReturn{
		Field1: FixedTupleArray{
			Pairs: [2]Simple{
				{Value: big.NewInt(111), Name: "pair1"},
				{Value: big.NewInt(222), Name: "pair2"},
			},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getFixedTupleArray"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetFixedTupleArrayReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// Fixed Bytes Inside Tuple Return Test
// =============================================================================

func TestComplexReturnsFixedBytesInTuple(t *testing.T) {
	args := &GetFixedBytesInTupleReturn{
		Field1: FixedBytesInTuple{
			Hash:   [32]byte{0x01, 0x02, 0x03, 0x04, 0x05},
			Prefix: [4]byte{0xde, 0xad, 0xbe, 0xef},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getFixedBytesInTuple"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetFixedBytesInTupleReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// All Types Return Test
// =============================================================================

func TestComplexReturnsAllTypes(t *testing.T) {
	args := &GetAllTypesReturnReturn{
		Field1: AllTypes{
			U: big.NewInt(123456),
			I: big.NewInt(-7890),
			A: common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			B: true,
			H: [32]byte{0xff, 0xee, 0xdd},
			D: []byte{0xca, 0xfe, 0xba, 0xbe},
			S: "all types string",
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getAllTypesReturn"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetAllTypesReturnReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// Matrix Struct Return Test (multi-dimensional array of tuples)
// =============================================================================

func TestComplexReturnsMatrixStruct(t *testing.T) {
	args := &GetMatrixStructReturn{
		Field1: MatrixStruct{
			Id: big.NewInt(99),
			Rows: [][]Simple{
				{
					{Value: big.NewInt(1), Name: "r0c0"},
					{Value: big.NewInt(2), Name: "r0c1"},
				},
				{
					{Value: big.NewInt(3), Name: "r1c0"},
				},
			},
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getMatrixStruct"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Test decoding
	var decoded GetMatrixStructReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, args, &decoded)
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestComplexReturnsEmptyReturn(t *testing.T) {
	// Test empty return (no outputs)
	args := &GetEmptyReturnReturn{}

	encoded, err := args.Encode()
	require.NoError(t, err)
	require.Len(t, encoded, 0)

	var decoded GetEmptyReturnReturn
	n, err := decoded.Decode(encoded)
	require.NoError(t, err)
	require.Equal(t, 0, n)
	require.Equal(t, args, &decoded)

	// Verify go-ethereum compatibility (both should be empty)
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getEmptyReturn"].Outputs.Pack()
	require.NoError(t, err)
	require.Len(t, goEthEncoded, 0)
	require.Len(t, encoded, len(goEthEncoded))
}

func TestComplexReturnsEmptySlice(t *testing.T) {
	// Test empty dynamic array of tuples
	args := &GetSimpleArrayReturn{
		Field1: []Simple{},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test with go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getSimpleArray"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	var decoded GetSimpleArrayReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)
	require.Empty(t, decoded.Field1)
}

func TestComplexReturnsEmptyString(t *testing.T) {
	// Test empty string in deeply nested struct with zero numeric values
	args := &GetNestedLevel1Return{
		Field1: NestedLevel1{
			Inner: Simple{
				Value: big.NewInt(0),
				Name:  "",
			},
			Extra: big.NewInt(0),
		},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	// Verify encoding matches go-ethereum
	goEthEncoded, err := ComplexReturnsABIDef.Methods["getNestedLevel1"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	// Round-trip decode and verify the decoded values are semantically correct
	var decoded GetNestedLevel1Return
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)

	// Verify the values are correct (zero big.Int may differ structurally, check semantically)
	require.Equal(t, int64(0), decoded.Field1.Inner.Value.Int64())
	require.Equal(t, "", decoded.Field1.Inner.Name)
	require.Equal(t, int64(0), decoded.Field1.Extra.Int64())

	// Re-encode decoded value and verify it matches original encoding
	reeEncoded, err := decoded.Encode()
	require.NoError(t, err)
	require.Equal(t, encoded, reeEncoded)
}

func TestComplexReturnsLargeDynamicArrays(t *testing.T) {
	// Test with larger arrays in return
	items := make([]TupleWithArray, 5)
	for i := 0; i < 5; i++ {
		values := make([]*big.Int, 3)
		for j := 0; j < 3; j++ {
			values[j] = big.NewInt(int64(i*100 + j))
		}
		items[i] = TupleWithArray{
			Id:     big.NewInt(int64(i)),
			Values: values,
		}
	}

	args := &GetNestedTupleArrayReturnReturn{
		Field1: NestedTupleArray{
			Label: "large test",
			Items: items,
		},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	goEthEncoded, err := ComplexReturnsABIDef.Methods["getNestedTupleArrayReturn"].Outputs.Pack(args.Field1)
	require.NoError(t, err)
	require.Equal(t, encoded, goEthEncoded)

	var decoded GetNestedTupleArrayReturnReturn
	_, err = decoded.Decode(encoded)
	require.NoError(t, err)

	// Verify the decoded values semantically (big.Int zero values may differ structurally)
	require.Equal(t, args.Field1.Label, decoded.Field1.Label)
	require.Len(t, decoded.Field1.Items, len(args.Field1.Items))
	for i := range args.Field1.Items {
		require.Equal(t, 0, args.Field1.Items[i].Id.Cmp(decoded.Field1.Items[i].Id))
		require.Len(t, decoded.Field1.Items[i].Values, len(args.Field1.Items[i].Values))
		for j := range args.Field1.Items[i].Values {
			require.Equal(t, 0, args.Field1.Items[i].Values[j].Cmp(decoded.Field1.Items[i].Values[j]))
		}
	}

	// Re-encode decoded value and verify it matches original encoding
	reeEncoded, err := decoded.Encode()
	require.NoError(t, err)
	require.Equal(t, encoded, reeEncoded)
}
