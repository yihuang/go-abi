package abi

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
)

var ABISpec abi.ABI

/* Human readable abi
[
  'struct Tuple{uint256 field0;}',
  'struct DynamicTuple{string field0; uint256 field1;}',
  'function transfer(string memo, address to, uint256 amount, uint256[2] staticArray, Tuple[] dynamicArray, DynamicTuple[][] dynamicArray2, int256 negative)'
]
*/

func init() {
	var err error
	ABISpec, err = abi.JSON(strings.NewReader(`[
  {
    "type": "function",
    "name": "transfer",
    "stateMutability": "nonpayable",
    "inputs": [
      {
        "name": "memo",
        "type": "string"
      },
      {
        "name": "to",
        "type": "address"
      },
      {
        "name": "amount",
        "type": "uint256"
      },
      {
        "name": "staticArray",
        "type": "uint256[2]"
      },
      {
        "name": "dynamicArray",
        "type": "tuple[]",
        "internalType": "struct Tuple[]",
        "components": [
          {
            "name": "field0",
            "type": "uint256"
          }
        ]
      },
      {
        "name": "dynamicArray2",
        "type": "tuple[][]",
        "internalType": "struct DynamicTuple[][]",
        "components": [
          {
            "name": "field0",
            "type": "string"
          },
          {
            "name": "field1",
            "type": "uint256"
          }
        ]
      },
      {
        "name": "negative",
        "type": "int256"
      }
    ],
    "outputs": []
  }
]`))
	if err != nil {
		panic(err)
	}
}

func TestAgainstGeth(t *testing.T) {
	addr := common.BigToAddress(big.NewInt(1))
	amount := big.NewInt(1000)
	negative := big.NewInt(-100)
	memo := "hello world"
	staticArray := [2]*big.Int{big.NewInt(10), big.NewInt(20)}
	dynamicArray := []Tuple0{{big.NewInt(5)}, {big.NewInt(15)}}
	dynamicArray2 := [][]DynamicTuple{
		{{"dynamic1", big.NewInt(500)}, {"dynamic2", big.NewInt(600)}},
		{{"dynamic1", big.NewInt(500)}, {"dynamic2", big.NewInt(600)}},
	}

	ethbz, err := ABISpec.Pack("transfer",
		memo,
		addr,
		amount,
		staticArray,
		dynamicArray,
		dynamicArray2,
		negative,
	)
	require.NoError(t, err)

	args := TransferCall{
		Memo:          memo,
		To:            addr,
		Amount:        amount,
		StaticArray:   staticArray,
		DynamicArray:  dynamicArray,
		DynamicArray2: dynamicArray2,
		Negative:      negative,
	}
	bz, err := args.Encode()
	require.NoError(t, err)

	require.Equal(t, hex.EncodeToString(ethbz[4:]), hex.EncodeToString(bz))
	// require.Equal(t, ethbz[4:], bz)
}
