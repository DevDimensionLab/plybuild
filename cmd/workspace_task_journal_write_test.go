package cmd

import (
	"strings"
	"testing"
)

// A waiting step is recorded with flags: its start names the reason, the
// dependency and the next actor, and the help shows how (backlog c961beef).
func TestJournalStepStartRecordsAWait(t *testing.T) {
	d, _, _, _, _ := goalStartFixture(t)
	start := []string{"task", "journal", "step", "start", "moved", "--id", "w1", "--kind", "waiting", "--title", "Wait for answer", "--actor", "dev-1"}
	if _, err := runWorkspaceView(t, d, start...); err == nil || !strings.Contains(err.Error(), "--reason, --dependency and --next-actor") {
		t.Fatalf("waiting without its fields: %v", err)
	}
	other := []string{"task", "journal", "step", "start", "moved", "--id", "i1", "--kind", "implementation", "--title", "Implement", "--actor", "dev-1", "--reason", "x"}
	if _, err := runWorkspaceView(t, d, other...); err == nil || !strings.Contains(err.Error(), "only for --kind waiting") {
		t.Fatalf("waiting fields on another kind: %v", err)
	}
	if _, err := runWorkspaceView(t, d, append(start, "--reason", "Storage choice", "--dependency", "user answer", "--next-actor", "planner-a")...); err != nil {
		t.Fatal(err)
	}
	show, err := runWorkspaceView(t, d, "task", "journal", "show", "moved", "--view", "timeline")
	if err != nil || !strings.Contains(show, "waiting: Storage choice; dependency: user answer; next actor: planner-a") {
		t.Fatalf("show %v:\n%s", err, show)
	}
	help, _ := runWorkspaceView(t, d, "task", "journal", "step", "start", "--help")
	for _, want := range []string{"coordination, waiting", "--next-actor"} {
		if !strings.Contains(help, want) {
			t.Fatalf("help lacks %q:\n%s", want, help)
		}
	}
}
