//go:build goexperiment.simd

package internal

import (
	"math/bits"
	"simd/archsimd"
	"slices"
	"unsafe"
)

const lowerHex = "0123456789abcdef"
const chunkSize = 16
const simdThreshold = 32

// AppendString appends the JSON representation of s to dst.
func AppendString(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+2)

	if len(s) < simdThreshold {
		return appendStringScalar(dst, s)
	}

	src := unsafe.Slice(unsafe.StringData(s), len(s))

	control := archsimd.BroadcastUint8x16(0x20)
	quote := archsimd.BroadcastUint8x16('"')
	backslash := archsimd.BroadcastUint8x16('\\')

	dst = append(dst, '"')

	start := 0
	i := 0

	// Process using SIMD for chunks of 16 bytes at a time.
	for ; i+chunkSize <= len(src); i += chunkSize {
		chunk := archsimd.LoadUint8x16Slice(src[i:])

		maskBits := chunk.
			Less(control).
			Or(chunk.Equal(quote)).
			Or(chunk.Equal(backslash)).
			ToBits()

		for maskBits != 0 {
			j := i + bits.TrailingZeros16(maskBits)

			dst = append(dst, s[start:j]...)
			dst = appendEscapedByte(dst, s[j])
			start = j + 1

			maskBits &= maskBits - 1
		}
	}

	// Process any remaining bytes that don't fit into a full SIMD chunk.
	for ; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' {
			continue
		}

		dst = append(dst, s[start:i]...)
		dst = appendEscapedByte(dst, c)
		start = i + 1
	}

	dst = append(dst, s[start:]...)
	return append(dst, '"')
}

// AppendStringHTML appends the HTML-safe JSON representation of s to dst.
func AppendStringHTML(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+2)

	if len(s) < simdThreshold {
		return appendStringHTMLScalar(dst, s)
	}

	src := unsafe.Slice(unsafe.StringData(s), len(s))

	control := archsimd.BroadcastUint8x16(0x20)
	quote := archsimd.BroadcastUint8x16('"')
	backslash := archsimd.BroadcastUint8x16('\\')
	lessThan := archsimd.BroadcastUint8x16('<')
	greaterThan := archsimd.BroadcastUint8x16('>')
	ampersand := archsimd.BroadcastUint8x16('&')

	dst = append(dst, '"')

	start := 0
	i := 0

	// Process using SIMD for chunks of 16 bytes at a time.
	for ; i+chunkSize <= len(src); i += chunkSize {
		chunk := archsimd.LoadUint8x16Slice(src[i:])

		maskBits := chunk.
			Less(control).
			Or(chunk.Equal(quote)).
			Or(chunk.Equal(backslash)).
			Or(chunk.Equal(lessThan)).
			Or(chunk.Equal(greaterThan)).
			Or(chunk.Equal(ampersand)).
			ToBits()

		for maskBits != 0 {
			j := i + bits.TrailingZeros16(maskBits)

			dst = append(dst, s[start:j]...)
			dst = appendEscapedByte(dst, s[j])
			start = j + 1

			maskBits &= maskBits - 1
		}
	}

	// Process any remaining bytes that don't fit into a full SIMD chunk.
	for ; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
			continue
		}

		dst = append(dst, s[start:i]...)
		dst = appendEscapedByte(dst, c)
		start = i + 1
	}

	dst = append(dst, s[start:]...)
	return append(dst, '"')
}

func appendStringScalar(dst []byte, s string) []byte {
	dst = append(dst, '"')

	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' {
			continue
		}

		dst = append(dst, s[start:i]...)
		dst = appendEscapedByte(dst, c)
		start = i + 1
	}

	dst = append(dst, s[start:]...)
	return append(dst, '"')
}

func appendStringHTMLScalar(dst []byte, s string) []byte {
	dst = append(dst, '"')

	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
			continue
		}

		dst = append(dst, s[start:i]...)
		dst = appendEscapedByte(dst, c)
		start = i + 1
	}

	dst = append(dst, s[start:]...)
	return append(dst, '"')
}

func appendEscapedByte(dst []byte, c byte) []byte {
	switch c {
	case '\\', '"':
		return append(dst, '\\', c)
	case '\b':
		return append(dst, '\\', 'b')
	case '\f':
		return append(dst, '\\', 'f')
	case '\n':
		return append(dst, '\\', 'n')
	case '\r':
		return append(dst, '\\', 'r')
	case '\t':
		return append(dst, '\\', 't')
	default:
		return append(dst, '\\', 'u', '0', '0', lowerHex[c>>4], lowerHex[c&0x0f])
	}
}
