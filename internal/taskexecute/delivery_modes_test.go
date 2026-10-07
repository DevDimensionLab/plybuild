package taskexecute

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestExecuteFreezesEachDeliveryAgreementInPreviewIntentAndRequest(t *testing.T) {
	for _, mode := range []string{workspace.DeliveryPullRequest, workspace.DeliveryLocalEpic, workspace.DeliveryLocalBranch} {
		t.Run(mode, func(t *testing.T) {
			a := workspace.DeliveryAgreement{Mode: mode}
			if mode == workspace.DeliveryPullRequest {
				a.TargetRef, a.GitHubRepository, a.Remote = "refs/heads/main", "fixture/repository", "origin"
			}
			d, in, root, _ := launcherFixtureForProvider(t, "claude", a)
			in.Check = true
			preview, err := Execute(d, in)
			if err != nil {
				t.Fatal(err)
			}
			if preview.Goal.Delivery == nil || preview.Goal.Delivery.Mode != mode || preview.Goal.Delivery.SourceRef == "" {
				t.Fatal("execute preview omitted selected mode or bound source")
			}
			if _, err = os.Stat(preview.Goal.WorktreePath); !os.IsNotExist(err) {
				t.Fatal("delivery preview created a Task worktree")
			}
			in.Check = false
			d.Fault = func(point string) error {
				if point == "workflow_after_reservation" {
					return errors.New("fixture stop before provider startup")
				}
				return nil
			}
			started, err := Execute(d, in)
			if err == nil || !strings.Contains(err.Error(), "fixture stop") {
				t.Fatalf("execution did not reach reserved frozen request: %v", err)
			}
			var request taskrun.WorkflowRequest
			before, err := readPreservedJSON(started.RequestPath, &request)
			if err != nil {
				t.Fatal(err)
			}
			want, _ := taskrun.Canonical(preview.Goal.Delivery)
			got, _ := taskrun.Canonical(request.Delivery.Agreement)
			if !bytes.Equal(want, got) {
				t.Fatal("request changed selected mode/source/target")
			}
			if err = workflowhandoff.ValidateDeliveryMandate(request.HandoffDraft, request.Delivery.Agreement); err != nil {
				t.Fatal(err)
			}
			var intent launchIntent
			if _, err = readPreservedJSON(filepath.Join(filepath.Dir(started.RequestPath), "intent.json"), &intent); err != nil {
				t.Fatal(err)
			}
			intentAgreement, _ := taskrun.Canonical(intent.Goal.Delivery)
			if !bytes.Equal(want, intentAgreement) {
				t.Fatal("intent failed to freeze agreement")
			}
			// Replacing an editable source draft after launch cannot mutate the
			// preserved goal, execution Spec or the exact already reserved request.
			fixtureJSON(t, filepath.Join(root, "goal.json"), map[string]any{"delivery": map[string]any{"mode": "pull_request", "target_ref": "refs/heads/changed"}})
			d.Fault = nil
			in.Check = true
			existing, err := Execute(d, in)
			if err != nil || existing.State != "existing" {
				t.Fatalf("frozen retry changed contract: %+v %v", existing, err)
			}
			after, err := os.ReadFile(started.RequestPath)
			if err != nil || !bytes.Equal(before, after) {
				t.Fatal("started execution request was rewritten")
			}
		})
	}
}
