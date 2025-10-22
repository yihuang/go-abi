package testdata

import (
	"testing"

	"github.com/test-go/testify/require"
)

type Encodable interface {
	EncodeTo(data []byte) (int, error)
}

type Decodable[T any] interface {
	Encode() ([]byte, error)
	DecodeFrom(data []byte) error

	*T
}

func DecodeRoundTrip[T any, PT Decodable[T]](t *testing.T, orig PT) {
	data, err := orig.Encode()
	require.NoError(t, err)

	var decoded T
	err = PT(&decoded).DecodeFrom(data)
	require.NoError(t, err)

	require.Equal(t, orig, &decoded)
}
