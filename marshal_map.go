package jsonexperiment

import (
	"reflect"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

func createMapMarshalFn(typ reflect.Type) marshalFn {
	keyType := typ.Key()
	keyKind := keyType.Kind()

	valueType := typ.Elem()
	valueKind := valueType.Kind()

	// fast path for common map types
	if keyKind == reflect.String {
		switch valueKind {
		case reflect.String:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				if flags&marshalFlagEscapeHTML != 0 {
					return internal.AppendStringStringMapHTML(dst, *(*map[string]string)(unsafe.Pointer(&ptr))), nil
				}
				return internal.AppendStringStringMap(dst, *(*map[string]string)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Int:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				if flags&marshalFlagEscapeHTML != 0 {
					return internal.AppendStringIntMapHTML(dst, *(*map[string]int)(unsafe.Pointer(&ptr))), nil
				}
				return internal.AppendStringIntMap(dst, *(*map[string]int)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Uint:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				if flags&marshalFlagEscapeHTML != 0 {
					return internal.AppendStringUintMapHTML(dst, *(*map[string]uint)(unsafe.Pointer(&ptr))), nil
				}
				return internal.AppendStringUintMap(dst, *(*map[string]uint)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Bool:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				if flags&marshalFlagEscapeHTML != 0 {
					return internal.AppendStringBoolMapHTML(dst, *(*map[string]bool)(unsafe.Pointer(&ptr))), nil
				}
				return internal.AppendStringBoolMap(dst, *(*map[string]bool)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Float32:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				if flags&marshalFlagEscapeHTML != 0 {
					return internal.AppendStringFloat32MapHTML(dst, *(*map[string]float32)(unsafe.Pointer(&ptr)))
				}
				return internal.AppendStringFloat32Map(dst, *(*map[string]float32)(unsafe.Pointer(&ptr)))
			}
		case reflect.Float64:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				if flags&marshalFlagEscapeHTML != 0 {
					return internal.AppendStringFloat64MapHTML(dst, *(*map[string]float64)(unsafe.Pointer(&ptr)))
				}
				return internal.AppendStringFloat64Map(dst, *(*map[string]float64)(unsafe.Pointer(&ptr)))
			}
		case reflect.Interface:
			return marshalMapStringInterface
		case reflect.Slice:
			if fn := createMapStringSliceMarshalFn(valueType.Elem().Kind()); fn != nil {
				return fn
			}
		}

		return createMapStringValueMarshalFn(typ, valueType, getOrCreateMarshalFn(valueType))
	}

	if keyKind == reflect.Int {
		switch valueKind {
		case reflect.String:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				if flags&marshalFlagEscapeHTML != 0 {
					return internal.AppendIntStringMapHTML(dst, *(*map[int]string)(unsafe.Pointer(&ptr))), nil
				}
				return internal.AppendIntStringMap(dst, *(*map[int]string)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Int:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendIntIntMap(dst, *(*map[int]int)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Uint:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendIntUintMap(dst, *(*map[int]uint)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Bool:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendIntBoolMap(dst, *(*map[int]bool)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Float32:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendIntFloat32Map(dst, *(*map[int]float32)(unsafe.Pointer(&ptr)))
			}
		case reflect.Float64:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendIntFloat64Map(dst, *(*map[int]float64)(unsafe.Pointer(&ptr)))
			}
		case reflect.Interface:
			return marshalMapIntInterface
		case reflect.Slice:
			if fn := createMapIntSliceMarshalFn(valueType.Elem().Kind()); fn != nil {
				return fn
			}
		}

		return createMapIntValueMarshalFn(typ, valueType, getOrCreateMarshalFn(valueType))
	}

	if keyKind == reflect.Uint {
		switch valueKind {
		case reflect.String:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				if flags&marshalFlagEscapeHTML != 0 {
					return internal.AppendUintStringMapHTML(dst, *(*map[uint]string)(unsafe.Pointer(&ptr))), nil
				}
				return internal.AppendUintStringMap(dst, *(*map[uint]string)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Int:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendUintIntMap(dst, *(*map[uint]int)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Uint:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendUintUintMap(dst, *(*map[uint]uint)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Bool:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendUintBoolMap(dst, *(*map[uint]bool)(unsafe.Pointer(&ptr))), nil
			}
		case reflect.Float32:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendUintFloat32Map(dst, *(*map[uint]float32)(unsafe.Pointer(&ptr)))
			}
		case reflect.Float64:
			return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
				return internal.AppendUintFloat64Map(dst, *(*map[uint]float64)(unsafe.Pointer(&ptr)))
			}
		case reflect.Interface:
			return marshalMapUintInterface
		case reflect.Slice:
			if fn := createMapUintSliceMarshalFn(valueType.Elem().Kind()); fn != nil {
				return fn
			}
		}

		return createMapUintValueMarshalFn(typ, valueType, getOrCreateMarshalFn(valueType))
	}

	return unsupportedTypeMarshalFn(typ)
}

//
// Marshal functions for map[K][]V (where V is a primitive type)
//

func createMapStringSliceMarshalFn(kind reflect.Kind) marshalFn {
	switch kind {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringBoolSliceMapHTML(dst, *(*map[string][]bool)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringBoolSliceMap(dst, *(*map[string][]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int8)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int16)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int32)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int64)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint8)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint16)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint32)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint64)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uintptr)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringFloat32SliceMapHTML(dst, *(*map[string][]float32)(unsafe.Pointer(&ptr)))
			}
			return internal.AppendStringFloat32SliceMap(dst, *(*map[string][]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringFloat64SliceMapHTML(dst, *(*map[string][]float64)(unsafe.Pointer(&ptr)))
			}
			return internal.AppendStringFloat64SliceMap(dst, *(*map[string][]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringStringSliceMapHTML(dst, *(*map[string][]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringStringSliceMap(dst, *(*map[string][]string)(unsafe.Pointer(&ptr))), nil
		}
	}
	return nil
}

func createMapIntSliceMarshalFn(kind reflect.Kind) marshalFn {
	switch kind {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntBoolSliceMap(dst, *(*map[int][]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[int][]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[int][]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[int][]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[int][]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[int][]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[int][]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[int][]uint8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[int][]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[int][]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[int][]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[int][]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntFloat32SliceMap(dst, *(*map[int][]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntFloat64SliceMap(dst, *(*map[int][]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendIntStringSliceMapHTML(dst, *(*map[int][]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendIntStringSliceMap(dst, *(*map[int][]string)(unsafe.Pointer(&ptr))), nil
		}
	}

	return nil
}

func createMapUintSliceMarshalFn(kind reflect.Kind) marshalFn {
	switch kind {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintBoolSliceMap(dst, *(*map[uint][]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[uint][]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[uint][]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[uint][]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[uint][]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[uint][]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[uint][]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[uint][]uint8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[uint][]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[uint][]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[uint][]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[uint][]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintFloat32SliceMap(dst, *(*map[uint][]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintFloat64SliceMap(dst, *(*map[uint][]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendUintStringSliceMapHTML(dst, *(*map[uint][]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendUintStringSliceMap(dst, *(*map[uint][]string)(unsafe.Pointer(&ptr))), nil
		}
	}

	return nil
}

//
// Reflect-based marshal functions for map[K]V (map with arbitrary key and value types)
//

func createMapStringValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		iter := value.MapRange()
		keyTarget := reflect.New(typ.Key()).Elem()
		valTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for iter.Next() {
			keyTarget.SetIterKey(iter)
			valTarget.SetIterValue(iter)

			if flags&marshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, keyTarget.String())
			} else {
				dst = internal.AppendString(dst, keyTarget.String())
			}
			dst = append(dst, ':')

			valPtr := unsafe.Pointer(valTarget.UnsafeAddr())
			if valueIsMap {
				valPtr = *(*unsafe.Pointer)(valPtr)
			}

			var err error
			if dst, err = valueFn(dst, valPtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func createMapIntValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		iter := value.MapRange()
		keyTarget := reflect.New(typ.Key()).Elem()
		valTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for iter.Next() {
			keyTarget.SetIterKey(iter)
			valTarget.SetIterValue(iter)

			dst = append(dst, '"')
			dst = internal.AppendInt(dst, keyTarget.Int())
			dst = append(dst, '"', ':')

			valPtr := unsafe.Pointer(valTarget.UnsafeAddr())
			if valueIsMap {
				valPtr = *(*unsafe.Pointer)(valPtr)
			}

			var err error
			if dst, err = valueFn(dst, valPtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func createMapUintValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		iter := value.MapRange()
		keyTarget := reflect.New(typ.Key()).Elem()
		valTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for iter.Next() {
			keyTarget.SetIterKey(iter)
			valTarget.SetIterValue(iter)

			dst = append(dst, '"')
			dst = internal.AppendUint(dst, keyTarget.Uint())
			dst = append(dst, '"', ':')

			valPtr := unsafe.Pointer(valTarget.UnsafeAddr())
			if valueIsMap {
				valPtr = *(*unsafe.Pointer)(valPtr)
			}

			var err error
			if dst, err = valueFn(dst, valPtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

//
// Marshal functions for map[K]interface{} (map with interface{} values)
//

func marshalMapStringInterface(
	dst []byte,
	ptr unsafe.Pointer,
	flags marshalFlags,
) ([]byte, error) {
	m := *(*map[string]any)(unsafe.Pointer(&ptr))

	if m == nil {
		return internal.AppendNull(dst), nil
	} else if len(m) == 0 {
		return append(dst, "{}"...), nil
	}

	dst = append(dst, '{')
	for key, value := range m {
		if flags&marshalFlagEscapeHTML != 0 {
			dst = internal.AppendStringHTML(dst, key)
		} else {
			dst = internal.AppendString(dst, key)
		}
		dst = append(dst, ':')

		var err error
		if dst, err = marshalInterface(dst, value, flags); err != nil {
			return dst, err
		}

		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

func marshalMapIntInterface(
	dst []byte,
	ptr unsafe.Pointer,
	flags marshalFlags,
) ([]byte, error) {
	m := *(*map[int]any)(unsafe.Pointer(&ptr))

	if m == nil {
		return internal.AppendNull(dst), nil
	} else if len(m) == 0 {
		return append(dst, "{}"...), nil
	}

	dst = append(dst, '{')
	for key, value := range m {
		dst = append(dst, '"')
		dst = internal.AppendInt(dst, key)
		dst = append(dst, '"', ':')

		var err error
		if dst, err = marshalInterface(dst, value, flags); err != nil {
			return dst, err
		}

		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

func marshalMapUintInterface(
	dst []byte,
	ptr unsafe.Pointer,
	flags marshalFlags,
) ([]byte, error) {
	m := *(*map[uint]any)(unsafe.Pointer(&ptr))

	if m == nil {
		return internal.AppendNull(dst), nil
	} else if len(m) == 0 {
		return append(dst, "{}"...), nil
	}

	dst = append(dst, '{')
	for key, value := range m {
		dst = append(dst, '"')
		dst = internal.AppendUint(dst, key)
		dst = append(dst, '"', ':')

		var err error
		if dst, err = marshalInterface(dst, value, flags); err != nil {
			return dst, err
		}

		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}
