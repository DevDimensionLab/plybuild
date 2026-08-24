# Agent Session: Migrate Exported Append Open

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T010457+0200-migrate-file-open-file`
Created: `2026-08-25T01:04:57+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `628510ed60aed5506ca3dc9b0263d089422854fca8f9399bac6f7a1450f1ce69`
Previous: [2026-08-25T003130+0200-migrate-file-open.md](2026-08-25T003130+0200-migrate-file-open.md)
Next: [2026-08-25T013438+0200-migrate-file-delete-single.md](2026-08-25T013438+0200-migrate-file-delete-single.md)
Outcome: completed P3 move 18 at `c2f3597`; Q1.3 is 49 of 60 with zero comparable ratchet regressions

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route exported `file.OpenFile`'s
append-open operation through the existing filesystem adapter. Preserve its
exact missing-file selection, empty-file creation, flags, mode, returned file
and errors, the completed file-read, directory-create, file-create,
file-overwrite, file-existence, Bitbucket, Wpost, and supervisor moves, and
every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)` operation in exported
`OpenFile` in `pkg/file/file.go` and the narrow operation required in the
existing filesystem adapter. It does not authorize `pkg/shell`'s
`os.OpenFile`, another file operation or Q1.3 flow, a new adapter family or
public API, mutation harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/file` implementation and tests,
all callers of `file.OpenFile`, the filesystem adapter and every complete
filesystem recording double, `.quality/inventory`, and the continuity and
quality-lift designs. Regenerate ignored reports outside the measured tree or
remove them before a clean audit.

Implementation commit `27c0d1a` has 157 tests across 16 of 25 packages. Q0.6
has 18 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, Q1.3 is 50 violations of 60 production effect sites with
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
implementation commit for the exported append-open boundary, its narrow
filesystem adapter operation, recording contracts, and measured planning
notes. Then perform the normal separate handoff-only commit. Do not push,
merge, publish, distribute, remove the worktree, stash inherited changes,
revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows, the
exported `file.OpenFile(fileName string) (*os.File, error)` signature, `Exists`,
`CreateFile`, `Open`, and every existing caller's observable behavior.

Use one complete, resolvable private filesystem dependency value and a safe
zero value; do not store a function-valued effect dependency. Production must
select `filesystem.System()`. Preserve the legacy sequence exactly: probe
existence first; only a missing result creates exact empty bytes with mode
`0644`; a create error returns immediately without append-open; existing paths
and non-missing stat errors skip creation; and the final operation receives
the complete path, flags `os.O_APPEND|os.O_WRONLY`, and mode `0644`. Return the
dependency's exact `*os.File` and error pair. Do not broaden into another file
operation, Bitbucket, HTTP/process effects, config behavior, Wpost, clock,
server, P4 adapters, P5 harnesses, dependencies, Docker, cloud, distribution,
or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, `pkg/file`, all callers of
`file.OpenFile`, and `internal/adapter/filesystem`. Read the relevant tests,
every filesystem test double that implements the complete adapter, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with recording append-open contracts. Prove the complete
   dependency and path; exact flags and mode; exact returned file pointer and
   dependency error; missing, existing, and non-missing-stat-error sequencing;
   exact empty creation and create-error short circuit; the safe zero value
   with no developer path access; production selection of
   `filesystem.System()`; and non-empty recorded populations. Perform no real
   filesystem mutation.

2. Add only the smallest zero-value-safe `OpenFile` operation to the existing
   filesystem adapter and its system implementation, updating complete
   recording doubles mechanically. Keep exported `file.OpenFile` as the
   production wrapper and use one complete private dependency to drive its
   existence, conditional empty write, and append-open sequence. Remove only
   its direct `os.OpenFile`. Do not move `pkg/shell`'s `os.OpenFile` or any
   other operation, and do not add a new adapter family or public API.

3. Run focused file/filesystem contracts and all callers' relevant tests, the
   launcher contract from `/bin/bash`, Make preflight meta-contracts, API/CLI
   compatibility, full Go tests and race/vet, all four host acceptance flows,
   the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 50 of 60
   to 49 of 60 while Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact values;
   the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent
`file.OpenFile` move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
