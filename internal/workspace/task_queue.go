package workspace

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func queueReasons(code, message string) []QueueReason { return []QueueReason{{code, message}} }
func queueErrorReasons(e error) []QueueReason {
	var ce *TaskContentError
	if errors.As(e, &ce) {
		return queueReasons(ce.Code, ce.Detail)
	}
	return queueReasons("task_queue_observation_unknown", e.Error())
}
func queueTaskTitle(d Dependencies, root string, r WorkItemRegistry, t TaskRecord) (string, error) {
	p := taskContentState(r, t.ID).ProblemHead
	if p == nil {
		return t.Title, nil
	}
	m, e := readRegisteredTaskManifest(d.TaskContent, root, r, p.ManifestSHA256)
	if e != nil {
		return "", e
	}
	return contentString(contentFields(m), "title"), nil
}
func queueTaskNames(qid string, id TaskID, title, parent string) (string, string) {
	title = strings.ToLower(title)
	slug := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(title, "-")
	slug = strings.Trim(slug, "-")
	if len(slug) > 40 {
		slug = strings.Trim(slug[:40], "-")
	}
	if slug == "" {
		slug = "task"
	}
	h := sha256.Sum256([]byte(qid + "\n" + string(id)))
	leaf := slug + "-" + hex.EncodeToString(h[:])[:12]
	return "ply-task/" + leaf, filepath.Join(filepath.Dir(parent), "ply-task-"+leaf)
}
func queueEntryEvaluation(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, tgt QueueTarget, qid string, entry QueueEntry, rank int) (QueuePending, *WorkspaceTaskPreparePlan) {
	return queueMixedEntryEvaluation(d, root, r, projects, tgt, qid, TaskGoalQueueEntry{TaskID: entry.TaskID, Selection: entry.Selection}, rank)
}
func queueMixedEntryEvaluation(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, tgt QueueTarget, qid string, entry TaskGoalQueueEntry, rank int) (QueuePending, *WorkspaceTaskPreparePlan) {
	if entry.Goal != nil {
		return queueGoalEvaluation(d, root, r, entry, rank), nil
	}
	t, _ := findTask(r, entry.TaskID)
	row := QueuePending{Rank: rank, TaskID: entry.TaskID, Selection: entry.Selection, State: "blocked", Reasons: []QueueReason{}}
	block := func(e error) (QueuePending, *WorkspaceTaskPreparePlan) {
		row.Reasons = queueErrorReasons(e)
		if strings.Contains(row.Reasons[0].Code, "unknown") {
			row.State = "unknown"
		}
		return row, nil
	}
	if t == nil {
		return block(queueError("task_queue_task_missing", "queued Task is missing"))
	}
	title, e := queueTaskTitle(d, root, r, *t)
	if e != nil {
		return block(e)
	}
	row.Title = title
	if entry.Selection == nil {
		return block(queueError("task_queue_selection_required", "explicitly bind a selected solution in queue set"))
	}
	sm, e := readRegisteredTaskManifest(d.TaskContent, root, r, entry.Selection.ManifestSHA256)
	if e != nil {
		return block(e)
	}
	sol := contentFields(contentFields(sm)["solution"])
	sid := contentString(sol, "spec_id")
	sr := valueRevision(sol["spec"])
	if contentString(contentFields(sm), "action") != "select" || sr == nil {
		return block(queueError("task_queue_selection_required", "queue entry must bind a select event"))
	}
	row.SpecID = &sid
	row.SpecRevision = &sr.Revision
	spec, e := registeredTaskSpec(d, root, r, t.ID, sid, sr)
	if e != nil {
		return block(e)
	}
	basis := contentFields(contentFields(spec)["implementation_basis"])
	oid := contentString(basis, "parent_oid")
	row.ParentOID = &oid
	ep, b := queueTargetBinding(r, tgt)
	for _, v := range epicBases(r, ep, b) {
		if v.OID == oid {
			n := v.Revision
			row.BaseRevision = &n
		}
	}
	if t.Worktree != nil || t.WorktreeState != WorkItemUnbound {
		return block(queueError("task_queue_task_bound", "Task already owns a worktree or intent; adoption is not supported"))
	}
	if op, _ := findOperation(r, t.ID); op != nil {
		return block(queueError("task_queue_task_bound", "Task has a preserved intent"))
	}
	state := taskContentState(r, t.ID)
	if !contentTypedEqual(entry.Selection, state.SelectionEvent) {
		return block(queueError("task_queue_selection_stale", "queue binding is not the current human selection"))
	}
	eval, e := currentTaskSpec(d, root, r, projects, *t, false)
	if e != nil {
		return block(e)
	}
	if eval.SelectionFreshness != "current" || eval.ContentIntegrity != "valid" {
		return block(queueError("task_queue_selection_stale", "problem or assessment no longer matches the selection"))
	}
	if e = taskSpecStructurallyReady(spec); e != nil {
		return block(e)
	}
	if e = validateSpecImplementationBasis(d, projects, r, *t, contentFields(spec)["implementation_basis"], true); e != nil {
		return block(e)
	}
	if len(queueUnresolved(r, tgt)) > 0 {
		return block(queueError("task_queue_effect_pending", "resolve outstanding target effects first"))
	}
	_, repo, e := projectAndRepo(projects, t.ProjectID, t.RepoID)
	if e != nil {
		return block(e)
	}
	branch, path := queueTaskNames(qid, t.ID, title, tgt.ParentLocator)
	if e = d.WorkGit.ValidateBranch(branch); e != nil {
		return block(e)
	}
	if _, exists, e := canonicalTaskTarget(d.Files, root, path); e != nil {
		return block(e)
	} else if exists {
		return block(queueError("task_queue_path_conflict", "proposed Task path already exists"))
	}
	current := currentEpicBinding(r, ep, b)
	if e = precheckParentAndResources(d, repo, current, "refs/heads/"+branch, path); e != nil {
		return block(e)
	}
	for _, o := range r.WorktreeOperations {
		if o.GitCommonDir == tgt.GitCommonDir && (o.SourceRef == "refs/heads/"+branch || o.TargetLocator == path) {
			return block(queueError("task_queue_resource_conflict", "proposed resource belongs to an existing intent"))
		}
	}
	obs, e := queueParentObservation(d, tgt)
	if e != nil {
		return block(e)
	}
	workspace, e := queueWorkspace(d, root)
	if e != nil {
		return block(e)
	}
	base := currentEpicBase(r, ep, b)
	p := WorkspaceTaskPreparePlan{Kind: "WorkspaceTaskPreparePlan@1", SchemaVersion: 1, Workspace: workspace, Target: tgt, RegistryUpgrade: queueUpgrade(r), RegistryFormatVersion: r.FormatVersion, QueueID: qid, TaskID: t.ID, Selection: *entry.Selection, Problem: *valueRevision(sol["problem"]), SpecID: sid, Spec: *sr, Assessment: *valueDecision(sol["assessment"]), RequiredInputs: eval.RequiredInputs, BaseRevision: base.Revision, ParentOID: base.OID, ParentTree: base.Tree, Branch: branch, WorktreePath: path, ObservedParent: obs}
	row.State = "ready"
	row.Reasons = []QueueReason{}
	return row, &p
}
func currentPreparation(d Dependencies, r WorkItemRegistry, p TaskPreparation) QueueCurrent {
	state := "preparing"
	reasons := []QueueReason{}
	if p.Outcome != nil {
		state = "prepared"
	} else if op, _ := findOperation(r, p.Plan.TaskID); op != nil {
		switch op.LastObservation.Classification {
		case "no_effect":
			state = "blocked"
			reasons = queueReasons("task_preparation_no_effect", "preserved intent has no confirmed worktree; apply the same digest to recover")
		case "partial_or_unknown":
			state = "unknown"
			reasons = queueReasons("task_preparation_unknown", "preserved effect requires explicit inspection and recovery")
		}
	}
	return QueueCurrent{PreparationID: p.ID, TaskID: p.Plan.TaskID, State: state, Disposition: "current", WorktreePath: p.Plan.WorktreePath, Reasons: reasons}
}
func buildQueueReadback(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, in QueueTargetInput, ready bool) (WorkspaceTaskQueueReadback, error) {
	var out WorkspaceTaskQueueReadback
	if e := validateTaskContentClosure(d.TaskContent, root, r, false); e != nil {
		return out, e
	}
	if e := validateQueueClosure(d, root, r); e != nil {
		return out, e
	}
	t, _, e := resolveQueueTarget(d, root, r, projects, in)
	if e != nil {
		return out, e
	}
	w, e := queueWorkspace(d, root)
	if e != nil {
		return out, e
	}
	id := queueID(root, t)
	q := foldQueue(r, id)
	out = WorkspaceTaskQueueReadback{Kind: "WorkspaceTaskQueueReadback@1", SchemaVersion: 1, Workspace: w, Target: t, Registry: QueueRegistryReadback{r.FormatVersion, r.RawSHA256, queueUpgrade(r)}, QueueID: id, Revision: q.Revision, Pending: []QueuePending{}, Reasons: []QueueReason{}}
	if q.Current != nil {
		p := findPreparation(r, *q.Current)
		c := currentPreparation(d, r, *p)
		spec, err := readRegisteredTaskManifest(d.TaskContent, root, r, p.Plan.Spec.ManifestSHA256)
		if err != nil {
			return out, err
		}
		c.Delivery, err = DeliveryAgreementFromSpec(spec)
		if err != nil {
			return out, err
		}
		out.Current = &c
	}
	for i, v := range q.Pending {
		row, _ := queueMixedEntryEvaluation(d, root, r, projects, t, id, v, i+1)
		if !ready || row.State == "ready" {
			out.Pending = append(out.Pending, row)
		}
	}
	if len(out.Pending) == 0 {
		out.Reasons = queueReasons("task_queue_empty", "no pending Tasks match this view")
	}
	return out, nil
}
func ListTaskQueue(d Dependencies, in QueueTargetInput, ready bool) (WorkspaceTaskQueueReadback, error) {
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return WorkspaceTaskQueueReadback{}, e
	}
	var out WorkspaceTaskQueueReadback
	e = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error {
		var e error
		out, e = buildQueueReadback(d, root, s.Registry, s.Projects, in, ready)
		return e
	})
	return out, e
}
func SetTaskQueue(d Dependencies, file string) (WorkspaceTaskQueueReadback, error) {
	var out WorkspaceTaskQueueReadback
	if !filepath.IsAbs(file) {
		return out, queueError("task_queue_invalid_input", "draft path must be absolute")
	}
	physical, e := filepath.EvalSymlinks(file)
	if e != nil || physical != file {
		return out, queueError("task_queue_invalid_input", "draft must have a physical path without symlinks")
	}
	st, e := os.Lstat(file)
	if e != nil {
		return out, e
	}
	if !st.Mode().IsRegular() || st.Size() > 256<<10 {
		return out, queueError("task_queue_invalid_input", "draft must be a regular file of at most 256 KiB")
	}
	raw, e := os.ReadFile(file)
	if e != nil {
		return out, e
	}
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return out, e
	}
	if e = queueDraftRule()(v); e != nil {
		return out, queueError("task_queue_invalid_input", e.Error())
	}
	draft := decodeQueueDraft(v)
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return out, e
	}
	in := QueueTargetInput{draft.ProjectID, draft.RepoID, draft.EpicID}
	request := queueRequest(v)
	e = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			t, _, e := resolveQueueTarget(d, root, r, projects, in)
			if e != nil {
				return e
			}
			id := queueID(root, t)
			for _, event := range r.TaskQueueEvents {
				if event.Kind == "set" && contentString(contentFields(queueRequestValue(event.Request)), "publication_key") == draft.PublicationKey {
					if event.QueueID != id || event.RequestSHA256 != queueDigest(request) {
						return queueError("task_queue_publication_conflict", "publication key is already bound to different inputs")
					}
					out, e = buildQueueReadback(d, root, r, projects, in, false)
					return e
				}
			}
			q := foldQueue(r, id)
			if q.Revision != draft.ExpectedRevision || q.Revision == 2147483647 {
				return queueError("task_queue_revision_conflict", "queue revision changed or is exhausted")
			}
			for _, entry := range draft.Entries {
				task, _ := findTask(r, entry.TaskID)
				if task == nil || task.ProjectID != t.ProjectID || task.RepoID != t.RepoID || task.ParentEpicID != t.EpicID {
					return queueError("task_queue_target_conflict", "entry belongs to another target")
				}
				if task.Worktree != nil || task.WorktreeState != WorkItemUnbound || preparationForTask(r, task.ID) != nil {
					return queueError("task_queue_task_bound", "a bound, current or terminal Task cannot be placed in the pending queue")
				}
				if op, _ := findOperation(r, task.ID); op != nil {
					return queueError("task_queue_task_bound", "Task already has an intent")
				}
				if e = validateQueueSelection(r, task.ID, entry.Selection); e != nil {
					return e
				}
				if entry.Goal != nil {
					if _, e = loadTaskGoal(d, root, r, task.ID, *entry.Goal); e != nil {
						return e
					}
				}
				if entry.Selection != nil {
					m, err := readRegisteredTaskManifest(d.TaskContent, root, r, entry.Selection.ManifestSHA256)
					if err != nil {
						return err
					}
					if contentString(contentFields(m), "action") != "select" {
						return queueError("task_queue_invalid_input", "entry must reference a select event")
					}
				}
			}
			if e = upgradeQueueRegistry(d, root, &r, draft.RegistryUpgrade); e != nil {
				return e
			}
			r.TaskQueueEvents = append(r.TaskQueueEvents, TaskQueueEvent{id, q.Revision + 1, "set", queueDigest(request), request, queueNow(d)})
			sortQueueRegistry(&r)
			if e = session.Publish(r); e != nil {
				observed, readErr := session.Snapshot()
				if readErr == nil {
					out, readErr = buildQueueReadback(d, root, observed, projects, in, false)
				}
				if readErr != nil {
					return fmt.Errorf("queue publication failed (%w); after-state unknown (%v); inspect queue list before retry", e, readErr)
				}
				return fmt.Errorf("queue publication failed (%w); registry before %s, observed revision %d, registry after %s; inspect queue list before retry", e, r.RawSHA256, out.Revision, out.Registry.SHA256)
			}
			r, e = session.Snapshot()
			if e != nil {
				return e
			}
			out, e = buildQueueReadback(d, root, r, projects, in, false)
			return e
		})
	})
	return out, e
}
func validateQueueClosure(d Dependencies, root string, r WorkItemRegistry) error {
	if r.FormatVersion != 4 {
		return nil
	}
	for _, ev := range r.TaskQueueEvents {
		if ev.Kind != "set" {
			continue
		}
		draft := decodeQueueDraft(queueRequestValue(ev.Request))
		t := QueueTarget{ProjectID: draft.ProjectID, RepoID: draft.RepoID, EpicID: draft.EpicID}
		for _, entry := range draft.Entries {
			if !d.TaskContent.includesTask(entry.TaskID) {
				continue
			}
			if entry.Goal != nil {
				if _, e := loadTaskGoal(d, root, r, entry.TaskID, *entry.Goal); e != nil {
					return e
				}
			}
			if entry.Selection != nil {
				m, e := readRegisteredTaskManifest(d.TaskContent, root, r, entry.Selection.ManifestSHA256)
				if e != nil {
					return e
				}
				if contentString(contentFields(m), "action") != "select" {
					return fmt.Errorf("queued decision is not a select event")
				}
			}
		}
		if ev.QueueID != queueID(root, t) {
			return queueError("task_queue_integrity_conflict", "queue identity differs from workspace")
		}
		if draft.RegistryUpgrade != nil && d.TaskContent.includesTarget(r, draft.EpicID, draft.RepoID) {
			if _, e := d.TaskContent.Read(root, "backups", draft.RegistryUpgrade.RegistrySHA256); e != nil {
				return e
			}
		}
	}
	for _, p := range r.TaskPreparations {
		if !d.TaskContent.includesTask(p.Plan.TaskID) {
			continue
		}
		if p.Plan.Workspace.Root != root {
			return queueError("task_queue_integrity_conflict", "preparation workspace changed")
		}
		if p.Plan.RegistryUpgrade != nil {
			if _, e := d.TaskContent.Read(root, "backups", p.Plan.RegistryUpgrade.RegistrySHA256); e != nil {
				return e
			}
		}
		b := TaskSpecBasis{TaskID: p.Plan.TaskID, Problem: p.Plan.Problem, SpecID: p.Plan.SpecID, Spec: p.Plan.Spec, Assessment: p.Plan.Assessment, Selection: p.Plan.Selection, Dependencies: []string{}}
		eval, e := loadTaskSpecContent(d, root, r, b)
		if e != nil {
			return e
		}
		basis := contentFields(contentFields(eval.Spec)["implementation_basis"])
		if contentString(basis, "parent_oid") != p.Plan.ParentOID || contentString(basis, "parent_tree") != p.Plan.ParentTree || contentString(basis, "start_oid") != p.Plan.ParentOID || contentString(basis, "start_tree") != p.Plan.ParentTree {
			return queueError("task_queue_integrity_conflict", "prepare plan differs from selected implementation basis")
		}
		if !contentTypedEqual(eval.RequiredInputs, p.Plan.RequiredInputs) {
			return queueError("task_queue_integrity_conflict", "preserved preparation inputs differ")
		}
	}
	for _, u := range r.EpicBaseUpdates {
		if !d.TaskContent.includesTarget(r, u.Plan.Target.EpicID, u.Plan.Target.RepoID) {
			continue
		}
		if u.Plan.Workspace.Root != root {
			return queueError("epic_base_integrity_conflict", "base workspace changed")
		}
		if u.Plan.RegistryUpgrade != nil {
			if _, e := d.TaskContent.Read(root, "backups", u.Plan.RegistryUpgrade.RegistrySHA256); e != nil {
				return e
			}
		}
	}
	return nil
}
func preparationIdentityKnown(d Dependencies, repo RepoRecord, p TaskPreparation) bool {
	o, e := d.WorkGit.ObserveWorktree(p.Plan.WorktreePath)
	if e != nil {
		return false
	}
	ref, e := d.WorkGit.ObserveRef(repo, "refs/heads/"+p.Plan.Branch)
	return e == nil && o.Locator == p.Plan.WorktreePath && o.GitCommonDir == p.Plan.Target.GitCommonDir && o.Ref == "refs/heads/"+p.Plan.Branch && o.InventoryMatch && ref.Exists && ref.CheckedOutAt == p.Plan.WorktreePath && ref.OID == o.OID
}
func preparationNoEffect(d Dependencies, repo RepoRecord, p TaskPreparation) bool {
	ref, e := d.WorkGit.ObserveRef(repo, "refs/heads/"+p.Plan.Branch)
	if e != nil || ref.Exists {
		return false
	}
	_, e = d.Files.Lstat(p.Plan.WorktreePath)
	if !errors.Is(e, fs.ErrNotExist) {
		return false
	}
	inventory, e := d.WorkGit.ListWorktrees(repo)
	if e != nil {
		return false
	}
	_, exists := worktreeEntry(inventory, p.Plan.WorktreePath)
	return !exists
}
func CloseTaskQueue(d Dependencies, in QueueTargetInput, id string, expected int, reason, kind string) (WorkspaceTaskQueueReadback, error) {
	var out WorkspaceTaskQueueReadback
	if kind != "advance" && kind != "release" {
		return out, WorkInvalidArguments("unknown queue transition")
	}
	if e := validateWorkText("reason", reason, 2000); e != nil {
		return out, e
	}
	if expected < 0 || expected >= 2147483647 {
		return out, WorkInvalidArguments("expected revision must be between 0 and 2147483646")
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return out, e
	}
	request := map[string]any{"preparation_id": id, "expected_revision": expected, "reason": reason}
	e = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			t, repo, e := resolveQueueTarget(d, root, r, projects, in)
			if e != nil {
				return e
			}
			qid := queueID(root, t)
			for _, ev := range r.TaskQueueEvents {
				if ev.QueueID == qid && ev.Kind == kind && ev.RequestSHA256 == queueDigest(request) {
					out, e = buildQueueReadback(d, root, r, projects, in, false)
					return e
				}
			}
			q := foldQueue(r, qid)
			if q.Revision != expected || q.Current == nil || *q.Current != id {
				return queueError("task_queue_revision_conflict", "current preparation or queue revision differs")
			}
			p := findPreparation(r, id)
			known := p.Outcome != nil && preparationIdentityKnown(d, repo, *p)
			if kind == "release" && p.Outcome == nil {
				known = preparationNoEffect(d, repo, *p)
			}
			if !known {
				return queueError("task_queue_effect_unknown", "resolve the preserved effect before closing this preparation")
			}
			r.TaskQueueEvents = append(r.TaskQueueEvents, TaskQueueEvent{qid, q.Revision + 1, kind, queueDigest(request), request, queueNow(d)})
			sortQueueRegistry(&r)
			if e = session.Publish(r); e != nil {
				observed, readErr := session.Snapshot()
				if readErr == nil {
					out, readErr = buildQueueReadback(d, root, observed, projects, in, false)
				}
				if readErr != nil {
					return fmt.Errorf("queue publication failed (%w); after-state unknown (%v); inspect queue list before retry", e, readErr)
				}
				return fmt.Errorf("queue publication failed (%w); registry before %s, observed revision %d, registry after %s; inspect queue list before retry", e, r.RawSHA256, out.Revision, out.Registry.SHA256)
			}
			r, e = session.Snapshot()
			if e != nil {
				return e
			}
			out, e = buildQueueReadback(d, root, r, projects, in, false)
			return e
		})
	})
	return out, e
}
func MarshalTaskQueue(value any) ([]byte, error) { return contentCanonical(value) }
func TaskQueueText(r WorkspaceTaskQueueReadback) string {
	s := fmt.Sprintf("Task queue %s, revision %d.\nTarget: %s / %s / %s\nWorkplace: %s\n", r.QueueID, r.Revision, r.Target.ProjectID, r.Target.RepoID, r.Target.EpicID, r.Target.ParentLocator)
	if r.Current != nil {
		s += fmt.Sprintf("Current: %s (%s), %s, %s\n", r.Current.TaskID, r.Current.State, r.Current.PreparationID, r.Current.WorktreePath)
		if r.Current.Delivery != nil {
			s += fmt.Sprintf("Delivery: %s to %s\n", r.Current.Delivery.Mode, r.Current.Delivery.TargetRef)
		}
	}
	for _, p := range r.Pending {
		s += fmt.Sprintf("%d. %s — %s [%s]\n", p.Rank, p.TaskID, p.Title, p.State)
		if p.Delivery != nil {
			s += fmt.Sprintf("   Delivery: %s to %s\n", p.Delivery.Mode, p.Delivery.TargetRef)
		}
		for _, reason := range p.Reasons {
			s += "   " + reason.Code + ": " + reason.Message + "\n"
		}
	}
	for _, reason := range r.Reasons {
		s += reason.Message + "\n"
	}
	return s + fmt.Sprintf("Next action: ply workspace task prepare --next --project %s --repo %s --epic %s\n", r.Target.ProjectID, r.Target.RepoID, r.Target.EpicID)
}

func sortQueueReasons(in []QueueReason) []QueueReason {
	sort.Slice(in, func(i, j int) bool { return in[i].Code < in[j].Code })
	out := []QueueReason{}
	for _, v := range in {
		if len(out) == 0 || out[len(out)-1].Code != v.Code {
			out = append(out, v)
		}
	}
	return out
}
