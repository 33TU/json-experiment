package jsonexperiment

import (
	"reflect"
	"unsafe"
)

// reflectMapIterPrefix mirrors the leading fields of reflect.MapIter and
// internal/runtime/maps.Iter. The latter guarantees that its current key and
// element pointers are its first two fields.
type reflectMapIterPrefix struct {
	value reflect.Value
	key   unsafe.Pointer
	elem  unsafe.Pointer
}

//go:noescape
//go:linkname runtimeTypedmemmove runtime.typedmemmove
func runtimeTypedmemmove(typ, dst, src unsafe.Pointer)

func createMapMarshalFn(typ reflect.Type) marshalFn {
	keyType := typ.Key()
	valueType := typ.Elem()
	unsortedFn := createUnsortedMapMarshalFn(typ, keyType, valueType, keyType.Kind(), valueType.Kind())
	sortedFn := createSortedMapMarshalFn(typ, keyType, valueType, keyType.Kind())

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		if flags&MarshalFlagSortMapKeys != 0 {
			return sortedFn(dst, ptr, flags)
		}
		return unsortedFn(dst, ptr, flags)
	}
}

func createSortedMapMarshalFn(typ, keyType, valueType reflect.Type, keyKind reflect.Kind) marshalFn {
	if fn := tryCreateSortedMapValueMarshalFn(typ, keyType, valueType, keyKind, getOrCreateMarshalFn(valueType)); fn != nil {
		return fn
	}

	return unsupportedTypeMarshalFn(typ)
}

func createUnsortedMapMarshalFn(typ, keyType, valueType reflect.Type, keyKind, valueKind reflect.Kind) marshalFn {
	if fn := tryCreateMapPrimitiveMarshalFn(keyType, valueType, keyKind, valueKind); fn != nil {
		return fn
	}

	if fn := tryCreateMapValueMarshalFn(typ, keyType, valueType, keyKind, getOrCreateMarshalFn(valueType)); fn != nil {
		return fn
	}

	return unsupportedTypeMarshalFn(typ)
}
