package abi

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
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
