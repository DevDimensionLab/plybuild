package planningrepo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

func existing(git Git, target string) error {
	info, err := os.Lstat(target)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect planning path %s: %w", target, err)
	}
	if info.IsDir() && info.Mode()&fs.ModeSymlink == 0 && recognized(git, target) {
		return fmt.Errorf("existing planning repository at %s; use the existing plan, creation stopped", target)
	}
	return fmt.Errorf("planning path conflict at %s: path already exists; creation stopped", target)
}

func recognized(git Git, target string) bool {
	top, err := git(target, "rev-parse", "--show-toplevel")
	if err != nil {
		return false
	}
	top, err = physicalDirectory(top)
	if err != nil || top != target {
		return false
	}
	manifestPath := filepath.Join(target, ManifestName)
	if regular(manifestPath) {
		content, err := os.ReadFile(manifestPath)
		var manifest Manifest
		if err == nil && json.Unmarshal(content, &manifest) == nil && validManifest(manifest, target) {
			return true
		}
	}
	// Recognition never reads planning prose or treats it as instructions.
	for _, name := range []string{"README.md", "AGENTS.md", "goal.md", "next-task.md", "status.md"} {
		if !regular(filepath.Join(target, name)) {
			return false
		}
	}
	return true
}

func regular(p string) bool {
	info, err := os.Lstat(p)
	return err == nil && info.Mode().IsRegular()
}

var hashPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)

func validManifest(m Manifest, target string) bool {
	if m.SchemaVersion != 1 || m.PlanningPath != target || m.RootPath != filepath.Dir(target) ||
		m.Template.ID != "agentic/planning" || m.Template.Version == "" ||
		(m.Language != "en" && m.Language != "nb") || len(m.Repositories) == 0 || len(m.Files) == 0 {
		return false
	}
	seen := map[string]bool{}
	for i, repo := range m.Repositories {
		if !filepath.IsAbs(repo.Path) || !filepath.IsAbs(repo.GitCommonDir) ||
			filepath.Clean(repo.Path) != repo.Path || filepath.Clean(repo.GitCommonDir) != repo.GitCommonDir ||
			seen[repo.GitCommonDir] || (i > 0 && m.Repositories[i-1].Path >= repo.Path) {
			return false
		}
		seen[repo.GitCommonDir] = true
	}
	for i, file := range m.Files {
		if !safeFile(file.Path) || file.Path == ManifestName || !hashPattern.MatchString(file.SHA256) ||
			(i > 0 && m.Files[i-1].Path >= file.Path) {
			return false
		}
	}
	return true
}
