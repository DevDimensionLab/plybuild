package workspaceview

import (
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestAttentionSortsSeverityThenActualClockThenIDWithUnknownTimeLast(t *testing.T) {
	s := progressFixture()
	process := map[workspace.TaskID]ProcessFacts{}
	s.Registry.Tasks = nil
	for _, id := range []workspace.TaskID{"unknown", "later", "earlier", "conflict"} {
		s.Registry.Tasks = append(s.Registry.Tasks, workspace.TaskRecord{ID: id, ParentEpicID: "epic", ProjectID: "project", RepoID: "repo", WorktreeState: workspace.WorkItemUnbound})
		var since *string
		if id == "later" {
			since = stringValue("2026-10-06T08:00:00.5Z")
		}
		if id == "earlier" {
			since = stringValue("2026-10-06T08:00:00Z")
		}
		if id != "conflict" {
			process[id] = ProcessFacts{HasProgress: true, Freshness: "fresh", Waiting: []ProcessWaiting{{Kind: "recorded_wait", Actor: "human", Reason: "Choose", SinceUTC: since}}}
		}
	}
	first, second := fixtureResult("passed"), fixtureResult("passed")
	first.ID, second.ID = "first", "second"
	first.TaskID, second.TaskID = "conflict", "conflict"
	s.Registry.TaskResults = []workspace.TaskResultRecord{first, second}
	out, err := BuildAttention(s, workspace.TaskListFilters{}, process)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"conflict", "earlier", "later", "unknown"}
	if len(out.Items) != len(want) {
		t.Fatalf("items=%+v", out.Items)
	}
	for i, id := range want {
		if out.Items[i].SubjectID != id {
			t.Fatalf("sort=%+v", out.Items)
		}
	}
}

func TestUnknownProcessHistoryDoesNotInventProgressOrHideRecordedFacts(t *testing.T) {
	s := progressFixture()
	process := map[workspace.TaskID]ProcessFacts{"task": {Freshness: "unknown", Reasons: []string{"journal_unavailable"}}}
	p, err := BuildTaskProgress(s, "task", process)
	if err != nil || p.HasProgress || p.State != "nothing" || p.Freshness != "unknown" || len(p.Reasons) != 1 {
		t.Fatalf("unknown became invented work: %+v %v", p, err)
	}
	s.Registry.TaskResults = []workspace.TaskResultRecord{fixtureResult("passed")}
	p, err = BuildTaskProgress(s, "task", process)
	if err != nil || !p.HasProgress || p.TechnicalGate == nil || *p.TechnicalGate != "passed" || p.Freshness != "unknown" {
		t.Fatalf("optional source hid native facts: %+v %v", p, err)
	}
}
