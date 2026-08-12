package jsonexperiment

import (
	"reflect"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

var anyType = reflect.TypeFor[any]()

//
// (map[K primitive]V primitive) and
// (map[K primitive][]V primitive) marshal functions
//

func tryCreateMapPrimitiveMarshalFn(keyType, valueType reflect.Type, keyKind, valueKind reflect.Kind) marshalFn {
	if keyType.Implements(stdTextMarshalerType) || implementsMarshaler(valueType) {
		return nil
	}

	switch keyKind {
	case reflect.String:
		return tryCreateMapPrimitiveStringMarshalFn(valueType, valueKind)
	case reflect.Int:
		return tryCreateMapPrimitiveIntMarshalFn[int](valueType, valueKind)
	case reflect.Int8:
		return tryCreateMapPrimitiveIntMarshalFn[int8](valueType, valueKind)
	case reflect.Int16:
		return tryCreateMapPrimitiveIntMarshalFn[int16](valueType, valueKind)
	case reflect.Int32:
		return tryCreateMapPrimitiveIntMarshalFn[int32](valueType, valueKind)
	case reflect.Int64:
		return tryCreateMapPrimitiveIntMarshalFn[int64](valueType, valueKind)
	case reflect.Uint:
		return tryCreateMapPrimitiveUintMarshalFn[uint](valueType, valueKind)
	case reflect.Uint8:
		return tryCreateMapPrimitiveUintMarshalFn[uint8](valueType, valueKind)
	case reflect.Uint16:
		return tryCreateMapPrimitiveUintMarshalFn[uint16](valueType, valueKind)
	case reflect.Uint32:
		return tryCreateMapPrimitiveUintMarshalFn[uint32](valueType, valueKind)
	case reflect.Uint64:
		return tryCreateMapPrimitiveUintMarshalFn[uint64](valueType, valueKind)
	case reflect.Uintptr:
		return tryCreateMapPrimitiveUintMarshalFn[uintptr](valueType, valueKind)
	}

	return nil
}

func tryCreateMapPrimitiveStringMarshalFn(valueType reflect.Type, valueKind reflect.Kind) marshalFn {
	if valueType == anyType {
		return createMapPrimitiveStringAnyMarshalFn()
	}

	switch valueKind {
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntMapHTML(dst, *(*map[string]int)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntMap(dst, *(*map[string]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntMapHTML(dst, *(*map[string]int8)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntMap(dst, *(*map[string]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntMapHTML(dst, *(*map[string]int16)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntMap(dst, *(*map[string]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntMapHTML(dst, *(*map[string]int32)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntMap(dst, *(*map[string]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntMapHTML(dst, *(*map[string]int64)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntMap(dst, *(*map[string]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintMapHTML(dst, *(*map[string]uint)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintMap(dst, *(*map[string]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintMapHTML(dst, *(*map[string]uint8)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintMap(dst, *(*map[string]uint8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintMapHTML(dst, *(*map[string]uint16)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintMap(dst, *(*map[string]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintMapHTML(dst, *(*map[string]uint32)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintMap(dst, *(*map[string]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintMapHTML(dst, *(*map[string]uint64)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintMap(dst, *(*map[string]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintMapHTML(dst, *(*map[string]uintptr)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintMap(dst, *(*map[string]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringBoolMapHTML(dst, *(*map[string]bool)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringBoolMap(dst, *(*map[string]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringFloat32MapHTML(dst, *(*map[string]float32)(unsafe.Pointer(&ptr)))
			}
			return internal.AppendStringFloat32Map(dst, *(*map[string]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringFloat64MapHTML(dst, *(*map[string]float64)(unsafe.Pointer(&ptr)))
			}
			return internal.AppendStringFloat64Map(dst, *(*map[string]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringStringMapHTML(dst, *(*map[string]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringStringMap(dst, *(*map[string]string)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Slice:
		if fn := tryCreateMapPrimitiveStringSliceMarshalFn(valueType); fn != nil {
			return fn
		}
	}

	return nil
}

func tryCreateMapPrimitiveStringSliceMarshalFn(valueType reflect.Type) marshalFn {
	elemType := valueType.Elem()
	elemKind := elemType.Kind()

	if hasAddressableMarshaler(elemType) {
		return nil
	}

	switch elemKind {
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int8)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int16)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int32)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringIntSliceMapHTML(dst, *(*map[string][]int64)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringIntSliceMap(dst, *(*map[string][]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringByteSliceBase64MapHTML(dst, *(*map[string][]byte)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringByteSliceBase64Map(dst, *(*map[string][]byte)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint16)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint32)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uint64)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringUintSliceMapHTML(dst, *(*map[string][]uintptr)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringUintSliceMap(dst, *(*map[string][]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringBoolSliceMapHTML(dst, *(*map[string][]bool)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringBoolSliceMap(dst, *(*map[string][]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringFloat32SliceMapHTML(dst, *(*map[string][]float32)(unsafe.Pointer(&ptr)))
			}
			return internal.AppendStringFloat32SliceMap(dst, *(*map[string][]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringFloat64SliceMapHTML(dst, *(*map[string][]float64)(unsafe.Pointer(&ptr)))
			}
			return internal.AppendStringFloat64SliceMap(dst, *(*map[string][]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringStringSliceMapHTML(dst, *(*map[string][]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendStringStringSliceMap(dst, *(*map[string][]string)(unsafe.Pointer(&ptr))), nil
		}
	}

	return nil
}

func tryCreateMapPrimitiveIntMarshalFn[K internal.SignedInteger](valueType reflect.Type, valueKind reflect.Kind) marshalFn {
	if valueType == anyType {
		return createMapPrimitiveIntAnyMarshalFn[K]()
	}

	switch valueKind {
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntMap(dst, *(*map[K]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntMap(dst, *(*map[K]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntMap(dst, *(*map[K]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntMap(dst, *(*map[K]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntMap(dst, *(*map[K]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintMap(dst, *(*map[K]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintMap(dst, *(*map[K]uint8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintMap(dst, *(*map[K]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintMap(dst, *(*map[K]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintMap(dst, *(*map[K]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintMap(dst, *(*map[K]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntBoolMap(dst, *(*map[K]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntFloat32Map(dst, *(*map[K]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntFloat64Map(dst, *(*map[K]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendIntStringMapHTML(dst, *(*map[K]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendIntStringMap(dst, *(*map[K]string)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Slice:
		if fn := tryCreateMapPrimitiveIntSliceMarshalFn[K](valueType); fn != nil {
			return fn
		}
	}

	return nil
}

func tryCreateMapPrimitiveIntSliceMarshalFn[K internal.SignedInteger](valueType reflect.Type) marshalFn {
	elemType := valueType.Elem()
	elemKind := elemType.Kind()

	if hasAddressableMarshaler(elemType) {
		return nil
	}

	switch elemKind {
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[K][]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[K][]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[K][]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[K][]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntIntSliceMap(dst, *(*map[K][]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[K][]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntByteSliceBase64Map(dst, *(*map[K][]byte)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[K][]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[K][]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[K][]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntUintSliceMap(dst, *(*map[K][]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntBoolSliceMap(dst, *(*map[K][]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntFloat32SliceMap(dst, *(*map[K][]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntFloat64SliceMap(dst, *(*map[K][]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendIntStringSliceMapHTML(dst, *(*map[K][]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendIntStringSliceMap(dst, *(*map[K][]string)(unsafe.Pointer(&ptr))), nil
		}
	}

	return nil
}

func tryCreateMapPrimitiveUintMarshalFn[K internal.UnsignedInteger](valueType reflect.Type, valueKind reflect.Kind) marshalFn {
	if valueType == anyType {
		return createMapPrimitiveUintAnyMarshalFn[K]()
	}

	switch valueKind {
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntMap(dst, *(*map[K]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntMap(dst, *(*map[K]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntMap(dst, *(*map[K]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntMap(dst, *(*map[K]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntMap(dst, *(*map[K]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintMap(dst, *(*map[K]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintMap(dst, *(*map[K]uint8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintMap(dst, *(*map[K]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintMap(dst, *(*map[K]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintMap(dst, *(*map[K]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintMap(dst, *(*map[K]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintBoolMap(dst, *(*map[K]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintFloat32Map(dst, *(*map[K]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintFloat64Map(dst, *(*map[K]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendUintStringMapHTML(dst, *(*map[K]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendUintStringMap(dst, *(*map[K]string)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Slice:
		if fn := tryCreateMapPrimitiveUintSliceMarshalFn[K](valueType); fn != nil {
			return fn
		}
	}

	return nil
}

func tryCreateMapPrimitiveUintSliceMarshalFn[K internal.UnsignedInteger](valueType reflect.Type) marshalFn {
	elemType := valueType.Elem()
	elemKind := elemType.Kind()

	if hasAddressableMarshaler(elemType) {
		return nil
	}

	switch elemKind {
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[K][]int)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[K][]int8)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[K][]int16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[K][]int32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintIntSliceMap(dst, *(*map[K][]int64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[K][]uint)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintByteSliceBase64Map(dst, *(*map[K][]byte)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[K][]uint16)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[K][]uint32)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[K][]uint64)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintUintSliceMap(dst, *(*map[K][]uintptr)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintBoolSliceMap(dst, *(*map[K][]bool)(unsafe.Pointer(&ptr))), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintFloat32SliceMap(dst, *(*map[K][]float32)(unsafe.Pointer(&ptr)))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintFloat64SliceMap(dst, *(*map[K][]float64)(unsafe.Pointer(&ptr)))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendUintStringSliceMapHTML(dst, *(*map[K][]string)(unsafe.Pointer(&ptr))), nil
			}
			return internal.AppendUintStringSliceMap(dst, *(*map[K][]string)(unsafe.Pointer(&ptr))), nil
		}
	}

	return nil
}

//
// (map[K primitive]any) marshal functions
//

func createMapPrimitiveStringAnyMarshalFn() marshalFn {
	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		m := *(*map[string]any)(unsafe.Pointer(&ptr))

		if m == nil {
			return internal.AppendNull(dst), nil
		} else if len(m) == 0 {
			return append(dst, "{}"...), nil
		}

		dst = append(dst, '{')
		for key, value := range m {
			if flags&MarshalFlagEscapeHTML != 0 {
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
}

func createMapPrimitiveIntAnyMarshalFn[K internal.SignedInteger]() marshalFn {
	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		m := *(*map[K]any)(unsafe.Pointer(&ptr))

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
}

func createMapPrimitiveUintAnyMarshalFn[K internal.UnsignedInteger]() marshalFn {
	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		m := *(*map[K]any)(unsafe.Pointer(&ptr))

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
}
