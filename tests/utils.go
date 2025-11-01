package tests

import (
	"errors"
	"io"
	"testing"

	"github.com/test-go/testify/require"
	"github.com/yihuang/go-abi"
)

type Tuple[T any] interface {
	abi.Tuple

	*T
}

func DecodeRoundTrip[T any, PT Tuple[T]](t *testing.T, orig PT) {
	data, err := orig.Encode()
	require.NoError(t, err)

	var decoded T
	_, err = PT(&decoded).Decode(data)
	require.NoError(t, err)

	require.Equal(t, orig, &decoded)

	// test ErrUnexpectedEOF
	for i := 0; i < len(data); i++ {
		_, err = PT(&decoded).Decode(data[:i])
		require.Error(t, err)
		require.True(t, errors.Is(err, io.ErrUnexpectedEOF))
	}
}
