package generator

import (
	"fmt"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/yihuang/go-abi"
)

// genIntEncoding generates encoding for integer types
func (g *Generator) genIntEncoding(t ethabi.Type) {
	// Optimize small integer types to avoid big.Int overhead
	if t.Size <= 64 {
		g.genSmallIntEncoding(t)
	} else {
		g.genBigIntEncoding(t)
	}
}

// genSmallIntEncoding generates optimized encoding for small integer types
func (g *Generator) genSmallIntEncoding(t ethabi.Type) {
	// For small integers, we can use direct binary encoding without big.Int
	// Use the closest native integer type that fits
	var nativeType string

	if t.T == ethabi.IntTy {
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
			g.genBigIntEncoding(t)
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
			g.genBigIntEncoding(t)
			return
		}
	}

	if t.T == ethabi.IntTy {
		// Signed: use the appropriate integer size and sign extend
		switch nativeType {
		case "int8":
			g.L("\tif value < 0 {")
			g.L("\t\tfor i := 0; i < 31; i++ { buf[i] = 0xff }")
			g.L("\t}")
			g.L("\tbuf[31] = byte(value)")
		case "int16":
			g.L("\tif value < 0 {")
			g.L("\t\tfor i := 0; i < 30; i++ { buf[i] = 0xff }")
			g.L("\t}")
			g.L("\tbinary.BigEndian.PutUint16(buf[30:32], uint16(value))")
		case "int32":
			g.L("\tif value < 0 {")
			g.L("\t\tfor i := 0; i < 28; i++ { buf[i] = 0xff }")
			g.L("\t}")
			g.L("\tbinary.BigEndian.PutUint32(buf[28:32], uint32(value))")
		case "int64":
			g.L("\tif value < 0 {")
			g.L("\t\tfor i := 0; i < 24; i++ { buf[i] = 0xff }")
			g.L("\t}")
			g.L("\tbinary.BigEndian.PutUint64(buf[24:32], uint64(value))")
		}
	} else {
		// Unsigned: use the appropriate integer size
		switch nativeType {
		case "uint8":
			g.L("\tbuf[31] = byte(value)")
		case "uint16":
			g.L("\tbinary.BigEndian.PutUint16(buf[30:32], uint16(value))")
		case "uint32":
			g.L("\tbinary.BigEndian.PutUint32(buf[28:32], uint32(value))")
		case "uint64":
			g.L("\tbinary.BigEndian.PutUint64(buf[24:32], uint64(value))")
		}
	}

	g.L("\treturn 32, nil")
}

// genBigIntEncoding generates encoding for big.Int types
func (g *Generator) genBigIntEncoding(t ethabi.Type) {
	signed := "false"
	if t.T == ethabi.IntTy {
		signed = "true"
	}

	g.L("\tif err := %sEncodeBigInt(value, buf[:32], %s); err != nil {", g.StdPrefix, signed)
	g.L("\t\treturn 0, err")
	g.L("\t}")
	g.L("\treturn 32, nil")
}

// genAddressEncoding generates encoding for address types
func (g *Generator) genAddressEncoding() {
	g.L("\tcopy(buf[12:32], value[:])")
	g.L("\treturn 32, nil")
}

// genBoolEncoding generates encoding for boolean types
func (g *Generator) genBoolEncoding() {
	g.L("\tif value {")
	g.L("\t\tbuf[31] = 1")
	g.L("\t}")
	g.L("\treturn 32, nil")
}

// genStringEncoding generates encoding for string types
func (g *Generator) genStringEncoding() {
	g.L("\t// Encode length")
	g.L("\tbinary.BigEndian.PutUint64(buf[24:32], uint64(len(value)))")
	g.L("\t")
	g.L("\t// Encode data")
	g.L("\tcopy(buf[32:], []byte(value))")
	g.L("\t")
	g.L("\treturn 32 + %sPad32(len(value)), nil", g.StdPrefix)
}

// genBytesEncoding generates encoding for bytes types
func (g *Generator) genBytesEncoding() {
	g.L("\t// Encode length")
	g.L("\tbinary.BigEndian.PutUint64(buf[24:32], uint64(len(value)))")
	g.L("\t")
	g.L("\t// Encode data")
	g.L("\tcopy(buf[32:], value)")
	g.L("\t")
	g.L("\treturn 32 + %sPad32(len(value)), nil", g.StdPrefix)
}

// genFixedBytesEncoding generates encoding for fixed bytes types
func (g *Generator) genFixedBytesEncoding(t ethabi.Type) {
	g.L("\tcopy(buf[:%d], value[:])", t.Size)
	g.L("\treturn %d, nil", t.Size)
}

