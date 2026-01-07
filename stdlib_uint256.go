//go:build uint256

package abi

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

//go:generate go run ./cmd -var StdlibABI -output=stdlib_uint256.abi.go -stdlib -uint256 -buildtag=uint256

var StdlibABI = []string{
	"function basic(bool,address,bytes32,string,bytes,bool[],address[],bytes32[],string[],bytes[]) returns ()",
	"function ints(uint8,int8,uint16,int16,uint24,int24,uint32,int32,uint40,int40,uint48,int48,uint56,int56,uint64,int64,uint72,int72,uint80,int80,uint88,int88,uint96,int96,uint104,int104,uint112,int112,uint120,int120,uint128,int128,uint136,int136,uint144,int144,uint152,int152,uint160,int160,uint168,int168,uint176,int176,uint184,int184,uint192,int192,uint200,int200,uint208,int208,uint216,int216,uint224,int224,uint232,int232,uint240,int240,uint248,int248,uint256,int256,uint8[],int8[],uint16[],int16[],uint24[],int24[],uint32[],int32[],uint40[],int40[],uint48[],int48[],uint56[],int56[],uint64[],int64[],uint72[],int72[],uint80[],int80[],uint88[],int88[],uint96[],int96[],uint104[],int104[],uint112[],int112[],uint120[],int120[],uint128[],int128[],uint136[],int136[],uint144[],int144[],uint152[],int152[],uint160[],int160[],uint168[],int168[],uint176[],int176[],uint184[],int184[],uint192[],int192[],uint200[],int200[],uint208[],int208[],uint216[],int216[],uint224[],int224[],uint232[],int232[],uint240[],int240[],uint248[],int248[],uint256[],int256[]) returns ()",
	"function bytes(bytes1,bytes2,bytes3,bytes4,bytes5,bytes6,bytes7,bytes8,bytes9,bytes10,bytes11,bytes12,bytes13,bytes14,bytes15,bytes16,bytes17,bytes18,bytes19,bytes20,bytes21,bytes22,bytes23,bytes24,bytes25,bytes26,bytes27,bytes28,bytes29,bytes30,bytes31,bytes32,bytes1[],bytes2[],bytes3[],bytes4[],bytes5[],bytes6[],bytes7[],bytes8[],bytes9[],bytes10[],bytes11[],bytes12[],bytes13[],bytes14[],bytes15[],bytes16[],bytes17[],bytes18[],bytes19[],bytes20[],bytes21[],bytes22[],bytes23[],bytes24[],bytes25[],bytes26[],bytes27[],bytes28[],bytes29[],bytes30[],bytes31[],bytes32[]) returns ()",
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
	for _, method := range abi.Methods {
		for _, input := range method.Inputs {
			stdlibTypes[GenTypeIdentifier(input.Type)] = struct{}{}
		}
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
