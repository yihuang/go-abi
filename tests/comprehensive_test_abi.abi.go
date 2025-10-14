package testdata

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/yihuang/go-abi"
	"math/big"
)

// Tuple_1064c5d1 represents a tuple type

var _ abi.Tuple = Tuple_1064c5d1{}

const Tuple_1064c5d1StaticSize = 32

type Tuple_1064c5d1 struct {
	Level2 Tuple_e05eda74
}

// EncodedSize returns the total encoded size of Tuple_1064c5d1
func (t Tuple_1064c5d1) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += t.Level2.EncodedSize() // dynamic tuple

	return Tuple_1064c5d1StaticSize + dynamicSize
}

// EncodeTo encodes Tuple_1064c5d1 to ABI bytes in the provided buffer
func (t Tuple_1064c5d1) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := Tuple_1064c5d1StaticSize // Start dynamic data after static section

	// Level2 (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Level2 (dynamic)

	n, err := t.Level2.EncodeTo(buf[dynamicOffset:])
	if err != nil {
		return 0, err
	}
	dynamicOffset += n

	return dynamicOffset, nil
}

// Encode encodes Tuple_1064c5d1 to ABI bytes
func (t Tuple_1064c5d1) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Tuple_54b20b3a represents a tuple type

var _ abi.Tuple = Tuple_54b20b3a{}

const Tuple_54b20b3aStaticSize = 64

type Tuple_54b20b3a struct {
	Value       *big.Int
	Description string
}

// EncodedSize returns the total encoded size of Tuple_54b20b3a
func (t Tuple_54b20b3a) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.Description)) // length + padded string data

	return Tuple_54b20b3aStaticSize + dynamicSize
}

// EncodeTo encodes Tuple_54b20b3a to ABI bytes in the provided buffer
func (t Tuple_54b20b3a) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := Tuple_54b20b3aStaticSize // Start dynamic data after static section

	// Value (static)

	if err := abi.EncodeBigInt(t.Value, buf[0:32], false); err != nil {
		return 0, err
	}

	// Description (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Description (dynamic)

	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Description)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], []byte(t.Description))
	dynamicOffset += abi.Pad32(len(t.Description))

	return dynamicOffset, nil
}

// Encode encodes Tuple_54b20b3a to ABI bytes
func (t Tuple_54b20b3a) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Tuple_5cee9471 represents a tuple type

var _ abi.Tuple = Tuple_5cee9471{}

const Tuple_5cee9471StaticSize = 32

type Tuple_5cee9471 struct {
	Level1 Tuple_1064c5d1
}

// EncodedSize returns the total encoded size of Tuple_5cee9471
func (t Tuple_5cee9471) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += t.Level1.EncodedSize() // dynamic tuple

	return Tuple_5cee9471StaticSize + dynamicSize
}

// EncodeTo encodes Tuple_5cee9471 to ABI bytes in the provided buffer
func (t Tuple_5cee9471) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := Tuple_5cee9471StaticSize // Start dynamic data after static section

	// Level1 (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Level1 (dynamic)

	n, err := t.Level1.EncodeTo(buf[dynamicOffset:])
	if err != nil {
		return 0, err
	}
	dynamicOffset += n

	return dynamicOffset, nil
}

// Encode encodes Tuple_5cee9471 to ABI bytes
func (t Tuple_5cee9471) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Tuple_8a486b93 represents a tuple type

var _ abi.Tuple = Tuple_8a486b93{}

const Tuple_8a486b93StaticSize = 96

type Tuple_8a486b93 struct {
	Name     string
	Emails   []string
	Metadata Tuple_dc8f1c28
}

// EncodedSize returns the total encoded size of Tuple_8a486b93
func (t Tuple_8a486b93) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.Name)) // length + padded string data
	dynamicSize += 32 + 32*len(t.Emails)       // length + offset pointers for dynamic elements
	for _, elem := range t.Emails {
		dynamicSize += 32 + abi.Pad32(len(elem)) // length + padded string data
	}
	dynamicSize += t.Metadata.EncodedSize() // dynamic tuple

	return Tuple_8a486b93StaticSize + dynamicSize
}

