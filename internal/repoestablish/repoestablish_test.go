package repoestablish

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/gitrepo"
)

// world is a stateful fake of every adapter; the wrapper is a real directory.
type world struct {
	t         *testing.T
	workspace string
	setup     ProjectSetup
	calls     []string
	resolved  []Resolved
	repos     map[string]gitrepo.Repo
	created   bool
	remotes   map[string]string
	defined   bool
	wrapper   string
	fail      map[string]error
	before    map[string]func()
}

func newWorld(t *testing.T) *world {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	projectRoot := filepath.Join(root, "projects")
	if err := os.Mkdir(projectRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	return &world{t: t, workspace: root, setup: ProjectSetup{Root: projectRoot, GitHubOwner: "Org"},
		repos: map[string]gitrepo.Repo{}, remotes: map[string]string{}, fail: map[string]error{}, before: map[string]func(){}}
}

func (w *world) request() Request {
	return Request{Workspace: w.workspace, Project: "ply", RepoName: "code-map", Actor: "ply-epic"}
}

func (w *world) enter(name string) error {
	w.calls = append(w.calls, name)
	if hook := w.before[name]; hook != nil {
		hook()
	}
	return w.fail[name]
}

func (w *world) deps() Dependencies {
	return Dependencies{
		Prepare: func(ctx context.Context, req Request) (ProjectSetup, error) {
			if err := w.enter("prepare"); err != nil {
				return ProjectSetup{}, err
			}
			return w.setup, nil
		},
		Check: func(ctx context.Context, r Resolved) error {
			w.resolved = append(w.resolved, r)
			return w.enter("check")
		},
		Observe: func(ctx context.Context, dir string) (gitrepo.Repo, error) {
			if err := w.enter("observe " + dir); err != nil {
				return gitrepo.Repo{}, err
			}
			return w.repos[dir], nil
		},
		Init: func(ctx context.Context, dir, branch string) (bool, error) {
			if err := w.enter("init " + dir + " " + branch); err != nil {
				return false, err
			}
			if w.repos[dir].State == gitrepo.Root {
				return false, nil
			}
			w.repos[dir] = gitrepo.Repo{State: gitrepo.Root, Dir: dir, Toplevel: dir, Branch: branch}
			return true, nil
		},
		GitHub: func(ctx context.Context, owner, name string) (string, bool, error) {
			if err := w.enter("github " + owner + "/" + name); err != nil {
				return "", false, err
			}
			changed := !w.created
			w.created = true
			return "git@github.com:" + owner + "/" + name + ".git", changed, nil
		},
		Remote: func(ctx context.Context, dir, name, url string) (bool, error) {
			if err := w.enter("remote " + dir + " " + name + " " + url); err != nil {
				return false, err
			}
			changed := w.remotes[name] == ""
			w.remotes[name] = url
			return changed, nil
		},
		Define: func(ctx context.Context, r Resolved) (bool, error) {
			w.resolved = append(w.resolved, r)
			if err := w.enter("define"); err != nil {
				return false, err
			}
			changed := !w.defined
			w.defined = true
			return changed, nil
		},
		SetWrapper: func(ctx context.Context, r Resolved) (bool, error) {
			w.resolved = append(w.resolved, r)
			if err := w.enter("metadata"); err != nil {
				return false, err
			}
			changed := w.wrapper == ""
			w.wrapper = r.Wrapper
			return changed, nil
		},
	}
}

func allSteps(changed bool) []StepResult {
	steps := []StepResult{}
	for _, step := range []Step{StepWrapper, StepInit, StepGitHub, StepRemote, StepDefine, StepMetadata} {
		steps = append(steps, StepResult{step, changed})
	}
	return steps
}

func requireStop(t *testing.T, err error, step Step, kind error) {
	t.Helper()
	var stopped *Error
	if !errors.As(err, &stopped) || stopped.Step != step || !errors.Is(err, kind) {
		t.Fatalf("error %v, want %v at %s", err, kind, step)
	}
}

func requireNoPath(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("%s exists: %v", path, err)
	}
}

