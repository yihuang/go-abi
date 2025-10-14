package abi

import (
	"cmp"
	"encoding/binary"
	"errors"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common/math"
)

type Tuple interface {
	EncodedSize() int
	EncodeTo(buf []byte) (int, error)
}

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

func ArrayEncodedSize[T Tuple](items []T, itemStaticSize int, isDynamic bool) int {
	if !isDynamic {
		return len(items) * itemStaticSize
	}

	// offsets
	size := len(items) * 32
	// dynamic parts
	for _, item := range items {
		size += item.EncodedSize()
	}

	return size
}

// array is like a tuple of same element types
func ArrayEncodeTo[T Tuple](buf []byte, items []T, staticSize int, isDynamic bool) (int, error) {
	if !isDynamic {
		var offset int
		for _, item := range items {
			if _, err := item.EncodeTo(buf[offset : offset+staticSize]); err != nil {
				return 0, err
			}
			offset += staticSize
		}
		return offset, nil
	}

	dynOffset := len(items) * 32
	var offset int
	for _, item := range items {
		// write offset
		binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynOffset))
		offset += 32

		// write dynamic part
		n, err := item.EncodeTo(buf[dynOffset:])
		if err != nil {
			return 0, err
		}
		dynOffset += n
	}

	return dynOffset, nil
}

func SortedMapKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
