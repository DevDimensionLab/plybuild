package workflowhandoff

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

var (
	handoffIDPattern = regexp.MustCompile(`^hnd_[0-9a-f]{32}$`)
	keyPattern       = regexp.MustCompile(`^[a-z0-9][a-z0-9._/-]*$`)
	digestPattern    = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	oidPattern       = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)
)

func ParseHandoffID(value string) (HandoffID, error) {
	if !handoffIDPattern.MatchString(value) {
		return "", InvalidArguments(fmt.Sprintf("handoff ID %q must be hnd_ followed by 32 lowercase hexadecimal characters", value))
	}
	return HandoffID(value), nil
}

func validateKey(name, value string) error {
	if len(value) < 1 || len(value) > 128 || !keyPattern.MatchString(value) || strings.Contains(value, "//") || strings.HasSuffix(value, "/") {
		return fmt.Errorf("%s must be 1..128 lowercase ASCII key characters", name)
	}
	for _, segment := range strings.Split(value, "/") {
		if segment == ".." {
			return fmt.Errorf("%s contains a parent segment", name)
		}
	}
	return nil
}

func validateRepoPath(value string) error {
	if value == "" || value == "." || strings.ContainsAny(value, "\\\x00") || strings.HasPrefix(value, "/") || path.Clean(value) != value {
		return fmt.Errorf("path %q must be a clean repository-relative slash path", value)
	}
	for _, segment := range strings.Split(value, "/") {
		if segment == ".." {
			return fmt.Errorf("path %q contains a parent segment", value)
		}
	}
	return nil
}

func validatePlainText(name, value string, minimum, maximum int) error {
	if !utf8.ValidString(value) || utf8.RuneCountInString(value) < minimum || utf8.RuneCountInString(value) > maximum || strings.TrimSpace(value) != value {
		return fmt.Errorf("%s must contain %d..%d UTF-8 codepoints without surrounding whitespace", name, minimum, maximum)
	}
	for _, character := range value {
		if unicode.IsControl(character) {
			return fmt.Errorf("%s contains a control character", name)
		}
	}
	return nil
}

func exactObject(value canonicaljson.Value, context string, names ...string) (map[string]canonicaljson.Value, error) {
	object, ok := value.(canonicaljson.Object)
	if !ok {
		return nil, fmt.Errorf("%s must be an object", context)
	}
	expected := make(map[string]struct{}, len(names))
	for _, name := range names {
		expected[name] = struct{}{}
	}
	if len(object) != len(expected) {
		return nil, fmt.Errorf("%s has missing or unknown fields", context)
	}
	fields := make(map[string]canonicaljson.Value, len(object))
	for _, member := range object {
		if _, known := expected[member.Name]; !known {
			return nil, fmt.Errorf("%s has unknown field %q", context, member.Name)
		}
		if _, duplicate := fields[member.Name]; duplicate {
			return nil, fmt.Errorf("%s has duplicate field %q", context, member.Name)
		}
		fields[member.Name] = member.Value
	}
	for _, name := range names {
		if _, found := fields[name]; !found {
			return nil, fmt.Errorf("%s is missing field %q", context, name)
		}
	}
	return fields, nil
}

func capabilityProof(secret string, payload []byte) (string, error) {
	if len(secret) != 64 {
		return "", fmt.Errorf("capability secret must contain 64 lowercase hexadecimal characters")
	}
	key, err := hex.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("capability secret: %w", err)
	}
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func stringField(fields map[string]canonicaljson.Value, name, context string) (string, error) {
	value, ok := fields[name].(string)
	if !ok {
		return "", fmt.Errorf("%s.%s must be a string", context, name)
	}
	return value, nil
}

func intField(fields map[string]canonicaljson.Value, name, context string) (int64, error) {
	value, ok := fields[name].(int64)
	if !ok {
		return 0, fmt.Errorf("%s.%s must be an integer", context, name)
	}
	return value, nil
}

func arrayField(fields map[string]canonicaljson.Value, name, context string) ([]canonicaljson.Value, error) {
	value, ok := fields[name].([]canonicaljson.Value)
	if !ok {
		return nil, fmt.Errorf("%s.%s must be an array", context, name)
	}
	return value, nil
}

func validateDigest(value string) bool { return digestPattern.MatchString(value) }
func validateOID(value string) bool    { return oidPattern.MatchString(value) }

const (
	maxHandoffBytes = 256 << 10
	maxStartBytes   = 128 << 10
	maxResultBytes  = 1 << 20
)

func decodeHandoffDraftV1(input []byte) (handoffDraft, error) {
	return decodeHandoffDraftWithBudget(input, false)
}

func decodeHandoffDraftWithBudget(input []byte, delivery bool) (handoffDraft, error) {
	if len(input) > maxHandoffBytes {
		return handoffDraft{}, classified(ErrorPayloadTooLarge, "handoff draft exceeds 256 KiB", nil)
	}
	value, err := canonicaljson.DecodeStrict(input)
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	fields, err := exactObject(value, "handoff draft",
		"kind", "schema_version", "format", "format_version", "canonicalization",
		"publication_key", "activity_key", "goal", "recipient", "binding_request", "inputs", "authority",
		"budget", "procedure", "verifiers", "stop_conditions", "reporting")
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	if err := validateEnvelope(fields, "ply.workflow.handoff-draft", "handoff draft"); err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	publicationKey, err := stringField(fields, "publication_key", "handoff draft")
	if err != nil || validateKey("publication_key", publicationKey) != nil {
		return handoffDraft{}, schemaError("handoff draft", firstError(err, validateKey("publication_key", publicationKey)))
	}
	activityKey, err := stringField(fields, "activity_key", "handoff draft")
	if err != nil || validateKey("activity_key", activityKey) != nil {
		return handoffDraft{}, schemaError("handoff draft", firstError(err, validateKey("activity_key", activityKey)))
	}
	goal, goalTitle, err := validateGoal(fields["goal"])
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	recipient, principalID, err := validateRecipient(fields["recipient"])
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	binding, err := validateBindingRequest(fields["binding_request"])
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	inputs, err := validateInputs(fields["inputs"])
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	procedure, procedureIDs, err := validateProcedure(fields["procedure"])
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	verifiers, verifierIDs, err := validateVerifiers(fields["verifiers"])
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	authority, err := validateHandoffAuthority(fields["authority"], procedureIDs, verifierIDs, delivery)
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	budget, maxRounds, err := validateHandoffBudget(fields["budget"], delivery)
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	stopConditions, err := validateStopConditions(fields["stop_conditions"])
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	reporting, err := validateReporting(fields["reporting"])
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	canonical, err := canonicaljson.Marshal(value)
	if err != nil {
		return handoffDraft{}, schemaError("handoff draft", err)
	}
	digest := digestBytes(canonical)
	return handoffDraft{
		Value: value.(canonicaljson.Object), Goal: goal, Recipient: recipient, Authority: authority, Budget: budget,
		Reporting: reporting, Binding: binding, Inputs: inputs, Procedure: procedure, Verifiers: verifiers,
		StopConditions: stopConditions, PublicationKey: publicationKey, ActivityKey: activityKey,
		GoalTitle: goalTitle, Canonical: canonical, Digest: digest, MaxRounds: maxRounds, PrincipalID: principalID,
	}, nil
}

func validateStoredHandoffV1(value canonicaljson.Object) error {
	return validateStoredHandoffWithBudget(value, false)
}

