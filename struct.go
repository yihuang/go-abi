package abi

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
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
}

func StructFromInputs(method abi.Method) Struct {
	fields := make([]StructField, 0, len(method.Inputs))
	for _, input := range method.Inputs {
		fields = append(fields, StructFieldFromArgument(input))
	}
	return Struct{
		Name:   fmt.Sprintf("%sCall", Title.String(method.Name)),
		Fields: fields,
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
	}
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
