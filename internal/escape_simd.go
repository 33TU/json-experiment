//go:build goexperiment.simd

package internal

import (
	"math/bits"
	"simd/archsimd"
	"unsafe"
)

// findStringEscapeFrom returns the byte index of the first character requiring JSON escaping at or after start, or -1 if no escaping is required.
func findStringEscapeFrom(s string, start int) int {
	for i := start; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c == '"' || c == '\\' {
			return i
		}
	}

	return -1
}

// findStringEscape returns the byte index of the first character requiring JSON escaping, or -1 if no escaping is required.
func findStringEscape(s string) int {
	const chunkSize = 16
	const simdThreshold = 32

	if len(s) < simdThreshold {
		return findStringEscapeFrom(s, 0)
	}

	src := unsafe.Slice(unsafe.StringData(s), len(s))

	control := archsimd.BroadcastUint8x16(0x20)
	quote := archsimd.BroadcastUint8x16('"')
	backslash := archsimd.BroadcastUint8x16('\\')

	i := 0
	for ; i+chunkSize <= len(src); i += chunkSize {
		chunk := archsimd.LoadUint8x16Slice(src[i:])

		mask := chunk.
			Less(control).
			Or(chunk.Equal(quote)).
			Or(chunk.Equal(backslash))

		if maskBits := mask.ToBits(); maskBits != 0 {
			return i + bits.TrailingZeros16(maskBits)
		}
	}

	return findStringEscapeFrom(s, i)
}

// findStringEscapeHTMLFrom returns the byte index of the first character requiring JSON or HTML-safe escaping at or after start, or -1.
func findStringEscapeHTMLFrom(s string, start int) int {
	for i := start; i < len(s); i++ {
		switch c := s[i]; c {
		case '"', '\\', '<', '>', '&':
			return i
		default:
			if c < 0x20 {
				return i
			}
		}
	}

	return -1
}

// findStringEscapeHTML returns the byte index of the first character requiring JSON or HTML-safe escaping, or -1 if no escaping is required.
func findStringEscapeHTML(s string) int {
	const chunkSize = 16
	const simdThreshold = 32

	if len(s) < simdThreshold {
		return findStringEscapeHTMLFrom(s, 0)
	}

	src := unsafe.Slice(unsafe.StringData(s), len(s))

	control := archsimd.BroadcastUint8x16(0x20)
	quote := archsimd.BroadcastUint8x16('"')
	backslash := archsimd.BroadcastUint8x16('\\')
	lessThan := archsimd.BroadcastUint8x16('<')
	greaterThan := archsimd.BroadcastUint8x16('>')
	ampersand := archsimd.BroadcastUint8x16('&')

	i := 0
	for ; i+chunkSize <= len(src); i += chunkSize {
		chunk := archsimd.LoadUint8x16Slice(src[i:])

		mask := chunk.
			Less(control).
			Or(chunk.Equal(quote)).
			Or(chunk.Equal(backslash)).
			Or(chunk.Equal(lessThan)).
			Or(chunk.Equal(greaterThan)).
			Or(chunk.Equal(ampersand))

		if maskBits := mask.ToBits(); maskBits != 0 {
			return i + bits.TrailingZeros16(maskBits)
		}
	}

	return findStringEscapeHTMLFrom(s, i)
}
