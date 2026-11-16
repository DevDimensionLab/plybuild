package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/inbox"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

var updateInboxFixtures = flag.Bool("update-inbox-fixtures", false, "rewrite schemas/ply-inbox-v1/capabilities-excerpt.json")

const inboxFixtureDirectory = "../schemas/ply-inbox-v1"

type inboxFixture struct{ backlogFixture }

func newInboxFixture(t *testing.T, projects ...string) inboxFixture {
	return inboxFixture{newBacklogFixture(t, projects...)}
}

func (f inboxFixture) run(stdin string, args ...string) (string, error) {
	command := newInboxCommandWithDirectory(inboxDirectory(f.show))
	command.SilenceUsage, command.SilenceErrors = true, true // as RootCmd under ExecuteE
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetIn(strings.NewReader(stdin))
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), err
}

func (f inboxFixture) success(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	out, err := f.run(stdin, args...)
	if err != nil {
		t.Fatalf("inbox %v: %v", args, err)
	}
	return out
}

func inboxRequest(i int, body string) string {
	return fmt.Sprintf(`{"kind":"InboxPutRequest@1","conversation_id":"conv_%032x","message_id":"msg_%032x","replied_at_utc":"2026-11-11T14:0%d:00Z","title":null,"body":%q,"truncated":false,"full_sha256":null,"answer_source":"codex.item_completed","agent_label":"quality","provider":"codex","write_kind":"live"}`, i, i, i, body)
}

func decodeInboxJSON(t *testing.T, raw string, target any) {
	t.Helper()
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		t.Fatalf("decode %q: %v", raw, err)
	}
}

func inboxErrorOf(t *testing.T, raw string) inbox.ErrorResult {
	t.Helper()
	var result inbox.ErrorResult
	decodeInboxJSON(t, raw, &result)
	if result.Kind != inbox.ErrorKind || result.SchemaVersion != 1 {
		t.Fatalf("error result %+v", result)
	}
	return result
}

var scope = []string{"--project", "ply", "--owner", "owner-1"}

func withScope(args ...string) []string { return append(append([]string{}, args...), scope...) }

// C2: every verb in text and json.
func TestInboxVerbsTextAndJSON(t *testing.T) {
	f := newInboxFixture(t, "ply")
	file := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(file, []byte(inboxRequest(1, "# Første svar\nmer tekst")), 0o600); err != nil {
		t.Fatal(err)
	}
	out := f.success(t, "", withScope("put", "--file", file)...)
	if !strings.HasPrefix(out, "created conv_00000000000000000000000000000001/msg_00000000000000000000000000000001 sha256:") || strings.Count(out, "\n") != 1 {
		t.Fatalf("put text %q", out)
	}
	var put inbox.PutResult
	decodeInboxJSON(t, f.success(t, "", withScope("put", "--file", file, "--format", "json")...), &put)
	if put.Kind != inbox.PutResultKind || put.Effect != inbox.EffectUnchanged || put.Item.DerivedTitle != "Første svar" {
		t.Fatalf("put json %+v", put)
	}
	if _, err := os.Stat(filepath.Join(f.roots["ply"], "data", "inbox", "owner-1", "messages")); err != nil {
		t.Fatalf("default data path not used: %v", err)
	}
	ids := []string{"--conversation", put.Item.ConversationID, "--message", put.Item.MessageID}

	out = f.success(t, "", withScope(append([]string{"show"}, ids...)...)...)
	if !strings.Contains(out, "Første svar") || !strings.Contains(out, "Read: no") || !strings.Contains(out, "mer tekst") {
		t.Fatalf("show text %q", out)
	}
	out = f.success(t, "", withScope(append([]string{"mark", "--value", "read"}, ids...)...)...)
	if out != "marked "+put.Item.ConversationID+"/"+put.Item.MessageID+" read\n" {
		t.Fatalf("mark text %q", out)
	}
	var mark inbox.MarkResult
	decodeInboxJSON(t, f.success(t, "", withScope(append([]string{"mark", "--value", "read", "--format", "json"}, ids...)...)...), &mark)
	if mark.Effect != inbox.EffectNone || !mark.Read {
		t.Fatalf("mark json %+v", mark)
	}
	var item inbox.Item
	decodeInboxJSON(t, f.success(t, "", withScope(append([]string{"show", "--format", "json"}, ids...)...)...), &item)
	if item.Kind != inbox.ItemKind || item.Read == nil || !*item.Read || item.DisplayTitle != "Første svar" {
		t.Fatalf("show json %+v", item)
	}

	f.success(t, inboxRequest(2, "Andre"), withScope("put", "--file", "-")...)
	out = f.success(t, "", withScope("list", "--limit", "1")...)
	if !strings.Contains(out, "2 items, 1 unread (filename_unverified)") || !strings.Contains(out, "unread  conv_00000000000000000000000000000002/msg_00000000000000000000000000000002  Andre") || !strings.Contains(out, "Next cursor: ") {
		t.Fatalf("list text %q", out)
	}
	var list inbox.List
	decodeInboxJSON(t, f.success(t, "", withScope("list", "--format", "json", "--limit", "1")...), &list)
	if list.Kind != inbox.ListKind || list.TotalCount != 2 || list.NextCursor == nil || len(list.Items) != 1 {
		t.Fatalf("list json %+v", list)
	}
	decodeInboxJSON(t, f.success(t, "", withScope("list", "--format", "json", "--cursor", *list.NextCursor)...), &list)
	if len(list.Items) != 1 || list.Items[0].MessageID != put.Item.MessageID || list.NextCursor != nil {
		t.Fatalf("list page 2 %+v", list)
	}
}

