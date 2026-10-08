package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

type ReconsiderationEvidence struct {
	Reason  string              `json:"reason"`
	Started bool                `json:"started"`
	Prior   *PRMergeObservation `json:"prior_request,omitempty"`
	Current *PRMergeObservation `json:"current_pr,omitempty"`
}

type ReconsiderationReceipt struct {
	Kind          string                              `json:"kind"`
	SchemaVersion int                                 `json:"schema_version"`
	IntegrationID string                              `json:"integration_id"`
	PlanSHA256    string                              `json:"plan_sha256"`
	Decision      Decision                            `json:"decision"`
	DecisionFile  taskrun.FileBinding                 `json:"decision_file"`
	Evidence      ReconsiderationEvidence             `json:"no_effect_evidence"`
	Native        *taskrun.HumanIntegrationAcceptance `json:"native_negative_decision"`
	Proof         *taskrun.FileBinding                `json:"proof,omitempty"`
	RecordedAtUTC string                              `json:"recorded_at_utc"`
}

type ReconsiderationPreview struct {
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	ID            string                   `json:"id"`
	PlanSHA256    string                   `json:"plan_sha256"`
	BasisSHA256   string                   `json:"basis_sha256"`
	CandidateOID  string                   `json:"candidate_oid"`
	Allowed       bool                     `json:"allowed"`
	Evidence      *ReconsiderationEvidence `json:"no_effect_evidence,omitempty"`
	Reason        string                   `json:"reason,omitempty"`
	NextAction    string                   `json:"next_action"`
}

// reconsiderEvidence establishes that the selected operation has not produced
// its integration effect. An interrupted intent is never equivalent to absence.
func (s *Service) reconsiderEvidence(r Receipt) (ReconsiderationEvidence, error) {
	var e ReconsiderationEvidence
	if r.Integrated || r.InstallationStarted || r.Closeout != nil || r.Plan.Flow == "legacy_reconcile" || r.State == "completed" || r.Reconsideration != nil || r.Decision.Human.Answer != "pass" {
		return e, fmt.Errorf("reconsider_forbidden: only an unintegrated, uncancelled passed plan can be reconsidered")
	}
	if !r.IntegrationStarted {
		e.Reason = "never_dispatched"
		if r.Plan.Flow == "github_pr_merge" {
			if r.Plan.PR == nil {
				return e, fmt.Errorf("reconsider_forbidden: the exact PR binding is missing")
			}
			current, err := s.options.PR.Observe(*r.Plan.PR)
			if err != nil {
				return e, err
			}
			if !exactRetryPR(*r.Plan.PR, current) || (current.State != "ready" && current.State != "blocked") {
				return e, fmt.Errorf("effect_unknown: reconsideration requires a fresh exact unmerged PR observation")
			}
			e.Current = &current
		}
		return e, nil
	}
	if r.Plan.Flow != "github_pr_merge" || r.Plan.PR == nil || r.PR == nil {
		return e, fmt.Errorf("effect_unknown: a started local or unknown integration cannot be cancelled")
	}
	target := *r.Plan.PR
	target.OperationID = r.PR.OperationID
	current, err := s.options.PR.Observe(target)
	if err != nil {
		return e, err
	}
	if !exactRetryPR(target, current) || (current.State != "ready" && current.State != "blocked") {
		return e, fmt.Errorf("effect_unknown: cancellation requires a fresh exact unmerged PR observation")
	}
	prior := *r.PR
	if !prior.NoEffect {
		if !current.NoEffect {
			return e, fmt.Errorf("effect_unknown: an open PR is not proof that the outstanding request had no effect")
		}
		prior = current
	}
	if !exactRetryPR(target, prior) || (prior.State != "ready" && prior.State != "blocked") {
		return e, fmt.Errorf("reconsider_forbidden: previous rejection differs from the approved PR")
	}
	e = ReconsiderationEvidence{Reason: "request_rejected", Started: true, Prior: &prior, Current: &current}
	return e, nil
}

