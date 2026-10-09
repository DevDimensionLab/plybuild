package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/taskrun"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func TestRegistryFailureProducesExactlyOneJSONDiagnostic(t *testing.T) {
	valid := "format_version: 1\nepics: []\ntasks: []\nworktree_operations: []\n"
	for _, fixture := range []struct{ raw, code string }{
		{strings.Replace(valid, "version: 1", "version: 6", 1), "registry_unsupported_format"},
		{valid + "authorization_override: true\n", "registry_unknown_field"},
		{valid + "broken: [", "registry_malformed_yaml"},
		{strings.Replace(valid, "epics: []", "epics: null", 1), "registry_invalid_content"},
	} {
		for _, args := range [][]string{
			{"workflow", "execute", "--check", "--format", "json"},
			{"workspace", "task", "prepare", "--next", "--format", "json"},
			{"workspace", "epic", "list", "--format", "json"},
		} {
			t.Run(fixture.code+"/"+strings.Join(args, " "), func(t *testing.T) {
				d, dir := viewFixture(t)
				for _, lock := range []string{"projects.lock", "work-items.lock"} {
					if err := os.WriteFile(filepath.Join(dir, ".ply", lock), nil, 0o600); err != nil {
						t.Fatal(err)
					}
				}
				path := filepath.Join(dir, ".ply", "work-items.yaml")
				before := []byte(fixture.raw)
				if err := os.WriteFile(path, before, 0o600); err != nil {
					t.Fatal(err)
				}
				r := &cobra.Command{Use: "ply"}
				w := &cobra.Command{Use: "workflow"}
				w.AddCommand(newWorkflowExecuteCommand(taskrun.SystemDependencies(d)))
				r.AddCommand(w, newWorkspaceCommand(d))
				original := RootCmd
				RootCmd = r
				defer func() { RootCmd = original }()
				var out, stderr bytes.Buffer
				r.SetOut(&out)
				r.SetErr(&stderr)
				r.SetArgs(args)
				runErr := ExecuteE()
				if runErr == nil {
					t.Fatal("unsupported format succeeded")
				}
				var diagnostic workspace.RegistryDiagnostic
				if err := json.Unmarshal(stderr.Bytes(), &diagnostic); err != nil {
					t.Fatalf("stderr is not one JSON diagnostic: %v\n%s\nCommand error: %v", err, stderr.String(), runErr)
				}
				if out.Len() > 0 && !json.Valid(out.Bytes()) {
					t.Fatalf("stdout is not a single command result: %s", out.String())
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
