// Package projectdata owns the layout of a Project's data: every category
// is the directory <wrapper>/data/<category>. No other package builds paths
// under data/. Resolving a path creates nothing; the existing part of the
// path must stay inside the wrapper after symlinks are followed. Errors are
// *Error, and their details carry the Project ID but never a path.
package projectdata

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ClassPath is the class of every data path error.
const ClassPath = "workspace_project_path_error"

// Error is a data path that cannot be resolved inside the wrapper.
type Error struct{ Class, Detail string }

func (e *Error) Error() string { return e.Class + ": " + e.Detail }

// Area is the data of one Project.
type Area struct{ project, wrapper string }

// For returns the data area of project, whose wrapper is an absolute path.
func For(project, wrapper string) Area { return Area{project: project, wrapper: wrapper} }

// AgentStatus is the directory of the Project's agent statuses.
func (a Area) AgentStatus() (string, error) { return a.resolve("agent-status") }

// Assignments is the directory of the Project's Assignments.
func (a Area) Assignments() (string, error) { return a.resolve("assignments") }

// Backlog is the directory of the Project's backlog.
func (a Area) Backlog() (string, error) { return a.resolve("backlog") }

// BacklogArchive is the directory of the Project's backlog archive.
func (a Area) BacklogArchive() (string, error) { return a.resolve("backlog-archive") }

// Inbox is the directory of the Project's agent reply Inbox.
func (a Area) Inbox() (string, error) { return a.resolve("inbox") }

// Intelligence is the directory of the Project's intelligence.
func (a Area) Intelligence() (string, error) { return a.resolve("intelligence") }

// Plans is the directory of the Project's plans.
func (a Area) Plans() (string, error) { return a.resolve("plans") }

// Publications is the directory of the Project's publications.
func (a Area) Publications() (string, error) { return a.resolve("publications") }

// Work is the directory of the Project's work: its work-item registry, the
// registry lock and the task-content store the registry references.
func (a Area) Work() (string, error) { return a.resolve("work") }

func (a Area) resolve(category string) (string, error) {
	path := filepath.Join(a.wrapper, "data", category)
	wrapper, err := filepath.EvalSymlinks(a.wrapper)
	if err != nil {
		return "", a.fail("cannot resolve wrapper of project %s")
	}
	existing := path
	for {
		if _, err := os.Lstat(existing); err == nil {
			break
		} else if !errors.Is(err, fs.ErrNotExist) {
			return "", a.fail("cannot inspect data path " + category + " of project %s")
		}
		existing = filepath.Dir(existing)
	}
	resolved, err := filepath.EvalSymlinks(existing)
	if err != nil {
		return "", a.fail("cannot resolve data path " + category + " of project %s")
	}
	if inside, err := filepath.Rel(wrapper, resolved); err != nil || inside == ".." || strings.HasPrefix(inside, ".."+string(filepath.Separator)) {
		return "", a.fail("data path " + category + " of project %s leaves the project wrapper")
	}
	return path, nil
}

func (a Area) fail(format string) error {
	return &Error{Class: ClassPath, Detail: fmt.Sprintf(format, a.project)}
}
