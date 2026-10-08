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

type cwdFiles struct {
	workspace.FileSystem
	cwd string
}

func (f cwdFiles) Getwd() (string, error) { return f.cwd, nil }

type workspaceObserver struct{ dependencies workspace.Dependencies }

func (o workspaceObserver) ObserveContaining() (workflowhandoff.WorkspaceSnapshot, error) {
	v, e := workspace.ObserveContaining(o.dependencies)
	if e != nil {
		return workflowhandoff.WorkspaceSnapshot{}, e
	}
	return o.ObserveRoot(v.Root)
}
func (o workspaceObserver) ObserveRoot(root string) (workflowhandoff.WorkspaceSnapshot, error) {
	v, e := workspace.ObserveRoot(o.dependencies, root)
	if e != nil {
		return workflowhandoff.WorkspaceSnapshot{}, e
	}
	p, r, e := o.dependencies.Projects.Snapshot(root)
	return workflowhandoff.WorkspaceSnapshot{Observation: v, Projects: p, Repositories: r}, e
}
func (s *Service) context(cwd string) (taskrun.Dependencies, string, error) {
	d := s.dependencies
	abs, e := filepath.Abs(cwd)
	if e != nil {
		return d, "", e
	}
	physical, e := filepath.EvalSymlinks(abs)
	if e != nil {
		return d, "", e
	}
	d.Workspace.Files = cwdFiles{d.Workspace.Files, physical}
	d.CWD = func() (string, error) { return physical, nil }
	d.Workflow.TaskWorkspace = &d.Workspace
	d.Workflow.Workspace = workspaceObserver{d.Workspace}
	o, e := workspace.ObserveContaining(d.Workspace)
	return d, o.Root, e
}

func validateSelection(in Selection) error {
	if in.Reconcile {
		if in.TaskID == "" || in.DeliveryID != "" || in.Method != "" || in.InstallProfilePath != "" {
			return fmt.Errorf("reconcile requires one --task and no Delivery, merge method or installation")
		}
		_, e := workspace.ParseTaskID(in.TaskID)
		return e
	}
	if in.TaskID != "" {
		return fmt.Errorf("--task requires --reconcile; use an exact Delivery ID for integration")
	}
	if in.Method != "" && in.Method != "merge" && in.Method != "squash" {
		return fmt.Errorf("merge method must be merge or squash")
	}
	for _, p := range []string{in.InstallProfilePath, in.NotificationRoutePath} {
		if p != "" && !filepath.IsAbs(p) {
			return fmt.Errorf("profile and route paths must be absolute")
		}
	}
	return nil
}
func blocked(p *Preview, reason, next string) {
	p.State = "blocked"
	p.Reasons = append(p.Reasons, reason)
	p.NextAction = next
}

func sourceContains(source, path string) bool {
	if source == "" || path == "" {
		return false
	}
	rel, err := filepath.Rel(source, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func validateControlLocations(out *Preview, p Plan) {
	for _, binding := range []*taskrun.FileBinding{p.InstallFile, p.NotificationRoute} {
		if binding != nil && sourceContains(p.TaskResult.SourceLocator, binding.Locator) {
			blocked(out, "control_input_inside_source: "+binding.Locator, "Place the selected profile or route outside the Task source before confirming the plan; recovery must retain these exact control inputs.")
		}
	}
	if p.NotificationChoice != nil && sourceContains(p.TaskResult.SourceLocator, p.NotificationChoice.StateRoot) {
		blocked(out, "control_output_inside_source: notification state", "Choose a route whose durable notification state is outside the Task source.")
	}
	if p.Install != nil && sourceContains(p.TaskResult.SourceLocator, p.Install.ArtifactPath) {
		blocked(out, "control_output_inside_source: installed artifact", "Choose an installed artifact outside the Task source so successful cleanup preserves it.")
	}
}

func integrationReceipts(root string) ([]Receipt, error) {
	files, e := os.ReadDir(filepath.Join(storeRoot(root), "records"))
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	var receipts []Receipt
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".publish-") {
			continue
		}
		r, e := readReceipt(root, strings.TrimSuffix(file.Name(), ".json"))
		if e != nil {
			return nil, e
		}
		receipts = append(receipts, r)
	}
	return receipts, nil
}

