package gitrepo

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/adapter/process"
)

type step struct {
	path, args     string
	stdout, stderr string
	code           int
	err            error
}

type fakeRunner struct {
	t     *testing.T
	steps []step
	calls []process.Command
	after func(process.Command)
}

func (f *fakeRunner) Run(c process.Command) error {
	f.t.Helper()
	i := len(f.calls)
	f.calls = append(f.calls, c)
	if i >= len(f.steps) {
		f.t.Fatalf("unexpected call %d: %v", i, c.Args)
	}
	s := f.steps[i]
	if got, want := strings.Join(c.Args, " "), "-C "+s.path+" "+s.args; c.Name != "git" || got != want {
		f.t.Fatalf("call %d = %s %s, want git %s", i, c.Name, got, want)
	}
	if c.Dir != "" || c.Context == nil || !filteredEnv(c.Env) {
		f.t.Fatalf("call %d request has Dir %q, Context %v, Env %v", i, c.Dir, c.Context, c.Env)
	}
	_, _ = io.WriteString(c.Stdout, s.stdout)
	_, _ = io.WriteString(c.Stderr, s.stderr)
	if f.after != nil {
		f.after(c)
	}
	if s.code != 0 {
		return exitError(f.t, s.code)
	}
	return s.err
}

func filteredEnv(env []string) bool {
	joined := "\n" + strings.Join(env, "\n") + "\n"
	return strings.Contains(joined, "\nLC_ALL=C\n") && strings.Contains(joined, "\nGIT_TERMINAL_PROMPT=0\n") &&
		!strings.Contains(joined, "\nGIT_DIR=")
}

func exitError(t *testing.T, code int) error {
	err := exec.Command("/bin/sh", "-c", "exit "+strconv.Itoa(code)).Run()
	if _, ok := err.(*exec.ExitError); !ok {
		t.Fatalf("could not produce exit %d: %v", code, err)
	}
	return err
}

func cat(groups ...[]step) []step {
	var all []step
	for _, group := range groups {
		all = append(all, group...)
	}
	return all
}

func refOK(path, branch string) []step {
	return []step{{path: path, args: "check-ref-format --branch " + branch, stdout: branch + "\n"}}
}

func inside(path, stdout string) []step {
	return []step{{path: path, args: "rev-parse --is-inside-work-tree", stdout: stdout + "\n"}}
}

func insideFails(path, stderr string, code int, err error) []step {
	return []step{{path: path, args: "rev-parse --is-inside-work-tree", stderr: stderr, code: code, err: err}}
}

func absent(path string) []step {
	return insideFails(path, "fatal: not a git repository (or any of the parent directories): .git\n", 128, nil)
}

func toplevel(path, top string) []step {
	return []step{{path: path, args: "rev-parse --path-format=absolute --show-toplevel", stdout: top + "\n"}}
}

func commonDir(path, stdout string, code int) []step {
	return []step{{path: path, args: "rev-parse --path-format=absolute --git-common-dir", stdout: stdout, code: code}}
}

func head(path, stdout string, code int) []step {
	return []step{{path: path, args: "symbolic-ref -q --short HEAD", stdout: stdout, code: code}}
}

func root(path, branch string) []step {
	return cat(inside(path, "true"), toplevel(path, path), head(path, branch+"\n", 0))
}

func initStep(path, branch string, code int) []step {
	return []step{{path: path, args: "init -q --initial-branch=" + branch + " --template=", code: code, stderr: "init diagnostic\n"}}
}

func tempBase(t *testing.T) string {
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return base
}

