package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func adoptedEpicBase(e EpicRecord, b EpicRepoBinding) EpicBaseVersion {
	w := b.Worktree
	return EpicBaseVersion{ProjectID: e.ProjectID, RepoID: b.RepoID, EpicID: e.ID, Revision: 1, WorktreeID: w.ID, Locator: w.Locator, Ref: w.Ref, GitCommonDir: b.GitCommonDir, OID: w.OID, Tree: w.Tree}
}
func epicBases(r WorkItemRegistry, e EpicRecord, b EpicRepoBinding) []EpicBaseVersion {
	if r.FormatVersion < 4 {
		return []EpicBaseVersion{adoptedEpicBase(e, b)}
	}
	out := []EpicBaseVersion{}
	for _, v := range r.EpicBaseVersions {
		if v.EpicID == e.ID && v.RepoID == b.RepoID {
			out = append(out, v)
		}
	}
	return out
}
func currentEpicBase(r WorkItemRegistry, e EpicRecord, b EpicRepoBinding) EpicBaseVersion {
	v := epicBases(r, e, b)
	return v[len(v)-1]
}
func currentEpicBinding(r WorkItemRegistry, e EpicRecord, b EpicRepoBinding) EpicRepoBinding {
	v := currentEpicBase(r, e, b)
	b.Worktree.OID = v.OID
	b.Worktree.Tree = v.Tree
	return b
}
func historicalEpicBaseMatches(r WorkItemRegistry, epic EpicID, repo RepoID, wt WorktreeID, ref, oid, tree string) bool {
	e, _ := findEpic(r, epic)
	if e == nil {
		return false
	}
	b, _ := findEpicRepo(*e, repo)
	if b == nil {
		return false
	}
	for _, v := range epicBases(r, *e, *b) {
		if v.WorktreeID == wt && v.Ref == ref && v.OID == oid && v.Tree == tree {
			return true
		}
	}
	return false
}
func historicalEpicBinding(r WorkItemRegistry, e EpicRecord, b EpicRepoBinding, oid, tree string) (EpicRepoBinding, bool) {
	if !historicalEpicBaseMatches(r, e.ID, b.RepoID, b.Worktree.ID, b.Worktree.Ref, oid, tree) {
		return b, false
	}
	b.Worktree.OID = oid
	b.Worktree.Tree = tree
	return b, true
}
func queueError(code, detail string) error { return contentError(code, detail, nil) }
func queueNow(d Dependencies) string       { return d.WorkClock.Now().UTC().Format(time.RFC3339Nano) }
func queueDigest(v any) string {
	b, e := contentCanonical(v)
	if e != nil {
		return ""
	}
	return digestTaskBytes(b)
}
func queueID(root string, t QueueTarget) string {
	return "que_" + strings.TrimPrefix(queueDigest(map[string]any{"workspace_root": root, "project_id": t.ProjectID, "repo_id": t.RepoID, "epic_id": t.EpicID}), "sha256:")
}
func queueWorkspace(d Dependencies, root string) (IntegrationPlanWorkspace, error) {
	b, e := d.Files.ReadFile(filepath.Join(root, MarkerDirectory, MarkerFile))
	return IntegrationPlanWorkspace{Root: root, MarkerSHA256: digestTaskBytes(b)}, e
}
func queueUpgrade(r WorkItemRegistry) *TaskRegistryUpgrade {
	if r.FormatVersion == 4 {
		return nil
	}
	return &TaskRegistryUpgrade{FromVersion: r.FormatVersion, RegistrySHA256: r.RawSHA256}
}
func upgradeQueueRegistry(d Dependencies, root string, r *WorkItemRegistry, u *TaskRegistryUpgrade) error {
	if r.FormatVersion == 4 {
		if u != nil {
			return queueError("registry_upgrade_conflict", "format 4 requires a null upgrade")
		}
		return nil
	}
	if u == nil {
		return queueError("registry_upgrade_required", "explicit format and registry digest are required")
	}
	if !contentTypedEqual(u, queueUpgrade(*r)) {
		return queueError("registry_upgrade_conflict", "registry changed since the upgrade was confirmed")
	}
	b, e := os.ReadFile(workItemsPath(root))
	if e != nil {
		return e
	}
	if digestTaskBytes(b) != u.RegistrySHA256 {
		return queueError("registry_upgrade_conflict", "registry bytes changed")
	}
	if _, e = d.TaskContent.Publish(root, "backups", b); e != nil {
		return e
	}
	if r.FormatVersion < 3 {
		upgradeTaskContentRegistry(r)
	}
	r.FormatVersion = 4
	r.EpicBaseVersions = []EpicBaseVersion{}
	r.EpicBaseUpdates = []EpicBaseUpdate{}
	r.TaskQueueEvents = []TaskQueueEvent{}
	r.TaskPreparations = []TaskPreparation{}
	for _, ep := range r.Epics {
		for _, b := range ep.RepoBindings {
			r.EpicBaseVersions = append(r.EpicBaseVersions, adoptedEpicBase(ep, b))
		}
	}
	sortQueueRegistry(r)
	return nil
}

