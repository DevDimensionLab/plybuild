package taskrun

import (
	"errors"

	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Exact DeliveryState wire shape from the immutable controller baseline
// 680fecc. Its decoder rejects unknown fields, so a negative human answer must
// return ownership using this shape while retaining new facts in event files.
type legacyDeliveryStateBeforeHumanIntegration struct {
	Phase                string                `json:"phase"`
	OwnerClaim           string                `json:"owner_claim"`
	PermissionState      string                `json:"permission_state"`
	ActualPolicySHA256   *string               `json:"actual_policy_sha256"`
	ActualPolicyEvidence []Evidence            `json:"actual_policy_evidence"`
	Events               []DeliveryEventRecord `json:"events"`
	Candidates           []DeliveryCandidate   `json:"candidates"`
	Attempt              *DeliveryAttempt      `json:"attempt"`
	LastEventSHA256      *string               `json:"last_event_sha256"`
}

func humanIntegrationCandidate(t *testing.T, human bool) (deliveryFixture, WorkflowRun) {
	t.Helper()
	a := workspace.DeliveryAgreement{Mode: workspace.DeliveryLocalBranch}
	if human {
		a.SchemaVersion = 2
		a.IntegrationOwner = workspace.IntegrationOwnerHuman
	}
	f := newDeliveryFixture(t, "codex", deliveryFixtureOptions{Agreement: &a, ParentRef: "main"})
	f.R.Delivery.Goal = FileBinding{filepath.Join(f.R.WorkspaceRoot, ".ply", "task-content", "v1", "manifests", "sha256", strings.TrimPrefix(f.Prepared.Goal.Spec.ManifestSHA256, "sha256:")+".json"), f.Prepared.Goal.Spec.ManifestSHA256}
	f.File = writeAny(t, f.R.WorkspaceRoot, "delivery-request.json", f.R)
	o := deliveryTestAccept(t, f, workflowTestStart(t, f.workflowFixture))
	deliveryTestCandidate(t, f, "human integration candidate")
	o, e := WorkflowDeliveryVerify(f.D, f.R.WorkspaceRoot, o.RunID, o.Paths.Context)
	if e != nil {
		t.Fatal(e)
	}
	return f, o
}

type releaseDuringIntegrationGit struct {
	workspace.TaskIntegrationGit
	beforeMerge func()
}

func (g releaseDuringIntegrationGit) MergeFastForward(in workspace.IntegrationMergeInput) workspace.GitCommandOutcome {
	g.beforeMerge()
	return g.TaskIntegrationGit.MergeFastForward(in)
}

type closeoutRetirementFailureStore struct{ workspace.WorkItemStore }

func (s closeoutRetirementFailureStore) WithLock(root string, action func(workspace.WorkItemStoreSession) error) error {
	return s.WorkItemStore.WithLock(root, func(session workspace.WorkItemStoreSession) error {
		return action(closeoutRetirementFailureSession{session})
	})
}

type closeoutRetirementFailureSession struct{ workspace.WorkItemStoreSession }

func (s closeoutRetirementFailureSession) Publish(registry workspace.WorkItemRegistry) error {
	for _, task := range registry.Tasks {
		if task.WorktreeState == workspace.WorkItemRetired {
			return errors.New("synthetic retirement publication interruption")
		}
	}
	return s.WorkItemStoreSession.Publish(registry)
}
