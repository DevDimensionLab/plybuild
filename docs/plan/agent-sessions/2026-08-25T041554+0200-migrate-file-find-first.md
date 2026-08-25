# Agent Session: Migrate Exported File Find First Walk

Status: NEXT
Session ID: `2026-08-25T041554+0200-migrate-file-find-first`
Created: `2026-08-25T04:15:54+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `0acefb5faed1519369214202e8cb04299df5788a766429168e54981dbb736dda`
Previous: [2026-08-25T034048+0200-migrate-file-render-create.md](2026-08-25T034048+0200-migrate-file-render-create.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route exported
`file.FindFirst`'s recursive walk through the existing filesystem adapter.
Preserve the exact root path, traversal and callback order, suffix matching,
first-match stop, named results, returned errors, and `io.EOF` normalization;
the completed render, clear-directory, file-move, recursive-delete,
single-file-delete, append-open, file-read, directory-create, file-create,
file-overwrite, file-existence, Bitbucket, Wpost, and supervisor moves; and
every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`filepath.Walk(dir, callback)` operation in exported `FindFirst` in
`pkg/file/file.go`, the narrow zero-value-safe `Walk` operation required in the
existing filesystem adapter, and the mechanical update of every complete
filesystem recording double. It does not authorize `FindAll`,
`GrepRecursive`, IDE cleanup, another file operation or Q1.3 flow, a new
adapter family or public API, mutation harnesses, or later roadmap
implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/file` implementation and tests,
all callers of `file.FindFirst`, the filesystem adapter and every complete
filesystem recording double, `.quality/inventory`, and the continuity and
quality-lift designs. Regenerate ignored reports outside the measured tree or
remove them before a clean audit.

Implementation commit `224a691` has 194 tests across 16 of 25 packages. Q0.6
has 18 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, Q1.3 is 43 violations of 58 production effect sites with
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
implementation commit for the exported find-first walk boundary, its narrow
filesystem adapter operation, recording contracts, complete-double updates,
and measured planning notes. Then perform the normal separate handoff-only
commit. Do not push, merge, publish, distribute, remove the worktree, stash
inherited changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows, the
exported
`file.FindFirst(fileSuffix string, dir string) (result string, err error)`
signature, its config and Spring callers, `FindAll`, `GrepRecursive`, IDE
cleanup, `Render`, `ClearDir`, `Move`, and every caller's observable behavior.

Use one complete, resolvable private filesystem dependency value and a safe
zero value; do not store a function-valued effect dependency. Production must
select `filesystem.System()`. Delegate the exact root path and callback through
the narrow adapter `Walk` operation without changing the callback body. Keep
ignoring the callback's `fi` and `errIn` inputs, match only
`strings.HasSuffix(path, fileSuffix)`, assign the exact matching path before
returning `io.EOF`, stop on the first matching callback, convert only a final
walk error equal to `io.EOF` to nil, and preserve every other exact error and
named result. Do not normalize, wrap, retry, preflight, sort, add filtering,
inspect callback metadata, or add selection behavior. Do not broaden into
another file operation, Bitbucket, HTTP/process effects, config behavior,
Wpost, clock, server, P4 adapters, P5 harnesses, dependencies, Docker, cloud,
distribution, or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, `pkg/file`, all callers of
`file.FindFirst`, and `internal/adapter/filesystem`. Read the relevant tests,
every filesystem test double that implements the complete adapter, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with adapter and recording find-first contracts. Prove the
   complete dependency, production selection of `filesystem.System()`, exact
   root and callback delivery, ordered callback paths, ignored callback
   metadata and errors, exact suffix matching, first-match result and stop,
   no-match empty result and nil error, exact non-EOF walk errors and current
   partial result, final `io.EOF` normalization, the safe zero value with no
   developer path access, and non-empty recorded walk and callback populations.
   Perform no real filesystem mutation in these file contracts.

2. Extend only the existing filesystem adapter with zero-value-safe `Walk` and
   its single system `filepath.Walk` delegation. Keep exported `file.FindFirst`
   as the production wrapper and use one complete private dependency for the
   walk. Remove only its direct `filepath.Walk`. Update every complete
   filesystem recording double mechanically. Do not move `FindAll` or another
   operation, or add another adapter family or public API.

3. Run focused file/filesystem contracts and both caller packages' relevant
   tests, the launcher contract from `/bin/bash`, Make preflight meta-contracts,
   API/CLI compatibility, full Go tests and race/vet, all four host acceptance
   flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 43 of 58
   to 42 of 58 while Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact values;
   the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent
`file.FindFirst` move or record an exact resumable state. Rewrite the rolling
handover, record the measured P3 result, answer this archive, create one linked
NEXT archive for the next coherent P3 effect move, replace only the launcher's
mutable regions, run the launcher contract, and make the separate handoff-only
commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
