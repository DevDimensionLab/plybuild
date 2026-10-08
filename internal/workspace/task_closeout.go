package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const taskCloseoutsFile = "task-closeouts.json"

// TaskCloseoutOwnership is supplied by the native authority service after it
// validates an explicit release of writing ownership. A clock or quiet process
// is never release evidence. This primitive additionally checks the binding and
// the preserved bytes; it does not manufacture the authority service's decision.
type TaskCloseoutOwnership struct {
	Released        bool         `json:"released"`
	TaskID          TaskID       `json:"task_id"`
	TaskResultID    TaskResultID `json:"task_result_id"`
	ResultOID       string       `json:"result_oid"`
	ResultTree      string       `json:"result_tree"`
	EvidenceLocator string       `json:"evidence_locator"`
	EvidenceSHA256  string       `json:"evidence_sha256"`
}

// TaskCloseoutPRIntegration carries the native adapter's observed exact merge.
// The caller validates repository/head/base against its approved Delivery plan.
// Squash is deliberately supported without inventing local ancestry evidence.
type TaskCloseoutPRIntegration struct {
	Repository      string `json:"repository"`
	Number          int    `json:"number"`
	HeadOID         string `json:"head_oid"`
	BaseRef         string `json:"base_ref"`
	MergeOID        string `json:"merge_oid"`
	URL             string `json:"url"`
	MergedAtUTC     string `json:"merged_at_utc"`
	EvidenceLocator string `json:"evidence_locator"`
	EvidenceSHA256  string `json:"evidence_sha256"`
}

type TaskCloseoutInput struct {
	OperationID        string
	TaskID             TaskID
	TaskResultID       TaskResultID
	ExpectedResultOID  string
	ExpectedResultTree string
	TargetRef          string
	Keep               bool
	Ownership          TaskCloseoutOwnership
	ObservedPR         *TaskCloseoutPRIntegration
	RequireIntegration bool // Preview only. Apply always requires an observed effect.
}

type TaskCloseoutPlan struct {
	Kind              string                `json:"kind"`
	SchemaVersion     int                   `json:"schema_version"`
	WorkspaceRoot     string                `json:"workspace_root"`
	TaskID            TaskID                `json:"task_id"`
	TaskResultID      TaskResultID          `json:"task_result_id"`
	ResultOID         string                `json:"result_oid"`
	ResultTree        string                `json:"result_tree"`
	Source            TaskWorktreeBinding   `json:"source"`
	RepositoryLocator string                `json:"repository_locator"`
	ReturnLocator     string                `json:"return_locator"`
	TargetRef         string                `json:"target_ref"`
	Keep              bool                  `json:"keep"`
	Ownership         TaskCloseoutOwnership `json:"ownership"`
	ArchiveRef        string                `json:"archive_ref"`
}

type TaskCloseoutPreview struct {
	Plan                TaskCloseoutPlan     `json:"plan"`
	PlanSHA256          string               `json:"plan_sha256"`
	Ready               bool                 `json:"ready"`
	Reasons             []string             `json:"reasons"`
	Details             []string             `json:"details"`
	NextAction          string               `json:"next_action"`
	IntegrationObserved bool                 `json:"integration_observed"`
	Receipt             *TaskCloseoutReceipt `json:"receipt"`
}

type TaskCloseoutReceipt struct {
	Kind                string                     `json:"kind"`
	SchemaVersion       int                        `json:"schema_version"`
	OperationID         string                     `json:"operation_id"`
	Plan                TaskCloseoutPlan           `json:"plan"`
	PlanSHA256          string                     `json:"plan_sha256"`
	State               string                     `json:"state"`          // cleanup_pending or complete
	ResourceState       string                     `json:"resource_state"` // retained, removed, retired, or kept
	Phase               string                     `json:"phase"`
	NativeIntegrationID IntegrationResultID        `json:"native_integration_id,omitempty"`
	ObservedPR          *TaskCloseoutPRIntegration `json:"observed_pr,omitempty"`
	Retention           *TaskRetentionManifest     `json:"retention"`
	WorktreeRemoved     bool                       `json:"worktree_removed"`
	BranchRemoved       bool                       `json:"branch_removed"`
	LifecycleCompleted  bool                       `json:"lifecycle_completed"`
	CreatedAtUTC        string                     `json:"created_at_utc"`
	UpdatedAtUTC        string                     `json:"updated_at_utc"`
	RemovedAtUTC        string                     `json:"removed_at_utc,omitempty"`
	CompletedAtUTC      string                     `json:"completed_at_utc,omitempty"`
}

