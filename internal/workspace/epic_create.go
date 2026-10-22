package workspace

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
)

// EpicCreateInput selects an ordinary parent explicitly; it never infers one
// from the registered locator or represents it as an owned WorktreeID.
type EpicCreateInput struct {
	EpicID            EpicID    `yaml:"epic_id" json:"epic_id"`
	Title             string    `yaml:"title" json:"title"`
	ProjectID         ProjectID `yaml:"project_id" json:"project_id"`
	RepoID            RepoID    `yaml:"repo_id" json:"repo_id"`
	ParentRef         string    `yaml:"parent_ref" json:"parent_ref"`
	ExpectedParentOID string    `yaml:"expected_parent_oid" json:"expected_parent_oid"`
	Branch            string    `yaml:"branch" json:"branch"`
	Path              string    `yaml:"path" json:"path"`
}

type EpicCreateAttempt struct {
	Exit   int    `yaml:"exit" json:"exit"`
	Stdout string `yaml:"stdout" json:"stdout"`
	Stderr string `yaml:"stderr" json:"stderr"`
	Error  string `yaml:"error" json:"error"`
}

type EpicCreateOperation struct {
	Attempts      []EpicCreateAttempt         `yaml:"attempts" json:"attempts"`
	ID            WorktreeOperationID         `yaml:"id" json:"id"`
	WorktreeID    WorktreeID                  `yaml:"worktree_id" json:"worktree_id"`
	Request       EpicCreateInput             `yaml:"request" json:"request"`
	GitCommonDir  string                      `yaml:"git_common_dir" json:"git_common_dir"`
	ParentLocator string                      `yaml:"parent_locator" json:"parent_locator"`
	ParentTree    string                      `yaml:"parent_tree" json:"parent_tree"`
	IntentDigest  string                      `yaml:"intent_digest" json:"intent_digest"`
	State         string                      `yaml:"state" json:"state"`
	Observations  []WorktreeObservationRecord `yaml:"observations" json:"observations"`
}

type EpicCreateResult struct {
	Kind       string              `json:"kind"`
	Workspace  string              `json:"workspace"`
	Outcome    string              `json:"outcome"`
	Operation  EpicCreateOperation `json:"operation"`
	Epic       *EpicRecord         `json:"epic"`
	NextAction string              `json:"next_action"`
}

func validateEpicCreateInput(in EpicCreateInput) error {
	if _, err := ParseEpicAdoptInput(string(in.EpicID), in.Title, string(in.ProjectID), string(in.RepoID), in.Path, in.ParentRef, in.ExpectedParentOID); err != nil {
		return err
	}
	_, err := ParseTaskWorktreeCreateInput("epic-create", in.Branch, in.Path, in.ExpectedParentOID)
	return err
}
func epicCreateDigest(op EpicCreateOperation) string {
	op.State = ""
	op.Observations = nil
	op.Attempts = nil
	op.IntentDigest = ""
	return queueDigest(op)
}
func findEpicCreate(r WorkItemRegistry, id EpicID) (*EpicCreateOperation, int) {
	for i := range r.EpicCreateOperations {
		if r.EpicCreateOperations[i].Request.EpicID == id {
			return &r.EpicCreateOperations[i], i
		}
	}
	return nil, -1
}
func epicCreateReservation(r WorkItemRegistry, common, ref, path string, id EpicID) error {
	for _, op := range r.EpicCreateOperations {
		if op.Request.EpicID == id || op.GitCommonDir == common && (op.Request.ParentRef == ref || "refs/heads/"+op.Request.Branch == ref || pathOverlap(op.Request.Path, path) || op.ParentLocator == path) && op.State != "completed" {
			return workError(ErrorWorkIdentityConflict, fmt.Sprintf("Epic create %s reserves this identity, parent, branch or path; inspect `ply workspace epic show %s --format json` and retry its exact request", op.ID, op.Request.EpicID), nil)
		}
	}
	return nil
}
func pathOverlap(a, b string) bool { return pathWithin(a, b) || pathWithin(b, a) }

