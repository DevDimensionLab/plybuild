package template

import (
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/sirupsen/logrus"
)

type templatePartialFailureLog struct {
	level   logrus.Level
	message string
}

type templatePartialFailureLogHook struct {
	entries []*logrus.Entry
}

func (hook *templatePartialFailureLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *templatePartialFailureLogHook) Fire(entry *logrus.Entry) error {
	copy := *entry
	hook.entries = append(hook.entries, &copy)
	return nil
}

func TestMergeTemplatesReportsMergeFailureExactlyAndContinuesToLaterTemplate(t *testing.T) {
	root := filepath.Join(t.TempDir(), "merge templates partial failure")
	if err := testutil.CopyFSOutsideWorkingTree(root, fstest.MapFS{
		".":                         &fstest.MapFile{Mode: fs.ModeDir | 0o755},
		"failing-template":          &fstest.MapFile{Data: []byte("regular source file\n"), Mode: 0o644},
		"later-template":            &fstest.MapFile{Mode: fs.ModeDir | 0o755},
		"later-template/.gitignore": &fstest.MapFile{Data: []byte("# empty template\n"), Mode: 0o644},
		"target":                    &fstest.MapFile{Mode: fs.ModeDir | 0o755},
		"target/.keep":              &fstest.MapFile{Data: []byte("target\n"), Mode: 0o644},
	}); err != nil {
		t.Fatalf("create MergeTemplates fixture: %v", err)
	}
	failingSource := filepath.Join(root, "failing-template")
	laterSource := filepath.Join(root, "later-template")
	targetPath := filepath.Join(root, "target")
	templates := []config.CloudTemplate{
		{Name: "failing-template", Project: config.Project{Path: failingSource}},
		{Name: "later-template", Project: config.Project{Path: laterSource}},
	}
	if len(templates) == 0 {
		t.Fatal("MergeTemplates partial-failure population is empty")
	}
	previousLogger := log
	hook := &templatePartialFailureLogHook{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(io.Discard)
	logger.AddHook(hook)
	SetLogger(logger)
	t.Cleanup(func() { SetLogger(previousLogger) })

	MergeTemplates(templates, config.Project{Path: targetPath})

	assertTemplatePartialFailureLogs(t, hook, []templatePartialFailureLog{
		{level: logrus.InfoLevel, message: "applying Template failing-template"},
		{level: logrus.WarnLevel, message: "directory cannot be deeper than filePath"},
		{level: logrus.InfoLevel, message: "applying Template later-template"},
		{level: logrus.InfoLevel, message: "merging template later-template into " + targetPath},
	})
}

func assertTemplatePartialFailureLogs(t *testing.T, hook *templatePartialFailureLogHook, want []templatePartialFailureLog) {
	t.Helper()
	if len(want) == 0 {
		t.Fatal("MergeTemplates log expectation population is empty")
	}
	if len(hook.entries) == 0 {
		t.Fatal("MergeTemplates log population is empty")
	}
	var selected []templatePartialFailureLog
	for _, entry := range hook.entries {
		for _, expected := range want {
			if entry.Level == expected.level && entry.Message == expected.message {
				selected = append(selected, templatePartialFailureLog{level: entry.Level, message: entry.Message})
			}
		}
	}
	if len(selected) == 0 {
		t.Fatal(errors.New("MergeTemplates exact selected log population is empty"))
	}
	if !reflect.DeepEqual(selected, want) {
		t.Fatalf("MergeTemplates selected logs were %#v, want exact failure and continuation %#v", selected, want)
	}
}
