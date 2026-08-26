package spring

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"testing/fstest"

	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/sirupsen/logrus"
)

type springPartialFailureLogHook struct {
	entries []*logrus.Entry
}

func (hook *springPartialFailureLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *springPartialFailureLogHook) Fire(entry *logrus.Entry) error {
	copy := *entry
	hook.entries = append(hook.entries, &copy)
	return nil
}

func (hook *springPartialFailureLogHook) assertedWarnings() ([]string, error) {
	if len(hook.entries) == 0 {
		return nil, errors.New("recorded Spring DeleteDemoFiles log population is empty")
	}
	var warnings []string
	for _, entry := range hook.entries {
		if entry.Level == logrus.WarnLevel {
			warnings = append(warnings, entry.Message)
		}
	}
	if len(warnings) == 0 {
		return nil, errors.New("recorded Spring DeleteDemoFiles warning population is empty")
	}
	return warnings, nil
}

func TestDeleteDemoFilesReportsDeletionFailuresExactlyAndContinues(t *testing.T) {
	t.Run("test deletion failure continues to fixed demo files", func(t *testing.T) {
		testDirectory := filepath.Join("src", "test", "kotlin", "BrokenTest.kt")
		root := copyDeleteDemoFixture(t, fstest.MapFS{
			filepath.ToSlash(filepath.Join(testDirectory, "child.txt")): &fstest.MapFile{Data: []byte("child\n"), Mode: 0o644},
			"HELP.md":  &fstest.MapFile{Data: []byte("help\n"), Mode: 0o644},
			"mvnw":     &fstest.MapFile{Data: []byte("wrapper\n"), Mode: 0o755},
			"mvnw.cmd": &fstest.MapFile{Data: []byte("wrapper\n"), Mode: 0o644},
		})
		hook := captureDeleteDemoWarnings(t)

		DeleteDemoFiles(root, config.ProjectConfiguration{MavenProjectConfiguration: config.MavenProjectConfiguration{Language: "kotlin"}})

		assertDeleteDemoWarnings(t, hook, []string{"Unable to delete testfile: " + filepath.Join(root, testDirectory)})
		assertDeleteDemoFilesAbsent(t, root, []string{"HELP.md", "mvnw", "mvnw.cmd"})
		if _, err := os.Stat(filepath.Join(root, testDirectory)); err != nil {
			t.Fatalf("failed test-directory deletion did not preserve its non-empty directory: %v", err)
		}
	})

	t.Run("fixed demo deletion failures continue through later entries", func(t *testing.T) {
		root := copyDeleteDemoFixture(t, fstest.MapFS{
			"src/test/kotlin/CompleteTest.kt": &fstest.MapFile{Data: []byte("class CompleteTest\n"), Mode: 0o644},
			"HELP.md/child.txt":               &fstest.MapFile{Data: []byte("help child\n"), Mode: 0o644},
			"mvnw":                            &fstest.MapFile{Data: []byte("wrapper\n"), Mode: 0o755},
			"mvnw.cmd/child.txt":              &fstest.MapFile{Data: []byte("wrapper child\n"), Mode: 0o644},
		})
		hook := captureDeleteDemoWarnings(t)

		DeleteDemoFiles(root, config.ProjectConfiguration{MavenProjectConfiguration: config.MavenProjectConfiguration{Language: "kotlin"}})

		helpPath := filepath.Join(root, "HELP.md")
		commandPath := filepath.Join(root, "mvnw.cmd")
		assertDeleteDemoWarnings(t, hook, []string{
			(&os.PathError{Op: "remove", Path: helpPath, Err: syscall.ENOTEMPTY}).Error(),
			(&os.PathError{Op: "remove", Path: commandPath, Err: syscall.ENOTEMPTY}).Error(),
		})
		assertDeleteDemoFilesAbsent(t, root, []string{"mvnw"})
		preservedPaths := []string{helpPath, commandPath}
		if len(preservedPaths) == 0 {
			t.Fatal("DeleteDemoFiles preserved-path expectation population is empty")
		}
		for _, path := range preservedPaths {
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("failed fixed demo deletion did not preserve %s: %v", path, err)
			}
		}
	})
}

