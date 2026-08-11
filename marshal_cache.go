package jsonexperiment

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// marshalFn appends the JSON encoding of the value at ptr to dst.
// For maps, ptr is the map data pointer rather than a pointer to map storage.
// Flags configure behavior that may vary between marshal operations.
// A marshalFn must not retain ptr after returning.
type marshalFn func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error)

type addressableMarshalFnCacheKey struct {
	typ reflect.Type
}

// Normal entries use reflect.Type keys. Addressable entries use
// addressableMarshalFnCacheKey to keep the two dispatch contexts distinct.
var marshalFnCache sync.Map

// getOrCreateMarshalFn returns a cached marshal function for typ.
func getOrCreateMarshalFn(typ reflect.Type) marshalFn {
	return getOrCreateMarshalFnForKey(typ, typ, false)
}

// getOrCreateAddressableMarshalFn returns a cached marshal function for a value
// whose storage is addressable. The function accepts a pointer to that storage,
// including for map values.
func getOrCreateAddressableMarshalFn(typ reflect.Type) marshalFn {
	return getOrCreateMarshalFnForKey(addressableMarshalFnCacheKey{typ}, typ, true)
}

func getOrCreateMarshalFnForKey(key any, typ reflect.Type, addressable bool) marshalFn {
	if cached, ok := marshalFnCache.Load(key); ok {
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
		func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			wg.Wait()
			return fn(dst, ptr, flags)
		},
	)

	cached, loaded := marshalFnCache.LoadOrStore(key, placeholder)
	if loaded {
		return cached.(marshalFn)
	}

	fn = createMarshalFn(typ, addressable)
	wg.Done()

	marshalFnCache.Store(key, fn)
	return fn
}

func createMarshalFn(typ reflect.Type, addressable bool) marshalFn {
	if addressable {
		return createAddressableMarshalFn(typ)
	}

	switch {
	case typ.Implements(marshalerAppendType):
		return createMarshalerAppendFn(typ)
	case typ.Implements(stdMarshalerType):
		return createStdMarshalerFn(typ)
	case typ.Implements(stdTextMarshalerType):
		return createStdTextMarshalerFn(typ)
	}

	kind := typ.Kind()

	if fn := tryCreatePrimitiveMarshalFn(kind); fn != nil {
		return fn
	}

	switch kind {
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
		return createStructMarshalFn(typ, false)
	default:
		return unsupportedTypeMarshalFn(typ)
	}
}

func createAddressableMarshalFn(typ reflect.Type) marshalFn {
	if implementsMarshaler(typ) {
		return marshalStorageFn(typ, getOrCreateMarshalFn(typ))
	}

	pointerType := reflect.PointerTo(typ)
	if implementsMarshaler(pointerType) {
		pointerFn := getOrCreateMarshalFn(pointerType)
		return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
			return pointerFn(dst, noescape(unsafe.Pointer(&ptr)), flags)
		}
	}

	switch typ.Kind() {
	case reflect.Array:
		return createArrayDefaultMarshalFn(typ, true)
	case reflect.Struct:
		return createStructMarshalFn(typ, true)
	default:
		return marshalStorageFn(typ, getOrCreateMarshalFn(typ))
	}
}

func hasAddressableMarshaler(typ reflect.Type) bool {
	return implementsMarshaler(typ) || implementsMarshaler(reflect.PointerTo(typ))
}

func hasPointerOnlyMarshaler(typ reflect.Type) bool {
	return !implementsMarshaler(typ) && implementsMarshaler(reflect.PointerTo(typ))
}

func marshalStorageFn(typ reflect.Type, fn marshalFn) marshalFn {
	if typ.Kind() != reflect.Map {
		return fn
	}

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		return fn(dst, *(*unsafe.Pointer)(ptr), flags)
	}
}

func unsupportedTypeMarshalFn(typ reflect.Type) marshalFn {
	return func(dst []byte, _ unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
		return dst, fmt.Errorf("jsonexperiment: unsupported type %s", typ)
	}
}
