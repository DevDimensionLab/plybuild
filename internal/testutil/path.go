package testutil

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var ErrRepositoryWrite = errors.New("refusing to write inside repository")

// WouldLeakIntoWorkingTree reports whether path resolves inside the caller's repository.
func WouldLeakIntoWorkingTree(path string) bool {
	resolvedPath, err := resolveCandidatePath(path)
	if err != nil {
		return true
	}
	cwd, err := os.Getwd()
	if err != nil {
		return true
	}
	repositoryRoot, err := findRepositoryRoot(cwd)
	if err != nil {
		return true
	}
	resolvedRoot, err := filepath.EvalSymlinks(repositoryRoot)
	if err != nil {
		return true
	}
	relative, err := filepath.Rel(resolvedRoot, resolvedPath)
	if err != nil {
		return true
	}
	return relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator)))
}

// WriteFileOutsideWorkingTree writes only when path resolves outside the repository.
func WriteFileOutsideWorkingTree(path string, contents []byte, permission fs.FileMode) error {
	if err := ensureOutsideWorkingTree(path); err != nil {
		return err
	}
	return os.WriteFile(path, contents, permission)
}

// CopyFSOutsideWorkingTree copies a fixture tree without requiring os.CopyFS (Go 1.23+).
func CopyFSOutsideWorkingTree(destination string, source fs.FS) error {
	if err := ensureOutsideWorkingTree(destination); err != nil {
		return err
	}

	return fs.WalkDir(source, ".", func(sourcePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		target := destination
		if sourcePath != "." {
			target = filepath.Join(destination, filepath.FromSlash(sourcePath))
		}
		if err := ensureOutsideWorkingTree(target); err != nil {
			return err
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		if entry.Type()&fs.ModeSymlink != 0 || !info.Mode().IsRegular() {
			return fmt.Errorf("unsupported fixture entry %s (%s)", sourcePath, info.Mode())
		}
		contents, err := fs.ReadFile(source, sourcePath)
		if err != nil {
			return err
		}
		return WriteFileOutsideWorkingTree(target, contents, info.Mode().Perm())
	})
}

func ensureOutsideWorkingTree(path string) error {
	if WouldLeakIntoWorkingTree(path) {
		return fmt.Errorf("%w: %s", ErrRepositoryWrite, path)
	}
	return nil
}

func resolveCandidatePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("path is empty")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if info, lstatErr := os.Lstat(absPath); lstatErr == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("final path is a symlink: %s", absPath)
		}
	} else if !errors.Is(lstatErr, fs.ErrNotExist) {
		return "", lstatErr
	}

	current := absPath
	var missing []string
	for {
		resolved, evalErr := filepath.EvalSymlinks(current)
		if evalErr == nil {
			for i := len(missing) - 1; i >= 0; i-- {
				resolved = filepath.Join(resolved, missing[i])
			}
			return filepath.Clean(resolved), nil
		}
		if !errors.Is(evalErr, fs.ErrNotExist) {
			return "", evalErr
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("cannot resolve existing parent for %s", absPath)
		}
		missing = append(missing, filepath.Base(current))
		current = parent
	}
}

func findRepositoryRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", fmt.Errorf("make repository search path absolute: %w", err)
	}
	moduleRoot := ""
	for {
		if _, statErr := os.Stat(filepath.Join(dir, ".git")); statErr == nil {
			return dir, nil
		} else if !os.IsNotExist(statErr) {
			return "", fmt.Errorf("inspect repository marker .git: %w", statErr)
		}
		if moduleRoot == "" {
			if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
				moduleRoot = dir
			} else if !os.IsNotExist(statErr) {
				return "", fmt.Errorf("inspect repository marker go.mod: %w", statErr)
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			if moduleRoot != "" {
				return moduleRoot, nil
			}
			return "", fmt.Errorf("could not find repository root above %s", start)
		}
		dir = parent
	}
}