// C2: json errors are InboxError@1 on stdout; details carry no input or paths.
func TestInboxErrorsAreTypedAndLeakNothing(t *testing.T) {
	f := newInboxFixture(t, "ply")
	f.success(t, inboxRequest(1, "Original"), withScope("put", "--file", "-")...)
	out, err := f.run(inboxRequest(1, "Hemmelig endring"), withScope("put", "--file", "-", "--format", "json")...)
	result := inboxErrorOf(t, out)
	if err == nil || ExitCode(err) != 1 || !strings.HasPrefix(err.Error(), "inbox_conflict: ") || result.Class != inbox.ClassConflict || result.Retryable ||
		!reflect.DeepEqual(result.ConflictingFields, []string{"body"}) || result.ExistingPayloadSHA256 == "" || result.SubmittedPayloadSHA256 == "" || strings.Contains(out, "Hemmelig") {
		t.Fatalf("conflict %v %q", err, out)
	}
	if out, err := f.run(inboxRequest(1, "Hemmelig endring"), withScope("put", "--file", "-")...); err == nil || out != "" {
		t.Fatalf("text error wrote stdout %q", out)
	}

	cases := []struct {
		name  string
		args  []string
		class string
	}{
		{"flag error", withScope("list", "--format", "json", "--limit", "secret-limit"), inbox.ClassInvalidArguments},
		{"unknown flag", withScope("list", "--format", "json", "--secret-flag=x"), inbox.ClassInvalidArguments},
		{"missing owner", []string{"list", "--format", "json", "--project", "ply"}, inbox.ClassInvalidArguments},
		{"missing project", []string{"list", "--format", "json", "--owner", "owner-1"}, inbox.ClassInvalidArguments},
		{"dot owner", []string{"list", "--format", "json", "--project", "ply", "--owner", ".."}, inbox.ClassInvalidArguments},
		{"owner checked before project", []string{"list", "--format", "json", "--project", "unregistered", "--owner", ".."}, inbox.ClassInvalidArguments},
		{"limit range", withScope("list", "--format", "json", "--limit", "51"), inbox.ClassInvalidArguments},
		{"bad cursor", withScope("list", "--format", "json", "--cursor", "secret-cursor"), inbox.ClassInvalidArguments},
		{"bad value", withScope("mark", "--format", "json", "--conversation", fmt.Sprintf("conv_%032x", 1), "--message", fmt.Sprintf("msg_%032x", 1), "--value", "secret-value"), inbox.ClassInvalidArguments},
		{"not found", withScope("show", "--format", "json", "--conversation", fmt.Sprintf("conv_%032x", 9), "--message", fmt.Sprintf("msg_%032x", 9)), inbox.ClassItemNotFound},
		{"positional", withScope("list", "--format", "json", "secret-arg"), inbox.ClassInvalidArguments},
		{"unknown project", []string{"list", "--format", "json", "--project", "secret-project", "--owner", "owner-1"}, string(workspace.ErrorProjectNotFound)},
		{"missing file flag", withScope("put", "--format", "json"), inbox.ClassInvalidArguments},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := f.run("", c.args...)
			if err == nil || ExitCode(err) != 1 || !strings.HasPrefix(err.Error(), c.class+": ") {
				t.Fatalf("error %v", err)
			}
			result := inboxErrorOf(t, out)
			if result.Class != c.class || result.Retryable || result.ConflictingFields != nil {
				t.Fatalf("result %+v", result)
			}
			if c.class != string(workspace.ErrorProjectNotFound) && strings.Contains(out+err.Error(), "secret") {
				t.Fatalf("raw input leaked: %q %v", out, err)
			}
			if strings.Contains(out+err.Error(), f.roots["ply"]) {
				t.Fatalf("path leaked: %q", out)
			}
		})
	}
}

