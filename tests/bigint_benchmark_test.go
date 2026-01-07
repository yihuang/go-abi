//go:build !uint256

package tests

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yihuang/go-abi"
)

func newBigInt(v int64) *big.Int {
	return big.NewInt(v)
}

func newBigIntMax() *big.Int {
	amount := new(big.Int)
	amount.SetString("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 16)
	return amount
}

func newBigIntSlice(count int) []*big.Int {
	result := make([]*big.Int, count)
	for i := 0; i < count; i++ {
		result[i] = big.NewInt(int64(i * 1000000000000000000))
	}
	return result
}

func BenchmarkBigInt_Transfer_Encode(b *testing.B) {
	call := NewTransferCall(testAddress, newBigInt(1000000000000000000))
	benchEncode(b, call)
}

func BenchmarkBigInt_Transfer_EncodeTo(b *testing.B) {
	call := NewTransferCall(testAddress, newBigInt(1000000000000000000))
	benchEncodeTo(b, call)
}

func BenchmarkBigInt_Transfer_Decode(b *testing.B) {
	call := NewTransferCall(testAddress, newBigInt(1000000000000000000))
	encoded, _ := call.Encode()
	benchDecode(b, encoded, func() abi.Decode { return &TransferCall{} })
}

func BenchmarkBigInt_TransferBatch_Encode(b *testing.B) {
	recipients := make([]common.Address, 10)
	for i := range recipients {
		recipients[i] = testAddress
	}
	call := NewTransferBatchCall(recipients, newBigIntSlice(10))
	benchEncode(b, call)
}

func BenchmarkBigInt_TransferBatch_EncodeTo(b *testing.B) {
	recipients := make([]common.Address, 10)
	for i := range recipients {
		recipients[i] = testAddress
	}
	call := NewTransferBatchCall(recipients, newBigIntSlice(10))
	benchEncodeTo(b, call)
}

func BenchmarkBigInt_TransferBatch_Decode(b *testing.B) {
	recipients := make([]common.Address, 10)
	for i := range recipients {
		recipients[i] = testAddress
	}
	call := NewTransferBatchCall(recipients, newBigIntSlice(10))
	encoded, _ := call.Encode()
	benchDecode(b, encoded, func() abi.Decode { return &TransferBatchCall{} })
}

func BenchmarkBigInt_LargeValue_Encode(b *testing.B) {
	call := NewTransferCall(testAddress, newBigIntMax())
	benchEncode(b, call)
}

func BenchmarkBigInt_LargeValue_Decode(b *testing.B) {
	call := NewTransferCall(testAddress, newBigIntMax())
	encoded, _ := call.Encode()
	benchDecode(b, encoded, func() abi.Decode { return &TransferCall{} })
}
