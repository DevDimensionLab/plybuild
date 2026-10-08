package workflowhandoff

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/devdimensionlab/plybuild/internal/canonicaljson"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// DeliveryHumanAttestation is authored from actual human input by the caller.
// Empty templates are not valid answers and cannot qualify a candidate.
type DeliveryHumanAttestation struct {
	Kind           string `json:"kind"`
	SchemaVersion  int    `json:"schema_version"`
	TaskID         string `json:"task_id"`
	TaskResultID   string `json:"task_result_id"`
	ResultOID      string `json:"result_oid"`
	ResultTree     string `json:"result_tree"`
	Outcome        string `json:"outcome"`
	ActorClaim     string `json:"actor_claim"`
	StartSurface   string `json:"start_surface"`
	StartedAtUTC   string `json:"started_at_utc"`
	CompletedAtUTC string `json:"completed_at_utc"`
	Answer         string `json:"answer"`
	Observation    string `json:"observation"`
}

// RecordDeliveryHumanQA preserves an explicit candidate-bound human answer.
// Actor claims have the same local provenance model as RecordTaskHumanQA; this
// API neither authenticates a human nor turns a technical pass into human QA.
func RecordDeliveryHumanQA(d Dependencies, taskID string, result workspace.TaskResultRecord, outcome, evidencePath string) (workspace.TaskHumanQARecord, error) {
	var empty workspace.TaskHumanQARecord
	if d.TaskWorkspace == nil {
		return empty, fmt.Errorf("Task workspace dependency missing")
	}
	raw, e := deliveryReadFile(d.Files, evidencePath, 256<<10)
	if e != nil {
		return empty, e
	}
	attestation, e := ValidateDeliveryHumanAttestation(raw, taskID, result, outcome)
	if e != nil {
		return empty, e
	}
	actor := workspace.TaskHumanActorRecord{ActorClaim: attestation.ActorClaim, StartSurface: attestation.StartSurface, StartedAtUTC: attestation.StartedAtUTC, CompletedAtUTC: attestation.CompletedAtUTC}
	s, e := d.Store.ReadByLocator(result.HandoffLocator)
	if e != nil {
		return empty, e
	}
	if taskSpecVersion(s.Handoff.Value) != 3 || deliveryMode(s.Handoff.Value) != "candidate" || s.Handoff.SHA256 != result.HandoffSHA256 {
		return empty, fmt.Errorf("human QA requires the qualified delivery candidate")
	}
	area := filepath.Join(s.Handoff.ReplyRoot, "staging", "human-qa-"+digestBytes(raw)[7:])
	if e = deliveryEnsureDir(d.Files, area); e != nil {
		return empty, e
	}
	path := filepath.Join(area, "attestation.json")
	if e = deliveryWriteOnce(d.Files, path, raw); e != nil {
		return empty, e
	}
	draft := deliveryObject(map[string]any{"kind": "WorkspaceTaskHumanQARecordDraft@1", "schema_version": 1, "format": "json", "format_version": 1, "canonicalization": "RFC8785", "publication_key": "delivery-qa/" + string(result.ID) + "/" + digestBytes(raw)[7:], "task_id": result.TaskID, "task_result_id": result.ID, "result_oid": result.ResultOID, "result_tree": result.ResultTree, "outcome": outcome, "actor": actor, "evidence": []workspace.TaskHumanQAEvidenceRecord{{ID: "human-attestation", Role: "report", Locator: path, SHA256: digestBytes(raw), SizeBytes: int64(len(raw))}}, "observation": attestation.Observation, "accepted_residual_risks": []workspace.TaskResidualRiskRecord{}})
	draftRaw, e := canonicaljson.Marshal(draft)
	if e != nil {
		return empty, e
	}
	draftPath := filepath.Join(area, "qa-draft.json")
	if e = deliveryWriteOnce(d.Files, draftPath, draftRaw); e != nil {
		return empty, e
	}
	wd := workspace.WithTaskContentScope(*d.TaskWorkspace, workspace.TaskID(taskID))
	wd.HandoffEvidence = NewTaskHandoffEvidenceReader(d)
	recorded, e := workspace.RecordTaskHumanQA(wd, workspace.TaskHumanQARecordInput{TaskID: result.TaskID, File: draftPath})
	return recorded.Record, e
}

