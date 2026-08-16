package jsonexperiment

import (
	"reflect"
	"sync"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

//
// (map[K primitive]V value) marshal functions
//

type mapValueState struct {
	keyTarget   reflect.Value
	valueTarget reflect.Value
}

func tryCreateMapValueMarshalFn(
	typ reflect.Type,
	keyType reflect.Type,
	valueType reflect.Type,
	keyKind reflect.Kind,
	valueFn marshalFn,
) marshalFn {
	if keyType.Implements(stdTextMarshalerType) {
		return createMapTextValueMarshalFn(typ, keyType, valueType, valueFn)
	}

	switch keyKind {
	case reflect.String:
		return createMapPrimitiveStringValueMarshalFn(typ, valueType, valueFn)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return createMapPrimitiveIntValueMarshalFn(typ, valueType, valueFn)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return createMapPrimitiveUintValueMarshalFn(typ, valueType, valueFn)
	}

	return nil
}

func createMapPrimitiveStringValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map
	valueTypePointer := internal.InterfaceData(valueType)
	statePool := createMapValueStatePool(nil, valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		state := statePool.Get().(*mapValueState)
		defer putMapValueState(statePool, state)
		valueTargetPointer := unsafe.Pointer(state.valueTarget.UnsafeAddr())

		dst = append(dst, '{')
		for iter := value.MapRange(); iter.Next(); {
			iterPrefix := (*reflectMapIterPrefix)(unsafe.Pointer(iter))
			runtimeTypedmemmove(valueTypePointer, valueTargetPointer, iterPrefix.elem)
			key := *(*string)(iterPrefix.key)

			if flags&MarshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, key)
			} else {
				dst = internal.AppendString(dst, key)
			}
			dst = append(dst, ':')

			valuePtr := mapValuePointer(valueTargetPointer, valueIsMap)
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

func createMapPrimitiveIntValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	switch typ.Key().Kind() {
	case reflect.Int:
		return createMapPrimitiveIntValueMarshalFnForKey[int](typ, valueType, valueFn)
	case reflect.Int8:
		return createMapPrimitiveIntValueMarshalFnForKey[int8](typ, valueType, valueFn)
	case reflect.Int16:
		return createMapPrimitiveIntValueMarshalFnForKey[int16](typ, valueType, valueFn)
	case reflect.Int32:
		return createMapPrimitiveIntValueMarshalFnForKey[int32](typ, valueType, valueFn)
	case reflect.Int64:
		return createMapPrimitiveIntValueMarshalFnForKey[int64](typ, valueType, valueFn)
	}

	return nil
}

