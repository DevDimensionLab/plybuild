# Agent Session: Migrate Exported Directory Creation

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T000314+0200-migrate-file-directory`
Created: `2026-08-25T00:03:14+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5ed3d0551037823d990e12e0cec3094a214dc793cee463ad92d430ffdea98bf5`
Previous: [2026-08-24T233519+0200-migrate-file-create.md](2026-08-24T233519+0200-migrate-file-create.md)
Next: [2026-08-25T003130+0200-migrate-file-open.md](2026-08-25T003130+0200-migrate-file-open.md)
Outcome: exported directory creation moved behind the filesystem adapter at `e054082`; the clean gate passed

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route exported
`file.CreateDirectory` through the existing filesystem adapter. Preserve its
missing-only creation selection, directory mode, exact legacy error behavior,
the completed file-create, file-overwrite, file-existence, Bitbucket, Wpost,
and supervisor moves, and every P2A contract with zero comparable ratchet
regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`os.Stat` and `os.MkdirAll` in exported `CreateDirectory` in
`pkg/file/file.go`; it does not authorize `OpenFile`'s `os.OpenFile`, another
file operation or Q1.3 flow, a new adapter family or public API, mutation
harnesses, or later roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/file` implementation and tests,
all callers of `file.CreateDirectory`, the filesystem adapter and contracts,
`.quality/inventory`, and the continuity and quality-lift designs. Regenerate
ignored reports outside the measured tree or remove them before a clean audit.

Implementation commit `60e5aac` has 148 tests across 16 of 25 packages. Q0.6
has 18 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, Q1.3 is 54 violations of 64 production effect sites with
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
implementation commit for the exported directory-create boundary, its
recording contracts, and measured planning notes. Then perform the normal
separate handoff-only commit. Do not push, merge, publish, distribute, remove
the worktree, stash inherited changes, revert user work, or run destructive
Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows, the
exported `file.CreateDirectory(path string) error` signature, `file.CreateFile`,
`file.OpenFile`, and every existing caller's observable behavior.

Use a complete, resolvable private filesystem dependency value and a safe zero
value; do not store a function-valued effect dependency. Production must
select `filesystem.System()`. Preserve the legacy behavior exactly: create
only when stat reports a missing path; return nil for existing paths and other
stat errors; if recursive creation fails, return the original missing-path
stat error rather than the creation error. Do not broaden into another file
operation, Bitbucket, HTTP/process effects, config behavior, Wpost, clock,
server, P4 adapters, P5 harnesses, dependencies, Docker, cloud, distribution,
or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, `pkg/file`, all callers of
`file.CreateDirectory`, and `internal/adapter/filesystem`. Read the relevant
tests and the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before
editing. Before the full gate, read the complete launcher contract, relevant
Make meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with recording directory-create contracts. Prove the complete
   path and dependency value; missing-path selection of exactly one
   `MkdirAll` with mode `0755`; existing-path and non-missing stat-error
   selection of no creation and nil; preservation of the original missing
   stat error when creation fails; the safe zero value with no mutation;
   production selection of `filesystem.System()`; and non-empty recorded stat
   and creation populations. Exercise no developer path and perform no real
   filesystem mutation.

2. Implement the smallest private, complete, resolvable, zero-value-safe
   dependency boundary that satisfies those contracts. Keep exported
   `file.CreateDirectory` as the production wrapper and delegate stat and
   recursive creation to the unchanged filesystem adapter. Remove only its
   direct `os.Stat` and `os.MkdirAll`; do not move `OpenFile`'s `os.OpenFile`
   or another operation, and do not change the adapter contract unless the
   focused tests prove it necessary.

3. Run focused file/filesystem contracts, all callers' relevant tests, the
   launcher contract from `/bin/bash`, Make preflight meta-contracts, API/CLI
   compatibility, full Go tests and race/vet, all four host acceptance flows,
   the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 54 of 64
   to 52 of 62 while Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact values;
   the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent
`file.CreateDirectory` move or record an exact resumable state. Rewrite the
rolling handover, record the measured P3 result, answer this archive, create
one linked NEXT archive for the next coherent P3 effect move, replace only the
launcher's mutable regions, run the launcher contract, and make the separate
handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
