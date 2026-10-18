package humaninput

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// hiFixture is a physical workspace under /private/tmp, outside any real
// workspace, with project ply (repo ply, Epic and Task) and project other.
type hiFixture struct {
	t                  *testing.T
	root, repo, origin string
	otherOrigin        string
	d                  Dependencies
	now                time.Time
}

func newFixture(t *testing.T) *hiFixture {
	t.Helper()
	tmp, err := os.MkdirTemp("/private/tmp", "hi01-test.")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmp) })
	root, err := filepath.EvalSymlinks(tmp)
	if err != nil || root != tmp {
		t.Fatalf("fixture root is not physical: %s %v", root, err)
	}
	f := &hiFixture{t: t, root: root, repo: filepath.Join(root, "ply", "main"), origin: filepath.Join(root, "ply", "main", "src"), otherOrigin: filepath.Join(root, "other", "main"), now: time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)}
	for _, repo := range []string{f.repo, f.otherOrigin} {
		if err := os.MkdirAll(repo, 0o755); err != nil {
			t.Fatal(err)
		}
		git(t, repo, "init", "-b", "main")
		git(t, repo, "config", "user.name", "HI01 Test")
		git(t, repo, "config", "user.email", "hi01@example.invalid")
		if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("base\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		git(t, repo, "add", "README.md")
		git(t, repo, "commit", "-m", "base")
	}
	if err := os.Mkdir(f.origin, 0o755); err != nil {
		t.Fatal(err)
	}
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	w := workspace.SystemDependencies()
	if _, err := workspace.Init(w); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.AddProject(w, workspace.ProjectAddInput{ProjectID: "ply", Name: "Ply", Wrapper: filepath.Join(root, "ply"), Repos: []workspace.RepoInput{{RepoID: "ply", Path: f.repo}}}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.AddProject(w, workspace.ProjectAddInput{ProjectID: "other", Name: "Other", Wrapper: filepath.Join(root, "other"), Repos: []workspace.RepoInput{{RepoID: "other", Path: f.otherOrigin}}}); err != nil {
		t.Fatal(err)
	}
	ref, oid := gitOut(t, f.repo, "symbolic-ref", "HEAD"), gitOut(t, f.repo, "rev-parse", "HEAD")
	if _, err := workspace.AdoptEpic(w, workspace.EpicAdoptInput{EpicID: "epic", Title: "Fixture Epic", ProjectID: "ply", RepoID: "ply", Worktree: f.repo, Ref: ref, ExpectedOID: oid}); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.CreateTask(w, workspace.TaskCreateInput{TaskID: "task", Title: "Fixture Task", Description: "Fixture need", ParentEpicID: "epic", ProjectID: "ply", RepoID: "ply"}); err != nil {
		t.Fatal(err)
	}
	f.d = Dependencies{Workspace: w, Now: func() time.Time { f.now = f.now.Add(time.Second); return f.now }, Poll: 5 * time.Millisecond}
	return f
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	if out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
}
func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}

func (f *hiFixture) choice(key string) QuestionDraft {
	return QuestionDraft{
		Kind: DraftKind, SchemaVersion: SchemaVersion, PublicationKey: key, ProjectID: "ply", RepoID: nil, TaskID: nil,
		ActivityID: "ordinary/conversation-1", ActivityTitle: "Document the new export format", Provider: "claude", SessionID: ptr("session-1"),
		OriginCWD: f.origin, Text: "Which output should this task produce?", Reason: "The goal leaves this product choice open.",
		AnswerKind: AnswerChoice, Options: []Option{{ID: "short", Label: "Short version", Consequence: "One page."}, {ID: "full", Label: "Full version", Consequence: "All sections."}},
	}
}
func (f *hiFixture) text(key string) QuestionDraft {
	q := f.choice(key)
	q.AnswerKind, q.Options, q.Text = AnswerText, []Option{}, "What exact heading should the README use?"
	return q
}
func (f *hiFixture) ask(q QuestionDraft) Question {
	f.t.Helper()
	out, err := AskDraft(f.d, f.root, q)
	if err != nil {
		f.t.Fatalf("ask: %v", err)
	}
	return out
}
func (f *hiFixture) preview(id, request string, v AnswerValue) Preview {
	f.t.Helper()
	out, err := PreviewAnswer(f.d, f.root, id, request, v)
	if err != nil {
		f.t.Fatalf("preview: %v", err)
	}
	return out
}
func (f *hiFixture) answer(id, request, option string) AnswerRequest {
	f.t.Helper()
	p := f.preview(id, request, AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr(option)})
	if s, err := Apply(f.d, p.Request); err != nil || s.State != StatusRecorded {
		f.t.Fatalf("apply: %+v %v", s, err)
	}
	return p.Request
}

func codeOf(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	if err == nil {
		return ""
	}
	return "untyped: " + err.Error()
}
func wantCode(t *testing.T, err error, code string) {
	t.Helper()
	if got := codeOf(err); got != code {
		t.Fatalf("want %s, got %s", code, got)
	}
}

// snapshot captures every path, type, mode and byte below root.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := os.Lstat(p)
		if err != nil {
			return err
		}
		v := info.Mode().String()
		if info.Mode().IsRegular() {
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			v += ":" + string(b)
		}
		out[p] = v
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func sameSnapshot(t *testing.T, before, after map[string]string) {
	t.Helper()
	var diff []string
	for k, v := range before {
		if after[k] != v {
			diff = append(diff, "changed/removed "+k)
		}
	}
	for k := range after {
		if _, ok := before[k]; !ok {
			diff = append(diff, "created "+k)
		}
	}
	sort.Strings(diff)
	if len(diff) > 0 {
		t.Fatalf("read surface changed the workspace: %v", diff)
	}
}
func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
