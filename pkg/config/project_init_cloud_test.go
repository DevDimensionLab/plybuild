package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/testutil"
	"github.com/mitchellh/go-homedir"
)

type recordingProjectCloudConfigOpener struct {
	cloud        CloudConfig
	profilePaths []string
}

func (opener *recordingProjectCloudConfigOpener) Open(profilePath string) CloudConfig {
	opener.profilePaths = append(opener.profilePaths, profilePath)
	return opener.cloud
}

func (opener *recordingProjectCloudConfigOpener) dependencies() projectCloudConfigDependencies {
	return projectCloudConfigDependencies{Opener: opener}
}

func (opener *recordingProjectCloudConfigOpener) assertedProfilePaths() ([]string, error) {
	if len(opener.profilePaths) == 0 {
		return nil, errors.New("recorded project cloud-config profile-path population is empty")
	}
	return opener.profilePaths, nil
}

type projectCloudConfigIdentity struct {
	GitCloudConfig
}

func TestProjectCloudConfigOpenerPreservesProductionSelectionAndCacheMapping(t *testing.T) {
	dependencies := systemProjectCloudConfigDependencies()
	if dependencies.Opener == nil {
		t.Fatal("project cloud-config opener selected an incomplete dependency")
	}
	if reflect.TypeOf(dependencies.Opener) != reflect.TypeOf(gitProjectCloudConfigOpener{}) {
		t.Fatalf("project cloud-config opener dependency is %T, want %T", dependencies.Opener, gitProjectCloudConfigOpener{})
	}

	profile := filepath.Join(t.TempDir(), `complete profile path with spaces and backslash\segment-ø`)
	cloud := openProjectCloudConfig(dependencies, profile)
	gitCloud, ok := cloud.(GitCloudConfig)
	if !ok {
		t.Fatalf("production project cloud-config opener returned %T, want GitCloudConfig", cloud)
	}
	wantCache := filepath.Join(profile, "cloud-config")
	if gitCloud.Implementation().Dir() != wantCache {
		t.Fatalf("production project cloud-config cache path was %q, want %q", gitCloud.Implementation().Dir(), wantCache)
	}
}

func TestProjectCloudConfigHelperDeliversCompletePathOnceAndReturnsExactIdentity(t *testing.T) {
	profile := filepath.Join(t.TempDir(), `complete resolved profile path with spaces and backslash\segment-β`)
	cloud := &projectCloudConfigIdentity{}
	opener := &recordingProjectCloudConfigOpener{cloud: cloud}

	actual := openProjectCloudConfig(opener.dependencies(), profile)

	paths, err := opener.assertedProfilePaths()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(paths, []string{profile}) {
		t.Fatalf("project cloud-config opener profile paths were %#v, want one complete path %q", paths, profile)
	}
	if actual != cloud {
		t.Fatalf("project cloud-config helper returned %T %#v, want exact identity %T %#v", actual, actual, cloud, cloud)
	}
}

func TestProjectCloudConfigDependenciesZeroValueReturnsNilWithoutTouchingPath(t *testing.T) {
	profile := filepath.Join(t.TempDir(), "zero project opener must remain absent")
	if cloud := openProjectCloudConfig(projectCloudConfigDependencies{}, profile); cloud != nil {
		t.Fatalf("zero project cloud-config opener returned %T %#v, want nil", cloud, cloud)
	}
	if _, err := os.Lstat(profile); !os.IsNotExist(err) {
		t.Fatalf("zero project cloud-config opener touched profile path %q: %v", profile, err)
	}
}

func TestInitProjectFromDirectoryPreservesActiveProfileAndMigrationCacheMapping(t *testing.T) {
	tests := []struct {
		name         string
		prepare      func(*testing.T, string) string
		wantMigrated bool
	}{
		{
			name: "active profile",
			prepare: func(t *testing.T, home string) string {
				profile := "active profile with spaces"
				writeProjectCloudTestFile(t, filepath.Join(home, ".ply", "profiles", ".active_profile"), []byte(profile))
				return filepath.Join(home, ".ply", "profiles", profile)
			},
		},
		{
			name: "missing profile migrates through legacy home",
			prepare: func(t *testing.T, home string) string {
				return filepath.Join(home, ".co-pilot", "profiles", "default")
			},
			wantMigrated: true,
		},
	}
	if len(tests) == 0 {
		t.Fatal("project profile cache-mapping test-case population is empty")
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			home := isolateProjectCloudTestHome(t)
			wantProfile := test.prepare(t, home)
			projectDir := t.TempDir()
			writeProjectCloudTestConfiguration(t, projectDir, "")

			project, err := InitProjectFromDirectory(projectDir)
			if err != nil {
				t.Fatalf("initialize project with %s: %v", test.name, err)
			}
			assertProjectCloudCachePath(t, project.CloudConfig, filepath.Join(wantProfile, "cloud-config"))
			if test.wantMigrated {
				activeProfile := filepath.Join(home, ".co-pilot", "profiles", ".active_profile")
				contents, err := os.ReadFile(activeProfile)
				if err != nil || string(contents) != "default" {
					t.Fatalf("migrated active profile was (%q, %v), want exact default", contents, err)
				}
			}
		})
	}
}

