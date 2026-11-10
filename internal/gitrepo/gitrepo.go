// Package gitrepo observes and initializes one local Git repository through an
// explicit process runner, without workspace or configuration side effects.
package gitrepo

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

var (
	ErrInvalid  = errors.New("invalid")
	ErrConflict = errors.New("conflict")
	ErrUnknown  = errors.New("unknown")
)

type State int

const (
	Absent State = iota
	Root
	Nested
	NotWorktree
)

func (s State) String() string {
	return [...]string{"absent", "root", "nested", "not-worktree"}[s]
}

// Repo is an observed repository; an empty Branch means a detached HEAD.
type Repo struct {
	State                 State
	Dir, Toplevel, Branch string
}

const readbackTimeout = 2 * time.Second

func fail(kind error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", kind, fmt.Sprintf(format, args...))
}

// Observe reports the repository state of the absolute dir without requiring a commit.
func Observe(ctx context.Context, runner process.Runner, dir string) (Repo, error) {
	if ctx == nil || runner == nil || !filepath.IsAbs(dir) {
		return Repo{}, fail(ErrInvalid, "context, runner and absolute directory are required")
	}
	path := dir
	if _, err := os.Lstat(dir); err != nil {
		path = filepath.Dir(dir)
	}
	return observe(ctx, runner, dir, path)
}

func observe(ctx context.Context, runner process.Runner, dir, path string) (Repo, error) {
	repo := Repo{Dir: dir}
	inside, stderr, code, err := git(ctx, runner, path, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		if code == 128 && ctx.Err() == nil && strings.Contains(stderr, "not a git repository") {
			return repo, nil
		}
		return repo, fail(ErrUnknown, "observe %s: %v: %s", path, err, strings.TrimSpace(stderr))
	}
	switch inside {
	case "false":
		repo.State = NotWorktree
		return repo, nil
	case "true":
	default:
		return repo, fail(ErrUnknown, "observe %s: unexpected output %q", path, inside)
	}
	repo.Toplevel, stderr, _, err = git(ctx, runner, path, "rev-parse", "--path-format=absolute", "--show-toplevel")
	if err != nil {
		return repo, fail(ErrUnknown, "toplevel %s: %v: %s", path, err, strings.TrimSpace(stderr))
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return repo, fail(ErrUnknown, "resolve %s: %v", path, err)
	}
	if path != dir || repo.Toplevel != resolved {
		repo.State = Nested
		return repo, nil
	}
	repo.State = Root
	branch, stderr, code, err := git(ctx, runner, path, "symbolic-ref", "-q", "--short", "HEAD")
	if err != nil && code != 1 {
		return repo, fail(ErrUnknown, "branch %s: %v: %s", path, err, strings.TrimSpace(stderr))
	}
	repo.Branch = branch
	return repo, nil
}

const (
	OpInsideWorkTree = "inside-work-tree"
	OpToplevel       = "toplevel"
	OpCommonDir      = "common-dir"
)

// Identity is Git's raw answer for a directory; Worktree false means Git said "false".
type Identity struct {
	Worktree            bool
	Toplevel, CommonDir string
}

// Error is an ErrUnknown Identify step failure; Err is the original run error, nil when malformed.
type Error struct {
	Op  string
	Err error
	msg string
}

func (e *Error) Error() string        { return e.msg }
func (e *Error) Is(target error) bool { return target == ErrUnknown }
func (e *Error) Unwrap() error        { return e.Err }

// Identify reads worktree status, toplevel and common dir of the absolute dir
// without classifying them. A common-dir failure keeps the partial Identity.
func Identify(ctx context.Context, runner process.Runner, dir string) (Identity, error) {
	if ctx == nil || runner == nil || !filepath.IsAbs(dir) {
		return Identity{}, fail(ErrInvalid, "context, runner and absolute directory are required")
	}
	var identity Identity
	inside, err := identityLine(ctx, runner, dir, OpInsideWorkTree, "--is-inside-work-tree")
	switch {
	case err != nil:
		return identity, err
	case inside == "false":
		return identity, nil
	case inside != "true":
		return identity, &Error{Op: OpInsideWorkTree, msg: fmt.Sprintf("%v: %s %s: malformed", ErrUnknown, OpInsideWorkTree, dir)}
	}
	identity.Worktree = true
	if identity.Toplevel, err = identityLine(ctx, runner, dir, OpToplevel, "--path-format=absolute", "--show-toplevel"); err != nil {
		return identity, err
	}
	identity.CommonDir, err = identityLine(ctx, runner, dir, OpCommonDir, "--path-format=absolute", "--git-common-dir")
	return identity, err
}

