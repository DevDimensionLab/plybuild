package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestMain(m *testing.M) { os.Exit(Isolate(m)) }

func TestIsolateKeepsHomeTemporaryAndWorkspaceInsideOneDirectory(t *testing.T) {
	ceiling := os.Getenv("PLY_WORKSPACE_CEILING")
	if ceiling == "" || !filepath.IsAbs(ceiling) {
		t.Fatalf("PLY_WORKSPACE_CEILING = %q", ceiling)
	}
	for _, name := range []string{"HOME", "TMPDIR", "GOTMPDIR", "XDG_CONFIG_HOME", "XDG_STATE_HOME", "XDG_CACHE_HOME", "XDG_DATA_HOME"} {
		if dir := os.Getenv(name); !strings.HasPrefix(dir, ceiling+string(filepath.Separator)) {
			t.Fatalf("%s = %q, not inside %s", name, dir, ceiling)
		}
	}
	if home, err := os.UserHomeDir(); err != nil || !strings.HasPrefix(home, ceiling) {
		t.Fatalf("home %q, %v", home, err)
	}
	if dir := t.TempDir(); !strings.HasPrefix(dir, os.Getenv("TMPDIR")) {
		t.Fatalf("t.TempDir %s is outside TMPDIR", dir)
	}
	for _, name := range []string{"GOCACHE", "GOPATH", "GOMODCACHE"} {
		if dir := os.Getenv(name); dir == "" || strings.HasPrefix(dir, ceiling) {
			t.Fatalf("%s = %q, want the caller's shared Go directory", name, dir)
		}
	}
}

// TestEveryTestPackageIsIsolated keeps every test binary of the module behind
// Isolate, so that no test can reach a real workspace, store or home.
func TestEveryTestPackageIsIsolated(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	isolated, tested := map[string]bool{}, map[string]bool{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != root && (strings.HasPrefix(d.Name(), ".") || d.Name() == "testdata" || d.Name() == "docs") {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}
		dir, _ := filepath.Rel(root, filepath.Dir(path))
		tested[dir] = true
		if source := sourceOf(t, path); strings.Contains(source, "func TestMain(m *testing.M)") && strings.Contains(source, "Isolate(m)") {
			isolated[dir] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for dir := range tested {
		if !isolated[dir] {
			t.Errorf("%s: no TestMain that runs testutil.Isolate(m)", dir)
		}
	}
}

func sourceOf(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
