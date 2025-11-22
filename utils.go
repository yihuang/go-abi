package abi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
)

const (
	// max values for all unsigned small integers of all bytes
	MaxUint8  = math.MaxUint8
	MaxUint16 = math.MaxUint16
	MaxUint24 = 1<<24 - 1
	MaxUint32 = math.MaxUint32
	MaxUint40 = 1<<40 - 1
	MaxUint48 = 1<<48 - 1
	MaxUint56 = 1<<56 - 1
	MaxUint64 = math.MaxUint64

	// min values for all signed small integers of all bytes
	MinInt8  = math.MinInt8
	MinInt16 = math.MinInt16
	MinInt24 = -1 << 23
	MinInt32 = math.MinInt32
	MinInt40 = -1 << 39
	MinInt48 = -1 << 47
	MinInt56 = -1 << 55
	MinInt64 = math.MinInt64

	// max values for all signed small integers of all bytes
	MaxInt8  = math.MaxInt8
	MaxInt16 = math.MaxInt16
	MaxInt24 = 1<<23 - 1
	MaxInt32 = math.MaxInt32
	MaxInt40 = 1<<39 - 1
	MaxInt48 = 1<<47 - 1
	MaxInt56 = 1<<55 - 1
	MaxInt64 = math.MaxInt64
)

var (
	// sign extension padding bytes
	PaddingBytes8  = bytes.Repeat([]byte{0xff}, 31)
	PaddingBytes16 = bytes.Repeat([]byte{0xff}, 30)
	PaddingBytes32 = bytes.Repeat([]byte{0xff}, 28)
	PaddingBytes64 = bytes.Repeat([]byte{0xff}, 24)
)

var (
	tt256      = new(big.Int).Lsh(common.Big1, 256)
	MaxUint256 = new(big.Int).Sub(tt256, common.Big1)
)

func Pad32(n int) int {
	return (n + 31) / 32 * 32
}

// DecodeUint is common utility to decode a small unsigned integer value from 32 bytes
// the caller must pass correct maxValue for the target type T
func DecodeUint[T int | uint8 | uint16 | uint32 | uint64](data []byte, maxValue uint64) (T, error) {
	var n uint256.Int
	n.SetBytes32(data)

	result, overflow := n.Uint64WithOverflow()
	if overflow || result > maxValue {
		return 0, ErrDirtyPadding
	}

	return T(result), nil
}

func DecodeInt[T int8 | int16 | int32 | int64](data []byte, minValue, maxValue int64) (T, error) {
	var n uint256.Int
	n.SetBytes32(data)

	i64 := int64(n.Uint64())

	// check sign extension in higher bytes
	if i64 < 0 {
		// should be all 1s
		if n[1]&n[2]&n[3] != ^uint64(0) {
			return 0, ErrDirtyPadding
		}
	} else {
		// should be all 0s
		if n[1]|n[2]|n[3] != 0 {
			return 0, ErrDirtyPadding
		}
	}

	if i64 < minValue || i64 > maxValue {
		return 0, ErrDirtyPadding
	}

	return T(i64), nil
}

func DecodeSize(data []byte) (int, error) {
	v, err := DecodeUint[int](data, math.MaxInt)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func EncodeBigInt(n *big.Int, buf []byte, signed bool) error {
	if n.Sign() < 0 {
		if !signed {
			return ErrNegativeValue
		}

		// convert to 256 bit two's complement
		n = new(big.Int).And(n, MaxUint256)
	}

	l := (n.BitLen() + 7) / 8
	if l > 32 {
		return ErrIntegerTooLarge
	}
	n.FillBytes(buf[32-l : 32])
	return nil
}

func DecodeBigInt(data []byte, signed bool) (*big.Int, error) {
	if len(data) < 32 {
		return nil, io.ErrUnexpectedEOF
	}

	ret := new(big.Int).SetBytes(data[:32])
	if signed && data[0]&0x80 != 0 {
		ret.Sub(ret, tt256)
	}

	return ret, nil
}

func EncodeEvent(event Event) ([]common.Hash, []byte, error) {
	topics, err := event.EncodeTopics()
	if err != nil {
		return nil, nil, err
	}

	data, err := event.Encode()
	if err != nil {
		return nil, nil, err
	}

	return topics, data, nil
}

func DecodeEvent(event Event, topics []common.Hash, data []byte) error {
	if err := event.DecodeTopics(topics); err != nil {
		return err
	}

	_, err := event.Decode(data)
	return err
}

// GenTypeIdentifier generates a unique identifier for any ABI type
// This is used to create unique function names for encoding/decoding
func GenTypeIdentifier(t ethabi.Type) string {
	switch t.T {
	case ethabi.UintTy:
		if t.Size <= 64 {
			return fmt.Sprintf("Uint%d", t.Size)
		}
		return "Uint256"
	case ethabi.IntTy:
		if t.Size <= 64 {
			return fmt.Sprintf("Int%d", t.Size)
		}
		return "Int256"
	case ethabi.AddressTy:
		return "Address"
	case ethabi.BoolTy:
		return "Bool"
	case ethabi.StringTy:
		return "String"
	case ethabi.BytesTy:
		return "Bytes"
	case ethabi.FixedBytesTy:
		return fmt.Sprintf("Bytes%d", t.Size)
	case ethabi.SliceTy:
		return fmt.Sprintf("%sSlice", GenTypeIdentifier(*t.Elem))
	case ethabi.ArrayTy:
		return fmt.Sprintf("%sArray%d", GenTypeIdentifier(*t.Elem), t.Size)
	case ethabi.TupleTy:
		return TupleStructName(t) // Reuse existing tuple identifier logic
	default:
		panic("unsupported ABI type for identifier generation: " + t.String())
	}
}

// GenTupleIdentifier generates a unique identifier for a tuple type
func GenTupleIdentifier(t ethabi.Type) string {
	// Create a signature based on tuple element types
	types := make([]string, len(t.TupleElems))
	for i, elem := range t.TupleElems {
		types[i] = elem.String()
	}

	sig := fmt.Sprintf("(%v)", strings.Join(types, ","))
	id := crypto.Keccak256([]byte(sig))
	return "Tuple" + hex.EncodeToString(id)[:8] // Use first 8 chars for readability
}

// TupleStructName generates a unique struct name for a tuple type
func TupleStructName(t ethabi.Type) string {
	if t.TupleRawName != "" {
		return t.TupleRawName
	}

	// Use the tuple's string representation as the basis for the struct name
	// This creates a deterministic name based on the tuple structure
	return GenTupleIdentifier(t)
}
