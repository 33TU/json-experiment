package internal

import (
	"slices"
	"sync"
	"unicode/utf8"
	"unsafe"
)

type utf8ScratchBuffer struct {
	b []byte
}

var utf8ScratchPool = sync.Pool{
	New: func() any {
		return &utf8ScratchBuffer{
			b: make([]byte, 0, 1024),
		}
	},
}

// ReplaceInvalidUTF8 replaces invalid UTF-8 in buf[start:] with the Unicode
// replacement character. Valid input is returned unchanged without using the scratch pool.
func ReplaceInvalidUTF8(buf []byte, start int) []byte {
	corrected, scratch := replaceInvalidUTF8(buf[start:])
	if scratch == nil {
		return buf
	}

	buf = append(buf[:start], corrected...)

	scratch.b = corrected[:0]
	utf8ScratchPool.Put(scratch)
	return buf
}

func acquireUTF8Scratch(size int) *utf8ScratchBuffer {
	scratch := utf8ScratchPool.Get().(*utf8ScratchBuffer)
	scratch.b = slices.Grow(scratch.b[:0], size)
	return scratch
}

func replaceInvalidUTF8(src []byte) ([]byte, *utf8ScratchBuffer) {
	const (
		byteHighs = 0x8080808080808080
		runeError = string(utf8.RuneError)
	)

	data := unsafe.Pointer(unsafe.SliceData(src))
	start := 0
	i := 0

	var scratch *utf8ScratchBuffer
	remainderKnownInvalid := false

	for i < len(src) {
		b0 := src[i]
		if b0 < utf8.RuneSelf {
			i++

			// Skip subsequent ASCII thirty-two bytes at a time.
			for i+32 <= len(src) {
				word := *(*uint64)(unsafe.Add(data, i)) |
					*(*uint64)(unsafe.Add(data, i+8)) |
					*(*uint64)(unsafe.Add(data, i+16)) |
					*(*uint64)(unsafe.Add(data, i+24))
				if word&byteHighs != 0 {
					break
				}
				i += 32
			}

			for i+16 <= len(src) {
				word := *(*uint64)(unsafe.Add(data, i)) |
					*(*uint64)(unsafe.Add(data, i+8))
				if word&byteHighs != 0 {
					break
				}
				i += 16
			}

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
			continue
		}

		if !remainderKnownInvalid {
			if validUTF8(src[i:]) {
				return src, nil
			}
			remainderKnownInvalid = true
		}

		size := 0
		switch {
		case b0 >= 0xc2 && b0 < 0xe0:
			if i+1 < len(src) && isUTF8Continuation(src[i+1]) {
				size = 2
			}
		case b0 >= 0xe0 && b0 < 0xf0:
			if i+2 < len(src) {
				b1 := src[i+1]
				if isUTF8Continuation(b1) && isUTF8Continuation(src[i+2]) &&
					(b0 != 0xe0 || b1 >= 0xa0) &&
					(b0 != 0xed || b1 < 0xa0) {
					size = 3
				}
			}
		case b0 >= 0xf0 && b0 < 0xf5:
			if i+3 < len(src) {
				b1 := src[i+1]
				if isUTF8Continuation(b1) && isUTF8Continuation(src[i+2]) && isUTF8Continuation(src[i+3]) &&
					(b0 != 0xf0 || b1 >= 0x90) &&
					(b0 != 0xf4 || b1 < 0x90) {
					size = 4
				}
			}
		}
		if size != 0 {
			i += size
			continue
		}

		if scratch == nil {
			scratch = acquireUTF8Scratch(len(src) + 2)
		}
		scratch.b = append(scratch.b, src[start:i]...)
		scratch.b = append(scratch.b, runeError...)
		i++
		start = i
	}

	if scratch == nil {
		return src, nil
	}
	scratch.b = append(scratch.b, src[start:]...)
	return scratch.b, scratch
}

func isUTF8Continuation(b byte) bool {
	return b >= 0x80 && b <= 0xbf
}
