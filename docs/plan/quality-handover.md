# Quality Upgrade Handover

Generated: 2026-08-26T16:58:09+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.6 product implementation:
  `8eeeb2a8fe923bd644b65f3c77c6fae805b0ee17`.
- Preceding P4.5 product implementation:
  `64c618944bd36dad86d34fc058bf8dcf036917cf`.
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
`docs/plan/agent-sessions/2026-08-26T165809+0200-cover-webservice-api-package.md`.
The P4.6 webservice-templates archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`bc9103a6d9a6f8378be9929bd9d3bc296262795be6a99461192db0616a716fc1`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 6

Product commit `8eeeb2a` adds only
`pkg/webservice/templates/templates_test.go`. All four production template
sources, exported untyped string constants, private fragments, API callers,
HTML, form behavior, and observable bytes remain unchanged. Three focused
contracts bring the suite to 429 tests across 24 of 27 packages. The package
contains no statements, so Go truthfully reports `[no statements]` coverage.

The contracts prove exact complete lengths and SHA-256 identities for Generate,
Upgrade, and every private header/body/footer fragment; exact shared-header +
private-body + shared-footer composition and boundaries; form actions, input
names, template actions and counts, CDN references; successful caller-selected
`text/template` and `html/template` parsing; exported constant names and untyped
assignability; and rejection of empty table populations before iteration.

The tests make no request, open no socket, launch no program, wait on no clock,
mutate no global API options, channel, logger, or production state, and write no
fixture. Existing unquoted Generate values and `text/template` parsing remain
legacy behavior; preferred quoting or HTML escaping is a separate product
decision requiring explicit authority.

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
Q3.9, and all 228 numeric debt leaves. P4.6 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `8eeeb2a` reports:

- Measured commit tree: `9c7322b185f485b4202cec0691ce9fd39c8dbdc0`.
- Structured scorecard SHA-256:
  `213df9ece4eb7c231ec3faa8f76e9180e9215a6ef0fc7d83039911940224c7e7`.
- Absolute L0: 8 of 8.
- 429 test functions, zero skipped; 24 of 27 packages have tests.
- Q0.6: 27 guarded safe-writer sites, 22 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 3 of 27 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 0 direct external sites outside five valid declared adapters of 27
  production effect sites.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 90 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 13 documented non-passing criteria, never 2.
- Comparable ratchets: six improved, two held, zero regressed, zero
  non-comparable.
- Clean identity: zero dirty paths.

The focused seven-criterion scorecard under
`/private/tmp/ply-p4-templates-focused-audit.FWPBwe` has SHA-256
`aaa5507955ca84ef255432c998ee2c3092776caeb737754736e78752bc49a57a`.
Pinned tools remain under `/private/tmp/ply-p4-tools` and
`/private/tmp/ply-p358-tools`. Authoritative output is
`target/quality-audit/scorecard.json`.

## Next Objective

Continue P4 with one test-only coverage move for `pkg/webservice/api`, reducing
the remaining untested package population without changing production code.
Characterize exact Generate and Upgrade GET rendering, POST form-to-project
mutations and ordered slice appends, exact `OK` responses, and one callback per
POST while restoring every package global and leaking no goroutine.

Use only in-memory `httptest` requests/recorders and a test-owned callback
channel, with explicit non-empty guards and no fixture writes. Keep the move
inside one new `pkg/webservice/api` test file, run cmd and webservice caller
tests, and expect Q1.1 to improve from 3/27 to 2/27 while exact Q1.3 holds at
0/27. Do not change API handlers/globals/channels, templates, server behavior,
config, Spring, resources, inventory, adapters, seams, scanners, audit
apparatus, mutation harnesses, manual evidence, another untested package, or
P5-P8.

## Verification Notes

- Focused templates, API, and webservice caller tests pass; the constants-only
  templates package truthfully reports `[no statements]` coverage.
- Actual API/CLI, CLI surface, fresh executable subprocess, and compatibility
  meta-contracts pass with reports and caches outside the worktree.
- Complete build and uncached tests, race, vet, exact pinned lint,
  `make preflight`, `make test`, Make contracts, and all 62 launcher controls
  pass.
- The first `make test` invocation hit the documented nested partial-raw-log
  signal-fixture race at launcher control 26; its immediate unchanged complete
  rerun passed all controls. All isolated focused, preflight, compatibility,
  race, vet, lint, and acceptance runs passed directly.
- All production-script meta-contracts, all four host acceptance flows, and the
  repaired 15-control audit meta-suite pass.
- Empty-HOME `go test ./... -count=2` passes with isolated writable state under
  `/private/tmp/ply-p4-templates-hermetic.Yk3uCL` and the existing module cache.
- Full clean audit: expected exit 1, 13 non-passing criteria, six improved, two
  held, zero regressed, zero non-comparable, and zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and
checkpoint gate, both design documents, `.quality/inventory`, all three API
files, every caller, the complete template sources/contracts, relevant
server/cmd/config/Spring types and tests, representative handler/global/
callback/non-empty contracts, the import-aware scanner, API/CLI contracts, and
the T15 repair and baseline reproduction README before editing. Confirm branch,
HEAD, clean status, reciprocal archive links, launcher `--check`, and exact
product commit `8eeeb2a`.

Make one focused webservice-API test implementation commit, then the normal
separate continuity commit. Stop before production API/template/server/config/
Spring/resources changes, another untested package, manual evidence, P5-P8,
publication, or distribution. Do not push, merge, stash, revert, launch a
successor, or remove the worktree.