type taskCloseoutStore struct {
	Kind          string                `json:"kind"`
	SchemaVersion int                   `json:"schema_version"`
	Receipts      []TaskCloseoutReceipt `json:"receipts"`
}

func readTaskCloseouts(root string) (taskCloseoutStore, error) {
	s := taskCloseoutStore{Kind: "WorkspaceTaskCloseoutStore@1", SchemaVersion: 1, Receipts: []TaskCloseoutReceipt{}}
	b, exists, err := readWorkspaceMetadata(root, taskCloseoutsFile)
	if err != nil || !exists {
		return s, err
	}
	if err = decodeWorkspaceMetadata(b, &s); err != nil {
		return s, err
	}
	if s.Kind != "WorkspaceTaskCloseoutStore@1" || s.SchemaVersion != 1 || s.Receipts == nil {
		return s, fmt.Errorf("invalid closeout store")
	}
	seen := map[TaskID]bool{}
	for _, r := range s.Receipts {
		if seen[r.Plan.TaskID] || r.Kind != "WorkspaceTaskCloseoutReceipt@1" || r.SchemaVersion != 1 || r.PlanSHA256 != closeoutPlanDigest(r.Plan) || r.Plan.WorkspaceRoot != root || !validTaskUTC(r.CreatedAtUTC) || !validTaskUTC(r.UpdatedAtUTC) || (r.State != "cleanup_pending" && r.State != "complete") || (r.State == "complete" && (!r.LifecycleCompleted || r.Retention == nil || (!r.Plan.Keep && (!r.WorktreeRemoved || !r.BranchRemoved || r.ResourceState != "retired")) || (r.Plan.Keep && r.ResourceState != "kept"))) {
			return s, fmt.Errorf("invalid or conflicting closeout receipt")
		}
		seen[r.Plan.TaskID] = true
	}
	return s, nil
}

func (s *taskCloseoutStore) receipt(id TaskID) *TaskCloseoutReceipt {
	for i := range s.Receipts {
		if s.Receipts[i].Plan.TaskID == id {
			return &s.Receipts[i]
		}
	}
	return nil
}

func (s *taskCloseoutStore) publish(root string) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return (workspaceMetadataPublisher{}).publish(root, taskCloseoutsFile, append(b, '\n'))
}

func closeoutPlanDigest(p TaskCloseoutPlan) string {
	b, _ := json.Marshal(p)
	return digestTaskBytes(b)
}

func ReadTaskCloseout(d Dependencies, id TaskID) (*TaskCloseoutReceipt, error) {
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return nil, err
	}
	return ReadTaskCloseoutAt(root, id)
}

// ReadTaskCloseoutAt reads only preserved facts and works after source removal.
func ReadTaskCloseoutAt(root string, id TaskID) (*TaskCloseoutReceipt, error) {
	if _, err := ParseTaskID(string(id)); err != nil {
		return nil, err
	}
	s, err := readTaskCloseouts(root)
	if err != nil {
		return nil, err
	}
	return s.receipt(id), nil
}

func ReadTaskCloseoutsAt(root string) ([]TaskCloseoutReceipt, error) {
	s, err := readTaskCloseouts(root)
	return s.Receipts, err
}

func PreviewTaskCloseout(d Dependencies, in TaskCloseoutInput) (TaskCloseoutPreview, error) {
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return TaskCloseoutPreview{}, err
	}
	if d.WorkItems == nil || d.Projects == nil {
		return TaskCloseoutPreview{}, fmt.Errorf("closeout stores unavailable")
	}
	r, err := d.WorkItems.Snapshot(root)
	if err != nil {
		return TaskCloseoutPreview{}, err
	}
	return previewTaskCloseout(d, root, r, in)
}

