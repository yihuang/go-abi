//go:build uint256

package tests

import "github.com/holiman/uint256"

func newUint256FromInt64(v int64) *uint256.Int { return uint256.NewInt(uint64(v)) }