func createMapPrimitiveIntValueMarshalFnForKey[K internal.SignedInteger](
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map
	valueTypePointer := internal.InterfaceData(valueType)
	statePool := createMapValueStatePool(nil, valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		state := statePool.Get().(*mapValueState)
		defer putMapValueState(statePool, state)
		valueTargetPointer := unsafe.Pointer(state.valueTarget.UnsafeAddr())

		dst = append(dst, '{')
		for iter := value.MapRange(); iter.Next(); {
			iterPrefix := (*reflectMapIterPrefix)(unsafe.Pointer(iter))
			runtimeTypedmemmove(valueTypePointer, valueTargetPointer, iterPrefix.elem)

			dst = append(dst, '"')
			dst = internal.AppendInt(dst, int64(*(*K)(iterPrefix.key)))
			dst = append(dst, '"', ':')

			valuePtr := mapValuePointer(valueTargetPointer, valueIsMap)
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

func createMapPrimitiveUintValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	switch typ.Key().Kind() {
	case reflect.Uint:
		return createMapPrimitiveUintValueMarshalFnForKey[uint](typ, valueType, valueFn)
	case reflect.Uint8:
		return createMapPrimitiveUintValueMarshalFnForKey[uint8](typ, valueType, valueFn)
	case reflect.Uint16:
		return createMapPrimitiveUintValueMarshalFnForKey[uint16](typ, valueType, valueFn)
	case reflect.Uint32:
		return createMapPrimitiveUintValueMarshalFnForKey[uint32](typ, valueType, valueFn)
	case reflect.Uint64:
		return createMapPrimitiveUintValueMarshalFnForKey[uint64](typ, valueType, valueFn)
	case reflect.Uintptr:
		return createMapPrimitiveUintValueMarshalFnForKey[uintptr](typ, valueType, valueFn)
	}

	return nil
}

func createMapPrimitiveUintValueMarshalFnForKey[K internal.UnsignedInteger](
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map
	valueTypePointer := internal.InterfaceData(valueType)
	statePool := createMapValueStatePool(nil, valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		state := statePool.Get().(*mapValueState)
		defer putMapValueState(statePool, state)
		valueTargetPointer := unsafe.Pointer(state.valueTarget.UnsafeAddr())

		dst = append(dst, '{')
		for iter := value.MapRange(); iter.Next(); {
			iterPrefix := (*reflectMapIterPrefix)(unsafe.Pointer(iter))
			runtimeTypedmemmove(valueTypePointer, valueTargetPointer, iterPrefix.elem)

			dst = append(dst, '"')
			dst = internal.AppendUint(dst, uint64(*(*K)(iterPrefix.key)))
			dst = append(dst, '"', ':')

			valuePtr := mapValuePointer(valueTargetPointer, valueIsMap)
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

//
// text-marshaled map key value functions
//

func createMapTextValueMarshalFn(
	typ reflect.Type,
	keyType reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	keyIsPointer := keyType.Kind() == reflect.Pointer
	valueIsMap := valueType.Kind() == reflect.Map
	keyTypePointer := internal.InterfaceData(keyType)
	valueTypePointer := internal.InterfaceData(valueType)
	statePool := createMapValueStatePool(keyType, valueType)

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		state := statePool.Get().(*mapValueState)
		defer putMapValueState(statePool, state)
		keyTargetPointer := unsafe.Pointer(state.keyTarget.UnsafeAddr())
		valueTargetPointer := unsafe.Pointer(state.valueTarget.UnsafeAddr())

		dst = append(dst, '{')
		for iter := value.MapRange(); iter.Next(); {
			iterPrefix := (*reflectMapIterPrefix)(unsafe.Pointer(iter))
			runtimeTypedmemmove(keyTypePointer, keyTargetPointer, iterPrefix.key)
			runtimeTypedmemmove(valueTypePointer, valueTargetPointer, iterPrefix.elem)

			var err error
			if dst, err = appendMapTextKey(dst, state.keyTarget, keyIsPointer, flags); err != nil {
				return dst, err
			}
			dst = append(dst, ':')

			valuePtr := mapValuePointer(valueTargetPointer, valueIsMap)
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func appendMapTextKey(
	dst []byte,
	key reflect.Value,
	keyIsPointer bool,
	flags MarshalFlags,
) ([]byte, error) {
	text, err := resolveMapTextKey(key, keyIsPointer)
	if err != nil {
		return dst, err
	}
	if flags&MarshalFlagEscapeHTML != 0 {
		return internal.AppendStringHTML(dst, text), nil
	}
	return internal.AppendString(dst, text), nil
}

func resolveMapTextKey(key reflect.Value, keyIsPointer bool) (string, error) {
	if keyIsPointer && key.IsNil() {
		return "", nil
	}

	text, err := key.Interface().(StdTextMarshaler).MarshalText()
	return string(text), err
}

//
// helpers
//

func createMapValueStatePool(keyType, valueType reflect.Type) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			state := &mapValueState{
				valueTarget: reflect.New(valueType).Elem(),
			}
			if keyType != nil {
				state.keyTarget = reflect.New(keyType).Elem()
			}
			return state
		},
	}
}

func putMapValueState(pool *sync.Pool, state *mapValueState) {
	if state.keyTarget.IsValid() {
		state.keyTarget.SetZero()
	}
	state.valueTarget.SetZero()
	pool.Put(state)
}

func mapValuePointer(ptr unsafe.Pointer, valueIsMap bool) unsafe.Pointer {
	if valueIsMap {
		return *(*unsafe.Pointer)(ptr)
	}
	return ptr
}
