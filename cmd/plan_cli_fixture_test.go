package cmd

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// planCLIFixture runs a built ply binary with an isolated Git configuration.
type planCLIFixture struct{ binary, root, config string }

func newPlanCLI(t *testing.T, binary string) planCLIFixture {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	f := planCLIFixture{binary: binary, root: root, config: filepath.Join(root, "gitconfig")}
	mustPlanWrite(t, f.config, "[user]\nname = Plan CLI Test\nemail = cli@example.invalid\n", 0o600)
	return f
}

func (f planCLIFixture) environment(overrides map[string]string) []string {
	env := []string{}
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "GIT_") || key == "EMAIL" {
			continue
		}
		if _, ok := overrides[key]; !ok {
			env = append(env, entry)
		}
	}
	env = append(env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+f.config, "GIT_OPTIONAL_LOCKS=0")
	for key, value := range overrides {
		env = append(env, key+"="+value)
	}
	return env
}

func (f planCLIFixture) run(t *testing.T, cwd string, overrides map[string]string, args ...string) (string, string, int) {
	t.Helper()
	c := exec.Command(f.binary, args...)
	c.Dir = cwd
	c.Env = f.environment(overrides)
	var out, stderr bytes.Buffer
	c.Stdout = &out
	c.Stderr = &stderr
	err := c.Run()
	exit := 0
	if err != nil {
		var e *exec.ExitError
		if !errors.As(err, &e) {
			t.Fatal(err)
		}
		exit = e.ExitCode()
	}
	return out.String(), stderr.String(), exit
}

func (f planCLIFixture) success(t *testing.T, cwd string, args ...string) string {
	t.Helper()
	out, stderr, exit := f.run(t, cwd, nil, args...)
	if exit != 0 {
		t.Fatalf("CLI %v exited %d\n%s\n%s", args, exit, out, stderr)
	}
	return out
}

func (f planCLIFixture) failure(t *testing.T, cwd string, env map[string]string, want string, args ...string) string {
	t.Helper()
	out, stderr, exit := f.run(t, cwd, env, args...)
	if exit == 0 || out != "" || !strings.Contains(stderr, want) {
		t.Fatalf("CLI %v: exit=%d stdout=%q stderr=%q; want error %q", args, exit, out, stderr, want)
	}
	return stderr
}

func (f planCLIFixture) git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	a := append([]string{"-c", "core.hooksPath=" + os.DevNull, "-c", "commit.gpgsign=false", "-C", dir}, args...)
	c := exec.Command("git", a...)
	c.Env = f.environment(nil)
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func mustPlanWrite(t *testing.T, p, s string, mode fs.FileMode) {
	t.Helper()
	if err := os.WriteFile(p, []byte(s), mode); err != nil {
		t.Fatal(err)
	}
}
