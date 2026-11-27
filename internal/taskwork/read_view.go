package taskwork

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/repohistory"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// ReadViewBasis is a request-local registration snapshot. It deliberately does
// not validate every historical content object or observe Git. Detailed reads
// and mutations retain their stronger, separate validation contracts.
type ReadViewBasisSnapshot struct {
	Workspace  workspace.WorkspaceObservation
	Projects   workspace.ProjectSnapshot
	Registries []ProjectRegistry // each Project's registry; content is read through it
	Registry   WorkItemRegistry  // all registries merged, for listing
	Titles     map[TaskID]TaskListTitle
	Sources    map[string]string
}

// TaskJournal is the journal basis of a Task in its own Project's registry.
func (b ReadViewBasisSnapshot) TaskJournal(task TaskRecord) TaskJournalBasis {
	owner, _ := registryOf(b.Registries, task.ProjectID)
	return TaskJournalBasis{Workspace: b.Workspace, Task: task, Registry: owner.Registry, Work: owner.Work}
}

func ReadViewBasis(d Dependencies) (ReadViewBasisSnapshot, error) {
	var out ReadViewBasisSnapshot
	root, err := containingRegistrationWorkspace(d)
	if err != nil {
		return out, err
	}
	ws, err := workspace.ObserveWorkspaceRoot(d.Files, root)
	if err != nil {
		return out, err
	}
	if d.WorkItems == nil || d.Projects == nil {
		return out, workspace.WorkError(workspace.ErrorWorkIO, "work-item and project readers are required", nil)
	}
	out.Workspace, out.Sources = ws, map[string]string{}
	definitionsPath := filepath.Join(ws.Root, workspace.DefinitionsFile)
	out.Sources[definitionsPath], err = readViewSource(d, definitionsPath)
	if err != nil {
		return out, err
	}
	out.Registries, err = SnapshotWorkspace(d, ws.Root)
	if err != nil {
		return out, err
	}
	out.Registry = MergeRegistries(out.Registries)
	lifecycles, err := ReadWorkItemLifecycle(d, ws.Root)
	if err != nil {
		return out, err
	}
	projects, repos, err := registryProjects(ws.Root, out.Registry, lifecycles)
	if err != nil {
		return out, err
	}
	out.Projects = workspace.ProjectSnapshot{Projects: projects, Repos: repos}
	out.Sources[workspace.MarkerPath(ws.Root)] = ws.MarkerSHA256
	out.Titles = make(map[TaskID]TaskListTitle, len(out.Registry.Tasks))
	for _, owner := range out.Registries {
		if owner.Registry.RawSHA256 == "" {
			out.Sources[owner.Work.Registry()] = "absent"
		} else {
			out.Sources[owner.Work.Registry()] = owner.Registry.RawSHA256
		}
		for _, source := range IntegrationEvidenceSources(owner.Work, owner.Registry) {
			out.Sources[source.Locator] = source.SHA256
		}
		for _, task := range owner.Registry.Tasks {
			title := readTaskListTitle(d, owner.Work, owner.Registry, task)
			out.Titles[task.ID] = title
			if title.Source == "problem_revision" && title.Status == "available" {
				if head := taskContentState(owner.Registry, task.ID).ProblemHead; head != nil {
					out.Sources[taskContentPath(owner.Work, "manifests", head.ManifestSHA256)] = head.ManifestSHA256
				}
			}
		}
	}
	return out, nil
}

// CheckReadViewBasis detects register drift without decoding the registers a
// second time, taking locks, or claiming a transaction across independent files.
func CheckReadViewBasis(d Dependencies, b ReadViewBasisSnapshot) (string, []string) {
	for _, owner := range b.Registries {
		if err := newSystemWorkItemStore().verifyIntegrationEvidence(owner.Work, owner.Registry); err != nil {
			return "unknown", []string{"registration_sources_changed_or_unavailable"}
		}
	}
	reasons := []string{}
	for path, expected := range b.Sources {
		actual, err := readViewSource(d, path)
		if err != nil || actual != expected {
			reasons = append(reasons, "registration_sources_changed_or_unavailable")
		}
	}
	if len(reasons) != 0 {
		return "unknown", sortedReasons(reasons)
	}
	return "fresh", reasons
}

func readViewSource(d Dependencies, path string) (string, error) {
	info, err := d.Files.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "absent", nil
	}
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return "", workspace.WorkError(workspace.ErrorWorkStoreConflict, "view source must remain a regular file: "+path, nil)
	}
	b, err := d.Files.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "absent", nil
	}
	if err != nil {
		return "", err
	}
	return digestTaskBytes(b), nil
}

