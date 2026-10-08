package taskrun

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

const deliveryContinuityLegacyCapability = "delivery-owner-v2/task-spec-v3/task-content-v2/scoped-publications-v1"
const deliveryContinuityCapability = deliveryContinuityLegacyCapability + "/accepted-launch-history-v1"

// DeliveryContinuationReaderCapability names this reader's bounded transition
// policy. It is compatibility evidence, never runtime permission evidence.
const DeliveryContinuationReaderCapability = deliveryContinuityCapability

type DeliveryContinuationInput struct {
	RunID               string
	ContextPath         string
	ControlExecutable   Executable
	RuntimeEvidencePath string
}

type DeliveryContinuationPreview struct {
	Kind                  string                       `json:"kind"`
	SchemaVersion         int                          `json:"schema_version"`
	RunID                 string                       `json:"run_id"`
	State                 string                       `json:"state"`
	Generation            int                          `json:"generation"`
	Directory             string                       `json:"directory"`
	ControlExecutable     Executable                   `json:"control_executable"`
	Context               FileBinding                  `json:"context"`
	Proof                 *FileBinding                 `json:"proof,omitempty"`
	Confirmation          string                       `json:"confirmation"`
	ReaderCapability      string                       `json:"reader_capability"`
	NextAction            string                       `json:"next_action"`
	HistoricalLauncher    *DeliveryLauncherObservation `json:"historical_launcher,omitempty"`
	RuntimeObservation    *FileBinding                 `json:"runtime_observation,omitempty"`
	RuntimeRefreshAllowed bool                         `json:"runtime_refresh_allowed,omitempty"`
}

type deliveryContinuationRecord struct {
	Kind                    string       `json:"kind"`
	SchemaVersion           int          `json:"schema_version"`
	Generation              int          `json:"generation"`
	RunID                   string       `json:"run_id"`
	RequestSHA256           string       `json:"request_sha256"`
	SessionID               string       `json:"session_id"`
	NativeSessionID         string       `json:"native_session_id"`
	BeforeState             FileBinding  `json:"before_state"`
	SourceControl           Executable   `json:"source_control"`
	PreviousControl         Executable   `json:"previous_control"`
	ControlExecutable       Executable   `json:"control_executable"`
	Context                 FileBinding  `json:"context"`
	Acceptance              FileBinding  `json:"acceptance"`
	NativeStartSHA256       string       `json:"native_start_sha256"`
	ActualPolicySHA256      string       `json:"actual_policy_sha256"`
	DeliveryAgreementSHA256 string       `json:"delivery_agreement_sha256,omitempty"`
	ReaderCapability        string       `json:"reader_capability"`
	Confirmation            string       `json:"confirmation"`
	RecordedAtUTC           string       `json:"recorded_at_utc"`
	RuntimeObservation      *FileBinding `json:"runtime_observation,omitempty"`
}

func deliveryContinuityError(code, detail string) error {
	return &Error{Code: "delivery_" + code, Detail: detail, Exit: 4}
}

func deliveryContinuationGeneration(s workflowState) int {
	if s.Continuation == nil {
		return 0
	}
	n, _ := strconv.Atoi(filepath.Base(filepath.Dir(s.Continuation.Locator)))
	return n
}

func deliveryContinuationDirectory(s workflowState, generation int) string {
	return filepath.Join(s.Result.Paths.RunRoot, "delivery", "continuations", fmt.Sprintf("%03d", generation))
}

// Every transition retains the previous state and control. The runtime resolver
// changes only the control executable; provider, accepted policy, native owner,
// context, mandate and request retain their original identities.
func workflowContinuationRecord(s workflowState) (*deliveryContinuationRecord, error) {
	return workflowContinuationRecordMode(s, true)
}

