//go:build !uint256

package abi

import (
	"io"
	"math/big"
)

// Uint256 is the Go type used to represent ABI uint256 values. It is an alias
// (selected by build tag) so generated code can be byte-identical between the
// default math/big build and the holiman/uint256 build. Take/return *Uint256
// in generated code; runtime helpers (EncodeUint256, DecodeUint256) are also
// tag-gated to operate on the matching concrete type.
type Uint256 = big.Int

// WriteUint256 writes value as a 32-byte big-endian unsigned integer into
// buf[:32]. Returns io.ErrShortBuffer if buf is shorter than 32 bytes,
// or errors from EncodeBigInt for negative or oversized values.
func WriteUint256(value *Uint256, buf []byte) error {
	if len(buf) < 32 {
		return io.ErrShortBuffer
	}
	return EncodeBigInt(value, buf[:32], false)
}

// ReadUint256 reads a 32-byte big-endian unsigned integer from data[:32].
// Caller must ensure len(data) >= 32.
func ReadUint256(data []byte) *Uint256 {
	return new(Uint256).SetBytes(data[:32])
}
