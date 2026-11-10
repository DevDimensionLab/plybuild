package githubrepo

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

// fakeGHScript logs argv and the gh-targeting environment (never auth values),
// then answers NOT_FOUND until "repo create" succeeds, and private empty after.
const fakeGHScript = `#!/bin/sh
{ printf 'argv'; printf ' %s' "$@"; printf '\n'
  env | grep -E '^(GH_HOST|GH_REPO|GH_DEBUG|GH_PROMPT_DISABLED|NO_COLOR|GH_NO_UPDATE_NOTIFIER|GH_NO_EXTENSION_UPDATE_NOTIFIER|GH_CONFIG_DIR)=' | sort
} >> "$FAKE_GH_DIR/log"
if [ "$1" = repo ]; then touch "$FAKE_GH_DIR/created"; echo https://github.com/Org/svc; exit 0; fi
if [ -e "$FAKE_GH_DIR/created" ]; then
  echo '{"data":{"repository":{"nameWithOwner":"Org/svc","visibility":"PRIVATE","isEmpty":true,"sshUrl":"git@github.com:Org/svc.git"}}}'
  exit 0
fi
echo '{"data":{"repository":null},"errors":[{"type":"NOT_FOUND","path":["repository"],"message":"not found"}]}'
exit 1
`

func TestEnsurePrivateWithFakeGHExecutableKeepsConfigAndIgnoresInheritedTarget(t *testing.T) {
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	bin, configDir, cwd := filepath.Join(base, "bin"), filepath.Join(base, "gh-config"), filepath.Join(base, "cwd")
	for _, dir := range []string{bin, configDir, cwd} {
		if err := os.Mkdir(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write(t, filepath.Join(bin, "gh"), fakeGHScript, 0o755)
	write(t, filepath.Join(configDir, "config.yml"), "git_protocol: ssh\n", 0o600)
	write(t, filepath.Join(configDir, "hosts.yml"), "github.com:\n    user: fake\n", 0o600)
	gitConfig := filepath.Join(base, "gitconfig")
	write(t, gitConfig, "", 0o644)
	t.Setenv("HOME", base)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CONFIG_GLOBAL", gitConfig)
	for _, args := range [][]string{{"init", "-q"}, {"remote", "add", "origin", "git@github.com:Other/x.git"}} {
		if output, err := exec.Command("git", append([]string{"-C", cwd}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, output)
		}
	}
	files := []string{filepath.Join(configDir, "config.yml"), filepath.Join(configDir, "hosts.yml"), filepath.Join(cwd, ".git", "config")}
	before := contents(t, files)
	t.Chdir(cwd)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("FAKE_GH_DIR", base)
	t.Setenv("GH_CONFIG_DIR", configDir)
	t.Setenv("GH_REPO", "evil/x")
	t.Setenv("GH_HOST", "example.com")

	repo, created, err := EnsurePrivate(context.Background(), process.SystemRunner().Runner, "Org", "svc")

	if err != nil || !created || repo.SSHURL != "git@github.com:Org/svc.git" {
		t.Fatalf("EnsurePrivate = %#v, %v, %v", repo, created, err)
	}
	env := "GH_CONFIG_DIR=" + configDir + "\nGH_HOST=github.com\nGH_NO_EXTENSION_UPDATE_NOTIFIER=1\nGH_NO_UPDATE_NOTIFIER=1\nGH_PROMPT_DISABLED=1\nNO_COLOR=1\n"
	read := "argv " + strings.Join(readArgs("Org", "svc"), " ") + "\n" + env
	if log := contents(t, []string{filepath.Join(base, "log")})[0]; log != read+"argv repo create Org/svc --private\n"+env+read {
		t.Fatalf("fake gh log:\n%s", log)
	}
	if after := contents(t, files); strings.Join(after, "\x00") != strings.Join(before, "\x00") {
		t.Fatalf("gh or Git config changed: %q -> %q", before, after)
	}
}

func write(t *testing.T, path, content string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func contents(t *testing.T, paths []string) []string {
	t.Helper()
	var result []string
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		result = append(result, string(content))
	}
	return result
}
