# Agent Session: Migrate Web Server Effects

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T145252+0200-migrate-web-server-effects`
Created: `2026-08-26T14:52:52+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `58f61de9e370f4904361e44146cc016d27e54b0f27752631faada84d87fdcaac`
Previous: [2026-08-26T141928+0200-migrate-kibana-retry-sleep.md](2026-08-26T141928+0200-migrate-kibana-retry-sleep.md)
Next: [2026-08-26T153936+0200-cover-sorting-package.md](2026-08-26T153936+0200-cover-sorting-package.md)
Outcome: Product `9601d29` added the declared server boundary and exact webservice start/stop recording contracts; exact Q1.3 measured 0/27 and the clean full audit exited 1 for 13 documented findings with zero ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one focused server move: introduce the declared server adapter
and migrate only webservice `StartWebServer`'s direct `ListenAndServe` and
`StopWebServer`'s direct `Shutdown` calls to that boundary. Preserve exact
server construction, identity, address, handlers, registration order, start,
shutdown, context, timeout, cancellation, logging, ignored error, caller, and
public behavior; every completed checkpoint; and zero comparable ratchet
regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P4.2 product commit `4c0d97a` passed its
truthful complete clock-sleep and Kibana retry checkpoint.

This mission authorizes only focused server-operation and webservice start/stop
recording contracts; one narrow `internal/adapter/server` boundary for the
exact existing `ListenAndServe` and `Shutdown` operations on the exact supplied
server; deterministic safe-zero no-ops; the complete private caller-owned
webservice dependencies and exact production selectors needed to record public
`StartWebServer` and `StopWebServer`; and replacing only their two direct server
method calls with those adapter operations.

It does not authorize changing the package-level server construction or
identity, `:7999` address, port, scheme, host, endpoint URI, handler paths or
callbacks, registration order, `StartWebServer` or `StopWebServer` signatures,
five-second timeout, background parent, deferred cancellation, logging rule,
ignored shutdown error, public APIs, browser or blocking flows, another
webservice flow, clock, process, HTTP-client, or filesystem adapters, Kibana,
Spring, Maven, Unzip, inventory, scanner, audit apparatus, mutation harnesses,
the loopback-address policy, or later P4 and P5-P8 implementation.

# Measurements At Start

Clean product commit `4c0d97a` has 401 tests across 20 of 26 packages. Q0.6 has
26 guarded safe-writer sites, 21 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 6/26, Q1.2 is zero, exact Q1.3 is 2/27 and
both remaining violations are the server calls in `pkg/webservice/api.go`,
Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is zero phrases across 86 Markdown
files.

The authoritative clean full audit exits 1 for 14 documented non-passing
criteria, never 2, with L0 8/8, five improved, two held, zero regressed, one
non-comparable ratchet, and zero dirty paths. Its scorecard SHA-256 is
`be0142a87d2ccb5ff1a3d9414ba9d376843f7ab8b92e5ba0f75f512913e0b0db`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for this single web-server effect move and its recording
contracts, then make the normal separate continuity-only commit. Do not push,
merge, publish, distribute, remove the worktree, stash inherited changes,
revert user work, or run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Each server adapter production
operation must call the corresponding method on the exact supplied server once
and return its exact error, without replacement, wrapping, fallback, recovery,
normalization, logging, cleanup, retry, address rewriting, handler mutation, or
global state. Safe-zero dependencies must perform no network operation and
return a deterministic nil error. Preserve complete caller-owned dependency
fields in every affected recording double and reject empty recorded
populations.

Preserve the exact package-level `*http.Server` value and `:7999` address.
`StartWebServer` must register the same four path/callback pairs in the same
order before one listen attempt; it must print the exact returned non-nil error
through the same logger rule and print nothing for nil. `StopWebServer` must
create one context from `context.Background()` with exactly `5*time.Second`,
defer its cancellation, pass that exact context to one shutdown attempt, and
continue to ignore the exact shutdown error. Preserve every caller and public
return path.

Focused tests must open no socket, make no request, launch no external program,
perform no real wait, and write no repository fixture. They must not depend on
ambient ports or mutate the production server identity. Leave all completed
clock, Kibana, Spring, filesystem, process, HTTP, Maven, browser,
plugin-diagrams, profile, shell, and Unzip behavior unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact commit `4c0d97a`. Read the rolling handover, this
archive, the P4 entry and checkpoint gate, both design documents,
`.quality/inventory`, complete `pkg/webservice/api.go` and `init.go`, every
`StartWebServer`/`StopWebServer` caller, all webservice contracts, the exact
server construction and four handlers, the completed clock and Kibana move,
every complete adapter source/test and dependency-preserving recording double,
the import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README.

# Three Moves

1. Start red with focused server-adapter and webservice recording contracts.
   Prove both safe-zero operations, exact server identity, exact context and
   error delivery, direct system delegation, complete dependency delivery and
   preservation, four exact handler registrations in order before listen,
   one listen attempt, nil/non-nil logging behavior, exact background-derived
   five-second context and deferred cancel, one shutdown attempt, ignored
   shutdown error, non-empty populations, public composition, and absence of
   another operation.
2. Add only the authorized server adapter and the complete private webservice
   dependency-taking compositions and production selectors needed to record
   both public functions; then replace only the direct `ListenAndServe` and
   `Shutdown` calls. Preserve the server value, address, handlers, timeout,
   logging, ignored error, signatures, callers, and every return path.
   Regenerate exact Q1.3 without changing the scanner or broadening the move to
   force a number.
3. Run focused server/webservice and relevant caller tests; API/CLI and
   subprocess compatibility; launcher and Make contracts; complete tests,
   race, and vet; all four host acceptance flows; the 15-control audit
   meta-suite; focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 measurement; full
   clean audit; and empty-HOME count-2. Expect Q0.6 to hold at 26 guarded sites
   and exact Q1.3 to improve from 2/27 to 0/27. The full audit may exit 1 for
   documented findings but never 2, and comparable ratchets must not regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent web-server move
or record an exact resumable state. Rewrite the rolling handover, record the
measured P4 result, answer this archive, create exactly one reciprocally linked
NEXT archive for the next coherent authorized roadmap move, replace only the
launcher's mutable regions, run launcher and handoff contracts, and make the
separate `docs: prepare next agent session` commit. Do not launch a real
successor, push, merge, publish, distribute, stash, revert, or remove the
worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
