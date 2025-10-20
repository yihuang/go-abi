package abi

import (
	"cmp"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

func EncodeBigInt(n *big.Int, buf []byte, signed bool) error {
	if n.Sign() < 0 {
		if !signed {
			return errors.New("negative integer for unsigned type")
		}

		// slow path for negative value
		copy(buf, math.U256Bytes(n))
		return nil
	}

	l := (n.BitLen() + 7) / 8
	if l > 32 {
		return errors.New("integer too large")
	}
	n.FillBytes(buf[32-l : 32])
	return nil
}

func Pad32(n int) int {
	return (n + 31) / 32 * 32
}

func SortedMapKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// IsDynamicType returns true if the type is dynamic.
// The following types are called “dynamic”:
// * bytes
// * string
// * T[] for any T
// * T[k] for any dynamic T and any k >= 0
// * (T1,...,Tk) if Ti is dynamic for some 1 <= i <= k
func IsDynamicType(t abi.Type) bool {
	if t.T == abi.TupleTy {
		for _, elem := range t.TupleElems {
			if IsDynamicType(*elem) {
				return true
			}
		}
		return false
	}
	return t.T == abi.StringTy || t.T == abi.BytesTy || t.T == abi.SliceTy || (t.T == abi.ArrayTy && IsDynamicType(*t.Elem))
}

// GetTypeSize returns the size that this type needs to occupy.
// We distinguish static and dynamic types. Static types are encoded in-place
// and dynamic types are encoded at a separately allocated location after the
// current block.
// So for a static variable, the size returned represents the size that the
// variable actually occupies.
// For a dynamic variable, the returned size is fixed 32 bytes, which is used
// to store the location reference for actual value storage.
func GetTypeSize(t abi.Type) int {
	if t.T == abi.ArrayTy && !IsDynamicType(*t.Elem) {
		// Recursively calculate type size if it is a nested array
		if t.Elem.T == abi.ArrayTy || t.Elem.T == abi.TupleTy {
			return t.Size * GetTypeSize(*t.Elem)
		}
		return t.Size * 32
	} else if t.T == abi.TupleTy && !IsDynamicType(t) {
		total := 0
		for _, elem := range t.TupleElems {
			total += GetTypeSize(*elem)
		}
		return total
	}
	return 32
}

func GetTupleSize(elems []*abi.Type) int {
	total := 0
	for _, elem := range elems {
		total += GetTypeSize(*elem)
	}
	return total
}

// GenTupleIdentifier generates a unique identifier for a tuple type
func GenTupleIdentifier(t abi.Type) string {
	// Create a signature based on tuple element types
	types := make([]string, len(t.TupleElems))
	for i, elem := range t.TupleElems {
		types[i] = elem.String()
	}
	sig := fmt.Sprintf("(%v)", strings.Join(types, ","))

	id := crypto.Keccak256([]byte(sig))
	return "Tuple_" + hex.EncodeToString(id)[:8] // Use first 8 chars for readability
}

// TupleStructName generates a unique struct name for a tuple type
func TupleStructName(t abi.Type) string {
	if t.TupleRawName != "" {
		return t.TupleRawName
	}

	// Use the tuple's string representation as the basis for the struct name
	// This creates a deterministic name based on the tuple structure
	return GenTupleIdentifier(t)
}
