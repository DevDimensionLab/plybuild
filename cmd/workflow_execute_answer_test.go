package cmd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestSL01AnswerCLICommandsAndCapabilities(t *testing.T) {
	for _, path := range [][]string{{"question"}, {"answer", "preview"}, {"answer", "apply"}, {"answer", "status"}} {
		t.Run(strings.Join(path, "/"), func(t *testing.T) {
			command := newWorkflowExecuteCommand(taskrun.Dependencies{})
			found, remaining, err := command.Find(path)
			if err != nil || len(remaining) != 0 || found.Name() != path[len(path)-1] {
				t.Fatalf("required public command %v is absent: remaining=%v err=%v", path, remaining, err)
			}
			var out bytes.Buffer
			found.SetOut(&out)
			if err := found.Help(); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(out.String(), "same-owner") || !strings.Contains(out.String(), "--format") {
				t.Fatalf("help omits same-owner contract or output selector: %s", out.String())
			}
		})
	}
	want := map[string]string{
		"workflow.execute.product-question":    "read",
		"workflow.execute.bound-answer":        "read",
		"workflow.execute.bound-answer.apply":  "approved_product_answer",
		"workflow.execute.bound-answer.status": "read",
	}
	counts := map[string]int{}
	for _, op := range coreCapabilities(capabilityBuild{}).WorkflowExtensions {
		effect, selected := want[op.ID]
		if !selected {
			continue
		}
		counts[op.ID]++
		if op.Mode != "week_codex_same_owner_turn_v1" || op.Effect != effect || len(op.ResultSchemas) == 0 {
			t.Errorf("incorrect bound-answer capability: %#v", op)
		}
		for _, schema := range op.ResultSchemas {
			if schema.SchemaVersion == nil || *schema.SchemaVersion != 2 {
				t.Errorf("answer capability advertises an incompatible schema: %#v", schema)
			}
		}
	}
	for id := range want {
		if counts[id] != 1 {
			t.Errorf("capability %s occurs %d times; want one", id, counts[id])
		}
	}
}

func sl01Python(t *testing.T) string {
	t.Helper()
	py := os.Getenv("PLY_SL01_SCHEMA_PYTHON")
	if py == "" {
		py = os.Getenv("PLY_MR05_SCHEMA_PYTHON")
	}
	if py == "" || !filepath.IsAbs(py) {
		t.Fatal("PLY_SL01_SCHEMA_PYTHON (or existing PLY_MR05_SCHEMA_PYTHON) must name an absolute Python with jsonschema; no install or skip")
	}
	return py
}

func sl01ValidateSchema(t *testing.T, schema string, raw []byte) {
	t.Helper()
	sl01ValidateSchemaDefinition(t, schema, "", raw)
}

func sl01ValidateSchemaDefinition(t *testing.T, schema, definition string, raw []byte) {
	t.Helper()
	c := exec.Command(sl01Python(t), "-I", "-B", "-c", `import json,sys,jsonschema
s=json.load(open(sys.argv[1])); v=json.load(sys.stdin)
if sys.argv[2]: s={'$schema':s['$schema'],'$defs':s['$defs'],'$ref':'#/$defs/'+sys.argv[2]}
jsonschema.Draft202012Validator(s,format_checker=jsonschema.FormatChecker()).validate(v)
`, schema, definition)
	c.Stdin = bytes.NewReader(raw)
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("actual public output violates schema: %v %s", err, out)
	}
}

