package workspace

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SourceWriterObservation reports bounded native observations, never an absolute lease.
type SourceWriterObservation struct {
	State      string `yaml:"state" json:"state"`
	Method     string `yaml:"method" json:"method"`
	Limitation string `yaml:"limitation" json:"limitation"`
}
type TaskAdoptionInput struct {
	TaskID      TaskID           `json:"task_id"`
	Target      QueueTargetInput `json:"target"`
	Path        string           `json:"path"`
	Ref         string           `json:"ref"`
	ExpectedOID string           `json:"expected_oid"`
}
type TaskAdoptionPlan struct {
	Workspace     string                  `yaml:"workspace" json:"workspace"`
	Input         TaskAdoptionInput       `yaml:"input" json:"input"`
	Goal          TaskGoalRef             `yaml:"goal" json:"goal"`
	Target        QueueTarget             `yaml:"target" json:"target"`
	Source        GitWorktreeObservation  `yaml:"source" json:"source"`
	Parent        PlanWorktreeObservation `yaml:"parent" json:"parent"`
	BaseRevision  int                     `yaml:"base_revision" json:"base_revision"`
	QueueID       string                  `yaml:"queue_id" json:"queue_id"`
	QueueRevision int                     `yaml:"queue_revision" json:"queue_revision"`
	Predecessor   string                  `yaml:"predecessor" json:"predecessor"`
}
type TaskAdoptionIntent struct {
	Kind          string                  `yaml:"kind" json:"kind"`
	Plan          TaskAdoptionPlan        `yaml:"plan" json:"plan"`
	Confirmation  string                  `yaml:"confirmation" json:"confirmation"`
	Writers       SourceWriterObservation `yaml:"writers" json:"writers"`
	RecordedAtUTC string                  `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}
type TaskAdoptionRelease struct {
	Kind           string              `yaml:"kind" json:"kind"`
	OperationID    WorktreeOperationID `yaml:"operation_id" json:"operation_id"`
	IntentDigest   string              `yaml:"intent_digest" json:"intent_digest"`
	ExpectedHead   string              `yaml:"expected_head" json:"expected_head"`
	Reason         string              `yaml:"reason" json:"reason"`
	ExecutionFence string              `yaml:"execution_fence" json:"execution_fence"`
	Confirmation   string              `yaml:"confirmation" json:"confirmation"`
	RecordedAtUTC  string              `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}
type TaskAdoptionObservation struct {
	Phase         string `yaml:"phase" json:"phase"`
	Reason        string `yaml:"reason" json:"reason"`
	RecordedAtUTC string `yaml:"recorded_at_utc" json:"recorded_at_utc"`
}
type TaskAdoptionReadback struct {
	Origin        string                   `json:"origin"`
	CurrentSource *GitWorktreeObservation  `json:"current_source"`
	CurrentParent *PlanWorktreeObservation `json:"current_parent"`
	Kind          string                   `json:"kind"`
	SchemaVersion int                      `json:"schema_version"`
	Phase         string                   `json:"phase"`
	Plan          TaskAdoptionPlan         `json:"plan"`
	Confirmation  string                   `json:"confirmation"`
	Operation     *WorktreeOperationRecord `json:"operation"`
	GitEffects    []string                 `json:"git_effects"`
	AgentStarted  bool                     `json:"agent_started"`
	Qualified     bool                     `json:"qualified"`
	MetadataOnly  bool                     `json:"metadata_only"`
	WriterAbsence string                   `json:"writer_absence"`
	Reasons       []string                 `json:"reasons"`
	NextArgv      []string                 `json:"next_argv"`
	CheckCwd      string                   `json:"check_cwd"`
}

