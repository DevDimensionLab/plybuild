//go:build darwin

package taskrun

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Provider/owner bootstrap is modeled by the existing fixture. All report
// consumption below runs the actual supplied candidate CLI. The deterministic
// under-read write is covered separately, not injected into this subprocess.
func TestAL23CandidateCLIRejectsFinalBytesForMixedDigest(t *testing.T) {
	binary := automaticCLIBinary(t)
	f, o := automaticCLIFixture(t, binary, true, 30)
	root := f.R.WorkspaceRoot
	evidence := al05Write(t, filepath.Join(root, "evidence.log"), []byte("BBBBbbbb"), 0644)
	evidence.SHA256 = hash([]byte("AAAAbbbb")) // actual consumed bytes in the read regression
	log := al05Write(t, filepath.Join(root, "verifier.log"), []byte("ok\n"), 0644)
	file := al05WriteReport(t, filepath.Join(root, "report.json"), al05Report(o, "al23-progress", []FileBinding{evidence}, log), 0644)
	run := func() (WorkflowRun, string, error) {
		args := []string{"workflow", "execute", "report", o.RunID, "--context", o.Paths.Context, "--file", file, "--format", "json"}
		c := exec.Command(binary, args...)
		c.Dir = f.Prepared.Preparation.Preparation.Plan.WorktreePath
		var stdout, stderr bytes.Buffer
		c.Stdout, c.Stderr = &stdout, &stderr
		err := c.Run()
		exit := -1
		if c.ProcessState != nil {
			exit = c.ProcessState.ExitCode()
		}
		t.Logf("argv=%q cwd=%q exit=%d\nstdout:\n%s\nstderr:\n%s", append([]string{binary}, args...), c.Dir, exit, stdout.String(), stderr.String())
		var got WorkflowRun
		if err == nil {
			if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
				t.Fatalf("stdout not one JSON value: %v", err)
			}
		}
		return got, stderr.String(), err
	}
	if _, stderr, err := run(); err == nil || !strings.Contains(stderr, "bytes changed since they were bound") || !strings.Contains(stderr, inputReportEvidence) || !strings.Contains(stderr, evidence.Locator) {
		t.Fatalf("wrong evidence consumed: %v %s", err, stderr)
	}
	s, err := workflowRead(root, o.RunID)
	if err != nil || len(s.Result.Delivery.Events) != 0 {
		t.Fatalf("rejected report changed history: %v", err)
	}
	// Bind the completed writer's bytes anew: the same owner continues without
	// a new permission, consistency assertion, or product decision.
	evidence.SHA256 = hash([]byte("BBBBbbbb"))
	al05WriteReport(t, file, al05Report(o, "al23-progress", []FileBinding{evidence}, log), 0644)
	for i := 0; i < 2; i++ {
		got, stderr, err := run()
		if err != nil || len(got.Delivery.Events) != 1 || len(got.Delivery.Candidates) != 0 || !strings.Contains(stderr, "level=warning") {
			t.Fatalf("stable report/retry: %v %s", err, stderr)
		}
	}
}