func validateStoredHandoffWithBudget(value canonicaljson.Object, delivery bool) error {
	fields, err := exactObject(value, "handoff", "kind", "schema_version", "format", "format_version", "canonicalization", "identity", "source_draft_sha256", "goal", "recipient", "workspace_binding", "project_binding", "target_binding", "inputs", "authority", "budget", "procedure", "verifiers", "stop_conditions", "reporting", "reply_capability")
	if err != nil {
		return err
	}
	if err := validateEnvelope(fields, "ply.workflow.handoff", "handoff"); err != nil {
		return err
	}
	identity, err := exactObject(fields["identity"], "identity", "activity_id", "run_id", "handoff_id", "start_receipt_id", "terminal_result_id", "publication_key", "activity_key", "created_at_utc")
	if err != nil {
		return err
	}
	patterns := map[string]string{"activity_id": `^act_[0-9a-f]{32}$`, "run_id": `^run_[0-9a-f]{32}$`, "handoff_id": `^hnd_[0-9a-f]{32}$`, "start_receipt_id": `^rcp_[0-9a-f]{32}$`, "terminal_result_id": `^res_[0-9a-f]{32}$`}
	for name, pattern := range patterns {
		text, err := stringField(identity, name, "identity")
		if err != nil || !regexp.MustCompile(pattern).MatchString(text) {
			return fmt.Errorf("invalid identity.%s", name)
		}
	}
	for _, name := range []string{"publication_key", "activity_key"} {
		text, err := stringField(identity, name, "identity")
		if err != nil || validateKey("identity."+name, text) != nil {
			return fmt.Errorf("invalid identity.%s", name)
		}
	}
	created, _ := stringField(identity, "created_at_utc", "identity")
	if !validateUTC(created) {
		return fmt.Errorf("invalid identity.created_at_utc")
	}
	source, err := stringField(fields, "source_draft_sha256", "handoff")
	if err != nil || !validateDigest(source) {
		return fmt.Errorf("invalid source_draft_sha256")
	}
	if _, _, err := validateGoal(fields["goal"]); err != nil {
		return err
	}
	if _, _, err := validateRecipient(fields["recipient"]); err != nil {
		return err
	}
	workspaceFields, err := exactObject(fields["workspace_binding"], "workspace_binding", "root", "marker_format_version", "marker_sha256")
	if err != nil {
		return err
	}
	if _, err := absoluteCleanString(workspaceFields, "root", "workspace_binding"); err != nil {
		return err
	}
	version, err := intField(workspaceFields, "marker_format_version", "workspace_binding")
	if err != nil || version != 1 {
		return fmt.Errorf("invalid workspace marker version")
	}
	marker, _ := stringField(workspaceFields, "marker_sha256", "workspace_binding")
	if !validateDigest(marker) {
		return fmt.Errorf("invalid workspace marker digest")
	}
	projectFields, err := exactObject(fields["project_binding"], "project_binding", "project_id", "repo_id", "registered_locator", "registered_git_common_dir")
	if err != nil {
		return err
	}
	for _, name := range []string{"project_id", "repo_id"} {
		text, err := stringField(projectFields, name, "project_binding")
		if err != nil || validateKey(name, text) != nil {
			return fmt.Errorf("invalid project binding")
		}
	}
	for _, name := range []string{"registered_locator", "registered_git_common_dir"} {
		if _, err := absoluteCleanString(projectFields, name, "project_binding"); err != nil {
			return err
		}
	}
	if err := validateTargetBinding(fields["target_binding"]); err != nil {
		return err
	}
	if _, err := validateInputs(fields["inputs"]); err != nil {
		return err
	}
	procedure, procedureIDs, err := validateProcedure(fields["procedure"])
	if err != nil {
		return err
	}
	_ = procedure
	verifiers, verifierIDs, err := validateVerifiers(fields["verifiers"])
	if err != nil {
		return err
	}
	_ = verifiers
	if _, err := validateHandoffAuthority(fields["authority"], procedureIDs, verifierIDs, delivery); err != nil {
		return err
	}
	if _, _, err := validateHandoffBudget(fields["budget"], delivery); err != nil {
		return err
	}
	if _, err := validateStopConditions(fields["stop_conditions"]); err != nil {
		return err
	}
	if _, err := validateReporting(fields["reporting"]); err != nil {
		return err
	}
	reply, err := exactObject(fields["reply_capability"], "reply_capability", "capability_id", "secret", "reply_root", "allowed_operations")
	if err != nil {
		return err
	}
	capability, _ := stringField(reply, "capability_id", "reply_capability")
	secret, _ := stringField(reply, "secret", "reply_capability")
	if !regexp.MustCompile(`^cap_[0-9a-f]{32}$`).MatchString(capability) || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(secret) {
		return fmt.Errorf("invalid reply capability")
	}
	if _, err := absoluteCleanString(reply, "reply_root", "reply_capability"); err != nil {
		return err
	}
	operations, err := stringArray(reply["allowed_operations"], "reply_capability.allowed_operations")
	if err != nil || !equalStrings(operations, []string{"submit-result", "submit-start"}) {
		return fmt.Errorf("invalid reply capability operations")
	}
	return nil
}

func validateTargetBinding(value canonicaljson.Value) error {
	fields, err := exactObject(value, "target_binding", "worktree", "git_common_dir", "ref", "oid", "tree", "status_policy")
	if err != nil {
		return err
	}
	for _, name := range []string{"worktree", "git_common_dir"} {
		if _, err := absoluteCleanString(fields, name, "target_binding"); err != nil {
			return err
		}
	}
	ref, err := stringField(fields, "ref", "target_binding")
	if err != nil || !strings.HasPrefix(ref, "refs/heads/") || strings.Contains(ref, "..") {
		return fmt.Errorf("invalid target binding ref")
	}
	for _, name := range []string{"oid", "tree"} {
		oid, err := stringField(fields, name, "target_binding")
		if err != nil || !validateOID(oid) {
			return fmt.Errorf("invalid target binding %s", name)
		}
	}
	_, err = validateStatusPolicy(fields["status_policy"])
	return err
}

func validateEnvelope(fields map[string]canonicaljson.Value, kind, context string) error {
	actualKind, err := stringField(fields, "kind", context)
	if err != nil || actualKind != kind {
		return fmt.Errorf("%s.kind must be %q", context, kind)
	}
	format, err := stringField(fields, "format", context)
	if err != nil || format != "json" {
		return fmt.Errorf("%s.format must be json", context)
	}
	canonicalization, err := stringField(fields, "canonicalization", context)
	if err != nil || canonicalization != "RFC8785" {
		return fmt.Errorf("%s.canonicalization must be RFC8785", context)
	}
	for _, name := range []string{"schema_version", "format_version"} {
		version, err := intField(fields, name, context)
		if err != nil || version != 1 {
			return fmt.Errorf("%s.%s must be 1", context, name)
		}
	}
	return nil
}

func validateGoal(value canonicaljson.Value) (canonicaljson.Object, string, error) {
	fields, err := exactObject(value, "goal", "title", "recipient_role", "objective", "done_when")
	if err != nil {
		return nil, "", err
	}
	limits := map[string]int{"title": 120, "recipient_role": 120, "objective": 8000, "done_when": 8000}
	for _, name := range []string{"title", "recipient_role", "objective", "done_when"} {
		text, err := stringField(fields, name, "goal")
		if err != nil || validatePlainText("goal."+name, text, 1, limits[name]) != nil {
			return nil, "", firstError(err, validatePlainText("goal."+name, text, 1, limits[name]))
		}
	}
	return value.(canonicaljson.Object), fields["title"].(string), nil
}

func validateRecipient(value canonicaljson.Value) (canonicaljson.Object, string, error) {
	fields, err := exactObject(value, "recipient", "principal_id", "principal_kind", "runtime_constraints")
	if err != nil {
		return nil, "", err
	}
	principalID, err := stringField(fields, "principal_id", "recipient")
	if err != nil || validateKey("recipient.principal_id", principalID) != nil {
		return nil, "", firstError(err, validateKey("recipient.principal_id", principalID))
	}
	kind, err := stringField(fields, "principal_kind", "recipient")
	if err != nil || kind != "human_started_agent" {
		return nil, "", fmt.Errorf("recipient.principal_kind must be human_started_agent")
	}
	constraints, err := sortedStringArray(fields["runtime_constraints"], "recipient.runtime_constraints", false)
	if err != nil {
		return nil, "", err
	}
	for _, constraint := range constraints {
		if err := validatePlainText("runtime constraint", constraint, 1, 2000); err != nil {
			return nil, "", err
		}
	}
	return value.(canonicaljson.Object), principalID, nil
}

func validateBindingRequest(value canonicaljson.Value) (bindingRequest, error) {
	fields, err := exactObject(value, "binding_request", "project_id", "repo_id", "target_worktree", "target_ref", "expected_oid", "status_policy")
	if err != nil {
		return bindingRequest{}, err
	}
	projectID, err := stringField(fields, "project_id", "binding_request")
	if err != nil || validateKey("project_id", projectID) != nil {
		return bindingRequest{}, firstError(err, validateKey("project_id", projectID))
	}
	repoID, err := stringField(fields, "repo_id", "binding_request")
	if err != nil || validateKey("repo_id", repoID) != nil {
		return bindingRequest{}, firstError(err, validateKey("repo_id", repoID))
	}
	worktree, err := absoluteCleanString(fields, "target_worktree", "binding_request")
	if err != nil {
		return bindingRequest{}, err
	}
	ref, err := stringField(fields, "target_ref", "binding_request")
	if err != nil || !strings.HasPrefix(ref, "refs/heads/") || strings.Contains(ref, "..") {
		return bindingRequest{}, fmt.Errorf("binding_request.target_ref must be a full branch ref")
	}
	oid, err := stringField(fields, "expected_oid", "binding_request")
	if err != nil || !validateOID(oid) {
		return bindingRequest{}, fmt.Errorf("binding_request.expected_oid must be a Git OID")
	}
	status, err := validateStatusPolicy(fields["status_policy"])
	if err != nil {
		return bindingRequest{}, err
	}
	return bindingRequest{ProjectID: workspaceProjectID(projectID), RepoID: workspaceRepoID(repoID), TargetWorktree: worktree, TargetRef: ref, ExpectedOID: oid, Status: status}, nil
}

// These conversions keep workspace identifier validation at the workspace adapter boundary.
func workspaceProjectID(value string) workspace.ProjectID { return workspace.ProjectID(value) }
func workspaceRepoID(value string) workspace.RepoID       { return workspace.RepoID(value) }

