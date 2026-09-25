# Agent Session: Migrate Spring Archive Clock

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T134455+0200-migrate-spring-archive-clock`
Created: `2026-08-26T13:44:55+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d7555a234bd2a57b20e3074617a6cf2b5bf36e644e0b1e7881500614dcd52662`
Previous: [2026-08-26T130843+0200-narrow-maven-standard-output-provenance.md](2026-08-26T130843+0200-narrow-maven-standard-output-provenance.md)
Next: [2026-08-26T141928+0200-migrate-kibana-retry-sleep.md](2026-08-26T141928+0200-migrate-kibana-retry-sleep.md)
Outcome: Product commit `91422eb` added the exact safe-zero/system clock adapter and migrated only Spring's private archive-path timestamp selection. The clean checkpoint has 394 tests, improves exact Q1.3 from 4/27 to 3/27, and has zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Begin P4 with one focused clock move: introduce the declared clock adapter and
migrate only Spring's private archive-path timestamp selection from direct
`time.Now().Unix()` to that boundary. Preserve exact directory, timestamp,
path, error, dependency, caller, and public behavior; completed P3; every prior
checkpoint; and zero comparable ratchet regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P3.66 product commit `c627a3e` passed its
truthful complete checkpoint.

This mission authorizes only focused clock-adapter and Spring archive-path
recording contracts; one narrow `internal/adapter/clock` dependency that
returns the exact existing system time; one exact safe-zero clock read; adding
that dependency to Spring's private archive-path dependency and production
selector; and replacing only the one direct `time.Now()` selection in private
`archivePathWithDependencies` with that adapter read.

It does not authorize changing another Spring flow or public API; Kibana sleep;
web server construction, start, shutdown, address, timeout, or handlers;
process, HTTP, or filesystem adapters; Maven; Unzip; browser; plugin diagrams;
profile; shell Run or Git; inventory; scanner; audit apparatus; mutation
harnesses; or later P4 and P5-P8 implementation.

# Measurements At Start

Clean product commit `c627a3e` has 388 tests across 19 of 25 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/25, Q1.2 is zero, exact Q1.3 is 4/27 with
clock and server absent, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero
phrases across 84 Markdown files.

The authoritative clean full audit exits 1 for 15 documented findings, never
2, with L0 8/8, five improved, two held, zero regressed, one non-comparable
ratchet, and zero dirty paths. Its scorecard SHA-256 is
`1023a95ee7e4d7a10e95ae647821892dae71b8815f35062ee5e4b8cbd9c207d0`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single Spring archive-clock move and its
recording contracts, then make the normal separate continuity-only commit. Do
not push, merge, publish, distribute, remove the worktree, stash inherited
changes, revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. The clock adapter must return
the exact `time.Now` value without fallback, truncation, rounding, timezone
conversion, caching, monotonic rewriting, logging, cleanup, retry, or global
state. Its zero dependency must be safe and deterministic. Preserve the
existing complete caller-owned Spring archive-path dependency value; the
private production selector must retain the exact complete filesystem system
dependency and add only the exact system clock.

Preserve the exact working-directory request and its direct error before any
clock read, one clock read only after directory success, exact `Unix()` second,
the `spring-%d.zip` name and slash composition, arbitrary directory bytes,
pre-epoch and subsecond injected values, named returns, caller behavior, and
every public path. Do not alter archive creation, download, unzip, deletion,
logging, or error behavior.

Focused tests must launch no external program, touch no network, and write no
repository fixture. Preserve complete caller-owned dependency fields in every
affected recording double. Leave all completed filesystem, process, HTTP,
Maven, browser, plugin-diagrams, profile, shell, and Unzip behavior unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `c627a3e`. Read the rolling handover, this
archive, the P4 entry and checkpoint gate, both design documents,
`.quality/inventory`, complete Spring archive-path code/contracts and every
caller, both remaining direct clock sites, every complete adapter source/test
and dependency-preserving recording double, the P3.66 Maven contracts, the
import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README.

# Three Moves

1. Start red with focused clock-adapter and Spring archive-path recording
   contracts. Prove safe zero behavior, exact system-time identity within a
   bounded before/after observation, complete dependency delivery and
   preservation, working-directory-before-clock order, no clock read after a
   directory error, one exact read after success, arbitrary injected times,
   exact path and returns, non-empty population, public composition, and
   absence of another operation.
2. Add only the authorized narrow clock adapter, the private Spring clock
   dependency and exact production selection, then replace only private
   `archivePathWithDependencies`' direct `time.Now()` with the adapter read.
   Preserve the filesystem dependency, `archivePath`, every caller, and the
   complete return path. Regenerate exact Q1.3 without changing the scanner or
   broadening the move to force a number.
3. Run focused clock/Spring and relevant caller tests; API/CLI and subprocess
   compatibility; launcher and Make contracts; complete tests, race, and vet;
   all four host acceptance flows; the 15-control audit meta-suite; focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full clean audit; and
   empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites. The full audit
   may exit 1 for documented findings but never 2, and comparable ratchets must
   not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Spring archive-
clock move or record an exact resumable state. Rewrite the rolling handover,
record the measured P4 result, answer this archive, create exactly one
reciprocally linked NEXT archive for the next coherent authorized roadmap move,
replace only the launcher's mutable regions, run launcher and handoff contracts,
and make the separate `docs: prepare next agent session` commit. Do not launch a
real successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
