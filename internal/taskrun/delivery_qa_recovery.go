package taskrun

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

type DeliveryQARecoveryInput struct {
	RunID, ContextPath, RuntimeEvidencePath, AttemptID string
}

type DeliveryQARecoveryPreview struct {
	Kind          string       `json:"kind"`
	SchemaVersion int          `json:"schema_version"`
	RunID         string       `json:"run_id"`
	AttemptID     string       `json:"attempt_id"`
	State         string       `json:"state"`
	Confirmation  string       `json:"confirmation"`
	NextAction    string       `json:"next_action"`
	Receipt       *FileBinding `json:"receipt,omitempty"`
}

func WorkflowPreviewDeliveryQARecovery(d Dependencies, root string, in DeliveryQARecoveryInput) (DeliveryQARecoveryPreview, error) {
	return deliveryQARecovery(d, root, in, "", false)
}

func WorkflowRecoverDeliveryQA(d Dependencies, root string, in DeliveryQARecoveryInput, confirmation string) (DeliveryQARecoveryPreview, error) {
	return deliveryQARecovery(d, root, in, confirmation, true)
}

type deliveryQARequest struct {
	Kind          string                 `json:"kind"`
	RunID         string                 `json:"run_id"`
	RequestSHA256 string                 `json:"request_sha256"`
	Candidate     workspace.TaskResultID `json:"candidate"`
	CandidateOID  string                 `json:"candidate_oid"`
	CandidateTree string                 `json:"candidate_tree"`
	Outcome       string                 `json:"outcome"`
	Evidence      FileBinding            `json:"evidence"`
}

type deliveryQARejectionReceipt struct {
	Kind               string            `json:"kind"`
	SchemaVersion      int               `json:"schema_version"`
	RunID              string            `json:"run_id"`
	RequestSHA256      string            `json:"request_sha256"`
	Context            FileBinding       `json:"context"`
	Acceptance         FileBinding       `json:"acceptance"`
	Attempt            DeliveryAttempt   `json:"attempt"`
	Candidate          DeliveryCandidate `json:"candidate"`
	Request            FileBinding       `json:"request"`
	Evidence           FileBinding       `json:"evidence"`
	StagedAttestation  FileBinding       `json:"staged_attestation"`
	StagedDraft        FileBinding       `json:"staged_draft"`
	PublicationKey     string            `json:"publication_key"`
	RuntimeObservation FileBinding       `json:"runtime_observation"`
	OriginalControl    Executable        `json:"original_control"`
	RepairExecutable   Executable        `json:"repair_executable"`
	RegistrySHA256     string            `json:"registry_sha256"`
	Reason             string            `json:"reason"`
}

