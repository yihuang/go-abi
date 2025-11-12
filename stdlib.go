package abi

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

//go:generate go run ./cmd -var StdlibABI -output=stdlib.abi.go -stdlib

var StdlibABI = []string{
	"function stdlib(bool,address,bytes32,string,bytes,uint8,int8,uint16,int16,uint24,int24,uint32,int32,uint40,int40,uint48,int48,uint56,int56,uint64,int64,uint72,int72,uint80,int80,uint88,int88,uint96,int96,uint104,int104,uint112,int112,uint120,int120,uint128,int128,bool[],address[],bytes32[],string[],bytes[],uint8[],int8[],uint16[],int16[],uint24[],int24[],uint32[],int32[],uint40[],int40[],uint48[],int48[],uint56[],int56[],uint64[],int64[],uint72[],int72[],uint80[],int80[],uint88[],int88[],uint96[],int96[],uint104[],int104[],uint112[],int112[],uint120[],int120[],uint128[],int128[]) returns ()",
}

var stdlibTypes map[string]struct{}

func init() {
	bz, err := ParseHumanReadableABI(StdlibABI)
	if err != nil {
		panic(err)
	}
	abi, err := ethabi.JSON(bytes.NewReader(bz))
	if err != nil {
		panic(err)
	}

	stdlibTypes = make(map[string]struct{})
	for _, input := range abi.Methods["stdlib"].Inputs {
		stdlibTypes[GenTypeIdentifier(input.Type)] = struct{}{}
	}
}

// GenStdlibSignature generates the standard library function signature.
//
// run this in a go playground and copy the result to the StdlibABI variable above.
func GenStdlibSignature() string {
	primitives := []string{
		"bool",
		"address",
		"bytes32",
		"string",
		"bytes",
	}

	// common integers
	commonInts := []int{8, 16, 24, 32, 40, 48, 56, 64, 72, 80, 88, 96, 104, 112, 120, 128}

	for _, i := range commonInts {
		primitives = append(primitives, fmt.Sprintf("uint%d", i))
		primitives = append(primitives, fmt.Sprintf("int%d", i))
	}

	types := slices.Clone(primitives)
	for _, p := range primitives {
		types = append(types, p+"[]")
	}

	return fmt.Sprintf("function stdlib(%s) returns ()", strings.Join(types, ","))
}

func IsStdlibType(ident string) bool {
	_, ok := stdlibTypes[ident]
	return ok
}
