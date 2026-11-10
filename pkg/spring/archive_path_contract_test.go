package spring

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

	"github.com/devdimensionlab/plybuild/internal/adapter/clock"
	"github.com/devdimensionlab/plybuild/internal/adapter/filesystem"
	"github.com/devdimensionlab/plybuild/pkg/file"
)

func (*recordingSpringArchivePathFilesystem) OpenZipEntry(*zip.File) (io.ReadCloser, error) {
	return nil, errors.New("unexpected Spring archive entry open")
}

func (*recordingSpringArchivePathFilesystem) OpenZipReader(string) (*zip.ReadCloser, error) {
	return nil, errors.New("unexpected Spring archive open")
}

func (recording *recordingSpringArchivePathFilesystem) Close(filesystem.File) error {
	return recording.unexpected("close")
}

func (recording *recordingSpringArchivePathFilesystem) CloseReader(io.Closer) error {
	return recording.unexpected("close reader")
}

type recordingSpringArchivePathFilesystem struct {
	workingDirectory         string
	workingDirectoryErr      error
	workingDirectoryAttempts int
	sequence                 *[]string
	unexpectedOperations     []string
}

var _ filesystem.FileSystem = (*recordingSpringArchivePathFilesystem)(nil)

func (recording *recordingSpringArchivePathFilesystem) WorkingDirectory() (string, error) {
	recording.workingDirectoryAttempts++
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "working-directory")
	}
	return recording.workingDirectory, recording.workingDirectoryErr
}

func (recording *recordingSpringArchivePathFilesystem) ReadFile(string) ([]byte, error) {
	return nil, recording.unexpected("read")
}

func (recording *recordingSpringArchivePathFilesystem) ReadDir(string) ([]fs.FileInfo, error) {
	return nil, recording.unexpected("read directory")
}

func (recording *recordingSpringArchivePathFilesystem) ReadDirEntries(string) ([]fs.DirEntry, error) {
	return nil, recording.unexpected("read directory entries")
}

func (recording *recordingSpringArchivePathFilesystem) Stat(string) (fs.FileInfo, error) {
	return nil, recording.unexpected("stat")
}

func (recording *recordingSpringArchivePathFilesystem) Mkdir(string, fs.FileMode) error {
	return recording.unexpected("mkdir")
}

func (recording *recordingSpringArchivePathFilesystem) MkdirAll(string, fs.FileMode) error {
	return recording.unexpected("recursive mkdir")
}

func (recording *recordingSpringArchivePathFilesystem) WriteFile(string, []byte, fs.FileMode) error {
	return recording.unexpected("write file")
}

func (recording *recordingSpringArchivePathFilesystem) OpenFile(string, int, fs.FileMode) (*os.File, error) {
	return nil, recording.unexpected("open file")
}

func (recording *recordingSpringArchivePathFilesystem) Remove(string) error {
	return recording.unexpected("remove")
}

func (recording *recordingSpringArchivePathFilesystem) RemoveAll(string) error {
	return recording.unexpected("recursive remove")
}

func (recording *recordingSpringArchivePathFilesystem) Rename(string, string) error {
	return recording.unexpected("rename")
}

func (recording *recordingSpringArchivePathFilesystem) Glob(string) ([]string, error) {
	return nil, recording.unexpected("glob")
}

func (recording *recordingSpringArchivePathFilesystem) Walk(string, filepath.WalkFunc) error {
	return recording.unexpected("walk")
}

func (recording *recordingSpringArchivePathFilesystem) Create(string) (filesystem.File, error) {
	return nil, recording.unexpected("create")
}

func (recording *recordingSpringArchivePathFilesystem) Copy(filesystem.File, io.Reader) (int64, error) {
	return 0, recording.unexpected("copy")
}

func (recording *recordingSpringArchivePathFilesystem) unexpected(operation string) error {
	recording.unexpectedOperations = append(recording.unexpectedOperations, operation)
	return errors.New("unexpected Spring archive-path " + operation)
}

func (recording *recordingSpringArchivePathFilesystem) dependencies(
	dependencies archivePathDependencies,
) archivePathDependencies {
	dependencies.Files.FileSystem = recording
	return dependencies
}

func (recording *recordingSpringArchivePathFilesystem) assertedWorkingDirectoryAttempts() (int, error) {
	if recording.workingDirectoryAttempts == 0 {
		return 0, errors.New("recorded Spring archive-path working-directory population is empty")
	}
	return recording.workingDirectoryAttempts, nil
}

