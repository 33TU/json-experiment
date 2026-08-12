package jsonexperiment_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

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

type testByte byte

type allMarshalers struct{}

var (
	_ jsonexperiment.MarshalerAppend = allMarshalers{}
	_ json.Marshaler                 = allMarshalers{}
)

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

type appendInt int

func (appendInt) MarshalJSONAppend(dst []byte) ([]byte, error) {
	return append(dst, `"append"`...), nil
}

type jsonString string

func (jsonString) MarshalJSON() ([]byte, error) {
	return []byte(`"json"`), nil
}

type jsonInt int

func (jsonInt) MarshalJSON() ([]byte, error) {
	return []byte(`"custom"`), nil
}

type pointerJSONInt int

func (*pointerJSONInt) MarshalJSON() ([]byte, error) {
	return []byte(`"pointer"`), nil
}

type pointerJSONMap map[string]int

func (*pointerJSONMap) MarshalJSON() ([]byte, error) {
	return []byte(`"pointer-map"`), nil
}

type jsonIntSlice []int

func (jsonIntSlice) MarshalJSON() ([]byte, error) {
	return []byte(`"custom-slice"`), nil
}

type textMapKey int

func (textMapKey) MarshalText() ([]byte, error) {
	return []byte("text-key"), nil
}

type stringTextMapKey string

func (stringTextMapKey) MarshalText() ([]byte, error) {
	return []byte("ignored-text-key"), nil
}

type pointerTextMapKey int

func (*pointerTextMapKey) MarshalText() ([]byte, error) {
	return []byte("pointer-text-key"), nil
}

type errorTextMapKey int

func (errorTextMapKey) MarshalText() ([]byte, error) {
	return nil, errors.New("map key error")
}

type sortedTextMapKey string

func (k sortedTextMapKey) MarshalText() ([]byte, error) {
	return []byte(k), nil
}

type jsonMap map[string]int

func (jsonMap) MarshalJSON() ([]byte, error) {
	return []byte(`"custom-map"`), nil
}

type valueZeroInt int

func (v valueZeroInt) IsZero() bool {
	return v == 1
}

type pointerZeroInt int

func (v *pointerZeroInt) IsZero() bool {
	return v != nil && *v == 2
}

type neverZeroInt int

func (neverZeroInt) IsZero() bool {
	return false
}

type zeroSlice []int

func (v zeroSlice) IsZero() bool {
	return len(v) == 1 && v[0] == 1
}

type recursiveNode struct {
	Value int            `json:"value"`
	Next  *recursiveNode `json:"next"`
}

type textBool bool

func (textBool) MarshalText() ([]byte, error) {
	return []byte("<text>"), nil
}