func validateStatusPolicy(value canonicaljson.Value) (StatusPolicy, error) {
	object, ok := value.(canonicaljson.Object)
	if !ok {
		return StatusPolicy{}, fmt.Errorf("status_policy must be an object")
	}
	modeValue, found := objectMember(object, "mode")
	mode, ok := modeValue.(string)
	if !found || !ok {
		return StatusPolicy{}, fmt.Errorf("status_policy.mode must be a string")
	}
	if mode == "clean" {
		if _, err := exactObject(object, "status_policy", "mode"); err != nil {
			return StatusPolicy{}, err
		}
		return StatusPolicy{Mode: mode, Value: object}, nil
	}
	if mode != "exact_manifest" {
		return StatusPolicy{}, fmt.Errorf("unknown status_policy mode %q", mode)
	}
	fields, err := exactObject(object, "status_policy", "mode", "entries")
	if err != nil {
		return StatusPolicy{}, err
	}
	values, err := arrayField(fields, "entries", "status_policy")
	if err != nil || len(values) == 0 {
		return StatusPolicy{}, fmt.Errorf("exact_manifest entries must be a non-empty array")
	}
	entries := make([]StatusEntry, 0, len(values))
	last := ""
	states := setOf("unmodified", "modified", "added", "deleted", "type_changed", "untracked")
	for _, entryValue := range values {
		entryFields, err := exactObject(entryValue, "status entry", "path", "index_state", "worktree_state", "kind", "sha256")
		if err != nil {
			return StatusPolicy{}, err
		}
		entryPath, err := stringField(entryFields, "path", "status entry")
		if err != nil || validateRepoPath(entryPath) != nil || (last != "" && entryPath <= last) {
			return StatusPolicy{}, fmt.Errorf("status entries must have unique sorted clean paths")
		}
		last = entryPath
		indexState, err := stringField(entryFields, "index_state", "status entry")
		if err != nil || !states[indexState] {
			return StatusPolicy{}, fmt.Errorf("invalid status index_state")
		}
		worktreeState, err := stringField(entryFields, "worktree_state", "status entry")
		if err != nil || !states[worktreeState] {
			return StatusPolicy{}, fmt.Errorf("invalid status worktree_state")
		}
		kind, err := stringField(entryFields, "kind", "status entry")
		if err != nil || (kind != "regular" && kind != "symlink") {
			return StatusPolicy{}, fmt.Errorf("invalid status kind")
		}
		digest := ""
		if entryFields["sha256"] != nil {
			digest, err = stringField(entryFields, "sha256", "status entry")
			if err != nil || !validateDigest(digest) {
				return StatusPolicy{}, fmt.Errorf("invalid status digest")
			}
		}
		deleted := indexState == "deleted" || worktreeState == "deleted"
		if deleted != (entryFields["sha256"] == nil) {
			return StatusPolicy{}, fmt.Errorf("deleted status entries alone use null sha256")
		}
		entries = append(entries, StatusEntry{Path: entryPath, IndexState: indexState, WorktreeState: worktreeState, Kind: kind, SHA256: digest})
	}
	return StatusPolicy{Mode: mode, Entries: entries, Value: object}, nil
}

func validateInputs(value canonicaljson.Value) ([]inputSpec, error) {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return nil, fmt.Errorf("inputs must be an array")
	}
	inputs := make([]inputSpec, 0, len(values))
	last := ""
	roles := setOf("spec", "design", "handoff", "repository_instruction", "verification_fixture", "other")
	media := regexp.MustCompile(`^[\x21-\x7e]+/[\x21-\x7e]+$`)
	for _, inputValue := range values {
		fields, err := exactObject(inputValue, "input", "id", "role", "locator", "sha256", "size_bytes", "media_type", "git_binding")
		if err != nil {
			return nil, err
		}
		id, err := stringField(fields, "id", "input")
		if err != nil || validateKey("input.id", id) != nil || (last != "" && id <= last) {
			return nil, fmt.Errorf("inputs must have unique sorted valid IDs")
		}
		last = id
		role, err := stringField(fields, "role", "input")
		if err != nil || !roles[role] {
			return nil, fmt.Errorf("invalid input role")
		}
		locator, err := absoluteCleanString(fields, "locator", "input")
		if err != nil {
			return nil, err
		}
		digest, err := stringField(fields, "sha256", "input")
		if err != nil || !validateDigest(digest) {
			return nil, fmt.Errorf("invalid input digest")
		}
		size, err := intField(fields, "size_bytes", "input")
		if err != nil || size < 0 || size > 256<<20 {
			return nil, fmt.Errorf("input size_bytes is out of range")
		}
		mediaType, err := stringField(fields, "media_type", "input")
		if err != nil || len(mediaType) > 127 || !media.MatchString(mediaType) {
			return nil, fmt.Errorf("invalid input media_type")
		}
		var binding *GitBinding
		if fields["git_binding"] != nil {
			bindingFields, err := exactObject(fields["git_binding"], "input.git_binding", "git_common_dir", "ref", "oid", "blob")
			if err != nil {
				return nil, err
			}
			common, err := absoluteCleanString(bindingFields, "git_common_dir", "input.git_binding")
			if err != nil {
				return nil, err
			}
			ref, _ := stringField(bindingFields, "ref", "input.git_binding")
			oid, _ := stringField(bindingFields, "oid", "input.git_binding")
			blob, _ := stringField(bindingFields, "blob", "input.git_binding")
			if !strings.HasPrefix(ref, "refs/") || !validateOID(oid) || !validateOID(blob) || len(oid) != len(blob) {
				return nil, fmt.Errorf("invalid input Git binding")
			}
			binding = &GitBinding{GitCommonDir: common, Ref: ref, OID: oid, Blob: blob}
		}
		inputs = append(inputs, inputSpec{ID: id, Role: role, Locator: locator, SHA256: digest, SizeBytes: size, MediaType: mediaType, GitBinding: binding})
	}
	return inputs, nil
}

func validateProcedure(value canonicaljson.Value) ([]canonicaljson.Value, map[string]bool, error) {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return nil, nil, fmt.Errorf("procedure must be an array")
	}
	seen := map[string]bool{}
	for _, item := range values {
		fields, err := exactObject(item, "procedure entry", "id", "instruction", "required_before")
		if err != nil {
			return nil, nil, err
		}
		id, _ := stringField(fields, "id", "procedure entry")
		instruction, _ := stringField(fields, "instruction", "procedure entry")
		if validateKey("procedure.id", id) != nil || seen[id] || validatePlainText("procedure.instruction", instruction, 1, 4000) != nil {
			return nil, nil, fmt.Errorf("invalid or duplicate procedure entry")
		}
		required, err := sortedStringArray(fields["required_before"], "procedure.required_before", true)
		if err != nil {
			return nil, nil, err
		}
		for _, dependency := range required {
			if !seen[dependency] {
				return nil, nil, fmt.Errorf("procedure dependency %q is not earlier", dependency)
			}
		}
		seen[id] = true
	}
	return values, seen, nil
}

func validateVerifiers(value canonicaljson.Value) ([]canonicaljson.Value, map[string]bool, error) {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return nil, nil, fmt.Errorf("verifiers must be an array")
	}
	seen := map[string]bool{}
	last := ""
	for _, item := range values {
		fields, err := exactObject(item, "verifier", "id", "argv", "cwd", "env", "expected_exit", "stop_on_failure", "evidence")
		if err != nil {
			return nil, nil, err
		}
		id, _ := stringField(fields, "id", "verifier")
		if validateKey("verifier.id", id) != nil || seen[id] || (last != "" && id <= last) {
			return nil, nil, fmt.Errorf("verifier IDs must be valid and sorted")
		}
		last = id
		argv, err := stringArray(fields["argv"], "verifier.argv")
		if err != nil || len(argv) < 1 || len(argv) > 64 || argv[0] == "" {
			return nil, nil, fmt.Errorf("verifier argv must contain 1..64 plain-text elements")
		}
		for _, argument := range argv {
			if validatePlainText("verifier argument", argument, 1, 2000) != nil {
				return nil, nil, fmt.Errorf("invalid verifier argument")
			}
		}
		cwd, err := absoluteCleanString(fields, "cwd", "verifier")
		if err != nil || cwd == "" {
			return nil, nil, fmt.Errorf("invalid verifier cwd")
		}
		if err := validateEnv(fields["env"]); err != nil {
			return nil, nil, err
		}
		exit, err := intField(fields, "expected_exit", "verifier")
		if err != nil || exit < 0 || exit > 255 {
			return nil, nil, fmt.Errorf("verifier expected_exit is out of range")
		}
		if _, ok := fields["stop_on_failure"].(bool); !ok {
			return nil, nil, fmt.Errorf("verifier stop_on_failure must be boolean")
		}
		evidence, err := exactObject(fields["evidence"], "verifier.evidence", "capture_stdout", "capture_stderr", "classification", "binding")
		if err != nil {
			return nil, nil, err
		}
		if _, ok := evidence["capture_stdout"].(bool); !ok {
			return nil, nil, fmt.Errorf("capture_stdout must be boolean")
		}
		if _, ok := evidence["capture_stderr"].(bool); !ok {
			return nil, nil, fmt.Errorf("capture_stderr must be boolean")
		}
		classification, _ := stringField(evidence, "classification", "verifier.evidence")
		binding, _ := stringField(evidence, "binding", "verifier.evidence")
		if classification != "workspace_internal" || (binding != "target_oid" && binding != "sha256") {
			return nil, nil, fmt.Errorf("invalid verifier evidence")
		}
		seen[id] = true
	}
	return values, seen, nil
}

