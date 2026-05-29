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

// genViewDecodeFunction emits DecodeXxxView. Construction reads only the
// dynamic offsets from the header and validates monotonic + in-bounds; no
// recursion into children. Caller must pass a slice covering exactly the
// tuple's payload. Views are returned by value so leaf-access chains
// don't allocate per level.
func (g *Generator) genViewDecodeFunction(s Struct) {
	staticSize := GetTupleSize(s.Types())
	dynamicCount := countDynamicFields(s)

	g.L("")
	g.L("// Decode%sView creates a lazy view of %s.", s.Name, s.Name)
	g.L("func Decode%sView(data []byte) (%sView, int, error) {", s.Name, s.Name)

	g.L("\tif len(data) < %d {", staticSize)
	g.L("\t\treturn %sView{}, 0, io.ErrUnexpectedEOF", s.Name)
	g.L("\t}")

	if dynamicCount > 0 {
		g.L("\tvar offsets [%d]int", dynamicCount)
		g.L("\tprevOffset := %d - 1 // floor for monotonic + in-bounds check", staticSize)

		var dynamicIdx int
		var staticOffset int
		for _, field := range s.Fields {
			if !IsDynamicType(*field.Type) {
				staticOffset += GetTypeSize(*field.Type)
				continue
			}

			g.L("")
			g.L("\t{")
			g.L("\t\toff, err := %sDecodeSize(data[%d:])", g.StdPrefix, staticOffset)
			g.L("\t\tif err != nil {")
			g.L("\t\t\treturn %sView{}, 0, err", s.Name)
			g.L("\t\t}")
			g.L("\t\tif off <= prevOffset || off > len(data) {")
			g.L("\t\t\treturn %sView{}, 0, %sErrInvalidOffsetForDynamicField", s.Name, g.StdPrefix)
			g.L("\t\t}")
			g.L("\t\toffsets[%d] = off", dynamicIdx)
			g.L("\t\tprevOffset = off")
			g.L("\t}")

			dynamicIdx++
			staticOffset += 32
		}

		g.L("")
		g.L("\treturn %sView{", s.Name)
		g.L("\t\tdata: data,")
		g.L("\t\toffsets: offsets,")
		g.L("\t}, len(data), nil")
	} else {
		g.L("\treturn %sView{", s.Name)
		g.L("\t\tdata: data[:%d],", staticSize)
		g.L("\t}, %d, nil", staticSize)
	}
	g.L("}")
}

// genViewGetters generates getter methods for all fields
func (g *Generator) genViewGetters(s Struct) {
	totalDynamic := countDynamicFields(s)

	var dynamicIdx int
	var staticOffset int
	for _, field := range s.Fields {
		isLastDynamic := IsDynamicType(*field.Type) && dynamicIdx == totalDynamic-1
		g.genViewGetter(s.Name, field, staticOffset, dynamicIdx, isLastDynamic)

		if !IsDynamicType(*field.Type) {
			staticOffset += GetTypeSize(*field.Type)
		} else {
			dynamicIdx++
			staticOffset += 32
		}
	}
}

// genViewGetter emits a single getter. Dynamic fields get a sub-slice
// bounded by the next dynamic offset (or len(data) for the last), so the
// child decoder doesn't need to re-walk to find its extent.
func (g *Generator) genViewGetter(structName string, field StructField, staticOffset int, dynamicIdx int, isLastDynamic bool) {
	returnType := g.viewElemReturnType(*field.Type)

	g.L("")
	g.L("// %s returns the %s field", field.Name, field.Type.String())
	g.L("func (v %sView) %s() (%s, error) {", structName, field.Name, returnType)

	if !IsDynamicType(*field.Type) {
		g.genViewGetBody(*field.Type, fmt.Sprintf("v.data[%d:]", staticOffset))
	} else {
		var dataRef string
		if isLastDynamic {
			dataRef = fmt.Sprintf("v.data[v.offsets[%d]:]", dynamicIdx)
		} else {
			dataRef = fmt.Sprintf("v.data[v.offsets[%d]:v.offsets[%d]]", dynamicIdx, dynamicIdx+1)
		}
		g.genViewGetBody(*field.Type, dataRef)
	}

	g.L("}")
}

