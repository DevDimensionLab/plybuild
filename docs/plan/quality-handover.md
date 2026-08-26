# Quality Upgrade Handover

Generated: 2026-08-26T09:46:45+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.58 product implementation:
  `4376e05afa889a260f6805aff086714e27c7b62a`.
- Preceding P3.57 product implementation:
  `9164e066ace0ee936b6fd0c27e440d04d4e5e3b3`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 58 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, 9164e06, and 4376e05 in
roadmap order. The separate operational continuity implementation is 1b85711.
The focused T15 audit-apparatus repair is ce736a2; neither operational commit
changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T094645+0200-migrate-shell-unzip-entry-open.md`.
The P3.58 Maven archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`01d203669b16b811b43fa302783b41bb66c7866f86ed30749884fa3f1a4d3876`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 58

Only Maven command's conditional standard-output capability selection changed
at product commit `4376e05`. `process.Dependencies` carries its runner and
standard-output writer together. `process.System()` selects the existing system
runner and exact `os.Stdout` identity together. The narrow `process.Stdout`
helper returns the injected writer only when enabled and nil otherwise.

`pkg/maven/command.go` retains the global logrus level predicate but routes the
selection through the complete process dependency instead of calling
`logger.StdOut()`. Every command retains its exact executable and ordered
argument bytes, `project.Path` directory, nil stdin and stderr, synchronous
Start false, single process attempt, exact dependency error, and exact log text
before execution. Debug and trace retain exact production `os.Stdout`; panic,
fatal, error, warn, and info retain nil stdout. The zero dependency remains a
safe no-op with nil output.

Five focused contracts bring the suite to 380 tests. They prove exact system
`os.Stdout`, enabled injected-writer identity, disabled and zero-dependency nil
output, all global logrus level states, the complete Maven command, one request,
exact error, complete dependency preservation, exact log ordering, rejection
of an empty population, and absence of another process request. Relevant
process recorders start from the complete system dependency and replace only
the runner. The tests launch no external program, touch no network, and write
no repository fixture. `maven.RunOn`, its callback, `Repository`, `pkg/logger`,
plugin diagrams, shell Unzip, public API/CLI, scanner, inventory, audit
apparatus, and every completed effect are unchanged.

## T15 Apparatus Checkpoint

The P3.54 T15 repair remains unchanged. Separate verified execution and
structured replicas preserve all historical test effects while measuring a
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
Q3.9, and all 228 numeric debt leaves. Product move 58 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `4376e05` reports:

- Measured commit tree: `550e10792b81c5b1fa3bc4714a451a1afb586541`.
- Structured scorecard SHA-256:
  `09cb1f364cbf55243f557d77c400b4befa7672d5b0a9378ad9e7d7cd6447ce32`.
- Absolute L0: 8 of 8.
- 380 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 15 direct external sites outside five declared adapters of 37
  production effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 76 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The Maven caller-side effect moved behind the process dependency. The unchanged
import-aware scanner now follows the new `os.Stdout` capability through the
existing `process.System()` selections, so the exact Q1.3 population changed
from 7/29 to 15/37 without a scanner, inventory, or comparable-ratchet change.

The authoritative audit output is `target/quality-audit/scorecard.json`.
Focused audit output is under `/private/tmp/ply-p358-focused-audit`; pinned
verification tools, caches, and compatibility reports are under `/private/tmp`.
The disposable empty-HOME root remains
`/private/tmp/ply-p358-hermetic.HllBTL`. The exact P3.57 Q1.3 reproduction is
under `/private/tmp/ply-p357-recheck.Zt2epV`.

## Next Objective

Complete one focused P3 production-effect move: route only shell Unzip's direct
archive-entry `zip.File.Open()` selection in `pkg/shell/command.go` through the
existing filesystem dependency.

Start red with filesystem-adapter and shell recording contracts for exact
archive-entry identity, exact returned reader and error, safe zero dependency,
the complete successful file-entry sequence, failure placement after the
destination open, suppression of copy/close/later traversal on open failure,
complete dependency preservation, a non-empty entry-open population, and no
unrelated filesystem request. Use only guarded temporary archives outside the
repository.

Then add only a narrow archive-entry open operation to
`internal/adapter/filesystem` and replace only `f.Open()` in shell Unzip.
Preserve source/destination, archive open and deferred close, iteration, log and
zip-slip order, partial filenames, directory/parent creation, output flags and
modes, ignored copy error, output/entry close order, exact errors, public API,
and safe-zero behavior. Leave process/Maven, every other Unzip operation,
another caller, inventory, scanner, audit apparatus, clock/server, mutation
harnesses, and P4-P8 unchanged.

## Verification Notes

- Focused red contracts failed only for the missing process stdout dependency
  capability and helper; focused green tests use recording boundaries.
- Focused process/Maven/logger and relevant caller tests: pass.
- API/CLI and fresh subprocess compatibility: pass with pinned tools and
  generated reports outside the worktree.
- Complete and uncached tests, race, vet, pinned lint, build, complete preflight,
  Make contracts, and all four host acceptance flows: pass.
- Complete preflight passes with pinned tools on PATH and the established Maven
  mutation contract names preserved.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass under `/private/tmp/ply-p358-hermetic.HllBTL`.
- Focused audit: expected exit 1, four improved, two held, zero regressed, one
  non-comparable, zero dirty paths.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete shell Unzip code/tests and
callers, the filesystem adapter and every complete filesystem double/caller,
the import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README before editing. Confirm branch, HEAD, clean status,
reciprocal archive links, launcher `--check`, and exact product commit
`4376e05`.

Make one focused shell Unzip entry-open implementation commit, then the normal
separate continuity commit. Stop before another Unzip operation, process/Maven,
plugin diagrams, clock/server work, mutation harnesses, P4-P8, audit or scanner
changes, inventory changes, publication, or distribution. Do not push, merge,
stash, revert, launch a successor, or remove the worktree.
