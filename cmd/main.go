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
		outputFile    = flag.String("output", "", "Output Go file")
		packageName   = flag.String("package", os.Getenv("GOPACKAGE"), "Package name for generated code")
		varName       = flag.String("var", "", "Variable name containing human-readable ABI (for Go source files)")
		extTuplesFlag = flag.String("external-tuples", "", "External tuple mappings in format 'key1=value1,key2=value2'")
		imports       = flag.String("imports", "", "Additional import paths, comma-separated")
	)
	flag.Parse()

	opts := []generator.Option{
		generator.PackageName(*packageName),
	}

	if *imports != "" {
		importPaths := strings.Split(*imports, ",")
		opts = append(opts, generator.ExtraImports(importPaths))
	}

	// Parse external tuples if provided
	if *extTuplesFlag != "" {
		extTuples := parseExternalTuples(*extTuplesFlag)
		opts = append(opts, generator.ExternalTuples(extTuples))
	}

	generator.Command(
		*inputFile,
		*outputFile,
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
