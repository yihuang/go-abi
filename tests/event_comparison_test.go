package tests

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/test-go/testify/require"
)

func TestEventTopicComparison(t *testing.T) {
	// Test event signature calculation
	t.Run("Event signature calculation", func(t *testing.T) {
		// Calculate event signatures using go-ethereum
		transferSig := "Transfer(address,address,uint256)"
		complexEventSig := "Complex(string,uint256[],address)"
		userCreatedSig := "UserCreated((address,string,uint256),address)"

		// Calculate expected topic hashes
		expectedTransferTopic := crypto.Keccak256Hash([]byte(transferSig))
		expectedComplexEventTopic := crypto.Keccak256Hash([]byte(complexEventSig))
		expectedUserCreatedTopic := crypto.Keccak256Hash([]byte(userCreatedSig))

		// Compare with our generated topics
		if TransferEventTopic != expectedTransferTopic {
			t.Errorf("Transfer topic mismatch:\nGot:  %x\nWant: %x", TransferEventTopic, expectedTransferTopic)
		}

		if ComplexEventTopic != expectedComplexEventTopic {
			t.Errorf("ComplexEvent topic mismatch:\nGot:  %x\nWant: %x", ComplexEventTopic, expectedComplexEventTopic)
		}

		if UserCreatedEventTopic != expectedUserCreatedTopic {
			t.Errorf("UserCreated topic mismatch:\nGot:  %x\nWant: %x", UserCreatedEventTopic, expectedUserCreatedTopic)
		}
	})

	t.Run("Address encoding in topics", func(t *testing.T) {
		// Test that addresses are properly padded in topics
		address := common.HexToAddress("0x1234567890123456789012345678901234567890")

		// Create a simple event with just one indexed address
		event := TransferEventIndexed{
			From: address,
			To:   common.Address{}, // Zero address
		}

		// Encode topics
		topics, err := event.EncodeTopics()
		require.NoError(t, err)

		// Check that the address is properly padded (12 bytes of zeros + 20 bytes of address)
		fromTopic := topics[1]
		expectedFromTopic := [32]byte{}
		copy(expectedFromTopic[12:], address[:])

		if fromTopic != expectedFromTopic {
			t.Errorf("Address encoding mismatch:\nGot:  %x\nWant: %x", fromTopic, expectedFromTopic)
		}
	})
}

func TestEventDataEncodingComparison(t *testing.T) {
	t.Run("Transfer data encoding", func(t *testing.T) {
		// Create TransferData
		data := TransferEventData{
			Value: big.NewInt(1000000000000000000),
		}

		// Encode using our implementation
		ourEncoded, err := data.Encode()
		if err != nil {
			t.Fatalf("Failed to encode data: %v", err)
		}

		// Test that our encoding matches expected ABI encoding for uint256
		// A uint256 should be 32 bytes in ABI encoding
		if len(ourEncoded) != 32 {
			t.Fatalf("Data length mismatch: got %d, want %d", len(ourEncoded), 32)
		}

		// Verify the encoding matches the big.Int bytes (right-padded to 32 bytes)
		expectedBytes := data.Value.Bytes()
		if len(expectedBytes) > 32 {
			t.Fatalf("Value too large for uint256")
		}

		// Check that the encoding is correct (big-endian, right-padded)
		for i := 0; i < 32-len(expectedBytes); i++ {
			if ourEncoded[i] != 0 {
				t.Errorf("Padding byte %d mismatch: got %02x, want 00", i, ourEncoded[i])
			}
		}

		for i := 0; i < len(expectedBytes); i++ {
			if ourEncoded[32-len(expectedBytes)+i] != expectedBytes[i] {
				t.Errorf("Data byte %d mismatch: got %02x, want %02x", i, ourEncoded[32-len(expectedBytes)+i], expectedBytes[i])
			}
		}
	})
}
