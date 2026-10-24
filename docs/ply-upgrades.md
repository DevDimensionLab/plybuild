# PV-01: Ply upgrades without binary locking (TC0E6E1FE)

Ordinary installed `ply` handles prepared startup and existing Run callbacks.
Historical Ply path/hash fields remain serialized observations; replacing or
removing those files does not grant or revoke authority. No version, build,
source-commit or replacement compatibility gate is introduced. An executable
must exist when the OS actually launches it.

## Consumer inventory

| Consumer | Change | Retained checks |
| --- | --- | --- |
| Named Goal preview/apply | Remove historical aggregate comparison and checksum prerequisite | Current named Goal/Spec, native authority, request identity, no duplicates and consistency during startup |
| TaskRun preview/start and terminal accept/report | Keep original Ply observations in confirmation and context, without reading old bytes or comparing callback executable paths | Provider executable, policy evidence, preparation, context, Task cwd, request/session and immutable acceptance |
| Herdr report and Delivery startup/callbacks | Remove Ply byte/path checks; current CLI reads old requests/contexts unchanged | Mandate, actual authority, live owner/native session, Herdr, stops, retirement and frozen inputs |
| Stored Delivery continuation | Validate immutable transition history without requiring either historical controller file | Before-state, original request/context, accepted policy, session and record digests |
| Review, verify, delivery and closeout | Shared runtime freshness no longer requires old Ply; generated callbacks use ordinary `ply` | Independent reviewer identity, exact candidate/review/test bytes, unknown-effect reservations and delivery target |
| Worker supervision, attempts and product-answer continuation | Child processes use current `Dependencies.Executable`; runtime controls no longer hash historical Ply | Provider/policy, parent/worker ownership locks, request/generation, queue, candidate, unknown attempts and explicit opt-in |
| Week scheduler, controls and detached recovery | Historical Ply no longer blocks progress; native adapter uses current invocation; new service command resolves installed Ply through saved PATH | Contract/workset/deadline, parent/child control, environment, queue and effect ownership |
| Existing week service records | Read the original serialized command as history; current CLI can run under the existing locks without rewriting an old launchd entry | All non-runtime plist settings remain compared; no provider is killed or uncertain input replayed |
| Capacity-stop reconciliation | Recognize the actual literal CLI invocation and exact preserved native rejection without a historical executable path prerequisite | Native source/session/cwd/completion, Run/context/stop, literal grammar, exact rejection, and unresolved-effect checks |
| Parallel registry consumer inventory | Classify actual remaining effects without comparing archived CLI bytes with reader hash | Run/Task evidence, release/integration state and unresolved effects |
| Build/release/install and tested product candidate | Unchanged: these hashes identify actual tested/installed product bytes | Exact candidate and provenance checks remain mandatory |

Recovery/resume protocols still validate their immutable records and current
provider/session observations. Their historical metadata remains readable. A Ply
upgrade alone does not require invoking those protocols. Old binaries already
running cannot gain new code retroactively; ordinary upgraded CLI is the entry
for existing Runs. No global installation or provider configuration changes are
part of this Task.

## Human clarification: current named Goal startup

The explicit human clarification on 24 October removes the aggregate preview
binding as a startup prerequisite for both old and new requests. Startup uses
the currently registered named Goal/Spec and normal current-state authority.
An old opaque request works without saved preview basis, old binaries or
installation history. The checksum is optional information. No historical
pre-preview comparison, recovery lookup or replacement stale-preview gate runs.

The original request and prior native records remain preserved. Replays keep
the first request's Run and session even if the informational checksum changes;
changing selected inputs under that same request ID is still rejected. A prior
unknown attempt is observed, never replayed. Current Goal validity, Task/cwd,
actual runtime authority, known refusals, unknown effects, no duplicate input and
exact tested-candidate checks remain in force. The existing consistency guard
captures current state at startup and rejects changes during that same startup.

`TestPlyUpgradeLegacyOpaqueRequestWithoutBasisOrHistory` reads an old request,
replaces Ply, starts once and checks byte-identical original input. Tests also
cover optional checksums, checksum-only replay, current named Goal/policy changes,
invalid current selectors and interrupted/concurrent startup. The human answer
and normative clarification were attached to a native progress report; original
Spec publication and review1 were preserved unchanged.

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
Review1 also identified the historical aggregate-preview requirement. The
explicit human clarification above settles that product choice; removal is
covered by separate RED-GREEN consumer tests.

## Review round 2 corrections

The independent review confirmed the human clarification and all round 1 fixes.
Its full native-transcript capacity reproducer failed when only the executable
word changed from a historical controller to installed Ply. The corrected
consumer keeps the actual invocation, exact rejected resolve-stop command and
native output, source/session/cwd, completion and unknown-effect evidence, while
allowing ordinary installed `ply` or its literal path. Shell substitutions,
assignment prefixes, unrelated failures and wrong Run/context/stop still fail.
The same reproducer covers this behavior in `TestPlyUpgradeCapacityRejectionUsesInstalledCLI`.
Current callback and week-service documentation consistently uses installed Ply.
