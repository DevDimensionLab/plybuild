package workflowhandoff

import (
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

func Create(dependencies Dependencies, input CreateInput) (CreateResult, error) {
	return createWithActivity(dependencies, input, "", nil)
}

func createWithActivity(dependencies Dependencies, input CreateInput, existingActivityID string, prepared *CreateStoreInput) (CreateResult, error) {
	if err := requireDependencies(dependencies); err != nil {
		return CreateResult{}, err
	}
	draftPath, contents, err := readDraft(dependencies.Files, input.DraftPath, maxHandoffBytes, "handoff draft")
	_ = draftPath
	if err != nil {
		return CreateResult{}, err
	}
	draft, err := decodeHandoffDraft(contents)
	if err != nil {
		return CreateResult{}, err
	}
	observed, err := observeCreateBindings(dependencies, draft)
	if err != nil {
		return CreateResult{}, err
	}
	// The second observation is the last read-only binding check before Store.Create.
	reobserved, err := observeCreateBindings(dependencies, draft)
	if err != nil {
		return CreateResult{}, err
	}
	if !sameCreateObservation(observed, reobserved) {
		return CreateResult{}, classified(ErrorTargetConflict, "workspace, project, target, or input changed during create", nil)
	}

	activityID := existingActivityID
	if activityID == "" {
		activityID, err = randomID(dependencies.Random, "act_")
		if err != nil {
			return CreateResult{}, err
		}
	}
	runID, err := randomID(dependencies.Random, "run_")
	if err != nil {
		return CreateResult{}, err
	}
	handoffID, err := randomID(dependencies.Random, "hnd_")
	if err != nil {
		return CreateResult{}, err
	}
	receiptID, err := randomID(dependencies.Random, "rcp_")
	if err != nil {
		return CreateResult{}, err
	}
	resultID, err := randomID(dependencies.Random, "res_")
	if err != nil {
		return CreateResult{}, err
	}
	capabilityID, err := randomID(dependencies.Random, "cap_")
	if err != nil {
		return CreateResult{}, err
	}
	secretBytes := make([]byte, 32)
	if _, err := io.ReadFull(dependencies.Random, secretBytes); err != nil {
		return CreateResult{}, classified(ErrorIO, "generate reply capability: "+err.Error(), err)
	}
	secret := hex.EncodeToString(secretBytes)
	createdAt := dependencies.Clock.Now().UTC().Format("2006-01-02T15:04:05.999999999Z07:00")
	identityValue := canonicaljson.Object{{Name: "activity_id", Value: activityID}, {Name: "run_id", Value: runID}, {Name: "handoff_id", Value: handoffID}, {Name: "start_receipt_id", Value: receiptID}, {Name: "terminal_result_id", Value: resultID}, {Name: "publication_key", Value: draft.PublicationKey}, {Name: "activity_key", Value: draft.ActivityKey}, {Name: "created_at_utc", Value: createdAt}}
	replyRoot := filepath.Join(storeRoot(observed.workspace.Observation.Root), "runs", runID, "reply")
	inputsValue, _ := objectMember(draft.Value, "inputs")
	handoffValue := envelope("ply.workflow.handoff",
		canonicaljson.Member{Name: "identity", Value: identityValue},
		canonicaljson.Member{Name: "source_draft_sha256", Value: draft.Digest},
		canonicaljson.Member{Name: "goal", Value: draft.Goal}, canonicaljson.Member{Name: "recipient", Value: draft.Recipient},
		canonicaljson.Member{Name: "workspace_binding", Value: canonicaljson.Object{{Name: "root", Value: observed.workspace.Observation.Root}, {Name: "marker_format_version", Value: int64(observed.workspace.Observation.MarkerFormatVersion)}, {Name: "marker_sha256", Value: observed.workspace.Observation.MarkerSHA256}}},
		canonicaljson.Member{Name: "project_binding", Value: canonicaljson.Object{{Name: "project_id", Value: string(draft.Binding.ProjectID)}, {Name: "repo_id", Value: string(draft.Binding.RepoID)}, {Name: "registered_locator", Value: observed.repository.Locator}, {Name: "registered_git_common_dir", Value: observed.repository.GitCommonDir}}},
		canonicaljson.Member{Name: "target_binding", Value: targetValue(observed.target)},
		canonicaljson.Member{Name: "inputs", Value: inputsValue}, canonicaljson.Member{Name: "authority", Value: draft.Authority}, canonicaljson.Member{Name: "budget", Value: draft.Budget},
		canonicaljson.Member{Name: "procedure", Value: draft.Procedure}, canonicaljson.Member{Name: "verifiers", Value: draft.Verifiers}, canonicaljson.Member{Name: "stop_conditions", Value: draft.StopConditions}, canonicaljson.Member{Name: "reporting", Value: draft.Reporting},
		canonicaljson.Member{Name: "reply_capability", Value: canonicaljson.Object{
			{Name: "capability_id", Value: capabilityID},
			{Name: "secret", Value: secret},
			{Name: "reply_root", Value: replyRoot},
			{Name: "allowed_operations", Value: []canonicaljson.Value{"submit-result", "submit-start"}},
		}},
	)
	handoffBytes, err := canonicaljson.Marshal(handoffValue)
	if err != nil {
		return CreateResult{}, schemaError("handoff", err)
	}
	if len(handoffBytes) > maxHandoffBytes {
		return CreateResult{}, classified(ErrorPayloadTooLarge, "canonical handoff exceeds 256 KiB", nil)
	}
	document := handoffDocument{Value: handoffValue, Bytes: handoffBytes, SHA256: digestBytes(handoffBytes), Identity: identity{ActivityID: activityID, RunID: runID, HandoffID: handoffID, StartReceiptID: receiptID, TerminalResultID: resultID, PublicationKey: draft.PublicationKey, ActivityKey: draft.ActivityKey, CreatedAtUTC: createdAt}, SourceDraftSHA256: draft.Digest, GoalTitle: draft.GoalTitle, Target: observed.target, Workspace: observed.workspace.Observation, ProjectID: draft.Binding.ProjectID, RepoID: draft.Binding.RepoID, RegisteredLocator: observed.repository.Locator, RegisteredGitCommonDir: observed.repository.GitCommonDir, ReplyCapabilityID: capabilityID, ReplySecret: secret, ReplyRoot: replyRoot, MaxRounds: draft.MaxRounds}
	storeInput := CreateStoreInput{WorkspaceRoot: observed.workspace.Observation.Root, PublicationKey: draft.PublicationKey, SourceDraftSHA256: draft.Digest, Handoff: document, DeferActivityPublication: prepared != nil, Revalidate: func() error {
		latest, err := observeCreateBindings(dependencies, draft)
		if err != nil {
			return err
		}
		if !sameCreateObservation(reobserved, latest) {
			return classified(ErrorTargetConflict, "workspace, project, target, or input changed before publication", nil)
		}
		return nil
	}}
	if prepared != nil {
		handoffName := "handoff." + strings.TrimPrefix(document.SHA256, "sha256:") + ".json"
		storeInput.Handoff.Locator = filepath.Join(storeRoot(observed.workspace.Observation.Root), "runs", runID, handoffName)
		*prepared = storeInput
		return CreateResult{HandoffID: HandoffID(handoffID), Purpose: document.GoalTitle, Worktree: document.Target.Worktree, Locator: storeInput.Handoff.Locator, Created: true}, nil
	}
	stored, err := dependencies.Store.Create(storeInput)
	if err != nil {
		return CreateResult{}, err
	}
	h := stored.Snapshot.Handoff
	return CreateResult{HandoffID: HandoffID(h.Identity.HandoffID), Purpose: h.GoalTitle, Worktree: h.Target.Worktree, Locator: h.Locator, Created: stored.Created}, nil
}

type createObservation struct {
	workspace  WorkspaceSnapshot
	repository workspace.RepoRecord
	target     TargetObservation
	inputs     []observedInput
}
type observedInput struct {
	ID, Locator, SHA256 string
	Size                int64
	Git                 *InputGitObservation
}

func observeCreateBindings(dependencies Dependencies, draft handoffDraft) (createObservation, error) {
	snapshot, err := dependencies.Workspace.ObserveContaining()
	if err != nil {
		return createObservation{}, err
	}
	projectFound := false
	member := false
	for _, project := range snapshot.Projects {
		if project.ID == draft.Binding.ProjectID {
			projectFound = true
			for _, id := range project.RepoIDs {
				if id == draft.Binding.RepoID {
					member = true
				}
			}
			break
		}
	}
	if !projectFound {
		return createObservation{}, classified(ErrorProjectNotFound, fmt.Sprintf("project %s is not registered", draft.Binding.ProjectID), nil)
	}
	if !member {
		return createObservation{}, classified(ErrorRepoNotFound, fmt.Sprintf("repository %s is not a member of project %s", draft.Binding.RepoID, draft.Binding.ProjectID), nil)
	}
	var repository workspace.RepoRecord
	found := false
	for _, candidate := range snapshot.Repositories {
		if candidate.ID == draft.Binding.RepoID {
			repository = candidate
			found = true
			break
		}
	}
	if !found {
		return createObservation{}, classified(ErrorRepoNotFound, fmt.Sprintf("repository %s is not registered", draft.Binding.RepoID), nil)
	}
	target, err := dependencies.Git.ObserveTarget(TargetRequest{Worktree: draft.Binding.TargetWorktree, Ref: draft.Binding.TargetRef})
	if err != nil {
		return createObservation{}, err
	}
	if target.GitCommonDir != repository.GitCommonDir || target.OID != draft.Binding.ExpectedOID || !equalStatus(target.Status, draft.Binding.Status) {
		return createObservation{}, classified(ErrorTargetConflict, "target common directory, OID, or status does not match the draft", nil)
	}
	inputs := make([]observedInput, 0, len(draft.Inputs))
	for _, input := range draft.Inputs {
		locator, err := physicalRegularFile(dependencies.Files, input.Locator)
		if err != nil {
			return createObservation{}, classified(ErrorSchemaInvalid, fmt.Sprintf("input %s: %v", input.ID, err), err)
		}
		digest, size, err := hashFile(dependencies.Files, locator)
		if err != nil {
			return createObservation{}, classified(ErrorIO, fmt.Sprintf("read input %s: %v", input.ID, err), err)
		}
		if digest != input.SHA256 || size != input.SizeBytes {
			return createObservation{}, classified(ErrorTargetConflict, fmt.Sprintf("input %s bytes do not match draft", input.ID), nil)
		}
		observed := observedInput{ID: input.ID, Locator: locator, SHA256: digest, Size: size}
		if input.GitBinding != nil {
			gitObservation, err := dependencies.Git.ObserveInputGit(InputGitRequest{Locator: locator, Binding: *input.GitBinding})
			if err != nil {
				return createObservation{}, err
			}
			if !gitObservation.MatchesExpected {
				return createObservation{}, classified(ErrorTargetConflict, fmt.Sprintf("input %s Git binding does not match", input.ID), nil)
			}
			observed.Git = &gitObservation
		}
		inputs = append(inputs, observed)
	}
	return createObservation{workspace: snapshot, repository: repository, target: target, inputs: inputs}, nil
}

func SubmitStart(dependencies Dependencies, input SubmitInput) (SubmitResult, error) {
	if err := requireDependencies(dependencies); err != nil {
		return SubmitResult{}, err
	}
	snapshot, err := dependencies.Store.ReadByLocator(input.HandoffLocator)
	if err != nil {
		return SubmitResult{}, err
	}
	_, contents, err := readDraft(dependencies.Files, input.DraftPath, maxStartBytes, "start receipt draft")
	if err != nil {
		if IsClass(err, ErrorPayloadTooLarge) {
			if system, ok := dependencies.Store.(*systemStore); ok {
				_ = system.publishReplyRejection(snapshot, "start", "rejected", nil, ErrorPayloadTooLarge, err.Error())
			}
		}
		return SubmitResult{}, err
	}
	draft, err := decodeStartDraft(contents)
	if err != nil {
		if system, ok := dependencies.Store.(*systemStore); ok {
			_ = system.publishReplyRejection(snapshot, "start", "rejected", contents, ErrorSchemaInvalid, err.Error())
		}
		return SubmitResult{}, err
	}
	if err := validateStartBinding(dependencies, snapshot, draft); err != nil {
		if system, ok := dependencies.Store.(*systemStore); ok {
			_ = system.publishReplyRejection(snapshot, "start", "conflict", draft.Canonical, classOf(err), err.Error())
		}
		return SubmitResult{}, err
	}
	final, bytes, err := addCapabilityProof(draft.Value, "ply.workflow.start-receipt", snapshot.Handoff.ReplyCapabilityID, snapshot.Handoff.ReplySecret)
	if err != nil {
		return SubmitResult{}, schemaError("start receipt", err)
	}
	document := acceptedDocument{Value: final, Bytes: bytes, SHA256: digestBytes(bytes), DocumentID: draft.ReceiptID, Outcome: draft.Acceptance}
	stored, err := dependencies.Store.SubmitStart(SubmitStoreInput{Snapshot: snapshot, Phase: "start", Document: document, Revalidate: func() error { return validateStartBinding(dependencies, snapshot, draft) }})
	if err != nil {
		return SubmitResult{}, err
	}
	return SubmitResult{Phase: "start", DocumentID: stored.Document.DocumentID, Locator: stored.Document.Locator, SHA256: stored.Document.SHA256, Created: stored.Created}, nil
}

func SubmitResultDocument(dependencies Dependencies, input SubmitInput) (SubmitResult, error) {
	if err := requireDependencies(dependencies); err != nil {
		return SubmitResult{}, err
	}
	snapshot, err := dependencies.Store.ReadByLocator(input.HandoffLocator)
	if err != nil {
		return SubmitResult{}, err
	}
	_, contents, err := readDraft(dependencies.Files, input.DraftPath, maxResultBytes, "terminal result draft")
	if err != nil {
		if IsClass(err, ErrorPayloadTooLarge) {
			if system, ok := dependencies.Store.(*systemStore); ok {
				_ = system.publishReplyRejection(snapshot, "terminal", "rejected", nil, ErrorPayloadTooLarge, err.Error())
			}
		}
		return SubmitResult{}, err
	}
	if snapshot.Start == nil {
		if system, ok := dependencies.Store.(*systemStore); ok {
			_ = system.publishReplyRejection(snapshot, "terminal", "rejected", contents, ErrorStartRequired, "a terminal result requires an accepted start receipt")
		}
		return SubmitResult{}, classified(ErrorStartRequired, "a terminal result requires an accepted start receipt", nil)
	}
	draft, err := decodeTerminalDraft(contents)
	if err != nil {
		if system, ok := dependencies.Store.(*systemStore); ok {
			_ = system.publishReplyRejection(snapshot, "terminal", "rejected", contents, ErrorSchemaInvalid, err.Error())
		}
		return SubmitResult{}, err
	}
	finalDraft, artifacts, err := validateTerminalBinding(dependencies, snapshot, draft)
	if err != nil {
		if system, ok := dependencies.Store.(*systemStore); ok {
			_ = system.publishReplyRejection(snapshot, "terminal", "conflict", draft.Canonical, classOf(err), err.Error())
		}
		return SubmitResult{}, err
	}
	final, bytes, err := addCapabilityProof(finalDraft, "ply.workflow.terminal-result", snapshot.Handoff.ReplyCapabilityID, snapshot.Handoff.ReplySecret)
	if err != nil {
		return SubmitResult{}, schemaError("terminal result", err)
	}
	document := acceptedDocument{Value: final, Bytes: bytes, SHA256: digestBytes(bytes), DocumentID: draft.ResultID, Outcome: draft.Outcome, Summary: draft.Summary, Meaning: draft.Meaning}
	stored, err := dependencies.Store.SubmitResult(SubmitStoreInput{Snapshot: snapshot, Phase: "terminal", Document: document, Artifacts: artifacts, Revalidate: func() error {
		rechecked, recheckedArtifacts, err := validateTerminalBinding(dependencies, snapshot, draft)
		if err != nil {
			return err
		}
		if !canonicalEqual(finalDraft, rechecked) || !equalManagedArtifacts(artifacts, recheckedArtifacts) {
			return classified(ErrorTargetConflict, "target or managed artifacts changed before terminal publication", nil)
		}
		return nil
	}})
	if err != nil {
		return SubmitResult{}, err
	}
	return SubmitResult{Phase: "terminal", DocumentID: stored.Document.DocumentID, Locator: stored.Document.Locator, SHA256: stored.Document.SHA256, Created: stored.Created}, nil
}

func Cancel(dependencies Dependencies, input ControlInput) (ControlResult, error) {
	if err := validateControl(dependencies, input, false); err != nil {
		return ControlResult{}, err
	}
	snapshot, err := dependencies.Workspace.ObserveContaining()
	if err != nil {
		return ControlResult{}, err
	}
	stored, err := dependencies.Store.Cancel(ControlStoreInput{WorkspaceRoot: snapshot.Observation.Root, HandoffID: input.HandoffID, Reason: input.Reason, State: "cancelled"})
	if err != nil {
		return ControlResult{}, err
	}
	return ControlResult{HandoffID: input.HandoffID, Reason: input.Reason, Locator: stored.Snapshot.Handoff.Locator}, nil
}
func Abandon(dependencies Dependencies, input ControlInput) (ControlResult, error) {
	if err := validateControl(dependencies, input, true); err != nil {
		return ControlResult{}, err
	}
	snapshot, err := dependencies.Workspace.ObserveContaining()
	if err != nil {
		return ControlResult{}, err
	}
	stored, err := dependencies.Store.Abandon(ControlStoreInput{WorkspaceRoot: snapshot.Observation.Root, HandoffID: input.HandoffID, Reason: input.Reason, State: "abandoned_unknown"})
	if err != nil {
		return ControlResult{}, err
	}
	return ControlResult{HandoffID: input.HandoffID, Reason: input.Reason, Locator: stored.Snapshot.Handoff.Locator}, nil
}
func Supersede(dependencies Dependencies, input SupersedeInput) (CreateResult, error) {
	if err := requireDependencies(dependencies); err != nil {
		return CreateResult{}, err
	}
	if err := validateReason(input.Reason); err != nil {
		return CreateResult{}, err
	}
	workspaceSnapshot, err := dependencies.Workspace.ObserveContaining()
	if err != nil {
		return CreateResult{}, err
	}
	old, err := dependencies.Store.ReadByID(workspaceSnapshot.Observation.Root, input.HandoffID)
	if err != nil {
		return CreateResult{}, err
	}
	_, contents, err := readDraft(dependencies.Files, input.DraftPath, maxHandoffBytes, "handoff draft")
	if err != nil {
		return CreateResult{}, err
	}
	draft, err := decodeHandoffDraft(contents)
	if err != nil {
		return CreateResult{}, err
	}
	if draft.ActivityKey != old.Handoff.Identity.ActivityKey {
		return CreateResult{}, classified(ErrorConflict, "replacement activity_key does not match the existing activity", nil)
	}
	var replacement CreateStoreInput
	result, err := createWithActivity(dependencies, CreateInput{DraftPath: input.DraftPath}, old.Handoff.Identity.ActivityID, &replacement)
	if err != nil {
		return CreateResult{}, err
	}
	stored, err := dependencies.Store.Supersede(SupersedeStoreInput{WorkspaceRoot: workspaceSnapshot.Observation.Root, OldHandoffID: input.HandoffID, Reason: input.Reason, Create: replacement})
	if err != nil {
		return CreateResult{}, err
	}
	handoff := stored.Snapshot.Handoff
	return CreateResult{HandoffID: HandoffID(handoff.Identity.HandoffID), Purpose: handoff.GoalTitle, Worktree: handoff.Target.Worktree, Locator: handoff.Locator, Created: stored.Created && result.Created}, nil
}

func validateStartBinding(dependencies Dependencies, snapshot Snapshot, draft startDraft) error {
	if draft.ReceiptID != snapshot.Handoff.Identity.StartReceiptID {
		return classified(ErrorCapabilityInvalid, "receipt_id does not match the handoff", nil)
	}
	binding, _ := objectMember(draft.Value, "binding")
	bf, _ := exactObject(binding, "binding", "activity_id", "run_id", "handoff_id", "handoff_sha256")
	if objectMapString(bf, "activity_id") != snapshot.Handoff.Identity.ActivityID || objectMapString(bf, "run_id") != snapshot.Handoff.Identity.RunID || objectMapString(bf, "handoff_id") != snapshot.Handoff.Identity.HandoffID || objectMapString(bf, "handoff_sha256") != snapshot.Handoff.SHA256 {
		return classified(ErrorCapabilityInvalid, "start binding does not match the handoff", nil)
	}
	principalFields, _ := exactObject(draft.Principal, "principal", "expected_principal_id", "human_start_principal", "start_surface", "session_id", "runtime_id", "model_id", "started_at_utc")
	expected := handoffNestedString(snapshot.Handoff.Value, "recipient", "principal_id")
	if objectMapString(principalFields, "expected_principal_id") != expected {
		return classified(ErrorCapabilityInvalid, "expected principal does not match the handoff", nil)
	}
	if !validateUTC(objectMapString(principalFields, "started_at_utc")) {
		return schemaError("start receipt draft", fmt.Errorf("principal.started_at_utc is invalid"))
	}
	contract, _ := objectMember(draft.Value, "contract_digests")
	cf, _ := exactObject(contract, "contract_digests", "authority_sha256", "budget_sha256", "verifiers_sha256", "stop_conditions_sha256")
	for field, section := range map[string]string{"authority_sha256": "authority", "budget_sha256": "budget", "verifiers_sha256": "verifiers", "stop_conditions_sha256": "stop_conditions"} {
		value, _ := objectMember(snapshot.Handoff.Value, section)
		bytes, _ := canonicaljson.Marshal(value)
		if objectMapString(cf, field) != digestBytes(bytes) {
			return classified(ErrorCapabilityInvalid, "contract digest mismatch for "+section, nil)
		}
	}
	workspaceSnapshot, err := dependencies.Workspace.ObserveRoot(snapshot.Handoff.Workspace.Root)
	if err != nil {
		return err
	}
	repository, projectMember := findBoundRepository(workspaceSnapshot, snapshot.Handoff.ProjectID, snapshot.Handoff.RepoID)
	projectMatch := projectMember && repository.Locator == snapshot.Handoff.RegisteredLocator && repository.GitCommonDir == snapshot.Handoff.RegisteredGitCommonDir
	target, err := dependencies.Git.ObserveTarget(TargetRequest{Worktree: snapshot.Handoff.Target.Worktree, Ref: snapshot.Handoff.Target.Ref})
	if err != nil {
		return err
	}
	workspaceExpected := workspaceSnapshot.Observation == snapshot.Handoff.Workspace
	targetExpected := sameTarget(target, snapshot.Handoff.Target)
	observedWorkspaceValue, _ := objectMember(draft.Value, "observed_workspace")
	observedWorkspace := observedWorkspaceValue.(canonicaljson.Object)
	workspaceReported := objectString(observedWorkspace, "root") == workspaceSnapshot.Observation.Root && objectInt(observedWorkspace, "marker_format_version") == int64(workspaceSnapshot.Observation.MarkerFormatVersion) && objectString(observedWorkspace, "marker_sha256") == workspaceSnapshot.Observation.MarkerSHA256
	if objectBool(observedWorkspace, "matches_expected") != workspaceExpected || !workspaceReported {
		return classified(ErrorConflict, "observed workspace does not match the current binding", nil)
	}
	observedProjectValue, _ := objectMember(draft.Value, "observed_project")
	observedProject := observedProjectValue.(canonicaljson.Object)
	projectReported := projectMember && objectString(observedProject, "project_id") == string(snapshot.Handoff.ProjectID) && objectString(observedProject, "repo_id") == string(snapshot.Handoff.RepoID) && objectString(observedProject, "registered_locator") == repository.Locator && objectString(observedProject, "registered_git_common_dir") == repository.GitCommonDir
	if objectBool(observedProject, "matches_expected") != projectMatch || !projectReported {
		return classified(ErrorConflict, "observed project does not match the current binding", nil)
	}
	observedTargetValue, _ := objectMember(draft.Value, "observed_target")
	observedTarget := observedTargetValue.(canonicaljson.Object)
	reportedStatus, _ := func() (StatusPolicy, error) {
		value, _ := objectMember(observedTarget, "status_policy")
		return validateStatusPolicy(value)
	}()
	targetReported := objectString(observedTarget, "worktree") == target.Worktree && objectString(observedTarget, "git_common_dir") == target.GitCommonDir && objectString(observedTarget, "ref") == target.Ref && objectString(observedTarget, "oid") == target.OID && objectString(observedTarget, "tree") == target.Tree && equalStatus(reportedStatus, target.Status)
	if objectBool(observedTarget, "matches_expected") != targetExpected || !targetReported {
		return classified(ErrorConflict, "observed target does not match the current binding", nil)
	}
	inputsExpected, err := validateObservedInputsAtStart(dependencies, snapshot, draft.Value)
	if err != nil {
		return err
	}
	sandboxValue, _ := objectMember(draft.Value, "sandbox")
	sandbox := sandboxValue.(canonicaljson.Object)
	sandboxErr := validateSandboxContract(dependencies.Files, snapshot, sandbox)
	sandboxActual := sandboxErr == nil
	if objectBool(sandbox, "matches_contract") != sandboxActual {
		detail := "sandbox matches_contract does not match reported roots"
		if sandboxErr != nil {
			detail += ": " + sandboxErr.Error()
		}
		return classified(ErrorConflict, detail, sandboxErr)
	}
	matches := workspaceExpected && projectMatch && targetExpected && inputsExpected && sandboxActual
	if draft.Acceptance == "started" && !matches {
		return classified(ErrorConflict, "start acceptance does not match the observed binding and sandbox", nil)
	}
	return nil
}

func validateTerminalBinding(dependencies Dependencies, snapshot Snapshot, draft terminalDraft) (canonicaljson.Object, []managedArtifact, error) {
	if draft.ResultID != snapshot.Handoff.Identity.TerminalResultID {
		return nil, nil, classified(ErrorCapabilityInvalid, "result_id does not match the handoff", nil)
	}
	bindingValue, _ := objectMember(draft.Value, "binding")
	binding, _ := exactObject(bindingValue, "binding", "activity_id", "run_id", "handoff_id", "handoff_sha256")
	if objectMapString(binding, "activity_id") != snapshot.Handoff.Identity.ActivityID || objectMapString(binding, "run_id") != snapshot.Handoff.Identity.RunID || objectMapString(binding, "handoff_id") != snapshot.Handoff.Identity.HandoffID || objectMapString(binding, "handoff_sha256") != snapshot.Handoff.SHA256 {
		return nil, nil, classified(ErrorCapabilityInvalid, "terminal binding does not match the handoff", nil)
	}
	startValue, _ := objectMember(draft.Value, "start_binding")
	start, _ := exactObject(startValue, "start_binding", "receipt_id", "start_receipt_sha256")
	if objectMapString(start, "receipt_id") != snapshot.Start.DocumentID || objectMapString(start, "start_receipt_sha256") != snapshot.Start.SHA256 {
		return nil, nil, classified(ErrorCapabilityInvalid, "terminal start binding does not match accepted start", nil)
	}
	startPrincipal, _ := objectMember(snapshot.Start.Value, "principal")
	startBytes, _ := canonicaljson.Marshal(startPrincipal)
	terminalBytes, _ := canonicaljson.Marshal(draft.Principal)
	if subtle.ConstantTimeCompare(startBytes, terminalBytes) != 1 {
		return nil, nil, classified(ErrorCapabilityInvalid, "terminal principal/session does not match accepted start", nil)
	}
	target, err := dependencies.Git.ObserveTarget(TargetRequest{Worktree: snapshot.Handoff.Target.Worktree, Ref: snapshot.Handoff.Target.Ref})
	if err != nil {
		return nil, nil, err
	}
	finalValue, _ := objectMember(draft.Value, "final_target")
	final, _ := exactObject(finalValue, "final_target", "worktree", "git_common_dir", "ref", "oid", "tree", "status_policy", "matches_expected")
	claimedMatches, _ := final["matches_expected"].(bool)
	reportedStatus, _ := validateStatusPolicy(final["status_policy"])
	actualMatches := objectMapString(final, "worktree") == target.Worktree && objectMapString(final, "git_common_dir") == target.GitCommonDir && objectMapString(final, "ref") == target.Ref && objectMapString(final, "oid") == target.OID && objectMapString(final, "tree") == target.Tree && equalStatus(reportedStatus, target.Status)
	if claimedMatches != actualMatches {
		return nil, nil, classified(ErrorConflict, "final target matches_expected does not match observation", nil)
	}
	if snapshot.Start.Outcome != "started" && draft.Outcome == "complete" {
		return nil, nil, classified(ErrorConflict, "a failed start cannot report a complete terminal outcome", nil)
	}
	if err := validateObservedEffectsAgainstHandoff(snapshot.Handoff.Value, draft.Value); err != nil {
		return nil, nil, classified(ErrorConflict, err.Error(), err)
	}
	if err := validateVerifierResultsAgainstHandoff(snapshot.Handoff.Value, draft.Value); err != nil {
		return nil, nil, classified(ErrorConflict, err.Error(), err)
	}
	managed := []managedArtifact{}
	rewritten := make([]canonicaljson.Value, 0, len(draft.Artifacts))
	var total int64
	lastID := ""
	for _, artifactValue := range draft.Artifacts {
		object, ok := artifactValue.(canonicaljson.Object)
		if !ok {
			return nil, nil, schemaError("terminal result draft", fmt.Errorf("artifact must be an object"))
		}
		kind := objectString(object, "kind")
		id := objectString(object, "artifact_id")
		if validateKey("artifact_id", id) != nil || (lastID != "" && id <= lastID) {
			return nil, nil, schemaError("terminal result draft", fmt.Errorf("artifact IDs must be sorted and unique"))
		}
		lastID = id
		if kind == "managed" {
			fields, err := exactObject(object, "managed artifact", "artifact_id", "kind", "description", "media_type", "classification", "size_bytes", "sha256", "locator")
			if err != nil {
				return nil, nil, schemaError("terminal result draft", err)
			}
			locator, err := absoluteCleanString(fields, "locator", "managed artifact")
			if err != nil {
				return nil, nil, schemaError("terminal result draft", err)
			}
			staging := filepath.Join(snapshot.Handoff.ReplyRoot, "staging") + string(filepath.Separator)
			if !strings.HasPrefix(locator, staging) {
				return nil, nil, classified(ErrorCapabilityInvalid, "managed artifact is outside reply staging", nil)
			}
			physical, err := physicalRegularFile(dependencies.Files, locator)
			if err != nil || physical != locator {
				return nil, nil, classified(ErrorCapabilityInvalid, "managed artifact is not a physical regular file", err)
			}
			digest, size, err := hashFile(dependencies.Files, locator)
			if err != nil {
				return nil, nil, classified(ErrorIO, "read managed artifact: "+err.Error(), err)
			}
			claimedSize, _ := intField(fields, "size_bytes", "managed artifact")
			claimedDigest, _ := stringField(fields, "sha256", "managed artifact")
			classification, _ := stringField(fields, "classification", "managed artifact")
			if size > 64<<20 || claimedSize != size || claimedDigest != digest || classification != "workspace_internal" {
				return nil, nil, classified(ErrorConflict, "managed artifact metadata does not match bytes", nil)
			}
			total += size
			if total > 256<<20 {
				return nil, nil, classified(ErrorPayloadTooLarge, "managed artifacts exceed 256 MiB", nil)
			}
			bytes, _ := dependencies.Files.ReadFile(locator)
			destination := filepath.Join(snapshot.Handoff.ReplyRoot, "artifacts", "sha256", strings.TrimPrefix(digest, "sha256:"))
			managed = append(managed, managedArtifact{ID: id, Locator: destination, SHA256: digest, Bytes: bytes})
			rewritten = append(rewritten, replaceObjectMember(object, "locator", destination))
		} else if kind == "withheld" {
			if _, err := exactObject(object, "withheld artifact", "artifact_id", "kind", "description", "reason", "media_type", "size_bytes", "sha256"); err != nil {
				return nil, nil, schemaError("terminal result draft", err)
			}
			rewritten = append(rewritten, object)
		} else {
			return nil, nil, schemaError("terminal result draft", fmt.Errorf("unknown artifact kind"))
		}
	}
	return replaceObjectMember(draft.Value, "artifacts", rewritten), managed, nil
}

func validateControl(dependencies Dependencies, input ControlInput, ack bool) error {
	if err := requireDependencies(dependencies); err != nil {
		return err
	}
	if _, err := ParseHandoffID(string(input.HandoffID)); err != nil {
		return err
	}
	if err := validateReason(input.Reason); err != nil {
		return err
	}
	if ack && !input.AcknowledgeEffectsUnknown {
		return InvalidArguments("--acknowledge-effects-unknown must be true")
	}
	return nil
}
func validateReason(reason string) error {
	if err := validatePlainText("reason", reason, 1, 600); err != nil {
		return InvalidArguments(err.Error())
	}
	return nil
}
func requireDependencies(dependencies Dependencies) error {
	if dependencies.Files == nil || dependencies.Workspace == nil || dependencies.Git == nil || dependencies.Clock == nil || dependencies.Random == nil || dependencies.Store == nil {
		return classified(ErrorIO, "dependencies: filesystem, workspace, Git, clock, random, and store are required", nil)
	}
	return nil
}
func readDraft(files FileSystem, input string, limit int, label string) (string, []byte, error) {
	if files == nil {
		return "", nil, classified(ErrorIO, "filesystem dependency is required", nil)
	}
	cwd, err := files.Getwd()
	if err != nil {
		return "", nil, classified(ErrorIO, "get working directory: "+err.Error(), err)
	}
	path := input
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	path = filepath.Clean(path)
	info, err := files.Lstat(path)
	if err != nil {
		return "", nil, classified(ErrorIO, fmt.Sprintf("inspect %s: %v", path, err), err)
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return "", nil, classified(ErrorSchemaInvalid, fmt.Sprintf("%s %s is not a regular non-symlink file", label, path), nil)
	}
	physical, err := files.EvalSymlinks(path)
	if err != nil || physical != path {
		return "", nil, classified(ErrorSchemaInvalid, fmt.Sprintf("%s path must be physical", label), err)
	}
	if info.Size() > int64(limit) {
		return "", nil, classified(ErrorPayloadTooLarge, fmt.Sprintf("%s exceeds %d bytes", label, limit), nil)
	}
	contents, err := files.ReadFile(path)
	if err != nil {
		return "", nil, classified(ErrorIO, fmt.Sprintf("read %s: %v", path, err), err)
	}
	if len(contents) > limit {
		return "", nil, classified(ErrorPayloadTooLarge, fmt.Sprintf("%s exceeds %d bytes", label, limit), nil)
	}
	return path, contents, nil
}
func physicalRegularFile(files FileSystem, path string) (string, error) {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return "", fmt.Errorf("path must be absolute and clean")
	}
	info, err := files.Lstat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Mode()&fs.ModeSymlink != 0 {
		return "", fmt.Errorf("path is not a regular non-symlink file")
	}
	physical, err := files.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	if physical != path {
		return "", fmt.Errorf("path is not physical")
	}
	return physical, nil
}
func randomID(random io.Reader, prefix string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := io.ReadFull(random, bytes); err != nil {
		return "", classified(ErrorIO, "generate "+prefix+" identity: "+err.Error(), err)
	}
	return prefix + hex.EncodeToString(bytes), nil
}
func targetValue(target TargetObservation) canonicaljson.Object {
	return canonicaljson.Object{{Name: "worktree", Value: target.Worktree}, {Name: "git_common_dir", Value: target.GitCommonDir}, {Name: "ref", Value: target.Ref}, {Name: "oid", Value: target.OID}, {Name: "tree", Value: target.Tree}, {Name: "status_policy", Value: target.Status.Value}}
}
func equalStatus(left, right StatusPolicy) bool {
	a, _ := canonicaljson.Marshal(left.Value)
	b, _ := canonicaljson.Marshal(right.Value)
	return string(a) == string(b)
}
func canonicalEqual(left, right canonicaljson.Value) bool {
	a, errA := canonicaljson.Marshal(left)
	b, errB := canonicaljson.Marshal(right)
	return errA == nil && errB == nil && string(a) == string(b)
}
func sameCreateObservation(left, right createObservation) bool {
	if left.workspace.Observation != right.workspace.Observation || left.repository != right.repository || left.target.Worktree != right.target.Worktree || left.target.GitCommonDir != right.target.GitCommonDir || left.target.Ref != right.target.Ref || left.target.OID != right.target.OID || left.target.Tree != right.target.Tree || !equalStatus(left.target.Status, right.target.Status) || len(left.inputs) != len(right.inputs) {
		return false
	}
	for index := range left.inputs {
		if left.inputs[index].ID != right.inputs[index].ID || left.inputs[index].SHA256 != right.inputs[index].SHA256 || left.inputs[index].Size != right.inputs[index].Size {
			return false
		}
	}
	return true
}
func handoffNestedString(handoff canonicaljson.Object, objectName, fieldName string) string {
	value, _ := objectMember(handoff, objectName)
	object, _ := value.(canonicaljson.Object)
	return objectString(object, fieldName)
}
func objectMapString(fields map[string]canonicaljson.Value, name string) string {
	value, _ := fields[name].(string)
	return value
}
func objectInt(object canonicaljson.Object, name string) int64 {
	value, _ := objectMember(object, name)
	number, _ := value.(int64)
	return number
}
func objectBool(object canonicaljson.Object, name string) bool {
	value, _ := objectMember(object, name)
	flag, _ := value.(bool)
	return flag
}

