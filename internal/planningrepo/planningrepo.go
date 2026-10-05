// Package planningrepo creates a standalone local planning repository. Validation
// and rendering finish before the exclusive destination mkdir. Effects are never
// retried or rolled back: a failed creation leaves its actual state for inspection.
package planningrepo

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/plantemplate"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const ManifestName = ".ply-planning.json"

type Input struct {
	Project, Root, Language, Goal string
	Repositories                  []string
	Check                         bool
}

type Repository struct {
	Path         string `json:"path"`
	GitCommonDir string `json:"git_common_dir"`
}

type Template struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

type GitState struct {
	Branch string `json:"branch"`
	OID    string `json:"oid"`
	Tree   string `json:"tree"`
}

type Result struct {
	Kind          string       `json:"kind"`
	SchemaVersion int          `json:"schema_version"`
	State         string       `json:"state"`
	RootPath      string       `json:"root_path"`
	PlanningPath  string       `json:"planning_path"`
	Language      string       `json:"language"`
	Template      Template     `json:"template"`
	Repositories  []Repository `json:"repositories"`
	Files         []string     `json:"files"`
	Git           *GitState    `json:"git"`
}

type FileHash struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type Manifest struct {
	SchemaVersion int          `json:"schema_version"`
	Template      Template     `json:"template"`
	Language      string       `json:"language"`
	RootPath      string       `json:"root_path"`
	PlanningPath  string       `json:"planning_path"`
	Repositories  []Repository `json:"repositories"`
	Files         []FileHash   `json:"files"`
}

// Dependencies exposes only the boundaries needed for deterministic creation
// faults. In particular, there is no cleanup, retry, remote or provider adapter.
type Dependencies struct {
	Git       Git
	Project   func(workspace.ProjectID) (workspace.ProjectResult, error)
	Render    func(plantemplate.Input) (plantemplate.Bundle, error)
	Mkdir     func(string, fs.FileMode) error
	WriteFile func(string, []byte) error
}

func SystemDependencies() Dependencies {
	return Dependencies{
		Git: systemGit,
		Project: func(id workspace.ProjectID) (workspace.ProjectResult, error) {
			return workspace.ShowProject(workspace.SystemDependencies(), id)
		},
		Render: plantemplate.Render, Mkdir: os.Mkdir, WriteFile: writeExclusive,
	}
}

func Validate(input Input) error {
	if (input.Project == "") == (len(input.Repositories) == 0) {
		return fmt.Errorf("select exactly one of --project ID or repeated --repo PATH")
	}
	if input.Project == "" && input.Root == "" {
		return fmt.Errorf("--root PATH is required with --repo; the root is never inferred from cwd")
	}
	if input.Language != "en" && input.Language != "nb" {
		return fmt.Errorf("unsupported language %q: use en or nb", input.Language)
	}
	if input.Project != "" {
		_, err := workspace.ParseProjectID(input.Project)
		return err
	}
	return nil
}

