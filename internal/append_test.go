package internal_test

import (
	"bytes"
	"encoding/json"
	"math"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/33TU/json-experiment/internal"
)

func TestAppendPrimitive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  []byte
		want string
	}{
		{"null", internal.AppendNull(nil), "null"},
		{"bool true", internal.AppendBool(nil, true), "true"},
		{"bool false", internal.AppendBool(nil, false), "false"},
		{"int", internal.AppendInt(nil, int64(math.MinInt64)), "-9223372036854775808"},
		{"uint", internal.AppendUint(nil, uint64(math.MaxUint64)), "18446744073709551615"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if string(tt.got) != tt.want {
				t.Fatalf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestAppendFloat(t *testing.T) {
	t.Parallel()

	t.Run("float32", func(t *testing.T) {
		values := []float32{0, float32(math.Copysign(0, -1)), 1.25, 1e-7, 1e20}
		for _, value := range values {
			got, err := internal.AppendFloat32(nil, value)
			if err != nil {
				t.Fatalf("AppendFloat32(%v): %v", value, err)
			}
			want, err := json.Marshal(value)
			if err != nil {
				t.Fatalf("json.Marshal(%v): %v", value, err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("AppendFloat32(%v) = %q, want %q", value, got, want)
			}
		}
	})

	t.Run("float64", func(t *testing.T) {
		values := []float64{0, math.Copysign(0, -1), 1.25, 1e-7, 1e20}
		for _, value := range values {
			got, err := internal.AppendFloat64(nil, value)
			if err != nil {
				t.Fatalf("AppendFloat64(%v): %v", value, err)
			}
			want, err := json.Marshal(value)
			if err != nil {
				t.Fatalf("json.Marshal(%v): %v", value, err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("AppendFloat64(%v) = %q, want %q", value, got, want)
			}
		}
	})

	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := internal.AppendFloat64(nil, value); err == nil {
			t.Fatalf("AppendFloat64(%v) error = nil", value)
		}
	}
}

func TestAppendString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"empty", "", `""`},
		{"plain", "hello world", `"hello world"`},
		{"quote and slash", `hello "\world`, `"hello \"\\world"`},
		{"controls", "\b\f\n\r\t\x00", `"\b\f\n\r\t\u0000"`},
		{"unicode", "Hello, 世界", `"Hello, 世界"`},
		{"no HTML escaping", "<>&", `"<>&"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := internal.AppendString([]byte("prefix:"), tt.value)
			want := "prefix:" + tt.want
			if string(got) != want {
				t.Fatalf("AppendString(%q) = %q, want %q", tt.value, got, want)
			}
		})
	}
}

func TestAppendStringHTML(t *testing.T) {
	t.Parallel()

	values := []string{"hello world", `<script>&"`, "\b\f\n\r\t\x00"}
	for _, value := range values {
		got := internal.AppendStringHTML(nil, value)
		want, err := json.Marshal(value)
		if err != nil {
			t.Fatalf("json.Marshal(%q): %v", value, err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("AppendStringHTML(%q) = %q, want %q", value, got, want)
		}
	}
}

func TestAppendQuotedString(t *testing.T) {
	t.Parallel()

	values := []string{"", "hello world", `quote: " slash: \\`, "\b\f\n\r\t\x00", "Hello, 世界", "<>&"}
	for _, value := range values {
		inner := internal.AppendString(nil, value)

		var want bytes.Buffer
		encoder := json.NewEncoder(&want)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(string(inner)); err != nil {
			t.Fatalf("encoding quoted string %q: %v", value, err)
		}

		got := internal.AppendQuotedString(nil, value)
		if !bytes.Equal(got, bytes.TrimSuffix(want.Bytes(), []byte{'\n'})) {
			t.Fatalf("AppendQuotedString(%q) = %q, want %q", value, got, want.Bytes())
		}
	}
}

func TestAppendQuotedStringHTML(t *testing.T) {
	t.Parallel()

	values := []string{"hello world", `<script>&"`, "\b\f\n\r\t\x00"}
	for _, value := range values {
		inner := internal.AppendStringHTML(nil, value)
		want, err := json.Marshal(string(inner))
		if err != nil {
			t.Fatalf("encoding HTML-safe quoted string %q: %v", value, err)
		}

		got := internal.AppendQuotedStringHTML(nil, value)
		if !bytes.Equal(got, want) {
			t.Fatalf("AppendQuotedStringHTML(%q) = %q, want %q", value, got, want)
		}
	}
}

func TestAppendSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		got   []byte
		value any
	}{
		{"bool", internal.AppendBoolSlice(nil, []bool{true, false}), []bool{true, false}},
		{"int", internal.AppendIntSlice(nil, []int{-1, 0, 1}), []int{-1, 0, 1}},
		{"uint", internal.AppendUintSlice(nil, []uint{0, 1, 2}), []uint{0, 1, 2}},
		{"string", internal.AppendStringSlice(nil, []string{"one", "two"}), []string{"one", "two"}},
		{"nil", internal.AppendIntSlice[int](nil, nil), []int(nil)},
		{"empty", internal.AppendIntSlice(nil, []int{}), []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			want, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("json.Marshal(%T): %v", tt.value, err)
			}
			if !bytes.Equal(tt.got, want) {
				t.Fatalf("got %q, want %q", tt.got, want)
			}
		})
	}
}

