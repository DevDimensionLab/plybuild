package workspace

import "fmt"

type WorkItemErrorClass string

const (
	ErrorWorkInvalidArguments                  WorkItemErrorClass = "workspace_work_invalid_arguments"
	ErrorWorkWorkspaceConflict                 WorkItemErrorClass = "workspace_work_workspace_conflict"
	ErrorWorkNotFound                          WorkItemErrorClass = "workspace_work_not_found"
	ErrorWorkProjectConflict                   WorkItemErrorClass = "workspace_work_project_conflict"
	ErrorWorkIdentityConflict                  WorkItemErrorClass = "workspace_work_identity_conflict"
	ErrorWorkGitObservation                    WorkItemErrorClass = "workspace_work_git_observation_error"
	ErrorWorkGitEffect                         WorkItemErrorClass = "workspace_work_git_effect_error"
	ErrorWorkParentDirty                       WorkItemErrorClass = "workspace_work_parent_dirty"
	ErrorWorkParentStale                       WorkItemErrorClass = "workspace_work_parent_stale"
	ErrorWorkPathConflict                      WorkItemErrorClass = "workspace_work_path_conflict"
	ErrorWorkRefConflict                       WorkItemErrorClass = "workspace_work_ref_conflict"
	ErrorWorkReconciliation                    WorkItemErrorClass = "workspace_work_reconciliation_required"
	ErrorWorkStoreConflict                     WorkItemErrorClass = "workspace_work_store_conflict"
	ErrorWorkIO                                WorkItemErrorClass = "workspace_work_io_error"
	ErrorTaskResultSchemaInvalid               WorkItemErrorClass = "workspace_task_result_schema_invalid"
	ErrorTaskResultEvidenceConflict            WorkItemErrorClass = "workspace_task_result_evidence_conflict"
	ErrorTaskResultConflict                    WorkItemErrorClass = "workspace_task_result_conflict"
	ErrorTaskQASchemaInvalid                   WorkItemErrorClass = "workspace_task_qa_schema_invalid"
	ErrorTaskQAEvidenceConflict                WorkItemErrorClass = "workspace_task_qa_evidence_conflict"
	ErrorTaskQAConflict                        WorkItemErrorClass = "workspace_task_qa_conflict"
	ErrorTaskIntegrationBlocked                WorkItemErrorClass = "workspace_task_integration_blocked"
	ErrorTaskIntegrationConflict               WorkItemErrorClass = "workspace_task_integration_conflict"
	ErrorTaskIntegrationConfirmationMismatch   WorkItemErrorClass = "workspace_task_integration_confirmation_mismatch"
	ErrorTaskIntegrationReconciliationRequired WorkItemErrorClass = "workspace_task_integration_reconciliation_required"
	ErrorTaskIntegrationUnsupportedRepository  WorkItemErrorClass = "workspace_task_integration_unsupported_repository"
)

type WorkItemError struct {
	Class  WorkItemErrorClass
	Detail string
	Err    error
}

func (err *WorkItemError) Error() string {
	detail := err.Detail
	if detail == "" && err.Err != nil {
		detail = err.Err.Error()
	}
	return fmt.Sprintf("%s: %s", err.Class, detail)
}
func (err *WorkItemError) Unwrap() error { return err.Err }

func WorkError(class WorkItemErrorClass, detail string, err error) error {
	return &WorkItemError{Class: class, Detail: detail, Err: err}
}

func WorkInvalidArguments(detail string) error {
	return WorkError(ErrorWorkInvalidArguments, detail, nil)
}