func TestEnsureInitStates(t *testing.T) {
	t.Setenv("GIT_DIR", "/outside/repo/.git")
	base := tempBase(t)
	empty, full, file, link := filepath.Join(base, "empty"), filepath.Join(base, "full"), filepath.Join(base, "file"), filepath.Join(base, "link")
	for _, err := range []error{os.Mkdir(empty, 0o755), os.Mkdir(full, 0o755), os.WriteFile(filepath.Join(full, "kept"), nil, 0o644),
		os.WriteFile(file, nil, 0o644), os.Symlink(empty, link)} {
		if err != nil {
			t.Fatal(err)
		}
	}
	created, nested := filepath.Join(base, "created"), filepath.Join(base, "nested")
	tests := []struct {
		name, dir, branch string
		steps             []step
		changed           bool
		want              error
	}{
		{"missing dir under absent parent", created, "main", cat(refOK(base, "main"), absent(base), initStep(created, "main", 0), root(created, "main")), true, nil},
		{"missing parent", filepath.Join(base, "none", "child"), "main", nil, false, ErrInvalid},
		{"missing dir under repo", nested, "main", cat(refOK(base, "main"), inside(base, "true"), toplevel(base, "/elsewhere")), false, ErrConflict},
		{"final symlink", link, "main", nil, false, ErrConflict},
		{"not a directory", file, "main", nil, false, ErrConflict},
		{"non-empty absent", full, "main", cat(refOK(full, "main"), absent(full)), false, ErrConflict},
		{"root same branch", full, "main", cat(refOK(full, "main"), root(full, "main")), false, nil},
		{"root other branch", full, "main", cat(refOK(full, "main"), root(full, "dev")), false, ErrConflict},
		{"root detached", full, "main", cat(refOK(full, "main"), inside(full, "true"), toplevel(full, full), head(full, "", 1)), false, ErrConflict},
		{"nested", full, "main", cat(refOK(full, "main"), inside(full, "true"), toplevel(full, base)), false, ErrConflict},
		{"not worktree", full, "main", cat(refOK(full, "main"), inside(full, "false")), false, ErrConflict},
		{"git missing", full, "main", cat(refOK(full, "main"), insideFails(full, "", 0, exec.ErrNotFound)), false, ErrUnknown},
		{"dubious ownership", full, "main", cat(refOK(full, "main"), insideFails(full, "fatal: detected dubious ownership in repository\n", 128, nil)), false, ErrUnknown},
		{"other exit 128", full, "main", cat(refOK(full, "main"), insideFails(full, "fatal: bad config line 1\n", 128, nil)), false, ErrUnknown},
		{"deadline", full, "main", cat(refOK(full, "main"), insideFails(full, "", 0, context.DeadlineExceeded)), false, ErrUnknown},
		{"unexpected stdout", full, "main", cat(refOK(full, "main"), inside(full, "maybe")), false, ErrUnknown},
		{"head failure", full, "main", cat(refOK(full, "main"), inside(full, "true"), toplevel(full, full), head(full, "", 128)), false, ErrUnknown},
		{"ref rejected", full, "bad..ref", []step{{path: full, args: "check-ref-format --branch bad..ref", code: 1}}, false, ErrInvalid},
		{"ref expanded", full, "@{-1}", []step{{path: full, args: "check-ref-format --branch @{-1}", stdout: "other\n"}}, false, ErrInvalid},
		{"ref check failure", full, "main", []step{{path: full, args: "check-ref-format --branch main", code: 128}}, false, ErrUnknown},
		{"init failure reads back once", empty, "main", cat(refOK(empty, "main"), absent(empty), initStep(empty, "main", 1), absent(empty)), false, ErrUnknown},
		{"readback mismatch", empty, "main", cat(refOK(empty, "main"), absent(empty), initStep(empty, "main", 0), root(empty, "dev")), false, ErrUnknown},
		{"empty absent", empty, "main", cat(refOK(empty, "main"), absent(empty), initStep(empty, "main", 0), root(empty, "main")), true, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, existedErr := os.Lstat(test.dir)
			runner := &fakeRunner{t: t, steps: test.steps}

			repo, changed, err := EnsureInit(context.Background(), runner, test.dir, test.branch)

			if !errors.Is(err, test.want) || (test.want == nil) != (err == nil) || changed != test.changed {
				t.Fatalf("EnsureInit = %v, %v; want changed %v, error %v", changed, err, test.changed, test.want)
			}
			if len(runner.calls) != len(test.steps) {
				t.Fatalf("made %d calls, want %d", len(runner.calls), len(test.steps))
			}
			if err == nil && (repo.State != Root || repo.Toplevel != test.dir || repo.Branch != test.branch) {
				t.Fatalf("repo = %#v", repo)
			}
			if _, statErr := os.Lstat(test.dir); err != nil && os.IsNotExist(existedErr) && !os.IsNotExist(statErr) {
				t.Fatalf("failed EnsureInit created %s", test.dir)
			}
		})
	}
}

