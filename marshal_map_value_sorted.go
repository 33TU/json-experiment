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

type sortedMapValueState struct {
	keyTarget reflect.Value
	values    reflect.Value
}

type sortedMapIntegerKey struct {
	text   [32]byte
	length uint8
}

var sortedMapStringKeysPool = sync.Pool{
	New: func() any {
		keys := make([]string, 0)
		return &keys
	},
}

var sortedMapIntKeysPool = sync.Pool{
	New: func() any {
		keys := make([]sortedMapIntegerKey, 0)
		return &keys
	},
}

var sortedMapUintKeysPool = sync.Pool{
	New: func() any {
		keys := make([]sortedMapIntegerKey, 0)
		return &keys
	},
}

var sortedMapIndexesPool = sync.Pool{
	New: func() any {
		indexes := make([]int, 0)
		return &indexes
	},
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
	valueStatePool := createSortedMapValueStatePool(typ.Key(), valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapValueState(valueStatePool, length)
		defer putSortedMapValueState(valueStatePool, state)

		keysPtr, keys := getSortedMapSlice[string](&sortedMapStringKeysPool, length)
		defer putSortedMapSlice(&sortedMapStringKeysPool, keysPtr, keys)

		indexesPtr, indexes := getSortedMapSlice[int](&sortedMapIndexesPool, length)
		defer putSortedMapSlice(&sortedMapIndexesPool, indexesPtr, indexes)

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			state.keyTarget.SetIterKey(iter)
			state.values.Index(i).SetIterValue(iter)
			keys = append(keys, state.keyTarget.String())
			indexes = append(indexes, i)
			i++
		}

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
	valueStatePool := createSortedMapValueStatePool(typ.Key(), valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapValueState(valueStatePool, length)
		defer putSortedMapValueState(valueStatePool, state)

		keysPtr, keys := getSortedMapSlice[sortedMapIntegerKey](&sortedMapIntKeysPool, length)
		defer putSortedMapSlice(&sortedMapIntKeysPool, keysPtr, keys)

		indexesPtr, indexes := getSortedMapSlice[int](&sortedMapIndexesPool, length)
		defer putSortedMapSlice(&sortedMapIndexesPool, indexesPtr, indexes)

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
	valueStatePool := createSortedMapValueStatePool(typ.Key(), valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapValueState(valueStatePool, length)
		defer putSortedMapValueState(valueStatePool, state)

		keysPtr, keys := getSortedMapSlice[sortedMapIntegerKey](&sortedMapUintKeysPool, length)
		defer putSortedMapSlice(&sortedMapUintKeysPool, keysPtr, keys)

		indexesPtr, indexes := getSortedMapSlice[int](&sortedMapIndexesPool, length)
		defer putSortedMapSlice(&sortedMapIndexesPool, indexesPtr, indexes)

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
	valueStatePool := createSortedMapValueStatePool(keyType, valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapValueState(valueStatePool, length)
		defer putSortedMapValueState(valueStatePool, state)

		namesPtr, names := getSortedMapSlice[string](&sortedMapStringKeysPool, length)
		defer putSortedMapSlice(&sortedMapStringKeysPool, namesPtr, names)

		indexesPtr, indexes := getSortedMapSlice[int](&sortedMapIndexesPool, length)
		defer putSortedMapSlice(&sortedMapIndexesPool, indexesPtr, indexes)

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			state.keyTarget.SetIterKey(iter)
			state.values.Index(i).SetIterValue(iter)
			name, err := resolveMapTextKey(state.keyTarget, keyIsPointer)
			if err != nil {
				return dst, err
			}
			names = append(names, name)
			indexes = append(indexes, i)
			i++
		}

		slices.SortFunc(indexes, func(a, b int) int {
			return strings.Compare(names[a], names[b])
		})

		dst = append(dst, '{')
		for _, index := range indexes {
			name := names[index]

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

func getSortedMapSlice[T any](pool *sync.Pool, capacity int) (*[]T, []T) {
	ptr := pool.Get().(*[]T)
	values := *ptr
	if cap(values) < capacity {
		values = make([]T, 0, capacity)
	}
	return ptr, values[:0]
}

func putSortedMapSlice[T any](pool *sync.Pool, ptr *[]T, values []T) {
	clear(values)
	*ptr = values[:0]
	pool.Put(ptr)
}

func createSortedMapValueStatePool(keyType, valueType reflect.Type) *sync.Pool {
	valueSliceType := reflect.SliceOf(valueType)
	return &sync.Pool{
		New: func() any {
			return &sortedMapValueState{
				keyTarget: reflect.New(keyType).Elem(),
				values:    reflect.New(valueSliceType).Elem(),
			}
		},
	}
}

func getSortedMapValueState(pool *sync.Pool, length int) *sortedMapValueState {
	state := pool.Get().(*sortedMapValueState)
	if state.values.Cap() < length {
		state.values.Set(reflect.MakeSlice(state.values.Type(), length, length))
	} else {
		state.values.SetLen(length)
	}
	return state
}

func putSortedMapValueState(pool *sync.Pool, state *sortedMapValueState) {
	state.keyTarget.SetZero()
	state.values.Clear()
	state.values.SetLen(0)
	pool.Put(state)
}
