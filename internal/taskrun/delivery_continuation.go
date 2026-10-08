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

const deliveryContinuityCapability = "delivery-owner-v2/task-spec-v3/task-content-v2/scoped-publications-v1"

// DeliveryContinuationReaderCapability names this reader's bounded transition
// policy. It is compatibility evidence, never runtime permission evidence.
const DeliveryContinuationReaderCapability = deliveryContinuityCapability

type DeliveryContinuationInput struct {
	RunID             string
	ContextPath       string
	ControlExecutable Executable
}

type DeliveryContinuationPreview struct {
	Kind              string       `json:"kind"`
	SchemaVersion     int          `json:"schema_version"`
	RunID             string       `json:"run_id"`
	State             string       `json:"state"`
	Generation        int          `json:"generation"`
	Directory         string       `json:"directory"`
	ControlExecutable Executable   `json:"control_executable"`
	Context           FileBinding  `json:"context"`
	Proof             *FileBinding `json:"proof,omitempty"`
	Confirmation      string       `json:"confirmation"`
	ReaderCapability  string       `json:"reader_capability"`
	NextAction        string       `json:"next_action"`
}

type deliveryContinuationRecord struct {
	Kind                    string      `json:"kind"`
	SchemaVersion           int         `json:"schema_version"`
	Generation              int         `json:"generation"`
	RunID                   string      `json:"run_id"`
	RequestSHA256           string      `json:"request_sha256"`
	SessionID               string      `json:"session_id"`
	NativeSessionID         string      `json:"native_session_id"`
	BeforeState             FileBinding `json:"before_state"`
	SourceControl           Executable  `json:"source_control"`
	PreviousControl         Executable  `json:"previous_control"`
	ControlExecutable       Executable  `json:"control_executable"`
	Context                 FileBinding `json:"context"`
	Acceptance              FileBinding `json:"acceptance"`
	NativeStartSHA256       string      `json:"native_start_sha256"`
	ActualPolicySHA256      string      `json:"actual_policy_sha256"`
	DeliveryAgreementSHA256 string      `json:"delivery_agreement_sha256,omitempty"`
	ReaderCapability        string      `json:"reader_capability"`
	Confirmation            string      `json:"confirmation"`
	RecordedAtUTC           string      `json:"recorded_at_utc"`
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
	if r.Kind != "PlyDeliveryContinuation@1" || r.SchemaVersion != 1 || r.Generation != n || r.RunID != s.Result.RunID || r.RequestSHA256 != s.Result.RequestSHA256 || r.SessionID != s.Result.SessionID || r.NativeSessionID != s.Result.Transport.AgentSessionID || r.ReaderCapability != deliveryContinuityCapability || r.BeforeState.Locator != filepath.Join(dir, "before-state.json") || r.ControlExecutable.Path != filepath.Join(dir, "ply-control") || r.ControlExecutable.SHA256 != r.SourceControl.SHA256 || r.Context != (FileBinding{s.Result.Paths.Context, s.ContextSHA256}) || s.Acceptance == nil || r.Acceptance != *s.Acceptance || s.StartSHA256 == nil || r.NativeStartSHA256 != *s.StartSHA256 || s.Result.Delivery == nil || s.Result.Delivery.ActualPolicySHA256 == nil || r.ActualPolicySHA256 != *s.Result.Delivery.ActualPolicySHA256 {
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
	if err = deliveryContinuationAuthority(old); err != nil {
		return nil, err
	}
	if !equal(old.Result.Delivery.ActualPolicyEvidence, s.Result.Delivery.ActualPolicyEvidence) {
		return nil, deliveryContinuityError("continuation_authority", "actual policy evidence differs from the accepted transition history")
	}
	prior, err := workflowEffectiveRuntime(old)
	if err != nil {
		return nil, err
	}
	if prior.PlyExecutable != r.PreviousControl {
		return nil, deliveryContinuityError("continuation_integrity", "prior control binding differs from the immutable transition history")
	}
	if r.Confirmation != digest([]any{old, r.SourceControl, r.ControlExecutable, deliveryContinuityCapability}) {
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
	if !deliveryRun(s.Request) || s.Result.Delivery == nil || s.Request.Envelope != deliveryEnv("herdr-run-request") {
		return deliveryContinuityError("continuation_unsupported", "continue supports accepted delivery-owner request schema 2; this contract needs a separately supported transition")
	}
	if _, err := workflowhandoff.ValidateDeliveryTaskRunDraft(s.Request.HandoffDraft); err != nil {
		return deliveryContinuityError("continuation_unsupported", "required mandate semantics are unsupported: "+err.Error())
	}
	if err := workflowhandoff.ValidateDeliveryMandate(s.Request.HandoffDraft, s.Request.Delivery.Agreement); err != nil {
		return deliveryContinuityError("continuation_unsupported", "required delivery agreement is unsupported: "+err.Error())
	}
	if s.Acceptance == nil || s.StartSHA256 == nil || s.Result.Delivery.PermissionState != "recipient_confirmed_contract" || s.Result.Delivery.OwnershipRelease != nil || s.Result.Delivery.Phase == "completed" {
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
	if err := workflowFresh(d, s, false); err != nil {
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
	if _, err = workflowAgentGet(d, s, false); err != nil {
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
		if r.ControlExecutable.SHA256 == in.ControlExecutable.SHA256 {
			p.State, p.Generation, p.Directory, p.ControlExecutable, p.Proof, p.Confirmation = "existing", r.Generation, filepath.Dir(s.Continuation.Locator), r.ControlExecutable, s.Continuation, r.Confirmation
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
	p.Confirmation = digest([]any{s, in.ControlExecutable, p.ControlExecutable, deliveryContinuityCapability})
	p.State, p.NextAction = "ready", "Run the same continue command without --check to preserve a compatible control for this same accepted session."
	return p, nil
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
		before, err := readFile(filepath.Join(s.Result.Paths.RunRoot, "state.json"), 8<<20, true)
		if err != nil {
			return err
		}
		previous, err := workflowEffectiveRuntime(*s)
		if err != nil {
			return err
		}
		r := deliveryContinuationRecord{Kind: "PlyDeliveryContinuation@1", SchemaVersion: 1, Generation: p.Generation, RunID: in.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, NativeSessionID: s.Result.Transport.AgentSessionID, BeforeState: FileBinding{filepath.Join(p.Directory, "before-state.json"), hash(before)}, SourceControl: in.ControlExecutable, PreviousControl: previous.PlyExecutable, ControlExecutable: p.ControlExecutable, Context: p.Context, Acceptance: *s.Acceptance, NativeStartSHA256: *s.StartSHA256, ActualPolicySHA256: *s.Result.Delivery.ActualPolicySHA256, ReaderCapability: deliveryContinuityCapability, Confirmation: p.Confirmation, RecordedAtUTC: d.Now().UTC().Format(time.RFC3339Nano)}
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
