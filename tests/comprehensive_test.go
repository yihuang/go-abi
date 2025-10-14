package testdata

import (
	"math/big"
	"strings"
	"testing"

	"github.com/test-go/testify/require"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	_ "embed"
)

var (
	//go:embed comprehensive_test_abi.json
	ComprehensiveABIJson string
	ComprehensiveABIDef  abi.ABI
)

func init() {
	var err error
	ComprehensiveABIDef, err = abi.JSON(strings.NewReader(ComprehensiveABIJson))
	if err != nil {
		panic(err)
	}
}

//go:generate go run github.com/yihuang/go-abi/cmd -input comprehensive_test_abi.json -output comprehensive_test_abi.abi.go

func TestComprehensiveSmallIntegers(t *testing.T) {
	args := &TestSmallIntegersArgs{
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
	goEthEncoded, err := ComprehensiveABIDef.Pack("testSmallIntegers",
		args.U8, args.U16, args.U32, args.U64,
		args.I8, args.I16, args.I32, args.I64)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)
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

	args := &TestFixedArraysArgs{
		Addresses: addresses,
		Uints:     uints,
		Bytes32s:  bytes32s,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveABIDef.Pack("testFixedArrays",
		args.Addresses, args.Uints, args.Bytes32s)
	require.NoError(t, err)

	require.Equal(t, len(encoded), len(goEthEncoded))
}

func TestComprehensiveNestedDynamicArrays(t *testing.T) {
	matrix := [][]*big.Int{
		{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
		{big.NewInt(4), big.NewInt(5)},
		{big.NewInt(6)},
	}
	addressMatrix := [][]common.Address{
		{
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
		},
		{
			common.HexToAddress("0x3333333333333333333333333333333333333333"),
		},
	}

	args := &TestNestedDynamicArraysArgs{
		Matrix:        matrix,
		AddressMatrix: addressMatrix,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveABIDef.Pack("testNestedDynamicArrays",
		args.Matrix, args.AddressMatrix)
	require.NoError(t, err)

	require.Equal(t, len(encoded), len(goEthEncoded))
}

func TestComprehensiveComplexDynamicTuples(t *testing.T) {
	users := []Tuple_e9afb3e4{
		{
			Id: big.NewInt(1),
			Profile: Tuple_8a486b93{
				Name:   "User 1",
				Emails: []string{"user1@example.com", "user1@gmail.com"},
				Metadata: Tuple_dc8f1c28{
					CreatedAt: big.NewInt(1234567890),
					Tags:      []string{"tag1", "tag2", "tag3"},
				},
			},
		},
		{
			Id: big.NewInt(2),
			Profile: Tuple_8a486b93{
				Name:   "User 2",
				Emails: []string{"user2@example.com"},
				Metadata: Tuple_dc8f1c28{
					CreatedAt: big.NewInt(9876543210),
					Tags:      []string{"tag4"},
				},
			},
		},
	}

	args := &TestComplexDynamicTuplesArgs{
		Users: users,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveABIDef.Pack("testComplexDynamicTuples",
		args.Users)
	require.NoError(t, err)

	require.Equal(t, len(encoded), len(goEthEncoded))
}

func TestComprehensiveMixedTypes(t *testing.T) {
	fixedData := [32]byte{0x01, 0x02, 0x03}
	dynamicData := []byte{0x04, 0x05, 0x06, 0x07}
	flag := true
	count := uint8(42)
	items := []Tuple_de3c4b6f{
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

	args := &TestMixedTypesArgs{
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
	goEthEncoded, err := ComprehensiveABIDef.Pack("testMixedTypes",
		args.FixedData, args.DynamicData, args.Flag, args.Count, args.Items)
	require.NoError(t, err)

	require.Equal(t, len(encoded), len(goEthEncoded))
}

func TestComprehensiveDeeplyNested(t *testing.T) {
	data := Tuple_5cee9471{
		Level1: Tuple_1064c5d1{
			Level2: Tuple_e05eda74{
				Level3: Tuple_54b20b3a{
					Value:       big.NewInt(999),
					Description: "Deeply nested value",
				},
			},
		},
	}

	args := &TestDeeplyNestedArgs{
		Data: data,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ComprehensiveABIDef.Pack("testDeeplyNested",
		args.Data)
	require.NoError(t, err)

	require.Equal(t, len(encoded), len(goEthEncoded))
}
