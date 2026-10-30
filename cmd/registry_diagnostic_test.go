package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func TestRegistryFailureProducesExactlyOneJSONDiagnostic(t *testing.T) {
	valid := "format_version: 1\nepics: []\ntasks: []\nworktree_operations: []\n"
	for _, fixture := range []struct{ name, raw, code string }{
		{"future format", strings.Replace(valid, "version: 1", "version: 99", 1), "registry_unsupported_format"},
		{"supported format missing lifecycle", strings.Replace(valid, "version: 1", "version: 6", 1), "registry_invalid_content"},
		{"unknown authority", valid + "authorization_override: true\n", "registry_unknown_field"},
		{"malformed YAML", valid + "broken: [", "registry_malformed_yaml"},
		{"null epics", strings.Replace(valid, "epics: []", "epics: null", 1), "registry_invalid_content"},
	} {
		for _, args := range [][]string{
			{"workspace", "task", "prepare", "--next", "--format", "json"},
			{"workspace", "epic", "list", "--format", "json"},
		} {
			t.Run(fixture.name+"/"+strings.Join(args, " "), func(t *testing.T) {
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
				var diagnostic workspace.RegistryDiagnostic
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
