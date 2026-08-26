# Quality Upgrade Handover

Generated: 2026-08-26T12:12:14+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.63 product implementation:
  `1bce06f2e0eac3f89dad530e004c89ba123666e7`.
- Preceding P3.62 product implementation:
  `b0d324a46617d4ca4ee12fe030f367e5094ea9a9`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 63 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, 9164e06, 4376e05,
9f56714, 8503601, 5b678ab, b0d324a, and 1bce06f in roadmap order. The separate
operational continuity implementation is 1b85711. The focused T15
audit-apparatus repair is ce736a2; neither operational commit changes a Go
quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T121214+0200-narrow-browser-launcher-system-capability.md`.
The P3.63 plugin-diagrams capability archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`798d847cc6c9bbd806660730cd7f4729b9145bf07ac8c843ffafd1896cfe875d`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 63

Only the three private plugin-diagrams production process-capability selections
changed at product commit `1bce06f`.
`systemPluginDiagramsExportDependencies()`,
`systemPluginDiagramsGraphvizDependencies()`, and
`systemPluginDiagramsOpenDependencies()` now carry the exact existing
`process.SystemRunner()` instead of `process.System()`, so their complete
process dependencies have nil `Stdout`. The Graphviz selector's complete exact
system filesystem dependency remains unchanged. Every execution helper,
command and filesystem field, return path, and other caller is unchanged.

The ignored synchronous `structurizr-cli export` request and error, workspace
and temporary-directory arguments, deletion-before-export and
discovery-after-export placement, and discovery result and error behavior are
unchanged. The exact synchronous `dot` request, ordered arguments, distinct
stdout and stderr buffers and bytes, direct process error,
no-write-after-process-error behavior, exact `0644` stdout write, ignored write
error, and process-then-write order are unchanged. The ignored synchronous
macOS `open` request and error, output path, placement after each successful
conversion, and continued iteration are unchanged. Safe zero dependencies,
non-empty populations, exact one-request attempts, absence of another request
or write, Cobra registration, flags, and public behavior remain intact.

Four strengthened focused top-level contracts keep the suite at 388 tests.
They prove exact system-runner identity with nil process-dependency stdout for
all three selectors, the unchanged exact system filesystem dependency,
complete requests, arguments, streams, bytes, writes, modes, attempts and
sequences, ignored and direct errors, safe zero behavior, complete caller-owned
dependency preservation, non-empty populations, and the absence of another
request or write. Three recording doubles now preserve complete caller-owned
process and filesystem dependency values. They launch no external program,
touch no network, and write no repository fixture.

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
Q3.9, and all 228 numeric debt leaves. Product move 63 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `1bce06f` reports:

- Measured commit tree: `ebd29be397a8fbdda257ad94e3c8863fbaae84fa`.
- Structured scorecard SHA-256:
  `5cf5c2c171749e42dee5dfea6f5986b09dce96d00510560f4e625c7a6e4e3643`.
- Absolute L0: 8 of 8.
- 388 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 8 direct external sites outside five declared adapters of 30 production
  effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 81 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The three plugin-diagrams selectors feed one Structurizr caller composition.
The unchanged scanner no longer follows `os.Stdout` through that runner-only
composition, so exact Q1.3 improves from 9/31 to 8/30. The authoritative audit
output is
`target/quality-audit/scorecard.json`. Focused audit output is under
`/private/tmp/ply-p363-gate.O0jcm3/focused-audit`; pinned tools are under
`/private/tmp/ply-p358-tools`; the complete gate root is
`/private/tmp/ply-p363-gate.O0jcm3`. The disposable empty-HOME root is
`/private/tmp/ply-p363-hermetic.3faQvu`.

## Next Objective

Complete one focused P3 provenance move: narrow only the private browser-
launcher production selector's process dependency from `process.System()` to
the exact runner-only `process.SystemRunner()` so public `OpenBrowser` no
longer inherits the unused process-dependency `os.Stdout` capability.

Start red by strengthening the browser-launcher selector and recording
contracts. Prove exact system-runner identity, nil process-dependency stdout,
unchanged exact `runtime.GOOS`, complete caller-owned dependency preservation,
the exact asynchronous Linux, Windows, and macOS commands and URL arguments,
nil command streams, direct start errors, unsupported-platform error without a
process request, safe zero behavior, a non-empty population, and absence of
another request.

Then change only that private selector call. Preserve `openBrowser`, exported
`OpenBrowser`, every command field and return path, runtime platform selection,
its caller, server and clock behavior, plugin diagrams, Maven stdout, profile
editor, shell Run and Git, Unzip, every other process caller, inventory,
scanner, audit apparatus, and P4-P8.

## Verification Notes

- The three strengthened selector contracts first failed only because
  `Process.Stdout` was exact `os.Stdout` instead of nil.
- Focused plugin-diagrams/process/filesystem and every relevant process caller
  package: pass.
- API/CLI and fresh subprocess compatibility: pass with reports outside the
  worktree.
- Build, complete and uncached tests, race, vet, pinned lint, complete
  preflight, Make and production-script contracts, and all four host acceptance
  flows: pass.
- An initial launcher invocation and one complete preflight invocation hit the
  established nested partial-raw-log signal-fixture flake. Their immediate
  standalone or complete unchanged reruns passed all 62 controls. The first
  preflight also required the documented writable lint cache under the sandbox.
- Continuity validation hit the same signal fixture once in the outer run and
  once in its nested source-archive run; the final unchanged run passed all 62
  controls.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass under `/private/tmp/ply-p363-hermetic.3faQvu`.
- The documented comma-list focused audit measured all seven requested
  criteria, exited 1 as expected, and reported zero regressions.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete process and filesystem
adapter code/tests and every complete process double/caller, complete browser-
launcher code and contracts and every caller, P3.63 plugin-diagrams contracts,
P3.62 profile-editor contracts, P3.61 shell Git contracts, P3.60 shell Run
contracts, P3.58 Maven stdout contracts, the import-aware scanner, API/CLI
contracts, and the T15 repair and baseline reproduction README before editing.
Confirm branch, HEAD, clean status, reciprocal archive links, launcher
`--check`, and exact product commit `1bce06f`.

Make one focused browser-launcher system-capability implementation commit, then
the normal separate continuity commit. Stop before server/clock work, plugin
diagrams, Maven, profile, shell Run or Git, Unzip, another process caller,
mutation harnesses, P4-P8, audit or scanner changes, inventory changes,
publication, or distribution. Do not push, merge, stash, revert, launch a
successor, or remove the worktree.
