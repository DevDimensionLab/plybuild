# Agent Session: Migrate Local Config Update Create

Status: NEXT
Session ID: `2026-08-25T125300+0200-migrate-local-config-update-create`
Created: `2026-08-25T12:53:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `e8aafbce4daa8480329b1d8a5714573b2bff6de3b0805f4e1b4dab6a3b254ecc`
Previous: [2026-08-25T122224+0200-migrate-local-config-touch-create.md](2026-08-25T122224+0200-migrate-local-config-touch-create.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct
`os.Create(configFilePath)` in `LocalConfigDir.UpdateLocalConfig` through the
existing filesystem adapter. Preserve exact directory and file-path selection,
the caller-supplied config, YAML marshaling, logging, create result and error,
the completed update write, write-error lifecycle, final Close result, every
caller, the completed local-config directory-stat, touch-write, update-write,
and touch-create moves, every completed filesystem, file, template, config,
Maven, tips, structurizr, Bitbucket, Wpost, and supervisor move, and every P2A
contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.Create(configFilePath)` operation in `LocalConfigDir.UpdateLocalConfig` in
`pkg/config/local.go` and its use of the existing zero-value-safe filesystem
`Create` operation. It authorizes focused config recording contracts and a
dedicated file double for that create boundary and its preserved Close result.
It does not authorize changing `CheckOrCreateConfigDir`, `FilePath`, the
caller-supplied config, YAML marshaling, logging, the completed update
`WriteFile`, its bytes/mode/error, Close selection or result, `TouchFile`,
another config method, an adapter or established complete-double extension,
another production effect, a function-valued effect dependency, a new adapter
family or public API, mutation harnesses, or later roadmap implementation.

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

Implementation commit `b232dda` has 286 tests across 17 of 25 packages. Q0.6
has 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe direct
test writes. Q1.1 is 8 of 25, Q1.2 is 0, Q1.3 is 25 violations of 42 production
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
implementation commit for the local-config update-create boundary, package
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
`Create(name string) (filesystem.File, error)` operation with a safe zero
value; do not add or change an adapter operation, extend an established
complete recording double, or store a function-valued effect dependency.
Production must select `filesystem.System()` only for a private complete
local-config-update-create composition.

Preserve exact `CheckOrCreateConfigDir()` evaluation and early error, exact
`FilePath()` result, the exact caller-supplied `config`,
`yaml.Marshal(&config)` and early error, and exact update log in their existing
order. Make one adapter create attempt with the exact `configFilePath` and
preserve its exact `filesystem.File` and error at the private boundary.
Preserve the enclosing lifecycle: a non-nil create error returns exactly
before the completed write; create success invokes the existing update-write
composition once with the exact path, bytes, and `0644`; a non-nil write error
returns exactly without Close; write success returns the exact final
`f.Close()` result.

Do not clean or join paths, preflight, retry, normalize, validate, create an
extra file, change truncation or permission behavior, write through the
created file, defer or add Close, close after a write error, wrap an error,
move the write, change YAML/logging, combine create and write dependencies, or
broaden into directory creation, touch, read, print, existence, context,
profile, Maven, structurizr, tips, template, file, HTTP, Wpost, clock, server,
P4, P5, dependency, Docker, distribution, or publication work.

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

1. Start red with focused local-config-update-create recording contracts.
   Prove the complete dependency, production selection of
   `filesystem.System()`, exact arbitrary file paths, non-empty recorded create
   populations, one attempt, exact `filesystem.File` and error identity for
   nil, successful, arbitrary-error, and unusual combined results, safe zero
   behavior without developer-path access, the dedicated file double's exact
   Close result, and no unrelated adapter operation. Production-composition
   contracts must not mutate the real filesystem or run another command.

2. Reuse only existing filesystem adapter `Create`. Keep exported
   `LocalConfigDir.UpdateLocalConfig(config LocalConfiguration) error` as the
   production entry, select `filesystem.System()` only for a private complete
   local-config-update-create composition, and replace only its direct
   `os.Create`. Do not change prior evaluation, completed WriteFile evaluation,
   bytes, mode, errors, final Close selection/result, another config effect,
   adapter, complete double, method or caller, inventory, or public API.

3. Run focused config/filesystem contracts and relevant process, command,
   context, file, template, Maven, tips, and structurizr caller packages'
   tests, the launcher contract from `/bin/bash`, Make preflight meta-contracts,
   API/CLI compatibility, full Go tests and race/vet, all four host acceptance
   flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.1 to hold at 8 of 25 and Q1.3 to
   move nominally from 25 of 42 to 23 of 40 because the concrete `*os.File`
   Close classification leaves with the direct create while the same lifecycle
   call remains. Q0.6, Q1.2, Q1.4, and exact Q2.1 must hold. Regenerate exact
   values; the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent local-config
update-create move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
