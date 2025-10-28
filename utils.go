package abi

import (
	"errors"
	"io"
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
