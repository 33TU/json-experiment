package jsonexperiment

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// marshalFn appends the JSON encoding of the value at ptr to dst.
// For maps, ptr is the map data pointer rather than a pointer to map storage.
// A marshalFn must not retain ptr after returning.
type marshalFn func(dst []byte, ptr unsafe.Pointer) ([]byte, error)

// map[reflect.Type]marshalFn
var marshalFnCache sync.Map

// getOrCreateMarshalFn returns a cached marshal function for a non-primitive type.
func getOrCreateMarshalFn(typ reflect.Type) marshalFn {
	if cached, ok := marshalFnCache.Load(typ); ok {
		return cached.(marshalFn)
	}

	var (
		wg sync.WaitGroup
		fn marshalFn
	)
	wg.Add(1)

	// Recursive type construction receives this placeholder.
	// e.g `type Node struct { Next *Node }`
	placeholder := marshalFn(
		func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
			wg.Wait()
			return fn(dst, ptr)
		},
	)

	cached, loaded := marshalFnCache.LoadOrStore(typ, placeholder)
	if loaded {
		return cached.(marshalFn)
	}

	fn = createMarshalFn(typ)
	wg.Done()

	marshalFnCache.Store(typ, fn)
	return fn
}

func createMarshalFn(typ reflect.Type) marshalFn {
	switch typ.Kind() {
	case reflect.Pointer:
		return createPointerMarshalFn(typ)
	case reflect.Interface:
		return createInterfaceMarshalFn(typ)
	case reflect.Array:
		return createArrayMarshalFn(typ)
	case reflect.Slice:
		return createSliceMarshalFn(typ)
	case reflect.Map:
		return createMapMarshalFn(typ)
	case reflect.Struct:
		return createStructMarshalFn(typ)
	default:
		return unsupportedTypeMarshalFn(typ)
	}
}

func unsupportedTypeMarshalFn(typ reflect.Type) marshalFn {
	return func(dst []byte, _ unsafe.Pointer) ([]byte, error) {
		return dst, fmt.Errorf("jsonexperiment: unsupported type %s", typ)
	}
}
