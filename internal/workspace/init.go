package workspace

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspacedef"
	"github.com/devdimensionlab/plybuild/pkg/config"
	"gopkg.in/yaml.v3"
)

const (
	MarkerDirectory = ".ply"
	MarkerFile      = "workspace.yaml"
	FormatVersion   = 1
)

// MarkerPath is the path of the workspace marker of the workspace at root.
func MarkerPath(root string) string { return filepath.Join(root, MarkerDirectory, MarkerFile) }

const (
	reasonReservedPath  = "reserved path is not a regular directory"
	reasonMarkerMissing = "workspace marker is missing"
	reasonMarkerType    = "workspace marker is not a regular file"
	reasonMarkerInvalid = "workspace marker is invalid"
	reasonVersion       = "unsupported format_version"
	reasonRootMismatch  = "workspace root does not match physical directory"
)

type Marker struct {
	FormatVersion int    `yaml:"format_version"`
	Root          string `yaml:"root"`
}

type InitResult struct {
	Root    string
	Created bool
}

type ErrorClass string

const (
	ErrorInvalidArguments ErrorClass = "workspace_init_invalid_arguments"
	ErrorPath             ErrorClass = "workspace_init_path_error"
	ErrorConflict         ErrorClass = "workspace_init_conflict"
	ErrorIO               ErrorClass = "workspace_init_io_error"
)

type Error struct {
	Class     ErrorClass
	Operation string
	Path      string
	Reason    string
	Err       error
}

func (err *Error) Error() string {
	detail := err.Reason
	if detail == "" && err.Err != nil {
		detail = err.Err.Error()
	}
	switch err.Class {
	case ErrorConflict:
		return fmt.Sprintf("%s: %s: %s", err.Class, err.Path, detail)
	case ErrorPath:
		return fmt.Sprintf("%s: %s: %s", err.Class, err.Operation, detail)
	case ErrorIO:
		return fmt.Sprintf("%s: %s: %s: %s", err.Class, err.Operation, err.Path, detail)
	default:
		return fmt.Sprintf("%s: %s", err.Class, detail)
	}
}

func (err *Error) Unwrap() error { return err.Err }

// Dependencies are the effects of the register. RepoUses comes from the
// work-item store; SystemDependencies leaves it unset.
type Dependencies struct {
	ProjectDataRoot func(string) (string, error)
	Files           FileSystem
	Repos           RepoObserver
	Projects        ProjectStore
	ProjectLocks    ProjectSnapshotLocker
	RepoUses        RepoUses
}

func InvalidArguments(detail string) error {
	return &Error{Class: ErrorInvalidArguments, Reason: detail}
}

func SystemDependencies() Dependencies {
	files := systemFileSystem{}
	projects := newSystemProjectStore()
	return Dependencies{
		ProjectDataRoot: config.ProjectDataRoot,
		Files:           files,
		Repos:           newSystemRepoObserver(files),
		Projects:        projects,
		ProjectLocks:    projects,
	}
}

type localState int

const (
	localAbsent localState = iota
	localCompatible
)

// Init makes the working directory a workspace: the .ply marker and, when
// missing, ply-workspace.json without Projects.
func Init(dependencies Dependencies) (InitResult, error) {
	result, err := initMarker(dependencies)
	if err != nil {
		return result, err
	}
	if err := workspacedef.Create(result.Root); err != nil {
		return result, ioError("create", filepath.Join(result.Root, DefinitionsFile), err)
	}
	return result, nil
}

