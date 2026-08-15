//go:build goexperiment.simd

package internal

import (
	"simd/archsimd"
	"unicode/utf8"
)

const utf8AVX2Threshold = 384

var (
	utf8FirstHigh = [16]byte{
		0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
		0x80, 0x80, 0x80, 0x80, 0x21, 0x01, 0x15, 0x49,
	}
	utf8FirstLow = [16]byte{
		0xe7, 0xa3, 0x83, 0x83, 0x8b, 0xcb, 0xcb, 0xcb,
		0xcb, 0xcb, 0xcb, 0xcb, 0xcb, 0xdb, 0xcb, 0xcb,
	}
	utf8SecondHigh = [16]byte{
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0xe6, 0xae, 0xba, 0xba, 0x01, 0x01, 0x01, 0x01,
	}
	utf8FirstHighAVX2 = [32]byte{
		0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
		0x80, 0x80, 0x80, 0x80, 0x21, 0x01, 0x15, 0x49,
		0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
		0x80, 0x80, 0x80, 0x80, 0x21, 0x01, 0x15, 0x49,
	}
	utf8FirstLowAVX2 = [32]byte{
		0xe7, 0xa3, 0x83, 0x83, 0x8b, 0xcb, 0xcb, 0xcb,
		0xcb, 0xcb, 0xcb, 0xcb, 0xcb, 0xdb, 0xcb, 0xcb,
		0xe7, 0xa3, 0x83, 0x83, 0x8b, 0xcb, 0xcb, 0xcb,
		0xcb, 0xcb, 0xcb, 0xcb, 0xcb, 0xdb, 0xcb, 0xcb,
	}
	utf8SecondHighAVX2 = [32]byte{
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0xe6, 0xae, 0xba, 0xba, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0xe6, 0xae, 0xba, 0xba, 0x01, 0x01, 0x01, 0x01,
	}
)

func validUTF8(src []byte) bool {
	if len(src) < 32 {
		return utf8.Valid(src)
	}

	original := src
	r, size := utf8.DecodeRune(src)
	if r == utf8.RuneError && size == 1 {
		return false
	}
	src = src[size:]

	if len(original) >= utf8AVX2Threshold && archsimd.X86.AVX2() {
		return validUTF8AVX2(original, src, size)
	}
	return validUTF8SIMD(original, src)
}

func validUTF8AVX2(original, src []byte, processed int) bool {
	const width = 32

	firstHighTable := archsimd.LoadUint8x32(&utf8FirstHighAVX2)
	firstLowTable := archsimd.LoadUint8x32(&utf8FirstLowAVX2)
	secondHighTable := archsimd.LoadUint8x32(&utf8SecondHighAVX2)
	lowNibble := archsimd.BroadcastUint8x32(0x0f)
	continuationBit := archsimd.BroadcastUint8x32(0x80)
	thirdThreshold := archsimd.BroadcastUint8x32(0xdf)
	fourthThreshold := archsimd.BroadcastUint8x32(0xef)
	zero := archsimd.BroadcastUint8x32(0)

	var previous archsimd.Uint8x32
	for len(src) >= width {
		input := archsimd.LoadUint8x32Slice(src)

		// AVX2 byte shifts operate independently on 128-bit lanes. Arrange the
		// preceding lane beside each current lane so UTF-8 sequences spanning
		// either the vector or lane boundary retain their prior three bytes.
		prior := input.SetLo(previous.GetHi()).SetHi(input.GetLo())
		previous1 := input.ConcatShiftBytesRightGrouped(15, prior)
		previous2 := input.ConcatShiftBytesRightGrouped(14, prior)
		previous3 := input.ConcatShiftBytesRightGrouped(13, prior)

		previous1High := previous1.AsUint16x16().ShiftAllRight(4).AsUint8x32().And(lowNibble)
		inputHigh := input.AsUint16x16().ShiftAllRight(4).AsUint8x32().And(lowNibble)
		special := firstHighTable.PermuteOrZeroGrouped(previous1High.AsInt8x32()).
			And(firstLowTable.PermuteOrZeroGrouped(previous1.And(lowNibble).AsInt8x32())).
			And(secondHighTable.PermuteOrZeroGrouped(inputHigh.AsInt8x32()))

		mustContinue := previous2.SubSaturated(thirdThreshold).
			Or(previous3.SubSaturated(fourthThreshold)).
			Greater(zero)
		required := continuationBit.Masked(mustContinue)
		if !required.Xor(special).IsZero() {
			archsimd.ClearAVXUpperBits()
			return false
		}

		previous = input
		src = src[width:]
		processed += width
	}

	archsimd.ClearAVXUpperBits()

	tailStart := processed - 1
	for tailStart > 0 && original[tailStart]&0xc0 == 0x80 {
		tailStart--
	}
	return utf8.Valid(original[tailStart:])
}

