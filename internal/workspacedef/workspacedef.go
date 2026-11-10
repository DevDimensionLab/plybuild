// Package workspacedef reads and updates the workspace definitions file
// ply-workspace.json under the workspace's projects.lock. Identifier, path
// and Git rules belong to the caller's check.
package workspacedef

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/atomicfile"
	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/filelock"
)

// File lives beside .ply, independently of disposable runtime data.
const File = "ply-workspace.json"

// Absent is the basis digest of a missing definitions file.
const Absent = "absent"

type Definitions struct {
	FormatVersion int       `json:"format_version"`
	Projects      []Project `json:"projects"`
}
type Project struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Wrapper   string            `json:"wrapper"`
	DataPaths map[string]string `json:"data_paths,omitempty"`
	Repos     []Repo            `json:"repos"`
	Epics     []Epic            `json:"epics"`
}
type Repo struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}
type Epic struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Repo  string `json:"repo"`
	Path  string `json:"path"`
	Ref   string `json:"ref"`
}

// Effect is what an Update did to the definitions file.
type Effect string

const (
	NotAttempted Effect = "not_attempted"
	Completed    Effect = "completed"
	NoEffect     Effect = "no_effect"
	Unknown      Effect = "unknown"
)

type Result struct {
	Effect  Effect
	Changed bool
}

type Kind string

const (
	KindUsage      Kind = "usage"
	KindNotRegular Kind = "not_regular"
	KindInvalid    Kind = "invalid"
	KindSchema     Kind = "schema"
	KindFormat     Kind = "format"
	KindMissing    Kind = "missing"
	KindNotFound   Kind = "not_found"
	KindConflict   Kind = "conflict"
	KindPublish    Kind = "publish"
	KindReadback   Kind = "readback"
)

type Error struct {
	Kind   Kind
	Path   string
	Detail string
	Err    error
}

func (e *Error) Error() string {
	message := string(e.Kind)
	if e.Path != "" {
		message += " " + e.Path
	}
	if e.Detail != "" {
		message += ": " + e.Detail
	}
	if e.Err != nil {
		message += ": " + e.Err.Error()
	}
	return message
}

func (e *Error) Unwrap() error { return e.Err }

// Read returns the definitions and the digest of their raw bytes, or
// (nil, Absent, nil) when the file does not exist. OS errors are returned as is.
func Read(root string) (*Definitions, string, error) {
	path := filepath.Join(root, File)
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, Absent, nil
	}
	if err != nil {
		return nil, "", err
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return nil, "", &Error{Kind: KindNotRegular, Path: path, Detail: "definitions must be a regular file"}
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	defs, err := decode(path, raw)
	if err != nil {
		return nil, "", err
	}
	return defs, digest(raw), nil
}

func decode(path string, raw []byte) (*Definitions, error) {
	if _, err := canonicaljson.DecodeStrict(raw); err != nil {
		return nil, &Error{Kind: KindInvalid, Path: path, Detail: "invalid definitions", Err: err}
	}
	var out Definitions
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&out); err != nil {
		return nil, &Error{Kind: KindSchema, Path: path, Detail: err.Error()}
	}
	if out.FormatVersion != 1 || out.Projects == nil {
		return nil, &Error{Kind: KindFormat, Path: path, Detail: "expected format_version 1 and projects array"}
	}
	return &out, nil
}