func operationStartOID(o WorktreeOperationRecord) string {
	if o.Adoption != nil {
		return o.Adoption.Plan.Source.OID
	}
	return o.ParentOID
}
func operationStartTree(o WorktreeOperationRecord) string {
	if o.Adoption != nil {
		return o.Adoption.Plan.Source.Tree
	}
	return o.ParentTree
}
func operationOrigin(o WorktreeOperationRecord) string {
	if o.Adoption != nil {
		return "adopted"
	}
	return "ply_created"
}
func validateStoredAdoption(r WorkItemRegistry, o WorktreeOperationRecord) error {
	a := o.Adoption
	p := a.Plan
	for _, v := range o.AdoptionObservations {
		if (v.Phase != "blocked" && v.Phase != "unknown") || v.Reason == "" || contentUTC(v.RecordedAtUTC) != nil {
			return fmt.Errorf("invalid adoption observation")
		}
	}
	t, _ := findTask(r, o.TaskID)
	if t == nil || t.ProjectID != o.ProjectID || t.RepoID != o.RepoID || t.ParentEpicID != o.ParentEpicID || t.GitCommonDir != o.GitCommonDir || p.Target.ProjectID != o.ProjectID || p.Target.RepoID != o.RepoID || p.Target.EpicID != o.ParentEpicID || p.Target.ParentWorktreeID != o.ParentWorktreeID || p.Target.ParentRef != o.ParentRef || p.Target.GitCommonDir != o.GitCommonDir || !historicalEpicBaseMatches(r, o.ParentEpicID, o.RepoID, o.ParentWorktreeID, o.ParentRef, o.ParentOID, o.ParentTree) || p.QueueID != queueID(p.Workspace, p.Target) || contentPath(p.Workspace) != nil || validateRegisteredGoalRef(r, o.TaskID, &p.Goal) != nil || validateQueuePlanObservation(p.Parent, p.Target, o.ParentOID, o.ParentTree) != nil {
		return fmt.Errorf("adoption scope or chosen basis differs")
	}
	if a.Kind != "WorkspaceTaskAdoptionIntent@1" || a.Confirmation != queueDigest(p) || p.Input.TaskID != o.TaskID || p.Input.Path != o.TargetLocator || p.Input.Ref != o.SourceRef || p.Source.OID != p.Input.ExpectedOID || p.Source.Locator != o.TargetLocator || p.Source.Ref != o.SourceRef || p.Source.GitCommonDir != o.GitCommonDir || !p.Source.Clean || !p.Source.InventoryMatch || p.Parent.OID != o.ParentOID || p.Parent.Tree != o.ParentTree || p.Source.OID == o.ParentOID || p.Source.Tree == o.ParentTree || a.Writers.State != "absent" || contentUTC(a.RecordedAtUTC) != nil {
		return fmt.Errorf("invalid adopted operation %s", o.ID)
	}
	if o.Execution != nil {
		p := o.Execution
		if o.Release != nil || p.SourceMode != "adopted" || p.AdoptionID != o.ID || p.Confirmation != goalExecutionDigest(*p) || p.Goal.Goal != o.Adoption.Plan.Goal || p.StartOID != operationStartOID(o) || p.StartTree != operationStartTree(o) {
			return fmt.Errorf("invalid adopted execution fence")
		}
	}
	if o.Release != nil {
		v := o.Release
		if v.Kind != "WorkspaceTaskAdoptionRelease@1" || v.OperationID != o.ID || v.IntentDigest != o.IntentDigest || !digestPattern.MatchString(v.ExpectedHead) || v.ExecutionFence != "unstarted" || contentText(2000)(v.Reason) != nil || contentUTC(v.RecordedAtUTC) != nil || v.Confirmation != adoptionReleaseDigest(*v) {
			return fmt.Errorf("invalid adoption release %s", o.ID)
		}
		for _, p := range r.TaskPreparations {
			if p.OperationID == o.ID {
				return fmt.Errorf("released adoption has preparation")
			}
		}
	}
	return nil
}
func AdoptTaskWorktree(d Dependencies, in TaskAdoptionInput, apply bool, confirmation string) (out TaskAdoptionReadback, err error) {
	if _, e := ParseTaskID(string(in.TaskID)); e != nil {
		return out, e
	}
	if !validFullBranchRef(in.Ref) || !validOIDText(in.ExpectedOID) {
		return out, WorkInvalidArguments("full local source ref and expected OID required")
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return out, e
	}
	d, bind := beginGoalOperationContent(d)
	err = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			if e = bind(root, r, in.TaskID); e != nil {
				return e
			}
			// Recovery is by preserved confirmation and input before any fresh observation.
			var old *WorktreeOperationRecord
			for i := range r.WorktreeOperations {
				o := &r.WorktreeOperations[i]
				if o.Adoption != nil && o.Adoption.Confirmation == confirmation && contentTypedEqual(o.Adoption.Plan.Input, in) {
					old = o
					break
				}
			}
			if apply && old != nil && (old.Release != nil || old.State == "ready") {
				out = adoptionReadback(old.Adoption.Plan, old)
				if old.Release == nil {
					if current, e := d.WorkGit.ObserveWorktree(old.TargetLocator); e == nil {
						out.CurrentSource = &current
					} else {
						out.Reasons = append(out.Reasons, "source_observation_unknown")
					}
					if current, e := queueParentObservation(d, old.Adoption.Plan.Target); e == nil {
						out.CurrentParent = &current
					} else {
						out.Reasons = append(out.Reasons, "parent_observation_unknown")
					}
				}
				return nil
			}
			own := WorktreeOperationID("")
			if old != nil {
				own = old.ID
			}
			p, e := observeAdoption(d, root, r, projects, in, own)
			if e != nil {
				if apply && old != nil {
					out, e = recordAdoptionObservation(d, session, r, *old, e)
				}
				return e
			}
			out = adoptionReadback(p, nil)
			if !apply {
				return nil
			}
			if confirmation != queueDigest(p) {
				e := queueError("adoption_preview_stale", "Goal, source, base, queue or adoption generation changed; preview again")
				if old != nil {
					out, e = recordAdoptionObservation(d, session, r, *old, e)
				}
				return e
			}
			if old == nil {
				wt, e := d.WorkIDs.NewWorktreeID()
				if e != nil {
					return e
				}
				id, e := d.WorkIDs.NewOperationID()
				if e != nil {
					return e
				}
				writers, e := adoptionWriters(d, in.Path)
				if e != nil {
					return e
				}
				op := WorktreeOperationRecord{ID: id, TaskID: in.TaskID, WorktreeID: wt, State: "creating", ProjectID: p.Target.ProjectID, RepoID: p.Target.RepoID, GitCommonDir: p.Target.GitCommonDir, ParentEpicID: p.Target.EpicID, ParentWorktreeID: p.Target.ParentWorktreeID, ParentRef: p.Target.ParentRef, ParentOID: p.Parent.OID, ParentTree: p.Parent.Tree, SourceRef: in.Ref, TargetLocator: in.Path, LastObservation: emptyCreateObservation("no_effect"), Adoption: &TaskAdoptionIntent{Kind: "WorkspaceTaskAdoptionIntent@1", Plan: p, Confirmation: confirmation, Writers: writers, RecordedAtUTC: queueNow(d)}}
				op.IntentDigest, e = intentDigest(op)
				if e != nil {
					return e
				}
				t, ti := findTask(r, in.TaskID)
				if t == nil {
					return fmt.Errorf("Task missing")
				}
				r.Tasks[ti].WorktreeState = WorkItemCreating
				r.WorktreeOperations = append(r.WorktreeOperations, op)
				sortWorkRegistry(&r)
				if e = session.Publish(r); e != nil {
					out.Phase = "unknown"
					out.Reasons = append(out.Reasons, e.Error())
					return e
				}
				r, e = session.Snapshot()
				if e != nil {
					return e
				}
				old, _ = adoptionOperation(r, id)
			}
			out = adoptionReadback(p, old)
			fresh, e := observeAdoption(d, root, r, projects, in, old.ID)
			if e != nil {
				out, e = recordAdoptionObservation(d, session, r, *old, e)
				return e
			}
			if !contentTypedEqual(p, fresh) {
				return queueError("adoption_preview_stale", "reserved adoption differs; inspect or release this generation")
			}
			_, repo, e := projectAndRepo(projects, old.ProjectID, old.RepoID)
			if e != nil {
				return e
			}
			ep, b := queueTargetBinding(r, p.Target)
			current := currentEpicBinding(r, ep, b)
			kind, observation := observeTaskCreate(d, repo, current, *old)
			if kind != "exact_effect" {
				out.Phase = "unknown"
				return queueError("adoption_source_unknown", "reserved source observation differs; inspect same operation or release before execution")
			}
			_, ti := findTask(r, in.TaskID)
			_, oi := adoptionOperation(r, old.ID)
			if e = finalizeReady(session, &r, ti, oi, ep, *old, observation); e != nil {
				out.Phase = "unknown"
				out.Reasons = append(out.Reasons, e.Error())
				return e
			}
			o, _ := adoptionOperation(r, old.ID)
			out = adoptionReadback(p, o)
			return nil
		})
	})
	return out, err
}

