package file

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
)

type recordingFileExistenceFilesystem struct {
	statPaths []string
	statInfo  fs.FileInfo
	statErr   error
}

type recordedFileOverwrite struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingFileOverwriteFilesystem struct {
	writes   []recordedFileOverwrite
	writeErr error
}

type existingFileInfo struct {
	name string
}

func (info existingFileInfo) Name() string  { return info.name }
func (existingFileInfo) Size() int64        { return 0 }
func (existingFileInfo) Mode() fs.FileMode  { return 0644 }
func (existingFileInfo) ModTime() time.Time { return time.Time{} }
func (existingFileInfo) IsDir() bool        { return false }
func (existingFileInfo) Sys() interface{}   { return nil }

func (*recordingFileExistenceFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected file-existence read")
}

func (recording *recordingFileExistenceFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.statPaths = append(recording.statPaths, path)
	return recording.statInfo, recording.statErr
}

func (*recordingFileExistenceFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected file-existence mkdir")
}

func (*recordingFileExistenceFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected file-existence write")
}

func (*recordingFileExistenceFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected file-existence create")
}

func (*recordingFileExistenceFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected file-existence copy")
}

func (recording *recordingFileExistenceFilesystem) dependencies() existsDependencies {
	return existsDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileExistenceFilesystem) assertedStatPaths() ([]string, error) {
	if len(recording.statPaths) == 0 {
		return nil, errors.New("recorded file-existence stat population is empty")
	}
	return recording.statPaths, nil
}

func (*recordingFileOverwriteFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected file-overwrite read")
}

func (*recordingFileOverwriteFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected file-overwrite stat")
}

func (*recordingFileOverwriteFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected file-overwrite mkdir")
}

func (recording *recordingFileOverwriteFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedFileOverwrite{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (*recordingFileOverwriteFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected file-overwrite create")
}

func (*recordingFileOverwriteFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected file-overwrite copy")
}

func (recording *recordingFileOverwriteFilesystem) dependencies() overwriteDependencies {
	return overwriteDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileOverwriteFilesystem) assertedWrites() ([]recordedFileOverwrite, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded file-overwrite population is empty")
	}
	return recording.writes, nil
}

func TestExistsSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemExistsDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Exists selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Exists filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestExistsPreservesCompletePathAndLegacyStatResults(t *testing.T) {
	path := "/complete file-existence/path with spaces/file.txt"
	otherError := errors.New("permission denied while probing complete path")
	tests := []struct {
		name     string
		info     fs.FileInfo
		err      error
		expected bool
	}{
		{
			name: "missing path is absent",
			err:  &os.PathError{Op: "stat", Path: path, Err: fs.ErrNotExist},
		},
		{
			name:     "existing path is present",
			info:     existingFileInfo{name: "file.txt"},
			expected: true,
		},
		{
			name:     "another stat error is present",
			err:      otherError,
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileExistenceFilesystem{statInfo: test.info, statErr: test.err}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("file-existence dependency lost its complete filesystem value: %#v", dependencies)
			}

			if actual := exists(dependencies, path); actual != test.expected {
				t.Fatalf("exists result was %t, want %t for stat error %v", actual, test.expected, test.err)
			}
			statPaths, err := recording.assertedStatPaths()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(statPaths, []string{path}) {
				t.Fatalf("file-existence dependency received paths %#v, want %#v", statPaths, []string{path})
			}
		})
	}
}

func TestExistsDependenciesDefaultToSafeAbsent(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist.txt")

	if exists(existsDependencies{}, target) {
		t.Fatal("zero-value file-existence dependency reported an unprobed path as present")
	}
}

func TestRecordedFileExistenceRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileExistenceFilesystem{}

	if _, err := recording.assertedStatPaths(); err == nil {
		t.Fatal("empty recorded file-existence stat population passed")
	}
}

func TestOverwriteSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemOverwriteDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("Overwrite selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("Overwrite filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestOverwritePreservesCompleteLinesPathBytesModeAndDependencyError(t *testing.T) {
	path := "/complete file-overwrite/path with spaces/file.txt"
	writeError := errors.New("complete file-overwrite dependency error")
	tests := []struct {
		name      string
		lines     []string
		wantBytes []byte
		writeErr  error
	}{
		{name: "empty lines", lines: []string{}, wantBytes: []byte{}},
		{name: "single line", lines: []string{"single complete line"}, wantBytes: []byte("single complete line")},
		{
			name:      "multiple lines and dependency error",
			lines:     []string{"first complete line", "", "third complete line", ""},
			wantBytes: []byte("first complete line\n\nthird complete line\n"),
			writeErr:  writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileOverwriteFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("file-overwrite dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := overwrite(dependencies, test.lines, path)
			if !errors.Is(err, test.writeErr) {
				t.Fatalf("overwrite error was %v, want %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			want := []recordedFileOverwrite{{path: path, data: test.wantBytes, mode: 0644}}
			if !reflect.DeepEqual(writes, want) {
				t.Fatalf("recorded file-overwrite differs:\n got: %#v\nwant: %#v", writes, want)
			}
		})
	}
}

func TestOverwriteDependenciesDefaultToSafeNoMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist.txt")

	err := overwrite(overwriteDependencies{}, []string{"must not be written"}, target)

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe file-overwrite dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe file-overwrite dependency default mutated the target: %v", statErr)
	}
}

func TestRecordedFileOverwriteRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileOverwriteFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded file-overwrite population passed")
	}
}

func TestRelPath(t *testing.T) {
	relPath, err := RelPath(
		"/home/user/.ply/cloud-config/templates/flyway-demo",
		"/home/user/.ply/cloud-config/templates/flyway-demo/src/main/kotlin/no/ply/template/demo/flyway/Queue.kt")

	expected := "src/main/kotlin/no/ply/template/demo/flyway/Queue.kt"

	if err != nil {
		t.Errorf("%v\n", err)
	}

	if relPath != expected {
		t.Errorf("%s is not %s", relPath, expected)
	}
}
