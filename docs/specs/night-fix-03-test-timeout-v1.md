# NF03A test-timeout evidence boundary v1 — T6C88FCE3

> Historisk dokument etter fjerningen av den kompilerte prosesskjeden.
> Execute-, worker/week-, delivery- og integrasjonsorkestreringen som omtales
> her er ikke en kjørbar oppskrift i dagens kilde. Fjernede kildefiler er
> bevart som historiske stier. Se [gjeldende direkte operasjoner](../../README.md#human-owned-process).

Decision: **(c) no justified small implementation**. The historical coverage-run
timeout has no established cause in the inspected evidence. Do not start a timeout
fix developer from this Spec. Root w9:pE decides whether to close this backlog
track as investigated/unfixed and continue the selected batch.

This static authoring Task changes only this document. Source basis is clean Epic
commit `52ce0496f839129ca251fc3016ef1edb3de7cdec`, tree
`af84f139d4385a6732ad1e7ef1e22ef9d01b032d`. Historical master below identifies an
analysis input; the sole delivery target remains `refs/heads/ply_agentic_workflow_support`.
No Go tests, coverage, runtime/OS probes or denial replays were performed here.

## Historical evidence and missing facts

The [timeout backlog](/Users/perottochristensen/github/ply/planning/koe/ply-crap-gyldig-coverage-og-timeout.md),
[analysis summary](/Users/perottochristensen/github/ply/planning/reports/crap-backlog-20261022/backlog-plan.md)
and [provenance](/Users/perottochristensen/github/ply/planning/reports/crap-backlog-20261022/provenance.json)
agree on the following boundary:

| Fact | Supported conclusion |
| --- | --- |
| Handoff reports ten-minute timeout in `internal/taskrun`. | A reported timeout, without an identified test or demonstrated failure mechanism. |
| Profile `/tmp/ply-coverage.out`, recorded SHA256 `06d2143f5f8b177ff0ca9a96ce61c8716f3de98ccff0692bfcddbe8fe778272a`, 2,513,596 bytes, `mode: set`. | Identifies the supplied profile in the old analysis; it does not establish successful test completion. The profile was not re-read here. |
| Analysis source `c401914e04e286e4f8d78d6fde5d54203db67518`, tree `35ada8d07eedcced91242e72cb99882c16d42fb9`; analysis context darwin/arm64, cgo=true. | Identifies analysis context only. Provenance explicitly leaves the profile's source revision and original test build context unknown. |
| Original test argv/cwd, environment, stdout/stderr, exit, failing test name and timeout stack are not supplied by these records. | Cannot distinguish slow completed subtests, aggregate workload, blocked subprocess/OS observation, deadlock or an outer runner timeout. No original raw-log locator is provided for this coverage run. |

The provenance's `commands` are **ply-crap analysis commands**, whose exit 0 is
not the exit of the coverage-producing Go test. `complete=true` means no mapping
gaps were discovered. Zero counters do not prove missing tests. The interrupted,
unbound profile remains unsuitable for baseline coverage, CRAP prioritization or
before/after improvement claims; NF01/NF02 focused passes do not repair it.

Other timeout history must retain its own identity:

- The [process-journal backlog](/Users/perottochristensen/github/ply/planning/koe/ply-task-prosessjournal-og-observer.md)
  records race timeout under concurrent heavy testing on 3 October. Its
  [result control](/Users/perottochristensen/github/ply/planning/design/evidence/73-task-process-journal-result-control-v1.json)
  records a default-10m to 18m adjustment and event `taskrun-race-timeout-adjusted`
  with exit 0, followed by `race-reviewed-all`. The first event spans
  15:49:33–15:59:39 UTC. These are preserved control-record claims for Journal
  candidate `503a58f9826dd71a384c7a5df449f97bc246c4d1`, not a current timing test
  or proof of the later coverage timeout's cause. Exact referenced streams
  `/private/tmp/ply-task-process-journal-delivery-v1/streams/taskrun-race-timeout-adjusted.json`,
  `race-reviewed-all.json` and `full-tests-reviewed.json` were absent at their
  named paths during this inspection; no raw argv/stack was recovered from them.
- The [agentloop learning](/Users/perottochristensen/github/ply/planning/koe/ply-agentloop-kontrollskatt-og-pflyt-friksjonsreduksjon.md)
  records FS02's whole-package `-timeout 3600s`, CL02's interrupted 50-minute
  selection, AL04's two 600s selections narrowed to 15 named tests, and AL10's
  broad `TestWorkflow` selection with a shorter outer Bash limit. These support
  bounded selection and truthful exits as process lessons. They do not bind any
  of those runs to this profile, prove current package runtime, or prove deadlock.

This is the evidence boundary of the named supplied records and direct locators,
not a claim that no recoverable log exists anywhere. No private root archive,
foreign Run or provider transcript was searched.

## Current source audit and hypotheses

| Existing boundary on the selected base | Implication for a later investigation |
| --- | --- |
| [Makefile](../../Makefile) `test` runs shell prerequisites, then `go test -v -cover ./...`; it declares no explicit Go timeout or test shards. `acceptance` invokes install/status/upgrade/build scripts. | Both are broader than this diagnosis. Neither proves the original argv; install acceptance is outside this docs role. |
| [Scoped acceptance](../../scripts/accept-scoped-continuation.sh) already selects a package and anchored test names with `-count=1`. | Existing selection is sufficient machinery for a later bounded check. Its DEP02 prefix is not a validated timeout-fix selection; do not copy it blindly. No new scheduler is justified. |
| deliveryCLIBinary (`internal/taskrun/delivery_cli_test.go`, historisk kilde) uses absolute `PLY_DELIVERY_TEST_BINARY`, otherwise builds a fresh binary in `t.TempDir()` on each helper call. | Repeated builds are a plausible cost when the variable is absent. Original variable value, call count and cost are unknown; this is no measured bottleneck. |
| automaticCLI (`internal/taskrun/delivery_automatic_cli_test.go`, historisk kilde) and cleanupCLIAt (`internal/taskrun/delivery_default_cleanup_test.go`, historisk kilde) synchronously run CLI children without a command context deadline. [observePathWriterLiveness](../../internal/taskrun/delivery_writer_liveness.go) falls back to real `lsof`/`ps` when `LocalProcesses` is nil. | Parent dependency closures do not cross CLI process boundaries. A later selection needs a pre-start launch/observer audit; these paths alone establish neither an actual hang nor a new bug. |
| NF02's [receiptCLIQualified/receiptControlledRuntime/receiptCLICall](../../internal/taskrun/delivery_receipt_test.go) install adapters and a pre-exec closeout tripwire before startup, use a 45s child context plus 1s WaitDelay, and call modeled closeout in process. | Delivered receipt isolation is already present. Defensive SourceWriters and StartupProcessTree have asserted zero calls in this ordinary non-adopted path; only reached branches count as evidence. Do not commission NF02 again. |
| NF01's nf01ReserveOnly (`internal/taskrun/registry_nf01_test.go`, historisk kilde) interrupts at the Files seam after reservation and before StartWorker. | Preserve the delivered reservation boundary; old detached-supervisor uncertainty is not a new timeout diagnosis or permission to replay it. |
| [TestSystemLocalProcessesObservesRealWriters](../../internal/taskrun/delivery_writer_liveness_test.go) starts real shell children and invokes the system observer; automatic interruption tests (`internal/taskrun/delivery_automatic_execution_unix_test.go`, historisk kilde) signal their own process. | Full-package execution includes actual OS/signal behavior. An apparently focused prefix is not an isolation guarantee. No such tests ran here. |

Aggregate suite size, repeated CLI builds and subprocess waits remain separate
current hypotheses. None identifies the historical failing test. No source-proven
small fixture correction or proven-too-large exact selection was established;
paths (a) and (b) would invent certainty at this point. Raising a timeout,
skipping tests or introducing infrastructure would not resolve this evidence gap.

## Next owner and conditional later work

Root w9:pE owns the batch disposition; the **Ply test/CI maintainer** owns any
separately selected evidence investigation. First obtain from the original
coverage-run owner an exact preserved argv/cwd, source OID/tree/dirty, toolchain,
build context/env, start/end, raw stdout/stderr/exit and named test/timeout stack.
If no locator can be supplied, record those fields unknown and stop this small
fix track. Do not replace missing historical evidence with a guessed reproduction
or ask for a broad suite just to fill it.

If concrete evidence later identifies a current bounded failure, write a new
implementation Spec for that exact test/helper and unchanged behavior contract.
Before execution, audit dependencies through startup, adoption/recovery, writer
inspection and every subprocess. Install controlled dependencies before start;
each owned child needs an explicit bounded handle and wait. Refuse an unmodeled
OS/provider boundary before launch and preserve genuine denials/unknown effects.
Reuse NF02's callable-body tripwire and NF01's pre-launch Files seam where relevant;
a helper name, route label or non-nil observer is not proof that its body ran.

For that later Spec, require a safe RED on the actual affected branch, a focused
GREEN preserving native state and behavior, and actual call/zero-call assertions.
Use existing anchored `-run` selections and an absolute clean make candidate;
declare exact test names and justified Go/outer/child limits before starting.
Preserve unpiped process exits and raw output, assert that expected tests actually
executed, and treat unfinished/unknown outcomes truthfully. A pure selection
correction instead needs the exact old/new inventory and successful bounded
completion; it must not claim a product RED/GREEN. No numerical speedup is promised.

Full module coverage is a separate larger evidence job: freeze one source/context,
instrumented inventory and coverage mode; accept only completed successful
contributions, with exact shard union and documented merge if sharding is chosen.
Exclude the old profile and report missing packages honestly. This Spec authorizes
neither that job nor a new coverage/CRAP score.

The immediate user benefit is avoiding an unjustified code change and another
expensive unisolated test run. A later small patch's effort cannot be estimated
until a failing branch is identified; evidence recovery may end with no fix.
Prefer that bounded outcome before proposing infrastructure.

## Acceptance and handoff

This Task's acceptance is bounded docs validation: required evidence/decision
sections, readable source links, this single file as the full base-to-candidate
diff, a clean candidate and `git diff --check`. The private make-built candidate
is solely for native docs verification; its identity must match the clean commit.
These checks establish document scope and consistency, not tested product behavior.

Root starts separate native Spec review. The original owner holds verify until
root releases it, then verifies without `--review` against that exact reviewed
candidate. Hold the single local Epic integration for root's later release and
fresh acceptance/target check. No installation; source cleanup remains pending
for root, with no postdelivery follow/ps. Root owns learning and the aggregate
notification; this child sends no Slack.

Learning: a ten-minute timeout and a parseable profile are not a causal diagnosis.
Keep older timeout controls, current source hypotheses and missing raw evidence
separate before selecting a developer or claiming a coverage improvement.
