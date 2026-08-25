package cmd

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"github.com/devdimensionlab/plybuild/pkg/file"
	"github.com/devdimensionlab/plybuild/pkg/tips"
)

type recordingTipsShowReadFilesystem struct {
	readPaths            []string
	readData             []byte
	readErr              error
	unexpectedOperations []string
}

func (recording *recordingTipsShowReadFilesystem) WorkingDirectory() (string, error) {
	return "", recording.unexpected("working directory")
}

func (recording *recordingTipsShowReadFilesystem) ReadFile(path string) ([]byte, error) {
	recording.readPaths = append(recording.readPaths, path)
	return recording.readData, recording.readErr
}

func (recording *recordingTipsShowReadFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingTipsShowReadFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingTipsShowReadFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingTipsShowReadFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingTipsShowReadFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingTipsShowReadFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingTipsShowReadFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingTipsShowReadFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingTipsShowReadFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingTipsShowReadFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingTipsShowReadFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingTipsShowReadFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingTipsShowReadFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingTipsShowReadFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingTipsShowReadFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected tips-show " + operation)
}

func (recording *recordingTipsShowReadFilesystem) dependencies() tipsShowReadDependencies {
	return tipsShowReadDependencies{Files: filesystem.Dependencies{FileSystem: recording}}
}

func (recording *recordingTipsShowReadFilesystem) assertedReadPaths() ([]string, error) {
	if len(recording.readPaths) == 0 {
		return nil, errors.New("recorded tips-show read population is empty")
	}
	return recording.readPaths, nil
}

func TestTipsShowReadSelectsCompleteSystemFilesystemDependencies(t *testing.T) {
	dependencies := systemTipsShowReadDependencies()
	systemFiles := filesystem.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("tips-show selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("tips-show filesystem dependency is %T, want %T", dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
}

func TestTipsShowReadPreservesExactCallerComposedPathBytesAndDependencyError(t *testing.T) {
	root := `/complete cloud root//with spaces/../backslash\segment-ø`
	name := `../arbitrary tip/name\with spaces-β`
	cloudConfig := config.GitCloudConfig{Impl: config.DirConfig{Path: root}}
	tipsPath := file.Path("%s/%s.md", tips.LocalDir(cloudConfig), name)
	readError := errors.New("complete tips-show read dependency error")
	tests := []struct {
		name     string
		readData []byte
		readErr  error
	}{
		{name: "empty bytes", readData: []byte{}},
		{name: "representative bytes", readData: []byte("# Complete tip\n\nexact body\n")},
		{name: "non-ASCII and arbitrary bytes", readData: []byte{0x00, 't', 0xc3, 0xb8, 0xff, '\n'}},
		{
			name:     "partial bytes and dependency error",
			readData: []byte("partial source bytes"),
			readErr:  readError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recording := &recordingTipsShowReadFilesystem{readData: test.readData, readErr: test.readErr}
			dependencies := recording.dependencies()
			if dependencies.Files.FileSystem != recording {
				t.Fatalf("tips-show dependency lost its complete filesystem value: %#v", dependencies)
			}

			data, err := readTipsShowSource(dependencies, tipsPath)

			if err != test.readErr {
				t.Fatalf("tips-show read error was %v, want exact dependency error %v", err, test.readErr)
			}
			if !reflect.DeepEqual(data, test.readData) {
				t.Fatalf("tips-show source bytes were %#v, want exact %#v", data, test.readData)
			}
			readPaths, populationErr := recording.assertedReadPaths()
			if populationErr != nil {
				t.Fatal(populationErr)
			}
			if !reflect.DeepEqual(readPaths, []string{tipsPath}) {
				t.Fatalf("tips-show dependency received paths %#v, want exact caller-composed %#v", readPaths, []string{tipsPath})
			}
			if len(recording.unexpectedOperations) != 0 {
				t.Fatalf("tips-show invoked unrelated filesystem operations: %#v", recording.unexpectedOperations)
			}
		})
	}
}

func TestTipsShowReadDependenciesDefaultToSafeNoDeveloperPathAccess(t *testing.T) {
	tipsPath := "/developer/home/project/must-not-be-accessed/tips/unsafe.md"

	data, err := readTipsShowSource(tipsShowReadDependencies{}, tipsPath)

	if data != nil || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe tips-show dependency default returned (%#v, %v), want (nil, %v)", data, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedTipsShowReadRejectsEmptyReadPopulation(t *testing.T) {
	recording := &recordingTipsShowReadFilesystem{}

	if _, err := recording.assertedReadPaths(); err == nil {
		t.Fatal("empty recorded tips-show read population passed")
	}
}
