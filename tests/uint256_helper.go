//go:build !uint256

package tests

import "math/big"

func newUint256FromInt64(v int64) *big.Int { return big.NewInt(v) }
