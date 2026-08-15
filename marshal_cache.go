package jsonexperiment

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
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

const (
	marshalFnFastCacheBits = 8
	marshalFnFastCacheSize = 1 << marshalFnFastCacheBits
)

type marshalFnFastCacheEntry struct {
	typ unsafe.Pointer
	fn  marshalFn
}

type marshalFnCacheValue struct {
	placeholder marshalFn
	ready       atomic.Pointer[marshalFnFastCacheEntry]
}

// Completed codecs use a direct-mapped atomic cache on the hot path. Separate
// tables keep normal and addressable codecs distinct. Collisions fall through
// to marshalFnCache and repopulate the slot with the most recently used codec.
var marshalFnFastCaches [2][marshalFnFastCacheSize]atomic.Pointer[marshalFnFastCacheEntry]

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
	cacheIndex := 0
	if addressable {
		cacheIndex = 1
	}
	fastKey := reflectTypePointer(typ)
	fastSlot := &marshalFnFastCaches[cacheIndex][marshalFnFastCacheIndex(fastKey)]
	if entry := fastSlot.Load(); entry != nil && entry.typ == fastKey {
		return entry.fn
	}

	if cached, ok := marshalFnCache.Load(key); ok {
		return loadMarshalFnCacheValue(cached.(*marshalFnCacheValue), fastSlot)
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
	cacheValue := &marshalFnCacheValue{placeholder: placeholder}

	cached, loaded := marshalFnCache.LoadOrStore(key, cacheValue)
	if loaded {
		return loadMarshalFnCacheValue(cached.(*marshalFnCacheValue), fastSlot)
	}

	fn = createMarshalFn(typ, addressable)
	ready := &marshalFnFastCacheEntry{typ: fastKey, fn: fn}
	cacheValue.ready.Store(ready)
	wg.Done()

	fastSlot.Store(ready)
	return fn
}

func loadMarshalFnCacheValue(
	value *marshalFnCacheValue,
	fastSlot *atomic.Pointer[marshalFnFastCacheEntry],
) marshalFn {
	if ready := value.ready.Load(); ready != nil {
		fastSlot.Store(ready)
		return ready.fn
	}
	return value.placeholder
}

func marshalFnFastCacheIndex(typ unsafe.Pointer) int {
	const fibonacci = uint64(11400714819323198485)
	return int(uint64(uintptr(typ)) * fibonacci >> (64 - marshalFnFastCacheBits))
}

func reflectTypePointer(typ reflect.Type) unsafe.Pointer {
	type reflectTypeInterface struct {
		tab  unsafe.Pointer
		data unsafe.Pointer
	}

	return (*reflectTypeInterface)(unsafe.Pointer(&typ)).data
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
