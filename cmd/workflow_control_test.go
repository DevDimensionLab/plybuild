package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

type mr05CLI struct {
	t                       *testing.T
	binary, root, home, tmp string
}

func mr05Candidate(t *testing.T) *mr05CLI {
	t.Helper()
	bin := os.Getenv("PLY_CANDIDATE_BINARY")
	if bin == "" || !filepath.IsAbs(bin) {
		t.Fatal("PLY_CANDIDATE_BINARY is mandatory; no build or fallback is permitted")
	}
	if s, e := os.Stat(bin); e != nil || s.Mode().Perm()&0111 == 0 {
		t.Fatal("exact candidate is not executable", e)
	}
	r, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	if e = os.Chmod(r, 0700); e != nil {
		t.Fatal(e)
	}
	f := &mr05CLI{t, bin, r, filepath.Join(r, "home"), filepath.Join(r, "tmp")}
	for _, p := range []string{f.home, f.tmp} {
		if e = os.Mkdir(p, 0700); e != nil {
			t.Fatal(e)
		}
	}
	return f
}
func (f *mr05CLI) run(args ...string) ([]byte, string, error) {
	f.t.Helper()
	c := exec.Command(f.binary, args...)
	c.Dir = f.root
	c.Env = []string{"HOME=" + f.home, "TMPDIR=" + f.tmp, "PATH=" + os.Getenv("PATH"), "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1"}
	var errout bytes.Buffer
	c.Stderr = &errout
	out, e := c.Output()
	return out, errout.String(), e
}
func mr05Digest(v any) string {
	b, e := taskrun.Canonical(v)
	if e != nil {
		panic(e)
	}
	h := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(h[:])
}
func (f *mr05CLI) oldRun() string {
	r := taskrun.WorkflowRequest{Envelope: taskrun.Envelope{Kind: "ply.workflow.herdr-run-request", SchemaVersion: 1}, WorkspaceRoot: f.root, RequestKey: "legacy-control-fixture", Runtime: taskrun.Runtime{Mode: "interactive", Provider: "codex"}, HandoffDraft: json.RawMessage(`{"private_secret":"must-not-leak-control-fixture"}`)}
	mr05EmptySlices(reflect.ValueOf(&r).Elem())
	id := "wfr_" + strings.TrimPrefix(mr05Digest([]string{f.root, r.RequestKey}), "sha256:")
	run := taskrun.WorkflowRun{RunID: id, SessionID: "ply:" + id, RequestSHA256: mr05Digest(r), Paths: taskrun.WorkflowPaths{RunRoot: filepath.Join(f.root, ".ply", "workflow-runs", "runs", id)}}
	mr05EmptySlices(reflect.ValueOf(&run).Elem())
	observed := taskrun.Observed{}
	mr05EmptySlices(reflect.ValueOf(&observed).Elem())
	p := filepath.Join(f.root, ".ply", "workflow-runs", "requests", id+".json")
	if e := os.MkdirAll(filepath.Dir(p), 0700); e != nil {
		f.t.Fatal(e)
	}
	b, _ := taskrun.Canonical(map[string]any{"request": r, "observed": observed, "result": run, "phase": "", "context_sha256": "", "acceptance": nil, "start_draft": nil, "start_sha256": nil, "records": []any{}})
	if e := os.WriteFile(p, b, 0600); e != nil {
		f.t.Fatal(e)
	}
	return id
}

func mr05EmptySlices(v reflect.Value) {
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			mr05EmptySlices(v.Field(i))
		}
	case reflect.Slice:
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 0))
		}
	}
}
func mr05Snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	e := filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		i, e := d.Info()
		if e != nil {
			return e
		}
		v := i.Mode().String()
		if !d.IsDir() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			h := sha256.Sum256(b)
			v += hex.EncodeToString(h[:])
		}
		out[p] = v
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	return out
}
func mr05Schema(t *testing.T, kind string, raw []byte) {
	t.Helper()
	py := os.Getenv("PLY_MR05_SCHEMA_PYTHON")
	if py == "" {
		t.Fatal("PLY_MR05_SCHEMA_PYTHON is mandatory; no skip or install")
	}
	schema, e := filepath.Abs("../schemas/workflow-control-v1.schema.json")
	if e != nil {
		t.Fatal(e)
	}
	c := exec.Command(py, "-I", "-B", "-c", `import json,sys,jsonschema
s=json.load(open(sys.argv[1])); v=json.load(sys.stdin)
jsonschema.Draft202012Validator(s,format_checker=jsonschema.FormatChecker()).validate(v)
assert v['kind']==sys.argv[2]
`, schema, kind)
	c.Stdin = bytes.NewReader(raw)
	if out, e := c.CombinedOutput(); e != nil {
		t.Fatalf("actual response violates schema: %v %s", e, out)
	}
}