// genViewMaterialize emits MaterializeTo (per-field decode that reuses the
// view's pre-validated offsets) and Materialize (allocator wrapper). Skipping
// the redundant DecodeSize + canonical-offset checks per dynamic field saves
// ~5 ns per dynamic field versus calling the materializing tuple decoder.
func (g *Generator) genViewMaterialize(s Struct) {
	totalDynamic := countDynamicFields(s)

	g.L("")
	g.L("// MaterializeTo decodes the view into a caller-owned dst.")
	g.L("// Skips offset re-parsing/validation by reusing the view's offsets.")
	g.L("// Returns ErrNilDestination if dst is nil.")
	g.L("func (v %sView) MaterializeTo(dst *%s) error {", s.Name, s.Name)
	g.L("\tif dst == nil {")
	g.L("\t\treturn %sErrNilDestination", g.StdPrefix)
	g.L("\t}")
	g.L("\tvar err error")

	var dynamicIdx int
	var staticOffset int
	for _, field := range s.Fields {
		var dataRef string
		if !IsDynamicType(*field.Type) {
			dataRef = fmt.Sprintf("v.data[%d:]", staticOffset)
			staticOffset += GetTypeSize(*field.Type)
		} else {
			if dynamicIdx == totalDynamic-1 {
				dataRef = fmt.Sprintf("v.data[v.offsets[%d]:]", dynamicIdx)
			} else {
				dataRef = fmt.Sprintf("v.data[v.offsets[%d]:v.offsets[%d]]", dynamicIdx, dynamicIdx+1)
			}
			dynamicIdx++
			staticOffset += 32
		}
		if field.Type.T == ethabi.TupleTy {
			g.L("\tif _, err = dst.%s.Decode(%s); err != nil {", field.Name, dataRef)
			g.L("\t\treturn err")
			g.L("\t}")
		} else {
			g.L("\tif dst.%s, _, err = %s; err != nil {", field.Name, g.genDecodeCall(*field.Type, dataRef))
			g.L("\t\treturn err")
			g.L("\t}")
		}
	}
	g.L("\treturn nil")
	g.L("}")

	g.L("")
	g.L("// Materialize allocates and decodes the view into a new %s. Use", s.Name)
	g.L("// MaterializeTo with a caller-owned dst to skip the destination alloc;")
	g.L("// inner dynamic fields (bytes, strings, slices) still allocate.")
	g.L("func (v %sView) Materialize() (*%s, error) {", s.Name, s.Name)
	g.L("\tresult := &%s{}", s.Name)
	g.L("\tif err := v.MaterializeTo(result); err != nil {")
	g.L("\t\treturn nil, err")
	g.L("\t}")
	g.L("\treturn result, nil")
	g.L("}")
}

// genViewRaw generates Raw method
func (g *Generator) genViewRaw(s Struct) {
	g.L("")
	g.L("// Raw returns the underlying encoded bytes")
	g.L("func (v %sView) Raw() []byte {", s.Name)
	g.L("\treturn v.data")
	g.L("}")
}

// genAllViews generates all View types for tuples
func (g *Generator) genAllViews(abiDef ethabi.ABI) {
	tupleTypes := make(map[string]Struct)

	walkAbiTypes(abiDef, func(t ethabi.Type) {
		if t.T != ethabi.TupleTy {
			return
		}
		name := abi.TupleStructName(t)
		if _, exists := tupleTypes[name]; exists {
			return
		}
		if _, isExternal := g.Options.ExternalTuples[name]; isExternal {
			return
		}
		tupleTypes[name] = StructFromTuple(t)
	})

	// Add the Call/Return structs for each method.
	for _, method := range abiDef.Methods {
		if len(method.Inputs) > 0 {
			callName := fmt.Sprintf("%sCall", Title.String(method.Name))
			tupleTypes[callName] = StructFromArguments(callName, method.Inputs)
		}
		if len(method.Outputs) > 0 {
			returnName := fmt.Sprintf("%sReturn", Title.String(method.Name))
			tupleTypes[returnName] = StructFromArguments(returnName, method.Outputs)
		}
	}

	// Generate views for all tuples in sorted order
	for _, name := range SortedMapKeys(tupleTypes) {
		s := tupleTypes[name]
		// No View for structs referencing external tuples — we can't
		// generate getters for types defined outside this package.
		if structReferencesExternalTuple(s, g.Options.ExternalTuples) {
			g.L("")
			g.L("// %sView is not generated: %s references an ExternalTuples type.", s.Name, s.Name)
			continue
		}
		g.genViewStruct(s)
		g.genViewDecodeFunction(s)
		g.genViewGetters(s)
		g.genViewMaterialize(s)
		g.genViewRaw(s)
	}
}
