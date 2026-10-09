package cmd

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/devdimensionlab/plybuild/internal/workspaceview"
	"github.com/spf13/cobra"
)

type capabilityBuild struct {
	Version     string  `json:"version"`
	VCSRevision *string `json:"vcs_revision"`
	VCSModified *bool   `json:"vcs_modified"`
}
type capabilitySchema struct {
	Kind          string `json:"kind"`
	SchemaVersion *int   `json:"schema_version"`
}
type capabilityOperation struct {
	ID            string             `json:"id"`
	Command       []string           `json:"command"`
	Mode          string             `json:"mode"`
	Selectors     []string           `json:"selectors"`
	Formats       []string           `json:"formats"`
	Filters       []string           `json:"filters"`
	FilterPolicy  string             `json:"filter_policy"`
	ResultSchemas []capabilitySchema `json:"result_schemas"`
	Effect        string             `json:"effect"`
}
type capabilityCatalog struct {
	Kind               string                `json:"kind"`
	SchemaVersion      int                   `json:"schema_version"`
	Build              capabilityBuild       `json:"build"`
	Coverage           string                `json:"coverage"`
	Operations         []capabilityOperation `json:"operations"`
	ReadExtensions     []capabilityOperation `json:"read_extensions"`
	WorkflowExtensions []capabilityOperation `json:"workflow_extensions,omitempty"`
}

func embeddedCapabilityBuild(info *debug.BuildInfo) capabilityBuild {
	build := capabilityBuild{Version: version}
	if info == nil {
		return build
	}
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if setting.Value != "" {
				value := setting.Value
				build.VCSRevision = &value
			}
		case "vcs.modified":
			if value, err := strconv.ParseBool(setting.Value); err == nil {
				build.VCSModified = &value
			}
		}
	}
	return build
}

// The original, closed core catalog stays stable. New read modes are advertised
// additively in read_extensions, so existing @1 consumers keep their exact core
// IDs, modes and coverage while newer consumers can discover richer projections.
func coreCapabilities(build capabilityBuild) capabilityCatalog {
	one, three := workspace.CoreReadSchemaVersion, 3
	operation := func(id, mode, kind string, selectors, filters []string, policy string) capabilityOperation {
		return capabilityOperation{id, append([]string{"workspace"}, strings.Split(id, ".")...), mode, selectors, []string{"text", "json"}, filters, policy, []capabilitySchema{{kind, &one}}, "read"}
	}
	ops := []capabilityOperation{
		operation("project.list", "default", workspace.ProjectListReadbackKind, []string{}, []string{}, "none"),
		operation("project.show", "default", workspace.ProjectReadbackKind, []string{}, []string{}, "none"),
		operation("task.list", "ready", "WorkspaceTaskQueueReadback@1", []string{"--ready"}, []string{"project", "repo", "epic"}, "all_or_none"),
		operation("task.list", "registered", workspace.TaskListReadbackKind, []string{}, []string{"project", "repo", "epic"}, "independent_and"),
		operation("task.show", "default", "WorkspaceTaskReadback@1", []string{}, []string{}, "none"),
	}
	ops[4].ResultSchemas = []capabilitySchema{{"WorkspaceTaskReadback@1", nil}, {"WorkspaceTaskIntegrationReadback@1", &one}, {"WorkspaceTaskReadback@3", &three}}
	extensions := []capabilityOperation{
		operation("attention", "default", "WorkspaceAttentionReadback@1", []string{}, []string{"project", "repo", "epic"}, "independent_and"),
		operation("epic.list", "default", "WorkspaceEpicListReadback@1", []string{}, []string{"project"}, "independent_and"),
		operation("epic.lifecycle.show", "default", "WorkspaceWorkItemLifecycleReadback@1", []string{}, []string{}, "none"),
		operation("journal.recent", "default", "WorkspaceActivityReadback@1", []string{}, []string{"project", "since", "limit"}, "independent_and"),
		operation("overview", "default", "WorkspaceOverviewReadback@1", []string{}, []string{"project"}, "independent_and"),
		operation("project.metadata.show", "default", "WorkspaceProjectMetadataReadback@1", []string{}, []string{}, "none"),
		operation("run.list", "default", "WorkspaceRunListReadback@1", []string{}, []string{"project", "active"}, "independent_and"),
		operation("status", "default", "WorkspaceStatusReadback@1", []string{}, []string{}, "none"),
		operation("task.lifecycle.show", "default", "WorkspaceWorkItemLifecycleReadback@1", []string{}, []string{}, "none"),
		operation("task.list", "progress", "WorkspaceTaskProgressListReadback@1", []string{"--progress"}, []string{"project", "repo", "epic"}, "independent_and"),
		operation("task.show", "progress", "WorkspaceTaskProgressReadback@1", []string{"--progress"}, []string{}, "none"),
		operation("worktree.list", "default", "WorkspaceWorktreeListReadback@1", []string{}, []string{"project", "repo"}, "independent_and"),
	}
	extensions = append(extensions, capabilityOperation{
		ID: "workflow.status", Command: []string{"workflow", "status"}, Mode: "default",
		Selectors: []string{"--all"}, Formats: []string{"text", "json"},
		Filters: []string{"project", "repo", "epic"}, FilterPolicy: "independent_and",
		ResultSchemas: []capabilitySchema{{workspaceview.WorkflowStatusKind, &one}}, Effect: "read",
	})
	extensions = append(extensions, capabilityOperation{
		ID: "workflow.trace", Command: []string{"workflow", "trace"}, Mode: "default",
		Selectors: []string{"--details"}, Formats: []string{"text", "json"},
		Filters: []string{}, FilterPolicy: "none",
		ResultSchemas: []capabilitySchema{{"WorkflowTraceReadback@1", &one}}, Effect: "read",
	})
	sort.Slice(extensions, func(i, j int) bool {
		if extensions[i].ID != extensions[j].ID {
			return extensions[i].ID < extensions[j].ID
		}
		return extensions[i].Mode < extensions[j].Mode
	})
	two := 2
	workflow := []capabilityOperation{
		{ID: "workflow.execute.verify", Command: []string{"workflow", "execute", "verify"}, Mode: "candidate_acceptance", Selectors: []string{"--context", "--review", "--candidate-binary", "--reuse"}, Formats: []string{"text", "json"}, Filters: []string{}, FilterPolicy: "none", ResultSchemas: []capabilitySchema{{"ply.workflow.run", &two}, {"PlyDeliveryVerification@2", &two}}, Effect: "candidate_test_and_evidence"},
		{ID: "workflow.execute.integrate", Command: []string{"workflow", "execute", "integrate"}, Mode: "frozen_acceptance", Selectors: []string{"--context"}, Formats: []string{"text", "json"}, Filters: []string{}, FilterPolicy: "none", ResultSchemas: []capabilitySchema{{"ply.workflow.run", &two}}, Effect: "authorized_local_integration"},
	}
	sort.Slice(workflow, func(i, j int) bool { return workflow[i].ID < workflow[j].ID })
	return capabilityCatalog{"PlyCapabilities@1", workspace.CoreReadSchemaVersion, build, "workspace-core-read", ops, extensions, workflow}
}