func digest(raw []byte) string {
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// AddRepo appends {id, path} to project's repos. An identical member is no
// change; the same id with another path is a conflict.
func AddRepo(d *Definitions, project, id, path string) (bool, error) {
	for i := range d.Projects {
		p := &d.Projects[i]
		if p.ID != project {
			continue
		}
		for _, r := range p.Repos {
			if r.ID == id && r.Path == path {
				return false, nil
			}
			if r.ID == id {
				return false, &Error{Kind: KindConflict, Detail: fmt.Sprintf("repository %s of Project %s is already defined with path %s", id, project, r.Path)}
			}
		}
		p.Repos = append(p.Repos, Repo{ID: id, Path: path})
		return true, nil
	}
	return false, &Error{Kind: KindNotFound, Detail: fmt.Sprintf("Project %s is not defined", project)}
}

// RemoveRepo removes the member {id, path} from project's repos. An absent id
// is no change; the same id with another path is a conflict.
func RemoveRepo(d *Definitions, project, id, path string) (bool, error) {
	for i := range d.Projects {
		p := &d.Projects[i]
		if p.ID != project {
			continue
		}
		for j, r := range p.Repos {
			if r.ID == id && r.Path != path {
				return false, &Error{Kind: KindConflict, Detail: fmt.Sprintf("repository %s of Project %s is defined with path %s, not %s", id, project, r.Path, path)}
			}
			if r.ID == id {
				p.Repos = append(append([]Repo{}, p.Repos[:j]...), p.Repos[j+1:]...)
				return true, nil
			}
		}
		return false, nil
	}
	return false, &Error{Kind: KindNotFound, Detail: fmt.Sprintf("Project %s is not defined", project)}
}

// hooks is a private, per-call test seam; production has no mutable hooks.
type hooks struct {
	beforeFresh func()
	replace     func(string, []byte, fs.FileMode) error
}

func systemHooks() hooks {
	return hooks{beforeFresh: func() {}, replace: atomicfile.Replace}
}

// Update applies mutate to a copy of the definitions whose digest is basis
// while holding <root>/.ply/projects.lock. check runs on the current
// definitions and on a changed candidate before anything is written; it must
// not take projects.lock. Lock and release failures are filelock errors.
func Update(root, basis string, mutate func(*Definitions) (bool, error), check func(Definitions) error) (Result, error) {
	return update(root, basis, mutate, check, systemHooks())
}

func update(root, basis string, mutate func(*Definitions) (bool, error), check func(Definitions) error, h hooks) (Result, error) {
	if mutate == nil || check == nil {
		return Result{Effect: NotAttempted}, &Error{Kind: KindUsage, Detail: "mutate and check are required"}
	}
	result := Result{Effect: NotAttempted}
	err := WithLock(root, func() error {
		var err error
		result, err = updateLocked(root, basis, mutate, check, h)
		return err
	})
	return result, err
}

// WithLock runs fn while holding <root>/.ply/projects.lock, the lock Update
// takes. fn must not take projects.lock. Lock and release failures are
// filelock errors; a release failure after fn failed is errors.Join(fnErr, releaseErr).
// fn's own error is returned unchanged.
func WithLock(root string, fn func() error) error {
	if fn == nil {
		return &Error{Kind: KindUsage, Detail: "fn is required"}
	}
	lock, err := filelock.Acquire(filepath.Join(root, ".ply", "projects.lock"), filelock.Options{})
	if err != nil {
		return err
	}
	err = fn()
	if releaseErr := lock.Release(); releaseErr != nil {
		if err != nil {
			return errors.Join(err, releaseErr)
		}
		return releaseErr
	}
	return err
}

func updateLocked(root, basis string, mutate func(*Definitions) (bool, error), check func(Definitions) error, h hooks) (Result, error) {
	path := filepath.Join(root, File)
	notAttempted := Result{Effect: NotAttempted}
	current, currentDigest, err := Read(root)
	if err != nil {
		return notAttempted, err
	}
	if current == nil {
		return notAttempted, &Error{Kind: KindMissing, Path: path, Detail: "definitions do not exist"}
	}
	if currentDigest != basis {
		return notAttempted, &Error{Kind: KindConflict, Path: path, Detail: "definitions changed since they were read"}
	}
	if err := check(clone(*current)); err != nil {
		return notAttempted, err
	}
	candidate := clone(*current)
	changed, err := mutate(&candidate)
	if err != nil {
		return notAttempted, err
	}
	if !changed {
		return Result{Effect: Completed}, nil
	}
	if err := check(clone(candidate)); err != nil {
		return notAttempted, err
	}
	raw, err := json.MarshalIndent(candidate, "", "  ")
	if err != nil {
		return notAttempted, err
	}
	raw = append(raw, '\n')

	// A fresh raw read narrows the window against editors that ignore the lock;
	// it is not an atomic compare-and-swap.
	h.beforeFresh()
	info, err := os.Lstat(path)
	if err != nil {
		return notAttempted, err
	}
	fresh, err := os.ReadFile(path)
	if err != nil {
		return notAttempted, err
	}
	if digest(fresh) != basis {
		return notAttempted, &Error{Kind: KindConflict, Path: path, Detail: "definitions changed before they were replaced"}
	}
	replaced := h.replace(path, raw, info.Mode().Perm())
	_, after, readErr := Read(root)
	switch {
	case readErr == nil && after == digest(raw):
		return Result{Effect: Completed, Changed: true}, nil
	case readErr == nil && after == basis:
		return Result{Effect: NoEffect}, &Error{Kind: KindPublish, Path: path, Detail: "definitions were not replaced", Err: replaced}
	default:
		cause := replaced
		if cause == nil {
			cause = readErr
		} else if readErr != nil {
			cause = errors.Join(replaced, readErr)
		}
		return Result{Effect: Unknown}, &Error{Kind: KindReadback, Path: path, Detail: "definitions hold unexpected content after replacement", Err: cause}
	}
}

func clone(d Definitions) Definitions {
	out := Definitions{FormatVersion: d.FormatVersion, Projects: make([]Project, len(d.Projects))}
	for i, p := range d.Projects {
		// Keep nil and empty slices apart: they marshal as null and [].
		if p.Repos != nil {
			p.Repos = append([]Repo{}, p.Repos...)
		}
		if p.Epics != nil {
			p.Epics = append([]Epic{}, p.Epics...)
		}
		if p.DataPaths != nil {
			paths := make(map[string]string, len(p.DataPaths))
			for k, v := range p.DataPaths {
				paths[k] = v
			}
			p.DataPaths = paths
		}
		out.Projects[i] = p
	}
	return out
}