type recordingSpringArchivePathClock struct {
	now      time.Time
	attempts int
	sequence *[]string
}

func (recording *recordingSpringArchivePathClock) Now() time.Time {
	recording.attempts++
	if recording.sequence != nil {
		*recording.sequence = append(*recording.sequence, "clock")
	}
	return recording.now
}

func (recording *recordingSpringArchivePathClock) dependencies(
	dependencies archivePathDependencies,
) archivePathDependencies {
	dependencies.Clock.Clock = recording
	return dependencies
}

func (recording *recordingSpringArchivePathClock) assertedAttempts() (int, error) {
	if recording.attempts == 0 {
		return 0, errors.New("recorded Spring archive-path clock population is empty")
	}
	return recording.attempts, nil
}

func TestArchivePathSelectsCompleteSystemFilesystemAndClockDependencies(t *testing.T) {
	dependencies := systemArchivePathDependencies()
	systemFiles := filesystem.System()
	systemClock := clock.System()

	if dependencies.Files.FileSystem == nil {
		t.Fatal("archivePath selected an incomplete filesystem dependency")
	}
	if reflect.TypeOf(dependencies.Files.FileSystem) != reflect.TypeOf(systemFiles.FileSystem) {
		t.Fatalf("archivePath filesystem dependency is %T, want %T",
			dependencies.Files.FileSystem, systemFiles.FileSystem)
	}
	if dependencies.Clock.Clock == nil {
		t.Fatal("archivePath selected an incomplete clock dependency")
	}
	if reflect.TypeOf(dependencies.Clock.Clock) != reflect.TypeOf(systemClock.Clock) {
		t.Fatalf("archivePath clock dependency is %T, want %T",
			dependencies.Clock.Clock, systemClock.Clock)
	}
}

func TestSpringArchivePathRecordingDoublesPreserveCompleteCallerDependencies(t *testing.T) {
	callerFiles := &recordingSpringArchivePathFilesystem{}
	callerClock := &recordingSpringArchivePathClock{}
	initial := archivePathDependencies{
		Files: filesystem.Dependencies{FileSystem: callerFiles},
		Clock: clock.Dependencies{Clock: callerClock},
	}
	filesRecording := &recordingSpringArchivePathFilesystem{}
	clockRecording := &recordingSpringArchivePathClock{}

	withFiles := filesRecording.dependencies(initial)
	withClock := clockRecording.dependencies(initial)

	if withFiles.Files.FileSystem != filesRecording || withFiles.Clock.Clock != callerClock {
		t.Fatalf("filesystem recorder lost the complete caller-owned dependency: %#v", withFiles)
	}
	if withClock.Files.FileSystem != callerFiles || withClock.Clock.Clock != clockRecording {
		t.Fatalf("clock recorder lost the complete caller-owned dependency: %#v", withClock)
	}
}

func TestArchivePathPreservesArbitraryDirectoryInjectedUnixSecondOrderAndExactReturns(t *testing.T) {
	arbitraryDirectory := string([]byte{
		'/', 'a', 'r', 'b', 'i', 't', 'r', 'a', 'r', 'y', '/', '/', 'w', 'o', 'r', 'k', 'i', 'n', 'g', ' ',
		'd', 'i', 'r', '/', '.', '.', '/', 'b', 'a', 'c', 'k', 's', 'l', 'a', 's', 'h', '\\', 0xc3, 0xb8, '\n', 0x00,
	})
	tests := []struct {
		name             string
		workingDirectory string
		now              time.Time
	}{
		{
			name:             "arbitrary directory and positive subsecond",
			workingDirectory: arbitraryDirectory,
			now:              time.Unix(1700000000, 999999999),
		},
		{
			name:             "pre-epoch subsecond keeps negative Unix floor",
			workingDirectory: `/pre-epoch//directory/../segment\-ø`,
			now:              time.Unix(-1, 500000000),
		},
		{
			name:             "arbitrary timezone keeps absolute Unix second",
			workingDirectory: "relative directory remains relative",
			now: time.Date(2031, time.April, 5, 6, 7, 8, 123456789,
				time.FixedZone("arbitrary injected zone", 9*60*60+17*60)),
		},
	}

	if len(tests) == 0 {
		t.Fatal("TestArchivePathPreservesArbitraryDirectoryInjectedUnixSecondOrderAndExactReturns test-case population is empty")
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sequence := []string{}
			files := &recordingSpringArchivePathFilesystem{
				workingDirectory: test.workingDirectory,
				sequence:         &sequence,
			}
			clockRecording := &recordingSpringArchivePathClock{now: test.now, sequence: &sequence}
			dependencies := clockRecording.dependencies(files.dependencies(archivePathDependencies{}))
			if dependencies.Files.FileSystem != files || dependencies.Clock.Clock != clockRecording {
				t.Fatalf("archive-path dependency delivery is incomplete: %#v", dependencies)
			}

			archive, err := archivePathWithDependencies(dependencies)

			if err != nil {
				t.Fatalf("archive-path composition returned an error: %v", err)
			}
			want := file.Path("%s/spring-%d.zip", test.workingDirectory, test.now.Unix())
			if archive != want {
				t.Fatalf("archive path was %q, want exact directory and injected Unix-second path %q", archive, want)
			}
			assertOneSpringArchivePathWorkingDirectoryAttempt(t, files)
			assertOneSpringArchivePathClockAttempt(t, clockRecording)
			if !reflect.DeepEqual(sequence, []string{"working-directory", "clock"}) {
				t.Fatalf("archive-path effect order was %#v, want directory before clock", sequence)
			}
			if len(files.unexpectedOperations) != 0 {
				t.Fatalf("archive-path composition invoked unrelated filesystem operations: %#v",
					files.unexpectedOperations)
			}
		})
	}
}

