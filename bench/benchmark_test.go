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

func BenchmarkMarshalLargeStruct(b *testing.B) {
	type address struct {
		Line1      string  `json:"line1"`
		Line2      string  `json:"line2,omitempty"`
		City       string  `json:"city"`
		Region     string  `json:"region"`
		PostalCode string  `json:"postal_code"`
		Country    string  `json:"country"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
	}

	type userProfile struct {
		ID          uint64            `json:"id"`
		Username    string            `json:"username"`
		DisplayName string            `json:"display_name"`
		Email       string            `json:"email"`
		Verified    bool              `json:"verified"`
		Roles       []string          `json:"roles"`
		Preferences map[string]string `json:"preferences"`
		Address     address           `json:"address"`
	}

	type lineItem struct {
		ID          uint64            `json:"id"`
		SKU         string            `json:"sku"`
		Name        string            `json:"name"`
		Description string            `json:"description"`
		Quantity    int               `json:"quantity"`
		UnitPrice   float64           `json:"unit_price"`
		Discount    float64           `json:"discount,omitempty"`
		Taxable     bool              `json:"taxable"`
		Tags        []string          `json:"tags"`
		Attributes  map[string]string `json:"attributes"`
	}

	type auditEvent struct {
		ID        uint64            `json:"id"`
		Timestamp string            `json:"timestamp"`
		Actor     string            `json:"actor"`
		Action    string            `json:"action"`
		Success   bool              `json:"success"`
		Details   map[string]string `json:"details"`
	}

	itemNames := []string{
		"Workstation", "Monitor", "Keyboard", "Mouse",
		"Docking station", "Headset", "Webcam", "Security key",
	}
	items := make([]lineItem, len(itemNames))
	for i, name := range itemNames {
		items[i] = lineItem{
			ID:          uint64(i + 1),
			SKU:         "BENCH-LARGE-STRUCT-SKU",
			Name:        name,
			Description: "A realistic line item with enough text to exercise string scanning and escaping.",
			Quantity:    i + 1,
			UnitPrice:   99.95 + float64(i)*125.50,
			Discount:    float64(i%3) * 2.5,
			Taxable:     i%2 == 0,
			Tags:        []string{"hardware", "benchmark", "priority"},
			Attributes: map[string]string{
				"color":     "graphite",
				"warehouse": "eu-north-1",
				"warranty":  "three years",
			},
		}
	}

	value := struct {
		RequestID      string            `json:"request_id"`
		Sequence       uint64            `json:"sequence"`
		GeneratedAt    string            `json:"generated_at"`
		Environment    string            `json:"environment"`
		Region         string            `json:"region"`
		Success        bool              `json:"success"`
		Owner          userProfile       `json:"owner"`
		BillingAddress address           `json:"billing_address"`
		Items          []lineItem        `json:"items"`
		Events         []auditEvent      `json:"events"`
		Labels         []string          `json:"labels"`
		Features       map[string]bool   `json:"features"`
		Counters       map[string]int64  `json:"counters"`
		Metadata       map[string]string `json:"metadata"`
		Warnings       []string          `json:"warnings,omitempty"`
		Checksum       string            `json:"checksum"`
	}{
		RequestID:   "req_01J5BENCHMARKLARGESTRUCT",
		Sequence:    math.MaxUint64,
		GeneratedAt: "2026-08-14T12:34:56.789123Z",
		Environment: "production",
		Region:      "eu-north-1",
		Success:     true,
		Owner: userProfile{
			ID:          math.MaxUint64,
			Username:    "benchmark-owner",
			DisplayName: "Benchmark Owner 世界",
			Email:       "benchmark.owner@example.com",
			Verified:    true,
			Roles:       []string{"administrator", "billing", "auditor"},
			Preferences: map[string]string{
				"language": "en-FI",
				"theme":    "dark",
				"timezone": "Europe/Helsinki",
			},
			Address: address{
				Line1: "123 Benchmark Avenue", City: "Helsinki", Region: "Uusimaa",
				PostalCode: "00100", Country: "FI", Latitude: 60.1699, Longitude: 24.9384,
			},
		},
		BillingAddress: address{
			Line1: "456 Allocation-Free Street", Line2: "Suite 32", City: "Espoo",
			Region: "Uusimaa", PostalCode: "02100", Country: "FI", Latitude: 60.2055, Longitude: 24.6559,
		},
		Items: items,
		Events: []auditEvent{
			{ID: 1, Timestamp: "2026-08-14T12:30:00Z", Actor: "benchmark-owner", Action: "order.created", Success: true, Details: map[string]string{"source": "api", "version": "v4"}},
			{ID: 2, Timestamp: "2026-08-14T12:31:00Z", Actor: "inventory-service", Action: "inventory.reserved", Success: true, Details: map[string]string{"warehouse": "eu-north-1", "items": "8"}},
			{ID: 3, Timestamp: "2026-08-14T12:32:00Z", Actor: "payment-service", Action: "payment.authorized", Success: true, Details: map[string]string{"currency": "EUR", "provider": "benchmark-pay"}},
			{ID: 4, Timestamp: "2026-08-14T12:33:00Z", Actor: "fulfillment-service", Action: "shipment.queued", Success: true, Details: map[string]string{"priority": "express", "carrier": "benchmark-logistics"}},
		},
		Labels: []string{"large-struct", "performance", "json", "v4", "production"},
		Features: map[string]bool{
			"audit_log": true, "discounts": true, "international_shipping": true,
		},
		Counters: map[string]int64{
			"attempts": 1, "items": 8, "notifications": 3, "retries": 0,
		},
		Metadata: map[string]string{
			"client": "json-experiment", "commit": "benchmark-large-struct",
			"host": "desktop-9950x3d", "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
		},
		Warnings: []string{"Synthetic benchmark payload; do not use as a production order."},
		Checksum: "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
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
			reference, err := json.Marshal(tt.value)
			if err != nil {
				b.Fatal(err)
			}
			encodedBytes := int64(len(reference))

			b.Run("marshal_append", func(b *testing.B) {
				var result []byte
				b.SetBytes(encodedBytes)
				b.ReportAllocs()
				for b.Loop() {
					result, _ = jsonexperiment.MarshalAppendWithFlags(result[:0], tt.value, flags)
				}
				marshalResult = result
			})

			b.Run("marshal", func(b *testing.B) {
				var result []byte
				b.SetBytes(encodedBytes)
				b.ReportAllocs()
				for b.Loop() {
					result, _ = jsonexperiment.MarshalWithFlags(tt.value, flags)
				}
				marshalResult = result
			})

			b.Run("encoding_json", func(b *testing.B) {
				var result []byte
				b.SetBytes(encodedBytes)
				b.ReportAllocs()
				for b.Loop() {
					result, _ = json.Marshal(tt.value)
				}
				marshalResult = result
			})

			b.Run("encoding_json_v2_write", func(b *testing.B) {
				buf := bytes.NewBuffer(nil)
				allowInvalidUTF8 := jsontext.AllowInvalidUTF8(true)

				b.SetBytes(encodedBytes)
				b.ReportAllocs()
				for b.Loop() {
					buf.Reset()
					_ = jsonv2.MarshalWrite(buf, tt.value, allowInvalidUTF8)
				}

				runtime.KeepAlive(buf.Bytes())
			})

			b.Run("sonic_json", func(b *testing.B) {
				var result []byte
				b.SetBytes(encodedBytes)
				b.ReportAllocs()
				for b.Loop() {
					result, _ = sonic.ConfigStd.Marshal(tt.value)
				}
				marshalResult = result
			})

			b.Run("sonic_encode_into", func(b *testing.B) {
				var result []byte
				b.SetBytes(encodedBytes)
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
	appendResult, err := jsonexperiment.MarshalAppend(nil, value)
	if err != nil {
		b.Fatal(err)
	}
	ownedResult, err := jsonexperiment.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}
	ownedBytes := int64(len(ownedResult))

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.SetBytes(int64(len(appendResult)))
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], value)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.SetBytes(ownedBytes)
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
	reference, err := json.Marshal(value)
	if err != nil {
		b.Fatal(err)
	}
	encodedBytes := int64(len(reference))

	b.Run("marshal_append", func(b *testing.B) {
		var result []byte
		b.SetBytes(encodedBytes)
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.MarshalAppend(result[:0], value)
		}
		marshalResult = result
	})

	b.Run("marshal", func(b *testing.B) {
		var result []byte
		b.SetBytes(encodedBytes)
		b.ReportAllocs()
		for b.Loop() {
			result, _ = jsonexperiment.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("encoding_json", func(b *testing.B) {
		var result []byte
		b.SetBytes(encodedBytes)
		b.ReportAllocs()
		for b.Loop() {
			result, _ = json.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("encoding_json_v2_write", func(b *testing.B) {
		buf := bytes.NewBuffer(nil)

		b.SetBytes(encodedBytes)
		b.ReportAllocs()
		for b.Loop() {
			buf.Reset()
			_ = jsonv2.MarshalWrite(buf, value)
		}

		runtime.KeepAlive(buf.Bytes())
	})

	b.Run("sonic_json", func(b *testing.B) {
		var result []byte
		b.SetBytes(encodedBytes)
		b.ReportAllocs()
		for b.Loop() {
			result, _ = sonicJson.Marshal(value)
		}
		marshalResult = result
	})

	b.Run("sonic_encode_into", func(b *testing.B) {
		var result []byte
		b.SetBytes(encodedBytes)
		b.ReportAllocs()
		for b.Loop() {
			result = result[:0]
			_ = sonicEncoder.EncodeInto(&result, value, 0)
		}
		marshalResult = result
	})

	runtime.KeepAlive(marshalResult)
}
