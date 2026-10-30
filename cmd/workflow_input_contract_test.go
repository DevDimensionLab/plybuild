package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/humaninput"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// hiBundle is absolute because fixtures change the working directory.
var hiBundle = func() string {
	p, err := filepath.Abs("../schemas/workflow-human-input-v1")
	if err != nil {
		panic(err)
	}
	return p
}()

// hiWorkspace creates a physical workspace with project ply (repo ply) at root,
// which must be a new directory under /private/tmp outside any workspace.
func hiWorkspace(t *testing.T, root string) (humaninput.Dependencies, string) {
	t.Helper()
	repo := filepath.Join(root, "ply", "main")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init", "-b", "main"}, {"config", "user.name", "HI01"}, {"config", "user.email", "hi01@example.invalid"}, {"commit", "--allow-empty", "-m", "base"}} {
		if out, err := exec.Command("git", append([]string{"-C", repo}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	old, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(old) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	w := workspace.SystemDependencies()
	if _, err := workspace.Init(w); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.AddProject(w, workspace.ProjectAddInput{ProjectID: "ply", Name: "Ply", Wrapper: filepath.Join(root, "ply"), Repos: []workspace.RepoInput{{RepoID: "ply", Path: repo}}}); err != nil {
		t.Fatal(err)
	}
	clock := time.Date(2026, 10, 18, 18, 0, 0, 0, time.UTC)
	return humaninput.Dependencies{Workspace: w, Now: func() time.Time { clock = clock.Add(time.Second); return clock }, Poll: time.Millisecond}, repo
}

// hiRun executes one public command in process and returns stdout and error.
func hiRun(t *testing.T, d humaninput.Dependencies, args ...string) (json.RawMessage, error) {
	t.Helper()
	c := newWorkflowInputCommand(d)
	var out, errout bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&errout)
	c.SetArgs(args)
	c.SilenceErrors = true
	_, err := c.ExecuteC()
	return json.RawMessage(bytes.TrimSpace(out.Bytes())), err
}

func hiField(t *testing.T, raw json.RawMessage, path ...string) string {
	t.Helper()
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatalf("not JSON: %s", raw)
	}
	for _, p := range path {
		v = v.(map[string]any)[p]
	}
	s, _ := v.(string)
	return s
}

type hiExample struct {
	Name       string          `json:"name"`
	Definition string          `json:"definition"`
	Command    *string         `json:"command"`
	Exit       *int            `json:"exit"`
	Value      json.RawMessage `json:"value"`
}

type hiExamples struct {
	Status           string      `json:"status"`
	Schema           string      `json:"schema"`
	Workspace        string      `json:"workspace"`
	Valid            []hiExample `json:"valid"`
	SchemaInvalid    []hiExample `json:"schema_invalid"`
	CanonicalVectors []struct {
		Name          string `json:"name"`
		CanonicalUTF8 string `json:"canonical_utf8"`
		RequestSHA256 string `json:"request_sha256"`
	} `json:"canonical_request_vectors"`
}