// ValidateDeliveryHumanAttestation is read-only so a callback can reject an
// incomplete template before reserving an effect. It checks an explicit claim,
// not the identity of the human who supplied the answer.
func ValidateDeliveryHumanAttestation(raw []byte, taskID string, result workspace.TaskResultRecord, outcome string) (DeliveryHumanAttestation, error) {
	var empty DeliveryHumanAttestation
	if len(raw) > 256<<10 || taskID != string(result.TaskID) || !setOf("pass", "fail", "blocked")[outcome] {
		return empty, fmt.Errorf("human QA must name the exact Task and outcome")
	}
	v, e := canonicaljson.DecodeStrict(raw)
	if e != nil {
		return empty, e
	}
	m, e := exactObject(v, "human attestation", "kind", "schema_version", "task_id", "task_result_id", "result_oid", "result_tree", "outcome", "actor_claim", "start_surface", "started_at_utc", "completed_at_utc", "answer", "observation")
	if e != nil {
		return empty, e
	}
	version, e := intField(m, "schema_version", "human attestation")
	if e != nil || version != 1 || objectMapString(m, "kind") != "DeliveryHumanAttestation@1" || objectMapString(m, "task_id") != taskID || objectMapString(m, "task_result_id") != string(result.ID) || objectMapString(m, "result_oid") != result.ResultOID || objectMapString(m, "result_tree") != result.ResultTree || objectMapString(m, "outcome") != outcome || objectMapString(m, "answer") != outcome {
		return empty, fmt.Errorf("human answer is missing or belongs to a different candidate")
	}
	attestation := DeliveryHumanAttestation{Kind: objectMapString(m, "kind"), SchemaVersion: 1, TaskID: taskID, TaskResultID: string(result.ID), ResultOID: result.ResultOID, ResultTree: result.ResultTree, Outcome: outcome, ActorClaim: objectMapString(m, "actor_claim"), StartSurface: objectMapString(m, "start_surface"), StartedAtUTC: objectMapString(m, "started_at_utc"), CompletedAtUTC: objectMapString(m, "completed_at_utc"), Answer: outcome, Observation: objectMapString(m, "observation")}
	started, se := time.Parse(time.RFC3339Nano, attestation.StartedAtUTC)
	completed, ce := time.Parse(time.RFC3339Nano, attestation.CompletedAtUTC)
	if se != nil || ce != nil || completed.Before(started) || validatePlainText("actor", attestation.ActorClaim, 1, 256) != nil || validatePlainText("surface", attestation.StartSurface, 1, 256) != nil || validatePlainText("observation", attestation.Observation, 1, 2000) != nil {
		return empty, fmt.Errorf("human attestation provenance is incomplete")
	}
	return attestation, nil
}

type DeliveryOwnerAuthority struct {
	RunID, RequestSHA256, ActorClaim, PreparationID string
	Authorization                                   *workspace.DeliveryAuthorization
}

type DeliveryIntegrationResult struct {
	Integration workspace.TaskIntegrationResult
	Queue       workspace.WorkspaceTaskQueueReadback
	Base        workspace.EpicBaseResult
	Completed   bool
}

