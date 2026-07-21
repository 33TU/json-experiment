package jsonexperiment_test

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"runtime"
	"testing"

	jsonexperiment "github.com/33TU/json-experiment"
	"github.com/bytedance/sonic"
	sonicEncoder "github.com/bytedance/sonic/encoder"
)

type numberInterface interface {
	number() int
}

type testNumber int

func (n testNumber) number() int {
	return int(n)
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	intValue := 42
	intPointer := &intValue
	mapValue := map[string]int{"one": 1, "two": 2}
	mapPointer := &mapValue

	tests := []struct {
		name  string
		value any
	}{
		{"nil", nil},
		{"bool", true},
		{"int", int(-1)},
		{"int8", int8(-8)},
		{"int16", int16(-16)},
		{"int32", int32(-32)},
		{"int64", int64(math.MinInt64)},
		{"uint", uint(1)},
		{"uint8", uint8(8)},
		{"uint16", uint16(16)},
		{"uint32", uint32(32)},
		{"uint64", uint64(math.MaxUint64)},
		{"uintptr", uintptr(64)},
		{"float32", float32(1.25)},
		{"float64", float64(1e-7)},
		{"string", "quote: \" slash: \\ newline:\n unicode: 世界 <>&"},
		{"pointer", intPointer},
		{"pointer chain", &intPointer},
		{"nil pointer", (*int)(nil)},
		{"array", [3]int{-1, 0, 1}},
		{"slice", []string{"one", "two"}},
		{"non-empty interface slice", []numberInterface{testNumber(1), testNumber(2)}},
		{"nil slice", []int(nil)},
		{"empty slice", []int{}},
		{"string string map", map[string]string{"one": "first", "two": "second"}},
		{"string int map", map[string]int{"negative": -1, "positive": 1}},
		{"int string map", map[int]string{-1: "negative", 1: "positive"}},
		{"uint bool map", map[uint]bool{1: true, 2: false}},
		{"string any map", map[string]any{"bool": true, "int": 1, "slice": []int{1, 2}}},
		{"int any map", map[int]any{1: "one", 2: []int{2, 3}}},
		{"uint any map", map[uint]any{1: true, 2: "two"}},
		{"composite map value", map[string][]int{"numbers": {1, 2, 3}}},
		{"named int slice map", map[string][]testNumber{"numbers": {1, 2, 3}}},
		{"bool slice map", map[string][]bool{"values": {true, false}}},
		{"float32 slice map", map[string][]float32{"values": {-1.25, 1e-7}}},
		{"float64 slice map", map[string][]float64{"values": {-1.25, 1e-7}}},
		{"string slice map", map[string][]string{"values": {"one", "two"}}},
		{"int bool slice map", map[int][]bool{-1: {true, false}}},
		{"int int slice map", map[int][]int16{-1: {-2, 3}}},
		{"int uint slice map", map[int][]uint32{-1: {2, 3}}},
		{"int float32 slice map", map[int][]float32{-1: {-1.25, 1e-7}}},
		{"int float64 slice map", map[int][]float64{-1: {-1.25, 1e-7}}},
		{"int string slice map", map[int][]string{-1: {"one", "two"}}},
		{"uint bool slice map", map[uint][]bool{1: {true, false}}},
		{"uint int slice map", map[uint][]int16{1: {-2, 3}}},
		{"uint uint slice map", map[uint][]uint32{1: {2, 3}}},
		{"uint float32 slice map", map[uint][]float32{1: {-1.25, 1e-7}}},
		{"uint float64 slice map", map[uint][]float64{1: {-1.25, 1e-7}}},
		{"uint string slice map", map[uint][]string{1: {"one", "two"}}},
		{"pointer to map", mapPointer},
		{"pointer chain to map", &mapPointer},
		{"slice of maps", []map[string]int{{"one": 1}, nil, {"two": 2}}},
		{"array of maps", [2]map[string]int{{"one": 1}, {"two": 2}}},
		{"nested map", map[string]map[string]int{"outer": {"inner": 1}}},
		{"nil map", map[string]int(nil)},
		{"empty map", map[string]int{}},
		{"empty struct", struct{}{}},
		{"struct", struct {
			Bool    bool              `json:"bool"`
			Int     int               `json:"integer"`
			Uint    uint              `json:"unsigned"`
			Float32 float32           `json:"float32"`
			Float64 float64           `json:"float64"`
			String  string            `json:"string"`
			Pointer *int              `json:"pointer"`
			Slice   []int             `json:"slice"`
			Map     map[string]string `json:"map"`
			Any     any               `json:"any"`
			Ignored string            `json:"-"`
			hidden  string
		}{
			Bool:    true,
			Int:     -1,
			Uint:    1,
			Float32: 1.25,
			Float64: 1e-7,
			String:  "hello",
			Pointer: intPointer,
			Slice:   []int{-1, 0, 1},
			Map:     map[string]string{"key": "value"},
			Any:     []string{"one", "two"},
			Ignored: "ignored",
			hidden:  "hidden",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := jsonexperiment.Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal(%T): %v", tt.value, err)
			}

			want, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("json.Marshal(%T): %v", tt.value, err)
			}

			assertJSONEqual(t, got, want)
		})
	}
}

func TestMarshalAppend(t *testing.T) {
	t.Parallel()

	dst := make([]byte, 0, 64)
	dst = append(dst, "prefix:"...)

	got, err := jsonexperiment.MarshalAppend(dst, []int{1, 2, 3})
	if err != nil {
		t.Fatalf("MarshalAppend: %v", err)
	}

	if want := []byte("prefix:[1,2,3]"); !bytes.Equal(got, want) {
		t.Fatalf("MarshalAppend = %q, want %q", got, want)
	}
}