func TestInvalidInputsMakeNoCalls(t *testing.T) {
	dir := tempBase(t)
	runner := &fakeRunner{t: t}
	var nilContext context.Context
	for name, call := range map[string]func() error{
		"nil context":       func() error { _, _, err := EnsureInit(nilContext, runner, dir, "main"); return err },
		"nil runner":        func() error { _, _, err := EnsureInit(context.Background(), nil, dir, "main"); return err },
		"relative dir":      func() error { _, _, err := EnsureInit(context.Background(), runner, "relative", "main"); return err },
		"empty branch":      func() error { _, _, err := EnsureInit(context.Background(), runner, dir, ""); return err },
		"option branch":     func() error { _, _, err := EnsureInit(context.Background(), runner, dir, "-b"); return err },
		"observe nil ctx":   func() error { _, err := Observe(nilContext, runner, dir); return err },
		"observe relative":  func() error { _, err := Observe(context.Background(), runner, "relative"); return err },
		"identify nil run":  func() error { _, err := Identify(context.Background(), nil, dir); return err },
		"identify relative": func() error { _, err := Identify(context.Background(), runner, "relative"); return err },
	} {
		if err := call(); !errors.Is(err, ErrInvalid) {
			t.Errorf("%s returned %v, want ErrInvalid", name, err)
		}
	}
	if len(runner.calls) != 0 {
		t.Fatalf("invalid inputs made %d calls", len(runner.calls))
	}
}

func TestInitFailureAfterDeadlineReadsBackWithFreshBoundedContext(t *testing.T) {
	dir := tempBase(t)
	ctx, cancel := context.WithCancel(context.Background())
	runner := &fakeRunner{t: t, steps: cat(refOK(dir, "main"), absent(dir), initStep(dir, "main", 1), root(dir, "main"))}
	runner.steps[2].code, runner.steps[2].err = 0, context.Canceled
	runner.after = func(c process.Command) {
		if c.Args[2] == "init" {
			cancel()
		}
	}

	_, _, err := EnsureInit(ctx, runner, dir, "main")

	if !errors.Is(err, ErrUnknown) || !strings.Contains(err.Error(), "root") {
		t.Fatalf("EnsureInit returned %v, want ErrUnknown with observed state", err)
	}
	readback := runner.calls[3].Context
	if _, ok := readback.Deadline(); readback == ctx || !ok {
		t.Fatal("readback after cancellation did not use a fresh bounded context")
	}
}

func TestObserveUnbornRoot(t *testing.T) {
	dir := tempBase(t)
	runner := &fakeRunner{t: t, steps: root(dir, "main")}

	repo, err := Observe(context.Background(), runner, dir)

	if err != nil || !reflect.DeepEqual(repo, Repo{State: Root, Dir: dir, Toplevel: dir, Branch: "main"}) {
		t.Fatalf("Observe = %#v, %v", repo, err)
	}
}

func TestGitEnvFiltersInheritedGitVariables(t *testing.T) {
	got := gitEnv([]string{"PATH=/bin", "GIT_DIR=/x", "GIT_WORK_TREE=/y", "GIT_CONFIG_COUNT=1", "GIT_CONFIG_PARAMETERS='a=b'",
		"GIT_CONFIG_GLOBAL=/g", "GIT_CONFIG_SYSTEM=/s", "GIT_CONFIG_NOSYSTEM=1", "HOME=/h"})

	want := []string{"PATH=/bin", "GIT_CONFIG_GLOBAL=/g", "GIT_CONFIG_SYSTEM=/s", "GIT_CONFIG_NOSYSTEM=1", "HOME=/h", "LC_ALL=C", "GIT_TERMINAL_PROMPT=0"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("gitEnv = %v, want %v", got, want)
	}
}

