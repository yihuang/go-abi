package abi

import (
	"github.com/ethereum/go-ethereum/common"
)

type Encode interface {
	Encode() ([]byte, error)
}

type Decode interface {
	Decode([]byte) error
}

type EmptyTuple struct{}

func (e EmptyTuple) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (e *EmptyTuple) Decode(data []byte) error {
	return nil
}

type EmptyIndexed struct{}

func (e EmptyIndexed) EncodeTopics() []common.Hash {
	return nil
}

func (e *EmptyIndexed) DecodeTopics([]common.Hash) error {
	return nil
}
