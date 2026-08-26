# Quality Upgrade Handover

Generated: 2026-08-26T12:40:26+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.64 product implementation:
  `0b96f10073c771d8dd86360b916dd296220cebbe`.
- Preceding P3.63 product implementation:
  `1bce06f2e0eac3f89dad530e004c89ba123666e7`.
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
9f56714, 8503601, 5b678ab, b0d324a, 1bce06f, and 0b96f10 in roadmap order. The
separate operational continuity implementation is 1b85711. The focused T15
audit-apparatus repair is ce736a2; neither operational commit changes a Go
quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T124026+0200-narrow-unzip-output-file-provenance.md`.
The P3.64 browser-launcher capability archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`ce9587d3ed287014fd7012c23b620583b619ab5eb375c69d72fcfae79fddfc88`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 64

Only the private browser-launcher production process-capability selection
changed at product commit `0b96f10`.
`systemBrowserLauncherDependencies()` now carries the exact existing
`process.SystemRunner()` instead of `process.System()`, so its complete process
dependency has nil `Stdout`. Its separate exact `runtime.GOOS` selection,
`openBrowser`, exported `OpenBrowser`, its caller, command fields, return path,
and every other process caller remain unchanged.

The exact asynchronous `xdg-open` Linux request, `rundll32` Windows request
with ordered `url.dll,FileProtocolHandler` and URL arguments, and `open` macOS
request are unchanged. Each arbitrary URL byte, empty directory, nil command
streams, `Start: true`, direct process start error, one exact attempt, absence
of another request, unsupported-platform error without a process attempt, safe
zero dependency, and public behavior remain intact. There is no fallback,
retry, wrapping, logging, synchronous execution, cleanup, or global-state
change.

Five strengthened focused top-level contracts keep the suite at 388 tests.
They prove exact system-runner identity with nil process-dependency stdout,
unchanged exact runtime platform selection, complete platform requests, URL
arguments, stream identities, asynchronous attempts, direct errors, safe zero
behavior, complete caller-owned dependency preservation, non-empty population,
absence of another request, and the exported wrapper's exact composition. The
recording double now preserves the complete caller-owned browser dependency.
The contracts launch no external program, touch no network, and write no
repository fixture.

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
Q3.9, and all 228 numeric debt leaves. Product move 64 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `0b96f10` reports:

- Measured commit tree: `8a7180aa20082d504a47ce8bdb4a12d405c9bf50`.
- Structured scorecard SHA-256:
  `f817974080ab4b1641155e24e7d2b7f4d0fddf1ae7c06d7c4922320a70f5a842`.
- Absolute L0: 8 of 8.
- 388 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 7 direct external sites outside five declared adapters of 29 production
  effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 81 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The unchanged scanner no longer follows `os.Stdout` through the browser
launcher's runner-only composition, so exact Q1.3 improves from 8/30 to 7/29.
The authoritative audit output is
`target/quality-audit/scorecard.json`. Focused audit output is under
`/private/tmp/ply-p364-gate.UjofWr/focused-audit`; pinned tools are under
`/private/tmp/ply-p358-tools`; the complete gate root is
`/private/tmp/ply-p364-gate.UjofWr`. The disposable empty-HOME root is
`/private/tmp/ply-p364-hermetic.0k9qo3`.

## Next Objective

Complete one focused P3 provenance move: add one narrow filesystem adapter
helper that returns the existing `File` interface while delegating to the
complete existing `FileSystem.OpenFile` operation, then use it only for shell
Unzip's destination-file open. This removes concrete `*os.File` provenance
from the following adapter `Copy` and `Close` calls without changing their
execution.

Start red with focused filesystem and Unzip recording contracts. Prove exact
complete path, flags, mode, returned file identity and error, safe zero
behavior, non-empty population, complete caller-owned dependency preservation,
and the complete existing archive open, path validation, directory, entry,
copy, close, partial-result, error-precedence, and traversal behavior.

Then add only the interface-returning helper and replace only Unzip's one
destination `filesystem.OpenFile` call. Preserve the existing complete
`FileSystem` interface and system implementation, exported `filesystem.OpenFile`
and public `file.OpenFile`, every flag and mode, the same open/copy/close
operations and order, Maven, browser, plugin diagrams, every process caller,
clock/server work, inventory, scanner, audit apparatus, and P4-P8.

## Verification Notes

- The strengthened selector contract first failed only because `Process.Stdout`
  was exact `os.Stdout` instead of nil.
- Focused browser-launcher/process and every relevant caller package: pass.
- API/CLI and fresh subprocess compatibility: pass with reports outside the
  worktree.
- Build, complete and uncached tests, race, vet, pinned lint, complete
  preflight, Make and production-script contracts, and all four host acceptance
  flows: pass.
- One complete preflight invocation hit the established nested partial-raw-log
  signal-fixture flake. Its immediate complete unchanged rerun passed all 62
  controls. Compatibility used the exact pinned `apidiff` path after the first
  shell environment lacked that executable on `PATH`.
- Continuity validation hit the same signal fixture once; its immediate
  unchanged rerun passed all 62 controls.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass under `/private/tmp/ply-p364-hermetic.0k9qo3`.
- The documented comma-list focused audit measured all seven requested
  criteria, exited 1 as expected, and reported zero regressions.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete filesystem adapter
code/tests and every complete filesystem double/caller, complete shell Unzip
code and contracts and every caller, P3.64 browser-launcher contracts, P3.63
plugin-diagrams contracts, P3.62 profile-editor contracts, P3.61 shell Git
contracts, P3.60 shell Run contracts, P3.58 Maven stdout contracts, the
import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README before editing. Confirm branch, HEAD, clean status,
reciprocal archive links, launcher `--check`, and exact product commit
`0b96f10`.

Make one focused Unzip destination-file provenance implementation commit, then
the normal separate continuity commit. Stop before server/clock work, Maven,
browser, plugin diagrams, profile, shell Run or Git, another process caller,
other Unzip changes, mutation harnesses, P4-P8, audit or scanner changes,
inventory changes, publication, or distribution. Do not push, merge, stash,
revert, launch a successor, or remove the worktree.