func TestIdentifyReturnsRawIdentity(t *testing.T) {
	t.Setenv("GIT_DIR", "/outside/repo/.git")
	dir := tempBase(t)
	topFails := []step{{path: dir, args: "rev-parse --path-format=absolute --show-toplevel", code: 128}}
	tests := []struct {
		name  string
		steps []step
		want  Identity
		op    string
	}{
		{"root", cat(inside(dir, "true"), toplevel(dir, dir), commonDir(dir, dir+"/.git\n", 0)), Identity{true, dir, dir + "/.git"}, ""},
		{"raw other toplevel", cat(inside(dir, "true"), toplevel(dir, "/elsewhere"), commonDir(dir, "/elsewhere/.git\n", 0)), Identity{true, "/elsewhere", "/elsewhere/.git"}, ""},
		{"not worktree", inside(dir, "false"), Identity{}, ""},
		{"inside exit", insideFails(dir, "fatal: detected dubious ownership in repository\n", 128, nil), Identity{}, OpInsideWorkTree},
		{"inside launch", insideFails(dir, "", 0, exec.ErrNotFound), Identity{}, OpInsideWorkTree},
		{"inside unexpected", inside(dir, "maybe"), Identity{}, OpInsideWorkTree},
		{"inside two lines", inside(dir, "true\nextra"), Identity{}, OpInsideWorkTree},
		{"toplevel exit", cat(inside(dir, "true"), topFails), Identity{Worktree: true}, OpToplevel},
		{"toplevel empty", cat(inside(dir, "true"), toplevel(dir, "")), Identity{Worktree: true}, OpToplevel},
		{"common-dir exit keeps partial", cat(inside(dir, "true"), toplevel(dir, "/elsewhere"), commonDir(dir, "", 128)), Identity{true, "/elsewhere", ""}, OpCommonDir},
		{"common-dir unterminated", cat(inside(dir, "true"), toplevel(dir, dir), commonDir(dir, dir+"/.git", 0)), Identity{true, dir, ""}, OpCommonDir},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := &fakeRunner{t: t, steps: test.steps}

			identity, err := Identify(context.Background(), runner, dir)

			if identity != test.want || len(runner.calls) != len(test.steps) || (test.op == "") != (err == nil) {
				t.Fatalf("Identify = %#v, %v after %d calls", identity, err, len(runner.calls))
			}
			if test.op == "" {
				return
			}
			last := test.steps[len(test.steps)-1]
			var typed *Error
			var exitErr *exec.ExitError
			if !errors.As(err, &typed) || typed.Op != test.op || !errors.Is(err, ErrUnknown) ||
				!strings.HasPrefix(err.Error(), "unknown: "+test.op+" "+dir+": ") ||
				errors.As(err, &exitErr) != (last.code != 0) || (typed.Err == nil) != (last.code == 0 && last.err == nil) {
				t.Fatalf("Identify error = %#v (%v)", err, err)
			}
		})
	}
}

func TestSingleLine(t *testing.T) {
	for _, test := range []struct {
		output, want string
		ok           bool
	}{
		{"true\n", "true", true}, {"true\r\n", "true", true}, {"\n", "", true},
		{"true", "", false}, {"", "", false}, {"true\nextra\n", "", false}, {"a\rb\n", "", false},
	} {
		if got, ok := SingleLine([]byte(test.output)); got != test.want || ok != test.ok {
			t.Errorf("SingleLine(%q) = %q, %v", test.output, got, ok)
		}
	}
}

const (
	wantURL = "https://x-token:SECRET@example.invalid/a.git"
	leak    = "fatal: unable to access 'https://x:SECRET@example.invalid/b.git'\n"
)

func rec(scope, key, value string) string { return scope + "\x00" + key + "\n" + value + "\x00" }

func readConfig(path, name, stdout string, code int) []step {
	return []step{{path: path, args: `config --show-scope --null --get-regexp ^remote\.` + strings.ReplaceAll(name, ".", `\.`) + `\.`,
		stdout: stdout, stderr: map[bool]string{true: leak}[code > 1], code: code}}
}

func getURL(path, name, stdout string, code int) []step {
	stderr := map[int]string{0: "", 2: "error: No such remote '" + name + "'\n"}[code]
	if code != 0 && code != 2 {
		stderr = leak
	}
	return []step{{path: path, args: "remote get-url " + name, stdout: stdout, stderr: stderr, code: code}}
}

func addRemote(path, name string, code int) []step {
	return []step{{path: path, args: "remote add -- " + name + " " + wantURL, stderr: map[bool]string{true: leak}[code != 0], code: code}}
}