// hiJourney drives the actual commands through every public result and
// returns them as labeled examples.
func hiJourney(t *testing.T, d humaninput.Dependencies, root, repo string) []hiExample {
	t.Helper()
	var out []hiExample
	add := func(name, def string, args []string, raw json.RawMessage, err error) json.RawMessage {
		exit := 0
		if err != nil {
			exit = 1
		}
		line := "ply workflow input " + strings.Join(args, " ")
		line = strings.ReplaceAll(line, root, "$WORKSPACE")
		out = append(out, hiExample{name, def, &line, &exit, raw})
		return raw
	}
	run := func(name, def string, args ...string) json.RawMessage {
		raw, err := hiRun(t, d, args...)
		if len(raw) == 0 {
			t.Fatalf("%s printed nothing: %v", name, err)
		}
		return add(name, def, args, raw, err)
	}
	file := func(name string, v any) string {
		p := filepath.Join(root, "inputs", name)
		if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
			t.Fatal(err)
		}
		b, ok := v.([]byte)
		if !ok {
			b, _ = json.Marshal(v)
		}
		if err := os.WriteFile(p, b, 0o600); err != nil {
			t.Fatal(err)
		}
		return p
	}
	ws := []string{"--workspace", root}
	choiceDraft := humaninput.QuestionDraft{Kind: humaninput.DraftKind, SchemaVersion: 1, PublicationKey: "docs-export/format", ProjectID: "ply", RepoID: ptrTo("ply"), TaskID: nil, ActivityID: "conversation/docs-export", ActivityTitle: "Document the export format", Provider: "claude", SessionID: ptrTo("3087ad01-0000-4000-8000-000000000001"), OriginCWD: repo, Text: "Which README section should this task write?", Reason: "The goal leaves the scope of the README change open.", AnswerKind: humaninput.AnswerChoice, Options: []humaninput.Option{{ID: "short", Label: "Short usage note", Consequence: "One paragraph under Usage."}, {ID: "full", Label: "Full reference", Consequence: "A new section with every flag."}}}
	textDraft := choiceDraft
	textDraft.PublicationKey, textDraft.AnswerKind, textDraft.Options, textDraft.Text, textDraft.SessionID = "docs-export/heading", humaninput.AnswerText, []humaninput.Option{}, "What exact heading should the new section use?", nil
	textDraft.Provider = "codex"
	out = append(out, hiExample{"draft-choice", "question_draft", nil, nil, mustRaw(t, choiceDraft)}, hiExample{"draft-text", "question_draft", nil, nil, mustRaw(t, textDraft)})
	choice := run("ask-choice-waiting", "question", append([]string{"ask", "--file", file("choice.json", choiceDraft)}, ws...)...)
	cid := hiField(t, choice, "question_id")
	run("ask-choice-idempotent-replay", "question", append([]string{"ask", "--file", file("choice.json", choiceDraft)}, ws...)...)
	text := run("ask-text-waiting", "question", append([]string{"ask", "--file", file("text.json", textDraft)}, ws...)...)
	tid := hiField(t, text, "question_id")
	changed := choiceDraft
	changed.Text = "Which README section should this task write now?"
	run("error-publication-conflict", "error", append([]string{"ask", "--file", file("changed.json", changed)}, ws...)...)
	run("list-first-page", "list", append([]string{"list", "--project", "ply", "--limit", "1"}, ws...)...)
	page := out[len(out)-1].Value
	run("list-last-page", "list", append([]string{"list", "--project", "ply", "--limit", "1", "--cursor", hiField(t, page, "next_cursor")}, ws...)...)
	run("wait-timeout-still-waiting", "wait", append([]string{"wait", cid, "--timeout-seconds", "1"}, ws...)...)
	run("status-not-recorded", "status", append([]string{"answer", "status", "app-choice-0001", "--question-id", cid}, ws...)...)
	run("error-invalid-answer-option", "error", append([]string{"answer", "preview", cid, "--request-id", "app-choice-0001", "--option-id", "medium"}, ws...)...)
	pc := run("preview-choice", "preview", append([]string{"answer", "preview", cid, "--request-id", "app-choice-0001", "--option-id", "full"}, ws...)...)
	var preview humaninput.Preview
	if err := json.Unmarshal(pc, &preview); err != nil {
		t.Fatal(err)
	}
	out = append(out, hiExample{"request-choice", "request", nil, nil, mustRaw(t, preview.Request)})
	other := run("preview-choice-competing", "preview", append([]string{"answer", "preview", cid, "--request-id", "app-choice-0002", "--option-id", "short"}, ws...)...)
	var competing humaninput.Preview
	_ = json.Unmarshal(other, &competing)
	run("apply-recorded", "status", "answer", "apply", "--workspace", root, "--file", file("request-choice.json", preview.Request))
	run("apply-replay-recorded", "status", "answer", "apply", "--workspace", root, "--file", file("request-choice.json", preview.Request))
	run("apply-other-winner", "status", "answer", "apply", "--workspace", root, "--file", file("request-competing.json", competing.Request))
	conflict := preview.Request
	conflict.Answer.OptionID = ptrTo("short")
	conflict.BindingSHA256 = bindingFor(t, conflict)
	run("apply-workspace-mismatch-not-recorded", "status", "answer", "apply", "--workspace", filepath.Join(root, "ply"), "--file", file("request-choice.json", preview.Request))
	run("apply-same-id-other-bytes", "status", "answer", "apply", "--workspace", root, "--file", file("request-conflict.json", conflict))
	run("status-recorded", "status", append([]string{"answer", "status", "app-choice-0001", "--question-id", cid}, ws...)...)
	run("status-other-winner", "status", append([]string{"answer", "status", "app-choice-0002", "--question-id", cid}, ws...)...)
	run("error-already-answered", "error", append([]string{"answer", "preview", cid, "--request-id", "app-choice-0003", "--option-id", "short"}, ws...)...)
	run("wait-answered", "wait", append([]string{"wait", cid, "--timeout-seconds", "1"}, ws...)...)
	sha := preview.RequestSHA256
	run("acknowledge-received", "question", append([]string{"acknowledge", cid, "--request-id", "app-choice-0001", "--request-sha256", sha}, ws...)...)
	run("continue-reported", "question", append([]string{"continue", cid, "--request-id", "app-choice-0001", "--request-sha256", sha, "--summary", "Wrote the full reference section into README.md."}, ws...)...)
	exact := "Export format\r\n  (v1) \n"
	tp := run("preview-text-exact", "preview", append([]string{"answer", "preview", tid, "--request-id", "app-text-0001", "--text-file", file("answer.txt", []byte(exact))}, ws...)...)
	var tpreview humaninput.Preview
	_ = json.Unmarshal(tp, &tpreview)
	out = append(out, hiExample{"request-text", "request", nil, nil, mustRaw(t, tpreview.Request)})
	// The maximal text answer: 4000 codepoints and 16000 UTF-8 bytes.
	run("preview-text-4000-codepoints-16000-bytes-boundary", "preview", append([]string{"answer", "preview", tid, "--request-id", "app-text-max-0001", "--text-file", file("answer-max.txt", []byte(strings.Repeat("\U0001F600", 4000)))}, ws...)...)
	run("error-invalid-answer-blank-text", "error", append([]string{"answer", "preview", tid, "--request-id", "app-text-0002", "--text-file", file("blank.txt", []byte(" \n\t"))}, ws...)...)
	run("apply-text-recorded", "status", "answer", "apply", "--workspace", root, "--file", file("request-text.json", tpreview.Request))
	run("show-text-answered", "question", append([]string{"show", tid}, ws...)...)
	run("continue-without-receipt", "question", append([]string{"continue", tid, "--request-id", "app-text-0001", "--request-sha256", tpreview.RequestSHA256, "--summary", "Used the exact text answer."}, ws...)...)
	gone := choiceDraft
	gone.PublicationKey = "docs-export/obsolete"
	gq := run("ask-to-withdraw", "question", append([]string{"ask", "--file", file("gone.json", gone)}, ws...)...)
	gid := hiField(t, gq, "question_id")
	gp := run("preview-before-withdrawal", "preview", append([]string{"answer", "preview", gid, "--request-id", "app-gone-0001", "--option-id", "short"}, ws...)...)
	var gpreview humaninput.Preview
	_ = json.Unmarshal(gp, &gpreview)
	run("withdraw", "question", append([]string{"withdraw", gid, "--reason", "The goal already answers this."}, ws...)...)
	run("apply-after-withdrawal-not-recorded", "status", "answer", "apply", "--workspace", root, "--file", file("request-gone.json", gpreview.Request))
	run("wait-withdrawn", "wait", append([]string{"wait", gid, "--timeout-seconds", "1"}, ws...)...)
	run("list-all-states", "list", append([]string{"list", "--project", "ply"}, ws...)...)
	missing := "hiq_" + strings.Repeat("0", 32)
	run("status-unavailable", "status", append([]string{"answer", "status", "app-none-0001", "--question-id", missing}, ws...)...)
	run("error-not-found", "error", append([]string{"show", missing}, ws...)...)
	// Damaged stored state is reported, never repaired.
	answer := filepath.Join(root, ".ply", "human-input", "v1", "answers", tid+".json")
	if err := os.WriteFile(answer, []byte(`{"truncated":`), 0o600); err != nil {
		t.Fatal(err)
	}
	run("status-unknown", "status", append([]string{"answer", "status", "app-text-0001", "--question-id", tid}, ws...)...)
	run("error-state-invalid", "error", append([]string{"show", tid}, ws...)...)
	return out
}

