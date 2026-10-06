package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type CompanionRole string

const (
	CompanionPlanning           CompanionRole = "planning"
	CompanionDocs               CompanionRole = "docs"
	CompanionOther              CompanionRole = "other"
	ProjectCompanionAdd                       = "companion_add"
	ProjectCompanionRemove                    = "companion_remove"
	ProjectRepoWrapperSet                     = "repo_wrapper_set"
	ProjectRepoWrapperClear                   = "repo_wrapper_clear"
	ProjectMetadataFile                       = "project-metadata.json"
	ProjectMetadataReadbackKind               = "WorkspaceProjectMetadataReadback@1"
	ProjectMetadataMutationKind               = "WorkspaceProjectMetadataMutation@1"
)

type ProjectCompanion struct {
	Role    CompanionRole `json:"role"`
	Locator string        `json:"locator"`
}
type ProjectRepositoryMetadata struct {
	RepoID  RepoID  `json:"repo_id"`
	Wrapper *string `json:"wrapper"`
}
type ProjectMetadataEvent struct {
	ID            string         `json:"event_id"`
	Sequence      int            `json:"sequence"`
	Operation     string         `json:"operation"`
	ProjectID     ProjectID      `json:"project_id"`
	RepoID        *RepoID        `json:"repo_id"`
	Role          *CompanionRole `json:"role"`
	Locator       *string        `json:"locator"`
	ActorClaim    string         `json:"actor_claim"`
	RecordedAtUTC string         `json:"recorded_at_utc"`
}
type ProjectMetadataSnapshot struct {
	Revision   int
	Events     []ProjectMetadataEvent
	Locator    string
	SHA256     string
	Exists     bool
	companions map[ProjectID]map[string]ProjectCompanion
	wrappers   map[RepoID]string
}

func (s ProjectMetadataSnapshot) Companions(id ProjectID) []ProjectCompanion {
	items := []ProjectCompanion{}
	for _, companion := range s.companions[id] {
		items = append(items, companion)
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Role != items[j].Role {
			return items[i].Role < items[j].Role
		}
		return items[i].Locator < items[j].Locator
	})
	return items
}
func (s ProjectMetadataSnapshot) RepositoryWrapper(id RepoID) *string {
	if wrapper, ok := s.wrappers[id]; ok {
		return &wrapper
	}
	return nil
}

type ProjectMetadataInput struct {
	Operation        string
	ProjectID        ProjectID
	RepoID           RepoID
	Role             CompanionRole
	Locator          string
	ActorClaim       string
	ExpectedRevision *int
}

func (in ProjectMetadataInput) Validate() error {
	if _, err := ParseProjectID(string(in.ProjectID)); err != nil {
		return err
	}
	if !validTaskText(in.ActorClaim, 1, 256) {
		return ProjectInvalidArguments("--actor must contain 1..256 characters without surrounding whitespace or control characters")
	}
	if in.ExpectedRevision != nil && *in.ExpectedRevision < 0 {
		return ProjectInvalidArguments("--expected-revision must be zero or greater")
	}
	switch in.Operation {
	case ProjectCompanionAdd, ProjectCompanionRemove:
		if in.Role != CompanionPlanning && in.Role != CompanionDocs && in.Role != CompanionOther {
			return ProjectInvalidArguments("--role must be planning, docs, or other")
		}
		if in.RepoID != "" {
			return ProjectInvalidArguments("companion changes do not take a repository ID")
		}
	case ProjectRepoWrapperSet, ProjectRepoWrapperClear:
		if err := validateIdentifier("repository", string(in.RepoID)); err != nil {
			return err
		}
		if in.Role != "" {
			return ProjectInvalidArguments("repository wrapper changes do not take a companion role")
		}
	default:
		return ProjectInvalidArguments("unknown project metadata operation")
	}
	if in.Operation == ProjectRepoWrapperClear {
		if in.Locator != "" {
			return ProjectInvalidArguments("clearing a repository wrapper does not take a path")
		}
	} else if in.Locator == "" || !utf8.ValidString(in.Locator) || strings.IndexByte(in.Locator, 0) >= 0 {
		return ProjectInvalidArguments("metadata path must be nonempty valid UTF-8 and contain no NUL")
	}
	return nil
}