func TestAppendMap(t *testing.T) {
	t.Parallel()

	float32Values := map[string]float32{"negative": -1.25, "small": 1e-7}
	float32Result, float32Err := internal.AppendStringFloat32Map(nil, float32Values)

	float64Values := map[string]float64{"negative": -1.25, "small": 1e-7}
	float64Result, float64Err := internal.AppendStringFloat64Map(nil, float64Values)

	intFloat32Values := map[int8]float32{-1: -1.25, 2: 1e-7}
	intFloat32Result, intFloat32Err := internal.AppendIntFloat32Map(nil, intFloat32Values)

	intFloat64Values := map[int16]float64{-1: -1.25, 2: 1e-7}
	intFloat64Result, intFloat64Err := internal.AppendIntFloat64Map(nil, intFloat64Values)

	uintFloat32Values := map[uint8]float32{1: -1.25, 2: 1e-7}
	uintFloat32Result, uintFloat32Err := internal.AppendUintFloat32Map(nil, uintFloat32Values)

	uintFloat64Values := map[uint16]float64{1: -1.25, 2: 1e-7}
	uintFloat64Result, uintFloat64Err := internal.AppendUintFloat64Map(nil, uintFloat64Values)

	float32SliceValues := map[string][]float32{"values": {-1.25, 1e-7}}
	float32SliceResult, float32SliceErr := internal.AppendStringFloat32SliceMap(nil, float32SliceValues)

	float64SliceValues := map[string][]float64{"values": {-1.25, 1e-7}}
	float64SliceResult, float64SliceErr := internal.AppendStringFloat64SliceMap(nil, float64SliceValues)

	tests := []struct {
		name  string
		got   []byte
		err   error
		value any
	}{
		{"bool", internal.AppendStringBoolMap(nil, map[string]bool{"true": true, "false": false}), nil, map[string]bool{"true": true, "false": false}},
		{"int", internal.AppendStringIntMap(nil, map[string]int{"negative": -1, "zero": 0, "positive": 1}), nil, map[string]int{"negative": -1, "zero": 0, "positive": 1}},
		{"uint", internal.AppendStringUintMap(nil, map[string]uint{"zero": 0, "positive": 1}), nil, map[string]uint{"zero": 0, "positive": 1}},
		{"float32", float32Result, float32Err, float32Values},
		{"float64", float64Result, float64Err, float64Values},
		{"string", internal.AppendStringStringMap(nil, map[string]string{`key"`: "line\nvalue"}), nil, map[string]string{`key"`: "line\nvalue"}},
		{"bool slice", internal.AppendStringBoolSliceMap(nil, map[string][]bool{"values": {true, false}}), nil, map[string][]bool{"values": {true, false}}},
		{"int slice", internal.AppendStringIntSliceMap(nil, map[string][]int16{"values": {-1, 0, 1}}), nil, map[string][]int16{"values": {-1, 0, 1}}},
		{"uint slice", internal.AppendStringUintSliceMap(nil, map[string][]uint32{"values": {0, 1, 2}}), nil, map[string][]uint32{"values": {0, 1, 2}}},
		{"float32 slice", float32SliceResult, float32SliceErr, float32SliceValues},
		{"float64 slice", float64SliceResult, float64SliceErr, float64SliceValues},
		{"string slice", internal.AppendStringStringSliceMap(nil, map[string][]string{"values": {"one", "two"}}), nil, map[string][]string{"values": {"one", "two"}}},
		{"nil int slice map", internal.AppendStringIntSliceMap[int](nil, nil), nil, map[string][]int(nil)},
		{"empty uint slice map", internal.AppendStringUintSliceMap(nil, map[string][]uint{}), nil, map[string][]uint{}},
		{"nil", internal.AppendStringIntMap[int](nil, nil), nil, map[string]int(nil)},
		{"empty", internal.AppendStringIntMap(nil, map[string]int{}), nil, map[string]int{}},
		{"int key bool", internal.AppendIntBoolMap(nil, map[int8]bool{-1: true, 2: false}), nil, map[int8]bool{-1: true, 2: false}},
		{"int key int", internal.AppendIntIntMap(nil, map[int16]int32{-1: -2, 3: 4}), nil, map[int16]int32{-1: -2, 3: 4}},
		{"int key uint", internal.AppendIntUintMap(nil, map[int32]uint64{-1: 2, 3: 4}), nil, map[int32]uint64{-1: 2, 3: 4}},
		{"int key float32", intFloat32Result, intFloat32Err, intFloat32Values},
		{"int key float64", intFloat64Result, intFloat64Err, intFloat64Values},
		{"int key string", internal.AppendIntStringMap(nil, map[int64]string{-1: "negative", 2: "positive"}), nil, map[int64]string{-1: "negative", 2: "positive"}},
		{"uint key bool", internal.AppendUintBoolMap(nil, map[uint8]bool{1: true, 2: false}), nil, map[uint8]bool{1: true, 2: false}},
		{"uint key int", internal.AppendUintIntMap(nil, map[uint16]int32{1: -2, 3: 4}), nil, map[uint16]int32{1: -2, 3: 4}},
		{"uint key uint", internal.AppendUintUintMap(nil, map[uint32]uint64{1: 2, 3: 4}), nil, map[uint32]uint64{1: 2, 3: 4}},
		{"uint key float32", uintFloat32Result, uintFloat32Err, uintFloat32Values},
		{"uint key float64", uintFloat64Result, uintFloat64Err, uintFloat64Values},
		{"uint key string", internal.AppendUintStringMap(nil, map[uint64]string{1: "one", 2: "two"}), nil, map[uint64]string{1: "one", 2: "two"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if tt.err != nil {
				t.Fatalf("append map: %v", tt.err)
			}

			want, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("json.Marshal(%T): %v", tt.value, err)
			}

			var gotValue, wantValue any
			if err := json.Unmarshal(tt.got, &gotValue); err != nil {
				t.Fatalf("json.Unmarshal(got %q): %v", tt.got, err)
			}
			if err := json.Unmarshal(want, &wantValue); err != nil {
				t.Fatalf("json.Unmarshal(want %q): %v", want, err)
			}
			if !reflect.DeepEqual(gotValue, wantValue) {
				t.Fatalf("got %q, want JSON equivalent to %q", tt.got, want)
			}
		})
	}

	if _, err := internal.AppendStringFloat32Map(nil, map[string]float32{"invalid": float32(math.NaN())}); err == nil {
		t.Fatal("AppendStringFloat32Map with NaN returned nil error")
	}
	if _, err := internal.AppendStringFloat64Map(nil, map[string]float64{"invalid": math.Inf(1)}); err == nil {
		t.Fatal("AppendStringFloat64Map with infinity returned nil error")
	}
	if _, err := internal.AppendStringFloat32SliceMap(nil, map[string][]float32{"invalid": {float32(math.NaN())}}); err == nil {
		t.Fatal("AppendStringFloat32SliceMap with NaN returned nil error")
	}
	if _, err := internal.AppendStringFloat64SliceMap(nil, map[string][]float64{"invalid": {math.Inf(1)}}); err == nil {
		t.Fatal("AppendStringFloat64SliceMap with infinity returned nil error")
	}
}

