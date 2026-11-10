// Package reporemove undoes the local part of repoestablish for one
// repository: it clears the wrapper metadata, removes the definitions member
// and deletes the wrapper directory, only when the checks prove that no local
// data is lost. The remote repository is never contacted for a change. Every
// read and effect goes through Dependencies; the package has no output,
// flags, working directory or environment.
package reporemove

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// Request is the caller's claim. Preview is the default; Apply requires the
// previewed plan digest as Expect and an Actor.
type Request struct {
	Workspace, Project, RepoName, Remote, Actor string
	Apply                                       bool
	Expect                                      string
}

// Target is what Locate reads without effect.
type Target struct {
	ProjectRoot    string   // physical Project root
	Remote         string   // the Project defaults' remote; empty when not configured
	Path           string   // stored definitions path of the repository in the Project; empty when not defined
	WrapperPresent bool     // <ProjectRoot>/<RepoName> exists
	Quarantines    []string // existing <ProjectRoot>/.ply-repo-remove-<RepoName>-* paths
}

// Resolved is built by Run only, and is the sole input of Inspect and the effects.
type Resolved struct {
	Workspace, Project, ProjectRoot, RepoName string
	Wrapper, RepoDir, Leaf, Path              string // ProjectRoot/RepoName, Wrapper/Leaf, stored definitions path
	Remote, Actor                             string
}

// Check is one passed pre-effect check, as shown in preview.
type Check struct{ Name, Detail string }

// Plan is the result of all checks; Digest covers what apply must find unchanged.
type Plan struct {
	Checks            []Check
	RemoteURL         string // empty when the remote has no single URL
	MetadataRevision  int
	MetadataCleared   bool // the wrapper metadata is already cleared by an earlier run
	DefinitionsDigest string
	Verdict           string   // VerdictUnborn or VerdictProven
	GitFiles          []string // local configuration files deleted with .git, relative to it
	LocalDigest       string   // HEAD, refs and wrapper listing; apply requires them unchanged
	Digest            string
}

const DefaultRemote = "origin"

const (
	VerdictUnborn = "unborn"
	VerdictProven = "proven"
)

type Step string

const (
	StepPrepare  Step = "prepare"
	StepCheck    Step = "check"
	StepMetadata Step = "metadata"
	StepDefine   Step = "define"
	StepWrapper  Step = "wrapper"
)

type StepResult struct {
	Step    Step
	Changed bool
}

type Result struct {
	Resolved Resolved     // set once the target is resolved
	Plan     Plan         // set once all checks have passed
	Nothing  bool         // the repository is neither defined nor present: no effect, no error
	Steps    []StepResult // confirmed completed effect steps only, in fixed order
}

var ErrInvalid, ErrConflict, ErrUnknown = errors.New("invalid"), errors.New("conflict"), errors.New("unknown")

// Error names the step where Run stopped. ErrInvalid and ErrConflict mean that
// step had no effect; ErrUnknown means it may have had one; a bare context
// error means it was not started.
type Error struct {
	Step Step
	Err  error
}

func (e *Error) Error() string { return fmt.Sprintf("%s: %v", e.Step, e.Err) }
func (e *Error) Unwrap() error { return e.Err }

// Dependencies are function fields, all required. Locate and Inspect have no
// effect. Each effect function's error must wrap exactly one of the three
// sentinels.
type Dependencies struct {
	Locate        func(ctx context.Context, req Request) (Target, error)
	Inspect       func(ctx context.Context, r Resolved) (Plan, error)
	ClearWrapper  func(ctx context.Context, r Resolved, p Plan) (bool, error)
	Undefine      func(ctx context.Context, r Resolved, p Plan) (bool, error)
	RemoveWrapper func(ctx context.Context, r Resolved, p Plan) (bool, error)
}

func (d Dependencies) complete() bool {
	return d.Locate != nil && d.Inspect != nil && d.ClearWrapper != nil && d.Undefine != nil && d.RemoveWrapper != nil
}

func invalid(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrInvalid, fmt.Sprintf(format, args...))
}

func conflict(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrConflict, fmt.Sprintf(format, args...))
}

