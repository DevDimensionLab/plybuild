package taskjournal

import (
	"encoding/json"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func integrationJournalFixture(t *testing.T) (Service, string) {
	t.Helper()
	dest, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c := exec.Command("go", "test", "./internal/workspace", "-run", "^TestIntegrationEvidenceExportJournalFixture$", "-count=1")
	c.Dir = filepath.Clean("../..")
	c.Env = append(os.Environ(), "PLY_INTEGRATION_JOURNAL_FIXTURE="+dest)
	if b, err := c.CombinedOutput(); err != nil {
		t.Fatalf("synthetic integration fixture: %v: %s", err, b)
	}
	b, err := os.ReadFile(filepath.Join(dest, "fixture.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m struct {
		Root string `json:"root"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(m.Root) })
	source := filepath.Join(dest, "archive")
	err = filepath.WalkDir(source, func(path string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(m.Root, rel)
		info, err := e.Info()
		if err != nil {
			return err
		}
		if e.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if e.Type()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
	if err != nil {
		t.Fatal(err)
	}
	d := workspace.SystemDependencies()
	d.Files = cwdFiles{d.Files, m.Root}
	d.ProjectDataRoot = func(string) (string, error) { return filepath.Join(dest, "data"), nil }
	return New(d), m.Root
}
func corruptJournalEvidenceMode(t *testing.T, root string, directory bool) string {
	t.Helper()
	r, err := workspace.SystemDependencies().WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	sources := workspace.IntegrationEvidenceSources(root, r)
	if len(sources) == 0 {
		t.Fatal("missing referenced integration fixture")
	}
	path := sources[0].Locator
	mode := fs.FileMode(0600)
	if directory {
		path = filepath.Dir(path)
		mode = 0755
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
	return sources[0].Locator
}
func assertJournalEvidenceDowngraded(t *testing.T, s Snapshot, path string) {
	t.Helper()
	found := false
	for _, src := range s.Sources {
		if src.Locator == path {
			found = true
			if src.Status == "valid" {
				t.Fatal("invalid managed evidence remains valid")
			}
		}
	}
	if !found {
		t.Fatal("integration object omitted from source provenance")
	}
	if s.Current.Axes["integration"].BasisStatus == "current" {
		t.Fatal("invalid integration evidence remains current")
	}
	if len(s.Events) == 0 {
		t.Fatal("historical events disappeared")
	}
}
func TestIntegrationEvidenceJournalShowFinalFreshness(t *testing.T) {
	for _, directory := range []bool{false, true} {
		name := "object"
		if directory {
			name = "directory"
		}
		t.Run(name, func(t *testing.T) {
			s, root := integrationJournalFixture(t)
			healthy, err := s.Show("task", Options{})
			if err != nil {
				t.Fatal(err)
			}
			if healthy.Current.Axes["integration"].BasisStatus != "current" {
				t.Fatalf("healthy integration not current: %+v", healthy.Current.Axes["integration"])
			}
			path := ""
			s.fault = func(stage string) error {
				if stage == "before_freshness" {
					path = corruptJournalEvidenceMode(t, root, directory)
				}
				return nil
			}
			out, err := s.Show("task", Options{})
			if err != nil {
				t.Fatal(err)
			}
			assertJournalEvidenceDowngraded(t, out, path)
		})
	}
}
func TestIntegrationEvidenceJournalBatchCapturedSnapshot(t *testing.T) {
	for _, directory := range []bool{false, true} {
		name := "object"
		if directory {
			name = "directory"
		}
		t.Run(name, func(t *testing.T) {
			s, root := integrationJournalFixture(t)
			basis, err := workspace.ReadViewBasis(s.Workspace)
			if err != nil {
				t.Fatal(err)
			}
			in := EventBatchInput{Workspace: basis.Workspace, Projects: basis.Projects, Registry: basis.Registry}
			path := corruptJournalEvidenceMode(t, root, directory)
			batch, err := ReadEventBatch(s.Workspace, in)
			if err != nil {
				t.Fatal(err)
			}
			assertJournalEvidenceDowngraded(t, batch.Tasks["task"], path)
		})
	}
}

func TestIntegrationEvidenceJournalBatchFinalFreshness(t *testing.T) {
	for _, directory := range []bool{false, true} {
		name := "object"
		if directory {
			name = "directory"
		}
		t.Run(name, func(t *testing.T) {
			s, root := integrationJournalFixture(t)
			basis, err := workspace.ReadViewBasis(s.Workspace)
			if err != nil {
				t.Fatal(err)
			}
			in := EventBatchInput{Workspace: basis.Workspace, Projects: basis.Projects, Registry: basis.Registry}
			path := ""
			s.fault = func(stage string) error {
				if stage == "before_freshness" {
					path = corruptJournalEvidenceMode(t, root, directory)
				}
				return nil
			}
			batch, err := s.readEventBatch(in)
			if err != nil {
				t.Fatal(err)
			}
			assertJournalEvidenceDowngraded(t, batch.Tasks["task"], path)
			if len(batch.Reasons) == 0 {
				t.Fatal("batch omitted its freshness failure")
			}
		})
	}
}
func TestIntegrationEvidenceCacheKeepsManagedValidation(t *testing.T) {
	s, root := integrationJournalFixture(t)
	r, err := s.Workspace.WorkItems.Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	path := workspace.IntegrationEvidenceSources(root, r)[0].Locator
	cache := map[string]sourceRead{}
	if _, err := cachedSourceHash(cache, root, "artifact", path); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := cachedSourceHash(cache, root, workspace.IntegrationEvidenceSourceKind, path); err == nil {
		t.Fatal("generic cached hash bypassed managed validation")
	}
	if err := os.Chmod(path, 0400); err != nil {
		t.Fatal(err)
	}
	cache = map[string]sourceRead{}
	if _, err := cachedSourceHash(cache, root, workspace.IntegrationEvidenceSourceKind, path); err != nil {
		t.Fatal(err)
	}
	if _, err := cachedSourceHash(cache, root, "artifact", path); err != nil {
		t.Fatal(err)
	}
	if cache[path].kind != workspace.IntegrationEvidenceSourceKind {
		t.Fatal("generic source weakened final managed recheck")
	}
	if err := os.Chmod(path, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := journalSourceHash(root, cache[path].kind, path); err == nil {
		t.Fatal("final cached-source recheck accepted writable evidence")
	}
}
