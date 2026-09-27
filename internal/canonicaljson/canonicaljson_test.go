package canonicaljson

import (
	"reflect"
	"testing"
)

func TestDecodeStrictRejectsNonCanonicalDomainInputs(t *testing.T) {
	tests := []string{
		"\ufeff{}",
		"{\"a\":1,\"a\":2}",
		"{} {}",
		"1.0",
		"1e2",
		"-0",
		"9223372036854775808",
		`"\uD800"`,
		`"\uDC00"`,
	}
	for _, input := range tests {
		if _, err := DecodeStrict([]byte(input)); err == nil {
			t.Fatalf("DecodeStrict(%q) succeeded", input)
		}
	}
}

func TestDecodeStrictPreservesSupportedValues(t *testing.T) {
	got, err := DecodeStrict([]byte(" { \"n\" : -12, \"a\" : [true, null, \"x\"] } \n"))
	if err != nil {
		t.Fatal(err)
	}
	want := Object{
		{Name: "n", Value: int64(-12)},
		{Name: "a", Value: []Value{true, nil, "x"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decoded value = %#v, want %#v", got, want)
	}
}

func TestMarshalUsesRFC8785OrderingAndEscapes(t *testing.T) {
	value := Object{
		{Name: "\U00010000", Value: "line\n/\u2028"},
		{Name: "\ue000", Value: int64(-12)},
		{Name: "controls", Value: "\b\t\n\f\r\x00"},
	}
	got, err := Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	want := "{\"controls\":\"\\b\\t\\n\\f\\r\\u0000\",\"\U00010000\":\"line\\n/\u2028\",\"\ue000\":-12}"
	if string(got) != want {
		t.Fatalf("canonical bytes = %q, want %q", got, want)
	}
	if len(got) > 0 && got[len(got)-1] == '\n' {
		t.Fatal("canonical output has trailing LF")
	}
}

func TestMarshalRejectsDuplicateMembersAndUnsupportedTypes(t *testing.T) {
	if _, err := Marshal(Object{{Name: "a", Value: nil}, {Name: "a", Value: true}}); err == nil {
		t.Fatal("duplicate members accepted")
	}
	if _, err := Marshal(1); err == nil {
		t.Fatal("machine int accepted")
	}
}
