package workspace

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ReserveTaskCloseout preserves the source reservation after observed integration
// and before required installation. ownershipGuard must only read current native
// ownership: it runs while the WorkItems lock excludes a new Task startup.
func ReserveTaskCloseout(d Dependencies, in TaskCloseoutInput, expectedPlanSHA256, phase string, ownershipGuard func() error) (out TaskCloseoutReceipt, err error) {
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	if d.WorkItems == nil || d.WorkClock == nil || ownershipGuard == nil {
		return out, fmt.Errorf("closeout reservation requires stores and a current ownership guard")
	}
	if !validTaskText(in.OperationID, 1, 256) || !closeoutPreCleanupPhase(phase) {
		return out, WorkInvalidArguments("closeout reservation requires an operation ID and explicit pending phase")
	}
	in.RequireIntegration = true
	err = d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		if err = ownershipGuard(); err != nil {
			return err
		}
		preview, err := previewTaskCloseout(d, root, r, in)
		if err != nil {
			return err
		}
		if preview.PlanSHA256 != expectedPlanSHA256 {
			return closeoutError("plan_changed", "Confirm the exact closeout resources before reserving them.", nil)
		}
		if !preview.Ready {
			return closeoutError("cleanup_blocked", strings.Join(preview.Reasons, "; ")+". "+preview.NextAction, nil)
		}
		if preview.Receipt != nil {
			out = *preview.Receipt
			return nil
		}
		native, _, err := closeoutIntegrationEvidence(d, r, preview.Plan, in.ObservedPR)
		if err != nil {
			return err
		}
		s, err := readTaskCloseouts(root)
		if err != nil {
			return err
		}
		now := d.WorkClock.Now().UTC().Format(time.RFC3339Nano)
		out = TaskCloseoutReceipt{Kind: "WorkspaceTaskCloseoutReceipt@1", SchemaVersion: 1, OperationID: in.OperationID, Plan: preview.Plan, PlanSHA256: expectedPlanSHA256, State: "cleanup_pending", ResourceState: "retained", Phase: phase, NativeIntegrationID: native, ObservedPR: in.ObservedPR, CreatedAtUTC: now, UpdatedAtUTC: now}
		s.Receipts = append(s.Receipts, out)
		return s.publish(root)
	})
	return out, err
}

func closeoutPreCleanupPhase(phase string) bool {
	switch phase {
	case "installation_pending", "install_failed", "install_unknown", "awaiting_closeout":
		return true
	}
	return false
}

// SetTaskCloseoutPhase records the installer's actual observation. It cannot
// change the approved resources, claim completion, or rewind a cleanup effect.
func SetTaskCloseoutPhase(d Dependencies, id TaskID, operationID, expectedPlanSHA256, phase string) (out TaskCloseoutReceipt, err error) {
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	if d.WorkItems == nil || d.WorkClock == nil || !closeoutPreCleanupPhase(phase) {
		return out, WorkInvalidArguments("closeout phase must describe pending or observed installation")
	}
	err = d.WorkItems.WithLock(root, func(_ WorkItemStoreSession) error {
		s, err := readTaskCloseouts(root)
		if err != nil {
			return err
		}
		r := s.receipt(id)
		if r == nil || r.OperationID != operationID || r.PlanSHA256 != expectedPlanSHA256 {
			return closeoutError("plan_changed", "Resume the exact reserved closeout operation.", nil)
		}
		out = *r
		if r.State != "cleanup_pending" || !closeoutPreCleanupPhase(r.Phase) || r.WorktreeRemoved || r.BranchRemoved || r.Retention != nil {
			return closeoutError("plan_changed", "Cleanup has already advanced; its observed phase cannot be replaced.", nil)
		}
		if r.Phase == phase {
			return nil
		}
		r.Phase, r.UpdatedAtUTC = phase, d.WorkClock.Now().UTC().Format(time.RFC3339Nano)
		out = *r
		return s.publish(root)
	})
	return out, err
}

// TaskCloseoutPendingAction exposes a stable reason and concrete next action
// without reporting installation as cleanup or as an unobserved success.
func TaskCloseoutPendingAction(r TaskCloseoutReceipt) (reason, kind, detail string) {
	switch r.Phase {
	case "installation_pending":
		return "installation_pending", "resume_installation", "Integration is observed. Resume the preserved operation to finish its required installation before Task cleanup."
	case "install_failed":
		return "install_failed", "resume_installation", "Integration is observed, but required installation failed. Inspect the preserved installation result and resume the same operation; the source is retained."
	case "install_unknown":
		return "install_unknown", "inspect_installation", "Integration is observed, but installation outcome is unknown. Observe the installed result before retrying the same operation; the source is retained."
	default:
		return "cleanup_pending", "resume_closeout", "Resume the preserved integration operation to complete Task closeout."
	}
}

// TaskCloseoutSourceRemovalPending recognizes the interrupted interval after a
// reserved worktree removal but before its observation was persisted. It claims
// only pending recovery, never completion or a fabricated WorktreeRemoved fact.
func TaskCloseoutSourceRemovalPending(r TaskCloseoutReceipt) bool {
	if r.State != "cleanup_pending" || r.Phase != "worktree_remove_reserved" || r.WorktreeRemoved || r.Retention == nil {
		return false
	}
	if _, err := os.Lstat(r.Plan.Source.Locator); !os.IsNotExist(err) {
		return false
	}
	return verifyCloseoutRetention(Dependencies{}, r.Plan.WorkspaceRoot, r.Plan, *r.Retention) == nil
}

func validateCloseoutTargetAvailable(root, locator, ref, common string) error {
	rows, err := ReadTaskCloseoutsAt(root)
	if err != nil {
		return contentError("task_content_observation_unknown", "closeout reservations could not be observed", err)
	}
	for _, r := range rows {
		p := r.Plan.Source
		if p.Locator == locator || p.GitCommonDir == common && p.Ref == ref {
			return contentError("task_closeout_reserved", "The exact Task source belongs to integration closeout "+r.OperationID+"; resume that operation instead of starting another writer.", nil)
		}
	}
	return nil
}
