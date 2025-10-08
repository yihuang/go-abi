package main

import (
	"fmt"
	"log"
	"os"

	"github.com/huangyi/go-abi/pkg/abi"
)

func main() {
	// Example source code with ABI annotations
	source := `
package main

// abi:generate
// Transfer represents a token transfer
type Transfer struct {
	From  [20]byte
	To    [20]byte
	Value *big.Int
}

// abi:generate
// BalanceOf represents a balance query
type BalanceOf struct {
	Address [20]byte
}
`

	// Create generator
	generator := abi.NewGenerator()

	// Generate ABI bindings
	generated, err := generator.GenerateFromSource(source, "main")
	if err != nil {
		log.Fatal("Generation failed:", err)
	}

	// Write generated code to file
	if err := os.WriteFile("generated_bindings.go", []byte(generated), 0644); err != nil {
		log.Fatal("Writing generated code failed:", err)
	}

	fmt.Println("Generated bindings written to generated_bindings.go")
	fmt.Println("\nGenerated code preview:")
	fmt.Println(generated[:500] + "...")
}