func ptrTo(s string) *string { return &s }
func mustRaw(t *testing.T, v any) json.RawMessage {
	b, err := humaninput.Encode(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func bindingFor(t *testing.T, r humaninput.AnswerRequest) string {
	t.Helper()
	m := map[string]any{}
	if err := json.Unmarshal(mustRaw(t, r), &m); err != nil {
		t.Fatal(err)
	}
	delete(m, "binding_sha256")
	return "sha256:" + sha256Hex(mustRaw(t, m))
}
func sha256Hex(b []byte) string { x := sha256.Sum256(b); return hex.EncodeToString(x[:]) }

func hiSchemaInvalid(t *testing.T, valid []hiExample) []hiExample {
	t.Helper()
	pick := func(name string) map[string]any {
		for _, e := range valid {
			if e.Name == name {
				m := map[string]any{}
				_ = json.Unmarshal(e.Value, &m)
				return m
			}
		}
		t.Fatalf("no example %s", name)
		return nil
	}
	var out []hiExample
	mut := func(name, def, from string, change func(map[string]any)) {
		m := pick(from)
		change(m)
		out = append(out, hiExample{name, def, nil, nil, mustRaw(t, m)})
	}
	mut("question-unknown-field", "question", "ask-choice-waiting", func(m map[string]any) { m["task_title"] = "invented" })
	mut("question-missing-nullable", "question", "ask-choice-waiting", func(m map[string]any) { delete(m, "answer") })
	mut("question-waiting-agent-actor", "question", "ask-choice-waiting", func(m map[string]any) { m["next_action"].(map[string]any)["actor"] = "agent" })
	mut("status-recorded-without-winner", "status", "status-recorded", func(m map[string]any) { m["winner"] = nil })
	mut("status-unknown-with-digest", "status", "status-unknown", func(m map[string]any) { m["request_sha256"] = "sha256:" + strings.Repeat("a", 64) })
	mut("request-missing-binding", "request", "request-choice", func(m map[string]any) { delete(m, "binding_sha256") })
	mut("request-old-purpose", "request", "request-choice", func(m map[string]any) { m["purpose"] = "product_answer_v1" })
	mut("request-text-with-option", "request", "request-text", func(m map[string]any) { m["answer"].(map[string]any)["option_id"] = "short" })
	// JSON Schema maxLength counts codepoints: 4000 is the text limit; the
	// 16000 UTF-8 byte cap is native-only.
	mut("request-text-4001-codepoints", "request", "request-text", func(m map[string]any) {
		m["answer"].(map[string]any)["text"] = strings.Repeat("a", 4001)
	})
	mut("draft-one-option", "question_draft", "draft-choice", func(m map[string]any) { m["options"] = m["options"].([]any)[:1] })
	mut("draft-text-with-options", "question_draft", "draft-text", func(m map[string]any) { m["options"] = pick("draft-choice")["options"] })
	mut("list-limit-21", "list", "list-all-states", func(m map[string]any) { m["limit"] = 21 })
	mut("wait-unknown-outcome", "wait", "wait-answered", func(m map[string]any) { m["outcome"] = "timed_out" })
	return out
}

// TestHI07WriteContractExamples regenerates examples.json from the actual
// commands. It runs only with HI01_WRITE_EXAMPLES=1; the locked digests in
// TestHI07ContractBundleIsLockedAndMatchesActualOutput must then be updated.
func TestHI07WriteContractExamples(t *testing.T) {
	if os.Getenv("HI01_WRITE_EXAMPLES") != "1" {
		t.Skip("set HI01_WRITE_EXAMPLES=1 to regenerate the bundle examples")
	}
	root := "/private/tmp/hi01-example-workspace"
	if err := os.RemoveAll(root); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(root) })
	d, repo := hiWorkspace(t, root)
	valid := hiJourney(t, d, root, repo)
	var boundary humaninput.Preview
	for _, e := range valid {
		if e.Name == "preview-text-4000-codepoints-16000-bytes-boundary" {
			if err := json.Unmarshal(e.Value, &boundary); err != nil {
				t.Fatal(err)
			}
		}
	}
	if text := *boundary.Request.Answer.Text; utf8.RuneCountInString(text) != 4000 || len(text) != 16000 {
		t.Fatalf("boundary text has %d codepoints and %d bytes", utf8.RuneCountInString(text), len(text))
	}
	valid = append(valid, hiExample{"request-text-4000-codepoints-16000-bytes-boundary", "request", nil, nil, mustRaw(t, boundary.Request)})
	ex := hiExamples{Status: "Generated from the actual ply workflow input commands with a fixed clock; the workspace path is a fixture, not a registered workspace.", Schema: "workflow-human-input-v1.schema.json", Workspace: root, Valid: valid, SchemaInvalid: hiSchemaInvalid(t, valid)}
	for _, name := range []string{"request-choice", "request-text"} {
		for _, e := range valid {
			if e.Name == name {
				var r humaninput.AnswerRequest
				_ = json.Unmarshal(e.Value, &r)
				canonical := mustRaw(t, r)
				ex.CanonicalVectors = append(ex.CanonicalVectors, struct {
					Name          string `json:"name"`
					CanonicalUTF8 string `json:"canonical_utf8"`
					RequestSHA256 string `json:"request_sha256"`
				}{name, string(canonical), "sha256:" + sha256Hex(canonical)})
			}
		}
	}
	if errs := hiDigestErrors(t, ex); len(errs) != 0 {
		t.Fatalf("generated examples carry wrong digests:\n%s", strings.Join(errs, "\n"))
	}
	hiValidate(t, ex.Valid, true)
	hiValidate(t, ex.SchemaInvalid, false)
	b, err := json.MarshalIndent(ex, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hiBundle, "examples.json"), append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func hiValidate(t *testing.T, items []hiExample, wantValid bool) {
	t.Helper()
	body, err := json.Marshal(items)
	if err != nil {
		t.Fatal(err)
	}
	c := exec.Command(sl01Python(t), "-I", "-B", "-c", `import json,sys,jsonschema
s=json.load(open(sys.argv[1])); items=json.load(sys.stdin); want=sys.argv[2]=='1'
for item in items:
  sub={'$schema':s['$schema'],'$defs':s['$defs'],'$ref':'#/$defs/'+item['definition']}
  errors=list(jsonschema.Draft202012Validator(sub,format_checker=jsonschema.FormatChecker()).iter_errors(item['value']))
  whole=list(jsonschema.Draft202012Validator(s).iter_errors(item['value']))
  if want and (errors or whole): raise SystemExit(item['name']+': '+str(errors[:1] or whole[:1]))
  if not want and not errors: raise SystemExit(item['name']+' unexpectedly valid')
print(len(items),'checked')
`, filepath.Join(hiBundle, "workflow-human-input-v1.schema.json"), map[bool]string{true: "1", false: "0"}[wantValid])
	c.Stdin = bytes.NewReader(body)
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("schema validation: %v %s", err, out)
	}
}