func TestMarshal(t *testing.T) {
	t.Parallel()

	intValue := 42
	intPointer := &intValue
	mapValue := map[string]int{"one": 1, "two": 2}
	mapPointer := &mapValue
	jsonValue := jsonInt(123)
	jsonPointer := &jsonValue
	pointerJSONValue := pointerJSONInt(123)
	pointerJSONPointer := &pointerJSONValue
	pointerTextKey := pointerTextMapKey(1)
	recursiveValue := &recursiveNode{Value: 1, Next: &recursiveNode{Value: 2}}

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
		{"defined int", testNumber(42)},
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
		{"value marshaler pointer chain", &jsonPointer},
		{"pointer marshaler", pointerJSONPointer},
		{"pointer marshaler pointer chain", &pointerJSONPointer},
		{"nil pointer", (*int)(nil)},
		{"array", [3]int{-1, 0, 1}},
		{"slice", []string{"one", "two"}},
		{"byte slice", []byte{0, 1, 255}},
		{"named byte slice", []testByte{0, 1, 255}},
		{"non-empty interface slice", []numberInterface{testNumber(1), testNumber(2)}},
		{"nil slice", []int(nil)},
		{"empty slice", []int{}},
		{"string string map", map[string]string{"one": "first", "two": "second"}},
		{"string int map", map[string]int{"negative": -1, "positive": 1}},
		{"int string map", map[int]string{-1: "negative", 1: "positive"}},
		{"uint bool map", map[uint]bool{1: true, 2: false}},
		{"string any map", map[string]any{"bool": true, "int": 1, "slice": []int{1, 2}}},
		{"string json.Marshaler map", map[string]jsonInt{"value": 1}},
		{"int json.Marshaler map", map[int]jsonInt{-1: 1}},
		{"uint json.Marshaler map", map[uint]jsonInt{1: 1}},
		{"named map json.Marshaler", jsonMap{"value": 1}},
		{"named map json.Marshaler value", map[string]jsonMap{"value": {"nested": 1}}},
		{"string pointer-receiver value map", map[string]pointerJSONInt{"value": 1}},
		{"string pointer value map", map[string]*pointerJSONInt{"value": &pointerJSONValue}},
		{"string json.Marshaler slice map", map[string]jsonIntSlice{"value": {1, 2}}},
		{"string json.Marshaler element slice map", map[string][]jsonInt{"value": {1, 2}}},
		{"string pointer-receiver element slice map", map[string][]pointerJSONInt{"value": {1, 2}}},
		{"string byte slice map", map[string][]byte{"value": {1, 2, 3}}},
		{"string named byte slice map", map[string][]testByte{"value": {1, 2, 3}}},
		{"int byte slice map", map[int][]byte{-1: {1, 2, 3}}},
		{"uint byte slice map", map[uint][]byte{1: {1, 2, 3}}},
		{"int any map", map[int]any{1: "one", 2: []int{2, 3}}},
		{"uint any map", map[uint]any{1: true, 2: "two"}},
		{"int8 key map", map[int8]string{-1: "negative", 1: "positive"}},
		{"uintptr key map", map[uintptr]bool{1: true, 2: false}},
		{"text key map", map[textMapKey]int{1: 1}},
		{"string text key map", map[stringTextMapKey]int{"raw-key": 1}},
		{"pointer text key map", map[*pointerTextMapKey]int{&pointerTextKey: 1}},
		{"nil pointer text key map", map[*pointerTextMapKey]int{nil: 1}},
		{"json.Marshaler key map", map[jsonInt]string{1: "one"}},
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
		{"recursive struct pointer", recursiveValue},
		{"slice of maps", []map[string]int{{"one": 1}, nil, {"two": 2}}},
		{"array of maps", [2]map[string]int{{"one": 1}, {"two": 2}}},
		{"nested map", map[string]map[string]int{"outer": {"inner": 1}}},
		{"nil map", map[string]int(nil)},
		{"empty map", map[string]int{}},
		{"empty struct", struct{}{}},
		{"struct json.Marshaler field", struct {
			Value jsonInt `json:"value"`
		}{Value: 123}},
		{"struct string-tagged json.Marshaler field", struct {
			Value jsonInt `json:"value,string"`
		}{Value: 123}},
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

func TestMarshalMapKinds(t *testing.T) {
	t.Parallel()

	type item struct {
		Value int `json:"value"`
	}

	tests := []struct {
		name  string
		value any
	}{
		// Every key kind supported by encoding/json.
		{"string key", map[string]int{"one": 1}},
		{"int key", map[int]int{-1: 1}},
		{"int8 key", map[int8]int{-1: 1}},
		{"int16 key", map[int16]int{-1: 1}},
		{"int32 key", map[int32]int{-1: 1}},
		{"int64 key", map[int64]int{-1: 1}},
		{"uint key", map[uint]int{1: 1}},
		{"uint8 key", map[uint8]int{1: 1}},
		{"uint16 key", map[uint16]int{1: 1}},
		{"uint32 key", map[uint32]int{1: 1}},
		{"uint64 key", map[uint64]int{1: 1}},
		{"uintptr key", map[uintptr]int{1: 1}},
		{"text key", map[textMapKey]int{1: 1}},

		// Every scalar value kind.
		{"bool value", map[string]bool{"value": true}},
		{"int value", map[string]int{"value": -1}},
		{"int8 value", map[string]int8{"value": -1}},
		{"int16 value", map[string]int16{"value": -1}},
		{"int32 value", map[string]int32{"value": -1}},
		{"int64 value", map[string]int64{"value": -1}},
		{"uint value", map[string]uint{"value": 1}},
		{"uint8 value", map[string]uint8{"value": 1}},
		{"uint16 value", map[string]uint16{"value": 1}},
		{"uint32 value", map[string]uint32{"value": 1}},
		{"uint64 value", map[string]uint64{"value": 1}},
		{"uintptr value", map[string]uintptr{"value": 1}},
		{"float32 value", map[string]float32{"value": 1.25}},
		{"float64 value", map[string]float64{"value": 1.25}},
		{"string value", map[string]string{"value": "text"}},

		// Every primitive slice element kind, including []byte's base64 rule.
		{"bool slice", map[string][]bool{"value": {true, false}}},
		{"int slice", map[string][]int{"value": {-1, 1}}},
		{"int8 slice", map[string][]int8{"value": {-1, 1}}},
		{"int16 slice", map[string][]int16{"value": {-1, 1}}},
		{"int32 slice", map[string][]int32{"value": {-1, 1}}},
		{"int64 slice", map[string][]int64{"value": {-1, 1}}},
		{"uint slice", map[string][]uint{"value": {0, 1}}},
		{"byte slice", map[string][]byte{"value": {0, 1, 255}}},
		{"uint16 slice", map[string][]uint16{"value": {0, 1}}},
		{"uint32 slice", map[string][]uint32{"value": {0, 1}}},
		{"uint64 slice", map[string][]uint64{"value": {0, 1}}},
		{"uintptr slice", map[string][]uintptr{"value": {0, 1}}},
		{"float32 slice", map[string][]float32{"value": {-1.25, 1.25}}},
		{"float64 slice", map[string][]float64{"value": {-1.25, 1.25}}},
		{"string slice", map[string][]string{"value": {"one", "two"}}},
		{"int-key slice fast path", map[int][]int{1: {1, 2}}},
		{"uint-key slice fast path", map[uint][]int{1: {1, 2}}},

		// Generic value encoders.
		{"pointer value", map[string]*int{"value": nil}},
		{"array value", map[string][2]int{"value": {1, 2}}},
		{"map value", map[string]map[string]int{"value": {"nested": 1}}},
		{"struct value", map[string]item{"value": {Value: 1}}},
		{"interface value", map[string]any{"value": item{Value: 1}}},
		{"json marshaler value", map[string]jsonInt{"value": 1}},
		{"text marshaler value", map[string]textBool{"value": true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := jsonexperiment.Marshal(tt.value)
			want, wantErr := json.Marshal(tt.value)
			if (gotErr != nil) != (wantErr != nil) {
				t.Fatalf("Marshal(%T) error = %v, json.Marshal error = %v", tt.value, gotErr, wantErr)
			}
			if gotErr == nil {
				assertJSONEqual(t, got, want)
			}
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

func TestMarshalSortMapKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"string fast path", map[string]int{"z": 1, "a": 2}, `{"a":2,"z":1}`},
		{"int fast path", map[int]string{2: "two", 10: "ten", -1: "negative"}, `{"-1":"negative","10":"ten","2":"two"}`},
		{"uint fast path", map[uint]string{2: "two", 10: "ten"}, `{"10":"ten","2":"two"}`},
		{"text keys", map[sortedTextMapKey]int{"z": 1, "a": 2}, `{"a":2,"z":1}`},
		{"nil pointer text key", map[*pointerTextMapKey]int{nil: 1}, `{"":1}`},
		{"custom values", map[string]jsonInt{"z": 1, "a": 2}, `{"a":"custom","z":"custom"}`},
		{"nested maps", map[string]map[string]int{
			"z": {"z": 1, "a": 2},
			"a": {"z": 3, "a": 4},
		}, `{"a":{"a":4,"z":3},"z":{"a":2,"z":1}}`},
	}

	options := jsonexperiment.MarshalOptions{SortMapKeys: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := jsonexperiment.MarshalWithOptions(tt.value, options)
			if err != nil {
				t.Fatalf("MarshalWithOptions(%T): %v", tt.value, err)
			}
			if string(got) != tt.want {
				t.Fatalf("MarshalWithOptions(%T) = %s, want %s", tt.value, got, tt.want)
			}
		})
	}

	got, err := jsonexperiment.MarshalWithFlags(
		map[string]int{"z": 1, "a": 2},
		jsonexperiment.MarshalFlagSortMapKeys,
	)
	if err != nil {
		t.Fatalf("MarshalWithFlags: %v", err)
	}
	if want := `{"a":2,"z":1}`; string(got) != want {
		t.Fatalf("MarshalWithFlags = %s, want %s", got, want)
	}

	got, err = jsonexperiment.MarshalWithFlags(
		map[string]int{"&": 1, "A": 2},
		jsonexperiment.MarshalFlagSortMapKeys|jsonexperiment.MarshalFlagEscapeHTML,
	)
	if err != nil {
		t.Fatalf("MarshalWithFlags HTML: %v", err)
	}
	if want := `{"\u0026":1,"A":2}`; string(got) != want {
		t.Fatalf("MarshalWithFlags HTML = %s, want %s", got, want)
	}

	if _, err = jsonexperiment.MarshalWithFlags(
		map[errorTextMapKey]int{1: 1},
		jsonexperiment.MarshalFlagSortMapKeys,
	); err == nil {
		t.Fatal("MarshalWithFlags text key error = nil")
	}

	if _, err = jsonexperiment.MarshalWithFlags(
		map[float64]int{1: 1},
		jsonexperiment.MarshalFlagSortMapKeys,
	); err == nil {
		t.Fatal("MarshalWithFlags unsupported map key error = nil")
	}
}

