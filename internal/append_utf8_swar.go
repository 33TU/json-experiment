//go:build !goexperiment.simd

package internal

import (
	"unicode/utf8"
	"unsafe"
)

// AppendValidUTF8 appends src to dst, replacing each invalid UTF-8 byte with \ufffd.
func AppendValidUTF8(dst, src []byte) []byte {
	const byteHighs = 0x8080808080808080

	data := unsafe.Pointer(unsafe.SliceData(src))
	start := 0
	i := 0

	for i < len(src) {
		// Skip ASCII eight bytes at a time.
		for i+8 <= len(src) {
			word := *(*uint64)(unsafe.Add(data, i))
			if word&byteHighs != 0 {
				break
			}
			i += 8
		}

		// Process the remaining bytes one by one.
		for i < len(src) && src[i] < utf8.RuneSelf {
			i++
		}
		if i == len(src) {
			break
		}

		_, size := utf8.DecodeRune(src[i:])
		if size != 1 {
			i += size
			continue
		}

		dst = append(dst, src[start:i]...)
		dst = append(dst, `\ufffd`...)
		i++
		start = i
	}

	return append(dst, src[start:]...)
}
