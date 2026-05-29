package generator

import (
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

// countDynamicFields returns the number of dynamic fields in a struct
func countDynamicFields(s Struct) int {
	count := 0
	for _, field := range s.Fields {
		if IsDynamicType(*field.Type) {
			count++
		}
	}
	return count
}

// collectSliceTypes collects all unique slice types from an ABI for view generation
func collectSliceTypes(abiDef ethabi.ABI) []ethabi.Type {
	return collectTypesIf(abiDef, func(t ethabi.Type) bool { return t.T == ethabi.SliceTy })
}

// collectTypesIf returns the unique types in abiDef matching pred, keyed by
// GenTypeIdentifier and emitted in deterministic (sorted) order.
func collectTypesIf(abiDef ethabi.ABI, pred func(ethabi.Type) bool) []ethabi.Type {
	seen := make(map[string]ethabi.Type)
	walkAbiTypes(abiDef, func(t ethabi.Type) {
		if !pred(t) {
			return
		}
		id := abi.GenTypeIdentifier(t)
		if _, ok := seen[id]; !ok {
			seen[id] = t
		}
	})
	out := make([]ethabi.Type, 0, len(seen))
	for _, k := range SortedMapKeys(seen) {
		out = append(out, seen[k])
	}
	return out
}

// walkAbiTypes calls visit on every type appearing in the ABI's method
// inputs/outputs and event inputs, recursing into Slice/Array/Tuple children.
func walkAbiTypes(abiDef ethabi.ABI, visit func(ethabi.Type)) {
	var walk func(t ethabi.Type)
	walk = func(t ethabi.Type) {
		visit(t)
		switch t.T {
		case ethabi.SliceTy, ethabi.ArrayTy:
			walk(*t.Elem)
		case ethabi.TupleTy:
			for _, elem := range t.TupleElems {
				walk(*elem)
			}
		}
	}
	for _, m := range abiDef.Methods {
		for _, in := range m.Inputs {
			walk(in.Type)
		}
		for _, out := range m.Outputs {
			walk(out.Type)
		}
	}
	for _, e := range abiDef.Events {
		for _, in := range e.Inputs {
			walk(in.Type)
		}
	}
}

// viewElemReturnType is the Go return type for a value extracted from a view
// (a tuple field getter, or Get(i) on a slice/array view). Tuples become
// XxxView, slices become XxxSliceView (unless they're stdlib types in non-
// stdlib builds, which return their plain Go type to avoid duplicate decls),
// dynamic-elem fixed arrays become XxxArrayNView, and everything else returns
// its plain Go type.
func (g *Generator) viewElemReturnType(t ethabi.Type) string {
	switch t.T {
	case ethabi.TupleTy:
		return abi.TupleStructName(t) + "View"
	case ethabi.SliceTy:
		typeID := abi.GenTypeIdentifier(t)
		if !g.Options.Stdlib && abi.IsStdlibType(typeID) {
			return g.abiTypeToGoType(t)
		}
		return sliceViewTypeName(t)
	case ethabi.ArrayTy:
		if IsDynamicType(*t.Elem) {
			return arrayViewTypeName(t)
		}
		return g.abiTypeToGoType(t)
	default:
		return g.abiTypeToGoType(t)
	}
}

// genViewGetBody emits the body of a tuple-view getter or a slice/array-view
// Get method for one element at dataRef. Mirrors viewElemReturnType: composite
// types whose return is a view get DecodeXxxView; everything else falls back
// to the materializing decoder.
func (g *Generator) genViewGetBody(t ethabi.Type, dataRef string) {
	switch t.T {
	case ethabi.TupleTy:
		g.L("\tview, _, err := Decode%sView(%s)", abi.TupleStructName(t), dataRef)
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

// zeroValue returns the zero value literal for a Go type
func zeroValue(goType string) string {
	if strings.HasPrefix(goType, "*") || strings.HasPrefix(goType, "[]") {
		return "nil"
	}
	if strings.HasPrefix(goType, "[") && strings.Contains(goType, "]") {
		return goType + "{}" // Array type
	}
	switch goType {
	case "bool":
		return "false"
	case "string":
		return `""`
	default:
		if strings.HasPrefix(goType, "uint") || strings.HasPrefix(goType, "int") {
			return "0"
		}
		return goType + "{}" // Struct or other
	}
}

// typeReferencesExternalTuple checks if a type (or any nested type) references an external tuple
func typeReferencesExternalTuple(t ethabi.Type, externalTuples map[string]string) bool {
	switch t.T {
	case ethabi.TupleTy:
		name := abi.TupleStructName(t)
		if _, isExternal := externalTuples[name]; isExternal {
			return true
		}
		// Check nested fields
		for _, elem := range t.TupleElems {
			if typeReferencesExternalTuple(*elem, externalTuples) {
				return true
			}
		}
	case ethabi.SliceTy, ethabi.ArrayTy:
		return typeReferencesExternalTuple(*t.Elem, externalTuples)
	}
	return false
}

// structReferencesExternalTuple checks if a Struct references any external tuple
func structReferencesExternalTuple(s Struct, externalTuples map[string]string) bool {
	for _, field := range s.Fields {
		if typeReferencesExternalTuple(*field.Type, externalTuples) {
			return true
		}
	}
	return false
}
