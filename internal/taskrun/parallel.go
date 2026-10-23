package taskrun

import (
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"strings"
)

// Called under the existing workflow store lock before a session is bound or
// a reviewer return accepted. It inspects native bindings, never OS processes.
func parallelSessionAvailable(d Dependencies, s workflowState, provider, session string, excludeOwner bool, ownReview ...string) error {
	if s.Request.SchemaVersion != 3 {
		return nil
	}
	r, err := d.Workspace.WorkItems.SnapshotRegistrations(s.Request.WorkspaceRoot)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(filepath.Join(workflowRoot(s.Request.WorkspaceRoot), "requests"))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			return workflowError(4, "unrecognized workflow reservation")
		}
		id := strings.TrimSuffix(entry.Name(), ".json")
		other, err := workflowRead(s.Request.WorkspaceRoot, id)
		if err != nil {
			return err
		}
		if !workspace.PreparationIsActive(r, other.Request.PreparationID) {
			continue
		}
		if !(excludeOwner && id == s.Result.RunID) && other.Request.Runtime.Provider == provider && other.Result.Transport.AgentSessionID == session {
			return workflowError(4, "native session already belongs to active Run "+id+"; no Task prompt or reviewer authority may be shared")
		}
		rounds, err := deliveryReviewRounds(other)
		if err != nil {
			return err
		}
		for _, round := range rounds {
			if len(ownReview) > 0 && id == s.Result.RunID && round.Dir == ownReview[0] {
				continue
			}
			start := deliveryReviewStartState(round.Dir)
			if start.Agent != nil && start.Agent.Provider == provider && start.Agent.SessionID == session {
				return workflowError(4, "native session already belongs to active member's reviewer in Run "+id)
			}
		}
	}
	return nil
}

func validateParallelRequest(r WorkflowRequest) error {
	id, err := workflowhandoff.ParallelPreparation(r.HandoffDraft)
	if err != nil {
		return err
	}
	if r.SchemaVersion != 3 {
		if id != "" || r.Delivery != nil && (r.Delivery.AdmissionMode != "" || r.Delivery.SourceReconciliation != "") {
			return workflowError(2, "legacy request cannot acquire parallel authority")
		}
		return nil
	}
	if r.Delivery == nil || r.Delivery.AdmissionMode != workspace.AdmissionParallel || r.Delivery.SourceReconciliation != workspace.SourceReconciliationMergeParent || id != r.PreparationID || r.Worker != nil || r.ControlContract != nil || r.Delivery.Agreement == nil || r.Delivery.Agreement.Mode != workspace.DeliveryLocalEpic {
		return workflowError(2, "parallel request requires exact explicit interactive local Epic authority")
	}
	return nil
}

func deliveryEffectiveBasis(d Dependencies, s workflowState) (*workspace.EffectiveDeliveryBasis, error) {
	if s.Request.SchemaVersion != 3 {
		return nil, nil
	}
	var out workspace.EffectiveDeliveryBasis
	read := func(session *workspace.TaskSpecSession) error {
		var err error
		out, err = session.EffectiveDeliveryBasis(s.Request.PreparationID)
		if err != nil {
			return err
		}
		draft, err := workflowhandoff.ValidateDeliveryTaskRunDraft(s.Request.HandoffDraft)
		if err != nil {
			return err
		}
		_, err = session.ValidateDeliveryCandidateTargetBasis(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref, s.Observed.Target.GitCommonDir, draft.Basis, &out)
		return err
	}
	var err error
	if d.deliverySnapshot != nil {
		if d.deliverySnapshot.Root != s.Request.WorkspaceRoot {
			return nil, workflowError(4, "release snapshot belongs to another workspace")
		}
		err = read(d.deliverySnapshot)
	} else {
		err = workspace.WithTaskSpecSnapshot(d.Workspace, s.Request.WorkspaceRoot, read)
	}
	if err != nil {
		return nil, err
	}
	return &out, nil
}
func deliveryBasisDigest(d Dependencies, s workflowState) (string, error) {
	b, err := deliveryEffectiveBasis(d, s)
	if err != nil {
		return "", err
	}
	if b == nil {
		return "", nil
	}
	return digest(b), nil
}
func deliveryRequireBasis(d Dependencies, s workflowState, want *workspace.EffectiveDeliveryBasis) error {
	have, err := deliveryEffectiveBasis(d, s)
	if err != nil {
		return err
	}
	if !equal(have, want) {
		return workflowError(4, "effective generation changed during verification; fresh review and verification required")
	}
	return nil
}

// Native source reconciliation uses the same callback/owner/session authority.
// It cannot upgrade a historical Run or transfer another member's effects.
func WorkflowDeliveryReconcileSource(d Dependencies, root, id, context string, in workspace.TaskSourceReconciliationInput) (workspace.TaskSourceReconciliationReadback, error) {
	var out workspace.TaskSourceReconciliationReadback
	if err := containing(d, root); err != nil {
		return out, err
	}
	err := withStore(root, func() error {
		s, err := workflowRead(root, id)
		if err != nil {
			return err
		}
		if s.Request.SchemaVersion != 3 {
			return workflowError(4, "reconcile-source requires the originally frozen parallel Run")
		}
		if err = deliveryCallback(d, s, context, true); err != nil {
			return err
		}
		if a := s.Result.Delivery.Attempt; a != nil && a.State != "recorded" {
			return workflowError(4, "resolve the previous native effect before reconciliation")
		}
		rounds, err := deliveryReviewRounds(s)
		if err != nil {
			return err
		}
		for _, r := range rounds {
			if r.Result == nil && !r.Superseded {
				return workflowError(4, fmt.Sprintf("seal or close review round %d before reconciliation", r.Request.Round))
			}
		}
		in.PreparationID = s.Request.PreparationID
		in.RunID = id
		in.RequestSHA256 = s.Result.RequestSHA256
		in.MandateSHA256 = s.Result.Handoff.SHA256
		in.ReviewSHA256 = digest(rounds)
		out, err = workspace.ReconcileTaskSource(d.Workspace, in)
		return err
	})
	return out, err
}
