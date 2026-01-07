//go:build uint256

package tests

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

func TestUint256Transfer(t *testing.T) {
	// Test that uint256.Int is used for uint256 types
	amount := uint256.NewInt(1000000000000000000) // 1 ETH in wei
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	call := NewTransferCall(to, amount)

	// Encode
	encoded, err := call.EncodeWithSelector()
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	// Verify selector
	if !bytes.Equal(encoded[:4], TransferSelector[:]) {
		t.Errorf("Selector mismatch: got %x, want %x", encoded[:4], TransferSelector)
	}

	// Decode
	var decoded TransferCall
	_, err = decoded.Decode(encoded[4:])
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	// Verify values
	if decoded.To != to {
		t.Errorf("To mismatch: got %s, want %s", decoded.To.Hex(), to.Hex())
	}
	if decoded.Amount.Cmp(amount) != 0 {
		t.Errorf("Amount mismatch: got %s, want %s", decoded.Amount.String(), amount.String())
	}
}

func TestUint256BalanceOf(t *testing.T) {
	account := common.HexToAddress("0xabcdef0123456789abcdef0123456789abcdef01")

	call := NewBalanceOfCall(account)

	// Encode and verify no error
	_, err := call.EncodeWithSelector()
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	// Test return decoding
	balance := uint256.NewInt(999999999999999999)
	var ret BalanceOfReturn
	ret.Field1 = balance

	retEncoded, err := ret.Encode()
	if err != nil {
		t.Fatalf("Failed to encode return: %v", err)
	}

	var retDecoded BalanceOfReturn
	_, err = retDecoded.Decode(retEncoded)
	if err != nil {
		t.Fatalf("Failed to decode return: %v", err)
	}

	if retDecoded.Field1.Cmp(balance) != 0 {
		t.Errorf("Balance mismatch: got %s, want %s", retDecoded.Field1.String(), balance.String())
	}
}

func TestUint256SliceEncoding(t *testing.T) {
	recipients := []common.Address{
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
	}
	amounts := []*uint256.Int{
		uint256.NewInt(100),
		uint256.NewInt(200),
	}

	call := NewMultiTransferCall(recipients, amounts)

	// Encode
	encoded, err := call.EncodeWithSelector()
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	// Decode
	var decoded MultiTransferCall
	_, err = decoded.Decode(encoded[4:])
	if err != nil {
		t.Fatalf("Failed to decode: %v", err)
	}

	// Verify recipients
	if len(decoded.Recipients) != len(recipients) {
		t.Fatalf("Recipients length mismatch: got %d, want %d", len(decoded.Recipients), len(recipients))
	}
	for i := range recipients {
		if decoded.Recipients[i] != recipients[i] {
			t.Errorf("Recipient[%d] mismatch: got %s, want %s", i, decoded.Recipients[i].Hex(), recipients[i].Hex())
		}
	}

	// Verify amounts
	if len(decoded.Amounts) != len(amounts) {
		t.Fatalf("Amounts length mismatch: got %d, want %d", len(decoded.Amounts), len(amounts))
	}
	for i := range amounts {
		if decoded.Amounts[i].Cmp(amounts[i]) != 0 {
			t.Errorf("Amount[%d] mismatch: got %s, want %s", i, decoded.Amounts[i].String(), amounts[i].String())
		}
	}
}

func TestUint256EquivalentToBigInt(t *testing.T) {
	// Test that encoding produces identical bytes whether using big.Int or uint256.Int
	value := uint256.NewInt(123456789012345678)
	bigValue := new(big.Int).SetUint64(123456789012345678)

	// Encode uint256
	buf := make([]byte, 32)
	value.WriteToArray32((*[32]byte)(buf))

	// Encode big.Int
	bigBuf := make([]byte, 32)
	bigValue.FillBytes(bigBuf)

	if !bytes.Equal(buf, bigBuf) {
		t.Errorf("Encoding mismatch:\nuint256: %x\nbig.Int: %x", buf, bigBuf)
	}
}

func TestUint256MaxValue(t *testing.T) {
	// Test with max uint256 value
	maxValue := new(uint256.Int)
	maxValue.SetAllOne()

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	call := NewTransferCall(to, maxValue)

	encoded, err := call.EncodeWithSelector()
	if err != nil {
		t.Fatalf("Failed to encode max value: %v", err)
	}

	var decoded TransferCall
	_, err = decoded.Decode(encoded[4:])
	if err != nil {
		t.Fatalf("Failed to decode max value: %v", err)
	}

	if decoded.Amount.Cmp(maxValue) != 0 {
		t.Errorf("Max value mismatch: got %s, want %s", decoded.Amount.String(), maxValue.String())
	}
}
