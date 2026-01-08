package generator

import (
	"fmt"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

// genViewStruct generates a View struct definition for a tuple
func (g *Generator) genViewStruct(s Struct) {
	dynamicCount := countDynamicFields(s)

	g.L("")
	g.L("// %sView provides lazy access to %s ABI data", s.Name, s.Name)
	g.L("type %sView struct {", s.Name)
	g.L("\tdata []byte")
	if dynamicCount > 0 {
		g.L("\toffsets [%d]int // offsets for dynamic fields", dynamicCount)
	}
	g.L("}")
}

// genViewDecodeFunction generates DecodeXxxView function
func (g *Generator) genViewDecodeFunction(s Struct) {
	staticSize := GetTupleSize(s.Types())
	dynamicCount := countDynamicFields(s)

	g.L("")
	g.L("// Decode%sView creates a lazy view of %s from ABI bytes", s.Name, s.Name)
	g.L("func Decode%sView(data []byte) (*%sView, int, error) {", s.Name, s.Name)

	// Validate minimum size
	g.L("\tif len(data) < %d {", staticSize)
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")

	if dynamicCount > 0 {
		g.L("\tvar (")
		g.L("\t\terr error")
		g.L("\t\toffset int")
		g.L("\t\tn int")
		g.L("\t\toffsets [%d]int", dynamicCount)
		g.L("\t)")
		g.L("\tdynamicOffset := %d", staticSize)

		// Parse and validate offsets for dynamic fields
		var dynamicIdx int
		var staticOffset int
		for _, field := range s.Fields {
			if !IsDynamicType(*field.Type) {
				// Static field - just advance offset
				staticOffset += GetTypeSize(*field.Type)
			} else {
				// Dynamic field - parse and validate offset
				g.L("")
				g.L("\t// Parse offset for dynamic field %s", field.Name)
				g.L("\toffset, err = %sDecodeSize(data[%d:])", g.StdPrefix, staticOffset)
				g.L("\tif err != nil {")
				g.L("\t\treturn nil, 0, err")
				g.L("\t}")
				g.L("\tif offset != dynamicOffset {")
				g.L("\t\treturn nil, 0, %sErrInvalidOffsetForDynamicField", g.StdPrefix)
				g.L("\t}")
				g.L("\toffsets[%d] = offset", dynamicIdx)

				// Calculate size of dynamic field to advance dynamicOffset
				g.genViewFieldSizeCalc(*field.Type, "data[offset:]", "n")
				g.L("\tdynamicOffset += n")

				dynamicIdx++
				staticOffset += 32 // Offset pointer size
			}
		}

		g.L("")
		g.L("\treturn &%sView{", s.Name)
		g.L("\t\tdata: data[:dynamicOffset],")
		g.L("\t\toffsets: offsets,")
		g.L("\t}, dynamicOffset, nil")
	} else {
		// All-static tuple
		g.L("\treturn &%sView{", s.Name)
		g.L("\t\tdata: data[:%d],", staticSize)
		g.L("\t}, %d, nil", staticSize)
	}
	g.L("}")
}

// genViewFieldSizeCalc generates inline size calculation for a dynamic field
func (g *Generator) genViewFieldSizeCalc(t ethabi.Type, dataRef string, resultVar string) {
	switch t.T {
	case ethabi.StringTy, ethabi.BytesTy:
		g.L("\t{")
		g.L("\t\tlength, err := %sDecodeSize(%s)", g.StdPrefix, dataRef)
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")
		g.L("\t\t%s = 32 + %sPad32(length)", resultVar, g.StdPrefix)
		g.L("\t}")

	case ethabi.SliceTy:
		if !IsDynamicType(*t.Elem) {
			// Static element slice
			elemSize := GetTypeSize(*t.Elem)
			g.L("\t{")
			g.L("\t\tlength, err := %sDecodeSize(%s)", g.StdPrefix, dataRef)
			g.L("\t\tif err != nil {")
			g.L("\t\t\treturn nil, 0, err")
			g.L("\t\t}")
			g.L("\t\t%s = 32 + length * %d", resultVar, elemSize)
			g.L("\t}")
		} else {
			// Dynamic element slice - use decode function to get size
			// This works for both stdlib and custom slice types
			g.L("\t{")
			g.L("\t\t_, %s, err = %s", resultVar, g.genDecodeCall(t, dataRef))
			g.L("\t\tif err != nil {")
			g.L("\t\t\treturn nil, 0, err")
			g.L("\t\t}")
			g.L("\t}")
		}

	case ethabi.ArrayTy:
		if IsDynamicType(*t.Elem) {
			// Dynamic element array - need to calculate each element's size
			g.L("\t{")
			g.L("\t\tarrayOffset := 0")
			g.L("\t\texpectedDynOffset := %d", t.Size*32)
			g.L("\t\tfor i := 0; i < %d; i++ {", t.Size)
			g.L("\t\t\telemOff, err := %sDecodeSize(%s[arrayOffset:])", g.StdPrefix, dataRef)
			g.L("\t\t\tif err != nil {")
			g.L("\t\t\t\treturn nil, 0, err")
			g.L("\t\t\t}")
			g.L("\t\t\tif elemOff != expectedDynOffset {")
			g.L("\t\t\t\treturn nil, 0, %sErrInvalidOffsetForArrayElement", g.StdPrefix)
			g.L("\t\t\t}")
			g.L("\t\t\tarrayOffset += 32")
			g.L("\t\t\tvar elemSize int")
			g.genViewFieldSizeCalcInner(*t.Elem, fmt.Sprintf("%s[expectedDynOffset:]", dataRef), "elemSize")
			g.L("\t\t\texpectedDynOffset += elemSize")
			g.L("\t\t}")
			g.L("\t\t%s = expectedDynOffset", resultVar)
			g.L("\t}")
		} else {
			// Static element array
			g.L("\t%s = %d", resultVar, t.Size*GetTypeSize(*t.Elem))
		}

	case ethabi.TupleTy:
		// Nested tuple
		tupleName := abi.TupleStructName(t)
		g.L("\t{")
		g.L("\t\t_, %s, err = Decode%sView(%s)", resultVar, tupleName, dataRef)
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")
		g.L("\t}")
	}
}

// genViewFieldSizeCalcInner generates size calculation without extra braces (for use inside loops)
func (g *Generator) genViewFieldSizeCalcInner(t ethabi.Type, dataRef string, resultVar string) {
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
	}
}

