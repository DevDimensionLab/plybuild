# Quality Upgrade Handover

Generated: 2026-08-26T10:48:16+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.60 product implementation:
  `850360191468ca57a91f78fc71cbcd6d081ecdb7`.
- Preceding P3.59 product implementation:
  `9f56714234ef4aa07cd04b5ebbb4994d388ade4c`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 60 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, 9164e06, 4376e05,
9f56714, and 8503601 in roadmap order. The separate operational continuity
implementation is 1b85711. The focused T15 audit-apparatus repair is ce736a2;
neither operational commit changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T104816+0200-narrow-shell-git-system-capabilities.md`.
The P3.60 shell Run capability archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`578798bfc15a6925974fa4d1aecf1fe9285765da616de0dd3d6bf6808598ee21`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 60

Only public shell `Run`'s production process-capability selection changed at
product commit `8503601`. The process adapter now exports `SystemRunner()`,
which returns the exact existing `systemRunner` with nil stdout. Private
`systemRunDependencies()` selects it instead of `process.System()`. The complete
`process.System()` implementation remains the same runner plus exact
`os.Stdout`, so Maven debug stdout and every other process caller are unchanged.

The exported `Run(name string, args ...string) Output` signature, command name
and ordered arguments, debug log text and placement, distinct stdout and stderr
buffer identities, nil stdin and directory, synchronous execution, one exact
runner request, returned bytes, legacy nil `Output.Err` after a runner error,
safe zero dependency, and all return paths remain unchanged. There is no
fallback, retry, wrapping, extra logging, execution, cleanup, or global-state
change.

Two focused top-level contracts bring the suite to 386 tests. They prove exact
system-runner identity, nil stdout on `SystemRunner()`, unchanged exact
`os.Stdout` on `System()`, complete arbitrary command and buffer identities,
exact dependency-error behavior, safe zero behavior, complete caller-owned
dependency preservation, exact log-before-process ordering, a non-empty
population, and no second process request. Focused tests record in memory,
launch no external program, touch no network, and write no repository fixture.

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
Q3.9, and all 228 numeric debt leaves. Product move 60 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `8503601` reports:

- Measured commit tree: `480c44f20723b48c74243ab4b8c5618c8ddf81fd`.
- Structured scorecard SHA-256:
  `09ef869b92d3d77a6deb0e5e4f1cd9bdeec1d72596ace95648d6165cf10797e1`.
- Absolute L0: 8 of 8.
- 386 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 14 direct external sites outside five declared adapters of 36
  production effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 78 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The unchanged scanner no longer follows the `os.Stdout` field through shell
Run's runner-only dependency, so exact Q1.3 improves from 15/37 to 14/36. The
authoritative audit output is `target/quality-audit/scorecard.json`. Focused
audit output is under `/private/tmp/ply-p360-gate.BuxOfR/focused-audit`; pinned
tools are under `/private/tmp/ply-p358-tools`; the complete gate root is
`/private/tmp/ply-p360-gate.BuxOfR`. The disposable empty-HOME root is
`/private/tmp/ply-p360-hermetic.xS9Thp`.

## Next Objective

Complete one focused P3 provenance move: give the four public shell Git
execution wrappers one runner-only system process dependency so `GitClone`,
`GitPull`, `GitInit`, and `GitAddAndCommit` no longer inherit unused standard
output.

Start red with shell Git recording contracts for the exact system runner and
nil stdout, complete caller-owned process dependency preservation, exact Git
command names and ordered arguments, exact log text and placement, distinct
stdout/stderr buffer identities and bytes, nil stdin and directory, synchronous
execution, exact one/two request populations, legacy error behavior, safe zero
behavior, non-empty populations, and absence of another process request.

Then add one private complete production selector that returns the existing
`process.SystemRunner()` and replace only the four `process.System()` selections
in those public wrappers. Preserve `GitDirty`/`GitIsRepo`, `runGit`, public API,
Maven stdout, shell Run, profile, plugin diagrams, browser, Unzip, every other
caller, inventory, scanner, audit apparatus, and P4-P8.

## Verification Notes

- Focused red process contracts failed only for missing `SystemRunner`; focused
  shell selection failed only because it inherited `*os.File` stdout.
- Focused process/shell and every relevant process caller package: pass.
- API/CLI and fresh subprocess compatibility: pass with reports outside the
  worktree.
- Complete tests, race, vet, pinned lint, build, complete preflight, Make and
  production-script contracts, and all four host acceptance flows: pass.
- Complete Make runs hit the previously observed partial-raw-log signal-fixture
  flake; the immediate standalone launcher run and final full rerun passed all
  62 controls.
- The repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass under `/private/tmp/ply-p360-hermetic.xS9Thp`.
- Focused audit: expected exit 1, four improved, two held, zero regressed, one
  non-comparable, zero dirty paths.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete process adapter code and
tests and every complete process double/caller, complete shell Git code/tests
and callers, the P3.60 shell Run contracts, P3.58 Maven stdout contracts, the
import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README before editing. Confirm branch, HEAD, clean status,
reciprocal archive links, launcher `--check`, and exact product commit
`8503601`.

Make one focused shell Git system-capability implementation commit, then the
normal separate continuity commit. Stop before `GitDirty`/`GitIsRepo`, Maven,
profile, plugin diagrams, browser, Unzip, clock/server work, mutation harnesses,
P4-P8, audit or scanner changes, inventory changes, publication, or
distribution. Do not push, merge, stash, revert, launch a successor, or remove
the worktree.
