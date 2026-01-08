package generator

import (
	"fmt"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

// sliceViewTypeName returns the view type name for a slice
// e.g., []Item -> "ItemSliceView", []string -> "StringSliceView"
func sliceViewTypeName(t ethabi.Type) string {
	return abi.GenTypeIdentifier(t) + "View"
}

// genSliceView generates a SliceView type for a slice
func (g *Generator) genSliceView(t ethabi.Type) {
	if t.T != ethabi.SliceTy {
		panic("genSliceView called on non-slice type")
	}

	// Skip stdlib slice types when not generating stdlib
	// (they would conflict with other generated files in the same package)
	typeID := abi.GenTypeIdentifier(t)
	if !g.Options.Stdlib && abi.IsStdlibType(typeID) {
		return
	}

	typeName := sliceViewTypeName(t) // e.g., "ItemSliceView" for []Item
	elemType := g.sliceViewElemReturnType(*t.Elem)
	hasDynamicElem := IsDynamicType(*t.Elem)

	g.L("")
	g.L("// %s provides lazy indexed access to %s", typeName, t.String())
	g.L("type %s struct {", typeName)
	g.L("\tdata []byte")
	g.L("\tlength int")
	if hasDynamicElem {
		g.L("\toffsets []int // offset for each element")
	}
	g.L("}")

	// Generate DecodeXxxSliceView
	g.genSliceViewDecodeFunction(t, typeName)

	// Generate Len() method
	g.L("")
	g.L("// Len returns the number of elements")
	g.L("func (v *%s) Len() int {", typeName)
	g.L("\treturn v.length")
	g.L("}")

	// Generate Get(i) method
	g.genSliceViewGet(t, typeName, elemType)

	// Generate Raw() method
	g.L("")
	g.L("// Raw returns the underlying encoded bytes")
	g.L("func (v *%s) Raw() []byte {", typeName)
	g.L("\treturn v.data")
	g.L("}")

	// Generate Materialize() method
	baseElemType := g.abiTypeToGoType(*t.Elem)
	g.L("")
	g.L("// Materialize fully decodes all elements into a slice")
	g.L("func (v *%s) Materialize() ([]%s, error) {", typeName, baseElemType)
	g.L("\tresult, _, err := %s", g.genDecodeCall(t, "v.data"))
	g.L("\treturn result, err")
	g.L("}")
}

// sliceViewElemReturnType returns the return type for Get(i)
func (g *Generator) sliceViewElemReturnType(t ethabi.Type) string {
	switch t.T {
	case ethabi.TupleTy:
		return "*" + abi.TupleStructName(t) + "View"
	case ethabi.SliceTy:
		// If the slice type is a stdlib type (and we're not in stdlib mode),
		// return the regular Go type instead of the view type
		typeID := abi.GenTypeIdentifier(t)
		if !g.Options.Stdlib && abi.IsStdlibType(typeID) {
			return g.abiTypeToGoType(t)
		}
		return "*" + sliceViewTypeName(t) // e.g., "*ItemSliceView"
	default:
		return g.abiTypeToGoType(t)
	}
}

// genSliceViewDecodeFunction generates DecodeXxx function for slice views
func (g *Generator) genSliceViewDecodeFunction(t ethabi.Type, typeName string) {
	hasDynamicElem := IsDynamicType(*t.Elem)

	g.L("")
	g.L("// Decode%s creates a lazy view of %s", typeName, t.String())
	g.L("func Decode%s(data []byte) (*%s, int, error) {", typeName, typeName)

	// Read length
	g.L("\tif len(data) < 32 {")
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")
	g.L("\tlength, err := %sDecodeSize(data)", g.StdPrefix)
	g.L("\tif err != nil {")
	g.L("\t\treturn nil, 0, err")
	g.L("\t}")

	if !hasDynamicElem {
		// Static elements - simple size calculation
		elemSize := GetTypeSize(*t.Elem)
		g.L("\ttotalSize := 32 + length * %d", elemSize)
		g.L("\tif len(data) < totalSize {")
		g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
		g.L("\t}")
		g.L("\treturn &%s{", typeName)
		g.L("\t\tdata: data[:totalSize],")
		g.L("\t\tlength: length,")
		g.L("\t}, totalSize, nil")
	} else {
		// Dynamic elements - parse offset table and validate
		g.L("\tif length == 0 {")
		g.L("\t\treturn &%s{data: data[:32], length: 0, offsets: nil}, 32, nil", typeName)
		g.L("\t}")
		g.L("")
		g.L("\tminSize := 32 + length * 32")
		g.L("\tif len(data) < minSize {")
		g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
		g.L("\t}")
		g.L("")
		g.L("\toffsets := make([]int, length)")
		g.L("\tdynamicOffset := length * 32")
		g.L("\tfor i := 0; i < length; i++ {")
		g.L("\t\toffset, err := %sDecodeSize(data[32 + i*32:])", g.StdPrefix)
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")
		g.L("\t\tif offset != dynamicOffset {")
		g.L("\t\t\treturn nil, 0, %sErrInvalidOffsetForSliceElement", g.StdPrefix)
		g.L("\t\t}")
		g.L("\t\toffsets[i] = 32 + offset // Adjust offset to be from start of data")
		g.L("")

		// Calculate element size to advance dynamicOffset
		g.L("\t\t// Calculate element size")
		g.L("\t\tvar n int")
		g.genSliceElemSizeCalc(*t.Elem, "data[32+offset:]")
		g.L("\t\tdynamicOffset += n")

		g.L("\t}")
		g.L("")
		g.L("\ttotalSize := 32 + dynamicOffset")
		g.L("\treturn &%s{", typeName)
		g.L("\t\tdata: data[:totalSize],")
		g.L("\t\tlength: length,")
		g.L("\t\toffsets: offsets,")
		g.L("\t}, totalSize, nil")
	}
	g.L("}")
}

// genSliceElemSizeCalc generates size calculation for a slice element
func (g *Generator) genSliceElemSizeCalc(t ethabi.Type, dataRef string) {
	switch t.T {
	case ethabi.StringTy, ethabi.BytesTy:
		g.L("\t\tlength, err := %sDecodeSize(%s)", g.StdPrefix, dataRef)
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")
		g.L("\t\tn = 32 + %sPad32(length)", g.StdPrefix)

	case ethabi.SliceTy:
		// Use the decode function to get size (works for both stdlib and custom types)
		g.L("\t\t_, n, err = %s", g.genDecodeCall(t, dataRef))
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")

	case ethabi.TupleTy:
		tupleName := abi.TupleStructName(t)
		g.L("\t\t_, n, err = Decode%sView(%s)", tupleName, dataRef)
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")

	case ethabi.ArrayTy:
		if IsDynamicType(*t.Elem) {
			// Dynamic array - need to parse
			g.L("\t\t// Calculate dynamic array size")
			g.L("\t\tn = 0")
			g.L("\t\tarrDynOffset := %d", t.Size*32)
			g.L("\t\tfor j := 0; j < %d; j++ {", t.Size)
			g.L("\t\t\t_, err := %sDecodeSize(%s[n:])", g.StdPrefix, dataRef)
			g.L("\t\t\tif err != nil {")
			g.L("\t\t\t\treturn nil, 0, err")
			g.L("\t\t\t}")
			g.L("\t\t\tn += 32")
			g.L("\t\t\tvar elemN int")
			// Recursive call for inner element
			g.genSliceElemSizeCalcInner(*t.Elem, fmt.Sprintf("%s[arrDynOffset:]", dataRef), "elemN")
			g.L("\t\t\tarrDynOffset += elemN")
			g.L("\t\t}")
			g.L("\t\tn = arrDynOffset")
		} else {
			g.L("\t\tn = %d", t.Size*GetTypeSize(*t.Elem))
		}
	}
}

// genSliceElemSizeCalcInner generates size calculation for nested elements (different indentation)
func (g *Generator) genSliceElemSizeCalcInner(t ethabi.Type, dataRef string, resultVar string) {
	switch t.T {
	case ethabi.StringTy, ethabi.BytesTy:
		g.L("\t\t\tlength, err := %sDecodeSize(%s)", g.StdPrefix, dataRef)
		g.L("\t\t\tif err != nil {")
		g.L("\t\t\t\treturn nil, 0, err")
		g.L("\t\t\t}")
		g.L("\t\t\t%s = 32 + %sPad32(length)", resultVar, g.StdPrefix)

	case ethabi.SliceTy:
		// Use the decode function to get size (works for both stdlib and custom types)
		g.L("\t\t\t_, %s, err = %s", resultVar, g.genDecodeCall(t, dataRef))
		g.L("\t\t\tif err != nil {")
		g.L("\t\t\t\treturn nil, 0, err")
		g.L("\t\t\t}")

	case ethabi.TupleTy:
		tupleName := abi.TupleStructName(t)
		g.L("\t\t\t_, %s, err = Decode%sView(%s)", resultVar, tupleName, dataRef)
		g.L("\t\t\tif err != nil {")
		g.L("\t\t\t\treturn nil, 0, err")
		g.L("\t\t\t}")

	default:
		g.L("\t\t\t%s = %d", resultVar, GetTypeSize(t))
	}
}

// genSliceViewGet generates the Get(i) method
func (g *Generator) genSliceViewGet(t ethabi.Type, typeName string, elemType string) {
	hasDynamicElem := IsDynamicType(*t.Elem)
	zeroVal := zeroValue(elemType)

	g.L("")
	g.L("// Get returns element at index i")
	g.L("func (v *%s) Get(i int) (%s, error) {", typeName, elemType)
	g.L("\tif i < 0 || i >= v.length {")
	g.L("\t\treturn %s, %sErrViewIndexOutOfBounds", zeroVal, g.StdPrefix)
	g.L("\t}")

	if !hasDynamicElem {
		// Static elements - calculate offset directly
		elemSize := GetTypeSize(*t.Elem)
		g.L("\toffset := 32 + i * %d", elemSize)
		g.genSliceViewGetBody(*t.Elem, "v.data[offset:]")
	} else {
		// Dynamic elements - use pre-parsed offsets
		g.L("\toffset := v.offsets[i]")
		g.genSliceViewGetBody(*t.Elem, "v.data[offset:]")
	}

	g.L("}")
}

// genSliceViewGetBody generates the body of Get method
func (g *Generator) genSliceViewGetBody(t ethabi.Type, dataRef string) {
	switch t.T {
	case ethabi.TupleTy:
		tupleName := abi.TupleStructName(t)
		g.L("\tview, _, err := Decode%sView(%s)", tupleName, dataRef)
		g.L("\treturn view, err")

	case ethabi.SliceTy:
		// If the slice type is a stdlib type (and we're not in stdlib mode),
		// decode directly instead of using the view
		typeID := abi.GenTypeIdentifier(t)
		if !g.Options.Stdlib && abi.IsStdlibType(typeID) {
			g.L("\tvalue, _, err := %s", g.genDecodeCall(t, dataRef))
			g.L("\treturn value, err")
		} else {
			viewTypeName := sliceViewTypeName(t)
			g.L("\tview, _, err := Decode%s(%s)", viewTypeName, dataRef)
			g.L("\treturn view, err")
		}

	default:
		g.L("\tvalue, _, err := %s", g.genDecodeCall(t, dataRef))
		g.L("\treturn value, err")
	}
}

// genAllSliceViews generates SliceView types for all slice types
func (g *Generator) genAllSliceViews(abiDef ethabi.ABI) {
	sliceTypes := collectSliceTypes(abiDef)
	for _, sliceType := range sliceTypes {
		// Skip slice types whose element type references external tuples
		if sliceType.Elem != nil && typeReferencesExternalTuple(*sliceType.Elem, g.Options.ExternalTuples) {
			continue
		}
		g.genSliceView(sliceType)
	}
}
