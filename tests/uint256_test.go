//go:build uint256

package tests

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
)

func TestUint256Transfer(t *testing.T) {
	tests := []struct {
		name   string
		to     common.Address
		amount *uint256.Int
	}{
		{
			name:   "typical",
			to:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
			amount: uint256.NewInt(1000000000000000000),
		},
		{
			name:   "max",
			to:     common.HexToAddress("0xabcdef0123456789abcdef0123456789abcdef01"),
			amount: newUint256Max(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DecodeRoundTrip(t, &TransferCall{To: tt.to, Amount: tt.amount})
		})
	}
}

func TestUint256BalanceOf(t *testing.T) {
	DecodeRoundTrip(t, &BalanceOfCall{
		Account: common.HexToAddress("0xabcdef0123456789abcdef0123456789abcdef01"),
	})
}

func TestUint256BalanceOfReturn(t *testing.T) {
	tests := []struct {
		name  string
		value *uint256.Int
	}{
		{"typical", uint256.NewInt(999999999999999999)},
		{"max", newUint256Max()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DecodeRoundTrip(t, &BalanceOfReturn{Field1: tt.value})
		})
	}
}

func TestUint256MultiTransfer(t *testing.T) {
	maxValue := newUint256Max()
	halfMax := new(uint256.Int).Rsh(maxValue, 1)

	tests := []struct {
		name       string
		recipients []common.Address
		amounts    []*uint256.Int
	}{
		{
			name:       "empty",
			recipients: []common.Address{},
			amounts:    []*uint256.Int{},
		},
		{
			name: "typical",
			recipients: []common.Address{
				common.HexToAddress("0x1111111111111111111111111111111111111111"),
				common.HexToAddress("0x2222222222222222222222222222222222222222"),
				common.HexToAddress("0x3333333333333333333333333333333333333333"),
			},
			amounts: []*uint256.Int{
				uint256.NewInt(100),
				uint256.NewInt(200),
				uint256.NewInt(300),
			},
		},
		{
			name: "large_values",
			recipients: []common.Address{
				common.HexToAddress("0x1111111111111111111111111111111111111111"),
				common.HexToAddress("0x2222222222222222222222222222222222222222"),
			},
			amounts: []*uint256.Int{maxValue, halfMax},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DecodeRoundTrip(t, &MultiTransferCall{
				Recipients: tt.recipients,
				Amounts:    tt.amounts,
			})
		})
	}
}

func TestUint256TransferEvent(t *testing.T) {
	tests := []struct {
		name  string
		from  common.Address
		to    common.Address
		value *uint256.Int
	}{
		{
			name:  "typical",
			from:  common.HexToAddress("0x1111111111111111111111111111111111111111"),
			to:    common.HexToAddress("0x2222222222222222222222222222222222222222"),
			value: uint256.NewInt(1000000000000000000),
		},
		{
			name:  "max",
			from:  common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			to:    common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
			value: newUint256Max(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			EventDecodeRoundTrip(t, NewTransferEvent(tt.from, tt.to, tt.value))
		})
	}
}
