//go:build goexperiment.simd

package internal

import (
	"simd/archsimd"
	"unicode/utf8"
	"unsafe"
)

// AppendValidUTF8 appends src to dst, replacing each invalid UTF-8 byte with \ufffd.
func AppendValidUTF8(dst, src []byte) []byte {
	const (
		byteHighs     = 0x8080808080808080
		chunkSize     = 16
		simdThreshold = 32
	)

	data := unsafe.Pointer(unsafe.SliceData(src))
	start := 0
	i := 0

	var nonASCII archsimd.Uint8x16
	if len(src) >= simdThreshold {
		nonASCII = archsimd.BroadcastUint8x16(utf8.RuneSelf)
	}

	for i < len(src) {
		// Skip ASCII sixteen bytes at a time.
		for len(src) >= simdThreshold && i+chunkSize <= len(src) {
			chunk := archsimd.LoadUint8x16Slice(src[i:])
			if chunk.GreaterEqual(nonASCII).ToBits() != 0 {
				break
			}
			i += chunkSize
		}

		// Process the remaining full words with SWAR.
		for i+8 <= len(src) {
			word := *(*uint64)(unsafe.Add(data, i))
			if word&byteHighs != 0 {
				break
			}
			i += 8
		}

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