func previewTaskCloseout(d Dependencies, root string, r WorkItemRegistry, in TaskCloseoutInput) (TaskCloseoutPreview, error) {
	out := TaskCloseoutPreview{Reasons: []string{}, Details: []string{}, NextAction: "Confirm the closeout plan after integration and required installation."}
	if _, err := ParseTaskID(string(in.TaskID)); err != nil {
		return out, err
	}
	if !validOIDText(in.ExpectedResultOID) || !validOIDText(in.ExpectedResultTree) || !validFullBranchRef(in.TargetRef) {
		return out, WorkInvalidArguments("closeout requires exact candidate OID, tree and full target ref")
	}
	task, _ := findTask(r, in.TaskID)
	tr := findTaskResult(r, in.TaskResultID)
	if task == nil || task.Worktree == nil || tr == nil || tr.TaskID != task.ID || tr.TaskWorktreeID != task.Worktree.ID || tr.ResultOID != in.ExpectedResultOID || tr.ResultTree != in.ExpectedResultTree || tr.SourceLocator != task.Worktree.Locator || tr.SourceRef != task.Worktree.Ref || tr.GitCommonDir != task.GitCommonDir {
		return out, closeoutError("candidate_changed", "Select the exact registered TaskResult and source before closeout.", nil)
	}
	ep, _ := findEpic(r, task.ParentEpicID)
	parent, _ := findEpicRepo(*ep, task.RepoID)
	projects, repos, err := d.Projects.Snapshot(root)
	if err != nil {
		return out, err
	}
	_, repo, err := projectAndRepo(ProjectSnapshot{Projects: projects, Repos: repos}, task.ProjectID, task.RepoID)
	if err != nil || repo.GitCommonDir != task.GitCommonDir {
		return out, closeoutError("cleanup_blocked", "Registered repository identity differs.", err)
	}
	out.Plan = TaskCloseoutPlan{Kind: "WorkspaceTaskCloseoutPlan@1", SchemaVersion: 1, WorkspaceRoot: root, TaskID: task.ID, TaskResultID: tr.ID, ResultOID: tr.ResultOID, ResultTree: tr.ResultTree, Source: *task.Worktree, RepositoryLocator: repo.Locator, ReturnLocator: parent.Worktree.Locator, TargetRef: in.TargetRef, Keep: in.Keep, Ownership: in.Ownership, ArchiveRef: "refs/ply/archive/tasks/" + string(task.ID) + "/" + tr.ResultOID}
	out.PlanSHA256 = closeoutPlanDigest(out.Plan)
	store, err := readTaskCloseouts(root)
	if err != nil {
		return out, err
	}
	out.Receipt = store.receipt(task.ID)
	if out.Receipt != nil && out.Receipt.PlanSHA256 != out.PlanSHA256 {
		return out, closeoutError("plan_changed", "Resume the preserved closeout plan; its candidate or resources differ.", nil)
	}
	if out.Receipt != nil && out.Receipt.State == "complete" {
		out.IntegrationObserved = out.Receipt.NativeIntegrationID != "" || out.Receipt.ObservedPR != nil
		out.Ready = out.IntegrationObserved
		if err := verifyCloseoutRetention(d, root, out.Plan, *out.Receipt.Retention); err != nil {
			out.Ready = false
			out.Reasons = append(out.Reasons, "evidence_retention_failed")
			out.Details = append(out.Details, err.Error())
		}
		out.NextAction = "The Task is closed; inspect its preserved receipt and evidence."
		return out, nil
	}
	if in.ObservedPR == nil && out.Receipt != nil {
		in.ObservedPR = out.Receipt.ObservedPR
	}
	add := func(code string) { out.Reasons = append(out.Reasons, code) }
	owner := in.Ownership
	if !owner.Released || owner.TaskID != task.ID || owner.TaskResultID != tr.ID || owner.ResultOID != tr.ResultOID || owner.ResultTree != tr.ResultTree {
		add("owner_active")
	} else if _, err = readCloseoutEvidence(owner.EvidenceLocator, owner.EvidenceSHA256); err != nil {
		add("ownership_evidence_unavailable")
	}
	_, out.IntegrationObserved, err = closeoutIntegrationEvidence(d, r, out.Plan, in.ObservedPR)
	if err != nil {
		add("integration_evidence_invalid")
	}
	if in.RequireIntegration && !out.IntegrationObserved {
		add("integration_not_observed")
	}
	if err := closeoutProtectResources(d, r, repos, out.Plan); err != nil {
		add("cleanup_blocked")
		out.Details = append(out.Details, err.Error())
	}
	removed := out.Receipt != nil && (out.Receipt.WorktreeRemoved || out.Receipt.Phase == "worktree_remove_reserved" || out.Receipt.Phase == "branch_remove_reserved")
	if err := closeoutObserveSource(out.Plan, removed, out.Receipt != nil && (out.Receipt.BranchRemoved || out.Receipt.Phase == "branch_remove_reserved")); err != nil {
		add("cleanup_blocked")
		out.Details = append(out.Details, err.Error())
	}
	if out.Receipt == nil || out.Receipt.Retention == nil {
		if _, err := collectCloseoutRetention(d, root, r, *tr, in); err != nil {
			add("evidence_retention_failed")
			out.Details = append(out.Details, err.Error())
		}
	} else if err := verifyCloseoutRetention(d, root, out.Plan, *out.Receipt.Retention); err != nil {
		add("evidence_retention_failed")
		out.Details = append(out.Details, err.Error())
	}
	out.Reasons = sortedReasons(out.Reasons)
	out.Ready = len(out.Reasons) == 0
	if !out.Ready {
		out.NextAction = "Resolve the reported ownership, evidence or resource mismatch, then check the same Task again. " + strings.Join(out.Details, "; ")
	}
	return out, nil
}

