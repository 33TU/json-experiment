package internal

import "slices"

const lowerHex = "0123456789abcdef"

// AppendString appends the JSON representation of s to dst.
func AppendString(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+2)

	if firstEscape := findStringEscape(s); firstEscape >= 0 {
		return appendEscapedString(dst, s, firstEscape)
	}

	dst = append(dst, '"')
	dst = append(dst, s...)
	return append(dst, '"')
}

// AppendStringHTML appends the HTML-safe JSON representation of s to dst.
func AppendStringHTML(dst []byte, s string) []byte {
	dst = slices.Grow(dst, len(s)+2)

	if firstEscape := findStringEscapeHTML(s); firstEscape >= 0 {
		return appendEscapedStringHTML(dst, s, firstEscape)
	}

	dst = append(dst, '"')
	dst = append(dst, s...)
	return append(dst, '"')
}

func appendEscapedString(dst []byte, s string, firstEscape int) []byte {
	dst = append(dst, '"')
	start := 0

	for i := firstEscape; i < len(s); i++ {
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

func appendEscapedStringHTML(dst []byte, s string, firstEscape int) []byte {
	dst = append(dst, '"')
	start := 0

	for i := firstEscape; i < len(s); i++ {
		c := s[i]
		if c >= 0x20 && c != '"' && c != '\\' && c != '<' && c != '>' && c != '&' {
			continue
		}

		dst = append(dst, s[start:i]...)
		switch c {
		case '<', '>', '&':
			dst = append(dst, '\\', 'u', '0', '0', lowerHex[c>>4], lowerHex[c&0x0f])
		default:
			dst = appendEscapedByte(dst, c)
		}
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