func (s *Service) CheckReconsideration(cwd, id string) (ReconsiderationPreview, error) {
	r, err := s.Read(cwd, id)
	if err != nil {
		return ReconsiderationPreview{}, err
	}
	p := ReconsiderationPreview{Kind: "ply.integration.reconsideration-preview", SchemaVersion: 1, ID: r.ID, PlanSHA256: r.PlanSHA256, BasisSHA256: digest(r), CandidateOID: r.Plan.TaskResult.ResultOID, NextAction: "Answer fail or blocked with a reason to revoke only this unexecuted plan; no integration or cleanup will occur."}
	evidence, err := s.reconsiderEvidence(r)
	if err != nil {
		p.Reason = err.Error()
		p.NextAction = "Observe and finish the existing operation; unknown or integrated effects cannot be cancelled."
		return p, nil
	}
	p.Allowed = true
	p.Evidence = &evidence
	return p, nil
}

func (s *Service) Reconsider(cwd, id, basis string, human workflowhandoff.DeliveryHumanAttestation) (Receipt, error) {
	d, root, err := s.context(cwd)
	if err != nil {
		return Receipt{}, err
	}
	var r Receipt
	err = withLock(root, func() error {
		var err error
		r, err = readReceipt(root, id)
		if err != nil {
			return err
		}
		if r.Reconsideration != nil {
			return fmt.Errorf("reconsideration_exists: resume the operation to preserve the already recorded answer")
		}
		if digest(r) != basis {
			return fmt.Errorf("plan_changed: operation changed after the reconsideration preview")
		}
		if human.Answer != "fail" && human.Answer != "blocked" {
			return fmt.Errorf("negative_answer_required: reconsideration requires the actual fail or blocked answer")
		}
		raw, err := taskrun.Canonical(human)
		if err != nil {
			return err
		}
		if _, err = workflowhandoff.ValidateDeliveryHumanAttestation(raw, string(r.Plan.TaskResult.TaskID), r.Plan.TaskResult, human.Answer); err != nil {
			return err
		}
		evidence, err := s.reconsiderEvidence(r)
		if err != nil {
			return err
		}
		decision := Decision{Kind: "ply.integration.decision", SchemaVersion: 1, IntegrationID: r.ID, PlanSHA256: r.PlanSHA256, Human: human}
		path := filepath.Join(storeRoot(root), "reconsiderations", r.ID, strings.TrimPrefix(digest(decision), "sha256:"), "decision.json")
		if err = atomicJSON(path, decision, true); err != nil {
			return err
		}
		r.Reconsideration = &ReconsiderationReceipt{Kind: "ply.integration.reconsideration", SchemaVersion: 1, IntegrationID: r.ID, PlanSHA256: r.PlanSHA256, Decision: decision, DecisionFile: taskrun.FileBinding{Locator: path, SHA256: digest(decision)}, Evidence: evidence, RecordedAtUTC: s.now()}
		r.State, r.Reasons, r.NextAction = "reconsidering", []string{}, "Preserve the actual later negative decision before releasing this effect reservation."
		if err = s.save(&r); err != nil {
			return err
		}
		if err = s.fault("after_reconsideration_decision"); err != nil {
			return err
		}
		return s.finishReconsider(d, &r)
	})
	return r, err
}

func (s *Service) finishReconsider(d taskrun.Dependencies, r *Receipt) error {
	c := r.Reconsideration
	if c == nil {
		return fmt.Errorf("reconsideration is missing")
	}
	if err := validateReconsideration(*r, false); err != nil {
		return err
	}
	if r.Integrated || r.InstallationStarted || r.Closeout != nil {
		return fmt.Errorf("reconsider_forbidden: an integration or later effect is already observed")
	}
	if c.Native == nil {
		// Recover an interrupted original acceptance before recording its later
		// negative answer; neither operation grants an effect on this path.
		if r.HumanAcceptance == nil {
			if err := s.accept(d, r); err != nil {
				return err
			}
		}
		a, err := taskrun.AcceptHumanIntegration(d, r.Plan.Workspace, r.Plan.WorkflowRunID, taskrun.HumanIntegrationInput{PlanID: r.ID, PlanSHA256: r.PlanSHA256, Decision: c.DecisionFile, Attestation: c.Decision.Human})
		if err != nil {
			return err
		}
		if a.Accepted || a.HumanQA.Outcome != c.Decision.Human.Answer {
			return fmt.Errorf("native negative decision was not observed")
		}
		c.Native = &a
		if err := s.save(r); err != nil {
			return err
		}
	}
	if err := s.fault("after_reconsideration_native"); err != nil {
		return err
	}
	proof := *c
	proof.Proof = nil
	path := filepath.Join(filepath.Dir(c.DecisionFile.Locator), "revocation.json")
	if err := atomicJSON(path, proof, true); err != nil {
		return err
	}
	c.Proof = &taskrun.FileBinding{Locator: path, SHA256: digest(proof)}
	r.State, r.Reasons, r.NextAction = "cancelled", []string{}, "The later negative human answer revoked this unexecuted plan. Correct or revise the candidate/choices, release source ownership again, then preview and confirm a new plan."
	s.event(r, "integration.cancelled")
	return s.save(r)
}

