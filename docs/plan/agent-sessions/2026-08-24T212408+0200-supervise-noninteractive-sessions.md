# Agent Session: Supervise Noninteractive Sessions

Status: NEXT
Session ID: `2026-08-24T212408+0200-supervise-noninteractive-sessions`
Created: `2026-08-24T21:24:08+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `df65ec792a101c934fb0c8ee23ce7c00270b9a78d409e715b03b089bd3ddbcde`
Previous: [2026-08-24T200448+0200-migrate-wpost-http-filesystem.md](2026-08-24T200448+0200-migrate-wpost-http-filesystem.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused continuity move: change normal `codex-dev-start.sh`
execution from an interactive one-shot into a fail-closed, observable
non-interactive supervisor loop. Run one authorized archived mission per Codex
turn, continue only after a valid committed handoff, and stop when the full
authorized roadmap reaches COMPLETE or when any failure/no-progress condition
appears. Preserve the completed Wpost move and every P2A contract with zero
comparable ratchet regressions.

# Authorized Roadmap

The user explicitly requested this launcher loop. P3 remains active, P4-P8 are
queued in the machine-readable block in docs/plan/quality-upgrade.md, and the
launcher must remain NEXT until every authorized checkpoint is complete. This
is an operational continuity move before the next P3 production-effect move;
it does not authorize adapter, quality-denominator, mutation-harness, or later
roadmap implementation.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, this active
archive, every historical session archive, the roadmap queue, all of
`codex-dev-start.sh` and `test/codex_dev_start_test.sh`, Make launcher targets,
the continuity and quality-lift designs, local `codex exec --help`, and the
official non-interactive-mode documentation. Regenerate ignored reports
outside the measured tree or remove them before a clean audit.

Implementation commit `a7eb3ef` has 130 tests across 16 of 25 packages. Q0.6
has 18 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, Q1.3 is 58 violations of 68 production effect sites with
clock and server absent, Q1.4 is 7 of 8, and exact Q2.1 is 0 of 8 executable
harnesses. The clean gate passed, the full audit exited 1 for 16 documented
findings and never 2, and comparable ratchets were five improved, two held,
and zero regressed.

The launcher has 49 Bash 3.2 contracts, three public modes, one immutable
stable execution skeleton, two inert mutable regions, exact prompt/archive
bytes, and a connected single-NEXT archive graph. Normal execution still
invokes interactive `codex -c 'service_tier="default"' -C <repo> <prompt>`.
Local `codex-cli 0.149.0` exposes `codex exec --json`,
`--sandbox workspace-write`, `-C`, and `--output-last-message`. Official docs
describe JSONL progress and terminal turn events; characterize the actual
local contract with a recording executable rather than trusting prose alone.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for stable launcher execution, its contract tests, and
the continuity design. Then perform the normal separate handoff-only commit.
Do not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands. Do not launch a
real successor while developing, testing, or finalizing this move.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve `--check`,
`--print-prompt`, `--help`, unknown/multiple-argument rejection, attached
branch/root validation, the authorized queue, source and planning symlink
guards, archive filename/metadata/digest/reciprocal-history validation, exact
launcher/archive prompt bytes with terminal LF, normal service tier, inert
mutable data, source-archive/Docker-copy coverage, dirty-state measurement, and
external executable resolution. Preserve `cmd.Execute()`, `cmd.ExecuteE()`,
`cmd.RootCmd`, public API/CLI behavior, and all four host acceptance flows.

Change only normal start execution and the minimum validation/logging support
it needs. Do not move prompt authority, add `.agent-task`, resume old Codex
threads, infer success from final prose, inject Git filenames into prompts,
write logs or state inside the repository, auto-publish, or broaden into
Wpost/HTTP/filesystem behavior, remaining Q1.3 effects, P4 adapters, P5
harnesses, dependencies, Docker, cloud, server, or distribution.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/agent-session-continuity.md,
docs/design/quality-lift.md, `.quality/inventory`, the complete launcher and
launcher test, every session archive and its lifecycle fields, relevant Make
targets/meta-tests, the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 implementation in
`.quality/tools`, local `codex exec --help` and version output, and the official
non-interactive-mode documentation before editing. Read the P2A compatibility,
subprocess, and four host acceptance contracts before the full gate.

# Supervisor Contract

1. Normal invocation starts an explicit non-interactive Codex turn with the
   exact validated prompt, repository working directory, normal service tier,
   workspace-write sandbox, and JSONL output. It never uses interactive mode or
   resumes a prior thread.
2. Stream concise, useful progress to the terminal while preserving the raw
   event stream in a uniquely named log directory outside the worktree. Print
   the log location. Validate JSON rather than grepping substrings, preserve
   agent messages as data, and never execute event content.
3. Treat a turn as successful only when the Codex process exits zero, the
   stream is well-formed and has the characterized successful terminal event,
   and it contains no failure/error terminal condition. Signals, malformed or
   truncated streams, contradictory terminals, `turn.failed`, and non-zero
   exits stop non-zero with actionable diagnostics.
4. Snapshot HEAD and active session identity before each turn. After the agent
   process has ended, re-read the on-disk launcher rather than the stale source
   snapshot and validate the worktree, queue, complete archive graph, active
   prompt, and launcher contract. Agent final text is observable but cannot
   override repository evidence.
5. Continue only when the worktree is clean, HEAD changed, the former archive
   is answered, the session ID changed, and exactly one new valid NEXT tail is
   committed with reciprocal history. Refuse replay, unchanged HEAD/session,
   dirty or detached state, partial handoff, invalid graph, and contract drift.
6. When the freshly validated launcher and authorized queue reach COMPLETE,
   exit zero after reporting completion. COMPLETE remains impossible while P3
   or any P4-P8 checkpoint is active or queued.
7. Keep the agent-side instruction not to launch its successor: the parent
   supervisor alone may begin the next fresh process, and only after the prior
   process ended and the committed handoff passed every progression check.
8. Preserve safe interruption. A user signal stops the active child and
   supervisor without starting another turn, leaves the raw log inspectable,
   and does not mutate task authority.

# Three Moves

1. Start red by extending the recording launcher contracts. Prove exact
   `codex exec` argv and prompt bytes, workspace-write and service-tier
   selection, concise progress and raw external logging, a successful handoff
   followed by a second generation, eventual COMPLETE, and non-execution of
   event text. Prove fail-closed behavior for turn failure/error, malformed or
   truncated JSONL, non-zero exit, no progress, unchanged session, dirty or
   invalid handoff, failed post-turn contract, and signal interruption. Reject
   empty recorded event and loop populations.

2. Implement the smallest Bash 3.2-compatible supervisor and structured-event
   parser that satisfies those contracts. Reuse the existing validation and
   immutable-tail model. Refactor validation only as needed to re-snapshot the
   changed launcher between turns. Keep test-only controls unreachable from
   ambient production state and keep logs outside the measured tree.

3. Run the focused launcher contract from `/bin/bash`, Make launcher and
   preflight meta-contracts, API/CLI compatibility, full Go tests and race/vet,
   all four host acceptance flows, the audit meta-suite, focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurements, full clean checkpoint
   audit, and empty-HOME count-2. Expect Go quality denominators and P3 ratios
   to hold exactly at the start values, with only documented plan/archive
   counts changing. The full audit may exit 1 for documented findings but
   never 2. Never invoke the real Codex executable in a test.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent launcher move or
record an exact resumable state. Rewrite the rolling handover, record the
launcher measurements without relabeling P3 quality progress, answer this
archive, create one linked NEXT archive for the next coherent P3 effect move,
replace only the launcher's mutable regions, run the launcher contract, and
make the separate handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
