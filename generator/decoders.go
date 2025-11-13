package generator

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// genIntDecoding generates decoding for integer types
func (g *Generator) genIntDecoding(t abi.Type) {
	// Optimize small integer types to avoid big.Int overhead
	if t.Size <= 64 {
		g.genSmallIntDecoding(t)
	} else {
		g.genBigIntDecoding(t)
	}
}

// genSmallIntDecoding generates optimized decoding for small integer types
func (g *Generator) genSmallIntDecoding(t abi.Type) {
	// For small integers, we can use direct binary decoding without big.Int
	// Use the closest native integer type that fits
	var nativeType string

	if t.T == abi.IntTy {
		// Signed integers: use next larger signed type
		if t.Size <= 8 {
			nativeType = "int8"
		} else if t.Size <= 16 {
			nativeType = "int16"
		} else if t.Size <= 32 {
			nativeType = "int32"
		} else if t.Size <= 64 {
			nativeType = "int64"
		} else {
			// > 64 bits: use big.Int
			g.genBigIntDecoding(t)
			return
		}
	} else {
		// Unsigned integers: use next larger unsigned type
		if t.Size <= 8 {
			nativeType = "uint8"
		} else if t.Size <= 16 {
			nativeType = "uint16"
		} else if t.Size <= 32 {
			nativeType = "uint32"
		} else if t.Size <= 64 {
			nativeType = "uint64"
		} else {
			// > 64 bits: use big.Int
			g.genBigIntDecoding(t)
			return
		}
	}

	paddingBytes := 32 - t.Size/8
	if t.T == abi.IntTy {
		// Signed: validate sign extension based on actual ABI type size
		g.L("\t// Validate sign extension for int%d (padding bytes: %d)", t.Size, paddingBytes)
		g.L("\tif data[%d]&0x80 != 0 {", paddingBytes)
		g.L("\t\t// Negative value, check all padding bytes are 0xFF")
		g.L("\t\tfor i := 0; i < %d; i++ {", paddingBytes)
		g.L("\t\t\tif data[i] != 0xFF {")
		g.L("\t\t\t\treturn 0, 0, %sErrDirtyPadding", g.StdPrefix)
		g.L("\t\t\t}")
		g.L("\t\t}")
		g.L("\t} else {")
		g.L("\t\t// Non-negative value, check all padding bytes are zero")
		g.L("\t\tfor i := 0; i < %d; i++ {", paddingBytes)
		g.L("\t\t\tif data[i] != 0x00 {")
		g.L("\t\t\t\treturn 0, 0, %sErrDirtyPadding", g.StdPrefix)
		g.L("\t\t\t}")
		g.L("\t\t}")
		g.L("\t}")
	} else {
		// Unsigned: validate no extra bits are set
		g.L("\t// Validate no extra bits are set for uint%d (padding bytes: %d)", t.Size, paddingBytes)
		g.L("\tfor i := 0; i < %d; i++ {", paddingBytes)
		g.L("\t\tif data[i] != 0x00 {")
		g.L("\t\t\treturn 0, 0, %sErrDirtyPadding", g.StdPrefix)
		g.L("\t\t}")
		g.L("\t}")
	}

	// Decode using the appropriate native type
	switch nativeType {
	case "int8":
		g.L("\tresult := int8(data[31])")
	case "int16":
		g.L("\tresult := int16(binary.BigEndian.Uint16(data[30:32]))")
	case "int32":
		// Handle int24 and int32
		g.L("\t// Decode int%d using int32", t.Size)
		g.L("\tresult := int32(binary.BigEndian.Uint32(data[28:32]))")
	case "int64":
		// Handle int40, int48, int56 and int64
		g.L("\t// Decode int%d using int64", t.Size)
		g.L("\tresult := int64(binary.BigEndian.Uint64(data[24:32]))")
	case "uint8":
		g.L("\tresult := uint8(data[31])")
	case "uint16":
		g.L("\tresult := binary.BigEndian.Uint16(data[30:32])")
	case "uint32":
		g.L("\tresult := binary.BigEndian.Uint32(data[28:32])")
	case "uint64":
		g.L("\tresult := binary.BigEndian.Uint64(data[24:32])")
	default:
		panic("unsupported native type for small integer decoding")
	}

	g.L("\treturn result, 32, nil")
}

// genBigIntDecoding generates decoding for big.Int types
func (g *Generator) genBigIntDecoding(t abi.Type) {
	signed := "false"
	if t.T == abi.IntTy {
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
	g.L("\tfor i := 0; i < 31; i++ {")
	g.L("\t\tif data[i] != 0x00 {")
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
	g.L("\tresult := string(data[:length])")
	g.L("\treturn result, 32 + %sPad32(length), nil", g.StdPrefix)
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
	g.L("\tresult := make([]byte, length)")
	g.L("\tcopy(result, data[:length])")
	g.L("\treturn result, 32 + %sPad32(length), nil", g.StdPrefix)
}

// genFixedBytesDecoding generates decoding for fixed bytes types
func (g *Generator) genFixedBytesDecoding(t abi.Type) {
	g.L("\tvar result [%d]byte", t.Size)
	g.L("\tcopy(result[:], data[:%d])", t.Size)
	g.L("\treturn result, %d, nil", t.Size)
}

// genSliceDecoding generates decoding for slice types
func (g *Generator) genSliceDecoding(t abi.Type) {
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

		if t.Elem.T == abi.TupleTy {
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

		if t.Elem.T == abi.TupleTy {
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
func (g *Generator) genArrayDecoding(t abi.Type) {
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
		if t.Elem.T == abi.TupleTy {
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
