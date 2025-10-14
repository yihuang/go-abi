package abi

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"go/format"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

var Title = cases.Title(language.English, cases.NoLower)

// Generator handles ABI code generation
type Generator struct {
	buf bytes.Buffer

	PackageName string
	Imports     []string
}

// NewGenerator creates a new ABI code generator
func NewGenerator(packageName string) *Generator {
	return &Generator{
		PackageName: packageName,
		Imports: []string{
			"math/big",
			"github.com/ethereum/go-ethereum/common",
			"github.com/yihuang/go-abi",
		},
	}
}

func (g *Generator) L(format string, args ...any) {
	fmt.Fprintf(&g.buf, format, args...)
	fmt.Fprint(&g.buf, "\n")
}

// GenerateFromABI generates Go code from ABI JSON
func (g *Generator) GenerateFromABI(abiDef abi.ABI) (string, error) {
	// Write package declaration
	g.L("package %s", g.PackageName)

	// Check if we need encoding/binary import for optimized integer encoding
	// We always need it for offset/length encoding, and also for 8,16,32,64-bit integers
	needsBinaryImport := false
	for _, method := range abiDef.Methods {
		for _, input := range method.Inputs {
			// Check for integer types that need binary encoding
			if (input.Type.T == abi.UintTy || input.Type.T == abi.IntTy) &&
				(input.Type.Size == 8 || input.Type.Size == 16 || input.Type.Size == 32 || input.Type.Size == 64) {
				needsBinaryImport = true
				break
			}
			// Check for dynamic types that need offset/length encoding
			if isDynamicType(input.Type) {
				needsBinaryImport = true
				break
			}
		}
		if needsBinaryImport {
			break
		}
	}

	// Write imports
	imports := make([]string, len(g.Imports))
	copy(imports, g.Imports)
	if needsBinaryImport {
		imports = append(imports, "encoding/binary")
	}

	if len(imports) > 0 {
		g.L("import (")
		for _, imp := range imports {
			if strings.Contains(imp, "/") {
				g.L("\"%s\"", imp)
			} else {
				g.L("%s", imp)
			}
		}
		g.L(")")
	}

	// First, generate all tuple structs needed for this function
	var methods []abi.Method
	for _, name := range SortedMapKeys(abiDef.Methods) {
		methods = append(methods, abiDef.Methods[name])
	}

	if err := g.genTuples(methods); err != nil {
		return "", err
	}

	// Generate code for each function
	for _, method := range methods {
		if err := g.genFunction(method); err != nil {
			return "", fmt.Errorf("failed to generate function %s: %w", method.Name, err)
		}
	}

	// Format the generated code
	src := g.buf.Bytes()
	formatted, err := format.Source(src)
	if err != nil {
		return string(src), fmt.Errorf("failed to format generated code: %w", err)
	}

	return string(formatted), nil
}

func (g *Generator) genStruct(s Struct) error {
	g.L(`
var _ abi.Tuple = %s{}

const %sStaticSize = %d

type %s struct {
`, s.Name, s.Name, getTupleSize(s.Types()), s.Name)

	for _, f := range s.Fields {
		goType, err := abiTypeToGoType(*f.Type)
		if err != nil {
			return err
		}
		g.L("%s %s", f.Name, goType)
	}
	g.L("}")
	return nil
}

// genFunction generates Go code for a single function
func (g *Generator) genFunction(method abi.Method) error {
	if len(method.Inputs) == 0 {
		return nil
	}

	s := StructFromInputs(method)

	// Generate struct for function arguments
	g.L("// %s represents the arguments for %s function", s.Name, method.Name)

	if err := g.genStruct(s); err != nil {
		return err
	}

	g.genStructMethods(s)

	// function sepecific methods
	g.L(`
// EncodeWithSelector encodes %s arguments to ABI bytes including function selector
func (t %s) EncodeWithSelector() ([]byte, error) {
	result := make([]byte, 4 + t.EncodedSize())
	copy(result[:4], %sSelector[:])
	if _, err := t.EncodeTo(result[4:]); err != nil {
		return nil, err
	}
	return result, nil
}
`, method.Name, s.Name, s.Name)

	// Generate selector
	g.L("// %sSelector is the function selector for %s", s.Name, method.Sig)
	g.L("var %sSelector = [4]byte{0x%02x, 0x%02x, 0x%02x, 0x%02x}",
		s.Name,
		method.ID[0],
		method.ID[1],
		method.ID[2],
		method.ID[3])

	g.L(`
// Selector returns the function selector for %s
func (%s) Selector() [4]byte {
	return %sSelector
}
`, method.Name, s.Name, s.Name)

	return nil
}

