package taskrun

import (
	"testing"
)

// This exercises generated native argv and instructions with the injected Herdr
// fixture. It does not start a provider or inspect live OS/browser state.

func labelArg(t *testing.T, argv []string) string {
	t.Helper()
	for i, a := range argv {
		if a == "--label" && i+1 < len(argv) {
			return argv[i+1]
		}
	}
	t.Fatalf("no label in %v", argv)
	return ""
}

func TestTaskPresentationStabilityAndCompatibility(t *testing.T) {
	g := deliveryReceiptGoal{TaskID: "task-stable", Title: "An ordinary title"}
	code, _ := g.presentation()
	renamed := g
	renamed.Title = "Et nytt navn 日本語"
	same, _ := renamed.presentation()
	other := g
	other.TaskID = "task-other"
	different, _ := other.presentation()
	if same != code || different == code {
		t.Fatalf("Task-derived code changed with title or collided: %s %s %s", code, same, different)
	}
	// Historical/report-only startup still has its original run/review/finished
	// labels.
	s := workflowState{}
	s.Request.Herdr.TabLabel = "Historical fixture"
	if got := labelArg(t, workflowTabCreateArgv(s, "w-fixture")); got != "run Historical fixture" {
		t.Fatal(got)
	}
	if got := workflowTaskTabLabel(s, "review 2"); got != "review 2 Historical fixture" {
		t.Fatal(got)
	}
	if got := workflowTaskTabLabel(s, "finished"); got != "finished Historical fixture" {
		t.Fatal(got)
	}
}
