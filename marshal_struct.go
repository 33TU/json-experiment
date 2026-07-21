package jsonexperiment

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

type omitFn func(ptr unsafe.Pointer) bool

type structFieldFlags uint8

const (
	structFieldOmitEmpty structFieldFlags = 1 << iota
	structFieldOmitZero
	structFieldString
)

type structField struct {
	key     []byte
	keyHTML []byte
	marshal marshalFn
	omit    omitFn
	offset  uintptr
	flags   structFieldFlags
	isMap   bool
}

func parseJSONTag(tag string) (name string, flags structFieldFlags) {
	name, options, _ := strings.Cut(tag, ",")

	for options != "" {
		var option string
		option, options, _ = strings.Cut(options, ",")

		switch option {
		case "omitempty":
			flags |= structFieldOmitEmpty
		case "omitzero":
			flags |= structFieldOmitZero
		case "string":
			flags |= structFieldString
		}
	}

	return name, flags
}

func fieldOmitFn(typ reflect.Type, kind reflect.Kind, omitEmpty, omitZero bool) omitFn {
	if !omitEmpty && !omitZero {
		return nil
	}

	switch kind {
	case reflect.Bool:
		return func(ptr unsafe.Pointer) bool {
			return !*(*bool)(ptr)
		}
	case reflect.Int:
		return func(ptr unsafe.Pointer) bool {
			return *(*int)(ptr) == 0
		}
	case reflect.Int8:
		return func(ptr unsafe.Pointer) bool {
			return *(*int8)(ptr) == 0
		}
	case reflect.Int16:
		return func(ptr unsafe.Pointer) bool {
			return *(*int16)(ptr) == 0
		}
	case reflect.Int32:
		return func(ptr unsafe.Pointer) bool {
			return *(*int32)(ptr) == 0
		}
	case reflect.Int64:
		return func(ptr unsafe.Pointer) bool {
			return *(*int64)(ptr) == 0
		}
	case reflect.Uint:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint)(ptr) == 0
		}
	case reflect.Uint8:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint8)(ptr) == 0
		}
	case reflect.Uint16:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint16)(ptr) == 0
		}
	case reflect.Uint32:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint32)(ptr) == 0
		}
	case reflect.Uint64:
		return func(ptr unsafe.Pointer) bool {
			return *(*uint64)(ptr) == 0
		}
	case reflect.Uintptr:
		return func(ptr unsafe.Pointer) bool {
			return *(*uintptr)(ptr) == 0
		}
	case reflect.Float32:
		return func(ptr unsafe.Pointer) bool {
			return *(*float32)(ptr) == 0
		}
	case reflect.Float64:
		return func(ptr unsafe.Pointer) bool {
			return *(*float64)(ptr) == 0
		}
	case reflect.String:
		return func(ptr unsafe.Pointer) bool {
			return len(*(*string)(ptr)) == 0
		}
	case reflect.Slice:
		if omitEmpty {
			return func(ptr unsafe.Pointer) bool {
				return (*sliceHeader)(ptr).len == 0
			}
		}
		return func(ptr unsafe.Pointer) bool {
			return (*sliceHeader)(ptr).data == nil
		}
	case reflect.Map:
		if omitEmpty {
			return func(ptr unsafe.Pointer) bool {
				return reflect.NewAt(typ, ptr).Elem().Len() == 0
			}
		}
		return func(ptr unsafe.Pointer) bool {
			return *(*unsafe.Pointer)(ptr) == nil
		}
	case reflect.Pointer:
		return func(ptr unsafe.Pointer) bool {
			return *(*unsafe.Pointer)(ptr) == nil
		}

	case reflect.Interface:
		return func(ptr unsafe.Pointer) bool {
			return *(*unsafe.Pointer)(ptr) == nil // A nil interface has a nil type/itab word.
		}
	case reflect.Array:
		if omitEmpty && typ.Len() == 0 {
			return func(unsafe.Pointer) bool {
				return true
			}
		}
		if !omitZero {
			return nil
		}
		return func(ptr unsafe.Pointer) bool {
			return reflect.NewAt(typ, ptr).Elem().IsZero()
		}
	default:
		if !omitZero {
			return nil
		}
		return func(ptr unsafe.Pointer) bool {
			return reflect.NewAt(typ, ptr).Elem().IsZero()
		}
	}
}

