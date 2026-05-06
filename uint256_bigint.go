//go:build !uint256

package abi

import "math/big"

// Uint256 is the Go type used to represent ABI uint256 values. It is an alias
// (selected by build tag) so generated code can be byte-identical between the
// default math/big build and the holiman/uint256 build. Take/return *Uint256
// in generated code; runtime helpers (EncodeUint256, DecodeUint256) are also
// tag-gated to operate on the matching concrete type.
type Uint256 = big.Int
