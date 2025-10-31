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
		moduleName    = flag.String("module", "", "Output module, output file will be $module.abi.go")
		prefix        = flag.String("prefix", "", "Prefix for generated types and functions")
		packageName   = flag.String("package", os.Getenv("GOPACKAGE"), "Package name for generated code")
		varName       = flag.String("var", "", "Variable name containing human-readable ABI (for Go source files)")
		extTuplesFlag = flag.String("external-tuples", "", "External tuple mappings in format 'key1=value1,key2=value2'")
		imports       = flag.String("imports", "", "Additional import paths, comma-separated")
	)
	flag.Parse()

	opts := []generator.Option{
		generator.PackageName(*packageName),
		generator.ModuleName(*moduleName),
		generator.Prefix(*prefix),
	}

	if *imports != "" {
		paths := strings.Split(*imports, ",")
		var importSpecs []generator.ImportSpec
		for _, imp := range paths {
			importSpecs = append(importSpecs, parseImport(imp))
		}
		opts = append(opts, generator.ExtraImports(importSpecs))
	}

	// Parse external tuples if provided
	if *extTuplesFlag != "" {
		extTuples := parseExternalTuples(*extTuplesFlag)
		opts = append(opts, generator.ExternalTuples(extTuples))
	}

	generator.Command(
		*inputFile,
		*varName,
		opts...,
	)
}

// parseExternalTuples parses external tuple mappings from string format
// Format: "key1=value1,key2=value2"
func parseExternalTuples(s string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}

	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" && value != "" {
				result[key] = value
			}
		}
	}
	return result
}

// parseImport parses an import string that may contain an alias
// Examples:
//
//	"github.com/ethereum/go-ethereum/common" -> ImportSpec{Path: "github.com/ethereum/go-ethereum/common", Alias: ""}
//	"cmn=github.com/ethereum/go-ethereum/common" -> ImportSpec{Path: "github.com/ethereum/go-ethereum/common", Alias: "cmn"}
func parseImport(imp string) generator.ImportSpec {
	parts := strings.Split(imp, "=")
	var spec generator.ImportSpec

	switch len(parts) {
	case 2:
		spec = generator.ImportSpec{
			Alias: parts[0],
			Path:  parts[1],
		}
	case 1:
		spec = generator.ImportSpec{
			Path: parts[0],
		}
	default:
		panic("invalid import format " + imp)
	}
	return spec
}
