package internal

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
)

// AppendBoolSlice appends the JSON representation of values to dst.
func AppendBoolSlice(dst []byte, values []bool) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendBool(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendIntSlice appends the JSON representation of values to dst.
func AppendIntSlice[T SignedInteger](dst []byte, values []T) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendInt(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendUintSlice appends the JSON representation of values to dst.
func AppendUintSlice[T UnsignedInteger](dst []byte, values []T) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendUint(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendFloat32Slice appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendFloat32Slice(dst []byte, values []float32) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "[]"...), nil
	}

	dst = append(dst, '[')
	for _, value := range values {
		var err error
		dst, err = AppendFloat32(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst, nil
}

// AppendFloat64Slice appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendFloat64Slice(dst []byte, values []float64) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "[]"...), nil
	}

	dst = append(dst, '[')
	for _, value := range values {
		var err error
		dst, err = AppendFloat64(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst, nil
}

// AppendStringSlice appends the JSON representation of values to dst.
func AppendStringSlice(dst []byte, values []string) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendString(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendStringSliceHTML appends the HTML-safe JSON representation of values to dst.
func AppendStringSliceHTML(dst []byte, values []string) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendStringHTML(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendByteSliceBase64 appends src as a standard Base64-encoded JSON string.
// A nil src is encoded as null.
func AppendByteSliceBase64(dst, src []byte) []byte {
	if src == nil {
		return AppendNull(dst)
	}

	dst = append(dst, '"')
	dst = base64.StdEncoding.AppendEncode(dst, src)
	return append(dst, '"')
}

// AppendByteSliceBase64URL appends src as a URL-safe Base64-encoded JSON string.
// A nil src is encoded as null.
func AppendByteSliceBase64URL(dst, src []byte) []byte {
	if src == nil {
		return AppendNull(dst)
	}

	dst = append(dst, '"')
	dst = base64.URLEncoding.AppendEncode(dst, src)
	return append(dst, '"')
}

// AppendByteSliceBase32 appends src as a standard Base32-encoded JSON string.
// A nil src is encoded as null.
func AppendByteSliceBase32(dst, src []byte) []byte {
	if src == nil {
		return AppendNull(dst)
	}

	dst = append(dst, '"')
	dst = base32.StdEncoding.AppendEncode(dst, src)
	return append(dst, '"')
}

// AppendByteSliceBase32Hex appends src as an extended-hex Base32-encoded JSON string.
// A nil src is encoded as null.
func AppendByteSliceBase32Hex(dst, src []byte) []byte {
	if src == nil {
		return AppendNull(dst)
	}

	dst = append(dst, '"')
	dst = base32.HexEncoding.AppendEncode(dst, src)
	return append(dst, '"')
}

// AppendByteSliceBase16 appends src as a lowercase Base16-encoded JSON string.
// A nil src is encoded as null.
func AppendByteSliceBase16(dst, src []byte) []byte {
	if src == nil {
		return AppendNull(dst)
	}

	dst = append(dst, '"')
	dst = hex.AppendEncode(dst, src)
	return append(dst, '"')
}

// AppendByteSliceHex appends src as a lowercase hexadecimal JSON string.
// It is an alias for AppendByteSliceBase16.
func AppendByteSliceHex(dst, src []byte) []byte {
	return AppendByteSliceBase16(dst, src)
}