// EncodeTo encodes Tuple_8a486b93 to ABI bytes in the provided buffer
func (t Tuple_8a486b93) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := Tuple_8a486b93StaticSize // Start dynamic data after static section

	// Name (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Name (dynamic)

	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Name)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], []byte(t.Name))
	dynamicOffset += abi.Pad32(len(t.Name))

	// Emails (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Emails (dynamic)

	{
		// length
		binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Emails)))
		dynamicOffset += 32

		var written int

		// data with dynamic region
		{
			buf := buf[dynamicOffset:]
			dynamicOffset := len(t.Emails) * 32 // start after static region

			var offset int
			for _, item := range t.Emails {
				// write offsets
				binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynamicOffset))
				offset += 32

				// write data (dynamic)

				// length
				binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(item)))
				dynamicOffset += 32

				// data
				copy(buf[dynamicOffset:], []byte(item))
				dynamicOffset += abi.Pad32(len(item))

			}
			written = dynamicOffset
		}
		dynamicOffset += written

	}

	// Metadata (offset)
	binary.BigEndian.PutUint64(buf[64+24:64+32], uint64(dynamicOffset))

	// Metadata (dynamic)

	n, err := t.Metadata.EncodeTo(buf[dynamicOffset:])
	if err != nil {
		return 0, err
	}
	dynamicOffset += n

	return dynamicOffset, nil
}

// Encode encodes Tuple_8a486b93 to ABI bytes
func (t Tuple_8a486b93) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Tuple_dc8f1c28 represents a tuple type

var _ abi.Tuple = Tuple_dc8f1c28{}

const Tuple_dc8f1c28StaticSize = 64

type Tuple_dc8f1c28 struct {
	CreatedAt *big.Int
	Tags      []string
}

// EncodedSize returns the total encoded size of Tuple_dc8f1c28
func (t Tuple_dc8f1c28) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + 32*len(t.Tags) // length + offset pointers for dynamic elements
	for _, elem := range t.Tags {
		dynamicSize += 32 + abi.Pad32(len(elem)) // length + padded string data
	}

	return Tuple_dc8f1c28StaticSize + dynamicSize
}

// EncodeTo encodes Tuple_dc8f1c28 to ABI bytes in the provided buffer
func (t Tuple_dc8f1c28) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := Tuple_dc8f1c28StaticSize // Start dynamic data after static section

	// CreatedAt (static)

	if err := abi.EncodeBigInt(t.CreatedAt, buf[0:32], false); err != nil {
		return 0, err
	}

	// Tags (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Tags (dynamic)

	{
		// length
		binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Tags)))
		dynamicOffset += 32

		var written int

		// data with dynamic region
		{
			buf := buf[dynamicOffset:]
			dynamicOffset := len(t.Tags) * 32 // start after static region

			var offset int
			for _, item := range t.Tags {
				// write offsets
				binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynamicOffset))
				offset += 32

				// write data (dynamic)

				// length
				binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(item)))
				dynamicOffset += 32

				// data
				copy(buf[dynamicOffset:], []byte(item))
				dynamicOffset += abi.Pad32(len(item))

			}
			written = dynamicOffset
		}
		dynamicOffset += written

	}

	return dynamicOffset, nil
}

// Encode encodes Tuple_dc8f1c28 to ABI bytes
func (t Tuple_dc8f1c28) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Tuple_de3c4b6f represents a tuple type

var _ abi.Tuple = Tuple_de3c4b6f{}

const Tuple_de3c4b6fStaticSize = 96

type Tuple_de3c4b6f struct {
	Id     uint32
	Data   []byte
	Active bool
}

// EncodedSize returns the total encoded size of Tuple_de3c4b6f
func (t Tuple_de3c4b6f) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.Data)) // length + padded bytes data

	return Tuple_de3c4b6fStaticSize + dynamicSize
}

// EncodeTo encodes Tuple_de3c4b6f to ABI bytes in the provided buffer
func (t Tuple_de3c4b6f) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := Tuple_de3c4b6fStaticSize // Start dynamic data after static section

	// Id (static)
	binary.BigEndian.PutUint32(buf[0+28:0+32], uint32(t.Id))

	// Data (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Data (dynamic)

	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Data)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], t.Data)
	dynamicOffset += abi.Pad32(len(t.Data))

	// Active (static)

	if t.Active {
		buf[64+31] = 1
	}

	return dynamicOffset, nil
}

// Encode encodes Tuple_de3c4b6f to ABI bytes
func (t Tuple_de3c4b6f) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Tuple_e05eda74 represents a tuple type

var _ abi.Tuple = Tuple_e05eda74{}