// HI07: the bundle bytes are locked, its examples validate as labeled, and a
// fresh journey through the actual commands produces only schema-valid output
// of the same definitions.
func TestHI07ContractBundleIsLockedAndMatchesActualOutput(t *testing.T) {
	locked := map[string]string{
		"README.md":                           "f9711f041f15a348df11e2b11e26a96be9751c2bb2d5e59b14d6a07f491d5216",
		"workflow-human-input-v1.schema.json": "202fddcbee1af761b890d1a99a6ad8ea1343be826ea11f71a67836a8a2fbb46e",
		"examples.json":                       "8f21910ae0d5ff7f8311f0e0f79f8f2c53cec55759e90ddb3a8c88835f6944cd",
	}
	for name, want := range locked {
		b, err := os.ReadFile(filepath.Join(hiBundle, name))
		if err != nil {
			t.Fatal(err)
		}
		if got := sha256Hex(b); got != want {
			t.Errorf("%s bytes changed: sha256 %s; update the lock only with a reviewed contract change", name, got)
		}
	}
	raw, err := os.ReadFile(filepath.Join(hiBundle, "examples.json"))
	if err != nil {
		t.Fatal(err)
	}
	var ex hiExamples
	if err := json.Unmarshal(raw, &ex); err != nil {
		t.Fatal(err)
	}
	hiValidate(t, ex.Valid, true)
	hiValidate(t, ex.SchemaInvalid, false)
	for _, v := range ex.CanonicalVectors {
		var r humaninput.AnswerRequest
		if err := json.Unmarshal([]byte(v.CanonicalUTF8), &r); err != nil {
			t.Fatal(err)
		}
		if string(mustRaw(t, r)) != v.CanonicalUTF8 || "sha256:"+sha256Hex([]byte(v.CanonicalUTF8)) != v.RequestSHA256 || humaninput.RequestSHA256(r) != v.RequestSHA256 {
			t.Fatalf("canonical vector %s differs from native canonicalization", v.Name)
		}
	}
	root, err := os.MkdirTemp("/private/tmp", "hi01-contract.")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(root) })
	d, repo := hiWorkspace(t, root)
	all := hiJourney(t, d, root, repo)
	hiValidate(t, all, true)
	var fresh []hiExample
	for _, e := range all {
		if e.Command != nil {
			fresh = append(fresh, e)
		}
	}
	var observed []hiExample
	for _, e := range ex.Valid {
		if e.Command != nil {
			observed = append(observed, e)
		}
	}
	if len(fresh) != len(observed) {
		t.Fatalf("journey has %d results, examples %d", len(fresh), len(observed))
	}
	for i := range fresh {
		if fresh[i].Name != observed[i].Name || fresh[i].Definition != observed[i].Definition || *fresh[i].Exit != *observed[i].Exit {
			t.Fatalf("example %s no longer matches the actual journey", observed[i].Name)
		}
	}
}