var appendResult []byte

func BenchmarkAppendString(b *testing.B) {
	benchmarks := []struct {
		name  string
		value string
		html  bool
	}{
		{"plain", "hello world", false},
		{"plain_clean_huge", strings.Repeat("The quick brown fox jumps. ", 64), false},
		{"escaped", "hello \"world\"\npath\\to\\file", false},
		{"escape_at_end", strings.Repeat("a", 1024) + `"`, false},
		{"escape_at_middle", strings.Repeat("a", 512) + `"` + strings.Repeat("a", 512), false},
		{"HTML", `<script type="text/javascript">a && b</script>`, true},
		{"HTML_clean_huge", strings.Repeat("plain text without HTML characters ", 32), true},
		{"HTML_escape_at_end", strings.Repeat("a", 1024) + "<", true},
		{"HTML_escape_at_middle", strings.Repeat("a", 512) + "<" + strings.Repeat("a", 512), true},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.Run("internal", func(b *testing.B) {
				dst := make([]byte, 0, len(bm.value)*2+16)
				var result []byte
				b.ReportAllocs()
				for b.Loop() {
					if bm.html {
						result = internal.AppendStringHTML(dst[:0], bm.value)
					} else {
						result = internal.AppendString(dst[:0], bm.value)
					}
				}
				appendResult = result
			})

			b.Run("encoding_json", func(b *testing.B) {
				var result []byte
				b.ReportAllocs()
				for b.Loop() {
					result, _ = json.Marshal(bm.value)
				}
				appendResult = result
			})
		})
	}

	runtime.KeepAlive(appendResult)
}

