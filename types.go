package abi

import (
	"github.com/ethereum/go-ethereum/common"
)

type Encode interface {
	Encode() ([]byte, error)
}

type Decode interface {
	Decode([]byte) (int, error)
}

type EmptyTuple struct{}

func (e EmptyTuple) EncodedSize() int {
	return 0
}

func (e EmptyTuple) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (e EmptyTuple) EncodeTo(data []byte) (int, error) {
	return 0, nil
}

func (e *EmptyTuple) Decode(data []byte) (int, error) {
	return 0, nil
}

type EmptyIndexed struct{}

func (e EmptyIndexed) EncodeTopics() []common.Hash {
	return nil
}

func (e *EmptyIndexed) DecodeTopics([]common.Hash) error {
	return nil
}