func validateEnv(value canonicaljson.Value) error {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return fmt.Errorf("verifier.env must be an array")
	}
	namePattern := regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)
	secretPattern := regexp.MustCompile(`TOKEN|SECRET|PASSWORD|KEY|CREDENTIAL`)
	last := ""
	for _, item := range values {
		object, ok := item.(canonicaljson.Object)
		if !ok {
			return fmt.Errorf("verifier env entry must be an object")
		}
		nameValue, _ := objectMember(object, "name")
		name, ok := nameValue.(string)
		if !ok || !namePattern.MatchString(name) || (last != "" && name <= last) {
			return fmt.Errorf("verifier env names must be valid and sorted")
		}
		last = name
		sourceValue, _ := objectMember(object, "source")
		source, _ := sourceValue.(string)
		if source == "literal" {
			fields, err := exactObject(object, "literal env", "name", "source", "value")
			if err != nil {
				return err
			}
			if secretPattern.MatchString(name) {
				return fmt.Errorf("secret-like env must be inherited")
			}
			if _, err := stringField(fields, "value", "literal env"); err != nil {
				return err
			}
		} else if source == "environment" {
			fields, err := exactObject(object, "inherited env", "name", "source", "variable")
			if err != nil {
				return err
			}
			variable, _ := stringField(fields, "variable", "inherited env")
			if !namePattern.MatchString(variable) {
				return fmt.Errorf("invalid inherited env variable")
			}
		} else {
			return fmt.Errorf("unknown env source")
		}
	}
	return nil
}

func validateAuthority(value canonicaljson.Value, procedureIDs, verifierIDs map[string]bool) (canonicaljson.Object, error) {
	return validateHandoffAuthority(value, procedureIDs, verifierIDs, false)
}

func validateHandoffAuthority(value canonicaljson.Value, procedureIDs, verifierIDs map[string]bool, delivery bool) (canonicaljson.Object, error) {
	fields, err := exactObject(value, "authority", "allowed_effects", "forbidden_effects", "human_gates")
	if err != nil {
		return nil, err
	}
	allowed, err := arrayField(fields, "allowed_effects", "authority")
	if err != nil {
		return nil, err
	}
	effectTypes := setOf("filesystem_write", "git_ref_write", "git_index_write", "git_commit", "command_execute", "network", "remote_provider", "process_start", "install", "cleanup", "human_qa", "integration", "push", "pull_request", "merge", "release", "deploy")
	seenIDs := map[string]bool{}
	seenTypes := map[string]bool{}
	for index, effect := range allowed {
		ef, err := exactObject(effect, "allowed effect", "id", "type", "scope", "max_occurrences", "sequence")
		if err != nil {
			return nil, err
		}
		id, _ := stringField(ef, "id", "allowed effect")
		typeName, _ := stringField(ef, "type", "allowed effect")
		max, _ := intField(ef, "max_occurrences", "allowed effect")
		sequence, _ := intField(ef, "sequence", "allowed effect")
		unlimited := delivery && ef["max_occurrences"] == nil && setOf("filesystem_write", "git_ref_write", "git_index_write", "git_commit", "command_execute", "process_start", "install")[typeName]
		if validateKey("effect.id", id) != nil || seenIDs[id] || !effectTypes[typeName] || (!unlimited && (max < 1 || max > 99)) || sequence != int64(index+1) {
			return nil, fmt.Errorf("invalid allowed effect")
		}
		scope, _ := ef["scope"].(canonicaljson.Object)
		wholeTask := delivery && typeName == "filesystem_write" && objectString(scope, "kind") == "task_worktree"
		if wholeTask {
			if _, err := exactObject(scope, "Task worktree scope", "kind"); err != nil {
				return nil, err
			}
		} else if delivery && objectString(scope, "kind") == "delivery" {
			if err := validateDeliveryEffectScope(typeName, scope); err != nil {
				return nil, err
			}
		} else if err := validateEffectScope(typeName, ef["scope"], procedureIDs, verifierIDs); err != nil {
			return nil, err
		}
		seenIDs[id], seenTypes[typeName] = true, true
	}
	forbidden, err := arrayField(fields, "forbidden_effects", "authority")
	if err != nil {
		return nil, err
	}
	forbiddenTypes := map[string]bool{}
	for _, item := range forbidden {
		ff, err := exactObject(item, "forbidden effect", "type", "reason")
		if err != nil {
			return nil, err
		}
		typeName, _ := stringField(ff, "type", "forbidden effect")
		reason, _ := stringField(ff, "reason", "forbidden effect")
		if !effectTypes[typeName] || forbiddenTypes[typeName] || seenTypes[typeName] || validatePlainText("forbidden reason", reason, 1, 2000) != nil {
			return nil, fmt.Errorf("invalid or contradictory forbidden effect")
		}
		forbiddenTypes[typeName] = true
	}
	gates, err := sortedStringArray(fields["human_gates"], "human_gates", true)
	if err != nil {
		return nil, err
	}
	gateSet := setOf("human_task_qa", "local_integration", "workflow_delivery", "human_final_qa", "merge_decision", "post_merge_verification")
	for _, gate := range gates {
		if !gateSet[gate] {
			return nil, fmt.Errorf("invalid human gate")
		}
	}
	return value.(canonicaljson.Object), nil
}

func validateEffectScope(effectType string, value canonicaljson.Value, procedureIDs, verifierIDs map[string]bool) error {
	switch effectType {
	case "filesystem_write":
		fields, err := exactObject(value, "filesystem scope", "kind", "paths", "directory_prefixes")
		if err != nil {
			return err
		}
		kind, _ := stringField(fields, "kind", "filesystem scope")
		if kind != "filesystem" {
			return fmt.Errorf("filesystem scope kind mismatch")
		}
		paths, err := sortedStringArray(fields["paths"], "filesystem paths", true)
		if err != nil {
			return err
		}
		prefixes, ok := fields["directory_prefixes"].([]canonicaljson.Value)
		if !ok || len(paths)+len(prefixes) == 0 {
			return fmt.Errorf("filesystem scope is empty")
		}
		for _, value := range paths {
			if validateRepoPath(value) != nil {
				return fmt.Errorf("invalid filesystem path")
			}
		}
		last := ""
		for _, item := range prefixes {
			pf, err := exactObject(item, "directory prefix", "path", "recursive")
			if err != nil {
				return err
			}
			p, _ := stringField(pf, "path", "directory prefix")
			recursive, ok := pf["recursive"].(bool)
			if validateRepoPath(p) != nil || !ok || !recursive || (last != "" && p <= last) {
				return fmt.Errorf("invalid directory prefix")
			}
			last = p
		}
	case "git_ref_write", "git_index_write", "git_commit":
		fields, err := exactObject(value, "git scope", "kind", "ref", "expected_before_oid", "operation")
		if err != nil {
			return err
		}
		kind, _ := stringField(fields, "kind", "git scope")
		ref, _ := stringField(fields, "ref", "git scope")
		oid, _ := stringField(fields, "expected_before_oid", "git scope")
		operation, _ := stringField(fields, "operation", "git scope")
		expected := map[string]string{"git_ref_write": "update_ref", "git_index_write": "update_index", "git_commit": "create_commit"}[effectType]
		if kind != "git" || !strings.HasPrefix(ref, "refs/") || !validateOID(oid) || operation != expected {
			return fmt.Errorf("invalid Git effect scope")
		}
	case "command_execute":
		fields, err := exactObject(value, "command scope", "kind", "procedure_ids", "verifier_ids")
		if err != nil {
			return err
		}
		kind, _ := stringField(fields, "kind", "command scope")
		procedure, err := sortedStringArray(fields["procedure_ids"], "command procedure IDs", true)
		if err != nil {
			return err
		}
		verifiers, err := sortedStringArray(fields["verifier_ids"], "command verifier IDs", true)
		if err != nil || kind != "command" || len(procedure)+len(verifiers) == 0 {
			return fmt.Errorf("invalid command scope")
		}
		for _, id := range procedure {
			if !procedureIDs[id] {
				return fmt.Errorf("unknown procedure ID %q", id)
			}
		}
		for _, id := range verifiers {
			if !verifierIDs[id] {
				return fmt.Errorf("unknown verifier ID %q", id)
			}
		}
	default:
		fields, err := exactObject(value, "boolean scope", "kind", "allowed")
		if err != nil {
			return err
		}
		kind, _ := stringField(fields, "kind", "boolean scope")
		allowed, ok := fields["allowed"].(bool)
		if kind != "boolean" || !ok || !allowed {
			return fmt.Errorf("invalid boolean scope")
		}
	}
	return nil
}

