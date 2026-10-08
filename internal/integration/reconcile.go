package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/delivery"
	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type LegacyRelease struct {
	Kind           string                     `json:"kind"`
	SchemaVersion  int                        `json:"schema_version"`
	TaskResult     workspace.TaskResultRecord `json:"task_result"`
	Answer         string                     `json:"answer"`
	ActorClaim     string                     `json:"actor_claim"`
	StartedAtUTC   string                     `json:"started_at_utc"`
	CompletedAtUTC string                     `json:"completed_at_utc"`
	Observation    string                     `json:"observation"`
}

func legacyReleasePath(root string, result workspace.TaskResultRecord) string {
	return filepath.Join(storeRoot(root), "legacy-releases", string(result.ID)+".json")
}
func legacyFacts(d taskrun.Dependencies, root, task string) (workspace.WorkItemRegistry, workspace.TaskRecordedProgressFacts, error) {
	r, e := d.Workspace.WorkItems.Snapshot(root)
	if e != nil {
		return r, workspace.TaskRecordedProgressFacts{}, e
	}
	f := workspace.RecordedTaskProgressFacts(r)[workspace.TaskID(task)]
	if f.Ambiguous || f.TaskResult == nil || f.IntegrationResult == nil || (f.IntegrationResult.Outcome != "exact_effect" && f.IntegrationResult.Outcome != "already_integrated") || f.IntegrationResult.RecoveryStatus != "complete" {
		return r, f, fmt.Errorf("legacy_result_ambiguous: an exact unambiguous native integration is required; neither age nor a missing directory proves integration")
	}
	return r, f, nil
}
func readLegacyRelease(root string, tr workspace.TaskResultRecord) (taskrun.HumanIntegrationOwnership, error) {
	o := taskrun.HumanIntegrationOwnership{TaskResultID: tr.ID, CandidateOID: tr.ResultOID, CandidateTree: tr.ResultTree, Reason: "legacy_owner_release_required", NextAction: "Confirm that all writers have stopped with ply integration release --task " + string(tr.TaskID) + " before closeout."}
	path := legacyReleasePath(root, tr)
	var r LegacyRelease
	e := readJSON(path, &r)
	if os.IsNotExist(e) {
		return o, nil
	}
	if e != nil {
		return o, e
	}
	if r.Kind != "ply.integration.legacy-owner-release" || r.SchemaVersion != 1 || digest(r.TaskResult) != digest(tr) || r.Answer != "released" || strings.TrimSpace(r.ActorClaim) == "" || r.Observation == "" {
		return o, fmt.Errorf("legacy ownership release conflicts with the exact candidate")
	}
	if r.CompletedAtUTC < r.StartedAtUTC {
		return o, fmt.Errorf("legacy release provenance invalid")
	}
	o.Released = true
	o.Release = &taskrun.FileBinding{Locator: path, SHA256: digest(r)}
	o.Reason = "explicit_legacy_owner_release"
	o.NextAction = "Review the exact native integration and closeout plan."
	return o, nil
}

// ReleaseLegacy preserves the operator's explicit release, separately from QA.
// Absence of a native owner is not labelled as proof of process exit. Any known
// active or unresolved native owner blocks the release.
func (s *Service) ReleaseLegacy(cwd, task, actor, answer, started string) (LegacyRelease, error) {
	var out LegacyRelease
	if answer != "released" || strings.TrimSpace(actor) == "" {
		return out, fmt.Errorf("ownership_release_required: exact released answer and actor are required")
	}
	d, root, e := s.context(cwd)
	if e != nil {
		return out, e
	}
	e = withLock(root, func() error {
		if existing, err := activeTaskIntegration(root, workspace.TaskID(task)); err != nil {
			return err
		} else if existing != nil {
			return existingIntegrationError(*existing)
		}
		_, f, err := legacyFacts(d, root, task)
		if err != nil {
			return err
		}
		observation, err := taskrun.CheckLegacyTaskOwnership(d, root, *f.TaskResult)
		if err != nil {
			return err
		}
		b, err := taskrun.Canonical(observation)
		if err != nil {
			return err
		}
		out = LegacyRelease{Kind: "ply.integration.legacy-owner-release", SchemaVersion: 1, TaskResult: *f.TaskResult, Answer: answer, ActorClaim: actor, StartedAtUTC: started, CompletedAtUTC: s.now(), Observation: string(b)}
		path := legacyReleasePath(root, *f.TaskResult)
		var old LegacyRelease
		if err = readJSON(path, &old); err == nil {
			if digest(old.TaskResult) != digest(*f.TaskResult) {
				return fmt.Errorf("legacy release candidate changed")
			}
			out = old
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}
		return atomicJSON(path, out, true)
	})
	return out, e
}