const Tuple_e05eda74StaticSize = 32

type Tuple_e05eda74 struct {
	Level3 Tuple_54b20b3a
}

// EncodedSize returns the total encoded size of Tuple_e05eda74
func (t Tuple_e05eda74) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += t.Level3.EncodedSize() // dynamic tuple

	return Tuple_e05eda74StaticSize + dynamicSize
}

// EncodeTo encodes Tuple_e05eda74 to ABI bytes in the provided buffer
func (t Tuple_e05eda74) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := Tuple_e05eda74StaticSize // Start dynamic data after static section

	// Level3 (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Level3 (dynamic)

	n, err := t.Level3.EncodeTo(buf[dynamicOffset:])
	if err != nil {
		return 0, err
	}
	dynamicOffset += n

	return dynamicOffset, nil
}

// Encode encodes Tuple_e05eda74 to ABI bytes
func (t Tuple_e05eda74) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// Tuple_e9afb3e4 represents a tuple type

var _ abi.Tuple = Tuple_e9afb3e4{}

const Tuple_e9afb3e4StaticSize = 64

type Tuple_e9afb3e4 struct {
	Id      *big.Int
	Profile Tuple_8a486b93
}

// EncodedSize returns the total encoded size of Tuple_e9afb3e4
func (t Tuple_e9afb3e4) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += t.Profile.EncodedSize() // dynamic tuple

	return Tuple_e9afb3e4StaticSize + dynamicSize
}

// EncodeTo encodes Tuple_e9afb3e4 to ABI bytes in the provided buffer
func (t Tuple_e9afb3e4) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := Tuple_e9afb3e4StaticSize // Start dynamic data after static section

	// Id (static)

	if err := abi.EncodeBigInt(t.Id, buf[0:32], false); err != nil {
		return 0, err
	}

	// Profile (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Profile (dynamic)

	n, err := t.Profile.EncodeTo(buf[dynamicOffset:])
	if err != nil {
		return 0, err
	}
	dynamicOffset += n

	return dynamicOffset, nil
}

// Encode encodes Tuple_e9afb3e4 to ABI bytes
func (t Tuple_e9afb3e4) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// TestComplexDynamicTuplesArgs represents the arguments for testComplexDynamicTuples function

var _ abi.Tuple = TestComplexDynamicTuplesArgs{}

const TestComplexDynamicTuplesArgsStaticSize = 32

type TestComplexDynamicTuplesArgs struct {
	Users []Tuple_e9afb3e4
}

// EncodedSize returns the total encoded size of TestComplexDynamicTuplesArgs
func (t TestComplexDynamicTuplesArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + 32*len(t.Users) // length + offset pointers for dynamic elements
	for _, elem := range t.Users {
		dynamicSize += elem.EncodedSize() // dynamic tuple
	}

	return TestComplexDynamicTuplesArgsStaticSize + dynamicSize
}

// EncodeTo encodes TestComplexDynamicTuplesArgs to ABI bytes in the provided buffer
func (t TestComplexDynamicTuplesArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := TestComplexDynamicTuplesArgsStaticSize // Start dynamic data after static section

	// Users (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Users (dynamic)

	{
		// length
		binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Users)))
		dynamicOffset += 32

		var written int

		// data with dynamic region
		{
			buf := buf[dynamicOffset:]
			dynamicOffset := len(t.Users) * 32 // start after static region

			var offset int
			for _, item := range t.Users {
				// write offsets
				binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynamicOffset))
				offset += 32

				// write data (dynamic)

				n, err := item.EncodeTo(buf[dynamicOffset:])
				if err != nil {
					return 0, err
				}
				dynamicOffset += n

			}
			written = dynamicOffset
		}
		dynamicOffset += written

	}

	return dynamicOffset, nil
}