func validateReconsideration(r Receipt, requireProof bool) error {
	c := r.Reconsideration
	if c == nil {
		if requireProof {
			return fmt.Errorf("cancelled operation lacks its revocation proof")
		}
		return nil
	}
	if c.Kind != "ply.integration.reconsideration" || c.SchemaVersion != 1 || c.IntegrationID != r.ID || c.PlanSHA256 != r.PlanSHA256 || c.Decision.Kind != "ply.integration.decision" || c.Decision.SchemaVersion != 1 || c.Decision.IntegrationID != r.ID || c.Decision.PlanSHA256 != r.PlanSHA256 || (c.Decision.Human.Answer != "fail" && c.Decision.Human.Answer != "blocked") || r.Integrated || r.InstallationStarted || r.Closeout != nil {
		return fmt.Errorf("reconsideration identity or no-effect boundary differs")
	}
	raw, err := taskrun.Canonical(c.Decision.Human)
	if err != nil {
		return err
	}
	if _, err = workflowhandoff.ValidateDeliveryHumanAttestation(raw, string(r.Plan.TaskResult.TaskID), r.Plan.TaskResult, c.Decision.Human.Answer); err != nil {
		return err
	}
	if err = checkFile(&c.DecisionFile); err != nil {
		return err
	}
	if c.DecisionFile.SHA256 != digest(c.Decision) {
		return fmt.Errorf("negative decision bytes differ")
	}
	switch c.Evidence.Reason {
	case "never_dispatched":
		if c.Evidence.Started || r.IntegrationStarted {
			return fmt.Errorf("unstarted proof conflicts with the attempt journal")
		}
		if r.Plan.Flow == "github_pr_merge" && (r.Plan.PR == nil || c.Evidence.Current == nil || !exactRetryPR(*r.Plan.PR, *c.Evidence.Current) || (c.Evidence.Current.State != "ready" && c.Evidence.Current.State != "blocked")) {
			return fmt.Errorf("unstarted PR cancellation lacks a fresh exact unmerged observation")
		}
	case "request_rejected":
		if !c.Evidence.Started || !r.IntegrationStarted || r.Plan.PR == nil || c.Evidence.Prior == nil || !c.Evidence.Prior.NoEffect || !exactRetryPR(*r.Plan.PR, *c.Evidence.Prior) || (c.Evidence.Prior.State != "ready" && c.Evidence.Prior.State != "blocked") || c.Evidence.Current == nil || !exactRetryPR(*r.Plan.PR, *c.Evidence.Current) || (c.Evidence.Current.State != "ready" && c.Evidence.Current.State != "blocked") {
			return fmt.Errorf("reconsideration lacks exact proven rejected-request evidence")
		}
	default:
		return fmt.Errorf("reconsideration has no supported no-effect proof")
	}
	if c.Native != nil {
		n := c.Native
		if n.Kind != "PlyHumanIntegrationDecision@1" || n.SchemaVersion != 1 || n.RunID != r.Plan.WorkflowRunID || n.Accepted || n.PlanID != r.ID || n.PlanSHA256 != r.PlanSHA256 || digest(n.Decision) != digest(c.DecisionFile) || digest(n.Attestation) != digest(c.Decision.Human) || n.HumanQA.TaskID != r.Plan.TaskResult.TaskID || n.HumanQA.Outcome != c.Decision.Human.Answer || n.HumanQA.TaskResultID != r.Plan.TaskResult.ID || n.HumanQA.ResultOID != r.Plan.TaskResult.ResultOID || n.HumanQA.ResultTree != r.Plan.TaskResult.ResultTree {
			return fmt.Errorf("reconsideration native negative answer is not bound to this candidate")
		}
	}
	if requireProof || c.Proof != nil {
		if c.Proof == nil || c.Native == nil {
			return fmt.Errorf("cancelled operation lacks a completed native revocation")
		}
		if err = checkFile(c.Proof); err != nil {
			return err
		}
		copy := *c
		copy.Proof = nil
		if c.Proof.SHA256 != digest(copy) {
			return fmt.Errorf("revocation proof differs from the preserved negative decision")
		}
	}
	return nil
}

