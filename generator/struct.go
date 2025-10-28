package generator

import (
	"fmt"

	abi "github.com/ethereum/go-ethereum/accounts/abi"
)

type StructField struct {
	Type *abi.Type
	Name string
}

func StructFieldFromArgument(arg abi.Argument) StructField {
	return StructField{
		Type: &arg.Type,
		Name: Title.String(arg.Name),
	}
}

func StructFieldFromTupleElement(t abi.Type, index int) StructField {
	fieldName := t.TupleRawNames[index]
	if fieldName == "" {
		fieldName = fmt.Sprintf("Field%d", index+1)
	}
	return StructField{
		Type: t.TupleElems[index],
		Name: Title.String(fieldName),
	}
}

type Struct struct {
	Name   string
	Fields []StructField

	// The tuple type
	T abi.Type
}

func StructFromArguments(name string, args []abi.Argument) Struct {
	fields := make([]StructField, 0, len(args))
	types := make([]*abi.Type, 0, len(args))
	names := make([]string, 0, len(args))
	for i, input := range args {
		field := StructFieldFromArgument(input)
		if field.Name == "" {
			field.Name = fmt.Sprintf("Field%d", i+1)
		}
		fields = append(fields, field)
		types = append(types, field.Type)
		names = append(names, field.Name)
	}
	return Struct{
		Name:   name,
		Fields: fields,
		T:      abi.Type{T: abi.TupleTy, TupleElems: types, TupleRawNames: names, TupleRawName: name},
	}
}

func StructFromTuple(t abi.Type) Struct {
	fields := make([]StructField, 0, len(t.TupleElems))
	for i := range t.TupleElems {
		fields = append(fields, StructFieldFromTupleElement(t, i))
	}
	return Struct{
		Name:   TupleStructName(t),
		Fields: fields,
		T:      t,
	}
}

func StructFromEventData(event abi.Event) Struct {
	name := fmt.Sprintf("%sEventData", Title.String(event.Name))
	arguments := make([]abi.Argument, 0)
	for _, input := range event.Inputs {
		if input.Indexed {
			continue
		}
		arguments = append(arguments, input)
	}
	return StructFromArguments(name, arguments)
}

func (s Struct) Types() []*abi.Type {
	types := make([]*abi.Type, len(s.Fields))
	for i, field := range s.Fields {
		types[i] = field.Type
	}
	return types
}

func (s Struct) HasDynamicField() bool {
	for _, field := range s.Fields {
		if IsDynamicType(*field.Type) {
			return true
		}
	}
	return false
}