func closeoutIntegrationEvidence(d Dependencies, r WorkItemRegistry, p TaskCloseoutPlan, pr *TaskCloseoutPRIntegration) (IntegrationResultID, bool, error) {
	if pr != nil {
		if pr.Repository == "" || pr.Number <= 0 || pr.HeadOID != p.ResultOID || pr.BaseRef != p.TargetRef || !validOIDText(pr.MergeOID) || !strings.HasPrefix(pr.URL, "https://") || !validTaskUTC(pr.MergedAtUTC) {
			return "", false, fmt.Errorf("PR integration binding differs")
		}
		if _, err := readCloseoutEvidence(pr.EvidenceLocator, pr.EvidenceSHA256); err != nil {
			return "", false, err
		}
		return "", true, nil
	}
	for _, result := range r.IntegrationResults {
		a := findAuthority(r, result.AuthorityID)
		if a == nil || a.TaskID != p.TaskID || a.TaskResultID != p.TaskResultID || a.Plan.Task.ResultOID != p.ResultOID || a.Plan.Task.ResultTree != p.ResultTree || a.Plan.Task.SourceLocator != p.Source.Locator || a.Plan.Task.SourceRef != p.Source.Ref || a.Plan.Repository.GitCommonDir != p.Source.GitCommonDir || a.Plan.Epic.ParentRef != p.TargetRef {
			continue
		}
		if result.Outcome != "exact_effect" && result.Outcome != "already_integrated" {
			continue
		}
		if result.AfterObservation.OID != p.ResultOID || result.AfterObservation.Ref != p.TargetRef || result.AfterObservation.GitCommonDir != p.Source.GitCommonDir {
			continue
		}
		if err := closeoutCheckAncestor(p); err != nil {
			return "", false, err
		}
		return result.ID, true, nil
	}
	return "", false, nil
}

func closeoutError(code, next string, cause error) error {
	if cause != nil {
		next += " " + cause.Error()
	}
	return workError(ErrorWorkReconciliation, code+": "+next, cause)
}

