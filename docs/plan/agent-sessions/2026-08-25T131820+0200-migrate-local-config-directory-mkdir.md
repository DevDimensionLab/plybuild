# Agent Session: Migrate Local Config Directory Mkdir

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T131820+0200-migrate-local-config-directory-mkdir`
Created: `2026-08-25T13:18:20+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `9eb9230dee677c155cb0b4abbeb3db32499f08d47defef67e6f9997906e57f99`
Previous: [2026-08-25T125300+0200-migrate-local-config-update-create.md](2026-08-25T125300+0200-migrate-local-config-update-create.md)
Next: [2026-08-25T135854+0200-migrate-tips-list-read-dir.md](2026-08-25T135854+0200-migrate-tips-list-read-dir.md)
Outcome: completed at `b2d37cc`; Q1.3 improved from 23/40 to 22/40 with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`os.Mkdir(dir, 0755)` in `LocalConfigDir.CheckOrCreateConfigDir` through one new
exact filesystem adapter `Mkdir` operation. Preserve exact directory selection,
the completed filesystem stat and `os.IsNotExist` decision, one single-level
mkdir attempt, mode, error and return behavior, every caller, both completed
local-config create and write lifecycles, every completed filesystem, file,
template, config, Maven, tips, structurizr, Bitbucket, Wpost, and supervisor
move, and every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.Mkdir(dir, 0755)` operation in `LocalConfigDir.CheckOrCreateConfigDir` in
`pkg/config/local.go`; one distinct zero-value-safe internal filesystem adapter
`Mkdir(path string, mode fs.FileMode) error` operation with exact `os.Mkdir`
system semantics; its focused adapter contracts; the mandatory method addition
to every current complete `filesystem.FileSystem` test double; and focused
local-config directory-create recording contracts. It does not authorize
changing `MkdirAll`, the completed stat adapter operation or contracts,
`os.IsNotExist`, another config method or effect, another adapter operation or
family, a function-valued effect dependency, a public API, inventory, mutation
harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/config/local.go` and all config tests,
every caller and use of `OpenLocalConfig`, `LocalConfigDir`, `LocalConfigFile`,
`CheckOrCreateConfigDir`, `TouchFile`, `UpdateLocalConfig`, `Config`, `Print`,
and `Exists`, complete `pkg/context/context.go`, `cmd/profile.go`, their relevant
command and context tests, relevant filesystem/file/template/config/Maven/tips/
structurizr behavior, the complete filesystem and process adapters, every
complete `filesystem.FileSystem` implementation and recording double,
`.quality/inventory`, and the continuity and quality-lift designs. Regenerate
ignored reports outside the measured tree or remove them before a clean audit.

Implementation commit `6928a72` has 290 tests across 17 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 8 of 25, Q1.2 is 0, Q1.3 is 23 violations of 40 production
effect sites with clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0
of 8 executable harnesses. The clean gate passed, the full audit exited 1 for
15 documented findings and never 2, and comparable ratchets were five improved,
two held, and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the exact filesystem Mkdir operation, mandatory
complete-double compatibility, local-config directory-create boundary,
contracts, and measured planning notes. Then perform the normal separate
handoff-only commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, all command objects and registrations, public API/CLI behavior,
all four host acceptance flows, `LocalConfigFile`, `LocalConfigDir`,
`OpenLocalConfig`, every LocalConfig method, context/profile behavior,
ProjectConfig, CloudConfig, Maven, structurizr and tips behavior, and every
caller's observable behavior.

Add `Mkdir(string, fs.FileMode) error` to the existing complete
`filesystem.FileSystem` interface, add the matching zero-value-safe package
operation returning exact `filesystem.ErrNoFilesystem`, and map the system
implementation directly to `os.Mkdir`. Keep established `MkdirAll` unchanged.
Update all 34 current complete filesystem test doubles only as required to
remain complete: the new method must record it for its dedicated contracts or
reject it consistently as unrelated everywhere else. Do not store a
function-valued dependency, split the complete dependency, add recursive
behavior, or expose application API.

Preserve exact `localCfgDir.Implementation().Path` evaluation, the completed
one-attempt adapter `Stat` with that exact arbitrary string, its exact
`fs.FileInfo` and error, and the exact `os.IsNotExist(err)` decision. Nil stat,
existing-directory stat, and arbitrary non-not-exist stat errors must still
return nil without mkdir. Only a not-exist error selects
`filesystem.System()` for a private complete local-config-directory-create
composition and makes one adapter mkdir attempt with the exact `dir` and
`0755`. Preserve the exact mkdir error as the method result and nil on success.

Do not clean or join the path, preflight, retry, normalize, validate, create
parents, tolerate an existing directory, change permission semantics, wrap an
error, move the stat or predicate, merge stat and mkdir dependencies, or
broaden into local-config path, YAML, logging, create, write, Close, touch,
update, read, print, existence, context, profile, Maven, structurizr, tips,
template, file, HTTP, Wpost, process, clock, server, P4, P5, dependency, Docker,
distribution, or publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete
`pkg/config/local.go` and all config tests, every local-config caller, complete
`pkg/context/context.go`, `cmd/profile.go`, and relevant registration, command,
and context tests. Read complete `internal/adapter/filesystem` source and tests,
every complete filesystem implementation and double, relevant
filesystem/file/template/config/Maven/tips/structurizr code and tests,
`internal/adapter/process`, and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit
implementation before editing. Before the full gate, read the complete launcher
contract, relevant Make meta-tests, P2A API/CLI/subprocess contracts, and all
four host acceptance flows.

# Three Moves

1. Start red with focused filesystem adapter Mkdir contracts. Prove exact
   arbitrary paths, exact modes, non-empty recorded populations, one attempt,
   exact error identity, safe zero behavior without developer-path access, and
   system single-level `os.Mkdir` behavior in temporary directories: missing
   parents fail, existing directories fail, and no recursive creation occurs.
   Add focused local-config-directory-create recording contracts that prove the
   complete dependency, production `filesystem.System()` selection, exact path
   and `0755`, one attempt, exact nil/arbitrary error results, non-empty
   populations, and no unrelated operation. Production-composition contracts
   must not mutate the real filesystem or run another command.

2. Add only the exact adapter `Mkdir` operation and mandatory method
   implementations on all current complete test doubles. Keep exported
   `LocalConfigDir.CheckOrCreateConfigDir() error` as the production entry,
   select `filesystem.System()` only for a private complete directory-create
   composition, and replace only its direct `os.Mkdir`. Do not change prior
   stat/predicate evaluation, error selection, another config effect, adapter
   operation, method, caller, inventory, or public API.

3. Run focused config/filesystem contracts and relevant process, command,
   context, file, template, Maven, tips, and structurizr caller packages'
   tests, the launcher contract from `/bin/bash`, Make preflight meta-contracts,
   API/CLI compatibility, full Go tests and race/vet, all four host acceptance
   flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.1 to hold at 8 of 25 and Q1.3 to
   move nominally from 23 of 40 to 22 of 39 because only the direct mkdir leaves
   the population. Q0.6, Q1.2, Q1.4, and exact Q2.1 must hold. Regenerate exact
   values; the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent filesystem
Mkdir and local-config directory-create move or record an exact resumable
state. Rewrite the rolling handover, record the measured P3 result, answer this
archive, create one linked NEXT archive for the next coherent P3 effect move,
replace only the launcher's mutable regions, run the launcher contract, and
make the separate handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
