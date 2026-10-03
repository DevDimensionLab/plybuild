package workspace

import "path/filepath"

// TaskJournalBasis is a read-only projection. Optional content is validated
// separately so one unavailable document cannot hide independent registry facts.
type TaskJournalBasis struct {
	Workspace       WorkspaceObservation
	Task            TaskRecord
	Registry        WorkItemRegistry
	RegistryLocator string
}

func ReadTaskJournalBasis(d Dependencies, id TaskID) (TaskJournalBasis, error) {
	var out TaskJournalBasis
	if _, err := ParseTaskID(string(id)); err != nil {
		return out, err
	}
	ws, err := ObserveContaining(d)
	if err != nil {
		return out, err
	}
	p := filepath.Join(ws.Root, MarkerDirectory, WorkItemsFile)
	b, err := contentReadBounded(p, 32<<20, false)
	if err != nil {
		return out, err
	}
	r, err := decodeWorkItemRegistry(b)
	if err != nil {
		return out, err
	}
	t, _ := findTask(r, id)
	if t == nil {
		return out, workError(ErrorWorkNotFound, "Task is not registered", nil)
	}
	projects, repos, err := d.Projects.Snapshot(ws.Root)
	if err != nil {
		return out, err
	}
	_, repo, err := projectAndRepo(ProjectSnapshot{Projects: projects, Repos: repos}, t.ProjectID, t.RepoID)
	if err != nil {
		return out, err
	}
	if repo.GitCommonDir != t.GitCommonDir {
		return out, workError(ErrorWorkIdentityConflict, "Task repository binding differs", nil)
	}
	r.RawSHA256 = digestTaskBytes(b)
	return TaskJournalBasis{ws, *t, r, p}, nil
}

type TaskJournalSource struct{ Kind, Locator, SHA256 string }

// ReadTaskJournalContent exposes the versions read by the native validators,
// including document bytes. It never repairs or copies the content store.
func ReadTaskJournalContent(d Dependencies, b TaskJournalBasis, p TaskContentPublication) ([]TaskJournalSource, error) {
	root := b.Workspace.Root
	sources := []TaskJournalSource{
		{"task_content", taskContentPath(root, "manifests", p.OutcomeRef.ManifestSHA256), p.OutcomeRef.ManifestSHA256},
		{"task_content_request", taskContentPath(root, "requests", p.RequestSHA256), p.RequestSHA256},
	}
	if p.BackupSHA256 != nil {
		sources = append(sources, TaskJournalSource{"task_content_backup", taskContentPath(root, "backups", *p.BackupSHA256), *p.BackupSHA256})
	}
	o, err := readRegisteredTaskManifest(d.TaskContent, root, b.Registry, p.OutcomeRef.ManifestSHA256)
	if err != nil {
		return sources, err
	}
	for _, v := range contentArray(contentFields(o), "documents") {
		m := contentFields(v)
		digest := contentString(m, "sha256")
		sources = append(sources, TaskJournalSource{"task_content_document", taskContentPath(root, "objects", digest), digest})
	}
	return sources, verifyManifestDocuments(d.TaskContent, root, o)
}

func ValidateTaskJournalPreparation(d Dependencies, b TaskJournalBasis, id string) (TaskPreparation, error) {
	if err := validateQueueClosure(d, b.Workspace.Root, b.Registry); err != nil {
		return TaskPreparation{}, err
	}
	p := findPreparation(b.Registry, id)
	if p == nil || p.Plan.TaskID != b.Task.ID {
		return TaskPreparation{}, workError(ErrorWorkIdentityConflict, "preparation does not belong to Task", nil)
	}
	return *p, nil
}

// Use the established Git observer with process-local settings that disable
// executable hooks and optional writes, including a configured fsmonitor hook.
func ObserveTaskJournalWorktree(d Dependencies, path string) (GitWorktreeObservation, error) {
	observer := d.WorkGit
	if g, ok := observer.(*systemWorkItemGit); ok {
		copy := *g
		copy.run = func(args, env []string) GitCommandOutcome {
			prefix := []string{"-c", "core.hooksPath=/dev/null", "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false", "-c", "gc.auto=0", "-c", "maintenance.auto=false", "-c", "submodule.recurse=false"}
			return g.run(append(prefix, args...), env)
		}
		observer = &copy
	}
	if observer == nil {
		return GitWorktreeObservation{}, workError(ErrorWorkGitObservation, "Git observer unavailable", nil)
	}
	return observer.ObserveWorktree(path)
}

func ValidateTaskJournalCandidate(d Dependencies, b TaskJournalBasis, oid, tree string) error {
	projects, repos, err := d.Projects.Snapshot(b.Workspace.Root)
	if err != nil {
		return err
	}
	_, repo, err := projectAndRepo(ProjectSnapshot{Projects: projects, Repos: repos}, b.Task.ProjectID, b.Task.RepoID)
	if err != nil {
		return err
	}
	g, ok := d.WorkGit.(TaskContentGit)
	if !ok {
		return workError(ErrorWorkGitObservation, "candidate Git observer unavailable", nil)
	}
	actual, err := g.ObserveContentCommit(repo, oid)
	if err != nil {
		return err
	}
	if actual != tree {
		return workError(ErrorWorkIdentityConflict, "candidate commit and tree differ", nil)
	}
	return nil
}