func initMarker(dependencies Dependencies) (InitResult, error) {
	files := dependencies.Files
	if files == nil {
		return InitResult{}, ioError("dependencies", "filesystem", errors.New("filesystem dependency is required"))
	}

	root, err := physicalWorkingDirectory(files)
	if err != nil {
		return InitResult{}, err
	}
	result := InitResult{Root: root}

	state, err := inspectLocal(files, root)
	if err != nil {
		return result, err
	}
	if state == localCompatible {
		return result, nil
	}

	markerBytes, err := yaml.Marshal(Marker{FormatVersion: FormatVersion, Root: root})
	if err != nil {
		return result, ioError("marshal", MarkerPath(root), err)
	}
	staging, err := files.MkdirTemp(root, ".ply.init-*")
	if err != nil {
		return result, ioError("mkdir-temp", root, err)
	}
	published := false
	markerPath := filepath.Join(staging, MarkerFile)

	markerFile, err := files.OpenFile(markerPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return result, abortStaging(files, staging, "open", markerPath, err)
	}
	if err := writeMarker(markerFile, markerBytes); err != nil {
		return result, abortOpenStaging(files, markerFile, staging, "write", markerPath, err)
	}
	if err := markerFile.Chmod(0o644); err != nil {
		return result, abortOpenStaging(files, markerFile, staging, "chmod", markerPath, err)
	}
	if err := markerFile.Sync(); err != nil {
		return result, abortOpenStaging(files, markerFile, staging, "sync", markerPath, err)
	}
	if err := markerFile.Close(); err != nil {
		return result, abortStaging(files, staging, "close", markerPath, err)
	}
	if err := files.Chmod(staging, 0o755); err != nil {
		return result, abortStaging(files, staging, "chmod", staging, err)
	}

	destination := filepath.Join(root, MarkerDirectory)
	state, err = inspectLocal(files, root)
	if err != nil {
		return result, abortTypedStaging(files, staging, err)
	}
	if state == localCompatible {
		if err := cleanupStaging(files, staging); err != nil {
			return result, err
		}
		return result, nil
	}

	if err := files.Rename(staging, destination); err != nil {
		renameErr := err
		state, observedErr := inspectLocal(files, root)
		if observedErr == nil && state == localCompatible {
			if err := cleanupStaging(files, staging); err != nil {
				return result, err
			}
			return result, nil
		}
		if observedErr != nil {
			if typed, ok := observedErr.(*Error); ok && typed.Class == ErrorConflict {
				if cleanupErr := files.RemoveAll(staging); cleanupErr != nil {
					return result, ioError("cleanup", staging, cleanupErr)
				}
				return result, observedErr
			}
			renameErr = fmt.Errorf("%w; destination observation: %v", renameErr, observedErr)
		}
		return result, abortStaging(files, staging, "rename", destination, renameErr)
	}
	published = true
	_ = published // documents that cleanup must never target the published destination
	result.Created = true
	return result, nil
}

func physicalWorkingDirectory(files FileSystem) (string, error) {
	workingDirectory, err := files.Getwd()
	if err != nil {
		return "", pathError("getwd", err)
	}
	absolute, err := filepath.Abs(workingDirectory)
	if err != nil {
		return "", pathError("abs", err)
	}
	absolute = filepath.Clean(absolute)
	physical, err := files.EvalSymlinks(absolute)
	if err != nil {
		return "", pathError("eval-symlinks", err)
	}
	physical, err = filepath.Abs(physical)
	if err != nil {
		return "", pathError("abs", err)
	}
	physical = filepath.Clean(physical)
	info, err := files.Stat(physical)
	if err != nil {
		return "", pathError("stat", err)
	}
	if !info.IsDir() {
		return "", pathError("stat", errors.New("working directory is not a directory"))
	}
	return physical, nil
}

func inspectLocal(files FileSystem, root string) (localState, error) {
	reserved := filepath.Join(root, MarkerDirectory)
	info, err := files.Lstat(reserved)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return localAbsent, nil
		}
		return localAbsent, pathError("lstat", err)
	}
	if !info.IsDir() || info.Mode()&fs.ModeSymlink != 0 {
		return localAbsent, conflict(reserved, reasonReservedPath)
	}

	markerPath := filepath.Join(reserved, MarkerFile)
	markerInfo, err := files.Lstat(markerPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return localAbsent, conflict(markerPath, reasonMarkerMissing)
		}
		return localAbsent, pathError("lstat", err)
	}
	if !markerInfo.Mode().IsRegular() || markerInfo.Mode()&fs.ModeSymlink != 0 {
		return localAbsent, conflict(markerPath, reasonMarkerType)
	}
	if reason := readAndValidateMarker(files, markerPath, root); reason != "" {
		return localAbsent, conflict(markerPath, reason)
	}
	return localCompatible, nil
}

