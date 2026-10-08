package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/delivery"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// Decide accepts only the caller's actual answer to the displayed preview. It
// rechecks the entire plan under the effect lock before preserving that answer.
func (s *Service) Decide(cwd string, shown Preview, human workflowhandoff.DeliveryHumanAttestation) (Receipt, error) {
	var out Receipt
	d, root, e := s.context(cwd)
	if e != nil {
		return out, e
	}
	if shown.Plan.Workspace != root || digest(shown.Plan) != shown.PlanSHA256 || operationID(shown.PlanSHA256) != shown.ID {
		return out, fmt.Errorf("plan_changed: invalid displayed plan")
	}
	b, e := taskrun.Canonical(human)
	if e != nil {
		return out, e
	}
	if _, e = workflowhandoff.ValidateDeliveryHumanAttestation(b, string(shown.Plan.TaskResult.TaskID), shown.Plan.TaskResult, human.Answer); e != nil {
		return out, e
	}
	e = withLock(root, func() error {
		if old, err := readReceipt(root, shown.ID); err == nil {
			out = old
			return fmt.Errorf("operation_exists: resume integration %s; a decision is never overwritten", old.ID)
		} else if !os.IsNotExist(err) {
			return err
		}
		fresh, err := s.Preview(cwd, shown.Plan.Selection)
		if err != nil {
			return err
		}
		if fresh.PlanSHA256 != shown.PlanSHA256 {
			return fmt.Errorf("plan_changed: candidate, target, QA, ownership or effect choices changed; inspect a new plan")
		}
		if fresh.State != "ready" && human.Answer == "pass" {
			return fmt.Errorf("integration_blocked: %s", strings.Join(fresh.Reasons, "; "))
		}
		decision := Decision{Kind: "ply.integration.decision", SchemaVersion: 1, IntegrationID: shown.ID, PlanSHA256: shown.PlanSHA256, Human: human}
		path := filepath.Join(storeRoot(root), "decisions", shown.ID+".json")
		if err = atomicJSON(path, decision, true); err != nil {
			return err
		}
		out = Receipt{Kind: "ply.integration.receipt", SchemaVersion: 1, ID: shown.ID, Plan: shown.Plan, PlanSHA256: shown.PlanSHA256, Decision: decision, DecisionFile: taskrun.FileBinding{Locator: path, SHA256: digest(decision)}, State: "decision_recorded", Reasons: []string{}, NextAction: "Preserve native human judgment for this exact candidate and plan.", CreatedAtUTC: s.now(), NotificationState: "off", Events: []Event{}}
		if shown.Plan.NotificationRoute != nil {
			out.NotificationState = "not_attempted"
		}
		if err = s.save(&out); err != nil {
			return err
		}
		if human.Answer == "pass" {
			if err = reserve(out); err != nil {
				return s.stop(&out, "blocked", err.Error(), "Resume the integration that already reserved this candidate and target.")
			}
		}
		if err = s.accept(d, &out); err != nil {
			return s.stop(&out, "blocked", err.Error(), "Resolve the native human authority conflict, then resume this operation.")
		}
		if human.Answer != "pass" {
			return s.stop(&out, "human_"+human.Answer, "The actual human answered "+human.Answer+"; no integration, installation or cleanup was started.", "Address the human findings, then preview and confirm the resulting candidate again.")
		}
		return s.execute(d, cwd, &out, false)
	})
	return out, e
}

func (s *Service) accept(d taskrun.Dependencies, r *Receipt) error {
	if r.Plan.Flow == "legacy_reconcile" {
		return nil
	}
	a, e := taskrun.AcceptHumanIntegration(d, r.Plan.Workspace, r.Plan.WorkflowRunID, taskrun.HumanIntegrationInput{PlanID: r.ID, PlanSHA256: r.PlanSHA256, Decision: r.DecisionFile, Attestation: r.Decision.Human})
	if e != nil {
		return e
	}
	r.HumanAcceptance = &a
	return s.save(r)
}