// genSliceEncoding generates encoding for slice types
func (g *Generator) genSliceEncoding(t ethabi.Type) {
	g.L("\t// Encode length")
	g.L("\tbinary.BigEndian.PutUint64(buf[24:32], uint64(len(value)))")
	g.L("\tbuf = buf[32:]")
	g.L("\t")
	if !IsDynamicType(*t.Elem) {
		g.L("\t// Encode elements with static types")
		g.L("\tvar offset int")
		g.L("\tfor _, elem := range value {")
		g.L("\t\tn, err := %s", g.genEncodeCall(*t.Elem, "elem", "buf[offset:]"))
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn 0, err")
		g.L("\t\t}")
		g.L("\t\toffset += n")
		g.L("\t}")
		g.L("\t")
		g.L("\treturn offset + 32, nil")
	} else {
		g.L("\t// Encode elements with dynamic types")
		g.L("\tvar offset int")
		g.L("\tdynamicOffset := len(value)*32")
		g.L("\tfor _, elem := range value {")
		g.L("\t\t// Write offset for element")
		g.L("\t\toffset += 32")
		g.L("\t\tbinary.BigEndian.PutUint64(buf[offset-8:offset], uint64(dynamicOffset))")
		g.L("")
		g.L("\t\t// Write element at dynamic region")
		g.L("\t\tn, err := %s", g.genEncodeCall(*t.Elem, "elem", "buf[dynamicOffset:]"))
		g.L("\t\tif err != nil {")
		g.L("\t\t\treturn 0, err")
		g.L("\t\t}")
		g.L("\t\tdynamicOffset += n")
		g.L("\t}")
		g.L("\t")
		g.L("\treturn dynamicOffset + 32, nil")
	}
}

// genArrayEncoding generates encoding for array types
func (g *Generator) genArrayEncoding(t ethabi.Type) {
	if !IsDynamicType(*t.Elem) {
		g.L("\t// Encode fixed-size array with static elements")

		var offset int
		for i := 0; i < t.Size; i++ {
			ref := fmt.Sprintf("value[%d]", i)

			g.L("\t\tif _, err := %s; err != nil {", g.genEncodeCall(*t.Elem, ref, fmt.Sprintf("buf[%d:]", offset)))
			g.L("\t\t\treturn 0, err")
			g.L("\t\t}")

			offset += GetTypeSize(*t.Elem)
		}
		g.L("\t")
		g.L("\treturn %d, nil", offset)
	} else {
		g.L("\t// Encode fixed-size array with dynamic elements")

		var offset int

		g.L("\tvar (")
		g.L("\t\tn int")
		g.L("\t\terr error")
		g.L("\t)")

		g.L("\tdynamicOffset := 32 * %d", t.Size)
		for i := 0; i < t.Size; i++ {
			g.L("\tbinary.BigEndian.PutUint64(buf[%d+24:%d+32], uint64(dynamicOffset))", offset, offset)
			offset += 32

			ref := fmt.Sprintf("value[%d]", i)
			g.L("\tn, err = %s", g.genEncodeCall(*t.Elem, ref, "buf[dynamicOffset:]"))
			g.L("\tif err != nil {")
			g.L("\t\treturn 0, err")
			g.L("\t}")
			g.L("\tdynamicOffset += n")
			g.L("\t")
		}
		g.L("\t")
		g.L("\treturn dynamicOffset, nil")
	}
}

// genTupleEncoding generates encoding for tuple types
func (g *Generator) genTupleEncoding(t ethabi.Type) {
	g.L("\t// Encode tuple fields")
	g.L("\tdynamicOffset := %sStaticSize // Start dynamic data after static section", abi.TupleStructName(t))

	// Generate encoding for each tuple element
	if IsDynamicType(t) {
		g.L("\tvar (")
		g.L("\t\terr error")
		g.L("\t\tn int")
		g.L("\t)")
	}

	var offset int
	for i, elem := range t.TupleElems {
		// Generate field access - use meaningful field names if available
		fieldName := GoFieldName(t.TupleRawNames[i])
		if fieldName == "" {
			fieldName = fmt.Sprintf("Field%d", i+1)
		}
		g.L("\t// Field %s: %s", fieldName, elem.String())

		ref := "value." + fieldName
		if !IsDynamicType(*elem) {
			// Static field - encode directly
			g.L("\tif _, err := %s; err != nil {", g.genEncodeCall(*elem, ref, fmt.Sprintf("buf[%d:]", offset)))
			g.L("\t\treturn 0, err")
			g.L("\t}")
			offset += GetTypeSize(*elem)
		} else {
			// Dynamic field - encode offset pointer and data in dynamic section
			g.L("\t// Encode offset pointer")
			g.L("\tbinary.BigEndian.PutUint64(buf[%d+24:%d+32], uint64(dynamicOffset))", offset, offset)
			offset += 32

			g.L("\t// Encode dynamic data")
			g.L("\tn, err = %s", g.genEncodeCall(*elem, ref, "buf[dynamicOffset:]"))
			g.L("\tif err != nil {")
			g.L("\t\treturn 0, err")
			g.L("\t}")
			g.L("\tdynamicOffset += n")
		}
		g.L("")
	}

	g.L("\treturn dynamicOffset, nil")
}
