package internal_test

import (
	"reflect"
	"runtime"
	"testing"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

type testInterface interface {
	value() int
}

type testValue int

func (v testValue) value() int {
	return int(v)
}

type testPointer struct {
	n int
}

func (v *testPointer) value() int {
	return v.n
}

func TestInterfaceData(t *testing.T) {
	value := any(123)
	ptr := internal.InterfaceData(value)

	if got := *(*int)(ptr); got != 123 {
		t.Fatalf("InterfaceData() points to %d, want 123", got)
	}

	runtime.KeepAlive(value)
}

func TestInterfaceValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{"value", testValue(42)},
		{"pointer", &testPointer{n: 42}},
		{"typed nil", (*testPointer)(nil)},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := internal.InterfaceValue(
				internal.InterfaceType(tt.value),
				internal.InterfaceData(tt.value),
			)
			if !reflect.DeepEqual(got, tt.value) {
				t.Fatalf("InterfaceValue() = %#v (%T), want %#v (%T)", got, got, tt.value, tt.value)
			}

			runtime.KeepAlive(tt.value)
			runtime.KeepAlive(got)
		})
	}
}

func TestNonEmptyInterfaceValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value testInterface
	}{
		{"value", testValue(42)},
		{"pointer", &testPointer{n: 42}},
		{"typed nil", (*testPointer)(nil)},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := internal.NonEmptyInterfaceValue(unsafe.Pointer(&tt.value))
			if !reflect.DeepEqual(got, tt.value) {
				t.Fatalf("NonEmptyInterfaceValue() = %#v (%T), want %#v (%T)", got, got, tt.value, tt.value)
			}

			runtime.KeepAlive(tt.value)
		})
	}
}

func BenchmarkInterfaceData(b *testing.B) {
	value := any(123)
	var result unsafe.Pointer

	for b.Loop() {
		result = internal.InterfaceData(value)
	}

	runtime.KeepAlive(value)
	runtime.KeepAlive(result)
}

func BenchmarkInterfaceValue(b *testing.B) {
	value := any(testValue(123))
	typ := internal.InterfaceType(value)
	data := internal.InterfaceData(value)
	var result any

	for b.Loop() {
		result = internal.InterfaceValue(typ, data)
	}

	runtime.KeepAlive(value)
	runtime.KeepAlive(result)
}

func BenchmarkNonEmptyInterfaceValue(b *testing.B) {
	var value testInterface = testValue(123)
	var result any

	for b.Loop() {
		result = internal.NonEmptyInterfaceValue(unsafe.Pointer(&value))
	}

	runtime.KeepAlive(value)
	runtime.KeepAlive(result)
}