// C3: put reads stdin and rejects oversize, unknown keys and trailing data.
func TestInboxPutRequestLimits(t *testing.T) {
	f := newInboxFixture(t, "ply")
	valid := inboxRequest(1, "Hei")
	cases := map[string]string{
		"oversize":      strings.Replace(valid, `"Hei"`, `"`+strings.Repeat("x", 512*1024)+`"`, 1),
		"unknown key":   strings.Replace(valid, `"title":null`, `"title":null,"extra":true`, 1),
		"trailing data": valid + "\n{}",
	}
	for name, stdin := range cases {
		t.Run(name, func(t *testing.T) {
			out, err := f.run(stdin, withScope("put", "--file", "-", "--format", "json")...)
			if result := inboxErrorOf(t, out); err == nil || result.Class != inbox.ClassInvalidArguments {
				t.Fatalf("%v %+v", err, result)
			}
		})
	}
	if _, err := os.Stat(filepath.Join(f.roots["ply"], "data", "inbox")); !os.IsNotExist(err) {
		t.Fatal("rejected put created the store")
	}
	missing := filepath.Join(t.TempDir(), "secret-missing.json")
	out, err := f.run("", withScope("put", "--file", missing, "--format", "json")...)
	if result := inboxErrorOf(t, out); err == nil || result.Class != inbox.ClassInvalidArguments || strings.Contains(out, "secret") {
		t.Fatalf("missing file %v %q", err, out)
	}
}

