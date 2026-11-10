package planningrepo

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Git runs one local command. Its separate stderr lets discovery distinguish an
// absent repository from a launch, permission, or repository-corruption failure.
type Git func(directory string, arguments ...string) (string, error)

type gitError struct {
	err    error
	stderr string
}

func (e *gitError) Error() string { return fmt.Sprintf("%v: %s", e.err, strings.TrimSpace(e.stderr)) }
func (e *gitError) Unwrap() error { return e.err }

func systemGit(directory string, arguments ...string) (string, error) {
	args := []string{"-c", "core.hooksPath=" + os.DevNull, "-c", "core.fsmonitor=false",
		"-c", "core.attributesFile=" + os.DevNull, "-c", "core.autocrlf=false",
		"-c", "commit.gpgSign=false", "-c", "maintenance.auto=false", "-c", "gc.auto=0",
		"-C", directory}
	command := exec.Command("git", append(args, arguments...)...)
	// Repository-routing variables must never redirect an operation into a source
	// repository. Keep explicit identity and config-file choices (also useful for
	// callers with an isolated Git configuration), but not injected -c settings.
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "GIT_") && !identityEnvironment(key) &&
			key != "GIT_CONFIG_GLOBAL" && key != "GIT_CONFIG_SYSTEM" && key != "GIT_CONFIG_NOSYSTEM" {
			continue
		}
		if key != "LC_ALL" {
			command.Env = append(command.Env, entry)
		}
	}
	command.Env = append(command.Env, "LC_ALL=C", "GIT_OPTIONAL_LOCKS=0", "GIT_ATTR_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0")
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		return "", &gitError{err: err, stderr: stderr.String()}
	}
	return strings.TrimSuffix(stdout.String(), "\n"), nil
}

func identityEnvironment(key string) bool {
	switch key {
	case "GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL", "GIT_AUTHOR_DATE", "GIT_COMMITTER_NAME", "GIT_COMMITTER_EMAIL", "GIT_COMMITTER_DATE":
		return true
	}
	return false
}

func notRepository(err error) bool {
	var gitErr *gitError
	var exitErr *exec.ExitError
	return errors.As(err, &gitErr) && errors.As(err, &exitErr) && exitErr.ExitCode() == 128 &&
		strings.Contains(gitErr.stderr, "not a git repository")
}
