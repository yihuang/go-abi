package testdata

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/yihuang/go-abi"
	"math/big"
)

const UserStaticSize = 96

// User represents an ABI tuple
type User struct {
	Address common.Address
	Name    string
	Age     *big.Int
}

// EncodedSize returns the total encoded size of User
func (t User) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.Name)) // length + padded string data

	return UserStaticSize + dynamicSize
}

// EncodeTo encodes User to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t User) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := UserStaticSize // Start dynamic data after static section

	// Address (static)
	copy(buf[0+12:0+32], t.Address[:])

	// Name (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Name (dynamic)
	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Name)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], []byte(t.Name))
	dynamicOffset += abi.Pad32(len(t.Name))

	// Age (static)

	if err := abi.EncodeBigInt(t.Age, buf[64:96], false); err != nil {
		return 0, err
	}

	return dynamicOffset, nil
}

// Encode encodes User to ABI bytes
func (t User) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

const UserDataStaticSize = 64

// UserData represents an ABI tuple
type UserData struct {
	Id   *big.Int
	Data UserMetadata
}

// EncodedSize returns the total encoded size of UserData
func (t UserData) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += t.Data.EncodedSize() // dynamic tuple

	return UserDataStaticSize + dynamicSize
}

// EncodeTo encodes UserData to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t UserData) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := UserDataStaticSize // Start dynamic data after static section

	// Id (static)

	if err := abi.EncodeBigInt(t.Id, buf[0:32], false); err != nil {
		return 0, err
	}

	// Data (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Data (dynamic)
	n, err := t.Data.EncodeTo(buf[dynamicOffset:])
	if err != nil {
		return 0, err
	}
	dynamicOffset += n

	return dynamicOffset, nil
}

// Encode encodes UserData to ABI bytes
func (t UserData) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

const UserMetadataStaticSize = 64

// UserMetadata represents an ABI tuple
type UserMetadata struct {
	Key   [32]byte
	Value string
}

// EncodedSize returns the total encoded size of UserMetadata
func (t UserMetadata) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.Value)) // length + padded string data

	return UserMetadataStaticSize + dynamicSize
}

// EncodeTo encodes UserMetadata to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t UserMetadata) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := UserMetadataStaticSize // Start dynamic data after static section

	// Key (static)
	copy(buf[0:0+32], t.Key[:])

	// Value (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Value (dynamic)
	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Value)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], []byte(t.Value))
	dynamicOffset += abi.Pad32(len(t.Value))

	return dynamicOffset, nil
}

// Encode encodes UserMetadata to ABI bytes
func (t UserMetadata) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

const BalanceOfArgsStaticSize = 32

// BalanceOfArgs represents an ABI tuple
type BalanceOfArgs struct {
	Account common.Address
}

// EncodedSize returns the total encoded size of BalanceOfArgs
func (t BalanceOfArgs) EncodedSize() int {
	dynamicSize := 0

	return BalanceOfArgsStaticSize + dynamicSize
}

// EncodeTo encodes BalanceOfArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t BalanceOfArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := BalanceOfArgsStaticSize // Start dynamic data after static section

	// Account (static)
	copy(buf[0+12:0+32], t.Account[:])

	return dynamicOffset, nil
}

