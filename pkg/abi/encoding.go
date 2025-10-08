package abi

import (
	"encoding/binary"
	"errors"
	"math/big"
)

var (
	// ErrInvalidType is returned when an invalid type is encountered
	ErrInvalidType = errors.New("invalid type")

	// ErrInvalidValue is returned when an invalid value is provided
	ErrInvalidValue = errors.New("invalid value")

	// ErrInsufficientData is returned when there's not enough data to decode
	ErrInsufficientData = errors.New("insufficient data")
)

// WordSize is the size of an ABI word in bytes
const WordSize = 32

// encodeWord encodes a single 32-byte word
func encodeWord(data []byte) []byte {
	if len(data) > WordSize {
		panic("data exceeds word size")
	}

	result := make([]byte, WordSize)
	copy(result[WordSize-len(data):], data)
	return result
}

// encodeUint encodes an unsigned integer to ABI format
func encodeUint(value *big.Int, bits int) []byte {
	if value.Sign() < 0 {
		panic("negative value for unsigned integer")
	}

	// Convert to bytes (big-endian)
	bytes := value.Bytes()

	// Ensure we have exactly 32 bytes
	result := make([]byte, WordSize)
	if len(bytes) > WordSize {
		// Truncate if too large (shouldn't happen for valid values)
		copy(result, bytes[len(bytes)-WordSize:])
	} else {
		// Pad with zeros on the left
		copy(result[WordSize-len(bytes):], bytes)
	}

	return result
}

// encodeInt encodes a signed integer to ABI format
func encodeInt(value *big.Int, bits int) []byte {
	if value.Sign() >= 0 {
		// Positive number - same as unsigned
		return encodeUint(value, bits)
	}

	// Negative number - two's complement
	twosComplement := new(big.Int).Add(
		new(big.Int).Lsh(big.NewInt(1), uint(bits)),
		value,
	)

	return encodeUint(twosComplement, bits)
}

// decodeUint decodes an unsigned integer from ABI format
func decodeUint(data []byte, bits int) (*big.Int, error) {
	if len(data) < WordSize {
		return nil, ErrInsufficientData
	}

	value := new(big.Int).SetBytes(data[:WordSize])

	// Check if value fits in the specified bits
	maxValue := new(big.Int).Lsh(big.NewInt(1), uint(bits))
	if value.Cmp(maxValue) >= 0 {
		return nil, ErrInvalidValue
	}

	return value, nil
}

// decodeInt decodes a signed integer from ABI format
func decodeInt(data []byte, bits int) (*big.Int, error) {
	if len(data) < WordSize {
		return nil, ErrInsufficientData
	}

	value := new(big.Int).SetBytes(data[:WordSize])

	// Check if the highest bit is set (negative number)
	highestBit := new(big.Int).Lsh(big.NewInt(1), uint(bits-1))
	if value.Cmp(highestBit) >= 0 {
		// Negative number - convert from two's complement
		maxValue := new(big.Int).Lsh(big.NewInt(1), uint(bits))
		value.Sub(value, maxValue)
	}

	return value, nil
}

// encodeBool encodes a boolean to ABI format
func encodeBool(value bool) []byte {
	if value {
		return encodeUint(big.NewInt(1), 8)
	}
	return encodeUint(big.NewInt(0), 8)
}

// decodeBool decodes a boolean from ABI format
func decodeBool(data []byte) (bool, error) {
	value, err := decodeUint(data, 8)
	if err != nil {
		return false, err
	}

	if value.Cmp(big.NewInt(1)) == 0 {
		return true, nil
	} else if value.Cmp(big.NewInt(0)) == 0 {
		return false, nil
	}

	return false, ErrInvalidValue
}

// encodeAddress encodes an Ethereum address to ABI format
func encodeAddress(addr [20]byte) []byte {
	result := make([]byte, WordSize)
	copy(result[WordSize-20:], addr[:])
	return result
}

// decodeAddress decodes an Ethereum address from ABI format
func decodeAddress(data []byte) ([20]byte, error) {
	var addr [20]byte

	if len(data) < WordSize {
		return addr, ErrInsufficientData
	}

	// Address is right-aligned in the 32-byte word
	copy(addr[:], data[WordSize-20:WordSize])
	return addr, nil
}

// encodeBytes encodes bytes to ABI format
func encodeBytes(data []byte) []byte {
	// For dynamic types, we encode length first, then data
	length := uint64(len(data))

	// Calculate padded data length
	paddedLength := (length + WordSize - 1) / WordSize * WordSize

	result := make([]byte, WordSize+paddedLength)

	// Encode length
	binary.BigEndian.PutUint64(result[WordSize-8:WordSize], length)

	// Copy data
	copy(result[WordSize:], data)

	return result
}

// decodeBytes decodes bytes from ABI format
func decodeBytes(data []byte) ([]byte, uint, error) {
	if len(data) < WordSize {
		return nil, 0, ErrInsufficientData
	}

	// Read length
	length := binary.BigEndian.Uint64(data[WordSize-8 : WordSize])

	// Calculate total size needed
	totalSize := uint(WordSize) + ((uint(length) + WordSize - 1) / WordSize * WordSize)

	if uint(len(data)) < totalSize {
		return nil, 0, ErrInsufficientData
	}

	// Extract bytes
	bytes := make([]byte, length)
	copy(bytes, data[WordSize:WordSize+length])

	return bytes, totalSize, nil
}