package abi

import (
	"cmp"
	"errors"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common/math"
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
