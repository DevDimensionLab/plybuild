package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func TestCapabilitiesSkipsRootInitializationAndUsesEmbeddedMetadata(t *testing.T) {
	for _, info := range []*debug.BuildInfo{nil, {Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "abc123"}, {Key: "vcs.modified", Value: "false"}}}} {
		for _, format := range []string{"text", "json"} {
			initialized := false
			root := &cobra.Command{Use: "ply", PersistentPreRunE: func(*cobra.Command, []string) error { initialized = true; return errors.New("initializer called") }}
			root.PersistentFlags().Bool("json", false, "logging")
			root.AddCommand(newCapabilitiesCommand(func() capabilityBuild { return embeddedCapabilityBuild(info) }))
			var out, stderr bytes.Buffer
			root.SetOut(&out)
			root.SetErr(&stderr)
			root.SetArgs([]string{"capabilities", "--format", format, "--json"})
			if err := root.Execute(); err != nil || initialized || stderr.Len() != 0 {
				t.Fatalf("err=%v initialized=%v stderr=%s", err, initialized, &stderr)
			}
			if format == "json" {
				var value capabilityCatalog
				if err := json.Unmarshal(out.Bytes(), &value); err != nil {
					t.Fatal(err)
				}
				if value.Kind != "PlyCapabilities@1" || value.Coverage != "workspace-core-read" || len(value.Operations) != 5 {
					t.Fatalf("catalog=%+v", value)
				}
				if info == nil && (value.Build.VCSModified != nil || value.Build.VCSRevision != nil) {
					t.Fatalf("guessed metadata: %+v", value.Build)
				}
				if info != nil && (value.Build.VCSModified == nil || *value.Build.VCSModified || value.Build.VCSRevision == nil || *value.Build.VCSRevision != "abc123") {
					t.Fatalf("lost metadata: %+v", value.Build)
				}
				// Each advertised operation must actually resolve, expose its formats and selectors,
				// and list only flags supported by the dispatched command.
				tree := newWorkspaceCommand(workspace.Dependencies{})
				for _, op := range value.Operations {
					command, remaining, err := tree.Find(op.Command[1:])
					if err != nil || len(remaining) != 0 {
						t.Fatalf("dispatch %v: %v %v", op.Command, remaining, err)
					}
					if command.Flags().Lookup("format") == nil || !reflect.DeepEqual(op.Formats, []string{"text", "json"}) {
						t.Fatalf("formats %s", op.ID)
					}
					for _, flag := range append(append([]string{}, op.Filters...), op.Selectors...) {
						if command.Flags().Lookup(strings.TrimPrefix(flag, "--")) == nil {
							t.Fatalf("missing flag %s on %s", flag, op.ID)
						}
					}
				}
			} else if info == nil && (!strings.Contains(out.String(), "VCS revision: unknown\n") || !strings.Contains(out.String(), "VCS modified: unknown\n")) {
				t.Fatalf("text=%s", &out)
			}
		}
	}
}

func TestCapabilitiesInvalidArgumentsAndWriterFailure(t *testing.T) {
	for _, args := range [][]string{{"extra"}, {"--format", ""}, {"--format", "xml"}, {"--unknown"}, {"--format"}} {
		calls := 0
		c := newCapabilitiesCommand(func() capabilityBuild { calls++; return capabilityBuild{} })
		var out bytes.Buffer
		c.SetOut(&out)
		c.SetErr(&bytes.Buffer{})
		c.SetArgs(args)
		c.SilenceErrors = true
		c.SilenceUsage = true
		if err := c.Execute(); err == nil || !strings.HasPrefix(err.Error(), "capabilities_invalid_arguments:") || out.Len() != 0 || calls != 0 {
			t.Fatalf("args=%v err=%v out=%s calls=%d", args, err, &out, calls)
		}
	}
	for _, format := range []string{"text", "json"} {
		c := newCapabilitiesCommand(func() capabilityBuild { return capabilityBuild{Version: "test"} })
		c.SetOut(coreFailWriter{})
		c.SetErr(&bytes.Buffer{})
		c.SetArgs([]string{"--format", format})
		c.SilenceErrors = true
		c.SilenceUsage = true
		if err := c.Execute(); err == nil {
			t.Fatal("ignored stdout failure")
		}
	}
}

type coreFailWriter struct{}

func (coreFailWriter) Write([]byte) (int, error) { return 0, errors.New("stdout failed") }

func TestRegisteredTaskListValidatesBeforeServiceAndAcceptsIndependentFilters(t *testing.T) {
	for _, args := range [][]string{{"list", "extra"}, {"list", "--format", "yaml"}, {"list", "--epic", ""}, {"list", "--project", ""}, {"list", "--repo", ""}, {"list", "--repo", "Bad"}, {"list", "--project", "ok", "--epic", "BAD"}} {
		calls := 0
		out, err := executeTaskCommand(t, workspaceTaskServices{list: func(workspace.TaskListFilters) (workspace.TaskListResult, error) {
			calls++
			return workspace.TaskListResult{}, nil
		}}, args...)
		if err == nil || !strings.HasPrefix(err.Error(), "workspace_work_invalid_arguments:") || out != "" || calls != 0 {
			t.Fatalf("args=%v err=%v out=%q calls=%d", args, err, out, calls)
		}
	}
	for _, flags := range [][]string{{}, {"--project", "p"}, {"--repo", "r"}, {"--epic", "e"}, {"--project", "p", "--repo", "r", "--epic", "e"}} {
		for _, format := range []string{"text", "json"} {
			calls := 0
			services := workspaceTaskServices{list: func(f workspace.TaskListFilters) (workspace.TaskListResult, error) {
				calls++
				return workspace.TaskListResult{Workspace: "/fixture", Filters: f}, nil
			}}
			args := append([]string{"list", "--format", format}, flags...)
			out, err := executeTaskCommand(t, services, args...)
			if err != nil || calls != 1 {
				t.Fatalf("args=%v err=%v calls=%d", args, err, calls)
			}
			if format == "json" && (!strings.Contains(out, `"tasks":[]`) || !strings.Contains(out, `"kind":"WorkspaceTaskListReadback@1"`)) {
				t.Fatalf("out=%s", out)
			}
			if format == "text" {
				for i := 0; i < len(flags); i += 2 {
					if flags[i] == "--epic" && len(flags) == 2 {
						continue
					}
					if !strings.Contains(out, flags[i+1]) {
						t.Fatalf("scope absent: %s", out)
					}
				}
			}
		}
	}
}

func TestProjectReadFormatsValidateBeforeDependencies(t *testing.T) {
	for _, args := range [][]string{{"list", "--format", "xml"}, {"list", "extra"}, {"show", "p", "--format", ""}, {"show", "--format", "json"}, {"show", "Bad", "--format", "json"}} {
		calls := 0
		c := newWorkspaceProjectCommandWithServices(workspaceProjectServices{list: func() (workspace.ProjectListResult, error) { calls++; return workspace.ProjectListResult{}, nil }, show: func(workspace.ProjectID) (workspace.ProjectResult, error) {
			calls++
			return workspace.ProjectResult{}, nil
		}})
		var out bytes.Buffer
		c.SetOut(&out)
		c.SetErr(&bytes.Buffer{})
		c.SetArgs(args)
		c.SilenceErrors = true
		c.SilenceUsage = true
		err := c.Execute()
		if err == nil || !strings.HasPrefix(err.Error(), "workspace_project_invalid_arguments:") || out.Len() != 0 || calls != 0 {
			t.Fatalf("args=%v err=%v out=%s calls=%d", args, err, &out, calls)
		}
	}
}
