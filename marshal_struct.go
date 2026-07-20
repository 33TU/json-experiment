package jsonexperiment

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

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

func createStructMarshalFn(typ reflect.Type) marshalFn {
	fieldFns := make([]marshalFn, 0, typ.NumField())

	for i := range typ.NumField() {
		field := typ.Field(i)

		fieldName, ok := jsonFieldName(field)
		if !ok {
			continue
		}

		fieldKey := internal.AppendString(nil, fieldName)
		fieldKey = append(fieldKey, ':')

		fieldType := field.Type
		fieldKind := fieldType.Kind()
		fieldOffset := field.Offset

		var fieldFn marshalFn

		switch fieldKind {
		case reflect.Bool:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendBool(dst, *(*bool)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Int:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendInt(dst, *(*int)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Int8:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendInt(dst, *(*int8)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Int16:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendInt(dst, *(*int16)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Int32:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendInt(dst, *(*int32)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Int64:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendInt(dst, *(*int64)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Uint:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendUint(dst, *(*uint)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Uint8:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendUint(dst, *(*uint8)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Uint16:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendUint(dst, *(*uint16)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Uint32:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendUint(dst, *(*uint32)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Uint64:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendUint(dst, *(*uint64)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Uintptr:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendUint(dst, *(*uintptr)(unsafe.Add(ptr, fieldOffset))), nil
			}
		case reflect.Float32:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendFloat32(dst, *(*float32)(unsafe.Add(ptr, fieldOffset)))
			}
		case reflect.Float64:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendFloat64(dst, *(*float64)(unsafe.Add(ptr, fieldOffset)))
			}
		case reflect.String:
			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)
				return internal.AppendString(dst, *(*string)(unsafe.Add(ptr, fieldOffset))), nil
			}
		default:
			valueFn := getOrCreateMarshalFn(fieldType)
			fieldIsMap := fieldKind == reflect.Map

			fieldFn = func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
				dst = append(dst, fieldKey...)

				fieldPtr := unsafe.Add(ptr, fieldOffset)
				if fieldIsMap {
					fieldPtr = *(*unsafe.Pointer)(fieldPtr)
				}

				return valueFn(dst, fieldPtr)
			}
		}

		fieldFns = append(fieldFns, fieldFn)
	}

	if len(fieldFns) == 0 {
		return func(dst []byte, _ unsafe.Pointer) ([]byte, error) {
			return append(dst, "{}"...), nil
		}
	}

	return func(dst []byte, ptr unsafe.Pointer) ([]byte, error) {
		dst = append(dst, '{')
		for _, fieldFn := range fieldFns {
			var err error
			if dst, err = fieldFn(dst, ptr); err != nil {
				return dst, err
			}

			dst = append(dst, ',')
		}
		dst[len(dst)-1] = '}'

		return dst, nil
	}
}