// This is a separate installed maintenance operation. It never changes the
// callback control, publishes QA, or retries an effect. Only the exact known
// schema rejection is provably incapable of reaching registry publication.
func deliveryQARecovery(d Dependencies, root string, in DeliveryQARecoveryInput, confirmation string, apply bool) (out DeliveryQARecoveryPreview, err error) {
	out = DeliveryQARecoveryPreview{Kind: "PlyDeliveryQARecoveryPreview@1", SchemaVersion: 1, RunID: in.RunID, AttemptID: in.AttemptID, State: "blocked"}
	defer func() {
		if err != nil {
			out.State = "blocked"
			out.NextAction = err.Error()
		}
	}()
	if err = containing(d, root); err != nil {
		return out, err
	}
	if !strings.HasPrefix(in.AttemptID, "qa-") || !digestPattern.MatchString("sha256:"+strings.TrimPrefix(in.AttemptID, "qa-")) {
		return out, workflowError(2, "--attempt must name the exact original QA reservation")
	}
	if apply && !digestPattern.MatchString(confirmation) {
		return out, workflowError(2, "recovery requires the exact preview --confirm digest")
	}
	err = withStore(root, func() error {
		s, e := workflowRead(root, in.RunID)
		if e != nil {
			return e
		}
		if in.ContextPath == "" || in.ContextPath != s.Result.Paths.Context || !deliveryRun(s.Request) || s.Result.Delivery == nil {
			return workflowError(4, "recovery requires the exact delivery run and original context")
		}
		eventID := "recover-" + in.AttemptID
		// A completed repair is a read-only receipt lookup, even if a later QA
		// or integration is now in flight. Never retire that later reservation.
		for _, event := range s.Result.Delivery.Events {
			if event.ID != eventID {
				continue
			}
			var report DeliveryReport
			raw, e := workflowBound(event.Binding, 4<<20)
			if e != nil {
				return e
			}
			if e = decode(raw, 4<<20, &report); e != nil {
				return e
			}
			if event.Kind != "report" || report.EventID != eventID || validateDeliveryReport(report) != nil || workflowClaim(s, report.RunID, report.RequestSHA256, report.SessionID) != nil || len(report.Evidence) != 1 {
				return workflowError(4, "recovery report changed")
			}
			var receipt deliveryQARejectionReceipt
			raw, e = workflowBound(report.Evidence[0], 8<<20)
			if e != nil {
				return e
			}
			if e = decode(raw, 8<<20, &receipt); e != nil {
				return e
			}
			path := filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", in.AttemptID, "rejection-recovery.json")
			if receipt.Kind != "PlyDeliveryQARejectionRecovery@1" || receipt.SchemaVersion != 1 || receipt.RunID != in.RunID || receipt.RequestSHA256 != s.Result.RequestSHA256 || receipt.Attempt.ID != in.AttemptID || receipt.Context != (FileBinding{s.Result.Paths.Context, s.ContextSHA256}) || report.Evidence[0].Locator != path {
				return workflowError(4, "recovery receipt binding changed")
			}
			out.State, out.Confirmation, out.Receipt = "existing", report.Evidence[0].SHA256, &report.Evidence[0]
			out.NextAction = "The exact rejected QA reservation was already retired; preserve any later QA or integration attempt."
			return nil
		}
		observation, repair, control, e := deliveryQARecoveryOwner(d, s, in)
		if e != nil {
			return e
		}
		if len(s.Result.Delivery.Candidates) == 0 {
			return workflowError(4, "recovery requires an existing qualified candidate")
		}
		c := s.Result.Delivery.Candidates[len(s.Result.Delivery.Candidates)-1]
		a := s.Result.Delivery.Attempt
		if a == nil || a.ID != in.AttemptID || a.Kind != "human_qa" || a.State != "attempted" || a.Path != filepath.Join(s.Result.Paths.RunRoot, "delivery", "attempts", in.AttemptID) || a.CandidateOID != c.OID || a.CandidateTree != c.Tree || c.HumanQA != nil || c.Integration != nil || c.PullRequest != nil {
			return workflowError(4, "recovery requires only the current unpublished QA attempt on the unchanged candidate")
		}
		requestBinding, raw, e := deliveryReadBinding(filepath.Join(a.Path, "attempt.json"), 1<<20)
		if e != nil {
			return e
		}
		var request deliveryQARequest
		if e = decode(raw, 1<<20, &request); e != nil {
			return e
		}
		inputSHA := digest([]any{c.TaskResult.ID, request.Outcome, request.Evidence})
		if request.Kind != "PlyDeliveryHumanQARequest@1" || request.RunID != in.RunID || request.RequestSHA256 != s.Result.RequestSHA256 || request.Candidate != c.TaskResult.ID || request.CandidateOID != c.OID || request.CandidateTree != c.Tree || a.InputSHA256 != inputSHA || a.ID != "qa-"+inputSHA[7:] {
			return workflowError(4, "original QA request does not prove this reservation")
		}
		if _, e = workflowBound(request.Evidence, 1<<20); e != nil {
			return e
		}
		proof, e := workflowhandoff.ValidateRejectedDeliveryHumanQA(d.Workflow, string(c.TaskResult.TaskID), c.TaskResult, request.Outcome, request.Evidence.Locator)
		if e != nil {
			return e
		}
		attestation, _, e := deliveryReadBinding(proof.AttestationPath, 1<<20)
		if e != nil {
			return e
		}
		draft, _, e := deliveryReadBinding(proof.DraftPath, 1<<20)
		if e != nil {
			return e
		}
		evidence, e := workflowhandoff.NewTaskHandoffEvidenceReader(d.Workflow).ReadTaskEvidence(taskResultEvidenceRequest(c.TaskResult))
		if e != nil {
			return e
		}
		if !evidence.TaskSpecValid || !evidence.TaskRequirementsValid || c.TaskResult.TechnicalGate != "passed" && c.TaskResult.TechnicalGate != "good_enough_with_known_debt" {
			return workflowError(4, "candidate technical evidence changed")
		}
		if e = workspace.ValidateTaskResultEvidence(c.TaskResult, evidence); e != nil {
			return e
		}
		// Keep the native TaskRun -> project -> work-items lock order. The
		// state save occurs INSIDE the absence-check lock, not after its return.
		return d.Workspace.ProjectLocks.WithSnapshotLock(root, func(_ workspace.ProjectSnapshot) error {
			return d.Workspace.WorkItems.WithLock(root, func(session workspace.WorkItemStoreSession) error {
				registry, e := session.Snapshot()
				if e != nil {
					return e
				}
				found := false
				for _, r := range registry.TaskResults {
					if r.ID == c.TaskResult.ID {
						found = equal(r, c.TaskResult)
					}
				}
				if !found {
					return workflowError(4, "qualified TaskResult is absent or changed")
				}
				for _, r := range registry.HumanQARecords {
					if r.TaskResultID == c.TaskResult.ID || r.PublicationKey == proof.PublicationKey {
						return workflowError(4, "QA publication already exists; observe it instead of retiring an unknown effect")
					}
				}
				for _, r := range registry.IntegrationAuthorities {
					if r.TaskID == c.TaskResult.TaskID {
						return workflowError(4, "native integration authority already exists")
					}
				}
				for _, r := range registry.IntegrationIntents {
					if r.TaskID == c.TaskResult.TaskID {
						return workflowError(4, "native integration intent already exists")
					}
				}
				x, e := deliveryTarget(d, s)
				if e != nil {
					return e
				}
				if x.OID != c.OID || x.Tree != c.Tree || x.Locator != s.Observed.Target.WorktreeLocator {
					return workflowError(4, "candidate changed before recovery")
				}
				r := deliveryQARejectionReceipt{Kind: "PlyDeliveryQARejectionRecovery@1", SchemaVersion: 1, RunID: in.RunID, RequestSHA256: s.Result.RequestSHA256, Context: FileBinding{s.Result.Paths.Context, s.ContextSHA256}, Acceptance: *s.Acceptance, Attempt: *a, Candidate: c, Request: requestBinding, Evidence: request.Evidence, StagedAttestation: attestation, StagedDraft: draft, PublicationKey: proof.PublicationKey, RuntimeObservation: observation, OriginalControl: control, RepairExecutable: repair, RegistrySHA256: registry.RawSHA256, Reason: "The immutable staged draft is rejected before workspace publication solely because zero-offset timestamps use +00:00 instead of literal Z; no candidate QA or integration publication exists."}
				path := filepath.Join(a.Path, "rejection-recovery.json")
				var prior deliveryQARejectionReceipt
				if e = readValue(path, 8<<20, &prior); e == nil {
					// A lost reply can leave the receipt before state retirement.
					// Preserve that original observation; require a fresh current
					// observation separately on every unfinished recovery attempt.
					if e = deliveryRuntimeObservation(s, prior.RuntimeObservation); e != nil {
						return e
					}
					r.RuntimeObservation, r.RegistrySHA256 = prior.RuntimeObservation, prior.RegistrySHA256
					if !equal(r, prior) {
						return workflowError(4, "preserved recovery receipt differs from the exact rejected attempt")
					}
				} else if !os.IsNotExist(e) {
					return e
				}
				binding := FileBinding{path, digest(r)}
				out.State, out.Confirmation, out.Receipt = "ready", binding.SHA256, &binding
				out.NextAction = "Confirm this exact recovery to retire only the proven rejected reservation; then record the preserved actual answer through the unchanged bound control with a new timestamp-corrected evidence file."
				if !apply {
					return nil
				}
				if confirmation != binding.SHA256 {
					return workflowError(4, "recovery confirmation differs from the current exact preview")
				}
				if _, e = workflowKeep(d, path, r); e != nil {
					return e
				}
				if e = d.fault("delivery_qa_recovery_after_receipt"); e != nil {
					return e
				}
				report := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: in.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, EventID: eventID, PreviousEventSHA256: s.Result.Delivery.LastEventSHA256, Phase: "working", Summary: "Retired the proven pre-publication QA validation rejection.", Meaning: "The original human answer and candidate are unchanged. Recovery records no human verdict and grants no integration; use the ordinary bound QA callback with a new equivalent timestamp-corrected attestation.", Evidence: []FileBinding{binding}, VerifierResults: []DeliveryVerifierResult{}}
				if e = validateDeliveryReport(report); e != nil {
					return e
				}
				if _, e = deliveryAppendEvent(d, &s, eventID, "report", report); e != nil {
					return e
				}
				if e = d.fault("delivery_qa_recovery_after_event"); e != nil {
					return e
				}
				s.Result.Delivery.Attempt = nil
				s.Result.Delivery.Phase, s.Result.Round.State = "awaiting_human_qa", "awaiting_human_qa"
				out.NextAction = "The rejected reservation is retired. Record the preserved actual human answer through the unchanged bound control's ordinary QA callback with a new equivalent Z-timestamp attestation."
				s.Result.NextAction = WorkflowAction{"recipient", out.NextAction}
				if e = workflowSave(d, s); e != nil {
					return e
				}
				if e = d.fault("delivery_qa_recovery_after_state"); e != nil {
					return e
				}
				out.State = "recovered"
				return nil
			})
		})
	})
	return out, err
}

