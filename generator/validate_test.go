package generator

import (
	"strings"
	"testing"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

func TestValidateABI_RejectsMalformed(t *testing.T) {
	nilSlice := ethabi.Type{T: ethabi.SliceTy}                // Elem nil
	nilArray := ethabi.Type{T: ethabi.ArrayTy, Size: 3}       // Elem nil
	sliceOfNilSlice := ethabi.Type{T: ethabi.SliceTy, Elem: &nilSlice}
	tupleWithNilField := ethabi.Type{T: ethabi.TupleTy, TupleElems: []*ethabi.Type{nil}}
	tupleWithNilSliceField := ethabi.Type{T: ethabi.TupleTy, TupleElems: []*ethabi.Type{&nilSlice}}

	method := func(name string, in, out []ethabi.Type) ethabi.ABI {
		args := func(ts []ethabi.Type) ethabi.Arguments {
			out := make(ethabi.Arguments, len(ts))
			for i, t := range ts {
				out[i] = ethabi.Argument{Type: t}
			}
			return out
		}
		return ethabi.ABI{Methods: map[string]ethabi.Method{
			name: {Name: name, Inputs: args(in), Outputs: args(out)},
		}}
	}
	event := func(name string, in []ethabi.Type) ethabi.ABI {
		args := make(ethabi.Arguments, len(in))
		for i, t := range in {
			args[i] = ethabi.Argument{Type: t}
		}
		return ethabi.ABI{Events: map[string]ethabi.Event{
			name: {Name: name, Inputs: args},
		}}
	}

	cases := []struct {
		name     string
		abiDef   ethabi.ABI
		wantPath string // substring expected in the error path
	}{
		{"slice nil Elem in method input", method("foo", []ethabi.Type{nilSlice}, nil), "method foo input[0]"},
		{"array nil Elem in method output", method("foo", nil, []ethabi.Type{nilArray}), "method foo output[0]"},
		{"nil entry in tuple elems", method("foo", []ethabi.Type{tupleWithNilField}, nil), "TupleElems[0] is nil"},
		{"nested slice-of-slice with nil inner Elem", method("foo", []ethabi.Type{sliceOfNilSlice}, nil), "method foo input[0].elem"},
		{"tuple field is slice with nil Elem", method("foo", []ethabi.Type{tupleWithNilSliceField}, nil), "method foo input[0].field[0]"},
		{"malformed type in event input", event("Bar", []ethabi.Type{nilSlice}), "event Bar input[0]"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateABI(tc.abiDef)
			if err == nil {
				t.Fatalf("want error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantPath) {
				t.Fatalf("want error path containing %q, got %q", tc.wantPath, err)
			}
		})
	}
}

func TestValidateABI_AcceptsWellFormed(t *testing.T) {
	abiDef, err := ethabi.JSON(strings.NewReader(`[
		{"name":"foo","type":"function","inputs":[{"name":"v","type":"uint256[]"}],"outputs":[]},
		{"name":"bar","type":"function","inputs":[{"name":"a","type":"address[3][]"}],"outputs":[{"name":"b","type":"bool"}]},
		{
			"name":"baz","type":"function",
			"inputs":[{
				"name":"u","type":"tuple",
				"components":[
					{"name":"id","type":"uint64"},
					{"name":"tags","type":"string[]"}
				]
			}],
			"outputs":[]
		},
		{"name":"Evt","type":"event","inputs":[{"name":"x","type":"bytes"}]}
	]`))
	if err != nil {
		t.Fatalf("parse ABI: %v", err)
	}
	if err := validateABI(abiDef); err != nil {
		t.Fatalf("validateABI: %v", err)
	}
}