func adoptionOperation(r WorkItemRegistry, id WorktreeOperationID) (*WorktreeOperationRecord, int) {
	for i := range r.WorktreeOperations {
		if r.WorktreeOperations[i].ID == id {
			return &r.WorktreeOperations[i], i
		}
	}
	return nil, -1
}
func adoptionPredecessor(r WorkItemRegistry, task TaskID) string {
	consumed := map[string]bool{}
	for _, o := range r.WorktreeOperations {
		if o.TaskID == task && o.Adoption != nil && o.Release != nil {
			consumed[o.Adoption.Plan.Predecessor] = true
		}
	}
	for _, o := range r.WorktreeOperations {
		if o.TaskID == task && o.Release != nil {
			h := queueDigest(o.Release)
			if !consumed[h] {
				return h
			}
		}
	}
	return ""
}
func adoptionReleased(r WorkItemRegistry, task TaskID) bool {
	op, _ := findOperation(r, task)
	return op == nil && adoptionPredecessor(r, task) != ""
}
func adoptionReadback(p TaskAdoptionPlan, o *WorktreeOperationRecord) TaskAdoptionReadback {
	out := TaskAdoptionReadback{Origin: "adopted", Kind: "ply.workspace.task-worktree-adoption", SchemaVersion: 1, Phase: "preview", Plan: p, Confirmation: queueDigest(p), Operation: o, GitEffects: []string{}, Reasons: []string{}, NextArgv: []string{"ply", "workflow", "execute", "--spec", p.Goal.SpecID, "--check"}, CheckCwd: p.Target.ParentLocator}
	out.NextArgv = adoptionApplyArgv(p, out.Confirmation)
	if o != nil {
		out.NextArgv = adoptionApplyArgv(p, o.Adoption.Confirmation)
		out.Confirmation = o.Adoption.Confirmation
		out.Phase = "reserved"
		if len(o.AdoptionObservations) > 0 {
			out.Phase = o.AdoptionObservations[len(o.AdoptionObservations)-1].Phase
		}
		if o.State == "ready" {
			out.Phase = "adopted"
			out.NextArgv = []string{"ply", "workflow", "execute", "--spec", p.Goal.SpecID, "--check"}
		}
		if o.Release != nil {
			out.Phase = "released"
			out.MetadataOnly = true
			out.WriterAbsence = "not_required"
			out.NextArgv = []string{"ply", "workspace", "task", "worktree", "adopt", string(o.TaskID), "--project", string(o.ProjectID), "--repo", string(o.RepoID), "--epic", string(o.ParentEpicID), "--path", o.TargetLocator, "--ref", o.SourceRef, "--expected-oid", operationStartOID(*o)}
		}
	}
	return out
}
func adoptionWriters(d Dependencies, path string) (SourceWriterObservation, error) {
	if d.SourceWriters == nil {
		return SourceWriterObservation{State: "unknown"}, queueError("adoption_writer_unknown", "bounded source writer observer unavailable; restore supported inspection and recheck")
	}
	o, e := d.SourceWriters(path)
	if e != nil && o.State == "present" {
		return o, queueError("adoption_writer_present", e.Error())
	}
	if e != nil {
		return o, queueError("adoption_writer_unknown", "source inspection unavailable: "+e.Error())
	}
	if o.State != "absent" || o.Method == "" || o.Limitation == "" {
		return o, queueError("adoption_writer_"+o.State, "stop actual source writers or restore supported inspection, then recheck")
	}
	return o, nil
}
func observeAdoption(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, in TaskAdoptionInput, own WorktreeOperationID) (TaskAdoptionPlan, error) {
	var p TaskAdoptionPlan
	if r.FormatVersion < 7 {
		return p, registryDiagnostic("registry_upgrade_required", "Task adoption requires an explicit delivered-reader upgrade to format 7", r.FormatVersion)
	}
	t, _ := findTask(r, in.TaskID)
	if t == nil {
		return p, queueError("adoption_task_missing", "register the Task and its first Goal before adoption")
	}
	target, repo, e := resolveQueueTarget(d, root, r, projects, in.Target)
	if e != nil {
		return p, e
	}
	cwd, e := projectPhysicalWorkingDirectory(d.Files)
	if e != nil {
		return p, e
	}
	if cwd != target.ParentLocator {
		return p, queueError("adoption_return_checkout_required", "run adoption from the registered Epic return checkout: "+target.ParentLocator)
	}
	observedRepo, e := d.WorkGit.ObserveRepo(repo)
	if e != nil {
		return p, e
	}
	if observedRepo.GitCommonDir != repo.GitCommonDir {
		return p, queueError("adoption_repository_changed", "registered repository common directory differs")
	}
	if t.ProjectID != target.ProjectID || t.RepoID != target.RepoID || t.ParentEpicID != target.EpicID {
		return p, queueError("adoption_scope_conflict", "Task differs from selected target")
	}
	life, e := readWorkItemLifecycle(root)
	if e != nil {
		return p, e
	}
	if life.Task(t.ID) != LifecycleActive || life.Epic(t.ParentEpicID) != LifecycleActive {
		return p, queueError("adoption_inactive", "Task and Epic must be active")
	}
	if prep := preparationForTask(r, t.ID); prep != nil {
		return p, queueError("adoption_execution_exists", "Task already has preparation "+prep.ID)
	}
	for _, result := range r.TaskResults {
		if result.TaskID == t.ID {
			return p, queueError("adoption_history_conflict", "Task has result history")
		}
	}
	if op, _ := findOperation(r, t.ID); op != nil && op.ID != own {
		return p, queueError("adoption_task_bound", "inspect existing worktree operation "+string(op.ID))
	}
	if own == "" && (t.Worktree != nil || t.WorktreeState != WorkItemUnbound) {
		return p, queueError("adoption_task_bound", "Task is not unused")
	}
	qid := queueID(root, target)
	q := foldQueue(r, qid)
	if len(q.Active) != 0 {
		return p, queueError("adoption_effect_pending", "inspect current queue preparation")
	}
	for _, id := range queueUnresolved(r, target) {
		if id != string(own) {
			return p, queueError("adoption_effect_pending", "inspect unresolved target effect "+id)
		}
	}
	for _, ev := range r.TaskQueueEvents {
		if ev.Kind == "release" {
			for _, pr := range r.TaskPreparations {
				if pr.Plan.TaskID == t.ID {
					return p, queueError("adoption_history_conflict", "Task has queue release history")
				}
			}
		}
	}
	for _, v := range r.TaskSolutionSelections {
		if v.TaskID == t.ID {
			return p, queueError("adoption_history_conflict", "Task has solution selection history")
		}
	}
	for _, v := range r.TaskSpecRevisions {
		if v.TaskID == t.ID {
			m, e := readRegisteredTaskManifest(d.TaskContent, root, r, v.ManifestSHA256)
			if e != nil {
				return p, e
			}
			if !isTaskGoalSpec(m) {
				return p, queueError("adoption_history_conflict", "Task has legacy or execution Spec history")
			}
		}
	}
	metadata, e := readProjectMetadata(root)
	if e != nil {
		return p, e
	}
	for _, rows := range metadata.companions {
		for _, c := range rows {
			if pathsOverlap(in.Path, c.Locator) {
				return p, queueError("adoption_path_conflict", "source overlaps registered companion")
			}
		}
	}
	var goal *TaskGoalRef
	for _, entry := range q.Pending {
		if entry.TaskID == t.ID {
			goal = entry.Goal
		}
	}
	if goal == nil {
		return p, queueError("adoption_goal_required", "register an exact pending first Goal before adoption")
	}
	g, e := loadTaskGoal(d, root, r, t.ID, *goal)
	if e != nil {
		return p, e
	}
	if !contentTypedEqual(taskContentState(r, t.ID).SpecHead, &goal.Spec) || !contentTypedEqual(taskContentState(r, t.ID).ProblemHead, &g.Problem) {
		return p, queueError("adoption_goal_stale", "Goal or Problem changed")
	}
	if g.Delivery == nil || g.Delivery.Mode != DeliveryLocalEpic || !g.Delivery.AutomaticAcceptance() || g.Delivery.Acceptance.RequireHumanQA {
		return p, queueError("adoption_delivery_unsupported", "first cut requires automatic local Epic delivery")
	}
	if _, e = BindDeliveryAgreement(*g.Delivery, target, in.Ref); e != nil {
		return p, e
	}
	if e = contentPath(in.Path); e != nil {
		return p, WorkInvalidArguments("source must be an absolute physical path")
	}
	physical, e := filepath.EvalSymlinks(in.Path)
	if e != nil || physical != in.Path {
		return p, queueError("adoption_path_conflict", "source path must exist without symlink aliases")
	}
	if pathsOverlap(in.Path, repo.Locator) || pathsOverlap(in.Path, target.ParentLocator) {
		return p, queueError("adoption_path_conflict", "source overlaps registered repository or return worktree")
	}
	for _, ep := range r.Epics {
		for _, b := range ep.RepoBindings {
			if pathsOverlap(in.Path, b.Worktree.Locator) || b.GitCommonDir == t.GitCommonDir && b.Worktree.Ref == in.Ref {
				return p, queueError("adoption_resource_conflict", "source is an Epic target")
			}
		}
	}
	for _, o := range r.WorktreeOperations {
		if o.ID != own && o.Release == nil && (pathsOverlap(in.Path, o.TargetLocator) || o.GitCommonDir == t.GitCommonDir && o.SourceRef == in.Ref) {
			return p, queueError("adoption_resource_conflict", "source claimed by "+string(o.ID))
		}
	}
	if e = epicCreateReservation(r, t.GitCommonDir, in.Ref, in.Path, ""); e != nil {
		return p, e
	}
	source, e := observeAdoptionSource(d, in.Path, in.Ref)
	if e != nil {
		return p, e
	}
	ref, e := d.WorkGit.ObserveRef(repo, in.Ref)
	if e != nil {
		return p, e
	}
	if !source.Clean || !source.InventoryMatch || source.GitCommonDir != t.GitCommonDir || source.Locator != in.Path || source.Ref != in.Ref || source.OID != in.ExpectedOID || !ref.Exists || ref.CheckedOutAt != in.Path || ref.OID != source.OID || ref.Tree != source.Tree {
		return p, queueError("adoption_source_changed", "exact clean source/ref/OID/tree/inventory required")
	}
	inventory, e := d.WorkGit.ListWorktrees(repo)
	if e != nil {
		return p, e
	}
	count := 0
	for _, w := range inventory.Entries {
		if w.Ref == in.Ref || w.Locator == in.Path {
			count++
			if w.Ref != in.Ref || w.Locator != in.Path || w.OID != in.ExpectedOID || w.Bare || w.Detached || w.Locked || w.Prunable {
				return p, queueError("adoption_inventory_conflict", "source inventory differs")
			}
		}
	}
	if count != 1 {
		return p, queueError("adoption_inventory_conflict", "source must have one linked checkout")
	}
	parent, e := queueParentObservation(d, target)
	if e != nil {
		return p, e
	}
	ep, b := queueTargetBinding(r, target)
	base := currentEpicBase(r, ep, b)
	if pendingBaseUpdate(r, t.ParentEpicID, t.RepoID) || parent.OID != base.OID || parent.Tree != base.Tree {
		return p, queueError("adoption_base_changed", "refresh the Epic base explicitly before adoption")
	}
	ancestor, e := d.IntegrationGit.CheckAncestor(repo, base.OID, source.OID)
	if e != nil {
		return p, e
	}
	if !ancestor {
		return p, queueError("adoption_base_not_ancestor", "chosen Epic parent is not an ancestor of source")
	}
	if source.OID == base.OID || source.Tree == base.Tree {
		return p, queueError("adoption_no_useful_delta", "source must contain a distinct committed delta")
	}
	if _, e = adoptionWriters(d, in.Path); e != nil {
		return p, e
	}
	return TaskAdoptionPlan{Workspace: root, Input: in, Goal: *goal, Target: target, Source: source, Parent: parent, BaseRevision: base.Revision, QueueID: qid, QueueRevision: q.Revision, Predecessor: adoptionPredecessor(r, t.ID)}, nil
}
func pathsOverlap(a, b string) bool {
	return a == b || strings.HasPrefix(a, b+string(filepath.Separator)) || strings.HasPrefix(b, a+string(filepath.Separator))
}

