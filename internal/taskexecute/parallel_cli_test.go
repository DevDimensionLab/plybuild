package taskexecute

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Candidate CLI calls are read/write operations in private fixture registries.
// Startup effects use the injected fake-Herdr/native-observer seam exclusively.
func TestParallelPrivateCandidateCLIContracts(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires exact private make-built PLY_DELIVERY_TEST_BINARY")
	}
	if !filepath.IsAbs(binary) {
		t.Fatal("no installed/stale CLI fallback")
	}
	source, e := filepath.Abs("../..")
	if e != nil {
		t.Fatal(e)
	}
	python := os.Getenv("PLY_PARALLEL_SCHEMA_PYTHON")
	if !filepath.IsAbs(python) {
		t.Fatal("select an absolute Python with pinned jsonschema4.26.0")
	}
	d, in, root := parallelLauncherFixture(t, binary)
	n := 0
	call := func(want int, args ...string) []byte {
		t.Helper()
		n++
		c := exec.Command(binary, args...)
		c.Dir = filepath.Join(root, "epic")
		var stdout, stderr bytes.Buffer
		c.Stdout = &stdout
		c.Stderr = &stderr
		err := c.Run()
		exit := 0
		if err != nil {
			var ok bool
			ee, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatal(err)
			}
			exit = ee.ExitCode()
		}
		if area := os.Getenv("PLY_PARALLEL_EVIDENCE"); area != "" {
			if !filepath.IsAbs(area) {
				t.Fatal("evidence area must be absolute")
			}
			path := filepath.Join(area, fmt.Sprintf("candidate-cli-%02d", n))
			if e := os.MkdirAll(area, 0700); e != nil {
				t.Fatal(e)
			}
			for suffix, b := range map[string][]byte{".stdout": stdout.Bytes(), ".stderr": stderr.Bytes(), ".exit": []byte(fmt.Sprint(exit) + "\n"), ".argv.json": mustJSONCLI(t, append([]string{binary}, args...))} {
				if e := os.WriteFile(path+suffix, b, 0600); e != nil {
					t.Fatal(e)
				}
			}
		}
		if want == 0 && exit != 0 || want != 0 && exit == 0 {
			t.Fatalf("CLI %v exit=%d: %s %s", args, exit, stdout.String(), stderr.String())
		}
		return stdout.Bytes()
	}
	validate := func(name string, raw []byte) {
		t.Helper()
		c := exec.Command(python, "-I", "-B", "-c", `import importlib.metadata,json,jsonschema,sys
assert importlib.metadata.version('jsonschema')=='4.26.0'
s=json.load(open(sys.argv[1]));jsonschema.Draft202012Validator.check_schema(s)
jsonschema.Draft202012Validator(s,format_checker=jsonschema.FormatChecker()).validate(json.load(sys.stdin))`, filepath.Join(source, "schemas", "read-contract", name))
		c.Stdin = bytes.NewReader(raw)
		out, e := c.CombinedOutput()
		if e != nil {
			t.Fatalf("schema %s: %v %s", name, e, out)
		}
	}
	for _, entry := range []struct {
		args     []string
		fragment string
	}{{[]string{"workflow", "execute", "--help"}, "--parallel"}, {[]string{"workflow", "execute", "reconcile-source", "--help"}, "--continue"}, {[]string{"workspace", "registry", "activate-parallel", "--help"}, "--qualified-run"}} {
		if !strings.Contains(string(call(0, entry.args...)), entry.fragment) {
			t.Fatalf("help omits %s", entry.fragment)
		}
	}
	validate("workspace-task-queue-v2.schema.json", call(0, "workspace", "task", "queue", "list", "--schema-version", "2", "--format", "json"))
	if _, e = Execute(d, in); e != nil {
		t.Fatal(e)
	}
	in.SpecID = "observable-b"
	if _, e = Execute(d, in); e != nil {
		t.Fatal(e)
	}
	validate("workspace-task-queue-v2.schema.json", call(0, "workspace", "task", "queue", "list", "--schema-version", "2", "--format", "json"))
	validate("workspace-task-queue-v2.schema.json", call(0, "workspace", "task", "list", "--ready", "--schema-version", "2", "--format", "json"))
	validate("workflow-status-v2.schema.json", call(0, "workflow", "status", "--schema-version", "2", "--all", "--format", "json"))
	call(1, "workspace", "task", "queue", "list", "--format", "json")
	call(1, "workflow", "status", "--format", "json")
	before, e := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if e != nil {
		t.Fatal(e)
	}
	call(1, "workflow", "execute", "--next", "--parallel", "--check", "--format", "json")
	after, e := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if e != nil || !bytes.Equal(before, after) {
		t.Fatal("invalid CLI admission wrote registry")
	}
}

func mustJSONCLI(t *testing.T, v any) []byte {
	t.Helper()
	raw, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	return append(raw, '\n')
}

func TestParallelPrivateCandidatePreservesActualSerialFour(t *testing.T) {
	binary := os.Getenv("PLY_DELIVERY_TEST_BINARY")
	if binary == "" {
		t.Skip("requires exact private make-built candidate")
	}
	if !filepath.IsAbs(binary) {
		t.Fatal("candidate must be absolute")
	}
	d, in, root := parallelLauncherFixture(t, binary)
	r, e := d.Workspace.WorkItems.Snapshot(root)
	if e != nil {
		t.Fatal(e)
	}
	// Seed the same ordinary pre-activation format4 history; no parallel event,
	// Run or activation has been published in this isolated fixture.
	r.FormatVersion = 4
	r.ParallelBoundary = nil
	r.SourceReconciliations = nil
	r.ParallelActivation = nil
	raw, e := yaml.Marshal(r)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(root, ".ply", "work-items.yaml"), raw, 0600); e != nil {
		t.Fatal(e)
	}
	in.Parallel = false
	a, e := Execute(d, in)
	if e != nil || a.Run == nil || a.Run.Transport.AgentSessionID == "" {
		t.Fatalf("native serial A startup: %v", e)
	}
	before, e := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if e != nil {
		t.Fatal(e)
	}
	in.SpecID = "observable-b"
	if _, e = Execute(d, in); e == nil {
		t.Fatal("ordinary serial B admitted while A active")
	}
	after, e := os.ReadFile(filepath.Join(root, ".ply", "work-items.yaml"))
	if e != nil || !bytes.Equal(before, after) {
		t.Fatal("serial refusal changed shared registry")
	}
	q, e := workspace.ListTaskQueue(d.Workspace, in.Target, false)
	if e != nil || q.Registry.FormatVersion != 4 || q.Current == nil || q.Current.TaskID != "goal" || len(q.Pending) != 1 || q.Pending[0].TaskID != "goal-b" {
		t.Fatalf("serial format4 or pending changed: %+v %v", q, e)
	}
	c := exec.Command(binary, "workspace", "task", "create", "serial-fresh", "--title", "fresh", "--description", "ordinary registration stays format4", "--epic", "epic", "--project", "fixture", "--repo", "fixture", "--format", "json")
	c.Dir = filepath.Join(root, "epic")
	if raw, e = c.CombinedOutput(); e != nil {
		t.Fatalf("candidate ordinary registration: %v %s", e, raw)
	}
	r, e = d.Workspace.WorkItems.Snapshot(root)
	if e != nil || r.FormatVersion != 4 {
		t.Fatal("maximum reader capability changed ordinary registration format")
	}
}