func sameTarget(left, right TargetObservation) bool {
	return left.Worktree == right.Worktree && left.GitCommonDir == right.GitCommonDir && left.Ref == right.Ref && left.OID == right.OID && left.Tree == right.Tree && equalStatus(left.Status, right.Status)
}

func findBoundRepository(snapshot WorkspaceSnapshot, projectID workspace.ProjectID, repoID workspace.RepoID) (workspace.RepoRecord, bool) {
	member := false
	for _, project := range snapshot.Projects {
		if project.ID == projectID {
			for _, candidate := range project.RepoIDs {
				if candidate == repoID {
					member = true
				}
			}
		}
	}
	for _, repository := range snapshot.Repositories {
		if repository.ID == repoID {
			return repository, member
		}
	}
	return workspace.RepoRecord{}, false
}

func validateObservedInputsAtStart(dependencies Dependencies, snapshot Snapshot, start canonicaljson.Object) (bool, error) {
	handoffValue, _ := objectMember(snapshot.Handoff.Value, "inputs")
	expected, _ := handoffValue.([]canonicaljson.Value)
	reportedValue, _ := objectMember(start, "observed_inputs")
	reported, _ := reportedValue.([]canonicaljson.Value)
	if len(expected) != len(reported) {
		return false, classified(ErrorConflict, "observed input set does not match the handoff", nil)
	}
	all := true
	for index, expectedValue := range expected {
		spec := expectedValue.(canonicaljson.Object)
		actual := reported[index].(canonicaljson.Object)
		id := objectString(spec, "id")
		if objectString(actual, "id") != id {
			return false, classified(ErrorConflict, "observed input IDs do not match the handoff", nil)
		}
		locator := objectString(spec, "locator")
		physical, physicalErr := physicalRegularFile(dependencies.Files, locator)
		digest, size, hashErr := hashFile(dependencies.Files, locator)
		matches := physicalErr == nil && hashErr == nil && physical == locator && digest == objectString(spec, "sha256") && size == objectInt(spec, "size_bytes")
		if gitValue, found := objectMember(spec, "git_binding"); found && gitValue != nil {
			git := gitValue.(canonicaljson.Object)
			observation, err := dependencies.Git.ObserveInputGit(InputGitRequest{Locator: locator, Binding: GitBinding{GitCommonDir: objectString(git, "git_common_dir"), Ref: objectString(git, "ref"), OID: objectString(git, "oid"), Blob: objectString(git, "blob")}})
			if err != nil {
				return false, err
			}
			matches = matches && observation.MatchesExpected
		}
		reportedCurrent := objectString(actual, "locator") == locator && objectString(actual, "sha256") == digest && objectInt(actual, "size_bytes") == size
		if !reportedCurrent || objectBool(actual, "matches_expected") != matches {
			return false, classified(ErrorConflict, "observed input does not match current bytes", nil)
		}
		all = all && matches
	}
	return all, nil
}

