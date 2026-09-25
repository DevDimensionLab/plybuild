# Agent Session: Migrate Exported File Clear Directory

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T030337+0200-migrate-file-clear-dir`
Created: `2026-08-25T03:03:37+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `a35a4bf30534736dd7fd776ad13d05667b59eb28bd0da8a42bf707eb58049892`
Previous: [2026-08-25T023033+0200-migrate-file-move.md](2026-08-25T023033+0200-migrate-file-move.md)
Next: [2026-08-25T034048+0200-migrate-file-render-create.md](2026-08-25T034048+0200-migrate-file-render-create.md)
Outcome: completed in `838daa153a29d13603a9f3bc1d532321f6db64d9`; exported clear-directory selection and repeated removals now delegate the exact pattern, ordered paths, and exact errors through the existing filesystem adapter, with the measured P3.22 result recorded in the rolling plan and handover.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route exported
`file.ClearDir`'s glob selection and repeated recursive removals through the
existing filesystem adapter. Preserve the exact glob pattern, result and
removal order, exclusion decisions, logging, paths, and returned errors, the
completed file-move, recursive-delete, single-file-delete, append-open,
file-read, directory-create, file-create, file-overwrite, file-existence,
Bitbucket, Wpost, and supervisor moves, and every P2A contract with zero
comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`filepath.Glob(filepath.Join(dirPath, "*"))` selection and repeated
`os.RemoveAll(file)` operation in exported `ClearDir` in `pkg/file/file.go`,
the narrow `Glob` operation required in the existing filesystem adapter, and
reuse of that adapter's existing `RemoveAll`. It does not authorize another
file operation or Q1.3 flow, a new adapter family or public API, mutation
harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/file` implementation and tests,
all callers of `file.ClearDir`, the filesystem adapter and every complete
filesystem recording double, `.quality/inventory`, and the continuity and
quality-lift designs. Regenerate ignored reports outside the measured tree or
remove them before a clean audit.

Implementation commit `acda4e3` has 180 tests across 16 of 25 packages. Q0.6
has 18 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, Q1.3 is 46 violations of 60 production effect sites with
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
implementation commit for the exported clear-directory boundary, its narrow
filesystem adapter operation, recording contracts, and measured planning
notes. Then perform the normal separate handoff-only commit. Do not push,
merge, publish, distribute, remove the worktree, stash inherited changes,
revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows, the
exported `file.ClearDir(dirPath string, excludes []string) error` signature,
`Move`, `DeleteAll`, `DeleteSingleFile`, and every existing caller's
observable behavior.

Use one complete, resolvable private filesystem dependency value and a safe
zero value; do not store a function-valued effect dependency. Production must
select `filesystem.System()`. Delegate the exact joined glob pattern and each
selected removal path unchanged. Preserve glob result order, the complete
ordered exclusion loop using `strings.Contains(file, exclude)`, every existing
skip/removal log, the first removal-error short-circuit, nil success, and exact
glob or removal error identity. Do not normalize, sort, deduplicate, wrap,
retry, preflight, continue after a removal error, or add selection behavior.
Do not broaden into another file operation, Bitbucket, HTTP/process effects,
config behavior, Wpost, clock, server, P4 adapters, P5 harnesses, dependencies,
Docker, cloud, distribution, or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, `pkg/file`, all callers of
`file.ClearDir`, and `internal/adapter/filesystem`. Read the relevant tests,
every filesystem test double that implements the complete adapter, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with recording clear-directory contracts. Prove the complete
   dependency, exact glob pattern, ordered glob results, exact ordered removal
   paths, unchanged exclusion and logging behavior, exact glob and removal
   error identities, the first removal-error short-circuit, successful nil
   error, the safe zero value with no developer path access, production
   selection of `filesystem.System()`, and non-empty recorded glob and removal
   populations. Perform no real filesystem mutation in these file contracts.

2. Add only the smallest zero-value-safe `Glob` operation to the existing
   filesystem adapter and its system implementation, updating complete
   recording doubles mechanically and reusing existing adapter `RemoveAll`.
   Keep exported `file.ClearDir` as the production wrapper and use one complete
   private dependency to drive selection and removal. Remove only its direct
   `filepath.Glob` and `os.RemoveAll`. Do not move any other operation, and do
   not add a new adapter family or public API.

3. Run focused file/filesystem contracts and all callers' relevant tests, the
   launcher contract from `/bin/bash`, Make preflight meta-contracts, API/CLI
   compatibility, full Go tests and race/vet, all four host acceptance flows,
   the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 46 of 60
   to 44 of 60 while Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact values;
   the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent `file.ClearDir`
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P3 result, answer this archive, create one linked NEXT archive for
the next coherent P3 effect move, replace only the launcher's mutable regions,
run the launcher contract, and make the separate handoff-only commit
`docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