func applyCloseoutIntegrationReadback(out *WorkspaceTaskIntegrationReadback, receipt TaskCloseoutReceipt) {
	if receipt.ObservedPR != nil {
		out.Classification = "pr_merged"
	}
	_, action, reason := TaskCloseoutPendingAction(receipt)
	out.NextAction = IntegrationNextAction{Kind: action, Reason: reason, Argv: []string{"ply", "integration", "resume", receipt.OperationID}}
	out.RecoveryStatus = "cleanup_pending"
	if receipt.State == "complete" {
		out.RecoveryStatus = "complete"
		out.NextAction = IntegrationNextAction{Kind: "none", Reason: "The Task is completed; inspect its retained evidence and closeout receipt.", Argv: []string{}}
	}
	replaceReadbackMember(&out.Value, "classification", out.Classification)
	replaceReadbackMember(&out.Value, "recovery_status", out.RecoveryStatus)
	replaceReadbackMember(&out.Value, "next_action", structCanonical(out.NextAction))
}

func ApplyTaskCloseout(d Dependencies, in TaskCloseoutInput, expectedPlanSHA256 string) (TaskCloseoutReceipt, error) {
	return applyTaskCloseout(d, in, expectedPlanSHA256, nil)
}

// ApplyTaskCloseoutWithOwnershipGuard closes the read-to-effect startup race.
// The guard must be read-only and must not acquire taskrun or WorkItems locks.
func ApplyTaskCloseoutWithOwnershipGuard(d Dependencies, in TaskCloseoutInput, expectedPlanSHA256 string, ownershipGuard func() error) (TaskCloseoutReceipt, error) {
	if ownershipGuard == nil {
		return TaskCloseoutReceipt{}, fmt.Errorf("current ownership guard is required")
	}
	return applyTaskCloseoutGuarded(d, in, expectedPlanSHA256, ownershipGuard, nil)
}

// The single workspace registry lock covers all Task and target resources. Each
// destructive effect is reserved durably before Git and observed before advance.
func applyTaskCloseout(d Dependencies, in TaskCloseoutInput, expected string, fault func(string) error) (out TaskCloseoutReceipt, err error) {
	return applyTaskCloseoutGuarded(d, in, expected, nil, fault)
}