func (s *Service) previewReconcile(d taskrun.Dependencies, root string, in Selection) (Preview, error) {
	out := Preview{Kind: "ply.integration.preview", SchemaVersion: 1, State: "ready", Reasons: []string{}, NextAction: "Review the preserved native integration and confirm only the remaining closeout."}
	if existing, err := activeTaskIntegration(root, workspace.TaskID(in.TaskID)); err != nil {
		return out, err
	} else if existing != nil {
		return out, existingIntegrationError(*existing)
	}
	r, f, e := legacyFacts(d, root, in.TaskID)
	if e != nil {
		return out, e
	}
	tr := *f.TaskResult
	var authority *workspace.IntegrationAuthority
	for i := range r.IntegrationAuthorities {
		if r.IntegrationAuthorities[i].ID == f.IntegrationResult.AuthorityID {
			authority = &r.IntegrationAuthorities[i]
			break
		}
	}
	if authority == nil {
		return out, fmt.Errorf("legacy_result_ambiguous: native authority missing")
	}
	o, e := readLegacyRelease(root, tr)
	if e != nil {
		return out, e
	}
	if _, e = taskrun.CheckLegacyTaskOwnership(d, root, tr); e != nil {
		blocked(&out, "owner_active: "+e.Error(), "Resolve the actual owner before any legacy cleanup.")
	}
	p := Plan{Kind: "ply.integration.plan", SchemaVersion: 1, Workspace: root, Selection: in, TaskResult: tr, Flow: "legacy_reconcile", PreviousQA: f.HumanQA, Ownership: o, ExpectedParentOID: authority.Plan.Epic.ExpectedParentOID}
	// This describes an observed historical target, not a new Delivery agreement
	// or permission to integrate. Reconcile never invokes a merge primitive.
	p.Agreement = workspace.DeliveryAgreement{TargetRef: authority.Plan.Epic.ParentRef, TargetWorktree: authority.Plan.Epic.ParentLocator}
	p.EffectKey = digest(struct{ Common, OID, Target, Integration string }{tr.GitCommonDir, tr.ResultOID, p.Agreement.TargetRef, string(f.IntegrationResult.ID)})
	closePreview, e := workspace.PreviewTaskReconcile(d.Workspace, tr.TaskID, closeoutInput(p, "preview").Ownership, in.Keep)
	if e != nil {
		return out, e
	}
	p.Closeout = closePreview.Plan
	p.CloseoutSHA256 = closePreview.PlanSHA256
	if !o.Released {
		blocked(&out, o.Reason, o.NextAction)
	}
	if !closePreview.Ready {
		for _, reason := range closePreview.Reasons {
			blocked(&out, reason, closePreview.NextAction)
		}
	}
	if in.NotificationRoutePath != "" {
		route, e := LoadNotificationRoute(in.NotificationRoutePath)
		if e != nil {
			return out, e
		}
		p.NotificationRoute = &taskrun.FileBinding{Locator: route.Path, SHA256: route.SHA256}
		p.NotificationChoice = &route
	}
	validateControlLocations(&out, p)
	out.Plan = p
	out.PlanSHA256 = digest(p)
	out.ID = operationID(out.PlanSHA256)
	return out, nil
}

type ScanEntry struct {
	TaskID       workspace.TaskID    `json:"task_id,omitempty"`
	ProjectID    workspace.ProjectID `json:"project_id"`
	Worktree     string              `json:"worktree,omitempty"`
	State        string              `json:"state"`
	CandidateOID string              `json:"candidate_oid,omitempty"`
	Reasons      []string            `json:"reasons"`
	NextAction   string              `json:"next_action"`
}
type ScanResult struct {
	Kind          string      `json:"kind"`
	SchemaVersion int         `json:"schema_version"`
	Entries       []ScanEntry `json:"entries"`
}