func TestMarshalError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{"float32 NaN", float32(math.NaN())},
		{"float64 positive infinity", math.Inf(1)},
		{"float64 negative infinity", math.Inf(-1)},
		{"int float slice infinity", map[int][]float64{1: {math.Inf(1)}}},
		{"uint float slice NaN", map[uint][]float32{1: {float32(math.NaN())}}},
		{"unsupported channel", make(chan int)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := jsonexperiment.Marshal(tt.value); err == nil {
				t.Fatalf("Marshal(%T) error = nil", tt.value)
			}
		})
	}
}

func TestMarshalReturnsOwnedBytes(t *testing.T) {
	first, err := jsonexperiment.Marshal("first")
	if err != nil {
		t.Fatalf("Marshal(first): %v", err)
	}

	if _, err := jsonexperiment.Marshal("second"); err != nil {
		t.Fatalf("Marshal(second): %v", err)
	}

	if want := []byte(`"first"`); !bytes.Equal(first, want) {
		t.Fatalf("first result changed to %q, want %q", first, want)
	}
}

func assertJSONEqual(t *testing.T, got, want []byte) {
	t.Helper()

	var gotValue, wantValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("invalid JSON %q: %v", got, err)
	}
	if err := json.Unmarshal(want, &wantValue); err != nil {
		t.Fatalf("invalid reference JSON %q: %v", want, err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("got %q, want JSON equivalent to %q", got, want)
	}
}

var sonicJson = sonic.ConfigFastest

func BenchmarkMarshalMapInt(b *testing.B) {
	var marshalResult []byte

	values := map[string]int{
		"minimum":   math.MinInt,
		"negative":  -1_000_000,
		"minus_one": -1,
		"zero":      0,
		"one":       1,
		"positive":  1_000_000,
		"maximum":   math.MaxInt,
	}

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], values)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, values, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}

func BenchmarkMarshalMapIntSlice(b *testing.B) {
	var marshalResult []byte

	values := map[string][]int{
		"negative": {math.MinInt, -1_000_000, -1},
		"zero":     {0},
		"positive": {1, 1_000_000, math.MaxInt},
		"mixed":    {-2, -1, 0, 1, 2},
		"empty":    {},
		"nil":      nil,
	}

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], values)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, values, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}

func BenchmarkMarshalMapAny(b *testing.B) {
	var marshalResult []byte

	values := map[string]any{
		"bool":    true,
		"int":     int(math.MinInt),
		"uint":    uint(math.MaxUint),
		"float32": float32(1.25),
		"float64": 1e-7,
		"string":  "hello, world",
		"slice":   []int{-1, 0, 1},
		"nil":     nil,
	}

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], values)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, values, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}

func BenchmarkMarshalIntSlice(b *testing.B) {
	var marshalResult []byte

	values := []int{
		math.MinInt,
		-1_000_000,
		-1,
		0,
		1,
		1_000_000,
		math.MaxInt,
	}

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], values)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(values)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, values, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}

func BenchmarkMarshalFloat32(b *testing.B) {
	var marshalResult []byte
	value := float32(1.234567)

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], value)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, value, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}

func BenchmarkMarshalFloat64(b *testing.B) {
	var marshalResult []byte
	value := 1.2345678901234567

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], value)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, value, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}

func BenchmarkMarshalStruct(b *testing.B) {
	var marshalResult []byte

	value := struct {
		ID       uint64            `json:"id"`
		Name     string            `json:"name"`
		Active   bool              `json:"active"`
		Score    float64           `json:"score"`
		Count    int               `json:"count"`
		Tags     []string          `json:"tags"`
		Metadata map[string]string `json:"metadata"`
	}{
		ID:     math.MaxUint64,
		Name:   "benchmark value",
		Active: true,
		Score:  1.2345678901234567,
		Count:  math.MaxInt,
		Tags:   []string{"json", "benchmark", "performance"},
		Metadata: map[string]string{
			"environment": "production",
			"region":      "eu-north-1",
		},
	}

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], value)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, value, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}

func BenchmarkMarshalStructSlice(b *testing.B) {
	var marshalResult []byte

	type itemMetadata map[string]string

	type item struct {
		ID       uint64       `json:"id"`
		Name     string       `json:"name"`
		Quantity int          `json:"quantity"`
		Price    float64      `json:"price"`
		Active   bool         `json:"active"`
		Metadata itemMetadata `json:"metadata"`
	}

	value := struct {
		OrderID uint64 `json:"order_id"`
		Owner   string `json:"owner"`
		Items   []item `json:"items"`
	}{
		OrderID: math.MaxUint64,
		Owner:   "benchmark owner",
		Items: []item{
			{ID: 1, Name: "first item", Quantity: 1, Price: 1.1, Active: true, Metadata: itemMetadata{"color": "red", "size": "M"}},
			{ID: 2, Name: "second item", Quantity: 10, Price: 22.22, Active: true, Metadata: itemMetadata{"color": "blue", "size": "L"}},
			{ID: 3, Name: "third item", Quantity: math.MaxInt, Price: 333.333, Active: false, Metadata: itemMetadata{"color": "green", "size": "S"}},
			{ID: math.MaxUint64, Name: "final item", Quantity: -1, Price: 4444.4444, Active: true, Metadata: itemMetadata{"color": "black", "size": "XL"}},
		},
	}

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], value)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, value, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}
