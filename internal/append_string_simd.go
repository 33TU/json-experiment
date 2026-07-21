//go:build goexperiment.simd

package internal

import (
	"math/bits"
	"simd/archsimd"
	"slices"
	"unsafe"
)

const (
	lowerHex = "0123456789abcdef"

	byteOnes  = uint64(0x0101010101010101)
	byteHighs = uint64(0x8080808080808080)

	controlWord   = byteOnes * 0x20
	quoteWord     = byteOnes * uint64('"')
	backslashWord = byteOnes * uint64('\\')
	lessWord      = byteOnes * uint64('<')
	greaterWord   = byteOnes * uint64('>')
	ampersandWord = byteOnes * uint64('&')

	simdChunkSize = 16
	simdThreshold = 32
)

// AppendString appends the JSON representation of s to dst.
func AppendString(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+2)
	dst = append(dst, '"')

	data := unsafe.Pointer(unsafe.StringData(s))
	start := 0
	i := 0

	// SIMD: Process 16 bytes at a time using SIMD operations to detect special characters.
	if len(s) >= simdThreshold {
		control := archsimd.BroadcastUint8x16(0x20)
		quote := archsimd.BroadcastUint8x16('"')
		backslash := archsimd.BroadcastUint8x16('\\')

		src := unsafe.Slice((*byte)(data), len(s))

		for ; i+simdChunkSize <= len(src); i += simdChunkSize {
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
	}

	// SWAR: Process 8 bytes at a time using bitwise operations to detect special characters.
	for ; i+8 <= len(s); i += 8 {
		word := *(*uint64)(unsafe.Add(data, i))
		mask := hasByteLessThanWord(word, controlWord) |
			hasByteWord(word, quoteWord) |
			hasByteWord(word, backslashWord)

		if mask == 0 {
			continue
		}

		for j := i; j < i+8; j++ {
			c := s[j]
			if c >= 0x20 && c != '"' && c != '\\' {
				continue
			}

			dst = append(dst, s[start:j]...)
			dst = appendEscapedByte(dst, c)
			start = j + 1
		}
	}

	// Process the remaining bytes scalarly.
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

// AppendQuotedString appends the JSON representation of s encoded as a JSON string.
func AppendQuotedString(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+6)
	dst = append(dst, '"', '\\', '"')

	data := unsafe.Pointer(unsafe.StringData(s))
	start := 0
	i := 0

	// SIMD: Process 16 bytes at a time using SIMD operations to detect special characters.
	if len(s) >= simdThreshold {
		control := archsimd.BroadcastUint8x16(0x20)
		quote := archsimd.BroadcastUint8x16('"')
		backslash := archsimd.BroadcastUint8x16('\\')

		src := unsafe.Slice((*byte)(data), len(s))

		for ; i+simdChunkSize <= len(src); i += simdChunkSize {
			chunk := archsimd.LoadUint8x16Slice(src[i:])

			maskBits := chunk.
				Less(control).
				Or(chunk.Equal(quote)).
				Or(chunk.Equal(backslash)).
				ToBits()

			for maskBits != 0 {
				j := i + bits.TrailingZeros16(maskBits)

				dst = append(dst, s[start:j]...)
				dst = appendQuotedEscapedByte(dst, s[j])
				start = j + 1

				maskBits &= maskBits - 1
			}
		}
	}

	// SWAR: Process 8 bytes at a time using bitwise operations to detect special characters.
	for ; i+8 <= len(s); i += 8 {
		word := *(*uint64)(unsafe.Add(data, i))
		mask := hasByteLessThanWord(word, controlWord) |
			hasByteWord(word, quoteWord) |
			hasByteWord(word, backslashWord)

		if mask == 0 {
			continue
		}

		for j := i; j < i+8; j++ {
			c := s[j]
			if c >= 0x20 && c != '"' && c != '\\' {
				continue
			}

			dst = append(dst, s[start:j]...)
			dst = appendQuotedEscapedByte(dst, c)
			start = j + 1
		}
	}

	// Process the remaining bytes scalarly.
	for ; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' {
			continue
		}

		dst = append(dst, s[start:i]...)
		dst = appendQuotedEscapedByte(dst, c)
		start = i + 1
	}

	dst = append(dst, s[start:]...)
	return append(dst, '\\', '"', '"')
}

