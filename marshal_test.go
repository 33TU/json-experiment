package jsonexperiment_test

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"runtime"
	"strings"
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

func TestMarshalStringTag(t *testing.T) {
	t.Parallel()

	value := struct {
		Bool    bool    `json:"bool,string"`
		Int     int64   `json:"int,string"`
		Uint    uint64  `json:"uint,string"`
		Float32 float32 `json:"float32,string"`
		Float64 float64 `json:"float64,string"`
		String  string  `json:"string,string"`
	}{
		Bool:    true,
		Int:     math.MinInt64,
		Uint:    math.MaxUint64,
		Float32: 1.25,
		Float64: 1e-7,
		String:  "quote: \" slash: \\ newline:\n",
	}

	got, err := jsonexperiment.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	want, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	assertJSONEqual(t, got, want)
}

func TestMarshalOmitTags(t *testing.T) {
	t.Parallel()

	var innerPointer *int
	outerPointer := &innerPointer

	type nested struct {
		Value int `json:"value"`
	}
	type value struct {
		Bool              bool           `json:"bool,omitempty"`
		Int               int            `json:"int,omitempty"`
		String            string         `json:"string,omitempty"`
		Slice             []int          `json:"slice,omitempty"`
		Map               map[string]int `json:"map,omitempty"`
		Pointer           *int           `json:"pointer,omitempty"`
		NilDoublePointer  **int          `json:"nil_double_pointer,omitempty"`
		KeepDoublePointer **int          `json:"keep_double_pointer,omitempty"`
		ZeroStruct        nested         `json:"zero_struct,omitzero"`
		KeepStruct        nested         `json:"keep_struct,omitempty"`
		KeepSlice         []int          `json:"keep_slice,omitzero"`
	}
	testValue := value{
		Slice:             []int{},
		Map:               map[string]int{},
		KeepDoublePointer: outerPointer,
		KeepStruct:        nested{},
		KeepSlice:         []int{},
	}

	got, err := jsonexperiment.Marshal(testValue)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	want, err := json.Marshal(testValue)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	assertJSONEqual(t, got, want)
}

func TestMarshalEscapeHTML(t *testing.T) {
	t.Parallel()

	options := jsonexperiment.MarshalOptions{EscapeHTML: true}
	value := struct {
		Text   string              `json:"<text>"`
		Array  [1]string           `json:"array"`
		Slice  []string            `json:"slice"`
		Map    map[string]string   `json:"map"`
		Nested map[string][]string `json:"nested"`
		Any    any                 `json:"any"`
	}{
		Text:  "<>&",
		Array: [1]string{"<array>"},
		Slice: []string{"<slice>"},
		Map:   map[string]string{"<key>": "<value>"},
		Nested: map[string][]string{
			"<nested>": {"<item>"},
		},
		Any: map[string]any{"<any>": "<interface>"},
	}

	got, err := jsonexperiment.MarshalWithOptions(value, options)
	if err != nil {
		t.Fatalf("MarshalWithOptions: %v", err)
	}
	if bytes.ContainsAny(got, "<>&") {
		t.Fatalf("MarshalWithOptions left HTML characters unescaped: %s", got)
	}
	for _, escaped := range [][]byte{[]byte(`\u003c`), []byte(`\u003e`), []byte(`\u0026`)} {
		if !bytes.Contains(got, escaped) {
			t.Fatalf("MarshalWithOptions output %q does not contain %q", got, escaped)
		}
	}

	want, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	assertJSONEqual(t, got, want)

	html := "<pointer>"
	additional := []any{
		&html,
		map[string]int{"<key>": 1},
		map[int]string{1: "<value>"},
		map[uint]string{1: "<value>"},
		map[int][]string{1: {"<value>"}},
		map[uint][]string{1: {"<value>"}},
	}
	for _, value := range additional {
		encoded, err := jsonexperiment.MarshalWithOptions(value, options)
		if err != nil {
			t.Fatalf("MarshalWithOptions(%T): %v", value, err)
		}
		if bytes.ContainsAny(encoded, "<>&") {
			t.Fatalf("MarshalWithOptions(%T) left HTML characters unescaped: %s", value, encoded)
		}
	}

	got, err = jsonexperiment.MarshalAppendWithOptions([]byte("prefix:"), "<>&", options)
	if err != nil {
		t.Fatalf("MarshalAppendWithOptions: %v", err)
	}
	if want := []byte(`prefix:"\u003c\u003e\u0026"`); !bytes.Equal(got, want) {
		t.Fatalf("MarshalAppendWithOptions = %q, want %q", got, want)
	}

	got, err = jsonexperiment.Marshal("<>&")
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := []byte(`"<>&"`); !bytes.Equal(got, want) {
		t.Fatalf("Marshal default = %q, want %q", got, want)
	}
}

