# Agent Session: Cover Webservice Templates Package

Status: ANSWERED - HISTORY
Session ID: `2026-08-26T163426+0200-cover-webservice-templates-package`
Created: `2026-08-26T16:34:26+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `7e55f3c10405261604beccd8e9a33a4c64e7e8fbae2cd46c7baf4b66d958bc79`
Previous: [2026-08-26T160625+0200-cover-resources-package.md](2026-08-26T160625+0200-cover-resources-package.md)
Next: [2026-08-26T165809+0200-cover-webservice-api-package.md](2026-08-26T165809+0200-cover-webservice-api-package.md)
Outcome: Product `8eeeb2a` added three deterministic complete template byte/composition/token/parser contracts; Q1.1 improved to 3/27 and the clean full audit exited 1 for 13 documented findings with zero ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Continue P4 with one focused test-only coverage move for the untested package
`pkg/webservice/templates`. Add deterministic characterization contracts for
the complete existing Generate and Upgrade template bytes and caller parser
compatibility without changing production code, exported constants, callers,
HTML, template actions, network references, or observable semantics. Preserve
every completed checkpoint and produce zero comparable ratchet regressions.

# Authorized Roadmap

P3 is complete, P4 is active, and P5-P8 remain queued in the machine-readable
block in `docs/plan/quality-upgrade.md`. The launcher must remain NEXT until all
authorized checkpoints finish. P4.5 product commit `64c6189` added the complete
resources characterization and improved Q1.1 from 5/27 to 4/27 while exact
Q1.3 held at zero of 27.

This mission authorizes only focused tests in `pkg/webservice/templates` that
characterize the existing `Generate` and `Upgrade` constants and their private
header, body, and footer composition. It does not authorize changing any file
under `pkg/webservice/templates` except one new test file, changing
`pkg/webservice/api`, server behavior, handlers, forms, config, Spring,
resources, callers, scanners, audit apparatus, inventory, mutation harnesses,
acceptance scripts, public APIs, HTML, external CDN strings, logging, filesystem
or process behavior, another untested package, manual evidence, or P5-P8
implementation. If truthful tests expose a preferred HTML or escaping change,
preserve the existing bytes and record that product decision for later scope.

# Measurements At Start

Clean product commit `64c6189` has 426 tests across 23 of 27 packages. Q0.6 has
27 guarded safe-writer sites, 22 write and 5 copy, zero skipped tests, and zero
unsafe direct test writes. Q1.1 is 4/27, Q1.2 is zero, exact Q1.3 is 0/27 with
all five declared adapters valid, Q1.4 is 7/8, exact Q2.1 is 0/8, and Q3.4 is
zero phrases across 89 Markdown files.

The authoritative clean full audit exits 1 for 13 documented non-passing
criteria, never 2, with L0 8/8, six improved, two held, zero regressed, zero
non-comparable ratchets, and zero dirty paths. Its scorecard SHA-256 is
`2bbbc911efefe7427cd3413b73b0a96583a1111460e57f5cdfb9b145b6818a95`.
The repaired T15 proof passes all 15 controls and reproduces the exact stored
debt and instrument identities.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
test implementation commit for this single package-coverage move, then make the
normal separate continuity-only commit. Do not push, merge, publish,
distribute, remove the worktree, stash inherited changes, revert user work, or
run destructive Git commands.

Keep Go 1.18 and `/bin/bash` 3.2 compatibility. Preserve `Generate` and
`Upgrade` as exported untyped string constants with their exact names and
complete byte values. Characterize exact `header + generate + footer` and
`header + upgrade + footer` composition, leading and trailing bytes, shared
header/footer identity, exact complete payload digests or equivalent exhaustive
byte contracts, the existing form actions, input names, template actions, CDN
references, and successful parsing through the exact `text/template` engine for
Generate and `html/template` engine for Upgrade selected by their callers.
Do not improve quoting, escaping, markup, formatting, URLs, or parser choice.

Focused tests must make no request, open no socket, launch no external program,
perform no wait, mutate no global API options/channel/logger or production
state, and write no fixture. Reject empty table-driven populations before
iterating. Leave resources, sorting, server, clock, Kibana, Spring, Maven,
config, filesystem, process, HTTP, browser, plugin-diagrams, profile, shell,
Unzip, and every completed contract unchanged.

# Required Reading

Before editing, confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact product commit `64c6189`. Read the rolling
handover, this archive, the complete P4 entry and checkpoint gate, both design
documents, `.quality/inventory`, all four complete files under
`pkg/webservice/templates`, every template caller under `pkg/webservice/api`,
the complete relevant server and API-facing tests, representative exact-byte,
parser, and explicit non-empty-population contracts, the import-aware scanner,
API/CLI contracts, and the T15 repair and baseline reproduction README.

# Three Moves

1. Add focused `pkg/webservice/templates` contracts that prove complete exact
   public template bytes, exact private-fragment composition and boundaries,
   stable actions/names/template expressions/CDN references, successful caller-
   selected parser compatibility, and non-empty test populations. Do not change
   production bytes to make preferred markup or escaping pass.
2. Run the focused package and relevant `pkg/webservice/api` and
   `pkg/webservice` caller tests, inspect package coverage, and confirm that only
   the new templates test file changed. Regenerate focused
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1/Q3.4 without changing scanner or inventory.
   Expect Q1.1 to improve from 4/27 to 3/27 and Q1.3 to hold at 0/27.
3. Run API/CLI and subprocess compatibility; launcher and Make contracts;
   complete tests, race, vet, and pinned lint; all four host acceptance flows;
   the 15-control audit meta-suite; full clean audit; and empty-HOME count-2.
   The full audit may exit 1 for documented findings but never 2, and no
   comparable ratchet may regress.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent webservice-
templates coverage move or record an exact resumable state. Rewrite the rolling
handover, record the measured P4 result, answer this archive, create exactly one
reciprocally linked NEXT archive for the next coherent authorized roadmap move,
replace only the launcher's mutable regions, run launcher and handoff contracts,
and make the separate `docs: prepare next agent session` commit. Do not launch
a real successor, push, merge, publish, distribute, stash, revert, or remove
the worktree. COMPLETE remains invalid while P4 or P5-P8 is unfinished.
<!-- CODEX_SESSION_PROMPT_END -->
