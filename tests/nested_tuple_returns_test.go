package testdata

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/test-go/testify/require"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/yihuang/go-abi"
)

//go:generate go run ../cmd/main.go -var NestedTupleReturnsABI -output nested_tuple_returns.abi.go

// NestedTupleReturnsABI contains human-readable ABI definitions for testing nested tuples in function returns
var NestedTupleReturnsABI = []string{
	// Define structs for nested tuples
	"struct SimplePair { uint256 first; uint256 second }",
	"struct AddressStringPair { address addr; string str }",
	"struct ComplexNested { uint256 num; address addr; string str; bytes data }",
	"struct DeeplyNested { uint256 num; string str; bool flag; address addr; bytes32 hash }",
	"struct UserWithMetadata { string name; uint256 id; uint256 age; string metadata }",

	// Simple nested tuple returns
	"function getSimplePair() returns (SimplePair)",
	"function getAddressStringPair() returns (AddressStringPair)",

	// Complex nested tuple returns
	"function getComplexNested() returns (ComplexNested)",
	"function getDeeplyNested() returns (DeeplyNested)",

	// Array of nested tuples
	"function getTupleArray() returns (SimplePair[])",
	"function getNestedTupleArray() returns (ComplexNested[])",

	// Mixed nested tuples with inline structs
	"function getUserWithMetadata() returns (UserWithMetadata)",
	"function getUsersArray() returns (AddressStringPair[])",

	// Multiple return values with nested tuples
	"function getMultipleReturns() returns (uint256, AddressStringPair, bool)",
}

var NestedTupleReturnsABIDef ethabi.ABI

func init() {
	var err error
	abiJSON, err := abi.ParseHumanReadableABI(NestedTupleReturnsABI)
	if err != nil {
		panic(err)
	}
	NestedTupleReturnsABIDef, err = ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		panic(err)
	}
}

func TestNestedTupleReturnsSimplePair(t *testing.T) {
	// Test encoding and decoding of simple tuple return
	args := &GetSimplePairReturn{
		Result1: SimplePair{
			First:  big.NewInt(100),
			Second: big.NewInt(200),
		},
	}

	// Test encoding
	encoded, err := args.Encode()
	require.NoError(t, err)

	// Test decoding
	var decoded GetSimplePairReturn
	err = decoded.Decode(encoded)
	require.NoError(t, err)

	require.Equal(t, args, &decoded)

	// Test with go-ethereum
	goEthEncoded, err := NestedTupleReturnsABIDef.Methods["getSimplePair"].Outputs.Pack(args.Result1)
	require.NoError(t, err)

	// The return data should match our encoding
	require.Equal(t, encoded, goEthEncoded)
}

func TestNestedTupleReturnsAddressStringPair(t *testing.T) {
	args := &GetAddressStringPairReturn{
		Result1: AddressStringPair{
			Addr: common.HexToAddress("0x1111111111111111111111111111111111111111"),
			Str:  "test string",
		},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	var decoded GetAddressStringPairReturn
	err = decoded.Decode(encoded)
	require.NoError(t, err)

	require.Equal(t, args, &decoded)
}

func TestNestedTupleReturnsComplexNested(t *testing.T) {
	args := &GetComplexNestedReturn{
		Result1: ComplexNested{
			Num:  big.NewInt(999),
			Addr: common.HexToAddress("0x2222222222222222222222222222222222222222"),
			Str:  "test string",
			Data: []byte{0x01, 0x02, 0x03},
		},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	var decoded GetComplexNestedReturn
	err = decoded.Decode(encoded)
	require.NoError(t, err)

	require.Equal(t, args, &decoded)
}

func TestNestedTupleReturnsDeeplyNested(t *testing.T) {
	args := &GetDeeplyNestedReturn{
		Result1: DeeplyNested{
			Num:  big.NewInt(12345),
			Str:  "deeply nested string",
			Flag: true,
			Addr: common.HexToAddress("0x3333333333333333333333333333333333333333"),
			Hash: [32]byte{0x01, 0x02, 0x03},
		},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	var decoded GetDeeplyNestedReturn
	err = decoded.Decode(encoded)
	require.NoError(t, err)

	require.Equal(t, args, &decoded)
}

func TestNestedTupleReturnsTupleArray(t *testing.T) {
	args := &GetTupleArrayReturn{
		Result1: []SimplePair{
			{
				First:  big.NewInt(1),
				Second: big.NewInt(2),
			},
			{
				First:  big.NewInt(3),
				Second: big.NewInt(4),
			},
		},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	var decoded GetTupleArrayReturn
	err = decoded.Decode(encoded)
	require.NoError(t, err)

	require.Len(t, decoded.Result1, 2)
	require.Equal(t, args, &decoded)
}

func TestNestedTupleReturnsUserWithMetadata(t *testing.T) {
	args := &GetUserWithMetadataReturn{
		Result1: UserWithMetadata{
			Name:     "Test User",
			Id:       big.NewInt(123),
			Age:      big.NewInt(30),
			Metadata: "test metadata",
		},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	var decoded GetUserWithMetadataReturn
	err = decoded.Decode(encoded)
	require.NoError(t, err)

	require.Equal(t, args, &decoded)
}

func TestNestedTupleReturnsUsersArray(t *testing.T) {
	args := &GetUsersArrayReturn{
		Result1: []AddressStringPair{
			{
				Addr: common.HexToAddress("0x5555555555555555555555555555555555555555"),
				Str:  "User 1",
			},
			{
				Addr: common.HexToAddress("0x6666666666666666666666666666666666666666"),
				Str:  "User 2",
			},
		},
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	var decoded GetUsersArrayReturn
	err = decoded.Decode(encoded)
	require.NoError(t, err)

	require.Len(t, decoded.Result1, 2)
	require.Equal(t, args, &decoded)
}

func TestNestedTupleReturnsMultipleReturns(t *testing.T) {
	args := &GetMultipleReturnsReturn{
		Result1: big.NewInt(42),
		Result2: AddressStringPair{
			Addr: common.HexToAddress("0x4444444444444444444444444444444444444444"),
			Str:  "multiple return string",
		},
		Result3: true,
	}

	encoded, err := args.Encode()
	require.NoError(t, err)

	var decoded GetMultipleReturnsReturn
	err = decoded.Decode(encoded)
	require.NoError(t, err)

	require.Equal(t, args, &decoded)
}
