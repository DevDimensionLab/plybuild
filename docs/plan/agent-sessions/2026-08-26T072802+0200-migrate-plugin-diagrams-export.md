# Agent Session: Migrate Plugin Diagrams Export

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T072802+0200-migrate-plugin-diagrams-export`
Created: `2026-08-26T07:28:02+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `1b6cf593ab559a4650717479f248b7ecb49ce676574e1d46f4511e734a6b50f8`
Previous: [2026-08-26T061848+0200-resume-p354-audit-harness-repair.md](2026-08-26T061848+0200-resume-p354-audit-harness-repair.md)
Next: [2026-08-26T080352+0200-migrate-plugin-diagrams-open.md](2026-08-26T080352+0200-migrate-plugin-diagrams-open.md)
Outcome: product commit `3972fc6` routes only the ignored structurizr export request through the process adapter; 364 tests, Q1.3 9/31, the clean gate, T15, full audit, and empty-HOME count-2 pass with zero comparable regressions

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the ignored
`structurizr-cli export` process request in `cmd/plugin_diagrams.go` through the
existing `internal/adapter/process` boundary. Preserve its exact command name,
ordered arguments, empty working directory and streams, synchronous execution,
one attempt, ignored error, and continued file discovery. Preserve the later
Graphviz `dot` and macOS `open` requests, Cobra behavior, every completed move,
and zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.54 product commit `3bd07e9` and focused
audit-apparatus commit `ce736a2` passed their truthful complete checkpoint.

This mission authorizes only private plugin-diagrams export recording
contracts, the private dependency composition required to select
`process.System()` for production, and replacement of this one request:

```go
_ = structurizr.Run(exec.Command("structurizr-cli", "export", "-w", workspace, "-format", "dot", "-output", tempDirectory))
```

with the equivalent ignored `process.Execute` request. It does not authorize
the Graphviz `dot` request, macOS `open` request, another plugin-diagrams flow,
`structurizr.Run` or `RunWithOutputToFile`, the process adapter, another
adapter/caller, public API, inventory, scanner, audit apparatus, clock/server,
mutation harnesses, or P4-P8 implementation.

# Measurements At Start

Checkpoint commit `ce736a2` has the unchanged P3.54 product tree with 360 tests
across 19 of 25 packages. Q0.6 has 26 guarded safe-writer sites, 21 write and 5
copy, zero skipped tests, and zero unsafe direct test writes. Q1.1 is 6/25,
Q1.2 is zero, Q1.3 is 10/32 with clock and server absent, Q1.4 is 7/8, exact
Q2.1 is 0/8, and Q3.4 is zero phrases across 72 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`c7ff3bd2cf32f8615c0fb1329c980bafba9b5e988cf4e785f4b8dc608d3f4243`.
Both complete repaired old/new T15 proofs pass all 15 controls and reproduce
the exact stored debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single export-process move and its recording
contracts, then make the normal separate continuity-only commit. Do not push,
merge, publish, distribute, remove the worktree, stash inherited changes,
revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the exact
`structurizr-cli` name and arguments `export`, `-w`, the caller-supplied
workspace, `-format`, `dot`, `-output`, `.structurizr/`. Preserve empty `Dir`,
nil stdin/stdout/stderr, `Start: false`, one synchronous attempt, ignored
process error, and unconditional continuation to `file.FindAll` after that
attempt. A recording test must reject an empty request population.

Select `process.System()` only in the production composition. Pass the complete
private dependency value without fallback, retry, logging, another process
attempt, error wrapping, environment or working-directory behavior, output
capture, or cleanup. The zero dependency must remain a safe no-op and must not
launch an external tool. Focused tests must launch no real process, touch no
network, and write no repository fixture.

Do not change the preceding mandatory workspace lookup, exact
`.structurizr/` deletion, following dot-file discovery, iteration order, output
PNG path construction, printed text, `dot` execution/error return, `open`
execution/ignored error, Cobra command registration, flags/help, initialization,
`pkg/structurizr`, `internal/adapter/process`, API/CLI surface, inventory,
scanner, baseline, audit repair, or any completed product behavior.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commits `3bd07e9` and `ce736a2`. Read the rolling
handover, this archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete `cmd/plugin_diagrams.go`, all diagrams and Cobra
callers/contracts, complete `pkg/structurizr`, the process adapter and all
complete doubles/callers, the import-aware scanner, API/CLI contracts, and the
T15 repair and baseline reproduction README.

# Three Moves

1. Start red with only focused private plugin-diagrams export recording
   contracts. Prove exact complete `process.Command`, one request, exact ignored
   injected error, continuation to file discovery, system dependency selection,
   safe zero dependency, non-empty recorded population, and absence of unrelated
   process requests. Use a private helper boundary so tests invoke no real
   `structurizr-cli`, `dot`, or `open` program.

2. Add only the private process dependency composition and route the one export
   request through `process.Execute`. Keep the request position after deletion
   and before discovery, ignored assignment semantics, command data, following
   operations, and all direct `dot`/`open` code unchanged. Regenerate exact Q1.3
   without changing the scanner or broadening the move to force a number.

3. Run focused cmd/process/structurizr and relevant caller tests; API/CLI and
   subprocess compatibility; launcher and Make contracts; complete tests,
   race, and vet; the 15-control audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/
   Q1.4/Q2.1/Q3.4 measurement; full clean audit; and empty-HOME count-2. Expect
   Q0.6 to hold at 26 guarded sites. The full audit may exit 1 for documented
   findings but never 2, and comparable ratchets must not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent export-process
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P3 result, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent authorized roadmap move, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make
the separate `docs: prepare next agent session` commit. Do not launch a real
successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P3 or P4-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
