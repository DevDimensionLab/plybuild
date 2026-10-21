package taskrun

import (
	"testing"
	"time"
)

// Regression for review 1: public follow changes observations while the bound
// owner, original stop, source and coordinator's selected input stay the same.
func TestSelectedBackupRetirementRecoversAfterNativeFollow(t *testing.T) {
	for _, registered := range []bool{false, true} {
		t.Run(map[bool]string{false: "before_apply", true: "after_registration"}[registered], func(t *testing.T) {
			f, o := retireFixture(t)
			o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
			retireAsCoordinator(&f)
			in := BackupGoalInput{TaskID: "backup", SpecID: "backup-goal", Provider: "claude", Retire: &DeliveryRetireInput{Observer: "coordinator", Reason: "Selected replacement"}}
			p, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
			if err != nil {
				t.Fatal(err)
			}
			var applied BackupGoalResult
			if registered {
				applied, err = WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
				if err != nil {
					t.Fatal(err)
				}
			}
			f.D.Now = func() time.Time { return time.Date(2026, 10, 22, 1, 3, 0, 0, time.UTC) }
			shown, followErr := WorkflowFollow(f.D, f.R.WorkspaceRoot, o.RunID, 300*time.Millisecond)
			t.Logf("native follow observation: state=%s observation=%s at=%s error=%v", shown.Transport.State, shown.Transport.Observation, shown.Transport.ObservedAt, followErr)
			if equal(shown.Transport, o.Transport) {
				t.Fatal("fixture did not change transport observation")
			}
			q, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, false, "")
			if err != nil {
				t.Fatalf("readback after follow lost recovery: %v", err)
			}
			current, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
			if err != nil {
				t.Fatal(err)
			}
			if !equal(q.Retirement.Owner, current.Result.Transport) {
				t.Fatal("full current transport observation missing from readback")
			}
			if q.Confirmation != p.Confirmation || q.Retirement.Binding != p.Retirement.Binding {
				t.Fatal("observation changed semantic confirmation")
			}
			if q.Retirement.StopBinding != p.Retirement.StopBinding || q.Retirement.Owner.AgentSessionID != p.Retirement.Owner.AgentSessionID {
				t.Fatal("fixture changed semantic identity")
			}
			a, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
			if err != nil {
				t.Fatalf("identical apply after native follow: %v", err)
			}
			if !a.Applied || registered && a.Goal != applied.Goal {
				t.Fatal("recovery changed registered refs")
			}
			after := deliveryCLISnapshot(t, f.R.WorkspaceRoot)
			retry, err := WorkflowBackupGoal(f.D, f.R.WorkspaceRoot, o.RunID, in, true, p.Confirmation)
			if err != nil || retry.Goal != a.Goal || after != deliveryCLISnapshot(t, f.R.WorkspaceRoot) {
				t.Fatalf("retry duplicated effects: %v", err)
			}
		})
	}
}

func TestSelectedBackupRetirementBindsOwnerAndSourceIdentity(t *testing.T) {
	f, o := retireFixture(t)
	o = retireStopFixture(t, f, o, DeliveryEffectNotStarted)
	s, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
	if err != nil {
		t.Fatal(err)
	}
	selected := &DeliveryRetireInput{Observer: "coordinator", Reason: "Selected replacement"}
	original, _, err := backupRetirement(s, selected, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Incidental transport facts remain visible without entering the hash.
	observed := s
	pending := true
	observed.Result.Transport.State = "unknown"
	observed.Result.Transport.Observation = "cached"
	observed.Result.Transport.ObservedAt = "2026-10-22T02:00:00Z"
	observed.Result.Transport.LaunchPending = &pending
	observedPreview, _, err := backupRetirement(observed, selected, nil)
	if err != nil || observedPreview.Binding != original.Binding || !equal(observedPreview.Owner, observed.Result.Transport) {
		t.Fatalf("observation leaked into binding or was hidden: %v", err)
	}
	for _, tc := range []struct {
		name   string
		change func(*workflowState)
	}{
		{"workspace", func(s *workflowState) { s.Result.Transport.WorkspaceID = "another-workspace" }},
		{"tab", func(s *workflowState) { s.Result.Transport.TabID = "another-tab" }},
		{"pane", func(s *workflowState) { s.Result.Transport.PaneID = "another-pane" }},
		{"terminal", func(s *workflowState) { s.Result.Transport.TerminalID = "another-terminal" }},
		{"session", func(s *workflowState) { s.Result.Transport.AgentSessionID = "another-session" }},
		{"source-oid", func(s *workflowState) { s.Observed.Target.OID = "another-commit" }},
		{"source-ref", func(s *workflowState) { s.Observed.Target.Ref = "refs/heads/another-source" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// In-memory value copies exercise the binding; no persisted fixture or
			// native owner/session is patched or replaced.
			changed := s
			tc.change(&changed)
			p, _, err := backupRetirement(changed, selected, nil)
			if err == nil && p.Binding == original.Binding {
				t.Fatal("changed semantic identity retained confirmation")
			}
		})
	}
}
