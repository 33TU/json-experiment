package internal

import (
	"fmt"
	"math"
	"strconv"
)

type signedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type unsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// AppendNull appends a JSON null value to dst.
func AppendNull(dst []byte) []byte {
	return append(dst, "null"...)
}

// AppendBool appends the JSON representation of v to dst.
func AppendBool(dst []byte, v bool) []byte {
	if v {
		return append(dst, "true"...)
	}
	return append(dst, "false"...)
}

// AppendInt appends the base-10 JSON representation of v to dst.
func AppendInt[T signedInteger](dst []byte, v T) []byte {
	return strconv.AppendInt(dst, int64(v), 10)
}

// AppendUint appends the base-10 JSON representation of v to dst.
func AppendUint[T unsignedInteger](dst []byte, v T) []byte {
	return strconv.AppendUint(dst, uint64(v), 10)
}

// AppendFloat32 appends the JSON representation of v to dst.
// It returns an error if v is NaN or infinite.
func AppendFloat32(dst []byte, v float32) ([]byte, error) {
	f := float64(v)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return dst, fmt.Errorf("jsonpd: unsupported float value: %v", v)
	}

	abs := float32(math.Abs(f))
	format := byte('f')
	if abs != 0 && (abs < 1e-6 || abs >= 1e21) {
		format = 'e'
	}

	dst = strconv.AppendFloat(dst, f, format, -1, 32)
	if format == 'e' {
		dst = trimExponentZero(dst)
	}

	return dst, nil
}

// AppendFloat64 appends the JSON representation of v to dst.
// It returns an error if v is NaN or infinite.
func AppendFloat64(dst []byte, v float64) ([]byte, error) {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return dst, fmt.Errorf("jsonpd: unsupported float value: %v", v)
	}

	abs := math.Abs(v)
	format := byte('f')
	if abs != 0 && (abs < 1e-6 || abs >= 1e21) {
		format = 'e'
	}

	dst = strconv.AppendFloat(dst, v, format, -1, 64)
	if format == 'e' {
		dst = trimExponentZero(dst)
	}

	return dst, nil
}

func trimExponentZero(dst []byte) []byte {
	n := len(dst)

	// strconv can emit e-09 or e+09, while encoding/json emits e-9/e+9.
	if n >= 4 &&
		dst[n-4] == 'e' &&
		(dst[n-3] == '-' || dst[n-3] == '+') &&
		dst[n-2] == '0' {

		dst[n-2] = dst[n-1]
		dst = dst[:n-1]
	}

	return dst
}