func validUTF8SIMD(original, src []byte) bool {
	const width = 16

	firstHighTable := archsimd.LoadUint8x16(&utf8FirstHigh)
	firstLowTable := archsimd.LoadUint8x16(&utf8FirstLow)
	secondHighTable := archsimd.LoadUint8x16(&utf8SecondHigh)
	lowNibble := archsimd.BroadcastUint8x16(0x0f)
	continuationBit := archsimd.BroadcastUint8x16(0x80)
	thirdThreshold := archsimd.BroadcastUint8x16(0xdf)
	fourthThreshold := archsimd.BroadcastUint8x16(0xef)
	zero := archsimd.BroadcastUint8x16(0)

	var previous archsimd.Uint8x16
	for len(src) >= 2*width {
		input0 := archsimd.LoadUint8x16Slice(src)
		previous01 := input0.ConcatShiftBytesRight(15, previous)
		previous02 := input0.ConcatShiftBytesRight(14, previous)
		previous03 := input0.ConcatShiftBytesRight(13, previous)

		previous01High := previous01.AsUint16x8().ShiftAllRight(4).AsUint8x16().And(lowNibble)
		input0High := input0.AsUint16x8().ShiftAllRight(4).AsUint8x16().And(lowNibble)
		special0 := firstHighTable.PermuteOrZero(previous01High.AsInt8x16()).
			And(firstLowTable.PermuteOrZero(previous01.And(lowNibble).AsInt8x16())).
			And(secondHighTable.PermuteOrZero(input0High.AsInt8x16()))

		mustContinue0 := previous02.SubSaturated(thirdThreshold).
			Or(previous03.SubSaturated(fourthThreshold)).
			Greater(zero)
		required0 := continuationBit.Masked(mustContinue0)
		if !required0.Xor(special0).IsZero() {
			return false
		}

		input1 := archsimd.LoadUint8x16Slice(src[width:])
		previous11 := input1.ConcatShiftBytesRight(15, input0)
		previous12 := input1.ConcatShiftBytesRight(14, input0)
		previous13 := input1.ConcatShiftBytesRight(13, input0)

		previous11High := previous11.AsUint16x8().ShiftAllRight(4).AsUint8x16().And(lowNibble)
		input1High := input1.AsUint16x8().ShiftAllRight(4).AsUint8x16().And(lowNibble)
		special1 := firstHighTable.PermuteOrZero(previous11High.AsInt8x16()).
			And(firstLowTable.PermuteOrZero(previous11.And(lowNibble).AsInt8x16())).
			And(secondHighTable.PermuteOrZero(input1High.AsInt8x16()))

		mustContinue1 := previous12.SubSaturated(thirdThreshold).
			Or(previous13.SubSaturated(fourthThreshold)).
			Greater(zero)
		required1 := continuationBit.Masked(mustContinue1)
		if !required1.Xor(special1).IsZero() {
			return false
		}

		previous = input1
		src = src[2*width:]
	}

	for len(src) >= width {
		input := archsimd.LoadUint8x16Slice(src)
		previous1 := input.ConcatShiftBytesRight(15, previous)
		previous2 := input.ConcatShiftBytesRight(14, previous)
		previous3 := input.ConcatShiftBytesRight(13, previous)

		previous1High := previous1.AsUint16x8().ShiftAllRight(4).AsUint8x16().And(lowNibble)
		inputHigh := input.AsUint16x8().ShiftAllRight(4).AsUint8x16().And(lowNibble)
		special := firstHighTable.PermuteOrZero(previous1High.AsInt8x16()).
			And(firstLowTable.PermuteOrZero(previous1.And(lowNibble).AsInt8x16())).
			And(secondHighTable.PermuteOrZero(inputHigh.AsInt8x16()))

		mustContinue := previous2.SubSaturated(thirdThreshold).
			Or(previous3.SubSaturated(fourthThreshold)).
			Greater(zero)
		required := continuationBit.Masked(mustContinue)
		if !required.Xor(special).IsZero() {
			return false
		}

		previous = input
		src = src[width:]
	}

	if len(src) != 0 {
		input := archsimd.LoadUint8x16SlicePart(src)
		previous1 := input.ConcatShiftBytesRight(15, previous)
		previous2 := input.ConcatShiftBytesRight(14, previous)
		previous3 := input.ConcatShiftBytesRight(13, previous)

		previous1High := previous1.AsUint16x8().ShiftAllRight(4).AsUint8x16().And(lowNibble)
		inputHigh := input.AsUint16x8().ShiftAllRight(4).AsUint8x16().And(lowNibble)
		special := firstHighTable.PermuteOrZero(previous1High.AsInt8x16()).
			And(firstLowTable.PermuteOrZero(previous1.And(lowNibble).AsInt8x16())).
			And(secondHighTable.PermuteOrZero(inputHigh.AsInt8x16()))

		mustContinue := previous2.SubSaturated(thirdThreshold).
			Or(previous3.SubSaturated(fourthThreshold)).
			Greater(zero)
		required := continuationBit.Masked(mustContinue)
		if !required.Xor(special).IsZero() {
			return false
		}
	}

	last := len(original) - 1
	for last > 0 && original[last]&0xc0 == 0x80 {
		last--
	}
	return utf8.Valid(original[last:])
}