func adoptionReleaseDigest(r TaskAdoptionRelease) string {
	r.Confirmation = ""
	r.RecordedAtUTC = ""
	return queueDigest(r)
}

// ReleaseTaskAdoption changes native metadata only. No Git or writer observation
// participates: the registry lock is also execution's materialization fence.
func ReleaseTaskAdoption(d Dependencies, task TaskID, id WorktreeOperationID, reason string, apply bool, confirmation string) (out TaskAdoptionReadback, err error) {
	if contentText(2000)(reason) != nil || !operationIDPattern.MatchString(string(id)) {
		return out, WorkInvalidArguments("exact adoption ID and release reason required")
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return out, e
	}
	d = adoptionMetadataDependencies(d, task)
	err = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			o, oi := adoptionOperation(r, id)
			if o == nil || o.TaskID != task || o.Adoption == nil || o.Adoption.Plan.Workspace != root {
				return queueError("adoption_release_identity_conflict", "operation is not this Task's adoption")
			}
			out = adoptionReadback(o.Adoption.Plan, o)
			out.MetadataOnly = true
			out.WriterAbsence = "not_required"
			if o.Release != nil {
				if !apply || confirmation == o.Release.Confirmation && reason == o.Release.Reason {
					out.Confirmation = o.Release.Confirmation
					return nil
				}
				return queueError("adoption_release_conflict", "different release already recorded")
			}
			if o.Execution != nil {
				return queueError("adoption_execution_exists", "execution intent exists for this adoption; recover its preserved native execution request")
			}
			current, _ := findOperation(r, task)
			if current == nil || current.ID != id {
				return queueError("adoption_release_replaced", "adoption no longer owns this Task")
			}
			// Execution materialization must publish this generation's native
			// fence before any content, preparation, Run or start can exist.
			for _, v := range r.TaskSolutionSelections {
				if v.TaskID == task {
					return queueError("adoption_execution_exists", "native solution selection exists; inspect execution metadata")
				}
			}
			for _, p := range r.TaskPreparations {
				if p.Plan.TaskID == task {
					return queueError("adoption_execution_exists", "preparation exists: "+p.ID)
				}
			}
			for _, v := range r.TaskResults {
				if v.TaskID == task {
					return queueError("adoption_execution_exists", "Task result exists")
				}
			}
			t, ti := findTask(r, task)
			if t == nil || t.Worktree != nil && t.Worktree.OperationID != id {
				return queueError("adoption_release_replaced", "current Task binding differs")
			}
			_, repo, e := projectAndRepo(projects, o.ProjectID, o.RepoID)
			if e != nil {
				return e
			}
			if repo.GitCommonDir != o.GitCommonDir {
				return queueError("adoption_release_scope_changed", "registered repository identity differs")
			}
			cwd, e := projectPhysicalWorkingDirectory(d.Files)
			if e != nil {
				return e
			}
			if cwd != o.Adoption.Plan.Target.ParentLocator {
				return queueError("adoption_return_checkout_required", "run release from the registered Epic return checkout")
			}
			release := TaskAdoptionRelease{Kind: "WorkspaceTaskAdoptionRelease@1", OperationID: id, IntentDigest: o.IntentDigest, ExpectedHead: r.RawSHA256, Reason: reason, ExecutionFence: "unstarted"}
			release.Confirmation = adoptionReleaseDigest(release)
			out.Confirmation = release.Confirmation
			out.NextArgv = []string{"ply", "workspace", "task", "worktree", "release-adoption", string(id), "--reason", reason, "--apply", "--confirmation", release.Confirmation}
			if !apply {
				out.Phase = "preview"
				return nil
			}
			if confirmation != release.Confirmation {
				return queueError("adoption_release_stale", "metadata head or release request changed; inspect and preview again")
			}
			release.RecordedAtUTC = queueNow(d)
			r.WorktreeOperations[oi].Release = &release
			r.Tasks[ti].Worktree = nil
			r.Tasks[ti].WorktreeState = WorkItemUnbound
			if e = session.Publish(r); e != nil {
				out.Phase = "unknown"
				out.Reasons = append(out.Reasons, e.Error())
				return e
			}
			out = adoptionReadback(o.Adoption.Plan, &r.WorktreeOperations[oi])
			out.Confirmation = release.Confirmation
			return nil
		})
	})
	return out, err
}

