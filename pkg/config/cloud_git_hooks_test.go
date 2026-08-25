package config

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/file"
)

func (*recordingGitHookFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected git-hook archive open")
}

func (*recordingGitHookFilesystem) Close(filesystem.File) error {
	return errors.New("unexpected Git-hook-files close")
}

func (*recordingGitHookFilesystem) CloseReader(io.ReadCloser) error {
	return errors.New("unexpected Git-hook-files reader close")
}

type recordingGitHookFilesystem struct {
	readDirPaths     []string
	readDirEntries   []fs.FileInfo
	deliveredEntries []fs.FileInfo
	readDirErr       error
}

func (*recordingGitHookFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected Git-hook-files working directory")
}

func (*recordingGitHookFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected Git-hook-files read")
}

func (recording *recordingGitHookFilesystem) ReadDir(path string) ([]fs.FileInfo, error) {
	recording.readDirPaths = append(recording.readDirPaths, path)
	recording.deliveredEntries = append(recording.deliveredEntries, recording.readDirEntries...)
	return append([]fs.FileInfo{}, recording.readDirEntries...), recording.readDirErr
}

func (*recordingGitHookFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected Git-hook-files read directory entries")
}

func (*recordingGitHookFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected Git-hook-files stat")
}

func (*recordingGitHookFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected Git-hook-files mkdir")
}

func (*recordingGitHookFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected Git-hook-files mkdir")
}

func (*recordingGitHookFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected Git-hook-files write")
}

func (*recordingGitHookFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected Git-hook-files open file")
}

func (*recordingGitHookFilesystem) Remove(string) error {
	return errors.New("unexpected Git-hook-files remove")
}

func (*recordingGitHookFilesystem) RemoveAll(string) error {
	return errors.New("unexpected Git-hook-files recursive remove")
}

func (*recordingGitHookFilesystem) Rename(string, string) error {
	return errors.New("unexpected Git-hook-files rename")
}

func (*recordingGitHookFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected Git-hook-files glob")
}

func (*recordingGitHookFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected Git-hook-files walk")
}

func (*recordingGitHookFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected Git-hook-files create")
}

func (*recordingGitHookFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected Git-hook-files copy")
}

func (recording *recordingGitHookFilesystem) dependencies() gitHookFilesDependencies {
	return gitHookFilesDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingGitHookFilesystem) assertedReadDirPaths() ([]string, error) {
	if len(recording.readDirPaths) == 0 {
		return nil, errors.New("recorded Git-hook-files directory-read population is empty")
	}
	return recording.readDirPaths, nil
}

func (recording *recordingGitHookFilesystem) assertedDeliveredEntries() ([]fs.FileInfo, error) {
	if len(recording.deliveredEntries) == 0 {
		return nil, errors.New("recorded Git-hook-files entry population is empty")
	}
	return recording.deliveredEntries, nil
}

type gitHookFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	system  interface{}
}

func (info gitHookFileInfo) Name() string       { return info.name }
func (info gitHookFileInfo) Size() int64        { return info.size }
func (info gitHookFileInfo) Mode() fs.FileMode  { return info.mode }
func (info gitHookFileInfo) ModTime() time.Time { return info.modTime }
func (info gitHookFileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info gitHookFileInfo) Sys() interface{}   { return info.system }

func TestGitHookFilesSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemGitHookFilesDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Git-hook-files selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Git-hook-files filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestGitHookFilesPreservesReceiverRootDeliveredOrderMetadataClassificationAndFilenameResults(t *testing.T) {
	implementationRoot := "/complete cloud-config receiver/root with spaces"
	path := "complete Git-hook folder/with spaces"
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: implementationRoot}}
	firstTime := time.Unix(1700000001, 0)
	secondTime := time.Unix(1700000002, 0)
	entries := []fs.FileInfo{
		gitHookFileInfo{name: "z-first-hook", size: 41, mode: 0641, modTime: firstTime, system: "first metadata"},
		gitHookFileInfo{name: "ignored-hook-directory", size: 82, mode: fs.ModeDir | 0751, system: "directory metadata"},
		gitHookFileInfo{name: "a-second-hook", size: 83, mode: 0601, modTime: secondTime, system: "second metadata"},
		gitHookFileInfo{name: "link-hook", size: 84, mode: fs.ModeSymlink | 0777, system: "link metadata"},
		gitHookFileInfo{name: "ignored-second-directory", size: 85, mode: fs.ModeDir | 0705},
	}
	recording := &recordingGitHookFilesystem{readDirEntries: entries}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("Git-hook-files dependency lost its complete filesystem value: %#v", dependencies)
	}

	actual, err := gitCfg.gitHookFiles(dependencies, path)

	if err != nil {
		t.Fatalf("Git-hook-files returned an error: %v", err)
	}
	wantRoot := file.Path("%s/%s", gitCfg.Implementation().Dir(), path)
	paths, populationErr := recording.assertedReadDirPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(paths, []string{wantRoot}) {
		t.Fatalf("Git-hook-files directory-read roots were %#v, want exact receiver-derived root %#v", paths, []string{wantRoot})
	}
	deliveredEntries, populationErr := recording.assertedDeliveredEntries()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(deliveredEntries, entries) {
		t.Fatalf("Git-hook-files entries were %#v, want ordered metadata %#v", deliveredEntries, entries)
	}
	want := []string{"z-first-hook", "a-second-hook", "link-hook"}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("Git-hook-files results were %#v, want filename-only delivered order %#v", actual, want)
	}
}

func TestGitHookFilesReturnsNilAndExactReadError(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config read-error receiver"}}
	path := "Git-hook read-error folder"
	readError := errors.New("complete Git-hook-files directory-read dependency error")
	recording := &recordingGitHookFilesystem{
		readDirEntries: []fs.FileInfo{gitHookFileInfo{name: "must-not-be-returned", mode: 0644}},
		readDirErr:     readError,
	}

	actual, err := gitCfg.gitHookFiles(recording.dependencies(), path)

	if actual != nil || err != readError {
		t.Fatalf("Git-hook-files read-error result was (%#v, %v), want (nil, exact error %v)", actual, err, readError)
	}
	paths, populationErr := recording.assertedReadDirPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantRoot := file.Path("%s/%s", gitCfg.Implementation().Dir(), path)
	if !reflect.DeepEqual(paths, []string{wantRoot}) {
		t.Fatalf("Git-hook-files read-error roots were %#v, want %#v", paths, []string{wantRoot})
	}
	if _, populationErr := recording.assertedDeliveredEntries(); populationErr != nil {
		t.Fatal(populationErr)
	}
}

func TestGitHookFilesReturnsNilForEmptyAndAllDirectoryPopulations(t *testing.T) {
	tests := []struct {
		name    string
		entries []fs.FileInfo
	}{
		{name: "empty", entries: []fs.FileInfo{}},
		{name: "all directories", entries: []fs.FileInfo{
			gitHookFileInfo{name: "first-directory", mode: fs.ModeDir | 0755},
			gitHookFileInfo{name: "second-directory", mode: fs.ModeDir | 0701},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete nil-result receiver"}}
			recording := &recordingGitHookFilesystem{readDirEntries: test.entries}

			actual, err := gitCfg.gitHookFiles(recording.dependencies(), "Git-hook nil-result folder")

			if actual != nil || err != nil {
				t.Fatalf("Git-hook-files result was (%#v, %v), want (nil, nil)", actual, err)
			}
			if _, populationErr := recording.assertedReadDirPaths(); populationErr != nil {
				t.Fatal(populationErr)
			}
		})
	}
}

func TestGitHookFilesDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/developer/home/cloud-config/must-not-be-accessed"}}

	actual, err := gitCfg.gitHookFiles(gitHookFilesDependencies{}, "git-hooks")

	if actual != nil || !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe Git-hook-files dependency default returned (%#v, %v), want (nil, %v)", actual, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedGitHookFilesRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingGitHookFilesystem{}

	if _, err := recording.assertedReadDirPaths(); err == nil {
		t.Fatal("empty recorded Git-hook-files directory-read population passed")
	}
	if _, err := recording.assertedDeliveredEntries(); err == nil {
		t.Fatal("empty recorded Git-hook-files entry population passed")
	}
}
