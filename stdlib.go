package abi

import (
	"bytes"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

//go:generate go run ./cmd -var StdlibABI -module=stdlib -stdlib

var StdlibABI = []string{
	"function stdlib(uint8,uint16,uint32,uint64,uint256, int8,int16,int32,int64,int256, address,bytes32,string,bytes, uint8[],uint16[],uint32[],uint64[],uint256[], int8[],int16[],int32[],int64[],int256[], address[],bytes32[],string[],bytes[]) view returns ()",
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

func IsStdlibType(ident string) bool {
	_, ok := stdlibTypes[ident]
	return ok
}