func workflowContinuationRecordMode(s workflowState, policy bool) (*deliveryContinuationRecord, error) {
	if s.Continuation == nil {
		return nil, nil
	}
	n := deliveryContinuationGeneration(s)
	dir := deliveryContinuationDirectory(s, n)
	if n < 1 || n > 128 || s.Continuation.Locator != filepath.Join(dir, "continuation.json") {
		return nil, deliveryContinuityError("continuation_integrity", "continuation generation or proof path changed; inspect the preserved continuation")
	}
	raw, err := workflowBound(*s.Continuation, 1<<20)
	if err != nil {
		return nil, deliveryContinuityError("continuation_integrity", err.Error())
	}
	var r deliveryContinuationRecord
	if err = decode(raw, 1<<20, &r); err != nil {
		return nil, deliveryContinuityError("continuation_integrity", err.Error())
	}
	if r.Kind != "PlyDeliveryContinuation@1" || r.SchemaVersion != 1 || r.Generation != n || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.SessionID != s.Result.SessionID || r.NativeSessionID != s.Result.Transport.AgentSessionID || r.ReaderCapability != deliveryContinuityCapability && r.ReaderCapability != deliveryContinuityLegacyCapability || r.BeforeState.Locator != filepath.Join(dir, "before-state.json") || r.ControlExecutable.Path != filepath.Join(dir, "ply-control") || r.ControlExecutable.SHA256 != r.SourceControl.SHA256 || r.Context != (FileBinding{s.Result.Paths.Context, s.ContextSHA256}) || s.Acceptance == nil || r.Acceptance != *s.Acceptance || s.StartSHA256 == nil || r.NativeStartSHA256 != *s.StartSHA256 || s.Result.Delivery == nil || s.Result.Delivery.ActualPolicySHA256 == nil || r.ActualPolicySHA256 != *s.Result.Delivery.ActualPolicySHA256 {
		return nil, deliveryContinuityError("continuation_integrity", "continuation proof differs from the same accepted run, context, session or supported reader capability")
	}
	if _, err = time.Parse(time.RFC3339Nano, r.RecordedAtUTC); err != nil {
		return nil, deliveryContinuityError("continuation_integrity", "continuation timestamp is invalid")
	}
	agreement := ""
	if s.Request.Delivery != nil && s.Request.Delivery.Agreement != nil {
		agreement = workspace.DeliveryAgreementDigest(*s.Request.Delivery.Agreement)
	}
	if r.DeliveryAgreementSHA256 != agreement {
		return nil, deliveryContinuityError("continuation_authority", "continuation delivery agreement changed; no new authority is inferred")
	}
	before, err := workflowBound(r.BeforeState, 8<<20)
	if err != nil {
		return nil, deliveryContinuityError("continuation_integrity", err.Error())
	}
	var old workflowState
	if err = decode(before, 8<<20, &old); err != nil {
		return nil, err
	}
	if deliveryContinuationGeneration(old) != n-1 || !equal(old.Request, s.Request) || !equal(old.Observed, s.Observed) || old.Result.Handoff != s.Result.Handoff || old.Result.RunID != s.Result.RunID || old.Result.SessionID != s.Result.SessionID || old.Result.Transport.AgentSessionID != s.Result.Transport.AgentSessionID || old.Result.Transport.TabID != s.Result.Transport.TabID || old.Result.Transport.PaneID != s.Result.Transport.PaneID || old.Result.Transport.TerminalID != s.Result.Transport.TerminalID || old.Result.Paths.Context != s.Result.Paths.Context || old.ContextSHA256 != s.ContextSHA256 || !equal(old.Acceptance, s.Acceptance) || !equal(old.StartSHA256, s.StartSHA256) {
		return nil, deliveryContinuityError("continuation_integrity", "continuation history changed Task, worktree, mandate, context, accepted authority or session ownership")
	}
	if policy {
		if err = deliveryContinuationAuthority(old); err != nil {
			return nil, err
		}
	}
	if !equal(old.Result.Delivery.ActualPolicyEvidence, s.Result.Delivery.ActualPolicyEvidence) {
		return nil, deliveryContinuityError("continuation_authority", "actual policy evidence differs from the accepted transition history")
	}
	prior, err := workflowEffectiveRuntimeMode(old, policy)
	if err != nil {
		return nil, err
	}
	if prior.PlyExecutable != r.PreviousControl {
		return nil, deliveryContinuityError("continuation_integrity", "prior control binding differs from the immutable transition history")
	}
	if r.RuntimeObservation != nil {
		if r.ReaderCapability != deliveryContinuityCapability {
			return nil, deliveryContinuityError("continuation_unsupported", "historical reader proof cannot carry a current runtime observation")
		}
		if err = deliveryRuntimeObservationBinding(s, *r.RuntimeObservation, policy); err != nil {
			return nil, err
		}
	}
	if r.Confirmation != deliveryContinuationConfirmation(old, r.SourceControl, r.ControlExecutable, r.ReaderCapability, r.RuntimeObservation) {
		return nil, deliveryContinuityError("continuation_integrity", "compatibility proof confirmation differs from its preserved observations")
	}
	for _, control := range []Executable{r.PreviousControl, r.ControlExecutable} {
		if err = verifyExecutable(control); err != nil {
			return nil, deliveryContinuityError("continuation_integrity", "preserved continuation control: "+err.Error())
		}
	}
	return &r, nil
}

