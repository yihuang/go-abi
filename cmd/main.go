package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

var DefaultPackageName = os.Getenv("GOPACKAGE")

func init() {
	if DefaultPackageName == "" {
		DefaultPackageName = "testdata"
	}
}

func main() {
	var (
		inputFile   = flag.String("input", "", "Input ABI JSON file")
		outputFile  = flag.String("output", "", "Output Go file")
		packageName = flag.String("package", DefaultPackageName, "Package name for generated code")
	)
	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Input file is required")
	}

	// Read ABI JSON
	abiJSON, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	abiDef, err := ethabi.JSON(bytes.NewReader(abiJSON))
	if err != nil {
		log.Fatalf("Failed to parse ABI JSON: %v", err)
	}

	// Generate code
	generator := abi.NewGenerator(*packageName)
	generatedCode, err := generator.GenerateFromABI(abiDef)
	if err != nil {
		log.Printf("Raw generated code before formatting:%s\n", generatedCode)
		log.Fatalf("Failed to generate code: %v", err)
	}

	// Write output
	if *outputFile == "" {
		fmt.Println(generatedCode)
	} else {
		if err := os.WriteFile(*outputFile, []byte(generatedCode), 0644); err != nil {
			log.Fatalf("Failed to write output file: %v", err)
		}
		fmt.Printf("Generated code written to %s\n", *outputFile)
	}
}
