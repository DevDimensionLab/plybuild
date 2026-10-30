# NF05A early required Git-capability diagnosis v1 — TB48F9D18

> Historisk dokument etter fjerningen av den kompilerte prosesskjeden.
> Execute-, worker/week-, delivery- og integrasjonsorkestreringen som omtales
> her er ikke en kjørbar oppskrift i dagens kilde. Fjernede kildefiler er
> bevart som historiske stier. Se [gjeldende direkte operasjoner](../../README.md#human-owned-process).

Decision: **a small current diagnostic gap remains**. Commission the next
separate developer to validate the formats already observed by the parent
preview path and retain its diagnostic before confirmation. No new capability
platform or higher Git-version floor is needed.

Static basis: commit `4816203311b27bdcf8681fc2f63cbbb1930e0ff7`, tree
`b5bcd638d98fcc8214d3bab58bcdf4d8f8dca5cc`. Only this document changes.
No Go or CLI tests, actual Apple Git/PATH/permission probes, historical QA
replay, installation or downgrade ran in this specification Task.

## Historical evidence

The frozen backlog `koe/ply-task-queue-nattko.md`, addition 1 October 2026,
session 102, and `design/evidence/67-ws06-human-qa-v1-result-control.json`
record candidate `51354bb03dd63688a6643536111b0abac88bd441`. The QA script's
restricted PATH selected Apple Git 2.39.5. The retained command returned
`--show-ref-format\n` with exit 0; the preparation readback contained that
literal ref format, `freshness=fresh` and a confirmation digest. Apply exited 1
with `workspace_work_store_conflict`: the plan did not bind an exact clean Epic
observation. The control preserves unchanged registry/fixture refs/worktrees
and no created Task worktree. Human QA was incomplete, with no verdict.

That is historical safe rejection plus misleading earlier diagnosis, not a
current reproduced failure or permission to replay the old environment. QA v2's
choice of capable Git fixes the harness selection, not the product preview.
The terminal dash/encoding observation is separate and outside this Spec.

## Current source reconciliation

| Source boundary | Current static evidence |
| --- | --- |
| [work_item_git.go](../../internal/workspace/work_item_git.go), `systemTaskIntegrationGit.ObserveIntegrationWorktree` | Reads `rev-parse --show-object-format=storage` and `rev-parse --show-ref-format` through `read`/`integrationLine`. It checks command error and line shape, not allowed format values. It does not call `ObserveIntegrationRepository` or `supportedGitVersion`. |
| [project_git.go](../../internal/workspace/project_git.go), `singleGitLine`; `work_item_git.go`, `integrationLine` | An option echo with one terminating newline is accepted. Even a lone newline yields an accepted empty string. Missing terminator/multiple lines and command errors are rejected. Exit-zero alone therefore proves neither capability nor supported storage. |
| [epic_base.go](../../internal/workspace/epic_base.go), `queueParentObservation` | Combines `WorkGit.ObserveWorktree` with `IntegrationGit.ObserveIntegrationWorktree`; checks clean/symbolic/no operation and matching OID/tree, then returns `planObservation(full)` without checking formats. Base update planning also consumes this function. |
| [task_queue.go](../../internal/workspace/task_queue.go), `queueMixedEntryEvaluation` / `queueErrorReasons` | Calls `queueParentObservation` before returning a ready row and preparation plan. Errors become blocked/unknown reasons, but a successful option echo currently reaches ready. Base-free goal rows use `queueGoalEvaluation`, which is not an execution confirmation. |
| [task_preparation.go](../../internal/workspace/task_preparation.go), `buildPreparePlan` / `PrepareTask` | Builds the above plan, initializes freshness to fresh, and issues its digest. A named blocked row retains reasons; `--next` can discard them for generic `task_queue_no_ready`. Apply rebuilds a new plan before reservation. Existing active preparation recovery is a separate historical-request path. |
| task_goal_execution.go (`internal/workspace/task_goal_execution.go`, historisk kilde), `buildGoalExecutionPreview` / `buildGoalExecutionPreviewInTarget` / `PreviewTaskGoalExecution` / `PrepareTaskGoalExecution` | Fresh execution calls `queueParentObservation` before `goalExecutionDigest`. Apply rebuilds/revalidates before materializing execution content. The shared parent check can prevent both new confirmation paths. |
| [task_queue_store.go](../../internal/workspace/task_queue_store.go), `validateQueuePlanObservation` | Persistence already permits object format sha1 or sha256 and ref storage files only. Echo/empty/unknown/reftable are rejected; do not relax this safeguard. |
| [task_integration.go](../../internal/workspace/task_integration.go), `buildIntegrationPlan`; `work_item_git.go`, `ObserveIntegrationRepository` | Repository observation already requires Git 2.45+ using `supportedGitVersion`. Integration subsequently blocks unsupported object/ref formats, including reftable, and enforces OID lengths. That version gate is bypassed by parent preview, not absent globally. Leave that existing integration policy and structured reasons unchanged. |

This proves a remaining source path for incorrect new readiness/confirmation;
it does not establish a fresh runtime incident. Adding another version check
would duplicate policy and still fail to validate the actual output.

## Smallest implementation

Change the shared parent admission boundary in `epic_base.go`:
`queueParentObservation`, after its existing identity/cleanliness/drift checks
and before returning `planObservation(full)`, must accept only `ObjectFormat`
sha1/sha256 and `RefFormat` files. Use the existing named values/constants where
available; do not change either Git interface or add another command. Return
`queueError("task_queue_git_capability_unsupported", detail)` on other values.
The detail must identify the required command option, quote the observed value
safely, and state the allowed output: for example, Git observation for
`rev-parse --show-ref-format` returned `"--show-ref-format"`; this Task path
requires `"files"`. For valid but unsupported reftable, say storage is unsupported
by this Task path; do not claim Git itself is broken. Empty strings and unknown
values get the same explicit rejection. Check object before ref for deterministic
first failure. Preserve the original error and operation detail when the
observer itself fails; do not relabel repository/path/ref failures as capability
failures. Existing malformed-output and command-error diagnostics remain early.

In `task_preparation.go`, `buildPreparePlan`, retain parent capability rejection
for named and `--next` selection: no plan/confirmation, state blocked and
freshness unknown, with the specific capability reason. For an unavailable or
malformed observation keep the existing observation-unknown code/detail and
unknown state/freshness. When `--next` finds no eligible plan, retain the first
such diagnostic in queue order instead of replacing it with generic no-ready;
keep normal ranking, other reasons and empty-queue behavior. Do not overwrite
an existing preparation's historical request/confirmation during recovery.
`cmd/workspace_task_prepare.go`'s `outputPreparation` already renders reasons
and suggests apply only when a confirmation exists; no new output schema is
needed. Blocked preparation readbacks may retain the existing exit-zero read
convention: failure is explicit in state/reasons and absence of confirmation.
Goal execution returns its existing error path, with empty confirmation, before
materialization. No claim of successful readiness is allowed in either path.

Keep `work_item_git.go`, repository integration classification, the version
floor, registry schema and persistence validators unchanged. sha1 and sha256
with files remain supported, with all existing OID-length/CAS/ref/path/content
provenance checks intact. Reftable remains deliberately unsupported by Task
plans; do not add storage support or substitute files for an unknown response.
Preserve no-repository, dirty, lock, detached/ref mismatch, drift and legitimate
unknown effects. Read-only preview remains no-write: no intent/reservation,
metadata publication, ref/worktree mutation, installation or agent launch.

## Controlled RED/GREEN

The later developer writes and runs failing regression tests before production
changes. Use the existing `systemTaskIntegrationGit.run` seam in
`work_item_git.go`, reached via `Dependencies.IntegrationGit`; do not replace
`ObserveIntegrationWorktree` with a fake result or test only a new helper.
Existing patterns are in [task_source_preflight_test.go](../../internal/workspace/task_source_preflight_test.go).
Build fresh disposable selected-solution/goal fixtures using `queueSelect`,
`queueSetFixture`, `goalFixture`, and `goalQueueFixture` from
[task_queue_test.go](../../internal/workspace/task_queue_test.go) and
[task_goal_execution_test.go](../../internal/workspace/task_goal_execution_test.go).
Install interception only after fixture setup and before the consuming preview.

Proposed focused tests in a new workspace test file:

- `TestParentGitCapabilityPreparationPreview`: call real `PrepareTask` for
  named and next selectors. Intercept the exact parent `rev-parse` format
  arguments and return exit 0/Err nil/option echo, lone newline, empty bytes,
  multiline, or unknown format. Assert the intercept was reached, no usable
  plan/confirmation, the precise reason/state/freshness above and unchanged
  registry bytes/refs/worktree inventory. Include reftable as valid unsupported
  storage, not a parser error. Current echo/lone-newline cases must fail RED
  specifically because a preview confirmation is issued.
- `TestParentGitCapabilityGoalPreview`: call real `PreviewTaskGoalExecution`
  with the same intercepted parent formats. Assert the command was reached,
  error/empty confirmation and unchanged content/registry/ref inventory. RED
  must reach the preview defect, not fail during fixture construction.
- `TestParentGitCapabilityChangedBeforeApply`: obtain a supported fresh
  preparation/goal preview, switch only the seam output to echo, then call
  `PrepareTask` with Apply or `PrepareTaskGoalExecution` using that digest.
  Assert no new preparation/reservation/content/ref/worktree effects. Do not
  reuse historical Runs or active recovery fixtures.
- Positive subcases feed sha1/files and sha256/files with matching fixture
  OIDs/trees through the same real observation path, and assert valid new
  confirmation, unchanged preview state, and no-write command tripwires.
  Use fully controlled command responses for sha256 if the fixture uses sha1;
  never mislabel 40-character OIDs as sha256. Command failures carry realistic
  Exit and Err; preserve the existing no-repo/ref/dirty/drift error semantics.

Keep positive reads controlled: delegate only known reads to a pinned capable
Git in the disposable fixture, or supply all command responses in process.
Fail on unexpected commands and every write/launch; record argument/count
tripwires so a bypassed check cannot pass. No actual unsupported Git, PATH
mutation, downgrade, provider/OS/process probe, or old denial replay. Run only
the three named tests with `-count=1 -timeout=90s` and a 120-second outer bound;
then the existing directly affected preparation/goal/base/integration checks
selected by the developer and reviewer. No full-suite or version matrix is
required. Keep real RED, GREEN and fixture failures distinct.

## Acceptance and handoff

For the later implementation, a new reviewer must inspect the actual callable
path, negative/positive readbacks and no-effect evidence. A separately authorized
installed CLI smoke then uses a new isolated disposable workspace and a fixture
Git executable that models only the two format responses while delegating
allowlisted reads to a pinned capable Git. Exercise text and JSON preparation
preview, goal execute-preview, and a stale pre-change confirmation apply refusal.
Capture exact argv/output/exit and before/after registry/content/refs/inventory;
require the same diagnostics, no offered apply command, no confirmation and no
writes. Run each command with a 30-second bound and the smoke with 120 seconds;
never select real Apple Git or a live provider. Installation belongs to that
later authorized product Task, not this docs-only delivery.

This Task's bounded docs acceptance checks the decision, source pointers,
behavior/test/handoff requirements, single-document base diff and
`git diff --check`. A private make candidate is for native docs qualification
only and proves no Git-capability runtime behavior.

Root w9:pE starts separate native Spec review and commissions the later developer.
Hold verify until root releases it after PASS; the original owner then verifies
this exact committed candidate without `--review`. Hold the single authorized
local Epic integration for a separate root release and fresh target/acceptance
check. Preserve source cleanup pending; no postdelivery follow/process inspection.
Root owns learning and the aggregate notice; this child sends no Slack.

Learning: persistence can fail closed while the earlier plan still promises
readiness. Validate the values at the actual shared confirmation boundary;
repository integration's version check does not cover that caller.
