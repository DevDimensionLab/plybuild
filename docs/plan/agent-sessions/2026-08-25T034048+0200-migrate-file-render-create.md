# Agent Session: Migrate Exported File Render Create

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T034048+0200-migrate-file-render-create`
Created: `2026-08-25T03:40:48+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `146c04e4c81a255d74968fb84ec9c0f9955ac1481af68ba952bd2f15b161d9d6`
Previous: [2026-08-25T030337+0200-migrate-file-clear-dir.md](2026-08-25T030337+0200-migrate-file-clear-dir.md)
Next: [2026-08-25T041554+0200-migrate-file-find-first.md](2026-08-25T041554+0200-migrate-file-find-first.md)
Outcome: completed in `224a691e1e43112fd8df36a7b4c579d60246fb5c`; exported render input-read and output-create now share one complete filesystem dependency, preserving exact paths, order, bytes, errors, parse panic, safe zero behavior, and the existing lack of a close, with the measured P3.23 result recorded in the rolling plan and handover.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route exported `file.Render`'s
output-file creation through the existing filesystem adapter. Preserve the
exact input and output paths, input-read and template-processing order,
creation and truncation behavior, rendered bytes, returned errors, and lack of
an added close; the completed clear-directory, file-move, recursive-delete,
single-file-delete, append-open, file-read, directory-create, file-create,
file-overwrite, file-existence, Bitbucket, Wpost, and supervisor moves; and
every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.Create(outputFilePath)` operation in exported `Render` in
`pkg/file/render.go`, reuse of the existing filesystem adapter `Create`, and
reuse of its already-migrated `ReadFile` through the same private complete
dependency solely to preserve the input boundary and a safe zero value. It
does not authorize another adapter operation, file operation or Q1.3 flow, a
new adapter family or public API, adding file closure, mutation harnesses, or
later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/file` implementation and tests,
all callers of `file.Render`, the filesystem adapter and every complete
filesystem recording double, `.quality/inventory`, and the continuity and
quality-lift designs. Regenerate ignored reports outside the measured tree or
remove them before a clean audit.

Implementation commit `838daa1` has 187 tests across 16 of 25 packages. Q0.6
has 18 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, Q1.3 is 44 violations of 59 production effect sites with
clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0 of 8 executable
harnesses. The clean gate passed, the full audit exited 1 for 16 documented
findings and never 2, and comparable ratchets were five improved, two held,
and zero regressed.

The launcher has 62 Bash 3.2 contracts and supervises fresh non-interactive
JSONL turns with external raw logs. It continues only after a successful
structured stream and valid clean committed handoff. Preserve its stable
skeleton and do not launch a real successor while developing, testing, or
finalizing this move.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the exported render-create boundary, recording
contracts, and measured planning notes. Then perform the normal separate
handoff-only commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows, the
exported
`file.Render(inputFilePath string, outputFilePath string, r interface{}) error`
signature, both existing template callers, `ClearDir`, `Move`, `DeleteAll`,
`DeleteSingleFile`, and every caller's observable behavior.

Use one complete, resolvable private filesystem dependency value and a safe
zero value; do not store a function-valued effect dependency. Production must
select `filesystem.System()`. Delegate the exact input and output paths
unchanged, using the existing adapter `ReadFile` before the existing adapter
`Create`. Preserve the exact input-read error, `template.Must` parse-panic
behavior, creation only after successful read and parse, output creation and
truncation, exact rendered bytes and writer execution-error behavior, nil
success, and the existing lack of a close. Do not normalize, wrap, retry,
preflight, reorder, add cleanup, or add selection behavior. Do not broaden into
another file operation, Bitbucket, HTTP/process effects, config behavior,
Wpost, clock, server, P4 adapters, P5 harnesses, dependencies, Docker, cloud,
distribution, or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, `pkg/file`, all callers of
`file.Render`, and `internal/adapter/filesystem`. Read the relevant tests,
every filesystem test double that implements the complete adapter, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with recording render contracts. Prove the complete dependency,
   exact input and output paths, read-before-create ordering, no create or
   write after an input error, no create after a template parse panic, exact
   create-error identity, exact rendered bytes, unchanged writer execution
   errors, successful nil error, the safe zero value with no developer path
   access, production selection of `filesystem.System()`, and non-empty
   recorded read, create, and write populations. Perform no real filesystem
   mutation in these file contracts.

2. Keep exported `file.Render` as the production wrapper and use one complete
   private dependency for its existing adapter `ReadFile` and existing adapter
   `Create` composition. Remove only its direct `os.Create(outputFilePath)`.
   Do not extend the adapter, add a close, move another operation, or add a new
   adapter family or public API.

3. Run focused file/filesystem contracts and both callers' relevant tests, the
   launcher contract from `/bin/bash`, Make preflight meta-contracts, API/CLI
   compatibility, full Go tests and race/vet, all four host acceptance flows,
   the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 44 of 59
   to 43 of 58 while Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact values;
   the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent `file.Render`
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P3 result, answer this archive, create one linked NEXT archive for
the next coherent P3 effect move, replace only the launcher's mutable regions,
run the launcher contract, and make the separate handoff-only commit
`docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