func deliveryContinuationAuthority(s workflowState) error {
	if s.Result.Delivery != nil && (s.Result.Delivery.OwnershipRelease != nil || s.Result.Delivery.Phase == "completed") {
		return deliveryContinuityError("continuation_authority", "continue requires the active owner; released or completed mandates cannot acquire new authority")
	}
	return deliveryAcceptedAuthority(s)
}

func deliveryAcceptedAuthority(s workflowState) error {
	if !deliveryRun(s.Request) || s.Result.Delivery == nil || s.Request.Envelope != deliveryEnv("herdr-run-request") {
		return deliveryContinuityError("continuation_unsupported", "continue supports accepted delivery-owner request schema 2; this contract needs a separately supported transition")
	}
	if _, err := workflowhandoff.ValidateDeliveryTaskRunDraft(s.Request.HandoffDraft); err != nil {
		return deliveryContinuityError("continuation_unsupported", "required mandate semantics are unsupported: "+err.Error())
	}
	if err := workflowhandoff.ValidateDeliveryMandate(s.Request.HandoffDraft, s.Request.Delivery.Agreement); err != nil {
		return deliveryContinuityError("continuation_unsupported", "required delivery agreement is unsupported: "+err.Error())
	}
	if s.Acceptance == nil || s.StartSHA256 == nil || s.Result.Delivery.PermissionState != "recipient_confirmed_contract" {
		return deliveryContinuityError("continuation_authority", "continue requires the active owner's existing positive native acceptance; resolve missing authority with the original owner")
	}
	raw, err := workflowBound(*s.Acceptance, 1<<20)
	if err != nil {
		return err
	}
	var a DeliveryAcceptance
	if err = decode(raw, 1<<20, &a); err != nil {
		return err
	}
	if a.Envelope != deliveryEnv("run-acceptance") || a.Acceptance.Acceptance != "started" || a.RuntimeClaim.NativeSessionID == nil || *a.RuntimeClaim.NativeSessionID != s.Result.Transport.AgentSessionID || !equal(a.RuntimeClaim.EffectivePolicySHA256, s.Result.Delivery.ActualPolicySHA256) || !equal(a.DeliveryPermission.ActualPolicyEvidence, s.Result.Delivery.ActualPolicyEvidence) {
		return deliveryContinuityError("continuation_authority", "actual acceptance, policy or native session differs; requested launch facts grant no transition")
	}
	if err = workflowClaim(s, a.RunID, a.RequestSHA256, a.SessionID); err != nil {
		return err
	}
	if err = positiveDeliveryClaim(a.Acceptance, a.DeliveryPermission, s.Request); err != nil {
		return deliveryContinuityError("continuation_authority", err.Error())
	}
	return nil
}