func identityLine(ctx context.Context, runner process.Runner, dir, op string, args ...string) (string, error) {
	stdout, _, _, err := run(ctx, runner, dir, append([]string{"rev-parse"}, args...)...)
	if err != nil {
		return "", &Error{Op: op, Err: err, msg: fmt.Sprintf("%v: %s %s: %v", ErrUnknown, op, dir, err)}
	}
	value, ok := SingleLine(stdout)
	if !ok || value == "" {
		return "", &Error{Op: op, msg: fmt.Sprintf("%v: %s %s: malformed", ErrUnknown, op, dir)}
	}
	return value, nil
}

// SingleLine returns output without its one trailing LF or CRLF, or false
// when output is not exactly one terminated line.
func SingleLine(output []byte) (string, bool) {
	value := string(output)
	switch {
	case strings.HasSuffix(value, "\r\n"):
		value = strings.TrimSuffix(value, "\r\n")
	case strings.HasSuffix(value, "\n"):
		value = strings.TrimSuffix(value, "\n")
	default:
		return "", false
	}
	if strings.ContainsAny(value, "\r\n") {
		return "", false
	}
	return value, true
}

// EnsureInit initializes an empty or missing dir on branch, or verifies an
// existing matching repository. It never adopts files or changes branches.
func EnsureInit(ctx context.Context, runner process.Runner, dir, branch string) (Repo, bool, error) {
	if ctx == nil || runner == nil || !filepath.IsAbs(dir) || branch == "" || strings.HasPrefix(branch, "-") {
		return Repo{}, false, fail(ErrInvalid, "context, runner, absolute directory and branch are required")
	}
	dir = filepath.Clean(dir)
	parent := filepath.Dir(dir)
	path := dir
	info, err := os.Lstat(dir)
	switch {
	case os.IsNotExist(err):
		if info, err = os.Stat(parent); os.IsNotExist(err) {
			return Repo{}, false, fail(ErrInvalid, "parent %s does not exist", parent)
		} else if err != nil || !info.IsDir() {
			return Repo{}, false, fail(ErrConflict, "parent %s is not a directory", parent)
		}
		path = parent
	case err != nil:
		return Repo{}, false, fail(ErrUnknown, "inspect %s: %v", dir, err)
	case info.Mode()&os.ModeSymlink != 0 || !info.IsDir():
		return Repo{}, false, fail(ErrConflict, "%s is a symlink or not a directory", dir)
	}
	checked, stderr, code, err := git(ctx, runner, path, "check-ref-format", "--branch", branch)
	rejected := code == 128 && ctx.Err() == nil && strings.Contains(stderr, fmt.Sprintf("'%s' is not a valid branch name", branch))
	if err == nil && checked != branch || code == 1 || rejected {
		return Repo{}, false, fail(ErrInvalid, "branch %q is not a valid branch name", branch)
	} else if err != nil {
		return Repo{}, false, fail(ErrUnknown, "check branch %q: %v: %s", branch, err, strings.TrimSpace(stderr))
	}
	if path == parent {
		above, err := observe(ctx, runner, parent, parent)
		if err != nil {
			return Repo{}, false, err
		}
		if above.State != Absent {
			return Repo{}, false, fail(ErrConflict, "%s would be nested in a Git repository (%s)", dir, above.State)
		}
		if err := os.Mkdir(dir, 0o755); err != nil {
			return Repo{}, false, fail(ErrUnknown, "create %s: %v", dir, err)
		}
		return initialize(ctx, runner, dir, branch)
	}
	repo, err := observe(ctx, runner, dir, dir)
	if err != nil {
		return Repo{}, false, err
	}
	switch {
	case repo.State == Root && repo.Branch == branch:
		return repo, false, nil
	case repo.State == Root:
		return Repo{}, false, fail(ErrConflict, "%s is on branch %q, want %q", dir, repo.Branch, branch)
	case repo.State != Absent:
		return Repo{}, false, fail(ErrConflict, "%s is %s in a Git repository", dir, repo.State)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Repo{}, false, fail(ErrUnknown, "read %s: %v", dir, err)
	}
	if len(entries) > 0 {
		return Repo{}, false, fail(ErrConflict, "%s is not empty", dir)
	}
	return initialize(ctx, runner, dir, branch)
}

