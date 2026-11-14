package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/agentstatus"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func runAgentStatus(t *testing.T, cwd string, args ...string) (string, error) {
	t.Helper()
	t.Chdir(cwd)
	command := newAgentStatusCommand(workspace.SystemDependencies())
	command.SilenceUsage, command.SilenceErrors = true, true // as RootCmd under ExecuteE
	var stdout bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(args)
	err := command.Execute()
	return stdout.String(), err
}

func TestAgentStatusSetListClearPerProject(t *testing.T) {
	root := newAssignmentWorkspace(t)
	set, err := runAgentStatus(t, root, "set", "--project", "a", "--actor", "claude:aa24", "--pane", "w7:pCY",
		"--state", "waiting_background", "--ref", "PLAN-36/arch", "--for", "30m", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	if keys := planKeys(t, set); !reflect.DeepEqual(keys, []string{"kind", "schema_version", "project", "item"}) {
		t.Fatalf("set keys %v", keys)
	}
	var readback struct {
		Kind string           `json:"kind"`
		Item agentstatus.Item `json:"item"`
	}
	if err := json.Unmarshal([]byte(set), &readback); err != nil || readback.Kind != "PlyAgentStatusReadback@1" ||
		readback.Item.Actor != "claude:aa24" || readback.Item.Pane != "w7:pCY" || readback.Item.Ref != "PLAN-36/arch" || readback.Item.Expired {
		t.Fatalf("set %s %v", set, err)
	}

	listed, err := runAgentStatus(t, root, "list", "--project", "a", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	if keys := planKeys(t, listed); !reflect.DeepEqual(keys, []string{"kind", "schema_version", "project", "observed_at", "items", "invalid"}) {
		t.Fatalf("list keys %v", keys)
	}
	var list struct {
		Kind  string             `json:"kind"`
		Items []agentstatus.Item `json:"items"`
	}
	if err := json.Unmarshal([]byte(listed), &list); err != nil || list.Kind != "PlyAgentStatusList@1" || len(list.Items) != 1 || list.Items[0] != readback.Item {
		t.Fatalf("list %s %v", listed, err)
	}
	text, err := runAgentStatus(t, root, "list", "--project", "a")
	if err != nil || !strings.Contains(text, "claude:aa24 waiting_background w7:pCY PLAN-36/arch until ") {
		t.Fatalf("text %q %v", text, err)
	}

	other, err := runAgentStatus(t, root, "list", "--project", "b", "--format", "json")
	if err != nil || !strings.Contains(other, `"items": []`) {
		t.Fatalf("other project %s %v", other, err)
	}

	cleared, err := runAgentStatus(t, root, "clear", "--project", "a", "--actor", "claude:aa24", "--format", "json")
	if err != nil || !strings.Contains(cleared, `"removed": true`) {
		t.Fatalf("clear %s %v", cleared, err)
	}
	if listed, _ = runAgentStatus(t, root, "list", "--project", "a", "--format", "json"); !strings.Contains(listed, `"items": []`) {
		t.Fatalf("after clear %s", listed)
	}
}

func TestAgentStatusRefusesInvalidInput(t *testing.T) {
	root := newAssignmentWorkspace(t)
	for _, args := range [][]string{
		{"set", "--project", "a", "--actor", "claude:a", "--pane", "w1:p1", "--state", "working"},
		{"set", "--project", "a", "--actor", "claude:a", "--pane", "w1:p1", "--state", "working", "--for", "soon"},
		{"set", "--project", "a", "--actor", "claude:a", "--pane", "w1:p1", "--state", "idle", "--for", "1m"},
		{"clear", "--project", "a"},
		{"list", "--project", "a", "--format", "yaml"},
	} {
		_, err := runAgentStatus(t, root, args...)
		var e *agentstatus.Error
		if !errors.As(err, &e) || e.Class != agentstatus.ClassInvalidArguments {
			t.Fatalf("%v: %v", args, err)
		}
	}
}