// abiTypeToGoType converts ABI type to Go type
func abiTypeToGoType(abiType abi.Type) (string, error) {
	switch abiType.T {
	case abi.UintTy:
		// Use native Go types for common sizes to avoid big.Int allocations
		switch abiType.Size {
		case 8, 16, 32, 64:
			return fmt.Sprintf("uint%d", abiType.Size), nil
		default:
			return "*big.Int", nil
		}
	case abi.IntTy:
		// Use native Go types for common sizes to avoid big.Int allocations
		switch abiType.Size {
		case 8, 16, 32, 64:
			return fmt.Sprintf("int%d", abiType.Size), nil
		default:
			return "*big.Int", nil
		}
	case abi.AddressTy:
		return "common.Address", nil
	case abi.BoolTy:
		return "bool", nil
	case abi.StringTy:
		return "string", nil
	case abi.BytesTy:
		return "[]byte", nil
	case abi.FixedBytesTy:
		return fmt.Sprintf("[%d]byte", abiType.Size), nil
	case abi.SliceTy:
		// Dynamic arrays like uint256[]
		if abiType.Elem == nil {
			return "", fmt.Errorf("invalid slice type: nil element")
		}
		elemType, err := abiTypeToGoType(*abiType.Elem)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[]%s", elemType), nil
	case abi.ArrayTy:
		// Fixed-size arrays like uint256[10]
		if abiType.Elem == nil {
			return "", fmt.Errorf("invalid array type: nil element")
		}
		elemType, err := abiTypeToGoType(*abiType.Elem)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%d]%s", abiType.Size, elemType), nil
	case abi.TupleTy:
		// Handle tuple types - generate struct type name
		return tupleStructName(abiType), nil
	default:
		return "", fmt.Errorf("unsupported ABI type: %s", abiType.String())
	}
}

// tupleStructName generates a unique struct name for a tuple type
func tupleStructName(t abi.Type) string {
	if t.TupleRawName != "" {
		return t.TupleRawName
	}

	// Use the tuple's string representation as the basis for the struct name
	// This creates a deterministic name based on the tuple structure
	return genTupleIdentifier(t)
}

// genTupleIdentifier generates a unique identifier for a tuple type
func genTupleIdentifier(t abi.Type) string {
	// Create a signature based on tuple element types
	types := make([]string, len(t.TupleElems))
	for i, elem := range t.TupleElems {
		types[i] = elem.String()
	}
	sig := fmt.Sprintf("(%v)", strings.Join(types, ","))

	id := crypto.Keccak256([]byte(sig))
	return "Tuple_" + hex.EncodeToString(id)[:8] // Use first 8 chars for readability
}

// genTuples generates all tuple structs needed for a function
func (g *Generator) genTuples(methods []abi.Method) error {
	// Collect all tuple types from function inputs
	tupleTypes := make(map[string]abi.Type)

	var collectTuples func(t abi.Type)
	collectTuples = func(t abi.Type) {
		if t.T == abi.TupleTy {
			structName := tupleStructName(t)
			if _, exists := tupleTypes[structName]; !exists {
				tupleTypes[structName] = t
				// Recursively collect nested tuples
				for _, elem := range t.TupleElems {
					collectTuples(*elem)
				}
			}
		} else if t.T == abi.ArrayTy || t.T == abi.SliceTy {
			// Check array elements for tuples
			if t.Elem != nil {
				collectTuples(*t.Elem)
			}
		}
	}

	// Collect tuples from all methods
	for _, method := range methods {
		// Collect tuples from all inputs
		for _, input := range method.Inputs {
			collectTuples(input.Type)
		}
	}

	// Generate struct definitions for collected tuples
	for _, name := range SortedMapKeys(tupleTypes) {
		s := StructFromTuple(tupleTypes[name])
		g.L("// %s represents a tuple type", name)

		if err := g.genStruct(s); err != nil {
			return err
		}

		// Generate encode method for the tuple struct
		g.genStructMethods(s)
	}

	return nil
}

