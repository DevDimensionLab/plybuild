package cmd

import (
	"errors"

	"github.com/devdimensionlab/plybuild/internal/privatestore"
	"github.com/devdimensionlab/plybuild/internal/workflownotification"
)

// ExitCode maps command errors at the executable boundary.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	var notificationError *workflownotification.Error
	if errors.As(err, &notificationError) && notificationError != nil {
		return notificationError.Exit
	}
	var typed *privatestore.Error
	if errors.As(err, &typed) && typed != nil && typed.Exit >= 2 && typed.Exit <= 5 {
		return typed.Exit
	}
	return 1
}