// Encode encodes BalanceOfArgs to ABI bytes
func (t BalanceOfArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes balanceOf arguments to ABI bytes including function selector
func (t BalanceOfArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], BalanceOfArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// BalanceOfArgsSelector is the function selector for balanceOf(address)
var BalanceOfArgsSelector = [4]byte{0x70, 0xa0, 0x82, 0x31}

// Selector returns the function selector for balanceOf
func (BalanceOfArgs) Selector() [4]byte {
	return BalanceOfArgsSelector
}

const BatchProcessArgsStaticSize = 32

// BatchProcessArgs represents an ABI tuple
type BatchProcessArgs struct {
	Users []UserData
}

// EncodedSize returns the total encoded size of BatchProcessArgs
func (t BatchProcessArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + 32*len(t.Users) // length + offset pointers for dynamic elements
	for _, elem := range t.Users {
		dynamicSize += elem.EncodedSize() // dynamic tuple
	}

	return BatchProcessArgsStaticSize + dynamicSize
}

// EncodeTo encodes BatchProcessArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t BatchProcessArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := BatchProcessArgsStaticSize // Start dynamic data after static section

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

// Encode encodes BatchProcessArgs to ABI bytes
func (t BatchProcessArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes batchProcess arguments to ABI bytes including function selector
func (t BatchProcessArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], BatchProcessArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// BatchProcessArgsSelector is the function selector for batchProcess((uint256,(bytes32,string))[])
var BatchProcessArgsSelector = [4]byte{0xb7, 0x78, 0x31, 0x64}

// Selector returns the function selector for batchProcess
func (BatchProcessArgs) Selector() [4]byte {
	return BatchProcessArgsSelector
}

const GetBalancesArgsStaticSize = 320

// GetBalancesArgs represents an ABI tuple
type GetBalancesArgs struct {
	Accounts [10]common.Address
}

// EncodedSize returns the total encoded size of GetBalancesArgs
func (t GetBalancesArgs) EncodedSize() int {
	dynamicSize := 0

	return GetBalancesArgsStaticSize + dynamicSize
}

// EncodeTo encodes GetBalancesArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t GetBalancesArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := GetBalancesArgsStaticSize // Start dynamic data after static section

	// Accounts (static)

	// Encode fixed-size array t.Accounts
	{
		offset := 0
		for _, item := range t.Accounts {

			copy(buf[offset+12:offset+32], item[:])

		}
	}

	return dynamicOffset, nil
}

// Encode encodes GetBalancesArgs to ABI bytes
func (t GetBalancesArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes getBalances arguments to ABI bytes including function selector
func (t GetBalancesArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], GetBalancesArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// GetBalancesArgsSelector is the function selector for getBalances(address[10])
var GetBalancesArgsSelector = [4]byte{0x51, 0x68, 0x3d, 0x7d}

// Selector returns the function selector for getBalances
func (GetBalancesArgs) Selector() [4]byte {
	return GetBalancesArgsSelector
}

const ProcessUserDataArgsStaticSize = 32

// ProcessUserDataArgs represents an ABI tuple
type ProcessUserDataArgs struct {
	User User
}

// EncodedSize returns the total encoded size of ProcessUserDataArgs
func (t ProcessUserDataArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += t.User.EncodedSize() // dynamic tuple

	return ProcessUserDataArgsStaticSize + dynamicSize
}

// EncodeTo encodes ProcessUserDataArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t ProcessUserDataArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := ProcessUserDataArgsStaticSize // Start dynamic data after static section

	// User (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// User (dynamic)
	n, err := t.User.EncodeTo(buf[dynamicOffset:])
	if err != nil {
		return 0, err
	}
	dynamicOffset += n

	return dynamicOffset, nil
}

// Encode encodes ProcessUserDataArgs to ABI bytes
func (t ProcessUserDataArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes processUserData arguments to ABI bytes including function selector
func (t ProcessUserDataArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], ProcessUserDataArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// ProcessUserDataArgsSelector is the function selector for processUserData((address,string,uint256))
var ProcessUserDataArgsSelector = [4]byte{0xc4, 0x9b, 0x1a, 0xe8}

// Selector returns the function selector for processUserData
func (ProcessUserDataArgs) Selector() [4]byte {
	return ProcessUserDataArgsSelector
}

const SetDataArgsStaticSize = 64

// SetDataArgs represents an ABI tuple
type SetDataArgs struct {
	Key   [32]byte
	Value []byte
}

// EncodedSize returns the total encoded size of SetDataArgs
func (t SetDataArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.Value)) // length + padded bytes data

	return SetDataArgsStaticSize + dynamicSize
}

// EncodeTo encodes SetDataArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t SetDataArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := SetDataArgsStaticSize // Start dynamic data after static section

	// Key (static)
	copy(buf[0:0+32], t.Key[:])

	// Value (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Value (dynamic)
	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Value)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], t.Value)
	dynamicOffset += abi.Pad32(len(t.Value))

	return dynamicOffset, nil
}

