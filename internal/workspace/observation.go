package workspace

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
)

type WorkspaceObservation struct {
	Root                string
	MarkerFormatVersion int
	MarkerSHA256        string
}

func ObserveContaining(dependencies Dependencies) (WorkspaceObservation, error) {
	if dependencies.Files == nil {
		return WorkspaceObservation{}, projectError(ErrorProjectIO, "dependencies: filesystem is required", nil)
	}
	cwd, err := projectPhysicalWorkingDirectory(dependencies.Files)
	if err != nil {
		return WorkspaceObservation{}, err
	}
	root, err := locateContainingWorkspace(dependencies.Files, cwd)
	if err != nil {
		return WorkspaceObservation{}, err
	}
	return observeWorkspaceRoot(dependencies.Files, root)
}

func ObserveRoot(dependencies Dependencies, root string) (WorkspaceObservation, error) {
	if dependencies.Files == nil {
		return WorkspaceObservation{}, projectError(ErrorProjectIO, "dependencies: filesystem is required", nil)
	}
	physical, err := canonicalDirectory(dependencies.Files, "", root, "workspace root")
	if err != nil {
		return WorkspaceObservation{}, err
	}
	return observeWorkspaceRoot(dependencies.Files, physical)
}

func observeWorkspaceRoot(files FileSystem, root string) (WorkspaceObservation, error) {
	markerPath := filepath.Join(root, MarkerDirectory, MarkerFile)
	info, err := files.Lstat(markerPath)
	if errors.Is(err, fs.ErrNotExist) {
		return WorkspaceObservation{}, projectError(ErrorWorkspaceNotFound, fmt.Sprintf("no Ply workspace exists at %s", root), nil)
	}
	if err != nil {
		return WorkspaceObservation{}, projectError(ErrorProjectWorkspace, fmt.Sprintf("inspect %s: %v", markerPath, err), err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return WorkspaceObservation{}, projectError(ErrorProjectWorkspace, fmt.Sprintf("%s: %s", markerPath, reasonMarkerType), nil)
	}
	contents, err := files.ReadFile(markerPath)
	if err != nil {
		return WorkspaceObservation{}, projectError(ErrorProjectWorkspace, fmt.Sprintf("read %s: %v", markerPath, err), err)
	}
	marker, reason := decodeMarker(contents)
	if reason == "" && marker.FormatVersion != FormatVersion {
		reason = reasonVersion
	}
	if reason == "" && !markerRootMatches(files, marker.Root, root) {
		reason = reasonRootMismatch
	}
	if reason != "" {
		return WorkspaceObservation{}, projectError(ErrorProjectWorkspace, fmt.Sprintf("%s: %s", markerPath, reason), nil)
	}
	digest := sha256.Sum256(contents)
	return WorkspaceObservation{
		Root:                root,
		MarkerFormatVersion: marker.FormatVersion,
		MarkerSHA256:        fmt.Sprintf("sha256:%x", digest),
	}, nil
}
