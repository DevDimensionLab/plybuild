package workspace

import (
	"encoding/json"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/githubrepo"
)

const (
	ProjectListReadbackKind = "WorkspaceProjectListReadback@1"
	ProjectReadbackKind     = "WorkspaceProjectReadback@1"
	CoreReadSchemaVersion   = 1
)

type CoreReadWorkspace struct {
	Root string `json:"root"`
}
type projectReadScope struct {
	ProjectID *ProjectID `json:"project_id"`
}
type projectReadRecord struct {
	ID         ProjectID          `json:"project_id"`
	Name       string             `json:"name"`
	Wrapper    string             `json:"wrapper"`
	RepoCount  int                `json:"repo_count"`
	Companions []ProjectCompanion `json:"companions"`
}
type repositoryReadRecord struct {
	ID           RepoID  `json:"repo_id"`
	Locator      string  `json:"locator"`
	GitCommonDir string  `json:"git_common_dir"`
	Wrapper      *string `json:"wrapper"`
	Remote       *string `json:"remote"`
	GitHubURL    *string `json:"github_url"`
}

func projectRead(p ProjectRecord) projectReadRecord {
	return projectReadRecord{p.ID, p.Name, p.Wrapper, len(p.RepoIDs), sortedCompanions(p.Companions)}
}

func MarshalProjectList(result ProjectListResult) ([]byte, error) {
	projects := make([]projectReadRecord, 0, len(result.Projects))
	for _, p := range result.Projects {
		projects = append(projects, projectRead(p))
	}
	sort.Slice(projects, func(i, j int) bool { return projects[i].ID < projects[j].ID })
	return json.Marshal(struct {
		Kind          string              `json:"kind"`
		SchemaVersion int                 `json:"schema_version"`
		Workspace     CoreReadWorkspace   `json:"workspace"`
		Scope         projectReadScope    `json:"scope"`
		Projects      []projectReadRecord `json:"projects"`
	}{ProjectListReadbackKind, CoreReadSchemaVersion, CoreReadWorkspace{result.Workspace}, projectReadScope{}, projects})
}

func MarshalProject(result ProjectResult) ([]byte, error) {
	return json.Marshal(projectReadback(result))
}

// MarshalProjectMetadataMutation is project show with whether the change
// changed ply-workspace.json.
func MarshalProjectMetadataMutation(result ProjectMetadataMutation) ([]byte, error) {
	return json.Marshal(struct {
		projectReadbackRecord
		Changed bool `json:"changed"`
	}{projectReadback(result.ProjectResult), result.Changed})
}

type projectReadbackRecord struct {
	Kind          string                 `json:"kind"`
	SchemaVersion int                    `json:"schema_version"`
	Workspace     CoreReadWorkspace      `json:"workspace"`
	Scope         projectReadScope       `json:"scope"`
	Project       projectReadRecord      `json:"project"`
	Repositories  []repositoryReadRecord `json:"repositories"`
}

func projectReadback(result ProjectResult) projectReadbackRecord {
	repos := make([]repositoryReadRecord, 0, len(result.Repos))
	for _, r := range result.Repos {
		var remote, githubURL *string
		if r.Remote != "" {
			remote = &r.Remote
		}
		if owner, name, ok := githubrepo.FromRemote(r.Remote); ok {
			url := "https://github.com/" + owner + "/" + name
			githubURL = &url
		}
		var wrapper *string
		if r.Wrapper != "" {
			wrapper = &r.Wrapper
		}
		repos = append(repos, repositoryReadRecord{r.ID, r.Locator, r.GitCommonDir, wrapper, remote, githubURL})
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].ID < repos[j].ID })
	return projectReadbackRecord{ProjectReadbackKind, CoreReadSchemaVersion, CoreReadWorkspace{result.Workspace}, projectReadScope{&result.Project.ID}, projectRead(result.Project), repos}
}
