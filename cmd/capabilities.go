package cmd

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workspace"
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
	Kind          string                `json:"kind"`
	SchemaVersion int                   `json:"schema_version"`
	Build         capabilityBuild       `json:"build"`
	Coverage      string                `json:"coverage"`
	Operations    []capabilityOperation `json:"operations"`
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

// This intentionally catalogs only the four workspace core read operations.
// Its fixed IDs and modes are not a generic CLI or permission discovery API.
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
	return capabilityCatalog{"PlyCapabilities@1", workspace.CoreReadSchemaVersion, build, "workspace-core-read", ops}
}

func newCapabilitiesCommand(build func() capabilityBuild) *cobra.Command {
	var format string
	invalid := func(detail string) error { return fmt.Errorf("capabilities_invalid_arguments: %s", detail) }
	command := &cobra.Command{
		Use: "capabilities", Short: "Show supported workspace core read contracts",
		Long:    "Show the built-in catalog for project list/show and Task list/show, including the separate ready mode. Coverage is workspace-core-read, not the entire CLI. No workspace, profile, credentials or network access is required. Entries describe read contracts, not permission to start agents or change workflow state.",
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
