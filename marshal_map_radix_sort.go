package jsonexperiment

import (
	"bytes"
	"slices"
	"strings"
)

func sortedMapRadixDepthLimit(length int) int {
	depthLimit := 0
	for n := length; n > 0; n >>= 1 {
		depthLimit++
	}
	return depthLimit * 2
}

func sortSortedMapStringKeys(keys []sortedMapStringKey) {
	radixSortSortedMapStringKeys(keys, 0, sortedMapRadixDepthLimit(len(keys)))
}

func radixSortSortedMapStringKeys(keys []sortedMapStringKey, depth, depthLimit int) {
	const insertionSortThreshold = 12
	for len(keys) > insertionSortThreshold {
		if depthLimit == 0 {
			slices.SortFunc(keys, func(a, b sortedMapStringKey) int {
				return strings.Compare(a.text, b.text)
			})
			return
		}
		depthLimit--

		pivot := sortedMapStringRadixPivot(keys, depth)
		left, current, right := 0, 0, len(keys)
		for current < right {
			value := sortedMapStringKeyByteAt(keys[current].text, depth)
			switch {
			case value < pivot:
				keys[left], keys[current] = keys[current], keys[left]
				left++
				current++
			case value > pivot:
				right--
				keys[current], keys[right] = keys[right], keys[current]
			default:
				current++
			}
		}

		if left == 0 && right == len(keys) {
			if pivot < 0 {
				return
			}
			depth = sortedMapStringCommonPrefixDepth(keys, depth+1)
			continue
		}

		if pivot < 0 {
			if left > len(keys)-right {
				radixSortSortedMapStringKeys(keys[right:], depth, depthLimit)
				keys = keys[:left]
			} else {
				radixSortSortedMapStringKeys(keys[:left], depth, depthLimit)
				keys = keys[right:]
			}
			continue
		}

		leftLength := left
		middleLength := right - left
		rightLength := len(keys) - right
		switch max(leftLength, middleLength, rightLength) {
		case leftLength:
			radixSortSortedMapStringKeys(keys[left:right], depth+1, depthLimit)
			radixSortSortedMapStringKeys(keys[right:], depth, depthLimit)
			keys = keys[:left]
		case middleLength:
			radixSortSortedMapStringKeys(keys[:left], depth, depthLimit)
			radixSortSortedMapStringKeys(keys[right:], depth, depthLimit)
			keys = keys[left:right]
			depth++
		default:
			radixSortSortedMapStringKeys(keys[:left], depth, depthLimit)
			radixSortSortedMapStringKeys(keys[left:right], depth+1, depthLimit)
			keys = keys[right:]
		}
	}
	insertionSortSortedMapStringKeys(keys, depth)
}

func sortedMapStringCommonPrefixDepth(keys []sortedMapStringKey, depth int) int {
	prefix := keys[0].text
	commonEnd := len(prefix)
	for i := 1; i < len(keys) && commonEnd > depth; i++ {
		key := keys[i].text
		commonEnd = min(commonEnd, len(key))
		for j := depth; j < commonEnd; j++ {
			if prefix[j] != key[j] {
				commonEnd = j
				break
			}
		}
	}
	return commonEnd
}

func insertionSortSortedMapStringKeys(keys []sortedMapStringKey, depth int) {
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && sortedMapStringKeyLessFrom(keys[j].text, keys[j-1].text, depth); j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
}

func sortedMapStringKeyLessFrom(a, b string, depth int) bool {
	length := min(len(a), len(b))
	for i := depth; i < length; i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}
	return len(a) < len(b)
}

func sortedMapStringKeyByteAt(key string, depth int) int {
	if depth < len(key) {
		return int(key[depth])
	}
	return -1
}

func sortedMapStringRadixPivot(keys []sortedMapStringKey, depth int) int {
	middle := len(keys) / 2
	if len(keys) > 40 {
		step := len(keys) / 8
		return medianSortedMapRadixByte(
			medianSortedMapRadixByte(
				sortedMapStringKeyByteAt(keys[0].text, depth),
				sortedMapStringKeyByteAt(keys[step].text, depth),
				sortedMapStringKeyByteAt(keys[2*step].text, depth),
			),
			medianSortedMapRadixByte(
				sortedMapStringKeyByteAt(keys[middle-step].text, depth),
				sortedMapStringKeyByteAt(keys[middle].text, depth),
				sortedMapStringKeyByteAt(keys[middle+step].text, depth),
			),
			medianSortedMapRadixByte(
				sortedMapStringKeyByteAt(keys[len(keys)-1-2*step].text, depth),
				sortedMapStringKeyByteAt(keys[len(keys)-1-step].text, depth),
				sortedMapStringKeyByteAt(keys[len(keys)-1].text, depth),
			),
		)
	}
	return medianSortedMapRadixByte(
		sortedMapStringKeyByteAt(keys[0].text, depth),
		sortedMapStringKeyByteAt(keys[middle].text, depth),
		sortedMapStringKeyByteAt(keys[len(keys)-1].text, depth),
	)
}

func sortSortedMapIntegerIndexes(indexes []int, keys []sortedMapIntegerKey) {
	radixSortSortedMapIntegerIndexes(indexes, keys, 0, sortedMapRadixDepthLimit(len(indexes)))
}