type ProjectMetadataReadback struct {
	Kind          string                      `json:"kind"`
	SchemaVersion int                         `json:"schema_version"`
	Workspace     MetadataWorkspace           `json:"workspace"`
	ProjectID     ProjectID                   `json:"project_id"`
	Revision      int                         `json:"revision"`
	Companions    []ProjectCompanion          `json:"companions"`
	Repositories  []ProjectRepositoryMetadata `json:"repositories"`
	Events        []ProjectMetadataEvent      `json:"events"`
	Source        MetadataSource              `json:"source"`
}
type ProjectMetadataMutation struct {
	ProjectMetadataReadback
	Changed bool                  `json:"changed"`
	Event   *ProjectMetadataEvent `json:"event"`
}

func projectMetadataReadback(root string, project ProjectRecord, snapshot ProjectMetadataSnapshot) ProjectMetadataReadback {
	repos := []ProjectRepositoryMetadata{}
	repoIDs := map[RepoID]bool{}
	for _, id := range project.RepoIDs {
		repos = append(repos, ProjectRepositoryMetadata{id, snapshot.RepositoryWrapper(id)})
		repoIDs[id] = true
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].RepoID < repos[j].RepoID })
	events := []ProjectMetadataEvent{}
	for _, event := range snapshot.Events {
		if event.ProjectID == project.ID || event.RepoID != nil && repoIDs[*event.RepoID] {
			events = append(events, event)
		}
	}
	return ProjectMetadataReadback{ProjectMetadataReadbackKind, 1, MetadataWorkspace{root}, project.ID, snapshot.Revision, snapshot.Companions(project.ID), repos, events, metadataSource(snapshot.Locator, snapshot.SHA256, snapshot.Exists)}
}

// ReadProjectMetadata returns explicit registrations without checking whether the
// registered directories still exist. No filesystem scanning or write occurs.
func ReadProjectMetadata(d Dependencies, root string) (ProjectMetadataSnapshot, error) {
	observed, err := ObserveRoot(d, root)
	if err != nil {
		return ProjectMetadataSnapshot{}, err
	}
	return readProjectMetadata(observed.Root)
}
func ShowProjectMetadata(d Dependencies, id ProjectID) (ProjectMetadataReadback, error) {
	project, err := ShowProject(d, id)
	if err != nil {
		return ProjectMetadataReadback{}, err
	}
	return projectMetadataReadback(project.Workspace, project.Project, project.Metadata), nil
}