type effectSupersession struct {
	Kind                  string              `json:"kind"`
	SchemaVersion         int                 `json:"schema_version"`
	Key                   string              `json:"key"`
	PreviousIntegrationID string              `json:"previous_integration_id"`
	Revocation            taskrun.FileBinding `json:"revocation"`
	NextIntegrationID     string              `json:"next_integration_id"`
	PreviousSHA256        string              `json:"previous_sha256,omitempty"`
}

// reserveSupersession never rewrites the first owner. It follows an immutable
// proof chain and requires a fully preserved native negative decision for each
// retired owner before appending the next owner under the integration lock.
func reserveSupersession(r Receipt, initial effectReservation) error {
	if initial.Kind != "ply.integration.reservation" || initial.Key != r.Plan.EffectKey || !integrationID.MatchString(initial.IntegrationID) {
		return fmt.Errorf("integration reservation identity changed")
	}
	if r.State == "cancelled" || r.Reconsideration != nil {
		return fmt.Errorf("plan_revoked: a reconsidered operation cannot reclaim effect authority")
	}
	dir := filepath.Join(storeRoot(r.Plan.Workspace), "effects", strings.TrimPrefix(initial.Key, "sha256:")+"-supersessions")
	entries, err := os.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	current, previous, number := initial.IntegrationID, "", 0
	seen := map[string]bool{current: true}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".publish-") {
			continue
		}
		number++
		if entry.Name() != fmt.Sprintf("%06d.json", number) {
			return fmt.Errorf("supersession history is incomplete or ambiguous")
		}
		var proof effectSupersession
		if err = readJSON(filepath.Join(dir, entry.Name()), &proof); err != nil {
			return err
		}
		if proof.Kind != "ply.integration.effect-supersession" || proof.SchemaVersion != 1 || proof.Key != initial.Key || proof.PreviousIntegrationID != current || proof.PreviousSHA256 != previous || !integrationID.MatchString(proof.NextIntegrationID) || seen[proof.NextIntegrationID] {
			return fmt.Errorf("invalid effect supersession history")
		}
		old, err := readReceipt(r.Plan.Workspace, current)
		if err != nil {
			return err
		}
		if old.State != "cancelled" || old.Plan.EffectKey != initial.Key || old.Reconsideration == nil || old.Reconsideration.Proof == nil || digest(proof.Revocation) != digest(*old.Reconsideration.Proof) {
			return fmt.Errorf("supersession has no exact cancelled predecessor")
		}
		next, err := readReceipt(r.Plan.Workspace, proof.NextIntegrationID)
		if err != nil {
			return err
		}
		if next.Plan.EffectKey != initial.Key {
			return fmt.Errorf("supersession changes the reserved integration effect")
		}
		seen[proof.NextIntegrationID] = true
		current, previous = proof.NextIntegrationID, digest(proof)
	}
	if current == r.ID {
		return nil
	}
	old, err := readReceipt(r.Plan.Workspace, current)
	if err != nil {
		return err
	}
	if old.State != "cancelled" || old.Reconsideration == nil || old.Reconsideration.Proof == nil || old.Plan.EffectKey != r.Plan.EffectKey {
		return fmt.Errorf("effect_reserved: resume integration %s; a second candidate/target effect is not allowed", current)
	}
	proof := effectSupersession{Kind: "ply.integration.effect-supersession", SchemaVersion: 1, Key: r.Plan.EffectKey, PreviousIntegrationID: current, Revocation: *old.Reconsideration.Proof, NextIntegrationID: r.ID, PreviousSHA256: previous}
	return atomicJSON(filepath.Join(dir, fmt.Sprintf("%06d.json", number+1)), proof, true)
}