func (s *Service) Resume(cwd, id string, check, retryInstall bool, retryMerge ...bool) (Receipt, error) {
	d, root, e := s.context(cwd)
	if e != nil {
		return Receipt{}, e
	}
	r, e := readReceipt(root, id)
	if e != nil {
		return r, e
	}
	if check {
		return r, nil
	}
	e = withLock(root, func() error {
		var err error
		r, err = readReceipt(root, id)
		if err != nil {
			return err
		}
		if r.State == "cancelled" {
			return fmt.Errorf("plan_revoked: the later negative human decision revoked this operation; preview a new plan")
		}
		if r.Reconsideration != nil {
			return s.finishReconsider(d, &r)
		}
		if r.Decision.Human.Answer != "pass" {
			return fmt.Errorf("human_qa_required: this operation has no exact pass")
		}
		if err = reserve(r); err != nil {
			return err
		}
		if r.HumanAcceptance == nil && r.Plan.Flow != "legacy_reconcile" {
			if err = s.accept(d, &r); err != nil {
				return s.stop(&r, "blocked", err.Error(), "Resolve the preserved native decision before any effect.")
			}
		}
		if r.State == "completed" {
			return s.notify(&r)
		}
		if len(retryMerge) > 0 && retryMerge[0] {
			if err = s.currentQA(d, &r); err != nil {
				return err
			}
			if err = s.prepareMergeRetry(&r); err != nil {
				return err
			}
		}
		return s.execute(d, cwd, &r, retryInstall)
	})
	return r, e
}

func (s *Service) stop(r *Receipt, state, reason, next string) error {
	r.State = state
	r.Reasons = []string{reason}
	r.NextAction = next
	return s.save(r)
}
func checkFile(b *taskrun.FileBinding) error {
	if b == nil {
		return nil
	}
	raw, e := readBounded(b.Locator, 16<<20)
	if e != nil {
		return e
	}
	if byteHash(raw) != b.SHA256 {
		return fmt.Errorf("plan_changed: preserved profile, route or authority bytes changed: %s", b.Locator)
	}
	return nil
}
func (s *Service) currentQA(d taskrun.Dependencies, r *Receipt) error {
	if r.Plan.Flow == "legacy_reconcile" {
		return nil
	}
	if r.HumanAcceptance == nil || !r.HumanAcceptance.Accepted {
		return fmt.Errorf("human_qa_required: native acceptance is missing")
	}
	registry, e := d.Workspace.WorkItems.Snapshot(r.Plan.Workspace)
	if e != nil {
		return e
	}
	qa, e := workspace.LatestTaskHumanQA(registry.HumanQARecords, r.Plan.TaskResult.ID, r.Plan.TaskResult.ResultOID, r.Plan.TaskResult.ResultTree)
	if e != nil {
		return e
	}
	if qa == nil || qa.Outcome != "pass" || qa.ID != r.HumanAcceptance.HumanQA.ID {
		return fmt.Errorf("human_qa_changed: the confirmed plan's pass is no longer the latest exact candidate answer")
	}
	return nil
}