func UpdateProjectMetadata(d Dependencies, in ProjectMetadataInput) (ProjectMetadataMutation, error) {
	return updateProjectMetadata(d, in, workspaceMetadataPublisher{})
}
func updateProjectMetadata(d Dependencies, in ProjectMetadataInput, publisher workspaceMetadataPublisher) (ProjectMetadataMutation, error) {
	var result ProjectMetadataMutation
	if err := in.Validate(); err != nil {
		return result, err
	}
	if d.ProjectLocks == nil || d.WorkClock == nil {
		return result, projectError(ErrorProjectIO, "project lock and clock dependencies are required", nil)
	}
	ws, err := ObserveContaining(d)
	if err != nil {
		return result, err
	}
	cwd, err := projectPhysicalWorkingDirectory(d.Files)
	if err != nil {
		return result, err
	}
	err = d.ProjectLocks.WithSnapshotLock(ws.Root, func(registry ProjectSnapshot) error {
		var project *ProjectRecord
		for _, record := range registry.Projects {
			if record.ID == in.ProjectID {
				copy := record
				project = &copy
				break
			}
		}
		if project == nil {
			return projectError(ErrorProjectNotFound, fmt.Sprintf("project %s is not registered", in.ProjectID), nil)
		}
		if in.Locator != "" {
			path, err := projectMetadataPath(d.Files, cwd, in.Locator, in.Operation == ProjectCompanionRemove)
			if err != nil {
				return err
			}
			if !utf8.ValidString(path) {
				return ProjectInvalidArguments("physical metadata path must be valid UTF-8")
			}
			in.Locator = path
		}
		if in.RepoID != "" {
			_, repo, err := projectAndRepo(registry, in.ProjectID, in.RepoID)
			if err != nil {
				return err
			}
			if in.Operation == ProjectRepoWrapperSet {
				rel, err := filepath.Rel(in.Locator, repo.Locator)
				if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
					return ProjectInvalidArguments("repository wrapper must contain the registered repository locator")
				}
			}
		}
		snapshot, err := readProjectMetadata(ws.Root)
		if err != nil {
			return err
		}
		result.ProjectMetadataReadback = projectMetadataReadback(ws.Root, *project, snapshot)
		result.Kind = ProjectMetadataMutationKind
		if in.ExpectedRevision != nil && *in.ExpectedRevision != snapshot.Revision {
			return projectError(ErrorProjectConflict, fmt.Sprintf("project metadata revision is %d, expected %d; read metadata before retrying", snapshot.Revision, *in.ExpectedRevision), nil)
		}
		event := ProjectMetadataEvent{Sequence: snapshot.Revision + 1, Operation: in.Operation, ProjectID: in.ProjectID, ActorClaim: in.ActorClaim, RecordedAtUTC: d.WorkClock.Now().UTC().Format(time.RFC3339Nano)}
		if in.RepoID != "" {
			event.RepoID = &in.RepoID
		}
		if in.Role != "" {
			event.Role = &in.Role
		}
		if in.Locator != "" {
			event.Locator = &in.Locator
		}
		event.ID = projectMetadataEventID(event)
		changed, err := applyProjectMetadataEvent(&snapshot, event)
		if err != nil {
			return err
		}
		if !changed {
			return nil
		}
		events := append(append([]ProjectMetadataEvent{}, snapshot.Events...), event)
		bytes, err := encodeProjectMetadata(events)
		if err != nil {
			return err
		}
		if err := publisher.publish(ws.Root, ProjectMetadataFile, bytes); err != nil {
			return err
		}
		snapshot, err = readProjectMetadata(ws.Root)
		if err != nil {
			return err
		}
		result.ProjectMetadataReadback = projectMetadataReadback(ws.Root, *project, snapshot)
		result.Kind = ProjectMetadataMutationKind
		result.Changed, result.Event = true, &event
		return nil
	})
	return result, err
}

func projectMetadataPath(files FileSystem, cwd, input string, remove bool) (string, error) {
	if !remove {
		return canonicalDirectory(files, cwd, input, "metadata")
	}
	path := input
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	path = filepath.Clean(path)
	ancestor := path
	missing := []string{}
	for {
		physical, err := files.EvalSymlinks(ancestor)
		if err == nil {
			// Removed children still belong to their physical parent. Preserve
			// that identity when the caller repeats the alias used at add time.
			return filepath.Join(append([]string{physical}, missing...)...), nil
		}
		if !errors.Is(err, fs.ErrNotExist) || filepath.Dir(ancestor) == ancestor {
			return "", projectPathError("resolve metadata removal", path, err)
		}
		missing = append([]string{filepath.Base(ancestor)}, missing...)
		ancestor = filepath.Dir(ancestor)
	}
}
func projectMetadataEventID(event ProjectMetadataEvent) string {
	event.ID = ""
	bytes, _ := json.Marshal(event)
	return "pme_" + strings.TrimPrefix(digestTaskBytes(bytes), "sha256:")
}
