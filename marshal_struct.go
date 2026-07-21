package jsonexperiment

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

type structField struct {
	key     []byte
	keyHTML []byte
	offset  uintptr
	marshal marshalFn
	isMap   bool
}

func jsonFieldName(field reflect.StructField) (string, bool) {
	if !field.IsExported() {
		return "", false
	}

	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", false
	}

	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		name = field.Name
	}

	return name, true
}

func fieldMarshalFn(typ reflect.Type, kind reflect.Kind) marshalFn {
	switch kind {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendBool(dst, *(*bool)(ptr)), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int)(ptr)), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int8)(ptr)), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int16)(ptr)), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int32)(ptr)), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendInt(dst, *(*int64)(ptr)), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint)(ptr)), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint8)(ptr)), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint16)(ptr)), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint32)(ptr)), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uint64)(ptr)), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendUint(dst, *(*uintptr)(ptr)), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendFloat32(dst, *(*float32)(ptr))
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			return internal.AppendFloat64(dst, *(*float64)(ptr))
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendStringHTML(dst, *(*string)(ptr)), nil
			}
			return internal.AppendString(dst, *(*string)(ptr)), nil
		}
	default:
		return getOrCreateMarshalFn(typ)
	}
}

func createStructMarshalFn(typ reflect.Type) marshalFn {
	fields := make([]structField, 0, typ.NumField())

	for i := range typ.NumField() {
		field := typ.Field(i)

		fieldName, ok := jsonFieldName(field)
		if !ok {
			continue
		}

		// Precompute plain field key
		fieldKey := make([]byte, 0, len(fieldName)+4)
		if len(fields) != 0 {
			fieldKey = append(fieldKey, ',')
		}
		fieldKey = internal.AppendString(fieldKey, fieldName)
		fieldKey = append(fieldKey, ':')

		// Precompute HTML-escaped field key
		fieldKeyHTML := make([]byte, 0, len(fieldKey))
		if len(fields) != 0 {
			fieldKeyHTML = append(fieldKeyHTML, ',')
		}
		fieldKeyHTML = internal.AppendStringHTML(fieldKeyHTML, fieldName)
		fieldKeyHTML = append(fieldKeyHTML, ':')

		fieldKind := field.Type.Kind()
		fieldFn := fieldMarshalFn(field.Type, fieldKind)

		fields = append(fields, structField{
			key:     fieldKey,
			keyHTML: fieldKeyHTML,
			offset:  field.Offset,
			marshal: fieldFn,
			isMap:   fieldKind == reflect.Map,
		})
	}

	if len(fields) == 0 {
		return func(dst []byte, _ unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			return append(dst, "{}"...), nil
		}
	}

	return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
		escapeHTML := flags&marshalFlagEscapeHTML != 0

		dst = append(dst, '{')

		for _, field := range fields {
			if escapeHTML {
				dst = append(dst, field.keyHTML...)
			} else {
				dst = append(dst, field.key...)
			}

			fieldPtr := unsafe.Add(ptr, field.offset)
			if field.isMap {
				fieldPtr = *(*unsafe.Pointer)(fieldPtr)
			}

			var err error
			if dst, err = field.marshal(dst, fieldPtr, flags); err != nil {
				return dst, err
			}
		}

		return append(dst, '}'), nil
	}
}
