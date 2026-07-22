package internal

import "slices"

// This file adapts James Edward Anhalt III's integer-to-ASCII algorithm:
// https://github.com/jeaiii/itoa/blob/main/itoa/itoa_jeaiii.cpp
//
// Copyright (c) 2017 James Edward Anhalt III
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

//

// With a reused destination, this implementation is approximately 3–40% faster
// than strconv.AppendUint/strconv.AppendInt, depending on the number of decimal digits.
// The largest measured gains are typically for values containing 4–10 digits.

const jeaiiiDigitPairs = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

// AppendUintJeaiii appends the base-10 representation of value to dst.
func AppendUintJeaiii(dst []byte, value uint64) []byte {
	if value < 10 {
		return append(dst, byte(value)+'0')
	}

	if value < 100 {
		start := len(dst)
		dst = append(dst, 0, 0)
		writeDigitPair(dst[start:], uint32(value))
		return dst
	}

	if value>>32 == 0 {
		value32 := uint32(value)
		digits := decimalLen32(value32)
		start := len(dst)
		dst = slices.Grow(dst, digits)
		dst = dst[:start+digits]
		formatUint32Jeaiii(dst[start:], value32, digits)
		return dst
	}

	high := value / 100_000_000
	low := uint32(value % 100_000_000)

	if high>>32 == 0 {
		high32 := uint32(high)
		highDigits := decimalLen32(high32)
		start := len(dst)
		dst = slices.Grow(dst, highDigits+8)
		dst = dst[:start+highDigits+8]
		formatUint32Jeaiii(dst[start:], high32, highDigits)
		formatUint32Jeaiii(dst[start+highDigits:], low, 8)
		return dst
	}

	top := uint32(high / 100_000_000)
	middle := uint32(high % 100_000_000)
	topDigits := decimalLen32(top)
	start := len(dst)
	dst = slices.Grow(dst, topDigits+16)
	dst = dst[:start+topDigits+16]
	formatUint32Jeaiii(dst[start:], top, topDigits)
	formatUint32Jeaiii(dst[start+topDigits:], middle, 8)
	formatUint32Jeaiii(dst[start+topDigits+8:], low, 8)
	return dst
}

// AppendIntJeaiii appends the base-10 representation of value to dst.
func AppendIntJeaiii(dst []byte, value int64) []byte {
	unsigned := uint64(value)
	if value < 0 {
		dst = append(dst, '-')
		unsigned = 0 - unsigned
	}
	return AppendUintJeaiii(dst, unsigned)
}

func decimalLen32(value uint32) int {
	if value < 100 {
		if value < 10 {
			return 1
		}
		return 2
	}

	if value < 1_000_000 {
		if value < 10_000 {
			if value < 1_000 {
				return 3
			}
			return 4
		}
		if value < 100_000 {
			return 5
		}
		return 6
	}

	if value < 100_000_000 {
		if value < 10_000_000 {
			return 7
		}
		return 8
	}

	if value < 1_000_000_000 {
		return 9
	}

	return 10
}

func formatUint32Jeaiii(dst []byte, value uint32, digits int) {
	switch digits {
	case 1:
		dst[0] = byte(value) + '0'
	case 2:
		writeDigitPair(dst, value)
	case 3:
		t := 429496730 * uint64(value)
		writeDigitPair(dst, uint32(t>>32))
		t = 10 * uint64(uint32(t))
		dst[2] = byte(t>>32) + '0'
	case 4:
		t := 42949673 * uint64(value)
		writeDigitPair(dst, uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[2:], uint32(t>>32))
	case 5:
		t := 4294968 * uint64(value)
		writeDigitPair(dst, uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[2:], uint32(t>>32))
		t = 10 * uint64(uint32(t))
		dst[4] = byte(t>>32) + '0'
	case 6:
		t := 429497 * uint64(value)
		writeDigitPair(dst, uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[2:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[4:], uint32(t>>32))
	case 7:
		t := (2814749768 * uint64(value)) >> 16
		writeDigitPair(dst, uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[2:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[4:], uint32(t>>32))
		t = 10 * uint64(uint32(t))
		dst[6] = byte(t>>32) + '0'
	case 8:
		t := ((2251799815 * uint64(value)) >> 19) + 4
		writeDigitPair(dst, uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[2:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[4:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[6:], uint32(t>>32))
	case 9:
		t := ((3602879703 * uint64(value)) >> 23) + 4
		writeDigitPair(dst, uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[2:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[4:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[6:], uint32(t>>32))
		t = 10 * uint64(uint32(t))
		dst[8] = byte(t>>32) + '0'
	case 10:
		t := ((2882303762 * uint64(value)) >> 26) + 4
		writeDigitPair(dst, uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[2:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[4:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[6:], uint32(t>>32))
		t = 100 * uint64(uint32(t))
		writeDigitPair(dst[8:], uint32(t>>32))
	}
}

func writeDigitPair(dst []byte, value uint32) {
	index := value * 2
	dst[0] = jeaiiiDigitPairs[index]
	dst[1] = jeaiiiDigitPairs[index+1]
}