type QueueTargetInput struct {
	ProjectID ProjectID
	RepoID    RepoID
	EpicID    EpicID
}

func resolveQueueTarget(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, in QueueTargetInput) (QueueTarget, RepoRecord, error) {
	explicit := in.ProjectID != "" || in.RepoID != "" || in.EpicID != ""
	if explicit && (in.ProjectID == "" || in.RepoID == "" || in.EpicID == "") {
		return QueueTarget{}, RepoRecord{}, queueError("task_queue_context_required", "provide all of --project <id> --repo <id> --epic <id>")
	}
	cwd, e := projectPhysicalWorkingDirectory(d.Files)
	if e != nil {
		return QueueTarget{}, RepoRecord{}, e
	}
	actualRoot := ""
	if !explicit {
		for dir := cwd; dir == root || strings.HasPrefix(dir, root+string(filepath.Separator)); dir = filepath.Dir(dir) {
			if _, err := d.Files.Lstat(filepath.Join(dir, ".git")); err == nil {
				o, err := d.WorkGit.ObserveWorktree(dir)
				if err != nil {
					return QueueTarget{}, RepoRecord{}, queueError("task_queue_context_required", "inspect the current Git worktree or supply --project <id> --repo <id> --epic <id>")
				}
				actualRoot = o.Locator
				break
			}
			if dir == root {
				break
			}
		}
	}

	matches := []QueueTarget{}
	for _, ep := range r.Epics {
		for _, b := range ep.RepoBindings {
			if explicit && (ep.ID != in.EpicID || ep.ProjectID != in.ProjectID || b.RepoID != in.RepoID) {
				continue
			}
			if !explicit && actualRoot != b.Worktree.Locator {
				continue
			}
			if !explicit {
				o, err := d.WorkGit.ObserveWorktree(b.Worktree.Locator)
				if err != nil || o.Locator != b.Worktree.Locator || o.GitCommonDir != b.GitCommonDir {
					continue
				}
			}
			matches = append(matches, QueueTarget{ep.ProjectID, b.RepoID, ep.ID, b.GitCommonDir, b.Worktree.ID, b.Worktree.Locator, b.Worktree.Ref})
		}
	}
	if len(matches) != 1 {
		return QueueTarget{}, RepoRecord{}, queueError("task_queue_context_required", "provide --project <id> --repo <id> --epic <id> for exactly one registered Epic worktree")
	}
	t := matches[0]
	_, repo, e := projectAndRepo(projects, t.ProjectID, t.RepoID)
	if e != nil {
		return t, repo, e
	}
	physical, e := filepath.EvalSymlinks(t.ParentLocator)
	if e != nil || physical != t.ParentLocator || repo.GitCommonDir != t.GitCommonDir {
		return t, repo, queueError("task_queue_target_conflict", "physical Epic or repository anchor changed")
	}
	obs, e := d.WorkGit.ObserveWorktree(t.ParentLocator)
	if e != nil {
		return t, repo, e
	}
	if obs.Locator != t.ParentLocator || obs.GitCommonDir != t.GitCommonDir {
		return t, repo, queueError("task_queue_target_conflict", "Epic worktree identity changed")
	}
	return t, repo, nil
}
func queueTargetBinding(r WorkItemRegistry, t QueueTarget) (EpicRecord, EpicRepoBinding) {
	e, _ := findEpic(r, t.EpicID)
	b, _ := findEpicRepo(*e, t.RepoID)
	return *e, *b
}
func queueParentObservation(d Dependencies, t QueueTarget) (PlanWorktreeObservation, error) {
	o, e := d.WorkGit.ObserveWorktree(t.ParentLocator)
	if e != nil {
		return PlanWorktreeObservation{}, e
	}
	if !o.InventoryMatch || o.Locator != t.ParentLocator || o.Ref != t.ParentRef || o.GitCommonDir != t.GitCommonDir || !o.Clean {
		return PlanWorktreeObservation{}, queueError("epic_base_stale", "Epic must have the same clean worktree, branch and repository")
	}
	if d.IntegrationGit == nil {
		return PlanWorktreeObservation{}, queueError("task_queue_observation_unknown", "integration Git observer is unavailable")
	}
	full, e := d.IntegrationGit.ObserveIntegrationWorktree(t.ParentLocator, t.ParentRef)
	if e != nil {
		return PlanWorktreeObservation{}, e
	}
	if full.OID != o.OID || full.Tree != o.Tree || !full.Symbolic || !full.Clean || len(full.InProgress) != 0 {
		return PlanWorktreeObservation{}, queueError("task_queue_observation_unknown", "Epic observation changed or has an operation in progress")
	}
	return planObservation(full), nil
}
func queueUnresolved(r WorkItemRegistry, t QueueTarget) []string {
	out := []string{}
	for _, o := range r.WorktreeOperations {
		if o.ParentEpicID == t.EpicID && o.RepoID == t.RepoID && o.State != "ready" {
			p := preparationForTask(r, o.TaskID)
			if p != nil && preparationDisposition(r, p.ID) == "released" {
				continue
			}
			out = append(out, string(o.ID))
		}
	}
	for _, a := range r.IntegrationAuthorities {
		if a.Plan.Epic.EpicID == t.EpicID && a.Plan.Repository.RepoID == t.RepoID && !taskIntegrationAuthorityComplete(r, a.ID) {
			out = append(out, string(a.ID))
		}
	}
	for _, u := range r.EpicBaseUpdates {
		if u.Plan.Target == t && u.Phase == "intent" {
			out = append(out, u.ID)
		}
	}
	sort.Strings(out)
	return out
}
func pendingBaseUpdate(r WorkItemRegistry, epic EpicID, repo RepoID) bool {
	for _, u := range r.EpicBaseUpdates {
		if u.Plan.Target.EpicID == epic && u.Plan.Target.RepoID == repo && u.Phase == "intent" {
			return true
		}
	}
	return false
}

