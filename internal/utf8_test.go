package internal_test

import (
	"bytes"
	"testing"
	"unicode/utf8"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

func TestReplaceInvalidUTF8(t *testing.T) {
	tests := []struct {
		name string
		src  []byte
		want string
	}{
		{name: "empty"},
		{name: "ASCII", src: []byte("hello"), want: "hello"},
		{name: "Unicode", src: []byte("Hello, 世界"), want: "Hello, 世界"},
		{name: "replacement rune", src: []byte("a\uFFFDb"), want: "a\uFFFDb"},
		{name: "invalid byte", src: []byte{'a', 0xff, 'b'}, want: "a\uFFFDb"},
		{name: "malformed sequence", src: []byte{0xe2, 0x28, 0xa1}, want: "\uFFFD(\uFFFD"},
		{name: "unexpected continuations", src: []byte{0x80, 0x80, 0x80}, want: "\uFFFD\uFFFD\uFFFD"},
		{name: "overlong two byte", src: []byte{0xc0, 0x80}, want: "\uFFFD\uFFFD"},
		{name: "overlong three byte", src: []byte{0xe0, 0x80, 0x80}, want: "\uFFFD\uFFFD\uFFFD"},
		{name: "surrogate", src: []byte{0xed, 0xa0, 0x80}, want: "\uFFFD\uFFFD\uFFFD"},
		{name: "overlong four byte", src: []byte{0xf0, 0x80, 0x80, 0x80}, want: "\uFFFD\uFFFD\uFFFD\uFFFD"},
		{name: "out of range", src: []byte{0xf4, 0x90, 0x80, 0x80}, want: "\uFFFD\uFFFD\uFFFD\uFFFD"},
		{name: "valid boundaries", src: []byte{0xc2, 0x80, 0xe0, 0xa0, 0x80, 0xed, 0x9f, 0xbf, 0xf0, 0x90, 0x80, 0x80, 0xf4, 0x8f, 0xbf, 0xbf}, want: "\u0080\u0800\uD7FF\U00010000\U0010FFFF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := append([]byte("prefix:"), tt.src...)
			got := internal.ReplaceInvalidUTF8(buf, len("prefix:"))
			want := append([]byte("prefix:"), tt.want...)

			if !bytes.Equal(got, want) {
				t.Fatalf("ReplaceInvalidUTF8() = %q, want %q", got, want)
			}
		})
	}

	buf := []byte("prefix:valid 世界")
	got := internal.ReplaceInvalidUTF8(buf, len("prefix:"))
	if unsafe.SliceData(got) != unsafe.SliceData(buf) {
		t.Fatal("ReplaceInvalidUTF8 copied valid input")
	}

	state := uint64(0x9e3779b97f4a7c15)
	for length := range 64 {
		for range 100 {
			src := make([]byte, length)
			for i := range src {
				state ^= state << 7
				state ^= state >> 9
				state ^= state << 8
				src[i] = byte(state)
			}

			got := internal.ReplaceInvalidUTF8(append([]byte(nil), src...), 0)
			want := replaceInvalidUTF8Reference(src)
			if !bytes.Equal(got, want) {
				t.Fatalf("ReplaceInvalidUTF8(%x) = %x, want %x", src, got, want)
			}
		}
	}
}

func replaceInvalidUTF8Reference(src []byte) []byte {
	dst := make([]byte, 0, len(src))
	for len(src) > 0 {
		r, size := utf8.DecodeRune(src)
		if r == utf8.RuneError && size == 1 {
			dst = append(dst, string(utf8.RuneError)...)
		} else {
			dst = append(dst, src[:size]...)
		}
		src = src[size:]
	}
	return dst
}
