package taskrun

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
)

// These files describe input, not evidence. Runtime facts and human judgments
// stay unobserved until the recipient fills a copy from actual observations.
func workflowDeliveryGuide(d Dependencies, s workflowState) error {
	if !deliveryRun(s.Request) {
		return nil
	}
	dir := filepath.Join(s.Result.Paths.RunRoot, "delivery", "templates")
	a := DeliveryAcceptance{
		Acceptance:         Acceptance{Envelope: deliveryEnv("run-acceptance"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, RuntimeClaim: RuntimeClaim{NativeSessionID: ptr(s.Result.Transport.AgentSessionID)}, Sandbox: json.RawMessage("null"), Issues: json.RawMessage("[]")},
		DeliveryPermission: DeliveryPermissionAcceptance{LaunchContractSHA256: s.Request.Runtime.PermissionBinding.EffectivePolicySHA256, ActualPolicyEvidence: []Evidence{}},
	}
	r := DeliveryReport{Envelope: deliveryEnv("delivery-report"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, Phase: "working", Evidence: []FileBinding{}, VerifierResults: []DeliveryVerifierResult{}}
	review, err := workflowhandoff.DeliveryCandidateReviewTemplate("", "")
	if err != nil {
		return err
	}
	var draft struct {
		Binding struct {
			TaskID string `json:"task_id"`
		} `json:"task_spec_binding"`
	}
	if err := json.Unmarshal(s.Request.HandoffDraft, &draft); err != nil {
		return err
	}
	qa := workflowhandoff.DeliveryHumanAttestation{Kind: "DeliveryHumanAttestation@1", SchemaVersion: 1, TaskID: draft.Binding.TaskID}
	for _, f := range []struct {
		name  string
		value any
	}{{"acceptance.json", a}, {"report.json", r}, {"review.json", json.RawMessage(review)}, {"human-qa.json", qa}} {
		if _, err := workflowKeep(d, filepath.Join(dir, f.name), f.value); err != nil {
			return err
		}
	}
	control, id, context := ShellQuote(s.Request.Runtime.PlyExecutable.Path), s.Result.RunID, ShellQuote(s.Result.Paths.Context)
	guide := fmt.Sprintf(`# Delivery callback guide

This guide and the adjacent templates describe the frozen callback contract.
They are not observed acceptance, review, verification or human QA. Copy a
template to a new private file; never overwrite the templates or a submitted
input. Run callbacks from the bound Task worktree: %s.

Control executable: %s
Run: %s
Context: %s
Provider and native session observed by Herdr: %s / %s

If the selected goal delegates this pilot to a parent/root that owns
notifications, this internal agent sends no Slack messages. Return through the
bound callbacks and the same tab; the parent owns user questions and completion.
Otherwise honor the existing notification contract in the delivery instructions.

## Accept the actual runtime before Task writes

Copy templates/acceptance.json. The runtime_claim fields are actual observations:
runtime_id (the canonical provider codex or claude, not a CLI version), model_id
(null when unknown), profile_id (the actual profile/mode label, or null when
unavailable), effective_policy_sha256, native_session_id. Put CLI versions in
the policy evidence. The native session must equal the Herdr observation above.
Do not copy a requested model or policy into an actual claim without evidence.
delivery_permission.launch_contract_sha256 identifies the requested launch
contract. It is distinct from runtime_claim.effective_policy_sha256. Keep
the requested launch profile in the request; the actual label need not match
that selector. A differing label or null does not by itself establish or deny
authority. Actual policy evidence and contract scope must still be confirmed.
Preserve actual policy evidence as physical files, hash their exact bytes, and set
delivery_permission.actual_policy_evidence to objects with locator, sha256 and
role (effective_policy, permission_proof or runtime_contract), sorted by locator.
At least one evidence hash must equal the actual effective_policy_sha256.
Set permission_confirmed only after the actual policy covers the mandate.
The evidence is the recipient's preserved observation, not OS attestation.
An existing config/profile name or writable-directory check alone is not proof
of effective authority. Preserve the actual runtime policy/context available to
this session without secrets. Never alter permission configuration to make the
acceptance pass. While necessary authority is unresolved, report needs_input or
stopped without claiming acceptance; use a negative acceptance only to preserve
an actual negative decision. Unknown model alone is model_id=null. Selected
provider and session stay exact.

sandbox is an object with read_roots and write_roots (sorted absolute physical
paths), temp_root (absolute physical path) and matches_contract (boolean).
Report actual roots, including the workspace's .ply state/reply areas needed by
the callbacks and later candidate records when those roots are authorized.
Set acceptance to started only when actually confirmed, with issues=[].
Otherwise use rejected, conflict or unknown and issues=[{"type":"sandbox",
"detail":"actual reason"}] as appropriate. Issue types: binding, workspace,
project, target, input, contract, sandbox, principal, unknown. Sort issues by
type then detail. A negative receipt cannot grant target-write authority.

%s workflow run accept %s --context %s --file <copied-acceptance.json>

## Progress and necessary human questions

Copy templates/report.json. Set a unique event_id, summary, meaning and phase
(working, needs_input or stopped). needs_input requires question with the actual
necessary question; other phases use question=null. Read current delivery state
with workflow run show and copy delivery.last_event_sha256 into
previous_event_sha256 (null before the first delivery event). This predecessor
changes after verification, QA and integration as well as reports.
needs_input and stopped may be submitted before runtime acceptance so that an
actual permission blocker can be reported truthfully. They grant no Task-write,
verification, QA or integration authority. working requires positive acceptance.
Evidence entries use locator and sha256. Each verifier_results entry has id,
outcome (passed, failed, not_run or unknown), argv, cwd, exit, evidence and reason.
not_run and unknown require exit=null; never fabricate an unrun exit code.
Optional process_observation and preventive_followup use string or null.
Reporting progress does not qualify a candidate or a human verdict.

%s workflow execute report %s --context %s --file <copied-report.json>

## Verify a clean committed candidate

Implement the declared acceptance entrypoint at %s with meaningful executable
checks. Preserve an actual local review using templates/review.json. Bind
candidate_oid and candidate_tree to the clean committed candidate. Set the
reviewer_claim, actual reviewer_session_id (or null) and observed decision. A
passed decision with no open actionable findings is required for qualification;
an empty template is not a review. Findings/fixes use the existing native review
schema: id, severity (low, medium, high or critical), summary and evidence_ids
(sorted unique strings). Entries are sorted by id, with IDs unique across all
three arrays. An empty array is appropriate only when the actual review found none.
The callback itself executes /bin/sh against the declared acceptance entrypoint,
captures output, checks candidate stability and preserves a separate candidate.

%s workflow execute verify %s --context %s --review <copied-review.json>

## Actual human QA and local integration

Prepare the installed product journey and ask for the human's real pass, fail
or blocked answer. Keep the interactive session available. Copy
templates/human-qa.json, fill task_result_id, result_oid and result_tree from the
qualified current candidate, preserve actor_claim, start_surface and actual UTC
started_at_utc/completed_at_utc. outcome and answer must both equal the exact
pass/fail/blocked answer; observation preserves its context. Never substitute
technical success, a TTY, or your own answer for actual human judgment. This is
local provenance, not cryptographic human authentication. Candidate failure
keeps this same delivery owner: correct and qualify a new candidate.

%s workflow execute qa %s <pass|fail|blocked> --context %s --evidence <copied-human-qa.json>
%s workflow execute integrate %s --context %s

Integration requires actual pass on the exact candidate and unchanged parent.
Follow any unknown effect through the preserved same attempt. Never replay an
unknown start, verifier or integration and never replace the control executable.
`, s.Observed.Target.WorktreeLocator, s.Request.Runtime.PlyExecutable.Path, id, s.Result.Paths.Context, s.Request.Runtime.Provider, s.Result.Transport.AgentSessionID, control, id, context, control, id, context, s.Request.Delivery.AcceptancePath, control, id, context, control, id, context, control, id, context)
	return d.writeOnce(filepath.Join(s.Result.Paths.RunRoot, "delivery", "callback-guide.md"), []byte(guide))
}
