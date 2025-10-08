package abi

import (
	"math/big"
	"testing"
)

func TestUint256Encoding(t *testing.T) {
	uint256 := &Uint256{}

	// Test encoding
	value := big.NewInt(42)
	encoded, err := uint256.Encode(value)
	if err != nil {
		t.Fatalf("Encoding failed: %v", err)
	}

	if len(encoded) != WordSize {
		t.Fatalf("Expected encoded length %d, got %d", WordSize, len(encoded))
	}

	// Test decoding
	var decoded big.Int
	bytesRead, err := uint256.Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}

	if bytesRead != WordSize {
		t.Fatalf("Expected %d bytes read, got %d", WordSize, bytesRead)
	}

	if decoded.Cmp(value) != 0 {
		t.Fatalf("Expected %v, got %v", value, &decoded)
	}
}

func TestBoolEncoding(t *testing.T) {
	boolType := &Bool{}

	// Test true encoding
	encodedTrue, err := boolType.Encode(true)
	if err != nil {
		t.Fatalf("Encoding true failed: %v", err)
	}

	var decodedTrue bool
	bytesRead, err := boolType.Decode(encodedTrue, &decodedTrue)
	if err != nil {
		t.Fatalf("Decoding true failed: %v", err)
	}

	if bytesRead != WordSize {
		t.Fatalf("Expected %d bytes read, got %d", WordSize, bytesRead)
	}

	if !decodedTrue {
		t.Fatal("Expected true, got false")
	}

	// Test false encoding
	encodedFalse, err := boolType.Encode(false)
	if err != nil {
		t.Fatalf("Encoding false failed: %v", err)
	}

	var decodedFalse bool
	bytesRead, err = boolType.Decode(encodedFalse, &decodedFalse)
	if err != nil {
		t.Fatalf("Decoding false failed: %v", err)
	}

	if decodedFalse {
		t.Fatal("Expected false, got true")
	}
}

func TestAddressEncoding(t *testing.T) {
	addrType := &Address{}

	addr := [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	encoded, err := addrType.Encode(addr)
	if err != nil {
		t.Fatalf("Encoding failed: %v", err)
	}

	if len(encoded) != WordSize {
		t.Fatalf("Expected encoded length %d, got %d", WordSize, len(encoded))
	}

	var decoded [20]byte
	bytesRead, err := addrType.Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}

	if bytesRead != WordSize {
		t.Fatalf("Expected %d bytes read, got %d", WordSize, bytesRead)
	}

	if decoded != addr {
		t.Fatalf("Expected %v, got %v", addr, decoded)
	}
}

func TestBytesEncoding(t *testing.T) {
	bytesType := &Bytes{}

	testBytes := []byte("hello world")
	encoded, err := bytesType.Encode(testBytes)
	if err != nil {
		t.Fatalf("Encoding failed: %v", err)
	}

	// For dynamic types, encoded length should be at least WordSize + data length
	expectedMinLength := WordSize + len(testBytes)
	if len(encoded) < expectedMinLength {
		t.Fatalf("Expected at least %d bytes, got %d", expectedMinLength, len(encoded))
	}

	var decoded []byte
	bytesRead, err := bytesType.Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}

	if string(decoded) != string(testBytes) {
		t.Fatalf("Expected %s, got %s", testBytes, decoded)
	}

	// Verify bytesRead accounts for padding
	expectedBytesRead := WordSize + ((uint(len(testBytes)) + WordSize - 1) / WordSize * WordSize)
	if bytesRead != expectedBytesRead {
		t.Fatalf("Expected %d bytes read, got %d", expectedBytesRead, bytesRead)
	}
}

func TestStringEncoding(t *testing.T) {
	stringType := &String{}

	testString := "hello world"
	encoded, err := stringType.Encode(testString)
	if err != nil {
		t.Fatalf("Encoding failed: %v", err)
	}

	var decoded string
	bytesRead, err := stringType.Decode(encoded, &decoded)
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}

	if decoded != testString {
		t.Fatalf("Expected %s, got %s", testString, decoded)
	}

	// Verify bytesRead accounts for padding
	expectedBytesRead := WordSize + ((uint(len(testString)) + WordSize - 1) / WordSize * WordSize)
	if bytesRead != expectedBytesRead {
		t.Fatalf("Expected %d bytes read, got %d", expectedBytesRead, bytesRead)
	}
}