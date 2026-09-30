package workflowhandoff

import (
	"errors"
	"fmt"
)

type ErrorClass string

const (
	ErrorTaskSpec          ErrorClass = "task_spec_binding_conflict"
	ErrorInvalidArguments  ErrorClass = "workflow_handoff_invalid_arguments"
	ErrorWorkspaceNotFound ErrorClass = "workspace_not_found"
	ErrorWorkspaceConflict ErrorClass = "workflow_handoff_workspace_conflict"
	ErrorProjectNotFound   ErrorClass = "workflow_handoff_project_not_found"
	ErrorRepoNotFound      ErrorClass = "workflow_handoff_repo_not_found"
	ErrorTargetConflict    ErrorClass = "workflow_handoff_target_conflict"
	ErrorSchemaInvalid     ErrorClass = "workflow_handoff_schema_invalid"
	ErrorNotFound          ErrorClass = "workflow_handoff_not_found"
	ErrorStale             ErrorClass = "workflow_handoff_stale"
	ErrorStartRequired     ErrorClass = "workflow_handoff_start_required"
	ErrorCapabilityInvalid ErrorClass = "workflow_handoff_capability_invalid"
	ErrorConflict          ErrorClass = "workflow_handoff_conflict"
	ErrorPayloadTooLarge   ErrorClass = "workflow_handoff_payload_too_large"
	ErrorIO                ErrorClass = "workflow_handoff_io_error"
)

type Error struct {
	Class  ErrorClass
	Detail string
	Err    error
}

func (err *Error) Error() string {
	detail := err.Detail
	if detail == "" && err.Err != nil {
		detail = err.Err.Error()
	}
	return fmt.Sprintf("%s: %s", err.Class, detail)
}

func (err *Error) Unwrap() error { return err.Err }

func classified(class ErrorClass, detail string, err error) error {
	return &Error{Class: class, Detail: detail, Err: err}
}

func InvalidArguments(detail string) error { return classified(ErrorInvalidArguments, detail, nil) }

func IsClass(err error, class ErrorClass) bool {
	var typed *Error
	return errors.As(err, &typed) && typed.Class == class
}

func classOf(err error) ErrorClass {
	var typed *Error
	if errors.As(err, &typed) {
		return typed.Class
	}
	return ErrorIO
}
