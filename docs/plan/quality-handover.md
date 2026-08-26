# Quality Upgrade Handover

Generated: 2026-08-26T11:45:25+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.62 product implementation:
  `b0d324a46617d4ca4ee12fe030f367e5094ea9a9`.
- Preceding P3.61 product implementation:
  `5b678ab676192b835bb6a3a35cfbee4eced3105e`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 62 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, 9164e06, 4376e05,
9f56714, 8503601, 5b678ab, and b0d324a in roadmap order. The separate
operational continuity implementation is 1b85711. The focused T15
audit-apparatus repair is ce736a2; neither operational commit changes a Go
quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T114525+0200-narrow-plugin-diagrams-system-capabilities.md`.
The P3.62 profile-editor capability archive is answered history and links
reciprocally to that archive. The graph has exactly one NEXT tail. The archived
prompt SHA-256 is
`3c4d8c36dc53e9af1a1006043655afe6497d9f54e4bc9a393a41f460b0179bf9`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 62

Only the private profile-editor production process-capability selection changed
at product commit `b0d324a`. `systemProfileEditorDependencies()` now carries the
exact existing `process.SystemRunner()` instead of `process.System()`, so its
complete process dependency has nil `Stdout`. Its separate command-stream
fields remain exact `os.Stdin` and `os.Stdout`. `runProfileEditor`, profile
commands and configuration, Maven, shell Run and Git, plugin diagrams, browser,
Unzip, every other caller, and both process selectors are unchanged.

The arbitrary editor executable, single exact config-path argument, empty
directory, dependency-selected stdin and stdout, nil stderr, synchronous
execution, exact one-request population, direct runner error, safe zero
dependency, `EDITOR` selection, `vim` default, Cobra registration and flags,
later profile operations and errors, and every public path remain unchanged.
There is no fallback, retry, wrapping, additional logging, execution, cleanup,
or global-state change.

Four strengthened focused top-level contracts keep the suite at 388 tests. They
prove exact system-runner identity with nil process-dependency stdout, separate
exact production terminal streams, complete arbitrary editor/config-path
requests and injected stream identities, direct errors, safe zero behavior,
complete caller-owned dependency preservation, a non-empty one-request
population, and absence of another request. The recorder launches no external
program, touches no network, and writes no repository fixture.

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
Q3.9, and all 228 numeric debt leaves. Product move 62 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `b0d324a` reports:

- Measured commit tree: `86b97ce8e9cbdbc2259a21a28e87c3678552c06d`.
- Structured scorecard SHA-256:
  `e9d5e45b477cde1791a0e73b0d5cf1332036fec371ccf29f4d9b0827a838ab46`.
- Absolute L0: 8 of 8.
- 388 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 9 direct external sites outside five declared adapters of 31 production
  effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 80 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The unchanged scanner no longer follows `os.Stdout` through the private
profile-editor runner-only production dependency, so exact Q1.3 improves from
10/32 to 9/31. The authoritative audit output is
`target/quality-audit/scorecard.json`. Focused audit output is under
`/private/tmp/ply-p362-gate.yUjJr4/focused-audit`; pinned tools are under
`/private/tmp/ply-p358-tools`; the complete gate root is
`/private/tmp/ply-p362-gate.yUjJr4`. The disposable empty-HOME root is
`/private/tmp/ply-p362-hermetic.77JaRw`.

## Next Objective

Complete one focused P3 provenance move: narrow only the three existing private
plugin-diagrams production selectors' process dependencies from
`process.System()` to the exact runner-only `process.SystemRunner()` so their
one Structurizr caller composition no longer inherits the unused `os.Stdout`
capability.

Start red by strengthening the export, Graphviz, and open selector/recording
contracts. Prove exact system-runner identity and nil process-dependency stdout
for all three selectors; preserve the Graphviz selector's complete exact system
filesystem dependency; preserve complete caller-owned dependency fields; every
exact command, argument, stream, byte, write, mode, request/write count and
sequence; ignored export/open errors; direct Graphviz process error; ignored
Graphviz write error; safe zero behavior; non-empty populations; and absence of
another request or write.

Then change only those three private selector calls. Preserve
`runStructurizrDiagrams`, `convertStructurizrDiagram`,
`openStructurizrDiagram`, deletion/discovery/iteration and error placement,
public command behavior, Maven stdout, profile editor, shell Run and Git,
browser, Unzip, every other process/filesystem caller, inventory, scanner,
audit apparatus, and P4-P8.

## Verification Notes

- The strengthened profile-editor selector contract first failed only because
  `Process.Stdout` was exact `os.Stdout` instead of nil.
- Focused profile/process and every relevant process caller package: pass.
- API/CLI and fresh subprocess compatibility: pass with reports outside the
  worktree.
- Build, complete and uncached tests, race, vet, pinned lint, complete
  preflight, Make and production-script contracts, and all four host acceptance
  flows: pass.
- The first `make test` attempt hit the established nested partial-raw-log
  signal-fixture flake. The immediate standalone 62-control launcher run and
  the complete `make test` rerun both passed.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass under `/private/tmp/ply-p362-hermetic.77JaRw`.
- The first focused audit invocation repeated `--only`, whose last-one-wins
  semantics measured only Q3.4 and exited 0. The corrected documented comma-
  list invocation measured all seven requested criteria and exited 1 as
  expected, with zero regressions.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete process and filesystem
adapter code/tests and every complete process double/caller, complete
plugin-diagrams command code and export/Graphviz/open contracts and callers,
P3.62 profile-editor contracts, P3.61 shell Git contracts, P3.60 shell Run
contracts, P3.58 Maven stdout contracts, the import-aware scanner, API/CLI
contracts, and the T15 repair and baseline reproduction README before editing.
Confirm branch, HEAD, clean status, reciprocal archive links, launcher
`--check`, and exact product commit `b0d324a`.

Make one focused plugin-diagrams system-capability implementation commit, then
the normal separate continuity commit. Stop before browser, Maven, profile,
shell Run or Git, Unzip, another process/filesystem caller, clock/server work,
mutation harnesses, P4-P8, audit or scanner changes, inventory changes,
publication, or distribution. Do not push, merge, stash, revert, launch a
successor, or remove the worktree.