// hiDigestErrors checks the RFC8785 digests a strict consumer recomputes:
// every request's binding_sha256, every preview's request_sha256, every
// request_sha256 reported for a known request_id and every canonical vector.
func hiDigestErrors(t *testing.T, ex hiExamples) []string {
	t.Helper()
	var errs []string
	known := map[string]string{}
	digest := func(v any) string { return "sha256:" + sha256Hex(mustRaw(t, v)) }
	request := func(name string, r map[string]any) string {
		unbound := map[string]any{}
		for k, v := range r {
			if k != "binding_sha256" {
				unbound[k] = v
			}
		}
		if got := digest(unbound); r["binding_sha256"] != got {
			errs = append(errs, name+": binding_sha256 "+r["binding_sha256"].(string)+" != "+got)
		}
		sha := digest(r)
		id := r["request_id"].(string)
		if prior, ok := known[id]; ok && prior != sha {
			errs = append(errs, name+": request_id "+id+" names two requests")
		}
		known[id] = sha
		return sha
	}
	values := map[string]map[string]any{}
	for _, e := range ex.Valid {
		m := map[string]any{}
		if err := json.Unmarshal(e.Value, &m); err != nil {
			t.Fatal(err)
		}
		values[e.Name] = m
		switch e.Definition {
		case "request":
			request(e.Name, m)
		case "preview":
			if sha := request(e.Name, m["request"].(map[string]any)); m["request_sha256"] != sha {
				errs = append(errs, e.Name+": request_sha256 does not match its request")
			}
		}
	}
	var reported func(name string, v any)
	reported = func(name string, v any) {
		switch x := v.(type) {
		case map[string]any:
			if id, ok := x["request_id"].(string); ok {
				if sha, ok := x["request_sha256"].(string); ok && known[id] != "" && sha != known[id] {
					errs = append(errs, name+": request_sha256 for "+id+" does not match its request")
				}
			}
			for _, y := range x {
				reported(name, y)
			}
		case []any:
			for _, y := range x {
				reported(name, y)
			}
		}
	}
	for _, e := range ex.Valid {
		reported(e.Name, values[e.Name])
	}
	for _, v := range ex.CanonicalVectors {
		m, ok := values[v.Name]
		if !ok || string(mustRaw(t, m)) != v.CanonicalUTF8 || digest(m) != v.RequestSHA256 {
			errs = append(errs, "canonical vector "+v.Name+" does not match its valid request")
		}
	}
	return errs
}

