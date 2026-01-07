//go:build uint256

package tests

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/yihuang/go-abi"
)

func newUint256(v uint64) *uint256.Int {
	return uint256.NewInt(v)
}

func newUint256Max() *uint256.Int {
	max := new(uint256.Int)
	max.SetAllOne()
	return max
}

func newUint256Slice(count int) []*uint256.Int {
	result := make([]*uint256.Int, count)
	for i := 0; i < count; i++ {
		result[i] = uint256.NewInt(uint64(i * 1000000000000000000))
	}
	return result
}

func BenchmarkUint256_Transfer_Encode(b *testing.B) {
	call := NewTransferCall(testAddress, newUint256(1000000000000000000))
	benchEncode(b, call)
}

func BenchmarkUint256_Transfer_EncodeTo(b *testing.B) {
	call := NewTransferCall(testAddress, newUint256(1000000000000000000))
	benchEncodeTo(b, call)
}

func BenchmarkUint256_Transfer_Decode(b *testing.B) {
	call := NewTransferCall(testAddress, newUint256(1000000000000000000))
	encoded, _ := call.Encode()
	benchDecode(b, encoded, func() abi.Decode { return &TransferCall{} })
}

func BenchmarkUint256_MultiTransfer_Encode(b *testing.B) {
	recipients := make([]common.Address, 10)
	for i := range recipients {
		recipients[i] = testAddress
	}
	call := NewMultiTransferCall(recipients, newUint256Slice(10))
	benchEncode(b, call)
}

func BenchmarkUint256_MultiTransfer_EncodeTo(b *testing.B) {
	recipients := make([]common.Address, 10)
	for i := range recipients {
		recipients[i] = testAddress
	}
	call := NewMultiTransferCall(recipients, newUint256Slice(10))
	benchEncodeTo(b, call)
}

func BenchmarkUint256_MultiTransfer_Decode(b *testing.B) {
	recipients := make([]common.Address, 10)
	for i := range recipients {
		recipients[i] = testAddress
	}
	call := NewMultiTransferCall(recipients, newUint256Slice(10))
	encoded, _ := call.Encode()
	benchDecode(b, encoded, func() abi.Decode { return &MultiTransferCall{} })
}

func BenchmarkUint256_LargeValue_Encode(b *testing.B) {
	call := NewTransferCall(testAddress, newUint256Max())
	benchEncode(b, call)
}

func BenchmarkUint256_LargeValue_Decode(b *testing.B) {
	call := NewTransferCall(testAddress, newUint256Max())
	encoded, _ := call.Encode()
	benchDecode(b, encoded, func() abi.Decode { return &TransferCall{} })
}
