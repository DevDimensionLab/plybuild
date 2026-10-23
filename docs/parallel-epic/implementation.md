# Concurrent native Tasks and exact serial Epic return

T32BD249B implements the interactive native path in the frozen
[specification](spec.md). Reader support for format 8 does not activate it.
Ordinary serial work retains its actual persisted format, including format 4.
Installation and workspace activation are separate coordinator operations.
The original design and scenario requirements remain historical inputs.

## Operator path

After separately qualified activation, start each explicitly named Goal from
its registered Epic worktree:

```sh
ply workflow execute --spec GOAL_A --parallel
ply workflow execute --spec GOAL_B --parallel
ply workspace task queue list --schema-version 2 --format json
ply workflow status --schema-version 2 --all --format json
```

Each Task retains its own preparation, source ref/worktree, Run, frozen control,
provider session, context, review rounds and verification. Binding rejects a
session owned by an active member or its reviewer. Reviewer control contexts
belong to one exact Run/round/candidate. Named closure removes only that member.
Schema 1 refuses multi-active readback; schema 2 always emits `active`, including
`[]`. `current` is only a singleton projection. Pending entries distinguish
serial/parallel readiness; neither proves live provider liveness.

The exact return path serializes target effects by physical Git common directory
and full ref, including aliases. A return moves the Epic base. A sibling's old
proof cannot qualify return against that moved parent. From the sibling's own
worktree, use its frozen control and context:

```sh
CONTROL workflow execute reconcile-source RUN --context CONTEXT --check --format json
CONTROL workflow execute reconcile-source RUN --context CONTEXT --apply --confirm SHA256 --format json
```

The preview binds source OID/tree, new parent, original preparation,
Run/request/mandate, review history and preceding effective generation. Apply
reserves a native source operation before controlled Git. A conflict retains
that source and active membership. Stage a real resolution, then use
`--continue --operation ID`. `--abort --operation ID` refuses unexpected edits
against the recorded conflict fingerprint. Unknown effects require observation
of the same operation and cannot be replayed by a new attempt.
Controlled merge preflight refuses collisions with ignored local files before
creating a marker or attempting Git merge. The same owner may preserve or move
the local content and recheck the existing intent against its original source
and parent. Actual unknown merge effects remain nonreplayable.
A known remaining gap is cancellation of an unstarted intent after its parent
changes: current abort requires a pending merge, while fresh reconciliation
refuses the retained intent. Preserve that operation for coordinator review;
this negative does not prove complete stale-intent recovery.

Reconciliation appends an effective basis. Original Task parent, Spec, mandate,
controls and prior receipts remain unchanged. Even identical source OIDs need
fresh separate review and executed verification for a changed generation.
Discarding committed reconciliation ancestry is refused. A Task's unresolved
human fail/blocked verdict survives source changes and automatic qualification.
Only a later real pass on the exact result resolves it; old QA is never rebound.

## Activation and old consumers

Only `workspace registry activate-parallel` publishes format 8. Its default
check is read-only. Apply requires an exact confirmation over predecessor bytes,
target, clean identified format-8 reader from a completed native self-hosted CLI
Run with complete closeout, and finite native consumer inventory. It preserves
a byte-exact backup and one direct 4/5/6/7-to-8 boundary without intermediate
upgrades.

```sh
ply workspace registry activate-parallel --qualified-run RUN --coordinator CLAIM --check --format json
ply workspace registry activate-parallel --qualified-run RUN --coordinator CLAIM --apply --confirm SHA256 --format json
```

An incompatible resumable owner (PAR-A01) blocks activation. Fully terminal
history (PAR-A02) has no remaining effect reader. Delivered, automatically
released, cleanup-pending history (PAR-A03) may use a qualified compatible
coordinator to finish only its original closeout. Missing compatibility blocks
that path. Human integration/correction release (PAR-A04) retains all legitimate
remaining paths, including frozen correction control. Unknown inventory blocks
activation. Activation never migrates foreign controls, clears stops, installs
anything or proves writer absence.

Ordinary content/queue upgrades retain format-3/4 rules. Explicit release-backed
`upgrade --to-format 5`, `6` (Epic creation) and `7` (Task adoption) retain their
original permits/boundaries. Every production `CurrentRegistryFormat` use is a
reader/validation upper bound, never an ordinary upgrade target. Worker, week,
backup, adoption, legacy prepare, `--next` and restart remain exclusive and
cannot acquire parallel authority through recovery or a nil singleton.

## Consumer and recovery inventory

| Consumer | Native contract |
| --- | --- |
| Queue codec, validation and CAS | One append-only membership fold; every active ID excluded from pending. |
| Goal materialization/preparation | Explicit plan 2, exact Goal/ref, named recovery and publication retry identity. |
| First-Goal registration / queue set | Version-2 nested readback carries all active members and preserves pending refs/order. |
| Execute / initial-start recovery | Request/Run 3 and Handoff 4 bind mode at creation; old controls retain their meaning. |
| Receipt/report/stop/question/retirement | Exact Run/provider/session/context; no sibling transfer or stop erasure. |
| Review/verification/TaskResult | Source plus effective generation, sealed separate reviewer and executed acceptance. |
| Integration/queue/Epic base | Repository/ref ledger and CAS; base plan 2 names sorted active IDs. |
| Return recovery | Exact authority/result, named advance and own base version prove historical completion; no second merge after a later descendant return. |
| Closeout | Defer destructive shared-Git cleanup while another preparation is active; existing absence/lifecycle gates remain. |
| Registered queue / workflow status | Pure registered projections, explicit schema 2 and legacy multi-active refusal. |
| Exclusive entry/recovery routes | Preserved request/transport/admission checks prevent importing parallel authority. |

Source operations release metadata locks during human conflict resolution.
Release borrows its existing locked native snapshot for basis checks, preserving
Project-before-WorkItem ordering without reacquiring the same lock. No metadata
lock spans a human wait or an entire Task's source work.

## Verification and limits

Focused `TestParallel*` cases map to PAR-T01–16 using isolated native registry,
Git, callback, review, executable acceptance, return and recovery. Bounded serial
ledger/alias and upgrade tests exercise shared primitives. The private acceptance
entrypoint selects named tests, never the full suite. Raw stdout/stderr/exits,
behavior reds, setup failures and actual outcomes live in the private T32BD249B
return area; test names alone do not establish execution.

All provider/Herdr/process observations in these fixtures are simulated. No live
two-agent, browser, PWA or process-inspection smoke is inferred. Isolated current
cleanup uses CAS `update-ref -d` and preserved config bytes/inode. This does not
prove live Codex process-local trust for every Git implementation. Cleanup is
deferred while another active participant shares Git metadata. This adds no
terminal-status or generic trust-service capability.

Optional `make build PLY_BUILD_VERSION=v0.0.0-T32BD249B` leaves the normal default
unchanged. It lets the same private make-built image exercise real optional
upgrades whose existing boundaries require a `v`-prefixed reader. Commit/dirty
embedding is unchanged. Retain full build argv, `version --format json`,
`go version -m`, hash and source OID/tree. Candidate CLI tests use these exact
bytes, never an installed fallback. Installation, runtime smoke and live
activation remain parent-owned decisions after delivery.