func validateSandboxContract(files FileSystem, snapshot Snapshot, sandbox canonicaljson.Object) error {
	readValue, _ := objectMember(sandbox, "read_roots")
	readRoots, _ := stringArray(readValue, "sandbox.read_roots")
	writeValue, _ := objectMember(sandbox, "write_roots")
	writeRoots, _ := stringArray(writeValue, "sandbox.write_roots")
	tempRoot := objectString(sandbox, "temp_root")
	for _, root := range append(append([]string{}, readRoots...), append(writeRoots, tempRoot)...) {
		physical, err := files.EvalSymlinks(root)
		if err != nil || physical != root {
			return fmt.Errorf("sandbox root %s is not physical", root)
		}
	}
	info, err := files.Lstat(tempRoot)
	if err != nil || !info.IsDir() || info.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("sandbox temp_root must be a private directory")
	}
	requiredReads := []string{snapshot.Handoff.Locator, snapshot.Handoff.Target.Worktree}
	inputsValue, _ := objectMember(snapshot.Handoff.Value, "inputs")
	if inputs, ok := inputsValue.([]canonicaljson.Value); ok {
		for _, value := range inputs {
			requiredReads = append(requiredReads, objectString(value.(canonicaljson.Object), "locator"))
		}
	}
	for _, required := range requiredReads {
		covered := false
		for _, root := range readRoots {
			if root == required || strings.HasPrefix(required, root+string(filepath.Separator)) {
				covered = true
				break
			}
		}
		if !covered {
			return fmt.Errorf("sandbox read roots do not cover %s", required)
		}
	}
	writes := []string{snapshot.Handoff.ReplyRoot}
	authorityValue, _ := objectMember(snapshot.Handoff.Value, "authority")
	authority := authorityValue.(canonicaljson.Object)
	allowedValue, _ := objectMember(authority, "allowed_effects")
	for _, value := range allowedValue.([]canonicaljson.Value) {
		effect := value.(canonicaljson.Object)
		if setOf("filesystem_write", "git_ref_write", "git_index_write", "git_commit")[objectString(effect, "type")] {
			writes = append(writes, snapshot.Handoff.Target.Worktree)
			break
		}
	}
	sort.Strings(writes)
	writes = uniqueStrings(writes)
	if !equalStrings(writes, writeRoots) {
		return fmt.Errorf("sandbox write roots do not exactly match the contract")
	}
	return nil
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return values
	}
	result := values[:1]
	for _, value := range values[1:] {
		if value != result[len(result)-1] {
			result = append(result, value)
		}
	}
	return result
}

