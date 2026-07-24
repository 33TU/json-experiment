package jsonexperiment

import (
	"reflect"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

func createArrayMarshalFn(typ reflect.Type) marshalFn {
	arrayLen := typ.Len()

	if arrayLen == 0 {
		return func(dst []byte, _ unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return append(dst, "[]"...), nil
		}
	}

	elemType := typ.Elem()
	elemSize := elemType.Size()

	switch elemKind := elemType.Kind(); elemKind {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendBoolSlice(dst, unsafe.Slice((*bool)(ptr), arrayLen)), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, unsafe.Slice((*int)(ptr), arrayLen)), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, unsafe.Slice((*int8)(ptr), arrayLen)), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, unsafe.Slice((*int16)(ptr), arrayLen)), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, unsafe.Slice((*int32)(ptr), arrayLen)), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, unsafe.Slice((*int64)(ptr), arrayLen)), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, unsafe.Slice((*uint)(ptr), arrayLen)), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, unsafe.Slice((*uint8)(ptr), arrayLen)), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, unsafe.Slice((*uint16)(ptr), arrayLen)), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, unsafe.Slice((*uint32)(ptr), arrayLen)), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, unsafe.Slice((*uint64)(ptr), arrayLen)), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, unsafe.Slice((*uintptr)(ptr), arrayLen)), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendFloat32Slice(dst, unsafe.Slice((*float32)(ptr), arrayLen))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return internal.AppendFloat64Slice(dst, unsafe.Slice((*float64)(ptr), arrayLen))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringSliceHTML(dst, unsafe.Slice((*string)(ptr), arrayLen)), nil
			}
			return internal.AppendStringSlice(dst, unsafe.Slice((*string)(ptr), arrayLen)), nil
		}
	default:
		elemFn := getOrCreateMarshalFn(elemType)
		elemIsMap := elemKind == reflect.Map

		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			dst = append(dst, '[')
			for i := range arrayLen {
				elemPtr := unsafe.Add(ptr, uintptr(i)*elemSize)
				if elemIsMap {
					elemPtr = *(*unsafe.Pointer)(elemPtr)
				}

				var err error
				if dst, err = elemFn(dst, elemPtr, flags); err != nil {
					return dst, err
				}

				dst = append(dst, ',')
			}
			dst[len(dst)-1] = ']'

			return dst, nil
		}
	}
}
