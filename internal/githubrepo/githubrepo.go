// Package githubrepo ensures one private empty GitHub repository through the
// gh CLI and an explicit process runner, without local Git or configuration.
package githubrepo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

var (
	ErrInvalid  = errors.New("invalid")
	ErrConflict = errors.New("conflict")
	ErrUnknown  = errors.New("unknown")
)

// Repo is a visible private empty repository with its canonical observed SSH URL.
type Repo struct{ SSHURL string }

const (
	readbackTimeout = 10 * time.Second
	query           = "query($owner:String!,$name:String!){repository(owner:$owner,name:$name){nameWithOwner visibility isEmpty sshUrl}}"
)

var (
	ownerPattern = regexp.MustCompile(`^[A-Za-z0-9]+(-[A-Za-z0-9]+)*$`)
	namePattern  = regexp.MustCompile(`^[A-Za-z0-9_.][A-Za-z0-9_.-]{0,99}$`)
)

// observation is one read; Visible false means GitHub answered NOT_FOUND,
// which is never proof that the repository is absent.
type observation struct {
	Visible                           bool
	NameWithOwner, Visibility, SSHURL string
	IsEmpty                           bool
}

func fail(kind error, format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", kind, fmt.Sprintf(format, args...))
}

// Valid reports whether owner/name is a GitHub owner and repository name that
// EnsurePrivate accepts.
func Valid(owner, name string) bool {
	return len(owner) <= 39 && ownerPattern.MatchString(owner) && namePattern.MatchString(name) &&
		name != "." && name != ".." && !strings.HasSuffix(strings.ToLower(name), ".git")
}

// EnsurePrivate returns the private empty github.com repository owner/name,
// creating it once when it is not visible. It never edits, deletes, retries or
// renames, and its errors never contain gh output or a URL.
func EnsurePrivate(ctx context.Context, runner process.Runner, owner, name string) (Repo, bool, error) {
	if ctx == nil || runner == nil || !Valid(owner, name) {
		return Repo{}, false, fail(ErrInvalid, "context, runner and a valid GitHub owner/name are required")
	}
	slug := owner + "/" + name
	repo, err := read(ctx, runner, owner, name)
	if err != nil {
		return Repo{}, false, err
	}
	if repo.Visible {
		if reason := mismatch(slug, repo); reason != "" {
			return Repo{}, false, fail(ErrConflict, "GitHub repository %s %s", slug, reason)
		}
		return Repo{SSHURL: repo.SSHURL}, false, nil
	}
	_, code, createErr := gh(ctx, runner, "repo", "create", slug, "--private")
	readCtx, cancel := readbackContext(ctx)
	defer cancel()
	repo, err = read(readCtx, runner, owner, name)
	readback := "unknown"
	switch {
	case err != nil:
	case !repo.Visible:
		readback = "not visible"
	default:
		if readback = mismatch(slug, repo); readback == "" {
			readback = "private empty"
		}
	}
	if createErr != nil {
		return Repo{}, false, fmt.Errorf("%w; readback: %s", failed(ctx, "create "+slug, code), readback)
	}
	if readback != "private empty" {
		return Repo{}, false, fail(ErrUnknown, "readback %s: %s", slug, readback)
	}
	return Repo{SSHURL: repo.SSHURL}, true, nil
}

// mismatch returns the fixed reason a visible repository is not the private
// empty slug with its canonical SSH URL, or "" when it is.
func mismatch(slug string, repo observation) string {
	switch {
	case !strings.EqualFold(repo.NameWithOwner, slug):
		return "identity differs"
	case repo.Visibility != "PRIVATE":
		return "is " + repo.Visibility
	case !repo.IsEmpty:
		return "is not empty"
	case repo.SSHURL != "git@github.com:"+repo.NameWithOwner+".git":
		return "SSH URL differs"
	}
	return ""
}

func read(ctx context.Context, runner process.Runner, owner, name string) (observation, error) {
	op := "read " + owner + "/" + name
	stdout, code, err := gh(ctx, runner, "api", "graphql", "--hostname", "github.com",
		"-f", "query="+query, "-f", "owner="+owner, "-f", "name="+name)
	if err != nil && code != 1 || ctx.Err() != nil {
		return observation{}, failed(ctx, op, code)
	}
	repo, notFound, ok := decode(stdout)
	switch {
	case ok && err == nil && repo != nil:
		return *repo, nil
	case ok && err != nil && notFound:
		return observation{}, nil
	case err != nil:
		return observation{}, failed(ctx, op, code)
	}
	return observation{}, fail(ErrUnknown, "%s: malformed", op)
}

