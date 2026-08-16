package jsonexperiment

import (
	"bytes"
	"reflect"
	"slices"
	"strings"
	"sync"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

type sortedMapState[K any] struct {
	keyTarget reflect.Value
	values    reflect.Value
	keys      []K
	indexes   []int
}

type sortedMapIntegerKey struct {
	text   [32]byte
	length uint8
}

func tryCreateSortedMapValueMarshalFn(
	typ reflect.Type,
	keyType reflect.Type,
	valueType reflect.Type,
	keyKind reflect.Kind,
	valueFn marshalFn,
) marshalFn {
	if keyType.Implements(stdTextMarshalerType) {
		return createSortedMapTextValueMarshalFn(typ, keyType, valueType, valueFn)
	}

	switch keyKind {
	case reflect.String:
		return createSortedMapPrimitiveStringValueMarshalFn(typ, valueType, valueFn)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return createSortedMapPrimitiveIntValueMarshalFn(typ, valueType, valueFn)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return createSortedMapPrimitiveUintValueMarshalFn(typ, valueType, valueFn)
	}

	return nil
}

func createSortedMapPrimitiveStringValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map
	statePool := createSortedMapStatePool[string](typ.Key(), valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapState[string](statePool, length)
		defer putSortedMapState(statePool, state)
		keys := state.keys
		indexes := state.indexes

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			state.keyTarget.SetIterKey(iter)
			state.values.Index(i).SetIterValue(iter)
			keys = append(keys, state.keyTarget.String())
			indexes = append(indexes, i)
			i++
		}
		state.keys = keys
		state.indexes = indexes

		slices.SortFunc(indexes, func(a, b int) int {
			return strings.Compare(keys[a], keys[b])
		})

		dst = append(dst, '{')
		for _, index := range indexes {
			key := keys[index]

			if flags&MarshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, key)
			} else {
				dst = internal.AppendString(dst, key)
			}
			dst = append(dst, ':')

			valuePtr := mapValuePointer(state.values.Index(index), valueIsMap)
			var err error
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func createSortedMapPrimitiveIntValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map
	statePool := createSortedMapStatePool[sortedMapIntegerKey](typ.Key(), valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapState[sortedMapIntegerKey](statePool, length)
		defer putSortedMapState(statePool, state)
		keys := state.keys
		indexes := state.indexes

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			state.keyTarget.SetIterKey(iter)
			state.values.Index(i).SetIterValue(iter)
			key := sortedMapIntegerKey{}
			key.length = uint8(len(internal.AppendInt(key.text[:0], state.keyTarget.Int())))
			keys = append(keys, key)
			indexes = append(indexes, i)
			i++
		}
		state.keys = keys
		state.indexes = indexes

		slices.SortFunc(indexes, func(a, b int) int {
			return compareSortedMapIntegerKeys(&keys[a], &keys[b])
		})

		dst = append(dst, '{')
		for _, index := range indexes {
			key := &keys[index]

			dst = append(dst, '"')
			dst = append(dst, key.text[:key.length]...)
			dst = append(dst, '"', ':')

			valuePtr := mapValuePointer(state.values.Index(index), valueIsMap)
			var err error
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func createSortedMapPrimitiveUintValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map
	statePool := createSortedMapStatePool[sortedMapIntegerKey](typ.Key(), valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapState[sortedMapIntegerKey](statePool, length)
		defer putSortedMapState(statePool, state)
		keys := state.keys
		indexes := state.indexes

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			state.keyTarget.SetIterKey(iter)
			state.values.Index(i).SetIterValue(iter)
			key := sortedMapIntegerKey{}
			key.length = uint8(len(internal.AppendUint(key.text[:0], state.keyTarget.Uint())))
			keys = append(keys, key)
			indexes = append(indexes, i)
			i++
		}
		state.keys = keys
		state.indexes = indexes

		slices.SortFunc(indexes, func(a, b int) int {
			return compareSortedMapIntegerKeys(&keys[a], &keys[b])
		})

		dst = append(dst, '{')
		for _, index := range indexes {
			key := &keys[index]

			dst = append(dst, '"')
			dst = append(dst, key.text[:key.length]...)
			dst = append(dst, '"', ':')

			valuePtr := mapValuePointer(state.values.Index(index), valueIsMap)
			var err error
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func createSortedMapTextValueMarshalFn(
	typ reflect.Type,
	keyType reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	keyIsPointer := keyType.Kind() == reflect.Pointer
	valueIsMap := valueType.Kind() == reflect.Map
	statePool := createSortedMapStatePool[string](keyType, valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapState[string](statePool, length)
		defer putSortedMapState(statePool, state)

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			state.keyTarget.SetIterKey(iter)
			state.values.Index(i).SetIterValue(iter)
			name, err := resolveMapTextKey(state.keyTarget, keyIsPointer)
			if err != nil {
				return dst, err
			}
			state.keys = append(state.keys, name)
			state.indexes = append(state.indexes, i)
			i++
		}

		slices.SortFunc(state.indexes, func(a, b int) int {
			return strings.Compare(state.keys[a], state.keys[b])
		})

		dst = append(dst, '{')
		for _, index := range state.indexes {
			name := state.keys[index]

			if flags&MarshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, name)
			} else {
				dst = internal.AppendString(dst, name)
			}
			dst = append(dst, ':')

			valuePtr := mapValuePointer(state.values.Index(index), valueIsMap)
			var err error
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func compareSortedMapIntegerKeys(a, b *sortedMapIntegerKey) int {
	return bytes.Compare(a.text[:a.length], b.text[:b.length])
}

func createSortedMapStatePool[K any](keyType, valueType reflect.Type) *sync.Pool {
	valueSliceType := reflect.SliceOf(valueType)
	return &sync.Pool{
		New: func() any {
			return &sortedMapState[K]{
				keyTarget: reflect.New(keyType).Elem(),
				values:    reflect.New(valueSliceType).Elem(),
			}
		},
	}
}

func getSortedMapState[K any](pool *sync.Pool, length int) *sortedMapState[K] {
	state := pool.Get().(*sortedMapState[K])
	if state.values.Cap() < length {
		state.values.Set(reflect.MakeSlice(state.values.Type(), length, length))
	} else {
		state.values.SetLen(length)
	}
	if cap(state.keys) < length {
		state.keys = make([]K, 0, length)
	}
	if cap(state.indexes) < length {
		state.indexes = make([]int, 0, length)
	}
	return state
}

func putSortedMapState[K any](pool *sync.Pool, state *sortedMapState[K]) {
	state.keyTarget.SetZero()
	state.values.Clear()
	state.values.SetLen(0)
	clear(state.keys)
	state.keys = state.keys[:0]
	clear(state.indexes)
	state.indexes = state.indexes[:0]
	pool.Put(state)
}
