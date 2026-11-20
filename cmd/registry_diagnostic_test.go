package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskwork"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func TestRegistryFailureProducesExactlyOneJSONDiagnostic(t *testing.T) {
	fields := []string{"epics", "tasks", "worktree_operations", "task_results", "human_qa_records", "integration_authorities", "integration_intents", "integration_attempts", "integration_results", "task_problem_revisions", "task_spec_revisions", "task_spec_assessments", "task_solution_selections", "task_result_spec_bindings", "task_content_publications", "epic_base_versions", "epic_base_updates", "task_queue_events", "task_preparations", "epic_create_operations"}
	valid := fmt.Sprintf("format_version: %d\nregistry:\n", taskwork.RegistryFormat)
	for _, field := range fields {
		valid += "  " + field + ": []\n"
	}
	for _, fixture := range []struct{ name, raw, code string }{
		{"future format", strings.Replace(valid, "version: 11", "version: 99", 1), "registry_unsupported_format"},
		{"previous format", "format_version: 6\nepics: []\ntasks: []\nworktree_operations: []\n", "registry_unsupported_format"},
		{"missing collection", strings.Replace(valid, "  task_results: []\n", "", 1), "registry_invalid_content"},
		{"unknown authority", valid + "  authorization_override: true\n", "registry_unknown_field"},
		{"malformed YAML", valid + "broken: [", "registry_malformed_yaml"},
		{"null epics", strings.Replace(valid, "epics: []", "epics: null", 1), "registry_invalid_content"},
	} {
		for _, args := range [][]string{
			{"workspace", "epic", "list", "--format", "json"},
		} {
			t.Run(fixture.name+"/"+strings.Join(args, " "), func(t *testing.T) {
				d, dir := viewFixture(t)
				// One Project; reading its registry observes no Git.
				if err := os.MkdirAll(filepath.Join(dir, "project"), 0o755); err != nil {
					t.Fatal(err)
				}
				definitions := `{"format_version": 1, "projects": [{"id": "project", "name": "Project", "wrapper": "project", "repos": [{"id": "repo", "path": "project/repo"}]}]}`
				if err := os.WriteFile(filepath.Join(dir, workspace.DefinitionsFile), []byte(definitions), 0o600); err != nil {
					t.Fatal(err)
				}
				works, err := workspace.DefinedWork(dir)
				if err != nil {
					t.Fatal(err)
				}
				path := works[0].Registry()
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatal(err)
				}
				before := []byte(fixture.raw)
				if err := os.WriteFile(path, before, 0o600); err != nil {
					t.Fatal(err)
				}
				r := &cobra.Command{Use: "ply"}
				r.AddCommand(newWorkspaceCommand(d))
				original := RootCmd
				RootCmd = r
				defer func() { RootCmd = original }()
				var out, stderr bytes.Buffer
				r.SetOut(&out)
				r.SetErr(&stderr)
				r.SetArgs(args)
				runErr := ExecuteE()
				if runErr == nil {
					t.Fatal("invalid registry succeeded")
				}
				var diagnostic taskwork.RegistryDiagnostic
				if err := json.Unmarshal(stderr.Bytes(), &diagnostic); err != nil {
					t.Fatalf("stderr is not one JSON diagnostic: %v\n%s\nCommand error: %v", err, stderr.String(), runErr)
				}
				if out.Len() > 0 {
					t.Fatalf("registry rejection emitted stdout: %s", out.String())
				}
				if diagnostic.Code != fixture.code || diagnostic.NextAction == "" {
					t.Fatalf("missing machine-readable recovery: %+v; command error: %v", diagnostic, runErr)
				}
				after, _ := os.ReadFile(path)
				if !bytes.Equal(before, after) {
					t.Fatal("read-only rejection changed the registry")
				}
			})
		}
	}
}