func initialize(ctx context.Context, runner process.Runner, dir, branch string) (Repo, bool, error) {
	_, stderr, _, initErr := git(ctx, runner, dir, "init", "-q", "--initial-branch="+branch, "--template=")
	readCtx, cancel := readbackContext(ctx)
	defer cancel()
	repo, err := observe(readCtx, runner, dir, dir)
	if initErr != nil {
		return Repo{}, false, fail(ErrUnknown, "git init %s: %v: %s; readback %s branch %q: %v",
			dir, initErr, strings.TrimSpace(stderr), repo.State, repo.Branch, err)
	}
	if err != nil || repo.State != Root || repo.Branch != branch {
		return Repo{}, false, fail(ErrUnknown, "readback %s: %s branch %q: %v", dir, repo.State, repo.Branch, err)
	}
	return repo, true, nil
}

// readbackContext keeps ctx, or bounds a fresh context when ctx has ended.
func readbackContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx.Err() != nil {
		return context.WithTimeout(context.Background(), readbackTimeout)
	}
	return ctx, func() {}
}

var remoteName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`)

// ValidRemoteName reports whether name is a remote name that EnsureRemote accepts.
func ValidRemoteName(name string) bool {
	return remoteName.MatchString(name) && !strings.Contains(name, "..") && !strings.HasSuffix(name, ".") && !strings.HasSuffix(name, ".lock")
}

func validRemote(name, url string) bool {
	if !ValidRemoteName(name) || url == "" || strings.HasPrefix(url, "-") {
		return false
	}
	return strings.IndexFunc(url, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) < 0
}

// EnsureRemote adds the local remote name with url to the repository rooted at
// dir, or verifies that exactly that local url is already set. It never
// overwrites a remote, and its errors never contain a URL or Git output.
func EnsureRemote(ctx context.Context, runner process.Runner, dir, name, url string) (bool, error) {
	if ctx == nil || runner == nil || !filepath.IsAbs(dir) || !validRemote(name, url) {
		return false, fail(ErrInvalid, "context, runner, absolute directory, valid remote name and URL are required")
	}
	dir = filepath.Clean(dir)
	info, err := os.Lstat(dir)
	switch {
	case os.IsNotExist(err):
		return false, fail(ErrInvalid, "%s does not exist", dir)
	case err != nil:
		return false, fail(ErrUnknown, "inspect %s: failed", dir)
	case info.Mode()&os.ModeSymlink != 0 || !info.IsDir():
		return false, fail(ErrConflict, "%s is a symlink or not a directory", dir)
	}
	repo, err := Observe(ctx, runner, dir)
	if err != nil {
		return false, failed(ctx, "observe "+dir, -1)
	}
	if repo.State != Root {
		return false, fail(ErrConflict, "%s is %s, not a Git repository root", dir, repo.State)
	}
	missing, err := readRemote(ctx, runner, dir, name, url)
	if err != nil || !missing {
		return false, err
	}
	_, _, code, addErr := run(ctx, runner, dir, "remote", "add", "--", name, url)
	readCtx, cancel := readbackContext(ctx)
	defer cancel()
	missing, err = readRemote(readCtx, runner, dir, name, url)
	readback := "set"
	if err != nil {
		readback = err.Error()
	} else if missing {
		readback = "absent"
	}
	if addErr != nil {
		return false, fmt.Errorf("%w; readback: %s", failed(ctx, "remote add "+name+" in "+dir, code), readback)
	}
	if readback != "set" {
		return false, fail(ErrUnknown, "readback remote %s in %s: %s", name, dir, readback)
	}
	return true, nil
}

// readRemote reports whether the remote is absent, nil when it is set to
// exactly url in local config only, and ErrConflict for any other remote.
func readRemote(ctx context.Context, runner process.Runner, dir, name, url string) (bool, error) {
	stdout, stderr, code, err := run(ctx, runner, dir, "config", "--show-scope", "--null", "--get-regexp", `^remote\.`+regexp.QuoteMeta(name)+`\.`)
	if err != nil && (code != 1 || len(stdout) > 0 || stderr != "") {
		return false, failed(ctx, "read config "+dir, code)
	}
	var urls []string
	var records int
	reason := ""
	fields := strings.Split(string(stdout), "\x00")
	if len(fields)%2 == 0 || fields[len(fields)-1] != "" {
		return false, fail(ErrUnknown, "read config %s: malformed", dir)
	}
	for i := 0; i+1 < len(fields); i += 2 {
		key, value, ok := strings.Cut(fields[i+1], "\n")
		dot := strings.LastIndex(key, ".")
		switch fields[i] {
		case "local", "worktree", "global", "system", "command":
		default:
			ok = false
		}
		if !ok || !strings.HasPrefix(key, "remote.") || dot <= len("remote.") {
			return false, fail(ErrUnknown, "read config %s: malformed", dir)
		}
		if key[len("remote."):dot] != name {
			continue
		}
		records++
		switch variable := key[dot+1:]; {
		case fields[i] != "local":
			reason = "has non-local configuration"
		case variable == "pushurl":
			reason = "has a push URL"
		case variable == "url":
			urls = append(urls, value)
		}
	}
	output, stderr, code, err := run(ctx, runner, dir, "remote", "get-url", name)
	effective, ok := SingleLine(output)
	switch {
	// Absence is only Git's exact one-line diagnostic for name, LF-terminated, with no stdout.
	case err != nil && code == 2 && ctx.Err() == nil && len(output) == 0 && stderr == "error: No such remote '"+name+"'\n":
		if records == 0 {
			return true, nil
		}
		reason = "is incomplete"
	case err != nil:
		return false, failed(ctx, "get-url "+name+" in "+dir, code)
	case !ok || effective == "":
		return false, fail(ErrUnknown, "get-url %s in %s: malformed", name, dir)
	case records == 0:
		reason = "exists without config"
	case reason != "":
	case len(urls) != 1:
		reason = fmt.Sprintf("has %d URLs", len(urls))
	case urls[0] != url:
		reason = "has a different URL"
	case effective != urls[0]:
		reason = "URL is rewritten"
	}
	if reason != "" {
		return false, fail(ErrConflict, "Git remote %s in %s %s", name, dir, reason)
	}
	return false, nil
}

// failed reports a Git operation failure without its output.
func failed(ctx context.Context, op string, code int) error {
	switch {
	case ctx.Err() != nil:
		return fail(ErrUnknown, "%s: canceled", op)
	case code > 0:
		return fail(ErrUnknown, "%s: exit %d", op, code)
	}
	return fail(ErrUnknown, "%s: failed", op)
}

func git(ctx context.Context, runner process.Runner, path string, args ...string) (string, string, int, error) {
	stdout, stderr, code, err := run(ctx, runner, path, args...)
	return strings.TrimSpace(string(stdout)), stderr, code, err
}

// Read runs one read-only git command in dir with the conventions of this
// package, optional locks disabled, so that no index refresh is written, and
// for this call only no configured file system monitor, so that no hook runs.
// It returns raw stdout, stderr, the exit code (-1 when Git did not run) and the run error.
func Read(ctx context.Context, runner process.Runner, dir string, args ...string) ([]byte, string, int, error) {
	return runEnv(ctx, runner, dir, []string{"GIT_OPTIONAL_LOCKS=0"}, append([]string{"-c", "core.fsmonitor=false"}, args...)...)
}

func run(ctx context.Context, runner process.Runner, path string, args ...string) ([]byte, string, int, error) {
	return runEnv(ctx, runner, path, nil, args...)
}

func runEnv(ctx context.Context, runner process.Runner, path string, extra []string, args ...string) ([]byte, string, int, error) {
	var stdout, stderr bytes.Buffer
	err := runner.Run(process.Command{Name: "git", Args: append([]string{"-C", path}, args...),
		Context: ctx, Env: append(gitEnv(os.Environ()), extra...), Stdout: &stdout, Stderr: &stderr})
	code := 0
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code = exitErr.ExitCode()
	} else if err != nil {
		code = -1
	}
	return stdout.Bytes(), stderr.String(), code, err
}

// gitEnv drops inherited GIT_* variables except config-file selection.
func gitEnv(environ []string) []string {
	env := make([]string, 0, len(environ)+2)
	for _, entry := range environ {
		name := strings.SplitN(entry, "=", 2)[0]
		if strings.HasPrefix(name, "GIT_") && name != "GIT_CONFIG_GLOBAL" && name != "GIT_CONFIG_SYSTEM" && name != "GIT_CONFIG_NOSYSTEM" {
			continue
		}
		env = append(env, entry)
	}
	return append(env, "LC_ALL=C", "GIT_TERMINAL_PROMPT=0")
}
