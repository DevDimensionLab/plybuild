package taskrun

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// All sessions and writer observations are modeled through existing adapters.
// No provider conversation or physical process observation is performed.
func TestAL40RetiredContinuationRejectsBeforeValidation(t *testing.T) {
	for _, state := range []string{"new", "indexed", "unindexed"} {
		for _, reappeared := range []bool{false, true} {
			name := state + "/absent"
			if reappeared {
				name = state + "/reappeared"
			}
			t.Run(name, func(t *testing.T) {
				f, o := retireFixture(t)
				owner := f.D
				in := DeliveryContinuationInput{RunID: o.RunID, ContextPath: o.Paths.Context, ControlExecutable: f.R.Runtime.PlyExecutable}
				initial, err := WorkflowPreviewDeliveryContinuation(owner, f.R.WorkspaceRoot, in)
				if err != nil {
					t.Fatal("unretired preview setup", err)
				}
				if err = privateDir(initial.Directory); err != nil {
					t.Fatal(err)
				}
				raw, err := os.ReadFile(in.ControlExecutable.Path)
				if err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(initial.ControlExecutable.Path, raw, 0500); err != nil {
					t.Fatal(err)
				}
				if state != "new" {
					d := owner
					if state == "unindexed" {
						d.Fault = func(point string) error {
							if point == "delivery_after_continuation_proof" {
								return os.ErrInvalid
							}
							return nil
						}
					}
					result, err := WorkflowContinueDelivery(d, f.R.WorkspaceRoot, in, initial.Confirmation)
					if state == "indexed" && (err != nil || result.State != "continued") {
						t.Fatal("indexed continuation setup", err)
					}
					if state == "unindexed" {
						if err == nil {
							t.Fatal("publication fault absent")
						}
						if _, err = os.Stat(filepath.Join(initial.Directory, "continuation.json")); err != nil {
							t.Fatal("unindexed proof absent", err)
						}
					}
				}
				o = retireStopFixture(t, f, o, DeliveryEffectPartial)
				workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": true})
				f.D.LocalProcesses = func(string, string) ([]LocalProcess, error) { return []LocalProcess{{PID: os.Getpid()}}, nil }
				retireAsCoordinator(&f)
				if _, err = WorkflowRetireDelivery(f.D, f.R.WorkspaceRoot, o.RunID, DeliveryRetireInput{Observer: "AL40 continuation regression", Reason: "Resources only", RetainEffects: true}); err != nil {
					t.Fatal("retirement setup", err)
				}
				workflowTestModel(t, f.workflowFixture, map[string]any{"closed_transport": !reappeared})
				before, err := workflowRead(f.R.WorkspaceRoot, o.RunID)
				if err != nil {
					t.Fatal(err)
				}
				if (before.Continuation != nil) != (state == "indexed") {
					t.Fatal("unexpected continuation index")
				}
				if _, err = WorkflowShow(f.D, f.R.WorkspaceRoot, o.RunID); err != nil {
					t.Fatal("historical show", err)
				}
				source := &al40CountSource{TaskIntegrationGit: owner.Workspace.IntegrationGit}
				owner.Workspace.IntegrationGit = source
				cwdCalls, executableCalls := 0, 0
				owner.CWD = func() (string, error) { cwdCalls++; return before.Observed.Target.WorktreeLocator, nil }
				owner.Executable = func() (string, error) { executableCalls++; return in.ControlExecutable.Path, nil }
				owner.LocalProcesses = func(string, string) ([]LocalProcess, error) {
					t.Fatal("retired continuation inspected processes")
					return nil, nil
				}
				confirmation := initial.Confirmation
				for _, mutation := range []bool{false, true} {
					source.observations = 0
					cwdCalls = 0
					executableCalls = 0
					calls := workflowTestCalls(t, f.workflowFixture, "agent get") + workflowTestCalls(t, f.workflowFixture, "pane get")
					var got DeliveryContinuationPreview
					if mutation {
						got, err = WorkflowContinueDelivery(owner, f.R.WorkspaceRoot, in, confirmation)
					} else {
						got, err = WorkflowPreviewDeliveryContinuation(owner, f.R.WorkspaceRoot, in)
					}
					if got.Confirmation != "" {
						confirmation = got.Confirmation
					}
					if err == nil || !strings.Contains(err.Error(), "owner_retired") {
						t.Errorf("mutation=%t retired continuation reached %s: %v", mutation, got.State, err)
					}
					if source.observations != 0 || cwdCalls != 0 || executableCalls != 0 {
						t.Errorf("mutation=%t validation before denial: source=%d cwd=%d executable=%d", mutation, source.observations, cwdCalls, executableCalls)
					}
					if after := workflowTestCalls(t, f.workflowFixture, "agent get") + workflowTestCalls(t, f.workflowFixture, "pane get"); after != calls {
						t.Errorf("mutation=%t inspected retired session", mutation)
					}
					after, readErr := workflowRead(f.R.WorkspaceRoot, o.RunID)
					if readErr != nil {
						t.Fatal(readErr)
					}
					if !equal(before, after) {
						t.Errorf("mutation=%t changed retired ownership/history", mutation)
					}
				}
			})
		}
	}
}
