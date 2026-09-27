package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

var errInjected = errors.New("injected failure")

type recordingFileSystem struct {
	cwd             string
	operations      []string
	fail            map[string]error
	requestedFlags  []int
	requestedOpen   []fs.FileMode
	requestedChmod  map[string]fs.FileMode
	mutations       int
	shortWrite      bool
	renameHook      func(string, string) error
	removeAllCalled []string
}

func newRecordingFileSystem(root string) *recordingFileSystem {
	return &recordingFileSystem{
		cwd:            root,
		fail:           make(map[string]error),
		requestedChmod: make(map[string]fs.FileMode),
	}
}

func (filesystem *recordingFileSystem) record(operation, path string) error {
	entry := operation
	if path != "" {
		entry += ":" + path
	}
	filesystem.operations = append(filesystem.operations, entry)
	if err := filesystem.fail[entry]; err != nil {
		return err
	}
	return filesystem.fail[operation]
}

func (filesystem *recordingFileSystem) Getwd() (string, error) {
	if err := filesystem.record("getwd", ""); err != nil {
		return "", err
	}
	return filesystem.cwd, nil
}

func (filesystem *recordingFileSystem) EvalSymlinks(path string) (string, error) {
	if err := filesystem.record("eval", path); err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(path)
}

func (filesystem *recordingFileSystem) Lstat(path string) (fs.FileInfo, error) {
	if err := filesystem.record("lstat", path); err != nil {
		return nil, err
	}
	return os.Lstat(path)
}

func (filesystem *recordingFileSystem) Stat(path string) (fs.FileInfo, error) {
	if err := filesystem.record("stat", path); err != nil {
		return nil, err
	}
	return os.Stat(path)
}

func (filesystem *recordingFileSystem) ReadFile(path string) ([]byte, error) {
	if err := filesystem.record("read", path); err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func (filesystem *recordingFileSystem) MkdirTemp(directory, pattern string) (string, error) {
	filesystem.mutations++
	if err := filesystem.record("mkdir-temp", directory); err != nil {
		return "", err
	}
	return os.MkdirTemp(directory, pattern)
}

func (filesystem *recordingFileSystem) OpenFile(path string, flag int, mode fs.FileMode) (File, error) {
	filesystem.mutations++
	filesystem.requestedFlags = append(filesystem.requestedFlags, flag)
	filesystem.requestedOpen = append(filesystem.requestedOpen, mode)
	if err := filesystem.record("open", path); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, flag, mode)
	if err != nil {
		return nil, err
	}
	return &recordingFile{File: file, filesystem: filesystem}, nil
}

func (filesystem *recordingFileSystem) Chmod(path string, mode fs.FileMode) error {
	filesystem.mutations++
	filesystem.requestedChmod[path] = mode
	if err := filesystem.record("dir-chmod", path); err != nil {
		return err
	}
	return os.Chmod(path, mode)
}

func (filesystem *recordingFileSystem) Rename(oldPath, newPath string) error {
	filesystem.mutations++
	if err := filesystem.record("rename", newPath); err != nil {
		return err
	}
	if filesystem.renameHook != nil {
		return filesystem.renameHook(oldPath, newPath)
	}
	return os.Rename(oldPath, newPath)
}

func (filesystem *recordingFileSystem) RemoveAll(path string) error {
	filesystem.mutations++
	filesystem.removeAllCalled = append(filesystem.removeAllCalled, path)
	if err := filesystem.record("remove", path); err != nil {
		return err
	}
	return os.RemoveAll(path)
}

type recordingFile struct {
	File       *os.File
	filesystem *recordingFileSystem
}

func (file *recordingFile) Write(contents []byte) (int, error) {
	file.filesystem.mutations++
	if err := file.filesystem.record("write", file.File.Name()); err != nil {
		return 0, err
	}
	if file.filesystem.shortWrite {
		return file.File.Write(contents[:len(contents)-1])
	}
	return file.File.Write(contents)
}

