package cmd

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"testing/fstest"

	"github.com/devdimensionlab/mvn-pom-mutator/pkg/pom"
	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/config"
	plycontext "github.com/devdimensionlab/plybuild/pkg/context"
	"github.com/sirupsen/logrus"
)

type initMessageOnlyFormatter struct{}

func (initMessageOnlyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}

func TestInitCmdReportsEveryProjectFailureExactlyAndContinuesToLaterProjects(t *testing.T) {
	root := filepath.Join(t.TempDir(), "init partial failures")
	fixture := fstest.MapFS{
		".":                        &fstest.MapFile{Mode: fs.ModeDir | 0o755},
		"config-error":             &fstest.MapFile{Mode: fs.ModeDir | 0o755},
		"config-error/placeholder": &fstest.MapFile{Data: []byte("config directory\n"), Mode: 0o644},
		"pom-error":                &fstest.MapFile{Mode: fs.ModeDir | 0o755},
		"pom-error/placeholder":    &fstest.MapFile{Data: []byte("pom directory\n"), Mode: 0o644},
	}
	if err := testutil.CopyFSOutsideWorkingTree(root, fixture); err != nil {
		t.Fatalf("create init failure fixture: %v", err)
	}

	completeConfig := config.ProjectConfiguration{
		MavenProjectConfiguration: config.MavenProjectConfiguration{
			Artifact: config.Artifact{GroupId: "com.example", ArtifactId: "complete"},
			Language: "kotlin",
			Package:  "com.example.complete",
		},
		Name: "complete",
	}
	configErrorPath := filepath.Join(root, "config-error")
	pomErrorPath := filepath.Join(root, "pom-error")
	successConfigPath := filepath.Join(root, "later-project.json")
	successPomPath := filepath.Join(root, "later-pom.xml")

	previousContext := ctx
	previousLogger := log
	output := &bytes.Buffer{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(output)
	logger.SetFormatter(initMessageOnlyFormatter{})
	t.Cleanup(func() {
		ctx = previousContext
		log = previousLogger
	})
	log = logger
	ctx = plycontext.Context{Projects: []config.Project{
		{Path: filepath.Join(root, "missing-project-type")},
		{
			Path: filepath.Join(root, "missing-model"),
			Type: config.MavenProject{PomFile: filepath.Join(root, "missing-model.xml")},
		},
		{
			Path:       filepath.Join(root, "config-write-failure"),
			ConfigFile: configErrorPath,
			Config:     completeConfig,
			Type: config.MavenProject{
				PomFile:  filepath.Join(root, "config-write-failure.xml"),
				PomModel: &pom.Model{},
			},
		},
		{
			Path:       filepath.Join(root, "pom-write-failure"),
			ConfigFile: filepath.Join(root, "pom-write-failure.json"),
			Config:     completeConfig,
			Type: config.MavenProject{
				PomFile:  pomErrorPath,
				PomModel: &pom.Model{},
			},
		},
		{
			Path:       filepath.Join(root, "later-success"),
			ConfigFile: successConfigPath,
			Config:     completeConfig,
			Type: config.MavenProject{
				PomFile:  successPomPath,
				PomModel: &pom.Model{ModelVersion: "4.0.0"},
			},
		},
	}}

	initCmd.Run(initCmd, nil)

	wantLog := strings.Join([]string{
		"no project type defined for " + filepath.Join(root, "missing-project-type") + ":",
		"formating pom file " + filepath.Join(root, "missing-model.xml"),
		"project type and model is nil",
		"formating pom file " + filepath.Join(root, "config-write-failure.xml"),
		(&os.PathError{Op: "open", Path: configErrorPath, Err: syscall.EISDIR}).Error(),
		"formating pom file " + pomErrorPath,
		(&os.PathError{Op: "open", Path: pomErrorPath, Err: syscall.EISDIR}).Error(),
		"formating pom file " + successPomPath,
	}, "\n") + "\n"
	if output.String() != wantLog {
		t.Fatalf("init partial-failure log differed:\n got:\n%s\nwant:\n%s", output.String(), wantLog)
	}
	writtenPaths := []string{successConfigPath, successPomPath}
	if len(writtenPaths) == 0 {
		t.Fatal("later init write expectation population is empty")
	}
	for _, path := range writtenPaths {
		contents, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("later init project did not write %s after earlier failures: %v", path, err)
		}
		if len(contents) == 0 {
			t.Fatalf("later init project wrote empty file %s", path)
		}
	}
	if _, err := os.Stat(configErrorPath); err != nil || !isDirectory(configErrorPath) {
		t.Fatalf("config failure directory changed unexpectedly: %v", err)
	}
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode()&fs.ModeDir != 0
}
