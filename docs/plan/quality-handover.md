# Quality Upgrade Handover

Generated: 2026-08-26T13:08:43+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.65 product implementation:
  `40cece5337359d313a2fda3200c9d1fab3850c44`.
- Preceding P3.64 product implementation:
  `0b96f10073c771d8dd86360b916dd296220cebbe`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 64 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, 9164e06, 4376e05,
9f56714, 8503601, 5b678ab, b0d324a, 1bce06f, 0b96f10, and 40cece5 in roadmap
order. The separate operational continuity implementation is 1b85711. The
focused T15 audit-apparatus repair is ce736a2; neither operational commit
changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T130843+0200-narrow-maven-standard-output-provenance.md`.
The P3.65 Unzip output-file provenance archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`3c880bbf28373de5dd123a7cd01ae045fdef9518c6520002ea67b66d1e569ce0`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 65

Only the private filesystem adapter forwarding surface and shell Unzip's one
destination-file open call changed at product commit `40cece5`.
`filesystem.OpenFileAsFile` delegates exactly once to the complete existing
`FileSystem.OpenFile` operation with the exact path, flags, and mode while
returning its exact file through the existing `filesystem.File` interface and
preserving the exact error. The complete `FileSystem` interface, system
implementation, exported `filesystem.OpenFile`, public `file.OpenFile`, and
every other operation and caller remain unchanged.

`unzipWithDependencies` uses that helper only for its destination open. Exact
source and destination bytes, archive open and deferred close, zip-slip
boundary and error, directory metadata classification, joined paths, recursive
directories and modes, ordered filenames, destination flags, entry modes and
opens, ignored copy results, output and reader closes, attempts, operation
order, direct errors, partial filenames, error precedence, traversal, safe zero
behavior, system selection, exported `Unzip`, and its caller remain intact.
There is no selection, fallback, retry, wrapping, logging, cleanup, or global-
state change.

Two strengthened focused top-level contracts keep the suite at 388 tests. They
prove safe zero behavior, exact complete dependency delivery and preservation,
path, flags, mode, returned `filesystem.File` identity, exact error, non-empty
population, absence of another operation, the exact private helper placement,
and every existing Unzip archive, path, directory, entry, copy, close, order,
partial-result, error-precedence, and traversal guarantee. The contracts launch
no external program, touch no network, and write no repository fixture.

## T15 Apparatus Checkpoint

The P3.54 T15 repair remains unchanged. Separate verified execution and
structured replicas preserve historical test effects while measuring a
pristine exact `5635d50` checkout plus the declared inventory overlay. No output
path is excluded or deleted, and every effect stays in the retained manifest.
The apparatus-owned reproduction HOME includes the required directory-shaped
`.co-pilot/profiles/.active_profile` fixture with marker SHA-256
`cb95f24c35d3987f8aba51231aade19580ffe9242324804fcad2ccff350d1c9a`.

The complete 15-control proof passes and reproduces old scorecard
`d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`,
new scorecard
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`,
normalized stored raw body
`cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f`,
Q3.9, and all 228 numeric debt leaves. Product move 65 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `40cece5` reports:

- Measured commit tree: `c1009903d65b062786604f71d8f1fa89d2dfe61d`.
- Structured scorecard SHA-256:
  `1bae6abec78d2306c2205bdf1f96b7bffdabcef7e82a232a2aedafe536f3bf56`.
- Absolute L0: 8 of 8.
- 388 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 5 direct external sites outside five declared adapters of 27 production
  effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 83 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The unchanged scanner no longer follows concrete `*os.File` provenance through
Unzip's adapter `Copy` and `Close` requests, so exact Q1.3 improves from 7/29 to
5/27.
The authoritative audit output is
`target/quality-audit/scorecard.json`. Focused audit output is under
`/private/tmp/ply-p365-focused-audit`; pinned tools are under
`/private/tmp/ply-p358-tools`. The disposable empty-HOME root is
`/private/tmp/ply-p365-hermetic.sWqcSL`.

## Next Objective

Complete one focused P3 provenance move: retain Maven's exact debug-only system
standard output without giving public `RunOn` concrete `os.Stdout` provenance.
Add one narrow process-adapter helper for the exact existing
`Stdout(System(), true)` composition, give public `RunOn` a private production
selector carrying the exact existing system runner and standard output, and
leave every other process composition unchanged.

Start red with focused process-standard-output and Maven recording contracts.
Prove exact standard-output identity for enabled and disabled capture, exact
system runner identity with nil runner-owned stdout, safe zero behavior,
non-empty production population, complete caller-owned dependency preservation,
and the complete existing debug/non-debug command, logging, callback, error,
attempt, and public-wrapper behavior.

Then add only the narrow standard-output helper and Maven selector and replace
only public `RunOn`'s production dependency composition. Preserve the complete
process adapter interface and system implementation, exported `process.System`,
`SystemRunner`, and `Stdout`, private `runOnWithDependencies`, every command,
stream, callback, error, attempt, caller, return path, Unzip, browser, plugin
diagrams, profile, shell Run and Git, server/clock work, inventory, scanner,
audit apparatus, and P4-P8.

## Verification Notes

- The strengthened filesystem-open contract first failed because
  `OpenFileAsFile` did not exist; the strengthened Unzip placement contract
  found zero interface-returning calls and the one existing concrete call.
- Focused filesystem, shell Unzip, and relevant Spring caller packages: pass.
- API/CLI and fresh subprocess compatibility: pass with reports outside the
  worktree.
- Build, complete and uncached tests, race, vet, pinned lint, complete
  preflight, Make and production-script contracts, and all four host acceptance
  flows: pass.
- The complete preflight passed on its first invocation using exact pinned
  compatibility and lint tools outside the worktree.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass under `/private/tmp/ply-p365-hermetic.sWqcSL`.
- The documented comma-list focused audit measured all seven requested
  criteria, exited 1 as expected, held Q0.6 at 26 guarded sites, improved exact
  Q1.3 to 5 of 27, and reported zero comparable regressions.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete process adapter code/tests
and every complete process double/caller, complete Maven RunOn code and
contracts and every caller, P3.65 Unzip contracts, P3.64 browser-launcher
contracts, P3.63 plugin-diagrams contracts, P3.62 profile-editor contracts,
P3.61 shell Git contracts, P3.60 shell Run contracts, the import-aware scanner,
API/CLI contracts, and the T15 repair and baseline reproduction README before
editing. Confirm branch, HEAD, clean status, reciprocal archive links, launcher
`--check`, and exact product commit `40cece5`.

Make one focused Maven standard-output provenance implementation commit, then
the normal separate continuity commit. Stop before server/clock work, Unzip,
browser, plugin diagrams, profile, shell Run or Git, another process caller,
other Maven changes, mutation harnesses, P4-P8, audit or scanner changes,
inventory changes, publication, or distribution. Do not push, merge, stash,
revert, launch a successor, or remove the worktree.