func validateBudget(value canonicaljson.Value) (canonicaljson.Object, int64, error) {
	return validateHandoffBudget(value, false)
}

func validateHandoffBudget(value canonicaljson.Value, delivery bool) (canonicaljson.Object, int64, error) {
	fields, err := exactObject(value, "budget", "max_rounds", "round_definition")
	if err != nil {
		return nil, 0, err
	}
	max, err := intField(fields, "max_rounds", "budget")
	if delivery && fields["max_rounds"] == nil {
		max, err = 0, nil
	}
	if err != nil || (!delivery && (max < 1 || max > 99)) || (delivery && (max < 0 || max > 2147483647 || (max == 0 && fields["max_rounds"] != nil))) {
		return nil, 0, fmt.Errorf("budget.max_rounds must be 1..99")
	}
	rd, err := exactObject(fields["round_definition"], "round_definition", "unit", "command_retry_consumes_round", "retry_condition")
	if err != nil {
		return nil, 0, err
	}
	unit, _ := stringField(rd, "unit", "round_definition")
	retry, _ := stringField(rd, "retry_condition", "round_definition")
	consumes, ok := rd["command_retry_consumes_round"].(bool)
	wantUnit := "implementation_or_review_fix_iteration"
	if delivery {
		wantUnit = "candidate_verification_attempt"
	}
	if unit != wantUnit || !ok || consumes || retry != "only_if_no_effect_started" {
		return nil, 0, fmt.Errorf("invalid round_definition")
	}
	return value.(canonicaljson.Object), max, nil
}

func validateStopConditions(value canonicaljson.Value) ([]canonicaljson.Value, error) {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return nil, fmt.Errorf("stop_conditions must be an array")
	}
	allowed := setOf("product_decision_required", "scope_or_authority_expansion", "target_or_input_drift", "unknown_or_partial_effect", "unexpected_sensitive_data", "round_budget_exhausted", "repository_state_unexpected", "verifier_unrecoverable")
	required := setOf("product_decision_required", "scope_or_authority_expansion", "target_or_input_drift", "unknown_or_partial_effect", "unexpected_sensitive_data", "round_budget_exhausted")
	seen := map[string]bool{}
	for _, item := range values {
		fields, err := exactObject(item, "stop condition", "type", "description")
		if err != nil {
			return nil, err
		}
		typeName, _ := stringField(fields, "type", "stop condition")
		description, _ := stringField(fields, "description", "stop condition")
		if !allowed[typeName] || seen[typeName] || validatePlainText("stop description", description, 1, 2000) != nil {
			return nil, fmt.Errorf("invalid stop condition")
		}
		seen[typeName] = true
	}
	for typeName := range required {
		if !seen[typeName] {
			return nil, fmt.Errorf("missing required stop condition %s", typeName)
		}
	}
	return values, nil
}

func validateReporting(value canonicaljson.Value) (canonicaljson.Object, error) {
	fields, err := exactObject(value, "reporting", "summary_max_codepoints", "meaning_max_codepoints", "required_start_fields", "required_terminal_fields")
	if err != nil {
		return nil, err
	}
	summary, _ := intField(fields, "summary_max_codepoints", "reporting")
	meaning, _ := intField(fields, "meaning_max_codepoints", "reporting")
	if summary != 240 || meaning != 600 {
		return nil, fmt.Errorf("reporting length policy must be 240/600")
	}
	wantStart := []string{"acceptance", "binding", "contract_digests", "issues", "observed_inputs", "observed_project", "observed_target", "observed_workspace", "principal", "receipt_id", "sandbox"}
	wantTerminal := []string{"artifacts", "binding", "evidence_gaps", "final_target", "forbidden_effects_observed", "meaning", "observed_effects", "principal", "reported_outcome", "result_id", "review", "rounds_used", "start_binding", "stop_reasons", "summary", "verifier_results"}
	start, err := sortedStringArray(fields["required_start_fields"], "required_start_fields", true)
	if err != nil || !equalStrings(start, wantStart) {
		return nil, fmt.Errorf("required_start_fields does not match schema")
	}
	terminal, err := sortedStringArray(fields["required_terminal_fields"], "required_terminal_fields", true)
	if err != nil || !equalStrings(terminal, wantTerminal) {
		return nil, fmt.Errorf("required_terminal_fields does not match schema")
	}
	return value.(canonicaljson.Object), nil
}

func decodeStartDraftV1(input []byte) (startDraft, error) {
	if len(input) > maxStartBytes {
		return startDraft{}, classified(ErrorPayloadTooLarge, "start receipt draft exceeds 128 KiB", nil)
	}
	value, err := canonicaljson.DecodeStrict(input)
	if err != nil {
		return startDraft{}, schemaError("start receipt draft", err)
	}
	fields, err := exactObject(value, "start receipt draft", "kind", "schema_version", "format", "format_version", "canonicalization", "receipt_id", "binding", "principal", "observed_workspace", "observed_project", "observed_target", "observed_inputs", "contract_digests", "sandbox", "acceptance", "issues")
	if err != nil {
		return startDraft{}, schemaError("start receipt draft", err)
	}
	if err := validateEnvelope(fields, "ply.workflow.start-receipt-draft", "start receipt draft"); err != nil {
		return startDraft{}, schemaError("start receipt draft", err)
	}
	receiptID, _ := stringField(fields, "receipt_id", "start receipt draft")
	if !regexp.MustCompile(`^rcp_[0-9a-f]{32}$`).MatchString(receiptID) {
		return startDraft{}, schemaError("start receipt draft", fmt.Errorf("invalid receipt_id"))
	}
	if err := validateStartNested(fields); err != nil {
		return startDraft{}, schemaError("start receipt draft", err)
	}
	acceptance, _ := stringField(fields, "acceptance", "start receipt draft")
	if !setOf("started", "rejected", "conflict", "unknown")[acceptance] {
		return startDraft{}, schemaError("start receipt draft", fmt.Errorf("invalid acceptance"))
	}
	issues, _ := fields["issues"].([]canonicaljson.Value)
	if (acceptance == "started") != (len(issues) == 0) {
		return startDraft{}, schemaError("start receipt draft", fmt.Errorf("started alone has no issues"))
	}
	canonical, _ := canonicaljson.Marshal(value)
	return startDraft{Value: value.(canonicaljson.Object), Canonical: canonical, ReceiptID: receiptID, Acceptance: acceptance, Principal: fields["principal"].(canonicaljson.Object)}, nil
}