func TestRunResolvesOnceAndRunsStepsInOrderWithURLFlow(t *testing.T) {
	w := newWorld(t)
	result, err := Run(context.Background(), w.request(), w.deps())
	if err != nil {
		t.Fatal(err)
	}
	wrapper := filepath.Join(w.setup.Root, "code-map")
	repoDir := filepath.Join(wrapper, "main")
	url := "git@github.com:Org/code-map.git"
	want := Resolved{Workspace: w.workspace, Project: "ply", ProjectRoot: w.setup.Root, RepoName: "code-map", Wrapper: wrapper, RepoDir: repoDir,
		Leaf: "main", Branch: "main", Remote: "origin", GitHubOwner: "Org", Actor: "ply-epic"}
	if result.Resolved != want || result.SSHURL != url || !reflect.DeepEqual(result.Steps, allSteps(true)) {
		t.Fatalf("result %+v", result)
	}
	for _, r := range w.resolved {
		if r != want {
			t.Fatalf("adapter got %+v", r)
		}
	}
	wantCalls := []string{"prepare", "check", "observe " + wrapper, "init " + repoDir + " main", "github Org/code-map",
		"remote " + repoDir + " origin " + url, "define", "metadata"}
	if !reflect.DeepEqual(w.calls, wantCalls) {
		t.Fatalf("calls\n%q\nwant\n%q", w.calls, wantCalls)
	}
	if info, err := os.Lstat(wrapper); err != nil || !info.IsDir() {
		t.Fatalf("wrapper %v %v", info, err)
	}
}

func TestRunPrecedenceOverrideConfigConvention(t *testing.T) {
	for name, tc := range map[string]struct {
		overrides Overrides
		setup     ProjectSetup
		want      [4]string // owner, leaf, branch, remote
	}{
		"convention": {Overrides{}, ProjectSetup{GitHubOwner: "Cfg"}, [4]string{"Cfg", "main", "main", "origin"}},
		"config":     {Overrides{}, ProjectSetup{GitHubOwner: "Cfg", Leaf: "trunk", Branch: "dev", Remote: "up"}, [4]string{"Cfg", "trunk", "dev", "up"}},
		"override": {Overrides{GitHubOwner: "Flag", Leaf: "l", Branch: "b", Remote: "r"}, ProjectSetup{GitHubOwner: "Cfg", Leaf: "trunk", Branch: "dev", Remote: "up"},
			[4]string{"Flag", "l", "b", "r"}},
		"owner override only": {Overrides{GitHubOwner: "Flag"}, ProjectSetup{}, [4]string{"Flag", "main", "main", "origin"}},
	} {
		t.Run(name, func(t *testing.T) {
			w := newWorld(t)
			tc.setup.Root = w.setup.Root
			w.setup = tc.setup
			req := w.request()
			req.Overrides = tc.overrides
			result, err := Run(context.Background(), req, w.deps())
			r := result.Resolved
			if err != nil || [4]string{r.GitHubOwner, r.Leaf, r.Branch, r.Remote} != tc.want || r.RepoDir != filepath.Join(r.Wrapper, tc.want[1]) {
				t.Fatalf("resolved %+v, error %v", r, err)
			}
			if req.Overrides != tc.overrides {
				t.Fatalf("request changed: %+v", req)
			}
		})
	}
	w := newWorld(t)
	w.setup.GitHubOwner = ""
	result, err := Run(context.Background(), w.request(), w.deps())
	requireStop(t, err, StepPrepare, ErrInvalid)
	if !reflect.DeepEqual(w.calls, []string{"prepare"}) || result.Steps != nil || result.Resolved != (Resolved{}) {
		t.Fatalf("calls %q, result %+v", w.calls, result)
	}
	requireNoPath(t, filepath.Join(w.setup.Root, "code-map"))
}

func TestRunRejectsStaticRequestWithoutPrepare(t *testing.T) {
	w := newWorld(t)
	for name, edit := range map[string]func(*Request){
		"empty actor":        func(r *Request) { r.Actor = "" },
		"empty project":      func(r *Request) { r.Project = "" },
		"relative workspace": func(r *Request) { r.Workspace = "ws" },
		"unclean workspace":  func(r *Request) { r.Workspace += "/." },
		"empty repo":         func(r *Request) { r.RepoName = "" },
		"dot dot":            func(r *Request) { r.RepoName = ".." },
		"dot":                func(r *Request) { r.RepoName = "." },
		"separator":          func(r *Request) { r.RepoName = "a/b" },
		"absolute":           func(r *Request) { r.RepoName = "/abs" },
	} {
		req := w.request()
		edit(&req)
		_, err := Run(context.Background(), req, w.deps())
		if requireStop(t, err, StepPrepare, ErrInvalid); len(w.calls) != 0 {
			t.Fatalf("%s: calls %q", name, w.calls)
		}
	}
	for name, overrides := range map[string]Overrides{
		"leaf dot dot": {Leaf: ".."}, "leaf separator": {Leaf: "a/b"}, "branch dash": {Branch: "-x"}, "remote dash": {Remote: "-x"},
	} {
		req := w.request()
		req.Overrides = overrides
		w.calls = nil
		_, err := Run(context.Background(), req, w.deps())
		if requireStop(t, err, StepPrepare, ErrInvalid); !reflect.DeepEqual(w.calls, []string{"prepare"}) {
			t.Fatalf("%s: calls %q", name, w.calls)
		}
	}
	requireNoPath(t, filepath.Join(w.setup.Root, "code-map"))
}

