# Quality Upgrade Handover

Generated: 2026-08-26T13:44:55+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.66 product implementation:
  `c627a3ec0ec8577f4bd2aef8cd28be7ab970b5e0`.
- Preceding P3.65 product implementation:
  `40cece5337359d313a2fda3200c9d1fab3850c44`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 66 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, 9164e06, 4376e05,
9f56714, 8503601, 5b678ab, b0d324a, 1bce06f, 0b96f10, 40cece5, and
c627a3e in roadmap order. The separate operational continuity implementation
is 1b85711. The focused T15 audit-apparatus repair is ce736a2; neither
operational commit changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P4 is active and P5-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T134455+0200-migrate-spring-archive-clock.md`.
The P3.66 Maven standard-output provenance archive is answered history and
links reciprocally to that archive. The graph has exactly one NEXT tail.
The archived prompt SHA-256 is
`d7555a234bd2a57b20e3074617a6cf2b5bf36e644e0b1e7881500614dcd52662`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 66 And Exit

Only the process adapter's standard-output forwarding surface and Maven's
private production dependency selection changed at product commit `c627a3e`.
New `process.SystemStdout` returns the exact existing
`Stdout(System(), true)` composition. It performs no fallback, wrapping,
buffering, copying, logging, cleanup, retry, or global-state mutation. The
complete process dependency, runner, system selectors, execution path, and
implementation remain unchanged.

Private `systemRunOnDependencies` carries the exact existing
`process.SystemRunner()` and that exact system standard output. Public
`maven.RunOn` replaces only its complete `process.System()` selection with that
private selector; private `runOn`, every callback and caller, and the return
path are unchanged. Every arbitrary command and argument byte and order,
repository value, exact project directory, info log before execution,
conditional stdout selection for debug and trace only, nil stdin and stderr,
synchronous one-attempt execution, and direct process error remain exact.

Strengthened focused contracts keep the suite at 388 tests. They prove exact
`os.Stdout` identity through the existing process composition, exact system
runner identity, complete caller-owned dependency delivery and preservation,
every log and command field, the full log-level matrix, callback construction
and invocation timing, one exact attempt, direct error, safe zero behavior,
non-empty populations, exported public composition, and absence of another
process operation. They launch no external program, touch no network, and
write no repository fixture.

P3 is complete: process exits outside `main` remain zero, its production-effect
moves and recording contracts are complete, migrated call sites have left
Q1.3, and the CLI/API contracts remain compatible. P4 is active.

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
Q3.9, and all 228 numeric debt leaves. Product move 66 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `c627a3e` reports:

- Measured commit tree: `fc8ae247df072a17c61d3e68a8359dfb4f1d7765`.
- Structured scorecard SHA-256:
  `1023a95ee7e4d7a10e95ae647821892dae71b8815f35062ee5e4b8cbd9c207d0`.
- Absolute L0: 8 of 8.
- 388 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 4 direct external sites outside five declared adapters of 27 production
  effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 84 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The unchanged scanner no longer follows concrete `os.Stdout` provenance
through Maven's complete public `RunOn` dependency, so exact Q1.3 improves from
5/27 to 4/27. The authoritative audit output is
`target/quality-audit/scorecard.json`. Focused audit output is under
`/private/tmp/ply-p366-focused-audit`; pinned tools are under
`/private/tmp/ply-p358-tools`. The disposable empty-HOME root is
`/private/tmp/ply-p366-hermetic.bpHLfz`.

## Next Objective

Begin P4 with one focused clock move: introduce the declared
`internal/adapter/clock` boundary and migrate only Spring's private archive-path
timestamp selection from direct `time.Now().Unix()` to that boundary. Preserve
the exact current working-directory request and error short circuit, one
timestamp selection only after that success, exact Unix-second conversion,
path formatting, return values, callers, and public behavior.

Start red with focused clock-adapter and Spring archive-path recording
contracts. Prove safe zero behavior, exact system-time identity within a
bounded before/after observation, complete caller-owned dependency preservation,
working-directory-before-clock order, absence of a clock request after a
directory error, one exact timestamp request after success, arbitrary injected
times including pre-epoch and subsecond values, exact path bytes, direct error,
non-empty populations, public composition, and absence of another operation.

Then add only the narrow clock adapter and the private Spring dependency field
and selection needed for that one flow. Leave Kibana sleep, web server start and
shutdown, process/HTTP/filesystem adapters, every completed P3 flow, inventory,
scanner, audit apparatus, mutation harnesses, and the rest of P4-P8 unchanged.

## Verification Notes

- The strengthened contracts first failed because `SystemStdout`,
  `systemRunOnDependencies`, and the Maven-visible helper did not exist.
- Focused process/Maven and relevant caller packages: pass.
- API/CLI and fresh subprocess compatibility: pass with reports outside the
  worktree.
- Complete and uncached tests, race, vet, pinned lint, complete preflight, Make
  and production-script contracts, and all four host acceptance flows: pass.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Two complete Make attempts hit the established nested partial-raw-log signal
  fixture flake; the unchanged complete `make test` rerun passed.
- Empty-HOME count-2: pass under `/private/tmp/ply-p366-hermetic.bpHLfz`.
- The focused audit measured all seven requested criteria, exited 1 as expected,
  held Q0.6 at 26 guarded sites, improved exact Q1.3 to 4 of 27, and reported
  zero comparable regressions. Its scorecard SHA-256 is
  `29d71651ef7cac78507cfe837115c20a8f5b37efed3d939ffbc14f0acb40ef17`.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P4 entry and checkpoint gate,
both design documents, `.quality/inventory`, the complete clock-related Spring
archive-path code/contracts/callers, both remaining direct clock sites, every
adapter's source/tests and complete dependency-preserving doubles, the
import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README before editing. Confirm branch, HEAD, clean status,
reciprocal archive links, launcher `--check`, and exact product commit
`c627a3e`.

Make one focused Spring archive-clock implementation commit, then the normal
separate continuity commit. Stop before Kibana sleep, server work, another
Spring flow, process/HTTP/filesystem changes, inventory or scanner changes,
mutation harnesses, later P4 work, P5-P8, publication, or distribution. Do not
push, merge, stash, revert, launch a successor, or remove the worktree.
