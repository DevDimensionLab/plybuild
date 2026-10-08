package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
)

const taskContentReaderCapability = "problem@1,spec@1,goal@2,execution@2,assessment@1,selection@1,delivery@1,delivery@2"

// Draft publication remains strict invalid-input validation. Stored artifacts
// additionally distinguish a reader's unsupported semantics from corrupt bytes.
func decodeStoredTaskContent(raw []byte, operation string, published, internal bool) (canonicaljson.Object, error) {
	value, err := decodeTaskContent(raw, operation, published, internal)
	if err == nil {
		return value, nil
	}
	v, decodeErr := canonicaljson.DecodeStrict(raw)
	if decodeErr != nil {
		return nil, err
	}
	m := contentFields(v)
	schema := taskContentSchema(operation, published, internal)
	version := contentInt(m, "schema_version")
	if operation == "spec_record" && version == 2 {
		if isTaskGoalSpec(v) {
			schema = taskGoalSpecSchema(published)
		} else {
			schema = taskExecutionSpecSchema(published)
		}
	}
	if operation == "spec_record" {
		if _, ok := m["delivery"]; ok {
			schema["delivery"] = deliveryAgreementRule
		}
	}
	unknown := []string{}
	for key := range m {
		if _, supported := schema[key]; !supported {
			unknown = append(unknown, key)
		}
	}
	var unknownFields *contentUnknownFields
	if errors.As(err, &unknownFields) && len(unknown) == 0 {
		unknown = append(unknown, unknownFields.names...)
	}
	sort.Strings(unknown)
	unsupportedVersion := version > 1 && !(operation == "spec_record" && version == 2)
	if delivery := contentFields(m["delivery"]); contentInt(delivery, "schema_version") > 2 {
		unsupportedVersion = true
	}
	if len(unknown) > 0 || unsupportedVersion {
		detail := fmt.Sprintf("required %s schema %d exceeds reader capability %s", contentString(m, "kind"), version, taskContentReaderCapability)
		if len(unknown) > 0 {
			detail += "; unsupported fields: " + strings.Join(unknown, ", ")
		}
		return nil, &TaskContentError{Code: "task_content_reader_incompatible", Detail: detail, Schema: contentString(m, "kind"), Err: err}
	}
	return nil, err
}

func contentArtifactError(root, kind, digest, operation string, err error) error {
	if err == nil {
		return nil
	}
	code := "task_content_observation_unknown"
	var cause *TaskContentError
	schema := ""
	if errors.As(err, &cause) {
		code = cause.Code
		schema = cause.Schema
	}
	if errors.Is(err, fs.ErrNotExist) {
		code = "task_content_missing"
	}
	artifact := taskContentPath(root, kind, digest)
	next := "Inspect the required artifact and its preserved publication before retrying."
	if code == "task_content_reader_incompatible" {
		next = "Use a compatible installed Ply and the native active-run continuation; preserve the original control and publication."
	}
	return &TaskContentError{Code: code, Detail: fmt.Sprintf("%s required %s %s: %v; %s", operation, kind, digest, err, next), Operation: operation, Artifact: artifact, Schema: schema, ReaderCapability: taskContentReaderCapability, NextAction: next, Err: err}
}

// TaskContentDiagnosticValue projects only classified artifact diagnostics.
// Draft-input errors keep their established command behavior. No artifact bytes
// or unknown field values are copied into this readback.
func TaskContentDiagnosticValue(err error, operation string) canonicaljson.Object {
	var e *TaskContentError
	if !errors.As(err, &e) || e.Artifact == "" {
		return nil
	}
	cause := ""
	if e.Err != nil {
		cause = e.Err.Error()
	}
	return contentObject(map[string]canonicaljson.Value{
		"kind": "WorkspaceTaskContentDiagnostic@1", "schema_version": int64(1),
		"operation": operation, "code": e.Code, "detail": e.Detail,
		"artifact": e.Artifact, "artifact_schema": e.Schema,
		"reader_capability": e.ReaderCapability, "cause": cause, "next_action": e.NextAction,
	})
}