func BenchmarkAppendInt(b *testing.B) {
	const value = int64(math.MinInt64)

	b.Run("internal", func(b *testing.B) {
		dst := make([]byte, 0, 32)
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = internal.AppendInt(dst[:0], value)
		}
		appendResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		appendResult = result
	})

	runtime.KeepAlive(appendResult)
}

func BenchmarkAppendUint(b *testing.B) {
	const value = uint64(math.MaxUint64)

	b.Run("internal", func(b *testing.B) {
		dst := make([]byte, 0, 32)
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = internal.AppendUint(dst[:0], value)
		}
		appendResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		appendResult = result
	})

	runtime.KeepAlive(appendResult)
}

func BenchmarkAppendFloat32(b *testing.B) {
	const value = float32(math.MaxFloat32)

	b.Run("internal", func(b *testing.B) {
		dst := make([]byte, 0, 32)
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = internal.AppendFloat32(dst[:0], value)
		}
		appendResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		appendResult = result
	})

	runtime.KeepAlive(appendResult)
}

func BenchmarkAppendFloat64(b *testing.B) {
	const value = float64(math.MaxFloat64)

	b.Run("internal", func(b *testing.B) {
		dst := make([]byte, 0, 32)
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = internal.AppendFloat64(dst[:0], value)
		}
		appendResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		appendResult = result
	})

	runtime.KeepAlive(appendResult)
}

func BenchmarkAppendSlice(b *testing.B) {
	values := []int64{
		math.MinInt64,
		-1_000_000,
		-1,
		0,
		1,
		1_000_000,
		math.MaxInt64,
	}

	b.Run("internal", func(b *testing.B) {
		dst := make([]byte, 0, 128)
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = internal.AppendIntSlice(dst[:0], values)
		}
		appendResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(values)
		}
		appendResult = result
	})

	runtime.KeepAlive(appendResult)
}

func BenchmarkAppendMap(b *testing.B) {
	values := map[string]int64{
		"minimum":   math.MinInt64,
		"negative":  -1_000_000,
		"minus_one": -1,
		"zero":      0,
		"one":       1,
		"positive":  1_000_000,
		"maximum":   math.MaxInt64,
	}

	b.Run("internal", func(b *testing.B) {
		dst := make([]byte, 0, 192)
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result = internal.AppendStringIntMap(dst[:0], values)
		}
		appendResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(values)
		}
		appendResult = result
	})

	runtime.KeepAlive(appendResult)
}
