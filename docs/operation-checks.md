# Goal startup check scope (RO02)

The pre-change cost on the fixed RO01B six-Goal fixture is 79 physical content
reads for preview and 1705 for materialization/preparation, with 18 registry
decoder calls, 32 content closures and 26 queue closures during the latter. Parsing
and cross-observation generation reuse belong to later work.

| Check | Failure protected against | Actual dependency / phase | Cost and decision before implementation |
| --- | --- | --- | --- |
| Strict registry decoding, identities, lifecycle, cross references | Unknown authoritative fields, forged relationships or impossible transitions | Entire shared registry, every observation and publication | Retain every decoder call and structural check; shared mutations must preserve unrelated records. |
| External content publication closure | Missing/changed immutable Goal, Spec, request, documents, result authority | Selected Task before startup/materialization and every publication | Scope complete closure to the selected Task; unrelated history cannot authorize or affect this operation. Full Snapshot and explicit historical reads retain validation. |
| Queue content closure | Substituted selections, preparation basis or upgrade backup | Selected Task and affected shared target | Retain selected Task preparation/selection checks and affected target backup checks; retain structural validation of all queue events. |
| Recovery search | Retrying another preserved execution or hiding an interrupted outcome | Selector, target queue, matching Task execution histories | Resolve candidates from registry metadata before reading external manifests; never skip a matching interrupted execution. |
| Live target, queue revision, current problem/selection, owner and Git checks | Concurrent relevant change or effects on an incorrect destination | Current selected operation, under native locks and again before effects | Retain; content scoping grants no mutation authority. |
| Whole-registry confirmation hash | Concurrent shared changes | Investigate actual startup confirmation separately | Retain until a tested dependency projection proves all relevant bindings; no blanket hash omission. |

Intentional behavior: independent startup may succeed with unavailable unrelated
historical external content. Explicit full snapshots and historical content reads
still detect it. No audit service, cache, migration, or new user maintenance step
is introduced. Task scope write restrictions remain distinct from multi-Task
queue mutation, which must serialize the entire fresh locked registry.

## Implemented boundary and drift binding

Goal preview reads the structurally validated registry under the native snapshot
locks, resolves its target once, recovers matching interrupted requests, and
validates the selected Task's complete external closure. Preparation starts with
that same structural observation under its existing mutation lock, binds the
resolved Task before publication, then validates its content. Publication is
rejected while the operation has no bound Task. Later registry observations and
publications still obtain fresh verified bytes; this change adds no content
cache, parser shortcut or reusable confirmed generation. Existing Task write
fences remain independent and alternate stores retain their stronger checks.

Recovery candidates include **every pending Task and the current preparation**
in the resolved target. Both `--next` and explicit Spec selection preserve the
interrupted-request search before choosing new work. A known confirmation also
recovers its native `goal-execute/<digest>/spec` publication after queue removal.
A repeated default acceptance path is resolved from preserved queue metadata.
Custom acceptance paths conservatively include the target's prepared Task
histories because those paths have no registry index. Reading unavailable
content in these recovery candidates can still block: it might hide an earlier
interruption and cannot safely be treated as unrelated history.

The observed `confirmation differs from current preview` path hashes the raw
registry as part of prepared workflow observations. Bound-start also compares a
raw registry hash before its native effects. New Goal delivery startup now binds
a domain-tagged `TaskStartupDependencies@1` projection. It omits only **other
Tasks' content policies, problem/Spec/assessment/selection registrations and
content publications**. Structural cross-reference checks prohibit those records
from supplying the selected Task's basis; actual selected manifests/documents
are still read and verified. The selected current basis and all shared records
remain bound. Physical `registry_sha256` remains the observed file SHA in readback.
Legacy workflow confirmations and already recorded bound-start intents retain
their previous hash behavior; the optional dependency digest distinguishes new
bound-start intents without rewriting historical records.

This is intentionally conservative: **all** Tasks and Epics, ownership/worktree
operations, results, human QA, integration authorities/intents/attempts/results,
result-Spec links, all queue events/preparations, and all base versions/updates
remain in the dependency digest. An unrelated Task creation, lifecycle change,
queue edit or base update can still invalidate confirmation. Project/repository
bindings and fresh Git/runtime/policy checks also remain. This bounded change
eliminates content-only interference without claiming a minimal dependency graph
for every possible shared transition. Unknown outcomes and request identity
checks still take precedence over fresh startup.

## Verification and measurement method

The new tests fail on the exact delivered A archive at consumed preview,
preparation, bound-start and workflow-confirmation boundaries. Candidate tests
cover missing/corrupt unrelated history with both selectors, selected content
loss, unknown authoritative registry fields, queue drift before effect, full
historical detection, shared dependency digest sensitivity, and real unrelated
problem publication followed by startup using the earlier confirmation. Existing
interruption/replay/planner-reordering, scope-preservation and startup race tests
remain in the focused checks. The legacy confirmation regression found during
implementation is explicitly retained as `TestWorkflowCodexTrustLegacyAndConflict`.

