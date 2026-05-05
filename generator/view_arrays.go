package generator

import (
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

// arrayViewTypeName returns the view type name for a fixed-size array
// e.g., Profile[3] -> "ProfileArray3View".
func arrayViewTypeName(t ethabi.Type) string {
	return abi.GenTypeIdentifier(t) + "View"
}

// genArrayView emits an ArrayView for a fixed-size dynamic-element array.
// Static-element arrays use the materializing decoder; views can't beat
// O(N) reads with no allocations.
func (g *Generator) genArrayView(t ethabi.Type) {
	if t.T != ethabi.ArrayTy {
		panic("genArrayView called on non-array type")
	}
	if !IsDynamicType(*t.Elem) {
		return
	}

	typeName := arrayViewTypeName(t)
	elemType := g.arrayViewElemReturnType(*t.Elem)
	n := t.Size

	g.L("")
	g.L("// %s provides lazy indexed access to %s", typeName, t.String())
	g.L("type %s struct {", typeName)
	g.L("\tdata []byte")
	g.L("\toffsets [%d]int", n)
	g.L("}")

	g.genArrayViewDecodeFunction(t, typeName)

	g.L("")
	g.L("// Len returns the number of elements")
	g.L("func (v *%s) Len() int {", typeName)
	g.L("\treturn %d", n)
	g.L("}")

	g.genArrayViewGet(t, typeName, elemType)

	g.L("")
	g.L("// Raw returns the underlying encoded bytes")
	g.L("func (v *%s) Raw() []byte {", typeName)
	g.L("\treturn v.data")
	g.L("}")

	baseElemType := g.abiTypeToGoType(*t.Elem)
	g.L("")
	g.L("// Materialize fully decodes all elements into an array")
	g.L("func (v *%s) Materialize() ([%d]%s, error) {", typeName, n, baseElemType)
	g.L("\tresult, _, err := %s", g.genDecodeCall(t, "v.data"))
	g.L("\treturn result, err")
	g.L("}")
}

// arrayViewElemReturnType returns the return type for Get(i).
func (g *Generator) arrayViewElemReturnType(t ethabi.Type) string {
	switch t.T {
	case ethabi.TupleTy:
		return "*" + abi.TupleStructName(t) + "View"
	case ethabi.SliceTy:
		typeID := abi.GenTypeIdentifier(t)
		if !g.Options.Stdlib && abi.IsStdlibType(typeID) {
			return g.abiTypeToGoType(t)
		}
		return "*" + sliceViewTypeName(t)
	case ethabi.ArrayTy:
		if IsDynamicType(*t.Elem) {
			return "*" + arrayViewTypeName(t)
		}
		return g.abiTypeToGoType(t)
	default:
		return g.abiTypeToGoType(t)
	}
}

// genArrayViewDecodeFunction emits DecodeXxx for an array view. Reads the
// fixed-size offset table and validates monotonic + in-bounds (no per-element
// walks). Offsets in the EVM ABI fixed-array encoding are relative to data
// start (no length prefix), so the floor is N*32 and the bound is len(data).
func (g *Generator) genArrayViewDecodeFunction(t ethabi.Type, typeName string) {
	n := t.Size

	g.L("")
	g.L("// Decode%s creates a lazy view of %s", typeName, t.String())
	g.L("func Decode%s(data []byte) (*%s, int, error) {", typeName, typeName)

	g.L("\tif len(data) < %d {", n*32)
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")

	g.L("\tvar offsets [%d]int", n)
	g.L("\tprev := %d - 1 // floor: first offset must be >= N*32", n*32)
	g.L("\tfor i := 0; i < %d; i++ {", n)
	g.L("\t\toff, err := %sDecodeSize(data[i*32:])", g.StdPrefix)
	g.L("\t\tif err != nil {")
	g.L("\t\t\treturn nil, 0, err")
	g.L("\t\t}")
	g.L("\t\tif off <= prev || off > len(data) {")
	g.L("\t\t\treturn nil, 0, %sErrInvalidOffsetForArrayElement", g.StdPrefix)
	g.L("\t\t}")
	g.L("\t\toffsets[i] = off")
	g.L("\t\tprev = off")
	g.L("\t}")

	g.L("")
	g.L("\treturn &%s{", typeName)
	g.L("\t\tdata: data,")
	g.L("\t\toffsets: offsets,")
	g.L("\t}, len(data), nil")
	g.L("}")
}

// genArrayViewGet emits Get(i). Each child gets a sub-slice bounded by the
// next element's offset (or len(data) for the last).
func (g *Generator) genArrayViewGet(t ethabi.Type, typeName string, elemType string) {
	n := t.Size
	zeroVal := zeroValue(elemType)

	g.L("")
	g.L("// Get returns element at index i")
	g.L("func (v *%s) Get(i int) (%s, error) {", typeName, elemType)
	g.L("\tif i < 0 || i >= %d {", n)
	g.L("\t\treturn %s, %sErrViewIndexOutOfBounds", zeroVal, g.StdPrefix)
	g.L("\t}")
	g.L("\tstart := v.offsets[i]")
	g.L("\tvar end int")
	g.L("\tif i+1 < %d {", n)
	g.L("\t\tend = v.offsets[i+1]")
	g.L("\t} else {")
	g.L("\t\tend = len(v.data)")
	g.L("\t}")
	g.genArrayViewGetBody(*t.Elem, "v.data[start:end]")
	g.L("}")
}

// genArrayViewGetBody emits the body of Get for one element.
func (g *Generator) genArrayViewGetBody(t ethabi.Type, dataRef string) {
	switch t.T {
	case ethabi.TupleTy:
		tupleName := abi.TupleStructName(t)
		g.L("\tview, _, err := Decode%sView(%s)", tupleName, dataRef)
		g.L("\treturn view, err")
	case ethabi.SliceTy:
		typeID := abi.GenTypeIdentifier(t)
		if !g.Options.Stdlib && abi.IsStdlibType(typeID) {
			g.L("\tvalue, _, err := %s", g.genDecodeCall(t, dataRef))
			g.L("\treturn value, err")
		} else {
			g.L("\tview, _, err := Decode%s(%s)", sliceViewTypeName(t), dataRef)
			g.L("\treturn view, err")
		}
	case ethabi.ArrayTy:
		if IsDynamicType(*t.Elem) {
			g.L("\tview, _, err := Decode%s(%s)", arrayViewTypeName(t), dataRef)
			g.L("\treturn view, err")
		} else {
			g.L("\tvalue, _, err := %s", g.genDecodeCall(t, dataRef))
			g.L("\treturn value, err")
		}
	default:
		g.L("\tvalue, _, err := %s", g.genDecodeCall(t, dataRef))
		g.L("\treturn value, err")
	}
}

// collectArrayViewTypes collects fixed-size arrays with dynamic elements.
func collectArrayViewTypes(abiDef ethabi.ABI) []ethabi.Type {
	typeSet := make(map[string]ethabi.Type)

	var visit func(t ethabi.Type)
	visit = func(t ethabi.Type) {
		if t.T == ethabi.ArrayTy && IsDynamicType(*t.Elem) {
			id := abi.GenTypeIdentifier(t)
			if _, ok := typeSet[id]; !ok {
				typeSet[id] = t
			}
		}
		switch t.T {
		case ethabi.SliceTy, ethabi.ArrayTy:
			visit(*t.Elem)
		case ethabi.TupleTy:
			for _, elem := range t.TupleElems {
				visit(*elem)
			}
		}
	}

	for _, method := range abiDef.Methods {
		for _, in := range method.Inputs {
			visit(in.Type)
		}
		for _, out := range method.Outputs {
			visit(out.Type)
		}
	}
	for _, event := range abiDef.Events {
		for _, in := range event.Inputs {
			visit(in.Type)
		}
	}

	result := make([]ethabi.Type, 0, len(typeSet))
	for _, name := range SortedMapKeys(typeSet) {
		result = append(result, typeSet[name])
	}
	return result
}

// genAllArrayViews emits ArrayView types for all dynamic-element fixed arrays.
func (g *Generator) genAllArrayViews(abiDef ethabi.ABI) {
	for _, t := range collectArrayViewTypes(abiDef) {
		if typeReferencesExternalTuple(*t.Elem, g.Options.ExternalTuples) {
			continue
		}
		g.genArrayView(t)
	}
}
