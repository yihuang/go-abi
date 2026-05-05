package generator

import (
	"fmt"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

// validateABI rejects hand-constructed types with nil Elem or TupleElems
// entries. Types from ethabi.JSON are always well-formed.
func validateABI(abiDef ethabi.ABI) error {
	for _, name := range SortedMapKeys(abiDef.Methods) {
		m := abiDef.Methods[name]
		for i, in := range m.Inputs {
			if err := validateABIType(in.Type, fmt.Sprintf("method %s input[%d]", name, i)); err != nil {
				return err
			}
		}
		for i, out := range m.Outputs {
			if err := validateABIType(out.Type, fmt.Sprintf("method %s output[%d]", name, i)); err != nil {
				return err
			}
		}
	}
	for _, name := range SortedMapKeys(abiDef.Events) {
		e := abiDef.Events[name]
		for i, in := range e.Inputs {
			if err := validateABIType(in.Type, fmt.Sprintf("event %s input[%d]", name, i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateABIType(t ethabi.Type, path string) error {
	switch t.T {
	case ethabi.SliceTy, ethabi.ArrayTy:
		if t.Elem == nil {
			return fmt.Errorf("malformed ABI type at %s: %s missing Elem", path, t.String())
		}
		return validateABIType(*t.Elem, path+".elem")
	case ethabi.TupleTy:
		for i, elem := range t.TupleElems {
			if elem == nil {
				return fmt.Errorf("malformed ABI type at %s: TupleElems[%d] is nil", path, i)
			}
			if err := validateABIType(*elem, fmt.Sprintf("%s.field[%d]", path, i)); err != nil {
				return err
			}
		}
	}
	return nil
}