// genSize generates size calculation logic for a type
func (g *Generator) genSize(t abi.Type, acc string, ref string) {
	if !isDynamicType(t) {
		g.L("%s += %d // static element %s", acc, getTypeSize(t), ref)
		return
	}

	switch t.T {
	case abi.StringTy:
		g.L("%s += 32 + abi.Pad32(len(%s)) // length + padded string data", acc, ref)

	case abi.BytesTy:
		g.L("%s += 32 + abi.Pad32(len(%s)) // length + padded bytes data", acc, ref)

	case abi.SliceTy:
		if isDynamicType(*t.Elem) {
			// Dynamic array with dynamic elements
			g.L("%s += 32 + 32 * len(%s) // length + offset pointers for dynamic elements", acc, ref)
			g.L("for _, elem := range %s {", ref)
			g.genSize(*t.Elem, acc, "elem")
			g.L("}")
		} else {
			// Dynamic array with static elements
			g.L("%s += 32 + %d * len(%s) // length + static elements", acc, getTypeSize(*t.Elem), ref)
		}

	case abi.ArrayTy:
		// Fixed size array of dynamic element types
		g.L("for _, elem := range %s {", ref)
		g.genSize(*t.Elem, acc, "elem")
		g.L("}")

	case abi.TupleTy:
		// Dynamic tuple, just call tuple struct method
		g.L("%s += %s.EncodedSize() // dynamic tuple", acc, ref)

	default:
		panic("impossible")
	}
}

// genEncodedSize generates the size calculation logic without selector
func (g *Generator) genEncodedSize(s Struct) {
	g.L(`
// EncodedSize returns the total encoded size of %s
func (t %s) EncodedSize() int {
	dynamicSize := 0
`, s.Name, s.Name)

	for _, f := range s.Fields {
		if !isDynamicType(*f.Type) {
			continue
		}
		g.genSize(*f.Type, "dynamicSize", "t."+f.Name)
	}

	g.L(`
	return %sStaticSize + dynamicSize
}`, s.Name)
}

// genEncodedTo the `EncodeTo(buf []byte) (int, error)` method
func (g *Generator) genEncodedTo(s Struct) {
	g.L(`
// EncodeTo encodes %s to ABI bytes in the provided buffer
func (t %s) EncodeTo(buf []byte) (int, error) {
	dynamicOffset := %sStaticSize // Start dynamic data after static section
`, s.Name, s.Name, s.Name)

	var offset int
	for _, f := range s.Fields {
		if !isDynamicType(*f.Type) {
			g.L("// %s (static)", f.Name)
			offset = g.genStaticItem("t."+f.Name, *f.Type, offset)
			continue
		}

		g.L(`
	// %s (offset)
	binary.BigEndian.PutUint64(buf[%d+24:%d+32], uint64(dynamicOffset))
`, f.Name, offset, offset)

		// Generate encoding for dynamic element
		g.L("// %s (dynamic)", f.Name)
		g.genDynamicItem(fmt.Sprintf("t.%s", f.Name), *f.Type)

		offset += 32
	}

	g.L(`
	return dynamicOffset, nil
}
`)
}

// genStructMethods generates an Encode method for tuple structs
func (g *Generator) genStructMethods(s Struct) {
	g.genEncodedSize(s)
	g.genEncodedTo(s)

	g.L(`// Encode encodes %s to ABI bytes
func (t %s) Encode() ([]byte, error) {
	buf := make([]byte, t.EncodedSize())
	if _, err := t.EncodeTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}
`, s.Name, s.Name)
}

func (g *Generator) genIntOffset(ref string, t abi.Type) {
	// Check if we can use native Go types for optimization
	switch t.Size {
	case 8:
		if t.T == abi.IntTy {
			// int8 - sign extend to 32 bytes
			g.L(`
if %s < 0 {
	for i := 0; i < 31; i++ { buf[offset+i] = 0xff }
}
buf[offset+31] = byte(%s)
`, ref, ref)
		} else {
			// uint8 - zero extend to 32 bytes
			g.L("buf[offset+31] = byte(%s)", ref)
		}
	case 16:
		if t.T == abi.IntTy {
			// int16 - sign extend to 32 bytes
			g.L(`
if %s < 0 {
	for i := 0; i < 30; i++ { buf[offset+i] = 0xff }
}
binary.BigEndian.PutUint16(buf[offset+30:offset+32], uint16(%s))
`, ref, ref)
		} else {
			// uint16 - zero extend to 32 bytes
			g.L("binary.BigEndian.PutUint16(buf[offset+30:offset+32], uint16(%s))", ref)
		}
	case 24, 32:
		if t.T == abi.IntTy {
			// int32 - sign extend to 32 bytes
			g.L(`
if %s < 0 {
	for i := 0; i < 28; i++ { buf[offset+i] = 0xff }
}
binary.BigEndian.PutUint32(buf[offset+28:offset+32], uint32(%s))
`, ref, ref)
		} else {
			// uint32 - zero extend to 32 bytes
			g.L("binary.BigEndian.PutUint32(buf[offset+28:offset+32], uint32(%s))", ref)
		}
	case 40, 48, 56, 64:
		if t.T == abi.IntTy {
			// int64 - sign extend to 32 bytes
			g.L(`
if %s < 0 {
	for i := 0; i < 24; i++ { buf[offset+i] = 0xff }
}
binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(%s))
`, ref, ref)
		} else {
			// uint64 - zero extend to 32 bytes
			g.L("binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(%s))", ref)
		}
	default:
		signed := "false"
		if t.T == abi.IntTy {
			signed = "true"
		}

		g.L(`
if err := abi.EncodeBigInt(%s, buf[offset:offset+32], %s); err != nil {
	return 0, err
}
`, ref, signed)
	}
}

