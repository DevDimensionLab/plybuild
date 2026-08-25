# Agent Session: Migrate Local Config Directory Stat

Status: NEXT
Session ID: `2026-08-25T115639+0200-migrate-local-config-directory-stat`
Created: `2026-08-25T11:56:39+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `8f07b3dc0a3f8ee22d09f23e658b595051cdc6f0198a5186a1d9f62e627d3e12`
Previous: [2026-08-25T112922+0200-migrate-local-config-update-write.md](2026-08-25T112922+0200-migrate-local-config-update-write.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`os.Stat(dir)` in `LocalConfigDir.CheckOrCreateConfigDir` through the existing
filesystem adapter. Preserve exact directory selection, stat evaluation,
`os.IsNotExist` branching, directory creation, mode and errors, every caller,
both completed local-config write moves, every completed filesystem, file,
template, config, Maven, tips, structurizr, Bitbucket, Wpost, and supervisor
move, and every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.Stat(dir)` operation in `LocalConfigDir.CheckOrCreateConfigDir` in
`pkg/config/local.go` and its use of the existing zero-value-safe filesystem
`Stat` operation. It authorizes focused config recording contracts for that
stat boundary. It does not authorize changing `Implementation`, `FilePath`,
`os.IsNotExist`, `os.Mkdir`, directory mode, error selection, `TouchFile`,
`UpdateLocalConfig`, another config method, an adapter or established
complete-double extension, another production effect, a function-valued
effect dependency, a new adapter family or public API, mutation harnesses, or
later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, complete `pkg/config/local.go` and all config tests,
every caller and use of `OpenLocalConfig`, `LocalConfigDir`, `LocalConfigFile`,
`CheckOrCreateConfigDir`, `TouchFile`, `UpdateLocalConfig`, `Config`, `Print`,
and `Exists`, complete `pkg/context/context.go`, `cmd/profile.go`, their relevant
command and context tests, relevant filesystem/file/template/config/Maven/tips/
structurizr behavior, the filesystem and process adapters and relevant
recording doubles, `.quality/inventory`, and the continuity and quality-lift
designs. Regenerate ignored reports outside the measured tree or remove them
before a clean audit.

Implementation commit `fd45ebf` has 278 tests across 17 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 8 of 25, Q1.2 is 0, Q1.3 is 28 violations of 45 production
effect sites with clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0
of 8 executable harnesses. The clean gate passed, the full audit exited 1 for
15 documented findings and never 2, and comparable ratchets were five
improved, two held, and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the local-config directory-stat boundary, package
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

Use the existing complete filesystem dependency and existing
`Stat(name string) (fs.FileInfo, error)` operation with a safe zero value; do
not add or change an adapter operation, extend an established complete
recording double, or store a function-valued effect dependency. Production
must select `filesystem.System()` only for a private complete
local-config-directory-stat composition.

Preserve exact `localCfgDir.Implementation().Path` evaluation and the resulting
arbitrary `dir` string. Make one adapter stat attempt with that exact string and
preserve its exact `fs.FileInfo` and error at the private boundary. Preserve
the enclosing exact `os.IsNotExist(err)` decision: only a not-exist error
attempts exact `os.Mkdir(dir, 0755)`, returns its exact non-nil error, and
otherwise returns nil. Preserve nil stat results, existing-directory results,
not-exist errors, arbitrary other stat errors, unusual paths, and every later
TouchFile and UpdateLocalConfig effect and lifecycle.

Do not clean or join paths, preflight, retry, normalize, validate, change the
not-exist predicate, broaden stat errors, create an extra directory, inspect or
change permissions, replace or move `os.Mkdir`, change mode, wrap or return an
otherwise ignored stat error, move file creation, writing, or closing, or
broaden into touch, update, read, print, existence, context, profile, Maven,
structurizr, tips, template, file, HTTP, Wpost, clock, server, P4, P5,
dependency, Docker, distribution, or publication work.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, complete
`pkg/config/local.go` and all config tests, every local-config caller, complete
`pkg/context/context.go`, `cmd/profile.go`, and relevant registration, command,
and context tests. Read relevant filesystem/file/template/config/Maven/tips/
structurizr code and tests, `internal/adapter/filesystem`,
`internal/adapter/process`, their relevant complete recording doubles, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with focused local-config-directory-stat recording contracts.
   Prove the complete dependency, production selection of
   `filesystem.System()`, exact arbitrary directory paths, non-empty recorded
   stat populations, one attempt, exact `fs.FileInfo` and error identity for
   nil, existing, not-exist, and arbitrary error results, safe zero behavior
   without developer-path access, and no unrelated adapter operation.
   Production-composition contracts must not mutate the real filesystem or run
   another command.

2. Reuse only the existing filesystem adapter `Stat` operation. Keep exported
   `LocalConfigDir.CheckOrCreateConfigDir() error` as the production entry,
   select `filesystem.System()` only for a private complete
   local-config-directory-stat composition, and replace only its direct
   `os.Stat` with the adapter operation. Do not change directory evaluation,
   `os.IsNotExist`, `os.Mkdir` evaluation, mode or error, another config effect,
   adapter, complete double, method or caller, inventory, or public API.

3. Run focused config/filesystem contracts and relevant process, command,
   context, file, template, Maven, tips, and structurizr caller packages'
   tests, the launcher contract from `/bin/bash`, Make preflight
   meta-contracts, API/CLI compatibility, full Go tests and race/vet, all four
   host acceptance flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.1 to hold at 8 of 25 and Q1.3 to
   move nominally from 28 of 45 to 27 of 44 while Q0.6, Q1.2, Q1.4, and exact
   Q2.1 hold. Regenerate exact values; the full audit may exit 1 for documented
   findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent local-config
directory-stat move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
