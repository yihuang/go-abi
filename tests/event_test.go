package tests

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
)

func TestEventIndexedEncodingDecoding(t *testing.T) {
	t.Run("Transfer event", func(t *testing.T) {
		// Create a Transfer event
		transfer := NewTransferEvent(
			common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F2"),
			common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F3"),
			big.NewInt(1000000000000000000),
		)

		// Encode topics
		topics, err := transfer.EncodeTopics()
		require.NoError(t, err)

		// Verify topics count
		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics, got %d", len(topics))
		}

		// Verify first topic is event signature
		if topics[0] != TransferEventTopic {
			t.Fatalf("First topic should be event signature")
		}

		// Decode topics back
		var decodedTransfer TransferEventIndexed
		if err := decodedTransfer.DecodeTopics(topics); err != nil {
			t.Fatalf("Failed to decode topics: %v", err)
		}

		// Verify decoded values
		if decodedTransfer.From != transfer.From {
			t.Errorf("From address mismatch: got %s, want %s", decodedTransfer.From, transfer.From)
		}
		if decodedTransfer.To != transfer.To {
			t.Errorf("To address mismatch: got %s, want %s", decodedTransfer.To, transfer.To)
		}
		// Note: Value is not in topics, it's in the data
	})

	t.Run("ComplexEvent event", func(t *testing.T) {
		// Create a ComplexEvent
		complexEvent := ComplexEventIndexed{
			Sender: common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F4"),
		}

		// Encode topics
		topics, err := complexEvent.EncodeTopics()
		require.NoError(t, err)

		// Verify topics count
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics, got %d", len(topics))
		}

		// Verify first topic is event signature
		if topics[0] != ComplexEventTopic {
			t.Fatalf("First topic should be event signature")
		}

		// Decode topics back
		var decodedComplexEvent ComplexEventIndexed
		if err := decodedComplexEvent.DecodeTopics(topics); err != nil {
			t.Fatalf("Failed to decode topics: %v", err)
		}

		// Verify decoded values
		if decodedComplexEvent.Sender != complexEvent.Sender {
			t.Errorf("Sender address mismatch: got %s, want %s", decodedComplexEvent.Sender, complexEvent.Sender)
		}
	})

	t.Run("UserCreated event", func(t *testing.T) {
		// Create a UserCreated event
		userCreated := UserCreatedEventIndexed{
			Creator: common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F6"),
		}

		// Encode topics
		topics, err := userCreated.EncodeTopics()
		require.NoError(t, err)

		// Verify topics count
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics, got %d", len(topics))
		}

		// Verify first topic is event signature
		if topics[0] != UserCreatedEventTopic {
			t.Fatalf("First topic should be event signature")
		}

		// Decode topics back
		var decodedUserCreated UserCreatedEventIndexed
		if err := decodedUserCreated.DecodeTopics(topics); err != nil {
			t.Fatalf("Failed to decode topics: %v", err)
		}

		// Verify decoded values
		if decodedUserCreated.Creator != userCreated.Creator {
			t.Errorf("Creator address mismatch: got %s, want %s", decodedUserCreated.Creator, userCreated.Creator)
		}
	})
}

func TestEventDataEncoding(t *testing.T) {
	// Create TransferData
	DecodeRoundTrip(t, &TransferEventData{
		Value: big.NewInt(1000000000000000000),
	})

	// Create ComplexEventData
	DecodeRoundTrip(t, &ComplexEventData{
		Message: "Test message for encoding",
		Numbers: []*big.Int{big.NewInt(100), big.NewInt(200), big.NewInt(300)},
	})

	// Create UserCreatedData
	DecodeRoundTrip(t, &UserCreatedEventData{
		User: User{
			Address: common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F7"),
			Name:    "Test User Name",
			Age:     big.NewInt(30),
		},
	})
}