// Encode encodes SetDataArgs to ABI bytes
func (t SetDataArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes setData arguments to ABI bytes including function selector
func (t SetDataArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], SetDataArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// SetDataArgsSelector is the function selector for setData(bytes32,bytes)
var SetDataArgsSelector = [4]byte{0x7f, 0x23, 0x69, 0x0c}

// Selector returns the function selector for setData
func (SetDataArgs) Selector() [4]byte {
	return SetDataArgsSelector
}

const SetMessageArgsStaticSize = 32

// SetMessageArgs represents an ABI tuple
type SetMessageArgs struct {
	Message string
}

// EncodedSize returns the total encoded size of SetMessageArgs
func (t SetMessageArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.Message)) // length + padded string data

	return SetMessageArgsStaticSize + dynamicSize
}

// EncodeTo encodes SetMessageArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t SetMessageArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := SetMessageArgsStaticSize // Start dynamic data after static section

	// Message (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Message (dynamic)
	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Message)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], []byte(t.Message))
	dynamicOffset += abi.Pad32(len(t.Message))

	return dynamicOffset, nil
}

// Encode encodes SetMessageArgs to ABI bytes
func (t SetMessageArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes setMessage arguments to ABI bytes including function selector
func (t SetMessageArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], SetMessageArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// SetMessageArgsSelector is the function selector for setMessage(string)
var SetMessageArgsSelector = [4]byte{0x36, 0x8b, 0x87, 0x72}

// Selector returns the function selector for setMessage
func (SetMessageArgs) Selector() [4]byte {
	return SetMessageArgsSelector
}

const SmallIntegersArgsStaticSize = 256

// SmallIntegersArgs represents an ABI tuple
type SmallIntegersArgs struct {
	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64
	I8  int8
	I16 int16
	I32 int32
	I64 int64
}

// EncodedSize returns the total encoded size of SmallIntegersArgs
func (t SmallIntegersArgs) EncodedSize() int {
	dynamicSize := 0

	return SmallIntegersArgsStaticSize + dynamicSize
}

// EncodeTo encodes SmallIntegersArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t SmallIntegersArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := SmallIntegersArgsStaticSize // Start dynamic data after static section

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

// Encode encodes SmallIntegersArgs to ABI bytes
func (t SmallIntegersArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes smallIntegers arguments to ABI bytes including function selector
func (t SmallIntegersArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], SmallIntegersArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// SmallIntegersArgsSelector is the function selector for smallIntegers(uint8,uint16,uint32,uint64,int8,int16,int32,int64)
var SmallIntegersArgsSelector = [4]byte{0x98, 0x83, 0xfe, 0x4a}

// Selector returns the function selector for smallIntegers
func (SmallIntegersArgs) Selector() [4]byte {
	return SmallIntegersArgsSelector
}

const TransferArgsStaticSize = 64

// TransferArgs represents an ABI tuple
type TransferArgs struct {
	To     common.Address
	Amount *big.Int
}

// EncodedSize returns the total encoded size of TransferArgs
func (t TransferArgs) EncodedSize() int {
	dynamicSize := 0

	return TransferArgsStaticSize + dynamicSize
}

// EncodeTo encodes TransferArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t TransferArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := TransferArgsStaticSize // Start dynamic data after static section

	// To (static)
	copy(buf[0+12:0+32], t.To[:])
	// Amount (static)

	if err := abi.EncodeBigInt(t.Amount, buf[32:64], false); err != nil {
		return 0, err
	}

	return dynamicOffset, nil
}

// Encode encodes TransferArgs to ABI bytes
func (t TransferArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes transfer arguments to ABI bytes including function selector
func (t TransferArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], TransferArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// TransferArgsSelector is the function selector for transfer(address,uint256)
var TransferArgsSelector = [4]byte{0xa9, 0x05, 0x9c, 0xbb}

// Selector returns the function selector for transfer
func (TransferArgs) Selector() [4]byte {
	return TransferArgsSelector
}

const TransferBatchArgsStaticSize = 64

// TransferBatchArgs represents an ABI tuple
type TransferBatchArgs struct {
	Recipients []common.Address
	Amounts    []*big.Int
}

// EncodedSize returns the total encoded size of TransferBatchArgs
func (t TransferBatchArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + 32*len(t.Recipients) // length + static elements
	dynamicSize += 32 + 32*len(t.Amounts)    // length + static elements

	return TransferBatchArgsStaticSize + dynamicSize
}

// EncodeTo encodes TransferBatchArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t TransferBatchArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := TransferBatchArgsStaticSize // Start dynamic data after static section

	// Recipients (offset)
	binary.BigEndian.PutUint64(buf[0+24:0+32], uint64(dynamicOffset))

	// Recipients (dynamic)
	{
		// length
		binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Recipients)))
		dynamicOffset += 32

		// data without dynamic region
		buf := buf[dynamicOffset:]
		var offset int
		for _, item := range t.Recipients {

			copy(buf[offset+12:offset+32], item[:])

			offset += 32
		}
		dynamicOffset += offset

	}

	// Amounts (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Amounts (dynamic)
	{
		// length
		binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Amounts)))
		dynamicOffset += 32

		// data without dynamic region
		buf := buf[dynamicOffset:]
		var offset int
		for _, item := range t.Amounts {

			if err := abi.EncodeBigInt(item, buf[offset:offset+32], false); err != nil {
				return 0, err
			}

			offset += 32
		}
		dynamicOffset += offset

	}

	return dynamicOffset, nil
}

