package bench_test

import (
	"bytes"
	"encoding/json"
	"math"
	"runtime"
	"strings"
	"testing"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	jsonexperiment "github.com/33TU/json-experiment"
	"github.com/bytedance/sonic"
	sonicEncoder "github.com/bytedance/sonic/encoder"
)

type allMarshalers struct{}

func (allMarshalers) MarshalJSONAppend(dst []byte) ([]byte, error) {
	return append(dst, `"append"`...), nil
}

func (allMarshalers) MarshalJSON() ([]byte, error) {
	return []byte(`"json"`), nil
}

type jsonMarshaller struct{}

func (jsonMarshaller) MarshalJSON() ([]byte, error) {
	return []byte(`"json"`), nil
}

type textMarshaler struct{}

func (textMarshaler) MarshalText() ([]byte, error) {
	return []byte("<text>"), nil
}

type jsonInt int

func (jsonInt) MarshalJSON() ([]byte, error) {
	return []byte(`"custom"`), nil
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

func BenchmarkMarshalOmitZero(b *testing.B) {
	value := struct {
		ZeroBool   bool           `json:"zero_bool,omitzero"`
		ZeroInt    int            `json:"zero_int,omitzero"`
		ZeroString string         `json:"zero_string,omitzero"`
		ZeroSlice  []int          `json:"zero_slice,omitzero"`
		ZeroMap    map[string]int `json:"zero_map,omitzero"`
		KeepInt    int            `json:"keep_int,omitzero"`
		KeepString string         `json:"keep_string,omitzero"`
		KeepSlice  []int          `json:"keep_slice,omitzero"`
	}{
		KeepInt:    math.MaxInt,
		KeepString: "benchmark value",
		KeepSlice:  []int{1, 2, 3},
	}

	benchmarkMarshalValue(b, value)
}

func BenchmarkMarshalUTF8(b *testing.B) {
	flags := jsonexperiment.MarshalOptions{ValidateString: true}.Flags()

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
					result, _ = jsonexperiment.MarshalAppendWithFlags(result[:0], tt.value, flags)
				}
				marshalResult = result
			})

			b.Run("marshal", func(b *testing.B) {
				var result []byte
				b.ReportAllocs()
				for b.Loop() {
					result, _ = jsonexperiment.MarshalWithFlags(tt.value, flags)
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

			b.Run("encoding_json_v2_write", func(b *testing.B) {
				buf := bytes.NewBuffer(nil)
				allowInvalidUTF8 := jsontext.AllowInvalidUTF8(true)

				b.ReportAllocs()
				for b.Loop() {
					buf.Reset()
					_ = jsonv2.MarshalWrite(buf, tt.value, allowInvalidUTF8)
				}

				runtime.KeepAlive(buf.Bytes())
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

func BenchmarkTextMarshaler(b *testing.B) {
	value := textMarshaler{}

	benchmarkMarshalValue(b, value)
}

func BenchmarkJsonMarshaler(b *testing.B) {
	value := jsonMarshaller{}

	benchmarkMarshalValue(b, value)
}

func BenchmarkJsonIntMarshaler(b *testing.B) {
	value := jsonInt(123)

	benchmarkMarshalValue(b, value)
}

func BenchmarkAllMarshalers(b *testing.B) {
	var marshalResult []byte

	value := allMarshalers{}

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

	runtime.KeepAlive(marshalResult)
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

	b.Run("encoding_json_v2_write", func(b *testing.B) {
		buf := bytes.NewBuffer(nil)

		b.ReportAllocs()
		for b.Loop() {
			buf.Reset()
			_ = jsonv2.MarshalWrite(buf, value)
		}

		runtime.KeepAlive(buf.Bytes())
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
