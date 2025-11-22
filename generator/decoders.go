package generator

import (
	"fmt"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

// genIntDecoding generates decoding for integer types
func (g *Generator) genIntDecoding(t ethabi.Type) {
	// Optimize small integer types to avoid big.Int overhead
	if t.Size <= 64 {
		g.genSmallIntDecoding(t)
	} else {
		g.genBigIntDecoding(t)
	}
}

// genSmallIntDecoding generates optimized decoding for small integer types
func (g *Generator) genSmallIntDecoding(t ethabi.Type) {
	if t.Size%8 != 0 {
		panic(fmt.Sprintf("unsupported size %d for small integer decoding", t.Size))
	}

	// For small integers, we can use direct binary decoding without big.Int
	// Use the closest native integer type that fits
	var nativeType string
	if t.T == ethabi.IntTy {
		nativeType = fmt.Sprintf("int%d", nativeSize(t.Size))
	} else {
		nativeType = fmt.Sprintf("uint%d", nativeSize(t.Size))
	}

	if t.T == ethabi.IntTy {
		g.L("\tresult, err := %sDecodeInt[%s](data, %sMinInt%d, %sMaxInt%d)", g.StdPrefix, nativeType, g.StdPrefix, t.Size, g.StdPrefix, t.Size)
	} else {
		g.L("\tresult, err := %sDecodeUint[%s](data, %sMaxUint%d)", g.StdPrefix, nativeType, g.StdPrefix, t.Size)
	}
	g.L("\tif err != nil {")
	g.L("\t\treturn 0, 0, err")
	g.L("\t}")

	g.L("\treturn result, 32, nil")
}

// genBigIntDecoding generates decoding for big.Int types
func (g *Generator) genBigIntDecoding(t ethabi.Type) {
	signed := "false"
	if t.T == ethabi.IntTy {
		signed = "true"
	}

	g.L("\tresult, err := %sDecodeBigInt(data[:32], %s)", g.StdPrefix, signed)
	g.L("\tif err != nil {")
	g.L("\t\treturn nil, 0, err")
	g.L("\t}")
	g.L("\treturn result, 32, nil")
}

// genAddressDecoding generates decoding for address types
func (g *Generator) genAddressDecoding() {
	g.L("\tvar result common.Address")
	g.L("\tfor i := 0; i < 12; i++ {")
	g.L("\t\tif data[i] != 0x00 {")
	g.L("\t\t\treturn result, 0, %sErrDirtyPadding", g.StdPrefix)
	g.L("\t\t}")
	g.L("\t}")
	g.L("\tcopy(result[:], data[12:32])")
	g.L("\treturn result, 32, nil")
}

// genBoolDecoding generates decoding for boolean types
func (g *Generator) genBoolDecoding() {
	g.L("\t// Validate boolean encoding - only 0 or 1 are valid")
	g.L("\tfor _, i := range data[:31] {")
	g.L("\t\tif i != 0 {")
	g.L("\t\t\treturn false, 0, %sErrDirtyPadding", g.StdPrefix)
	g.L("\t\t}")
	g.L("\t}")
	g.L("\tswitch data[31] {")
	g.L("\tcase 0x01:")
	g.L("\t\treturn true, 32, nil")
	g.L("\tcase 0x00:")
	g.L("\t\treturn false, 32, nil")
	g.L("\tdefault:")
	g.L("\t\treturn false, 0, %sErrDirtyPadding", g.StdPrefix)
	g.L("\t}")
}

// genStringDecoding generates decoding for string types
func (g *Generator) genStringDecoding() {
	g.L("\t// Decode length")

	g.L("\tif len(data) < 32 {")
	g.L("\t\treturn \"\", 0, io.ErrUnexpectedEOF")
	g.L("\t}")

	g.L("\tlength, err := %sDecodeSize(data)", g.StdPrefix)
	g.L("\tif err != nil {")
	g.L("\t\treturn \"\", 0, err")
	g.L("\t}")
	g.L("\tdata = data[32:]")

	g.L("\tpaddedLength := %sPad32(length)", g.StdPrefix)
	g.L("\tif len(data) < paddedLength {")
	g.L("\t\treturn \"\", 0, io.ErrUnexpectedEOF")
	g.L("\t}")

	g.L("\t// check padding bytes")
	g.L("\tfor i := length; i < paddedLength; i++ {")
	g.L("\t\tif data[i] != 0x00 {")
	g.L("\t\t\treturn \"\", 0, %sErrDirtyPadding", g.StdPrefix)
	g.L("\t\t}")
	g.L("\t}")

	g.L("")
	g.L("\t// Decode data")
	g.L("\treturn string(data[:length]), 32 + %sPad32(length), nil", g.StdPrefix)
}

// genBytesDecoding generates decoding for bytes types
func (g *Generator) genBytesDecoding() {
	g.L("\t// Decode length")

	g.L("\tif len(data) < 32 {")
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")

	g.L("\tlength, err := %sDecodeSize(data)", g.StdPrefix)
	g.L("\tif err != nil {")
	g.L("\t\treturn nil, 0, err")
	g.L("\t}")
	g.L("\tdata = data[32:]")

	g.L("\tpaddedLength := %sPad32(length)", g.StdPrefix)
	g.L("\tif len(data) < paddedLength {")
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")

	g.L("\t// check padding bytes")
	g.L("\tfor i := length; i < paddedLength; i++ {")
	g.L("\t\tif data[i] != 0x00 {")
	g.L("\t\t\treturn nil, 0, %sErrDirtyPadding", g.StdPrefix)
	g.L("\t\t}")
	g.L("\t}")

	g.L("")
	g.L("\t// Decode data")
	g.L("\treturn data[:length], 32 + %sPad32(length), nil", g.StdPrefix)
}

// genFixedBytesDecoding generates decoding for fixed bytes types
func (g *Generator) genFixedBytesDecoding(t ethabi.Type) {
	// Validate padding bytes
	g.L("\t// Validate padding bytes for fixed bytes[%d]", t.Size)
	g.L("\tfor i := %d; i < 32; i++ {", t.Size)
	g.L("\t\tif data[i] != 0x00 {")
	g.L("\t\t\treturn [%d]byte{}, 0, %sErrDirtyPadding", t.Size, g.StdPrefix)
	g.L("\t\t}")
	g.L("\t}")
	g.L("\tvar result [%d]byte", t.Size)
	g.L("\tcopy(result[:], data[:%d])", t.Size)
	g.L("\treturn result, %d, nil", t.Size)
}

// genSliceDecoding generates decoding for slice types
func (g *Generator) genSliceDecoding(t ethabi.Type) {
	g.L("\t// Decode length")

	g.L("\tif len(data) < 32 {")
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")

	g.L("\tlength, err := %sDecodeSize(data)", g.StdPrefix)
	g.L("\tif err != nil {")
	g.L("\t\treturn nil, 0, err")
	g.L("\t}")

	g.L("\tdata = data[32:]")
	g.L("\t\tif length > len(data) || length * %d > len(data) {", GetTypeSize(*t.Elem))
	g.L("\t\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t\t}")

	g.L("\tvar (")
	g.L("\t\tn int")
	g.L("\t\toffset int")
	g.L("\t)")

	goType := g.abiTypeToGoType(*t.Elem)
	if !IsDynamicType(*t.Elem) {
		g.L("\t// Decode elements with static types")
		g.L("\tresult := make([]%s, length)", goType)
		g.L("\tfor i := 0; i < length; i++ {")

		if t.Elem.T == ethabi.TupleTy {
			g.L("\t\tn, err = result[i].Decode(data[offset:])")
		} else {
			g.L("\t\tresult[i], n, err = %s", g.genDecodeCall(*t.Elem, "data[offset:]"))
		}

		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")
		g.L("\t\toffset += n")
		g.L("\t}")
		g.L("\treturn result, offset + 32, nil")
	} else {
		g.L("\t// Decode elements with dynamic types")
		g.L("\tresult := make([]%s, length)", goType)
		g.L("\tdynamicOffset := length * 32")
		g.L("\tfor i := 0; i < length; i++ {")
		g.L("\t\ttmp, err := %sDecodeSize(data[offset:])", g.StdPrefix)
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")
		g.L("\t\toffset += 32")
		g.L("")
		g.L("\t\tif dynamicOffset != tmp {")
		g.L("\t\t\treturn nil, 0, %sErrInvalidOffsetForSliceElement", g.StdPrefix)
		g.L("\t\t}")

		if t.Elem.T == ethabi.TupleTy {
			g.L("\t\tn, err = result[i].Decode(data[dynamicOffset:])")
		} else {
			g.L("\t\tresult[i], n, err = %s", g.genDecodeCall(*t.Elem, "data[dynamicOffset:]"))
		}

		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn nil, 0, err")
		g.L("\t\t}")
		g.L("\t\tdynamicOffset += n")
		g.L("\t}")
		g.L("\treturn result, dynamicOffset + 32, nil")
	}
}

// genArrayDecoding generates decoding for array types
func (g *Generator) genArrayDecoding(t ethabi.Type) {
	goType := g.abiTypeToGoType(*t.Elem)
	typeSize := GetTypeSize(*t.Elem)

	if !IsDynamicType(*t.Elem) {
		g.L("\t// Decode fixed-size array with static elements")
		g.L("\tvar (")
		g.L("\t\tresult [%d]%s", t.Size, goType)
		g.L("\t\terr error")
		g.L("\t)")
		g.L("\tif len(data) < %d {", t.Size*typeSize)
		g.L("\t\treturn result, 0, io.ErrUnexpectedEOF")
		g.L("\t}")

		var offset int
		for i := 0; i < t.Size; i++ {
			g.L("\t// Element %d", i)
			g.L("\tresult[%d], _, err = %s", i, g.genDecodeCall(*t.Elem, fmt.Sprintf("data[%d:]", offset)))
			g.L("\tif err != nil {")
			g.L("\t\treturn result, 0, err")
			g.L("\t}")
			offset += typeSize
		}
		g.L("\treturn result, %d, nil", offset)
	} else {
		g.L("\t// Decode fixed-size array with dynamic elements")
		g.L("\tvar result [%d]%s", t.Size, goType)

		g.L("\tif len(data) < %d {", t.Size*32)
		g.L("\t\treturn result, 0, io.ErrUnexpectedEOF")
		g.L("\t}")

		g.L("\tvar (")
		g.L("\t\tn int")
		g.L("\t\terr error")
		g.L("\t\ttmp int")
		g.L("\t)")
		g.L("\toffset := 0")
		g.L("\tdynamicOffset := %d", t.Size*32)
		g.L("\tfor i := 0; i < %d; i++ {", t.Size)
		g.L("\t\ttmp, err = %sDecodeSize(data[offset:])", g.StdPrefix)
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn result, 0, err")
		g.L("\t\t}")
		g.L("\t\toffset += 32")
		g.L("")
		g.L("\t\tif dynamicOffset != tmp {")
		g.L("\t\t\treturn result, 0, %sErrInvalidOffsetForArrayElement", g.StdPrefix)
		g.L("\t\t}")
		if t.Elem.T == ethabi.TupleTy {
			g.L("\t\tn, err = result[i].Decode(data[dynamicOffset:])")
		} else {
			g.L("\t\tresult[i], n, err = %s", g.genDecodeCall(*t.Elem, "data[dynamicOffset:]"))
		}
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn result, 0, err")
		g.L("\t\t}")
		g.L("\t\tdynamicOffset += n")
		g.L("\t}")
		g.L("\treturn result, dynamicOffset, nil")
	}
}
