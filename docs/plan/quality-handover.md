# Quality Upgrade Handover

Generated: 2026-08-26T15:39:36+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.3 product implementation:
  `9601d29174346b0bbcfb6762c602591c78ba4908`.
- Preceding P4.2 product implementation:
  `4c0d97ad2ac1a6e5faee6ec6f240eeed43c4b434`.
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
`docs/plan/agent-sessions/2026-08-26T153936+0200-cover-sorting-package.md`.
The P4.3 web-server archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`b99cfa23f141551c169e130d4b6438aae212a419a914b07d7d44590f10447ee8`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 3

Product commit `9601d29` adds only the declared `internal/adapter/server`
boundary and migrates webservice's two remaining direct server operations. The
adapter exposes complete listen and shutdown operations plus a caller-owned
server selector. Its safe-zero paths select no server, perform no network
operation, and return deterministic nil. Production calls the corresponding
method once on the exact selected `*http.Server`, passes the exact shutdown
context, and returns the exact error without wrapping, replacement, recovery,
retry, logging, cleanup, address rewriting, handler mutation, or global state.

The exact package-level server construction remains
`&http.Server{Addr: fmt.Sprintf(":%d", port)}` with `port == 7999`. Private
start and stop dependency values select that same pointer through
`selectWebServer`; dereferencing happens only inside the declared adapter. The
unchanged import-aware scanner therefore records the two concrete system calls
inside the adapter and zero direct effects outside all five valid adapters.

Public `StartWebServer` still registers `/ui/generate`, `/api/generate`,
`/ui/upgrade`, and `/api/upgrade` with the same callbacks in that order before
one listen attempt. It prints the exact non-nil listen error through
`log.Print` and prints nothing for nil. Public `StopWebServer` still derives one
context from `context.Background()` with exactly `5*time.Second`, defers its
cancellation, passes that exact context to one shutdown attempt, and ignores
the exact shutdown error. Signatures, callers, handlers, browser and blocking
flows, address, timeout, and public return behavior are unchanged.

Thirteen net new focused tests bring the suite to 414 tests. Adapter contracts
prove safe-zero operation, exact selector/server/context/error delivery,
single selection and operation counts, direct system delegation, complete
dependency identity, no unrelated operation, and rejection of empty recorded
populations. Webservice contracts prove exact production selectors, complete
caller-owned dependency preservation by every recorder, four exact handler
registrations before listen, nil/non-nil logging, exact
background/timeout/shutdown/cancel order, ignored shutdown error, safe-zero
compositions, public composition, and non-empty populations. They open no
socket, make no request, wait on no clock, launch no program, mutate no
production server, and write no repository fixture.

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
Q3.9, and all 228 numeric debt leaves. P4.3 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `9601d29` reports:

- Measured commit tree: `fd78ef577b886a423d2964f4352503d2aac8992b`.
- Structured scorecard SHA-256:
  `010923f3718d99723c99c76a6b541669cb7f99d7a4a8eb119044705be005cdfa`.
- Absolute L0: 8 of 8.
- 414 test functions, zero skipped; 21 of 27 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 27 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 0 direct external sites outside five valid declared adapters of 27
  production effect sites.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 87 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 13 documented non-passing criteria, never 2.
- Comparable ratchets: six improved, two held, zero regressed, zero
  non-comparable.
- Clean identity: zero dirty paths.

The focused seven-criterion scorecard under
`/private/tmp/ply-p4-server-focused-audit-final` has SHA-256
`d52e2d34017dac46134a4eacc2fca628e4e39bed995995ab7837217796ca65ed`.
Pinned tools remain under `/private/tmp/ply-p4-tools` and
`/private/tmp/ply-p358-tools`. Authoritative output is
`target/quality-audit/scorecard.json`.

## Next Objective

Continue P4 with one test-only coverage move for `pkg/sorting`, reducing the
remaining untested logic-package population without changing production code.
Characterize exact `DependencySort` length, in-place swap, and less-than
behavior; exact concatenated scope/group/artifact keys; all scope weights;
case-sensitive and empty sort-key containment; matching and nonmatching group
weights; group and artifact tie-breaks; full in-place sort results; and the
legacy lexical consequence of formatting unknown weight 100 as text. Preserve
surprising existing order rather than normalizing it.

Use explicit non-empty guards for table-driven populations. Keep the move
inside a new `pkg/sorting` test file, run relevant config/Maven callers, and
expect Q1.1 to improve from 6/27 to 5/27 while exact Q1.3 holds at 0/27. Do not
change `pkg/sorting/sort.go`, config/Maven production, inventory, adapters,
seams, scanners, audit apparatus, mutation harnesses, manual evidence, another
untested package, the loopback-address policy, or P5-P8.

## Verification Notes

- Red focused contracts initially failed because the server adapter and private
  webservice dependency compositions did not exist.
- Focused server/webservice and relevant caller tests: pass without network or
  production server mutation.
- Actual API/CLI, CLI surface, and fresh subprocess compatibility: pass with
  reports and caches outside the worktree.
- Complete tests, race, vet, pinned lint, `make test`, complete preflight, Make
  and production-script contracts, and all four host acceptance flows: pass.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Repeated `make test` attempts encountered the documented nested
  partial-raw-log signal-fixture race before an unchanged complete run passed.
- One focused-audit attempt used invalid `--only` syntax and failed closed;
  the comma-separated exact rerun produced the recorded 0/27 scorecard.
- One final preflight attempt supplied script-level rather than Make-level
  pinned tool variables; the corrected complete rerun passed.
- Empty-HOME `go test ./... -count=2`: pass with isolated writable state under
  `/private/tmp` and the existing module cache.
- Full clean audit: expected exit 1, 13 non-passing criteria, six improved, two
  held, zero regressed, zero non-comparable, and zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and
checkpoint gate, both design documents, `.quality/inventory`, complete sorting
source and callers, relevant config/Maven tests, representative non-empty
table contracts, the import-aware scanner, API/CLI contracts, and the T15
repair and baseline reproduction README before editing. Confirm branch, HEAD,
clean status, reciprocal archive links, launcher `--check`, and exact product
commit `9601d29`.

Make one focused sorting-test implementation commit, then the normal separate
continuity commit. Stop before production sorting changes, another untested
package, loopback/address-policy changes, manual evidence, P5-P8, publication,
or distribution. Do not push, merge, stash, revert, launch a successor, or
remove the worktree.
