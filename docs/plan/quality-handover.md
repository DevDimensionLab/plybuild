# Quality Upgrade Handover

Generated: 2026-08-26T11:19:24+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.61 product implementation:
  `5b678ab676192b835bb6a3a35cfbee4eced3105e`.
- Preceding P3.60 product implementation:
  `850360191468ca57a91f78fc71cbcd6d081ecdb7`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 61 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, 9164e06, 4376e05,
9f56714, 8503601, and 5b678ab in roadmap order. The separate operational
continuity implementation is 1b85711. The focused T15 audit-apparatus repair is
ce736a2; neither operational commit changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T111924+0200-narrow-profile-editor-system-capability.md`.
The P3.61 shell Git capability archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`120d496bfbbd83ec455b6de007e916ecc7b6ebb971566f4ed563ab8b239e0163`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 61

Only the production process-capability selections used by the four public shell
Git execution wrappers changed at product commit `5b678ab`. Private
`systemGitDependencies()` returns the exact existing `process.SystemRunner()`;
`GitClone`, `GitPull`, `GitInit`, and `GitAddAndCommit` select it instead of
`process.System()`. `runGit`, `GitDirty`, `GitIsRepo`, public shell Run, Maven,
every other caller, `process.SystemRunner()`, and `process.System()` are
unchanged.

The exported signatures, executable `git`, exact ordered arguments, debug log
text and placement, distinct stdout and stderr buffer identities and bytes, nil
stdin and directory, synchronous execution, exact one-request and two-request
populations, existing add/commit request selection, direct runner-error
behavior including nil returned `Output.Err`, safe zero dependencies, and all
return paths remain unchanged. There is no fallback, retry, wrapping,
additional logging, execution, cleanup, or global-state change.

Two focused top-level contracts bring the suite to 388 tests. They prove exact
system-runner identity with nil dependency stdout, complete arbitrary commands
and buffer identities, exact request sequences and log/process placement,
legacy error behavior, safe zero behavior for every wrapper, complete
caller-owned dependency preservation, non-empty populations, and no additional
process request. Focused tests record in memory, launch no external program,
touch no network, and write no repository fixture.

## T15 Apparatus Checkpoint

The P3.54 T15 repair remains unchanged. Separate verified execution and
structured replicas preserve historical test effects while measuring a
pristine exact `5635d50` checkout plus the declared inventory overlay. No
output path is excluded or deleted, and every effect stays in the retained
manifest. The apparatus-owned reproduction HOME includes the required
directory-shaped `.co-pilot/profiles/.active_profile` fixture with marker
SHA-256
`cb95f24c35d3987f8aba51231aade19580ffe9242324804fcad2ccff350d1c9a`.

The complete 15-control proof passes and reproduces old scorecard
`d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`,
new scorecard
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`,
normalized stored raw body
`cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f`,
Q3.9, and all 228 numeric debt leaves. Product move 61 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `5b678ab` reports:

- Measured commit tree: `010970642433199467db2fa3b6eaaf7c9ab7bece`.
- Structured scorecard SHA-256:
  `0632a71b59dc2793c2cf308d0514569b698cdead3f02fef83a87cbe029fccfff`.
- Absolute L0: 8 of 8.
- 388 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 10 direct external sites outside five declared adapters of 32
  production effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 79 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The unchanged scanner no longer follows the `os.Stdout` field through the four
shell Git runner-only production dependencies, so exact Q1.3 improves from
14/36 to 10/32. The authoritative audit output is
`target/quality-audit/scorecard.json`. Focused audit output is under
`/private/tmp/ply-p361-gate/focused-audit`; pinned tools are under
`/private/tmp/ply-p358-tools`; the complete gate root is
`/private/tmp/ply-p361-gate`. The disposable empty-HOME root is
`/private/tmp/ply-p361-hermetic.ZiSud0`.

## Next Objective

Complete one focused P3 provenance move: narrow only the existing private
profile-editor production selector's process dependency from `process.System()`
to the exact runner-only `process.SystemRunner()` so it no longer inherits an
unused dependency standard-output capability.

Start red by strengthening the existing profile-editor recording contracts.
Prove exact system-runner identity and nil process-dependency stdout; separate
exact `os.Stdin` and `os.Stdout` command streams; complete caller-owned
dependency preservation; exact arbitrary editor and config-path arguments;
empty directory, nil stderr, synchronous one-request execution; exact runner
error; safe zero behavior; a non-empty population; and no second process
request.

Then change only `systemProfileEditorDependencies()` to select
`process.SystemRunner()`. Preserve `runProfileEditor`, public command behavior,
environment and default-editor selection, command streams and fields, profile
flows, Maven stdout, shell Run and Git, plugin diagrams, browser, Unzip, every
other process caller, inventory, scanner, audit apparatus, and P4-P8.

## Verification Notes

- Focused shell Git contracts first failed only because the private selector
  was absent.
- Focused process/shell and every relevant process caller package: pass.
- API/CLI and fresh subprocess compatibility: pass with reports outside the
  worktree.
- Complete tests, race, vet, pinned lint, build, complete preflight, Make and
  production-script contracts, and all four host acceptance flows: pass.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass under `/private/tmp/ply-p361-hermetic.ZiSud0`.
- Focused audit: expected exit 1, four improved, two held, zero regressed, one
  non-comparable, zero dirty paths.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete process adapter code and
tests and every complete process double/caller, complete profile command and
editor contracts and callers, P3.61 shell Git contracts, P3.60 shell Run
contracts, P3.58 Maven stdout contracts, the import-aware scanner, API/CLI
contracts, and the T15 repair and baseline reproduction README before editing.
Confirm branch, HEAD, clean status, reciprocal archive links, launcher
`--check`, and exact product commit `5b678ab`.

Make one focused profile-editor system-capability implementation commit, then
the normal separate continuity commit. Stop before Maven, shell Run or Git,
plugin diagrams, browser, Unzip, another process caller, clock/server work,
mutation harnesses, P4-P8, audit or scanner changes, inventory changes,
publication, or distribution. Do not push, merge, stash, revert, launch a
successor, or remove the worktree.
