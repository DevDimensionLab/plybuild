package workflowhandoff

import (
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

func TestParseHandoffID(t *testing.T) {
	want := HandoffID("hnd_0123456789abcdef0123456789abcdef")
	got, err := ParseHandoffID(string(want))
	if err != nil || got != want {
		t.Fatalf("ParseHandoffID = %q, %v", got, err)
	}
	for _, value := range []string{"", "hnd_ABC", "run_0123456789abcdef0123456789abcdef", "hnd_0123"} {
		if _, err := ParseHandoffID(value); err == nil {
			t.Fatalf("ParseHandoffID(%q) succeeded", value)
		}
	}
}

func TestScalarValidation(t *testing.T) {
	for _, value := range []string{"a", "a.b/c-d_1"} {
		if err := validateKey("key", value); err != nil {
			t.Fatalf("valid key %q: %v", value, err)
		}
	}
	for _, value := range []string{"", "A", "a//b", "a/../b", strings.Repeat("a", 129)} {
		if err := validateKey("key", value); err == nil {
			t.Fatalf("invalid key %q accepted", value)
		}
	}
	for _, value := range []string{"a", "a/b", "unicode/æ"} {
		if err := validateRepoPath(value); err != nil {
			t.Fatalf("valid path %q: %v", value, err)
		}
	}
	for _, value := range []string{"", ".", "../a", "a/../b", "/a", `a\b`} {
		if err := validateRepoPath(value); err == nil {
			t.Fatalf("invalid path %q accepted", value)
		}
	}
}

func TestExactObjectRejectsUnknownMissingAndDuplicateFields(t *testing.T) {
	valid := canonicaljson.Object{{Name: "a", Value: "x"}, {Name: "b", Value: int64(1)}}
	if _, err := exactObject(valid, "test", "a", "b"); err != nil {
		t.Fatal(err)
	}
	for _, object := range []canonicaljson.Object{
		{{Name: "a", Value: "x"}},
		{{Name: "a", Value: "x"}, {Name: "b", Value: int64(1)}, {Name: "c", Value: true}},
		{{Name: "a", Value: "x"}, {Name: "a", Value: "y"}},
	} {
		if _, err := exactObject(object, "test", "a", "b"); err == nil {
			t.Fatalf("exactObject(%#v) succeeded", object)
		}
	}
}

func TestCapabilityProofGoldenVector(t *testing.T) {
	secret := strings.Repeat("00", 32)
	proof, err := capabilityProof(secret, []byte("{}"))
	if err != nil {
		t.Fatal(err)
	}
	if proof != "22f8eea909400af98adf3681a9f31923ef6b7fcba4abb553d92823a3e9d5c25e" {
		t.Fatalf("proof = %s", proof)
	}
}