func Create(d Dependencies, input Input) (Result, error) {
	if err := Validate(input); err != nil {
		return Result{}, err
	}
	paths := append([]string(nil), input.Repositories...)
	expected := map[string]string{}
	if input.Project != "" {
		project, err := d.Project(workspace.ProjectID(input.Project))
		if err != nil {
			return Result{}, err
		}
		if input.Root == "" {
			input.Root = project.Project.Wrapper
		}
		for _, repo := range project.Repos {
			paths = append(paths, repo.Locator)
			expected[repo.Locator] = repo.GitCommonDir
		}
	}
	root, err := physicalDirectory(input.Root)
	if err != nil {
		return Result{}, fmt.Errorf("root %q: %w", input.Root, err)
	}
	rootInfo, err := os.Lstat(root)
	if err != nil {
		return Result{}, fmt.Errorf("inspect root %s: %w", root, err)
	}
	target := filepath.Join(root, "planning")
	if err := existing(d.Git, target); err != nil {
		return Result{}, err
	}
	if len(paths) == 0 {
		return Result{}, fmt.Errorf("no repository members selected for %s", target)
	}
	repos := make([]Repository, 0, len(paths))
	seen := map[string]bool{}
	for _, p := range paths {
		repo, gitDir, err := observeRepository(d.Git, p)
		if err != nil {
			return Result{}, err
		}
		if previous, registered := expected[p]; registered && previous != repo.GitCommonDir {
			return Result{}, fmt.Errorf("repository identity changed at %s: registered %s, observed %s", p, previous, repo.GitCommonDir)
		}
		if seen[repo.GitCommonDir] {
			return Result{}, fmt.Errorf("duplicate repository identity %s selected via %s", repo.GitCommonDir, repo.Path)
		}
		seen[repo.GitCommonDir] = true
		for _, forbidden := range []string{repo.Path, repo.GitCommonDir, gitDir} {
			if within(target, forbidden) {
				return Result{}, fmt.Errorf("planning path %s is inside selected repository or Git directory %s", target, forbidden)
			}
		}
		repos = append(repos, repo)
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].Path < repos[j].Path })
	if _, err := d.Git(root, "rev-parse", "--is-inside-work-tree"); err == nil {
		return Result{}, fmt.Errorf("planning path %s is inside an existing Git worktree or Git directory", target)
	} else if !notRepository(err) {
		return Result{}, fmt.Errorf("inspect Git placement at %s: %w", root, err)
	}
	templateRepos := make([]plantemplate.Repository, 0, len(repos))
	for _, repo := range repos {
		templateRepos = append(templateRepos, plantemplate.Repository{Path: repo.Path, GitCommonDir: repo.GitCommonDir})
	}
	bundle, err := d.Render(plantemplate.Input{Language: input.Language, RootPath: root, PlanningPath: target, Goal: input.Goal, Repositories: templateRepos})
	if err != nil {
		return Result{}, fmt.Errorf("render planning repository at %s: %w", target, err)
	}
	manifest := Manifest{SchemaVersion: 1, Template: Template{bundle.TemplateID, bundle.TemplateVersion}, Language: input.Language, RootPath: root, PlanningPath: target, Repositories: repos, Files: []FileHash{}}
	files := append([]plantemplate.File(nil), bundle.Files...)
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	for i, file := range files {
		if !safeFile(file.Path) || file.Path == ManifestName || (i > 0 && files[i-1].Path == file.Path) {
			return Result{}, fmt.Errorf("unsafe or duplicate rendered path %q for %s", file.Path, target)
		}
		manifest.Files = append(manifest.Files, FileHash{file.Path, fmt.Sprintf("sha256:%x", sha256.Sum256(file.Content))})
	}
	if len(files) == 0 {
		return Result{}, fmt.Errorf("empty planning template for %s", target)
	}
	content, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return Result{}, fmt.Errorf("render manifest at %s: %w", target, err)
	}
	files = append(files, plantemplate.File{Path: ManifestName, Content: append(content, '\n')})
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	result := Result{Kind: "PlanningRepoCreation@1", SchemaVersion: 1, State: "preview", RootPath: root, PlanningPath: target, Language: input.Language, Template: manifest.Template, Repositories: repos, Files: []string{}}
	for _, file := range files {
		result.Files = append(result.Files, file.Path)
	}
	for _, ident := range []string{"GIT_AUTHOR_IDENT", "GIT_COMMITTER_IDENT"} {
		if _, err := d.Git(root, "-c", "user.useConfigOnly=true", "var", ident); err != nil {
			return Result{}, fmt.Errorf("Git identity %s is unavailable for %s: %w", ident, target, err)
		}
	}
	if err := unchangedDirectory(root, rootInfo); err != nil {
		return Result{}, err
	}
	if input.Check {
		return result, nil
	}
	if err := d.Mkdir(target, 0o755); err != nil {
		if conflict := existing(d.Git, target); conflict != nil {
			return Result{}, conflict
		}
		return Result{}, fmt.Errorf("create directory %s: %w", target, err)
	}
	partial := func(stage string, err error) (Result, error) {
		return Result{}, fmt.Errorf("planning repository creation failed at %s during %s; partial state preserved, no retry: %w", target, stage, err)
	}
	targetInfo, err := os.Lstat(target)
	if err != nil {
		return partial("inspect destination", err)
	}
	var gitDirectoryInfo fs.FileInfo
	checkDestination := func() error {
		if err := unchangedDirectory(root, rootInfo); err != nil {
			return err
		}
		if err := unchangedDirectory(target, targetInfo); err != nil {
			return err
		}
		if gitDirectoryInfo != nil {
			return unchangedDirectory(filepath.Join(target, ".git"), gitDirectoryInfo)
		}
		return nil
	}
	// Recheck after each effect as well as before the next one. A path replaced
	// by another actor must not redirect later file or Git operations.
	gitAtTarget := func(args ...string) (string, error) {
		if err := checkDestination(); err != nil {
			return "", err
		}
		out, err := d.Git(target, args...)
		if err != nil {
			return "", err
		}
		return out, checkDestination()
	}
	if _, err := gitAtTarget("init", "--initial-branch=main", "--template="); err != nil {
		return partial("git init", err)
	}
	gitDirectoryInfo, err = os.Lstat(filepath.Join(target, ".git"))
	if err != nil {
		return partial("verify Git directory", err)
	}
	for _, binding := range []struct{ argument, expected string }{
		{"--show-toplevel", target},
		{"--git-common-dir", filepath.Join(target, ".git")},
	} {
		observed, err := gitAtTarget("rev-parse", "--path-format=absolute", binding.argument)
		if err != nil {
			return partial("verify standalone Git repository", err)
		}
		if observed != binding.expected {
			return partial("verify standalone Git repository", fmt.Errorf("Git path changed: expected %s, observed %s", binding.expected, observed))
		}
	}
	directories := map[string]bool{".": true}
	for _, file := range files {
		if err := checkDestination(); err != nil {
			return partial("verify destination before "+file.Path, err)
		}
		dir := path.Dir(file.Path)
		if !directories[dir] {
			if err := d.Mkdir(filepath.Join(target, filepath.FromSlash(dir)), 0o755); err != nil {
				return partial("mkdir "+dir, err)
			}
			directories[dir] = true
		}
		parent := filepath.Join(target, filepath.FromSlash(dir))
		if physical, err := physicalDirectory(parent); err != nil || physical != parent {
			return partial("verify parent for "+file.Path, fmt.Errorf("directory changed at %s", parent))
		}
		if err := d.WriteFile(filepath.Join(target, filepath.FromSlash(file.Path)), file.Content); err != nil {
			return partial("write "+file.Path, err)
		}
		if err := checkDestination(); err != nil {
			return partial("verify destination after "+file.Path, err)
		}
	}
	if _, err := gitAtTarget(append([]string{"add", "--force", "--"}, result.Files...)...); err != nil {
		return partial("git add", err)
	}
	if _, err := gitAtTarget("-c", "user.useConfigOnly=true", "commit", "--no-gpg-sign", "-m", "Initialize planning repository"); err != nil {
		return partial("git commit", err)
	}
	git := &GitState{}
	for _, item := range []struct {
		args        []string
		destination *string
	}{
		{[]string{"symbolic-ref", "--short", "HEAD"}, &git.Branch},
		{[]string{"rev-parse", "HEAD"}, &git.OID},
		{[]string{"rev-parse", "HEAD^{tree}"}, &git.Tree},
	} {
		value, err := gitAtTarget(item.args...)
		if err != nil {
			return partial("verify commit", err)
		}
		*item.destination = value
	}
	status, err := gitAtTarget("status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return partial("verify clean worktree", err)
	}
	if status != "" || git.Branch != "main" {
		return partial("verify clean main branch", fmt.Errorf("unexpected branch %q or nonempty status %q", git.Branch, status))
	}
	result.State, result.Git = "created", git
	return result, nil
}

