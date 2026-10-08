package cmd

import (
	"fmt"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
	"github.com/spf13/cobra"
)

func taskContentCommandError(cmd *cobra.Command, format string, err error) error {
	if format != "json" || err == nil {
		return err
	}
	value := workspace.TaskContentDiagnosticValue(err, cmd.CommandPath())
	if value == nil {
		return err
	}
	raw, marshalErr := canonicaljson.Marshal(value)
	if marshalErr != nil {
		return marshalErr
	}
	if _, writeErr := cmd.OutOrStdout().Write(append(raw, '\n')); writeErr != nil {
		return fmt.Errorf("%w; original operation failed: %v", writeErr, err)
	}
	return err
}
