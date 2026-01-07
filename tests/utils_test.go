package tests

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/yihuang/go-abi"
)

var testAddress = common.HexToAddress("0x1234567890123456789012345678901234567890")

func benchEncode(b *testing.B, call abi.Tuple) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := call.Encode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchEncodeTo(b *testing.B, call abi.Tuple) {
	buf := make([]byte, call.EncodedSize())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := call.EncodeTo(buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchDecode(b *testing.B, encoded []byte, newCall func() abi.Decode) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		call := newCall()
		_, err := call.Decode(encoded)
		if err != nil {
			b.Fatal(err)
		}
	}
}