func (s *Service) execute(d taskrun.Dependencies, cwd string, r *Receipt, retryInstall bool) error {
	if e := s.currentQA(d, r); e != nil {
		return s.stop(r, "blocked", e.Error(), "Review the latest candidate judgment and confirm a new plan if appropriate.")
	}
	for _, b := range []*taskrun.FileBinding{r.Plan.InstallFile, r.Plan.NotificationRoute, r.Plan.Ownership.Release} {
		if e := checkFile(b); e != nil {
			return s.stop(r, "blocked", e.Error(), "Restore the approved bytes; changes require a new visible decision.")
		}
	}
	if r.Plan.Install != nil && (r.Installation == nil || r.Installation.State != "verified") {
		ready, e := s.controller(d, r)
		if e != nil {
			return e
		}
		if !ready {
			return nil
		}
	}
	if !r.Integrated {
		if r.Plan.Flow == "legacy_reconcile" {
			r.Integrated = true
			r.IntegratedOID = r.Plan.TaskResult.ResultOID
			s.event(r, "integration.reconciled")
			if e := s.save(r); e != nil {
				return e
			}
		} else {
			_, _, e := delivery.NewService(d, delivery.Options{}).ValidateForIntegration(cwd, r.Plan.DeliveryID)
			if e != nil {
				return s.stop(r, "blocked", "candidate_changed: "+e.Error(), "Requalify the candidate; no unobserved integration is repeated.")
			}
			if r.Plan.Flow == "github_pr_merge" {
				if e = s.merge(r); e != nil {
					return e
				}
			} else {
				if e = s.local(d, r); e != nil {
					return e
				}
			}
			if !r.Integrated {
				return nil
			}
			if r.Local != nil && !r.Local.Completed {
				return s.notify(r)
			}
		}
	}
	if r.Local != nil && !r.Local.Completed {
		if e := s.local(d, r); e != nil {
			return e
		}
		if r.Local == nil || !r.Local.Completed {
			return s.notify(r)
		}
	}
	var pendingCloseout *workspace.TaskCloseoutReceipt
	if r.Closeout == nil {
		in, e := closeoutEvidence(r)
		if e != nil {
			return e
		}
		phase := "awaiting_closeout"
		if r.Plan.Install != nil && (r.Installation == nil || r.Installation.State != "verified") {
			phase = "installation_pending"
		}
		guard := func() error {
			_, e := taskrun.CheckLegacyTaskOwnership(d, r.Plan.Workspace, r.Plan.TaskResult)
			return e
		}
		reserved, e := workspace.ReserveTaskCloseout(d.Workspace, in, r.Plan.CloseoutSHA256, phase, guard)
		if e != nil {
			return s.stop(r, "cleanup_pending", e.Error(), "Resolve the current Task ownership or closeout blocker; resume this same Integration ID afterward.")
		}
		if reserved.State == "complete" {
			r.Closeout = &reserved
			if e = s.save(r); e != nil {
				return e
			}
		} else {
			pendingCloseout = &reserved
		}
	}
	if r.Plan.Install != nil {
		if e := s.install(r, retryInstall); e != nil {
			return e
		}
		if r.Closeout == nil && pendingCloseout != nil && strings.HasPrefix(pendingCloseout.Phase, "install") {
			phase := "awaiting_closeout"
			if r.Installation == nil || r.Installation.State != "verified" {
				phase = "install_unknown"
				if r.Installation != nil && r.Installation.State == "install_failed" {
					phase = "install_failed"
				}
			}
			if _, e := workspace.SetTaskCloseoutPhase(d.Workspace, r.Plan.TaskResult.TaskID, r.ID, r.Plan.CloseoutSHA256, phase); e != nil {
				return e
			}
		}
		if r.Installation == nil || r.Installation.State != "verified" {
			return s.notify(r)
		}
	}
	if r.Closeout == nil {
		in, e := closeoutEvidence(r)
		if e != nil {
			return e
		}
		r.State = "closing"
		r.NextAction = "Preserve evidence and close only the confirmed Task resources."
		if e := s.save(r); e != nil {
			return e
		}
		if e := s.fault("before_closeout"); e != nil {
			return e
		}
		guard := func() error {
			_, e := taskrun.CheckLegacyTaskOwnership(d, r.Plan.Workspace, r.Plan.TaskResult)
			return e
		}
		closed, e := workspace.ApplyTaskCloseoutWithOwnershipGuard(d.Workspace, in, r.Plan.CloseoutSHA256, guard)
		if e != nil {
			if err := s.stop(r, "cleanup_pending", e.Error(), "Resolve the exact cleanup blocker, then resume this Integration ID; the integration will not repeat."); err != nil {
				return err
			}
			return s.notify(r)
		}
		if e = s.fault("after_closeout"); e != nil {
			return e
		}
		if closed.State != "complete" {
			return s.stop(r, "cleanup_pending", "Native closeout remains incomplete.", "Resume this operation to finish the preserved closeout steps.")
		}
		r.Closeout = &closed
		if e = s.save(r); e != nil {
			return e
		}
	}
	s.event(r, "integration.closed")
	if e := s.fault("before_completed"); e != nil {
		return e
	}
	r.State = "completed"
	r.Reasons = []string{}
	r.NextAction = "The selected Task integration and closeout are complete."
	if e := s.save(r); e != nil {
		return e
	}
	return s.notify(r)
}

