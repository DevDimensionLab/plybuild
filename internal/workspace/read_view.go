package workspace

import (
	"errors"
	"io/fs"
	"path/filepath"
)

// ReadViewBasis is a request-local registration snapshot. It deliberately does
// not validate every historical content object or observe Git. Detailed reads
// and mutations retain their stronger, separate validation contracts.
type ReadViewBasisSnapshot struct {
	Workspace WorkspaceObservation
	Projects  ProjectSnapshot
	Registry  WorkItemRegistry
	Titles    map[TaskID]TaskListTitle
	Sources   map[string]string
}

func ReadViewBasis(d Dependencies) (ReadViewBasisSnapshot, error) {
	var out ReadViewBasisSnapshot
	ws, err := ObserveContaining(d)
	if err != nil {
		return out, err
	}
	if d.WorkItems == nil || d.Projects == nil {
		return out, workError(ErrorWorkIO, "work-item and project readers are required", nil)
	}
	out.Workspace, out.Sources = ws, map[string]string{}
	projectsPath := filepath.Join(ws.Root, MarkerDirectory, ProjectsFile)
	out.Sources[projectsPath], err = readViewSource(d, projectsPath)
	if err != nil {
		return out, err
	}
	projects, repos, err := d.Projects.Snapshot(ws.Root)
	if err != nil {
		return out, err
	}
	out.Projects = ProjectSnapshot{Projects: projects, Repos: repos}
	out.Registry, err = d.WorkItems.SnapshotRegistrations(ws.Root)
	if err != nil {
		return out, err
	}
	itemsPath := filepath.Join(ws.Root, MarkerDirectory, WorkItemsFile)
	if out.Registry.RawSHA256 == "" {
		out.Sources[itemsPath] = "absent"
	} else {
		out.Sources[itemsPath] = out.Registry.RawSHA256
	}
	out.Sources[filepath.Join(ws.Root, MarkerDirectory, MarkerFile)] = ws.MarkerSHA256
	out.Titles = make(map[TaskID]TaskListTitle, len(out.Registry.Tasks))
	for _, task := range out.Registry.Tasks {
		title := readTaskListTitle(d, ws.Root, out.Registry, task)
		out.Titles[task.ID] = title
		if title.Source == "problem_revision" && title.Status == "available" {
			if head := taskContentState(out.Registry, task.ID).ProblemHead; head != nil {
				out.Sources[taskContentPath(ws.Root, "manifests", head.ManifestSHA256)] = head.ManifestSHA256
			}
		}
	}
	return out, nil
}

// CheckReadViewBasis detects register drift without decoding the registers a
// second time, taking locks, or claiming a transaction across independent files.
func CheckReadViewBasis(d Dependencies, b ReadViewBasisSnapshot) (string, []string) {
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
		return "", workError(ErrorWorkStoreConflict, "view source must remain a regular file: "+path, nil)
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

func readTaskListTitle(d Dependencies, root string, registry WorkItemRegistry, task TaskRecord) TaskListTitle {
	title := task.Title
	result := TaskListTitle{Title: &title, Source: "registration", Status: "available"}
	head := taskContentState(registry, task.ID).ProblemHead
	if head == nil {
		return result
	}
	result = TaskListTitle{Source: "problem_revision", Status: "unavailable"}
	m, err := readRegisteredTaskManifest(d.TaskContent, root, registry, head.ManifestSHA256)
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
		if id == "" {
			if preparation := findPreparation(r, contentString(fields, "preparation_id")); preparation != nil {
				id = preparation.Plan.Target.EpicID
			}
		}
		if id != "" {
			out[id] = recordedLatestTime(out[id], event.RecordedAtUTC)
		}
	}
	return out
}
