package humaninput

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func sessionsOf(sources []Source) []string {
	out := []string{}
	for _, s := range sources {
		out = append(out, s.Provider+":"+*s.SessionID)
	}
	return out
}

// Lifecycle errors are decided before terminal state, so a damaged fact after
// a validated identity blocks only that session, even when the partial view is
// already continued. Unidentifiable questions and sessionless open questions
// leave coverage unknown. One terminal question never hides another's block.
func TestReadQuestionBlockersScopesDamageToValidatedSession(t *testing.T) {
	for _, tc := range []struct {
		name     string
		setup    func(f *hiFixture, q Question)
		blocked  []string
		unknown  bool
		partial  bool
		listFail bool
	}{
		{"waiting", func(*hiFixture, Question) {}, []string{"claude:session-1"}, false, false, false},
		{"continued", func(f *hiFixture, q Question) { f.continued(q.QuestionID) }, []string{}, false, false, false},
		{"withdrawn", func(f *hiFixture, q Question) {
			_, err := Withdraw(f.d, f.root, q.QuestionID, "No longer needed.")
			mustDo(f.t, err)
		}, []string{}, false, false, false},
		{"damaged-answer", func(f *hiFixture, q Question) {
			f.answer(q.QuestionID, "request-1", "short")
			f.damage(answerPath(f.root, q.QuestionID))
		}, []string{"claude:session-1"}, false, true, true},
		{"damaged-receipt", func(f *hiFixture, q Question) {
			f.continued(q.QuestionID)
			f.damage(receiptPath(f.root, q.QuestionID))
		}, []string{"claude:session-1"}, false, true, true},
		{"damaged-continuation", func(f *hiFixture, q Question) {
			f.continued(q.QuestionID)
			f.damage(continuationPath(f.root, q.QuestionID))
		}, []string{"claude:session-1"}, false, true, true},
		{"invalid-withdrawal-after-continued", func(f *hiFixture, q Question) {
			f.continued(q.QuestionID)
			f.damage(withdrawalPath(f.root, q.QuestionID))
		}, []string{"claude:session-1"}, false, true, true},
		{"invalid-question-digest", func(f *hiFixture, q Question) {
			raw, err := os.ReadFile(questionPath(f.root, q.QuestionID))
			mustDo(f.t, err)
			mustDo(f.t, os.WriteFile(questionPath(f.root, q.QuestionID), bytes.Replace(raw, []byte("Which output"), []byte("Which outpuT"), 1), 0o600))
		}, []string{}, true, true, true},
		{"unexpected-inventory-entry", func(f *hiFixture, q Question) {
			f.damage(questionPath(f.root, "unexpected"))
		}, []string{}, true, true, true},
		{"open-without-session", func(f *hiFixture, q Question) {
			d := f.choice("probe/unbound")
			d.SessionID = nil
			f.ask(d)
		}, []string{"claude:session-1"}, true, false, false},
		{"damaged-without-session", func(f *hiFixture, q Question) {
			d := f.choice("probe/unbound")
			d.SessionID = nil
			u := f.ask(d)
			f.continued(u.QuestionID)
			f.damage(continuationPath(f.root, u.QuestionID))
		}, []string{"claude:session-1"}, true, true, true},
		{"terminal-and-damaged-same-session", func(f *hiFixture, q Question) {
			f.continued(q.QuestionID)
			d := f.choice("probe/second")
			d.ProjectID, d.OriginCWD = "other", f.otherOrigin
			s := f.ask(d)
			f.answer(s.QuestionID, "request-2", "short")
			f.damage(answerPath(f.root, s.QuestionID))
		}, []string{"claude:session-1"}, false, true, true},
		{"terminal-and-active-same-session", func(f *hiFixture, q Question) {
			f.continued(q.QuestionID)
			f.ask(f.choice("probe/second"))
		}, []string{"claude:session-1"}, false, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			q := f.ask(f.choice("probe/first"))
			tc.setup(f, q)
			before := snapshot(t, f.root)
			blocked, unknown, partial := ReadQuestionBlockers(f.root)
			sameSnapshot(t, before, snapshot(t, f.root))
			got := sessionsOf(blocked)
			if !reflect.DeepEqual(got, tc.blocked) || unknown != tc.unknown || partial != tc.partial {
				t.Fatalf("blocked=%v unknown=%v partial=%v; want %v %v %v", got, unknown, partial, tc.blocked, tc.unknown, tc.partial)
			}
			_, err := List(f.d, f.root, ListQuery{ProjectID: "ply"})
			if (err != nil) != tc.listFail {
				t.Fatalf("public list contract changed: %v", err)
			}
		})
	}
}

func (f *hiFixture) continued(id string) {
	f.t.Helper()
	r := f.answer(id, "request-"+id, "short")
	_, err := Acknowledge(f.d, f.root, id, r.RequestID, digestOf(r))
	mustDo(f.t, err)
	_, err = Continue(f.d, f.root, id, r.RequestID, digestOf(r), "Used answer.")
	mustDo(f.t, err)
}

func (f *hiFixture) damage(path string) {
	f.t.Helper()
	mustDo(f.t, os.MkdirAll(filepath.Dir(path), 0o700))
	mustDo(f.t, os.WriteFile(path, []byte(`{"truncated":`), 0o600))
}