func TestSL01AnswerCLILockedSchemaAndVectors(t *testing.T) {
	base := "../schemas/workflow-bound-answer-v2"
	for name, want := range map[string]string{
		"README.md":                            "8b972f65ce0b4ea72299d802309cc3c5ba254874d2460007af8b5bdfcc996a41",
		"workflow-bound-answer-v2.schema.json": "bba59f011f8d1217cd183304ecdb14a9311801289e3102e69aa91f8feb977075",
		"examples.json":                        "6f994001df8bac2f054c0014a4d2b68dbe32a08fc3de9835f4e99784065d3dd0",
	} {
		body, err := os.ReadFile(filepath.Join(base, name))
		if err != nil {
			t.Fatal(err)
		}
		hash := sha256.Sum256(body)
		if hex.EncodeToString(hash[:]) != want {
			t.Fatalf("locked wire bytes changed: %s", name)
		}
	}
	c := exec.Command(sl01Python(t), "-I", "-B", "-c", `import json,sys,jsonschema
s=json.load(open(sys.argv[1])); x=json.load(open(sys.argv[2]))
jsonschema.Draft202012Validator.check_schema(s)
for group in ['valid','invalid']:
 for item in x[group]:
  sub={'$schema':s['$schema'],'$defs':s['$defs'],'$ref':'#/$defs/'+item['definition']}
  errors=list(jsonschema.Draft202012Validator(sub,format_checker=jsonschema.FormatChecker()).iter_errors(item['value']))
  assert bool(errors)==(group=='invalid'),(group,item['name'],str(errors))
print('16 valid and 14 schema-invalid vectors checked')
`, filepath.Join(base, "workflow-bound-answer-v2.schema.json"), filepath.Join(base, "examples.json"))
	if out, err := c.CombinedOutput(); err != nil {
		t.Fatalf("locked schema vectors: %v %s", err, out)
	}
	var examples struct {
		CanonicalRequestUTF8   string `json:"canonical_request_utf8"`
		CanonicalRequestSHA256 string `json:"canonical_request_sha256"`
		Valid                  []struct {
			Name       string
			Definition string
			Value      json.RawMessage
		} `json:"valid"`
	}
	body, err := os.ReadFile(filepath.Join(base, "examples.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(body, &examples); err != nil {
		t.Fatal(err)
	}
	checked := false
	for _, item := range examples.Valid {
		var typed any
		switch item.Definition {
		case "question_read":
			typed = &taskrun.ProductQuestionReadback{}
		case "preview":
			typed = &taskrun.BoundAnswerPreview{}
		case "request":
			typed = &taskrun.BoundAnswerRequest{}
		case "status":
			typed = &taskrun.BoundAnswerStatus{}
		default:
			t.Fatalf("unhandled wire definition: %s", item.Definition)
		}
		decoder := json.NewDecoder(bytes.NewReader(item.Value))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(typed); err != nil {
			t.Fatalf("public type rejects valid wire %s: %v", item.Name, err)
		}
		typedCanonical, err := taskrun.Canonical(typed)
		if err != nil {
			t.Fatal(err)
		}
		wireCanonical, err := taskrun.Canonical(item.Value)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(typedCanonical, wireCanonical) {
			t.Fatalf("public type changed locked fields/nulls: %s", item.Name)
		}
		if item.Name != "request-complete" {
			continue
		}
		var value any
		decoder = json.NewDecoder(bytes.NewReader(item.Value))
		decoder.UseNumber()
		if err := decoder.Decode(&value); err != nil {
			t.Fatal(err)
		}
		canonical, err := taskrun.Canonical(value)
		if err != nil {
			t.Fatal(err)
		}
		hash := sha256.Sum256(canonical)
		if string(canonical) != examples.CanonicalRequestUTF8 || "sha256:"+hex.EncodeToString(hash[:]) != examples.CanonicalRequestSHA256 {
			t.Fatal("native canonical request differs from the locked consumer vector")
		}
		checked = true
	}
	if !checked {
		t.Fatal("canonical request fixture missing")
	}
	catalog, err := json.Marshal(coreCapabilities(capabilityBuild{}))
	if err != nil {
		t.Fatal(err)
	}
	sl01ValidateSchema(t, "../schemas/read-contract/ply-capabilities-v1.schema.json", catalog)
	phase := taskrun.PublicControlState("awaiting_product_answer")
	if phase != "awaiting_product_answer" {
		t.Fatalf("question wait must remain distinct from global pause: %s", phase)
	}
	observation, err := json.Marshal(taskrun.ControlObservation{State: phase, Source: "worker_state", Freshness: "current", ObservedAt: "2026-10-16T18:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	sl01ValidateSchemaDefinition(t, "../schemas/workflow-control-v1.schema.json", "observation", observation)
}

func sl01Candidate(t *testing.T) *mr05CLI {
	t.Helper()
	bin := ""
	for _, key := range []string{"PLY_DELIVERY_TEST_BINARY", "PLY_TASK_RUN_TEST_BINARY", "PLY_CANDIDATE_BINARY"} {
		if next := os.Getenv(key); next != "" {
			if bin != "" && bin != next {
				t.Fatalf("candidate executable bindings disagree: %s", key)
			}
			bin = next
		}
	}
	if bin == "" || !filepath.IsAbs(bin) {
		t.Fatal("an exact packaged CLI is mandatory via PLY_DELIVERY_TEST_BINARY, PLY_TASK_RUN_TEST_BINARY or PLY_CANDIDATE_BINARY; no implicit build")
	}
	info, err := os.Stat(bin)
	if err != nil || !info.Mode().IsRegular() || info.Mode().Perm()&0111 == 0 {
		t.Fatalf("candidate is not an executable regular file: %v", err)
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	f := &mr05CLI{t: t, binary: bin, root: root, home: filepath.Join(root, "home"), tmp: filepath.Join(root, "tmp")}
	for _, p := range []string{f.home, f.tmp} {
		if err := os.Mkdir(p, 0700); err != nil {
			t.Fatal(err)
		}
	}
	return f
}

func sl01Run(t *testing.T, f *mr05CLI, args ...string) ([]byte, string, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, f.binary, args...)
	c.Dir = f.root
	c.Env = []string{"HOME=" + f.home, "TMPDIR=" + f.tmp, "PATH=" + os.Getenv("PATH"), "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1"}
	var stderr bytes.Buffer
	c.Stderr = &stderr
	out, err := c.Output()
	if ctx.Err() != nil {
		t.Fatalf("CLI did not finish within its read-only/process test bound: %v", args)
	}
	return out, stderr.String(), err
}

func TestSL01AnswerCLICandidateReadOnlyAndSemanticErrors(t *testing.T) {
	f := sl01Candidate(t)
	run := f.oldRun()
	before := mr05Snapshot(t, f.root)
	schema := "../schemas/workflow-bound-answer-v2/workflow-bound-answer-v2.schema.json"
	for _, args := range [][]string{
		{"workflow", "execute", "question", run, "--workspace", f.root, "--format", "json"},
		{"workflow", "execute", "answer", "preview", run, "--workspace", f.root, "--question-id", "pqw_" + strings.Repeat("2", 32), "--question-revision", "1", "--option-id", "task_first", "--request-id", "sl01-cli-read-001", "--format", "json"},
		{"workflow", "execute", "answer", "status", "sl01-cli-missing-001", "--workspace", f.root, "--format", "json"},
	} {
		out, stderr, err := sl01Run(t, f, args...)
		if len(out) == 0 {
			t.Fatalf("public read lost structured status: %v %s", err, stderr)
		}
		sl01ValidateSchema(t, schema, out)
		var value map[string]any
		if err := json.Unmarshal(out, &value); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(out), "must-not-leak-control-fixture") {
			t.Fatal("private source leaked in public answer readback")
		}
		if value["kind"] == "ply.workflow.bound-answer-status" {
			if value["state"] != "not_recorded" || value["answer_recorded"] != nil || value["execution_effect"] != "unknown" {
				t.Fatalf("absence falsely establishes no effect: %s", out)
			}
		} else if value["eligibility"] != "unsupported" {
			t.Fatalf("old Run was implicitly upgraded: %s", out)
		}
	}
	if after := mr05Snapshot(t, f.root); !reflect.DeepEqual(before, after) {
		t.Fatal("read/preview/status changed fixture files")
	}
	request := sl01RequestFixture(t)
	request["workspace"], request["run_id"] = f.root, run
	file := filepath.Join(f.root, "approved-request.json")
	body, _ := json.Marshal(request)
	if err := os.WriteFile(file, body, 0600); err != nil {
		t.Fatal(err)
	}
	before = mr05Snapshot(t, f.root)
	out, stderr, err := sl01Run(t, f, "workflow", "execute", "answer", "apply", "--file", file, "--format", "json")
	if err == nil || !strings.Contains(stderr, "bound_answer_") {
		t.Fatalf("unapproved/old basis was not rejected: %v %s %s", err, stderr, out)
	}
	if len(out) == 0 {
		t.Fatalf("semantic rejection lost known request identity: %s", stderr)
	}
	sl01ValidateSchema(t, schema, out)
	var value map[string]any
	if err := json.Unmarshal(out, &value); err != nil {
		t.Fatal(err)
	}
	if value["request_id"] != request["request_id"] || value["state"] != "rejected" {
		t.Fatalf("rejection lost exact identity: %s", out)
	}
	if after := mr05Snapshot(t, f.root); !reflect.DeepEqual(before, after) {
		t.Fatal("rejected first apply changed fixture")
	}
}

func sl01RequestFixture(t *testing.T) map[string]any {
	t.Helper()
	body, err := os.ReadFile("../schemas/workflow-bound-answer-v2/examples.json")
	if err != nil {
		t.Fatal(err)
	}
	var x struct {
		Valid []struct {
			Name  string
			Value map[string]any
		}
	}
	if err := json.Unmarshal(body, &x); err != nil {
		t.Fatal(err)
	}
	for _, item := range x.Valid {
		if item.Name == "request-complete" {
			return item.Value
		}
	}
	t.Fatal("request fixture missing")
	return nil
}

func TestSL01AnswerCLICandidateRejectsMalformedRequests(t *testing.T) {
	f := sl01Candidate(t)
	request := sl01RequestFixture(t)
	request["workspace"] = f.root
	body, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	mutate := func(fn func(map[string]any)) []byte {
		var v map[string]any
		if err := json.Unmarshal(body, &v); err != nil {
			t.Fatal(err)
		}
		fn(v)
		x, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return x
	}
	cases := map[string][]byte{
		"unknown":            mutate(func(v map[string]any) { v["command"] = "resume" }),
		"nested_unknown":     mutate(func(v map[string]any) { v["validity"].(map[string]any)["unexpected"] = true }),
		"missing":            mutate(func(v map[string]any) { delete(v, "binding") }),
		"null":               mutate(func(v map[string]any) { v["binding"] = nil }),
		"nested_null":        mutate(func(v map[string]any) { v["validity"] = nil }),
		"old_version":        mutate(func(v map[string]any) { v["schema_version"] = 1 }),
		"old_mode":           mutate(func(v map[string]any) { v["mode"] = "week_codex_same_thread_v1" }),
		"wrong_action":       mutate(func(v map[string]any) { v["action"] = "resume" }),
		"relative_workspace": mutate(func(v map[string]any) { v["workspace"] = "relative" }),
		"dot_workspace":      mutate(func(v map[string]any) { v["workspace"] = f.root + "/./home" }),
		"bad_time":           mutate(func(v map[string]any) { v["validity"].(map[string]any)["expires_at_utc"] = "2026-99-99T18:05:00Z" }),
		"boolean_revision":   mutate(func(v map[string]any) { v["question_revision"] = true }),
		"duplicate":          []byte(strings.Replace(string(body), `"request_id":`, `"request_id":"different-request-001","request_id":`, 1)),
		"duplicate_nested":   []byte(strings.Replace(string(body), `"issued_at_utc":`, `"issued_at_utc":"2026-10-16T18:00:00Z","issued_at_utc":`, 1)),
		"fraction":           []byte(strings.Replace(string(body), `"question_revision":1`, `"question_revision":1.0`, 1)),
		"exponent":           []byte(strings.Replace(string(body), `"question_revision":1`, `"question_revision":1e0`, 1)),
		"negative_zero":      []byte(strings.Replace(string(body), `"question_revision":1`, `"question_revision":-0`, 1)),
		"bom":                append([]byte{0xef, 0xbb, 0xbf}, body...),
		"invalid_utf8":       append(append([]byte{}, body...), 0xff),
		"surrogate":          []byte(strings.Replace(string(body), `"task_first"`, `"\ud800"`, 1)),
		"two_values":         append(append([]byte{}, body...), []byte(" {}")...),
		"too_large":          []byte(strings.Repeat(" ", 64<<10) + string(body)),
	}
	for name, raw := range cases {
		t.Run(name, func(t *testing.T) {
			file := filepath.Join(f.root, name+".json")
			if err := os.WriteFile(file, raw, 0600); err != nil {
				t.Fatal(err)
			}
			before := mr05Snapshot(t, f.root)
			out, stderr, err := sl01Run(t, f, "workflow", "execute", "answer", "apply", "--file", file, "--format", "json")
			if err == nil || !strings.Contains(stderr, "bound_answer_invalid_request") || len(out) != 0 {
				t.Fatalf("malformed request produced effect claims or lacked stable error: %v stdout=%s stderr=%s", err, out, stderr)
			}
			if after := mr05Snapshot(t, f.root); !reflect.DeepEqual(before, after) {
				t.Fatal("malformed request changed files")
			}
		})
	}
	regular := filepath.Join(f.root, "regular.json")
	if err := os.WriteFile(regular, body, 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(f.root, "request-link.json")
	if err := os.Symlink(regular, link); err != nil {
		t.Fatal(err)
	}
	ancestor := filepath.Join(f.root, "directory-link")
	if err := os.Symlink(f.root, ancestor); err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"regular.json", link, filepath.Join(ancestor, "regular.json"), f.home, filepath.Join(f.root, "not-created.json")} {
		out, stderr, err := sl01Run(t, f, "workflow", "execute", "answer", "apply", "--file", file, "--format", "json")
		if err == nil || !strings.Contains(stderr, "bound_answer_invalid_request") || len(out) != 0 {
			t.Fatalf("unsafe file accepted: %s err=%v stdout=%s stderr=%s", file, err, out, stderr)
		}
	}
}

func TestSL01AnswerCLIRejectsArgumentsBeforeWorkspaceAccess(t *testing.T) {
	run := "wfr_" + strings.Repeat("1", 64)
	for _, args := range [][]string{
		{"question"}, {"question", run}, {"question", run, "--workspace", "relative"},
		{"question", run, "--workspace", "/private/tmp/x/../y"},
		{"question", run, "--workspace", "/private/tmp", "--format", "yaml"},
		{"question", run, "--workspace", "/private/tmp", "--unknown"},
		{"answer"}, {"answer", "unexpected"}, {"answer", "preview", run},
		{"answer", "preview", run, "--workspace", "/private/tmp", "--question-revision", "2"},
		{"answer", "apply"}, {"answer", "apply", "unexpected", "--file", "/not-created"},
		{"answer", "status"}, {"answer", "status", "bad", "--workspace", "/private/tmp"},
		{"answer", "status", "request-001", "--workspace", "/private/tmp", "--format", "yaml"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			command := newWorkflowExecuteCommand(taskrun.Dependencies{})
			var out, stderr bytes.Buffer
			command.SetOut(&out)
			command.SetErr(&stderr)
			command.SetArgs(args)
			err := command.Execute()
			if err == nil || !strings.Contains(err.Error(), "bound_answer_invalid_arguments") {
				t.Fatalf("invalid input lacks stable answer error: %v; stdout=%q", err, out.String())
			}
			if out.Len() != 0 {
				t.Fatalf("invalid input fabricated stdout: %s", out.String())
			}
		})
	}
}
