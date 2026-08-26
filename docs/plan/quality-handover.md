# Quality Upgrade Handover

Generated: 2026-08-26T14:19:28+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.1 product implementation:
  `91422eb7fe0cd925a3a5512934ea08dda50c45c3`.
- Preceding P3.66 product implementation:
  `c627a3ec0ec8577f4bd2aef8cd28be7ab970b5e0`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

All earlier focused P3 implementation commits and their clean checkpoints
remain recorded in `docs/plan/quality-upgrade.md`. The separate operational
continuity implementation is `1b85711`, and the focused T15 repair is
`ce736a2`; neither changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active and P5-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T141928+0200-migrate-kibana-retry-sleep.md`.
The P4.1 Spring archive-clock archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`6a88ef2bfa7dc1aa912ccac6ef78ab2f1d0974813be9c512290087ff79982116`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 1

Only a narrow current-time adapter and Spring's private archive-path dependency
selection changed at product commit `91422eb`. New `clock.Now` returns the exact
injected `time.Time`. Its zero dependency returns the deterministic exact zero
time. `clock.System` selects a private system implementation whose single read
directly returns `time.Now()` without fallback, truncation, rounding, timezone
conversion, caching, monotonic rewriting, logging, cleanup, retry, or global
state.

Private `archivePathDependencies` retains the exact complete
`filesystem.System()` value and adds only the exact complete `clock.System()`
value. The existing working-directory request and its direct error still occur
before any clock read. After directory success, exactly one clock read supplies
the unchanged `Unix()` second to the unchanged `spring-%d.zip` name and slash
composition. Named returns, arbitrary directory bytes, pre-epoch and subsecond
conversion, `archivePath`, callers, archive creation, download, unzip,
deletion, logging, errors, and every public path are unchanged.

Four new clock tests and two new Spring archive-path tests bring the suite to
394 tests. The focused contracts prove deterministic safe zero behavior,
direct system-time identity within a bounded before/after observation, exact
arbitrary injected times, complete caller-owned dependency delivery and
preservation, directory-before-clock order, no clock read after a directory
error, one exact read after success, exact path and returns, non-empty
populations, public composition, and absence of another operation. They launch
no external program, touch no network, and write no repository fixture.

## T15 Apparatus Checkpoint

The P3.54 T15 repair remains unchanged. Separate verified execution and
structured replicas preserve historical test effects while measuring a
pristine exact `5635d50` checkout plus the declared inventory overlay. No
output path is excluded or deleted, and every effect stays in the retained
manifest. The complete 15-control proof passes and reproduces old scorecard
`d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`,
new scorecard
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`,
normalized stored raw body
`cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f`,
Q3.9, and all 228 numeric debt leaves. P4.1 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `91422eb` reports:

- Measured commit tree: `eedf149499ead042a28352a4c35b9d827f41a58f`.
- Structured scorecard SHA-256:
  `60bc1cf4994c63444d9b841f1dfd971d99283b21aad02cf3a859305c41fc62c4`.
- Absolute L0: 8 of 8.
- 394 test functions, zero skipped; 20 of 26 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 26 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 3 direct external sites outside five declared adapters of 27
  production effect sites; clock is valid and only server is missing.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 85 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 14 documented non-passing criteria, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The unchanged scanner recognizes the exact declared clock adapter and no
longer counts Spring's private `time.Now` selection, so exact Q1.3 improves
from 4/27 to 3/27. The authoritative audit output is
`target/quality-audit/scorecard.json`. Focused audit output is under
`/private/tmp/ply-p4-focused-audit`; its scorecard SHA-256 is
`1e38244f5ec922d540c85c91fa3fb9f66c68fd616b2e42b6343b1977d4732b03`.
Pinned tools are under `/private/tmp/ply-p4-tools` and
`/private/tmp/ply-p358-tools`. The disposable empty-HOME root is
`/private/tmp/ply-p4-hermetic.3xPZLb`.

## Next Objective

Continue P4 with the remaining direct clock site: extend the existing clock
adapter with an exact duration sleep and migrate only public Kibana `POST`'s
fixed `time.Sleep(15 * time.Second)` retry delay. Preserve the complete existing
`internalPOST` operation and exact request, first response selection, retry
condition, print, duration, retry count and ordering, first-result discard,
second error and response, callers, and public behavior.

Start red with focused clock-sleep and Kibana retry recording contracts. Prove
safe zero no-op behavior, exact arbitrary duration delivery, direct system
delegation, complete caller-owned dependency preservation, one first request,
no print/sleep/retry for non-empty hits, the exact print then one 15-second
sleep then one retry for empty hits, retry despite a first error, exact second
return values, at most two requests, non-empty populations, public composition,
and absence of another operation. Focused tests must perform no real sleep and
touch no network.

Then add only that sleep capability and the private complete Kibana retry
dependency/selector needed to record public `POST`. Leave `internalPOST`,
`internalPost`, HTTP behavior, Spring, web server construction/start/shutdown,
all completed P3/P4 flows, inventory, scanner, audit apparatus, mutation
harnesses, and the rest of P4-P8 unchanged.

## Verification Notes

- The red focused contracts failed because the production clock package/API and
  Spring dependency field did not exist.
- Focused clock/Spring and all relevant caller packages: pass.
- API/CLI and fresh subprocess compatibility: pass with reports outside the
  worktree; the exact pinned `apidiff` was rebuilt outside the repository after
  the first environment lacked it.
- Complete and uncached tests, race, vet, pinned lint, complete preflight, Make
  and production-script contracts, and all four host acceptance flows: pass.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- The first preflight attempt used the sandbox-blocked default lint cache; its
  isolated-cache rerun passed completely.
- The first complete `make test` attempt hit the established nested partial-
  raw-log signal-fixture flake after 25 launcher controls; the immediate
  unchanged complete rerun passed all 62 controls and the full target.
- During handoff verification, one direct and one nested second-generation
  launcher run hit that same signal-fixture race; the third unchanged
  `make test-agent-start` run passed all 62 controls.
- Empty-HOME count-2: pass under `/private/tmp/ply-p4-hermetic.3xPZLb`.
- The focused audit measured all seven requested criteria, exited 1 as
  expected, held Q0.6 at 26 guarded sites, improved exact Q1.3 to 3 of 27, and
  reported zero comparable regressions.
- Full clean audit: expected exit 1, 14 non-passing criteria, five improved,
  two held, zero regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P4 entry and checkpoint gate,
both design documents, `.quality/inventory`, the complete clock adapter and
Kibana POST/retry code/contracts/callers, the completed Spring archive-clock
contracts, both remaining server call sites, every adapter's source/tests and
complete dependency-preserving doubles, the import-aware scanner, API/CLI
contracts, and the T15 repair and baseline reproduction README before editing.
Confirm branch, HEAD, clean status, reciprocal archive links, launcher
`--check`, and exact product commit `91422eb`.

Make one focused Kibana retry-sleep implementation commit, then the normal
separate continuity commit. Stop before server work, another Kibana flow,
Spring or HTTP behavior changes, inventory or scanner changes, mutation
harnesses, later P4 work, P5-P8, publication, or distribution. Do not push,
merge, stash, revert, launch a successor, or remove the worktree.
