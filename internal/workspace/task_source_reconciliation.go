package workspace

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Each reconciliation is a native source operation. Its immutable plan binds
// original identity and an exact parent; the terminal receipt adds an effective
// generation without modifying the preparation or the frozen Spec.
type TaskSourceReconciliation struct {
	// Positive native pre-effect fact, revoked durably before invoking Git.
	// Absence on historical records never establishes that an effect is unstarted.
	Unstarted      bool                         `yaml:"unstarted,omitempty" json:"unstarted,omitempty"`
	ConflictSHA256 string                       `yaml:"conflict_sha256" json:"conflict_sha256"`
	ID             string                       `yaml:"id" json:"id"`
	PlanSHA256     string                       `yaml:"plan_sha256" json:"plan_sha256"`
	Plan           TaskSourceReconciliationPlan `yaml:"plan" json:"plan"`
	State          string                       `yaml:"state" json:"state"`
	AfterOID       string                       `yaml:"after_oid" json:"after_oid"`
	AfterTree      string                       `yaml:"after_tree" json:"after_tree"`
	ConflictPaths  []string                     `yaml:"conflict_paths" json:"conflict_paths"`
	RecordedAtUTC  string                       `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}

type TaskSourceReconciliationPlan struct {
	AttemptOrdinal int             `yaml:"attempt_ordinal" json:"attempt_ordinal"`
	TaskID         TaskID          `yaml:"task_id" json:"task_id"`
	PreparationID  string          `yaml:"preparation_id" json:"preparation_id"`
	RunID          string          `yaml:"run_id" json:"run_id"`
	RequestSHA256  string          `yaml:"request_sha256" json:"request_sha256"`
	MandateSHA256  string          `yaml:"mandate_sha256" json:"mandate_sha256"`
	Generation     int             `yaml:"generation" json:"generation"`
	PreviousSHA256 string          `yaml:"previous_sha256" json:"previous_sha256"`
	BeforeOID      string          `yaml:"before_oid" json:"before_oid"`
	BeforeTree     string          `yaml:"before_tree" json:"before_tree"`
	PreviousBase   EpicBaseVersion `yaml:"previous_base" json:"previous_base"`
	Parent         EpicBaseVersion `yaml:"parent" json:"parent"`
	ReviewSHA256   string          `yaml:"review_sha256" json:"review_sha256"`
}

type TaskSourceReconciliationInput struct {
	PreparationID, RunID, RequestSHA256, MandateSHA256, ReviewSHA256 string
	Apply                                                            bool
	Confirmation, OperationID, Action                                string
}

type TaskSourceReconciliationReadback struct {
	State        string                        `json:"state"`
	Plan         *TaskSourceReconciliationPlan `json:"plan"`
	Confirmation string                        `json:"confirmation"`
	Operation    *TaskSourceReconciliation     `json:"operation"`
	NextAction   string                        `json:"next_action"`
}

// EffectiveDeliveryBasis is a native generation, not an override of the frozen
// Task parent. Its digest is carried by every new review/candidate/result.
type EffectiveDeliveryBasis struct {
	PreparationID        string          `yaml:"preparation_id" json:"preparation_id"`
	Generation           int             `yaml:"generation" json:"generation"`
	ReconciliationSHA256 string          `yaml:"reconciliation_sha256" json:"reconciliation_sha256"`
	Parent               EpicBaseVersion `yaml:"parent" json:"parent"`
}

func EffectiveTaskDeliveryBasis(r WorkItemRegistry, preparation string) (EffectiveDeliveryBasis, error) {
	p := findPreparation(r, preparation)
	if p == nil || p.Plan.SchemaVersion != 2 {
		return EffectiveDeliveryBasis{}, queueError("parallel_authority_required", "preparation has no frozen source reconciliation authority")
	}
	var parent EpicBaseVersion
	for _, b := range r.EpicBaseVersions {
		if b.EpicID == p.Plan.Target.EpicID && b.RepoID == p.Plan.Target.RepoID && b.Revision == p.Plan.BaseRevision {
			parent = b
		}
	}
	out := EffectiveDeliveryBasis{PreparationID: preparation, Parent: parent}
	for _, op := range r.SourceReconciliations {
		if op.Plan.PreparationID != preparation || op.State == "aborted" {
			continue
		}
		if op.State != "committed" {
			return out, queueError("source_reconciliation_pending", "observe the same native reconciliation "+op.ID)
		}
		out.Generation = op.Plan.Generation
		out.ReconciliationSHA256 = queueDigest(op)
		out.Parent = op.Plan.Parent
	}
	return out, nil
}
func (s *TaskSpecSession) EffectiveDeliveryBasis(preparation string) (EffectiveDeliveryBasis, error) {
	return EffectiveTaskDeliveryBasis(s.Registry, preparation)
}

func validateReconciledSourceAncestry(d Dependencies, r WorkItemRegistry, repo RepoRecord, b EffectiveDeliveryBasis, oid string) error {
	for _, op := range r.SourceReconciliations {
		if op.State != "committed" || queueDigest(op) != b.ReconciliationSHA256 {
			continue
		}
		ok, err := d.IntegrationGit.CheckAncestor(repo, op.AfterOID, oid)
		if err != nil {
			return err
		}
		if !ok {
			return queueError("source_reconciliation_ancestry", "source discarded the committed native reconciliation history")
		}
	}
	return nil
}

func reconciliationMember(r WorkItemRegistry, in TaskSourceReconciliationInput) (*TaskPreparation, error) {
	p := findPreparation(r, in.PreparationID)
	if r.FormatVersion != ParallelRegistryFormat || p == nil || p.Plan.SchemaVersion != 2 || p.Outcome == nil || !foldQueue(r, p.Plan.QueueID).contains(p.ID) || in.RunID == "" || !digestPattern.MatchString(in.RequestSHA256) || !digestPattern.MatchString(in.MandateSHA256) || !digestPattern.MatchString(in.ReviewSHA256) {
		return nil, queueError("source_reconciliation_authority", "exact active parallel preparation and accepted Run authority required")
	}
	return p, nil
}
func sourceReconciliationPlan(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, in TaskSourceReconciliationInput) (TaskSourceReconciliationPlan, error) {
	var plan TaskSourceReconciliationPlan
	p, err := reconciliationMember(r, in)
	if err != nil {
		return plan, err
	}
	basis, err := EffectiveTaskDeliveryBasis(r, p.ID)
	if err != nil {
		return plan, err
	}
	ep, binding := queueTargetBinding(r, p.Plan.Target)
	base := currentEpicBase(r, ep, binding)
	if base.Revision == basis.Parent.Revision {
		return plan, queueError("source_reconciliation_unnecessary", "source already has the registered effective parent")
	}
	if pendingBaseUpdate(r, ep.ID, binding.RepoID) || len(unresolvedTargetIntegrations(r, p.Plan.Target)) > 0 {
		return plan, queueError("shared_target_effect_pending", "resolve the existing target attempt before source reconciliation")
	}
	parent, err := queueParentObservation(d, p.Plan.Target)
	if err != nil {
		return plan, err
	}
	if parent.OID != base.OID || parent.Tree != base.Tree {
		return plan, queueError("epic_base_changed", "registered and observed parent differ")
	}
	source, err := d.IntegrationGit.ObserveIntegrationWorktree(p.Plan.WorktreePath, "refs/heads/"+p.Plan.Branch)
	if err != nil {
		return plan, err
	}
	if !source.Clean || len(source.InProgress) > 0 || source.GitCommonDir != p.Plan.Target.GitCommonDir || source.Ref != "refs/heads/"+p.Plan.Branch {
		return plan, queueError("source_reconciliation_dirty", "commit the owned source before reconcile-source")
	}
	_, repo, err := projectAndRepo(projects, p.Plan.Target.ProjectID, p.Plan.Target.RepoID)
	if err != nil {
		return plan, err
	}
	for _, pair := range [][2]string{{basis.Parent.OID, base.OID}, {p.Plan.SourceStartOID(), source.OID}} {
		ok, e := d.IntegrationGit.CheckAncestor(repo, pair[0], pair[1])
		if e != nil {
			return plan, e
		}
		if !ok {
			return plan, queueError("source_reconciliation_ancestry", "original source and preceding parent ancestry must be preserved")
		}
	}
	if err = validateReconciledSourceAncestry(d, r, repo, basis, source.OID); err != nil {
		return plan, err
	}
	ordinal := 1
	for _, old := range r.SourceReconciliations {
		if old.Plan.PreparationID == p.ID {
			ordinal++
		}
	}
	plan = TaskSourceReconciliationPlan{AttemptOrdinal: ordinal, TaskID: p.Plan.TaskID, PreparationID: p.ID, RunID: in.RunID, RequestSHA256: in.RequestSHA256, MandateSHA256: in.MandateSHA256, Generation: basis.Generation + 1, PreviousSHA256: basis.ReconciliationSHA256, BeforeOID: source.OID, BeforeTree: source.Tree, PreviousBase: basis.Parent, Parent: base, ReviewSHA256: in.ReviewSHA256}
	return plan, nil
}

// ReconcileTaskSource reserves a native source operation before Git, releases
// metadata locks while conflicts await the owner, and observes the same effect
// after interruption. The taskrun callback validates owner/session/stops first.
func ReconcileTaskSource(d Dependencies, in TaskSourceReconciliationInput) (TaskSourceReconciliationReadback, error) {
	var out TaskSourceReconciliationReadback
	if in.Action == "cancel-unstarted" && (!in.Apply || in.OperationID == "") {
		return out, WorkInvalidArguments("cancel-unstarted requires the exact --operation ID")
	}
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	locked := func(fn func(WorkItemStoreSession, WorkItemRegistry, ProjectSnapshot) error) error {
		return d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
			return d.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
				r, e := s.Snapshot()
				if e != nil {
					return e
				}
				return fn(s, r, projects)
			})
		})
	}
	var op TaskSourceReconciliation
	var prep TaskPreparation
	err = locked(func(s WorkItemStoreSession, r WorkItemRegistry, projects ProjectSnapshot) error {
		p, e := reconciliationMember(r, in)
		if e != nil {
			return e
		}
		prep = *p
		for _, old := range r.SourceReconciliations {
			if old.Plan.PreparationID != p.ID {
				continue
			}
			if in.OperationID == old.ID || in.Confirmation == old.PlanSHA256 {
				if old.Plan.RunID != in.RunID || old.Plan.RequestSHA256 != in.RequestSHA256 || old.Plan.MandateSHA256 != in.MandateSHA256 {
					return queueError("source_reconciliation_owner_conflict", "operation belongs to another exact Run")
				}
				op = old
				out = sourceReconciliationReadback(old)
				return nil
			}
		}
		if in.OperationID != "" || in.Action != "" {
			return queueError("source_reconciliation_not_found", "continue/abort/cancel-unstarted requires the same recorded operation")
		}
		plan, e := sourceReconciliationPlan(d, root, r, projects, in)
		if e != nil {
			return e
		}
		digest := queueDigest(plan)
		out = TaskSourceReconciliationReadback{State: "preview", Plan: &plan, Confirmation: digest, NextAction: "Apply this exact reconcile-source confirmation from the accepted owner context."}
		if !in.Apply {
			return nil
		}
		if in.Confirmation != digest {
			return queueError("source_reconciliation_stale", "source, parent, review or effective basis changed before apply")
		}
		op = TaskSourceReconciliation{Unstarted: true, ID: "src_" + strings.TrimPrefix(digest, "sha256:"), PlanSHA256: digest, Plan: plan, State: "intent", ConflictPaths: []string{}, RecordedAtUTC: queueNow(d)}
		r.SourceReconciliations = append(r.SourceReconciliations, op)
		return s.Publish(r)
	})
	if err != nil || !in.Apply {
		return out, err
	}
	if op.State == "committed" || op.State == "aborted" {
		return sourceReconciliationReadback(op), nil
	}
	if op.State == "effect_unknown" {
		return sourceReconciliationReadback(op), queueError("reconciliation_effect_unknown", "observe this operation; no Git effect will be replayed")
	}
	adapter, ok := d.IntegrationGit.(TaskSourceReconciliationGit)
	if !ok {
		return out, queueError("source_reconciliation_adapter_unknown", "native source Git adapter unavailable")
	}
	if d.SourceReconciliationFault != nil {
		if e := d.SourceReconciliationFault("after_intent", op.ID); e != nil {
			return sourceReconciliationReadback(op), e
		}
	}
	err = locked(func(s WorkItemStoreSession, r WorkItemRegistry, projects ProjectSnapshot) error {
		if _, e := reconciliationMember(r, in); e != nil {
			return e
		}
		for i, old := range r.SourceReconciliations {
			if old.ID != op.ID {
				continue
			}
			if old.State == "committed" || old.State == "aborted" {
				op = old
				return nil
			}
			if old.State == "effect_unknown" {
				op = old
				return queueError("reconciliation_effect_unknown", "observe the exact existing operation")
			}
			if !contentTypedEqual(old, op) {
				return queueError("source_reconciliation_conflict", "refresh this same operation before applying changed state")
			}
			// A crash or lost reply after this publication cannot leave a
			// cancellation proof attached to an attempted or ambiguous effect.
			// Abort/continue require an already pending merge and cannot start
			// an intent. Their pre-effect refusal retains its positive fact.
			if old.Unstarted && in.Action == "" {
				r.SourceReconciliations[i].Unstarted = false
				if e := s.Publish(r); e != nil {
					return e
				}
				op = r.SourceReconciliations[i]
			}
			observed, e := adapter.ReconcileSource(prep.Plan, old, in.Action)
			if e != nil {
				// Restore only original positive proof with the owning adapter's
				// typed read-only outcome. A generic error/flag cannot invent it.
				var noEffect *sourceReconciliationNoEffectError
				if old.Unstarted && observed.Unstarted && errors.As(e, &noEffect) {
					r.SourceReconciliations[i] = observed
					if publishErr := s.Publish(r); publishErr != nil {
						return publishErr
					}
					op = observed
				}
				return e
			}
			if d.SourceReconciliationFault != nil {
				if e = d.SourceReconciliationFault("after_git", op.ID); e != nil {
					return e
				}
			}
			r.SourceReconciliations[i] = observed
			op = observed
			return s.Publish(r)
		}
		return queueError("source_reconciliation_not_found", "recorded operation disappeared")
	})
	out = sourceReconciliationReadback(op)
	if err == nil && op.State == "effect_unknown" {
		err = queueError("reconciliation_effect_unknown", "source observation cannot prove this merge; do not replay")
	}
	return out, err
}
func sourceReconciliationReadback(op TaskSourceReconciliation) TaskSourceReconciliationReadback {
	action := "Request fresh separate native review, then verify this exact effective generation before return."
	if op.State == "source_conflict" {
		action = "Resolve and stage the recorded conflict paths in this source, then reconcile-source --continue; --abort refuses unexpected edits."
		if !digestPattern.MatchString(op.ConflictSHA256) {
			action = "The original conflict baseline was not published; abort is refused. Preserve source edits, resolve and stage the conflict, then use native --continue."
		}
	}
	if op.State == "aborted" {
		action = "Request a fresh reconcile-source preview for the current parent; this aborted operation adds no effective generation."
	}
	if op.State == "intent" || op.State == "effect_unknown" {
		action = "Observe/resume only this exact native operation; unknown effects cannot be replayed."
	}
	if op.State == "intent" && op.Unstarted {
		action = "Recheck the same known unstarted intent, or explicitly --cancel-unstarted --operation ID before a fresh parent preview. --abort is only for a known pending merge."
	}
	return TaskSourceReconciliationReadback{State: op.State, Plan: &op.Plan, Confirmation: op.PlanSHA256, Operation: &op, NextAction: action}
}
func unresolvedTargetIntegrations(r WorkItemRegistry, t QueueTarget) []string {
	out := []string{}
	for _, i := range r.IntegrationIntents {
		if i.GitCommonDir == t.GitCommonDir && i.ParentRef == t.ParentRef && i.State != "completed" {
			resolved := false
			for _, v := range r.IntegrationResults {
				if v.AuthorityID == i.AuthorityID {
					resolved = true
				}
			}
			if !resolved {
				out = append(out, string(i.ID))
			}
		}
	}
	return out
}
func validateSourceReconciliations(r WorkItemRegistry) error {
	if r.FormatVersion != ParallelRegistryFormat {
		return nil
	}
	byPreparation := map[string]EffectiveDeliveryBasis{}
	seen := map[string]bool{}
	pending := map[string]bool{}
	ordinals := map[string]int{}
	for _, op := range r.SourceReconciliations {
		p := findPreparation(r, op.Plan.PreparationID)
		plan := op.Plan
		if p == nil || p.Plan.SchemaVersion != 2 || p.Plan.TaskID != plan.TaskID || seen[op.ID] || op.ID != "src_"+strings.TrimPrefix(op.PlanSHA256, "sha256:") || op.PlanSHA256 != queueDigest(plan) || contentUTC(op.RecordedAtUTC) != nil || op.ConflictPaths == nil || !sort.StringsAreSorted(op.ConflictPaths) || plan.RunID == "" || !digestPattern.MatchString(plan.RequestSHA256) || !digestPattern.MatchString(plan.MandateSHA256) || !digestPattern.MatchString(plan.ReviewSHA256) || !validOIDText(plan.BeforeOID) || !validOIDText(plan.BeforeTree) {
			return fmt.Errorf("invalid source reconciliation binding")
		}
		seen[op.ID] = true
		if op.Unstarted && op.State != "intent" {
			return fmt.Errorf("unstarted proof on a started or terminal source operation")
		}
		ordinals[p.ID]++
		if plan.AttemptOrdinal != ordinals[p.ID] {
			return fmt.Errorf("broken reconciliation attempt ordinal")
		}
		basis, ok := byPreparation[p.ID]
		if !ok {
			var parent EpicBaseVersion
			for _, b := range r.EpicBaseVersions {
				if b.EpicID == p.Plan.Target.EpicID && b.RepoID == p.Plan.Target.RepoID && b.Revision == p.Plan.BaseRevision {
					parent = b
				}
			}
			basis = EffectiveDeliveryBasis{PreparationID: p.ID, Parent: parent}
		}
		if pending[p.ID] || plan.Generation != basis.Generation+1 || plan.PreviousSHA256 != basis.ReconciliationSHA256 || !contentTypedEqual(plan.PreviousBase, basis.Parent) || plan.Parent.Revision <= basis.Parent.Revision || !historicalEpicBaseMatches(r, p.Plan.Target.EpicID, p.Plan.Target.RepoID, p.Plan.Target.ParentWorktreeID, p.Plan.Target.ParentRef, plan.Parent.OID, plan.Parent.Tree) {
			return fmt.Errorf("broken source reconciliation generation")
		}
		switch op.State {
		case "committed":
			if !validOIDText(op.AfterOID) || !validOIDText(op.AfterTree) {
				return fmt.Errorf("missing merge outcome")
			}
			basis.Generation = plan.Generation
			basis.Parent = plan.Parent
			basis.ReconciliationSHA256 = queueDigest(op)
		case "aborted":
			if op.AfterOID != plan.BeforeOID || op.AfterTree != plan.BeforeTree {
				return fmt.Errorf("abort did not restore exact pre-state")
			}
		case "intent", "source_conflict", "effect_unknown":
			pending[p.ID] = true
		default:
			return fmt.Errorf("invalid reconciliation state")
		}
		byPreparation[p.ID] = basis
	}
	return nil
}
func validateSourceReconciliationTransition(old, next WorkItemRegistry) error {
	for _, op := range old.SourceReconciliations {
		found := false
		for _, n := range next.SourceReconciliations {
			if n.ID != op.ID {
				continue
			}
			found = true
			if !contentTypedEqual(n.Plan, op.Plan) || n.PlanSHA256 != op.PlanSHA256 || n.RecordedAtUTC != op.RecordedAtUTC {
				return fmt.Errorf("immutable reconciliation intent changed")
			}
			if op.State == "source_conflict" && n.ConflictSHA256 != op.ConflictSHA256 {
				return fmt.Errorf("original source conflict evidence changed")
			}
			if op.State == "committed" || op.State == "aborted" || op.State == "effect_unknown" {
				if !contentTypedEqual(n, op) {
					return fmt.Errorf("terminal reconciliation changed")
				}
			}
		}
		if !found {
			return fmt.Errorf("source reconciliation removed")
		}
	}
	return nil
}

func ValidateEffectiveDeliveryBasis(r WorkItemRegistry, b EffectiveDeliveryBasis) error {
	p := findPreparation(r, b.PreparationID)
	if p == nil || p.Plan.SchemaVersion != 2 {
		return fmt.Errorf("effective basis has no parallel preparation")
	}
	if b.Generation == 0 {
		for _, base := range r.EpicBaseVersions {
			if base.EpicID == p.Plan.Target.EpicID && base.RepoID == p.Plan.Target.RepoID && base.Revision == p.Plan.BaseRevision && contentTypedEqual(base, b.Parent) && b.ReconciliationSHA256 == "" {
				return nil
			}
		}
	} else {
		for _, op := range r.SourceReconciliations {
			if op.Plan.PreparationID == p.ID && op.State == "committed" && op.Plan.Generation == b.Generation && queueDigest(op) == b.ReconciliationSHA256 && contentTypedEqual(op.Plan.Parent, b.Parent) {
				return nil
			}
		}
	}
	return fmt.Errorf("effective basis is not a native committed generation")
}
func effectiveTaskParent(r WorkItemRegistry, t TaskRecord, result TaskResultID) (TaskRecord, error) {
	link := taskResultSpecLink(r, result)
	if link == nil || link.EffectiveBasis == nil {
		return t, nil
	}
	b := *link.EffectiveBasis
	if err := ValidateEffectiveDeliveryBasis(r, b); err != nil {
		return t, err
	}
	p := findPreparation(r, b.PreparationID)
	if p.Plan.TaskID != t.ID {
		return t, fmt.Errorf("foreign effective basis")
	}
	copy := *t.Worktree
	t.Worktree = &copy
	t.Worktree.ParentOID = b.Parent.OID
	t.Worktree.ParentTree = b.Parent.Tree
	return t, nil
}
