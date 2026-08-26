# Quality Upgrade Handover

Generated: 2026-08-26T14:52:52+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.2 product implementation:
  `4c0d97ad2ac1a6e5faee6ec6f240eeed43c4b434`.
- Preceding P4.1 product implementation:
  `91422eb7fe0cd925a3a5512934ea08dda50c45c3`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

All earlier focused P3 and P4 implementation commits and their clean
checkpoints remain recorded in `docs/plan/quality-upgrade.md`. The separate
operational continuity implementation is `1b85711`, and the focused T15 repair
is `ce736a2`; neither changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active and P5-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T145252+0200-migrate-web-server-effects.md`.
The P4.2 Kibana retry-sleep archive is answered history and links reciprocally
to that archive. The graph has exactly one NEXT tail. The archived prompt
SHA-256 is
`58f61de9e370f4904361e44146cc016d27e54b0f27752631faada84d87fdcaac`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 2

Only exact duration sleep in the existing clock adapter and public Kibana
`POST`'s fixed retry delay changed at product commit `4c0d97a`. New
`clock.Sleep` passes the exact injected `time.Duration` once. Its zero
dependency is a deterministic no-op. `clock.System` preserves the complete
current-time dependency and adds a sleeper implemented by the same private
system clock; its production sleep delegates the exact duration directly to
`time.Sleep` once without fallback, normalization, clamping, rounding,
conversion, caching, logging, cleanup, retry, or global state. Existing
`clock.Now`, Spring's complete clock value, and every earlier clock contract are
unchanged.

Private `postDependencies` owns the complete existing `internalPOST` operation
and complete clock value. Its production selector chooses exactly
`internalPOST` plus `clock.System()`, and public `POST` composes only that
selector with the private dependency-taking operation. The first request still
happens before selection. A non-empty first hit list returns the exact first
error and response with no print, sleep, or retry. An empty list, even with a
non-nil first error, prints `sleep and retry`, sleeps once for exactly 15
seconds, retries the same complete request once, and returns the exact second
error and response. It discards the first result and never makes a third
request. `internalPOST`, `internalPost`, HTTP behavior, parsing, responses,
callers, and public APIs are unchanged.

Seven net new focused tests bring the suite to 401 tests. The contracts prove
safe zero behavior, arbitrary exact duration delivery, direct system-sleep
delegation, complete caller-owned dependency delivery and preservation,
first-request ordering, both selection branches, retry despite the first
error, exact print/sleep/retry order and duration, exact return identities, at
most two requests, non-empty populations, public composition, and absence of
another operation. They perform no real sleep, launch no external program,
touch no network, and write no repository fixture.

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
Q3.9, and all 228 numeric debt leaves. P4.2 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `4c0d97a` reports:

- Measured commit tree: `2bff9b81b137130ba026772a37526bbb24595c32`.
- Structured scorecard SHA-256:
  `be0142a87d2ccb5ff1a3d9414ba9d376843f7ab8b92e5ba0f75f512913e0b0db`.
- Absolute L0: 8 of 8.
- 401 test functions, zero skipped; 20 of 26 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 26 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 2 direct external sites outside five declared adapters of 27
  production effect sites; both are in `pkg/webservice/api.go`, and only the
  server adapter is missing.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 86 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 14 documented non-passing criteria, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The unchanged scanner recognizes the clock adapter sleep and no longer counts
Kibana's retry delay, so exact Q1.3 improves from 3/27 to 2/27. The
authoritative audit output is `target/quality-audit/scorecard.json`. Focused
audit output is under `/private/tmp/ply-p4-focused-audit.64QbvU`; its scorecard
SHA-256 is
`10cf5c8891a57a0ac883b9b685afd8ad719d14adffadf5dcfa5ff0db6631e63a`.
Pinned tools are under `/private/tmp/ply-p4-tools` and
`/private/tmp/ply-p358-tools`. The disposable empty-HOME root is
`/private/tmp/ply-p4-hermetic.zAFimZ`.

## Next Objective

Continue P4 with the remaining direct-effect group: add the declared
`internal/adapter/server` boundary and migrate only webservice
`StartWebServer`'s `ListenAndServe` and `StopWebServer`'s `Shutdown` calls.
Preserve the exact package-level server pointer, `:7999` address, four handler
paths and callbacks in order, one listen attempt, nil/non-nil error logging,
background-derived five-second timeout context, deferred cancellation, one
shutdown attempt, ignored shutdown error, callers, and public APIs.

Start red with focused server-adapter and webservice recording contracts. Prove
safe zero no-op behavior without a socket, exact receiver/context/error
delivery, direct production delegation, complete dependency preservation,
handler-before-listen ordering, exact logging, timeout/cancel behavior, ignored
shutdown error, non-empty populations, public composition, and absence of
another operation. Then add only the adapter and private production composition
needed for those two calls. Leave server construction and address policy,
browser and blocking flows, clock, Kibana, Spring, all completed P3/P4 flows,
inventory, scanner, audit apparatus, mutation harnesses, and the rest of P4-P8
unchanged.

## Verification Notes

- The red focused contracts failed because clock sleep, the Kibana retry
  dependency/selector, and the private dependency-taking `post` did not exist.
- Focused clock/Kibana and all relevant caller packages: pass.
- Actual API/CLI, CLI surface, and fresh subprocess compatibility: pass with
  reports and caches outside the worktree.
- Complete and uncached tests, race, vet, pinned lint, complete preflight, Make
  and production-script contracts, and all four host acceptance flows: pass.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- One preflight launcher run hit the documented partial-raw-log signal-fixture
  race at control 25; its immediate unchanged complete rerun passed all 62
  controls and the full preflight.
- The first continuity launcher-contract run hit that race only in its nested
  second-generation copy; its immediate unchanged complete rerun passed all 62
  controls.
- Empty-HOME count-2: pass under `/private/tmp/ply-p4-hermetic.zAFimZ`.
- The focused audit measured all seven requested criteria, exited 1 as
  expected, held Q0.6 at 26 guarded sites, improved exact Q1.3 to 2 of 27, and
  reported zero comparable regressions.
- Full clean audit: expected exit 1, 14 non-passing criteria, five improved,
  two held, zero regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P4 entry and checkpoint gate,
both design documents, `.quality/inventory`, complete webservice start/stop and
caller code/contracts, all complete adapter sources/tests and dependency-
preserving doubles, the completed clock/Kibana contracts, the import-aware
scanner, API/CLI contracts, and the T15 repair and baseline reproduction README
before editing. Confirm branch, HEAD, clean status, reciprocal archive links,
launcher `--check`, and exact product commit `4c0d97a`.

Make one focused server-effect implementation commit, then the normal separate
continuity commit. Stop before address-policy or loopback changes, another
webservice flow, untested-package work, manual L1 evidence, P5-P8, publication,
or distribution. Do not push, merge, stash, revert, launch a successor, or
remove the worktree.