func readTaskListTitle(d Dependencies, work workspace.ProjectWork, registry WorkItemRegistry, task TaskRecord) TaskListTitle {
	title := task.Title
	result := TaskListTitle{Title: &title, Source: "registration", Status: "available"}
	head := taskContentState(registry, task.ID).ProblemHead
	if head == nil {
		return result
	}
	result = TaskListTitle{Source: "problem_revision", Status: "unavailable"}
	m, err := readRegisteredTaskManifest(d.TaskContent, work, registry, head.ManifestSHA256)
	if err != nil {
		return result
	}
	for _, publication := range registry.TaskContentPublications {
		if publication.OutcomeRef.ManifestSHA256 == head.ManifestSHA256 {
			fields := contentFields(m)
			if contentString(fields, "task_id") != string(task.ID) || contentString(fields, "publication_key") != publication.PublicationKey || contentString(fields, "source_draft_sha256") != publication.IntentSHA256 {
				return result
			}
			break
		}
	}
	title = contentString(contentFields(m), "title")
	result.Title, result.Status = &title, "available"
	return result
}

// RegisteredEpicBases returns the latest recorded base per repository. It never
// claims that the checkout or branch still matches this historical registration.
func RegisteredEpicBases(r WorkItemRegistry, epic EpicRecord) []EpicBaseVersion {
	bases := make([]EpicBaseVersion, 0, len(epic.RepoBindings))
	for _, binding := range epic.RepoBindings {
		bases = append(bases, currentEpicBase(r, epic, binding))
	}
	return bases
}

// RegisteredEpicActivityTimes includes native Epic-base and queue events whose
// timestamps and ownership are carried by the validated registry itself.
func RegisteredEpicActivityTimes(r WorkItemRegistry) map[EpicID]*string {
	out := map[EpicID]*string{}
	for _, update := range r.EpicBaseUpdates {
		id := update.Plan.Target.EpicID
		out[id] = recordedLatestTime(out[id], update.RecordedAtUTC)
	}
	for _, event := range r.TaskQueueEvents {
		fields := contentFields(queueRequestValue(event.Request))
		id := EpicID(contentString(fields, "epic_id"))
		if id != "" {
			out[id] = recordedLatestTime(out[id], event.RecordedAtUTC)
		}
	}
	return out
}

// registryProjects resolves the defined Projects statically, with the Git
// common directories the registry recorded for their repositories.
func registryProjects(root string, registry WorkItemRegistry, lifecycles WorkItemLifecycles) ([]workspace.ProjectRecord, []workspace.RepoRecord, error) {
	defs, _, err := workspace.ReadDefinitions(root)
	if err != nil {
		return nil, nil, err
	}
	if defs == nil {
		return nil, nil, workspace.MissingDefinitions(root)
	}
	identities, err := repohistory.CommonDirs(registryBindings(registry, lifecycles))
	if err != nil {
		return nil, nil, workspace.WorkInvalidArguments("conflicting historical repository identity")
	}
	projects := []workspace.ProjectRecord{}
	repos := []workspace.RepoRecord{}
	seen := map[workspace.RepoID]bool{}
	path := func(value string) string {
		if filepath.IsAbs(value) {
			return filepath.Clean(value)
		}
		return filepath.Join(root, value)
	}
	for _, p := range defs.Projects {
		project := workspace.ProjectRecord{ID: workspace.ProjectID(p.ID), Name: p.Name, Wrapper: path(p.Wrapper)}
		for _, r := range p.Repos {
			project.RepoIDs = append(project.RepoIDs, workspace.RepoID(r.ID))
			if !seen[workspace.RepoID(r.ID)] {
				repos = append(repos, workspace.RepoRecord{ID: workspace.RepoID(r.ID), Locator: path(r.Path), GitCommonDir: identities[r.ID]})
				seen[workspace.RepoID(r.ID)] = true
			}
		}
		projects = append(projects, project)
	}
	return projects, repos, nil
}

// registryBindings projects the repository identities that the live Epic and
// Task records of registry carry; finished work is history and names none.
func registryBindings(registry WorkItemRegistry, lifecycles WorkItemLifecycles) []repohistory.Binding {
	var bindings []repohistory.Binding
	for _, use := range repoUses(registry, lifecycles) {
		if use.Live && (use.Kind == workspace.RepoUseEpic || use.Kind == workspace.RepoUseTask) {
			bindings = append(bindings, repohistory.Binding{RepoID: string(use.RepoID), GitCommonDir: use.GitCommonDir})
		}
	}
	return bindings
}