func TestRunFullRerunIsUnchangedToTheEnd(t *testing.T) {
	w := newWorld(t)
	if _, err := Run(context.Background(), w.request(), w.deps()); err != nil {
		t.Fatal(err)
	}
	w.calls = nil
	result, err := Run(context.Background(), w.request(), w.deps())
	if err != nil || !reflect.DeepEqual(result.Steps, allSteps(false)) {
		t.Fatalf("rerun %+v, error %v", result, err)
	}
	// The wrapper now exists, so nesting is observed at the repository directory.
	if w.calls[2] != "observe "+result.Resolved.RepoDir || w.calls[len(w.calls)-1] != "metadata" {
		t.Fatalf("calls %q", w.calls)
	}
}

func TestRunWithMissingDependencyMakesNoCallsAndNoDirectory(t *testing.T) {
	for _, field := range []string{"Prepare", "Check", "Observe", "Init", "GitHub", "Remote", "Define", "SetWrapper"} {
		w := newWorld(t)
		d := w.deps()
		reflect.ValueOf(&d).Elem().FieldByName(field).Set(reflect.Zero(reflect.ValueOf(d).FieldByName(field).Type()))
		_, err := Run(context.Background(), w.request(), d)
		if requireStop(t, err, StepPrepare, ErrInvalid); len(w.calls) != 0 {
			t.Fatalf("%s: calls %q", field, w.calls)
		}
		requireNoPath(t, filepath.Join(w.setup.Root, "code-map"))
	}
	w := newWorld(t)
	var noContext context.Context
	if _, err := Run(noContext, w.request(), w.deps()); !errors.Is(err, ErrInvalid) || len(w.calls) != 0 {
		t.Fatalf("nil context: %v, calls %q", err, w.calls)
	}
}

func TestRunPrecheckStopsWithoutEffect(t *testing.T) {
	for name, tc := range map[string]struct {
		arrange func(*world)
		step    Step
		kind    error
	}{
		"symlink wrapper": {func(w *world) {
			if err := os.Symlink(w.workspace, filepath.Join(w.setup.Root, "code-map")); err != nil {
				w.t.Fatal(err)
			}
		}, StepCheck, ErrConflict},
		"file wrapper": {func(w *world) {
			if err := os.WriteFile(filepath.Join(w.setup.Root, "code-map"), nil, 0o644); err != nil {
				w.t.Fatal(err)
			}
		}, StepCheck, ErrConflict},
		"missing project root": {func(w *world) { w.setup.Root = filepath.Join(w.workspace, "gone") }, StepCheck, ErrInvalid},
		"project root symlink": {func(w *world) {
			link := filepath.Join(w.workspace, "alias")
			if err := os.Symlink(w.setup.Root, link); err != nil {
				w.t.Fatal(err)
			}
			w.setup.Root = link
		}, StepCheck, ErrInvalid},
		"wrapper outside workspace": {func(w *world) { w.setup.Root = filepath.Dir(w.workspace) }, StepPrepare, ErrInvalid},
		"relative project root":     {func(w *world) { w.setup.Root = "projects" }, StepPrepare, ErrInvalid},
		"nested wrapper": {func(w *world) {
			w.repos[filepath.Join(w.setup.Root, "code-map")] = gitrepo.Repo{State: gitrepo.Nested}
		}, StepCheck, ErrConflict},
		"nested repository": {func(w *world) {
			wrapper := filepath.Join(w.setup.Root, "code-map")
			if err := os.Mkdir(wrapper, 0o755); err != nil {
				w.t.Fatal(err)
			}
			w.repos[filepath.Join(wrapper, "main")] = gitrepo.Repo{State: gitrepo.Nested}
		}, StepCheck, ErrConflict},
		"repository on other branch": {func(w *world) {
			wrapper := filepath.Join(w.setup.Root, "code-map")
			if err := os.Mkdir(wrapper, 0o755); err != nil {
				w.t.Fatal(err)
			}
			w.repos[filepath.Join(wrapper, "main")] = gitrepo.Repo{State: gitrepo.Root, Branch: "dev"}
		}, StepCheck, ErrConflict},
		"check conflict":     {func(w *world) { w.fail["check"] = fmt.Errorf("%w: metadata differs", ErrConflict) }, StepCheck, ErrConflict},
		"check unclassified": {func(w *world) { w.fail["check"] = errors.New("boom") }, StepCheck, ErrInvalid},
		"check unknown":      {func(w *world) { w.fail["check"] = fmt.Errorf("%w: probe", ErrUnknown) }, StepCheck, ErrInvalid},
		"prepare conflict":   {func(w *world) { w.fail["prepare"] = fmt.Errorf("%w: x", ErrConflict) }, StepPrepare, ErrConflict},
		"prepare other":      {func(w *world) { w.fail["prepare"] = errors.New("boom") }, StepPrepare, ErrInvalid},
		"observe failure": {func(w *world) {
			w.fail["observe "+filepath.Join(w.setup.Root, "code-map")] = fmt.Errorf("%w: git", ErrUnknown)
		}, StepCheck, ErrInvalid},
	} {
		t.Run(name, func(t *testing.T) {
			w := newWorld(t)
			tc.arrange(w)
			result, err := Run(context.Background(), w.request(), w.deps())
			requireStop(t, err, tc.step, tc.kind)
			if errors.Is(err, ErrUnknown) || result.Steps != nil {
				t.Fatalf("error %v, steps %+v", err, result.Steps)
			}
			for _, call := range w.calls {
				if call != "prepare" && call != "check" && !strings.HasPrefix(call, "observe ") {
					t.Fatalf("effect call %q", call)
				}
			}
		})
	}
}