func remoteAbsent(path, name string) []step {
	return cat(readConfig(path, name, "", 1), getURL(path, name, "", 2))
}

func remoteSet(path, name string) []step {
	return cat(readConfig(path, name, rec("local", "remote."+name+".url", wantURL), 0), getURL(path, name, wantURL+"\n", 0))
}

func TestEnsureRemote(t *testing.T) {
	dir := tempBase(t)
	origin := rec("local", "remote.origin.url", wantURL)
	other := rec("local", "remote.origin.url", "https://x:OLDSECRET@example.invalid/old.git")
	conflict := func(records string) []step {
		return cat(root(dir, "main"), readConfig(dir, "origin", records, 0), getURL(dir, "origin", wantURL+"\n", 0))
	}
	malformed := func(records string) []step { return cat(root(dir, "main"), readConfig(dir, "origin", records, 0)) }
	missing := func(stdout, stderr string) []step {
		return cat(root(dir, "main"), readConfig(dir, "origin", "", 1), []step{{path: dir, args: "remote get-url origin", stdout: stdout, stderr: stderr, code: 2}})
	}
	tests := []struct {
		name, remote string
		steps        []step
		changed      bool
		want         error
	}{
		{"absent adds and reads back", "origin", cat(root(dir, "main"), remoteAbsent(dir, "origin"), addRemote(dir, "origin", 0), remoteSet(dir, "origin")), true, nil},
		{"same url keeps other metadata", "origin", conflict(origin + rec("local", "remote.origin.prune", "true") + rec("local", "remote.origin.tagopt", "--no-tags")), false, nil},
		{"only other subsection is absent", "origin", cat(root(dir, "main"), readConfig(dir, "origin", rec("local", "remote.origin.x.url", wantURL), 0), getURL(dir, "origin", "", 2),
			addRemote(dir, "origin", 0), remoteSet(dir, "origin")), true, nil},
		{"dotted name is quoted", "up.stream", cat(root(dir, "main"), remoteAbsent(dir, "up.stream"), addRemote(dir, "up.stream", 0), remoteSet(dir, "up.stream")), true, nil},
		{"two urls", "origin", conflict(origin + origin), false, ErrConflict},
		{"pushurl", "origin", conflict(origin + rec("local", "remote.origin.pushurl", wantURL)), false, ErrConflict},
		{"global url", "origin", conflict(rec("global", "remote.origin.url", wantURL)), false, ErrConflict},
		{"worktree scope", "origin", conflict(origin + rec("worktree", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*")), false, ErrConflict},
		{"raw differs", "origin", cat(root(dir, "main"), readConfig(dir, "origin", other, 0), getURL(dir, "origin", "https://x:OLDSECRET@example.invalid/old.git\n", 0)), false, ErrConflict},
		{"effective rewritten", "origin", cat(root(dir, "main"), readConfig(dir, "origin", origin, 0), getURL(dir, "origin", "https://x:NEWSECRET@rewritten.invalid/a.git\n", 0)), false, ErrConflict},
		{"legacy remote", "origin", cat(root(dir, "main"), readConfig(dir, "origin", "", 1), getURL(dir, "origin", wantURL+"\n", 0)), false, ErrConflict},
		{"fetch only", "origin", cat(root(dir, "main"), readConfig(dir, "origin", rec("local", "remote.origin.fetch", "x"), 0), getURL(dir, "origin", "", 2)), false, ErrConflict},
		{"record without value", "origin", malformed("local\x00remote.origin.url\x00"), false, ErrUnknown},
		{"unterminated record", "origin", malformed(strings.TrimSuffix(origin, "\x00")), false, ErrUnknown},
		{"unknown scope", "origin", malformed(rec("blob", "remote.origin.url", wantURL)), false, ErrUnknown},
		{"malformed other record", "origin", malformed(origin + "local\x00remote.origin.x.url\x00"), false, ErrUnknown},
		{"config rc 1 with output", "origin", cat(root(dir, "main"), readConfig(dir, "origin", origin, 1)), false, ErrUnknown},
		{"config rc 128", "origin", cat(root(dir, "main"), readConfig(dir, "origin", "", 128)), false, ErrUnknown},
		{"get-url two lines", "origin", cat(root(dir, "main"), readConfig(dir, "origin", origin, 0), getURL(dir, "origin", wantURL+"\n"+wantURL+"\n", 0)), false, ErrUnknown},
		{"get-url empty", "origin", cat(root(dir, "main"), readConfig(dir, "origin", origin, 0), getURL(dir, "origin", "\n", 0)), false, ErrUnknown},
		{"missing with stdout", "origin", missing(wantURL+"\n", "error: No such remote 'origin'\n"), false, ErrUnknown},
		{"missing other name", "origin", missing("", "error: No such remote 'origin2'\n"), false, ErrUnknown},
		{"missing extra diagnostic", "origin", missing("", "error: No such remote 'origin'\nwarning: "+leak), false, ErrUnknown},
		{"get-url rc 128", "origin", cat(root(dir, "main"), readConfig(dir, "origin", origin, 0), getURL(dir, "origin", "", 128)), false, ErrUnknown},
		{"add fails reads back once", "origin", cat(root(dir, "main"), remoteAbsent(dir, "origin"), addRemote(dir, "origin", 3), remoteSet(dir, "origin")), false, ErrUnknown},
		{"add ok readback drift", "origin", cat(root(dir, "main"), remoteAbsent(dir, "origin"), addRemote(dir, "origin", 0),
			readConfig(dir, "origin", other, 0), getURL(dir, "origin", "https://x:OLDSECRET@example.invalid/old.git\n", 0)), false, ErrUnknown},
		{"add ok readback absent", "origin", cat(root(dir, "main"), remoteAbsent(dir, "origin"), addRemote(dir, "origin", 0), remoteAbsent(dir, "origin")), false, ErrUnknown},
		{"nested", "origin", cat(inside(dir, "true"), toplevel(dir, "/elsewhere")), false, ErrConflict},
		{"absent repository", "origin", absent(dir), false, ErrConflict},
		{"not worktree", "origin", inside(dir, "false"), false, ErrConflict},
		{"observe failure", "origin", insideFails(dir, leak, 128, nil), false, ErrUnknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := &fakeRunner{t: t, steps: test.steps}

			changed, err := EnsureRemote(context.Background(), runner, dir, test.remote, wantURL)

			if !errors.Is(err, test.want) || (test.want == nil) != (err == nil) || changed != test.changed || len(runner.calls) != len(test.steps) {
				t.Fatalf("EnsureRemote = %v, %v after %d calls; want %v, %v", changed, err, len(runner.calls), test.changed, test.want)
			}
			if err != nil && (strings.Contains(err.Error(), "SECRET") || strings.Contains(err.Error(), "example.invalid") || !strings.Contains(err.Error(), dir)) {
				t.Fatalf("error %q leaks a URL or lacks the directory", err)
			}
		})
	}
}

func TestEnsureRemoteRejectsInvalidInputsWithoutCalls(t *testing.T) {
	base := tempBase(t)
	file, link := filepath.Join(base, "file"), filepath.Join(base, "link")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(base, link); err != nil {
		t.Fatal(err)
	}
	type input struct {
		ctx            context.Context
		runner         process.Runner
		dir, name, url string
		want           error
	}
	var nilContext context.Context
	tests := []input{
		{nilContext, &fakeRunner{t: t}, base, "origin", wantURL, ErrInvalid},
		{context.Background(), nil, base, "origin", wantURL, ErrInvalid},
		{context.Background(), &fakeRunner{t: t}, "relative", "origin", wantURL, ErrInvalid},
		{context.Background(), &fakeRunner{t: t}, filepath.Join(base, "missing"), "origin", wantURL, ErrInvalid},
		{context.Background(), &fakeRunner{t: t}, file, "origin", wantURL, ErrConflict},
		{context.Background(), &fakeRunner{t: t}, link, "origin", wantURL, ErrConflict},
	}
	for _, name := range []string{"", "-x", ".x", "a..b", "x.lock", "x.", "a b", "a/b", strings.Repeat("a", 65)} {
		tests = append(tests, input{context.Background(), &fakeRunner{t: t}, base, name, wantURL, ErrInvalid})
	}
	for _, url := range []string{"", "-u", "--upload-pack=x", "a b", "a\tb", "a\nb", "a\x7fb"} {
		tests = append(tests, input{context.Background(), &fakeRunner{t: t}, base, "origin", url, ErrInvalid})
	}
	for _, test := range tests {
		changed, err := EnsureRemote(test.ctx, test.runner, test.dir, test.name, test.url)

		if changed || !errors.Is(err, test.want) || (test.url != "" && strings.Contains(err.Error(), test.url)) {
			t.Errorf("EnsureRemote(%q, %q, %q) = %v, %v; want %v", test.dir, test.name, test.url, changed, err, test.want)
		}
		if runner, ok := test.runner.(*fakeRunner); ok && len(runner.calls) != 0 {
			t.Errorf("EnsureRemote(%q, %q) made %d calls", test.dir, test.name, len(runner.calls))
		}
	}
}

func TestEnsureRemoteReadsBackWithFreshBoundedContextAfterCancellation(t *testing.T) {
	dir := tempBase(t)
	ctx, cancel := context.WithCancel(context.Background())
	runner := &fakeRunner{t: t, steps: cat(root(dir, "main"), remoteAbsent(dir, "origin"), addRemote(dir, "origin", 0), remoteSet(dir, "origin"))}
	runner.after = func(c process.Command) {
		if c.Args[2] == "remote" && c.Args[3] == "add" {
			cancel()
		}
	}

	changed, err := EnsureRemote(ctx, runner, dir, "origin", wantURL)

	if err != nil || !changed {
		t.Fatalf("EnsureRemote = %v, %v; want added", changed, err)
	}
	for _, readback := range runner.calls[len(runner.calls)-2:] {
		if _, ok := readback.Context.Deadline(); readback.Context == ctx || !ok {
			t.Fatal("readback after cancellation did not use a fresh bounded context")
		}
	}
}

func TestEnsureRemoteReportsCanceledWithoutGitOutput(t *testing.T) {
	dir := tempBase(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	runner := &fakeRunner{t: t, steps: insideFails(dir, leak, 0, context.Canceled)}

	_, err := EnsureRemote(ctx, runner, dir, "origin", wantURL)

	if want := "unknown: observe " + dir + ": canceled"; !errors.Is(err, ErrUnknown) || err.Error() != want || len(runner.calls) != 1 {
		t.Fatalf("EnsureRemote = %v after %d calls, want %q", err, len(runner.calls), want)
	}
}

func TestValidRemoteNameMatchesEnsureRemoteNameRule(t *testing.T) {
	for _, name := range []string{"origin", "up-stream", "a.b_c"} {
		if !ValidRemoteName(name) {
			t.Fatalf("%q rejected", name)
		}
	}
	for _, name := range []string{"", "-x", ".x", "a..b", "x.lock", "x.", "a b", "a/b", strings.Repeat("a", 65)} {
		if ValidRemoteName(name) {
			t.Fatalf("%q accepted", name)
		}
	}
}

func TestReadDisablesOptionalLocksAndKeepsRawOutput(t *testing.T) {
	t.Setenv("GIT_DIR", "/elsewhere/.git")
	t.Setenv("GIT_OPTIONAL_LOCKS", "1")
	runner := &fakeRunner{t: t, steps: []step{{path: "/repo", args: "-c core.fsmonitor=false status --porcelain=v2 -z", stdout: "# branch.head main\x00", stderr: "notice: x\n", code: 3}}}

	stdout, stderr, code, err := Read(context.Background(), runner, "/repo", "status", "--porcelain=v2", "-z")

	if string(stdout) != "# branch.head main\x00" || stderr != "notice: x\n" || code != 3 || err == nil {
		t.Fatalf("Read = %q, %q, %d, %v", stdout, stderr, code, err)
	}
	env := "\n" + strings.Join(runner.calls[0].Env, "\n") + "\n"
	if !strings.Contains(env, "\nGIT_OPTIONAL_LOCKS=0\n") || strings.Contains(env, "\nGIT_OPTIONAL_LOCKS=1\n") {
		var git []string
		for _, entry := range runner.calls[0].Env {
			if strings.HasPrefix(entry, "GIT_") {
				git = append(git, entry)
			}
		}
		t.Fatalf("GIT_ env %v", git)
	}
}