func radixSortSortedMapIntegerIndexes(indexes []int, keys []sortedMapIntegerKey, depth, depthLimit int) {
	const insertionSortThreshold = 12
	for len(indexes) > insertionSortThreshold {
		if depthLimit == 0 {
			slices.SortFunc(indexes, func(a, b int) int {
				aKey := &keys[a]
				bKey := &keys[b]
				return bytes.Compare(aKey.text[:aKey.length], bKey.text[:bKey.length])
			})
			return
		}
		depthLimit--

		pivot := sortedMapIntegerRadixPivot(indexes, keys, depth)
		left, current, right := 0, 0, len(indexes)
		for current < right {
			value := sortedMapIntegerKeyByteAt(&keys[indexes[current]], depth)
			switch {
			case value < pivot:
				indexes[left], indexes[current] = indexes[current], indexes[left]
				left++
				current++
			case value > pivot:
				right--
				indexes[current], indexes[right] = indexes[right], indexes[current]
			default:
				current++
			}
		}

		if pivot < 0 {
			if left > len(indexes)-right {
				radixSortSortedMapIntegerIndexes(indexes[right:], keys, depth, depthLimit)
				indexes = indexes[:left]
			} else {
				radixSortSortedMapIntegerIndexes(indexes[:left], keys, depth, depthLimit)
				indexes = indexes[right:]
			}
			continue
		}

		leftLength := left
		middleLength := right - left
		rightLength := len(indexes) - right
		switch max(leftLength, middleLength, rightLength) {
		case leftLength:
			radixSortSortedMapIntegerIndexes(indexes[left:right], keys, depth+1, depthLimit)
			radixSortSortedMapIntegerIndexes(indexes[right:], keys, depth, depthLimit)
			indexes = indexes[:left]
		case middleLength:
			radixSortSortedMapIntegerIndexes(indexes[:left], keys, depth, depthLimit)
			radixSortSortedMapIntegerIndexes(indexes[right:], keys, depth, depthLimit)
			indexes = indexes[left:right]
			depth++
		default:
			radixSortSortedMapIntegerIndexes(indexes[:left], keys, depth, depthLimit)
			radixSortSortedMapIntegerIndexes(indexes[left:right], keys, depth+1, depthLimit)
			indexes = indexes[right:]
		}
	}
	insertionSortSortedMapIntegerIndexes(indexes, keys, depth)
}

func insertionSortSortedMapIntegerIndexes(indexes []int, keys []sortedMapIntegerKey, depth int) {
	for i := 1; i < len(indexes); i++ {
		for j := i; j > 0 && sortedMapIntegerKeyLessFrom(&keys[indexes[j]], &keys[indexes[j-1]], depth); j-- {
			indexes[j], indexes[j-1] = indexes[j-1], indexes[j]
		}
	}
}

func sortedMapIntegerKeyByteAt(key *sortedMapIntegerKey, depth int) int {
	if depth < int(key.length) {
		return int(key.text[depth])
	}
	return -1
}

func sortedMapIntegerRadixPivot(indexes []int, keys []sortedMapIntegerKey, depth int) int {
	middle := len(indexes) / 2
	if len(indexes) > 40 {
		step := len(indexes) / 8
		return medianSortedMapRadixByte(
			medianSortedMapRadixByte(
				sortedMapIntegerKeyByteAt(&keys[indexes[0]], depth),
				sortedMapIntegerKeyByteAt(&keys[indexes[step]], depth),
				sortedMapIntegerKeyByteAt(&keys[indexes[2*step]], depth),
			),
			medianSortedMapRadixByte(
				sortedMapIntegerKeyByteAt(&keys[indexes[middle-step]], depth),
				sortedMapIntegerKeyByteAt(&keys[indexes[middle]], depth),
				sortedMapIntegerKeyByteAt(&keys[indexes[middle+step]], depth),
			),
			medianSortedMapRadixByte(
				sortedMapIntegerKeyByteAt(&keys[indexes[len(indexes)-1-2*step]], depth),
				sortedMapIntegerKeyByteAt(&keys[indexes[len(indexes)-1-step]], depth),
				sortedMapIntegerKeyByteAt(&keys[indexes[len(indexes)-1]], depth),
			),
		)
	}
	return medianSortedMapRadixByte(
		sortedMapIntegerKeyByteAt(&keys[indexes[0]], depth),
		sortedMapIntegerKeyByteAt(&keys[indexes[middle]], depth),
		sortedMapIntegerKeyByteAt(&keys[indexes[len(indexes)-1]], depth),
	)
}

func sortedMapIntegerKeyLessFrom(a, b *sortedMapIntegerKey, depth int) bool {
	length := min(a.length, b.length)
	for i := uint8(depth); i < length; i++ {
		if a.text[i] != b.text[i] {
			return a.text[i] < b.text[i]
		}
	}
	return a.length < b.length
}

func medianSortedMapRadixByte(a, b, c int) int {
	if a > b {
		a, b = b, a
	}
	if c < a {
		return a
	}
	if c > b {
		return b
	}
	return c
}
