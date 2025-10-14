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
	//go:embed test_abi.json
	ABIJson string
	ABIDef  abi.ABI
)

func init() {
	var err error
	ABIDef, err = abi.JSON(strings.NewReader(ABIJson))
	if err != nil {
		panic(err)
	}
}

//go:generate go run github.com/yihuang/go-abi/cmd -input test_abi.json -output test_abi.abi.go

func TestTransferEncoding(t *testing.T) {
	to := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9D7B6f7e5c3a3")
	amount := big.NewInt(1000)

	// Get our generated encoding
	args := &TransferArgs{
		To:     to,
		Amount: amount,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ABIDef.Pack("transfer", to, amount)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)
}

func TestSetMessageEncoding(t *testing.T) {
	message := "Hello, World!"

	// Get our generated encoding
	args := &SetMessageArgs{
		Message: message,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ABIDef.Pack("setMessage", args.Message)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)
}

func TestUpdateProfileEncoding(t *testing.T) {
	user := common.HexToAddress("0x1234567890123456789012345678901234567890")
	name := "Test User"
	age := big.NewInt(25)

	// Get our generated encoding
	args := &UpdateProfileArgs{
		User: user,
		Name: name,
		Age:  age,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ABIDef.Pack("updateProfile", args.User, args.Name, args.Age)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)
}

func TestProcessUserDataEncoding(t *testing.T) {
	user := Tuple_b53c1574{
		Address: common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Name:    "Test User",
		Age:     big.NewInt(25),
	}

	// Get our generated encoding
	args := &ProcessUserDataArgs{
		User: user,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ABIDef.Pack("processUserData", user)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)
}

func TestBatchProcessEncoding(t *testing.T) {
	users := []Tuple_1821f6d7{
		{
			Id: big.NewInt(1),
			Data: Tuple_4a9d2179{
				Key:   [32]byte{0x01, 0x02, 0x03},
				Value: "First user",
			},
		},
		{
			Id: big.NewInt(2),
			Data: Tuple_4a9d2179{
				Key:   [32]byte{0x04, 0x05, 0x06},
				Value: "Second user",
			},
		},
	}

	// Get our generated encoding
	args := &BatchProcessArgs{
		Users: users,
	}

	// Test encoding with selector
	encoded, err := args.EncodeWithSelector()
	require.NoError(t, err)

	// Get go-ethereum encoding
	goEthEncoded, err := ABIDef.Pack("batchProcess", users)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)
}

func TestSmallIntegers(t *testing.T) {
	// Get our generated encoding
	args := &SmallIntegersArgs{
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
	goEthEncoded, err := ABIDef.Pack("smallIntegers",
		args.U8, args.U16, args.U32, args.U64,
		args.I8, args.I16, args.I32, args.I64,
	)
	require.NoError(t, err)

	require.Equal(t, encoded, goEthEncoded)
}