func TestDeleteDemoFilesReportsDiscoveryFailureExactlyAndContinuesToFixedDemoFiles(t *testing.T) {
	root := copyDeleteDemoFixture(t, fstest.MapFS{
		"HELP.md":  &fstest.MapFile{Data: []byte("help\n"), Mode: 0o644},
		"mvnw":     &fstest.MapFile{Data: []byte("wrapper\n"), Mode: 0o755},
		"mvnw.cmd": &fstest.MapFile{Data: []byte("wrapper\n"), Mode: 0o644},
	})
	hook := captureDeleteDemoWarnings(t)
	sentinel := errors.New("complete DeleteDemoFiles discovery failure")
	type discoveryCall struct {
		fileSuffix string
		directory  string
	}
	var calls []discoveryCall
	findFirst := func(fileSuffix string, directory string) (string, error) {
		calls = append(calls, discoveryCall{fileSuffix: fileSuffix, directory: directory})
		return "", sentinel
	}

	deleteDemoFiles(findFirst, root, config.ProjectConfiguration{
		MavenProjectConfiguration: config.MavenProjectConfiguration{Language: "kotlin"},
	})

	wantCalls := []discoveryCall{{
		fileSuffix: ".kt",
		directory:  filepath.Join(root, "src", "test", "kotlin"),
	}}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("DeleteDemoFiles discovery calls were %#v, want exact lookup %#v", calls, wantCalls)
	}
	assertDeleteDemoWarnings(t, hook, []string{"Unable to find testfile, fileSuffix=.kt"})
	assertDeleteDemoFilesAbsent(t, root, []string{"HELP.md", "mvnw", "mvnw.cmd"})
}

func copyDeleteDemoFixture(t *testing.T, fixture fstest.MapFS) string {
	t.Helper()
	if len(fixture) == 0 {
		t.Fatal("DeleteDemoFiles fixture population is empty")
	}
	completeFixture := fstest.MapFS{
		".": &fstest.MapFile{Mode: fs.ModeDir | 0o755},
	}
	for path, entry := range fixture {
		parts := strings.Split(filepath.ToSlash(path), "/")
		if len(parts) > 1 {
			for index := 1; index < len(parts); index++ {
				completeFixture[strings.Join(parts[:index], "/")] = &fstest.MapFile{Mode: fs.ModeDir | 0o755}
			}
		}
		completeFixture[path] = entry
	}
	root := filepath.Join(t.TempDir(), "delete demo files")
	if err := testutil.CopyFSOutsideWorkingTree(root, completeFixture); err != nil {
		t.Fatalf("copy DeleteDemoFiles fixture: %v", err)
	}
	return root
}

func captureDeleteDemoWarnings(t *testing.T) *springPartialFailureLogHook {
	t.Helper()
	previousLogger := log
	hook := &springPartialFailureLogHook{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.AddHook(hook)
	SetLogger(logger)
	t.Cleanup(func() { SetLogger(previousLogger) })
	return hook
}

func assertDeleteDemoWarnings(t *testing.T, hook *springPartialFailureLogHook, want []string) {
	t.Helper()
	if len(want) == 0 {
		t.Fatal("DeleteDemoFiles warning expectation population is empty")
	}
	actual, err := hook.assertedWarnings()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("DeleteDemoFiles warnings were %#v, want exact ordered warnings %#v", actual, want)
	}
}

func assertDeleteDemoFilesAbsent(t *testing.T, root string, paths []string) {
	t.Helper()
	if len(paths) == 0 {
		t.Fatal("DeleteDemoFiles absence expectation population is empty")
	}
	for _, path := range paths {
		if _, err := os.Stat(filepath.Join(root, path)); !errors.Is(err, fs.ErrNotExist) {
			t.Fatalf("later DeleteDemoFiles item %s was not deleted after an earlier failure: %v", path, err)
		}
	}
}
