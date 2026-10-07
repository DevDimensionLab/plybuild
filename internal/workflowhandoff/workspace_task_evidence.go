package workflowhandoff

import (
	"fmt"
	"sort"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type workspaceTaskEvidenceReader struct{ dependencies Dependencies }

func NewTaskHandoffEvidenceReader(dependencies Dependencies) workspace.TaskHandoffEvidenceReader {
	return &workspaceTaskEvidenceReader{dependencies: dependencies}
}

func (reader *workspaceTaskEvidenceReader) ReadTaskEvidence(request workspace.TaskHandoffEvidenceRequest) (workspace.TaskHandoffEvidence, error) {
	if err := requireDependencies(reader.dependencies); err != nil {
		return workspace.TaskHandoffEvidence{}, err
	}
	snapshot, err := reader.dependencies.Store.ReadByLocator(request.HandoffLocator)
	if err != nil {
		return workspace.TaskHandoffEvidence{}, err
	}
	if snapshot.Start == nil || snapshot.Terminal == nil {
		return workspace.TaskHandoffEvidence{}, fmt.Errorf("handoff does not have accepted start and terminal evidence")
	}
	if snapshot.Handoff.Identity.ActivityID != request.ActivityID || snapshot.Handoff.Identity.RunID != request.RunID || snapshot.Handoff.Identity.HandoffID != request.HandoffID || snapshot.Handoff.SHA256 != request.HandoffSHA256 || snapshot.Start.DocumentID != request.StartReceiptID || snapshot.Start.Locator != request.StartReceiptLocator || snapshot.Start.SHA256 != request.StartReceiptSHA256 || snapshot.Terminal.DocumentID != request.TerminalResultID || snapshot.Terminal.Locator != request.TerminalResultLocator || snapshot.Terminal.SHA256 != request.TerminalResultSHA256 {
		return workspace.TaskHandoffEvidence{}, fmt.Errorf("handoff document identities, locators, or digests differ from the request")
	}
	inspection, err := Inspect(reader.dependencies, InspectInput{HandoffLocator: request.HandoffLocator, Format: "json"})
	if err != nil {
		return workspace.TaskHandoffEvidence{}, err
	}
	inspectionDigest := digestBytes(inspection.Bytes)
	if inspectionDigest != request.InspectionSHA256 {
		return workspace.TaskHandoffEvidence{}, fmt.Errorf("inspection digest differs from the request")
	}
	state, policyReasons, evidenceReasons := deriveState(snapshot)
	targetValue, _ := objectMember(snapshot.Terminal.Value, "final_target")
	target, _ := targetValue.(canonicaljson.Object)
	result := workspace.TaskHandoffEvidence{
		ActivityID: snapshot.Handoff.Identity.ActivityID, RunID: snapshot.Handoff.Identity.RunID, HandoffID: snapshot.Handoff.Identity.HandoffID, HandoffLocator: snapshot.Handoff.Locator, HandoffSHA256: snapshot.Handoff.SHA256,
		StartReceiptID: snapshot.Start.DocumentID, StartReceiptLocator: snapshot.Start.Locator, StartReceiptSHA256: snapshot.Start.SHA256, TerminalResultID: snapshot.Terminal.DocumentID, TerminalResultLocator: snapshot.Terminal.Locator, TerminalResultSHA256: snapshot.Terminal.SHA256, InspectionSHA256: inspectionDigest, ReportedOutcome: snapshot.Terminal.Outcome,
		TargetWorktree: objectString(target, "worktree"), TargetRef: objectString(target, "ref"), ResultOID: objectString(target, "oid"), ResultTree: objectString(target, "tree"), GitCommonDir: objectString(target, "git_common_dir"),
		SchemaValid: true, DigestValid: true, LifecycleValid: true, BindingValid: true, CapabilityValid: true, PrincipalSessionValid: true, PolicyValid: len(policyReasons) == 0, EvidenceCoverageValid: len(evidenceReasons) == 0,
		EvidenceGaps: append([]string(nil), evidenceReasons...),
	}
	_ = state
	result.ExpectedVerifierIDs = workspaceExpectedVerifierIDs(snapshot)
	result.VerifierResults = workspaceVerifierProjection(snapshot)
	result.Review = workspaceReviewProjection(snapshot)
	result.Artifacts = workspaceArtifactProjection(snapshot, result.VerifierResults, result.Review)
	result.AcceptedStartOutcome = snapshot.Start.Outcome
	basis, eval, specErr := historicalHandoffTaskSpec(reader.dependencies, snapshot)
	result.TaskSpecBasis = basis
	result.TaskSpecValid = specErr == nil
	if specErr != nil {
		result.TaskSpecValid = false
	}
	if coverageErr := taskRequirementsCoverage(reader.dependencies, snapshot, basis, eval); coverageErr != nil {
		result.EvidenceCoverageValid = false
		result.EvidenceGaps = append(result.EvidenceGaps, "task_spec_requirement_evidence: "+coverageErr.Error())
	} else {
		result.TaskRequirementsValid = true
	}
	if reader.dependencies.DeliveryAuthorization != nil {
		result.DeliveryAuthorization, err = reader.dependencies.DeliveryAuthorization(request)
		if err != nil {
			return workspace.TaskHandoffEvidence{}, err
		}
	}
	return result, nil
}

func workspaceExpectedVerifierIDs(snapshot Snapshot) []string {
	value, ok := objectMember(snapshot.Handoff.Value, "verifiers")
	if !ok {
		return []string{}
	}
	items, _ := value.([]canonicaljson.Value)
	ids := make([]string, 0, len(items))
	for _, item := range items {
		object, _ := item.(canonicaljson.Object)
		ids = append(ids, objectString(object, "id"))
	}
	sort.Strings(ids)
	return ids
}

func workspaceVerifierProjection(snapshot Snapshot) []workspace.TaskVerifierResultRecord {
	value, ok := objectMember(snapshot.Terminal.Value, "verifier_results")
	if !ok {
		return []workspace.TaskVerifierResultRecord{}
	}
	list, _ := value.([]canonicaljson.Value)
	result := make([]workspace.TaskVerifierResultRecord, 0, len(list))
	expected := map[string]int64{}
	if raw, ok := objectMember(snapshot.Handoff.Value, "verifiers"); ok {
		for _, v := range raw.([]canonicaljson.Value) {
			o := v.(canonicaljson.Object)
			exit, _ := objectMember(o, "expected_exit")
			expected[objectString(o, "id")], _ = exit.(int64)
		}
	}
	for _, item := range list {
		o := item.(canonicaljson.Object)
		argvValue, _ := objectMember(o, "argv")
		argvRaw, _ := argvValue.([]canonicaljson.Value)
		argv := make([]string, len(argvRaw))
		for i, v := range argvRaw {
			argv[i], _ = v.(string)
		}
		exitValue, _ := objectMember(o, "exit")
		exit, _ := exitValue.(int64)
		stdout := nullableStringPointer(o, "stdout_artifact_id")
		stderr := nullableStringPointer(o, "stderr_artifact_id")
		id := objectString(o, "verifier_id")
		outcome := "failed"
		if want, ok := expected[id]; ok && want == exit {
			outcome = "passed"
		}
		result = append(result, workspace.TaskVerifierResultRecord{VerifierID: id, Outcome: outcome, Argv: argv, CWD: objectString(o, "cwd"), Exit: int(exit), BoundOIDOrSHA256: objectString(o, "bound_oid_or_sha256"), StdoutArtifactID: stdout, StderrArtifactID: stderr})
	}
	return result
}

func workspaceReviewProjection(snapshot Snapshot) workspace.TaskReviewRecord {
	value, ok := objectMember(snapshot.Terminal.Value, "review")
	if !ok {
		return workspace.TaskReviewRecord{Findings: []workspace.TaskReviewEntry{}, Fixes: []workspace.TaskReviewEntry{}, OpenActionableFindings: []workspace.TaskReviewEntry{}}
	}
	review := value.(canonicaljson.Object)
	return workspace.TaskReviewRecord{Findings: workspaceReviewEntries(review, "findings"), Fixes: workspaceReviewEntries(review, "fixes"), OpenActionableFindings: workspaceReviewEntries(review, "open_actionable_findings")}
}
func workspaceReviewEntries(review canonicaljson.Object, name string) []workspace.TaskReviewEntry {
	value, _ := objectMember(review, name)
	list, _ := value.([]canonicaljson.Value)
	result := make([]workspace.TaskReviewEntry, 0, len(list))
	for _, item := range list {
		o := item.(canonicaljson.Object)
		evidenceValue, _ := objectMember(o, "evidence_ids")
		raw, _ := evidenceValue.([]canonicaljson.Value)
		ids := make([]string, len(raw))
		for i, v := range raw {
			ids[i], _ = v.(string)
		}
		result = append(result, workspace.TaskReviewEntry{ID: objectString(o, "id"), Severity: objectString(o, "severity"), Summary: objectString(o, "summary"), EvidenceArtifactIDs: ids})
	}
	return result
}

func workspaceArtifactProjection(snapshot Snapshot, verifiers []workspace.TaskVerifierResultRecord, review workspace.TaskReviewRecord) []workspace.TaskArtifactRecord {
	roles := map[string]string{}
	for _, v := range verifiers {
		if v.StdoutArtifactID != nil {
			roles[*v.StdoutArtifactID] = "verifier_stdout"
		}
		if v.StderrArtifactID != nil {
			roles[*v.StderrArtifactID] = "verifier_stderr"
		}
	}
	for _, list := range [][]workspace.TaskReviewEntry{review.Findings, review.Fixes, review.OpenActionableFindings} {
		for _, entry := range list {
			for _, id := range entry.EvidenceArtifactIDs {
				roles[id] = "review"
			}
		}
	}
	value, _ := objectMember(snapshot.Terminal.Value, "artifacts")
	list, _ := value.([]canonicaljson.Value)
	result := []workspace.TaskArtifactRecord{}
	for _, item := range list {
		o := item.(canonicaljson.Object)
		if objectString(o, "kind") != "managed" {
			continue
		}
		sizeValue, _ := objectMember(o, "size_bytes")
		size, _ := sizeValue.(int64)
		id := objectString(o, "artifact_id")
		role := roles[id]
		if role == "" {
			role = "other"
		}
		result = append(result, workspace.TaskArtifactRecord{ArtifactID: id, Role: role, Locator: objectString(o, "locator"), SHA256: objectString(o, "sha256"), SizeBytes: size})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ArtifactID < result[j].ArtifactID })
	return result
}
func nullableStringPointer(object canonicaljson.Object, name string) *string {
	value, ok := objectMember(object, name)
	if !ok || value == nil {
		return nil
	}
	text, ok := value.(string)
	if !ok {
		return nil
	}
	return &text
}
