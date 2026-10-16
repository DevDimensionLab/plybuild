//go:build darwin

package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func sl01RetryFixture(t *testing.T, legacy bool, retries int) *SL01NativeFixture {
	t.Helper()
	w := &WorkerContract{SchemaVersion: 2, MaxAttempts: 1, AgentTimeoutSeconds: 30, MaxProviderRetries: retries, ProductQuestions: &ProductQuestionContract{1, BoundAnswerMode, 1}}
	if legacy {
		w.SchemaVersion, w.ProductQuestions = 1, nil
	}
	f := NewSL01NativeFixture(t, func(d Dependencies, root, _ string, _ workspace.TaskGoalRef) *ControlContract {
		if legacy {
			return nil
		}
		c := NewControlContract("week", "week_"+strings.Repeat("a", 32), "sha256:"+strings.Repeat("b", 64))
		gate := ParentControlGate{Contract: c, DesiredAction: "continue", Deadline: time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339Nano)}
		if e := d.writeValue(ParentControlPath(root, *c.WeekID), gate); e != nil {
			t.Fatal(e)
		}
		return c
	}, w)
	if !legacy {
		s, e := workflowRead(f.Root, f.RunID)
		if e != nil {
			t.Fatal(e)
		}
		gate, e := ReadParentControl(f.Root, s.Request.ControlContract)
		if e != nil {
			t.Fatal(e)
		}
		gate.ActiveRun = f.RunID
		if e = f.D.replaceValue(ParentControlPath(f.Root, *s.Request.ControlContract.WeekID), gate); e != nil {
			t.Fatal(e)
		}
	}
	return f
}

func sl01OwnAfterRestart(t *testing.T, f *SL01NativeFixture) error {
	t.Helper()
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if e = f.D.WorkerOwner.lock.Close(); e != nil {
		t.Fatal(e)
	}
	f.D.WorkerOwner = nil
	lock, e := workerLock(filepath.Join(workerDirectory(s), "owner.lock"), false)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { lock.Close() })
	return workerClaimOwner(&f.D, s, lock)
}

