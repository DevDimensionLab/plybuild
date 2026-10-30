package taskrun

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/mitchellh/go-homedir"
)

func TestMain(m *testing.M) {
	code := isolatedTestMain(m)
	os.Exit(code)
}

func isolatedTestMain(m *testing.M) int {
	if os.Getenv("PLY_ISOLATION_CHILD") == "1" {
		return m.Run()
	}
	output, err := exec.Command("go", "env", "-json", "GOCACHE", "GOMODCACHE", "GOPATH").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	var cache map[string]string
	if err = json.Unmarshal(output, &cache); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	for key, value := range cache {
		if err = os.Setenv(key, value); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}
	testHome, err := os.MkdirTemp("", "ply-taskrun-tests-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer os.RemoveAll(testHome)
	testHome, err = filepath.EvalSymlinks(testHome)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	testTemp, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err = os.Setenv("TMPDIR", testTemp); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err = os.Setenv("HOME", testHome); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	homedir.Reset()
	return m.Run()
}

func TestFixtureConfigurationIsIsolated(t *testing.T) {
	dataRoot, err := config.ProjectDataRoot("ply")
	if err != nil || dataRoot != "" {
		t.Fatalf("test inherited a configured project data root: %q %v", dataRoot, err)
	}
	if os.Getenv("PLY_ISOLATION_CHILD") == "1" {
		if os.Getenv("HOME") != os.Getenv("PLY_ISOLATION_EXPECTED_HOME") {
			t.Fatal("child did not inherit the fixture profile")
		}
		return
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	child := exec.Command(executable, "-test.run=^TestFixtureConfigurationIsIsolated$")
	child.Env = append(os.Environ(), "PLY_ISOLATION_CHILD=1", "PLY_ISOLATION_EXPECTED_HOME="+os.Getenv("HOME"))
	if output, err := child.CombinedOutput(); err != nil {
		t.Fatalf("child inherited configuration: %v\n%s", err, output)
	}
	f := workflowTestFixture(t)
	o := workflowTestStart(t, f)
	relative, err := filepath.Rel(f.R.WorkspaceRoot, o.Paths.RunRoot)
	if err != nil || filepath.IsAbs(relative) || relative == ".." || len(relative) >= 3 && relative[:3] == "../" {
		t.Fatalf("fixture escaped its workspace: %q %q %v", f.R.WorkspaceRoot, o.Paths.RunRoot, err)
	}
}
