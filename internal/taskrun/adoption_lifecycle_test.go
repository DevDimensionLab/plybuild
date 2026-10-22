package taskrun

import (
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Native lifecycle with injected Herdr/provider/process observations. This proves
// software semantics, not actual OS writer absence or provider acceptance.
func TestAdoptionNativeGoalReviewVerifyIntegrateClose(t *testing.T) {
	f, o, binary := reviewRoundFixture(t, "adopted")
	r0, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	foreign := func(r workspace.WorkItemRegistry) string {
		var rows []any
		for _, v := range r.Tasks {
			if v.ID != "task" {
				rows = append(rows, v)
			}
		}
		for _, v := range r.TaskProblemRevisions {
			if v.TaskID != "task" {
				rows = append(rows, v)
			}
		}
		return digest(rows)
	}
	originalForeign := foreign(r0)
	q, e := workspace.ListTaskQueue(f.D.Workspace, workspace.QueueTargetInput{ProjectID: "ply", RepoID: "ply", EpicID: "epic"}, false)
	if e != nil {
		t.Fatal(e)
	}
	if len(q.Pending) != 2 || q.Pending[0].TaskID != "foreign-u" || q.Pending[1].TaskID != "foreign-v" {
		t.Fatal("execution did not retain foreign pending order")
	}
	plan := f.Prepared.Preparation.Preparation.Plan
	if plan.SourceMode != "adopted" || plan.StartOID == plan.ParentOID || plan.StartTree == plan.ParentTree || plan.WorktreePath != filepath.Join(filepath.Dir(f.Parent), "useful-existing") {
		t.Fatalf("lost S/B: %+v", plan)
	}
	s, e := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if s.Observed.Target.OID != plan.StartOID {
		t.Fatal("fresh handoff did not start at S")
	}
	candidate, _ := reviewRoundCandidate(t, f, "fixed candidate includes inherited delta")
	script, e := os.ReadFile(f.AcceptancePath)
	if e != nil {
		t.Fatal(e)
	}
	script = append(script, []byte("test \"$(cat inherited.txt)\" = 'useful inherited delta'\n")...)
	if e = os.WriteFile(f.AcceptancePath, script, 0600); e != nil {
		t.Fatal(e)
	}
	coordinator := reviewRoundAsCoordinator(f)
	reviewRoundIn(t, "w-fixture:root")
	requested, e := WorkflowDeliveryRequestReview(coordinator, f.R.WorkspaceRoot, o.RunID, DeliveryReviewRequestInput{Coordinator: "modeled adoption coordinator"})
	if e != nil {
		t.Fatal(e)
	}
	round := reviewRoundLatest(t, requested)
	brief, e := os.ReadFile(round.Brief.Locator)
	if e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(string(brief), "Review/diff base: "+plan.ParentOID) || !strings.Contains(string(brief), "Adopted start: "+plan.StartOID) {
		t.Fatalf("generated reviewer brief omits frozen B..S inherited scope: %s", brief)
	}
	passed := reviewRoundReturnFile(t, f, "adoption-pass.json", `{"decision":"passed","summary":"Modeled review assessed inherited and new delta","findings":[],"evidence_gaps":[]}`)
	reviewRoundIn(t, round.ReviewerPane)
	if _, e = WorkflowDeliveryReviewReturn(coordinator, f.R.WorkspaceRoot, o.RunID, round.Context.Locator, passed); e != nil {
		t.Fatal(e)
	}
	if o, e = WorkflowDeliveryVerifyWithBinary(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context, "", binary); e != nil {
		t.Fatal(e)
	}
	if o.Delivery.Candidates[0].OID != candidate {
		t.Fatal("qualified wrong C")
	}
	o, e = WorkflowDeliveryIntegrate(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	cleanupOwnerHandoff(t, f, o, e)
	if gitOutput(t, f.Parent, "rev-parse", "HEAD") != candidate {
		t.Fatal("did not integrate C")
	}
	reflog := gitOutput(t, f.Parent, "reflog", "show", "--format=%H", "refs/heads/epic")
	cleanupOwnerLeaves(t, f)
	reviewRoundIn(t, "w-fixture:root")
	coordinator.LocalProcesses = func(string, string) ([]LocalProcess, error) {
		return []LocalProcess{{PID: os.Getpid(), CWD: f.Parent}}, nil
	}
	o, e = WorkflowFollow(coordinator, f.R.WorkspaceRoot, o.RunID, 30*time.Second)
	if e != nil {
		t.Fatal(e)
	}
	assertCleanupRemoved(t, f, o, reflog)
	r, e := f.D.Workspace.WorkItems.Snapshot(f.R.WorkspaceRoot)
	if e != nil {
		t.Fatal(e)
	}
	if foreign(r) != originalForeign {
		t.Fatal("delivery rewrote foreign native history")
	}
	if len(r.WorktreeOperations) != 1 || r.WorktreeOperations[0].Adoption == nil {
		t.Fatal("lost immutable adoption history")
	}
	shown, e := workspace.ShowTask(f.D.Workspace, "task")
	if e != nil {
		t.Fatal(e)
	}
	if len(shown.Adoptions) != 1 {
		t.Fatal("adoption no longer readable after cleanup")
	}
}