// decode accepts exactly one JSON object followed by EOF holding either a
// visible repository without errors, or a null repository with errors.
// notFound reports exactly one NOT_FOUND error on path ["repository"].
func decode(stdout []byte) (repo *observation, notFound, ok bool) {
	decoder := json.NewDecoder(bytes.NewReader(stdout))
	var top, data, fields map[string]json.RawMessage
	if decoder.Decode(&top) != nil || top == nil || decoder.Decode(&struct{}{}) != io.EOF ||
		json.Unmarshal(top["data"], &data) != nil || len(data) != 1 || data["repository"] == nil {
		return nil, false, false
	}
	for key := range top {
		if key != "data" && key != "errors" {
			return nil, false, false
		}
	}
	if string(data["repository"]) == "null" {
		var errs []map[string]json.RawMessage
		if json.Unmarshal(top["errors"], &errs) != nil || errs == nil {
			return nil, false, false
		}
		for _, entry := range errs {
			for key := range entry {
				if key != "type" && key != "path" && key != "message" && key != "locations" {
					return nil, false, false
				}
			}
		}
		var kind string
		var path []string
		notFound = len(errs) == 1 && json.Unmarshal(errs[0]["type"], &kind) == nil && kind == "NOT_FOUND" &&
			json.Unmarshal(errs[0]["path"], &path) == nil && len(path) == 1 && path[0] == "repository"
		return nil, notFound, true
	}
	if top["errors"] != nil || json.Unmarshal(data["repository"], &fields) != nil || len(fields) != 4 {
		return nil, false, false
	}
	var nameWithOwner, visibility, sshURL *string
	var isEmpty *bool
	for key, target := range map[string]interface{}{"nameWithOwner": &nameWithOwner, "visibility": &visibility, "sshUrl": &sshURL, "isEmpty": &isEmpty} {
		if fields[key] == nil || json.Unmarshal(fields[key], target) != nil {
			return nil, false, false
		}
	}
	if nameWithOwner == nil || *nameWithOwner == "" || sshURL == nil || *sshURL == "" || isEmpty == nil || visibility == nil ||
		*visibility != "PRIVATE" && *visibility != "PUBLIC" && *visibility != "INTERNAL" {
		return nil, false, false
	}
	return &observation{true, *nameWithOwner, *visibility, *sshURL, *isEmpty}, false, true
}

// readbackContext keeps ctx, or bounds a fresh context when ctx has ended.
func readbackContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx.Err() != nil {
		return context.WithTimeout(context.Background(), readbackTimeout)
	}
	return ctx, func() {}
}

// failed reports a gh operation failure without its output.
func failed(ctx context.Context, op string, code int) error {
	switch {
	case ctx.Err() != nil:
		return fail(ErrUnknown, "%s: canceled", op)
	case code > 0:
		return fail(ErrUnknown, "%s: exit %d", op, code)
	}
	return fail(ErrUnknown, "%s: failed", op)
}

func gh(ctx context.Context, runner process.Runner, args ...string) ([]byte, int, error) {
	var stdout bytes.Buffer
	err := runner.Run(process.Command{Name: "gh", Args: args, Context: ctx, Env: ghEnv(os.Environ()),
		Stdout: &stdout, Stderr: io.Discard})
	code := 0
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code = exitErr.ExitCode()
	} else if err != nil {
		code = -1
	}
	return stdout.Bytes(), code, err
}

// ghEnv drops inherited variables that redirect gh or Git or change its
// output, and pins gh to github.com without prompts, color or notifiers.
func ghEnv(environ []string) []string {
	pinned := []string{"GH_HOST=github.com", "GH_PROMPT_DISABLED=1", "NO_COLOR=1", "GH_NO_UPDATE_NOTIFIER=1", "GH_NO_EXTENSION_UPDATE_NOTIFIER=1"}
	dropped := map[string]bool{"GH_REPO": true, "GH_DEBUG": true, "DEBUG": true, "GH_FORCE_TTY": true, "CLICOLOR_FORCE": true,
		"GH_PAGER": true, "PAGER": true, "GIT_DIR": true, "GIT_WORK_TREE": true}
	for _, entry := range pinned {
		dropped[strings.SplitN(entry, "=", 2)[0]] = true
	}
	env := make([]string, 0, len(environ)+len(pinned))
	for _, entry := range environ {
		name := strings.SplitN(entry, "=", 2)[0]
		if !dropped[name] && !strings.HasPrefix(name, "GIT_CONFIG") {
			env = append(env, entry)
		}
	}
	return append(env, pinned...)
}
