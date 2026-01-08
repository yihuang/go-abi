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
	typeSet := make(map[string]ethabi.Type)

	var collectSlices func(t ethabi.Type)
	collectSlices = func(t ethabi.Type) {
		if t.T == ethabi.SliceTy {
			typeID := abi.GenTypeIdentifier(t)
			if _, exists := typeSet[typeID]; !exists {
				typeSet[typeID] = t
			}
		}

		// Recurse into nested types
		switch t.T {
		case ethabi.SliceTy, ethabi.ArrayTy:
			if t.Elem != nil {
				collectSlices(*t.Elem)
			}
		case ethabi.TupleTy:
			for _, elem := range t.TupleElems {
				collectSlices(*elem)
			}
		}
	}

	// Collect from all method inputs and outputs
	for _, method := range abiDef.Methods {
		for _, input := range method.Inputs {
			collectSlices(input.Type)
		}
		for _, output := range method.Outputs {
			collectSlices(output.Type)
		}
	}

	// Collect from all events
	for _, event := range abiDef.Events {
		for _, input := range event.Inputs {
			collectSlices(input.Type)
		}
	}

	// Convert to sorted slice for deterministic output
	result := make([]ethabi.Type, 0, len(typeSet))
	for _, name := range SortedMapKeys(typeSet) {
		result = append(result, typeSet[name])
	}
	return result
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
		if t.Elem != nil {
			return typeReferencesExternalTuple(*t.Elem, externalTuples)
		}
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
