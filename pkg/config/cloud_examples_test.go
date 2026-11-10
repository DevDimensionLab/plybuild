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

func (*recordingExamplesFilesystem) OpenZipEntry(*zip.File) (io.ReadCloser, error) {
	return nil, errors.New("unexpected examples archive entry open")
}

func (*recordingExamplesFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected examples archive open")
}

func (*recordingExamplesFilesystem) Close(filesystem.File) error {
	return errors.New("unexpected Examples close")
}

func (*recordingExamplesFilesystem) CloseReader(io.Closer) error {
	return errors.New("unexpected Examples reader close")
}

type recordingExamplesFilesystem struct {
	readDirPaths     []string
	readDirEntries   []fs.FileInfo
	deliveredEntries []fs.FileInfo
	readDirErr       error
}

func (*recordingExamplesFilesystem) WorkingDirectory() (string, error) {
	return "", errors.New("unexpected Examples working directory")
}

func (*recordingExamplesFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected Examples read")
}

func (recording *recordingExamplesFilesystem) ReadDir(path string) ([]fs.FileInfo, error) {
	recording.readDirPaths = append(recording.readDirPaths, path)
	recording.deliveredEntries = append(recording.deliveredEntries, recording.readDirEntries...)
	return append([]fs.FileInfo{}, recording.readDirEntries...), recording.readDirErr
}

func (*recordingExamplesFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, errors.New("unexpected Examples read directory entries")
}

func (*recordingExamplesFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected Examples stat")
}

func (*recordingExamplesFilesystem) Mkdir(string, fs.FileMode) error {
	return errors.New("unexpected Examples mkdir")
}

func (*recordingExamplesFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected Examples mkdir")
}

func (*recordingExamplesFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected Examples write")
}

func (*recordingExamplesFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, errors.New("unexpected Examples open file")
}

func (*recordingExamplesFilesystem) Remove(string) error {
	return errors.New("unexpected Examples remove")
}

func (*recordingExamplesFilesystem) RemoveAll(string) error {
	return errors.New("unexpected Examples recursive remove")
}

func (*recordingExamplesFilesystem) Rename(string, string) error {
	return errors.New("unexpected Examples rename")
}

func (*recordingExamplesFilesystem) Glob(string) ([]string, error) {
	return nil, errors.New("unexpected Examples glob")
}

func (*recordingExamplesFilesystem) Walk(string, filepath.WalkFunc) error {
	return errors.New("unexpected Examples walk")
}

func (*recordingExamplesFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected Examples create")
}

func (*recordingExamplesFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected Examples copy")
}

func (recording *recordingExamplesFilesystem) dependencies() examplesDependencies {
	return examplesDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingExamplesFilesystem) assertedReadDirPaths() ([]string, error) {
	if len(recording.readDirPaths) == 0 {
		return nil, errors.New("recorded Examples directory-read population is empty")
	}
	return recording.readDirPaths, nil
}

func (recording *recordingExamplesFilesystem) assertedDeliveredEntries() ([]fs.FileInfo, error) {
	if len(recording.deliveredEntries) == 0 {
		return nil, errors.New("recorded Examples entry population is empty")
	}
	return recording.deliveredEntries, nil
}

type examplesFileInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	system  interface{}
}

func (info examplesFileInfo) Name() string       { return info.name }
func (info examplesFileInfo) Size() int64        { return info.size }
func (info examplesFileInfo) Mode() fs.FileMode  { return info.mode }
func (info examplesFileInfo) ModTime() time.Time { return info.modTime }
func (info examplesFileInfo) IsDir() bool        { return info.mode.IsDir() }
func (info examplesFileInfo) Sys() interface{}   { return info.system }

func TestExamplesSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemExamplesDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Examples selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Examples filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestExamplesPreservesReceiverRootDeliveredOrderMetadataClassificationAndDirectoryNameResults(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config receiver/root with spaces"}}
	firstTime := time.Unix(1700000101, 0)
	secondTime := time.Unix(1700000102, 0)
	entries := []fs.FileInfo{
		examplesFileInfo{name: "z-first-example", size: 41, mode: fs.ModeDir | 0751, modTime: firstTime, system: "first metadata"},
		examplesFileInfo{name: "ignored-regular-file.json", size: 82, mode: 0641, system: "file metadata"},
		examplesFileInfo{name: "a-second-example", size: 83, mode: fs.ModeDir | 0705, modTime: secondTime, system: "second metadata"},
		examplesFileInfo{name: "ignored-symbolic-link", size: 84, mode: fs.ModeSymlink | 0777, system: "link metadata"},
		examplesFileInfo{name: "middle-example", size: 85, mode: fs.ModeDir | 0711},
	}
	recording := &recordingExamplesFilesystem{readDirEntries: entries}
	dependencies := recording.dependencies()
	if dependencies.Files.FileSystem != recording {
		t.Fatalf("Examples dependency lost its complete filesystem value: %#v", dependencies)
	}

	actual, err := gitCfg.examples(dependencies)

	if err != nil {
		t.Fatalf("Examples returned an error: %v", err)
	}
	wantRoot := file.Path("%s/examples", gitCfg.Implementation().Dir())
	paths, populationErr := recording.assertedReadDirPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(paths, []string{wantRoot}) {
		t.Fatalf("Examples directory-read roots were %#v, want exact receiver-derived root %#v", paths, []string{wantRoot})
	}
	deliveredEntries, populationErr := recording.assertedDeliveredEntries()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	if !reflect.DeepEqual(deliveredEntries, entries) {
		t.Fatalf("Examples entries were %#v, want ordered metadata %#v", deliveredEntries, entries)
	}
	want := []string{"z-first-example", "a-second-example", "middle-example"}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("Examples results were %#v, want directory-name-only delivered order %#v", actual, want)
	}
}

func TestExamplesReturnsNilAndExactReadError(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete cloud-config read-error receiver"}}
	readError := errors.New("complete Examples directory-read dependency error")
	recording := &recordingExamplesFilesystem{
		readDirEntries: []fs.FileInfo{examplesFileInfo{name: "must-not-be-returned", mode: fs.ModeDir | 0755}},
		readDirErr:     readError,
	}

	actual, err := gitCfg.examples(recording.dependencies())

	if actual != nil || err != readError {
		t.Fatalf("Examples read-error result was (%#v, %v), want (nil, exact error %v)", actual, err, readError)
	}
	paths, populationErr := recording.assertedReadDirPaths()
	if populationErr != nil {
		t.Fatal(populationErr)
	}
	wantRoot := file.Path("%s/examples", gitCfg.Implementation().Dir())
	if !reflect.DeepEqual(paths, []string{wantRoot}) {
		t.Fatalf("Examples read-error roots were %#v, want %#v", paths, []string{wantRoot})
	}
	if _, populationErr := recording.assertedDeliveredEntries(); populationErr != nil {
		t.Fatal(populationErr)
	}
}

func TestExamplesReturnsNilForEmptyAndAllFilePopulations(t *testing.T) {
	tests := []struct {
		name    string
		entries []fs.FileInfo
	}{
		{name: "empty", entries: []fs.FileInfo{}},
		{name: "all files", entries: []fs.FileInfo{
			examplesFileInfo{name: "first-file.json", mode: 0644},
			examplesFileInfo{name: "second-link", mode: fs.ModeSymlink | 0777},
		}},
	}

	if len(tests) == 0 {
		t.Fatal("TestExamplesReturnsNilForEmptyAndAllFilePopulations test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/complete Examples nil-result receiver"}}
			recording := &recordingExamplesFilesystem{readDirEntries: test.entries}

			actual, err := gitCfg.examples(recording.dependencies())

			if actual != nil || err != nil {
				t.Fatalf("Examples result was (%#v, %v), want (nil, nil)", actual, err)
			}
			if _, populationErr := recording.assertedReadDirPaths(); populationErr != nil {
				t.Fatal(populationErr)
			}
		})
	}
}

func TestExamplesDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	gitCfg := GitCloudConfig{Impl: DirConfig{Path: "/developer/home/cloud-config/must-not-be-accessed"}}

	actual, err := gitCfg.examples(examplesDependencies{})

	if actual != nil || !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe Examples dependency default returned (%#v, %v), want (nil, %v)", actual, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedExamplesRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingExamplesFilesystem{}

	if _, err := recording.assertedReadDirPaths(); err == nil {
		t.Fatal("empty recorded Examples directory-read population passed")
	}
	if _, err := recording.assertedDeliveredEntries(); err == nil {
		t.Fatal("empty recorded Examples entry population passed")
	}
}