// Run previews or applies the removal of one repository. Apply recomputes the
// plan and stops on any difference from the previewed digest; it then clears
// the wrapper metadata, removes the definitions member and deletes the
// wrapper, in that order, and rolls nothing back.
func Run(ctx context.Context, req Request, d Dependencies) (Result, error) {
	var result Result
	if ctx == nil || !d.complete() {
		return result, &Error{StepPrepare, invalid("context and all dependencies are required")}
	}
	if err := validateRequest(req); err != nil {
		return result, &Error{StepPrepare, err}
	}
	if err := ctx.Err(); err != nil {
		return result, &Error{StepPrepare, err}
	}
	target, err := d.Locate(ctx, req)
	if err != nil {
		return result, &Error{StepPrepare, withoutEffect(err)}
	}
	r, err := resolve(req, target)
	if err != nil {
		return result, &Error{StepPrepare, err}
	}
	result.Resolved = r
	switch {
	case len(target.Quarantines) > 0:
		return result, &Error{StepCheck, conflict("earlier removal left %s; inspect it", strings.Join(target.Quarantines, ", "))}
	case target.Path == "" && !target.WrapperPresent:
		result.Nothing = true
		return result, nil
	case target.Path == "":
		return result, &Error{StepCheck, conflict("repository %s is not defined in Project %s but wrapper %s exists; its repository can no longer be proven. Run 'ply repo create %s' to re-adopt a matching repository, then remove it again",
			r.RepoName, r.Project, r.Wrapper, r.RepoName)}
	}
	if err := ctx.Err(); err != nil {
		return result, &Error{StepCheck, err}
	}
	plan, err := d.Inspect(ctx, r)
	if err != nil {
		return result, &Error{StepCheck, withoutEffect(err)}
	}
	result.Plan = plan
	if !req.Apply {
		return result, nil
	}
	if plan.Digest != req.Expect {
		return result, &Error{StepCheck, conflict("plan digest is %s, expected %s; preview again", plan.Digest, req.Expect)}
	}
	steps := []struct {
		step Step
		run  func(context.Context, Resolved, Plan) (bool, error)
	}{{StepMetadata, d.ClearWrapper}, {StepDefine, d.Undefine}, {StepWrapper, d.RemoveWrapper}}
	for _, s := range steps {
		if err := ctx.Err(); err != nil {
			return result, &Error{s.step, err}
		}
		changed, err := s.run(ctx, r, plan)
		if err != nil {
			return result, &Error{s.step, withPossibleEffect(err)}
		}
		result.Steps = append(result.Steps, StepResult{s.step, changed})
	}
	return result, nil
}

// withoutEffect keeps an ErrInvalid or ErrConflict from a reading step and
// classifies anything else, including ErrUnknown, as ErrInvalid.
func withoutEffect(err error) error {
	switch {
	case errors.Is(err, ErrUnknown):
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	case errors.Is(err, ErrInvalid), errors.Is(err, ErrConflict):
		return err
	}
	return fmt.Errorf("%w: %w", ErrInvalid, err)
}

// withPossibleEffect keeps a classified effect error; an unclassified one,
// even a context error, may have had an effect.
func withPossibleEffect(err error) error {
	if errors.Is(err, ErrUnknown) || errors.Is(err, ErrInvalid) || errors.Is(err, ErrConflict) {
		return err
	}
	return fmt.Errorf("%w: %w", ErrUnknown, err)
}

// element reports whether value is exactly one relative path element.
func element(value string) bool {
	return value != "" && value != "." && value != ".." && !strings.ContainsRune(value, '/') &&
		!strings.ContainsRune(value, filepath.Separator) && !filepath.IsAbs(value)
}

func validateRequest(req Request) error {
	switch {
	case !filepath.IsAbs(req.Workspace) || filepath.Clean(req.Workspace) != req.Workspace:
		return invalid("workspace %q must be an absolute clean path", req.Workspace)
	case req.Project == "":
		return invalid("Project is required")
	case !element(req.RepoName):
		return invalid("repository name %q must be one path element", req.RepoName)
	case req.Apply && req.Expect == "":
		return invalid("--apply requires --expect with the previewed plan digest")
	case req.Apply && req.Actor == "":
		return invalid("actor is required; set --actor or PLY_ACTOR")
	case strings.HasPrefix(req.Remote, "-"):
		return invalid("remote must not begin with -")
	}
	return nil
}

// resolve takes the leaf only from the stored definitions path, which must be
// <Project root>/<repository>/<leaf> relative to the workspace.
func resolve(req Request, target Target) (Resolved, error) {
	r := Resolved{Workspace: req.Workspace, Project: req.Project, ProjectRoot: target.ProjectRoot, RepoName: req.RepoName,
		Path: target.Path, Remote: req.Remote, Actor: req.Actor}
	if r.Remote == "" {
		r.Remote = target.Remote
	}
	if r.Remote == "" {
		r.Remote = DefaultRemote
	}
	if !filepath.IsAbs(r.ProjectRoot) || filepath.Clean(r.ProjectRoot) != r.ProjectRoot {
		return Resolved{}, invalid("Project root %q must be an absolute clean path", r.ProjectRoot)
	}
	r.Wrapper = filepath.Join(r.ProjectRoot, r.RepoName)
	if target.Path == "" {
		return r, nil
	}
	prefix, err := filepath.Rel(r.Workspace, r.Wrapper)
	if err != nil {
		return Resolved{}, invalid("wrapper %s is not relative to workspace %s", r.Wrapper, r.Workspace)
	}
	r.Leaf = strings.TrimPrefix(target.Path, filepath.ToSlash(prefix)+"/")
	if r.Leaf == target.Path || !element(r.Leaf) {
		return Resolved{}, invalid("repository %s is defined at %s, not at %s/<leaf> as ply repo create defines it", r.RepoName, target.Path, filepath.ToSlash(prefix))
	}
	r.RepoDir = filepath.Join(r.Wrapper, r.Leaf)
	return r, nil
}
