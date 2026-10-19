//go:build !windows

package taskrun

import "testing"

func TestAR02EmptyDiagnosticPreservesConcreteOutcome(t *testing.T) {
	f := newAW01Fixture(t)
	f.s.Result.Worker.Iteration = f.s.Request.Worker.MaxAttempts
	f.save()
	f.own()
	cause := workerCorrection(f.d, f.s, "Concrete completed review still requires changes")
	if cause == nil {
		t.Fatal("fixture has not exhausted its attempts")
	}
	if err := workerBlocked(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID, cause); err != nil {
		t.Fatal(err)
	}
	if err := workerStateUpdate(f.d, f.s.Request.WorkspaceRoot, f.s.Result.RunID, func(s *workflowState) error { s.Result.Worker.Reason = ""; return nil }); err != nil {
		t.Fatal(err)
	}
	f.reload()
	if !f.s.Result.Worker.QueueReleaseAdvice() {
		t.Fatal("clearing a diagnostic lost the concrete outcome")
	}
}

func TestAR02LegacyReleasePolicy(t *testing.T) {
	for _, row := range []struct {
		reason string
		want   bool
	}{
		{"worker reached its explicit limit of 3 developer attempts", true},
		{"developer blocked: requirement cannot be satisfied", true},
		{"review did not give an unambiguous pass: blocked; insufficient evidence", true},
		{"reviewer blocked: never a native release outcome", false},
		{"provider usage-limit retry bound exhausted", false},
		{"attempt completion unknown", false},
	} {
		t.Run(row.reason, func(t *testing.T) {
			w := WorkerStatus{Phase: "blocked", Stage: "reviewer", Job: &WorkerJob{Kind: "reviewer"}, Reason: row.reason}
			if w.QueueReleaseAdvice() != row.want {
				t.Fatal("changed historical native release policy")
			}
		})
	}
}
