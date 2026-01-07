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
	} else if t.T == ethabi.UintTy && g.Options.UseUint256 {
		g.genUint256Encoding()
	} else {
		g.genBigIntEncoding(t)
	}
}

// genUint256Encoding generates encoding for holiman/uint256.Int types
func (g *Generator) genUint256Encoding() {
	g.L("\tvalue.WriteToArray32((*[32]byte)(buf[:32]))")
	g.L("\treturn 32, nil")
}

// genSmallIntEncoding generates optimized encoding for small integer types
func (g *Generator) genSmallIntEncoding(t ethabi.Type) {
	if t.Size%8 != 0 {
		panic(fmt.Sprintf("unsupported size %d for small integer decoding", t.Size))
	}

	// For small integers, we can use direct binary decoding without big.Int
	// Use the closest native integer type that fits
	size := nativeSize(t.Size)
	var nativeType string
	if t.T == ethabi.IntTy {
		nativeType = fmt.Sprintf("int%d", size)
	} else {
		nativeType = fmt.Sprintf("uint%d", size)
	}

	if t.T == ethabi.IntTy {
		// Signed: use the appropriate integer size and sign extend
		switch nativeType {
		case "int8":
			g.L("\tbuf[31] = byte(value)")
		case "int16":
			g.L("\tbinary.BigEndian.PutUint16(buf[30:32], uint16(value))")
		case "int32":
			g.L("\tbinary.BigEndian.PutUint32(buf[28:32], uint32(value))")
		case "int64":
			g.L("\tbinary.BigEndian.PutUint64(buf[24:32], uint64(value))")
		}

		g.L("\tif value < 0 {")
		g.L("\t\tcopy(buf, %sPaddingBytes%d)", g.StdPrefix, size)
		g.L("\t}")
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

// =============================================================================
// Packed Encoding Generators
// =============================================================================

// genPackedIntEncoding generates packed encoding for integer types (no padding)
func (g *Generator) genPackedIntEncoding(t ethabi.Type) {
	byteSize := t.Size / 8

	// Buffer length validation
	g.L("\tif len(buf) < %d {", byteSize)
	g.L("\t\treturn 0, io.ErrShortBuffer")
	g.L("\t}")

	if byteSize <= 8 {
		// For sizes <= 8 bytes, use native integer types
		switch byteSize {
		case 1:
			g.L("\tbuf[0] = byte(value)")
		case 2:
			g.L("\tbinary.BigEndian.PutUint16(buf[:2], uint16(value))")
		case 3:
			// 3 bytes: write as big-endian from uint32
			g.L("\tbuf[0] = byte(value >> 16)")
			g.L("\tbuf[1] = byte(value >> 8)")
			g.L("\tbuf[2] = byte(value)")
		case 4:
			g.L("\tbinary.BigEndian.PutUint32(buf[:4], uint32(value))")
		case 5:
			// 5 bytes: write as big-endian from uint64
			g.L("\tbuf[0] = byte(value >> 32)")
			g.L("\tbuf[1] = byte(value >> 24)")
			g.L("\tbuf[2] = byte(value >> 16)")
			g.L("\tbuf[3] = byte(value >> 8)")
			g.L("\tbuf[4] = byte(value)")
		case 6:
			// 6 bytes: write as big-endian from uint64
			g.L("\tbuf[0] = byte(value >> 40)")
			g.L("\tbuf[1] = byte(value >> 32)")
			g.L("\tbuf[2] = byte(value >> 24)")
			g.L("\tbuf[3] = byte(value >> 16)")
			g.L("\tbuf[4] = byte(value >> 8)")
			g.L("\tbuf[5] = byte(value)")
		case 7:
			// 7 bytes: write as big-endian from uint64
			g.L("\tbuf[0] = byte(value >> 48)")
			g.L("\tbuf[1] = byte(value >> 40)")
			g.L("\tbuf[2] = byte(value >> 32)")
			g.L("\tbuf[3] = byte(value >> 24)")
			g.L("\tbuf[4] = byte(value >> 16)")
			g.L("\tbuf[5] = byte(value >> 8)")
			g.L("\tbuf[6] = byte(value)")
		case 8:
			g.L("\tbinary.BigEndian.PutUint64(buf[:8], uint64(value))")
		}
	} else if t.T == ethabi.UintTy && g.Options.UseUint256 {
		g.L("\tvalue.WriteToArray32((*[32]byte)(buf[:32]))")
	} else {
		// For sizes > 8 bytes (big.Int), use EncodeBigInt
		if t.T == ethabi.IntTy {
			g.L("\tif err := %sEncodeBigInt(value, buf[:%d], true); err != nil {", g.StdPrefix, byteSize)
		} else {
			g.L("\tif err := %sEncodeBigInt(value, buf[:%d], false); err != nil {", g.StdPrefix, byteSize)
		}
		g.L("\t\treturn 0, err")
		g.L("\t}")
	}

	g.L("\treturn %d, nil", byteSize)
}

// genPackedAddressEncoding generates packed encoding for address (20 bytes, no padding)
func (g *Generator) genPackedAddressEncoding() {
	g.L("\tif len(buf) < 20 {")
	g.L("\t\treturn 0, io.ErrShortBuffer")
	g.L("\t}")
	g.L("\tcopy(buf[:20], value[:])")
	g.L("\treturn 20, nil")
}

// genPackedBoolEncoding generates packed encoding for bool (1 byte)
func (g *Generator) genPackedBoolEncoding() {
	g.L("\tif len(buf) < 1 {")
	g.L("\t\treturn 0, io.ErrShortBuffer")
	g.L("\t}")
	g.L("\tif value {")
	g.L("\t\tbuf[0] = 1")
	g.L("\t} else {")
	g.L("\t\tbuf[0] = 0")
	g.L("\t}")
	g.L("\treturn 1, nil")
}

// genPackedFixedBytesEncoding generates packed encoding for fixed bytes (no padding)
func (g *Generator) genPackedFixedBytesEncoding(t ethabi.Type) {
	g.L("\tif len(buf) < %d {", t.Size)
	g.L("\t\treturn 0, io.ErrShortBuffer")
	g.L("\t}")
	g.L("\tcopy(buf[:%d], value[:])", t.Size)
	g.L("\treturn %d, nil", t.Size)
}

// genPackedArrayEncoding generates packed encoding for fixed-size arrays
func (g *Generator) genPackedArrayEncoding(t ethabi.Type) {
	elemSize := GetPackedTypeSize(*t.Elem)
	totalSize := t.Size * elemSize

	g.L("\tif len(buf) < %d {", totalSize)
	g.L("\t\treturn 0, io.ErrShortBuffer")
	g.L("\t}")
	g.L("\t// Encode fixed-size array elements sequentially (no padding)")
	g.L("\tvar offset int")
	g.L("\tfor i := 0; i < %d; i++ {", t.Size)
	g.L("\t\tn, err := %s", g.genPackedEncodeCall(*t.Elem, "value[i]", "buf[offset:]"))
	g.L("\t\tif err != nil {")
	g.L("\t\t\treturn 0, err")
	g.L("\t\t}")
	g.L("\t\toffset += n")
	g.L("\t}")
	g.L("\treturn %d, nil", t.Size*elemSize)
}

// genPackedTupleEncoding generates packed encoding for tuple types
func (g *Generator) genPackedTupleEncoding(t ethabi.Type) {
	g.L("\t// Encode tuple fields sequentially (packed, no dynamic section)")
	if len(t.TupleElems) > 0 {
		g.L("\tvar (")
		g.L("\t\toffset int")
		g.L("\t\tn int")
		g.L("\t\terr error")
		g.L("\t)")
	}

	for i, elem := range t.TupleElems {
		fieldName := GoFieldName(t.TupleRawNames[i])
		if fieldName == "" {
			fieldName = fmt.Sprintf("Field%d", i+1)
		}

		ref := "value." + fieldName
		g.L("\t// Field %s: %s", fieldName, elem.String())
		g.L("\tn, err = %s", g.genPackedEncodeCall(*elem, ref, "buf[offset:]"))
		g.L("\tif err != nil {")
		g.L("\t\treturn 0, err")
		g.L("\t}")
		g.L("\toffset += n")
		g.L("")
	}

	g.L("\treturn offset, nil")
}
