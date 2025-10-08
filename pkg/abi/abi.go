package abi

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// ABI provides the main interface for encoding and decoding ABI data
type ABI struct {
	// Function selectors and their signatures
	Functions map[string]Function
	// Event signatures
	Events map[string]Event
}

// Function represents an ABI function
type Function struct {
	Name   string
	Inputs []Argument
	Outputs []Argument
}

// Event represents an ABI event
type Event struct {
	Name      string
	Inputs    []Argument
	Anonymous bool
}

// Argument represents a function or event argument
type Argument struct {
	Name string
	Type Type
}

// NewABI creates a new ABI instance
func NewABI() *ABI {
	return &ABI{
		Functions: make(map[string]Function),
		Events:    make(map[string]Event),
	}
}

// AddFunction adds a function to the ABI
func (a *ABI) AddFunction(name string, inputs, outputs []Argument) {
	selector := a.functionSelector(name, inputs)
	a.Functions[selector] = Function{
		Name:   name,
		Inputs: inputs,
		Outputs: outputs,
	}
}

// AddEvent adds an event to the ABI
func (a *ABI) AddEvent(name string, inputs []Argument, anonymous bool) {
	signature := a.eventSignature(name, inputs)
	a.Events[signature] = Event{
		Name:      name,
		Inputs:    inputs,
		Anonymous: anonymous,
	}
}

// EncodeFunctionCall encodes a function call with the given arguments
func (a *ABI) EncodeFunctionCall(name string, args ...interface{}) ([]byte, error) {
	// Find function by name
	var function Function
	var selector string

	for sel, fn := range a.Functions {
		if fn.Name == name {
			function = fn
			selector = sel
			break
		}
	}

	if function.Name == "" {
		return nil, fmt.Errorf("function %s not found", name)
	}

	// Validate argument count
	if len(args) != len(function.Inputs) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(function.Inputs), len(args))
	}

	// Start with function selector
	result, err := hex.DecodeString(selector[2:]) // Remove "0x" prefix
	if err != nil {
		return nil, fmt.Errorf("invalid selector: %w", err)
	}

	// Encode arguments
	for i, arg := range function.Inputs {
		encoded, err := arg.Type.Encode(args[i])
		if err != nil {
			return nil, fmt.Errorf("encoding argument %d: %w", i, err)
		}
		result = append(result, encoded...)
	}

	return result, nil
}

// DecodeFunctionResult decodes function return data
func (a *ABI) DecodeFunctionResult(name string, data []byte, outputs ...interface{}) error {
	// Find function by name
	var function Function

	for _, fn := range a.Functions {
		if fn.Name == name {
			function = fn
			break
		}
	}

	if function.Name == "" {
		return fmt.Errorf("function %s not found", name)
	}

	// Validate output count
	if len(outputs) != len(function.Outputs) {
		return fmt.Errorf("expected %d outputs, got %d", len(function.Outputs), len(outputs))
	}

	offset := 0
	for i, output := range function.Outputs {
		bytesRead, err := output.Type.Decode(data[offset:], outputs[i])
		if err != nil {
			return fmt.Errorf("decoding output %d: %w", i, err)
		}
		offset += int(bytesRead)
	}

	return nil
}

// functionSelector computes the 4-byte function selector
func (a *ABI) functionSelector(name string, inputs []Argument) string {
	signature := name + "("
	for i, input := range inputs {
		if i > 0 {
			signature += ","
		}
		signature += input.Type.TypeName()
	}
	signature += ")"

	// In a real implementation, we would compute the Keccak-256 hash
	// and take the first 4 bytes. For now, we'll use a placeholder.
	return "0x" + hex.EncodeToString([]byte(signature[:8]))
}

// eventSignature computes the event signature
func (a *ABI) eventSignature(name string, inputs []Argument) string {
	signature := name + "("
	for i, input := range inputs {
		if i > 0 {
			signature += ","
		}
		signature += input.Type.TypeName()
	}
	signature += ")"

	// In a real implementation, we would compute the Keccak-256 hash
	return signature
}

// Example usage functions

// ExampleUint256 creates an example ABI with uint256 operations
func ExampleUint256() *ABI {
	abi := NewABI()

	// Add a simple transfer function
	abi.AddFunction("transfer", []Argument{
		{Name: "to", Type: &Address{}},
		{Name: "value", Type: &Uint256{}},
	}, []Argument{
		{Name: "success", Type: &Bool{}},
	})

	// Add a Transfer event
	abi.AddEvent("Transfer", []Argument{
		{Name: "from", Type: &Address{}},
		{Name: "to", Type: &Address{}},
		{Name: "value", Type: &Uint256{}},
	}, false)

	return abi
}

// EncodeSimpleTransfer encodes a simple transfer call
func EncodeSimpleTransfer(to [20]byte, value *big.Int) ([]byte, error) {
	abi := ExampleUint256()
	return abi.EncodeFunctionCall("transfer", to, value)
}

// DecodeSimpleTransferResult decodes a transfer result
func DecodeSimpleTransferResult(data []byte) (bool, error) {
	var success bool
	abi := ExampleUint256()
	err := abi.DecodeFunctionResult("transfer", data, &success)
	return success, err
}