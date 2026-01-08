//go:build !uint256

package tests

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEventIndexedEncodingDecoding(t *testing.T) {
	t.Run("Transfer event", func(t *testing.T) {
		// Create a Transfer event
		transfer := NewTransferEvent(
			common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F2"),
			common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F3"),
			big.NewInt(1000000000000000000),
		)
		EventDecodeRoundTrip(t, transfer)
	})

	t.Run("ComplexEvent event", func(t *testing.T) {
		// Create a ComplexEvent
		complexEvent := NewComplexEvent(
			"hello world",
			[]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
			common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F4"),
		)
		EventDecodeRoundTrip(t, complexEvent)
	})

	t.Run("UserCreated event", func(t *testing.T) {
		// Create a UserCreated event
		userCreated := NewUserCreatedEvent(
			User{
				Address: common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F5"),
				Name:    "Alice",
				Age:     big.NewInt(28),
			},
			common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F6"),
		)
		EventDecodeRoundTrip(t, userCreated)
	})
}