// C1: inbox_extensions is additive and pinned by the fixture.
func TestInboxCapabilitiesAreAdditive(t *testing.T) {
	catalog := coreCapabilities(capabilityBuild{Version: "test"})
	core, err := json.Marshal(struct {
		Operations     []capabilityOperation `json:"operations"`
		ReadExtensions []capabilityOperation `json:"read_extensions"`
	}{catalog.Operations, catalog.ReadExtensions})
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(filepath.Join("testdata", "inbox", "capabilities-core-before.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(append(core, '\n'), before) {
		t.Fatalf("operations/read_extensions changed:\n%s", core)
	}
	if catalog.Kind != "PlyCapabilities@1" || catalog.SchemaVersion != 1 {
		t.Fatal("catalog identity changed")
	}
	one := 1
	want := []capabilityOperation{}
	for _, verb := range []struct {
		name, effect, kind string
		selectors          []string
	}{
		{"put", "inbox_message_record", inbox.PutResultKind, []string{"--workspace", "--project", "--owner", "--file"}},
		{"show", "read", inbox.ItemKind, []string{"--workspace", "--project", "--owner", "--conversation", "--message"}},
		{"list", "read", inbox.ListKind, []string{"--workspace", "--project", "--owner", "--limit", "--cursor"}},
		{"mark", "inbox_read_mark", inbox.MarkResultKind, []string{"--workspace", "--project", "--owner", "--conversation", "--message", "--value"}},
	} {
		want = append(want, capabilityOperation{ID: "inbox." + verb.name, Command: []string{"inbox", verb.name}, Mode: "inbox_v1", Selectors: verb.selectors,
			Formats: []string{"text", "json"}, Filters: []string{}, FilterPolicy: "none",
			ResultSchemas: []capabilitySchema{{verb.kind, &one}, {inbox.ErrorKind, &one}}, Effect: verb.effect})
	}
	if !reflect.DeepEqual(catalog.InboxExtensions, want) {
		t.Fatalf("inbox_extensions %+v", catalog.InboxExtensions)
	}

	excerpt, err := encodeInboxExcerpt(catalog)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(inboxFixtureDirectory, "capabilities-excerpt.json")
	if *updateInboxFixtures {
		if err := os.WriteFile(path, excerpt, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	fixture, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(fixture, excerpt) {
		t.Fatalf("capabilities-excerpt.json differs (%v):\n%s", err, excerpt)
	}

	// plyapp-server 5174d61 decodes the catalogue with this shape and no
	// DisallowUnknownFields; the additive array must not disturb it.
	type plyappSchema struct {
		Kind    string `json:"kind"`
		Version *int   `json:"schema_version"`
	}
	type plyappEntry struct {
		ID           string         `json:"id"`
		Command      []string       `json:"command"`
		Mode         string         `json:"mode"`
		Selectors    []string       `json:"selectors"`
		Formats      []string       `json:"formats"`
		Filters      []string       `json:"filters"`
		FilterPolicy string         `json:"filter_policy"`
		Schemas      []plyappSchema `json:"result_schemas"`
		Effect       string         `json:"effect"`
	}
	type plyappCatalogue struct {
		Extensions []plyappEntry `json:"read_extensions,omitempty"`
		Kind       string        `json:"kind"`
		Version    int           `json:"schema_version"`
		Coverage   string        `json:"coverage"`
		Entries    []plyappEntry `json:"operations"`
	}
	full, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}
	var decoded plyappCatalogue
	if err := json.Unmarshal(full, &decoded); err != nil || decoded.Kind != "PlyCapabilities@1" || len(decoded.Entries) != 5 || len(decoded.Extensions) != len(catalog.ReadExtensions) {
		t.Fatalf("plyapp decode %v %+v", err, decoded)
	}
	var inboxOnly struct {
		InboxExtensions []plyappEntry `json:"inbox_extensions"`
	}
	if err := json.Unmarshal(fixture, &inboxOnly); err != nil || len(inboxOnly.InboxExtensions) != 4 || inboxOnly.InboxExtensions[0].Schemas[0].Kind != inbox.PutResultKind {
		t.Fatalf("fixture decode %v %+v", err, inboxOnly)
	}
}

func TestCapabilitiesJSONAlwaysHasInboxExtensions(t *testing.T) {
	command := newCapabilitiesCommand(func() capabilityBuild { return capabilityBuild{Version: "test"} })
	var out bytes.Buffer
	command.SetOut(&out)
	command.SetArgs([]string{"--format", "json"})
	if err := command.Execute(); err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(out.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	var entries []capabilityOperation
	if err := json.Unmarshal(raw["inbox_extensions"], &entries); err != nil || len(entries) != 4 {
		t.Fatalf("inbox_extensions %s", raw["inbox_extensions"])
	}
}

func encodeInboxExcerpt(catalog capabilityCatalog) ([]byte, error) {
	var buffer bytes.Buffer
	err := inbox.Encode(&buffer, struct {
		Kind            string                `json:"kind"`
		InboxExtensions []capabilityOperation `json:"inbox_extensions"`
	}{catalog.Kind, catalog.InboxExtensions})
	return buffer.Bytes(), err
}
