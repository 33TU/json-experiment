package jsonexperiment

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

type sortedMapEntry struct {
	name string
	key  reflect.Value
}

func tryCreateSortedMapValueMarshalFn(
	typ reflect.Type,
	keyType reflect.Type,
	valueType reflect.Type,
	keyKind reflect.Kind,
	valueFn marshalFn,
) marshalFn {
	if keyType.Implements(stdTextMarshalerType) {
		return createSortedMapTextValueMarshalFn(typ, keyType, valueType, valueFn)
	}

	switch keyKind {
	case reflect.String:
		return createSortedMapPrimitiveStringValueMarshalFn(typ, valueType, valueFn)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return createSortedMapPrimitiveIntValueMarshalFn(typ, valueType, valueFn)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return createSortedMapPrimitiveUintValueMarshalFn(typ, valueType, valueFn)
	}

	return nil
}

func createSortedMapPrimitiveStringValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		keys := value.MapKeys()
		slices.SortFunc(keys, func(a, b reflect.Value) int {
			return strings.Compare(a.String(), b.String())
		})

		valueTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for _, key := range keys {
			valueTarget.Set(value.MapIndex(key))

			if flags&MarshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, key.String())
			} else {
				dst = internal.AppendString(dst, key.String())
			}
			dst = append(dst, ':')

			valuePtr := mapValuePointer(valueTarget, valueIsMap)
			var err error
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func createSortedMapPrimitiveIntValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		keys := value.MapKeys()
		entries := make([]sortedMapEntry, len(keys))
		for i, key := range keys {
			entries[i] = sortedMapEntry{name: strconv.FormatInt(key.Int(), 10), key: key}
		}
		slices.SortFunc(entries, compareSortedMapEntries)

		valueTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for _, entry := range entries {
			valueTarget.Set(value.MapIndex(entry.key))

			dst = append(dst, '"')
			dst = append(dst, entry.name...)
			dst = append(dst, '"', ':')

			valuePtr := mapValuePointer(valueTarget, valueIsMap)
			var err error
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func createSortedMapPrimitiveUintValueMarshalFn(
	typ reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		keys := value.MapKeys()
		entries := make([]sortedMapEntry, len(keys))
		for i, key := range keys {
			entries[i] = sortedMapEntry{name: strconv.FormatUint(key.Uint(), 10), key: key}
		}
		slices.SortFunc(entries, compareSortedMapEntries)

		valueTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for _, entry := range entries {
			valueTarget.Set(value.MapIndex(entry.key))

			dst = append(dst, '"')
			dst = append(dst, entry.name...)
			dst = append(dst, '"', ':')

			valuePtr := mapValuePointer(valueTarget, valueIsMap)
			var err error
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func createSortedMapTextValueMarshalFn(
	typ reflect.Type,
	keyType reflect.Type,
	valueType reflect.Type,
	valueFn marshalFn,
) marshalFn {
	keyIsPointer := keyType.Kind() == reflect.Pointer
	valueIsMap := valueType.Kind() == reflect.Map

	return func(dst []byte, ptr unsafe.Pointer, flags MarshalFlags) ([]byte, error) {
		value := reflect.NewAt(typ, noescape(unsafe.Pointer(&ptr))).Elem()
		if value.IsNil() {
			return internal.AppendNull(dst), nil
		} else if value.Len() == 0 {
			return append(dst, "{}"...), nil
		}

		keys := value.MapKeys()
		entries := make([]sortedMapEntry, len(keys))
		for i, key := range keys {
			name, err := resolveMapTextKey(key, keyIsPointer)
			if err != nil {
				return dst, err
			}
			entries[i] = sortedMapEntry{name: name, key: key}
		}
		slices.SortFunc(entries, compareSortedMapEntries)

		valueTarget := reflect.New(valueType).Elem()

		dst = append(dst, '{')
		for _, entry := range entries {
			valueTarget.Set(value.MapIndex(entry.key))

			if flags&MarshalFlagEscapeHTML != 0 {
				dst = internal.AppendStringHTML(dst, entry.name)
			} else {
				dst = internal.AppendString(dst, entry.name)
			}
			dst = append(dst, ':')

			valuePtr := mapValuePointer(valueTarget, valueIsMap)
			var err error
			if dst, err = valueFn(dst, valuePtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}

func compareSortedMapEntries(a, b sortedMapEntry) int {
	return strings.Compare(a.name, b.name)
}