func CreateEpic(d Dependencies, in EpicCreateInput) (EpicCreateResult, error) {
	out := EpicCreateResult{Kind: "WorkspaceEpicCreate@1"}
	if e := validateEpicCreateInput(in); e != nil {
		return out, e
	}
	if e := requireWorkMutationDependencies(d); e != nil {
		return out, e
	}
	// Reuse the native Git adapter, with a bounded environment for every read and write.
	if g, ok := d.WorkGit.(*systemWorkItemGit); ok {
		copy := *g
		copy.run = runEpicCreateGit
		d.WorkGit = &copy
	}
	if e := d.WorkGit.ValidateBranch(in.Branch); e != nil {
		return out, e
	}
	root, e := containingWorkItemWorkspace(d)
	if e != nil {
		return out, e
	}
	out.Workspace = root
	cwd, e := projectPhysicalWorkingDirectory(d.Files)
	if e != nil {
		return out, e
	}
	target, exists, e := canonicalTaskTarget(d.Files, cwd, in.Path)
	if e != nil {
		return out, e
	}
	in.Path = target
	e = d.ProjectLocks.WithSnapshotLock(root, func(projects ProjectSnapshot) error {
		return d.WorkItems.WithLock(root, func(s WorkItemStoreSession) error {
			r, e := s.Snapshot()
			if e != nil {
				return e
			}
			if r.FormatVersion < 6 {
				return registryDiagnostic("registry_upgrade_required", "Epic creation requires format 6; ordinary existing operations remain available", r.FormatVersion)
			}
			op, index := findEpicCreate(r, in.EpicID)
			if op != nil {
				if !reflect.DeepEqual(op.Request, in) {
					return workError(ErrorWorkIdentityConflict, "Epic create request differs from the preserved intent; inspect epic show --format json", nil)
				}
				out.Operation = *op
				if op.State == "completed" {
					epic, _ := findEpic(r, in.EpicID)
					out.Epic = epic
					out.Outcome = "already_created"
					out.NextAction = "Inspect current readiness with ply workspace epic show; creation is historical and will not repair or reset Git."
					return nil
				}
				if op.State == "reconciliation_required" {
					return epicCreatePending(&out, *op, "Preserved partial/unknown effect requires inspection; no supported automatic retry until producer reconciliation.")
				}
			}
			_, repo, e := projectAndRepo(projects, in.ProjectID, in.RepoID)
			if e != nil {
				return e
			}
			observedRepo, e := d.WorkGit.ObserveRepo(repo)
			if e != nil {
				return e
			}
			if observedRepo.GitCommonDir != repo.GitCommonDir || observedRepo.Locator != repo.Locator {
				return workError(ErrorWorkProjectConflict, "registered repository changed; inspect Project binding", nil)
			}
			if g, ok := d.WorkGit.(*systemWorkItemGit); ok {
				if e = g.checkEpicMaterialization(repo, in.ExpectedParentOID); e != nil {
					return e
				}
			}
			if op == nil {
				if ep, _ := findEpic(r, in.EpicID); ep != nil {
					return workError(ErrorWorkIdentityConflict, "Epic ID already registered; use epic show or select a new ID", nil)
				}
				if exists {
					return workError(ErrorWorkPathConflict, "target already exists; select a vacant leaf without moving or deleting it", nil)
				}
				if e = epicCreateReservation(r, repo.GitCommonDir, "refs/heads/"+in.Branch, target, in.EpicID); e != nil {
					return e
				}
				parent, e := d.WorkGit.ObserveRef(repo, in.ParentRef)
				if e != nil {
					return e
				}
				if !parent.Exists || parent.OID != in.ExpectedParentOID || parent.CheckedOutAt == "" {
					return workError(ErrorWorkParentStale, "explicit parent must match expected OID and have one ordinary symbolic checkout; inspect Git worktree list", nil)
				}
				if e = epicCreateReservation(r, repo.GitCommonDir, in.ParentRef, parent.CheckedOutAt, in.EpicID); e != nil {
					return e
				}
				for _, ep := range r.Epics {
					for _, b := range ep.RepoBindings {
						if b.GitCommonDir == repo.GitCommonDir && (b.Worktree.Ref == in.ParentRef || b.Worktree.Locator == parent.CheckedOutAt || b.Worktree.Ref == "refs/heads/"+in.Branch || pathOverlap(b.Worktree.Locator, target)) {
							return workError(ErrorWorkIdentityConflict, "parent or target belongs to an Epic; select an ordinary parent and vacant target", nil)
						}
					}
				}
				for _, t := range r.WorktreeOperations {
					if t.GitCommonDir == repo.GitCommonDir && (t.SourceRef == in.ParentRef || t.SourceRef == "refs/heads/"+in.Branch || pathOverlap(t.TargetLocator, target) || t.TargetLocator == parent.CheckedOutAt) {
						return workError(ErrorWorkIdentityConflict, "parent or target reserved by a Task; inspect task show", nil)
					}
				}
				inventory, e := d.WorkGit.ListWorktrees(repo)
				if e != nil {
					return e
				}
				for _, entry := range inventory.Entries {
					if pathWithin(entry.Locator, target) {
						return workError(ErrorWorkPathConflict, "target is nested inside a Git worktree; select a sibling path", nil)
					}
				}
				if pathOverlap(repo.GitCommonDir, target) || pathOverlap(filepath.Join(root, MarkerDirectory), target) {
					return workError(ErrorWorkPathConflict, "target overlaps Git administration or .ply; select a vacant sibling path", nil)
				}
				id, e := d.WorkIDs.NewOperationID()
				if e != nil {
					return e
				}
				wid, e := d.WorkIDs.NewWorktreeID()
				if e != nil {
					return e
				}
				created := EpicCreateOperation{ID: id, WorktreeID: wid, Request: in, GitCommonDir: repo.GitCommonDir, ParentLocator: parent.CheckedOutAt, ParentTree: parent.Tree, State: "pending", Attempts: []EpicCreateAttempt{}, Observations: []WorktreeObservationRecord{}}
				created.IntentDigest = epicCreateDigest(created)
				class, obs := observeEpicCreate(d, repo, created)
				if class != "no_effect" {
					return workError(ErrorWorkParentStale, "parent or target is not exact and vacant; inspect parent status, branch and inventory before retrying", nil)
				}
				created.Observations = append(created.Observations, obs)
				r.EpicCreateOperations = append(r.EpicCreateOperations, created)
				sort.Slice(r.EpicCreateOperations, func(i, j int) bool { return r.EpicCreateOperations[i].ID < r.EpicCreateOperations[j].ID })
				if e = s.Publish(r); e != nil {
					return e
				}
				op, index = findEpicCreate(r, in.EpicID)
			}
			class, obs := observeEpicCreate(d, repo, *op)
			if class == "no_effect" || class == "branch_only" {
				input := GitCreateInput{Repository: repo, Branch: in.Branch, SourceRef: "refs/heads/" + in.Branch, TargetLocator: target, ParentOID: in.ExpectedParentOID}
				var outcome GitCommandOutcome
				if class == "branch_only" {
					outcome = d.WorkGit.AddWorktreeForExistingBranch(input)
				} else {
					outcome = d.WorkGit.CreateBranchAndWorktree(input)
				}
				attempt := EpicCreateAttempt{Exit: outcome.Exit, Stdout: string(outcome.Stdout), Stderr: string(outcome.Stderr)}
				if outcome.Err != nil {
					attempt.Error = outcome.Err.Error()
				}
				r.EpicCreateOperations[index].Attempts = append(r.EpicCreateOperations[index].Attempts, attempt)
				// The actual effect, including a lost/nonzero response, determines recovery.
				class, obs = observeEpicCreate(d, repo, *op)
			}
			r.EpicCreateOperations[index].Observations = append(r.EpicCreateOperations[index].Observations, obs)
			switch class {
			case "exact_effect":
				epic := EpicRecord{ID: in.EpicID, Title: in.Title, ProjectID: in.ProjectID, RepoBindings: []EpicRepoBinding{{RepoID: in.RepoID, GitCommonDir: op.GitCommonDir, Worktree: EpicWorktreeBinding{ID: op.WorktreeID, OwnerKind: "epic", OwnerID: in.EpicID, Origin: "ply_created", Locator: target, Ref: "refs/heads/" + in.Branch, OID: in.ExpectedParentOID, Tree: op.ParentTree, GitCommonDir: op.GitCommonDir}}}}
				r.Epics = append(r.Epics, epic)
				r.EpicBaseVersions = append(r.EpicBaseVersions, adoptedEpicBase(epic, epic.RepoBindings[0]))
				r.EpicCreateOperations[index].State = "completed"
				sortWorkRegistry(&r)
				sortQueueRegistry(&r)
				if e = s.Publish(r); e != nil {
					return e
				}
				out.Operation = r.EpicCreateOperations[index]
				out.Epic = &epic
				out.Outcome = "created"
				return nil
			case "no_effect", "branch_only":
				if e = s.Publish(r); e != nil {
					return e
				}
				return epicCreatePending(&out, r.EpicCreateOperations[index], "Inspect epic show --format json, then repeat the identical create request to resume the observed no-effect or branch-only state.")
			default:
				r.EpicCreateOperations[index].State = "reconciliation_required"
				if e = s.Publish(r); e != nil {
					return e
				}
				return epicCreatePending(&out, r.EpicCreateOperations[index], "Inspect epic show --format json; preserve the partial/unknown effect for producer reconciliation without cleanup or blind retry.")
			}
		})
	})
	return out, e
}
func epicCreatePending(out *EpicCreateResult, op EpicCreateOperation, next string) error {
	out.Operation = op
	out.Outcome = op.State
	out.NextAction = next
	return workError(ErrorWorkReconciliation, next, nil)
}
func observeEpicCreate(d Dependencies, repo RepoRecord, op EpicCreateOperation) (string, WorktreeObservationRecord) {
	surrogate := WorktreeOperationRecord{ParentRef: op.Request.ParentRef, ParentOID: op.Request.ExpectedParentOID, ParentTree: op.ParentTree, GitCommonDir: op.GitCommonDir, SourceRef: "refs/heads/" + op.Request.Branch, TargetLocator: op.Request.Path}
	binding := EpicRepoBinding{Worktree: EpicWorktreeBinding{Locator: op.ParentLocator}}
	class, obs := observeTaskCreate(d, repo, binding, surrogate)
	if g, ok := d.WorkGit.(*systemWorkItemGit); ok {
		if e := g.checkEpicParentOperation(op.ParentLocator); e != nil {
			obs.Classification = "partial_or_unknown"
			obs.Reasons = append(obs.Reasons, "parent_observation_unknown")
			sort.Strings(obs.Reasons)
			obs.Reasons = uniqueReasonStrings(obs.Reasons)
			return "partial_or_unknown", obs
		}
	}
	return class, obs
}

