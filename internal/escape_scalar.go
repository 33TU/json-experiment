//go:build !goexperiment.simd

package internal

// findStringEscape returns the byte index of the first character requiring JSON escaping, or -1 if no escaping is required.
func findStringEscape(s string) int {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c == '"' || c == '\\' {
			return i
		}
	}
	return -1
}

// findStringEscapeHTML returns the byte index of the first character requiring JSON or HTML-safe escaping, or -1 if no escaping is required.
func findStringEscapeHTML(s string) int {
	for i := 0; i < len(s); i++ {
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