func validateStartNested(fields map[string]canonicaljson.Value) error {
	shapes := map[string][]string{
		"binding":            {"activity_id", "run_id", "handoff_id", "handoff_sha256"},
		"principal":          {"expected_principal_id", "human_start_principal", "start_surface", "session_id", "runtime_id", "model_id", "started_at_utc"},
		"observed_workspace": {"root", "marker_format_version", "marker_sha256", "matches_expected"},
		"observed_project":   {"project_id", "repo_id", "registered_locator", "registered_git_common_dir", "matches_expected"},
		"observed_target":    {"worktree", "git_common_dir", "ref", "oid", "tree", "status_policy", "matches_expected"},
		"contract_digests":   {"authority_sha256", "budget_sha256", "verifiers_sha256", "stop_conditions_sha256"},
		"sandbox":            {"read_roots", "write_roots", "temp_root", "matches_contract"},
	}
	for name, shape := range shapes {
		if _, err := exactObject(fields[name], name, shape...); err != nil {
			return err
		}
	}
	if err := validateReceiptBinding(fields["binding"]); err != nil {
		return err
	}
	if err := validatePrincipal(fields["principal"]); err != nil {
		return err
	}
	workspaceFields, _ := exactObject(fields["observed_workspace"], "observed_workspace", "root", "marker_format_version", "marker_sha256", "matches_expected")
	if _, err := absoluteCleanString(workspaceFields, "root", "observed_workspace"); err != nil {
		return err
	}
	if version, err := intField(workspaceFields, "marker_format_version", "observed_workspace"); err != nil || version != 1 {
		return fmt.Errorf("observed_workspace.marker_format_version must be 1")
	}
	if digest, err := stringField(workspaceFields, "marker_sha256", "observed_workspace"); err != nil || !validateDigest(digest) {
		return fmt.Errorf("invalid observed workspace digest")
	}
	if _, ok := workspaceFields["matches_expected"].(bool); !ok {
		return fmt.Errorf("observed_workspace.matches_expected must be boolean")
	}
	projectFields, _ := exactObject(fields["observed_project"], "observed_project", "project_id", "repo_id", "registered_locator", "registered_git_common_dir", "matches_expected")
	for _, name := range []string{"project_id", "repo_id"} {
		text, err := stringField(projectFields, name, "observed_project")
		if err != nil || validateKey(name, text) != nil {
			return fmt.Errorf("invalid observed project %s", name)
		}
	}
	for _, name := range []string{"registered_locator", "registered_git_common_dir"} {
		if _, err := absoluteCleanString(projectFields, name, "observed_project"); err != nil {
			return err
		}
	}
	if _, ok := projectFields["matches_expected"].(bool); !ok {
		return fmt.Errorf("observed_project.matches_expected must be boolean")
	}
	if err := validateObservedTarget(fields["observed_target"]); err != nil {
		return err
	}
	digests, _ := exactObject(fields["contract_digests"], "contract_digests", "authority_sha256", "budget_sha256", "verifiers_sha256", "stop_conditions_sha256")
	for _, name := range []string{"authority_sha256", "budget_sha256", "verifiers_sha256", "stop_conditions_sha256"} {
		digest, err := stringField(digests, name, "contract_digests")
		if err != nil || !validateDigest(digest) {
			return fmt.Errorf("invalid contract digest %s", name)
		}
	}
	if err := validateSandbox(fields["sandbox"]); err != nil {
		return err
	}
	inputs, ok := fields["observed_inputs"].([]canonicaljson.Value)
	if !ok {
		return fmt.Errorf("observed_inputs must be an array")
	}
	last := ""
	for _, item := range inputs {
		input, err := exactObject(item, "observed input", "id", "locator", "sha256", "size_bytes", "matches_expected")
		if err != nil {
			return err
		}
		id, _ := stringField(input, "id", "observed input")
		if validateKey("observed input id", id) != nil || (last != "" && id <= last) {
			return fmt.Errorf("observed inputs must be sorted")
		}
		last = id
		if _, err := absoluteCleanString(input, "locator", "observed input"); err != nil {
			return err
		}
		digest, err := stringField(input, "sha256", "observed input")
		if err != nil || !validateDigest(digest) {
			return fmt.Errorf("invalid observed input digest")
		}
		size, err := intField(input, "size_bytes", "observed input")
		if err != nil || size < 0 {
			return fmt.Errorf("invalid observed input size")
		}
		if _, ok := input["matches_expected"].(bool); !ok {
			return fmt.Errorf("observed input matches_expected must be boolean")
		}
	}
	issues, ok := fields["issues"].([]canonicaljson.Value)
	if !ok {
		return fmt.Errorf("issues must be an array")
	}
	lastIssue := ""
	for _, item := range issues {
		issue, err := exactObject(item, "issue", "type", "detail")
		if err != nil {
			return err
		}
		typeName, _ := stringField(issue, "type", "issue")
		detail, _ := stringField(issue, "detail", "issue")
		if !setOf("binding", "workspace", "project", "target", "input", "contract", "sandbox", "principal", "unknown")[typeName] || validatePlainText("issue detail", detail, 1, 2000) != nil {
			return fmt.Errorf("invalid issue")
		}
		combined := typeName + "\x00" + detail
		if lastIssue != "" && combined <= lastIssue {
			return fmt.Errorf("issues must be sorted")
		}
		lastIssue = combined
	}
	return nil
}

func decodeTerminalDraft(input []byte) (terminalDraft, error) {
	if len(input) > maxResultBytes {
		return terminalDraft{}, classified(ErrorPayloadTooLarge, "terminal result draft exceeds 1 MiB", nil)
	}
	value, err := canonicaljson.DecodeStrict(input)
	if err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	fields, err := exactObject(value, "terminal result draft", "kind", "schema_version", "format", "format_version", "canonicalization", "result_id", "binding", "start_binding", "principal", "reported_outcome", "rounds_used", "stop_reasons", "summary", "meaning", "final_target", "observed_effects", "verifier_results", "review", "artifacts", "evidence_gaps", "forbidden_effects_observed")
	if err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	if err := validateEnvelope(fields, "ply.workflow.terminal-result-draft", "terminal result draft"); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	resultID, _ := stringField(fields, "result_id", "terminal result draft")
	if !regexp.MustCompile(`^res_[0-9a-f]{32}$`).MatchString(resultID) {
		return terminalDraft{}, schemaError("terminal result draft", fmt.Errorf("invalid result_id"))
	}
	outcome, _ := stringField(fields, "reported_outcome", "terminal result draft")
	if !setOf("complete", "blocked", "conflict", "budget_exhausted", "unknown")[outcome] {
		return terminalDraft{}, schemaError("terminal result draft", fmt.Errorf("invalid reported_outcome"))
	}
	rounds, _ := intField(fields, "rounds_used", "terminal result draft")
	if rounds < 0 || rounds > 99 {
		return terminalDraft{}, schemaError("terminal result draft", fmt.Errorf("rounds_used is out of range"))
	}
	summary, _ := stringField(fields, "summary", "terminal result draft")
	meaning, _ := stringField(fields, "meaning", "terminal result draft")
	if validatePlainText("summary", summary, 1, 240) != nil || validatePlainText("meaning", meaning, 1, 600) != nil {
		return terminalDraft{}, schemaError("terminal result draft", fmt.Errorf("invalid summary or meaning"))
	}
	for name, shape := range map[string][]string{"binding": {"activity_id", "run_id", "handoff_id", "handoff_sha256"}, "start_binding": {"receipt_id", "start_receipt_sha256"}, "principal": {"expected_principal_id", "human_start_principal", "start_surface", "session_id", "runtime_id", "model_id", "started_at_utc"}, "final_target": {"worktree", "git_common_dir", "ref", "oid", "tree", "status_policy", "matches_expected"}, "review": {"findings", "fixes", "open_actionable_findings"}} {
		if _, err := exactObject(fields[name], name, shape...); err != nil {
			return terminalDraft{}, schemaError("terminal result draft", err)
		}
	}
	if err := validateReceiptBinding(fields["binding"]); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	startBinding, _ := exactObject(fields["start_binding"], "start_binding", "receipt_id", "start_receipt_sha256")
	receiptID, _ := stringField(startBinding, "receipt_id", "start_binding")
	digest, _ := stringField(startBinding, "start_receipt_sha256", "start_binding")
	if !regexp.MustCompile(`^rcp_[0-9a-f]{32}$`).MatchString(receiptID) || !validateDigest(digest) {
		return terminalDraft{}, schemaError("terminal result draft", fmt.Errorf("invalid start binding"))
	}
	if err := validatePrincipal(fields["principal"]); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	if err := validateObservedTarget(fields["final_target"]); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	stopReasons, ok := fields["stop_reasons"].([]canonicaljson.Value)
	if !ok {
		return terminalDraft{}, schemaError("terminal result draft", fmt.Errorf("stop_reasons must be an array"))
	}
	if (outcome == "complete") != (len(stopReasons) == 0) {
		return terminalDraft{}, schemaError("terminal result draft", fmt.Errorf("complete alone has no stop reasons"))
	}
	if err := validateStopReasons(stopReasons, outcome); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	if err := validateObservedEffects(fields["observed_effects"]); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	if err := validateVerifierResults(fields["verifier_results"]); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	if err := validateReview(fields["review"]); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	artifacts, ok := fields["artifacts"].([]canonicaljson.Value)
	if !ok {
		return terminalDraft{}, schemaError("terminal result draft", fmt.Errorf("artifacts must be an array"))
	}
	if err := validateArtifacts(artifacts); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	if err := validateGapEntries(fields["evidence_gaps"], "evidence_gaps"); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	if err := validateGapEntries(fields["forbidden_effects_observed"], "forbidden_effects_observed"); err != nil {
		return terminalDraft{}, schemaError("terminal result draft", err)
	}
	canonical, _ := canonicaljson.Marshal(value)
	return terminalDraft{Value: value.(canonicaljson.Object), Canonical: canonical, ResultID: resultID, Outcome: outcome, RoundsUsed: rounds, Summary: summary, Meaning: meaning, Principal: fields["principal"].(canonicaljson.Object), Artifacts: artifacts}, nil
}

func validateReceiptBinding(value canonicaljson.Value) error {
	fields, err := exactObject(value, "binding", "activity_id", "run_id", "handoff_id", "handoff_sha256")
	if err != nil {
		return err
	}
	patterns := map[string]string{"activity_id": `^act_[0-9a-f]{32}$`, "run_id": `^run_[0-9a-f]{32}$`, "handoff_id": `^hnd_[0-9a-f]{32}$`}
	for name, pattern := range patterns {
		text, err := stringField(fields, name, "binding")
		if err != nil || !regexp.MustCompile(pattern).MatchString(text) {
			return fmt.Errorf("invalid binding.%s", name)
		}
	}
	digest, err := stringField(fields, "handoff_sha256", "binding")
	if err != nil || !validateDigest(digest) {
		return fmt.Errorf("invalid binding.handoff_sha256")
	}
	return nil
}