func newCapabilitiesCommand(build func() capabilityBuild) *cobra.Command {
	var format string
	invalid := func(detail string) error { return fmt.Errorf("capabilities_invalid_arguments: %s", detail) }
	command := &cobra.Command{
		Use: "capabilities", Short: "Show supported workspace core read contracts",
		Long:    "Show the stable workspace-core-read catalog and additive read_extensions for workflow status and trace, progress, attention, lifecycle, Epics, activity, runs, worktrees and change digests. Each operation advertises its result kind and decoder version. This is not the entire CLI. No workspace, profile, credentials or network access is required. Entries describe read contracts, not permission to start agents or change workflow state.",
		Example: "  ply capabilities\n  ply capabilities --format json",
		// Do not inherit the root's profile initializer, even with global logging flags.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return invalid(fmt.Sprintf("expected no positional arguments, got %d", len(args)))
			}
			if format != "text" && format != "json" {
				return invalid("--format must be text or json")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			catalog := coreCapabilities(build())
			if format == "json" {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(catalog)
			}
			revision, modified := "unknown", "unknown"
			if catalog.Build.VCSRevision != nil {
				revision = *catalog.Build.VCSRevision
			}
			if catalog.Build.VCSModified != nil {
				modified = strconv.FormatBool(*catalog.Build.VCSModified)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Build version: %s\nVCS revision: %s\nVCS modified: %s\nCoverage: %s (project list/show and Task list/show)\nOperations:\n", catalog.Build.Version, revision, modified, catalog.Coverage); err != nil {
				return err
			}
			for _, op := range catalog.Operations {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  %s (%s): %s; formats: %s\n", op.ID, op.Mode, strings.Join(op.Command, " "), strings.Join(op.Formats, ", ")); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Read extensions:"); err != nil {
				return err
			}
			for _, op := range catalog.ReadExtensions {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  %s (%s): %s; formats: %s\n", op.ID, op.Mode, strings.Join(op.Command, " "), strings.Join(op.Formats, ", ")); err != nil {
					return err
				}
			}
			for _, op := range catalog.WorkflowExtensions {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  %s (%s): %s; effect: %s\n", op.ID, op.Mode, strings.Join(op.Command, " "), op.Effect); err != nil {
					return err
				}
			}
			return nil
		},
	}
	command.Flags().StringVar(&format, "format", "text", "output format (text or json)")
	command.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error { return invalid(err.Error()) })
	return command
}

func init() {
	RootCmd.AddCommand(newCapabilitiesCommand(func() capabilityBuild {
		info, _ := debug.ReadBuildInfo()
		return embeddedCapabilityBuild(info)
	}))
}
