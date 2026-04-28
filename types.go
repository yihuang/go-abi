package abi

import (
	"github.com/ethereum/go-ethereum/common"
)

type Encode interface {
	EncodedSize() int
	Encode() ([]byte, error)
	EncodeTo([]byte) (int, error)
}

type Decode interface {
	// Decode returns io.UnexpectedEOF if data is too short.
	Decode([]byte) (int, error)
}

type Tuple interface {
	Encode
	Decode
}

type PackedEncode interface {
	PackedEncodedSize() int
	PackedEncode() ([]byte, error)
	PackedEncodeTo([]byte) (int, error)
}

type PackedDecode interface {
	PackedDecode([]byte) (int, error)
}

type PackedTuple interface {
	PackedEncode
	PackedDecode
}

type Method interface {
	Tuple

	EncodeWithSelector() ([]byte, error)

	GetMethodName() string
	GetMethodID() uint32
	GetMethodSelector() [4]byte
}

type Event interface {
	// indexed fields
	EncodeTopics() ([]common.Hash, error)
	DecodeTopics([]common.Hash) error

	// data fields
	Tuple

	// metadata
	GetEventName() string
	GetEventID() common.Hash
}

type EmptyTuple struct{}

func (EmptyTuple) EncodedSize() int {
	return 0
}

func (EmptyTuple) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (EmptyTuple) EncodeTo(data []byte) (int, error) {
	return 0, nil
}

func (EmptyTuple) Decode(data []byte) (int, error) {
	return 0, nil
}

func (EmptyTuple) PackedEncodedSize() int {
	return 0
}

func (EmptyTuple) PackedEncode() ([]byte, error) {
	return []byte{}, nil
}

func (EmptyTuple) PackedEncodeTo(data []byte) (int, error) {
	return 0, nil
}

func (EmptyTuple) PackedDecode(data []byte) (int, error) {
	return 0, nil
}

type EmptyIndexed struct{}

func (EmptyIndexed) EncodeTopics() ([]common.Hash, error) {
	return nil, nil
}

func (EmptyIndexed) DecodeTopics([]common.Hash) error {
	return nil
}
