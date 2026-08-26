# Agent Session: Migrate Plugin Diagrams Dot

Status: NEXT
Session ID: `2026-08-26T084000+0200-migrate-plugin-diagrams-dot`
Created: `2026-08-26T08:40:00+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5c260abb2033cc27c5c7ee8ee57d6f485cc4e44ed753bf926b408d99d3a5ba05`
Previous: [2026-08-26T080352+0200-migrate-plugin-diagrams-open.md](2026-08-26T080352+0200-migrate-plugin-diagrams-open.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete one focused P3 production-effect move: route only the direct Graphviz
`dot` conversion request in `cmd/plugin_diagrams.go` through the existing
`internal/adapter/process` boundary while preserving the conversion's captured
stdout, discarded stderr, exact output write, error behavior, and following
ignored macOS open attempt. Preserve the completed export and open moves, Cobra
behavior, every completed move, and zero comparable ratchet regressions.

# Authorized Roadmap

P3 remains active, P4-P8 remain queued in the machine-readable block in
`docs/plan/quality-upgrade.md`, and the launcher must remain NEXT until all
authorized checkpoints finish. P3.56 product commit `6f4bc37` passed its
truthful complete checkpoint.

This mission authorizes only private plugin-diagrams Graphviz recording
contracts, the private process/filesystem dependency composition required to
preserve the existing conversion output behavior, and replacement of this one
call:

```go
err = structurizr.RunWithOutputToFile(exec.Command("dot", file, "-Tpng"), outputPngFile)
```

with an equivalent private helper using `process.Execute` for the direct `dot`
request and the existing filesystem adapter for the exact successful stdout
write. It does not authorize the completed export or `open` requests, another
plugin-diagrams flow, `structurizr.Run` or `RunWithOutputToFile`, any adapter,
another caller, public API, inventory, scanner, audit apparatus, clock/server,
mutation harnesses, or P4-P8 implementation.

# Measurements At Start

Clean product commit `6f4bc37` has 369 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, Q1.3 is 8/30 with clock
and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero phrases
across 74 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`efaaee797a34d182d03d6830d0af4c342b7e744180775942119a77c7e0f88090`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single Graphviz-process move and its recording
contracts, then make the normal separate continuity-only commit. Do not push,
merge, publish, distribute, remove the worktree, stash inherited changes,
revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve exact executable
`dot`, ordered arguments `file`, `-Tpng`, empty `Dir`, nil stdin, distinct
stdout/stderr buffers, `Start: false`, and one synchronous attempt for each
discovered file. Preserve the exact direct process error and suppress the output
write and following open attempt on that error. On process success, preserve the
exact stdout bytes, `outputPngFile`, mode `0644`, one write attempt, ignored
write error, discarded stderr bytes, and following one ignored open attempt.
Preserve iteration after a successful conversion/open sequence. Recording tests
must reject empty process and write populations.

Select `process.System()` and `filesystem.System()` only in the production
composition. Pass the complete private dependency value without fallback,
retry, logging, another process or write attempt, error wrapping, environment or
working-directory behavior, output transformation, cleanup, or surfaced write
error. The zero dependency must remain safe and must not launch a tool or write
a file. Focused tests must launch no real `structurizr-cli`, `dot`, or `open`
program, touch no network, and write no repository fixture.

Do not change mandatory workspace lookup, `.structurizr/` deletion, completed
export `process.Execute` request, dot-file discovery, iteration order, output
PNG construction, printed text, completed open `process.Execute` request,
Cobra registration/flags/help, initialization, `pkg/structurizr`, either
adapter, API/CLI surface, inventory, scanner, baseline, audit repair, or any
completed product behavior.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `6f4bc37`. Read the rolling handover, this
archive, the P3 tail and checkpoint gate, both design documents,
`.quality/inventory`, complete plugin-diagrams and structurizr code/tests and
Cobra callers, the process and filesystem adapters and every complete relevant
double/caller, the import-aware scanner, API/CLI contracts, and the T15 repair
and baseline reproduction README.

# Three Moves

1. Start red with only focused private plugin-diagrams Graphviz recording
   contracts. Prove the exact complete process command and stream identities,
   one request, exact process error precedence, no write/open after process
   failure, exact successful output bytes/path/mode, ignored exact write error,
   following open placement, complete system dependency selection, safe zero
   dependency, non-empty populations, and absence of export, open, or another
   unrelated process/write request. Use a private helper boundary so focused
   tests invoke no external program.
2. Add only the private Graphviz process/filesystem dependency composition and
   route the one direct `dot` request through `process.Execute`. Preserve the
   existing successful output write through `filesystem.WriteFile`, including
   its ignored result, and keep the following completed open helper unchanged.
   Leave `pkg/structurizr` and its exported functions untouched. Regenerate exact
   Q1.3 without changing the scanner or broadening the move to force a number.
3. Run focused cmd/process/filesystem/structurizr and relevant caller tests;
   API/CLI and subprocess compatibility; launcher and Make contracts; complete
   tests, race, and vet; the 15-control audit meta-suite; focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full clean audit; and
   empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites. The full audit
   may exit 1 for documented findings but never 2, and comparable ratchets must
   not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Graphviz-process
move or record an exact resumable state. Rewrite the rolling handover, record
the measured P3 result, answer this archive, create exactly one reciprocally
linked NEXT archive for the next coherent authorized roadmap move, replace only
the launcher's mutable regions, run launcher and handoff contracts, and make the
separate `docs: prepare next agent session` commit. Do not launch a real
successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P3 or P4-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
