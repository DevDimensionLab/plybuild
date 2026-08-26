# Quality Upgrade Handover

Generated: 2026-08-26T19:12:51+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.10 product implementation:
  `c382caa38be167fe17f847370ad8a12270644de3`.
- Preceding P4.9 product implementation:
  `746a5abbdc5b879201bc83690b3c633831dca81b`.
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
`docs/plan/agent-sessions/2026-08-26T191251+0200-record-p4-manual-l1-evidence.md`.
The P4.10 loopback archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`a37605edc2c426e6ea0984c0f2436977c0249158d020523fd56c262a699dade0`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 10

Product commit `c382caa` changes only the private package-level interactive
HTTP server address from `:7999` to `127.0.0.1:7999`. Port 7999, the exact
`*http.Server` identity, selector and adapter, four handler paths and order,
listen/log behavior, background shutdown context, five-second timeout, deferred
cancellation, ignored shutdown error, endpoint URIs, browser request, callback
blocking, public signatures, and all callers remain exact. Tests do not open a
socket, start the server, or launch a browser.

Four focused contracts bring the suite to 453 tests across all 27 packages.
They pin the exact IPv4 loopback address and unchanged start/stop selection;
endpoint URIs and server/browser/callback composition; standalone and project
blocking flow; interactive build option delivery and validation order; and
interactive upgrade guards and callback selection.

The regular non-executable `scripts/mutate-interactive-build` seam driver
contains `8. interactive server binds only to loopback` exactly once, declares
no `mutate` or `mutate2` operation, and runs 21 named command, webservice,
server-adapter, API, and callback contracts with isolated Go state. Its
reciprocal meta-contract proves its mode and source shape, exact label, absence
of premature mutation operations, successful JSON test execution, non-empty
population, and every named test run. All eight declared P4 seams now have
executable swap coverage.

## T15 Apparatus Checkpoint

The P3.54 T15 repair remains unchanged. Separate verified execution and
structured replicas preserve historical test effects while measuring a
pristine exact `5635d50` checkout plus the declared inventory overlay. The
complete 15-control proof passes and reproduces old scorecard
`d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`,
new scorecard
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`,
normalized stored raw body
`cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f`,
Q3.9, and all 228 numeric debt leaves. P4.10 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `c382caa` reports:

- Measured commit tree: `6e44c014c73fb8c99c6c90f20a64a41f314fce43`.
- Structured scorecard SHA-256:
  `b4c7ddc7d163a833663f8b7cd876e94aadf4c1a875b945252a46441e354c6b72`.
- Absolute L0: 8 of 8.
- 453 test functions, zero skipped; all 27 packages have tests.
- Q0.6: 29 guarded safe-writer sites, 23 write and 6 copy, with zero unsafe
  direct test writes.
- Q0.8: all 13 production scripts have reciprocal meta-tests.
- Q1.1: 0 of 27 packages has no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 0 direct external sites outside five valid declared adapters of 27
  production effect sites.
- Q1.4: 8 of 8 declared seams covered.
- Q1.6, Q1.7, and Q1.9: UNMEASURABLE because current-project manual receipts
  have not yet been supplied.
- Exact Q2.1: 0 of 8 subjects has an executable harness.
- Q3.4: 0 state-claim phrases across 94 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 13 documented non-passing criteria, never 2.
- Comparable ratchets: six improved, two held, zero regressed, zero
  non-comparable.
- Clean identity: zero dirty paths.

The focused eight-criterion scorecard under
`/private/tmp/ply-p4-interactive-focused-audit.b6Hgfs` has SHA-256
`b63afac3b5917dd45d74e19ace4fbf588fe3c871356db4abf1d0c428dea6f703`.
The full clean audit is under
`/private/tmp/ply-p4-interactive-full-audit.YuyZ2i`. Pinned tools remain under
`/private/tmp/ply-p4-tools` and `/private/tmp/ply-p358-tools`. The empty-HOME
count-2 proof is under `/private/tmp/ply-p4-interactive-hermetic.RFfdQW`.

## Next Objective

Finish P4 only if exhaustive source-and-test review can truthfully support
schema-2 manual receipts for Q1.6, Q1.7, and Q1.9. Enumerate a non-empty exact
population for each row: dependency structs and their default/recording doubles
for Q1.6; per-item or partial-failure command paths and content assertions for
Q1.7; and test-iterated collections plus executable empty-population assertions
for Q1.9. Record the commands, subjects, exact counts, and exclusions. If any
required truth fails, do not claim PASS; record the exact failing population
and keep P4 active.

The measured-tree implementation deliberately counts ignored and untracked
repository bytes. Therefore, build the commit-bound manual-evidence document
outside the worktree, never at `.quality/manual-evidence.json`, and pass it
explicitly with `--manual-evidence`. Bind it to the exact clean current commit,
tree, status digest, inventory, and audit instruments; include only the three
authorized L1 receipts and canonical evidence digests. A valid clean full audit
must make Q1.6, Q1.7, and Q1.9 PASS while preserving every automated verdict
and ratchet. Record the external document SHA-256, receipt digests, populations,
scorecard identity, and audit exit in tracked continuity. Do not weaken the
validator, scanner, inventory, baseline, contracts, or product code to obtain a
receipt.

## Verification Notes

- The initial exact-address test failed only on the old `:7999` value; it passed
  after the one-line production change and exact start/stop assertion updates.
- Focused command, webservice, API, server-adapter, driver, coverage,
  shuffle-count-10, and vet checks pass. The driver meta-test reports all 21
  named contracts.
- Actual API/CLI, CLI surface, fresh executable subprocess, and compatibility
  meta-contracts pass with reports and caches outside the worktree.
- Complete build and tests, race, vet, exact pinned lint, `make preflight`,
  `make test`, Make contracts, install, and all 62 launcher controls pass.
- All production-script meta-contracts, all four host acceptance flows, and the
  repaired 15-control audit meta-suite pass.
- Broader unchanged launcher invocations intermittently hit the documented
  partial-raw-log signal-fixture race at controls 26 and 50. Unchanged
  standalone launcher and complete Make reruns passed all 62 controls.
- Empty-HOME `go test ./... -count=2` passes with isolated writable state and
  the existing read-only module cache.
- Full clean audit: expected exit 1, 13 non-passing criteria, six improved, two
  held, zero regressed, zero non-comparable, and zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and
checkpoint gate, both design documents, `.quality/README.md`, the complete
schema-2 parser and negative meta-test, baseline reproduction materials, all
dependency-double, partial-failure, and non-empty-population contracts, the
full audit wrapper, and the T15 repair before recording evidence. Confirm
branch, HEAD, clean status, reciprocal archive links, launcher `--check`, and
exact product commit `c382caa`.

Do not make a product-code or test change for this evidence-only checkpoint.
Measure the clean current continuity head with an external receipt file, then
make the normal continuity-only commit that records the result and advances to
P5 only when every P4 exit is truthfully satisfied. Stop before an executable
P5 harness, T1-T10, P6-P8, publication, or distribution. Do not push, merge,
stash, revert, launch a successor, or remove the worktree.