func closeoutEvidence(r *Receipt) (workspace.TaskCloseoutInput, error) {
	in := closeoutInput(r.Plan, r.ID)
	in.RequireIntegration = true
	if r.PR != nil {
		path := filepath.Join(storeRoot(r.Plan.Workspace), "evidence", r.ID+"-merge.json")
		if e := atomicJSON(path, *r.PR, true); e != nil {
			return in, e
		}
		in.ObservedPR = &workspace.TaskCloseoutPRIntegration{Repository: r.Plan.PR.Repository, Number: r.Plan.PR.Number, HeadOID: r.Plan.TaskResult.ResultOID, BaseRef: r.Plan.Agreement.TargetRef, MergeOID: r.PR.MergeOID, URL: r.PR.URL, MergedAtUTC: r.PR.MergedAt, EvidenceLocator: path, EvidenceSHA256: digest(*r.PR)}
	}
	return in, nil
}

func (s *Service) local(d taskrun.Dependencies, r *Receipt) error {
	if !r.IntegrationStarted {
		r.IntegrationStarted = true
		r.State = "integrating"
		if e := s.save(r); e != nil {
			return e
		}
		if e := s.fault("before_local_effect"); e != nil {
			return e
		}
	}
	result, e := taskrun.CompleteHumanLocalDelivery(d, r.Plan.Workspace, r.Plan.WorkflowRunID, r.Plan.TaskResult.ID, r.PlanSHA256)
	if err := s.fault("after_local_effect"); err != nil {
		return err
	}
	r.Local = &result
	if result.Integration.Readback.RecoveryStatus == "complete" && result.Integration.Readback.ParentOID == r.Plan.TaskResult.ResultOID {
		r.Integrated = true
		r.IntegratedOID = r.Plan.TaskResult.ResultOID
		s.event(r, "integration.observed")
	}
	if e != nil || !result.Completed {
		reason := "native base/queue completion remains pending"
		if e != nil {
			reason = e.Error()
		}
		if err := s.stop(r, "integration_pending", reason, "Resume this operation to observe native integration and complete only its remaining base/queue updates."); err != nil {
			return err
		}
		return s.notify(r)
	}
	return s.save(r)
}

