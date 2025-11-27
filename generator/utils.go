package generator

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var Title = cases.Title(language.English, cases.NoLower)

func ToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		parts[i] = Title.String(part)
	}
	return strings.Join(parts, "")
}

func ToArgName(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func SortedMapKeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// IsDynamicType returns true if the type is dynamic.
// The following types are called “dynamic”:
// * bytes
// * string
// * T[] for any T
// * T[k] for any dynamic T and any k >= 0
// * (T1,...,Tk) if Ti is dynamic for some 1 <= i <= k
func IsDynamicType(t abi.Type) bool {
	if t.T == abi.TupleTy {
		for _, elem := range t.TupleElems {
			if IsDynamicType(*elem) {
				return true
			}
		}
		return false
	}
	return t.T == abi.StringTy || t.T == abi.BytesTy || t.T == abi.SliceTy || (t.T == abi.ArrayTy && IsDynamicType(*t.Elem))
}

// GetTypeSize returns the size that this type needs to occupy.
// We distinguish static and dynamic types. Static types are encoded in-place
// and dynamic types are encoded at a separately allocated location after the
// current block.
// So for a static variable, the size returned represents the size that the
// variable actually occupies.
// For a dynamic variable, the returned size is fixed 32 bytes, which is used
// to store the location reference for actual value storage.
func GetTypeSize(t abi.Type) int {
	if t.T == abi.ArrayTy && !IsDynamicType(*t.Elem) {
		// Recursively calculate type size if it is a nested array
		if t.Elem.T == abi.ArrayTy || t.Elem.T == abi.TupleTy {
			return t.Size * GetTypeSize(*t.Elem)
		}
		return t.Size * 32
	} else if t.T == abi.TupleTy && !IsDynamicType(t) {
		total := 0
		for _, elem := range t.TupleElems {
			total += GetTypeSize(*elem)
		}
		return total
	}
	return 32
}

func GetTupleSize(elems []*abi.Type) int {
	total := 0
	for _, elem := range elems {
		total += GetTypeSize(*elem)
	}
	return total
}

// GetPackedTypeSize returns the packed size of a type (without padding)
func GetPackedTypeSize(t abi.Type) int {
	switch t.T {
	case abi.UintTy, abi.IntTy:
		// Integer types use their natural size
		return (t.Size + 7) / 8 // Convert bits to bytes
	case abi.AddressTy:
		// Address is 20 bytes
		return 20
	case abi.BoolTy:
		// Boolean is 1 byte
		return 1
	case abi.FixedBytesTy:
		// Fixed bytes use their actual size
		return t.Size
	case abi.ArrayTy:
		// Fixed arrays multiply element size by count
		return t.Size * GetPackedTypeSize(*t.Elem)
	case abi.TupleTy:
		// Tuples sum the sizes of their elements
		total := 0
		for _, elem := range t.TupleElems {
			total += GetPackedTypeSize(*elem)
		}
		return total
	case abi.SliceTy, abi.BytesTy, abi.StringTy:
		// Dynamic types are not supported in packed format
		// This should be caught by validation before reaching here
		panic(fmt.Sprintf("dynamic type %s not supported in packed format", t.String()))
	default:
		panic(fmt.Sprintf("unsupported ABI type for packed size calculation: %s", t.String()))
	}
}

// RequiresLengthPrefix returns whether the type requires any sort of length
// prefixing.
func RequiresLengthPrefix(t abi.Type) bool {
	return t.T == abi.StringTy || t.T == abi.BytesTy || t.T == abi.SliceTy
}

func VisitABIType(t abi.Type, visit func(abi.Type)) {
	visit(t)
	switch t.T {
	case abi.TupleTy:
		for _, elem := range t.TupleElems {
			VisitABIType(*elem, visit)
		}
	case abi.ArrayTy, abi.SliceTy:
		VisitABIType(*t.Elem, visit)
	}
}

// GoFieldName converts abi field name to a valid Go field name
func GoFieldName(name string) string {
	name = strings.TrimPrefix(name, "_")
	return Title.String(name)
}

// ParseExternalTuples parses external tuple mappings from string format
// Format: "key1=value1,key2=value2"
func ParseExternalTuples(s string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}

	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" && value != "" {
				result[key] = value
			}
		}
	}
	return result
}

// ParseImport parses an import string that may contain an alias
// Examples:
//
//	"github.com/ethereum/go-ethereum/common" -> ImportSpec{Path: "github.com/ethereum/go-ethereum/common", Alias: ""}
//	"cmn=github.com/ethereum/go-ethereum/common" -> ImportSpec{Path: "github.com/ethereum/go-ethereum/common", Alias: "cmn"}
func ParseImport(imp string) ImportSpec {
	parts := strings.Split(imp, "=")
	var spec ImportSpec

	switch len(parts) {
	case 2:
		spec = ImportSpec{
			Alias: parts[0],
			Path:  parts[1],
		}
	case 1:
		spec = ImportSpec{
			Path: parts[0],
		}
	default:
		panic("invalid import format " + imp)
	}
	return spec
}

// IsPackedSupported returns true if the type can be encoded in packed format
func IsPackedSupported(t abi.Type) bool {
	switch t.T {
	case abi.UintTy, abi.IntTy, abi.AddressTy, abi.BoolTy, abi.FixedBytesTy:
		// Static primitive types are supported
		return true
	case abi.ArrayTy:
		// Fixed arrays are supported if their element type is supported
		return IsPackedSupported(*t.Elem)
	case abi.TupleTy:
		// Tuples are supported if all their elements are supported
		for _, elem := range t.TupleElems {
			if !IsPackedSupported(*elem) {
				return false
			}
		}
		return true
	case abi.SliceTy, abi.BytesTy, abi.StringTy:
		// Dynamic types are not supported in packed format
		return false
	default:
		return false
	}
}
