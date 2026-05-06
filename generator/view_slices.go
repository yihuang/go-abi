package generator

import (
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
	elemType := g.viewElemReturnType(*t.Elem)
	hasDynamicElem := IsDynamicType(*t.Elem)

	g.L("")
	g.L("// %s provides lazy indexed access to %s", typeName, t.String())
	g.L("type %s struct {", typeName)
	g.L("\tdata []byte")
	g.L("\tlength int")
	if hasDynamicElem {
		g.L("\toffsets []int // absolute offset into data for each element")
	}
	g.L("}")

	// Generate DecodeXxxSliceView
	g.genSliceViewDecodeFunction(t, typeName)

	// Generate Len() method
	g.L("")
	g.L("// Len returns the number of elements")
	g.L("func (v %s) Len() int {", typeName)
	g.L("\treturn v.length")
	g.L("}")

	// Generate Get(i) method
	g.genSliceViewGet(t, typeName, elemType)

	// Generate Raw() method
	g.L("")
	g.L("// Raw returns the underlying encoded bytes")
	g.L("func (v %s) Raw() []byte {", typeName)
	g.L("\treturn v.data")
	g.L("}")

	// Generate Materialize() method
	baseElemType := g.abiTypeToGoType(*t.Elem)
	g.L("")
	g.L("// Materialize fully decodes all elements into a slice")
	g.L("func (v %s) Materialize() ([]%s, error) {", typeName, baseElemType)
	g.L("\tresult, _, err := %s", g.genDecodeCall(t, "v.data"))
	g.L("\treturn result, err")
	g.L("}")
}

// genSliceViewDecodeFunction emits DecodeXxx for a slice view. Static-elem
// slices know their size up front; dynamic-elem slices read the offset table
// and validate monotonic + in-bounds, with no per-element walks.
func (g *Generator) genSliceViewDecodeFunction(t ethabi.Type, typeName string) {
	hasDynamicElem := IsDynamicType(*t.Elem)

	g.L("")
	g.L("// Decode%s creates a lazy view of %s", typeName, t.String())
	g.L("func Decode%s(data []byte) (%s, int, error) {", typeName, typeName)

	g.L("\tif len(data) < 32 {")
	g.L("\t\treturn %s{}, 0, io.ErrUnexpectedEOF", typeName)
	g.L("\t}")
	g.L("\tlength, err := %sDecodeSize(data)", g.StdPrefix)
	g.L("\tif err != nil {")
	g.L("\t\treturn %s{}, 0, err", typeName)
	g.L("\t}")

	if !hasDynamicElem {
		elemSize := GetTypeSize(*t.Elem)
		// Reject before arithmetic; length*N can overflow int.
		g.L("\tif length > len(data) || length*%d > len(data)-32 {", elemSize)
		g.L("\t\treturn %s{}, 0, io.ErrUnexpectedEOF", typeName)
		g.L("\t}")
		g.L("\ttotalSize := 32 + length * %d", elemSize)
		g.L("\treturn %s{", typeName)
		g.L("\t\tdata: data[:totalSize],")
		g.L("\t\tlength: length,")
		g.L("\t}, totalSize, nil")
	} else {
		g.L("\tif length == 0 {")
		g.L("\t\treturn %s{data: data[:32], length: 0, offsets: nil}, 32, nil", typeName)
		g.L("\t}")
		g.L("")
		g.L("\tif length > len(data) || length*32 > len(data)-32 {")
		g.L("\t\treturn %s{}, 0, io.ErrUnexpectedEOF", typeName)
		g.L("\t}")
		g.L("")
		g.L("\toffsets := make([]int, length)")
		// maxOff = len(data)-32 (not 32+off vs len(data); the latter overflows
		// when off is near MaxInt and would bypass validation).
		g.L("\tprev := length*32 - 1")
		g.L("\tmaxOff := len(data) - 32")
		g.L("\tfor i := 0; i < length; i++ {")
		g.L("\t\toff, err := %sDecodeSize(data[32+i*32:])", g.StdPrefix)
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn %s{}, 0, err", typeName)
		g.L("\t\t}")
		g.L("\t\tif off <= prev || off > maxOff {")
		g.L("\t\t\treturn %s{}, 0, %sErrInvalidOffsetForSliceElement", typeName, g.StdPrefix)
		g.L("\t\t}")
		g.L("\t\toffsets[i] = 32 + off // absolute offset from start of data")
		g.L("\t\tprev = off")
		g.L("\t}")
		g.L("")
		g.L("\treturn %s{", typeName)
		g.L("\t\tdata: data,")
		g.L("\t\tlength: length,")
		g.L("\t\toffsets: offsets,")
		g.L("\t}, len(data), nil")
	}
	g.L("}")
}

// genSliceViewGet emits Get(i). Dynamic-elem slices hand the child decoder
// a sub-slice bounded by the next element's offset (or len(data) for last).
func (g *Generator) genSliceViewGet(t ethabi.Type, typeName string, elemType string) {
	hasDynamicElem := IsDynamicType(*t.Elem)
	zeroVal := zeroValue(elemType)

	g.L("")
	g.L("// Get returns element at index i")
	g.L("func (v %s) Get(i int) (%s, error) {", typeName, elemType)
	g.L("\tif i < 0 || i >= v.length {")
	g.L("\t\treturn %s, %sErrViewIndexOutOfBounds", zeroVal, g.StdPrefix)
	g.L("\t}")

	if !hasDynamicElem {
		elemSize := GetTypeSize(*t.Elem)
		g.L("\toffset := 32 + i * %d", elemSize)
		g.genViewGetBody(*t.Elem, "v.data[offset:]")
	} else {
		g.L("\tstart := v.offsets[i]")
		g.L("\tvar end int")
		g.L("\tif i+1 < v.length {")
		g.L("\t\tend = v.offsets[i+1]")
		g.L("\t} else {")
		g.L("\t\tend = len(v.data)")
		g.L("\t}")
		g.genViewGetBody(*t.Elem, "v.data[start:end]")
	}

	g.L("}")
}

// genAllSliceViews generates SliceView types for all slice types
func (g *Generator) genAllSliceViews(abiDef ethabi.ABI) {
	sliceTypes := collectSliceTypes(abiDef)
	for _, sliceType := range sliceTypes {
		// Skip slice types whose element type references external tuples
		if typeReferencesExternalTuple(*sliceType.Elem, g.Options.ExternalTuples) {
			continue
		}
		g.genSliceView(sliceType)
	}
}