// An integrated Task can still have required installation or closeout work.
// Legacy reconciliation must never create a second owner for those effects.
func activeTaskIntegration(root string, task workspace.TaskID) (*Receipt, error) {
	receipts, e := integrationReceipts(root)
	if e != nil {
		return nil, e
	}
	for _, r := range receipts {
		if r.Plan.TaskResult.TaskID == task && r.Decision.Human.Answer == "pass" && r.State != "completed" && r.State != "cancelled" {
			return &r, nil
		}
	}
	return nil, nil
}

func existingIntegrationError(r Receipt) error {
	return fmt.Errorf("integration_in_progress: Task belongs to %s (%s); use ply integration resume %s", r.ID, r.State, r.ID)
}

func (s *Service) Preview(cwd string, in Selection) (Preview, error) {
	out := Preview{Kind: "ply.integration.preview", SchemaVersion: 1, State: "ready", Reasons: []string{}, NextAction: "Review this exact candidate and effect plan, then answer pass, fail or blocked."}
	if e := validateSelection(in); e != nil {
		return out, e
	}
	d, root, e := s.context(cwd)
	if e != nil {
		return out, e
	}
	if in.Reconcile {
		return s.previewReconcile(d, root, in)
	}
	deliveries := delivery.NewService(d, delivery.Options{})
	if in.DeliveryID == "" {
		list, err := deliveries.List(cwd, "", "")
		if err != nil {
			return out, err
		}
		if len(list.Deliveries) != 1 {
			return out, fmt.Errorf("delivery_selection_required: found %d Deliveries; use integration list and select --delivery explicitly", len(list.Deliveries))
		}
		in.DeliveryID = list.Deliveries[0].ID
	}
	r, auth, e := deliveries.ValidateForIntegration(cwd, in.DeliveryID)
	if e != nil {
		return out, fmt.Errorf("candidate_evidence_invalid: %w", e)
	}
	p := Plan{Kind: "ply.integration.plan", SchemaVersion: 1, Workspace: root, Selection: in, DeliveryID: r.ID, DeliverySHA256: r.ManifestSHA256, WorkflowRunID: r.Manifest.WorkflowRunID, TaskResult: r.Manifest.TaskResult, Agreement: r.Manifest.Agreement, ExpectedParentOID: auth.ExpectedParentOID, Flow: "local_return", EffectKey: r.Manifest.EffectKey, PreviousQA: r.HumanQA}
	if existing, err := activeTaskIntegration(root, p.TaskResult.TaskID); err != nil {
		return out, err
	} else if existing != nil {
		return out, existingIntegrationError(*existing)
	}
	p.Ownership, e = taskrun.PreviewHumanIntegration(d, root, p.WorkflowRunID, p.TaskResult.ID)
	if e != nil {
		return out, e
	}
	if !p.Ownership.Released {
		blocked(&out, "owner_active: "+p.Ownership.Reason, p.Ownership.NextAction)
	}
	if p.Agreement.Mode == workspace.DeliveryPullRequest {
		p.Flow = "github_pr_merge"
		if r.PR == nil || r.State != "delivered" || !r.NativeClosed || !r.PR.MetadataApplied {
			blocked(&out, "delivery_pending: this Delivery has no completed, exact PR publication", "Complete the existing PR Delivery before starting integration; this command does not publish a PR.")
		}
		if in.Method == "" {
			blocked(&out, "merge_method_required: select merge or squash explicitly", "Repeat with --method merge or --method squash after reviewing the repository policy.")
		}
		if r.PR != nil {
			if r.PR.HeadOID != p.TaskResult.ResultOID || r.PR.Repository != p.Agreement.GitHubRepository || r.PR.BaseRef != p.Agreement.TargetRef || r.PR.HeadRef != p.Agreement.SourceRef {
				return out, fmt.Errorf("PR receipt differs from frozen Delivery head/base")
			}
			p.PR = &PRMergeTarget{Worktree: p.TaskResult.SourceLocator, Repository: r.PR.Repository, Number: r.PR.Number, HeadRef: r.PR.HeadRef, HeadOID: r.PR.HeadOID, BaseRef: r.PR.BaseRef, Method: in.Method}
			if in.Method != "" {
				o, err := s.options.PR.Observe(*p.PR)
				out.PRObservation = &o
				if err != nil {
					blocked(&out, "pr_observation_failed: "+err.Error(), "Restore PR observation and preview again.")
				} else if o.State != "ready" && o.State != "merged" {
					blocked(&out, o.State+": "+o.Reason, o.NextAction)
				}
			}
		}
	} else {
		if in.Method != "" {
			return out, fmt.Errorf("local return does not take a PR merge method")
		}
		parent, err := d.Workspace.WorkGit.ObserveWorktree(p.Agreement.TargetWorktree)
		if err != nil {
			blocked(&out, "target_changed: "+err.Error(), "Restore the exact registered target and preview again.")
		} else if parent.Ref != p.Agreement.TargetRef || parent.OID != p.ExpectedParentOID || parent.GitCommonDir != p.TaskResult.GitCommonDir || !parent.Clean || !parent.InventoryMatch {
			blocked(&out, "target_changed: the registered return checkout is dirty, moved or has a different identity", "Requalify against the changed parent; this command never rebases or resets it.")
		}
	}
	if in.InstallProfilePath != "" {
		profile, err := LoadInstallProfile(in.InstallProfilePath)
		if err != nil {
			return out, err
		}
		if err = ValidateInstallProfile(profile, p.TaskResult.ResultOID); err != nil {
			return out, err
		}
		b, err := readBounded(in.InstallProfilePath, 1<<20)
		if err != nil {
			return out, err
		}
		p.Install = &profile
		p.InstallFile = &taskrun.FileBinding{Locator: in.InstallProfilePath, SHA256: byteHash(b)}
	}
	if in.NotificationRoutePath != "" {
		route, err := LoadNotificationRoute(in.NotificationRoutePath)
		if err != nil {
			return out, err
		}
		p.NotificationRoute = &taskrun.FileBinding{Locator: route.Path, SHA256: route.SHA256}
		p.NotificationChoice = &route
	}
	validateControlLocations(&out, p)
	closeInput := closeoutInput(p, "preview")
	closePreview, err := workspace.PreviewTaskCloseout(d.Workspace, closeInput)
	if err != nil {
		return out, err
	}
	p.Closeout = closePreview.Plan
	p.CloseoutSHA256 = closePreview.PlanSHA256
	if !closePreview.Ready {
		for _, reason := range closePreview.Reasons {
			blocked(&out, reason, closePreview.NextAction)
		}
	}
	out.Plan = p
	out.PlanSHA256 = digest(p)
	out.ID = operationID(out.PlanSHA256)
	return out, nil
}