// ReleaseTaskWorktreeAdoption resolves Task identity from the native operation;
// the locked mutation rechecks the same immutable operation/owner pair.
func ReleaseTaskWorktreeAdoption(d Dependencies, id WorktreeOperationID, reason string, apply bool, confirmation string) (TaskAdoptionReadback, error) {
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return TaskAdoptionReadback{}, e
	}
	d = adoptionMetadataDependencies(d, "")
	r, e := d.WorkItems.Snapshot(root)
	if e != nil {
		return TaskAdoptionReadback{}, e
	}
	o, _ := adoptionOperation(r, id)
	if o == nil || o.Adoption == nil {
		return TaskAdoptionReadback{}, queueError("adoption_not_found", "inspect the exact native adoption operation")
	}
	return ReleaseTaskAdoption(d, o.TaskID, id, reason, apply, confirmation)
}

// This narrow metadata view retains strict structural/transition validation and
// the same CAS lock. It omits external Goal/document reads only for release.
// Publication is additionally restricted to one exact unstarted release fold.
func adoptionMetadataDependencies(d Dependencies, task TaskID) Dependencies {
	var narrow func(WorkItemStore) WorkItemStore
	narrow = func(store WorkItemStore) WorkItemStore {
		switch s := store.(type) {
		case *systemWorkItemStore:
			copy := *s
			copy.adoptionMetadataOnly = true
			copy.adoptionReleaseTask = task
			return &copy
		case guardedTaskStartStore:
			s.WorkItemStore = narrow(s.WorkItemStore)
			return s
		default:
			return store
		}
	}
	d.WorkItems = narrow(d.WorkItems)
	return d
}
func validateAdoptionMetadataRelease(old, next WorkItemRegistry, task TaskID) error {
	if task == "" {
		return fmt.Errorf("metadata lookup cannot publish")
	}
	o, oi := findOperation(old, task)
	t, ti := findTask(old, task)
	if o == nil || t == nil || o.Adoption == nil || o.Execution != nil || o.Release != nil || oi >= len(next.WorktreeOperations) || ti >= len(next.Tasks) {
		return fmt.Errorf("metadata release lost active unstarted generation")
	}
	n := next.WorktreeOperations[oi]
	if n.ID != o.ID || n.Release == nil || n.Release.ExpectedHead != old.RawSHA256 || n.Release.OperationID != o.ID || next.Tasks[ti].Worktree != nil || next.Tasks[ti].WorktreeState != WorkItemUnbound {
		return fmt.Errorf("metadata release differs from exact ownership fold")
	}
	for _, p := range old.TaskPreparations {
		if p.Plan.TaskID == task {
			return fmt.Errorf("metadata release has preparation")
		}
	}
	for _, p := range old.TaskResults {
		if p.TaskID == task {
			return fmt.Errorf("metadata release has result")
		}
	}
	copy := cloneWorkItemRegistry(next)
	copy.WorktreeOperations[oi] = *o
	copy.Tasks[ti] = *t
	copy.RawSHA256 = old.RawSHA256
	if !contentTypedEqual(copy, old) {
		return fmt.Errorf("metadata release attempted changes outside its binding and release receipt")
	}
	return nil
}

