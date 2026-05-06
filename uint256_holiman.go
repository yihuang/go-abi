//go:build uint256

package abi

import (
	"io"

	"github.com/holiman/uint256"
)

// Uint256 is the Go type used to represent ABI uint256 values. It is an alias
// (selected by build tag) so generated code can be byte-identical between the
// default math/big build and the holiman/uint256 build. Take/return *Uint256
// in generated code; runtime helpers (EncodeUint256, DecodeUint256) are also
// tag-gated to operate on the matching concrete type.
type Uint256 = uint256.Int

// WriteUint256 writes value as a 32-byte big-endian unsigned integer into
// buf[:32]. Returns io.ErrShortBuffer if buf is shorter than 32 bytes;
// otherwise never errors (uint256.Int is bounded to [0, 2^256)).
func WriteUint256(value *Uint256, buf []byte) error {
	if len(buf) < 32 {
		return io.ErrShortBuffer
	}
	value.WriteToArray32((*[32]byte)(buf[:32]))
	return nil
}

// ReadUint256 reads a 32-byte big-endian unsigned integer from data[:32].
// Uses the SetBytes32 fast path. Returns io.ErrUnexpectedEOF if data is
// shorter than 32 bytes.
func ReadUint256(data []byte) (*Uint256, error) {
	if len(data) < 32 {
		return nil, io.ErrUnexpectedEOF
	}
	return new(Uint256).SetBytes32(data[:32]), nil
}
