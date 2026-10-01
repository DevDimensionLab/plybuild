package taskrun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// The compatibility script supplies a newly built CLI. This variable is used
// only by tests; it neither changes product validation nor launches a provider.
func TestR13FreshBinaryCallbacks(t *testing.T) {
	bin := os.Getenv("PLY_TASK_RUN_TEST_BINARY")
	if bin == "" {
		t.Skip("fresh CLI is exercised by make compat-cli")
	}
	d, r, f := fixture(t)
	r.Runtime.PlyExecutable = Executable{bin, hashFileTest(t, bin)}
	f = writeAny(t, r.WorkspaceRoot, "request.json", r)
	d.ContextPath = func() string { return filepath.Join(runPaths(r).TempRoot, "context.json") }
	call := func(cwd string, args ...string) ([]byte, error) {
		c := exec.Command(bin, args...)
		c.Dir = cwd
		c.Env = append(os.Environ(), "PLY_TASK_RUN_CONTEXT="+d.ContextPath())
		var out, errout bytes.Buffer
		c.Stdout = &out
		c.Stderr = &errout
		err := c.Run()
		if err != nil {
			return out.Bytes(), fmt.Errorf("%v: %v: %s", args, err, errout.String())
		}
		if errout.Len() != 0 {
			return nil, fmt.Errorf("stderr: %s", errout.String())
		}
		var doc map[string]any
		if e := json.Unmarshal(out.Bytes(), &doc); e != nil {
			return nil, fmt.Errorf("not one JSON document: %s", out.String())
		}
		if bytes.Contains(out.Bytes(), []byte(`"secret"`)) {
			return nil, fmt.Errorf("secret field leaked")
		}
		return out.Bytes(), nil
	}
	before := treeState(t, r.WorkspaceRoot)
	out, e := call(r.WorkspaceRoot, "workspace", "task", "run", "start", "--file", f, "--format", "json")
	if e != nil {
		t.Fatal(e)
	}
	if before != treeState(t, r.WorkspaceRoot) {
		t.Fatal("binary preview mutated workspace")
	}
	var preview Preview
	if e = json.Unmarshal(out, &preview); e != nil {
		t.Fatal(e)
	}
	runner := d.Runner.(*fakeRunner)
	runner.inside = func(s LaunchSpec) error {
		cp := writeAny(t, r.WorkspaceRoot, "claim.json", acceptance(t, d, r))
		if _, e = call(s.CWD, "workspace", "task", "run", "accept", RunID(r.RequestKey), "--file", cp, "--format", "json"); e != nil {
			return e
		}
		rp := writeAny(t, r.WorkspaceRoot, "return.json", semanticReport(t, d, r))
		if _, e = call(s.CWD, "workspace", "task", "run", "report", RunID(r.RequestKey), "--file", rp, "--format", "json"); e != nil {
			return e
		}
		for _, cmd := range []string{"show", "collect"} {
			if _, e = call(s.CWD, "workspace", "task", "run", cmd, RunID(r.RequestKey), "--format", "json"); e != nil {
				return e
			}
		}
		return nil
	}
	result, e := Start(d, f, *preview.Confirmation)
	if e != nil {
		t.Fatalf("%v; callback=%v", e, runner.lastErr)
	}
	if result.Collection.State != "qualified" {
		t.Fatalf("not qualified: %+v", result)
	}
	before = treeState(t, r.WorkspaceRoot)
	for _, cmd := range []string{"show", "collect"} {
		if _, e = call(r.WorkspaceRoot, "workspace", "task", "run", cmd, RunID(r.RequestKey), "--format", "json"); e != nil {
			t.Fatal(e)
		}
	}
	if before != treeState(t, r.WorkspaceRoot) {
		t.Fatal("binary readback wrote")
	}
	for _, args := range [][]string{{"start"}, {"start", "--file", f, "--apply", "--confirm", "x", "--format", "json"}, {"start", "--next"}, {"show", RunID(r.RequestKey), "extra"}, {"start", "--file", f, "--force"}} {
		c := exec.Command(bin, append([]string{"workspace", "task", "run"}, args...)...)
		c.Dir = r.WorkspaceRoot
		b, e := c.CombinedOutput()
		exit, ok := e.(*exec.ExitError)
		if !ok || exit.ExitCode() != 2 {
			t.Fatalf("usage status for %v: %v %s", args, e, b)
		}
	}
	t.Logf("Synthetic CLI journey; no provider or human QA. Binary SHA256 %s", r.Runtime.PlyExecutable.SHA256)
}
