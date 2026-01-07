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
	} else if t.T == ethabi.UintTy && g.Options.UseUint256 {
		g.genUint256Decoding()
	} else {
		g.genBigIntDecoding(t)
	}
}

// genUint256Decoding generates decoding for holiman/uint256.Int types
func (g *Generator) genUint256Decoding() {
	g.L("\tif len(data) < 32 {")
	g.L("\t\treturn nil, 0, io.ErrUnexpectedEOF")
	g.L("\t}")
	g.L("\tresult := new(uint256.Int)")
	g.L("\tresult.SetBytes32(data[:32])")
	g.L("\treturn result, 32, nil")
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

// =============================================================================
// Packed Decoding Generators
// =============================================================================

// genPackedIntDecoding generates packed decoding for integer types
func (g *Generator) genPackedIntDecoding(t ethabi.Type) {
	byteSize := t.Size / 8

	// Use appropriate zero value for error returns
	zeroValue := "0"
	if byteSize > 8 {
		zeroValue = "nil"
	}

	g.L("\tif len(data) < %d {", byteSize)
	g.L("\t\treturn %s, 0, io.ErrUnexpectedEOF", zeroValue)
	g.L("\t}")

	if byteSize <= 8 {
		// For sizes <= 8 bytes, use native integer types
		switch byteSize {
		case 1:
			if t.T == ethabi.IntTy {
				g.L("\treturn int8(data[0]), 1, nil")
			} else {
				g.L("\treturn data[0], 1, nil")
			}
		case 2:
			if t.T == ethabi.IntTy {
				g.L("\treturn int16(binary.BigEndian.Uint16(data[:2])), 2, nil")
			} else {
				g.L("\treturn binary.BigEndian.Uint16(data[:2]), 2, nil")
			}
		case 3:
			// 3 bytes: read as big-endian into uint32/int32
			if t.T == ethabi.IntTy {
				g.L("\tv := int32(data[0])<<16 | int32(data[1])<<8 | int32(data[2])")
				g.L("\tif data[0]&0x80 != 0 { v = v | (^int32(0) << 24) }") // sign extend upper 8 bits
				g.L("\treturn v, 3, nil")
			} else {
				g.L("\treturn uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2]), 3, nil")
			}
		case 4:
			if t.T == ethabi.IntTy {
				g.L("\treturn int32(binary.BigEndian.Uint32(data[:4])), 4, nil")
			} else {
				g.L("\treturn binary.BigEndian.Uint32(data[:4]), 4, nil")
			}
		case 5:
			// 5 bytes: read as big-endian into uint64/int64
			if t.T == ethabi.IntTy {
				g.L("\tv := int64(data[0])<<32 | int64(data[1])<<24 | int64(data[2])<<16 | int64(data[3])<<8 | int64(data[4])")
				g.L("\tif data[0]&0x80 != 0 { v = v | (^int64(0) << 40) }") // sign extend upper 24 bits
				g.L("\treturn v, 5, nil")
			} else {
				g.L("\treturn uint64(data[0])<<32 | uint64(data[1])<<24 | uint64(data[2])<<16 | uint64(data[3])<<8 | uint64(data[4]), 5, nil")
			}
		case 6:
			// 6 bytes: read as big-endian into uint64/int64
			if t.T == ethabi.IntTy {
				g.L("\tv := int64(data[0])<<40 | int64(data[1])<<32 | int64(data[2])<<24 | int64(data[3])<<16 | int64(data[4])<<8 | int64(data[5])")
				g.L("\tif data[0]&0x80 != 0 { v = v | (^int64(0) << 48) }") // sign extend upper 16 bits
				g.L("\treturn v, 6, nil")
			} else {
				g.L("\treturn uint64(data[0])<<40 | uint64(data[1])<<32 | uint64(data[2])<<24 | uint64(data[3])<<16 | uint64(data[4])<<8 | uint64(data[5]), 6, nil")
			}
		case 7:
			// 7 bytes: read as big-endian into uint64/int64
			if t.T == ethabi.IntTy {
				g.L("\tv := int64(data[0])<<48 | int64(data[1])<<40 | int64(data[2])<<32 | int64(data[3])<<24 | int64(data[4])<<16 | int64(data[5])<<8 | int64(data[6])")
				g.L("\tif data[0]&0x80 != 0 { v = v | (^int64(0) << 56) }") // sign extend upper 8 bits
				g.L("\treturn v, 7, nil")
			} else {
				g.L("\treturn uint64(data[0])<<48 | uint64(data[1])<<40 | uint64(data[2])<<32 | uint64(data[3])<<24 | uint64(data[4])<<16 | uint64(data[5])<<8 | uint64(data[6]), 7, nil")
			}
		case 8:
			if t.T == ethabi.IntTy {
				g.L("\treturn int64(binary.BigEndian.Uint64(data[:8])), 8, nil")
			} else {
				g.L("\treturn binary.BigEndian.Uint64(data[:8]), 8, nil")
			}
		}
	} else {
		// For sizes > 8 bytes
		if t.T == ethabi.UintTy && g.Options.UseUint256 {
			// Use uint256.Int for large unsigned integers when enabled
			g.genPackedLargeUintDecoding(t)
			return
		}
		// Use big.Int
		if t.T == ethabi.IntTy {
			g.L("\tresult, err := %sDecodeBigInt(data[:%d], true)", g.StdPrefix, byteSize)
			g.L("\tif err != nil {")
			g.L("\t\treturn nil, 0, err")
			g.L("\t}")
			g.L("\treturn result, %d, nil", byteSize)
		} else {
			g.L("\tresult := new(big.Int).SetBytes(data[:%d])", byteSize)
			g.L("\treturn result, %d, nil", byteSize)
		}
	}
}

// genPackedUint256Decoding generates packed decoding for holiman/uint256.Int types
func (g *Generator) genPackedUint256Decoding() {
	g.L("\tresult := new(uint256.Int)")
	g.L("\tresult.SetBytes32(data[:32])")
	g.L("\treturn result, 32, nil")
}

// genPackedLargeUintDecoding generates packed decoding for large unsigned integers using uint256.Int
func (g *Generator) genPackedLargeUintDecoding(t ethabi.Type) {
	byteSize := t.Size / 8
	g.L("\tresult := new(uint256.Int)")
	if byteSize == 32 {
		g.L("\tresult.SetBytes32(data[:32])")
	} else {
		g.L("\tresult.SetBytes(data[:%d])", byteSize)
	}
	g.L("\treturn result, %d, nil", byteSize)
}

// genPackedAddressDecoding generates packed decoding for address (20 bytes)
func (g *Generator) genPackedAddressDecoding() {
	g.L("\tif len(data) < 20 {")
	g.L("\t\treturn common.Address{}, 0, io.ErrUnexpectedEOF")
	g.L("\t}")
	g.L("\tvar result common.Address")
	g.L("\tcopy(result[:], data[:20])")
	g.L("\treturn result, 20, nil")
}

// genPackedBoolDecoding generates packed decoding for bool (1 byte)
func (g *Generator) genPackedBoolDecoding() {
	g.L("\tswitch data[0] {")
	g.L("\tcase 0x00:")
	g.L("\t\treturn false, 1, nil")
	g.L("\tcase 0x01:")
	g.L("\t\treturn true, 1, nil")
	g.L("\tdefault:")
	g.L("\t\treturn false, 0, %sErrDirtyPadding", g.StdPrefix)
	g.L("\t}")
}

// genPackedFixedBytesDecoding generates packed decoding for fixed bytes
func (g *Generator) genPackedFixedBytesDecoding(t ethabi.Type) {
	g.L("\tif len(data) < %d {", t.Size)
	g.L("\t\treturn [%d]byte{}, 0, io.ErrUnexpectedEOF", t.Size)
	g.L("\t}")
	g.L("\tvar result [%d]byte", t.Size)
	g.L("\tcopy(result[:], data[:%d])", t.Size)
	g.L("\treturn result, %d, nil", t.Size)
}

// genPackedArrayDecoding generates packed decoding for fixed-size arrays
func (g *Generator) genPackedArrayDecoding(t ethabi.Type) {
	goType := g.abiTypeToGoType(*t.Elem)
	elemSize := GetPackedTypeSize(*t.Elem)
	totalSize := t.Size * elemSize

	g.L("\tif len(data) < %d {", totalSize)
	g.L("\t\treturn [%d]%s{}, 0, io.ErrUnexpectedEOF", t.Size, goType)
	g.L("\t}")

	g.L("\tvar (")
	g.L("\t\tresult [%d]%s", t.Size, goType)
	g.L("\t\toffset int")
	g.L("\t\tn int")
	g.L("\t\terr error")
	g.L("\t)")

	g.L("\tfor i := 0; i < %d; i++ {", t.Size)
	if t.Elem.T == ethabi.TupleTy {
		g.L("\t\tn, err = result[i].PackedDecode(data[offset:])")
	} else {
		g.L("\t\tresult[i], n, err = %s", g.genPackedDecodeCall(*t.Elem, "data[offset:]"))
	}
	g.L("\t\tif err != nil {")
	g.L("\t\t\treturn result, 0, err")
	g.L("\t\t}")
	g.L("\t\toffset += n")
	g.L("\t}")
	g.L("\treturn result, %d, nil", totalSize)
}
