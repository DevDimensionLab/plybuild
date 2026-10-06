package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func TestWorkspaceProjectMetadataCommandsHaveExplicitEffectsAndJSON(t *testing.T) {
	for _, test := range []struct {
		args      []string
		operation string
		role      workspace.CompanionRole
		repo      workspace.RepoID
		locator   string
	}{
		{[]string{"companion", "add", "ply", "--role", "planning", "--path", "/work/planning", "--actor", "codex", "--expected-revision", "0", "--format", "json"}, workspace.ProjectCompanionAdd, workspace.CompanionPlanning, "", "/work/planning"},
		{[]string{"companion", "remove", "ply", "--role", "docs", "--path", "/work/docs", "--actor", "codex", "--expected-revision", "0", "--format", "json"}, workspace.ProjectCompanionRemove, workspace.CompanionDocs, "", "/work/docs"},
		{[]string{"repo-wrapper", "set", "ply", "--repo", "ply", "--wrapper", "/work/ply", "--actor", "codex", "--expected-revision", "0", "--format", "json"}, workspace.ProjectRepoWrapperSet, "", "ply", "/work/ply"},
		{[]string{"repo-wrapper", "clear", "ply", "--repo", "ply", "--actor", "codex", "--expected-revision", "0", "--format", "json"}, workspace.ProjectRepoWrapperClear, "", "ply", ""},
	} {
		calls := 0
		services := workspaceProjectMetadataServices{update: func(input workspace.ProjectMetadataInput) (workspace.ProjectMetadataMutation, error) {
			calls++
			if input.Operation != test.operation || input.ProjectID != "ply" || input.Role != test.role || input.RepoID != test.repo || input.Locator != test.locator || input.ActorClaim != "codex" || input.ExpectedRevision == nil || *input.ExpectedRevision != 0 {
				t.Fatalf("input = %#v", input)
			}
			return workspace.ProjectMetadataMutation{ProjectMetadataReadback: workspace.ProjectMetadataReadback{Kind: workspace.ProjectMetadataMutationKind, SchemaVersion: 1, ProjectID: "ply", Revision: 1, Companions: []workspace.ProjectCompanion{}, Repositories: []workspace.ProjectRepositoryMetadata{}, Events: []workspace.ProjectMetadataEvent{}}, Changed: true}, nil
		}}
		out, err := executeProjectMetadata(t, services, test.args...)
		var decoded map[string]any
		if err != nil || json.Unmarshal([]byte(out), &decoded) != nil || strings.Count(out, "\n") != 1 || decoded["kind"] != workspace.ProjectMetadataMutationKind || calls != 1 {
			t.Fatalf("output %q error %v calls %d", out, err, calls)
		}
	}
}

func TestWorkspaceProjectMetadataCommandsRejectInvalidInputBeforeServices(t *testing.T) {
	services := workspaceProjectMetadataServices{update: func(workspace.ProjectMetadataInput) (workspace.ProjectMetadataMutation, error) {
		t.Fatal("invalid input reached metadata service")
		return workspace.ProjectMetadataMutation{}, nil
	}, show: func(workspace.ProjectID) (workspace.ProjectMetadataReadback, error) {
		t.Fatal("invalid input reached show service")
		return workspace.ProjectMetadataReadback{}, nil
	}}
	for _, args := range [][]string{
		{"companion", "add", "ply", "--role", "backup", "--path", "/work", "--actor", "codex"},
		{"companion", "add", "ply", "--role", "planning", "--path", "/work"},
		{"companion", "remove", "ply", "--role", "docs", "--path", "/work", "--actor", "bad\nactor"},
		{"repo-wrapper", "set", "ply", "--repo", "ply", "--actor", "codex"},
		{"repo-wrapper", "clear", "ply", "--repo", "ply", "--actor", "codex", "--expected-revision", "-1"},
		{"metadata", "show", "ply", "--format", "yaml"},
		{"metadata", "show", "BAD"},
	} {
		if out, err := executeProjectMetadata(t, services, args...); err == nil || out != "" {
			t.Fatalf("accepted %q: %q %v", args, out, err)
		}
	}
}

func executeProjectMetadata(t *testing.T, services workspaceProjectMetadataServices, args ...string) (string, error) {
	t.Helper()
	command := &cobra.Command{Use: "project", SilenceUsage: true, SilenceErrors: true}
	command.AddCommand(newWorkspaceProjectMetadataCommandsWithServices(services)...)
	var out, errors bytes.Buffer
	command.SetOut(&out)
	command.SetErr(&errors)
	command.SetArgs(args)
	err := command.Execute()
	return out.String(), err
}
