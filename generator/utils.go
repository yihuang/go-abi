package generator

import (
	"cmp"
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

// nativeSize returns the closest native size for a given int size s
func nativeSize(s int) int {
	switch {
	case s <= 8:
		return 8
	case s <= 16:
		return 16
	case s <= 32:
		return 32
	default:
		return 64
	}
}