func validateObservedEffectsAgainstHandoff(handoff, result canonicaljson.Object) error {
	authorityValue, _ := objectMember(handoff, "authority")
	authority := authorityValue.(canonicaljson.Object)
	allowedValue, _ := objectMember(authority, "allowed_effects")
	allowed := allowedValue.([]canonicaljson.Value)
	reportedValue, _ := objectMember(result, "observed_effects")
	reported := reportedValue.([]canonicaljson.Value)
	if len(reported) != len(allowed) {
		return fmt.Errorf("observed effects do not cover every allowed effect")
	}
	for index, value := range allowed {
		expected := value.(canonicaljson.Object)
		actual := reported[index].(canonicaljson.Object)
		scope, _ := objectMember(expected, "scope")
		actualScope, _ := objectMember(actual, "scope")
		max := objectInt(expected, "max_occurrences")
		occurrences := objectInt(actual, "occurrences")
		within := occurrences <= max
		if objectString(actual, "effect_id") != objectString(expected, "id") || objectString(actual, "type") != objectString(expected, "type") || !canonicalEqual(scope, actualScope) || objectBool(actual, "within_authority") != within {
			return fmt.Errorf("observed effect %q does not match handoff authority", objectString(expected, "id"))
		}
	}
	return nil
}

func validateVerifierResultsAgainstHandoff(handoff, result canonicaljson.Object) error {
	verifiersValue, _ := objectMember(handoff, "verifiers")
	verifiers := verifiersValue.([]canonicaljson.Value)
	resultsValue, _ := objectMember(result, "verifier_results")
	results := resultsValue.([]canonicaljson.Value)
	if len(verifiers) != len(results) {
		return fmt.Errorf("verifier results do not cover every handoff verifier")
	}
	finalValue, _ := objectMember(result, "final_target")
	final := finalValue.(canonicaljson.Object)
	for index, value := range verifiers {
		verifier := value.(canonicaljson.Object)
		actual := results[index].(canonicaljson.Object)
		argv, _ := objectMember(verifier, "argv")
		actualArgv, _ := objectMember(actual, "argv")
		if objectString(actual, "verifier_id") != objectString(verifier, "id") || !canonicalEqual(argv, actualArgv) || objectString(actual, "cwd") != objectString(verifier, "cwd") {
			return fmt.Errorf("verifier result %q does not match the handoff", objectString(verifier, "id"))
		}
		evidenceValue, _ := objectMember(verifier, "evidence")
		evidence := evidenceValue.(canonicaljson.Object)
		bound := objectString(actual, "bound_oid_or_sha256")
		if objectString(evidence, "binding") == "target_oid" && bound != objectString(final, "oid") {
			return fmt.Errorf("verifier result %q has the wrong target binding", objectString(verifier, "id"))
		}
	}
	return nil
}

func equalManagedArtifacts(left, right []managedArtifact) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index].ID != right[index].ID || left[index].Locator != right[index].Locator || left[index].SHA256 != right[index].SHA256 || string(left[index].Bytes) != string(right[index].Bytes) {
			return false
		}
	}
	return true
}