func TestMR05ControlCandidateCLIAndSchema(t *testing.T) {
	f := mr05Candidate(t)
	raw, stderr, e := f.run("capabilities", "--format", "json")
	if e != nil {
		t.Fatal(e, stderr)
	}
	var caps struct {
		Operations []struct {
			ID, Mode, Effect string
			ResultSchemas    []struct {
				Kind string `json:"kind"`
			} `json:"result_schemas"`
		} `json:"workflow_extensions"`
	}
	if e = json.Unmarshal(raw, &caps); e != nil {
		t.Fatal(e)
	}
	for _, mode := range []string{"codex_headless_standalone_v1", "codex_headless_week_v1"} {
		for _, verb := range []string{"show", "preview", "apply", "status"} {
			found := 0
			for _, op := range caps.Operations {
				if op.ID == "workflow.control."+verb && op.Mode == mode {
					found++
					effect := "read"
					if verb == "apply" {
						effect = "approved_runtime_control"
					}
					if op.Effect != effect || len(op.ResultSchemas) == 0 {
						t.Fatal("wrong operation contract", op)
					}
				}
			}
			if found != 1 {
				t.Fatalf("missing/duplicate %s %s", mode, verb)
			}
		}
	}
	for _, tc := range []struct {
		args []string
		want string
	}{{[]string{"workflow", "control", "--help"}, "request-bound"}, {[]string{"workflow", "execute", "--help"}, "--bound-control"}, {[]string{"workflow", "week", "start", "--help"}, "runtime_control"}} {
		out, errout, e := f.run(tc.args...)
		if e != nil || !strings.Contains(string(out), tc.want) {
			t.Fatal("missing opt-in/help contract", tc.args, string(out), errout, e)
		}
	}
	id := f.oldRun()
	raw, stderr, e = f.run("workflow", "control", "show", id, "--workspace", f.root, "--format", "json")
	if e != nil {
		t.Fatal(e, stderr)
	}
	mr05Schema(t, "ply.workflow.control-capabilities", raw)
	raw, stderr, e = f.run("workflow", "control", "status", "unrecorded-0001", "--workspace", f.root, "--format", "json")
	if e != nil {
		t.Fatal(e, stderr)
	}
	mr05Schema(t, "ply.workflow.control-status", raw)
	raw, stderr, e = f.run("workflow", "control", "preview", id, "--action", "pause", "--request-id", "preview-0001", "--workspace", f.root, "--format", "json")
	if e == nil || !strings.Contains(stderr, "control_unavailable") {
		t.Fatal("legacy preview offered approval", e, stderr)
	}
	mr05Schema(t, "ply.workflow.control-preview", raw)
	py := os.Getenv("PLY_MR05_SCHEMA_PYTHON")
	script, e := filepath.Abs("../test/mr05-schema-examples.py")
	if e != nil {
		t.Fatal(e)
	}
	if out, e := exec.Command(py, "-I", "-B", script).CombinedOutput(); e != nil {
		t.Fatal("valid/invalid fixtures", e, string(out))
	}
}

func TestMR05ControlCandidateCLIReadOnlyErrors(t *testing.T) {
	t.Run("unavailable_week_json", mr05UnavailableWeekJSON)
	t.Run("public_observation_states", mr05PublicObservationStates)
	f := mr05Candidate(t)
	id := f.oldRun()
	before := mr05Snapshot(t, f.root)
	for _, args := range [][]string{{"workflow", "control", "show", id, "--week", "week_" + strings.Repeat("a", 32), "--workspace", f.root}, {"workflow", "control", "preview", id, "--action", "kill", "--request-id", "malformed-0001", "--workspace", f.root}, {"workflow", "control", "apply", "--file", filepath.Join(f.root, "missing.json")}, {"workflow", "control", "show", "--workspace", f.root}, {"workflow", "control", "status", "bad", "--workspace", f.root}} {
		out, errout, e := f.run(args...)
		if e == nil || !strings.Contains(errout, "control_") {
			t.Fatal("bad arguments lack stable error", args, e, errout)
		}
		if strings.Contains(string(out)+errout, "must-not-leak-control-fixture") {
			t.Fatal("secret leaked")
		}
	}
	for i := 0; i < 3; i++ {
		for _, args := range [][]string{{"workflow", "control", "show", id, "--workspace", f.root}, {"workflow", "control", "status", "missing-request-0001", "--workspace", f.root}, {"workflow", "control", "preview", id, "--workspace", f.root, "--action", "pause", "--request-id", "old-preview-0001"}} {
			out, errout, _ := f.run(args...)
			if strings.Contains(string(out)+errout, "must-not-leak-control-fixture") {
				t.Fatal("read leaked secret")
			}
		}
	}
	if !reflect.DeepEqual(before, mr05Snapshot(t, f.root)) {
		t.Fatal("read-only candidate commands wrote files/config")
	}
	valid := `{"kind":"ply.workflow.control-request","schema_version":1,"request_id":"strict-request-0001","workspace":` + string(mustMR05JSON(f.root)) + `,"input":{"target":{"kind":"run","id":"` + id + `"},"action":"pause"},"validity":{"issued_at":"2026-10-10T08:00:00Z","expires_at":"2026-10-10T08:05:00Z"},"binding":"sha256:` + strings.Repeat("a", 64) + `"}`
	for _, bad := range []string{strings.Replace(valid, `"schema_version":1`, `"schema_version":99`, 1), strings.Replace(valid, `"schema_version":1`, `"schema_version":1,"schema_version":1`, 1), strings.Replace(valid, `"schema_version":1`, `"schema_version":"1"`, 1), strings.Replace(valid, `"action":"pause"`, `"action":"pause","extra":true`, 1), strings.Replace(valid, `"kind":"run"`, `"kind":"RUN"`, 1), strings.Replace(valid, `"request_id":"strict-request-0001"`, `"request_id":null`, 1)} {
		p := filepath.Join(f.root, "bad.json")
		if e := os.WriteFile(p, []byte(bad), 0600); e != nil {
			t.Fatal(e)
		}
		snapshot := mr05Snapshot(t, f.root)
		_, errout, e := f.run("workflow", "control", "apply", "--file", p)
		if e == nil || !strings.Contains(errout, "control_invalid_request") {
			t.Fatal("malformed request accepted", e, errout)
		}
		if !reflect.DeepEqual(snapshot, mr05Snapshot(t, f.root)) {
			t.Fatal("rejected parse mutated native state")
		}
	}
}
func mustMR05JSON(v any) []byte {
	b, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return b
}
