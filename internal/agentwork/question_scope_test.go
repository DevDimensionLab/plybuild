package agentwork

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/humaninput"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// questionWorkspace is a physical workspace with one registered project whose
// questions are published through the public humaninput surface.
func questionWorkspace(t *testing.T) (workspace.Dependencies, humaninput.Dependencies, string, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(root, "ply", "main")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"init", "-b", "main"}, {"config", "user.name", "F03"}, {"config", "user.email", "f03@example.invalid"}, {"commit", "--allow-empty", "-m", "base"}} {
		if out, err := exec.Command("git", append([]string{"-C", repo}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	d := workspace.SystemDependencies()
	d.Files = rootedFiles{d.Files, root}
	if _, err := workspace.Init(d); err != nil {
		t.Fatal(err)
	}
	if _, err := workspace.AddProject(d, workspace.ProjectAddInput{ProjectID: "ply", Name: "Ply", Wrapper: filepath.Join(root, "ply"), Repos: []workspace.RepoInput{{RepoID: "ply", Path: repo}}}); err != nil {
		t.Fatal(err)
	}
	return d, humaninput.Dependencies{Workspace: d, Now: time.Now}, root, repo
}

func askQuestion(t *testing.T, h humaninput.Dependencies, root, origin, key string, session *string) humaninput.Question {
	t.Helper()
	q, err := humaninput.AskDraft(h, root, humaninput.QuestionDraft{Kind: humaninput.DraftKind, SchemaVersion: humaninput.SchemaVersion, PublicationKey: key, ProjectID: "ply", ActivityID: "ordinary/f03", ActivityTitle: "F03 fixture", Provider: "claude", SessionID: session, OriginCWD: origin, Text: "Which output should this task produce?", Reason: "The goal leaves this choice open.", AnswerKind: humaninput.AnswerChoice, Options: []humaninput.Option{{ID: "short", Label: "Short", Consequence: "One page."}, {ID: "full", Label: "Full", Consequence: "All sections."}}})
	if err != nil {
		t.Fatal(err)
	}
	return q
}

// continueQuestion answers, acknowledges and continues q; it returns after the
// requested lifecycle step so a later fact can be damaged.
func continueQuestion(t *testing.T, h humaninput.Dependencies, root, id, until string) {
	t.Helper()
	p, err := humaninput.PreviewAnswer(h, root, id, "request-"+id, humaninput.AnswerValue{AnswerKind: humaninput.AnswerChoice, OptionID: ptr("short")})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := humaninput.Apply(h, root, p.Request); err != nil {
		t.Fatal(err)
	}
	if until == "answer" {
		return
	}
	sha := humaninput.RequestSHA256(p.Request)
	if _, err := humaninput.Acknowledge(h, root, id, p.Request.RequestID, sha); err != nil {
		t.Fatal(err)
	}
	if until == "receipt" {
		return
	}
	if _, err := humaninput.Continue(h, root, id, p.Request.RequestID, sha, "Used the answer."); err != nil {
		t.Fatal(err)
	}
}

func damage(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"truncated":`), 0o600); err != nil {
		t.Fatal(err)
	}
}

func treeSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
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

// F03: damaged history of a validated question on session B vetoes B only;
// an otherwise authorized status follow-up to unrelated A is sent exactly once.
// Unidentifiable questions and an own question still keep A from being sent.
func TestNativeQuestionHistoryVetoesOnlyItsVerifiedSession(t *testing.T) {
	for _, tc := range []struct {
		damage string
		aSent  bool
	}{
		{"answer", true},
		{"receipt", true},
		{"continuation", true},
		{"withdrawal-after-continued", true},
		{"question-digest", false},
		{"inventory-entry", false},
		{"open-question-without-session", false},
		{"own-open-question", false},
	} {
		t.Run(tc.damage, func(t *testing.T) {
			d, h, root, repo := questionWorkspace(t)
			b := askQuestion(t, h, root, repo, "f03/b", ptr("session-B"))
			store := filepath.Join(root, ".ply", "human-input", "v1")
			fact := func(kind string) string { return filepath.Join(store, kind, b.QuestionID+".json") }
			switch tc.damage {
			case "answer":
				continueQuestion(t, h, root, b.QuestionID, "answer")
				damage(t, fact("answers"))
			case "receipt":
				continueQuestion(t, h, root, b.QuestionID, "receipt")
				damage(t, fact("receipts"))
			case "continuation":
				continueQuestion(t, h, root, b.QuestionID, "continued")
				damage(t, fact("continuations"))
			case "withdrawal-after-continued":
				continueQuestion(t, h, root, b.QuestionID, "continued")
				damage(t, fact("withdrawals"))
			case "question-digest":
				raw, err := os.ReadFile(fact("questions"))
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(fact("questions"), bytes.Replace(raw, []byte("Which output"), []byte("Which outpuT"), 1), 0o600); err != nil {
					t.Fatal(err)
				}
			case "inventory-entry":
				damage(t, filepath.Join(store, "questions", "unexpected.json"))
			case "open-question-without-session":
				askQuestion(t, h, root, repo, "f03/unbound", nil)
			case "own-open-question":
				continueQuestion(t, h, root, b.QuestionID, "answer")
				damage(t, fact("answers"))
				askQuestion(t, h, root, repo, "f03/a", ptr("session-A"))
			}
			// The public list keeps rejecting damaged history.
			if _, err := humaninput.List(h, root, humaninput.ListQuery{ProjectID: "ply"}); (err == nil) != (tc.damage == "open-question-without-session") {
				t.Fatalf("public list contract changed: %v", err)
			}

			now := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
			a := Identity{Provider: "claude", Session: "session-A", Workspace: "w1", Tab: "tab-A", Pane: "pane-A", Terminal: "terminal-A"}
			bi := Identity{Provider: "claude", Session: "session-B", Workspace: "w1", Tab: "tab-B", Pane: "pane-B", Terminal: "terminal-B"}
			rt := &fakeRuntime{agents: []Observation{{Identity: a, Status: "idle", Ready: boolp(true)}, {Identity: bi, Status: "idle", Ready: boolp(true)}}}
			sent := map[Identity]int{}
			before := treeSnapshot(t, store)
			binary, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			qualifyTestCallback(t, root, binary)
			s := &Service{Root: root, Now: func() time.Time { return now }, Runtime: rt, Native: func() NativeSnapshot { return readNative(d, root) }, Send: func(i Identity, _ string) error { sent[i]++; return nil }}
			n := readNative(d, root)
			if tc.aSent && (!n.Partial || n.CoverageUnknown || !n.VetoCoverageComplete || n.Vetoes[bi.sessionKey()] != "native_question") {
				t.Fatalf("known B failure not scoped to B with diagnostics: %+v", n)
			}
			for _, target := range []Identity{a, bi} {
				draft := Draft{PublicationKey: "f03/" + target.Session, OwnerRef: "owner", MandateRef: "mandate", Registrar: Registrar{Kind: HumanRequested, ActorRef: "human"}, TargetRef: target.Ref(), WaitKind: "status", WaitRef: target.Ref(), NextCheckAtUTC: now.Format(time.RFC3339)}
				f, err := s.Register(draft)
				if err != nil {
					t.Fatal(err)
				}
				for i := 0; i < 2; i++ {
					if _, err := s.Dispatch(f.ID); err != nil {
						t.Fatal(err)
					}
				}
			}
			wantA := 0
			if tc.aSent {
				wantA = 1
			}
			if sent[a] != wantA || sent[bi] != 0 {
				t.Fatalf("sends A=%d B=%d; want A=%d B=0", sent[a], sent[bi], wantA)
			}
			if after := treeSnapshot(t, store); !reflect.DeepEqual(before, after) {
				t.Fatalf("native question reads changed the human-input store:\n%v\n%v", before, after)
			}
		})
	}
}
