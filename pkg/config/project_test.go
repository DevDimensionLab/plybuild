package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/mitchellh/go-homedir"
)

func TestProjectConfiguration_SourceMainPath(t *testing.T) {
	sourceConfig, err := InitProjectConfigurationFromDir("test/cloud-config/templates/test-template")
	if err != nil {
		t.Fatalf("load source project configuration: %v", err)
	}

	targetConfig, err := InitProjectConfigurationFromDir("test/target-app")
	if err != nil {
		t.Fatalf("load target project configuration: %v", err)
	}

	expectedSourceRoot := "src/main/kotlin/no/ply/template/test"
	if sourceConfig.SourceMainPath() != expectedSourceRoot {
		t.Errorf("expected %s, but got instead %s", expectedSourceRoot, sourceConfig.SourceMainPath())
	}

	expectedTargetRoot := "src/main/java/no/ply/template/target"
	if targetConfig.SourceMainPath() != expectedTargetRoot {
		t.Errorf("expected %s, but got instead %s", expectedTargetRoot, targetConfig.SourceMainPath())
	}
}

func TestSortAndWritePom_sort_enabled_by_default(t *testing.T) {
	project, pomFile := sortingProjectFromUnsortedPom(t, "test/sorting/sorted")

	if err := project.SortAndWritePom(); err != nil {
		t.Fatalf("sort and write POM: %v", err)
	}

	assertFilesEqual(t, "test/sorting/sorted/pom.xml", pomFile)
}

func TestSortAndWritePom_sort_enabled(t *testing.T) {
	project, pomFile := sortingProjectFromUnsortedPom(t, "test/sorting/sorted")

	projectConfig := ProjectConfiguration{
		Settings: ProjectSettings{
			DisableDependencySort: false,
		},
	}

	project.Config = projectConfig
	if err := project.SortAndWritePom(); err != nil {
		t.Fatalf("sort and write POM: %v", err)
	}

	assertFilesEqual(t, "test/sorting/sorted/pom.xml", pomFile)
}

func TestSortAndWritePom_sort_disabled(t *testing.T) {
	project, pomFile := sortingProjectFromUnsortedPom(t, "test/sorting/unsorted")

	projectConfig := ProjectConfiguration{
		Settings: ProjectSettings{
			DisableDependencySort: true,
		},
	}

	project.Config = projectConfig
	if err := project.SortAndWritePom(); err != nil {
		t.Fatalf("write POM with sorting disabled: %v", err)
	}

	assertFilesEqual(t, "test/sorting/pom.xml", pomFile)
}

func TestCreateProjectConfig(t *testing.T) {
	isolateConfigTestHome(t)
	projectDir := t.TempDir()
	copyConfigFixtureFile(t, "test/project-config/pom.xml", filepath.Join(projectDir, "pom.xml"))
	newConfig := filepath.Join(projectDir, "ply.json")
	if _, err := os.Stat(newConfig); !os.IsNotExist(err) {
		t.Fatalf("generated config must not exist before initialization, stat returned %v", err)
	}

	project, err := InitProjectFromDirectory(projectDir)
	if err != nil {
		t.Fatalf("load project fixture: %v", err)
	}

	err = project.InitProjectConfiguration()
	if err != nil {
		t.Fatalf("initialize project configuration: %v", err)
	}

	if err := project.Config.WriteTo(project.ConfigFile); err != nil {
		t.Fatalf("write project configuration: %v", err)
	}

	originConfig := "test/project-config/origin.ply.json"
	equal, err := file.Equal(originConfig, newConfig)
	if err != nil {
		t.Fatalf("compare generated project configuration: %v", err)
	}
	if !equal {
		t.Errorf("%s is not equal to %s", originConfig, newConfig)
	}
}

func TestProjectConfigFilePreservesExistingLegacyName(t *testing.T) {
	projectDir := t.TempDir()
	legacyConfig := filepath.Join(projectDir, legacyProjectConfigFileName)
	copyConfigFixtureFile(t, "test/project-config/origin.ply.json", legacyConfig)

	if actual := projectConfigFile(projectDir); actual != legacyConfig {
		t.Fatalf("legacy project config path was %q, want %q", actual, legacyConfig)
	}
}

func sortingProjectFromUnsortedPom(t *testing.T, fixtureDir string) (Project, string) {
	t.Helper()
	isolateConfigTestHome(t)

	projectDir := copyConfigFixtureDir(t, fixtureDir)
	pomFile := filepath.Join(projectDir, "pom.xml")
	copyConfigFixtureFile(t, "test/sorting/pom.xml", pomFile)

	project, err := InitProjectFromPomFile(pomFile)
	if err != nil {
		t.Fatalf("load project fixture: %v", err)
	}
	return project, pomFile
}

func isolateConfigTestHome(t *testing.T) {
	t.Helper()
	homedir.Reset()
	t.Setenv("HOME", t.TempDir())
	t.Cleanup(homedir.Reset)
}

func copyConfigFixtureDir(t *testing.T, source string) string {
	t.Helper()
	target := filepath.Join(t.TempDir(), filepath.Base(source))
	if err := testutil.CopyFSOutsideWorkingTree(target, os.DirFS(source)); err != nil {
		t.Fatalf("copy fixture directory %s: %v", source, err)
	}
	return target
}

func copyConfigFixtureFile(t *testing.T, source, target string) {
	t.Helper()
	contents, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read fixture %s: %v", source, err)
	}
	if err := testutil.WriteFileOutsideWorkingTree(target, contents, 0o644); err != nil {
		t.Fatalf("copy fixture %s to %s: %v", source, target, err)
	}
}

func assertFilesEqual(t *testing.T, expected, actual string) {
	t.Helper()
	equal, err := file.Equal(expected, actual)
	if err != nil {
		t.Fatalf("compare %s with %s: %v", expected, actual, err)
	}
	if !equal {
		t.Errorf("%s is not equal to %s", actual, expected)
	}
}
