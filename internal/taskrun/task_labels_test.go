package taskrun

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// This exercises generated native argv and instructions with the injected Herdr
// fixture. It does not start a provider or inspect live OS/browser state.
func TestNativeTaskCodeLabelsAndContext(t *testing.T) {
	f := receiptFixture(t, "claude")
	o := workflowTestStart(t, f.workflowFixture)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	goal, err := workflowBound(s.Request.Delivery.Goal, 4<<20)
	if err != nil {
		t.Fatal(err)
	}
	// The fixture's real Task ID is available before the Run/request is made.
	// Read the actual frozen Task identity rather than the Run ID.
	var g struct {
		TaskID string `json:"task_id"`
	}
	if err := json.Unmarshal(goal, &g); err != nil {
		t.Fatal(err)
	}
	taskID := g.TaskID
	code := fmt.Sprintf("T%X", sha256.Sum256([]byte(taskID)))[:9]
	calls := reviewRoundArgv(t, f, "tab create")
	if len(calls) != 1 {
		t.Fatalf("owner tab creates: %v", calls)
	}
	label := labelArg(t, calls[0])
	if !strings.HasPrefix(label, code+" owner") {
		t.Errorf("owner label = %q; want %s owner first", label, code)
	}
	if text := deliveryInstructions(s); !strings.Contains(text, "Use task code "+code+" in references and returns.") {
		t.Errorf("owner context lacks shared task code: %s", text)
	}
	// Regenerated owner/recovery argv and finished status retain the Task code.
	if again := labelArg(t, workflowTabCreateArgv(s, "w-fixture")); again != label {
		t.Errorf("resumed label changed: %q", again)
	}
	if finished := workflowTaskTabLabel(s, "finished"); !strings.HasPrefix(finished, code+" finished") {
		t.Errorf("finished label: %s", finished)
	}
	for _, round := range []int{1, 2} {
		plan, err := deliveryReviewPlan(f.D, s, round, "w-fixture")
		if err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("%s review %d", code, round)
		if label := labelArg(t, plan.TabArgv); !strings.HasPrefix(label, want) {
			t.Errorf("review label = %q; want %s first", label, want)
		}
		brief := deliveryReviewBrief(s, deliveryReviewRound{Round: round}, nil, f.R.WorkspaceRoot, "/fixture/ply", "/fixture/template", "/fixture/context", workspace.WorkspaceTaskPreparePlan{})
		if !strings.Contains(brief, "Use task code "+code+" in references and returns.") {
			t.Errorf("review brief lacks same code: %s", brief)
		}
	}
	// Explicit suffixes and Unicode use the same real generators, with an
	// in-memory synthetic goal binding; preserved fixture contracts stay intact.
	for _, code := range []string{"AL25B", "RO04R", "AL47"} {
		title := code + ": " + strings.Repeat("Blåbær 日本語 ", 15)
		path := writeAny(t, f.R.WorkspaceRoot, code+"-goal.json", map[string]string{"task_id": taskID, "title": title, "objective": "Full frozen objective"})
		copy := s
		copy.Request.Delivery = new(DeliveryContract)
		*copy.Request.Delivery = *s.Request.Delivery
		copy.Request.Delivery.Goal = FileBinding{path, hashFileTest(t, path)}
		label := labelArg(t, workflowTabCreateArgv(copy, "w-fixture"))
		if !strings.HasPrefix(label, code+" owner") || len([]rune(label)) > 40 || !utf8.ValidString(label) {
			t.Errorf("compact Unicode label: %q", label)
		}
		review, err := deliveryReviewPlan(f.D, copy, 2, "w-fixture")
		if err != nil {
			t.Fatal(err)
		}
		if got := labelArg(t, review.TabArgv); !strings.HasPrefix(got, code+" review 2") {
			t.Errorf("suffix review label: %q", got)
		}
		copy.Request.Herdr.TabLabel = "Renamed transport display"
		if got := workflowTaskTabLabel(copy, "finished"); !strings.HasPrefix(got, code+" finished") {
			t.Errorf("renamed finished label: %q", got)
		}
		prompt := deliveryInstructions(copy)
		if !strings.Contains(prompt, "Use task code "+code+" in references and returns.") || !strings.Contains(prompt, title) || !strings.Contains(prompt, "Full frozen objective") || !strings.Contains(prompt, o.RunID) {
			t.Errorf("full native context lost: %s", prompt)
		}
	}
}

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
	// labels. A headless contract never acquires presentation requirements.
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
	s.Request.Delivery = &DeliveryContract{StartupProtocol: MandateReceiptV1}
	s.Request.Worker = &WorkerContract{}
	if _, ok := deliveryTaskPresentation(s); ok {
		t.Fatal("headless gained tab presentation")
	}
}