func TestInitProjectFromDirectoryPreservesDecodedProfileNonOverride(t *testing.T) {
	home := isolateProjectCloudTestHome(t)
	activeProfile := "active-profile"
	writeProjectCloudTestFile(t, filepath.Join(home, ".ply", "profiles", ".active_profile"), []byte(activeProfile))
	projectDir := t.TempDir()
	decodedProfile := "decoded-project-profile"
	writeProjectCloudTestConfiguration(t, projectDir, decodedProfile)

	project, err := InitProjectFromDirectory(projectDir)
	if err != nil {
		t.Fatalf("initialize project with decoded profile: %v", err)
	}
	if project.Config.Profile != decodedProfile {
		t.Fatalf("returned project profile was %q, want decoded value %q", project.Config.Profile, decodedProfile)
	}
	wantCache := filepath.Join(home, ".ply", "profiles", activeProfile, "cloud-config")
	assertProjectCloudCachePath(t, project.CloudConfig, wantCache)
}

func TestInitProjectFromDirectoryReturnsConfigurationFailureBeforeCloudOpening(t *testing.T) {
	isolateProjectCloudTestHome(t)
	projectDir := t.TempDir()
	writeProjectCloudTestFile(t, filepath.Join(projectDir, "ply.json"), []byte(`{"applicationName":"ConfiguredApplication"}`))
	writeProjectCloudTestFile(t, filepath.Join(projectDir, "src", "marker.txt"), []byte("source directory without a language"))

	project, err := InitProjectFromDirectory(projectDir)

	if err == nil {
		t.Fatal("project configuration failure was lost")
	}
	if project.CloudConfig != nil {
		t.Fatalf("project configuration failure returned cloud config %T %#v, want nil before opening", project.CloudConfig, project.CloudConfig)
	}
}

func TestInitProjectFromDirectoryConstructsCloudConfigIndependentlyEachTime(t *testing.T) {
	home := isolateProjectCloudTestHome(t)
	activeProfilePath := filepath.Join(home, ".ply", "profiles", ".active_profile")
	projectDir := t.TempDir()
	writeProjectCloudTestConfiguration(t, projectDir, "")
	profiles := []string{"first-profile", "second-profile"}
	if len(profiles) == 0 {
		t.Fatal("repeated project construction profile population is empty")
	}

	projects := make([]Project, 0, len(profiles))
	for _, profile := range profiles {
		writeProjectCloudTestFile(t, activeProfilePath, []byte(profile))
		project, err := InitProjectFromDirectory(projectDir)
		if err != nil {
			t.Fatalf("initialize project for profile %q: %v", profile, err)
		}
		projects = append(projects, project)
	}

	for index, profile := range profiles {
		wantCache := filepath.Join(home, ".ply", "profiles", profile, "cloud-config")
		assertProjectCloudCachePath(t, projects[index].CloudConfig, wantCache)
	}
}

func TestRecordedProjectCloudConfigOpenerRejectsEmptyPopulation(t *testing.T) {
	if _, err := (&recordingProjectCloudConfigOpener{}).assertedProfilePaths(); err == nil {
		t.Fatal("empty recorded project cloud-config profile-path population passed")
	}
}

func isolateProjectCloudTestHome(t *testing.T) string {
	t.Helper()
	homedir.Reset()
	t.Cleanup(homedir.Reset)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg-config"))
	return home
}

func writeProjectCloudTestConfiguration(t *testing.T, projectDir, profile string) {
	t.Helper()
	contents := `{"applicationName":"ConfiguredApplication","language":"java"}`
	if profile != "" {
		contents = `{"applicationName":"ConfiguredApplication","language":"java","profile":"` + profile + `"}`
	}
	writeProjectCloudTestFile(t, filepath.Join(projectDir, "ply.json"), []byte(contents))
}

func writeProjectCloudTestFile(t *testing.T, path string, contents []byte) {
	t.Helper()
	if err := testutil.CopyFSOutsideWorkingTree(filepath.Dir(path), os.DirFS(t.TempDir())); err != nil {
		t.Fatalf("create disposable project cloud-config fixture directory %s: %v", filepath.Dir(path), err)
	}
	if err := testutil.WriteFileOutsideWorkingTree(path, contents, 0644); err != nil {
		t.Fatalf("write disposable project cloud-config fixture %s: %v", path, err)
	}
}

func assertProjectCloudCachePath(t *testing.T, cloud CloudConfig, want string) {
	t.Helper()
	gitCloud, ok := cloud.(GitCloudConfig)
	if !ok {
		t.Fatalf("project cloud config was %T %#v, want GitCloudConfig", cloud, cloud)
	}
	if actual := gitCloud.Implementation().Dir(); actual != want {
		t.Fatalf("project cloud-config cache path was %q, want %q", actual, want)
	}
}
