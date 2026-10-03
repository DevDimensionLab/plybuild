package taskrun

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

func TestWorkflowNativeTerminalDriftInvalidatesAcceptedShow(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestReady(t, f)
	o, e := WorkflowReviewRun(f.D, f.R.WorkspaceRoot, o.RunID, workflowTestReview(t, f, o, "accepted"))
	if e != nil {
		t.Fatal(e)
	}
	facts, e := workflowhandoff.ReadTaskRunReturnFacts(f.D.Workflow, f.Handoff)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.Chmod(facts.Terminal.Locator, 0600); e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(facts.Terminal.Locator, []byte("changed native terminal"), 0600); e != nil {
		t.Fatal(e)
	}
	o, e = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if o.FinalReturn.State == "accepted" || len(o.Reasons) == 0 {
		t.Fatal("accepted show hid native terminal artifact drift")
	}
}
func TestWorkflowFailedObservationPersistsUnknown(t *testing.T) {
	f := workflowTestFixture(t)
	o := workflowTestReady(t, f)
	workflowTestModel(t, f, map[string]any{"agent_session_id": "replacement-session"})
	if _, e := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, time.Second); e == nil {
		t.Fatal("changed session observed as successful")
	}
	o, e := WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if o.Transport.State != "unknown" || o.Round.State == "ready_for_review" {
		t.Fatal("cached transport concealed the failed fresh observation")
	}
}
func TestWorkflowReadRejectsAncestorSwapAfterPhysicalPrecheck(t *testing.T) {
	root, e := filepath.EvalSymlinks(t.TempDir())
	if e != nil {
		t.Fatal(e)
	}
	parent, outside := filepath.Join(root, "parent"), filepath.Join(root, "outside")
	os.Mkdir(parent, 0700)
	os.Mkdir(outside, 0700)
	path := filepath.Join(parent, "bound.json")
	os.WriteFile(path, []byte("original"), 0600)
	os.WriteFile(filepath.Join(outside, "bound.json"), []byte("outside"), 0600)
	if e = physical(path, false); e != nil {
		t.Fatal(e)
	}
	if e = os.Rename(parent, parent+"-saved"); e != nil {
		t.Fatal(e)
	}
	if e = os.Symlink(outside, parent); e != nil {
		t.Fatal(e)
	}
	f, e := openRead(path)
	if e == nil {
		f.Close()
		t.Fatal("file open followed a swapped ancestor after precheck")
	}
}