func validatePrincipal(value canonicaljson.Value) error {
	fields, err := exactObject(value, "principal", "expected_principal_id", "human_start_principal", "start_surface", "session_id", "runtime_id", "model_id", "started_at_utc")
	if err != nil {
		return err
	}
	for _, name := range []string{"expected_principal_id", "human_start_principal", "start_surface", "session_id", "runtime_id", "model_id"} {
		text, err := stringField(fields, name, "principal")
		if err != nil || validatePlainText("principal."+name, text, 1, 256) != nil {
			return fmt.Errorf("invalid principal.%s", name)
		}
	}
	started, err := stringField(fields, "started_at_utc", "principal")
	if err != nil || !validateUTC(started) {
		return fmt.Errorf("invalid principal.started_at_utc")
	}
	return nil
}

func validateObservedTarget(value canonicaljson.Value) error {
	fields, err := exactObject(value, "target observation", "worktree", "git_common_dir", "ref", "oid", "tree", "status_policy", "matches_expected")
	if err != nil {
		return err
	}
	for _, name := range []string{"worktree", "git_common_dir"} {
		if _, err := absoluteCleanString(fields, name, "target observation"); err != nil {
			return err
		}
	}
	ref, err := stringField(fields, "ref", "target observation")
	if err != nil || !strings.HasPrefix(ref, "refs/") || strings.Contains(ref, "..") {
		return fmt.Errorf("invalid target observation ref")
	}
	for _, name := range []string{"oid", "tree"} {
		oid, err := stringField(fields, name, "target observation")
		if err != nil || !validateOID(oid) {
			return fmt.Errorf("invalid target observation %s", name)
		}
	}
	if _, err := validateStatusPolicy(fields["status_policy"]); err != nil {
		return err
	}
	if _, ok := fields["matches_expected"].(bool); !ok {
		return fmt.Errorf("target observation matches_expected must be boolean")
	}
	return nil
}

func validateSandbox(value canonicaljson.Value) error {
	fields, err := exactObject(value, "sandbox", "read_roots", "write_roots", "temp_root", "matches_contract")
	if err != nil {
		return err
	}
	for _, name := range []string{"read_roots", "write_roots"} {
		roots, err := sortedStringArray(fields[name], "sandbox."+name, false)
		if err != nil {
			return err
		}
		for _, root := range roots {
			if !filepath.IsAbs(root) || filepath.Clean(root) != root {
				return fmt.Errorf("sandbox.%s contains a non-absolute or unclean path", name)
			}
		}
	}
	if _, err := absoluteCleanString(fields, "temp_root", "sandbox"); err != nil {
		return err
	}
	if _, ok := fields["matches_contract"].(bool); !ok {
		return fmt.Errorf("sandbox.matches_contract must be boolean")
	}
	return nil
}

func validateStopReasons(values []canonicaljson.Value, outcome string) error {
	allowed := setOf("product_decision_required", "scope_or_authority_expansion", "target_or_input_drift", "unknown_or_partial_effect", "unexpected_sensitive_data", "round_budget_exhausted", "repository_state_unexpected", "verifier_unrecoverable")
	for index, value := range values {
		fields, err := exactObject(value, "stop reason", "type", "detail", "primary")
		if err != nil {
			return err
		}
		kind, _ := stringField(fields, "type", "stop reason")
		detail, _ := stringField(fields, "detail", "stop reason")
		primary, ok := fields["primary"].(bool)
		if !allowed[kind] || validatePlainText("stop reason detail", detail, 1, 2000) != nil || !ok || primary != (index == 0) {
			return fmt.Errorf("invalid stop reason")
		}
		if index == 0 && outcome == "budget_exhausted" && kind != "round_budget_exhausted" {
			return fmt.Errorf("budget_exhausted requires round_budget_exhausted")
		}
		if index == 0 && outcome == "unknown" && kind != "unknown_or_partial_effect" {
			return fmt.Errorf("unknown requires unknown_or_partial_effect")
		}
	}
	return nil
}

func validateObservedEffects(value canonicaljson.Value) error {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return fmt.Errorf("observed_effects must be an array")
	}
	seen := map[string]bool{}
	for _, value := range values {
		fields, err := exactObject(value, "observed effect", "effect_id", "type", "scope", "occurrences", "within_authority")
		if err != nil {
			return err
		}
		id, _ := stringField(fields, "effect_id", "observed effect")
		kind, _ := stringField(fields, "type", "observed effect")
		occurrences, err := intField(fields, "occurrences", "observed effect")
		_, scopeOK := fields["scope"].(canonicaljson.Object)
		_, authorityOK := fields["within_authority"].(bool)
		if validateKey("effect_id", id) != nil || seen[id] || kind == "" || err != nil || occurrences < 0 || occurrences > 99 || !scopeOK || !authorityOK {
			return fmt.Errorf("invalid observed effect")
		}
		seen[id] = true
	}
	return nil
}

func validateVerifierResults(value canonicaljson.Value) error {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return fmt.Errorf("verifier_results must be an array")
	}
	last := ""
	for _, value := range values {
		fields, err := exactObject(value, "verifier result", "verifier_id", "argv", "cwd", "exit", "bound_oid_or_sha256", "stdout_artifact_id", "stderr_artifact_id")
		if err != nil {
			return err
		}
		id, _ := stringField(fields, "verifier_id", "verifier result")
		if validateKey("verifier_id", id) != nil || (last != "" && id <= last) {
			return fmt.Errorf("verifier results must be ID-sorted and unique")
		}
		last = id
		argv, err := stringArray(fields["argv"], "verifier result argv")
		if err != nil || len(argv) == 0 {
			return fmt.Errorf("invalid verifier result argv")
		}
		for _, arg := range argv {
			if validatePlainText("verifier argument", arg, 1, 2000) != nil {
				return fmt.Errorf("invalid verifier argument")
			}
		}
		if _, err := absoluteCleanString(fields, "cwd", "verifier result"); err != nil {
			return err
		}
		exit, err := intField(fields, "exit", "verifier result")
		if err != nil || exit < 0 || exit > 255 {
			return fmt.Errorf("invalid verifier exit")
		}
		bound, err := stringField(fields, "bound_oid_or_sha256", "verifier result")
		if err != nil || (!validateOID(bound) && !validateDigest(bound)) {
			return fmt.Errorf("invalid verifier binding")
		}
		for _, name := range []string{"stdout_artifact_id", "stderr_artifact_id"} {
			if fields[name] != nil {
				artifactID, ok := fields[name].(string)
				if !ok || validateKey(name, artifactID) != nil {
					return fmt.Errorf("invalid %s", name)
				}
			}
		}
	}
	return nil
}

func validateReview(value canonicaljson.Value) error {
	fields, err := exactObject(value, "review", "findings", "fixes", "open_actionable_findings")
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, listName := range []string{"findings", "fixes", "open_actionable_findings"} {
		values, ok := fields[listName].([]canonicaljson.Value)
		if !ok {
			return fmt.Errorf("review.%s must be an array", listName)
		}
		last := ""
		for _, value := range values {
			item, err := exactObject(value, "review entry", "id", "severity", "summary", "evidence_ids")
			if err != nil {
				return err
			}
			id, _ := stringField(item, "id", "review entry")
			severity, _ := stringField(item, "severity", "review entry")
			summary, _ := stringField(item, "summary", "review entry")
			if validateKey("review id", id) != nil || seen[id] || (last != "" && id <= last) || !setOf("low", "medium", "high", "critical")[severity] || validatePlainText("review summary", summary, 1, 600) != nil {
				return fmt.Errorf("invalid review entry")
			}
			if _, err := sortedStringArray(item["evidence_ids"], "review evidence_ids", true); err != nil {
				return err
			}
			seen[id] = true
			last = id
		}
	}
	return nil
}

func validateArtifacts(values []canonicaljson.Value) error {
	last := ""
	for _, value := range values {
		object, ok := value.(canonicaljson.Object)
		if !ok {
			return fmt.Errorf("artifact must be an object")
		}
		id := objectString(object, "artifact_id")
		kind := objectString(object, "kind")
		if validateKey("artifact_id", id) != nil || (last != "" && id <= last) {
			return fmt.Errorf("artifacts must be ID-sorted and unique")
		}
		last = id
		if kind == "managed" {
			fields, err := exactObject(object, "managed artifact", "artifact_id", "kind", "description", "media_type", "classification", "size_bytes", "sha256", "locator")
			if err != nil {
				return err
			}
			if err := validateArtifactCommon(fields); err != nil {
				return err
			}
			classification, _ := stringField(fields, "classification", "managed artifact")
			if classification != "workspace_internal" {
				return fmt.Errorf("invalid managed artifact classification")
			}
			if _, err := absoluteCleanString(fields, "locator", "managed artifact"); err != nil {
				return err
			}
		} else if kind == "withheld" {
			fields, err := exactObject(object, "withheld artifact", "artifact_id", "kind", "description", "reason", "media_type", "size_bytes", "sha256")
			if err != nil {
				return err
			}
			description, _ := stringField(fields, "description", "withheld artifact")
			reason, _ := stringField(fields, "reason", "withheld artifact")
			size, err := intField(fields, "size_bytes", "withheld artifact")
			if validatePlainText("artifact description", description, 1, 600) != nil || !setOf("sensitive", "too_large")[reason] || err != nil || size < 0 {
				return fmt.Errorf("invalid withheld artifact")
			}
			for _, name := range []string{"media_type", "sha256"} {
				if fields[name] != nil {
					text, ok := fields[name].(string)
					if !ok || (name == "sha256" && !validateDigest(text)) || (name == "media_type" && text == "") {
						return fmt.Errorf("invalid withheld artifact %s", name)
					}
				}
			}
		} else {
			return fmt.Errorf("unknown artifact kind")
		}
	}
	return nil
}

