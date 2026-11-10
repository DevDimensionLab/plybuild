package canonicaljson

import (
	"bytes"
	"crypto/sha256"
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

func TestCanonicalWhitespaceDigestAndInvalidUTF8Boundaries(t *testing.T) {
	compact, err := DecodeStrict([]byte(`{"a":[true,null,"æ"],"n":1}`))
	if err != nil {
		t.Fatal(err)
	}
	spaced, err := DecodeStrict([]byte(" \n { \"n\" : 1, \"a\" : [ true, null, \"æ\" ] } \t"))
	if err != nil {
		t.Fatal(err)
	}
	compactBytes, err := Marshal(compact)
	if err != nil {
		t.Fatal(err)
	}
	spacedBytes, err := Marshal(spaced)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(compactBytes, spacedBytes) || sha256.Sum256(compactBytes) != sha256.Sum256(spacedBytes) {
		t.Fatalf("equivalent JSON did not produce the same canonical bytes: %q != %q", compactBytes, spacedBytes)
	}
	if bytes.HasSuffix(compactBytes, []byte{'\n'}) {
		t.Fatal("canonical output has a trailing LF")
	}

	invalid := string([]byte{'x', 0xff, 'y'})
	if _, err := DecodeStrict([]byte{'"', 0xff, '"'}); err == nil {
		t.Fatal("DecodeStrict accepted invalid UTF-8")
	}
	if _, err := Marshal(invalid); err == nil {
		t.Fatal("Marshal accepted invalid UTF-8 string value")
	}
	if _, err := Marshal(Object{{Name: invalid, Value: nil}}); err == nil {
		t.Fatal("Marshal accepted invalid UTF-8 member name")
	}
}