func TestRunStopsAtEachEffectStepWithConfirmedStepsOnly(t *testing.T) {
	names := []string{"init", "github", "remote", "define", "metadata"}
	for index, name := range names {
		for _, kind := range []error{ErrConflict, ErrInvalid, ErrUnknown} {
			t.Run(fmt.Sprintf("%s %v", name, kind), func(t *testing.T) {
				w := newWorld(t)
				d := failAt(w.deps(), name, fmt.Errorf("%w: injected", kind))
				result, err := Run(context.Background(), w.request(), d)
				step := []Step{StepInit, StepGitHub, StepRemote, StepDefine, StepMetadata}[index]
				requireStop(t, err, step, kind)
				if want := allSteps(true)[:index+1]; !reflect.DeepEqual(result.Steps, want) {
					t.Fatalf("steps %+v, want %+v", result.Steps, want)
				}
				if (result.SSHURL != "") != (index > 1) {
					t.Fatalf("SSH URL %q at %s", result.SSHURL, name)
				}
			})
		}
	}
	t.Run("unclassified effect error is unknown", func(t *testing.T) {
		w := newWorld(t)
		d := failAt(w.deps(), "define", context.Canceled)
		_, err := Run(context.Background(), w.request(), d)
		requireStop(t, err, StepDefine, ErrUnknown)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cause lost: %v", err)
		}
	})
	t.Run("unknown in first effect step leaves no steps", func(t *testing.T) {
		w := newWorld(t)
		w.before["check"] = func() {
			if err := os.Chmod(w.setup.Root, 0o555); err != nil {
				t.Fatal(err)
			}
		}
		t.Cleanup(func() { _ = os.Chmod(w.setup.Root, 0o755) })
		result, err := Run(context.Background(), w.request(), w.deps())
		requireStop(t, err, StepWrapper, ErrUnknown)
		if len(result.Steps) != 0 || result.Resolved.Wrapper == "" {
			t.Fatalf("result %+v", result)
		}
	})
}

// failAt replaces one effect function by a call that fails after nothing.
func failAt(d Dependencies, name string, err error) Dependencies {
	switch name {
	case "init":
		d.Init = func(context.Context, string, string) (bool, error) { return false, err }
	case "github":
		d.GitHub = func(context.Context, string, string) (string, bool, error) { return "", false, err }
	case "remote":
		d.Remote = func(context.Context, string, string, string) (bool, error) { return false, err }
	case "define":
		d.Define = func(context.Context, Resolved) (bool, error) { return false, err }
	case "metadata":
		d.SetWrapper = func(context.Context, Resolved) (bool, error) { return false, err }
	}
	return d
}