func readAndValidateMarker(files FileSystem, markerPath, expectedRoot string) string {
	contents, err := files.ReadFile(markerPath)
	if err != nil {
		return reasonMarkerInvalid
	}
	marker, reason := decodeMarker(contents)
	if reason != "" {
		return reason
	}
	if marker.FormatVersion != FormatVersion {
		return reasonVersion
	}
	if !markerRootMatches(files, marker.Root, expectedRoot) {
		return reasonRootMismatch
	}
	return ""
}

func decodeMarker(contents []byte) (Marker, string) {
	decoder := yaml.NewDecoder(strings.NewReader(string(contents)))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		return Marker{}, reasonMarkerInvalid
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		return Marker{}, reasonMarkerInvalid
	}
	if len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
		return Marker{}, reasonMarkerInvalid
	}
	mapping := document.Content[0]
	if len(mapping.Content) != 4 {
		return Marker{}, reasonMarkerInvalid
	}
	seen := make(map[string]bool, 2)
	for index := 0; index < len(mapping.Content); index += 2 {
		key := mapping.Content[index]
		value := mapping.Content[index+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" || seen[key.Value] {
			return Marker{}, reasonMarkerInvalid
		}
		seen[key.Value] = true
		switch key.Value {
		case "format_version":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!int" {
				return Marker{}, reasonMarkerInvalid
			}
		case "root":
			if value.Kind != yaml.ScalarNode || value.Tag != "!!str" {
				return Marker{}, reasonMarkerInvalid
			}
		default:
			return Marker{}, reasonMarkerInvalid
		}
	}
	if !seen["format_version"] || !seen["root"] {
		return Marker{}, reasonMarkerInvalid
	}
	var marker Marker
	if err := mapping.Decode(&marker); err != nil {
		return Marker{}, reasonMarkerInvalid
	}
	return marker, ""
}

func markerRootMatches(files FileSystem, markerRoot, expectedRoot string) bool {
	if markerRoot == "" || !filepath.IsAbs(markerRoot) || filepath.Clean(markerRoot) != markerRoot {
		return false
	}
	physical, err := files.EvalSymlinks(markerRoot)
	if err != nil {
		return false
	}
	physical, err = filepath.Abs(physical)
	if err != nil {
		return false
	}
	physical = filepath.Clean(physical)
	if physical != markerRoot || physical != expectedRoot {
		return false
	}
	info, err := files.Stat(physical)
	return err == nil && info.IsDir()
}

func writeMarker(file File, contents []byte) error {
	written, err := file.Write(contents)
	if err != nil {
		return err
	}
	if written != len(contents) {
		return io.ErrShortWrite
	}
	return nil
}

func abortOpenStaging(files FileSystem, file File, staging, operation, path string, primary error) error {
	if closeErr := file.Close(); closeErr != nil {
		primary = fmt.Errorf("%w; close: %v", primary, closeErr)
	}
	return abortStaging(files, staging, operation, path, primary)
}

func abortTypedStaging(files FileSystem, staging string, primary error) error {
	if cleanupErr := files.RemoveAll(staging); cleanupErr != nil {
		if typed, ok := primary.(*Error); ok && typed.Class == ErrorIO {
			typed.Reason = fmt.Sprintf("%v; cleanup %s: %v", typed.Err, staging, cleanupErr)
			return typed
		}
		return ioError("cleanup", staging, cleanupErr)
	}
	return primary
}

func abortStaging(files FileSystem, staging, operation, path string, primary error) error {
	reason := primary.Error()
	if cleanupErr := files.RemoveAll(staging); cleanupErr != nil {
		reason = fmt.Sprintf("%s; cleanup %s: %v", reason, staging, cleanupErr)
	}
	return &Error{Class: ErrorIO, Operation: operation, Path: path, Reason: reason, Err: primary}
}

func cleanupStaging(files FileSystem, staging string) error {
	if err := files.RemoveAll(staging); err != nil {
		return ioError("cleanup", staging, err)
	}
	return nil
}

func pathError(operation string, err error) error {
	return &Error{Class: ErrorPath, Operation: operation, Err: err}
}

func ioError(operation, path string, err error) error {
	return &Error{Class: ErrorIO, Operation: operation, Path: path, Err: err}
}

func conflict(path, reason string) error {
	return &Error{Class: ErrorConflict, Path: path, Reason: reason}
}
