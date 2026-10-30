package taskrun

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/buildinfo"
)

// Uses the real public CLI built from this checkout, including linked worktrees.
func TestAL16SelfHostedBuildIdentity(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	git := func(args ...string) string {
		t.Helper()
		c := exec.Command("git", args...)
		c.Dir = repo
		b, err := c.Output()
		if err != nil {
			t.Fatal(err)
		}
		return strings.TrimSpace(string(b))
	}
	oid := git("rev-parse", "HEAD")
	dirty := git("status", "--porcelain") != ""
	binary := filepath.Join(t.TempDir(), "ply")
	c := exec.Command("make", "build", "PLY_BUILD_OUTPUT="+binary)
	c.Dir = repo
	if b, err := c.CombinedOutput(); err != nil {
		t.Fatalf("build: %v %s", err, b)
	}
	c = exec.Command(binary, "version", "--format", "json")
	c.Dir = t.TempDir()
	b, err := c.Output()
	if err != nil {
		t.Fatal(err)
	}
	var identity buildinfo.Identity
	if err = json.Unmarshal(b, &identity); err != nil {
		t.Fatal(err)
	}
	if identity.SourceCommit != oid || identity.Dirty != dirty {
		t.Fatalf("build-time identity: got %s dirty=%t; want %s dirty=%t", identity.SourceCommit, identity.Dirty, oid, dirty)
	}
	t.Logf("actual CLI source=%s dirty=%t (executed outside source)", identity.SourceCommit, identity.Dirty)
}

// Actual VCS-free CLI; native provider/review facts in the reused fixture are
// modeled test observations, not real product acceptance or human judgment.