func deliveryContinuationPreview(d Dependencies, s workflowState, in DeliveryContinuationInput) (DeliveryContinuationPreview, error) {
	p := DeliveryContinuationPreview{Kind: "PlyDeliveryContinuationPreview@1", SchemaVersion: 1, RunID: in.RunID, State: "blocked", ReaderCapability: deliveryContinuityCapability, Context: FileBinding{s.Result.Paths.Context, s.ContextSHA256}}
	if err := deliveryContinuationAuthority(s); err != nil {
		return p, err
	}
	if in.ContextPath == "" || in.ContextPath != s.Result.Paths.Context {
		return p, deliveryContinuityError("continuation_context", "continue requires the exact original private context; it never rewrites frozen context bytes")
	}
	runtime, err := workflowEffectiveRuntime(s)
	if err != nil {
		return p, err
	}
	p.HistoricalLauncher = deliveryLauncherObservation(runtime.Executable)
	if in.RuntimeEvidencePath != "" {
		binding, _, e := deliveryReadBinding(in.RuntimeEvidencePath, 1<<20)
		if e != nil {
			return p, deliveryContinuityError("runtime_observation", "read current owner runtime evidence: "+e.Error())
		}
		if e = deliveryRuntimeObservation(s, binding); e != nil {
			return p, e
		}
		p.RuntimeObservation = &binding
	} else if s.Continuation != nil {
		prior, e := workflowContinuationRecord(s)
		if e != nil {
			return p, e
		}
		if prior.ControlExecutable.SHA256 == in.ControlExecutable.SHA256 {
			p.RuntimeObservation = prior.RuntimeObservation
		}
	}
	if p.HistoricalLauncher.State != "present" && p.RuntimeObservation == nil {
		return p, deliveryContinuityError("runtime_observation_required", "historical provider launcher is "+p.HistoricalLauncher.State+" at "+runtime.Executable.Path+"; continuation requires --runtime-evidence from the original owner observing its current session and unchanged authority; do not restore or substitute the old launcher")
	}
	if err := workflowFreshRuntimeOperation(d, s, false, runtime, p.RuntimeObservation != nil); err != nil {
		return p, err
	}
	if err := workflowhandoff.ValidateDeliveryOwnerContinuity(d.Workflow, s.Result.Handoff.Locator); err != nil {
		return p, fmt.Errorf("continue required artifacts for reader %s: %w", deliveryContinuityCapability, err)
	}
	cwd, err := d.CWD()
	if err != nil {
		return p, err
	}
	if cwd != s.Observed.Target.WorktreeLocator {
		return p, deliveryContinuityError("continuation_session", "continue must run in the same exact Task worktree and live owner session")
	}
	if err = physical(cwd, false); err != nil {
		return p, err
	}
	target, err := d.Workspace.IntegrationGit.ObserveIntegrationWorktree(s.Observed.Target.WorktreeLocator, s.Observed.Target.Ref)
	if err != nil {
		return p, err
	}
	if target.Locator != s.Observed.Target.WorktreeLocator || target.Ref != s.Observed.Target.Ref || target.GitCommonDir != s.Observed.Target.GitCommonDir || len(target.InProgress) != 0 {
		return p, deliveryContinuityError("continuation_stale", "Task worktree identity or Git operation changed; inspect the original worktree before continuing")
	}
	if err = deliveryLiveOwner(d, s); err != nil {
		return p, deliveryContinuityError("continuation_session", "observe the original provider/session without restarting: "+err.Error())
	}
	actual, err := d.Executable()
	if err == nil {
		actual, err = filepath.EvalSymlinks(actual)
	}
	if err != nil || actual != in.ControlExecutable.Path {
		return p, deliveryContinuityError("continuation_control", "continuation source must be the actual invoking executable")
	}
	if err = verifyExecutable(in.ControlExecutable); err != nil {
		return p, err
	}
	if a := s.Result.Delivery.Attempt; a != nil && a.State == "attempted" {
		if a.Kind != "verification" {
			return p, deliveryContinuityError("effect_unknown", "observe unresolved "+a.Kind+" attempt "+a.ID+" through its native records before continuing; no effect is replayed")
		}
		var initial deliveryVerificationReceipt
		if err = readValue(filepath.Join(a.Path, "attempt.json"), 256<<10, &initial); err != nil {
			return p, deliveryContinuityError("effect_unknown", "verification reservation is unreadable; observe the same attempt without replaying it")
		}
		if _, err = deliveryRecoverVerification(s, initial.Review.Locator); err != nil {
			return p, deliveryContinuityError("effect_unknown", err.Error())
		}
	}
	if s.Continuation != nil {
		r, err := workflowContinuationRecord(s)
		if err != nil {
			return p, err
		}
		if r.ControlExecutable.SHA256 == in.ControlExecutable.SHA256 && (r.RuntimeObservation != nil || equal(r.RuntimeObservation, p.RuntimeObservation)) {
			p.State, p.Generation, p.Directory, p.ControlExecutable, p.Proof, p.Confirmation = "existing", r.Generation, filepath.Dir(s.Continuation.Locator), r.ControlExecutable, s.Continuation, r.Confirmation
			p.RuntimeObservation = r.RuntimeObservation
			p.NextAction = "Continue callbacks with " + ShellQuote(p.ControlExecutable.Path) + " and the unchanged --context " + ShellQuote(p.Context.Locator) + "."
			return p, nil
		}
	}
	p.Generation = deliveryContinuationGeneration(s) + 1
	if p.Generation > 128 {
		return p, deliveryContinuityError("continuation_unsupported", "continuation history limit reached; inspect preserved history before a separately supported transition")
	}
	p.Directory = deliveryContinuationDirectory(s, p.Generation)
	p.ControlExecutable = Executable{filepath.Join(p.Directory, "ply-control"), in.ControlExecutable.SHA256}
	// A durable proof is an already-validated transition whose final state
	// pointer may be missing. Observe it before checking freshness for a new
	// effect; recovery must retain intervening reports and its original bytes.
	if recovered, err := deliveryPendingContinuation(s, in, &p); err != nil || recovered {
		return p, err
	}
	if err = deliveryContinuationPreparation(s, p); err != nil {
		return p, err
	}
	p.RuntimeRefreshAllowed = true
	if p.RuntimeObservation != nil {
		if err = deliveryRuntimeObservationFresh(d, *p.RuntimeObservation); err != nil {
			return p, err
		}
	}
	p.Confirmation = deliveryContinuationConfirmation(s, in.ControlExecutable, p.ControlExecutable, deliveryContinuityCapability, p.RuntimeObservation)
	p.State, p.NextAction = "ready", "Run the same continue command without --check to preserve a compatible control for this same accepted session."
	return p, nil
}