// Encode encodes TestComplexDynamicTuplesArgs to ABI bytes
func (t TestComplexDynamicTuplesArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes testComplexDynamicTuples arguments to ABI bytes including function selector
func (t TestComplexDynamicTuplesArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], TestComplexDynamicTuplesArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// TestComplexDynamicTuplesArgsSelector is the function selector for testComplexDynamicTuples((uint256,(string,string[],(uint256,string[])))[])
var TestComplexDynamicTuplesArgsSelector = [4]byte{0xc0, 0x96, 0x4c, 0x93}

// Selector returns the function selector for testComplexDynamicTuples
func (TestComplexDynamicTuplesArgs) Selector() [4]byte {
	return TestComplexDynamicTuplesArgsSelector
}

// TestDeeplyNestedArgs represents the arguments for testDeeplyNested function

var _ abi.Tuple = TestDeeplyNestedArgs{}

const TestDeeplyNestedArgsStaticSize = 32

type TestDeeplyNestedArgs struct {
	Data Tuple_5cee9471
}

// EncodedSize returns the total encoded size of TestDeeplyNestedArgs
func (t TestDeeplyNestedArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += t.Data.EncodedSize() // dynamic tuple

	return TestDeeplyNestedArgsStaticSize + dynamicSize
}

// EncodeTo encodes TestDeeplyNestedArgs to ABI bytes in the provided buffer
func (t TestDeeplyNestedArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := TestDeeplyNestedArgsStaticSize // Start dynamic data after static section

	// Data (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Data (dynamic)

	n, err := t.Data.EncodeTo(buf[dynamicOffset:])
	if err != nil {
		return 0, err
	}
	dynamicOffset += n

	return dynamicOffset, nil
}

// Encode encodes TestDeeplyNestedArgs to ABI bytes
func (t TestDeeplyNestedArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes testDeeplyNested arguments to ABI bytes including function selector
func (t TestDeeplyNestedArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], TestDeeplyNestedArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// TestDeeplyNestedArgsSelector is the function selector for testDeeplyNested(((((uint256,string)))))
var TestDeeplyNestedArgsSelector = [4]byte{0x21, 0x75, 0xe8, 0x54}

// Selector returns the function selector for testDeeplyNested
func (TestDeeplyNestedArgs) Selector() [4]byte {
	return TestDeeplyNestedArgsSelector
}

// TestFixedArraysArgs represents the arguments for testFixedArrays function

var _ abi.Tuple = TestFixedArraysArgs{}

const TestFixedArraysArgsStaticSize = 320

type TestFixedArraysArgs struct {
	Addresses [5]common.Address
	Uints     [3]*big.Int
	Bytes32s  [2][32]byte
}

// EncodedSize returns the total encoded size of TestFixedArraysArgs
func (t TestFixedArraysArgs) EncodedSize() int {
	dynamicSize := 0

	return TestFixedArraysArgsStaticSize + dynamicSize
}

// EncodeTo encodes TestFixedArraysArgs to ABI bytes in the provided buffer
func (t TestFixedArraysArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := TestFixedArraysArgsStaticSize // Start dynamic data after static section

	// Addresses (static)

	// Encode fixed-size array t.Addresses
	{
		offset := 0
		for _, item := range t.Addresses {

			copy(buf[offset+12:offset+32], item[:])

		}
	}

	// Uints (static)

	// Encode fixed-size array t.Uints
	{
		offset := 160
		for _, item := range t.Uints {

			if err := abi.EncodeBigInt(item, buf[offset:offset+32], false); err != nil {
				return 0, err
			}

		}
	}

	// Bytes32s (static)

	// Encode fixed-size array t.Bytes32s
	{
		offset := 256
		for _, item := range t.Bytes32s {

			copy(buf[offset:offset+32], item[:])

		}
	}

	return dynamicOffset, nil
}

// Encode encodes TestFixedArraysArgs to ABI bytes
func (t TestFixedArraysArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes testFixedArrays arguments to ABI bytes including function selector
func (t TestFixedArraysArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], TestFixedArraysArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// TestFixedArraysArgsSelector is the function selector for testFixedArrays(address[5],uint256[3],bytes32[2])
var TestFixedArraysArgsSelector = [4]byte{0x23, 0xb8, 0x46, 0x5c}

// Selector returns the function selector for testFixedArrays
func (TestFixedArraysArgs) Selector() [4]byte {
	return TestFixedArraysArgsSelector
}

// TestMixedTypesArgs represents the arguments for testMixedTypes function

var _ abi.Tuple = TestMixedTypesArgs{}

const TestMixedTypesArgsStaticSize = 160

type TestMixedTypesArgs struct {
	FixedData   [32]byte
	DynamicData []byte
	Flag        bool
	Count       uint8
	Items       []Tuple_de3c4b6f
}

// EncodedSize returns the total encoded size of TestMixedTypesArgs
func (t TestMixedTypesArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.DynamicData)) // length + padded bytes data
	dynamicSize += 32 + 32*len(t.Items)               // length + offset pointers for dynamic elements
	for _, elem := range t.Items {
		dynamicSize += elem.EncodedSize() // dynamic tuple
	}

	return TestMixedTypesArgsStaticSize + dynamicSize
}

// EncodeTo encodes TestMixedTypesArgs to ABI bytes in the provided buffer
func (t TestMixedTypesArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := TestMixedTypesArgsStaticSize // Start dynamic data after static section

	// FixedData (static)
	copy(buf[0:0+32], t.FixedData[:])

	// DynamicData (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// DynamicData (dynamic)

	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.DynamicData)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], t.DynamicData)
	dynamicOffset += abi.Pad32(len(t.DynamicData))

	// Flag (static)

	if t.Flag {
		buf[64+31] = 1
	}

	// Count (static)
	buf[96+31] = byte(t.Count)

	// Items (offset)
	binary.BigEndian.PutUint64(buf[128+24:128+32], uint64(dynamicOffset))

	// Items (dynamic)

	{
		// length
		binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Items)))
		dynamicOffset += 32

		var written int

		// data with dynamic region
		{
			buf := buf[dynamicOffset:]
			dynamicOffset := len(t.Items) * 32 // start after static region

			var offset int
			for _, item := range t.Items {
				// write offsets
				binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynamicOffset))
				offset += 32

				// write data (dynamic)

				n, err := item.EncodeTo(buf[dynamicOffset:])
				if err != nil {
					return 0, err
				}
				dynamicOffset += n

			}
			written = dynamicOffset
		}
		dynamicOffset += written

	}

	return dynamicOffset, nil
}

