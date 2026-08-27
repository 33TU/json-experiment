package jsonexperiment

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/33TU/json-experiment/internal"
)

type numberInterface interface {
	number() int
}

type testNumber int

func (n testNumber) number() int {
	return int(n)
}

type testByte byte

type stringMapKey string

type testUint uint

type allMarshalers struct{}

var (
	_ MarshalerAppend = allMarshalers{}
	_ json.Marshaler  = allMarshalers{}
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

type mapMutatingField struct {
	parent *map[string]mapMutatingValue
	key    string
}

func (v mapMutatingField) MarshalJSONAppend(dst []byte) ([]byte, error) {
	delete(*v.parent, v.key)
	return append(dst, "null"...), nil
}

func (v mapMutatingField) MarshalJSON() ([]byte, error) {
	delete(*v.parent, v.key)
	return []byte("null"), nil
}

type mapMutatingValue struct {
	Mutator mapMutatingField `json:"mutator"`
	Value   int64            `json:"value"`
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
	timeValue := time.Date(2026, time.August, 16, 12, 34, 56, 789123456, time.FixedZone("EEST", 3*60*60))

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
		{"time.Time", timeValue},
		{"*time.Time", &timeValue},
		{"nil *time.Time", (*time.Time)(nil)},
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

			got, err := Marshal(tt.value)
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

func TestMarshalTimeError(t *testing.T) {
	t.Parallel()

	value := time.Date(10_000, time.January, 1, 0, 0, 0, 0, time.UTC)
	_, wantErr := value.AppendText(nil)
	if wantErr == nil {
		t.Fatal("time.Time.AppendText error = nil")
	}

	for _, tt := range []struct {
		name  string
		value any
	}{
		{"time.Time", value},
		{"*time.Time", &value},
	} {
		t.Run(tt.name, func(t *testing.T) {
			prefix := []byte("prefix")
			got, err := MarshalAppend(prefix, tt.value)
			if err == nil || err.Error() != wantErr.Error() {
				t.Fatalf("MarshalAppend error = %v, want %v", err, wantErr)
			}
			if string(got) != string(prefix) {
				t.Fatalf("MarshalAppend result = %q, want %q", got, prefix)
			}
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

			got, gotErr := Marshal(tt.value)
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

func TestMarshalMapValueCanMutateParentMap(t *testing.T) {
	newValues := func() map[string]mapMutatingValue {
		values := make(map[string]mapMutatingValue)
		values["key"] = mapMutatingValue{
			Mutator: mapMutatingField{parent: &values, key: "key"},
			Value:   42,
		}
		return values
	}

	values := newValues()
	got, err := MarshalAppend(nil, values)
	if err != nil {
		t.Fatalf("MarshalAppend: %v", err)
	}

	wantValues := newValues()
	want, err := json.Marshal(wantValues)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("MarshalAppend = %s, want %s", got, want)
	}
	if len(values) != 0 || len(wantValues) != 0 {
		t.Fatalf("map lengths = %d and %d, want 0", len(values), len(wantValues))
	}
}

func TestMarshalAppend(t *testing.T) {
	t.Parallel()

	dst := make([]byte, 0, 64)
	dst = append(dst, "prefix:"...)

	got, err := MarshalAppend(dst, []int{1, 2, 3})
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
	}{
		{"string fast path", map[string]int{"z": 1, "a": 2}},
		{"string int8 fast path", map[string]int8{"z": 1, "a": -2}},
		{"string int16 fast path", map[string]int16{"z": 1, "a": -2}},
		{"string int32 fast path", map[string]int32{"z": 1, "a": -2}},
		{"string int64 fast path", map[string]int64{"z": 1, "a": -2}},
		{"named string key fast path", map[stringMapKey]int{"z": 1, "a": 2}},
		{"named string int fast path", map[string]testNumber{"z": 1, "a": -2}},
		{"string values", map[string]string{"z": "last", "a": "first"}},
		{"slice values", map[string][]int{"z": {3}, "a": {1, 2}}},
		{"interface values", map[string]any{"z": []int{2}, "a": "first"}},
		{"int fast path", map[int]string{2: "two", 10: "ten", -1: "negative"}},
		{"int8 fast path", map[int8]string{2: "two", 10: "ten", -1: "negative"}},
		{"int16 fast path", map[int16]string{2: "two", 10: "ten", -1: "negative"}},
		{"int32 fast path", map[int32]string{2: "two", 10: "ten", -1: "negative"}},
		{"int64 fast path", map[int64]string{2: "two", 10: "ten", -1: "negative"}},
		{"named int key fast path", map[testNumber]string{2: "two", 10: "ten", -1: "negative"}},
		{"uint fast path", map[uint]string{2: "two", 10: "ten"}},
		{"uint8 fast path", map[uint8]string{2: "two", 10: "ten"}},
		{"uint16 fast path", map[uint16]string{2: "two", 10: "ten"}},
		{"uint32 fast path", map[uint32]string{2: "two", 10: "ten"}},
		{"uint64 fast path", map[uint64]string{2: "two", 10: "ten"}},
		{"uintptr fast path", map[uintptr]string{2: "two", 10: "ten"}},
		{"named uint key fast path", map[testUint]string{2: "two", 10: "ten"}},
		{"int64 key limits", map[int64]string{math.MinInt64: "min", math.MaxInt64: "max"}},
		{"uint64 key limit", map[uint64]string{0: "zero", math.MaxUint64: "max"}},
		{"text keys", map[sortedTextMapKey]int{"z": 1, "a": 2}},
		{"nil pointer text key", map[*pointerTextMapKey]int{nil: 1}},
		{"custom values", map[string]jsonInt{"z": 1, "a": 2}},
		{"nested maps", map[string]map[string]int{
			"z": {"z": 1, "a": 2},
			"a": {"z": 3, "a": 4},
		}},
	}

	options := MarshalOptions{SortMapKeys: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := MarshalWithOptions(tt.value, options)
			want, wantErr := json.Marshal(tt.value)
			if (gotErr != nil) != (wantErr != nil) {
				t.Fatalf("MarshalWithOptions(%T) error = %v, json.Marshal error = %v", tt.value, gotErr, wantErr)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("MarshalWithOptions(%T) = %s, json.Marshal = %s", tt.value, got, want)
			}
		})
	}

	value := map[string]int{"z": 1, "a": 2}
	got, gotErr := MarshalWithFlags(
		value,
		MarshalFlagSortMapKeys,
	)
	want, wantErr := json.Marshal(value)
	if (gotErr != nil) != (wantErr != nil) {
		t.Fatalf("MarshalWithFlags error = %v, json.Marshal error = %v", gotErr, wantErr)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("MarshalWithFlags = %s, json.Marshal = %s", got, want)
	}

	value = map[string]int{"&": 1, "A": 2}
	got, gotErr = MarshalWithFlags(
		value,
		MarshalFlagSortMapKeys|MarshalFlagEscapeHTML,
	)
	want, wantErr = json.Marshal(value)
	if (gotErr != nil) != (wantErr != nil) {
		t.Fatalf("MarshalWithFlags HTML error = %v, json.Marshal error = %v", gotErr, wantErr)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("MarshalWithFlags HTML = %s, json.Marshal = %s", got, want)
	}

	errorTextMap := map[errorTextMapKey]int{1: 1}
	_, gotErr = MarshalWithFlags(
		errorTextMap,
		MarshalFlagSortMapKeys,
	)
	_, wantErr = json.Marshal(errorTextMap)
	if (gotErr != nil) != (wantErr != nil) {
		t.Fatalf("MarshalWithFlags text key error = %v, json.Marshal error = %v", gotErr, wantErr)
	}

	unsupportedMap := map[float64]int{1: 1}
	_, gotErr = MarshalWithFlags(
		unsupportedMap,
		MarshalFlagSortMapKeys,
	)
	if gotErr == nil {
		t.Fatal("MarshalWithFlags unsupported map key error = nil")
	}
}

func TestMarshalInterfacePrecedence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value any
		flags MarshalFlags
		want  string
	}{
		{"MarshalerAppend", allMarshalers{}, 0, `"append"`},
		{"json.Marshaler", jsonMarshaller{}, 0, `"json"`},
		{"integer json.Marshaler", jsonInt(123), 0, `"custom"`},
		{"*MarshalerAppend", &allMarshalers{}, 0, `"append"`},
		{"*json.Marshaler", &jsonMarshaller{}, 0, `"json"`},
		{"encoding.TextMarshaler", textMarshaler{}, 0, `"<text>"`},
		{"*encoding.TextMarshaler", &textMarshaler{}, 0, `"<text>"`},
		{"encoding.TextMarshaler HTML", textMarshaler{}, MarshalFlagEscapeHTML, `"\u003ctext\u003e"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalWithFlags(tt.value, tt.flags)
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
			got, err := Marshal(tt.value)
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
			got, err := Marshal(tt.value)
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
			got, err := Marshal(tt.value)
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

	got, err := Marshal(value)
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
	got, err := Marshal(formattedBytes{
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
			if _, err := Marshal(tt.value); err == nil {
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

	got, err := Marshal(value)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}
	assertJSONEqual(t, got, []byte(`{"value":"custom"}`))
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

	got, err := Marshal(testValue)
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

	got, err := Marshal(testValue)
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

	options := MarshalOptions{EscapeHTML: true}
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

	got, err := MarshalWithOptions(value, options)
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
		encoded, err := MarshalWithOptions(value, options)
		if err != nil {
			t.Fatalf("MarshalWithOptions(%T): %v", value, err)
		}
		if bytes.ContainsAny(encoded, "<>&") {
			t.Fatalf("MarshalWithOptions(%T) left HTML characters unescaped: %s", value, encoded)
		}
	}

	got, err = MarshalAppendWithOptions([]byte("prefix:"), "<>&", options)
	if err != nil {
		t.Fatalf("MarshalAppendWithOptions: %v", err)
	}
	if want := []byte(`prefix:"\u003c\u003e\u0026"`); !bytes.Equal(got, want) {
		t.Fatalf("MarshalAppendWithOptions = %q, want %q", got, want)
	}

	got, err = Marshal("<>&")
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
	options := MarshalOptions{ValidateString: true}

	want, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}

	got, err := MarshalWithOptions(value, options)
	if err != nil {
		t.Fatalf("MarshalWithOptions: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("MarshalWithOptions = %q, want %q", got, want)
	}

	got, err = MarshalWithFlags(value, MarshalFlagValidateString)
	if err != nil {
		t.Fatalf("MarshalWithFlags: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("MarshalWithFlags = %q, want %q", got, want)
	}

	got, err = MarshalAppendWithOptions([]byte("prefix:\xff"), value, options)
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

	got, err = MarshalWithOptions(nested, options)
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

	got, err = Marshal(value)
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
			if _, err := Marshal(tt.value); err == nil {
				t.Fatalf("Marshal(%T) error = nil", tt.value)
			}
		})
	}
}

func TestMarshalReturnsOwnedBytes(t *testing.T) {
	first, err := Marshal("first")
	if err != nil {
		t.Fatalf("Marshal(first): %v", err)
	}

	if _, err := Marshal("second"); err != nil {
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

func TestSortSortedMapIndexesRandomized(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for size := 0; size <= 256; size++ {
		stringsToSort := make([]string, size)
		integerKeys := make([]sortedMapIntegerKey, size)
		for i := range size {
			prefix := stringsToSort[rng.Intn(i+1)]
			stringsToSort[i] = prefix + strconv.Itoa(rng.Intn(32))

			value := int64(rng.Uint64())
			integerKeys[i].length = uint8(len(internal.AppendInt(integerKeys[i].text[:0], value)))
		}

		assertSortedMapStringKeys(t, stringsToSort, sortSortedMapStringKeys)
		assertSortedMapIndexes(t, integerKeys, sortSortedMapIntegerIndexes, func(key sortedMapIntegerKey) string {
			return string(key.text[:key.length])
		})
	}
}

func TestSortSortedMapKeysAdversarial(t *testing.T) {
	ordered := make([]string, 128)
	for i := range ordered {
		ordered[i] = "key_" + strconv.Itoa(i)
	}
	slices.Sort(ordered)
	reversed := slices.Clone(ordered)
	slices.Reverse(reversed)

	commonPrefix := strings.Repeat("shared-prefix-", 256)
	longCommonPrefix := make([]string, 128)
	for i := range longCommonPrefix {
		longCommonPrefix[i] = commonPrefix + strconv.Itoa(i)
	}

	tests := []struct {
		name string
		keys []string
	}{
		{"empty", nil},
		{"one", []string{"one"}},
		{"insertion threshold minus one", slices.Clone(ordered[:11])},
		{"insertion threshold", slices.Clone(ordered[:12])},
		{"insertion threshold plus one", slices.Clone(ordered[:13])},
		{"sorted", ordered},
		{"reverse sorted", reversed},
		{"duplicates", []string{"same", "same", "same", "prefix", "prefix"}},
		{"prefix lengths", []string{"", "a", "aa", "aaa", "aab", "ab", "b"}},
		{"invalid UTF-8", []string{"\xff", "\xfe", "a\xff", "a\xfe", "a", ""}},
		{"long common prefix", longCommonPrefix},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertSortedMapStringKeys(t, tt.keys, sortSortedMapStringKeys)
		})
	}

	t.Run("forced fallback", func(t *testing.T) {
		assertSortedMapStringKeys(t, reversed, func(keys []sortedMapStringKey) {
			radixSortSortedMapStringKeys(keys, 0, 0)
		})
	})
}

func TestSortSortedMapIntegerIndexesFallback(t *testing.T) {
	keys := make([]sortedMapIntegerKey, 128)
	for i := range keys {
		value := int64(i*7919 - 500_000)
		keys[i].length = uint8(len(internal.AppendInt(keys[i].text[:0], value)))
	}
	assertSortedMapIndexes(t, keys, func(indexes []int, keys []sortedMapIntegerKey) {
		radixSortSortedMapIntegerIndexes(indexes, keys, 0, 0)
	}, func(key sortedMapIntegerKey) string {
		return string(key.text[:key.length])
	})
}

func TestReflectMapIterPrefix(t *testing.T) {
	type mapValue struct {
		Text    string
		Numbers []int
		Pointer *int
		Padding [256]byte
	}

	pointedValue := 42
	wantValue := mapValue{
		Text:    "value",
		Numbers: []int{1, 2, 3},
		Pointer: &pointedValue,
	}
	wantValue.Padding[255] = 1
	value := reflect.ValueOf(map[stringMapKey]mapValue{"key": wantValue})
	iter := value.MapRange()
	if !iter.Next() {
		t.Fatal("map iterator is empty")
	}

	prefix := (*reflectMapIterPrefix)(unsafe.Pointer(iter))
	if prefix.key == nil || prefix.elem == nil {
		t.Fatalf("iterator prefix pointers = (%p, %p), want non-nil", prefix.key, prefix.elem)
	}

	var reflectedKey stringMapKey
	reflect.ValueOf(&reflectedKey).Elem().SetIterKey(iter)
	if got := *(*stringMapKey)(prefix.key); got != reflectedKey {
		t.Fatalf("iterator key = %q, reflect key = %q", got, reflectedKey)
	}

	var reflectedValue mapValue
	reflect.ValueOf(&reflectedValue).Elem().SetIterValue(iter)
	if got := *(*mapValue)(prefix.elem); !reflect.DeepEqual(got, reflectedValue) {
		t.Fatalf("iterator value = %#v, reflect value = %#v", got, reflectedValue)
	}

	var copiedValue mapValue
	runtimeTypedmemmove(
		internal.InterfaceData(reflect.TypeFor[mapValue]()),
		unsafe.Pointer(&copiedValue),
		prefix.elem,
	)
	if !reflect.DeepEqual(copiedValue, wantValue) {
		t.Fatalf("typedmemmove value = %#v, want %#v", copiedValue, wantValue)
	}
}

func assertSortedMapStringKeys(
	t *testing.T,
	input []string,
	sortKeys func([]sortedMapStringKey),
) {
	t.Helper()
	keys := make([]sortedMapStringKey, len(input))
	want := slices.Clone(input)
	for i, key := range input {
		keys[i] = sortedMapStringKey{text: key, valueIndex: i}
	}
	rand.New(rand.NewSource(int64(len(keys)))).Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	sortKeys(keys)
	slices.Sort(want)
	for i, key := range keys {
		if key.text != want[i] {
			t.Fatalf("key %d = %q, want %q", i, key.text, want[i])
		}
		if key.text != input[key.valueIndex] {
			t.Fatalf("key %d value index = %d, which belongs to %q", i, key.valueIndex, input[key.valueIndex])
		}
	}
}

func assertSortedMapIndexes[K any](
	t *testing.T,
	keys []K,
	sortIndexes func([]int, []K),
	stringify func(K) string,
) {
	t.Helper()
	indexes := make([]int, len(keys))
	want := make([]string, len(keys))
	for i, key := range keys {
		indexes[i] = i
		want[i] = stringify(key)
	}
	rand.New(rand.NewSource(int64(len(keys)))).Shuffle(len(indexes), func(i, j int) {
		indexes[i], indexes[j] = indexes[j], indexes[i]
	})

	sortIndexes(indexes, keys)
	slices.Sort(want)
	for i, index := range indexes {
		if got := stringify(keys[index]); got != want[i] {
			t.Fatalf("key %d = %q, want %q", i, got, want[i])
		}
	}
}
