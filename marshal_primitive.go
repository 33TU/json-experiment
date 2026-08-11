package jsonexperiment

import (
	"reflect"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

var builtinScalarTypes = [...]reflect.Type{
	reflect.Bool:    reflect.TypeFor[bool](),
	reflect.Int:     reflect.TypeFor[int](),
	reflect.Int8:    reflect.TypeFor[int8](),
	reflect.Int16:   reflect.TypeFor[int16](),
	reflect.Int32:   reflect.TypeFor[int32](),
	reflect.Int64:   reflect.TypeFor[int64](),
	reflect.Uint:    reflect.TypeFor[uint](),
	reflect.Uint8:   reflect.TypeFor[uint8](),
	reflect.Uint16:  reflect.TypeFor[uint16](),
	reflect.Uint32:  reflect.TypeFor[uint32](),
	reflect.Uint64:  reflect.TypeFor[uint64](),
	reflect.Uintptr: reflect.TypeFor[uintptr](),
	reflect.Float32: reflect.TypeFor[float32](),
	reflect.Float64: reflect.TypeFor[float64](),
	reflect.String:  reflect.TypeFor[string](),
}

func isBuiltinScalar(typ reflect.Type, kind reflect.Kind) bool {
	return int(kind) < len(builtinScalarTypes) && typ == builtinScalarTypes[kind]
}

func tryCreatePrimitiveMarshalFn(kind reflect.Kind) marshalFn {
	switch kind {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendBool(dst, *(*bool)(ptr)), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int)(ptr)), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int8)(ptr)), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int16)(ptr)), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int32)(ptr)), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int64)(ptr)), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint)(ptr)), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint8)(ptr)), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint16)(ptr)), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint32)(ptr)), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint64)(ptr)), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uintptr)(ptr)), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendFloat32(dst, *(*float32)(ptr))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
			return internal.AppendFloat64(dst, *(*float64)(ptr))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			if flags&MarshalFlagEscapeHTML != 0 {
				return internal.AppendStringHTML(dst, *(*string)(ptr)), nil
			}
			return internal.AppendString(dst, *(*string)(ptr)), nil
		}
	}

	return nil
}
