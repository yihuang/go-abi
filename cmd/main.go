package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/imports"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

var DefaultPackageName = os.Getenv("GOPACKAGE")

func init() {
	if DefaultPackageName == "" {
		DefaultPackageName = "testdata"
	}
}

// parseHumanReadableABIFromFile parses a Go source file and extracts human-readable ABI from a variable
func parseHumanReadableABIFromFile(filename, varName string) (ethabi.ABI, error) {
	// Parse the Go source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return ethabi.ABI{}, fmt.Errorf("failed to parse Go file: %w", err)
	}

	// Find the specified variable
	var abiLines []string
	ast.Inspect(node, func(n ast.Node) bool {
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, name := range valueSpec.Names {
						if name.Name == varName {
							// Found the variable, extract its value
							if i < len(valueSpec.Values) {
								if lit, ok := valueSpec.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
									// Single string value
									abiLines = append(abiLines, strings.Trim(lit.Value, `"`))
								} else if compLit, ok := valueSpec.Values[i].(*ast.CompositeLit); ok {
									// Array/slice literal
									for _, elt := range compLit.Elts {
										if lit, ok := elt.(*ast.BasicLit); ok && lit.Kind == token.STRING {
											abiLines = append(abiLines, strings.Trim(lit.Value, `"`))
										}
									}
								}
							}
							return false // Stop searching
						}
					}
				}
			}
		}
		return true
	})

	if len(abiLines) == 0 {
		return ethabi.ABI{}, fmt.Errorf("variable %s not found or has no string value", varName)
	}

	// Parse human-readable ABI
	abiJSON, err := abi.ParseHumanReadableABI(abiLines)
	if err != nil {
		return ethabi.ABI{}, fmt.Errorf("failed to parse human-readable ABI: %w", err)
	}

	// Convert to go-ethereum ABI
	return ethabi.JSON(bytes.NewReader(abiJSON))
}

func main() {
	var (
		inputFile   = flag.String("input", "", "Input file (JSON ABI or Go source file)")
		outputFile  = flag.String("output", "", "Output Go file")
		packageName = flag.String("package", DefaultPackageName, "Package name for generated code")
		varName     = flag.String("var", "", "Variable name containing human-readable ABI (for Go source files)")
	)
	flag.Parse()

	var abiDef ethabi.ABI
	var err error

	if *inputFile == "" {
		// If no input file specified, use GOFILE environment variable
		goFile := os.Getenv("GOFILE")
		if goFile == "" {
			log.Fatal("-input flag is required or must be run via go generate")
		}
		*inputFile = goFile
	}

	// Determine input type by file extension
	if strings.HasSuffix(*inputFile, ".go") {
		// Go source file - requires -var flag
		if *varName == "" {
			log.Fatal("-var flag is required when input is a Go source file")
		}
		abiDef, err = parseHumanReadableABIFromFile(*inputFile, *varName)
		if err != nil {
			log.Fatalf("Failed to parse human-readable ABI from variable %s in file %s: %v", *varName, *inputFile, err)
		}
	} else if strings.HasSuffix(*inputFile, ".json") {
		// JSON ABI file
		abiJSON, err := os.ReadFile(*inputFile)
		if err != nil {
			log.Fatalf("Failed to read input file: %v", err)
		}

		abiDef, err = ethabi.JSON(bytes.NewReader(abiJSON))
		if err != nil {
			log.Fatalf("Failed to parse ABI JSON: %v", err)
		}
	} else {
		log.Fatalf("Unsupported input file type: %s (expected .go or .json)", *inputFile)
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
		return
	}

	opt := imports.Options{
		Comments: true,
	}
	formatted, err := imports.Process(*outputFile, []byte(generatedCode), &opt)
	if err != nil {
		log.Printf("Raw generated code before formatting:%s\n", generatedCode)
		log.Fatalf("failed to format generated code: %v", err)
	}

	if err := os.WriteFile(*outputFile, formatted, 0644); err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}
	fmt.Printf("Generated code written to %s\n", *outputFile)
}
