package humaninput

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// store is the human-input store of Project id in the fixture. Tests reach
// store paths only through it.
func (f *hiFixture) store(id workspace.ProjectID) store {
	f.t.Helper()
	works, err := workspace.DefinedWork(f.root)
	if err != nil {
		f.t.Fatal(err)
	}
	for _, work := range works {
		if work.Project == id {
			return store{work}
		}
	}
	f.t.Fatalf("no Project %s", id)
	return store{}
}

// Questions live in the work of their Project; ID-only operations find them
// there, and a store left in .ply by an earlier version is never read.
func TestQuestionsLiveInTheirProjectsWork(t *testing.T) {
	f := newFixture(t)
	old := filepath.Join(f.root, workspace.MarkerDirectory, "human-input", "v1", "questions")
	mustDo(t, os.MkdirAll(old, 0o700))
	ghost := questionIDFor(f.root, "other", "ghost")
	mustDo(t, os.WriteFile(filepath.Join(old, ghost+".json"), []byte("not a question"), 0o600))
	mine := f.ask(f.choice("mine"))
	draft := f.choice("theirs")
	draft.ProjectID, draft.OriginCWD = "other", f.otherOrigin
	theirs := f.ask(draft)
	for project, id := range map[workspace.ProjectID]string{"ply": mine.QuestionID, "other": theirs.QuestionID} {
		dir := f.store(project).dir()
		if !strings.HasPrefix(dir, filepath.Join(f.root, string(project), "data", "work")+string(filepath.Separator)) {
			t.Fatalf("store of %s is %s", project, dir)
		}
		if _, err := os.Stat(questionPath(f.store(project), id)); err != nil {
			t.Fatalf("question of %s is not in its work: %v", project, err)
		}
	}
	if _, err := os.Stat(questionPath(f.store("ply"), theirs.QuestionID)); !os.IsNotExist(err) {
		t.Fatalf("other's question is in ply's work: %v", err)
	}
	if _, err := Show(f.d, f.root, ghost); codeOf(err) != CodeNotFound {
		t.Fatalf("old .ply question was read: %v", err)
	}
	shown, err := Show(f.d, f.root, theirs.QuestionID)
	if err != nil || shown.Scope.ProjectID != "other" {
		t.Fatalf("show = %+v %v", shown, err)
	}
	p := f.preview(theirs.QuestionID, "request-other-1", AnswerValue{AnswerKind: AnswerChoice, OptionID: ptr("short")})
	if s, err := Apply(f.d, f.root, p.Request); err != nil || s.State != StatusRecorded {
		t.Fatalf("apply = %+v %v", s, err)
	}
	if _, err := os.Stat(answerPath(f.store("other"), theirs.QuestionID)); err != nil {
		t.Fatalf("answer is not beside its question: %v", err)
	}
	if q, err := Acknowledge(f.d, f.root, theirs.QuestionID, p.Request.RequestID, p.RequestSHA256); err != nil || q.State != StateReceived {
		t.Fatalf("acknowledge = %+v %v", q, err)
	}
	if s, err := Status(f.d, f.root, theirs.QuestionID, p.Request.RequestID); err != nil || s.State != StatusRecorded {
		t.Fatalf("status = %+v %v", s, err)
	}
	list, err := List(f.d, f.root, ListQuery{ProjectID: "other"})
	if err != nil || len(list.Questions) != 1 || list.Questions[0].QuestionID != theirs.QuestionID {
		t.Fatalf("list = %+v %v", list.Questions, err)
	}
	if b, _ := os.ReadFile(filepath.Join(old, ghost+".json")); string(b) != "not a question" {
		t.Fatal("old .ply store was changed")
	}
}