func (file *recordingFile) Chmod(mode fs.FileMode) error {
	file.filesystem.mutations++
	file.filesystem.requestedChmod[file.File.Name()] = mode
	if err := file.filesystem.record("file-chmod", file.File.Name()); err != nil {
		return err
	}
	return file.File.Chmod(mode)
}

func (file *recordingFile) Sync() error {
	file.filesystem.mutations++
	if err := file.filesystem.record("sync", file.File.Name()); err != nil {
		return err
	}
	return file.File.Sync()
}

func (file *recordingFile) Close() error {
	file.filesystem.mutations++
	recordedErr := file.filesystem.record("close", file.File.Name())
	closeErr := file.File.Close()
	if recordedErr != nil {
		return recordedErr
	}
	return closeErr
}

func TestInitPublishesCanonicalMarkerAtomically(t *testing.T) {
	root := physicalPath(t, t.TempDir())
	filesystem := newRecordingFileSystem(root)

	result, err := Init(Dependencies{Files: filesystem})
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}
	if result != (InitResult{Root: root, Created: true}) {
		t.Fatalf("result = %#v", result)
	}
	wantMarker := fmt.Sprintf("format_version: 1\nroot: %s\n", root)
	markerPath := filepath.Join(root, MarkerDirectory, MarkerFile)
	if contents, err := os.ReadFile(markerPath); err != nil || string(contents) != wantMarker {
		t.Fatalf("marker = %q, %v; want %q", contents, err, wantMarker)
	}
	if got := filesystem.operations[:3]; !reflect.DeepEqual(got, []string{
		"getwd",
		"eval:" + root,
		"stat:" + root,
	}) {
		t.Fatalf("path operation prefix = %#v", got)
	}
	if !reflect.DeepEqual(filesystem.requestedOpen, []fs.FileMode{0o600}) {
		t.Fatalf("OpenFile modes = %#v", filesystem.requestedOpen)
	}
	if !reflect.DeepEqual(filesystem.requestedFlags, []int{os.O_WRONLY | os.O_CREATE | os.O_EXCL}) {
		t.Fatalf("OpenFile flags = %#v", filesystem.requestedFlags)
	}
	openOperation := firstOperationWithPrefix(filesystem.operations, "open:")
	if openOperation == "" {
		t.Fatalf("marker open was not recorded: %#v", filesystem.operations)
	}
	staging := filepath.Dir(strings.TrimPrefix(openOperation, "open:"))
	markerChmodRecorded := false
	for path, mode := range filesystem.requestedChmod {
		if filepath.Base(path) == MarkerFile && mode == 0o644 {
			markerChmodRecorded = true
		}
	}
	if !markerChmodRecorded {
		t.Fatalf("marker chmod 0644 not recorded: %#v", filesystem.requestedChmod)
	}
	if mode := filesystem.requestedChmod[staging]; mode != 0o755 {
		t.Fatalf("staging directory chmod = %o, want 0755", mode)
	}
	wantOperations := []string{"getwd", "eval:" + root, "stat:" + root}
	for ancestor := filepath.Dir(root); ; ancestor = filepath.Dir(ancestor) {
		wantOperations = append(wantOperations, "lstat:"+filepath.Join(ancestor, MarkerDirectory))
		if ancestor == filepath.Dir(ancestor) {
			break
		}
	}
	stagedMarkerPath := filepath.Join(staging, MarkerFile)
	wantOperations = append(wantOperations,
		"lstat:"+filepath.Join(root, MarkerDirectory),
		"mkdir-temp:"+root,
		"open:"+stagedMarkerPath,
		"write:"+stagedMarkerPath,
		"file-chmod:"+stagedMarkerPath,
		"sync:"+stagedMarkerPath,
		"close:"+stagedMarkerPath,
		"dir-chmod:"+staging,
		"lstat:"+filepath.Join(root, MarkerDirectory),
		"rename:"+filepath.Join(root, MarkerDirectory),
	)
	if !reflect.DeepEqual(filesystem.operations, wantOperations) {
		t.Fatalf("publish operations = %#v, want %#v", filesystem.operations, wantOperations)
	}
	if runtime.GOOS != "windows" {
		assertMode(t, filepath.Join(root, MarkerDirectory), 0o755)
		assertMode(t, markerPath, 0o644)
	}
}

