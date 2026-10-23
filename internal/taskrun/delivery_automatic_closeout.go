package taskrun

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// ObserveAutomaticLocalDelivery reads an already recorded local integration.
// It grants no merge authority, including after the source owner has released.
func ObserveAutomaticLocalDelivery(d Dependencies, root, id string, resultID workspace.TaskResultID) (workflowhandoff.DeliveryIntegrationResult, bool, error) {
	s, err := workflowRead(root, id)
	if err != nil {
		return workflowhandoff.DeliveryIntegrationResult{}, false, err
	}
	out, handled, err := automaticIntegratedResult(s, resultID)
	if err == nil && handled {
		out.Closeout, err = workflowTaskCloseout(root, s)
	}
	return out, handled, err
}

func automaticIntegratedResult(s workflowState, resultID workspace.TaskResultID) (workflowhandoff.DeliveryIntegrationResult, bool, error) {
	var out workflowhandoff.DeliveryIntegrationResult
	if !deliveryRun(s.Request) || deliveryEffectiveAgreement(s) == nil || !deliveryEffectiveAgreement(s).AutomaticAcceptance() || deliveryEffectiveAgreement(s).HumanOwnedIntegration() || s.Result.Delivery == nil || len(s.Result.Delivery.Candidates) == 0 {
		return out, false, nil
	}
	if s.Result.Delivery.HumanIntegration != nil {
		return out, false, nil
	}
	c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	if c.TaskResult.ID != resultID || c.Integration == nil {
		return out, false, nil
	}
	if s.Result.Delivery.Phase != "closing" && s.Result.Delivery.Phase != "completed" {
		return out, true, workflowError(4, "automatic closeout requires preserved local integration completion")
	}
	raw, err := workflowBound(*c.Integration, 4<<20)
	if err != nil {
		return out, true, err
	}
	var observation struct {
		Kind          string                                    `json:"kind"`
		RunID         string                                    `json:"run_id"`
		RequestSHA256 string                                    `json:"request_sha256"`
		Candidate     string                                    `json:"candidate"`
		Result        workflowhandoff.DeliveryIntegrationResult `json:"result"`
		Error         string                                    `json:"error,omitempty"`
	}
	if err = json.Unmarshal(raw, &observation); err != nil {
		return out, true, err
	}
	result := observation.Result.Integration.Readback.TaskResult
	if observation.Kind != "PlyDeliveryIntegration@1" || observation.RunID != s.Result.RunID || observation.RequestSHA256 != s.Result.RequestSHA256 || observation.Candidate != c.Key || observation.Error != "" || !observation.Result.Completed || observation.Result.Integration.Readback.RecoveryStatus != "complete" || observation.Result.Integration.Readback.ParentOID != c.OID && !observation.Result.Integration.Readback.HistoricalCompletion || result == nil || !equal(*result, c.TaskResult) {
		return out, true, workflowError(4, "automatic closeout integration receipt differs from its exact candidate")
	}
	found := false
	for _, event := range s.Result.Delivery.Events {
		if event.Kind == "integration" && equal(event.Binding, *c.Integration) {
			found = true
		}
	}
	if !found {
		return out, true, workflowError(4, "automatic closeout integration receipt is absent from native history")
	}
	return observation.Result, true, nil
}

// automaticReleaseRetainReason is the release text of a Run whose frozen
// mandate forbids cleanup. Such Runs keep their worktree and branch.
const automaticReleaseRetainReason = "Automatic local integration, queue closure and Epic base update are observed. Release source writing ownership for native closeout while retaining its worktree and branch."

const automaticReleaseCleanupReason = "Automatic local integration, queue closure and Epic base update are observed. Release source writing ownership for native closeout; the source worktree is removed after its owner has left it."

// automaticCloseoutKeeps binds the cleanup default to the Run's own frozen
// mandate. Mandates without the cleanup effect, including every Run started
// before it existed, keep their retained source; a newer CLI never widens them.
func automaticCloseoutKeeps(s workflowState) (bool, error) {
	allowed, err := workflowhandoff.DeliveryCleanupAllowed(s.Request.HandoffDraft)
	return !allowed, err
}

// automaticOwnerLeft requires positive absence of every known Task writer
// before removal: the bound Herdr pane and agent for an interactive owner, and
// for every owner the local processes in the source or naming its session.
// ownerPresent reports that the bound interactive owner itself remains.
func automaticOwnerLeft(d Dependencies, s workflowState) (ownerPresent bool, err error) {
	observed := (taskWriters{d, s}).cleanup()
	return observed.herdrPresent, observed.err
}

