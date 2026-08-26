# Quality Upgrade Handover

Generated: 2026-08-26T10:17:55+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.59 product implementation:
  `9f56714234ef4aa07cd04b5ebbb4994d388ade4c`.
- Preceding P3.58 product implementation:
  `4376e05afa889a260f6805aff086714e27c7b62a`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 59 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, 9164e06, 4376e05, and
9f56714 in roadmap order. The separate operational continuity implementation
is 1b85711. The focused T15 audit-apparatus repair is ce736a2; neither
operational commit changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T101755+0200-narrow-shell-run-system-capability.md`.
The P3.59 shell Unzip entry-open archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`3e150735be39600469a295d93fe7c9167c717b27270b4c741dc48e4314e1fba7`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 59

Only shell Unzip's direct archive-entry reader selection changed at product
commit `9f56714`. The existing complete filesystem dependency now includes
`OpenZipEntry(*zip.File) (io.ReadCloser, error)`. Its public forwarding helper
returns the established `ErrNoFilesystem` for the zero dependency; its system
implementation invokes `Open` on the exact supplied archive entry. Private
`unzipWithDependencies` replaces only direct `f.Open()` with this helper.

The exported `Unzip` signature, archive source and destination, single archive
open, exact deferred archive-close identity and ignored result, traversal and
entry order, log text and placement, zip-slip validation, partial-filename
append point, directory handling, parent creation, destination flags and mode,
copy and close behavior, and every established return path remain unchanged.
Each file still opens its destination before one entry open, copies once from
the exact returned reader with ignored count and error, closes the output before
the entry reader, and suppresses later work at the same failure points. The
zero dependency performs no archive or entry open and no filesystem mutation.

Four focused top-level contracts bring the suite to 384 tests. The filesystem
contracts prove exact archive-entry identity, exact returned reader and error,
safe zero behavior, system delegation, and a non-empty population. The shell
recorder proves the complete successful file sequence, exact reader identity in
copy and close, entry-open failure after destination open, no copy, close, or
later traversal after that failure, complete dependency preservation, rejection
of an empty entry-open population, and absence of another filesystem request.
Every complete filesystem double and caller mechanically preserves the extended
dependency. Focused tests use guarded temporary archives, launch no process,
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
Q3.9, and all 228 numeric debt leaves. Product move 59 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `9f56714` reports:

- Measured commit tree: `18f303931958d94ba72b9ed861faa24b94e29925`.
- Structured scorecard SHA-256:
  `41b51871f6436851613af04c86af93b6fb7f7f0e7f74833c47d8e3d45ddfb10e`.
- Absolute L0: 8 of 8.
- 384 test functions, zero skipped; 19 of 25 packages have tests.
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

The unchanged scanner does not catalog `archive/zip.File.Open`, so exact Q1.3
holds at 15/37 even though the direct selection moved behind the filesystem
dependency. The authoritative audit output is
`target/quality-audit/scorecard.json`. Focused audit output is under
`/private/tmp/ply-p359-focused-audit.cpukOv`; pinned tools and caches are under
`/private/tmp/ply-p358-tools` and `/private/tmp/ply-p359-*`. The disposable
empty-HOME root is `/private/tmp/ply-p359-hermetic.QTiEca`.

## Next Objective

Complete one focused P3 provenance move: give public shell `Run` a runner-only
system process dependency so `systemRunDependencies()` no longer inherits the
unused standard-output capability added in P3.58.

Start red with process-adapter and shell recording contracts for exact system
runner identity, nil stdout on the narrow selector, unchanged exact `os.Stdout`
on `process.System()`, the complete arbitrary command and stdout/stderr buffer
identities, exact error, safe zero behavior, complete dependency preservation,
a non-empty request population, and no second process request.

Then add only the narrow runner-only production selector and use it only in
`systemRunDependencies()`. Preserve `Run`, exact command bytes, logs, buffer
identities, synchronous execution, error and output behavior. Leave
`process.System()`, Maven stdout, Git wrappers, profile, plugin diagrams,
browser launching, Unzip, every other caller, public API/CLI, inventory,
scanner, audit apparatus, clock/server, mutation harnesses, and P4-P8
unchanged.

## Verification Notes

- Focused red adapter contracts failed to compile only for missing
  `OpenZipEntry`; focused shell contracts failed because direct `f.Open()`
  ignored the injected reader/error and recorded no entry-open population.
- Focused filesystem/shell and all relevant caller packages: pass.
- API/CLI and fresh subprocess compatibility: pass with generated reports
  outside the worktree.
- Complete tests, race, vet, pinned lint, build, complete preflight, Make
  contracts, and all four host acceptance flows: pass.
- The initial preflight used the sandbox-blocked default golangci-lint cache;
  the isolated-cache rerun passed with zero lint issues.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass under `/private/tmp/ply-p359-hermetic.QTiEca`.
- Focused audit: expected exit 1, four improved, two held, zero regressed, one
  non-comparable, zero dirty paths.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete process adapter code and
tests and every complete process double/caller, complete shell Run code/tests
and callers, P3.58 Maven stdout contracts, the import-aware scanner, API/CLI
contracts, and the T15 repair and baseline reproduction README before editing.
Confirm branch, HEAD, clean status, reciprocal archive links, launcher
`--check`, and exact product commit `9f56714`.

Make one focused shell Run system-capability implementation commit, then the
normal separate continuity commit. Stop before another process caller, Git,
Maven, profile, plugin diagrams, browser, Unzip, clock/server work, mutation
harnesses, P4-P8, audit or scanner changes, inventory changes, publication, or
distribution. Do not push, merge, stash, revert, launch a successor, or remove
the worktree.
