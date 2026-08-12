package jsonexperiment

import (
	"reflect"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

//
// (map[K primitive]V value) marshal functions
//

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

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		keyTarget := reflect.New(typ.Key()).Elem()
		valueTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for iter := value.MapRange(); iter.Next(); {
			keyTarget.SetIterKey(iter)
			valueTarget.SetIterValue(iter)

			if flags&MarshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, keyTarget.String())
			} else {
				dst = internal.AppendString(dst, keyTarget.String())
			}
			dst = append(dst, ':')

			valuePtr := mapValuePointer(valueTarget, valueIsMap)
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
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		keyTarget := reflect.New(typ.Key()).Elem()
		valueTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for iter := value.MapRange(); iter.Next(); {
			keyTarget.SetIterKey(iter)
			valueTarget.SetIterValue(iter)

			dst = append(dst, '"')
			dst = internal.AppendInt(dst, keyTarget.Int())
			dst = append(dst, '"', ':')

			valuePtr := mapValuePointer(valueTarget, valueIsMap)
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
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		keyTarget := reflect.New(typ.Key()).Elem()
		valueTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for iter := value.MapRange(); iter.Next(); {
			keyTarget.SetIterKey(iter)
			valueTarget.SetIterValue(iter)

			dst = append(dst, '"')
			dst = internal.AppendUint(dst, keyTarget.Uint())
			dst = append(dst, '"', ':')

			valuePtr := mapValuePointer(valueTarget, valueIsMap)
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

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		keyTarget := reflect.New(keyType).Elem()
		valueTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for iter := value.MapRange(); iter.Next(); {
			keyTarget.SetIterKey(iter)
			valueTarget.SetIterValue(iter)

			var err error
			if dst, err = appendMapTextKey(dst, keyTarget, keyIsPointer, flags); err != nil {
				return dst, err
			}
			dst = append(dst, ':')

			valuePtr := mapValuePointer(valueTarget, valueIsMap)
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

func mapValuePointer(value reflect.Value, valueIsMap bool) unsafe.Pointer {
	ptr := unsafe.Pointer(value.UnsafeAddr())
	if valueIsMap {
		return *(*unsafe.Pointer)(ptr)
	}
	return ptr
}