// Encode encodes TestMixedTypesArgs to ABI bytes
func (t TestMixedTypesArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes testMixedTypes arguments to ABI bytes including function selector
func (t TestMixedTypesArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], TestMixedTypesArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// TestMixedTypesArgsSelector is the function selector for testMixedTypes(bytes32,bytes,bool,uint8,(uint32,bytes,bool)[])
var TestMixedTypesArgsSelector = [4]byte{0x85, 0x8a, 0xe6, 0x15}

// Selector returns the function selector for testMixedTypes
func (TestMixedTypesArgs) Selector() [4]byte {
	return TestMixedTypesArgsSelector
}

// TestNestedDynamicArraysArgs represents the arguments for testNestedDynamicArrays function

var _ abi.Tuple = TestNestedDynamicArraysArgs{}

const TestNestedDynamicArraysArgsStaticSize = 64

type TestNestedDynamicArraysArgs struct {
	Matrix        [][]*big.Int
	AddressMatrix [][]common.Address
}

// EncodedSize returns the total encoded size of TestNestedDynamicArraysArgs
func (t TestNestedDynamicArraysArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + 32*len(t.Matrix) // length + offset pointers for dynamic elements
	for _, elem := range t.Matrix {
		dynamicSize += 32 + 32*len(elem) // length + static elements
	}
	dynamicSize += 32 + 32*len(t.AddressMatrix) // length + offset pointers for dynamic elements
	for _, elem := range t.AddressMatrix {
		dynamicSize += 32 + 32*len(elem) // length + static elements
	}

	return TestNestedDynamicArraysArgsStaticSize + dynamicSize
}

// EncodeTo encodes TestNestedDynamicArraysArgs to ABI bytes in the provided buffer
func (t TestNestedDynamicArraysArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := TestNestedDynamicArraysArgsStaticSize // Start dynamic data after static section

	// Matrix (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Matrix (dynamic)

	{
		// length
		binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Matrix)))
		dynamicOffset += 32

		var written int

		// data with dynamic region
		{
			buf := buf[dynamicOffset:]
			dynamicOffset := len(t.Matrix) * 32 // start after static region

			var offset int
			for _, item := range t.Matrix {
				// write offsets
				binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynamicOffset))
				offset += 32

				// write data (dynamic)

				{
					// length
					binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(item)))
					dynamicOffset += 32

					// data without dynamic region
					buf := buf[dynamicOffset:]
					var offset int
					for _, item := range item {

						if err := abi.EncodeBigInt(item, buf[offset:offset+32], false); err != nil {
							return 0, err
						}

						offset += 32
					}
					dynamicOffset += offset

				}

			}
			written = dynamicOffset
		}
		dynamicOffset += written

	}

	// AddressMatrix (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// AddressMatrix (dynamic)

	{
		// length
		binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.AddressMatrix)))
		dynamicOffset += 32

		var written int

		// data with dynamic region
		{
			buf := buf[dynamicOffset:]
			dynamicOffset := len(t.AddressMatrix) * 32 // start after static region

			var offset int
			for _, item := range t.AddressMatrix {
				// write offsets
				binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynamicOffset))
				offset += 32

				// write data (dynamic)

				{
					// length
					binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(item)))
					dynamicOffset += 32

					// data without dynamic region
					buf := buf[dynamicOffset:]
					var offset int
					for _, item := range item {

						copy(buf[offset+12:offset+32], item[:])

						offset += 32
					}
					dynamicOffset += offset

				}

			}
			written = dynamicOffset
		}
		dynamicOffset += written

	}

	return dynamicOffset, nil
}

