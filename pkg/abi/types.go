package abi

// Type represents an ABI type
// This is the base interface for all ABI types
// Code generation will create concrete implementations
// that avoid reflection at runtime
type Type interface {
	// TypeName returns the canonical ABI type name
	TypeName() string

	// IsDynamic returns true if the type is dynamic (variable size)
	IsDynamic() bool

	// StaticSize returns the static size in bytes for fixed-size types
	// Returns 0 for dynamic types
	StaticSize() uint

	// Encode encodes the value to ABI format
	Encode(value interface{}) ([]byte, error)

	// Decode decodes ABI data to the value
	Decode(data []byte, value interface{}) (uint, error)
}

// Basic types that will be code-generated
type (
	Uint8   struct{}
	Uint16  struct{}
	Uint32  struct{}
	Uint64  struct{}
	Uint128 struct{}
	Uint256 struct{}

	Int8   struct{}
	Int16  struct{}
	Int32  struct{}
	Int64  struct{}
	Int128 struct{}
	Int256 struct{}

	Bool    struct{}
	Address struct{}
	Bytes   struct{}
	String  struct{}

	// Array types will be code-generated for specific sizes
	// FixedArray[T, N] for fixed arrays
	// DynamicArray[T] for dynamic arrays

	// Tuple types will be code-generated for specific structs
)

// Common ABI type names
const (
	TypeUint8   = "uint8"
	TypeUint16  = "uint16"
	TypeUint32  = "uint32"
	TypeUint64  = "uint64"
	TypeUint128 = "uint128"
	TypeUint256 = "uint256"

	TypeInt8   = "int8"
	TypeInt16  = "int16"
	TypeInt32  = "int32"
	TypeInt64  = "int64"
	TypeInt128 = "int128"
	TypeInt256 = "int256"

	TypeBool    = "bool"
	TypeAddress = "address"
	TypeBytes   = "bytes"
	TypeString  = "string"
)