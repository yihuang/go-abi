//go:build !uint256

package tests

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/test-go/testify/require"
	"github.com/yihuang/go-abi"
)

//go:generate go run ../cmd -var PackedTestABI -output packed.abi.go -prefix packed -buildtag=!uint256

// PackedTestABI contains human-readable ABI definitions for packed encoding testing
var PackedTestABI = []string{
	"function packedTransfer(address to, uint256 amount) returns (bool)",
	"function packedSmallInts(uint8 u8, uint16 u16, uint32 u32, uint64 u64, int8 i8, int16 i16, int32 i32, int64 i64) returns (bool)",
	"function packedBytes(bytes32 b32, bytes4 b4) returns (bool)",
	"function packedBool(bool a, bool b) returns (bool)",
	"function packedIntermediate(uint24 u24, uint40 u40, int24 i24, int40 i40) returns (bool)",
	"struct PackedStruct { address addr; uint256 value; bytes32 data }",
	"function packedStruct(PackedStruct s) returns (bool)",
}

var PackedTestABIDef ethabi.ABI

func init() {
	var err error
	abiJSON, err := abi.ParseHumanReadableABI(PackedTestABI)
	if err != nil {
		panic(err)
	}
	PackedTestABIDef, err = ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		panic(err)
	}
}

// TestPackedTransfer tests packed encoding for address + uint256
func TestPackedTransfer(t *testing.T) {
	to := common.HexToAddress("0x742d35Cc6634C0532925a3b8D4C9D7B6f7e5c3a3")
	amount := big.NewInt(1000)

	call := &PackedTransferCall{
		To:     to,
		Amount: amount,
	}

	// Test PackedEncodedSize - address (20) + uint256 (32) = 52 bytes
	require.Equal(t, 52, call.PackedEncodedSize())

	// Test PackedEncode
	encoded, err := call.PackedEncode()
	require.NoError(t, err)
	require.Len(t, encoded, 52)

	// First 20 bytes should be the address
	require.Equal(t, to[:], encoded[:20])

	// Next 32 bytes should be the uint256 value (big-endian)
	expectedAmount := make([]byte, 32)
	amount.FillBytes(expectedAmount)
	require.Equal(t, expectedAmount, encoded[20:52])

	// Test round-trip
	DecodePackedRoundTrip(t, call)
}

// TestPackedSmallInts tests packed encoding for small integer types
func TestPackedSmallInts(t *testing.T) {
	call := &PackedSmallIntsCall{
		U8:  uint8(0xAB),
		U16: uint16(0xABCD),
		U32: uint32(0xABCDEF12),
		U64: uint64(0xABCDEF1234567890),
		I8:  int8(-10),
		I16: int16(-1000),
		I32: int32(-100000),
		I64: int64(-10000000000),
	}

	// Size: 1+2+4+8+1+2+4+8 = 30 bytes
	require.Equal(t, 30, call.PackedEncodedSize())

	encoded, err := call.PackedEncode()
	require.NoError(t, err)
	require.Len(t, encoded, 30)

	// Verify encoding
	require.Equal(t, byte(0xAB), encoded[0])
	require.Equal(t, []byte{0xAB, 0xCD}, encoded[1:3])
	require.Equal(t, []byte{0xAB, 0xCD, 0xEF, 0x12}, encoded[3:7])

	// Test round-trip
	DecodePackedRoundTrip(t, call)
}

// TestPackedBytes tests packed encoding for fixed bytes types
func TestPackedBytes(t *testing.T) {
	call := &PackedBytesCall{
		B32: [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
			17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
		B4: [4]byte{0xDE, 0xAD, 0xBE, 0xEF},
	}

	// Size: 32 + 4 = 36 bytes
	require.Equal(t, 36, call.PackedEncodedSize())

	encoded, err := call.PackedEncode()
	require.NoError(t, err)
	require.Len(t, encoded, 36)

	// Verify encoding
	require.Equal(t, call.B32[:], encoded[:32])
	require.Equal(t, call.B4[:], encoded[32:36])

	// Test round-trip
	DecodePackedRoundTrip(t, call)
}

// TestPackedBool tests packed encoding for boolean types
func TestPackedBool(t *testing.T) {
	testCases := []struct {
		a, b     bool
		expected []byte
	}{
		{false, false, []byte{0, 0}},
		{true, false, []byte{1, 0}},
		{false, true, []byte{0, 1}},
		{true, true, []byte{1, 1}},
	}

	for _, tc := range testCases {
		call := &PackedBoolCall{A: tc.a, B: tc.b}

		// Size: 1 + 1 = 2 bytes
		require.Equal(t, 2, call.PackedEncodedSize())

		encoded, err := call.PackedEncode()
		require.NoError(t, err)
		require.Equal(t, tc.expected, encoded)

		// Test round-trip
		DecodePackedRoundTrip(t, call)
	}
}

// TestPackedIntermediate tests packed encoding for intermediate-sized integers (24, 40, 48, 56 bits)
func TestPackedIntermediate(t *testing.T) {
	call := &PackedIntermediateCall{
		U24: uint32(0xABCDEF),     // 3 bytes
		U40: uint64(0xABCDEF1234), // 5 bytes
		I24: int32(-1234),         // 3 bytes
		I40: int64(-123456789),    // 5 bytes
	}

	// Size: 3+5+3+5 = 16 bytes
	require.Equal(t, 16, call.PackedEncodedSize())

	encoded, err := call.PackedEncode()
	require.NoError(t, err)
	require.Len(t, encoded, 16)

	// Verify U24 encoding (big-endian, 3 bytes)
	require.Equal(t, []byte{0xAB, 0xCD, 0xEF}, encoded[0:3])

	// Verify U40 encoding (big-endian, 5 bytes)
	require.Equal(t, []byte{0xAB, 0xCD, 0xEF, 0x12, 0x34}, encoded[3:8])

	// Test round-trip
	DecodePackedRoundTrip(t, call)
}

// TestPackedStruct tests packed encoding for struct types
func TestPackedStruct(t *testing.T) {
	s := PackedStruct{
		Addr:  common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Value: big.NewInt(999999),
		Data:  [32]byte{1, 2, 3, 4, 5, 6, 7, 8},
	}

	call := &PackedStructCall{S: s}

	// Size: 20 (address) + 32 (uint256) + 32 (bytes32) = 84 bytes
	require.Equal(t, 84, call.PackedEncodedSize())

	encoded, err := call.PackedEncode()
	require.NoError(t, err)
	require.Len(t, encoded, 84)

	// Test round-trip
	DecodePackedRoundTrip(t, call)
}

// TestPackedCompareWithSolidityEncodePacked verifies our encoding matches Solidity's abi.encodePacked
func TestPackedCompareWithSolidityEncodePacked(t *testing.T) {
	// This test verifies known encodings from Solidity
	// abi.encodePacked(address(0x1234...1234), uint256(100)) produces:
	// 0x1234567890123456789012345678901234567890 + 0x0000...0064 (32 bytes)

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	amount := big.NewInt(100)

	call := &PackedTransferCall{
		To:     to,
		Amount: amount,
	}

	encoded, err := call.PackedEncode()
	require.NoError(t, err)

	// Known Solidity output for abi.encodePacked(address, uint256(100))
	expectedHex := "1234567890123456789012345678901234567890" +
		"0000000000000000000000000000000000000000000000000000000000000064"
	expected, err := hex.DecodeString(expectedHex)
	require.NoError(t, err)

	require.Equal(t, expected, encoded)
}
