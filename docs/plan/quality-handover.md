# Quality Upgrade Handover

Generated: 2026-08-26T16:06:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P4.4 product implementation:
  `cbb620a79d0f8797e4c7e2c41358a7f1b6026679`.
- Preceding P4.3 product implementation:
  `9601d29174346b0bbcfb6762c602591c78ba4908`.
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
`docs/plan/agent-sessions/2026-08-26T160625+0200-cover-resources-package.md`.
The P4.4 sorting archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`8d08e706ab7b69c3da7f7766dacd013556a6a1e576af143479a1f6a703128230`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P4 Move 4

Product commit `cbb620a` adds only `pkg/sorting/sort_test.go`. Production
`pkg/sorting/sort.go`, `DependencySort`, its exported slice and `SortKey`
fields, value receivers, callers, and observable ordering remain unchanged.
Nine focused contracts bring `pkg/sorting` to 100% statement coverage and the
suite to 423 tests across 22 of 27 packages.

The contracts prove exact nil, empty, and populated lengths; shared backing
slice behavior through copied value receivers; exact in-place swaps and full
sorts; strict less-than behavior for equal keys; the exact decimal scope,
colon, group-weight, colon, artifact key; every scope class; case-sensitive
substring matching; empty sort keys; matching `1-` and nonmatching `100-`
prefixes; group and artifact tie-breaks; and rejection of empty table
populations before iteration. Dependency fields outside scope, group ID, and
artifact ID remain ignored by the comparison key.

The implementation compares formatted strings rather than numeric tuples. The
unknown-scope key beginning `100:` therefore sorts before the empty-scope key
beginning `10:`. That surprising legacy behavior is now characterized exactly;
normalizing it is a later product decision requiring separate authority. The
focused tests make no request, open no socket, launch no program, wait on no
clock, mutate no global logger or production state, and write no fixture.

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
Q3.9, and all 228 numeric debt leaves. P4.4 did not modify the apparatus,
parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `cbb620a` reports:

- Measured commit tree: `18d1ff3d026dc358f7f854b338fd7d5e2d411c68`.
- Structured scorecard SHA-256:
  `4f7550ad8fd184f09237c5f8ae339a1d746ab8c2b366706d4cd9c5feb79ae0af`.
- Absolute L0: 8 of 8.
- 423 test functions, zero skipped; 22 of 27 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 5 of 27 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 0 direct external sites outside five valid declared adapters of 27
  production effect sites.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 88 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 13 documented non-passing criteria, never 2.
- Comparable ratchets: six improved, two held, zero regressed, zero
  non-comparable.
- Clean identity: zero dirty paths.

The focused seven-criterion scorecard under
`/private/tmp/ply-p4-sorting-focused-audit` has SHA-256
`7f6fbf4354a7b7040c350ba0786b468a054d2febae6b7a3c9031507e16471723`.
Pinned tools remain under `/private/tmp/ply-p4-tools` and
`/private/tmp/ply-p358-tools`. Authoritative output is
`target/quality-audit/scorecard.json`.

## Next Objective

Continue P4 with one test-only coverage move for `pkg/resources`, reducing the
remaining untested logic-package population without changing production code.
Characterize exact `CloudConfig.Implementation().Dir()` selection, `LocalDir`
path composition, `ResourceAsString` resource/filename composition, byte-to-
string preservation, empty content, nested names, and direct read-error return
behavior using only temporary fixtures outside the worktree.

Use explicit non-empty guards for table populations and the central safe writer
for any temporary setup. Keep the move inside a new `pkg/resources` test file,
run the template caller tests, and expect Q1.1 to improve from 5/27 to 4/27
while exact Q1.3 holds at 0/27. Do not change resources, config, file, or
template production; inventory, adapters, seams, scanners, audit apparatus,
mutation harnesses, manual evidence, another untested package, or P5-P8.

## Verification Notes

- Focused sorting, config, and Maven tests pass; `pkg/sorting` reports 100%
  statement coverage.
- Actual API/CLI, CLI surface, fresh executable subprocess, and compatibility
  meta-contracts pass with reports and caches outside the worktree.
- Complete build and uncached tests, race, vet, exact pinned lint, `make test`,
  Make contracts, and all 62 launcher controls pass.
- The first integrated preflight supplied `GOLANGCI_LINT` as a Make command-line
  override. Make propagated it into the fake-linter negative controls, which
  correctly failed because their missing-binary override was masked. Supplying
  the same pinned executable through the environment preserved those controls;
  the unchanged complete preflight then passed.
- All production-script meta-contracts, all four host acceptance flows, and the
  repaired 15-control audit meta-suite pass.
- Empty-HOME `go test ./... -count=2` passes with isolated writable state under
  `/private/tmp/ply-p4-sorting-hermetic.JwUPYB` and the existing module cache.
- Full clean audit: expected exit 1, 13 non-passing criteria, six improved, two
  held, zero regressed, zero non-comparable, and zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the complete P4 entry and
checkpoint gate, both design documents, `.quality/inventory`, complete resource
source and caller, relevant template tests, safe-writer contracts, the
import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README before editing. Confirm branch, HEAD, clean status,
reciprocal archive links, launcher `--check`, and exact product commit
`cbb620a`.

Make one focused resources-test implementation commit, then the normal separate
continuity commit. Stop before production resource/file/config/template changes,
another untested package, manual evidence, P5-P8, publication, or distribution.
Do not push, merge, stash, revert, launch a successor, or remove the worktree.
