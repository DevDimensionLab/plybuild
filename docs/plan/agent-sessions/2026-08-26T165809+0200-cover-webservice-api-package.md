# Agent Session: Cover Webservice API Package

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T165809+0200-cover-webservice-api-package`
Created: `2026-08-26T16:58:09+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `bc9103a6d9a6f8378be9929bd9d3bc296262795be6a99461192db0616a716fc1`
Previous: [2026-08-26T163426+0200-cover-webservice-templates-package.md](2026-08-26T163426+0200-cover-webservice-templates-package.md)
Next: [2026-08-26T172816+0200-cover-logger-package.md](2026-08-26T172816+0200-cover-logger-package.md)
Outcome: Product `4da64d2` added five deterministic in-memory Generate and Upgrade GET/POST handler contracts at 100% API coverage; Q1.1 improved to 2/27 and the clean full audit exited 1 for 13 documented findings with zero ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one focused test-only coverage move for the untested package
`pkg/webservice/api`. Add deterministic characterization contracts for the
complete existing Generate and Upgrade GET/POST handler behavior without
changing production code, globals, channels, forms, templates, callers, HTML,
callback timing, responses, or observable semantics. Preserve every completed
checkpoint and produce zero comparable ratchet regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P4.6 product commit `8eeeb2a` added the complete
webservice-template characterization and improved Q1.1 from 4/27 to 3/27 while
exact Q1.3 held at zero of 27.

This mission authorizes only focused tests in `pkg/webservice/api` that
characterize the existing `GetGenerate`, `PostGenerate`, `GetUpgrade`, and
`PostUpgrade` handlers plus the exact public data/globals they already use. It
does not authorize changing any production file, extracting a seam, changing
`GenerateOptions`, `GOptions`, `CallbackChannel`, `CurrentProject`, API routes,
server behavior, form parsing, response text, goroutine delivery, config,
Spring, templates, callers, scanners, audit apparatus, inventory, mutation
harnesses, acceptance scripts, public APIs, HTML, logging, filesystem or
process behavior, another untested package, manual evidence, or P5-P8
implementation. If truthful tests expose preferred validation, quoting,
escaping, error handling, or callback synchronization, preserve the existing
behavior and record that product decision for later scope.

# Measurements At Start

Clean product commit `8eeeb2a` has 429 tests across 24 of 27 packages. Q0.6 has
27 guarded safe-writer sites, 22 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 3/27, Q1.2 is zero, exact Q1.3 is 0/27 with
all five declared adapters valid, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is
zero phrases across 90 Markdown files.

The authoritative clean full audit exits 1 for 13 documented non-passing
criteria, never 2, with L0 8/8, six improved, two held, zero regressed, zero
non-comparable ratchets, and zero dirty paths. Its scorecard SHA-256 is
`213df9ece4eb7c231ec3faa8f76e9180e9215a6ef0fc7d83039911940224c7e7`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
test implementation commit for this single package-coverage move, then make the
normal separate continuity-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve the exact exported
handler signatures, `GenerateOptions` shape, package globals and identities,
and the unbuffered production callback channel. Characterize Generate GET
rendering through its existing `text/template` caller path and Upgrade GET
rendering through its existing `html/template` caller path, including exact
status/body bytes for representative existing options. Characterize Generate
POST's existing ParseForm behavior, first-value scalar selection, exact project
field overwrites, ordered append of all `templates` and `dependencies` values,
one asynchronous callback delivery, and exact `OK` response. Characterize
Upgrade POST's one asynchronous callback and exact `OK` response. Do not add
validation, recovery, synchronization, status changes, redirects, escaping, or
error propagation.

Focused tests may use only in-memory `httptest` requests/recorders and a
test-owned callback channel whose original identity is restored with every
other mutated package global. They must open no socket, launch no external
program, perform no sleep or timed wait, mutate no logger or production state
after cleanup, leak no goroutine, and write no fixture. Do not use parallel
tests around package globals. Reject empty table-driven and recorded callback
populations before iterating. Leave templates, resources, sorting, server,
clock, Kibana, Spring, Maven, config, filesystem, process, HTTP adapters,
browser, plugin-diagrams, profile, shell, Unzip, and every completed contract
unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact product commit `8eeeb2a`. Read the rolling
handover, this archive, the complete P4 entry and checkpoint gate, both design
documents, `.quality/inventory`, all three complete files under
`pkg/webservice/api`, every API caller, the complete template sources and new
template contracts, relevant server/cmd/config/Spring types and tests,
representative in-memory handler, exact-response, global restoration,
callback-population, and explicit non-empty contracts, the import-aware
scanner, API/CLI contracts, and the T15 repair and baseline reproduction README.

# Three Moves

1. Add one focused `pkg/webservice/api` test file proving exact GET rendering,
   exact POST form-to-project mutations and slice appends, exact response bytes,
   one callback per POST, restored globals, no leaked goroutine, and non-empty
   populations. Do not change production handlers to make preferred validation,
   escaping, synchronization, or errors pass.
2. Run the focused package and relevant `cmd` and `pkg/webservice` caller tests,
   inspect package coverage, and confirm that only the new API test file changed.
   Regenerate focused Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 without changing
   scanner or inventory. Expect Q1.1 to improve from 3/27 to 2/27 and Q1.3 to
   hold at 0/27.
3. Run API/CLI and subprocess compatibility; launcher and Make contracts;
   complete tests, race, vet, and pinned lint; all four host acceptance flows;
   the 15-control audit meta-suite; full clean audit; and empty-HOME count-2.
   The full audit may exit 1 for documented findings but never 2, and no
   comparable ratchet may regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent webservice-API
coverage move or record an exact resumable state. Rewrite the rolling handover,
record the measured P4 result, answer this archive, create exactly one
reciprocally linked NEXT archive for the next coherent authorized roadmap move,
replace only the launcher's mutable regions, run launcher and handoff contracts,
and make the separate `docs: prepare next agent session` commit. Do not launch
a real successor, push, merge, publish, distribute, stash, revert, or remove
the worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