func deliveryQARecoveryOwner(d Dependencies, s workflowState, in DeliveryQARecoveryInput) (FileBinding, Executable, Executable, error) {
	var observation FileBinding
	var repair, control Executable
	fail := func(e error) (FileBinding, Executable, Executable, error) { return observation, repair, control, e }
	if e := deliveryContinuationAuthority(s); e != nil {
		return fail(e)
	}
	if _, e := deliveryAuthorityState(d, s); e != nil {
		return fail(e)
	}
	b, _, e := deliveryReadBinding(in.RuntimeEvidencePath, 1<<20)
	if e != nil {
		return fail(e)
	}
	observation = b
	if e = deliveryRuntimeObservation(s, b); e != nil {
		return fail(e)
	}
	if e = deliveryRuntimeObservationFresh(d, b); e != nil {
		return fail(e)
	}
	runtime, e := workflowEffectiveRuntime(s)
	if e != nil {
		return fail(e)
	}
	control = runtime.PlyExecutable
	if e = workflowFreshRuntimeOperation(d, s, false, runtime, true); e != nil {
		return fail(e)
	}
	if e = workflowhandoff.ValidateDeliveryOwnerContinuity(d.Workflow, s.Result.Handoff.Locator); e != nil {
		return fail(e)
	}
	cwd, e := d.CWD()
	if e != nil {
		return fail(e)
	}
	if cwd != s.Observed.Target.WorktreeLocator {
		return fail(workflowError(4, "recovery requires the same original Task cwd"))
	}
	if e = physical(cwd, false); e != nil {
		return fail(e)
	}
	if e = deliveryLiveOwner(d, s); e != nil {
		return fail(e)
	}
	path, e := d.Executable()
	if e == nil {
		path, e = filepath.EvalSymlinks(path)
	}
	if e != nil {
		return fail(e)
	}
	raw, e := readFile(path, 256<<20, false)
	if e != nil {
		return fail(e)
	}
	repair = Executable{Path: path, SHA256: hash(raw)}
	if e = verifyExecutable(repair); e != nil {
		return fail(e)
	}
	return observation, repair, control, nil
}