func validateEpicCreates(r WorkItemRegistry) error {
	if r.FormatVersion < 6 {
		if len(r.EpicCreateOperations) > 0 || r.EpicCreateBoundary != nil {
			return fmt.Errorf("legacy format cannot contain Epic creation data")
		}
		return nil
	}
	seen := map[EpicID]bool{}
	ids := map[WorktreeOperationID]bool{}
	wids := map[WorktreeID]bool{}
	previous := WorktreeOperationID("")
	for _, op := range r.EpicCreateOperations {
		if validateEpicCreateInput(op.Request) != nil || !operationIDPattern.MatchString(string(op.ID)) || !worktreeIDPattern.MatchString(string(op.WorktreeID)) || seen[op.Request.EpicID] || ids[op.ID] || wids[op.WorktreeID] || op.ID <= previous || !validStoredPath(op.Request.Path) || !validStoredPath(op.ParentLocator) || !validStoredPath(op.GitCommonDir) || !validOIDText(op.ParentTree) || epicCreateDigest(op) != op.IntentDigest || len(op.Observations) == 0 {
			return fmt.Errorf("invalid Epic create intent %s", op.ID)
		}
		for _, obs := range op.Observations {
			state := "creating"
			if obs.Classification == "exact_effect" {
				state = "ready"
			} else if obs.Classification == "partial_or_unknown" {
				state = "reconciliation_required"
			}
			if e := validateObservation(obs, state); e != nil {
				return e
			}
		}
		if op.Attempts == nil {
			return fmt.Errorf("Epic create requires explicit attempts")
		}
		previous = op.ID
		seen[op.Request.EpicID] = true
		ids[op.ID] = true
		wids[op.WorktreeID] = true
		if op.State != "pending" && op.State != "completed" && op.State != "reconciliation_required" {
			return fmt.Errorf("invalid Epic create state")
		}
		for _, t := range r.WorktreeOperations {
			if t.ID == op.ID || t.WorktreeID == op.WorktreeID || t.GitCommonDir == op.GitCommonDir && (t.SourceRef == "refs/heads/"+op.Request.Branch || pathOverlap(t.TargetLocator, op.Request.Path)) {
				return fmt.Errorf("Epic create conflicts with Task intent")
			}
		}
		epic, _ := findEpic(r, op.Request.EpicID)
		if op.State == "completed" {
			if epic == nil || len(epic.RepoBindings) != 1 {
				return fmt.Errorf("completed create lacks exact Epic")
			}
			b := epic.RepoBindings[0]
			w := b.Worktree
			if epic.Title != op.Request.Title || epic.ProjectID != op.Request.ProjectID || b.RepoID != op.Request.RepoID || w.ID != op.WorktreeID || w.Origin != "ply_created" || w.Locator != op.Request.Path || w.Ref != "refs/heads/"+op.Request.Branch || w.OID != op.Request.ExpectedParentOID || w.Tree != op.ParentTree || w.GitCommonDir != op.GitCommonDir {
				return fmt.Errorf("completed create differs from Epic provenance")
			}
		} else if epic != nil {
			return fmt.Errorf("pending create cannot have an Epic")
		}
	}
	for _, ep := range r.Epics {
		for _, b := range ep.RepoBindings {
			if b.Worktree.Origin == "ply_created" && !seen[ep.ID] {
				return fmt.Errorf("created Epic lacks operation")
			}
		}
	}
	return nil
}
func validateEpicCreateTransition(old, next WorkItemRegistry) error {
	for _, a := range old.EpicCreateOperations {
		b, _ := findEpicCreate(next, a.Request.EpicID)
		if b == nil || a.IntentDigest != b.IntentDigest || epicCreateDigest(*b) != a.IntentDigest || len(b.Observations) < len(a.Observations) || !reflect.DeepEqual(a.Observations, b.Observations[:len(a.Observations)]) || len(b.Attempts) < len(a.Attempts) || !reflect.DeepEqual(a.Attempts, b.Attempts[:len(a.Attempts)]) {
			return fmt.Errorf("Epic creation intent and observations are immutable")
		}
		if a.State == "completed" && !reflect.DeepEqual(a, *b) || a.State == "reconciliation_required" && !reflect.DeepEqual(a, *b) {
			return fmt.Errorf("terminal Epic create operation is immutable")
		}
	}
	return nil
}
