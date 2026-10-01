package taskrun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		if len(args) == 2 && args[0] == "__generated_shell" {
			c = exec.Command("/bin/sh", "-c", args[1])
		}
		c.Dir = cwd
		c.Env = []string{}
		for _, v := range os.Environ() {
			if !strings.HasPrefix(v, "PLY_TASK_RUN_CONTEXT=") && !strings.HasPrefix(v, "PLY_QA_") {
				c.Env = append(c.Env, v)
			}
		}
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
		a := acceptance(t, d, r)
		a.SchemaVersion = 2
		a.RuntimeClaim.ModelID = nil
		cp := writeAny(t, r.WorkspaceRoot, "claim.json", a)
		j, err := readJournal(r.WorkspaceRoot, RunID(r.RequestKey))
		if err != nil {
			return err
		}
		instructions := recipientInstructions(r, *j.Binding)
		if strings.Count(instructions, "--context "+ShellQuote(s.ContextPath)) != 2 {
			return fmt.Errorf("generated callbacks lack exact context")
		}
		// Execute the actual generated callback command in a clean tool shell.
		generatedAccept := strings.Split(strings.Split(instructions, "Run: ")[1], "\n")[0]
		generatedAccept = strings.Replace(generatedAccept, "<absolute-private-claim.json>", ShellQuote(cp), 1) + " --format json"
		if _, e = call(s.CWD, "__generated_shell", generatedAccept); e != nil {
			return e
		}
		if e = os.MkdirAll(filepath.Dir(j.Binding.ReportPath), 0700); e != nil {
			return e
		}
		writeAny(t, filepath.Dir(j.Binding.ReportPath), filepath.Base(j.Binding.ReportPath), semanticReport(t, d, r))
		generatedReport := strings.Split(strings.Split(instructions, " and run: ")[1], "\n")[0] + " --format json"
		if _, e = call(s.CWD, "__generated_shell", generatedReport); e != nil {
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
	var pending *Error
	if !errors.As(e, &pending) || pending.Code != "task_run_task_status_unknown" || runner.lastErr != nil {
		t.Fatalf("%v; callback=%v", e, runner.lastErr)
	}
	if result.RuntimeFacts.ModelState != "unknown" || result.RuntimeFacts.ReportedModel != nil || result.RuntimeFacts.Source == nil {
		t.Fatal("unknown model provenance lost")
	}
	status := writeAny(t, r.WorkspaceRoot, "status.json", syntheticStatus(t, d, r, "inactive"))
	out, e = call(r.WorkspaceRoot, "workspace", "task", "run", "collect", RunID(r.RequestKey), "--task-status", status, "--check", "--format", "json")
	if e != nil {
		t.Fatal(e)
	}
	var collect CollectPreview
	if e = json.Unmarshal(out, &collect); e != nil {
		t.Fatal(e)
	}
	out, e = call(r.WorkspaceRoot, append(collect.NextArgv[1:], "--format", "json")...)
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(out, &result); e != nil {
		t.Fatal(e)
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
	for _, args := range [][]string{{"show", RunID(r.RequestKey), "--context", "/x"}, {"start", "--file", f, "--context", "/x"}, {"collect", RunID(r.RequestKey), "--context", "/x"}, {"accept", RunID(r.RequestKey), "--file", f, "--task-status", "/x"}, {"report", RunID(r.RequestKey), "--file", f, "--task-status", "/x"}, {"show", RunID(r.RequestKey), "--task-status", "/x"}, {"start"}, {"start", "--file", f, "--apply", "--confirm", "x", "--format", "json"}, {"start", "--next"}, {"show", RunID(r.RequestKey), "extra"}, {"start", "--file", f, "--force"}} {
		c := exec.Command(bin, append([]string{"workspace", "task", "run"}, args...)...)
		c.Dir = r.WorkspaceRoot
		b, e := c.CombinedOutput()
		exit, ok := e.(*exec.ExitError)
		if !ok || exit.ExitCode() != 2 {
			t.Fatalf("usage status for %v: %v %s", args, e, b)
		}
	}
	textCmd := exec.Command(bin, "workspace", "task", "run", "show", RunID(r.RequestKey))
	textCmd.Dir = r.WorkspaceRoot
	textOut, err := textCmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Requested model:", "Reported model: unknown (recipient claim)", "Client process: exited", "Task execution: inactive (source: human_observation)"} {
		if !bytes.Contains(textOut, []byte(want)) {
			t.Fatalf("missing %q in %s", want, textOut)
		}
	}
	t.Logf("Synthetic CLI journey; no provider or human QA. Binary SHA256 %s", r.Runtime.PlyExecutable.SHA256)
}
