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
		{[]string{"companion", "add", "ply", "--role", "planning", "--path", "/work/planning", "--format", "json"}, workspace.ProjectCompanionAdd, workspace.CompanionPlanning, "", "/work/planning"},
		{[]string{"companion", "remove", "ply", "--role", "docs", "--path", "/work/docs", "--format", "json"}, workspace.ProjectCompanionRemove, workspace.CompanionDocs, "", "/work/docs"},
		{[]string{"repo-wrapper", "set", "ply", "--repo", "ply", "--wrapper", "/work/ply", "--format", "json"}, workspace.ProjectRepoWrapperSet, "", "ply", "/work/ply"},
		{[]string{"repo-wrapper", "clear", "ply", "--repo", "ply", "--format", "json"}, workspace.ProjectRepoWrapperClear, "", "ply", ""},
	} {
		calls := 0
		services := workspaceProjectMetadataServices{update: func(input workspace.ProjectMetadataInput) (workspace.ProjectMetadataMutation, error) {
			calls++
			if input != (workspace.ProjectMetadataInput{Operation: test.operation, ProjectID: "ply", Role: test.role, RepoID: test.repo, Locator: test.locator}) {
				t.Fatalf("input = %#v", input)
			}
			return workspace.ProjectMetadataMutation{ProjectResult: workspace.ProjectResult{Workspace: "/work", Project: workspace.ProjectRecord{ID: "ply", Name: "Ply", Wrapper: "/work/ply", RepoIDs: []workspace.RepoID{"ply"}}}, Changed: true}, nil
		}}
		out, err := executeProjectMetadata(t, services, test.args...)
		var decoded map[string]any
		if err != nil || json.Unmarshal([]byte(out), &decoded) != nil || strings.Count(out, "\n") != 1 || decoded["kind"] != workspace.ProjectReadbackKind || decoded["changed"] != true || calls != 1 {
			t.Fatalf("output %q error %v calls %d", out, err, calls)
		}
	}
}

func TestWorkspaceProjectMetadataCommandsRejectInvalidInputBeforeServices(t *testing.T) {
	services := workspaceProjectMetadataServices{update: func(workspace.ProjectMetadataInput) (workspace.ProjectMetadataMutation, error) {
		t.Fatal("invalid input reached metadata service")
		return workspace.ProjectMetadataMutation{}, nil
	}}
	for _, args := range [][]string{
		{"companion", "add", "ply", "--role", "backup", "--path", "/work"},
		{"companion", "add", "ply", "--role", "planning"},
		{"companion", "add", "ply", "--role", "planning", "--path", "/work", "--actor", "codex"},
		{"companion", "remove", "ply", "--role", "docs", "--path", "/work", "--expected-revision", "0"},
		{"repo-wrapper", "set", "ply", "--repo", "ply"},
		{"repo-wrapper", "clear", "ply", "--repo", "ply", "--format", "yaml"},
		{"metadata", "show", "ply"},
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
