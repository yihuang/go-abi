package abi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

func Pad32(n int) int {
	return (n + 31) / 32 * 32
}

func EncodeBigInt(n *big.Int, buf []byte, signed bool) error {
	if n.Sign() < 0 {
		if !signed {
			return errors.New("negative integer for unsigned type")
		}

		// slow path for negative value
		// Use a copy of n to avoid modification from math.U256Bytes
		tmp := new(big.Int).Set(n)
		copy(buf, math.U256Bytes(tmp))
		return nil
	}

	l := (n.BitLen() + 7) / 8
	if l > 32 {
		return errors.New("integer too large")
	}
	n.FillBytes(buf[32-l : 32])
	return nil
}

func DecodeBigInt(data []byte, signed bool) (*big.Int, error) {
	if len(data) < 32 {
		return nil, io.ErrUnexpectedEOF
	}

	if signed && data[0]&0x80 != 0 {
		// negative number
		tmp := make([]byte, 32)
		for i := 0; i < 32; i++ {
			tmp[i] = ^data[i]
		}
		bigN := new(big.Int).SetBytes(tmp)
		bigN.Add(bigN, big.NewInt(1))
		bigN.Neg(bigN)
		return bigN, nil
	}

	bigN := new(big.Int).SetBytes(data[:32])
	return bigN, nil
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
