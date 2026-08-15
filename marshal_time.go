package jsonexperiment

import (
	"reflect"
	"time"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

var (
	timeType        = reflect.TypeFor[time.Time]()
	timePointerType = reflect.TypeFor[*time.Time]()
)

func createTimeMarshalFn() marshalFn {
	return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
		return appendTime(dst, *(*time.Time)(ptr))
	}
}

func createTimePointerMarshalFn() marshalFn {
	return func(dst []byte, ptr unsafe.Pointer, _ MarshalFlags) ([]byte, error) {
		value := *(**time.Time)(ptr)
		if value == nil {
			return internal.AppendNull(dst), nil
		}
		return appendTime(dst, *value)
	}
}

func appendTime(dst []byte, value time.Time) ([]byte, error) {
	original := dst
	dst = append(dst, '"')

	var err error
	dst, err = value.AppendText(dst)
	if err != nil {
		return original, err
	}

	return append(dst, '"'), nil
}
