package taskactor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// A receipt that reserved a start, or stopped before the agent, names no
// started agent: the agent is that of the highest Goal revision whose receipt
// reached an agent, else that of the receipt kept per Task.
func TestStartedAgentIsTheHighestRevisionThatReachedAnAgent(t *testing.T) {
	receipt := func(key, provider, reached string) string {
		return `{"key":"` + key + `","provider":"` + provider + `","input_sha256":"a","reached":"` + reached + `"}`
	}
	revision := func(n int, provider, reached string) map[string]string {
		return map[string]string{taskwork.GoalStartReceiptPath(workspace.ProjectWork{}, "task", "spec", n): receipt(taskwork.GoalStartKey("task", "spec", n), provider, reached)}
	}
	single := func(provider, reached string) map[string]string {
		return map[string]string{taskwork.TaskGoalStartReceiptPath(workspace.ProjectWork{}, "task"): receipt("task", provider, reached)}
	}
	for name, c := range map[string]struct {
		receipts []map[string]string
		want     string
	}{
		"no receipt":                              {nil, ""},
		"only a reservation":                      {[]map[string]string{revision(1, "codex", "reserved")}, ""},
		"later revision reserved":                 {[]map[string]string{revision(1, "codex", "prompt_submitted"), revision(2, "claude", "reserved")}, "codex:task"},
		"later revision stopped at its tab":       {[]map[string]string{revision(1, "codex", "prompt_submitted"), revision(2, "claude", "tab_created")}, "codex:task"},
		"later revision with an empty receipt":    {[]map[string]string{revision(1, "codex", "prompt_submitted"), {taskwork.GoalStartReceiptPath(workspace.ProjectWork{}, "task", "spec", 2): ""}}, "codex:task"},
		"later revision started its agent":        {[]map[string]string{revision(1, "codex", "prompt_submitted"), revision(2, "claude", "agent_started")}, "claude:task"},
		"later revision bound a session":          {[]map[string]string{revision(1, "codex", "prompt_submitted"), revision(2, "claude", "session_bound")}, "claude:task"},
		"receipt per Task beside a later start":   {[]map[string]string{single("codex", "prompt_submitted"), revision(1, "claude", "prompt_submitted")}, "claude:task"},
		"receipt per Task beside a reservation":   {[]map[string]string{single("codex", "prompt_submitted"), revision(1, "claude", "reserved")}, "codex:task"},
		"receipt per Task that only reserved":     {[]map[string]string{single("codex", "reserved")}, ""},
		"receipt of another Task under this name": {[]map[string]string{{taskwork.GoalStartReceiptPath(workspace.ProjectWork{}, "task", "spec", 1): receipt("other@spec@1", "codex", "prompt_submitted")}}, ""},
	} {
		t.Run(name, func(t *testing.T) {
			work := workspace.ProjectWork{Dir: t.TempDir()}
			for _, files := range c.receipts {
				for relative, content := range files {
					path := filepath.Join(work.Dir, relative)
					if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(path, []byte(content), 0600); err != nil {
						t.Fatal(err)
					}
				}
			}
			if got := StartedAgent(work, "task"); got != c.want {
				t.Fatalf("started agent %q, want %q", got, c.want)
			}
		})
	}
}