// Encode encodes TestNestedDynamicArraysArgs to ABI bytes
func (t TestNestedDynamicArraysArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes testNestedDynamicArrays arguments to ABI bytes including function selector
func (t TestNestedDynamicArraysArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], TestNestedDynamicArraysArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// TestNestedDynamicArraysArgsSelector is the function selector for testNestedDynamicArrays(uint256[][],address[][])
var TestNestedDynamicArraysArgsSelector = [4]byte{0x3d, 0xb1, 0xee, 0x06}

// Selector returns the function selector for testNestedDynamicArrays
func (TestNestedDynamicArraysArgs) Selector() [4]byte {
	return TestNestedDynamicArraysArgsSelector
}

// TestSmallIntegersArgs represents the arguments for testSmallIntegers function

var _ abi.Tuple = TestSmallIntegersArgs{}

const TestSmallIntegersArgsStaticSize = 256

type TestSmallIntegersArgs struct {
	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64
	I8  int8
	I16 int16
	I32 int32
	I64 int64
}

// EncodedSize returns the total encoded size of TestSmallIntegersArgs
func (t TestSmallIntegersArgs) EncodedSize() int {
	dynamicSize := 0

	return TestSmallIntegersArgsStaticSize + dynamicSize
}

// EncodeTo encodes TestSmallIntegersArgs to ABI bytes in the provided buffer
func (t TestSmallIntegersArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := TestSmallIntegersArgsStaticSize // Start dynamic data after static section

	// U8 (static)
	buf[0+31] = byte(t.U8)
	// U16 (static)
	binary.BigEndian.PutUint16(buf[32+30:32+32], uint16(t.U16))
	// U32 (static)
	binary.BigEndian.PutUint32(buf[64+28:64+32], uint32(t.U32))
	// U64 (static)
	binary.BigEndian.PutUint64(buf[96+24:96+32], uint64(t.U64))
	// I8 (static)

	if t.I8 < 0 {
		for i := 0; i < 31; i++ {
			buf[128+i] = 0xff
		}
	}
	buf[128+31] = byte(t.I8)

	// I16 (static)

	if t.I16 < 0 {
		for i := 0; i < 30; i++ {
			buf[160+i] = 0xff
		}
	}
	binary.BigEndian.PutUint16(buf[160+30:160+32], uint16(t.I16))

	// I32 (static)

	if t.I32 < 0 {
		for i := 0; i < 28; i++ {
			buf[192+i] = 0xff
		}
	}
	binary.BigEndian.PutUint32(buf[192+28:192+32], uint32(t.I32))

	// I64 (static)

	if t.I64 < 0 {
		for i := 0; i < 24; i++ {
			buf[224+i] = 0xff
		}
	}
	binary.BigEndian.PutUint64(buf[224+24:224+32], uint64(t.I64))

	return dynamicOffset, nil
}

// Encode encodes TestSmallIntegersArgs to ABI bytes
func (t TestSmallIntegersArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes testSmallIntegers arguments to ABI bytes including function selector
func (t TestSmallIntegersArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], TestSmallIntegersArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// TestSmallIntegersArgsSelector is the function selector for testSmallIntegers(uint8,uint16,uint32,uint64,int8,int16,int32,int64)
var TestSmallIntegersArgsSelector = [4]byte{0x29, 0x2b, 0xd2, 0x39}

// Selector returns the function selector for testSmallIntegers
func (TestSmallIntegersArgs) Selector() [4]byte {
	return TestSmallIntegersArgsSelector
}
