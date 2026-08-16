package jsonexperiment

import (
	"reflect"
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

type sortedMapStringKey struct {
	text       string
	valueIndex int
}

// reflectMapIterPrefix mirrors the leading fields of reflect.MapIter and
// internal/runtime/maps.Iter. The latter guarantees that its current key and
// element pointers are its first two fields.
type reflectMapIterPrefix struct {
	value reflect.Value
	key   unsafe.Pointer
	elem  unsafe.Pointer
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
	valueSize := valueType.Size()
	valueTypePointer := internal.InterfaceData(valueType)
	statePool := createSortedMapStatePool[sortedMapStringKey](typ.Key(), valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapState[sortedMapStringKey](statePool, length)
		defer putSortedMapState(statePool, state)
		keys := state.keys
		valuesPtr := state.values.UnsafePointer()

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			iterPrefix := (*reflectMapIterPrefix)(unsafe.Pointer(iter))
			runtimeTypedmemmove(valueTypePointer, unsafe.Add(valuesPtr, uintptr(i)*valueSize), iterPrefix.elem)
			key := *(*string)(iterPrefix.key)
			keys = append(keys, sortedMapStringKey{text: key, valueIndex: i})
			i++
		}
		state.keys = keys

		sortSortedMapStringKeys(keys)

		dst = append(dst, '{')
		for _, key := range keys {

			if flags&MarshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, key.text)
			} else {
				dst = internal.AppendString(dst, key.text)
			}
			dst = append(dst, ':')

			valuePtr := sortedMapValueAt(valuesPtr, key.valueIndex, valueSize, valueIsMap)
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
	switch typ.Key().Kind() {
	case reflect.Int:
		return createSortedMapPrimitiveIntValueMarshalFnForKey[int](typ, valueType, valueFn)
	case reflect.Int8:
		return createSortedMapPrimitiveIntValueMarshalFnForKey[int8](typ, valueType, valueFn)
	case reflect.Int16:
		return createSortedMapPrimitiveIntValueMarshalFnForKey[int16](typ, valueType, valueFn)
	case reflect.Int32:
		return createSortedMapPrimitiveIntValueMarshalFnForKey[int32](typ, valueType, valueFn)
	case reflect.Int64:
		return createSortedMapPrimitiveIntValueMarshalFnForKey[int64](typ, valueType, valueFn)
	}

	return nil
}

func createSortedMapPrimitiveIntValueMarshalFnForKey[K internal.SignedInteger](
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map
	valueSize := valueType.Size()
	valueTypePointer := internal.InterfaceData(valueType)
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
		valuesPtr := state.values.UnsafePointer()

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			iterPrefix := (*reflectMapIterPrefix)(unsafe.Pointer(iter))
			runtimeTypedmemmove(valueTypePointer, unsafe.Add(valuesPtr, uintptr(i)*valueSize), iterPrefix.elem)
			key := sortedMapIntegerKey{}
			key.length = uint8(len(internal.AppendInt(key.text[:0], int64(*(*K)(iterPrefix.key)))))
			keys = append(keys, key)
			indexes = append(indexes, i)
			i++
		}
		state.keys = keys
		state.indexes = indexes

		sortSortedMapIntegerIndexes(indexes, keys)

		dst = append(dst, '{')
		for _, index := range indexes {
			key := &keys[index]

			dst = append(dst, '"')
			dst = append(dst, key.text[:key.length]...)
			dst = append(dst, '"', ':')

			valuePtr := sortedMapValueAt(valuesPtr, index, valueSize, valueIsMap)
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
	switch typ.Key().Kind() {
	case reflect.Uint:
		return createSortedMapPrimitiveUintValueMarshalFnForKey[uint](typ, valueType, valueFn)
	case reflect.Uint8:
		return createSortedMapPrimitiveUintValueMarshalFnForKey[uint8](typ, valueType, valueFn)
	case reflect.Uint16:
		return createSortedMapPrimitiveUintValueMarshalFnForKey[uint16](typ, valueType, valueFn)
	case reflect.Uint32:
		return createSortedMapPrimitiveUintValueMarshalFnForKey[uint32](typ, valueType, valueFn)
	case reflect.Uint64:
		return createSortedMapPrimitiveUintValueMarshalFnForKey[uint64](typ, valueType, valueFn)
	case reflect.Uintptr:
		return createSortedMapPrimitiveUintValueMarshalFnForKey[uintptr](typ, valueType, valueFn)
	}

	return nil
}

func createSortedMapPrimitiveUintValueMarshalFnForKey[K internal.UnsignedInteger](
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map
	valueSize := valueType.Size()
	valueTypePointer := internal.InterfaceData(valueType)
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
		valuesPtr := state.values.UnsafePointer()

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			iterPrefix := (*reflectMapIterPrefix)(unsafe.Pointer(iter))
			runtimeTypedmemmove(valueTypePointer, unsafe.Add(valuesPtr, uintptr(i)*valueSize), iterPrefix.elem)
			key := sortedMapIntegerKey{}
			key.length = uint8(len(internal.AppendUint(key.text[:0], uint64(*(*K)(iterPrefix.key)))))
			keys = append(keys, key)
			indexes = append(indexes, i)
			i++
		}
		state.keys = keys
		state.indexes = indexes

		sortSortedMapIntegerIndexes(indexes, keys)

		dst = append(dst, '{')
		for _, index := range indexes {
			key := &keys[index]

			dst = append(dst, '"')
			dst = append(dst, key.text[:key.length]...)
			dst = append(dst, '"', ':')

			valuePtr := sortedMapValueAt(valuesPtr, index, valueSize, valueIsMap)
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
	valueSize := valueType.Size()
	statePool := createSortedMapStatePool[sortedMapStringKey](keyType, valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		length := value.Len()
		state := getSortedMapState[sortedMapStringKey](statePool, length)
		defer putSortedMapState(statePool, state)

		i := 0
		for iter := value.MapRange(); iter.Next(); {
			state.keyTarget.SetIterKey(iter)
			state.values.Index(i).SetIterValue(iter)
			name, err := resolveMapTextKey(state.keyTarget, keyIsPointer)
			if err != nil {
				return dst, err
			}
			state.keys = append(state.keys, sortedMapStringKey{text: name, valueIndex: i})
			i++
		}

		sortSortedMapStringKeys(state.keys)
		valuesPtr := state.values.UnsafePointer()

		dst = append(dst, '{')
		for _, key := range state.keys {

			if flags&MarshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, key.text)
			} else {
				dst = internal.AppendString(dst, key.text)
			}
			dst = append(dst, ':')

			valuePtr := sortedMapValueAt(valuesPtr, key.valueIndex, valueSize, valueIsMap)
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

//go:noescape
//go:linkname runtimeTypedmemmove runtime.typedmemmove
func runtimeTypedmemmove(typ, dst, src unsafe.Pointer)

func sortedMapValueAt(base unsafe.Pointer, index int, size uintptr, valueIsMap bool) unsafe.Pointer {
	ptr := unsafe.Add(base, uintptr(index)*size)
	if valueIsMap {
		return *(*unsafe.Pointer)(ptr)
	}
	return ptr
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