func TestInitUsesPhysicalSymlinkRoot(t *testing.T) {
	parent := t.TempDir()
	physical := filepath.Join(parent, "physical")
	alias := filepath.Join(parent, "alias")
	if err := os.Mkdir(physical, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(physical, alias); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	filesystem := newRecordingFileSystem(alias)
	result, err := Init(Dependencies{Files: filesystem})
	if err != nil {
		t.Fatal(err)
	}
	physical = physicalPath(t, physical)
	if result.Root != physical {
		t.Fatalf("root = %q, want %q", result.Root, physical)
	}
	if _, err := os.Stat(filepath.Join(physical, MarkerDirectory, MarkerFile)); err != nil {
		t.Fatalf("physical marker missing: %v", err)
	}
}

func TestInitCompatibleMarkerIsIdempotentWithoutMutation(t *testing.T) {
	root := physicalPath(t, t.TempDir())
	markerPath := filepath.Join(root, MarkerDirectory, MarkerFile)
	if err := os.Mkdir(filepath.Dir(markerPath), 0o755); err != nil {
		t.Fatal(err)
	}
	contents := []byte("# retained\nroot: " + strconv.Quote(root) + "\nformat_version: 1\n")
	if err := os.WriteFile(markerPath, contents, 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, MarkerDirectory, "future-data"), []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	wantTime := time.Unix(946684800, 0)
	if err := os.Chtimes(markerPath, wantTime, wantTime); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(markerPath)
	if err != nil {
		t.Fatal(err)
	}
	filesystem := newRecordingFileSystem(root)

	result, err := Init(Dependencies{Files: filesystem})
	if err != nil {
		t.Fatal(err)
	}
	if result.Created || result.Root != root {
		t.Fatalf("result = %#v", result)
	}
	if filesystem.mutations != 0 {
		t.Fatalf("idempotent init made %d mutating calls: %#v", filesystem.mutations, filesystem.operations)
	}
	after, err := os.Stat(markerPath)
	if err != nil {
		t.Fatal(err)
	}
	afterContents, _ := os.ReadFile(markerPath)
	if !reflect.DeepEqual(afterContents, contents) || after.Mode() != before.Mode() || !after.ModTime().Equal(before.ModTime()) {
		t.Fatalf("compatible marker changed: before=%v after=%v", before, after)
	}
}

func TestInitRejectsInvalidMarkers(t *testing.T) {
	tests := []struct {
		name   string
		marker func(string) string
		reason string
	}{
		{"unknown field", func(root string) string { return "format_version: 1\nroot: " + root + "\nname: no\n" }, reasonMarkerInvalid},
		{"duplicate field", func(root string) string { return "format_version: 1\nformat_version: 1\n" }, reasonMarkerInvalid},
		{"missing field", func(root string) string { return "format_version: 1\n" }, reasonMarkerInvalid},
		{"multiple documents", func(root string) string { return "format_version: 1\nroot: " + root + "\n---\n{}\n" }, reasonMarkerInvalid},
		{"wrong version type", func(root string) string { return "format_version: \"1\"\nroot: " + root + "\n" }, reasonMarkerInvalid},
		{"wrong root type", func(root string) string { return "format_version: 1\nroot: 3\n" }, reasonMarkerInvalid},
		{"unsupported version", func(root string) string { return "format_version: 2\nroot: " + root + "\n" }, reasonVersion},
		{"noncanonical root", func(root string) string { return "format_version: 1\nroot: " + root + "/.\n" }, reasonRootMismatch},
		{"different root", func(root string) string { return "format_version: 1\nroot: " + filepath.Dir(root) + "\n" }, reasonRootMismatch},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := physicalPath(t, t.TempDir())
			writeMarkerForTest(t, root, []byte(test.marker(root)))
			_, err := Init(Dependencies{Files: newRecordingFileSystem(root)})
			assertWorkspaceError(t, err, ErrorConflict, test.reason)
		})
	}
}

