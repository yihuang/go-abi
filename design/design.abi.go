package abi

import (
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/yihuang/go-abi"
)

const (
	TransferCallStaticSize = 32 * 8
	Tuple0StaticSize       = 32
	DynamicTupleStaticSize = 64
)

type Tuple0 struct {
	Field0 *big.Int
}

func (t Tuple0) EncodedSize() int {
	return Tuple0StaticSize
}

func (t Tuple0) DynamicSize() int {
	return 0
}

func (t Tuple0) EncodeTo(buf []byte) (int, error) {
	if err := abi.EncodeBigInt(t.Field0, buf, false); err != nil {
		return 0, err
	}
	return 32, nil
}

type DynamicTuple struct {
	Field0 string
	Field1 *big.Int
}

func (t DynamicTuple) EncodedSize() int {
	return DynamicTupleStaticSize + t.DynamicSize()
}

func (t DynamicTuple) DynamicSize() int {
	var dynamicSize int

	// Field0 (dynamic)
	dynamicSize += 32                       // length
	dynamicSize += abi.Pad32(len(t.Field0)) // data (padded to 32 bytes)

	return dynamicSize
}

func (t DynamicTuple) EncodeTo(buf []byte) (int, error) {
	dynOffset := DynamicTupleStaticSize

	// Encode Field0 (offset)
	binary.BigEndian.PutUint64(buf[24:32], uint64(dynOffset))

	// Encode Field0 (dynamic)
	{
		// length
		binary.BigEndian.PutUint64(buf[dynOffset+24:dynOffset+32], uint64(len(t.Field0)))
		// data
		copy(buf[dynOffset+32:], []byte(t.Field0))
		dynOffset += 32 + abi.Pad32(len(t.Field0))
	}

	// Encode Field1
	if err := abi.EncodeBigInt(t.Field1, buf[32:], false); err != nil {
		return 0, err
	}

	return dynOffset, nil
}

type TransferCall struct {
	Memo          string
	To            common.Address
	Amount        *big.Int
	StaticArray   [2]*big.Int
	DynamicArray  []Tuple0
	DynamicArray2 [][]DynamicTuple
	Negative      *big.Int
}

func (t TransferCall) DynamicSize() int {
	var dynamicSize int

	// Memo (dynamic)
	dynamicSize += 32                     // length
	dynamicSize += abi.Pad32(len(t.Memo)) // data (padded to 32 bytes)

	// DynamicArray
	dynamicSize += 32                                     // length
	dynamicSize += len(t.DynamicArray) * Tuple0StaticSize // data

	// DynamicArray2
	dynamicSize += 32 // length
	{                 // data
		// offsets
		dynamicSize += len(t.DynamicArray2) * 32
		// dynamic parts
		for _, item := range t.DynamicArray2 {
			dynamicSize += 32 // length
			{                 // data
				// offsets
				dynamicSize += len(item) * 32
				// dynamic parts
				for _, item2 := range item {
					dynamicSize += item2.EncodedSize()
				}
			}
		}
	}

	return dynamicSize
}

func (t TransferCall) EncodedSize() int {
	return TransferCallStaticSize + t.DynamicSize()
}

func (t TransferCall) EncodeTo(buf []byte) (int, error) {
	dynOffset := TransferCallStaticSize

	// Encode Memo (offset)
	binary.BigEndian.PutUint64(buf[24:32], uint64(dynOffset))

	// Encode Memo (dynamic)
	{
		// length
		binary.BigEndian.PutUint64(buf[dynOffset+24:dynOffset+32], uint64(len(t.Memo)))
		// data
		copy(buf[dynOffset+32:], []byte(t.Memo))
		dynOffset += 32 + abi.Pad32(len(t.Memo))
	}

	// Encode To
	copy(buf[44:64], t.To.Bytes())

	// Encode Amount
	if err := abi.EncodeBigInt(t.Amount, buf[64:], false); err != nil {
		return 0, err
	}

	// Encode StaticArray
	for i := 0; i < 2; i++ {
		if err := abi.EncodeBigInt(t.StaticArray[i], buf[96+32*i:], false); err != nil {
			return 0, err
		}
	}

	// Encode DynamicArray (offset)
	binary.BigEndian.PutUint64(buf[160+24:160+32], uint64(dynOffset))

	// Encode DynamicArray (dynamic)
	{
		// length
		binary.BigEndian.PutUint64(buf[dynOffset+24:dynOffset+32], uint64(len(t.DynamicArray)))
		dynOffset += 32
		// data
		written, err := abi.ArrayEncodeTo(buf[dynOffset:], t.DynamicArray, Tuple0StaticSize, false)
		if err != nil {
			return 0, err
		}
		dynOffset += written
	}

	// Encode DynamicArray2 (offset)
	binary.BigEndian.PutUint64(buf[192+24:192+32], uint64(dynOffset))

	// Encode DynamicArray2 (dynamic)
	{
		// length
		binary.BigEndian.PutUint64(buf[dynOffset+24:dynOffset+32], uint64(len(t.DynamicArray2)))
		dynOffset += 32

		var written int
		// data
		{
			buf := buf[dynOffset:]

			dynOffset := len(t.DynamicArray2) * 32
			var offset int
			for _, item := range t.DynamicArray2 {
				// write offset
				binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynOffset))
				offset += 32

				// write dynamic part
				// length
				binary.BigEndian.PutUint64(buf[dynOffset+24:dynOffset+32], uint64(len(t.DynamicArray2)))
				dynOffset += 32
				// data
				var written int
				{
					buf := buf[dynOffset:]

					dynOffset := len(item) * 32
					var offset int
					for _, item := range item {
						// write offset
						binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynOffset))
						offset += 32

						// write dynamic part
						n, err := item.EncodeTo(buf[dynOffset:])
						if err != nil {
							return 0, err
						}
						dynOffset += n
					}
					written = dynOffset
				}
				dynOffset += written
			}
			written = dynOffset
		}
		dynOffset += written
	}

	// Encode Negative
	if err := abi.EncodeBigInt(t.Negative, buf[224:], true); err != nil {
		return 0, err
	}

	return dynOffset, nil
}

func (t TransferCall) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}