func applyTaskCloseoutGuarded(d Dependencies, in TaskCloseoutInput, expected string, ownershipGuard func() error, fault func(string) error) (out TaskCloseoutReceipt, err error) {
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return out, err
	}
	if d.WorkItems == nil || d.WorkClock == nil {
		return out, fmt.Errorf("closeout dependencies unavailable")
	}
	if !validTaskText(in.OperationID, 1, 256) {
		return out, WorkInvalidArguments("closeout operation ID is required")
	}
	in.RequireIntegration = true
	err = d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
		r, err := session.Snapshot()
		if err != nil {
			return err
		}
		preview, err := previewTaskCloseout(d, root, r, in)
		if err != nil {
			return err
		}
		if preview.PlanSHA256 != expected {
			return closeoutError("plan_changed", "Check and confirm the current exact closeout resources.", nil)
		}
		if !preview.Ready {
			return closeoutError("cleanup_blocked", strings.Join(preview.Reasons, "; ")+". "+preview.NextAction, nil)
		}
		if preview.Receipt != nil && preview.Receipt.State == "complete" {
			out = *preview.Receipt
			return nil
		}
		if ownershipGuard != nil {
			if err = ownershipGuard(); err != nil {
				return err
			}
		}
		if preview.Receipt != nil && strings.HasPrefix(preview.Receipt.Phase, "install") {
			reason, _, detail := TaskCloseoutPendingAction(*preview.Receipt)
			return closeoutError(reason, detail, nil)
		}
		if in.ObservedPR == nil && preview.Receipt != nil {
			in.ObservedPR = preview.Receipt.ObservedPR
		}
		s, err := readTaskCloseouts(root)
		if err != nil {
			return err
		}
		receipt := s.receipt(in.TaskID)
		now := func() string { return d.WorkClock.Now().UTC().Format(time.RFC3339Nano) }
		if receipt == nil {
			native, _, err := closeoutIntegrationEvidence(d, r, preview.Plan, in.ObservedPR)
			if err != nil {
				return err
			}
			s.Receipts = append(s.Receipts, TaskCloseoutReceipt{Kind: "WorkspaceTaskCloseoutReceipt@1", SchemaVersion: 1, OperationID: in.OperationID, Plan: preview.Plan, PlanSHA256: expected, State: "cleanup_pending", ResourceState: "retained", Phase: "reserved", NativeIntegrationID: native, ObservedPR: in.ObservedPR, CreatedAtUTC: now(), UpdatedAtUTC: now()})
			receipt = &s.Receipts[len(s.Receipts)-1]
			if err := s.publish(root); err != nil {
				return err
			}
		}
		out = *receipt
		if receipt.State == "complete" {
			return nil
		}
		save := func(phase string) error {
			receipt.Phase, receipt.UpdatedAtUTC = phase, now()
			out = *receipt
			return s.publish(root)
		}
		check := func(stage string) error {
			if fault != nil {
				return fault(stage)
			}
			return nil
		}
		if receipt.Retention == nil {
			tr := findTaskResult(r, in.TaskResultID)
			entries, err := collectCloseoutRetention(d, root, r, *tr, in)
			if err != nil {
				return closeoutError("evidence_retention_failed", "Preserve readable evidence before removing the source.", err)
			}
			manifest, err := retainCloseoutEvidence(d, root, preview.Plan, entries)
			if err != nil {
				return closeoutError("evidence_retention_failed", "Preserve readable evidence before removing the source.", err)
			}
			receipt.Retention = &manifest
			if err := save("evidence_retained"); err != nil {
				return err
			}
		}
		if err := verifyCloseoutRetention(d, root, preview.Plan, *receipt.Retention); err != nil {
			return closeoutError("evidence_retention_failed", "Restore retained evidence before cleanup.", err)
		}
		if !in.Keep && !receipt.WorktreeRemoved {
			if err := closeoutObserveSource(preview.Plan, receipt.Phase == "worktree_remove_reserved", false); err != nil {
				return err
			}
			if err := save("worktree_remove_reserved"); err != nil {
				return err
			}
			if err := check("before-worktree-remove"); err != nil {
				return err
			}
			if err := closeoutRemoveWorktree(preview.Plan); err != nil {
				return err
			}
			if err := check("after-worktree-remove"); err != nil {
				return err
			}
			receipt.WorktreeRemoved, receipt.ResourceState, receipt.RemovedAtUTC = true, "removed", now()
			if err := save("worktree_removed"); err != nil {
				return err
			}
		}
		if !in.Keep && !receipt.BranchRemoved {
			if err := save("branch_remove_reserved"); err != nil {
				return err
			}
			if err := check("before-branch-remove"); err != nil {
				return err
			}
			if err := closeoutRemoveBranch(preview.Plan); err != nil {
				return err
			}
			if err := check("after-branch-remove"); err != nil {
				return err
			}
			receipt.BranchRemoved = true
			if err := save("branch_removed"); err != nil {
				return err
			}
		}
		if err := check("before-lifecycle"); err != nil {
			return err
		}
		if !in.Keep {
			_, i := findTask(r, in.TaskID)
			if r.Tasks[i].WorktreeState != WorkItemRetired {
				r.Tasks[i].WorktreeState = WorkItemRetired
				if _, err := publishWorkItemRegistryRecover(session, r); err != nil {
					return err
				}
			}
		}
		if err := completeCloseoutLifecycle(root, in.TaskID, now()); err != nil {
			return err
		}
		if err := check("after-lifecycle"); err != nil {
			return err
		}
		receipt.LifecycleCompleted, receipt.State, receipt.CompletedAtUTC = true, "complete", now()
		receipt.ResourceState = "retired"
		if in.Keep {
			receipt.ResourceState = "kept"
		}
		return save("complete")
	})
	return out, err
}

func completeCloseoutLifecycle(root string, id TaskID, now string) error {
	s, err := readWorkItemLifecycle(root)
	if err != nil {
		return err
	}
	if s.Task(id) == LifecycleCompleted {
		return nil
	}
	e := WorkItemLifecycleEvent{Sequence: s.Revision + 1, SubjectKind: "task", SubjectID: string(id), From: s.Task(id), To: LifecycleCompleted, ActorClaim: "ply integration closeout", RecordedAtUTC: now}
	e.ID = lifecycleEventID(e)
	b, err := encodeWorkItemLifecycle(append(s.Events, e))
	if err != nil {
		return err
	}
	return (workspaceMetadataPublisher{}).publish(root, WorkItemLifecycleFile, b)
}

