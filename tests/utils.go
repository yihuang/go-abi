package testdata

import (
	"testing"

	"github.com/test-go/testify/require"
	"github.com/yihuang/go-abi"
)

type Tuple[T any] interface {
	abi.Encode
	abi.Decode

	*T
}

func DecodeRoundTrip[T any, PT Tuple[T]](t *testing.T, orig PT) {
	data, err := orig.Encode()
	require.NoError(t, err)

	var decoded T
	err = PT(&decoded).Decode(data)
	require.NoError(t, err)

	require.Equal(t, orig, &decoded)
}