// IntegrateDeliveryCandidate reuses the native optimistic check/apply and
// one-attempt effect ledger. A moved parent stops; refreshing is a later feature.
func IntegrateDeliveryCandidate(d Dependencies, taskID string, resultID workspace.TaskResultID, qaID workspace.HumanQARecordID, expectedParentOID, resultOID string, owner DeliveryOwnerAuthority) (DeliveryIntegrationResult, error) {
	var out DeliveryIntegrationResult
	if d.TaskWorkspace == nil {
		return out, fmt.Errorf("Task workspace dependency missing")
	}
	wd := workspace.WithTaskContentScope(*d.TaskWorkspace, workspace.TaskID(taskID))
	wd.HandoffEvidence = NewTaskHandoffEvidenceReader(d)
	in, e := workspace.ParseTaskIntegrationInput(taskID, string(resultID), string(qaID), resultOID, expectedParentOID, "", false, "")
	if e != nil {
		return out, e
	}
	in.DeliveryOwner = &workspace.DeliveryIntegrationOwner{RunID: owner.RunID, RequestSHA256: owner.RequestSHA256, ActorClaim: owner.ActorClaim, PreparationID: owner.PreparationID}
	in.DeliveryAuthorization = owner.Authorization
	preview, e := workspace.CheckTaskIntegration(wd, in)
	out.Integration = preview
	if e != nil {
		return out, e
	}
	result := preview.Readback.TaskResult
	if result == nil || result.ResultOID != resultOID {
		return out, fmt.Errorf("integration candidate readback differs")
	}
	// Preserve the declared owner and the originally selected parent binding.
	s, e := d.Store.ReadByLocator(result.HandoffLocator)
	if e != nil {
		return out, e
	}
	if taskSpecVersion(s.Handoff.Value) != 3 || deliveryMode(s.Handoff.Value) != "candidate" {
		return out, fmt.Errorf("delivery integration requires candidate control evidence")
	}
	binding, _ := objectMember(s.Handoff.Value, "delivery_binding")
	b := binding.(canonicaljson.Object)
	parent, e := d.Store.ReadByLocator(objectString(b, "parent_handoff_locator"))
	if e != nil {
		return out, e
	}
	if parent.Handoff.SHA256 != objectString(b, "parent_handoff_sha256") || objectString(b, "owner_claim") != owner.ActorClaim || parent.Handoff.Target.OID != expectedParentOID {
		return out, fmt.Errorf("delivery owner or original parent binding differs")
	}
	if preview.Readback.RecoveryStatus != "complete" && preview.Readback.Classification != "ready" {
		return out, fmt.Errorf("delivery integration stopped: %s; changed parent requires refresh", preview.Readback.Classification)
	}
	plan, _ := objectMember(preview.Readback.Value, "plan")
	po, ok := plan.(canonicaljson.Object)
	if !ok {
		return out, fmt.Errorf("integration plan missing")
	}
	in.Apply = true
	in.Confirmation = objectString(po, "sha256")
	persistedAuthority, _ := objectMember(preview.Readback.Value, "authority")
	if preview.Readback.RecoveryStatus != "complete" || persistedAuthority == nil {
		out.Integration, e = workspace.ApplyTaskIntegration(wd, in)
		if e != nil {
			return out, e
		}
	}
	if out.Integration.Readback.RecoveryStatus != "complete" || out.Integration.Readback.ParentOID != resultOID {
		return out, fmt.Errorf("integration is not observed complete: %s", out.Integration.Readback.Classification)
	}
	qt := workspace.QueueTargetInput{ProjectID: result.ProjectID, RepoID: result.RepoID, EpicID: out.Integration.Readback.EpicID}
	out.Queue, e = CloseDeliveryTaskQueue(wd, qt, workspace.TaskID(taskID), owner.PreparationID, "Exact candidate integrated after recorded human pass.")
	if e != nil {
		return out, e
	}
	baseInput := workspace.EpicBaseInput{EpicID: qt.EpicID, RepoID: qt.RepoID}
	out.Base, e = workspace.UpdateEpicBase(wd, baseInput)
	if e != nil {
		return out, e
	}
	if out.Base.State != "unchanged" {
		if out.Base.Plan == nil || out.Base.Confirmation == nil || out.Base.Plan.NextOID != resultOID {
			return out, fmt.Errorf("Epic base update does not bind the integrated candidate")
		}
		baseInput.Apply = true
		baseInput.Confirmation = *out.Base.Confirmation
		out.Base, e = workspace.UpdateEpicBase(wd, baseInput)
		if e != nil {
			return out, e
		}
	}
	if out.Base.Freshness != "fresh" || out.Base.Observed == nil || out.Base.Observed.OID != resultOID || len(out.Base.Versions) == 0 || out.Base.Versions[len(out.Base.Versions)-1].OID != resultOID {
		return out, fmt.Errorf("Epic base update is not observed complete")
	}
	out.Completed = true
	return out, nil
}

// CloseDeliveryTaskQueue advances only its own preparation. On recovery, a
// prior native advance proves closure even if an unrelated Task is now current.
// An empty queue alone is never evidence that this Task was delivered.
func CloseDeliveryTaskQueue(d workspace.Dependencies, target workspace.QueueTargetInput, taskID workspace.TaskID, preparationID, reason string) (workspace.WorkspaceTaskQueueReadback, error) {
	d = workspace.WithTaskContentScope(d, taskID)
	out, err := workspace.ListTaskQueue(d, target, false)
	if err != nil {
		return out, err
	}
	if out.Current != nil && out.Current.PreparationID == preparationID && out.Current.TaskID == taskID {
		return workspace.CloseTaskQueue(d, target, preparationID, out.Revision, reason, "advance")
	}
	registry, err := d.WorkItems.Snapshot(out.Workspace.Root)
	if err != nil {
		return out, err
	}
	for _, event := range registry.TaskQueueEvents {
		if event.QueueID == out.QueueID && event.Kind == "advance" && event.Request["preparation_id"] == preparationID {
			return out, nil
		}
	}
	return out, fmt.Errorf("the exact delivery preparation is neither current nor observed advanced; preserve unrelated queue work and inspect its native history")
}