// This read-only observation distinguishes a known partial file publication
// from an unknown effect. Fresh evidence may finish the same generation; it
// never overwrites its original before-state or already-copied control.
func deliveryContinuationPreparation(s workflowState, p DeliveryContinuationPreview) error {
	if err := physical(p.Directory, true); err != nil {
		return err
	}
	entries, err := os.ReadDir(p.Directory)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		switch entry.Name() {
		case "ply-control":
			if err = verifyExecutable(p.ControlExecutable); err != nil {
				return deliveryContinuityError("continuation_integrity", "partial continuation control changed: "+err.Error())
			}
		case "before-state.json":
			raw, e := readFile(filepath.Join(p.Directory, entry.Name()), 8<<20, true)
			if e != nil {
				return e
			}
			if _, e = deliveryContinuationBeforeState(s, raw, p.Generation); e != nil {
				return e
			}
		default:
			return deliveryContinuityError("effect_unknown", "unrecognized partial continuation artifact; observe the preserved transition before refreshing runtime evidence")
		}
	}
	return nil
}

func deliveryContinuationBeforeState(s workflowState, raw []byte, generation int) (workflowState, error) {
	var before workflowState
	if err := decode(raw, 8<<20, &before); err != nil {
		return before, err
	}
	if deliveryContinuationGeneration(before) != generation-1 || !equal(before.Request, s.Request) || !equal(before.Observed, s.Observed) || before.Result.Handoff != s.Result.Handoff || before.Result.RunID != s.Result.RunID || before.Result.SessionID != s.Result.SessionID || before.Result.Transport.AgentSessionID != s.Result.Transport.AgentSessionID || before.Result.Transport.WorkspaceID != s.Result.Transport.WorkspaceID || before.Result.Transport.TabID != s.Result.Transport.TabID || before.Result.Transport.PaneID != s.Result.Transport.PaneID || before.Result.Transport.TerminalID != s.Result.Transport.TerminalID || before.Result.Paths.Context != s.Result.Paths.Context || before.ContextSHA256 != s.ContextSHA256 || !equal(before.Acceptance, s.Acceptance) || !equal(before.StartSHA256, s.StartSHA256) {
		return before, deliveryContinuityError("continuation_integrity", "pending before-state no longer binds the same accepted run, owner, context and history")
	}
	if err := deliveryContinuationAuthority(before); err != nil {
		return before, err
	}
	prior, err := workflowEffectiveRuntime(before)
	if err != nil {
		return before, err
	}
	current, err := workflowEffectiveRuntime(s)
	if err != nil {
		return before, err
	}
	if prior.PlyExecutable != current.PlyExecutable || !equal(before.Result.Delivery.ActualPolicyEvidence, s.Result.Delivery.ActualPolicyEvidence) {
		return before, deliveryContinuityError("continuation_integrity", "partial continuation changed its previous control or accepted authority")
	}
	return before, nil
}

