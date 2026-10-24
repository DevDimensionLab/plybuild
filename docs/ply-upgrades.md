# PV-01: Ply upgrades without binary locking (TC0E6E1FE)

Ordinary installed `ply` handles prepared startup and existing Run callbacks.
Historical Ply path/hash fields remain serialized observations; replacing or
removing those files does not grant or revoke authority. No version, build,
source-commit or replacement compatibility gate is introduced. An executable
must exist when the OS actually launches it.

## Consumer inventory

| Consumer | Change | Retained checks |
| --- | --- | --- |
| TaskRun preview/start and terminal accept/report | Keep original Ply observations in confirmation and context, without reading old bytes or comparing callback executable paths | Provider executable, policy evidence, preparation, context, Task cwd, request/session and immutable acceptance |
| Herdr report and Delivery startup/callbacks | Remove Ply byte/path checks; current CLI reads old requests/contexts unchanged | Mandate, actual authority, live owner/native session, Herdr, stops, retirement and frozen inputs |
| Stored Delivery continuation | Validate immutable transition history without requiring either historical controller file | Before-state, original request/context, accepted policy, session and record digests |
| Review, verify, delivery and closeout | Shared runtime freshness no longer requires old Ply; generated callbacks use ordinary `ply` | Independent reviewer identity, exact candidate/review/test bytes, unknown-effect reservations and delivery target |
| Worker supervision, attempts and product-answer continuation | Child processes use current `Dependencies.Executable`; runtime controls no longer hash historical Ply | Provider/policy, parent/worker ownership locks, request/generation, queue, candidate, unknown attempts and explicit opt-in |
| Week scheduler, controls and detached recovery | Historical Ply no longer blocks progress; native adapter uses current invocation; new service command resolves installed Ply through saved PATH | Contract/workset/deadline, parent/child control, environment, queue and effect ownership |
| Existing week service records | Read the original serialized command as history; current CLI can run under the existing locks without rewriting an old launchd entry | All non-runtime plist settings remain compared; no provider is killed or uncertain input replayed |
| Parallel registry consumer inventory | Classify actual remaining effects without comparing archived CLI bytes with reader hash | Run/Task evidence, release/integration state and unresolved effects |
| Build/release/install and tested product candidate | Unchanged: these hashes identify actual tested/installed product bytes | Exact candidate and provenance checks remain mandatory |

Recovery/resume protocols still validate their immutable records and current
provider/session observations. Their historical metadata remains readable. A Ply
upgrade alone does not require invoking those protocols. Old binaries already
running cannot gain new code retroactively; ordinary upgraded CLI is the entry
for existing Runs. No global installation or provider configuration changes are
part of this Task.

## Exact legacy representation limit

Before PV-01, bound-start preview@1 returned a request containing only kind,
schema_version, request_id, workspace, input and a SHA-256 binding. The binding
hashed the entire preview basis, including Ply path/hash. Preview saved nothing;
the request did not carry that basis. After any basis change, a digest alone
cannot identify which field changed or recover the previous value. Removing the
Ply component from an unknown preimage is impossible. Accepting arbitrary
mismatches would also accept changed mandates, policies or delivery agreements.

`TestPlyUpgradeLegacyOpaquePreviewLimit` constructs exactly that old digest,
replaces Ply, and proves apply rejects without creating an intent. This is a
pre-existing representation limit, not a new compatibility gate. There is no
claim of automatic support for an unretained legacy preview basis. Existing
retained intents accept their original digest and preserve request/Run identity
without rewriting their bytes (`TestPlyUpgradeLegacyRetainedIntentKeepsIdentity`).
New previews exclude only Ply identity from the binding and survive upgrade
between preview/apply, with one start on replay. All other inputs stay bound.

## Verification and evidence

The new consumer fixtures failed before production changes: prepared startup
rejected changed/missing Ply, owner report rejected changed/missing controllers,
stored continuation required deleted files, bound preview/apply rejected the
upgrade, and worker/week controls rejected their unchanged requests. RED logs
were preserved privately as `/private/tmp/TC0E6E1FE-red*.log`.

The committed entrypoint `test/pv01-upgrade-acceptance.sh` executes the new
upgrade fixtures and existing mandate/policy/session, unknown-verification and
exact-candidate guards. It checks the exact candidate's embedded source identity;
that is tested-product evidence, not a requirement on the CLI operating a Run.
The native wrapper supplies existing read-only Go caches, with isolated HOME
and short physical TMPDIR and no tool network. Separate checks cover relevant
CLI tests, `go vet` and `make build`.

Tests use real native command consumers with isolated Git/Herdr/provider fixtures,
including a built candidate CLI for review-return. They prove protocol behavior,
not a live provider start, human product judgment, installed smoke or delivery.
Broad taskrun/taskexecute diagnostic results remain separate; their known base
suite timeouts do not establish a regression pass or fail for PV-01. Root owns
independent native review and the later installed Ply smoke.

## Review round 1 corrections

Independent review of candidate `0378f024` requested a concrete legacy service
startup fix and remaining old-lock test corrections. The reviewer's isolated
service-start reproducer was adopted as `TestPlyUpgradeLegacyWeekServiceStart`:
it failed while synthetic launchctl attempted the removed historical executable.
The corrected start selects the current controller through the existing owner
locks, keeping the original plist and binding unchanged. The fix was followed
by its focused GREEN test and the reviewer's exact targeted command.

Worker preview/consume, stored continuation and parallel activation tests now
expect Ply-only changes to succeed while still testing provider, policy,
unsupported contract, unknown effects, remaining effects and other owner
isolation. Parallel owners retain separate requests, contexts, native sessions,
panes and Task sources while sharing ordinary installed Ply. CLI help points
to that same supported command. No broad suite was repeated for these fixes.
The old hash-only preview acceptance question is a separate coordinator/human
decision; these technical fixes do not claim to resolve it.