// HI01A: schema validity is not enough; a strict consumer recomputes every
// digest of a declared-valid request and must accept them all.
func TestHI07ValidRequestExamplesCarryNativeDigests(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(hiBundle, "examples.json"))
	if err != nil {
		t.Fatal(err)
	}
	var ex hiExamples
	if err := json.Unmarshal(raw, &ex); err != nil {
		t.Fatal(err)
	}
	if errs := hiDigestErrors(t, ex); len(errs) != 0 {
		t.Fatalf("declared-valid examples carry wrong digests:\n%s", strings.Join(errs, "\n"))
	}
	// The delivered HI01 counterexample: request-text with other text and a
	// stale binding stays schema-valid, but must be rejected here.
	for i, e := range ex.Valid {
		if e.Name == "request-text-4000-codepoints-16000-bytes-boundary" {
			m := map[string]any{}
			_ = json.Unmarshal(e.Value, &m)
			m["answer"].(map[string]any)["text"] = "Export format"
			stale := ex
			stale.Valid = append(append([]hiExample{}, ex.Valid[:i]...), hiExample{e.Name, e.Definition, nil, nil, mustRaw(t, m)})
			if errs := hiDigestErrors(t, stale); len(errs) == 0 {
				t.Fatal("a stale binding_sha256 was accepted")
			}
			return
		}
	}
	t.Fatal("no maximal text request example")
}