func (s *Service) merge(r *Receipt) error {
	target := *r.Plan.PR
	if r.PR != nil {
		target.OperationID = r.PR.OperationID
	}
	o, e := s.options.PR.Observe(target)
	if o.OperationID == "" {
		o.OperationID = target.OperationID
	}
	if e != nil {
		o.State, o.NoEffect = "effect_unknown", false
		if o.Reason == "" {
			o.Reason = e.Error()
		}
		r.MergeObservations = append(r.MergeObservations, o)
		if r.PR == nil {
			r.PR = &o
		}
		return s.stop(r, "effect_unknown", e.Error(), "Restore exact PR observation and resume this operation; no merge request will be repeated.")
	}
	r.MergeObservations = append(r.MergeObservations, o)
	if o.State == "merged" {
		r.PR = &o
		r.Integrated = true
		r.IntegratedOID = o.MergeOID
		s.event(r, "integration.observed")
		return s.save(r)
	}
	if r.IntegrationStarted {
		if r.PR != nil && r.PR.NoEffect && exactRetryPR(target, o) && (o.State == "ready" || o.State == "blocked") {
			return s.stop(r, "merge_blocked", r.PR.Reason, "The previous request is proven not to have merged. Resolve its blocker, then explicitly resume with --retry-merge, or reconsider this plan.")
		}
		r.PR = &o
		if o.NoEffect {
			return s.stop(r, "merge_blocked", o.Reason, "The previous request is proven not to have merged. Resolve its blocker, then explicitly resume with --retry-merge, or reconsider this plan.")
		}
		state := "effect_unknown"
		if o.State == "remote_pending" {
			state = o.State
		}
		return s.stop(r, state, o.Reason, "Resume this same Integration ID to observe the already reserved merge; do not start another merge request.")
	}
	if o.State != "ready" {
		r.PR = &o
		return s.stop(r, o.State, o.Reason, o.NextAction)
	}
	r.IntegrationStarted = true
	r.State = "integrating"
	r.PR = &o
	if e = s.save(r); e != nil {
		return e
	}
	if e = s.fault("before_pr_effect"); e != nil {
		return e
	}
	o, e = s.options.PR.Merge(target)
	r.PR = &o
	// Keep an acknowledged operation ID before any subsequent phase boundary.
	if err := s.save(r); err != nil {
		return err
	}
	if err := s.fault("after_pr_effect"); err != nil {
		return err
	}
	if e != nil {
		return s.stop(r, "effect_unknown", e.Error(), "Observe this same operation before any further action; the merge request will not be repeated.")
	}
	if o.State == "merged" {
		r.Integrated = true
		r.IntegratedOID = o.MergeOID
		s.event(r, "integration.observed")
		return s.save(r)
	}
	if o.NoEffect {
		return s.stop(r, "merge_blocked", o.Reason, "The request is proven not to have merged. Resolve its blocker, then explicitly resume with --retry-merge, or reconsider this plan.")
	}
	return s.stop(r, o.State, o.Reason, o.NextAction)
}

func (s *Service) install(r *Receipt, retry bool) error {
	if r.Installation != nil && r.Installation.State == "verified" {
		return nil
	}
	if r.InstallationStarted && (!retry || r.Installation == nil || r.Installation.State != "install_failed") {
		previous := r.Installation
		o, e := s.options.Installer.Observe(*r.Plan.Install)
		r.InstallationObservations = append(r.InstallationObservations, o)
		r.Installation = &o
		if previous != nil && previous.State == "install_failed" && (e != nil || o.State != "verified") {
			r.Installation = previous
			reason := o.Reason
			if e != nil {
				reason = e.Error()
			}
			return s.stop(r, "install_failed", "The last installation command failed; current verification has not confirmed the artifact: "+reason, "Restore the approved installation inputs, then explicitly resume with --retry-install.")
		}
		if e != nil {
			return s.stop(r, "install_unknown", e.Error(), "Observe the installed artifact before retry; keep the worktree and do not repeat integration.")
		}
		if o.State == "verified" {
			return s.save(r)
		}
		return s.stop(r, "install_unknown", o.Reason, "Verify or restore the approved artifact, then resume. An unknown install is never automatically repeated.")
	}
	r.InstallationStarted = true
	// A previous failed command authorizes only this explicit retry. Once its
	// intent is durable, an interruption must not reuse that old failure as
	// evidence about the new command's unknown outcome.
	r.Installation = nil
	r.State = "installing"
	if e := s.save(r); e != nil {
		return e
	}
	if e := s.fault("before_install_effect"); e != nil {
		return e
	}
	o, e := s.options.Installer.Install(*r.Plan.Install)
	if err := s.fault("after_install_effect"); err != nil {
		return err
	}
	r.Installation = &o
	r.InstallationAttempts = append(r.InstallationAttempts, o)
	if e != nil || o.State != "verified" {
		reason := o.Reason
		if e != nil {
			reason = e.Error()
		}
		if o.State == "effect_unknown" {
			return s.stop(r, "install_unknown", reason, "Keep the source and observe the installed artifact; the installation outcome is unknown and cannot be retried blindly.")
		}
		return s.stop(r, "install_failed", reason, "Keep the source; inspect the observed installation. A known failed install may be resumed with --retry-install; integration is already preserved.")
	}
	return s.save(r)
}