// Base updates publish metadata only. Git remains an independently observed input.
type EpicBaseInput struct {
	EpicID       EpicID
	RepoID       RepoID
	Apply        bool
	Confirmation string
	Upgrade      string
}
type EpicBaseResult struct {
	Observed        *PlanWorktreeObservation     `json:"observed"`
	PersistedBefore *QueueRegistryReadback       `json:"persisted_before"`
	PersistedAfter  *QueueRegistryReadback       `json:"persisted_after"`
	Kind            string                       `json:"kind"`
	SchemaVersion   int                          `json:"schema_version"`
	Workspace       string                       `json:"workspace"`
	Target          QueueTarget                  `json:"target"`
	Versions        []EpicBaseVersion            `json:"versions"`
	Plan            *WorkspaceEpicBaseUpdatePlan `json:"plan"`
	Confirmation    *string                      `json:"confirmation"`
	Operation       *EpicBaseUpdate              `json:"operation"`
	State           string                       `json:"state"`
	Freshness       string                       `json:"freshness"`
	Reasons         []QueueReason                `json:"reasons"`
	BackupSHA256    *string                      `json:"backup_sha256"`
}

func baseResult(d Dependencies, root string, r WorkItemRegistry, t QueueTarget) EpicBaseResult {
	ep, b := queueTargetBinding(r, t)
	v := epicBases(r, ep, b)
	out := EpicBaseResult{Kind: "WorkspaceEpicBaseReadback@1", SchemaVersion: 1, Workspace: root, Target: t, Versions: v, State: "recorded", Freshness: "unknown", Reasons: []QueueReason{}}
	if d.IntegrationGit != nil {
		if full, e := d.IntegrationGit.ObserveIntegrationWorktree(t.ParentLocator, t.ParentRef); e == nil {
			o := planObservation(full)
			out.Observed = &o
		}
	}
	if o, e := queueParentObservation(d, t); e == nil {
		last := v[len(v)-1]
		out.Freshness = freshness(o.OID == last.OID && o.Tree == last.Tree)
	}
	return out
}
func buildBasePlan(d Dependencies, root string, r WorkItemRegistry, projects ProjectSnapshot, in EpicBaseInput) (EpicBaseResult, error) {
	ep, _ := findEpic(r, in.EpicID)
	if ep == nil {
		return EpicBaseResult{}, queueError("epic_not_found", "Epic is not registered")
	}
	t, repo, e := resolveQueueTarget(d, root, r, projects, QueueTargetInput{ep.ProjectID, in.RepoID, in.EpicID})
	if e != nil {
		return EpicBaseResult{}, e
	}
	out := baseResult(d, root, r, t)
	obs, e := queueParentObservation(d, t)
	if e != nil {
		return out, e
	}
	b := out.Versions[len(out.Versions)-1]
	if obs.OID == b.OID {
		out.State = "unchanged"
		return out, nil
	}
	if b.Revision == 2147483647 {
		return out, queueError("epic_base_conflict", "base revision exhausted")
	}
	ancestor, e := d.IntegrationGit.CheckAncestor(repo, b.OID, obs.OID)
	if e != nil {
		return out, e
	}
	if !ancestor {
		return out, queueError("epic_base_conflict", "new Epic commit must descend from the current base")
	}
	unresolved := queueUnresolved(r, t)
	if len(unresolved) > 0 {
		return out, queueError("epic_base_blocked", "resolve pending operations before updating the base: "+strings.Join(unresolved, ", "))
	}
	w, e := queueWorkspace(d, root)
	if e != nil {
		return out, e
	}
	q := foldQueue(r, queueID(root, t))
	affected := []EpicAffectedSelection{}
	for _, task := range r.Tasks {
		if task.ParentEpicID == t.EpicID && task.RepoID == t.RepoID {
			affected = append(affected, EpicAffectedSelection{task.ID, taskContentState(r, task.ID).SelectionEvent})
		}
	}
	p := WorkspaceEpicBaseUpdatePlan{Kind: "WorkspaceEpicBaseUpdatePlan@1", SchemaVersion: 1, Workspace: w, Target: t, RegistryUpgrade: queueUpgrade(r), RegistryFormatVersion: r.FormatVersion, ExpectedRevision: b.Revision, PreviousOID: b.OID, PreviousTree: b.Tree, NextOID: obs.OID, NextTree: obs.Tree, ObservedParent: obs, AffectedSelections: affected, CurrentPreparationID: q.Current, UnresolvedOperations: unresolved}
	h := queueDigest(p)
	out.Plan = &p
	out.Confirmation = &h
	out.State = "preview"
	return out, nil
}
func ShowEpicBase(d Dependencies, in EpicBaseInput) (EpicBaseResult, error) {
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return EpicBaseResult{}, e
	}
	var out EpicBaseResult
	e = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error {
		ep, _ := findEpic(s.Registry, in.EpicID)
		if ep == nil {
			return queueError("epic_not_found", "Epic is not registered")
		}
		t, _, e := resolveQueueTarget(d, root, s.Registry, s.Projects, QueueTargetInput{ep.ProjectID, in.RepoID, in.EpicID})
		if e != nil {
			return e
		}
		out = baseResult(d, root, s.Registry, t)
		return nil
	})
	return out, e
}
func UpdateEpicBase(d Dependencies, in EpicBaseInput) (EpicBaseResult, error) {
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return EpicBaseResult{}, e
	}
	var out EpicBaseResult
	if !in.Apply {
		e = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error {
			var e error
			out, e = buildBasePlan(d, root, s.Registry, s.Projects, in)
			return e
		})
		return out, e
	}
	if !digestPattern.MatchString(in.Confirmation) {
		return out, queueError("epic_base_confirmation_required", "apply requires the exact preview digest")
	}
	e = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(session WorkItemStoreSession) error {
			r, e := session.Snapshot()
			if e != nil {
				return e
			}
			before := QueueRegistryReadback{r.FormatVersion, r.RawSHA256, queueUpgrade(r)}
			defer func() {
				out.PersistedBefore = &before
				if after, err := session.Snapshot(); err == nil {
					out.PersistedAfter = &QueueRegistryReadback{after.FormatVersion, after.RawSHA256, queueUpgrade(after)}
				}
			}()
			id := "ebu_" + strings.TrimPrefix(in.Confirmation, "sha256:")
			for i, u := range r.EpicBaseUpdates {
				if u.ID == id {
					if u.Plan.Target.EpicID != in.EpicID || u.Plan.Target.RepoID != in.RepoID {
						return queueError("epic_base_conflict", "confirmation belongs to another target")
					}
					out = baseResult(d, root, r, u.Plan.Target)
					out.Operation = &u
					out.State = u.Phase
					if u.Phase == "conflict" {
						return queueError("epic_base_conflict", "preserved request closed without a new base; inspect ply workspace epic base operation show "+u.ID)
					}
					if u.Phase != "intent" {
						return nil
					}
					return finishBaseUpdate(d, root, projects, session, &r, i, &out)
				}
			}
			out, e = buildBasePlan(d, root, r, projects, in)
			if e != nil {
				return e
			}
			if out.Confirmation == nil || *out.Confirmation != in.Confirmation {
				return queueError("epic_base_conflict", "preview changed; inspect a new preview")
			}
			p := *out.Plan
			if p.RegistryUpgrade != nil && in.Upgrade != p.RegistryUpgrade.RegistrySHA256 {
				return queueError("registry_upgrade_required", "confirm the exact registry upgrade digest")
			}
			if p.RegistryUpgrade == nil && in.Upgrade != "" {
				return queueError("registry_upgrade_conflict", "format 4 needs no upgrade")
			}
			if e = upgradeQueueRegistry(d, root, &r, p.RegistryUpgrade); e != nil {
				return e
			}
			r.EpicBaseUpdates = append(r.EpicBaseUpdates, EpicBaseUpdate{ID: id, PlanSHA256: in.Confirmation, Plan: p, Phase: "intent", RecordedAtUTC: queueNow(d)})
			sortQueueRegistry(&r)
			if e = session.Publish(r); e != nil {
				observeBasePublication(d, root, session, p.Target, id, &out)
				return fmt.Errorf("%w; inspect ply workspace epic base operation show %s", e, id)
			}
			for i := range r.EpicBaseUpdates {
				if r.EpicBaseUpdates[i].ID == id {
					return finishBaseUpdate(d, root, projects, session, &r, i, &out)
				}
			}
			return fmt.Errorf("missing published base request")
		})
	})
	return out, e
}
func finishBaseUpdate(d Dependencies, root string, projects ProjectSnapshot, session WorkItemStoreSession, r *WorkItemRegistry, index int, out *EpicBaseResult) error {
	u := r.EpicBaseUpdates[index]
	p := u.Plan
	ep, b := queueTargetBinding(*r, p.Target)
	previous := currentEpicBase(*r, ep, b)
	obs, e := queueParentObservation(d, p.Target)
	conflict := e != nil || !contentTypedEqual(obs, p.ObservedParent) || previous.Revision != p.ExpectedRevision || previous.OID != p.PreviousOID
	if !conflict {
		candidate := *r
		candidate.EpicBaseUpdates = []EpicBaseUpdate{}
		for _, v := range r.EpicBaseUpdates {
			if v.ID != u.ID {
				candidate.EpicBaseUpdates = append(candidate.EpicBaseUpdates, v)
			}
		}
		fresh, err := buildBasePlan(d, root, candidate, projects, EpicBaseInput{EpicID: p.Target.EpicID, RepoID: p.Target.RepoID})
		if err != nil || fresh.Plan == nil {
			conflict = true
		} else {
			fresh.Plan.RegistryUpgrade = p.RegistryUpgrade
			fresh.Plan.RegistryFormatVersion = p.RegistryFormatVersion
			conflict = !contentTypedEqual(*fresh.Plan, p)
		}
	}
	if conflict {
		u.Phase = "conflict"
		u.Outcome = &EpicBaseUpdateOutcome{Reasons: []QueueReason{{"epic_base_conflict", "target changed after intent publication"}}}
	} else {
		v := previous
		v.PreviousRevision = &previous.Revision
		v.Revision++
		v.OID = p.NextOID
		v.Tree = p.NextTree
		v.UpdateID = &u.ID
		r.EpicBaseVersions = append(r.EpicBaseVersions, v)
		u.Phase = "committed"
		u.Outcome = &EpicBaseUpdateOutcome{Revision: &v.Revision, Reasons: []QueueReason{}}
	}
	r.EpicBaseUpdates[index] = u
	sortQueueRegistry(r)
	out.Operation = &u
	out.State = "unknown"
	if e = session.Publish(*r); e != nil {
		observeBasePublication(d, root, session, p.Target, u.ID, out)
		return fmt.Errorf("%w; inspect ply workspace epic base operation show %s", e, u.ID)
	}
	actual, e := session.Snapshot()
	if e != nil {
		return e
	}
	*out = baseResult(d, root, actual, p.Target)
	out.Operation = &u
	out.State = u.Phase
	if p.RegistryUpgrade != nil {
		out.BackupSHA256 = &p.RegistryUpgrade.RegistrySHA256
	}
	if conflict {
		return queueError("epic_base_conflict", "preserved request closed without a new base; inspect a new preview")
	}
	return nil
}
func ShowEpicBaseOperation(d Dependencies, id string) (EpicBaseResult, error) {
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return EpicBaseResult{}, e
	}
	var out EpicBaseResult
	e = WithTaskSpecSnapshot(d, root, func(s *TaskSpecSession) error {
		for _, u := range s.Registry.EpicBaseUpdates {
			if u.ID == id {
				out = baseResult(d, root, s.Registry, u.Plan.Target)
				out.Operation = &u
				out.State = u.Phase
				if u.Phase == "intent" {
					out.State = "not_committed"
				}
				return nil
			}
		}
		return queueError("epic_base_operation_not_found", "base request is not recorded")
	})
	return out, e
}

func observeBasePublication(d Dependencies, root string, session WorkItemStoreSession, target QueueTarget, id string, out *EpicBaseResult) {
	r, e := session.Snapshot()
	if e != nil {
		out.State = "unknown"
		return
	}
	*out = baseResult(d, root, r, target)
	out.State = "not_recorded"
	out.PersistedAfter = &QueueRegistryReadback{r.FormatVersion, r.RawSHA256, queueUpgrade(r)}
	for _, u := range r.EpicBaseUpdates {
		if u.ID == id {
			out.Operation = &u
			out.State = u.Phase
			return
		}
	}
}