func TestInitRejectsReservedPathAndMarkerTypes(t *testing.T) {
	t.Run("reserved path file", func(t *testing.T) {
		root := physicalPath(t, t.TempDir())
		if err := os.WriteFile(filepath.Join(root, MarkerDirectory), []byte("owned"), 0o600); err != nil {
			t.Fatal(err)
		}
		_, err := Init(Dependencies{Files: newRecordingFileSystem(root)})
		assertWorkspaceError(t, err, ErrorConflict, reasonReservedPath)
	})
	t.Run("reserved path symlink", func(t *testing.T) {
		root := physicalPath(t, t.TempDir())
		target := filepath.Join(root, "reserved-target")
		if err := os.Mkdir(target, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(root, MarkerDirectory)); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}
		_, err := Init(Dependencies{Files: newRecordingFileSystem(root)})
		assertWorkspaceError(t, err, ErrorConflict, reasonReservedPath)
	})
	t.Run("marker missing", func(t *testing.T) {
		root := physicalPath(t, t.TempDir())
		if err := os.Mkdir(filepath.Join(root, MarkerDirectory), 0o755); err != nil {
			t.Fatal(err)
		}
		_, err := Init(Dependencies{Files: newRecordingFileSystem(root)})
		assertWorkspaceError(t, err, ErrorConflict, reasonMarkerMissing)
	})
	t.Run("marker symlink", func(t *testing.T) {
		root := physicalPath(t, t.TempDir())
		reserved := filepath.Join(root, MarkerDirectory)
		if err := os.Mkdir(reserved, 0o755); err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(root, "marker-target")
		if err := os.WriteFile(target, []byte("not followed"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(reserved, MarkerFile)); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}
		_, err := Init(Dependencies{Files: newRecordingFileSystem(root)})
		assertWorkspaceError(t, err, ErrorConflict, reasonMarkerType)
	})
	t.Run("marker directory", func(t *testing.T) {
		root := physicalPath(t, t.TempDir())
		if err := os.MkdirAll(filepath.Join(root, MarkerDirectory, MarkerFile), 0o755); err != nil {
			t.Fatal(err)
		}
		_, err := Init(Dependencies{Files: newRecordingFileSystem(root)})
		assertWorkspaceError(t, err, ErrorConflict, reasonMarkerType)
	})
}

func TestInitChecksAncestorsBeforeLocalState(t *testing.T) {
	parent := physicalPath(t, t.TempDir())
	child := filepath.Join(parent, "child")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	writeMarkerForTest(t, parent, markerBytesForTest(parent))
	if err := os.WriteFile(filepath.Join(child, MarkerDirectory), []byte("local conflict"), 0o600); err != nil {
		t.Fatal(err)
	}
	filesystem := newRecordingFileSystem(child)
	_, err := Init(Dependencies{Files: filesystem})
	assertWorkspaceError(t, err, ErrorNested, parent)
	if operationContains(filesystem.operations, "lstat:"+filepath.Join(child, MarkerDirectory)) {
		t.Fatalf("local state inspected before nested result: %#v", filesystem.operations)
	}
}