`TestDeliveryStartupCLI` is excluded from the final passing set: its removed
`resolve-stop --permission-confirmed` flag fails identically on delivered A.
The unchanged baseline failure is preserved separately. No full-suite pass is
claimed, and the unrelated FS01 cleanup debt is outside this change.

The benchmark uses the unchanged RO01B six-Goal fixture archive and test driver,
with baseline `e5c96d73ca04f21df21afbc47606bb3ba7f3d9b6`. It restores only
`/private/tmp/ro01b-measure/fixture`, never a live registry, and prepares synthetic
worktrees without starting an agent. Five interleaved warm samples per mode and
variant follow excluded warmups. Separate instrumented binaries count physical
managed reads and actual registry/content/queue closure calls; timing binaries
contain no counters. Memory figures are Go `TotalAlloc`/`Mallocs` deltas, with
heap-in-use retained in raw samples; they are not process RSS or peak memory.

The preliminary run reduced reads/allocations but regressed median latency:
preview 195.91 to 220.78 ms; preparation 2119.44 to 2321.78 ms. It added one
registry decoder call and closure pair (18/32/26 to 19/33/27). These samples and their
source/binary hashes are preserved in `ro02/measurement/preliminary`. Investigation
found repeated live target/Git resolution inside recovery and an extra complete
preview before the mutation lock. Both avoidable costs were removed: the planner
now consumes the already resolved target, and operation dependencies bind under
the original lock without another observation. Final timing runs after tests
settle; retained checks and remaining costs are reported with the final results.

## Final fixed-fixture results

| Mode / metric | Delivered A | B |
| --- | ---: | ---: |
| Preview median (range), ms | 199.06 (189.97–201.69) | 187.85 (185.15–195.66) |
| Preparation median (range), ms | 2083.42 (2054.67–2192.92) | 1899.57 (1881.38–1938.85) |
| Preview physical content reads | 79 | 54 |
| Preparation physical content reads | 1705 | 725 |
| Preview allocated bytes / allocations | 10,459,936 / 172,060 | 7,056,656 / 113,198 |
| Preparation allocated bytes / allocations | 305,385,320 / 5,480,520 | 179,340,328 / 3,259,962 |
| Preview decoder calls / content closures / queue closures | 1 / 1 / 1 | 1 / 1 / 1 |
| Preparation decoder calls / content closures / queue closures | 18 / 32 / 26 | 18 / 32 / 26 |

Median latency decreased 5.6% for preview and 8.8% for preparation; preview
ranges overlap, so five samples do not establish a general latency guarantee.
Allocated bytes decreased 32.5% and 41.3%; allocation counts decreased 34.2%
and 40.5%. The removed preliminary duplicate observation is visible in the
unchanged final decoder/closure counts. The 18 remaining decoder calls protect fresh
shared snapshots/publication/recovery; closures still validate the selected
Task and affected target on every such phase. These are retained costs, not
claims of completed parsing or confirmed-generation optimization.

The new dependency hash is additional work on startup paths outside the
workspace-only fixture driver. A separate measurement on the same frozen
registry, five warm batches of 100 calls with registry loading excluded, gives
0.180 ms median per digest (range 0.177–0.192 ms), approximately 144,002 allocated
bytes and 4,073 allocations per call. It performs no external content reads;
canonical serialization scales with retained shared metadata. This small added
cost buys a binding that permits unrelated content publication while keeping
selected basis and shared conflicts protected. Neither benchmark measures a
real Herdr/provider startup or promises the same cost on a larger registry.

Exact source archives, timing/counter binaries, fixture SHA, raw samples,
allocation/heap observations, preliminary regressions, baseline red tests and
final report are preserved separately under:
`/Users/perottochristensen/.local/share/ply/planning/request-optimization-20261020/ro02/measurement/`.
The RO01B evidence remains unchanged. The excluded receipt CLI cleanup journey
also reached the existing FS01 process-inspection boundary (`/bin/ps` denied by
the owner sandbox); its failure is retained. The replacement consumed CLI test
uses `PLY_DELIVERY_TEST_BINARY` for read-only `workflow run start --check`
previews around a real unrelated publication and verifies unchanged confirmation
before the native fixture startup. No cleanup or liveness assertion is inferred.

## Confirmed registry generations (RO04B)

A locked work-item session now owns a confirmed structural generation. It
reuses that generation while the physical file identity and native state remain
unchanged, and confirms its own publication after replacement and directory
sync. Each returned snapshot and retained publication is independently owned.
Current content scope and external content closure are still checked on every
observation/publication. Failed publications discard confirmation; recovery and
a new lock period observe physical bytes. No state crosses sessions.

The historical RO02 table above counts decoder calls, not actual YAML parser
invocations. RO03 introduced separate parser instrumentation; RO04B retains
separate physical-read, decoder, YAML-parser, closure and publication counts.
Workspace operation timings do not measure Herdr/provider startup, and separate
A–D percentage changes must not be added. The fresh RO04B per-stage and
cumulative measurements, exact archives, binaries, drivers and fixture binding
are preserved in `/private/tmp/ro04b-avw4iiq_`. Historical reports are unchanged.