// Encode encodes TransferBatchArgs to ABI bytes
func (t TransferBatchArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes transferBatch arguments to ABI bytes including function selector
func (t TransferBatchArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], TransferBatchArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// TransferBatchArgsSelector is the function selector for transferBatch(address[],uint256[])
var TransferBatchArgsSelector = [4]byte{0x3b, 0x3e, 0x67, 0x2f}

// Selector returns the function selector for transferBatch
func (TransferBatchArgs) Selector() [4]byte {
	return TransferBatchArgsSelector
}

const UpdateProfileArgsStaticSize = 96

// UpdateProfileArgs represents an ABI tuple
type UpdateProfileArgs struct {
	User common.Address
	Name string
	Age  *big.Int
}

// EncodedSize returns the total encoded size of UpdateProfileArgs
func (t UpdateProfileArgs) EncodedSize() int {
	dynamicSize := 0

	dynamicSize += 32 + abi.Pad32(len(t.Name)) // length + padded string data

	return UpdateProfileArgsStaticSize + dynamicSize
}

// EncodeTo encodes UpdateProfileArgs to ABI bytes in the provided buffer
// it panics if the buffer is not large enough
func (t UpdateProfileArgs) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := UpdateProfileArgsStaticSize // Start dynamic data after static section

	// User (static)
	copy(buf[0+12:0+32], t.User[:])

	// Name (offset)
	binary.BigEndian.PutUint64(buf[32+24:32+32], uint64(dynamicOffset))

	// Name (dynamic)
	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(t.Name)))
	dynamicOffset += 32

	// data
	copy(buf[dynamicOffset:], []byte(t.Name))
	dynamicOffset += abi.Pad32(len(t.Name))

	// Age (static)

	if err := abi.EncodeBigInt(t.Age, buf[64:96], false); err != nil {
		return 0, err
	}

	return dynamicOffset, nil
}

// Encode encodes UpdateProfileArgs to ABI bytes
func (t UpdateProfileArgs) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// EncodeWithSelector encodes updateProfile arguments to ABI bytes including function selector
func (t UpdateProfileArgs) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4+t.EncodedSize())
	copy(result[:4], UpdateProfileArgsSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateProfileArgsSelector is the function selector for updateProfile(address,string,uint256)
var UpdateProfileArgsSelector = [4]byte{0x6d, 0xe9, 0x52, 0x01}

// Selector returns the function selector for updateProfile
func (UpdateProfileArgs) Selector() [4]byte {
	return UpdateProfileArgsSelector
}