// Argument errors print nothing on stdout; semantic refusals print one result.
func TestHI07ArgumentErrorsPrintNoResult(t *testing.T) {
	root, err := os.MkdirTemp("/private/tmp", "hi01-args.")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(root) })
	d, _ := hiWorkspace(t, root)
	id := "hiq_" + strings.Repeat("0", 32)
	for _, args := range [][]string{
		{"show", "bad-id", "--workspace", root},
		{"show", id},
		{"show", id, "--workspace", "relative"},
		{"list", "--workspace", root},
		{"list", "--workspace", root, "--project", "ply", "--limit", "0"},
		{"wait", id, "--workspace", root, "--timeout-seconds", "601"},
		{"wait", id, "--workspace", root},
		{"answer", "preview", id, "--workspace", root, "--request-id", "app-0000001"},
		{"answer", "preview", id, "--workspace", root, "--request-id", "app-0000001", "--option-id", "a", "--text-file", "/x"},
		{"answer", "status", "app-0000001", "--workspace", root},
		{"continue", id, "--workspace", root, "--request-id", "app-0000001", "--request-sha256", "sha256:" + strings.Repeat("a", 64)},
		{"ask", "--workspace", root},
		{"answer", "apply"},
		{"answer", "apply", "--file", filepath.Join(root, "request.json")},
		{"unknown"},
	} {
		out, err := hiRun(t, d, args...)
		if err == nil || len(out) != 0 || !strings.HasPrefix(err.Error(), humaninput.CodeInvalidArguments) {
			t.Fatalf("%v: out=%s err=%v", args, out, err)
		}
	}
	if _, err := os.Lstat(filepath.Join(root, ".ply", "human-input")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("argument errors created human-input state")
	}
}