func adoptionApplyArgv(p TaskAdoptionPlan, confirmation string) []string {
	in := p.Input
	return []string{"ply", "workspace", "task", "worktree", "adopt", string(in.TaskID), "--project", string(p.Target.ProjectID), "--repo", string(p.Target.RepoID), "--epic", string(p.Target.EpicID), "--path", in.Path, "--ref", in.Ref, "--expected-oid", in.ExpectedOID, "--apply", "--confirmation", confirmation}
}
func (s *TaskSpecSession) ValidateTaskStartWriters(id TaskID) error {
	o, _ := findOperation(s.Registry, id)
	if o == nil || o.Adoption == nil {
		return nil
	}
	if o.Release != nil || o.Execution == nil {
		return queueError("adoption_generation_stale", "start requires this adoption's durable execution intent")
	}
	_, e := adoptionWriters(s.Dependencies, o.TargetLocator)
	return e
}

func observeAdoptionSource(d Dependencies, path, ref string) (GitWorktreeObservation, error) {
	source, e := d.WorkGit.ObserveWorktree(path)
	if e != nil {
		return source, queueError("adoption_source_observation_unknown", e.Error())
	}
	exact, e := d.IntegrationGit.ObserveIntegrationWorktree(path, ref)
	if e != nil {
		return source, queueError("adoption_source_observation_unknown", e.Error())
	}
	if exact.Locator != source.Locator || exact.Ref != source.Ref || exact.OID != source.OID || exact.Tree != source.Tree || exact.GitCommonDir != source.GitCommonDir || !exact.Symbolic || !exact.Clean || len(exact.StatusEntries) != 0 || len(exact.InProgress) != 0 {
		return source, queueError("adoption_source_changed", "source identity/status or in-progress Git operation differs")
	}
	return source, nil
}

