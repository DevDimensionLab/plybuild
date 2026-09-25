# Agent Session: Migrate Kibana Retry Sleep

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T141928+0200-migrate-kibana-retry-sleep`
Created: `2026-08-26T14:19:28+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `6a88ef2bfa7dc1aa912ccac6ef78ab2f1d0974813be9c512290087ff79982116`
Previous: [2026-08-26T134455+0200-migrate-spring-archive-clock.md](2026-08-26T134455+0200-migrate-spring-archive-clock.md)
Next: [2026-08-26T145252+0200-migrate-web-server-effects.md](2026-08-26T145252+0200-migrate-web-server-effects.md)
Outcome: Product commit `4c0d97a` added exact safe-zero/system duration sleep to the complete clock adapter and migrated only public Kibana `POST`'s fixed retry sleep. The clean checkpoint has 401 tests, improves exact Q1.3 from 3/27 to 2/27, and has zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one focused clock move: extend the declared clock adapter with
the exact existing duration sleep and migrate only Kibana public `POST`'s fixed
retry delay from direct `time.Sleep(15 * time.Second)` to that boundary.
Preserve exact request, first-result selection, print, duration, retry, error,
response, dependency, caller, and public behavior; every completed checkpoint;
and zero comparable ratchet regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P4.1 product commit `91422eb` passed its
truthful complete Spring archive-clock checkpoint.

This mission authorizes only focused clock-sleep and Kibana `POST` retry
recording contracts; one narrow clock sleep operation that delegates the exact
requested duration to the existing system sleep; one exact safe-zero no-op;
one complete private caller-owned Kibana retry dependency and exact production
selector retaining the existing `internalPOST` operation and adding only the
exact system clock; one private dependency-taking composition for public
`POST`; and replacing only its direct `time.Sleep` call with the adapter sleep.

It does not authorize changing `internalPOST`, `internalPost`, HTTP request
execution, parsing, response bodies, public APIs, another Kibana flow, Spring,
web server construction, start, shutdown, address, timeout, or handlers;
process, HTTP, or filesystem adapters; Maven; Unzip; browser; plugin diagrams;
profile; shell; inventory; scanner; audit apparatus; mutation harnesses; or
later P4 and P5-P8 implementation.

# Measurements At Start

Clean product commit `91422eb` has 394 tests across 20 of 26 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/26, Q1.2 is zero, exact Q1.3 is 3/27 with
the clock adapter present and only server absent, Q1.4 is 7/8, exact Q2.1 is
0/8, and Q3.4 is zero phrases across 85 Markdown files.

The authoritative clean full audit exits 1 for 14 documented non-passing
criteria, never 2, with L0 8/8, five improved, two held, zero regressed, one
non-comparable ratchet, and zero dirty paths. Its scorecard SHA-256 is
`60bc1cf4994c63444d9b841f1dfd971d99283b21aad02cf3a859305c41fc62c4`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single Kibana retry-sleep move and its
recording contracts, then make the normal separate continuity-only commit. Do
not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. The clock adapter's production
sleep must delegate the exact supplied `time.Duration` to `time.Sleep` once,
without fallback, normalization, clamping, rounding, conversion, caching,
logging, cleanup, retry, or global state. Its zero dependency must be a safe,
deterministic no-op. Preserve the complete current-time dependency and exact
`Now` behavior, Spring's complete clock value, and every existing clock
contract while adding only the sleep capability.

Preserve one first `internalPOST` attempt with the exact request. When its
returned hit list is non-empty, preserve no print, no sleep, no retry, and the
exact first error and response. When it is empty, preserve the exact
`sleep and retry` print, one exact 15-second sleep after the first result, one
second `internalPOST` attempt with the same complete request, and the exact
second error and response; discard the first result exactly as today and never
make a third attempt. This retry selection remains based only on hit-list
length, including when the first error is non-nil.

Focused tests must perform no real sleep, launch no external program, touch no
network, and write no repository fixture. Preserve complete caller-owned
dependency fields in every affected recording double and reject empty recorded
populations. Leave all completed clock-Now, Spring, filesystem, process, HTTP,
Maven, browser, plugin-diagrams, profile, shell, and Unzip behavior unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `91422eb`. Read the rolling handover, this
archive, the P4 entry and checkpoint gate, both design documents,
`.quality/inventory`, the complete clock adapter source/contracts, complete
Kibana `POST`/`internalPOST`/`internalPost` code/contracts and every caller,
the completed Spring archive-clock code/contracts, both remaining server call
sites, every complete adapter source/test and dependency-preserving recording
double, the import-aware scanner, API/CLI contracts, and the T15 repair and
baseline reproduction README.

# Three Moves

1. Start red with focused clock-sleep and Kibana retry recording contracts.
   Prove safe zero behavior, exact arbitrary duration delivery, exact direct
   system-sleep delegation, complete dependency delivery and preservation,
   first request before selection, no sleep or retry for non-empty hits, exact
   print/sleep/retry order for empty hits, retry despite a first error, exact
   second return values, at most two requests, non-empty populations, public
   composition, and absence of another operation.
2. Add only the authorized sleep capability to the existing clock adapter and
   the complete private Kibana retry dependency/production selector needed to
   record public `POST`; then replace only its direct `time.Sleep` with the
   adapter sleep. Preserve `internalPOST`, `internalPost`, the request and
   response types, all callers, the exact 15-second constant, and every return
   path. Regenerate exact Q1.3 without changing the scanner or broadening the
   move to force a number.
3. Run focused clock/Kibana and relevant caller tests; API/CLI and subprocess
   compatibility; launcher and Make contracts; complete tests, race, and vet;
   all four host acceptance flows; the 15-control audit meta-suite; focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full clean audit; and
   empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites and exact Q1.3
   to improve from 3/27 to 2/27. The full audit may exit 1 for documented
   findings but never 2, and comparable ratchets must not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Kibana retry-
sleep move or record an exact resumable state. Rewrite the rolling handover,
record the measured P4 result, answer this archive, create exactly one
reciprocally linked NEXT archive for the next coherent authorized roadmap move,
replace only the launcher's mutable regions, run launcher and handoff contracts,
and make the separate `docs: prepare next agent session` commit. Do not launch a
real successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