func closeoutInput(p Plan, operation string) workspace.TaskCloseoutInput {
	r := p.TaskResult
	o := workspace.TaskCloseoutOwnership{Released: p.Ownership.Released, TaskID: r.TaskID, TaskResultID: r.ID, ResultOID: r.ResultOID, ResultTree: r.ResultTree}
	if p.Ownership.Release != nil {
		o.EvidenceLocator = p.Ownership.Release.Locator
		o.EvidenceSHA256 = p.Ownership.Release.SHA256
	}
	return workspace.TaskCloseoutInput{OperationID: operation, TaskID: r.TaskID, TaskResultID: r.ID, ExpectedResultOID: r.ResultOID, ExpectedResultTree: r.ResultTree, TargetRef: p.Agreement.TargetRef, Keep: p.Selection.Keep, Ownership: o}
}

func (s *Service) Read(cwd, id string) (Receipt, error) {
	_, root, e := s.context(cwd)
	if e != nil {
		return Receipt{}, e
	}
	return readReceipt(root, id)
}
func (s *Service) List(cwd, project, task string) (ListResult, error) {
	out := ListResult{Kind: "ply.integration.list", SchemaVersion: 1, Deliveries: []delivery.Receipt{}, Integrations: []Receipt{}}
	d, root, e := s.context(cwd)
	if e != nil {
		return out, e
	}
	list, e := delivery.NewService(d, delivery.Options{}).List(cwd, task, "")
	if e != nil {
		return out, e
	}
	for _, r := range list.Deliveries {
		if project == "" || string(r.Manifest.ProjectID) == project {
			out.Deliveries = append(out.Deliveries, r)
		}
	}
	files, e := os.ReadDir(filepath.Join(storeRoot(root), "records"))
	if os.IsNotExist(e) {
		return out, nil
	}
	if e != nil {
		return out, e
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".publish-") {
			continue
		}
		r, err := readReceipt(root, strings.TrimSuffix(file.Name(), ".json"))
		if err != nil {
			return out, err
		}
		if (project == "" || string(r.Plan.TaskResult.ProjectID) == project) && (task == "" || string(r.Plan.TaskResult.TaskID) == task) {
			out.Integrations = append(out.Integrations, r)
		}
	}
	return out, nil
}