func validateArtifactCommon(fields map[string]canonicaljson.Value) error {
	description, _ := stringField(fields, "description", "artifact")
	mediaType, _ := stringField(fields, "media_type", "artifact")
	size, err := intField(fields, "size_bytes", "artifact")
	digest, _ := stringField(fields, "sha256", "artifact")
	if validatePlainText("artifact description", description, 1, 600) != nil || mediaType == "" || err != nil || size < 0 || size > 64<<20 || !validateDigest(digest) {
		return fmt.Errorf("invalid artifact metadata")
	}
	return nil
}

func validateGapEntries(value canonicaljson.Value, context string) error {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return fmt.Errorf("%s must be an array", context)
	}
	last := ""
	for _, value := range values {
		fields, err := exactObject(value, context+" entry", "type", "detail", "artifact_ids", "effect_ids")
		if err != nil {
			return err
		}
		kind, _ := stringField(fields, "type", context)
		detail, _ := stringField(fields, "detail", context)
		if validatePlainText(context+" type", kind, 1, 120) != nil || validatePlainText(context+" detail", detail, 1, 2000) != nil {
			return fmt.Errorf("invalid %s entry", context)
		}
		combined := kind + "\x00" + detail
		if last != "" && combined <= last {
			return fmt.Errorf("%s must be sorted", context)
		}
		last = combined
		if _, err := sortedStringArray(fields["artifact_ids"], context+" artifact_ids", true); err != nil {
			return err
		}
		if _, err := sortedStringArray(fields["effect_ids"], context+" effect_ids", true); err != nil {
			return err
		}
	}
	return nil
}

func validateFinalDocument(value canonicaljson.Value, kind string, draftFields []string) (canonicaljson.Object, error) {
	fields := append(append([]string{}, draftFields...), "capability_proof")
	objectFields, err := exactObject(value, kind, fields...)
	if err != nil {
		return nil, err
	}
	if err := validateEnvelope(objectFields, kind, kind); err != nil {
		return nil, err
	}
	proof, err := exactObject(objectFields["capability_proof"], "capability_proof", "capability_id", "algorithm", "value")
	if err != nil {
		return nil, err
	}
	algorithm, _ := stringField(proof, "algorithm", "capability_proof")
	proofValue, _ := stringField(proof, "value", "capability_proof")
	if algorithm != "hmac-sha256" || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(proofValue) {
		return nil, fmt.Errorf("invalid capability proof")
	}
	return value.(canonicaljson.Object), nil
}

var startDocumentFields = []string{"kind", "schema_version", "format", "format_version", "canonicalization", "receipt_id", "binding", "principal", "observed_workspace", "observed_project", "observed_target", "observed_inputs", "contract_digests", "sandbox", "acceptance", "issues"}
var terminalDocumentFields = []string{"kind", "schema_version", "format", "format_version", "canonicalization", "result_id", "binding", "start_binding", "principal", "reported_outcome", "rounds_used", "stop_reasons", "summary", "meaning", "final_target", "observed_effects", "verifier_results", "review", "artifacts", "evidence_gaps", "forbidden_effects_observed"}

func validateStoredAcceptedV1(value canonicaljson.Object, phase, capabilityID, secret string) error {
	kind := map[string]string{"start": "ply.workflow.start-receipt", "terminal": "ply.workflow.terminal-result"}[phase]
	fields := startDocumentFields
	if phase == "terminal" {
		fields = terminalDocumentFields
	}
	if kind == "" {
		return fmt.Errorf("unknown accepted phase %q", phase)
	}
	if _, err := validateFinalDocument(value, kind, fields); err != nil {
		return err
	}
	proofValue, _ := objectMember(value, "capability_proof")
	proof, _ := exactObject(proofValue, "capability_proof", "capability_id", "algorithm", "value")
	if objectMapString(proof, "capability_id") != capabilityID {
		return fmt.Errorf("capability ID does not match handoff")
	}
	payload := removeObjectMember(value, "capability_proof")
	payloadBytes, err := canonicaljson.Marshal(payload)
	if err != nil {
		return err
	}
	expected, err := capabilityProof(secret, payloadBytes)
	if err != nil {
		return err
	}
	if !hmac.Equal([]byte(objectMapString(proof, "value")), []byte(expected)) {
		return fmt.Errorf("capability proof does not match stored payload")
	}
	draftKind := map[string]string{"start": "ply.workflow.start-receipt-draft", "terminal": "ply.workflow.terminal-result-draft"}[phase]
	draft := replaceObjectMember(payload, "kind", draftKind)
	draftBytes, err := canonicaljson.Marshal(draft)
	if err != nil {
		return err
	}
	if phase == "start" {
		_, err = decodeStartDraft(draftBytes)
	} else {
		_, err = decodeTerminalDraft(draftBytes)
	}
	return err
}

func addCapabilityProof(draft canonicaljson.Object, finalKind, capabilityID, secret string) (canonicaljson.Object, []byte, error) {
	payload := replaceObjectMember(draft, "kind", finalKind)
	payloadBytes, err := canonicaljson.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	proof, err := capabilityProof(secret, payloadBytes)
	if err != nil {
		return nil, nil, err
	}
	final := append(canonicaljson.Object{}, payload...)
	final = append(final, canonicaljson.Member{Name: "capability_proof", Value: canonicaljson.Object{{Name: "capability_id", Value: capabilityID}, {Name: "algorithm", Value: "hmac-sha256"}, {Name: "value", Value: proof}}})
	bytes, err := canonicaljson.Marshal(final)
	return final, bytes, err
}

func schemaError(context string, err error) error {
	return classified(ErrorSchemaInvalid, fmt.Sprintf("%s: %v", context, err), err)
}
func digestBytes(input []byte) string {
	sum := sha256.Sum256(input)
	return fmt.Sprintf("sha256:%x", sum)
}
func firstError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
func setOf(values ...string) map[string]bool {
	result := make(map[string]bool, len(values))
	for _, value := range values {
		result[value] = true
	}
	return result
}
func objectMember(object canonicaljson.Object, name string) (canonicaljson.Value, bool) {
	for _, member := range object {
		if member.Name == name {
			return member.Value, true
		}
	}
	return nil, false
}
func replaceObjectMember(object canonicaljson.Object, name string, value canonicaljson.Value) canonicaljson.Object {
	result := make(canonicaljson.Object, 0, len(object))
	for _, member := range object {
		if member.Name == name {
			member.Value = value
		}
		result = append(result, member)
	}
	return result
}
func removeObjectMember(object canonicaljson.Object, name string) canonicaljson.Object {
	result := make(canonicaljson.Object, 0, len(object)-1)
	for _, member := range object {
		if member.Name != name {
			result = append(result, member)
		}
	}
	return result
}
func absoluteCleanString(fields map[string]canonicaljson.Value, name, context string) (string, error) {
	value, err := stringField(fields, name, context)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(value) || filepath.Clean(value) != value {
		return "", fmt.Errorf("%s.%s must be an absolute clean path", context, name)
	}
	return value, nil
}
func stringArray(value canonicaljson.Value, context string) ([]string, error) {
	values, ok := value.([]canonicaljson.Value)
	if !ok {
		return nil, fmt.Errorf("%s must be an array", context)
	}
	result := make([]string, 0, len(values))
	for _, item := range values {
		text, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("%s must contain strings", context)
		}
		result = append(result, text)
	}
	return result, nil
}
func sortedStringArray(value canonicaljson.Value, context string, allowEmpty bool) ([]string, error) {
	values, err := stringArray(value, context)
	if err != nil {
		return nil, err
	}
	if !allowEmpty && len(values) == 0 {
		return nil, fmt.Errorf("%s must not be empty", context)
	}
	if !sort.StringsAreSorted(values) {
		return nil, fmt.Errorf("%s must be sorted", context)
	}
	for index := 1; index < len(values); index++ {
		if values[index] == values[index-1] {
			return nil, fmt.Errorf("%s must be unique", context)
		}
	}
	return values, nil
}
func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return false
		}
	}
	return true
}
func validateUTC(value string) bool {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	return err == nil && strings.HasSuffix(value, "Z") && parsed.Location() == time.UTC
}