// genViewGetters generates getter methods for all fields
func (g *Generator) genViewGetters(s Struct) {
	var dynamicIdx int
	var staticOffset int

	for _, field := range s.Fields {
		g.genViewGetter(s.Name, field, staticOffset, dynamicIdx)

		if !IsDynamicType(*field.Type) {
			staticOffset += GetTypeSize(*field.Type)
		} else {
			dynamicIdx++
			staticOffset += 32
		}
	}
}

// genViewGetter generates a single getter method
func (g *Generator) genViewGetter(structName string, field StructField, staticOffset int, dynamicIdx int) {
	returnType := g.viewGetterReturnType(*field.Type)

	g.L("")
	g.L("// %s returns the %s field", field.Name, field.Type.String())
	g.L("func (v *%sView) %s() (%s, error) {", structName, field.Name, returnType)

	if !IsDynamicType(*field.Type) {
		// Static field - decode from fixed offset in data
		g.genViewGetterBody(*field.Type, fmt.Sprintf("v.data[%d:]", staticOffset))
	} else {
		// Dynamic field - decode from stored offset
		g.genViewGetterBody(*field.Type, fmt.Sprintf("v.data[v.offsets[%d]:]", dynamicIdx))
	}

	g.L("}")
}

// viewGetterReturnType returns the return type for a view getter
func (g *Generator) viewGetterReturnType(t ethabi.Type) string {
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

// genViewGetterBody generates the body of a getter method
func (g *Generator) genViewGetterBody(t ethabi.Type, dataRef string) {
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

// genViewMaterialize generates Materialize method
func (g *Generator) genViewMaterialize(s Struct) {
	g.L("")
	g.L("// Materialize fully decodes the view into %s", s.Name)
	g.L("func (v *%sView) Materialize() (*%s, error) {", s.Name, s.Name)
	g.L("\tresult := &%s{}", s.Name)
	g.L("\t_, err := result.Decode(v.data)")
	g.L("\tif err != nil {")
	g.L("\t\treturn nil, err")
	g.L("\t}")
	g.L("\treturn result, nil")
	g.L("}")
}

// genViewRaw generates Raw method
func (g *Generator) genViewRaw(s Struct) {
	g.L("")
	g.L("// Raw returns the underlying encoded bytes")
	g.L("func (v *%sView) Raw() []byte {", s.Name)
	g.L("\treturn v.data")
	g.L("}")
}

// genAllViews generates all View types for tuples
func (g *Generator) genAllViews(abiDef ethabi.ABI) {
	// Collect all tuple types
	tupleTypes := make(map[string]Struct)

	var collectTuples func(t ethabi.Type)
	collectTuples = func(t ethabi.Type) {
		if t.T == ethabi.TupleTy {
			name := abi.TupleStructName(t)
			if _, exists := tupleTypes[name]; !exists {
				if _, isExternal := g.Options.ExternalTuples[name]; !isExternal {
					tupleTypes[name] = StructFromTuple(t)
				}
			}
			// Recurse into tuple fields
			for _, elem := range t.TupleElems {
				collectTuples(*elem)
			}
		}
		// Recurse into other composite types
		switch t.T {
		case ethabi.SliceTy, ethabi.ArrayTy:
			if t.Elem != nil {
				collectTuples(*t.Elem)
			}
		}
	}

	// Collect from all methods - both the input tuple types and generate Call/Return views
	for _, method := range abiDef.Methods {
		for _, input := range method.Inputs {
			collectTuples(input.Type)
		}
		for _, output := range method.Outputs {
			collectTuples(output.Type)
		}

		// Also add the Call and Return structs themselves
		if len(method.Inputs) > 0 {
			callName := fmt.Sprintf("%sCall", Title.String(method.Name))
			callStruct := StructFromArguments(callName, method.Inputs)
			tupleTypes[callName] = callStruct
		}
		if len(method.Outputs) > 0 {
			returnName := fmt.Sprintf("%sReturn", Title.String(method.Name))
			returnStruct := StructFromArguments(returnName, method.Outputs)
			tupleTypes[returnName] = returnStruct
		}
	}

	// Collect from all events
	for _, event := range abiDef.Events {
		for _, input := range event.Inputs {
			collectTuples(input.Type)
		}
	}

	// Generate views for all tuples in sorted order
	for _, name := range SortedMapKeys(tupleTypes) {
		s := tupleTypes[name]
		// Skip structs that reference external tuples (their views would be incomplete)
		if structReferencesExternalTuple(s, g.Options.ExternalTuples) {
			continue
		}
		g.genViewStruct(s)
		g.genViewDecodeFunction(s)
		g.genViewGetters(s)
		g.genViewMaterialize(s)
		g.genViewRaw(s)
	}
}