func deliveryPendingContinuation(s workflowState, in DeliveryContinuationInput, p *DeliveryContinuationPreview) (bool, error) {
	binding, raw, err := deliveryReadBinding(filepath.Join(p.Directory, "continuation.json"), 1<<20)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, deliveryContinuityError("continuation_integrity", "observe pending continuation proof: "+err.Error())
	}
	var record deliveryContinuationRecord
	if err = decode(raw, 1<<20, &record); err != nil {
		return false, err
	}
	if record.SourceControl != in.ControlExecutable || !equal(record.RuntimeObservation, p.RuntimeObservation) {
		return false, deliveryContinuityError("effect_unknown", "a continuation proof is already preserved; resume its exact source control and runtime observation before requesting a different transition")
	}
	pending := s
	pending.Continuation = &binding
	if _, err = workflowContinuationRecord(pending); err != nil {
		return false, err
	}
	p.State, p.Proof = "ready", &binding
	p.Confirmation = digest([]any{s, binding, "resume-published-continuation"})
	if digest(s) == record.BeforeState.SHA256 {
		p.Confirmation = record.Confirmation
	}
	p.NextAction = "Resume the preserved continuation proof in this same session; retain any intervening report and do not publish another control."
	return true, nil
}

func WorkflowPreviewDeliveryContinuation(d Dependencies, root string, in DeliveryContinuationInput) (DeliveryContinuationPreview, error) {
	if err := containing(d, root); err != nil {
		return DeliveryContinuationPreview{}, err
	}
	s, err := workflowRead(root, in.RunID)
	if err != nil {
		return DeliveryContinuationPreview{}, err
	}
	return deliveryContinuationPreview(d, s, in)
}

