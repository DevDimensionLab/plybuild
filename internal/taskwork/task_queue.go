package taskwork

import (
	"errors"
	"fmt"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"os"
	"path/filepath"
	"sort"
)

func queueReasons(code, message string) []QueueReason { return []QueueReason{{code, message}} }
func queueErrorReasons(e error) []QueueReason {
	var ce *TaskContentError
	if errors.As(e, &ce) {
		return queueReasons(ce.Code, ce.Detail)
	}
	return queueReasons("task_queue_observation_unknown", e.Error())
}
func queueTaskTitle(d Dependencies, work workspace.ProjectWork, r WorkItemRegistry, t TaskRecord) (string, error) {
	p := taskContentState(r, t.ID).ProblemHead
	if p == nil {
		return t.Title, nil
	}
	m, e := readRegisteredTaskManifest(d.TaskContent, work, r, p.ManifestSHA256)
	if e != nil {
		return "", e
	}
	return contentString(contentFields(m), "title"), nil
}
func queuePendingEvaluation(d Dependencies, work workspace.ProjectWork, r WorkItemRegistry, entry TaskGoalQueueEntry, rank int) QueuePending {
	if entry.Goal != nil {
		return queueGoalEvaluation(d, work, r, entry, rank)
	}
	row := QueuePending{Rank: rank, TaskID: entry.TaskID, Selection: entry.Selection, State: "blocked", Reasons: queueReasons("task_queue_goal_required", "this entry names a selected solution; register a Goal for the Task")}
	if t, _ := findTask(r, entry.TaskID); t != nil {
		if title, e := queueTaskTitle(d, work, r, *t); e == nil {
			row.Title = title
		}
	}
	return row
}
func buildQueueReadback(d Dependencies, work workspace.ProjectWork, r WorkItemRegistry, projects workspace.ProjectSnapshot, in QueueTargetInput, ready bool) (WorkspaceTaskQueueReadback, error) {
	var out WorkspaceTaskQueueReadback
	t, _, e := resolveQueueTarget(d, work, r, projects, in)
	if e != nil {
		return out, e
	}
	d = withQueueContentScope(d, r, t)
	if e := validateTaskContentClosure(d.TaskContent, work, r, false); e != nil {
		return out, e
	}
	if e := validateQueueClosure(d, work, r); e != nil {
		return out, e
	}
	w, e := queueWorkspace(d, work.Workspace)
	if e != nil {
		return out, e
	}
	id := queueID(work.Workspace, t)
	q := foldQueue(r, id)
	out = WorkspaceTaskQueueReadback{Kind: "WorkspaceTaskQueueReadback@1", SchemaVersion: 1, Workspace: w, Target: t, Registry: registryReadback(r), QueueID: id, Revision: q.Revision, Pending: []QueuePending{}, Reasons: []QueueReason{}}
	for i, v := range q.Pending {
		row := queuePendingEvaluation(d, work, r, v, i+1)
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
	return ListTaskQueueVersion(d, in, ready, 1)
}

func ListTaskQueueVersion(d Dependencies, in QueueTargetInput, ready bool, version int) (WorkspaceTaskQueueReadback, error) {
	if version != 1 && version != 2 {
		return WorkspaceTaskQueueReadback{}, workspace.WorkInvalidArguments("queue schema-version must be 1 or 2")
	}

	work, e := queueTargetWork(d, in)
	if e != nil {
		return WorkspaceTaskQueueReadback{}, e
	}
	var out WorkspaceTaskQueueReadback
	e = WithTaskSpecSnapshot(d, work, func(s *TaskSpecSession) error {
		var e error
		out, e = buildQueueReadback(d, work, s.Registry, s.Projects, in, ready)
		return e
	})
	if e == nil && version == 2 {
		out.Kind = "WorkspaceTaskQueueReadback@2"
		out.SchemaVersion = 2
	}
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
	in := QueueTargetInput{draft.ProjectID, draft.RepoID, draft.EpicID}
	work, e := queueTargetWork(d, in)
	if e != nil {
		return out, e
	}
	request := queueRequest(v)
	e = d.ProjectLocks.WithSnapshotLock(work.Workspace, func(projects workspace.ProjectSnapshot) error {
		return d.WorkItems.WithLock(work, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			t, _, e := resolveQueueTarget(d, work, r, projects, in)
			if e != nil {
				return e
			}
			id := queueID(work.Workspace, t)
			for _, event := range r.TaskQueueEvents {
				if event.Kind == "set" && contentString(contentFields(queueRequestValue(event.Request)), "publication_key") == draft.PublicationKey {
					if event.QueueID != id || event.RequestSHA256 != queueDigest(request) {
						return queueError("task_queue_publication_conflict", "publication key is already bound to different inputs")
					}
					out, e = buildQueueReadback(d, work, r, projects, in, false)
					return e
				}
			}
			q := foldQueue(r, id)
			if q.Revision != draft.ExpectedRevision || q.Revision == 2147483647 {
				return queueError("task_queue_revision_conflict", "queue revision changed or is exhausted")
			}
			pending := make(map[TaskID]TaskGoalQueueEntry, len(q.Pending))
			for _, entry := range q.Pending {
				pending[entry.TaskID] = entry
			}
			for _, entry := range draft.Entries {
				task, _ := findTask(r, entry.TaskID)
				if task == nil || task.ProjectID != t.ProjectID || task.RepoID != t.RepoID || task.ParentEpicID != t.EpicID {
					return queueError("task_queue_target_conflict", "entry belongs to another target")
				}
				previous, retained := pending[entry.TaskID]
				if entry.Goal == nil {
					return queueError("task_queue_goal_required", "a queue entry names a registered Goal")
				}
				if !retained || !contentTypedEqual(previous, entry) {
					if task.Worktree != nil || task.WorktreeState != WorkItemUnbound {
						return queueError("task_queue_task_bound", "a bound, current or terminal Task cannot be placed in the pending queue")
					}
					if op, _ := findOperation(r, task.ID); op != nil {
						return queueError("task_queue_task_bound", "Task already has an intent")
					}
				}
				if _, e = loadTaskGoal(d, work, r, task.ID, *entry.Goal); e != nil {
					return e
				}
			}
			r.TaskQueueEvents = append(r.TaskQueueEvents, TaskQueueEvent{id, q.Revision + 1, "set", queueDigest(request), request, queueNow(d)})
			sortQueueRegistry(&r)
			if e = session.Publish(r); e != nil {
				observed, readErr := session.Snapshot()
				if readErr == nil {
					out, readErr = buildQueueReadback(d, work, observed, projects, in, false)
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
			out, e = buildQueueReadback(d, work, r, projects, in, false)
			return e
		})
	})
	return out, e
}
func validateQueueClosure(d Dependencies, work workspace.ProjectWork, r WorkItemRegistry) error {
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
				if _, e := loadTaskGoal(d, work, r, entry.TaskID, *entry.Goal); e != nil {
					return e
				}
			}
			if entry.Selection != nil {
				m, e := readRegisteredTaskManifest(d.TaskContent, work, r, entry.Selection.ManifestSHA256)
				if e != nil {
					return e
				}
				if contentString(contentFields(m), "action") != "select" {
					return fmt.Errorf("queued decision is not a select event")
				}
			}
		}
		if ev.QueueID != queueID(work.Workspace, t) {
			return queueError("task_queue_integrity_conflict", "queue identity differs from workspace")
		}
	}
	for _, u := range r.EpicBaseUpdates {
		if !d.TaskContent.includesTarget(r, u.Plan.Target.EpicID, u.Plan.Target.RepoID) {
			continue
		}
		if u.Plan.Workspace.Root != work.Workspace {
			return queueError("epic_base_integrity_conflict", "base workspace changed")
		}
	}
	return nil
}
func MarshalTaskQueue(value any) ([]byte, error) { return contentCanonical(value) }
func TaskQueueText(r WorkspaceTaskQueueReadback) string {
	s := fmt.Sprintf("Task queue %s, revision %d.\nTarget: %s / %s / %s\nWorkplace: %s\n", r.QueueID, r.Revision, r.Target.ProjectID, r.Target.RepoID, r.Target.EpicID, r.Target.ParentLocator)
	for _, p := range r.Pending {
		s += fmt.Sprintf("%d. %s — %s [%s]\n", p.Rank, p.TaskID, p.Title, p.State)
		for _, reason := range p.Reasons {
			s += "   " + reason.Code + ": " + reason.Message + "\n"
		}
	}
	for _, reason := range r.Reasons {
		s += reason.Message + "\n"
	}
	return s
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
