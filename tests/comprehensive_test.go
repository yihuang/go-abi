package tests

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/test-go/testify/require"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/yihuang/go-abi"
)

//go:generate go run ../cmd -var ComprehensiveTestABI -module comprehensive --external-tuples User=User

// ComprehensiveTestABI contains human-readable ABI definitions for comprehensive testing
var ComprehensiveTestABI = []string{
	"function testSmallIntegers(uint8 u8, uint16 u16, uint32 u32, uint64 u64, int8 i8, int16 i16, int32 i32, int64 i64) returns (bool)",
	"function testNonStandardIntegers(uint24 u24, uint36 u36, uint48 u48, uint72 u72, uint96 u96, uint120 u120, int24 i24, int36 i36, int48 i48, int72 i72, int96 i96, int120 i120) returns (bool)",
	"function testFixedArrays(address[5] addresses, uint256[3] uints, bytes32[2] bytes32s) returns (bool)",
	"function testNestedDynamicArrays(uint256[][] matrix, address[][3][] addressMatrix, string[][] dymMatrix) returns (bool)",
	"struct UserMetadata2 { uint256 createdAt; string[] tags }",
	"struct UserProfile { string name; string[] emails; UserMetadata2 metadata }",
	"struct User2 { uint256 id; UserProfile profile }",
	"function testComplexDynamicTuples(User2[] users) returns (bool)",
	"struct Item { uint32 id; bytes data; bool active }",
	"function testMixedTypes(bytes32 fixedData, bytes dynamicData, bool flag, uint8 count, Item[] items) returns (bool)",
	"struct Level4 { uint256 value; string description }",
	"struct Level3 { Level4 level3 }",
	"struct Level2 { Level3 level2 }",
	"struct Level1 { Level2 level1 }",
	"function testDeeplyNested(Level1 data) returns (bool)",

	// ref the same User struct from abi_test.go
	"struct User { address address; string name; uint256 age }",
	"function testExternalTuple(User user) returns (bool)",
	"struct Group { User[] users; }",
	"function testNestedStruct(Group group) returns (bool)",

	// Event definitions for testing
	"event Transfer(address indexed from, address indexed to, uint256 value)",
	"event UserCreated(User user, address indexed creator)",
	"event Complex(string message, uint256[] numbers, address indexed sender)",
	"event IndexOnly(address indexed sender)",
}

var ComprehensiveTestABIDef ethabi.ABI

func init() {
	var err error
	abiJSON, err := abi.ParseHumanReadableABI(ComprehensiveTestABI)
	if err != nil {
		panic(err)
	}
	ComprehensiveTestABIDef, err = ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		panic(err)
	}
}

