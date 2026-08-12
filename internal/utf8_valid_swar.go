//go:build !goexperiment.simd

package internal

import "unicode/utf8"

func validUTF8(src []byte) bool {
	return utf8.Valid(src)
}
