package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type viewCWD struct {
	workspace.FileSystem
	root string
}

func (f viewCWD) Getwd() (string, error) { return f.root, nil }

func viewFixture(t *testing.T) (workspace.Dependencies, string) {
	t.Helper()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	d := workspace.SystemDependencies()
	d.Files = viewCWD{d.Files, root}
	if _, err := workspace.Init(d); err != nil {
		t.Fatal(err)
	}
	return d, root
}

func runWorkspaceView(t *testing.T, d workspace.Dependencies, args ...string) (string, error) {
	t.Helper()
	c := newWorkspaceCommand(d)
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&bytes.Buffer{})
	c.SilenceErrors = true
	c.SilenceUsage = true
	c.SetArgs(args)
	err := c.Execute()
	return out.String(), err
}

func TestWorkspaceViewsRejectInvalidArgumentsBeforeDependencies(t *testing.T) {
	for _, args := range [][]string{
		{"status", "--format", "yaml"}, {"status", "extra"},
		{"attention", "--format", "yaml"}, {"attention", "--project", ""},
		{"journal", "recent", "--since", "tomorrow"}, {"journal", "recent", "--limit", "-1"},
		{"run", "list", "--project", "Bad"}, {"worktree", "list", "--repo", ""},
		{"task", "list", "--progress", "--ready"}, {"task", "list", "--progress", "--format", "yaml"},
		{"task", "show", "--progress", "Bad"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			out, err := runWorkspaceView(t, workspace.Dependencies{}, args...)
			if err == nil || out != "" {
				t.Fatalf("err=%v stdout=%q", err, out)
			}
			if strings.Contains(err.Error(), "dependencies") {
				t.Fatalf("accessed dependencies before validating input: %v", err)
			}
		})
	}
}

func TestWorkspaceViewsEmptyCollectionsAreJSONAndReadOnly(t *testing.T) {
	d, root := viewFixture(t)
	for _, tt := range []struct {
		args             []string
		kind, collection string
	}{
		{[]string{"task", "list", "--progress"}, "WorkspaceTaskProgressListReadback@1", "tasks"},
		{[]string{"epic", "list"}, "WorkspaceEpicListReadback@1", "epics"},
		{[]string{"attention"}, "WorkspaceAttentionReadback@1", "items"},
		{[]string{"overview"}, "WorkspaceOverviewReadback@1", "tasks"},
		{[]string{"journal", "recent"}, "WorkspaceActivityReadback@1", "events"},
		{[]string{"run", "list"}, "WorkspaceRunListReadback@1", "runs"},
		{[]string{"worktree", "list"}, "WorkspaceWorktreeListReadback@1", "worktrees"},
	} {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			out, err := runWorkspaceView(t, d, append(tt.args, "--format", "json")...)
			if err != nil {
				t.Fatal(err)
			}
			var v map[string]any
			if err := json.Unmarshal([]byte(out), &v); err != nil {
				t.Fatal(err)
			}
			if v["kind"] != tt.kind || v["schema_version"] != float64(1) {
				t.Fatalf("wrong envelope: %s", out)
			}
			rows, ok := v[tt.collection].([]any)
			if !ok || len(rows) != 0 {
				t.Fatalf("empty collection: %s", out)
			}
		})
	}
	entries, err := os.ReadDir(filepath.Join(root, ".ply"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "workspace.yaml" {
		t.Fatalf("reads created state: %v", entries)
	}
}

func TestWorkspaceStatusExposesChangeScope(t *testing.T) {
	d, _ := viewFixture(t)
	out, err := runWorkspaceView(t, d, "status", "--format", "json")
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatal(err)
	}
	if v["kind"] != "WorkspaceStatusReadback@1" || v["scope"] != "registered_ply_facts" || v["registry_sha256"] == nil {
		t.Fatalf("status: %s", out)
	}
}

func TestCapabilitiesRetainCoreAndAdvertiseReadExtensions(t *testing.T) {
	catalog := coreCapabilities(capabilityBuild{Version: "test"})
	encoded, err := json.Marshal(catalog)
	if err != nil {
		t.Fatal(err)
	}
	var v struct {
		Operations []capabilityOperation `json:"operations"`
		Extensions []capabilityOperation `json:"read_extensions"`
	}
	if err := json.Unmarshal(encoded, &v); err != nil {
		t.Fatal(err)
	}
	if len(v.Operations) != 5 {
		t.Fatal("legacy core catalog changed")
	}
	want := map[string]bool{"task.list:progress": false, "task.show:progress": false, "attention:default": false, "epic.list:default": false, "journal.recent:default": false, "run.list:default": false, "worktree.list:default": false, "status:default": false}
	for _, op := range v.Extensions {
		command, remaining, err := newWorkspaceCommand(workspace.Dependencies{}).Find(op.Command[1:])
		if err != nil || len(remaining) != 0 || command.Flags().Lookup("format") == nil {
			t.Fatalf("extension does not resolve: %+v, %v", op, err)
		}
		key := op.ID + ":" + op.Mode
		if _, ok := want[key]; ok {
			want[key] = true
		}
		if op.Effect != "read" {
			t.Fatalf("advertised a mutation: %+v", op)
		}
	}
	for id, present := range want {
		if !present {
			t.Fatalf("missing %s", id)
		}
	}
}
