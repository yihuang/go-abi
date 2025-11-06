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

	g.L("\tvar result %s", nativeType)

	if t.T == abi.IntTy {
		// Signed: use the appropriate integer size and sign extend
		switch nativeType {
		case "int8":
			g.L("\tresult = int8(data[31])")
			g.L("\tif data[0]&0x80 != 0 { // Check sign bit")
			g.L("\t\tresult = result | ^0x7f // Sign extend")
			g.L("\t}")
		case "int16":
			g.L("\tresult = int16(binary.BigEndian.Uint16(data[30:32]))")
			g.L("\tif data[0]&0x80 != 0 { // Check sign bit")
			g.L("\t\tresult = result | ^0x7fff // Sign extend")
			g.L("\t}")
		case "int32":
			g.L("\tresult = int32(binary.BigEndian.Uint32(data[28:32]))")
			g.L("\tif data[0]&0x80 != 0 { // Check sign bit")
			g.L("\t\tresult = result | ^0x7fffffff // Sign extend")
			g.L("\t}")
		case "int64":
			g.L("\tresult = int64(binary.BigEndian.Uint64(data[24:32]))")
			g.L("\tif data[0]&0x80 != 0 { // Check sign bit")
			g.L("\t\tresult = result | ^0x7fffffffffffffff // Sign extend")
			g.L("\t}")
		}
	} else {
		// Unsigned: use the appropriate integer size
		switch nativeType {
		case "uint8":
			g.L("\tresult = uint8(data[31])")
		case "uint16":
			g.L("\tresult = binary.BigEndian.Uint16(data[30:32])")
		case "uint32":
			g.L("\tresult = binary.BigEndian.Uint32(data[28:32])")
		case "uint64":
			g.L("\tresult = binary.BigEndian.Uint64(data[24:32])")
		}
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
	g.L("\tcopy(result[:], data[12:32])")
	g.L("\treturn result, 32, nil")
}

// genBoolDecoding generates decoding for boolean types
func (g *Generator) genBoolDecoding() {
	g.L("\tresult := data[31] != 0")
	g.L("\treturn result, 32, nil")
}

// genStringDecoding generates decoding for string types
func (g *Generator) genStringDecoding() {
	g.L("\t// Decode length")
	g.L("\tlength := int(binary.BigEndian.Uint64(data[24:32]))")
	g.L("\tif len(data) < 32 + %sPad32(length) {", g.StdPrefix)
	g.L("\t\treturn \"\", 0, io.ErrUnexpectedEOF")
	g.L("\t}")
	g.L("")
	g.L("\t// Decode data")
	g.L("\tresult := string(data[32:32+length])")
	g.L("\treturn result, 32 + %sPad32(length), nil", g.StdPrefix)
}

// genBytesDecoding generates decoding for bytes types
func (g *Generator) genBytesDecoding() {
	g.L("\t// Decode length")
	g.L("\tlength := int(binary.BigEndian.Uint64(data[24:32]))")
	g.L("\tif len(data) < 32 + %sPad32(length) {", g.StdPrefix)
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")
	g.L("")
	g.L("\t// Decode data")
	g.L("\tresult := make([]byte, length)")
	g.L("\tcopy(result, data[32:32+length])")
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
	g.L("\tlength := int(binary.BigEndian.Uint64(data[24:32]))")
	g.L("\tif len(data) < 32 {")
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")

	g.L("\tdata = data[32:]")
	g.L("\t\tif len(data) < %d * length {", GetTypeSize(*t.Elem))
	g.L("\t\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t\t}")

	g.L("\tvar (")
	g.L("\t\tn int")
	g.L("\t\terr error")
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
		g.L("\t\toffset += 32")
		g.L("\t\ttmp := int(binary.BigEndian.Uint64(data[offset-8:offset]))")
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
		g.L("\t)")
		g.L("\toffset := 0")
		g.L("\tdynamicOffset := %d", t.Size*32)
		g.L("\tfor i := 0; i < %d; i++ {", t.Size)
		g.L("\t\toffset += 32")
		g.L("\t\ttmp := int(binary.BigEndian.Uint64(data[offset-8:offset]))")
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