// PreviewTaskReconcile never chooses by timestamp. Only an unambiguous current
// native result with an exact integration receipt qualifies as legacy closeout.
func PreviewTaskReconcile(d Dependencies, id TaskID, ownership TaskCloseoutOwnership, keep bool) (TaskCloseoutPreview, error) {
	root, err := containingWorkItemWorkspace(d)
	if err != nil {
		return TaskCloseoutPreview{}, err
	}
	r, err := d.WorkItems.Snapshot(root)
	if err != nil {
		return TaskCloseoutPreview{}, err
	}
	fact := RecordedTaskProgressFacts(r)[id]
	if fact.Ambiguous || fact.TaskResult == nil || fact.IntegrationResult == nil || (fact.IntegrationResult.Outcome != "exact_effect" && fact.IntegrationResult.Outcome != "already_integrated") {
		return TaskCloseoutPreview{}, closeoutError("legacy_result_ambiguous", "Select and prove the exact historical native integration; no merge or QA is inferred.", nil)
	}
	a := findAuthority(r, fact.IntegrationResult.AuthorityID)
	if a == nil {
		return TaskCloseoutPreview{}, closeoutError("legacy_result_ambiguous", "Native integration authority is unavailable.", nil)
	}
	tr := fact.TaskResult
	return previewTaskCloseout(d, root, r, TaskCloseoutInput{TaskID: id, TaskResultID: tr.ID, ExpectedResultOID: tr.ResultOID, ExpectedResultTree: tr.ResultTree, TargetRef: a.Plan.Epic.ParentRef, Keep: keep, Ownership: ownership, RequireIntegration: true})
}

func pathWithin(parent, path string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func closeoutProtectResources(d Dependencies, r WorkItemRegistry, repos []RepoRecord, p TaskCloseoutPlan) error {
	s := p.Source
	if s.OwnerKind != "task" || s.OwnerID != p.TaskID || s.Origin != "ply_created" || s.GitCommonDir == "" || s.Locator == p.ReturnLocator || s.Ref == p.TargetRef || pathWithin(s.Locator, p.WorkspaceRoot) || pathWithin(s.Locator, s.GitCommonDir) || pathWithin(s.Locator, p.Ownership.EvidenceLocator) {
		return fmt.Errorf("source is not a distinct owned linked worktree")
	}
	for _, repo := range repos {
		if pathWithin(s.Locator, repo.Locator) {
			return fmt.Errorf("source contains a registered repository")
		}
	}
	metadata, err := readProjectMetadata(p.WorkspaceRoot)
	if err != nil {
		return err
	}
	for _, companions := range metadata.companions {
		for _, companion := range companions {
			if pathWithin(s.Locator, companion.Locator) || pathWithin(companion.Locator, s.Locator) {
				return fmt.Errorf("source overlaps a registered planning or companion directory")
			}
		}
	}
	for _, e := range r.Epics {
		for _, b := range e.RepoBindings {
			if pathWithin(s.Locator, b.Worktree.Locator) || pathWithin(b.Worktree.Locator, s.Locator) {
				return fmt.Errorf("source overlaps an Epic worktree")
			}
		}
	}
	for _, t := range r.Tasks {
		if t.ID != p.TaskID && t.Worktree != nil && (pathWithin(s.Locator, t.Worktree.Locator) || pathWithin(t.Worktree.Locator, s.Locator)) {
			return fmt.Errorf("source overlaps another Task")
		}
	}
	if !p.Keep {
		cwd, err := projectPhysicalWorkingDirectory(d.Files)
		if err != nil {
			return err
		}
		if pathWithin(s.Locator, cwd) {
			return fmt.Errorf("run from an existing checkout outside source before removal")
		}
	}
	// Existing paths must be physical; a missing source is only handled by the
	// receipt-aware observer. Even a dangling symlink is never an absent source.
	if _, err := os.Lstat(s.Locator); err == nil {
		if err := contentPhysical(s.Locator); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
