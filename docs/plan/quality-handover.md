# Quality Upgrade Handover

Generated: 2026-08-26T17:28:16+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.7 product implementation:
  `4da64d2468307d39ed3e8303ad9183fa18d71a75`.
- Preceding P4.6 product implementation:
  `8eeeb2a8fe923bd644b65f3c77c6fae805b0ee17`.
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
`docs/plan/agent-sessions/2026-08-26T172816+0200-cover-logger-package.md`.
The P4.7 webservice-API archive is answered history and links reciprocally to
that archive. The graph has exactly one NEXT tail. The archived prompt SHA-256
is `551a599d5cb7993bd7e215a9b0bdf985598d29e9ecac7305f49ede79ff4dabc6`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 7

Product commit `4da64d2` adds only
`pkg/webservice/api/handlers_test.go`. All three production API files, exported
handler signatures, `GenerateOptions`, package globals, channel identity and
capacity, routes, callers, forms, templates, HTML, callback timing, and response
semantics remain unchanged. Five focused contracts bring the suite to 434
tests across 25 of 27 packages and cover all four handlers at 100% statement
coverage.

Generate GET uses only an in-memory recorder/request and a static CloudConfig
test value. It executes the existing `text/template` path and pins status 200,
the exact 2,957-byte body with SHA-256
`9a9eb38ab026e751b06dee004e84d8db9a427454c5e33f4b864ebc83635399c1`,
one template-list request, and representative raw project, template, dependency
group, dependency ID, and dependency name bytes. Upgrade GET executes the
existing `html/template` path and pins status 200 plus the exact 788-byte body
with SHA-256
`ae9194d9a6734d1fe568692e11327fc1b31b3f0d797fd558fd6ceba3ffa58ecc`.
The raw unquoted Generate values are deliberate characterization; preferred
validation, quoting, or escaping remains a later product decision requiring
explicit authority.

Generate POST proves exact body parsing, body-over-query first-value selection
for all six scalar fields, preservation of every other project field, and
ordered duplicate-preserving append of all template and dependency values.
Both POST handlers return status 200 and exact `OK` before a receiver exists on
the test-owned unbuffered callback channel, then deliver exactly one `true`
callback. Upgrade POST retains its exact choice not to parse a malformed form
body. Each unbuffered rendezvous releases the only sender. The tests perform no
sleep, timed wait, socket, external program, logger mutation, or fixture write,
and leak no goroutine.

Every handler subtest restores `GOptions`, `CallbackChannel`, and
`CurrentProject`; its outer contract verifies restored pointer/channel
identities and values. Recorded callback and iterated table populations reject
empty inputs. No production code, scanner, inventory, audit apparatus, manual
evidence, mutation harness, or completed contract changed.

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
Q3.9, and all 228 numeric debt leaves. P4.7 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `4da64d2` reports:

- Measured commit tree: `d23edc06e2f958c43f2add2a43c540f48c7145e9`.
- Structured scorecard SHA-256:
  `0682df76893e0eac05302870ea986b1f18cab32e7513a7e0e91ab85da5f16c38`.
- Absolute L0: 8 of 8.
- 434 test functions, zero skipped; 25 of 27 packages have tests.
- Q0.6: 27 guarded safe-writer sites, 22 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 2 of 27 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 0 direct external sites outside five valid declared adapters of 27
  production effect sites.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 91 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 13 documented non-passing criteria, never 2.
- Comparable ratchets: six improved, two held, zero regressed, zero
  non-comparable.
- Clean identity: zero dirty paths.

The focused seven-criterion scorecard under
`/private/tmp/ply-p4-api-focused-audit.EVMVn0` has SHA-256
`54b314c71a0e6bce29b84647d347d15b643ecf3c8baf982876411ef7b72cb428`.
The full clean audit is under
`/private/tmp/ply-p4-api-full-audit.1qAqQS`. Pinned tools remain under
`/private/tmp/ply-p4-tools` and `/private/tmp/ply-p358-tools`.

## Next Objective

Continue P4 with one test-only coverage move for `pkg/logger`, reducing the
remaining untested package population without changing production code or
global semantics. Characterize the exact `Collector` hook levels and ordered
entry-pointer capture; the package logger's installed collector hook;
`DebugLogger`, `Context`, and `ExternalError`; JSON formatter selection;
field-logger state; `StdOut`; and `LogEntries`.

Keep the move inside one new `pkg/logger` test file. Do not use parallel tests
around package or logrus globals. Snapshot and restore `fieldLogger`,
`collector`, `log`, and every mutated standard-logrus value, preserving exact
pointer identities and leaving no state after cleanup. Use only in-memory
logrus values and buffers; do not write a fixture, open a socket, launch a
program, wait on a clock, or add a seam. Reject empty table and recorded-entry
populations before iterating. Expect Q1.1 to improve from 2/27 to 1/27 while
exact Q1.3 remains 0/27. Leave `pkg/context`, API handlers, templates, callers,
logging behavior, scanner, inventory, apparatus, mutation work, manual
evidence, and P5-P8 unchanged.

## Verification Notes

- Focused API, webservice, templates, and command caller tests pass; the API
  package reports 100% statement coverage and its focused race run passes.
- Actual API/CLI, CLI surface, fresh executable subprocess, and compatibility
  meta-contracts pass with reports and caches outside the worktree.
- Complete build and uncached tests, race, vet, exact pinned lint,
  `make preflight`, `make test`, Make contracts, and all 62 launcher controls
  pass.
- The first two standalone launcher invocations hit the documented partial-
  raw-log signal-fixture race at control 26. The third passed its primary
  signal control but hit the same race in the nested control-50 run. The fourth
  unchanged complete invocation passed all 62 controls.
- The first `make test` invocation hit the same launcher race at control 26;
  its unchanged complete rerun passed.
- All production-script meta-contracts, all four host acceptance flows, and the
  repaired 15-control audit meta-suite pass.
- Empty-HOME `go test ./... -count=2` passes with isolated writable state under
  `/private/tmp/ply-p4-api-hermetic.oKxSI3` and the existing module cache.
- Full clean audit: expected exit 1, 13 non-passing criteria, six improved, two
  held, zero regressed, zero non-comparable, and zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and
checkpoint gate, both design documents, `.quality/inventory`, both complete
logger production files and package initialization, all logger callers,
relevant logrus types and package-level state behavior, representative
global-restoration/log-capture/exact-error/stdout/non-empty contracts, the
import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README before editing. Confirm branch, HEAD, clean status,
reciprocal archive links, launcher `--check`, and exact product commit
`4da64d2`.

Make one focused logger test implementation commit, then the normal separate
continuity commit. Stop before production logger/context/API/template changes,
another untested package, manual evidence, P5-P8, publication, or distribution.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
