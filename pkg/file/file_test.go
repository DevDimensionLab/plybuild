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

type recordedFileCreate struct {
	path string
	data []byte
	mode fs.FileMode
}

type recordingFileCreateFilesystem struct {
	writes   []recordedFileCreate
	writeErr error
}

type recordedDirectoryCreate struct {
	path string
	mode fs.FileMode
}

type recordingDirectoryCreateFilesystem struct {
	statPaths []string
	creates   []recordedDirectoryCreate
	statInfo  fs.FileInfo
	statErr   error
	mkdirErr  error
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

func (*recordingFileCreateFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected file-create read")
}

func (*recordingFileCreateFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, errors.New("unexpected file-create stat")
}

func (*recordingFileCreateFilesystem) MkdirAll(string, fs.FileMode) error {
	return errors.New("unexpected file-create mkdir")
}

func (recording *recordingFileCreateFilesystem) WriteFile(path string, data []byte, mode fs.FileMode) error {
	recording.writes = append(recording.writes, recordedFileCreate{
		path: path,
		data: append([]byte{}, data...),
		mode: mode,
	})
	return recording.writeErr
}

func (*recordingFileCreateFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected file-create create")
}

func (*recordingFileCreateFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected file-create copy")
}

func (recording *recordingFileCreateFilesystem) dependencies() createFileDependencies {
	return createFileDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingFileCreateFilesystem) assertedWrites() ([]recordedFileCreate, error) {
	if len(recording.writes) == 0 {
		return nil, errors.New("recorded file-create population is empty")
	}
	return recording.writes, nil
}

func (*recordingDirectoryCreateFilesystem) ReadFile(string) ([]byte, error) {
	return nil, errors.New("unexpected directory-create read")
}

func (recording *recordingDirectoryCreateFilesystem) Stat(path string) (fs.FileInfo, error) {
	recording.statPaths = append(recording.statPaths, path)
	return recording.statInfo, recording.statErr
}

func (recording *recordingDirectoryCreateFilesystem) MkdirAll(path string, mode fs.FileMode) error {
	recording.creates = append(recording.creates, recordedDirectoryCreate{path: path, mode: mode})
	return recording.mkdirErr
}

func (*recordingDirectoryCreateFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return errors.New("unexpected directory-create write")
}

func (*recordingDirectoryCreateFilesystem) Create(string) (filesystem.File, error) {
	return nil, errors.New("unexpected directory-create create")
}

func (*recordingDirectoryCreateFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, errors.New("unexpected directory-create copy")
}

func (recording *recordingDirectoryCreateFilesystem) dependencies() createDirectoryDependencies {
	return createDirectoryDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingDirectoryCreateFilesystem) assertedStatPaths() ([]string, error) {
	if len(recording.statPaths) == 0 {
		return nil, errors.New("recorded directory-create stat population is empty")
	}
	return recording.statPaths, nil
}

func (recording *recordingDirectoryCreateFilesystem) assertedCreates() ([]recordedDirectoryCreate, error) {
	if len(recording.creates) == 0 {
		return nil, errors.New("recorded directory-create creation population is empty")
	}
	return recording.creates, nil
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

func TestCreateFileSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemCreateFileDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("CreateFile selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("CreateFile filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestCreateFilePreservesCompletePathContentBytesModeAndDependencyError(t *testing.T) {
	path := "/complete file-create/path with spaces/file.txt"
	writeError := errors.New("complete file-create dependency error")
	tests := []struct {
		name      string
		content   string
		wantBytes []byte
		writeErr  error
	}{
		{name: "empty content", content: "", wantBytes: []byte{}},
		{
			name:      "multiline content and dependency error",
			content:   "first complete line\n\nthird complete line\n",
			wantBytes: []byte("first complete line\n\nthird complete line\n"),
			writeErr:  writeError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingFileCreateFilesystem{writeErr: test.writeErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("file-create dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := createFile(dependencies, path, test.content)
			if !errors.Is(err, test.writeErr) {
				t.Fatalf("create-file error was %v, want %v", err, test.writeErr)
			}
			writes, populationErr := recording.assertedWrites()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			want := []recordedFileCreate{{path: path, data: test.wantBytes, mode: 0644}}
			if !reflect.DeepEqual(writes, want) {
				t.Fatalf("recorded file-create differs:\n got: %#v\nwant: %#v", writes, want)
			}
		})
	}
}

func TestCreateFileDependenciesDefaultToSafeNoMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist.txt")

	err := createFile(createFileDependencies{}, target, "must not be written")

	if !errors.Is(err, filesystem.ErrNoFilesystem) {
		t.Fatalf("safe file-create dependency default returned %v, want %v", err, filesystem.ErrNoFilesystem)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe file-create dependency default mutated the target: %v", statErr)
	}
}

func TestRecordedFileCreateRejectsEmptyPopulation(t *testing.T) {
	recording := &recordingFileCreateFilesystem{}

	if _, err := recording.assertedWrites(); err == nil {
		t.Fatal("empty recorded file-create population passed")
	}
}

func TestCreateDirectorySelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemCreateDirectoryDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("CreateDirectory selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("CreateDirectory filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestCreateDirectoryPreservesCompletePathMissingSelectionModeAndLegacyErrors(t *testing.T) {
	path := "/complete directory-create/path with spaces/nested"
	missingError := &os.PathError{Op: "stat", Path: path, Err: fs.ErrNotExist}
	mkdirError := errors.New("complete directory-create dependency error")
	tests := []struct {
		name      string
		mkdirErr  error
		wantError error
	}{
		{name: "missing path is created"},
		{
			name:      "creation failure returns original missing-path stat error",
			mkdirErr:  mkdirError,
			wantError: missingError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingDirectoryCreateFilesystem{statErr: missingError, mkdirErr: test.mkdirErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("directory-create dependency lost its complete filesystem value: %#v", dependencies)
			}

			err := createDirectoryWithDependencies(dependencies, path)
			if err != test.wantError {
				t.Fatalf("directory-create error was %v, want original stat error %v", err, test.wantError)
			}
			statPaths, populationErr := recording.assertedStatPaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(statPaths, []string{path}) {
				t.Fatalf("directory-create dependency received stat paths %#v, want %#v", statPaths, []string{path})
			}
			creates, populationErr := recording.assertedCreates()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			wantCreates := []recordedDirectoryCreate{{path: path, mode: 0755}}
			if !reflect.DeepEqual(creates, wantCreates) {
				t.Fatalf("recorded directory creation differs:\n got: %#v\nwant: %#v", creates, wantCreates)
			}
		})
	}
}

func TestCreateDirectorySkipsCreationForExistingPathAndOtherStatErrors(t *testing.T) {
	path := "/complete directory-create/path with spaces/existing"
	otherError := errors.New("permission denied while probing complete directory path")
	tests := []struct {
		name string
		info fs.FileInfo
		err  error
	}{
		{name: "existing path", info: existingFileInfo{name: "existing"}},
		{name: "non-missing stat error", err: otherError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingDirectoryCreateFilesystem{statInfo: test.info, statErr: test.err}

			if err := createDirectoryWithDependencies(recording.dependencies(), path); err != nil {
				t.Fatalf("directory-create returned %v, want nil", err)
			}
			statPaths, populationErr := recording.assertedStatPaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(statPaths, []string{path}) {
				t.Fatalf("directory-create dependency received stat paths %#v, want %#v", statPaths, []string{path})
			}
			if len(recording.creates) != 0 {
				t.Fatalf("directory-create selected unexpected creations: %#v", recording.creates)
			}
		})
	}
}

func TestCreateDirectoryDependenciesDefaultToSafeNoMutation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe default", "must-not-exist")

	if err := createDirectoryWithDependencies(createDirectoryDependencies{}, target); err != nil {
		t.Fatalf("safe directory-create dependency default returned %v, want nil", err)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("safe directory-create dependency default mutated the target: %v", statErr)
	}
}

func TestRecordedDirectoryCreateRejectsEmptyPopulations(t *testing.T) {
	recording := &recordingDirectoryCreateFilesystem{}

	if _, err := recording.assertedStatPaths(); err == nil {
		t.Fatal("empty recorded directory-create stat population passed")
	}
	if _, err := recording.assertedCreates(); err == nil {
		t.Fatal("empty recorded directory-create creation population passed")
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