// automaticCleanupPending selects an interactive automatic delivery whose
// integration is observed and released but whose source removal remains.
func automaticCleanupPending(s workflowState) bool {
	a := deliveryEffectiveAgreement(s)
	if !deliveryRun(s.Request) || headlessRun(s.Request) || a == nil || !a.AutomaticAcceptance() || a.HumanOwnedIntegration() || s.Result.Delivery == nil || s.Result.Delivery.Phase != "closing" || s.Result.Delivery.OwnershipRelease == nil {
		return false
	}
	keep, err := automaticCloseoutKeeps(s)
	return err == nil && !keep
}

func automaticCloseoutPendingAction(s workflowState, cause string) WorkflowAction {
	if headlessRun(s.Request) {
		return WorkflowAction{"recipient", "Local integration, queue closure and Epic base update are observed. Source cleanup is pending: " + cause + ". The worker finishes this same closeout after its integration attempt has ended; the Git integration is not repeated."}
	}
	return WorkflowAction{"coordinator", "Local integration, queue closure and Epic base update are observed. Source cleanup is pending: " + cause + ". The owner ends its session; the coordinator's ply workflow execute follow " + s.Result.RunID + " then finishes this same closeout from outside the source without repeating the Git integration."}
}

// ResumeAutomaticLocalDelivery only completes native closeout after an
// observed integration. Its release cannot authorize a second Git effect.
func ResumeAutomaticLocalDelivery(d Dependencies, root, id string, resultID workspace.TaskResultID) (workflowhandoff.DeliveryIntegrationResult, bool, error) {
	s, err := workflowRead(root, id)
	if err != nil {
		return workflowhandoff.DeliveryIntegrationResult{}, false, err
	}
	out, handled, err := automaticIntegratedResult(s, resultID)
	if err != nil || !handled {
		return out, handled, err
	}
	c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
	authority, err := deliveryAuthorityState(d, s)
	if err != nil {
		return out, true, err
	}
	authority.CandidateRunID = c.TaskResult.RunID
	keep, err := automaticCloseoutKeeps(s)
	if err != nil {
		return out, true, err
	}
	reason := automaticReleaseCleanupReason
	if keep {
		reason = automaticReleaseRetainReason
	}
	err = deliveryReleaseUpdate(d, root, id, func(current *workflowState, d Dependencies) error {
		if _, found, e := automaticIntegratedResult(*current, resultID); e != nil || !found {
			if e == nil {
				e = workflowError(4, "automatic closeout integration changed before owner release")
			}
			return e
		}
		if prior, e := deliveryRelease(*current, c); e != nil {
			return e
		} else if prior != nil {
			raw, e := workflowBound(*prior, 1<<20)
			var recorded deliveryOwnershipRelease
			if e == nil {
				e = decode(raw, 1<<20, &recorded)
			}
			if e != nil {
				return e
			}
			if recorded.Origin != "automatic_delivery_completed" {
				return workflowError(4, "human_integration_required: explicit source release belongs to the human integration flow")
			}
			return nil
		}
		path, e := d.Executable()
		if e != nil {
			return e
		}
		path, e = filepath.EvalSymlinks(path)
		if e != nil {
			return e
		}
		bytes, e := readFile(path, 128<<20, false)
		if e != nil {
			return e
		}
		if e = preserveDeliveryRelease(d, current, c, reason, "automatic_delivery_completed", &Executable{Path: path, SHA256: hash(bytes)}, nil, nil); e != nil {
			return e
		}
		current.Result.NextAction = WorkflowAction{"recipient", "Finish native Task closeout for the observed automatic delivery; the source is removed after its owner has left it. Resume the same operation after interruption."}
		if keep {
			current.Result.NextAction = WorkflowAction{"recipient", "Finish native Task closeout for the observed automatic delivery; retain the source worktree and branch. Resume the same operation after interruption."}
		}
		return nil
	})
	if err != nil {
		return out, true, err
	}
	if err = d.fault("automatic_closeout_after_release"); err != nil {
		return out, true, err
	}
	s, err = workflowRead(root, id)
	if err != nil {
		return out, true, err
	}
	release, err := deliveryRelease(s, c)
	if err != nil || release == nil {
		if err == nil {
			err = workflowError(4, "automatic closeout ownership release is missing")
		}
		return out, true, err
	}
	// The release closes owner callbacks, not the owner's processes. Removal
	// waits until the owner and every local Task writer have left.
	ownerPresent := false
	ownerLeft := func() error {
		if keep {
			return nil
		}
		present, e := automaticOwnerLeft(d, s)
		ownerPresent = present
		return e
	}
	ownership := workspace.TaskCloseoutOwnership{Released: true, TaskID: c.TaskResult.TaskID, TaskResultID: resultID, ResultOID: c.OID, ResultTree: c.Tree, EvidenceLocator: release.Locator, EvidenceSHA256: release.SHA256}
	guard := func() error {
		observed, err := CheckLegacyTaskOwnership(d, root, c.TaskResult)
		if err != nil {
			return err
		}
		if !observed.Ready {
			return workflowError(4, "automatic closeout owner remains unresolved: "+observed.Reason)
		}
		return ownerLeft()
	}
	var closeout workspace.TaskCloseoutReceipt
	closeoutErr := ownerLeft()
	// The interactive owner itself is the remaining writer: its delivery is
	// done and it hands the source to the coordinator; this is not a stop.
	handoff := closeoutErr != nil && ownerPresent
	if closeoutErr == nil {
		closeout, closeoutErr = workflowhandoff.CloseAutomaticDeliveryCandidate(d.Workflow, string(c.TaskResult.TaskID), resultID, authority, ownership, keep, guard)
	}
	if closeout.Kind != "" {
		out.Closeout = &closeout
	}
	if closeoutErr == nil {
		closeoutErr = d.fault("automatic_closeout_after_effect")
	}
	err = workflowUpdate(d, root, id, func(current *workflowState) error {
		if _, found, e := automaticIntegratedResult(*current, resultID); e != nil || !found {
			if e == nil {
				e = workflowError(4, "automatic closeout candidate changed before completion")
			}
			return e
		}
		resource := "retired"
		if keep {
			resource = "kept"
		}
		if closeoutErr != nil || closeout.State != "complete" || !closeout.LifecycleCompleted || closeout.ResourceState != resource {
			current.Result.Delivery.Phase, current.Result.Round.State, current.Result.FinalReturn.State = "closing", "closing", "closing"
			if keep {
				current.Result.NextAction = WorkflowAction{"recipient", "Local integration, queue closure and Epic base update are observed. Resume this same automatic delivery to finish native Task closeout; the source worktree and branch remain retained."}
			} else {
				cause := "native Task closeout is incomplete"
				if closeoutErr != nil {
					cause = closeoutErr.Error()
				}
				current.Result.NextAction = automaticCloseoutPendingAction(*current, cause)
			}
			return nil
		}
		_, e := deliveryAppendEvent(d, current, "automatic-closeout-"+c.Key, "closeout", map[string]any{"kind": "PlyDeliveryCloseout@1", "run_id": id, "request_sha256": s.Result.RequestSHA256, "candidate": c.Key, "result": closeout})
		if e != nil {
			return e
		}
		current.Result.Delivery.Phase, current.Result.Round.State, current.Result.FinalReturn.State = "completed", "completed", "completed"
		current.Result.NextAction = WorkflowAction{"user", fmt.Sprintf("The exact automatically accepted candidate is integrated in %s. Task, queue and Epic base are complete; its worktree and branch are retained.", authority.Agreement.TargetRef)}
		if !keep {
			current.Result.NextAction = WorkflowAction{"user", fmt.Sprintf("The exact automatically accepted candidate is integrated in %s. Task, queue and Epic base are complete; its released worktree and branch are removed and the candidate is archived at %s.", authority.Agreement.TargetRef, closeout.Plan.ArchiveRef)}
		}
		return nil
	})
	if err == nil && handoff {
		err = workflowError(5, "owner_handoff: delivered; local integration, queue and Epic base are observed and this owner's writing is finished. End this session and return to the coordinator, whose workflow execute follow removes the source; this is not a stop and integrate needs no retry ("+closeoutErr.Error()+")")
	} else if err == nil && closeoutErr != nil && !keep {
		err = workflowError(5, "automatic local integration is observed; source cleanup is pending: "+closeoutErr.Error())
	} else if err == nil {
		err = closeoutErr
	}
	if err == nil && (out.Closeout == nil || out.Closeout.State != "complete") {
		err = workflowError(5, "automatic local integration is observed; native Task closeout remains pending")
	}
	return out, true, err
}
