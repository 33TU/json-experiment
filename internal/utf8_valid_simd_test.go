//go:build goexperiment.simd

package internal

import (
	"bytes"
	"strconv"
	"testing"
	"unicode/utf8"
)

var validUTF8Result bool

func BenchmarkValidUTF8(b *testing.B) {
	for _, size := range []int{32, 64, 96, 128, 160, 256, 320, 384, 448, 512, 1024, 4096} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			src := bytes.Repeat([]byte("世"), size/3)
			src = append(src, bytes.Repeat([]byte{'a'}, size-len(src))...)
			var result bool
			b.SetBytes(int64(len(src)))
			for b.Loop() {
				result = validUTF8(src)
			}
			validUTF8Result = result
		})
	}
}

func TestValidUTF8SIMDMatchesStandard(t *testing.T) {
	state := uint64(0x9e3779b97f4a7c15)
	for length := 128; length <= 512; length++ {
		for range 32 {
			src := make([]byte, length)
			for i := range src {
				state ^= state << 7
				state ^= state >> 9
				state ^= state << 8
				src[i] = byte(state)
			}

			if got, want := validUTF8(src), utf8.Valid(src); got != want {
				t.Fatalf("validUTF8(%x) = %v, want %v", src, got, want)
			}
		}
	}

	runes := []byte("a¢€𐍈")
	for prefix := range 32 {
		src := append(bytes.Repeat([]byte{'a'}, prefix), bytes.Repeat(runes, 40)...)
		if !validUTF8(src) {
			t.Fatalf("validUTF8 rejected valid input with prefix length %d", prefix)
		}

		for i := range src {
			corrupt := append([]byte(nil), src...)
			corrupt[i] = 0xff
			if validUTF8(corrupt) {
				t.Fatalf("validUTF8 accepted corruption at prefix %d, offset %d", prefix, i)
			}
		}
	}
}