func validateAdoptionChains(r WorkItemRegistry) error {
	releases := map[string]WorktreeOperationRecord{}
	for _, o := range r.WorktreeOperations {
		if o.Release != nil {
			releases[queueDigest(o.Release)] = o
		}
	}
	successors := map[string]bool{}
	for _, o := range r.WorktreeOperations {
		if o.Adoption == nil {
			continue
		}
		key := string(o.TaskID) + "/" + o.Adoption.Plan.Predecessor
		if successors[key] {
			return fmt.Errorf("adoption history forks for Task %s", o.TaskID)
		}
		successors[key] = true
		seen := map[string]bool{}
		parent := o.Adoption.Plan.Predecessor
		for parent != "" {
			if seen[parent] {
				return fmt.Errorf("adoption predecessor cycle")
			}
			seen[parent] = true
			p, ok := releases[parent]
			if !ok || p.TaskID != o.TaskID || p.Adoption == nil {
				return fmt.Errorf("adoption predecessor is not this Task's proven release")
			}
			parent = p.Adoption.Plan.Predecessor
		}
	}
	return nil
}

func recordAdoptionObservation(d Dependencies, session WorkItemStoreSession, r WorkItemRegistry, o WorktreeOperationRecord, cause error) (TaskAdoptionReadback, error) {
	out := adoptionReadback(o.Adoption.Plan, &o)
	out.Phase = "blocked"
	if strings.Contains(cause.Error(), "unknown") {
		out.Phase = "unknown"
	}
	out.Reasons = append(out.Reasons, cause.Error())
	_, i := adoptionOperation(r, o.ID)
	if i < 0 {
		return out, fmt.Errorf("adoption observation lost its operation: %w", cause)
	}
	r.WorktreeOperations[i].AdoptionObservations = append(r.WorktreeOperations[i].AdoptionObservations, TaskAdoptionObservation{Phase: out.Phase, Reason: cause.Error(), RecordedAtUTC: queueNow(d)})
	out.Operation = &r.WorktreeOperations[i]
	if e := session.Publish(r); e != nil {
		out.Phase = "unknown"
		return out, fmt.Errorf("%w; observation publication unresolved: %v", cause, e)
	}
	return out, cause
}
