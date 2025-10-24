package main

import (
	"flag"
	"os"

	"github.com/yihuang/go-abi/generator"
)

func main() {
	var (
		inputFile   = flag.String("input", os.Getenv("GOFILE"), "Input file (JSON ABI or Go source file)")
		outputFile  = flag.String("output", "", "Output Go file")
		packageName = flag.String("package", os.Getenv("GOPACKAGE"), "Package name for generated code")
		varName     = flag.String("var", "", "Variable name containing human-readable ABI (for Go source files)")
	)
	flag.Parse()

	generator.Command(
		*inputFile,
		*outputFile,
		*varName,
		generator.PackageName(*packageName),
	)
}