func TestRunRerunAfterEachBoundaryResumes(t *testing.T) {
	names := []string{"wrapper", "init", "github", "remote", "define", "metadata"}
	for k := 1; k <= len(names); k++ {
		t.Run(names[k-1], func(t *testing.T) {
			w := newWorld(t)
			d := w.deps()
			if k < len(names) {
				d = failAt(d, names[k], fmt.Errorf("%w: stop after %d", ErrConflict, k))
			}
			if _, err := Run(context.Background(), w.request(), d); k < len(names) && err == nil {
				t.Fatal("first run did not stop")
			}
			result, err := Run(context.Background(), w.request(), w.deps())
			if err != nil {
				t.Fatal(err)
			}
			want := allSteps(false)
			for i := k; i < len(want); i++ {
				want[i].Changed = true
			}
			if !reflect.DeepEqual(result.Steps, want) {
				t.Fatalf("steps %+v, want %+v", result.Steps, want)
			}
		})
	}
}

func TestRunCancellationStopsBeforeTheNextStep(t *testing.T) {
	for _, tc := range []struct {
		cancelIn string
		next     Step
		last     string
	}{
		{"", StepPrepare, ""},
		{"prepare", StepCheck, "prepare"},
		{"init", StepGitHub, "init"},
		{"define", StepMetadata, "define"},
	} {
		t.Run(string(tc.next), func(t *testing.T) {
			w := newWorld(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.cancelIn == "" {
				cancel()
			}
			switch tc.cancelIn {
			case "prepare", "define":
				w.before[tc.cancelIn] = cancel
			case "init":
				w.before["init "+filepath.Join(w.setup.Root, "code-map", "main")+" main"] = cancel
			}
			result, err := Run(ctx, w.request(), w.deps())
			requireStop(t, err, tc.next, context.Canceled)
			if errors.Is(err, ErrUnknown) || errors.Is(err, ErrInvalid) || errors.Is(err, ErrConflict) {
				t.Fatalf("cancellation classified as sentinel: %v", err)
			}
			last := ""
			if len(w.calls) > 0 {
				last = strings.Fields(w.calls[len(w.calls)-1])[0]
			}
			if last != tc.last {
				t.Fatalf("calls %q", w.calls)
			}
			if tc.next == StepGitHub && !reflect.DeepEqual(result.Steps, allSteps(true)[:2]) {
				t.Fatalf("steps %+v", result.Steps)
			}
		})
	}
}

func TestRunWrapperMkdir(t *testing.T) {
	t.Run("existing directory is unchanged", func(t *testing.T) {
		w := newWorld(t)
		if err := os.Mkdir(filepath.Join(w.setup.Root, "code-map"), 0o755); err != nil {
			t.Fatal(err)
		}
		result, err := Run(context.Background(), w.request(), w.deps())
		if err != nil || result.Steps[0] != (StepResult{StepWrapper, false}) || !result.Steps[1].Changed {
			t.Fatalf("steps %+v, error %v", result.Steps, err)
		}
	})
	t.Run("symlink swapped in after precheck is a conflict", func(t *testing.T) {
		w := newWorld(t)
		wrapper := filepath.Join(w.setup.Root, "code-map")
		target := filepath.Join(w.workspace, "elsewhere")
		w.before["observe "+wrapper] = func() {
			if err := os.Mkdir(target, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, wrapper); err != nil {
				t.Fatal(err)
			}
		}
		result, err := Run(context.Background(), w.request(), w.deps())
		requireStop(t, err, StepWrapper, ErrConflict)
		if len(result.Steps) != 0 {
			t.Fatalf("steps %+v", result.Steps)
		}
	})
	t.Run("project root swapped after precheck is unknown", func(t *testing.T) {
		w := newWorld(t)
		root := w.setup.Root
		moved := filepath.Join(w.workspace, "moved")
		other := filepath.Join(w.workspace, "other")
		w.before["observe "+filepath.Join(root, "code-map")] = func() {
			if err := os.Rename(root, moved); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(other, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(other, root); err != nil {
				t.Fatal(err)
			}
		}
		result, err := Run(context.Background(), w.request(), w.deps())
		requireStop(t, err, StepWrapper, ErrUnknown)
		if len(result.Steps) != 0 {
			t.Fatalf("steps %+v", result.Steps)
		}
		if _, err := os.Lstat(filepath.Join(other, "code-map")); err != nil {
			t.Fatalf("mkdir went elsewhere as staged: %v", err)
		}
	})
}