func TestInitClassifiesAncestorFailures(t *testing.T) {
	t.Run("invalid marker", func(t *testing.T) {
		parent := physicalPath(t, t.TempDir())
		child := filepath.Join(parent, "child")
		if err := os.Mkdir(child, 0o755); err != nil {
			t.Fatal(err)
		}
		writeMarkerForTest(t, parent, []byte("invalid: true\n"))
		_, err := Init(Dependencies{Files: newRecordingFileSystem(child)})
		assertWorkspaceError(t, err, ErrorConflict, reasonMarkerInvalid)
	})
	t.Run("unreadable marker", func(t *testing.T) {
		parent := physicalPath(t, t.TempDir())
		child := filepath.Join(parent, "child")
		if err := os.Mkdir(child, 0o755); err != nil {
			t.Fatal(err)
		}
		writeMarkerForTest(t, parent, markerBytesForTest(parent))
		markerPath := filepath.Join(parent, MarkerDirectory, MarkerFile)
		filesystem := newRecordingFileSystem(child)
		filesystem.fail["read:"+markerPath] = fs.ErrPermission
		_, err := Init(Dependencies{Files: filesystem})
		assertWorkspaceError(t, err, ErrorConflict, reasonMarkerInvalid)
	})
	t.Run("observation error", func(t *testing.T) {
		parent := physicalPath(t, t.TempDir())
		child := filepath.Join(parent, "child")
		if err := os.Mkdir(child, 0o755); err != nil {
			t.Fatal(err)
		}
		filesystem := newRecordingFileSystem(child)
		filesystem.fail["lstat:"+filepath.Join(parent, MarkerDirectory)] = fs.ErrPermission
		_, err := Init(Dependencies{Files: filesystem})
		assertWorkspaceError(t, err, ErrorPath, "lstat")
	})
}

func TestInitClassifiesPathFailures(t *testing.T) {
	for _, operation := range []string{"getwd", "eval", "stat"} {
		t.Run(operation, func(t *testing.T) {
			root := physicalPath(t, t.TempDir())
			filesystem := newRecordingFileSystem(root)
			filesystem.fail[operation] = errInjected
			_, err := Init(Dependencies{Files: filesystem})
			assertWorkspaceError(t, err, ErrorPath, errInjected.Error())
		})
	}
}

func TestInitCleansOwnStagingOnPublicationFailures(t *testing.T) {
	tests := []struct {
		name       string
		operation  string
		shortWrite bool
	}{
		{"open", "open", false},
		{"write", "write", false},
		{"short write", "", true},
		{"file chmod", "file-chmod", false},
		{"sync", "sync", false},
		{"close", "close", false},
		{"directory chmod", "dir-chmod", false},
		{"rename", "rename", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := physicalPath(t, t.TempDir())
			filesystem := newRecordingFileSystem(root)
			if test.operation != "" {
				filesystem.fail[test.operation] = errInjected
			}
			filesystem.shortWrite = test.shortWrite
			_, err := Init(Dependencies{Files: filesystem})
			assertWorkspaceError(t, err, ErrorIO, "")
			if _, statErr := os.Lstat(filepath.Join(root, MarkerDirectory)); !errors.Is(statErr, fs.ErrNotExist) {
				t.Fatalf("authoritative destination exists after failure: %v", statErr)
			}
			staging, globErr := filepath.Glob(filepath.Join(root, ".ply.init-*"))
			if globErr != nil || len(staging) != 0 {
				t.Fatalf("staging left behind: %v, %v", staging, globErr)
			}
		})
	}
}

func TestInitPreservesPrimaryFailureWhenCleanupFails(t *testing.T) {
	root := physicalPath(t, t.TempDir())
	filesystem := newRecordingFileSystem(root)
	filesystem.fail["write"] = errInjected
	filesystem.fail["remove"] = errors.New("cleanup failed")
	_, err := Init(Dependencies{Files: filesystem})
	assertWorkspaceError(t, err, ErrorIO, "cleanup")
	if !errors.Is(err, errInjected) {
		t.Fatalf("primary error not preserved: %v", err)
	}
	if len(filesystem.removeAllCalled) != 1 || !strings.Contains(err.Error(), filesystem.removeAllCalled[0]) {
		t.Fatalf("staging provenance missing: %v; paths=%v", err, filesystem.removeAllCalled)
	}
}