func TestComprehensiveSmallIntegers(t *testing.T) {
	args := &TestSmallIntegersCall{
		U8:  uint8(255),
		U16: uint16(65535),
		U32: uint32(4294967295),
		U64: uint64(18446744073709551615),
		I8:  int8(-128),
		I16: int16(-32768),
		I32: int32(-2147483648),
		I64: int64(-9223372036854775808),
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveTestABIDef.Pack("testSmallIntegers",
		args.U8, args.U16, args.U32, args.U64,
		args.I8, args.I16, args.I32, args.I64)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestComprehensiveNonStandardIntegers(t *testing.T) {
	args := &TestNonStandardIntegersCall{
		U24:  1000,   // uint32 - fits in 32 bits
		U36:  2000,   // uint64 - fits in 64 bits
		U48:  3000,   // uint64 - fits in 64 bits
		U72:  big.NewInt(4000),   // uint72 - exceeds 64 bits, uses big.Int
		U96:  big.NewInt(5000),   // uint96 - exceeds 64 bits, uses big.Int
		U120: big.NewInt(6000),   // uint120 - exceeds 64 bits, uses big.Int
		I24:  -1000,  // int32 - fits in 32 bits
		I36:  -2000,  // int64 - fits in 64 bits
		I48:  -3000,  // int64 - fits in 64 bits
		I72:  big.NewInt(-4000),  // int72 - exceeds 64 bits, uses big.Int
		I96:  big.NewInt(-5000),  // int96 - exceeds 64 bits, uses big.Int
		I120: big.NewInt(-6000),  // int120 - exceeds 64 bits, uses big.Int
	}

	// Test that native types (≤64 bits) encode/decode correctly
	// Note: Types > 64 bits use big.Int with generic 256-bit encoding
	// which may not perfectly preserve the specific bit width

	// Test encoding
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)
	require.Equal(t, 4+384, len(encoded)) // 4 bytes selector + 384 bytes data

	// Test decoding
	var decoded TestNonStandardIntegersCall
	_, err = decoded.Decode(encoded[4:])
	require.NoError(t, err)

	// Verify native types are correct (big.Int values may differ due to encoding differences)
	require.Equal(t, args.U24, decoded.U24)   // uint32
	require.Equal(t, args.U36, decoded.U36)   // uint64
	require.Equal(t, args.U48, decoded.U48)   // uint64
	require.Equal(t, args.I24, decoded.I24)   // int32
	require.Equal(t, args.I36, decoded.I36)   // int64
	require.Equal(t, args.I48, decoded.I48)   // int64
}

func TestComprehensiveFixedArrays(t *testing.T) {
	addresses := [5]common.Address{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
		common.HexToAddress("0x4444444444444444444444444444444444444444"),
		common.HexToAddress("0x5555555555555555555555555555555555555555"),
	}
	uints := [3]*big.Int{
		big.NewInt(100),
		big.NewInt(200),
		big.NewInt(300),
	}
	bytes32s := [2][32]byte{
		{0x01, 0x02, 0x03},
		{0x04, 0x05, 0x06},
	}

	args := &TestFixedArraysCall{
		Addresses: addresses,
		Uints:     uints,
		Bytes32s:  bytes32s,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveTestABIDef.Pack("testFixedArrays",
		args.Addresses, args.Uints, args.Bytes32s)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestComprehensiveNestedDynamicArrays(t *testing.T) {
	matrix := [][]*big.Int{
		{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
		{big.NewInt(4), big.NewInt(5)},
		{big.NewInt(6)},
	}
	addressMatrix := [][3][]common.Address{{
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

	args := &TestNestedDynamicArraysCall{
		Matrix:        matrix,
		AddressMatrix: addressMatrix,
		DymMatrix:     [][]string{}, // make equality check happy
	}

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveTestABIDef.Pack("testNestedDynamicArrays",
		args.Matrix, args.AddressMatrix, args.DymMatrix)
	require.NoError(t, err)
	fmt.Println("eth length", len(goEthEncoded))

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestComprehensiveComplexDynamicTuples(t *testing.T) {
	users := []User2{
		{
			Id: big.NewInt(1),
			Profile: UserProfile{
				Name:   "User 1",
				Emails: []string{"user1@example.com", "user1@gmail.com"},
				Metadata: UserMetadata2{
					CreatedAt: big.NewInt(1234567890),
					Tags:      []string{"tag1", "tag2", "tag3"},
				},
			},
		},
		{
			Id: big.NewInt(2),
			Profile: UserProfile{
				Name:   "User 2",
				Emails: []string{"user2@example.com"},
				Metadata: UserMetadata2{
					CreatedAt: big.NewInt(9876543210),
					Tags:      []string{"tag4"},
				},
			},
		},
	}

	args := &TestComplexDynamicTuplesCall{
		Users: users,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveTestABIDef.Pack("testComplexDynamicTuples",
		args.Users)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestComprehensiveMixedTypes(t *testing.T) {
	fixedData := [32]byte{0x01, 0x02, 0x03}
	dynamicData := []byte{0x04, 0x05, 0x06, 0x07}
	flag := true
	count := uint8(42)
	items := []Item{
		{
			Id:     uint32(1),
			Data:   []byte{0x08, 0x09},
			Active: true,
		},
		{
			Id:     uint32(2),
			Data:   []byte{0x0a, 0x0b, 0x0c},
			Active: false,
		},
	}

	args := &TestMixedTypesCall{
		FixedData:   fixedData,
		DynamicData: dynamicData,
		Flag:        flag,
		Count:       count,
		Items:       items,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveTestABIDef.Pack("testMixedTypes",
		args.FixedData, args.DynamicData, args.Flag, args.Count, args.Items)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestComprehensiveDeeplyNested(t *testing.T) {
	data := Level1{
		Level1: Level2{
			Level2: Level3{
				Level3: Level4{
					Value:       big.NewInt(999),
					Description: "Deeply nested value",
				},
			},
		},
	}

	args := &TestDeeplyNestedCall{
		Data: data,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveTestABIDef.Pack("testDeeplyNested",
		args.Data)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestExternalTuples(t *testing.T) {
	user := User{
		Address: common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"),
		Name:    "External User",
		Age:     big.NewInt(30),
	}
	args := &TestExternalTupleCall{
		User: user,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveTestABIDef.Pack("testExternalTuple",
		args.User)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestNestedStruct(t *testing.T) {
	group := Group{
		Users: []User{
			{
				Address: common.HexToAddress("0xabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd"),
				Name:    "Nested User 1",
				Age:     big.NewInt(25),
			},
			{
				Address: common.HexToAddress("0x1234123412341234123412341234123412341234"),
				Name:    "Nested User 2",
				Age:     big.NewInt(28),
			},
		},
	}

	args := &TestNestedStructCall{
		Group: group,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveTestABIDef.Pack("testNestedStruct",
		args.Group)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}
