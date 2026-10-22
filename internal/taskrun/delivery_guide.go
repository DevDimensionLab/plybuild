package taskrun

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/devdimensionlab/plybuild/internal/workflowhandoff"
	"github.com/devdimensionlab/plybuild/internal/workspace"
)

// These files describe input, not evidence. Runtime facts and human judgments
// stay unobserved until the recipient fills a copy from actual observations.
func workflowDeliveryGuide(d Dependencies, s workflowState) error {
	if !deliveryRun(s.Request) {
		return nil
	}
	dir := filepath.Join(workflowDeliveryGuideDirectory(s), "templates")
	a := DeliveryAcceptance{
		Acceptance:         Acceptance{Envelope: deliveryEnv("run-acceptance"), RunID: s.Result.RunID, RequestSHA256: s.Result.RequestSHA256, SessionID: s.Result.SessionID, RuntimeClaim: RuntimeClaim{NativeSessionID: ptr(s.Result.Transport.AgentSessionID)}, Sandbox: json.RawMessage("null"), Issues: json.RawMessage("[]")},
		DeliveryPermission: DeliveryPermissionAcceptance{LaunchContractSHA256: s.Request.Runtime.PermissionBinding.EffectivePolicySHA256, ActualPolicyEvidence: []Evidence{}},
	}
	if s.Request.Delivery.Agreement != nil {
		a.DeliveryPermission.DeliveryAgreementSHA256 = workspace.DeliveryAgreementDigest(*s.Request.Delivery.Agreement)
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
		if mandateReceipt(s.Request) && f.name == "acceptance.json" {
			continue
		}
		if _, err := workflowKeep(d, filepath.Join(dir, f.name), f.value); err != nil {
			return err
		}
	}
	runtime, err := workflowEffectiveRuntime(s)
	if err != nil {
		return err
	}
	control, id, context := ShellQuote(runtime.PlyExecutable.Path), s.Result.RunID, ShellQuote(s.Result.Paths.Context)
	guide := fmt.Sprintf(`# Delivery callback guide

This guide and the adjacent templates describe the frozen callback contract.
They are not observed acceptance, review, verification or human QA. Copy a
template to a new private file; never overwrite the templates or a submitted
input. A delivery report, its evidence logs, the acceptance script and a review
file are ordinary inputs: Ply binds the bytes it reads by SHA-256, so an
ordinary 0644 file is accepted and a file writable by group or other users is
refused. Evidence logs and the review are re-read from their path against that
digest, so keep them unchanged. File modes do not prove who wrote a file;
submit only files you chose. Runtime acceptance, human QA and other
authenticating or control files keep the owner-only (0600) requirement.
Run callbacks from the bound Task worktree: %s.

Control executable: %s
Run: %s
Context: %s
Provider and native session observed by Herdr: %s / %s

If the selected goal delegates this pilot to a parent/root that owns
notifications, this internal agent sends no Slack messages. Return through the
bound callbacks and the same tab; the parent owns user questions and completion.
Otherwise honor the existing notification contract in the delivery instructions.

## Begin with actual observations before Task writes

Use the bound control's workflow execute begin command with this Run
and --context. Supply --provider, --native-session,
--policy (your effective runtime observation and its source), and repeat
--read-root / --write-root for actual physical roots. Include
--permission-confirmed only when the observed policy covers the selected
mandate, exact effects and targets. Optional --model and --profile are
actual observations; omit them when unknown. Ply creates private evidence and
temporary directories, computes hashes and submits the native acceptance. It
does not change provider permissions. Repeating identical observations reuses
the receipt. No handwritten mandate, acceptance JSON or hashing shell is needed.
Unknown or insufficient permission stops dependent work. Read the frozen goal
and mandate for scope; the launch request is not proof of runtime permission.

## Advanced acceptance file contract

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
Progress is optional; verifier_results may be empty. It is never a prerequisite
for starting separate review and does not qualify a candidate or a human verdict.
IDs are local stable declaration keys, not test names or a registered catalogue:
use 1..128 lowercase ASCII letters, digits, '.', '_', '/' or '-', starting with
a letter or digit; no '//', trailing '/' or '..' path segment, and no duplicates.
A reason is required: 1..2000 Unicode codepoints, no surrounding whitespace or
control characters. Explain the actual observation; correct fields without
inventing a pass or exit. Keep actual test names in argv and bound raw evidence.
For example, only after actually running the command from the stated cwd and
observing exit 0, replace the paths and digest with the real captured log binding:

    "verifier_results": [{
      "id": "focused-test", "outcome": "passed",
      "argv": ["go", "test", "./internal/taskrun", "-run", "^TestDeliveryReportVerifierFieldDiagnostics$", "-count=1"],
      "cwd": "/absolute/task", "exit": 0,
      "evidence": [{"locator": "/absolute/private/test.log", "sha256": "sha256:<actual-log-digest>"}],
      "reason": "The named focused test executed and exited 0; raw output is bound."
    }]

For an unexecuted check use not_run (or unknown when its result is unknown),
exit=null and a truthful reason; never copy the example's outcome as evidence.
Future coordinator review, installation and smoke checks remain future duties
until observed separately.

%s workflow execute report %s --context %s --file <copied-report.json>

## Verify a clean committed candidate

Implement the declared acceptance entrypoint at %s with meaningful executable
checks. Preserve an actual local review using templates/review.json. Bind
candidate_oid and candidate_tree to the clean committed candidate. Set the
reviewer_claim (1..256 Unicode codepoints, with no surrounding whitespace),
actual reviewer_session_id (or null) and observed decision. A
passed decision with no open actionable findings is required for qualification;
an empty template is not a review. Findings/fixes use the existing native review
schema: id, severity (low, medium, high or critical), summary and evidence_ids
(sorted unique managed artifact IDs). Available IDs are acceptance-script,
candidate-review, task-requirements, verification-receipt, verifier-stderr and
verifier-stdout. For an observation preserved in this review, cite candidate-review;
an arbitrary note name or file path is not a managed artifact ID. Keep all actual
findings and fixes with truthful summaries and supported evidence references;
never delete them to make validation pass. Entries are sorted by id, with IDs
unique across all three arrays. An empty array is appropriate only when the
actual review found none.

Native verification assembles and preserves task-requirements
from the frozen Spec requirements and this exact candidate's acceptance and review
evidence. Do not hand-author a second requirements/status artifact to satisfy
that managed ID. Managed artifacts become available during successful candidate verification;
their IDs refer to the native evidence set, not arbitrary files you create.
Keep actual evidence and truthful requirement outcomes;
a progress report or this guide does not qualify a candidate.

The callback itself executes /bin/sh against the declared acceptance entrypoint,
captures output, checks candidate stability and preserves a separate candidate.

Independent review is a round on this same Run: the coordinator runs the
installed "ply workflow execute review RUN_ID --coordinator CLAIM" for your
clean committed candidate; Ply starts a separate reviewer session in its own
Herdr tab, which returns with "workflow execute review-return". Keep that
candidate unchanged until the round returns. Findings
return to you: correct them, commit and ask the coordinator for a new round.
When a round passed exactly the current candidate, omit --review and Ply builds
the review record from it; earlier findings stay as fixes.

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
`, s.Observed.Target.WorktreeLocator, runtime.PlyExecutable.Path, id, s.Result.Paths.Context, s.Request.Runtime.Provider, s.Result.Transport.AgentSessionID, control, id, context, control, id, context, s.Request.Delivery.AcceptancePath, control, id, context, control, id, context, control, id, context)
	if a := s.Request.Delivery.Agreement; a != nil {
		guide += fmt.Sprintf("\n## Frozen delivery agreement\n\nMode: %s. Source: %s. Target: %s. Repository: %s. Local target: %s.\nAgreement SHA256: %s.\n\nFor positive acceptance, delivery_permission.delivery_agreement_sha256 must name\nthis exact agreement. Set delivery_permission.allowed_effects to %s only after\nthe actual runtime policy covers those exact effects and targets. This is an\nactual permission claim; templates do not establish it. No mode or target\nsubstitution is allowed. An altered agreement requires a revised Spec and a new\nauthorized execution.\n", a.Mode, a.SourceRef, a.TargetRef, a.GitHubRepository, a.TargetWorktree, workspace.DeliveryAgreementDigest(*a), workflowJSON(workflowhandoff.DeliveryAllowedEffects(*a)))
		if a.Mode == workspace.DeliveryPullRequest {
			guide = strings.ReplaceAll(guide, "## Actual human QA and local integration", "## Actual human QA and PR delivery")
			guide = strings.ReplaceAll(guide, control+" workflow execute integrate "+id+" --context "+context, "Use the delivery CLI to register, check and execute the exact pull_request agreement after actual candidate-bound human pass.")
			guide = strings.ReplaceAll(guide, "Integration requires actual pass on the exact candidate and unchanged parent.", "PR delivery requires actual pass on the exact candidate and the agreed remote repository, source and base. It stops before merge and leaves the local parent and Epic base unchanged.")
		} else if a.HumanOwnedIntegration() {
			start, end := strings.Index(guide, "## Actual human QA and local integration"), strings.Index(guide, "\n## Frozen delivery agreement")
			if start >= 0 && end > start {
				guide = guide[:start] + fmt.Sprintf(`## Candidate handoff to human integration

Register the exact technically qualified Delivery with the delivery CLI. Stop
writing the source, then explicitly release this candidate:

%s workflow execute release %s --context %s --reason 'Exact candidate qualified; source writing is released to human integration.'

The human starts ply integration from an existing return checkout, reviews the
visible plan, and supplies the actual pass/fail/blocked answer. The native command
preserves the answer and authorizes only that exact plan. The developer stops
before integration; execute integrate and delivery execute cannot bypass this
boundary. A negative human decision returns correction ownership: correct,
qualify and release the candidate again before a new human integration decision.
`, control, id, context) + guide[end:]
			}
		}
	}
	if a := s.Request.Delivery.Agreement; a != nil && a.AutomaticAcceptance() {
		guide = strings.ReplaceAll(guide, "--review <copied-review.json>", "--review <copied-review.json> --candidate-binary <absolute-private-executable>")
		guide += fmt.Sprintf("\n## Automatic acceptance instructions\n\nThis agreement explicitly selects automatic acceptance, owned by %s, with a %d-second timeout. Build the private candidate from the clean committed Task with `make build PLY_BUILD_OUTPUT=/absolute/private/ply` and pass that absolute path to verify. This uses existing buildinfo ldflags to embed the build-time Git commit and dirty status, including linked worktrees where Go may omit VCS metadata even with -buildvcs=true. Check `ply version --format json` and `go version -m /absolute/private/ply` before verification. Plain go build can omit VCS identity; without embedded metadata source_commit stays empty (unknown), and dirty=false alone does not establish a clean source. Native qualification OID/tree/hash remain separate evidence. Missing embedded identity adds no install or human gate. Never rebuild or alter the qualified bytes after verification. The script receives PLY_CANDIDATE_BINARY, PLY_CANDIDATE_OID and PLY_CANDIDATE_TREE, PATH and LC_ALL, and isolated HOME and TMPDIR. No provider credentials are inherited. The invocation freezes script and binary hashes, environment, candidate and executor identity before startup. A completed receipt is recovered without repeating the test; an unknown or active attempt remains blocked. Correct real failures in this Task, then qualify the new candidate.\n", a.Acceptance.ResponsibleActor, a.Acceptance.TimeoutSeconds)
		if !a.Acceptance.RequireHumanQA && !a.HumanOwnedIntegration() {
			start, end := strings.Index(guide, "## Actual human QA and local integration"), strings.Index(guide, "\n## Frozen delivery agreement")
			if start >= 0 && end > start {
				guide = guide[:start] + fmt.Sprintf("## Automatic local Epic delivery\n\nA successful exact candidate test and review satisfy the frozen automatic acceptance gate. Complete the authorized native local return with:\n\n%s workflow execute integrate %s --context %s\n\nObserve Git, Task queue and Epic base closure. Continue only the same reserved integration after interruption. Machine results never create human QA. Actual later human fail or blocked judgments still veto delivery.\n\nSource cleanup is the default. While you are still in the Task worktree, integrate ends with exit 5 and owner_handoff after the observed integration: the delivery is complete and your writing is finished. This is not a stop: do not retry, report a stop or edit the source. End this session and return to the coordinator, whose workflow execute follow removes the worktree. Any other exit 5 names the remaining cleanup cause.\n", control, id, context) + guide[end:]
			}
		}
	}
	if mandateReceipt(s.Request) {
		start := strings.Index(guide, "## Begin with actual observations")
		end := strings.Index(guide, "## Progress and necessary human questions")
		receipt := fmt.Sprintf(`## Receive the frozen mandate

Startup protocol: mandate_receipt_v1. Run this small idempotent receipt from
this exact Task and immutable control:

%s workflow execute check-in %s --context %s --provider %s --native-session %s

Ply checks the native mandate digest, provider/session, cwd and frozen inputs.
The receipt acknowledges the existing human mandate. It does not attest
sandbox coverage or future runtime approval. No handwritten acceptance file,
blanket permission confirmation or repeated human approval is required.
Optional --model, --profile, --policy, --read-root and --write-root preserve
sourced recipient observations. Omitted facts stay unknown; reported roots
have matches_contract=null and cannot create permission coverage.
Readback separates this receipt, observed facts, human scope, runtime decisions,
actual work and delivery. The launch profile is a request, not observed policy.

Runtime evaluates actions when they happen. Preserve any actual refusal with
its action, stage, source and known effect status through execute observe-stop;
the coordinator can report denied receipt/report transport without impersonating
the owner. After the cause is clarified, continue with execute resolve-stop,
naming the changed action you take now; no permission confirmation is needed and
the refused action stays refused. A known active stop or unknown prior attempt
blocks dependent work.
Repeating receipt cannot clear it or replay a previous action. Legacy frozen
runs retain begin and their original permission requirements.

`, control, id, context, s.Request.Runtime.Provider, ShellQuote(s.Result.Transport.AgentSessionID))
		guide = guide[:start] + receipt + guide[end:]
		guide = strings.ReplaceAll(guide, "working requires positive acceptance.", "working requires the bound mandate receipt.")
		start = strings.Index(guide, "For positive acceptance, delivery_permission.delivery_agreement_sha256")
		if start >= 0 {
			end = strings.Index(guide[start:], "authorized execution.") + start + len("authorized execution.")
			guide = guide[:start] + "The preserved human mandate authorizes only this exact agreement and effects. Native callbacks recheck the mandate, owner, candidate, target and active stops before effects. Receipt supplies no additional permissions; runtime evaluates each action." + guide[end:]
		}
	}

	return d.writeOnce(filepath.Join(workflowDeliveryGuideDirectory(s), "callback-guide.md"), []byte(guide))
}