// AppendStringHTML appends the HTML-safe JSON representation of s to dst.
func AppendStringHTML(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+2)
	dst = append(dst, '"')

	data := unsafe.Pointer(unsafe.StringData(s))
	start := 0
	i := 0

	// SIMD: Process 16 bytes at a time using SIMD operations to detect special characters.
	if len(s) >= simdThreshold {
		control := archsimd.BroadcastUint8x16(0x20)
		quote := archsimd.BroadcastUint8x16('"')
		backslash := archsimd.BroadcastUint8x16('\\')
		less := archsimd.BroadcastUint8x16('<')
		greater := archsimd.BroadcastUint8x16('>')
		ampersand := archsimd.BroadcastUint8x16('&')

		src := unsafe.Slice((*byte)(data), len(s))

		for ; i+simdChunkSize <= len(src); i += simdChunkSize {
			chunk := archsimd.LoadUint8x16Slice(src[i:])

			maskBits := chunk.
				Less(control).
				Or(chunk.Equal(quote)).
				Or(chunk.Equal(backslash)).
				Or(chunk.Equal(less)).
				Or(chunk.Equal(greater)).
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
	}

	// SWAR: Process 8 bytes at a time using bitwise operations to detect special characters.
	for ; i+8 <= len(s); i += 8 {
		word := *(*uint64)(unsafe.Add(data, i))
		mask := hasByteLessThanWord(word, controlWord) |
			hasByteWord(word, quoteWord) |
			hasByteWord(word, backslashWord) |
			hasByteWord(word, lessWord) |
			hasByteWord(word, greaterWord) |
			hasByteWord(word, ampersandWord)

		if mask == 0 {
			continue
		}

		for j := i; j < i+8; j++ {
			c := s[j]
			if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
				continue
			}

			dst = append(dst, s[start:j]...)
			dst = appendEscapedByte(dst, c)
			start = j + 1
		}
	}

	// Process the remaining bytes scalarly.
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

// AppendQuotedStringHTML appends the HTML-safe JSON representation of s encoded as a JSON string.
func AppendQuotedStringHTML(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+6)
	dst = append(dst, '"', '\\', '"')

	data := unsafe.Pointer(unsafe.StringData(s))
	start := 0
	i := 0

	// SIMD: Process 16 bytes at a time using SIMD operations to detect special characters.
	if len(s) >= simdThreshold {
		control := archsimd.BroadcastUint8x16(0x20)
		quote := archsimd.BroadcastUint8x16('"')
		backslash := archsimd.BroadcastUint8x16('\\')
		less := archsimd.BroadcastUint8x16('<')
		greater := archsimd.BroadcastUint8x16('>')
		ampersand := archsimd.BroadcastUint8x16('&')

		src := unsafe.Slice((*byte)(data), len(s))

		for ; i+simdChunkSize <= len(src); i += simdChunkSize {
			chunk := archsimd.LoadUint8x16Slice(src[i:])

			maskBits := chunk.
				Less(control).
				Or(chunk.Equal(quote)).
				Or(chunk.Equal(backslash)).
				Or(chunk.Equal(less)).
				Or(chunk.Equal(greater)).
				Or(chunk.Equal(ampersand)).
				ToBits()

			for maskBits != 0 {
				j := i + bits.TrailingZeros16(maskBits)

				dst = append(dst, s[start:j]...)
				dst = appendQuotedEscapedByte(dst, s[j])
				start = j + 1

				maskBits &= maskBits - 1
			}
		}
	}

	// SWAR: Process 8 bytes at a time using bitwise operations to detect special characters.
	for ; i+8 <= len(s); i += 8 {
		word := *(*uint64)(unsafe.Add(data, i))
		mask := hasByteLessThanWord(word, controlWord) |
			hasByteWord(word, quoteWord) |
			hasByteWord(word, backslashWord) |
			hasByteWord(word, lessWord) |
			hasByteWord(word, greaterWord) |
			hasByteWord(word, ampersandWord)

		if mask == 0 {
			continue
		}

		for j := i; j < i+8; j++ {
			c := s[j]
			if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
				continue
			}

			dst = append(dst, s[start:j]...)
			dst = appendQuotedEscapedByte(dst, c)
			start = j + 1
		}
	}

	// Process the remaining bytes scalarly.
	for ; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
			continue
		}

		dst = append(dst, s[start:i]...)
		dst = appendQuotedEscapedByte(dst, c)
		start = i + 1
	}

	dst = append(dst, s[start:]...)
	return append(dst, '\\', '"', '"')
}

func appendEscapedByte(dst []byte, c byte) []byte {
	switch c {
	case '\\':
		return append(dst, '\\', '\\')
	case '"':
		return append(dst, '\\', '"')
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

func appendQuotedEscapedByte(dst []byte, c byte) []byte {
	switch c {
	case '"':
		return append(dst, '\\', '\\', '\\', '"')
	case '\\':
		return append(dst, '\\', '\\', '\\', '\\')
	case '\b':
		return append(dst, '\\', '\\', 'b')
	case '\f':
		return append(dst, '\\', '\\', 'f')
	case '\n':
		return append(dst, '\\', '\\', 'n')
	case '\r':
		return append(dst, '\\', '\\', 'r')
	case '\t':
		return append(dst, '\\', '\\', 't')
	default:
		return append(dst, '\\', '\\', 'u', '0', '0', lowerHex[c>>4], lowerHex[c&0x0f])
	}
}

func hasZeroByte(x uint64) uint64 {
	return (x - byteOnes) &^ x & byteHighs
}

func hasByteWord(x, word uint64) uint64 {
	return hasZeroByte(x ^ word)
}

func hasByteLessThanWord(x, word uint64) uint64 {
	return (x - word) &^ x & byteHighs
}
