package jsonexperiment

import (
	"reflect"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

func createPointerMarshalFn(typ reflect.Type) marshalFn {
	pointerDepth := 0

	// Collapse the pointer chain and record how many dereferences are required.
	elemType := typ
	for elemType.Kind() == reflect.Pointer {
		pointerDepth++
		elemType = elemType.Elem()
	}

	switch elemType.Kind() {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendBool(dst, *(*bool)(ptr)), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendInt(dst, *(*int)(ptr)), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendInt(dst, *(*int8)(ptr)), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendInt(dst, *(*int16)(ptr)), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendInt(dst, *(*int32)(ptr)), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendInt(dst, *(*int64)(ptr)), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendUint(dst, *(*uint)(ptr)), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendUint(dst, *(*uint8)(ptr)), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendUint(dst, *(*uint16)(ptr)), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendUint(dst, *(*uint32)(ptr)), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendUint(dst, *(*uint64)(ptr)), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendUint(dst, *(*uintptr)(ptr)), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendFloat32(dst, *(*float32)(ptr))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendFloat64(dst, *(*float64)(ptr))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			return internal.AppendString(dst, *(*string)(ptr)), nil
		}
	default:
		elemFn := getOrCreateMarshalFn(elemType)
		elemIsMap := elemType.Kind() == reflect.Map

		return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			for range pointerDepth {
				ptr = *(*unsafe.Pointer)(ptr)
				if ptr == nil {
					return internal.AppendNull(dst), nil
				}
			}
			if elemIsMap {
				ptr = *(*unsafe.Pointer)(ptr)
			}
			return elemFn(dst, ptr)
		}
	}
}
