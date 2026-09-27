package workflowhandoff

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestHandoffDraftEnvelopeAndExactPayloadSizeBoundaries(t *testing.T) {
	value := minimalHandoffDraftValue("/tmp/wf01-target", "refs/heads/main", strings.Repeat("a", 40))
	canonical, err := canonicaljson.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if len(canonical) >= maxHandoffBytes {
		t.Fatalf("minimal fixture unexpectedly has %d bytes", len(canonical))
	}
	for _, size := range []int{len(canonical), maxHandoffBytes - 1, maxHandoffBytes} {
		input := append(append([]byte(nil), canonical...), bytes.Repeat([]byte{' '}, size-len(canonical))...)
		decoded, err := decodeHandoffDraft(input)
		if err != nil {
			t.Fatalf("decode size %d: %v", size, err)
		}
		if !bytes.Equal(decoded.Canonical, canonical) {
			t.Fatalf("size %d canonical bytes changed", size)
		}
	}
	over := append(append([]byte(nil), canonical...), bytes.Repeat([]byte{' '}, maxHandoffBytes+1-len(canonical))...)
	if _, err := decodeHandoffDraft(over); !IsClass(err, ErrorPayloadTooLarge) {
		t.Fatalf("limit+1 error = %v", err)
	}

	for _, mutation := range []struct {
		name  string
		value canonicaljson.Object
	}{
		{name: "wrong kind", value: replaceObjectMember(value, "kind", "ply.workflow.start-receipt-draft")},
		{name: "wrong schema version", value: replaceObjectMember(value, "schema_version", int64(2))},
		{name: "wrong format", value: replaceObjectMember(value, "format", "yaml")},
		{name: "wrong format version", value: replaceObjectMember(value, "format_version", int64(2))},
		{name: "wrong canonicalization", value: replaceObjectMember(value, "canonicalization", "none")},
		{name: "unknown field", value: append(append(canonicaljson.Object{}, value...), canonicaljson.Member{Name: "unknown", Value: true})},
		{name: "missing field", value: removeObjectMember(value, "reporting")},
	} {
		t.Run(mutation.name, func(t *testing.T) {
			input, err := canonicaljson.Marshal(mutation.value)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := decodeHandoffDraft(input); !IsClass(err, ErrorSchemaInvalid) {
				t.Fatalf("decode error = %v", err)
			}
		})
	}
}

func TestAllTwelveWorkflowDocumentShapesAreCanonicalGoldenBytes(t *testing.T) {
	dependencies, created, _ := completeRoundTrip(t)
	snapshot, err := dependencies.Store.ReadByLocator(created.Locator)
	if err != nil {
		t.Fatal(err)
	}
	principalValue, _ := objectMember(snapshot.Start.Value, "principal")
	principal := principalValue.(canonicaljson.Object)
	startResult := SubmitResult{DocumentID: snapshot.Start.DocumentID, SHA256: snapshot.Start.SHA256}
	values := []canonicaljson.Value{
		minimalHandoffDraftValue(snapshot.Handoff.Target.Worktree, snapshot.Handoff.Target.Ref, snapshot.Handoff.Target.OID),
		minimalStartDraftValue(t, snapshot, principal),
		minimalTerminalDraftValue(snapshot, principal, startResult),
		snapshot.Handoff.Value,
		snapshot.Start.Value,
		snapshot.Terminal.Value,
	}
	inspection, err := Inspect(dependencies, InspectInput{HandoffLocator: created.Locator, Format: "json"})
	if err != nil {
		t.Fatal(err)
	}
	allBytes := [][]byte{inspection.Bytes}
	for _, value := range values {
		contents, err := canonicaljson.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		allBytes = append(allBytes, contents)
	}
	if err := filepath.WalkDir(storeRoot(snapshot.WorkspaceRoot), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".json" && filepath.Ext(path) != ".ref" {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		allBytes = append(allBytes, contents)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	kinds := map[string]bool{}
	for index, contents := range allBytes {
		if bytes.HasSuffix(contents, []byte{'\n'}) {
			t.Fatalf("document %d has trailing LF", index)
		}
		value, err := canonicaljson.DecodeStrict(contents)
		if err != nil {
			t.Fatalf("document %d decode: %v", index, err)
		}
		roundTrip, err := canonicaljson.Marshal(value)
		if err != nil || !bytes.Equal(roundTrip, contents) {
			t.Fatalf("document %d is not byte-exact canonical JSON: %v", index, err)
		}
		object := value.(canonicaljson.Object)
		kinds[objectString(object, "kind")] = true
	}
	wantKinds := []string{
		"ply.workflow.handoff-draft", "ply.workflow.handoff", "ply.workflow.start-receipt-draft", "ply.workflow.start-receipt",
		"ply.workflow.terminal-result-draft", "ply.workflow.terminal-result", "ply.workflow.activity-event", "ply.workflow.reply-attempt-event",
		"ply.workflow.publication-ref", "ply.workflow.accepted-ref", "ply.workflow.activity-head-ref", "ply.workflow.handoff-inspection",
	}
	for _, kind := range wantKinds {
		if !kinds[kind] {
			t.Errorf("missing canonical golden shape %s", kind)
		}
	}
	if len(kinds) != len(wantKinds) {
		t.Fatalf("observed kinds = %#v", kinds)
	}
}
