package jsonexperiment

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

// noescape returns p while hiding it from escape analysis.
// The returned pointer must not outlive the value referenced by p.
//
//go:noescape
//go:linkname noescape runtime.noescape
func noescape(p unsafe.Pointer) unsafe.Pointer

func marshalInterface(dst []byte, v any, flags MarshalFlags) ([]byte, error) {
	if v == nil {
		return internal.AppendNull(dst), nil
	}

	typ := reflect.TypeOf(v)
	ptr := internal.InterfaceData(v)

	var err error

	switch kind := typ.Kind(); kind {
	case reflect.Bool:
		dst = internal.AppendBool(dst, *(*bool)(ptr))
	case reflect.Int:
		dst = internal.AppendInt(dst, *(*int)(ptr))
	case reflect.Int8:
		dst = internal.AppendInt(dst, *(*int8)(ptr))
	case reflect.Int16:
		dst = internal.AppendInt(dst, *(*int16)(ptr))
	case reflect.Int32:
		dst = internal.AppendInt(dst, *(*int32)(ptr))
	case reflect.Int64:
		dst = internal.AppendInt(dst, *(*int64)(ptr))
	case reflect.Uint:
		dst = internal.AppendUint(dst, *(*uint)(ptr))
	case reflect.Uint8:
		dst = internal.AppendUint(dst, *(*uint8)(ptr))
	case reflect.Uint16:
		dst = internal.AppendUint(dst, *(*uint16)(ptr))
	case reflect.Uint32:
		dst = internal.AppendUint(dst, *(*uint32)(ptr))
	case reflect.Uint64:
		dst = internal.AppendUint(dst, *(*uint64)(ptr))
	case reflect.Uintptr:
		dst = internal.AppendUint(dst, *(*uintptr)(ptr))
	case reflect.Float32:
		dst, err = internal.AppendFloat32(dst, *(*float32)(ptr))
	case reflect.Float64:
		dst, err = internal.AppendFloat64(dst, *(*float64)(ptr))
	case reflect.String:
		if flags&MarshalFlagEscapeHTML != 0 {
			dst = internal.AppendStringHTML(dst, *(*string)(ptr))
		} else {
			dst = internal.AppendString(dst, *(*string)(ptr))
		}
	case reflect.Pointer:
		fn := getOrCreateMarshalFn(typ)
		dst, err = fn(dst, noescape(unsafe.Pointer(&ptr)), flags)
	default:
		fn := getOrCreateMarshalFn(typ)
		dst, err = fn(dst, noescape(ptr), flags)
	}

	runtime.KeepAlive(v)

	return dst, err
}

func createInterfaceMarshalFn(typ reflect.Type) marshalFn {
	if typ.NumMethod() == 0 {
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return marshalInterface(dst, *(*any)(ptr), flags)
		}
	}

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		return marshalInterface(dst, internal.NonEmptyInterfaceValue(ptr), flags)
	}
}