func TestSL01ZeroProviderRetriesControlsActualLaunches(t *testing.T) {
	for _, role := range []string{"developer", "reviewer"} {
		for _, legacy := range []bool{false, true} {
			name := role + "/new_zero"
			if legacy {
				name = role + "/legacy_zero"
			}
			t.Run(name, func(t *testing.T) {
				f := sl01RetryFixture(t, legacy, 0)
				if role == "reviewer" {
					deliveryTestCandidate(t, f.f, "owned retry regression candidate")
					pretest, e := workflowKeep(f.D, filepath.Join(f.Root, "synthetic-prior-pretest.json"), map[string]string{"purpose": "seeded reviewer input for actual provider retry regression"})
					if e != nil {
						t.Fatal(e)
					}
					if e := workflowUpdate(f.D, f.Root, f.RunID, func(s *workflowState) error {
						w := s.Result.Worker
						w.Stage, w.Iteration = "reviewer", 1
						w.DeveloperSessions = []string{"synthetic-prior-developer"}
						w.CandidateOID = gitOutput(t, f.Task, "rev-parse", "HEAD")
						w.CandidateTree = gitOutput(t, f.Task, "rev-parse", "HEAD^{tree}")
						w.CandidateBinary = s.Request.Runtime.PlyExecutable.Path
						w.Pretest = &pretest
						return nil
					}); e != nil {
						t.Fatal(e)
					}
				} else {
					f.Step(t)
				}
				if e := os.WriteFile(filepath.Join(filepath.Dir(f.Turns), "usage-limit-role"), []byte(role), 0600); e != nil {
					t.Fatal(e)
				}
				f.Step(t) // actual owned process emits a failed usage-limit turn
				s, e := workflowRead(f.Root, f.RunID)
				if e != nil {
					t.Fatal(e)
				}
				original := *s.Result.Worker.Job
				r, e := workerRunReserved(f.D, s)
				if e != nil || r.Process == nil || r.Process.PID == nil || r.Process.StartIdentity == nil || !quiescent(*r.Process) || r.ProviderFailure == nil || r.ProviderFailure.Kind != "usage_limit" || f.TurnCount(t) != 1 {
					t.Fatalf("missing positive owned usage-limit result: %+v %v", r, e)
				}
				if e = workerStep(f.D, s); e != nil {
					if e = workerBlocked(f.D, f.Root, f.RunID, e); e != nil {
						t.Fatal(e)
					}
				}
				if e = sl01OwnAfterRestart(t, f); e != nil {
					t.Fatal(e)
				}
				s, e = workflowRead(f.Root, f.RunID)
				if e != nil {
					t.Fatal(e)
				}
				if !legacy {
					if e = WithWorkerParentGuard(f.Root, s.Request.ControlContract, func() error {
						gate, er := ReadParentControl(f.Root, s.Request.ControlContract)
						if er != nil {
							return er
						}
						gate.ResumeAllowed = true
						gate.Revision++
						return f.D.replaceValue(ParentControlPath(f.Root, *s.Request.ControlContract.WeekID), gate)
					}); e != nil {
						t.Fatal(e)
					}
				}
				if e = workerControlWrite(f.D, s, "resume"); e != nil {
					t.Fatal(e)
				}
				var control workerControl
				if e = readValue(filepath.Join(workerDirectory(s), "control.json"), 64<<10, &control); e != nil {
					t.Fatal(e)
				}
				if e = workerConsumeResume(f.D, s, control); e != nil {
					t.Fatal(e)
				}
				s, e = workflowRead(f.Root, f.RunID)
				if e != nil {
					t.Fatal(e)
				}
				if s.Result.Worker.RetryNotBefore != "" {
					now := time.Now().Add(10 * time.Minute)
					f.D.Now = func() time.Time { return now }
					if e = workerContinueAfterWait(f.D, s); e != nil {
						t.Fatal(e)
					}
					f.Step(t) // legacy keeps its actual second process launch
				} else if e = workerStep(f.D, s); e == nil || !strings.Contains(e.Error(), "retry bound 0 exhausted") {
					t.Fatalf("resumed zero budget did not remain blocked: %v", e)
				}
				want := 1
				if legacy {
					want = 2
				}
				if launches := f.TurnCount(t); launches != want {
					t.Fatalf("actual %s provider launches=%d, want=%d after usage limit and restart/resume", name, launches, want)
				}
				s, e = workflowRead(f.Root, f.RunID)
				if e != nil || !legacy && (s.Result.Worker.ProviderRetries != 0 || s.Result.Worker.RetryNotBefore != "" || *s.Result.Worker.Job != original) {
					t.Fatalf("zero budget changed its retained failure: %+v %v", s.Result.Worker, e)
				}
				t.Logf("actual owned %s launches=%d; same Run after new owner generation and explicit resume", name, want)
			})
		}
	}
}

func TestSL01FailedAnswerNeverUsesNormalProviderRetryBudget(t *testing.T) {
	f := sl01RetryFixture(t, false, 2)
	f.Step(t)
	f.Step(t)
	f.Step(t)
	q, e := ReadProductQuestion(f.D, f.Root, f.RunID)
	if e != nil || q.Eligibility != "ready" {
		t.Fatal(q, e)
	}
	p, e := PreviewBoundAnswer(f.D, f.Root, f.RunID, q.Question.QuestionID, 1, "task_first", "failed-continuation-01")
	if e != nil || p.Request == nil {
		t.Fatal(p, e)
	}
	if _, e = ApplyBoundAnswer(f.D, *p.Request); e != nil {
		t.Fatal(e)
	}
	f.Step(t)
	if e = os.WriteFile(filepath.Join(filepath.Dir(f.Turns), "usage-limit-role"), []byte("developer_answer"), 0600); e != nil {
		t.Fatal(e)
	}
	f.Step(t)
	s, e := workflowRead(f.Root, f.RunID)
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 2; i++ {
		if e = workerStep(f.D, s); e == nil || !strings.Contains(e.Error(), "no automatic retry") {
			t.Fatalf("failed answer used normal retry budget: %v", e)
		}
	}
	if e = sl01OwnAfterRestart(t, f); e == nil {
		t.Fatal("failed answer granted a replacement launch owner")
	}
	status, e := ApplyBoundAnswer(f.D, *p.Request)
	if e != nil || status.State != "outcome_unknown" || f.TurnCount(t) != 2 {
		t.Fatalf("answer replay changed the two actual launches: %+v %v", status, e)
	}
	t.Log("actual owned launches=2 (initial + failed answer); normal retry budget=2 did not authorize redispatch")
}