func (s *Service) notify(r *Receipt) error {
	if !r.Integrated || r.Plan.NotificationRoute == nil {
		return nil
	}
	route, e := LoadNotificationRoute(r.Plan.NotificationRoute.Locator)
	if e != nil {
		return s.notificationError(r, e)
	}
	if route.SHA256 != r.Plan.NotificationRoute.SHA256 {
		return s.notificationError(r, fmt.Errorf("plan_changed: notification route changed"))
	}
	if r.NotificationEvent == nil {
		var effect Event
		for _, ev := range r.Events {
			if ev.Kind == "integration.observed" || ev.Kind == "integration.reconciled" {
				effect = ev
				break
			}
		}
		remaining := []string{}
		if r.State != "completed" {
			remaining = append(remaining, r.State)
		}
		r.NotificationEvent = &NotificationEvent{Kind: "ply.integration.observed", SchemaVersion: 1, ID: effect.ID, IntegrationID: r.ID, TaskID: string(r.Plan.TaskResult.TaskID), DeliveryID: r.Plan.DeliveryID, CandidateOID: r.Plan.TaskResult.ResultOID, TargetRef: r.Plan.Agreement.TargetRef, IntegratedOID: r.IntegratedOID, ObservedAtUTC: effect.AtUTC, PRURL: effect.PRURL, Remaining: remaining}
		if e = s.save(r); e != nil {
			return e
		}
	}
	var observed NotificationObservation
	if r.NotificationStarted {
		observed, e = s.options.Notifier.Observe(route, *r.NotificationEvent)
	} else {
		r.NotificationStarted = true
		r.NotificationState = "reserved"
		if e = s.save(r); e != nil {
			return e
		}
		if e = s.fault("before_notification"); e != nil {
			return e
		}
		observed, e = s.options.Notifier.Send(route, *r.NotificationEvent)
		if err := s.fault("after_notification"); err != nil {
			return err
		}
	}
	r.Notification = &observed
	if e != nil {
		return s.notificationError(r, e)
	}
	r.NotificationState = observed.State
	if observed.State != "transport_acknowledged" {
		r.NextAction = "Integration status is preserved. Inspect the separate notification outcome; do not retry an unknown send."
	}
	return s.save(r)
}
func (s *Service) notificationError(r *Receipt, e error) error {
	r.NotificationState = "unknown"
	r.NextAction = "Integration status is preserved; notification needs attention: " + e.Error()
	return s.save(r)
}

func (s *Service) controller(d taskrun.Dependencies, r *Receipt) (bool, error) {
	current, e := d.Executable()
	if e != nil {
		return false, e
	}
	current, e = filepath.EvalSymlinks(current)
	if e != nil {
		return false, e
	}
	if r.Controller == nil {
		raw, err := readBounded(current, 256<<20)
		if err != nil {
			return false, err
		}
		path := filepath.Join(storeRoot(r.Plan.Workspace), "controllers", r.ID, "ply-control")
		if err = ensureDirectories(filepath.Dir(path)); err != nil {
			return false, err
		}
		if old, err := readBounded(path, 256<<20); err == nil {
			if byteHash(old) != byteHash(raw) {
				return false, fmt.Errorf("controller changed")
			}
		} else if os.IsNotExist(err) {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0500)
			if err != nil {
				return false, err
			}
			_, err = f.Write(raw)
			if err == nil {
				err = f.Sync()
			}
			closeErr := f.Close()
			if err != nil {
				return false, err
			}
			if closeErr != nil {
				return false, closeErr
			}
		} else {
			return false, err
		}
		r.Controller = &taskrun.FileBinding{Locator: path, SHA256: byteHash(raw)}
		if err = s.save(r); err != nil {
			return false, err
		}
	}
	if e = ValidateInstallController(*r.Plan.Install, current, r.Plan.TaskResult.SourceLocator, r.DecisionFile.Locator, r.Plan.InstallFile.Locator); e != nil {
		return false, s.stop(r, "controller_required", e.Error(), r.Controller.Locator+" integration resume "+r.ID)
	}
	return true, nil
}