func TestMarshalValidateString(t *testing.T) {
	t.Parallel()

	const value = "before\xffafter"
	options := jsonexperiment.MarshalOptions{ValidateString: true}

	want, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	got, err := jsonexperiment.MarshalWithOptions(value, options)
	if err != nil {
		t.Fatalf("MarshalWithOptions: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("MarshalWithOptions = %q, want %q", got, want)
	}

	got, err = jsonexperiment.MarshalWithFlags(value, jsonexperiment.MarshalFlagValidateString)
	if err != nil {
		t.Fatalf("MarshalWithFlags: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("MarshalWithFlags = %q, want %q", got, want)
	}

	got, err = jsonexperiment.MarshalAppendWithOptions([]byte("prefix:\xff"), value, options)
	if err != nil {
		t.Fatalf("MarshalAppendWithOptions: %v", err)
	}
	wantAppend := append([]byte("prefix:\xff"), want...)
	if !bytes.Equal(got, wantAppend) {
		t.Fatalf("MarshalAppendWithOptions = %q, want %q", got, wantAppend)
	}

	nested := struct {
		Value string            `json:"value"`
		Map   map[string]string `json:"map"`
	}{
		Value: value,
		Map:   map[string]string{"key\xff": "value\xfe"},
	}

	got, err = jsonexperiment.MarshalWithOptions(nested, options)
	if err != nil {
		t.Fatalf("MarshalWithOptions(nested): %v", err)
	}
	want, err = json.Marshal(nested)
	if err != nil {
		t.Fatalf("json.Marshal(nested): %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("MarshalWithOptions(nested) = %q, want %q", got, want)
	}

	got, err = jsonexperiment.Marshal(value)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := []byte("\"before\xffafter\""); !bytes.Equal(got, want) {
		t.Fatalf("Marshal without validation = %q, want %q", got, want)
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
	values := map[string]int{
		"minimum":   math.MinInt,
		"negative":  -1_000_000,
		"minus_one": -1,
		"zero":      0,
		"one":       1,
		"positive":  1_000_000,
		"maximum":   math.MaxInt,
	}

	benchmarkMarshalValue(b, values)
}

func BenchmarkMarshalMapIntSlice(b *testing.B) {
	values := map[string][]int{
		"negative": {math.MinInt, -1_000_000, -1},
		"zero":     {0},
		"positive": {1, 1_000_000, math.MaxInt},
		"mixed":    {-2, -1, 0, 1, 2},
		"empty":    {},
		"nil":      nil,
	}

	benchmarkMarshalValue(b, values)
}

func BenchmarkMarshalMapAny(b *testing.B) {
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

	benchmarkMarshalValue(b, values)
}

func BenchmarkMarshalIntSlice(b *testing.B) {
	values := []int{
		math.MinInt,
		-1_000_000,
		-1,
		0,
		1,
		1_000_000,
		math.MaxInt,
	}

	benchmarkMarshalValue(b, values)
}

func BenchmarkMarshalFloat32(b *testing.B) {
	value := float32(1.234567)

	benchmarkMarshalValue(b, value)
}

func BenchmarkMarshalFloat64(b *testing.B) {
	value := 1.2345678901234567

	benchmarkMarshalValue(b, value)
}

func BenchmarkMarshalStruct(b *testing.B) {
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

	benchmarkMarshalValue(b, value)
}

func BenchmarkMarshalStructSlice(b *testing.B) {
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

	benchmarkMarshalValue(b, value)
}

func BenchmarkMarshalStructQuoted(b *testing.B) {
	value := struct {
		Bool    bool    `json:"bool,string"`
		Int     int64   `json:"int,string"`
		Uint    uint64  `json:"uint,string"`
		Float32 float32 `json:"float32,string"`
		Float64 float64 `json:"float64,string"`
		String  string  `json:"string,string"`
	}{
		Bool:    true,
		Int:     math.MinInt64,
		Uint:    math.MaxUint64,
		Float32: 1.25,
		Float64: 1e-7,
		String:  `quoted "benchmark" value`,
	}

	benchmarkMarshalValue(b, value)
}

func BenchmarkMarshalOmits(b *testing.B) {
	value := struct {
		EmptyBool   bool           `json:"empty_bool,omitempty"`
		EmptyInt    int            `json:"empty_int,omitempty"`
		EmptyString string         `json:"empty_string,omitempty"`
		EmptySlice  []int          `json:"empty_slice,omitempty"`
		EmptyMap    map[string]int `json:"empty_map,omitempty"`
		KeepInt     int            `json:"keep_int,omitempty"`
		KeepString  string         `json:"keep_string,omitempty"`
		KeepSlice   []int          `json:"keep_slice,omitempty"`
	}{
		EmptySlice: []int{},
		EmptyMap:   map[string]int{},
		KeepInt:    math.MaxInt,
		KeepString: "benchmark value",
		KeepSlice:  []int{1, 2, 3},
	}

	benchmarkMarshalValue(b, value)
}

func BenchmarkMarshalUTF8(b *testing.B) {
	options := jsonexperiment.MarshalOptions{ValidateString: true}

	values := []struct {
		name  string
		value string
	}{
		{name: "valid_ascii", value: strings.Repeat("The quick brown fox jumps. ", 32)},
		{name: "valid_unicode", value: strings.Repeat("Hello, 世界. ", 32)},
		{name: "invalid_middle", value: strings.Repeat("a", 512) + "\xff" + strings.Repeat("b", 512)},
	}

	for _, tt := range values {
		b.Run(tt.name, func(b *testing.B) {
			var marshalResult []byte

			b.Run("marshal_append", func(b *testing.B) {
				var result []byte
				b.ReportAllocs()
				for b.Loop() {
					result, _ = jsonexperiment.MarshalAppendWithOptions(result[:0], tt.value, options)
				}
				marshalResult = result
			})

			b.Run("marshal", func(b *testing.B) {
				var result []byte
				b.ReportAllocs()
				for b.Loop() {
					result, _ = jsonexperiment.MarshalWithOptions(tt.value, options)
				}
				marshalResult = result
			})

			b.Run("encoding_json", func(b *testing.B) {
				var result []byte
				b.ReportAllocs()
				for b.Loop() {
					result, _ = json.Marshal(tt.value)
				}
				marshalResult = result
			})

			b.Run("sonic_json", func(b *testing.B) {
				var result []byte
				b.ReportAllocs()
				for b.Loop() {
					result, _ = sonic.ConfigStd.Marshal(tt.value)
				}
				marshalResult = result
			})

			b.Run("sonic_encode_into", func(b *testing.B) {
				var result []byte
				b.ReportAllocs()
				for b.Loop() {
					result = result[:0]
					_ = sonicEncoder.EncodeInto(&result, tt.value, sonicEncoder.ValidateString)
				}
				marshalResult = result
			})

			runtime.KeepAlive(marshalResult)
		})
	}
}

func BenchmarkMarshalUTF8Slice(b *testing.B) {
	values := []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		{Name: "valid_ascii", Value: strings.Repeat("The quick brown fox jumps. ", 32)},
		{Name: "valid_unicode", Value: strings.Repeat("Hello, 世界. ", 32)},
		{Name: "invalid_middle", Value: strings.Repeat("a", 512) + "\xff" + strings.Repeat("b", 512)},
	}

	benchmarkMarshalValue(b, values)
}

func benchmarkMarshalValue[T any](b *testing.B, value T) {
	var marshalResult []byte

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
