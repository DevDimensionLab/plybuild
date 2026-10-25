package taskrun

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The acceptance script supplies the exact packaged candidate CLI. The fixture
// journal is produced by the existing fake runner; the CLI never launches a
// provider here because the Run already exists.
func TestF06CandidateCLIReadsJournalWithDamagedCache(t *testing.T) {
	bin := os.Getenv("PLY_TASK_RUN_TEST_BINARY")
	if bin == "" {
		t.Skip("candidate CLI is exercised by the F06 acceptance script")
	}
	for _, mode := range []string{"missing", "malformed", "forged_authority"} {
		t.Run(mode, func(t *testing.T) {
			d, r, f := fixture(t)
			completeInside(t, &d, r, nil)
			p, e := PreviewStart(d, f)
			if e != nil {
				t.Fatal(e)
			}
			_, _ = Start(d, f, *p.(Preview).Confirmation)
			id := RunID(r.RequestKey)
			call := func(args ...string) []byte {
				t.Helper()
				c := exec.Command(bin, append([]string{"workspace", "task", "run"}, args...)...)
				c.Dir = r.WorkspaceRoot
				for _, v := range os.Environ() {
					if !strings.HasPrefix(v, "PLY_TASK_RUN_CONTEXT=") && !strings.HasPrefix(v, "PLY_QA_") {
						c.Env = append(c.Env, v)
					}
				}
				var out, errout bytes.Buffer
				c.Stdout, c.Stderr = &out, &errout
				if err := c.Run(); err != nil {
					t.Fatalf("%v: %v: %s", args, err, errout.String())
				}
				return out.Bytes()
			}
			reads := [][]string{{"show", id, "--format", "json"}, {"start", "--file", f, "--format", "json"}, {"collect", id, "--check", "--format", "json"}}
			var baseline [][]byte
			for _, args := range reads {
				baseline = append(baseline, call(args...))
			}
			j, e := readJournal(r.WorkspaceRoot, id)
			if e != nil {
				t.Fatal(e)
			}
			if !bytes.Contains(baseline[1], []byte(fmt.Sprintf("%q", id))) {
				t.Fatalf("start did not rediscover the Run: %s", baseline[1])
			}
			calls := d.Runner.(*fakeRunner).calls
			damageResultCache(t, r, mode, j.Result)
			state := treeState(t, r.WorkspaceRoot)
			for i, args := range reads {
				if got := call(args...); !bytes.Equal(got, baseline[i]) {
					t.Fatalf("%v changed with %s cache:\n%s\n%s", args, mode, got, baseline[i])
				}
			}
			if state != treeState(t, r.WorkspaceRoot) {
				t.Fatal("candidate readback wrote or repaired files")
			}
			events, e := os.ReadDir(filepath.Join(runPaths(r).RunRoot, "events"))
			if e != nil || len(events) != len(j.Events) || d.Runner.(*fakeRunner).calls != calls {
				t.Fatal("candidate rediscovery relaunched or appended events")
			}
			t.Logf("candidate CLI %s read and rediscovered %s with %s cache", bin, id, mode)
		})
	}
}
