package tests

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEventEncodingDecoding(t *testing.T) {
	t.Run("Transfer event", func(t *testing.T) {
		// Create a Transfer event
		transfer := Transfer{
			From:  common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F2"),
			To:    common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F3"),
			Value: big.NewInt(1000000000000000000), // 1 ETH
		}

		// Encode topics
		topics, err := transfer.EncodeTopics()
		if err != nil {
			t.Fatalf("Failed to encode topics: %v", err)
		}

		// Verify topics count
		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics, got %d", len(topics))
		}

		// Verify first topic is event signature
		if topics[0] != TransferTopic {
			t.Fatalf("First topic should be event signature")
		}

		// Decode topics back
		var decodedTransfer Transfer
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
		complexEvent := ComplexEvent{
			Message: "Test message",
			Numbers: []*big.Int{big.NewInt(100), big.NewInt(200), big.NewInt(300)},
			Sender:  common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F4"),
		}

		// Encode topics
		topics, err := complexEvent.EncodeTopics()
		if err != nil {
			t.Fatalf("Failed to encode topics: %v", err)
		}

		// Verify topics count
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics, got %d", len(topics))
		}

		// Verify first topic is event signature
		if topics[0] != ComplexEventTopic {
			t.Fatalf("First topic should be event signature")
		}

		// Decode topics back
		var decodedComplexEvent ComplexEvent
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
		userCreated := UserCreated{
			User: User{
				Address: common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F5"),
				Name:    "Test User",
				Age:     big.NewInt(25),
			},
			Creator: common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F6"),
		}

		// Encode topics
		topics, err := userCreated.EncodeTopics()
		if err != nil {
			t.Fatalf("Failed to encode topics: %v", err)
		}

		// Verify topics count
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics, got %d", len(topics))
		}

		// Verify first topic is event signature
		if topics[0] != UserCreatedTopic {
			t.Fatalf("First topic should be event signature")
		}

		// Decode topics back
		var decodedUserCreated UserCreated
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
	t.Run("TransferData encoding", func(t *testing.T) {
		// Create TransferData
		data := TransferData{
			Value: big.NewInt(1000000000000000000),
		}

		// Encode data
		encoded, err := data.Encode()
		if err != nil {
			t.Fatalf("Failed to encode data: %v", err)
		}

		// Decode data back
		var decodedData TransferData
		if err := decodedData.Decode(encoded); err != nil {
			t.Fatalf("Failed to decode data: %v", err)
		}

		// Verify decoded values
		if decodedData.Value.Cmp(data.Value) != 0 {
			t.Errorf("Value mismatch: got %s, want %s", decodedData.Value, data.Value)
		}
	})

	t.Run("ComplexEventData encoding", func(t *testing.T) {
		// Create ComplexEventData
		data := ComplexEventData{
			Message: "Test message for encoding",
			Numbers: []*big.Int{big.NewInt(100), big.NewInt(200), big.NewInt(300)},
		}

		// Encode data
		encoded, err := data.Encode()
		if err != nil {
			t.Fatalf("Failed to encode data: %v", err)
		}

		// Decode data back
		var decodedData ComplexEventData
		if err := decodedData.Decode(encoded); err != nil {
			t.Fatalf("Failed to decode data: %v", err)
		}

		// Verify decoded values
		if decodedData.Message != data.Message {
			t.Errorf("Message mismatch: got %s, want %s", decodedData.Message, data.Message)
		}

		if len(decodedData.Numbers) != len(data.Numbers) {
			t.Fatalf("Numbers length mismatch: got %d, want %d", len(decodedData.Numbers), len(data.Numbers))
		}

		for i, num := range data.Numbers {
			if decodedData.Numbers[i].Cmp(num) != 0 {
				t.Errorf("Number[%d] mismatch: got %s, want %s", i, decodedData.Numbers[i], num)
			}
		}
	})

	t.Run("UserCreatedData encoding", func(t *testing.T) {
		// Create UserCreatedData
		data := UserCreatedData{
			User: User{
				Address: common.HexToAddress("0x742d35Cc6634C0532925a3b8Dc9F2a5C3B8Dc9F7"),
				Name:    "Test User Name",
				Age:     big.NewInt(30),
			},
		}

		// Encode data
		encoded, err := data.Encode()
		if err != nil {
			t.Fatalf("Failed to encode data: %v", err)
		}

		// Decode data back
		var decodedData UserCreatedData
		if err := decodedData.Decode(encoded); err != nil {
			t.Fatalf("Failed to decode data: %v", err)
		}

		// Verify decoded values
		if decodedData.User.Address != data.User.Address {
			t.Errorf("User address mismatch: got %s, want %s", decodedData.User.Address, data.User.Address)
		}
		if decodedData.User.Name != data.User.Name {
			t.Errorf("User name mismatch: got %s, want %s", decodedData.User.Name, data.User.Name)
		}
		if decodedData.User.Age.Cmp(data.User.Age) != 0 {
			t.Errorf("User age mismatch: got %s, want %s", decodedData.User.Age, data.User.Age)
		}
	})
}