func TestInitHandlesConcurrentPublish(t *testing.T) {
	t.Run("identical marker", func(t *testing.T) {
		root := physicalPath(t, t.TempDir())
		filesystem := newRecordingFileSystem(root)
		filesystem.renameHook = func(oldPath, newPath string) error {
			if err := os.Mkdir(newPath, 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(newPath, MarkerFile), markerBytesForTest(root), 0o644); err != nil {
				return err
			}
			return errInjected
		}
		result, err := Init(Dependencies{Files: filesystem})
		if err != nil {
			t.Fatal(err)
		}
		if result.Created {
			t.Fatalf("concurrent loser reported creation: %#v", result)
		}
		if len(filesystem.removeAllCalled) != 1 {
			t.Fatalf("losing staging was not cleaned: %v", filesystem.removeAllCalled)
		}
	})
	t.Run("conflicting marker", func(t *testing.T) {
		root := physicalPath(t, t.TempDir())
		filesystem := newRecordingFileSystem(root)
		filesystem.renameHook = func(oldPath, newPath string) error {
			if err := os.Mkdir(newPath, 0o755); err != nil {
				return err
			}
			return errInjected
		}
		_, err := Init(Dependencies{Files: filesystem})
		assertWorkspaceError(t, err, ErrorConflict, reasonMarkerMissing)
		if len(filesystem.removeAllCalled) != 1 {
			t.Fatalf("losing staging was not cleaned: %v", filesystem.removeAllCalled)
		}
	})
}

func TestInitZeroDependenciesHasNoOSEffect(t *testing.T) {
	_, err := Init(Dependencies{})
	assertWorkspaceError(t, err, ErrorIO, "filesystem dependency is required")
}

func TestInvalidArgumentsIsTyped(t *testing.T) {
	err := InvalidArguments("bad input")
	assertWorkspaceError(t, err, ErrorInvalidArguments, "bad input")
}

func markerBytesForTest(root string) []byte {
	return []byte(fmt.Sprintf("format_version: 1\nroot: %s\n", root))
}

func writeMarkerForTest(t *testing.T, root string, contents []byte) {
	t.Helper()
	reserved := filepath.Join(root, MarkerDirectory)
	if err := os.MkdirAll(reserved, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(reserved, MarkerFile), contents, 0o644); err != nil {
		t.Fatal(err)
	}
}

func physicalPath(t *testing.T, path string) string {
	t.Helper()
	physical, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Clean(physical)
}

func assertMode(t *testing.T, path string, want fs.FileMode) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != want {
		t.Fatalf("%s mode = %o, want %o", path, got, want)
	}
}

func assertWorkspaceError(t *testing.T, err error, class ErrorClass, detail string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected %s error", class)
	}
	var workspaceErr *Error
	if !errors.As(err, &workspaceErr) {
		t.Fatalf("error type = %T, want *Error: %v", err, err)
	}
	if workspaceErr.Class != class {
		t.Fatalf("error class = %s, want %s: %v", workspaceErr.Class, class, err)
	}
	if detail != "" && !strings.Contains(err.Error(), detail) {
		t.Fatalf("error %q does not contain %q", err, detail)
	}
	if !strings.HasPrefix(err.Error(), string(class)+":") {
		t.Fatalf("error does not start with class: %v", err)
	}
}

func operationContains(operations []string, prefix string) bool {
	for _, operation := range operations {
		if strings.HasPrefix(operation, prefix) {
			return true
		}
	}
	return false
}

func firstOperationWithPrefix(operations []string, prefix string) string {
	for _, operation := range operations {
		if strings.HasPrefix(operation, prefix) {
			return operation
		}
	}
	return ""
}
