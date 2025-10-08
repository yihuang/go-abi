package abi

import (
	"math/big"
)

// Uint8 implementation
func (u *Uint8) TypeName() string   { return TypeUint8 }
func (u *Uint8) IsDynamic() bool    { return false }
func (u *Uint8) StaticSize() uint   { return WordSize }
func (u *Uint8) Encode(value interface{}) ([]byte, error) {
	val, ok := value.(*big.Int)
	if !ok {
		return nil, ErrInvalidValue
	}
	return encodeUint(val, 8), nil
}
func (u *Uint8) Decode(data []byte, value interface{}) (uint, error) {
	val, err := decodeUint(data, 8)
	if err != nil {
		return 0, err
	}

	if ptr, ok := value.(**big.Int); ok {
		*ptr = val
	} else if ptr, ok := value.(*big.Int); ok {
		ptr.Set(val)
	} else {
		return 0, ErrInvalidValue
	}

	return WordSize, nil
}

// Uint256 implementation
func (u *Uint256) TypeName() string { return TypeUint256 }
func (u *Uint256) IsDynamic() bool  { return false }
func (u *Uint256) StaticSize() uint { return WordSize }
func (u *Uint256) Encode(value interface{}) ([]byte, error) {
	val, ok := value.(*big.Int)
	if !ok {
		return nil, ErrInvalidValue
	}
	return encodeUint(val, 256), nil
}
func (u *Uint256) Decode(data []byte, value interface{}) (uint, error) {
	val, err := decodeUint(data, 256)
	if err != nil {
		return 0, err
	}

	if ptr, ok := value.(**big.Int); ok {
		*ptr = val
	} else if ptr, ok := value.(*big.Int); ok {
		ptr.Set(val)
	} else {
		return 0, ErrInvalidValue
	}

	return WordSize, nil
}

// Int256 implementation
func (i *Int256) TypeName() string { return TypeInt256 }
func (i *Int256) IsDynamic() bool  { return false }
func (i *Int256) StaticSize() uint { return WordSize }
func (i *Int256) Encode(value interface{}) ([]byte, error) {
	val, ok := value.(*big.Int)
	if !ok {
		return nil, ErrInvalidValue
	}
	return encodeInt(val, 256), nil
}
func (i *Int256) Decode(data []byte, value interface{}) (uint, error) {
	val, err := decodeInt(data, 256)
	if err != nil {
		return 0, err
	}

	if ptr, ok := value.(**big.Int); ok {
		*ptr = val
	} else if ptr, ok := value.(*big.Int); ok {
		ptr.Set(val)
	} else {
		return 0, ErrInvalidValue
	}

	return WordSize, nil
}

// Bool implementation
func (b *Bool) TypeName() string   { return TypeBool }
func (b *Bool) IsDynamic() bool    { return false }
func (b *Bool) StaticSize() uint   { return WordSize }
func (b *Bool) Encode(value interface{}) ([]byte, error) {
	val, ok := value.(bool)
	if !ok {
		return nil, ErrInvalidValue
	}
	return encodeBool(val), nil
}
func (b *Bool) Decode(data []byte, value interface{}) (uint, error) {
	val, err := decodeBool(data)
	if err != nil {
		return 0, err
	}

	if ptr, ok := value.(*bool); ok {
		*ptr = val
	} else {
		return 0, ErrInvalidValue
	}

	return WordSize, nil
}

// Address implementation
func (a *Address) TypeName() string { return TypeAddress }
func (a *Address) IsDynamic() bool  { return false }
func (a *Address) StaticSize() uint { return WordSize }
func (a *Address) Encode(value interface{}) ([]byte, error) {
	addr, ok := value.([20]byte)
	if !ok {
		return nil, ErrInvalidValue
	}
	return encodeAddress(addr), nil
}
func (a *Address) Decode(data []byte, value interface{}) (uint, error) {
	addr, err := decodeAddress(data)
	if err != nil {
		return 0, err
	}

	if ptr, ok := value.(*[20]byte); ok {
		*ptr = addr
	} else {
		return 0, ErrInvalidValue
	}

	return WordSize, nil
}

// Bytes implementation
func (b *Bytes) TypeName() string { return TypeBytes }
func (b *Bytes) IsDynamic() bool  { return true }
func (b *Bytes) StaticSize() uint { return 0 }
func (b *Bytes) Encode(value interface{}) ([]byte, error) {
	bytes, ok := value.([]byte)
	if !ok {
		return nil, ErrInvalidValue
	}
	return encodeBytes(bytes), nil
}
func (b *Bytes) Decode(data []byte, value interface{}) (uint, error) {
	bytes, bytesRead, err := decodeBytes(data)
	if err != nil {
		return 0, err
	}

	if ptr, ok := value.(*[]byte); ok {
		*ptr = bytes
	} else {
		return 0, ErrInvalidValue
	}

	return bytesRead, nil
}

// String implementation
func (s *String) TypeName() string { return TypeString }
func (s *String) IsDynamic() bool  { return true }
func (s *String) StaticSize() uint { return 0 }
func (s *String) Encode(value interface{}) ([]byte, error) {
	str, ok := value.(string)
	if !ok {
		return nil, ErrInvalidValue
	}
	return encodeBytes([]byte(str)), nil
}
func (s *String) Decode(data []byte, value interface{}) (uint, error) {
	bytes, bytesRead, err := decodeBytes(data)
	if err != nil {
		return 0, err
	}

	if ptr, ok := value.(*string); ok {
		*ptr = string(bytes)
	} else {
		return 0, ErrInvalidValue
	}

	return bytesRead, nil
}