func (g *Generator) genInt(ref string, t abi.Type, offset int) int {
	// Check if we can use native Go types for optimization
	switch t.Size {
	case 8:
		if t.T == abi.IntTy {
			// int8 - sign extend to 32 bytes
			g.L(`
if %s < 0 {
	for i := 0; i < 31; i++ { buf[%d+i] = 0xff }
}
buf[%d+31] = byte(%s)
`, ref, offset, offset, ref)
		} else {
			// uint8 - zero extend to 32 bytes
			g.L("buf[%d+31] = byte(%s)", offset, ref)
		}

	case 16:
		if t.T == abi.IntTy {
			// int16 - sign extend to 32 bytes
			g.L(`
if %s < 0 {
	for i := 0; i < 30; i++ { buf[%d+i] = 0xff }
}
binary.BigEndian.PutUint16(buf[%d+30:%d+32], uint16(%s))
`, ref, offset, offset, offset, ref)
		} else {
			// uint16 - zero extend to 32 bytes
			g.L("binary.BigEndian.PutUint16(buf[%d+30:%d+32], uint16(%s))", offset, offset, ref)
		}

	case 24, 32:
		if t.T == abi.IntTy {
			// int32 - sign extend to 32 bytes
			g.L(`
if %s < 0 {
	for i := 0; i < 28; i++ { buf[%d+i] = 0xff }
}
binary.BigEndian.PutUint32(buf[%d+28:%d+32], uint32(%s))
`, ref, offset, offset, offset, ref)
		} else {
			// uint32 - zero extend to 32 bytes
			g.L("binary.BigEndian.PutUint32(buf[%d+28:%d+32], uint32(%s))", offset, offset, ref)
		}

	case 40, 48, 56, 64:
		if t.T == abi.IntTy {
			// int64 - sign extend to 32 bytes
			g.L(`
if %s < 0 {
	for i := 0; i < 24; i++ { buf[%d+i] = 0xff }
}
binary.BigEndian.PutUint64(buf[%d+24:%d+32], uint64(%s))
`, ref, offset, offset, offset, ref)
		} else {
			// uint64 - zero extend to 32 bytes
			g.L("binary.BigEndian.PutUint64(buf[%d+24:%d+32], uint64(%s))", offset, offset, ref)
		}

	default:
		signed := "false"
		if t.T == abi.IntTy {
			signed = "true"
		}
		g.L(`
if err := abi.EncodeBigInt(%s, buf[%d:%d], %s); err != nil {
	return 0, err
}
`, ref, offset, offset+32, signed)

	}

	return offset + 32
}

// genStaticItemOffset generates encoding for a single tuple element in tuple Encode method
func (g *Generator) genStaticItemOffset(ref string, t abi.Type) {
	switch t.T {
	case abi.AddressTy:
		g.L("copy(buf[offset+12:offset+32], %s[:])", ref)

	case abi.UintTy, abi.IntTy:
		g.genIntOffset(ref, t)

	case abi.BoolTy:
		g.L(`
if %s {
	buf[offset+31] = 1
}
`, ref)

	case abi.FixedBytesTy:
		g.L("copy(buf[offset:offset+32], %s[:])", ref)

	case abi.ArrayTy:
		// Fixed-size array with static elements
		elemSize := getTypeSize(*t.Elem)
		g.L(`
// Encode fixed-size array %s
for _, item := range %s {
`, ref, ref)

		g.genStaticItemOffset("item", *t.Elem)

		g.L(`
	offset += %d
}
`, elemSize)

	case abi.TupleTy:
		// Nested static tuple - use the generated EncodeTo method
		g.L(`
// Encode nested tuple %s
if _, err := %s.EncodeTo(buf[offset:]); err != nil {
	return 0, err
}
`, ref, ref)

	default:
		panic("unknown static type")
	}
}

