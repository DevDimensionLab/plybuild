# Quality Upgrade Handover

Generated: 2026-08-26T09:13:01+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.57 product implementation:
  `9164e066ace0ee936b6fd0c27e440d04d4e5e3b3`.
- Preceding P3.56 product implementation:
  `6f4bc3737ba07a2a4614617d272ef44558d39b77`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 57 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, 6f4bc37, and 9164e06 in roadmap
order. The separate operational continuity implementation is 1b85711. The
focused T15 audit-apparatus repair is ce736a2; neither operational commit
changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T091301+0200-migrate-maven-debug-stdout.md`.
The P3.57 Graphviz archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`58e0daf53a3bef445f0f1fff667e8c2223f4ceb47c39f44526366e524e0fd13d`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 57

Only the direct Graphviz `dot` conversion and its existing successful stdout
write in private plugin-diagrams flow changed at product commit `9164e06`.
`structurizrCmd.RunE` retains mandatory workspace lookup and passes the complete
private Graphviz dependency value beside the completed export and open values.
Production selects `process.System()` and `filesystem.System()` explicitly.

For each discovered dot file, `convertStructurizrDiagram` calls
`process.Execute` once with:

- executable `dot`;
- ordered arguments containing the input file and `-Tpng`;
- empty `Dir`, nil stdin, distinct stdout/stderr buffers, and `Start: false`.

The exact process error returns before any write or open. On success, the exact
captured stdout bytes are passed once to `filesystem.WriteFile` at the existing
output PNG path with mode `0644`; the exact write result remains ignored,
captured stderr remains discarded, and the completed ignored open helper still
follows. Successful conversion/open sequences continue to later files.
Workspace lookup, deletion, export, discovery, iteration order, output paths,
printed text, Cobra behavior, public API, `pkg/structurizr`, both adapters,
scanner, inventory, audit apparatus, and every completed effect are unchanged.

Six focused private contracts bring the suite to 375 tests. They prove complete
system dependency selection, exact process command and stream identities, one
request, process-error precedence, suppression of write/open after failure,
exact successful bytes/path/mode, one ignored-error write, following open
placement, continued iteration, safe zero dependencies, rejection of empty
process and write populations, and absence of export, open, or unrelated
process/write requests. They use recorders and source AST, launch no external
program, touch no network, and write no repository fixture.

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
Q3.9, and all 228 numeric debt leaves. Product move 57 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `9164e06` reports:

- Measured commit tree: `78fef4303e29d2b8d79434727bb8f7d9e0fd070a`.
- Structured scorecard SHA-256:
  `9f584097422a17fd6eea005dda97182e042aad05a34837f0461ce0933248f700`.
- Absolute L0: 8 of 8.
- 375 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 7 direct external sites outside five declared adapters of 29 production
  effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 75 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative audit output is `target/quality-audit/scorecard.json`.
Focused audit output is under `/private/tmp/ply-p357-focused-audit`; pinned
verification tools, caches, and compatibility reports are under `/private/tmp`.
The disposable empty-HOME root remains
`/private/tmp/ply-p357-hermetic.6jtS8S`.

## Next Objective

Complete one focused P3 production-effect move: route only Maven command's
conditional `logger.StdOut()` capability selection in `pkg/maven/command.go`
through the existing process dependency.

Start red with process-adapter and Maven recording contracts for exact
production `os.Stdout` selection, debug-enabled injected writer identity,
debug-disabled and zero-dependency nil output, complete commands in both level
states, one request, exact dependency error, complete dependency preservation,
and non-empty populations. Invoke no real external program.

Then add only a narrow standard-output capability and conditional forwarding
helper to `internal/adapter/process`; select the runner and `os.Stdout` together
in `process.System()`, and replace only the Maven caller-side
`logger.StdOut()` selection. Preserve command name, ordered arguments,
`project.Path`, nil stdin/stderr, synchronous execution, logging order, returned
error, exported callback shape, and safe-zero behavior. Leave `pkg/logger`,
plugin diagrams, shell unzip provenance, public API/CLI, inventory, scanner,
audit apparatus, clock/server, mutation harnesses, and P4-P8 unchanged.

## Verification Notes

- Focused red contracts failed only for the missing private Graphviz dependency
  and helper; focused green tests launch no real process.
- Focused cmd/process/filesystem/structurizr and relevant caller tests: pass.
- API/CLI and fresh subprocess compatibility: pass with pinned tools and
  generated reports outside the worktree.
- Complete and uncached tests, race, vet, pinned lint, build, complete preflight,
  Make contracts, and all four host acceptance flows: pass.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass; its disposable HOME remains under `/private/tmp`.
- Focused audit: expected exit 1, four improved, two held, zero regressed, one
  non-comparable, zero dirty paths.
- Full clean audit: expected exit 1, 15 findings, five improved, two held, zero
  regressed, one non-comparable, zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete Maven/logger code/tests
and callers, the process adapter and every complete process double/caller, the
import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README before editing. Confirm branch, HEAD, clean status,
reciprocal archive links, launcher `--check`, and exact product commit
`9164e06`.

Make one focused Maven stdout-selection implementation commit, then the normal
separate continuity commit. Stop before another Maven effect, plugin diagrams,
shell provenance, clock/server work, mutation harnesses, P4-P8, audit or scanner
changes, inventory changes, publication, or distribution. Do not push, merge,
stash, revert, launch a successor, or remove the worktree.
