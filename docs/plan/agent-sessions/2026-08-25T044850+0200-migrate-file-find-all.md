# Agent Session: Migrate Exported File Find All Walk

Status: ANSWERED - HISTORY
Session ID: `2026-08-25T044850+0200-migrate-file-find-all`
Created: `2026-08-25T04:48:50+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `c8cba210ebcd957ada77e1c14c84f8b571501cb5b0436aa7d0a4536586a7a85a`
Previous: [2026-08-25T041554+0200-migrate-file-find-first.md](2026-08-25T041554+0200-migrate-file-find-first.md)
Next: [2026-08-25T052243+0200-migrate-file-grep-recursive.md](2026-08-25T052243+0200-migrate-file-grep-recursive.md)
Outcome: completed in `04cfe44bf6750d9309f774ffec21a6af576db0bd`; exported find-all now delegates its exact recursive walk through the existing zero-value-safe filesystem adapter, preserving traversal and callback order, suffix and exclusion decisions, ordered all-match accumulation, named results, exact errors, and final `io.EOF` normalization, with the measured P3.25 result recorded in the rolling plan and handover.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route exported `file.FindAll`'s
recursive walk through the existing filesystem adapter. Preserve the exact root
path, traversal and callback order, suffix and exclusion matching, ordered
all-match accumulation, named results, returned errors, and `io.EOF`
normalization; the completed find-first, render, clear-directory, file-move,
recursive-delete, single-file-delete, append-open, file-read, directory-create,
file-create, file-overwrite, file-existence, Bitbucket, Wpost, and supervisor
moves; and every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 are queued in the machine-readable block in
docs/plan/quality-upgrade.md, and the launcher must remain NEXT until every
authorized checkpoint is complete. This mission authorizes only the direct
`filepath.Walk(dir, callback)` operation in exported `FindAll` in
`pkg/file/file.go` and reuse of the existing zero-value-safe filesystem adapter
`Walk` operation. It does not authorize `FindFirst`, `SuffixIn`,
`GrepRecursive`, IDE cleanup, another file operation or Q1.3 flow, another
adapter operation, changes to complete filesystem recording doubles, a new
adapter family or public API, mutation harnesses, or later roadmap
implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, the roadmap queue, the complete `pkg/file` implementation and tests,
all callers of `file.FindAll`, the filesystem adapter and every complete
filesystem recording double, `.quality/inventory`, and the continuity and
quality-lift designs. Regenerate ignored reports outside the measured tree or
remove them before a clean audit.

Implementation commit `549685d` has 202 tests across 16 of 25 packages. Q0.6
has 20 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, Q1.3 is 42 violations of 58 production effect sites with
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
implementation commit for the exported find-all walk boundary, recording
contracts, and measured planning notes. Then perform the normal separate
handoff-only commit. Do not push, merge, publish, distribute, remove the
worktree, stash inherited changes, revert user work, or run destructive Git
commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the full 62-contract
launcher supervisor, exact prompt/archive bytes, the authorized queue and
reciprocal archive graph. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, all four host acceptance flows, the
exported
`file.FindAll(suffix string, excludes []string, dir string) (result []string, err error)`
signature, its cmd, config, context, template, and IDE callers, `FindFirst`,
`SuffixIn`, `GrepRecursive`, `Render`, `ClearDir`, `Move`, and every caller's
observable behavior.

Use one complete, resolvable private filesystem dependency value and a safe
zero value; do not store a function-valued effect dependency. Production must
select `filesystem.System()`. Delegate the exact root path and callback through
the existing adapter `Walk` operation without changing the callback body. Keep
ignoring the callback's `fi` and `errIn` inputs; match only
`strings.HasSuffix(path, suffix) && !SuffixIn(path, excludes)`; preserve
`SuffixIn`'s exact ordered `strings.Contains` behavior for every exclusion;
append every exact matching path in callback order; and always return nil from
the callback. Preserve a nil result when nothing matches, every accumulated
path with any final error, and every exact error. Convert only a final walk
error equal to `io.EOF` to nil. Do not normalize, wrap, retry, preflight, sort,
deduplicate, clean paths, inspect callback metadata, stop on a match, or add
selection behavior. Do not broaden into another file operation, Bitbucket,
HTTP/process effects, config behavior, Wpost, clock, server, P4 adapters, P5
harnesses, dependencies, Docker, cloud, distribution, or publication.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, `pkg/file`, all callers of
`file.FindAll`, and `internal/adapter/filesystem`. Read the relevant tests,
every filesystem test double that implements the complete adapter, and the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 audit implementation before editing.
Before the full gate, read the complete launcher contract, relevant Make
meta-tests, P2A API/CLI/subprocess contracts, and all four host acceptance
flows.

# Three Moves

1. Start red with recording find-all contracts. Prove the complete dependency,
   production selection of `filesystem.System()`, exact root and callback
   delivery, ordered callback paths, ignored callback metadata and errors,
   exact case-sensitive suffix decisions, exact ordered substring exclusion
   decisions including an empty exclusion, ordered accumulation of every
   match, no-match nil result and nil error, exact non-EOF walk errors with the
   current partial result, final `io.EOF` normalization, the safe zero value
   with no developer-path access, and non-empty recorded walk and callback
   populations. Perform no real filesystem mutation in these contracts.

2. Keep exported `file.FindAll` as the production wrapper, add one private
   complete filesystem dependency, and route only its exact root and unchanged
   callback through existing `filesystem.Walk`. Remove only its direct
   `filepath.Walk`. Do not extend the adapter or complete doubles, move
   `FindFirst`, `SuffixIn`, `GrepRecursive`, or another operation, or add
   another adapter family or public API.

3. Run focused file/filesystem contracts and all caller packages' relevant
   tests, the launcher contract from `/bin/bash`, Make preflight meta-contracts,
   API/CLI compatibility, full Go tests and race/vet, all four host acceptance
   flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Q1.3 to move nominally from 42 of 58
   to 41 of 57 while Q1.2, Q1.4, and exact Q2.1 hold. Regenerate exact values;
   the full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent `file.FindAll`
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P3 result, answer this archive, create one linked NEXT archive for
the next coherent P3 effect move, replace only the launcher's mutable regions,
run the launcher contract, and make the separate handoff-only commit
`docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