// genStaticItem generates encoding for a single tuple element in tuple Encode method
func (g *Generator) genStaticItem(ref string, elemType abi.Type, offset int) int {
	switch elemType.T {
	case abi.AddressTy:
		g.L("copy(buf[%d+12:%d+32], %s[:])", offset, offset, ref)
		offset += 32

	case abi.UintTy, abi.IntTy:
		offset = g.genInt(ref, elemType, offset)

	case abi.BoolTy:
		g.L(`
if %s {
	buf[%d+31] = 1
}
`, ref, offset)
		offset += 32

	case abi.FixedBytesTy:
		g.L("copy(buf[%d:%d+32], %s[:])", offset, offset, ref)
		offset += 32

	case abi.ArrayTy:
		// Fixed-size array with static elements
		g.L(`
// Encode fixed-size array %s
{
	offset := %d
	for _, item := range %s {
`, ref, offset, ref)

		g.genStaticItemOffset("item", *elemType.Elem)

		g.L(`
	}
}
`)

		offset += elemType.Size * getTypeSize(*elemType.Elem)

	case abi.TupleTy:
		// Nested static tuple - use the generated EncodeTo method
		g.L(`
// Encode nested tuple %s
if _, err := %s.EncodeTo(buf[%d:]); err != nil {
	return 0, err
}
`, ref, ref, offset)

		offset += getTypeSize(elemType)

	default:
		panic("unknown static type")
	}

	return offset
}

// genDynamicItem generate code to write element into the dynamic section
func (g *Generator) genDynamicItem(ref string, t abi.Type) {
	switch t.T {
	case abi.StringTy:
		g.L(`
// length
binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(%s)))
dynamicOffset += 32

// data
copy(buf[dynamicOffset:], []byte(%s))
dynamicOffset += abi.Pad32(len(%s))
`, ref, ref, ref)

	case abi.BytesTy:
		g.L(`
// length
binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(%s)))
dynamicOffset += 32

// data
copy(buf[dynamicOffset:], %s)
dynamicOffset += abi.Pad32(len(%s))
`, ref, ref, ref)

	case abi.TupleTy:
		g.L(`
n, err := %s.EncodeTo(buf[dynamicOffset:])
if err != nil {
	return 0, err
}
dynamicOffset += n
`, ref)

	case abi.SliceTy:
		g.L(`
{
	// length
	binary.BigEndian.PutUint64(buf[dynamicOffset+24:dynamicOffset+32], uint64(len(%s)))
	dynamicOffset += 32
`, ref)

		if isDynamicType(*t.Elem) {
			g.L(`
	var written int

	// data with dynamic region
	{
		buf := buf[dynamicOffset:]
		dynamicOffset := len(%s) * 32 // start after static region

		var offset int
		for _, item := range %s {
			// write offsets
			binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynamicOffset))
			offset += 32

			// write data (dynamic)
`, ref, ref)

			g.genDynamicItem("item", *t.Elem)

			g.L(`
		}
		written = dynamicOffset
	}
	dynamicOffset += written
`)
		} else {
			elemSize := getTypeSize(*t.Elem)
			g.L(`
	// data without dynamic region
	buf := buf[dynamicOffset:]
	var offset int
	for _, item := range %s {
`, ref)
			g.genStaticItemOffset("item", *t.Elem)
			g.L(`
		offset += %d
	}
	dynamicOffset += offset
`, elemSize)
		}

		g.L("}")

	case abi.ArrayTy:
		g.L(`
{
	var written int

	// data with dynamic region
	{
		buf := buf[dynamicOffset:]
		dynamicOffset := %d * 32 // start after static region

		var offset int
		for _, item := range %s {
			// write offsets
			binary.BigEndian.PutUint64(buf[offset+24:offset+32], uint64(dynamicOffset))
			offset += 32

			// write data (dynamic)
`, t.Size, ref)

		g.genDynamicItem("item", *t.Elem)

		g.L(`
		}
		written = dynamicOffset
	}
	dynamicOffset += written
}
`)

	default:
		panic("unknown dynamic type")
	}
}