func (s *Service) Scan(cwd, project string) (ScanResult, error) {
	out := ScanResult{Kind: "ply.integration.scan", SchemaVersion: 1, Entries: []ScanEntry{}}
	d, root, e := s.context(cwd)
	if e != nil {
		return out, e
	}
	registry, e := d.Workspace.WorkItems.Snapshot(root)
	if e != nil {
		return out, e
	}
	facts := workspace.RecordedTaskProgressFacts(registry)
	integrations, e := integrationReceipts(root)
	if e != nil {
		return out, e
	}
	pending := map[workspace.TaskID]Receipt{}
	for _, operation := range integrations {
		if operation.Decision.Human.Answer == "pass" && operation.State != "completed" && operation.State != "cancelled" {
			pending[operation.Plan.TaskResult.TaskID] = operation
		}
	}
	known := map[string]bool{}
	for _, ep := range registry.Epics {
		for _, repo := range ep.RepoBindings {
			known[repo.Worktree.Locator] = true
		}
	}
	projects, repos, e := d.Workspace.Projects.Snapshot(root)
	if e != nil {
		return out, e
	}
	for _, repo := range repos {
		known[repo.Locator] = true
	}
	for _, task := range registry.Tasks {
		if task.Worktree != nil {
			known[task.Worktree.Locator] = true
		}
		if project != "" && string(task.ProjectID) != project {
			continue
		}
		row := ScanEntry{TaskID: task.ID, ProjectID: task.ProjectID, State: "active", Reasons: []string{}, NextAction: "Continue the registered Task or inspect its native evidence."}
		if task.Worktree != nil {
			row.Worktree = task.Worktree.Locator
		}
		f := facts[task.ID]
		if f.TaskResult != nil {
			row.CandidateOID = f.TaskResult.ResultOID
		}
		closed, err := workspace.ReadTaskCloseoutAt(root, task.ID)
		if err != nil {
			return out, err
		}
		if operation, ok := pending[task.ID]; ok {
			row.State = operation.State
			row.Reasons = append(row.Reasons, operation.Reasons...)
			row.NextAction = "Resume the existing operation with ply integration resume " + operation.ID + "."
		} else if closed != nil {
			if closed.State == "complete" {
				row.State = closed.ResourceState
				row.NextAction = "Inspect retained history; no cleanup remains."
			} else {
				row.State = "cleanup_pending"
				row.NextAction = "Resume integration " + closed.OperationID + "."
			}
		} else if f.Ambiguous {
			row.State = "unknown"
			row.Reasons = []string{"legacy_result_ambiguous"}
			row.NextAction = "Resolve competing native results; no cleanup is authorized."
		} else if f.IntegrationResult != nil && f.IntegrationResult.RecoveryStatus == "complete" {
			row.State = "integrated_not_closed"
			row.NextAction = "Inspect the exact result with integration --task " + string(task.ID) + " --reconcile --check."
		}
		if _, owned := pending[task.ID]; row.Worktree != "" && closed == nil && !owned {
			if _, err = os.Lstat(row.Worktree); os.IsNotExist(err) {
				row.State = "missing_worktree"
				row.Reasons = append(row.Reasons, "absence_is_not_integration_evidence")
				row.NextAction = "Investigate the missing registered resource; do not infer completion."
			} else if err != nil {
				row.State = "unknown"
				row.Reasons = append(row.Reasons, err.Error())
			}
		}
		out.Entries = append(out.Entries, row)
	}
	for _, repo := range repos {
		var projectID workspace.ProjectID
		for _, p := range projects {
			for _, repoID := range p.RepoIDs {
				if repoID == repo.ID && (project == "" || string(p.ID) == project) {
					projectID = p.ID
				}
			}
		}
		if projectID == "" {
			continue
		}
		inventory, err := d.Workspace.WorkGit.ListWorktrees(repo)
		if err != nil {
			return out, err
		}
		for _, entry := range inventory.Entries {
			if !known[entry.Locator] {
				out.Entries = append(out.Entries, ScanEntry{ProjectID: projectID, Worktree: entry.Locator, State: "unknown", Reasons: []string{"no_registered_task_owner"}, NextAction: "Inspect the unowned worktree manually; scan never authorizes removal."})
			}
		}
	}
	return out, nil
}

func (s *Service) ReleaseExited(cwd, id string) (taskrun.WorkflowRun, error) {
	d, root, e := s.context(cwd)
	if e != nil {
		return taskrun.WorkflowRun{}, e
	}
	r, _, e := delivery.NewService(d, delivery.Options{}).ValidateForIntegration(cwd, id)
	if e != nil {
		return taskrun.WorkflowRun{}, e
	}
	return taskrun.ReleaseExitedDeliveryOwnership(d, root, r.Manifest.WorkflowRunID, r.Manifest.TaskResult.ID)
}

func (s *Service) PreviewOwnerRecovery(cwd, id string) (taskrun.HumanOwnerRecoveryPreview, error) {
	d, root, err := s.context(cwd)
	if err != nil {
		return taskrun.HumanOwnerRecoveryPreview{}, err
	}
	r, _, err := delivery.NewService(d, delivery.Options{}).ValidateForIntegration(cwd, id)
	if err != nil {
		return taskrun.HumanOwnerRecoveryPreview{}, err
	}
	return taskrun.PreviewHumanOwnerRecovery(d, root, r.Manifest.WorkflowRunID, r.Manifest.TaskResult.ID)
}

func (s *Service) RecoverOwner(cwd, id string, in taskrun.HumanOwnerRecoveryAnswer) (taskrun.WorkflowRun, error) {
	d, root, err := s.context(cwd)
	if err != nil {
		return taskrun.WorkflowRun{}, err
	}
	r, _, err := delivery.NewService(d, delivery.Options{}).ValidateForIntegration(cwd, id)
	if err != nil {
		return taskrun.WorkflowRun{}, err
	}
	return taskrun.RecoverHumanDeliveryOwnership(d, root, r.Manifest.WorkflowRunID, r.Manifest.TaskResult.ID, in)
}
