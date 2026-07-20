//go:build !goexperiment.simd

package internal

import "slices"

const lowerHex = "0123456789abcdef"

// AppendString appends the JSON representation of s to dst.
func AppendString(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+2)
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

// AppendStringHTML appends the HTML-safe JSON representation of s to dst.
func AppendStringHTML(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+2)
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
