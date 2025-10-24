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

//go:generate go run ../cmd -var TestABI -output test.abi.go

// TestABI contains human-readable ABI definitions for testing
var TestABI = []string{
	"function transfer(address to, uint256 amount) returns (bool)",
	"function balanceOf(address account) view returns (uint256)",
	"function setMessage(string message) returns (bool)",
	"function updateProfile(address user, string name, uint256 age) returns (bool)",
	"function transferBatch(address[] recipients, uint256[] amounts) returns (bool)",
	"function setData(bytes32 key, bytes value)",
	"function getBalances(address[10] accounts) view returns (uint256[10])",
	"struct User { address address; string name; uint256 age }",
	"function processUserData(User user) returns (bool)",
	"struct UserMetadata { bytes32 key; string value }",
	"struct UserData { uint256 id; UserMetadata data }",
	"function batchProcess(UserData[] users) returns (bool)",
	"function smallIntegers(uint8 u8, uint16 u16, uint32 u32, uint64 u64, int8 i8, int16 i16, int32 i32, int64 i64) returns (bool)",
	"function communityPool() view returns ((string denom, uint256 amount)[] coins)",
}

var TestABIDef ethabi.ABI

func init() {
	var err error
	abiJSON, err := abi.ParseHumanReadableABI(TestABI)
	if err != nil {
		panic(err)
	}
	TestABIDef, err = ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		panic(err)
	}
}

func TestTransferEncoding(t *testing.T) {
	to := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9D7B6f7e5c3a3")
	amount := big.NewInt(1000)

	// Get our generated encoding
	args := &TransferCall{
		To:     to,
		Amount: amount,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := TestABIDef.Pack("transfer", to, amount)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestSetMessageEncoding(t *testing.T) {
	message := "Hello, World!"

	// Get our generated encoding
	args := &SetMessageCall{
		Message: message,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := TestABIDef.Pack("setMessage", args.Message)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestUpdateProfileEncoding(t *testing.T) {
	user := common.HexToAddress("0x1234567890123456789012345678901234567890")
	name := "Test User"
	age := big.NewInt(25)

	// Get our generated encoding
	args := &UpdateProfileCall{
		User: user,
		Name: name,
		Age:  age,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := TestABIDef.Pack("updateProfile", args.User, args.Name, args.Age)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestProcessUserDataEncoding(t *testing.T) {
	user := User{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Name:    "Test User",
		Age:     big.NewInt(25),
	}

	// Get our generated encoding
	args := &ProcessUserDataCall{
		User: user,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := TestABIDef.Pack("processUserData", user)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestBatchProcessEncoding(t *testing.T) {
	users := []UserData{
		{
			Id: big.NewInt(1),
			Data: UserMetadata{
				Key:   [32]byte{0x01, 0x02, 0x03},
				Value: "First user",
			},
		},
		{
			Id: big.NewInt(2),
			Data: UserMetadata{
				Key:   [32]byte{0x04, 0x05, 0x06},
				Value: "Second user",
			},
		},
	}

	// Get our generated encoding
	args := &BatchProcessCall{
		Users: users,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := TestABIDef.Pack("batchProcess", users)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}

func TestSmallIntegers(t *testing.T) {
	// Get our generated encoding
	args := &SmallIntegersCall{
		U8:  uint8(10),
		U16: uint16(1000),
		U32: uint32(100000),
		U64: uint64(10000000000),
		I8:  int8(-10),
		I16: int16(-1000),
		I32: int32(-100000),
		I64: int64(-10000000000),
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := TestABIDef.Pack("smallIntegers",
		args.U8, args.U16, args.U32, args.U64,
		args.I8, args.I16, args.I32, args.I64,
	)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)

	DecodeRoundTrip(t, args)
}
