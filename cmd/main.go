package main

import (
	"flag"
	"os"
	"strings"

	"github.com/yihuang/go-abi/generator"
)

func main() {
	var (
		inputFile     = flag.String("input", os.Getenv("GOFILE"), "Input file (JSON ABI or Go source file)")
		outputFile    = flag.String("output", "", "Output file")
		prefix        = flag.String("prefix", "", "Prefix for generated types and functions")
		packageName   = flag.String("package", os.Getenv("GOPACKAGE"), "Package name for generated code")
		varName       = flag.String("var", "", "Variable name containing human-readable ABI (for Go source files)")
		extTuplesFlag = flag.String("external-tuples", "", "External tuple mappings in format 'key1=value1,key2=value2'")
		imports       = flag.String("imports", "", "Additional import paths, comma-separated")
		stdlib        = flag.Bool("stdlib", false, "Generate stdlib itself")
		artifactInput = flag.Bool("artifact-input", false, "Input file is a solc artifact JSON, will extract the abi field from it")
		packed        = flag.Bool("packed", false, "Generate packed encoding format (no padding, no dynamic types)")
	)
	flag.Parse()

	opts := []generator.Option{
		generator.PackageName(*packageName),
		generator.Prefix(*prefix),
		generator.Stdlib(*stdlib),
		generator.Packed(*packed),
	}

	if *imports != "" {
		paths := strings.Split(*imports, ",")
		var importSpecs []generator.ImportSpec
		for _, imp := range paths {
			importSpecs = append(importSpecs, generator.ParseImport(imp))
		}
		opts = append(opts, generator.ExtraImports(importSpecs))
	}

	// Parse external tuples if provided
	if *extTuplesFlag != "" {
		extTuples := generator.ParseExternalTuples(*extTuplesFlag)
		opts = append(opts, generator.ExternalTuples(extTuples))
	}

	generator.Command(
		*inputFile,
		*varName,
		*artifactInput,
		*outputFile,
		opts...,
	)
}
