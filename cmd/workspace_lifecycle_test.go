package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func TestWorkspaceOrganizingLifecycleCommandsValidateAndPreserveIdentity(t *testing.T) {
	for _, kind := range []string{"task", "epic"} {
		calls := 0
		services := workspaceLifecycleServices{
			show: func(gotKind, id string) (workspace.WorkItemLifecycleReadback, error) {
				calls++
				if gotKind != kind || id != "an-item" {
					t.Fatalf("wrong show identity %q %q", gotKind, id)
				}
				return workspace.WorkItemLifecycleReadback{Kind: workspace.WorkItemLifecycleReadbackKind, SchemaVersion: 1, SubjectKind: kind, SubjectID: id, Lifecycle: workspace.LifecycleActive, Events: []workspace.WorkItemLifecycleEvent{}}, nil
			},
			set: func(input workspace.WorkItemLifecycleInput) (workspace.WorkItemLifecycleMutation, error) {
				calls++
				if input.SubjectKind != kind || input.SubjectID != "an-item" || input.State != workspace.LifecycleArchived || input.ActorClaim != "codex" || input.ExpectedRevision == nil || *input.ExpectedRevision != 4 {
					t.Fatalf("wrong set input %#v", input)
				}
				return workspace.WorkItemLifecycleMutation{WorkItemLifecycleReadback: workspace.WorkItemLifecycleReadback{Kind: workspace.WorkItemLifecycleMutationKind, SchemaVersion: 1, SubjectKind: kind, SubjectID: "an-item", Lifecycle: workspace.LifecycleArchived, Revision: 5, Events: []workspace.WorkItemLifecycleEvent{}}, Changed: true}, nil
			},
		}
		for _, args := range [][]string{{"show", "an-item", "--format", "json"}, {"set", "an-item", "archived", "--actor", "codex", "--expected-revision", "4", "--format", "json"}} {
			output, err := executeOrganizingLifecycle(t, kind, services, args...)
			var value map[string]any
			if err != nil || json.Unmarshal([]byte(output), &value) != nil || strings.Count(output, "\n") != 1 || value["subject_kind"] != kind {
				t.Fatalf("output %q, error %v", output, err)
			}
		}
		if calls != 2 {
			t.Fatalf("calls = %d", calls)
		}
		for _, args := range [][]string{
			{"set", "an-item", "deleted", "--actor", "codex"},
			{"set", "an-item", "archived"},
			{"set", "an-item", "archived", "--actor", "codex\nother"},
			{"set", "an-item", "archived", "--actor", "codex", "--expected-revision", "-1"},
			{"show", "an-item", "--format", "yaml"},
			{"show", "INVALID"},
			{"show", "an-item", "unexpected"},
		} {
			if output, err := executeOrganizingLifecycle(t, kind, services, args...); err == nil || output != "" {
				t.Fatalf("invalid args %q accepted: %q %v", args, output, err)
			}
		}
		if calls != 2 {
			t.Fatalf("invalid input called services: %d", calls)
		}
	}
}

func executeOrganizingLifecycle(t *testing.T, kind string, services workspaceLifecycleServices, args ...string) (string, error) {
	t.Helper()
	command := newWorkspaceLifecycleCommandWithServices(kind, services)
	var out, errors bytes.Buffer
	command.SetOut(&out)
	command.SetErr(&errors)
	command.SetArgs(args)
	command.SilenceUsage, command.SilenceErrors = true, true
	err := command.Execute()
	return out.String(), err
}