func TestArchivePathReturnsEmptyPathAndExactWorkingDirectoryErrorAfterOneAttempt(t *testing.T) {
	directoryError := errors.New("complete Spring archive-path working-directory error")
	recording := &recordingSpringArchivePathFilesystem{
		workingDirectory:    "/partial directory must be discarded",
		workingDirectoryErr: directoryError,
		sequence:            &[]string{},
	}
	clockRecording := &recordingSpringArchivePathClock{now: time.Unix(1700000001, 0), sequence: recording.sequence}

	archive, err := archivePathWithDependencies(clockRecording.dependencies(recording.dependencies(archivePathDependencies{})))

	if archive != "" || err != directoryError {
		t.Fatalf("archive-path error result was (%q, %v), want (empty, exact %v)",
			archive, err, directoryError)
	}
	assertOneSpringArchivePathWorkingDirectoryAttempt(t, recording)
	if clockRecording.attempts != 0 {
		t.Fatalf("archive-path read the clock %d times after a directory error, want 0", clockRecording.attempts)
	}
	if !reflect.DeepEqual(*recording.sequence, []string{"working-directory"}) {
		t.Fatalf("archive-path directory-error sequence was %#v, want only working-directory", *recording.sequence)
	}
	if len(recording.unexpectedOperations) != 0 {
		t.Fatalf("archive-path error invoked unrelated filesystem operations: %#v",
			recording.unexpectedOperations)
	}
}

func TestArchivePathDependenciesDefaultToSafeEmptyResult(t *testing.T) {
	archive, err := archivePathWithDependencies(archivePathDependencies{})

	if archive != "" || err != filesystem.ErrNoFilesystem {
		t.Fatalf("safe archive-path dependency default returned (%q, %v), want (empty, %v)",
			archive, err, filesystem.ErrNoFilesystem)
	}
}

func TestRecordedSpringArchivePathRejectsEmptyWorkingDirectoryPopulation(t *testing.T) {
	recording := &recordingSpringArchivePathFilesystem{}

	if _, err := recording.assertedWorkingDirectoryAttempts(); err == nil {
		t.Fatal("empty recorded Spring archive-path working-directory population passed")
	}
}

func TestRecordedSpringArchivePathRejectsEmptyClockPopulation(t *testing.T) {
	recording := &recordingSpringArchivePathClock{}

	if _, err := recording.assertedAttempts(); err == nil {
		t.Fatal("empty recorded Spring archive-path clock population passed")
	}
}

func assertOneSpringArchivePathWorkingDirectoryAttempt(t *testing.T, recording *recordingSpringArchivePathFilesystem) {
	t.Helper()
	attempts, err := recording.assertedWorkingDirectoryAttempts()
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 1 {
		t.Fatalf("archive-path composition made %d working-directory attempts, want 1", attempts)
	}
}

func assertOneSpringArchivePathClockAttempt(t *testing.T, recording *recordingSpringArchivePathClock) {
	t.Helper()
	attempts, err := recording.assertedAttempts()
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 1 {
		t.Fatalf("archive-path composition made %d clock attempts, want 1", attempts)
	}
}