func TestMarshalInterfacePrecedence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		flags jsonexperiment.MarshalFlags
		want  string
	}{
		{"MarshalerAppend", allMarshalers{}, 0, `"append"`},
		{"json.Marshaler", jsonMarshaller{}, 0, `"json"`},
		{"integer json.Marshaler", jsonInt(123), 0, `"custom"`},
		{"*MarshalerAppend", &allMarshalers{}, 0, `"append"`},
		{"*json.Marshaler", &jsonMarshaller{}, 0, `"json"`},
		{"encoding.TextMarshaler", textMarshaler{}, 0, `"<text>"`},
		{"*encoding.TextMarshaler", &textMarshaler{}, 0, `"<text>"`},
		{"encoding.TextMarshaler HTML", textMarshaler{}, jsonexperiment.MarshalFlagEscapeHTML, `"\u003ctext\u003e"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonexperiment.MarshalWithFlags(tt.value, tt.flags)
			if err != nil {
				t.Fatalf("MarshalWithFlags: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("MarshalWithFlags = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMarshalCollectionElementInterfaces(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		want  string
	}{
		{"slice MarshalerAppend", []appendInt{1, 2}, `["append","append"]`},
		{"array MarshalerAppend", [2]appendInt{1, 2}, `["append","append"]`},
		{"slice json.Marshaler", []jsonString{"one", "two"}, `["json","json"]`},
		{"array json.Marshaler", [2]jsonString{"one", "two"}, `["json","json"]`},
		{"slice encoding.TextMarshaler", []textBool{true, false}, `["<text>","<text>"]`},
		{"array encoding.TextMarshaler", [2]textBool{true, false}, `["<text>","<text>"]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonexperiment.Marshal(tt.value)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if string(got) != tt.want {
				t.Fatalf("Marshal = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestMarshalCollectionPointerReceiver(t *testing.T) {
	t.Parallel()

	array := [2]pointerJSONInt{1, 2}
	structArray := struct {
		Values [2]pointerJSONInt `json:"values"`
	}{Values: array}
	tests := []struct {
		name  string
		value any
	}{
		{"slice", []pointerJSONInt{1, 2}},
		{"array", array},
		{"array pointer", &array},
		{"struct slice field", struct {
			Values []pointerJSONInt `json:"values"`
		}{Values: []pointerJSONInt{1, 2}}},
		{"struct array field", structArray},
		{"struct array field pointer", &structArray},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestMarshalMapPointerReceiver(t *testing.T) {
	t.Parallel()

	value := pointerJSONMap{"value": 1}
	pointer := &value
	array := [2]pointerJSONMap{value, value}
	tests := []struct {
		name  string
		value any
	}{
		{"pointer chain", &pointer},
		{"slice", []pointerJSONMap{value}},
		{"array", array},
		{"array pointer", &array},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestMarshalByteSliceFormat(t *testing.T) {
	t.Parallel()

	type formattedBytes struct {
		Base64    []byte `json:"base64,format:base64"`
		Base64URL []byte `json:"base64url,format:base64url"`
		Base32    []byte `json:"base32,format:base32"`
		Base32Hex []byte `json:"base32hex,format:base32hex"`
		Base16    []byte `json:"base16,format:base16"`
		Hex       []byte `json:"hex,format:hex"`
		Array     []byte `json:"array,format:array"`
		Nil       []byte `json:"nil,format:hex"`
	}

	value := []byte{0xfb, 0xff}
	got, err := jsonexperiment.Marshal(formattedBytes{
		Base64:    value,
		Base64URL: value,
		Base32:    value,
		Base32Hex: value,
		Base16:    value,
		Hex:       value,
		Array:     value,
	})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	want := []byte(
		`{"base64":"+/8=","base64url":"-_8=","base32":"7P7Q====",` +
			`"base32hex":"VFVG====","base16":"fbff","hex":"fbff",` +
			`"array":[251,255],"nil":null}`,
	)
	if !bytes.Equal(got, want) {
		t.Fatalf("Marshal = %s, want %s", got, want)
	}
}

func TestMarshalByteSliceFormatError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
	}{
		{
			"unknown format",
			struct {
				Data []byte `json:",format:unknown"`
			}{},
		},
		{
			"empty format",
			struct {
				Data []byte `json:",format:"`
			}{},
		},
		{
			"unsupported type",
			struct {
				Data string `json:",format:hex"`
			}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := jsonexperiment.Marshal(tt.value); err == nil {
				t.Fatal("Marshal error = nil")
			}
		})
	}
}

func TestMarshalFormatMarshalerPrecedence(t *testing.T) {
	t.Parallel()

	value := struct {
		Value jsonInt `json:"value,format:hex"`
	}{Value: 123}

	got, gotErr := jsonexperiment.Marshal(value)
	want, wantErr := jsonv2.Marshal(value)
	if (gotErr != nil) != (wantErr != nil) {
		t.Fatalf("Marshal error = %v, jsonv2.Marshal error = %v", gotErr, wantErr)
	}
	if gotErr == nil {
		assertJSONEqual(t, got, want)
	}
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

func TestMarshalOmitZeroMethod(t *testing.T) {
	t.Parallel()

	type value struct {
		ValueMethod   valueZeroInt   `json:"value_method,omitzero"`
		PointerMethod pointerZeroInt `json:"pointer_method,omitzero"`
		KeepZero      neverZeroInt   `json:"keep_zero,omitzero"`
		EmptyWins     neverZeroInt   `json:"empty_wins,omitempty,omitzero"`
		SliceMethod   zeroSlice      `json:"slice_method,omitzero"`
		KeepSlice     zeroSlice      `json:"keep_slice,omitzero"`
		Marshaler     jsonInt        `json:"marshaler,omitempty"`
	}
	testValue := value{
		ValueMethod:   1,
		PointerMethod: 2,
		SliceMethod:   zeroSlice{1},
		KeepSlice:     zeroSlice{},
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
		{"unsupported map key", map[float64]int{1: 1}},
		{"map key text marshaler error", map[errorTextMapKey]int{1: 1}},
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
