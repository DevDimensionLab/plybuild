package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/mitchellh/go-homedir"
)

func newMockCloudConfig(path string) (cfg config.GitCloudConfig) {
	cfg.Impl.Path = path
	return
}

func TestMergeTemplate_test_template(t *testing.T) {
	fixtures := copyTemplateFixtures(t)
	cfg := newMockCloudConfig(filepath.Join(fixtures, "cloud-config"))
	targetDir := filepath.Join(fixtures, "target-test-template")
	project, err := config.InitProjectFromDirectory(targetDir)
	if err != nil {
		t.Fatalf("load target project: %v", err)
	}
	cloudTemplate, err := cfg.Template("test-template")
	if err != nil {
		t.Fatalf("load cloud template: %v", err)
	}
	if err := MergeTemplate(cloudTemplate, project, false); err != nil {
		t.Fatalf("merge template: %v", err)
	}

	assertTemplateFileContent(t, filepath.Join(targetDir, "test.properties"), "key=val")
	assertTemplateFileContent(t, filepath.Join(targetDir, "textfile.txt"), "hello world")
}

func TestMergeTemplate_simple_template(t *testing.T) {
	fixtures := copyTemplateFixtures(t)
	cfg := newMockCloudConfig(filepath.Join(fixtures, "cloud-config"))
	targetDir := filepath.Join(fixtures, "target-simple-template")
	project, err := config.InitProjectFromDirectory(targetDir)
	if err != nil {
		t.Fatalf("load target project: %v", err)
	}
	cloudTemplate, err := cfg.Template("simple-template")
	if err != nil {
		t.Fatalf("load cloud template: %v", err)
	}
	if err := MergeTemplate(cloudTemplate, project, false); err != nil {
		t.Fatalf("merge template: %v", err)
	}

	generatedSource := filepath.Join(targetDir, "src/main/java/no/ply/template/target/DummyConfiguration.kt")
	contents, err := os.ReadFile(generatedSource)
	if err != nil {
		t.Fatalf("read generated source: %v", err)
	}
	if !strings.Contains(string(contents), "package no.ply.template.target") {
		t.Errorf("generated source has the wrong package:\n%s", contents)
	}
}

func TestReplacePathForSource(t *testing.T) {
	fixtures := copyTemplateFixtures(t)
	sourceDir := filepath.Join(fixtures, "cloud-config/templates/simple-template")
	targetDir := filepath.Join(fixtures, "target-simple-template")

	sourceConfig, err := config.InitProjectConfigurationFromDir(sourceDir)
	if err != nil {
		t.Fatalf("load source project configuration: %v", err)
	}

	targetConfig, err := config.InitProjectConfigurationFromDir(targetDir)
	if err != nil {
		t.Fatalf("load target project configuration: %v", err)
	}

	files, err := filteredFilesFromTemplate(sourceDir, getIgnores(sourceDir))
	if err != nil {
		t.Fatalf("list source template files: %v", err)
	}
	sourceFilesChecked := 0
	for _, f := range files {
		if strings.HasSuffix(f, ".kt") {
			sourceFilesChecked++
			sourceRelPath, err := file.RelPath(sourceDir, f)
			if err != nil {
				t.Fatalf("make source path relative: %v", err)
			}
			sourceRelPath = replacePathForSource(sourceRelPath, sourceConfig, targetConfig)

			expected := file.Path("src/main/java/no/ply/template/target/DummyConfiguration.kt")
			if sourceRelPath != expected {
				t.Errorf("expected rewritten path %s, got %s", expected, sourceRelPath)
			}
		}
	}
	if sourceFilesChecked != 1 {
		t.Fatalf("expected one Kotlin source path to exercise rewriting, checked %d", sourceFilesChecked)
	}
}

func copyTemplateFixtures(t *testing.T) string {
	t.Helper()
	homedir.Reset()
	t.Setenv("HOME", t.TempDir())
	t.Cleanup(homedir.Reset)

	target := filepath.Join(t.TempDir(), "fixtures")
	if err := testutil.CopyFSOutsideWorkingTree(target, os.DirFS("test")); err != nil {
		t.Fatalf("copy template fixtures: %v", err)
	}
	return target
}

func assertTemplateFileContent(t *testing.T, path, expected string) {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file %s: %v", path, err)
	}
	if string(contents) != expected {
		t.Errorf("generated file %s contained %q, want %q", path, contents, expected)
	}
}