func unchangedDirectory(p string, original fs.FileInfo) error {
	current, err := os.Lstat(p)
	if err != nil || !current.IsDir() || !os.SameFile(original, current) {
		return fmt.Errorf("directory changed at %s; creation stopped", p)
	}
	physical, err := physicalDirectory(p)
	if err != nil || physical != p {
		return fmt.Errorf("physical directory changed at %s; creation stopped", p)
	}
	return nil
}

func physicalDirectory(p string) (string, error) {
	if p == "" {
		return "", fmt.Errorf("directory path is required")
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	physical, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(physical)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", physical)
	}
	return physical, nil
}

func observeRepository(git Git, p string) (Repository, string, error) {
	physical, err := physicalDirectory(p)
	if err != nil {
		return Repository{}, "", fmt.Errorf("repository %q: %w", p, err)
	}
	top, err := git(physical, "rev-parse", "--show-toplevel")
	if err != nil {
		return Repository{}, "", fmt.Errorf("repository %s is not an accessible Git worktree: %w", physical, err)
	}
	top, err = physicalDirectory(top)
	if err != nil || top != physical {
		return Repository{}, "", fmt.Errorf("repository %s must be a Git worktree root", physical)
	}
	common, err := git(physical, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return Repository{}, "", fmt.Errorf("Git common directory for %s: %w", physical, err)
	}
	common, err = physicalDirectory(common)
	if err != nil {
		return Repository{}, "", fmt.Errorf("Git common directory for %s: %w", physical, err)
	}
	gitDir, err := git(physical, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return Repository{}, "", fmt.Errorf("Git directory for %s: %w", physical, err)
	}
	gitDir, err = physicalDirectory(gitDir)
	if err != nil {
		return Repository{}, "", fmt.Errorf("Git directory for %s: %w", physical, err)
	}
	return Repository{physical, common}, gitDir, nil
}

func within(child, parent string) bool {
	rel, err := filepath.Rel(parent, child)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func safeFile(p string) bool {
	return p != "" && p != "." && p != ".." && !strings.HasPrefix(p, "../") &&
		!strings.HasPrefix(p, "/") && !strings.ContainsAny(p, "\\\x00") &&
		path.Clean(p) == p && strings.Split(p, "/")[0] != ".git"
}

func writeExclusive(p string, content []byte) error {
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	_, writeErr := f.Write(content)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