// fieldMarshalFn returns a marshal function for a struct field based on its type and kind.
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
		return getOrCreateMarshalFn(typ) // unaffected by the "string" flag, fallback to default marshal function
	}
}

// fieldStringMarshalFn returns a marshal function for a struct field with the "string" flag.
func fieldStringMarshalFn(typ reflect.Type, kind reflect.Kind) marshalFn {
	switch kind {
	case reflect.Bool:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendBool(dst, *(*bool)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Int:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendInt(dst, *(*int)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Int8:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendInt(dst, *(*int8)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Int16:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendInt(dst, *(*int16)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Int32:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendInt(dst, *(*int32)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Int64:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendInt(dst, *(*int64)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Uint:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendUint(dst, *(*uint)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Uint8:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendUint(dst, *(*uint8)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Uint16:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendUint(dst, *(*uint16)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Uint32:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendUint(dst, *(*uint32)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Uint64:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendUint(dst, *(*uint64)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Uintptr:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			dst = internal.AppendUint(dst, *(*uintptr)(ptr))
			return append(dst, '"'), nil
		}
	case reflect.Float32:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			var err error
			dst, err = internal.AppendFloat32(dst, *(*float32)(ptr))
			if err != nil {
				return dst, err
			}
			return append(dst, '"'), nil
		}
	case reflect.Float64:
		return func(dst []byte, ptr unsafe.Pointer, _ marshalFlags) ([]byte, error) {
			dst = append(dst, '"')
			var err error
			dst, err = internal.AppendFloat64(dst, *(*float64)(ptr))
			if err != nil {
				return dst, err
			}
			return append(dst, '"'), nil
		}
	case reflect.String:
		return func(dst []byte, ptr unsafe.Pointer, flags marshalFlags) ([]byte, error) {
			if flags&marshalFlagEscapeHTML != 0 {
				return internal.AppendQuotedStringHTML(dst, *(*string)(ptr)), nil
			}
			return internal.AppendQuotedString(dst, *(*string)(ptr)), nil
		}
	default:
		return getOrCreateMarshalFn(typ)
	}
}

func createStructMarshalFn(typ reflect.Type) marshalFn {
	fields := make([]structField, 0, typ.NumField())

	for i := range typ.NumField() {
		field := typ.Field(i)
		fieldKind := field.Type.Kind()

		// Get the JSON field fieldName and check if the field should be ignored
		fieldName, flags := parseJSONTag(field.Tag.Get("json"))
		if fieldName == "-" {
			continue
		}
		if fieldName == "" {
			fieldName = field.Name
		}

		// Precompute plain field key
		fieldKey := make([]byte, 0, len(fieldName)+4)
		fieldKey = internal.AppendString(fieldKey, fieldName)
		fieldKey = append(fieldKey, ':')

		// Precompute HTML-escaped field key
		fieldKeyHTML := make([]byte, 0, len(fieldKey))
		fieldKeyHTML = internal.AppendStringHTML(fieldKeyHTML, fieldName)
		fieldKeyHTML = append(fieldKeyHTML, ':')

		// Determine the marshal function for the field's type
		var fieldFn marshalFn
		if flags&structFieldString != 0 {
			fieldFn = fieldStringMarshalFn(field.Type, fieldKind)
		} else {
			fieldFn = fieldMarshalFn(field.Type, fieldKind)
		}

		omitFn := fieldOmitFn(
			field.Type,
			fieldKind,
			flags&structFieldOmitEmpty != 0,
			flags&structFieldOmitZero != 0,
		)

		fields = append(fields, structField{
			key:     fieldKey,
			keyHTML: fieldKeyHTML,
			marshal: fieldFn,
			omit:    omitFn,
			offset:  field.Offset,
			flags:   flags,
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
			fieldPtr := unsafe.Add(ptr, field.offset)

			// Skip the field if it should be omitted based on the omit function
			if field.omit != nil && field.omit(fieldPtr) {
				continue
			}

			if escapeHTML {
				dst = append(dst, field.keyHTML...)
			} else {
				dst = append(dst, field.key...)
			}

			if field.isMap {
				fieldPtr = *(*unsafe.Pointer)(fieldPtr)
			}

			var err error
			if dst, err = field.marshal(dst, fieldPtr, flags); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}

		// Handle the case where no fields were marshaled (e.g., all fields were omitted)
		if dst[len(dst)-1] == '{' {
			return append(dst, '}'), nil
		}

		dst[len(dst)-1] = '}'
		return dst, nil
	}
}
