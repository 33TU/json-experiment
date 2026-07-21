package jsonexperiment

import (
	"reflect"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

// sliceHeader mirrors the runtime representation of a slice.
type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

func createSliceMarshalFn(typ reflect.Type) marshalFn {
	elemType := typ.Elem()
	elemSize := elemType.Size()

	switch elemKind := elemType.Kind(); elemKind {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendBoolSlice(dst, *(*[]bool)(ptr)), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, *(*[]int)(ptr)), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, *(*[]int8)(ptr)), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, *(*[]int16)(ptr)), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, *(*[]int32)(ptr)), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendIntSlice(dst, *(*[]int64)(ptr)), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, *(*[]uint)(ptr)), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, *(*[]uint8)(ptr)), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, *(*[]uint16)(ptr)), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, *(*[]uint32)(ptr)), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, *(*[]uint64)(ptr)), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUintSlice(dst, *(*[]uintptr)(ptr)), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendFloat32Slice(dst, *(*[]float32)(ptr))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendFloat64Slice(dst, *(*[]float64)(ptr))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringSliceHTML(dst, *(*[]string)(ptr)), nil
			}
			return internal.AppendStringSlice(dst, *(*[]string)(ptr)), nil
		}
	default:
		elemFn := getOrCreateMarshalFn(elemType)
		elemIsMap := elemKind == reflect.Map

		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			header := (*sliceHeader)(ptr)

			if header.data == nil {
				return internal.AppendNull(dst), nil
			}
			if header.len == 0 {
				return append(dst, "[]"...), nil
			}

			dst = append(dst, '[')
			for i := range header.len {
				elemPtr := unsafe.Add(header.data, uintptr(i)*elemSize)
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
