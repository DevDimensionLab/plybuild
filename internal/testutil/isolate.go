package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Isolate runs the tests of one test binary inside a fresh directory: HOME,
// the XDG directories, TMPDIR and GOTMPDIR (which t.TempDir prefers) point
// into it, and PLY_WORKSPACE_CEILING keeps workspace discovery inside it.
// No test, and no command a test starts, can then reach a real workspace,
// store or home, wherever the suite runs.
// The Go environment file, build cache and module cache stay the caller's.
// The directory is removed afterwards. Every test package calls it from
// TestMain: func TestMain(m *testing.M) { os.Exit(testutil.Isolate(m)) }.
//
// A test binary that an isolated test starts again, such as a helper process,
// keeps the directory and environment its parent gave it.
func Isolate(m *testing.M) int {
	if root := os.Getenv(isolatedVariable); root != "" {
		if os.Getenv("PLY_WORKSPACE_CEILING") == "" {
			if err := os.Setenv("PLY_WORKSPACE_CEILING", root); err != nil {
				fmt.Fprintf(os.Stderr, "testutil.Isolate: %v\n", err)
				return 1
			}
		}
		return m.Run()
	}
	root, err := isolate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "testutil.Isolate: %v\n", err)
		return 1
	}
	code := m.Run()
	if err := os.RemoveAll(root); err != nil {
		fmt.Fprintf(os.Stderr, "testutil.Isolate: remove %s: %v\n", root, err)
		if code == 0 {
			code = 1
		}
	}
	return code
}

func isolate() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	config, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	gopath := filepath.Join(home, "go")
	if value := os.Getenv("GOPATH"); value != "" {
		gopath = filepath.SplitList(value)[0]
	}
	for name, value := range map[string]string{
		"GOENV":      filepath.Join(config, "go", "env"),
		"GOCACHE":    filepath.Join(cache, "go-build"),
		"GOPATH":     gopath,
		"GOMODCACHE": filepath.Join(gopath, "pkg", "mod"),
	} {
		if os.Getenv(name) == "" {
			if err := os.Setenv(name, value); err != nil {
				return "", err
			}
		}
	}
	base, err := os.MkdirTemp("", "ply-test-")
	if err != nil {
		return "", err
	}
	root, err := filepath.EvalSymlinks(base)
	if err != nil {
		return "", err
	}
	for name, dir := range map[string]string{
		"HOME":            "home",
		"TMPDIR":          "tmp",
		"GOTMPDIR":        "tmp",
		"XDG_CONFIG_HOME": "home/.config",
		"XDG_STATE_HOME":  "home/.local/state",
		"XDG_CACHE_HOME":  "home/.cache",
		"XDG_DATA_HOME":   "home/.local/share",
	} {
		path := filepath.Join(root, filepath.FromSlash(dir))
		if err := os.MkdirAll(path, 0o700); err != nil {
			return "", err
		}
		if err := os.Setenv(name, path); err != nil {
			return "", err
		}
	}
	if err := os.Setenv(isolatedVariable, root); err != nil {
		return "", err
	}
	return root, os.Setenv("PLY_WORKSPACE_CEILING", root)
}

// isolatedVariable carries the directory to test binaries started from an
// isolated test. It has no PLY_ prefix because helpers drop those variables.
const isolatedVariable = "PLYTEST_ISOLATED_ROOT"
