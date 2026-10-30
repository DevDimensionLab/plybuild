package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

func TestHumanOwnedWorkflowRejectsAutomaticControllers(t *testing.T) {
	workflow := newWorkflowCommand(workflowhandoff.Dependencies{})
	manualStart := false
	for _, command := range workflow.Commands() {
		switch command.Name() {
		case "run":
			manualStart = true
		case "execute", "week", "worker", "control", "delivery":
			t.Errorf("automatic controller command %q is still available", command.Name())
		}
	}
	if !manualStart {
		t.Fatal("explicit workflow run start was removed")
	}
}

func TestHumanOwnedRunCannotStartRemovedDeliveryProcess(t *testing.T) {
	for _, version := range []string{"2", "3", "2.0", "3.0", "2e0", "3e0"} {
		t.Run(version, func(t *testing.T) {
			for _, mode := range []string{"default", "check", "apply"} {
				t.Run(mode, func(t *testing.T) {
					root, err := filepath.EvalSymlinks(t.TempDir())
					if err != nil {
						t.Fatal(err)
					}
					path := filepath.Join(root, "request.json")
					workspace, err := json.Marshal(root)
					if err != nil {
						t.Fatal(err)
					}
					request := []byte(`{"schema_version":` + version + `,"workspace_root":` + string(workspace) + `}`)
					if err := os.WriteFile(path, request, 0600); err != nil {
						t.Fatal(err)
					}
					command := newWorkflowRunCommand(taskrun.Dependencies{})
					var output bytes.Buffer
					command.SetOut(&output)
					command.SetErr(&output)
					args := []string{"start", "--file", path, "--format", "json"}
					switch mode {
					case "check":
						args = append(args, "--check")
					case "apply":
						t.Setenv("HERDR_ENV", "1")
						args = append(args, "--apply", "--confirm", "sha256:"+strings.Repeat("a", 64))
					}
					command.SetArgs(args)
					err = command.Execute()
					wantError := "delivery-owner start has been removed"
					if strings.ContainsAny(version, ".e") {
						wantError = "unsupported JSON number"
					}
					if err == nil || !strings.Contains(err.Error(), wantError) {
						t.Fatalf("removed start surface not rejected: %v\n%s", err, output.String())
					}
					entries, err := os.ReadDir(root)
					if err != nil || len(entries) != 1 || entries[0].Name() != "request.json" {
						t.Fatalf("rejected request changed workspace state: %v %v", entries, err)
					}
					preserved, err := os.ReadFile(path)
					if err != nil || !bytes.Equal(preserved, request) {
						t.Fatalf("rejected request changed input: %q %v", preserved, err)
					}
				})
			}
		})
	}
}
