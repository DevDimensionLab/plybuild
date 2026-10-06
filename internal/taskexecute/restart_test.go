package taskexecute

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
)

func TestExecuteRestartUsesSelectedSpecAndPreservesTask(t *testing.T) {
	for _, closed := range []bool{false, true} {
		t.Run(map[bool]string{false: "existing_shell", true: "closed_tab"}[closed], func(t *testing.T) {
			d, in, root, _ := launcherFixtureForProvider(t, "codex")
			in.Next, in.SpecID, in.Restart = false, "observable-goal", true
			d.StartupProcessTree = func(pid int) ([]taskrun.StartupProcess, error) {
				if closed {
					t.Fatal("closed terminal recovery inspected the process table")
				}
				return []taskrun.StartupProcess{{PID: pid, ParentPID: 1, ProcessGroupID: pid}}, nil
			}
			if err := os.WriteFile(in.Runtime.HerdrExecutable, []byte(recoveryLauncherHerdr), 0700); err != nil {
				t.Fatal(err)
			}
			d.Fault = func(point string) error {
				if point == "workflow_before_agent_send" {
					return errors.New("synthetic interrupted startup")
				}
				return nil
			}
			first, err := Execute(d, in)
			if err == nil || first.Run == nil {
				t.Fatalf("fixture did not preserve its interrupted start: %+v %v", first, err)
			}
			d.Fault = nil
			if closed {
				modelPath := filepath.Join(root, "bin", "recovery-model.json")
				body, err := os.ReadFile(modelPath)
				if err != nil {
					t.Fatal(err)
				}
				var model map[string]any
				if err = json.Unmarshal(body, &model); err != nil {
					t.Fatal(err)
				}
				model["closed"] = true
				body, err = json.Marshal(model)
				if err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(modelPath, body, 0600); err != nil {
					t.Fatal(err)
				}
				in.Runtime.HerdrWorkspace = ""
				t.Setenv("HERDR_WORKSPACE_ID", "w-reopened")
			}
			request, err := os.ReadFile(first.RequestPath)
			if err != nil {
				t.Fatal(err)
			}
			if err = os.WriteFile(in.Runtime.ProviderExecutable, []byte("#!/bin/sh\n# replacement provider in disposable fixture\nexit 0\n"), 0700); err != nil {
				t.Fatal(err)
			}
			registry, err := d.Workspace.WorkItems.Snapshot(root)
			if err != nil {
				t.Fatal(err)
			}
			before := recoveryFiles(t, first.Run.Paths.RunRoot)
			in.Restart, in.Check = true, true
			preview, err := Execute(d, in)
			if err != nil || preview.State != "ready" || preview.Recovery == nil || preview.Recovery.Preview == nil || preview.Recovery.Preview.RunID != first.Run.RunID {
				t.Fatalf("Spec restart did not preview its existing run: %+v %v", preview, err)
			}
			if closed && (preview.Recovery.Preview.TransportMode != "replace_missing_terminal" || preview.Recovery.Preview.DestinationWorkspaceID != "w-reopened") {
				t.Fatalf("Spec restart did not select the current Herdr workspace: %+v", preview.Recovery.Preview)
			}
			if after := recoveryFiles(t, first.Run.Paths.RunRoot); !reflect.DeepEqual(before, after) {
				t.Fatal("Spec restart --check changed the preserved run")
			}
			if _, err := os.Lstat(preview.Recovery.Preview.ControlExecutable.Path); !os.IsNotExist(err) {
				t.Fatal("Spec restart --check copied a control binary")
			}
			in.Check = false
			restarted, err := Execute(d, in)
			if err != nil || restarted.State != "started" || restarted.Run == nil || restarted.Run.RunID != first.Run.RunID || restarted.Run.StartupRecovery == nil {
				t.Fatalf("Spec restart did not start the same Task: %+v %v", restarted, err)
			}
			if got, err := os.ReadFile(first.RequestPath); err != nil || string(got) != string(request) {
				t.Fatal("Spec restart replaced the original request")
			}
			if after, err := d.Workspace.WorkItems.Snapshot(root); err != nil || after.RawSHA256 != registry.RawSHA256 {
				t.Fatal("Spec restart changed the goal or queue ownership")
			}
			modelPath := filepath.Join(root, "bin", "recovery-model.json")
			modelBefore, err := os.ReadFile(modelPath)
			if err != nil {
				t.Fatal(err)
			}
			var model struct{ Tabs, Starts, Prompts int }
			if err := json.Unmarshal(modelBefore, &model); err != nil {
				t.Fatal(err)
			}
			wantTabs := 1
			if closed {
				wantTabs = 2
			}
			if model.Tabs != wantTabs || model.Starts != 1 || model.Prompts != 1 {
				t.Fatalf("unexpected native fixture effects: %+v", model)
			}
			if closed && (restarted.Run.Transport.TerminalID != "replacement-terminal" || restarted.Run.Transport.WorkspaceID != "w-reopened") {
				t.Fatalf("replacement transport not bound: %+v", restarted.Run.Transport)
			}
			// The ordinary command remains a readback after any startup attempt.
			in.Restart = false
			again, err := Execute(d, in)
			if err != nil || again.State != "existing" || again.Run.RunID != first.Run.RunID {
				t.Fatalf("ordinary repeat did not inspect the same Task: %+v %v", again, err)
			}
			if modelAfter, err := os.ReadFile(modelPath); err != nil || string(modelBefore) != string(modelAfter) {
				t.Fatal("ordinary repeat touched Herdr")
			}
		})
	}
}
