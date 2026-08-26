# Agent Session: Migrate Plugin Diagrams Open

Status: NEXT
Session ID: `2026-08-26T080352+0200-migrate-plugin-diagrams-open`
Created: `2026-08-26T08:03:52+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `33599b9547262f76055dfa8d61c121e541051d6633a1ece6ad924a1b5eb40831`
Previous: [2026-08-26T072802+0200-migrate-plugin-diagrams-export.md](2026-08-26T072802+0200-migrate-plugin-diagrams-export.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the ignored macOS
`open` process request after each successful Graphviz conversion in
`cmd/plugin_diagrams.go` through the existing `internal/adapter/process`
boundary. Preserve its exact command name, argument, empty working directory
and streams, synchronous execution, one attempt, ignored error, and continued
iteration. Preserve the preceding export move, the Graphviz `dot` request,
Cobra behavior, every completed move, and zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.55 product commit `3972fc6` passed its
truthful complete checkpoint.

This mission authorizes only private plugin-diagrams open recording contracts,
the private dependency composition needed to select `process.System()` for
production, and replacement of this one request:

```go
_ = structurizr.Run(exec.Command("open", outputPngFile))
```

with the equivalent ignored `process.Execute` request. It does not authorize
the Graphviz `dot` request, another plugin-diagrams flow, the completed export
request, `structurizr.Run` or `RunWithOutputToFile`, the process adapter,
another adapter/caller, public API, inventory, scanner, audit apparatus,
clock/server, mutation harnesses, or P4-P8 implementation.

# Measurements At Start

Clean product commit `3972fc6` has 364 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, Q1.3 is 9/31 with clock
and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero phrases
across 73 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`ced6611a0990d55285007f9b9d19a990bbbe2bef9ee2fbe18084420fe9956fa6`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single open-process move and its recording
contracts, then make the normal separate continuity-only commit. Do not push,
merge, publish, distribute, remove the worktree, stash inherited changes,
revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve exact executable
`open`, exact single argument `outputPngFile`, empty `Dir`, nil
stdin/stdout/stderr, `Start: false`, one synchronous attempt after every
successful `dot` conversion, ignored process error, and unconditional
continuation to the next discovered file after that attempt. A recording test
must reject an empty request population.

Select `process.System()` only in the production composition. Pass the complete
private dependency value without fallback, retry, logging, another process
attempt, error wrapping, environment or working-directory behavior, output
capture, or cleanup. The zero dependency must remain a safe no-op and must not
launch an external tool. Focused tests must launch no real `structurizr-cli`,
`dot`, or `open` program, touch no network, and write no repository fixture.

Do not change the mandatory workspace lookup, `.structurizr/` deletion,
completed export `process.Execute` request, dot-file discovery, iteration order,
output PNG construction, printed text, direct Graphviz execution and returned
error, Cobra registration/flags/help, initialization, `pkg/structurizr`,
`internal/adapter/process`, API/CLI surface, inventory, scanner, baseline, audit
repair, or any completed product behavior.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `3972fc6`. Read the rolling handover, this
archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete plugin-diagrams and structurizr code/tests and
Cobra callers, the process adapter and every complete process double/caller,
the import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README.

# Three Moves

1. Start red with only focused private plugin-diagrams open recording contracts.
   Prove the exact complete `process.Command`, one request per helper invocation,
   exact ignored injected error semantics, placement only after a successful
   direct `dot` conversion, continued iteration, complete system dependency
   selection, safe zero dependency, non-empty population, and absence of export,
   `dot`, or another unrelated process request. Use a private helper boundary so
   focused tests invoke no external program.
2. Add only the private open-process dependency composition and route the one
   ignored `open` request through `process.Execute`. Keep it after the successful
   `dot` request and before the next iteration, and leave the export request and
   all direct `dot` code unchanged. Regenerate exact Q1.3 without changing the
   scanner or broadening the move to force a number.
3. Run focused cmd/process/structurizr and relevant caller tests; API/CLI and
   subprocess compatibility; launcher and Make contracts; complete tests, race,
   and vet; the 15-control audit meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/
   Q2.1/Q3.4 measurement; full clean audit; and empty-HOME count-2. Expect Q0.6
   to hold at 26 guarded sites. The full audit may exit 1 for documented
   findings but never 2, and comparable ratchets must not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent open-process move
or record an exact resumable state. Rewrite the rolling handover, record the
measured P3 result, answer this archive, create exactly one reciprocally linked
NEXT archive for the next coherent authorized roadmap move, replace only the
launcher's mutable regions, run launcher and handoff contracts, and make the
separate `docs: prepare next agent session` commit. Do not launch a real
successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P3 or P4-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
