# Quality Upgrade Handover

Generated: 2026-08-26T16:34:26+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.5 product implementation:
  `64c618944bd36dad86d34fc058bf8dcf036917cf`.
- Preceding P4.4 product implementation:
  `cbb620a79d0f8797e4c7e2c41358a7f1b6026679`.
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
`docs/plan/agent-sessions/2026-08-26T163426+0200-cover-webservice-templates-package.md`.
The P4.5 resources archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`7e55f3c10405261604beccd8e9a33a4c64e7e8fbae2cd46c7baf4b66d958bc79`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 5

Product commit `64c6189` adds only `pkg/resources/resources_test.go`.
Production `pkg/resources/resources.go`, `LocalDir(config.CloudConfig) string`,
`ResourceAsString(config.CloudConfig, string) (string, error)`, config and file
boundaries, callers, and observable filesystem behavior remain unchanged.
Three focused contracts bring `pkg/resources` to 100% statement coverage and
the suite to 426 tests across 23 of 27 packages.

The contracts prove exact one-time `Implementation` and `Dir` selection with no
`FilePath` call; representative arbitrary, empty, and trailing-separator root
bytes; the existing `file.Path` directory construction; literal resource/
filename slash composition; successful empty, arbitrary, non-UTF-8, and nested
resource reads; exact byte-to-string preservation; the direct missing
`*os.PathError`, exact uncleaned path bytes, and missing-file cause; and
rejection of empty table populations before iteration.

All fixture files live below `t.TempDir()`, their parent tree is prepared first,
and their writes go through the central safe writer. The tests make no request,
open no socket, launch no program, wait on no clock, mutate no global logger or
production state, and write no repository fixture. The retained duplicate and
mixed separators are legacy behavior; normalization remains a separate product
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
Q3.9, and all 228 numeric debt leaves. P4.5 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `64c6189` reports:

- Measured commit tree: `88da045036be8849a8a3ffb327e01bb5f23c81b5`.
- Structured scorecard SHA-256:
  `2bbbc911efefe7427cd3413b73b0a96583a1111460e57f5cdfb9b145b6818a95`.
- Absolute L0: 8 of 8.
- 426 test functions, zero skipped; 23 of 27 packages have tests.
- Q0.6: 27 guarded safe-writer sites, 22 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 4 of 27 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 0 direct external sites outside five valid declared adapters of 27
  production effect sites.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 89 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 13 documented non-passing criteria, never 2.
- Comparable ratchets: six improved, two held, zero regressed, zero
  non-comparable.
- Clean identity: zero dirty paths.

The focused seven-criterion scorecard under
`/private/tmp/ply-p4-resources-focused-audit` has SHA-256
`9eedfafb6977fae6a0e27b878174b49ecf3536b7d8997b2746df41f722f34430`.
Pinned tools remain under `/private/tmp/ply-p4-tools` and
`/private/tmp/ply-p358-tools`. Authoritative output is
`target/quality-audit/scorecard.json`.

## Next Objective

Continue P4 with one test-only coverage move for
`pkg/webservice/templates`, reducing the remaining untested package population
without changing production code. Characterize the complete exact `Generate`
and `Upgrade` bytes, their header/body/footer concatenation and boundaries,
their expected form/action/template tokens, and successful parsing through the
same `text/template` and `html/template` engines selected by their callers.

Use explicit non-empty guards for table populations and no fixture writes.
Keep the move inside a new `pkg/webservice/templates` test file, run the API and
webservice caller tests, and expect Q1.1 to improve from 4/27 to 3/27 while
exact Q1.3 holds at 0/27. Do not change template constants, API handlers,
server behavior, config, Spring, resources, inventory, adapters, seams,
scanners, audit apparatus, mutation harnesses, manual evidence, another
untested package, or P5-P8.

## Verification Notes

- Focused resources and template caller tests pass; `pkg/resources` reports
  100% statement coverage.
- Actual API/CLI, CLI surface, fresh executable subprocess, and compatibility
  meta-contracts pass with reports and caches outside the worktree.
- Complete build and uncached tests, race, vet, exact pinned lint,
  `make preflight`, `make test`, Make contracts, and all 62 launcher controls
  pass.
- The first focused coverage command and first preflight invocation reached
  sandbox-blocked default build/lint caches; unchanged isolated-cache reruns
  passed. The first `make test` invocation hit the documented nested partial-
  raw-log signal-fixture race; its immediate unchanged complete rerun passed.
- All production-script meta-contracts, all four host acceptance flows, and the
  repaired 15-control audit meta-suite pass.
- Empty-HOME `go test ./... -count=2` passes with isolated writable state under
  `/private/tmp/ply-p4-resources-hermetic.3ASoTT` and the existing module cache.
- Full clean audit: expected exit 1, 13 non-passing criteria, six improved, two
  held, zero regressed, zero non-comparable, and zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and
checkpoint gate, both design documents, `.quality/inventory`, all four template
constant sources, every caller and relevant webservice test, representative
exact-byte and non-empty-population contracts, the import-aware scanner,
API/CLI contracts, and the T15 repair and baseline reproduction README before
editing. Confirm branch, HEAD, clean status, reciprocal archive links, launcher
`--check`, and exact product commit `64c6189`.

Make one focused webservice-templates test implementation commit, then the
normal separate continuity commit. Stop before production template/API/server/
config/Spring/resources changes, another untested package, manual evidence,
P5-P8, publication, or distribution. Do not push, merge, stash, revert, launch
a successor, or remove the worktree.