func WorkflowContinueDelivery(d Dependencies, root string, in DeliveryContinuationInput, confirmation string) (DeliveryContinuationPreview, error) {
	var p DeliveryContinuationPreview
	if err := containing(d, root); err != nil {
		return p, err
	}
	err := workflowUpdate(d, root, in.RunID, func(s *workflowState) error {
		var err error
		p, err = deliveryContinuationPreview(d, *s, in)
		if err != nil {
			return err
		}
		if p.State == "existing" {
			return nil
		}
		if confirmation == "" || confirmation != p.Confirmation {
			return deliveryContinuityError("continuation_stale", "run observations changed before continuation; inspect and repeat the native command")
		}
		if err = verifyExecutable(p.ControlExecutable); err != nil {
			return err
		}
		if p.Proof != nil {
			s.Continuation = p.Proof
			if _, err = workflowContinuationRecord(*s); err != nil {
				return err
			}
			p.State = "continued"
			p.NextAction = "Continue callbacks with " + ShellQuote(p.ControlExecutable.Path) + " and the unchanged --context " + ShellQuote(p.Context.Locator) + "."
			return nil
		}
		beforePath := filepath.Join(p.Directory, "before-state.json")
		before, err := readFile(beforePath, 8<<20, true)
		if os.IsNotExist(err) {
			before, err = readFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, true)
		}
		if err != nil {
			return err
		}
		// A report may have been recorded after this immutable snapshot but
		// before proof publication. Keep both; the proof validator below checks
		// the exact original identity, authority and previous control chain.
		snapshot, err := deliveryContinuationBeforeState(*s, before, p.Generation)
		if err != nil {
			return err
		}
		previous, err := workflowEffectiveRuntime(*s)
		if err != nil {
			return err
		}
		r := deliveryContinuationRecord{Kind: "PlyDeliveryContinuation@1", SchemaVersion: 1, Generation: p.Generation, RunID: in.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, NativeSessionID: s.Result.Transport.AgentSessionID, BeforeState: FileBinding{filepath.Join(p.Directory, "before-state.json"), hash(before)}, SourceControl: in.ControlExecutable, PreviousControl: previous.PlyExecutable, ControlExecutable: p.ControlExecutable, Context: p.Context, Acceptance: *s.Acceptance, NativeStartSHA256: *s.StartSHA256, ActualPolicySHA256: *s.Result.Delivery.ActualPolicySHA256, ReaderCapability: deliveryContinuityCapability, Confirmation: p.Confirmation, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
		r.RuntimeObservation = p.RuntimeObservation
		r.Confirmation = deliveryContinuationConfirmation(snapshot, in.ControlExecutable, p.ControlExecutable, deliveryContinuityCapability, p.RuntimeObservation)
		if s.Request.Delivery.Agreement != nil {
			r.DeliveryAgreementSHA256 = workspace.DeliveryAgreementDigest(*s.Request.Delivery.Agreement)
		}
		if err = d.writeOnce(r.BeforeState.Locator, before); err != nil {
			return err
		}
		// A lost publication reply resumes the identical immutable proof, never
		// writes a new timestamp into an occupied slot.
		proofPath := filepath.Join(p.Directory, "continuation.json")
		if existing, readErr := readFile(proofPath, 1<<20, true); readErr == nil {
			var old deliveryContinuationRecord
			if err = decode(existing, 1<<20, &old); err != nil {
				return err
			}
			r.RecordedAtUTC = old.RecordedAtUTC
		} else if !os.IsNotExist(readErr) {
			return readErr
		}
		proof, err := workflowKeep(d, proofPath, r)
		if err != nil {
			return err
		}
		if err = d.fault("delivery_after_continuation_proof"); err != nil {
			return err
		}
		s.Continuation = &proof
		if _, err = workflowContinuationRecord(*s); err != nil {
			return err
		}
		p.State, p.Proof = "continued", &proof
		p.NextAction = "Continue callbacks with " + ShellQuote(p.ControlExecutable.Path) + " and the unchanged --context " + ShellQuote(p.Context.Locator) + "; reuse a completed verification only with verify --reuse ATTEMPT_ID."
		return nil
	})
	return p, err
}
