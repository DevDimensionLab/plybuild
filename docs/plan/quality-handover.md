# Quality Upgrade Handover

Generated: 2026-08-26T08:03:52+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.55 product implementation:
  `3972fc6500cdb68a3b48e12dda2c90c6c1e18783`.
- Preceding P3.54 product implementation:
  `3bd07e99c3170fd0f1c0dda076252c6eb5833db7`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 55 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, and 3972fc6 in roadmap order. The separate
operational continuity implementation is 1b85711. The focused T15
audit-apparatus repair is ce736a2; neither operational commit changes a Go
quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T080352+0200-migrate-plugin-diagrams-open.md`.
The P3.55 export archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`33599b9547262f76055dfa8d61c121e541051d6633a1ece6ad924a1b5eb40831`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 55

Only the ignored `structurizr-cli export` request in the private
plugin-diagrams flow changed at product commit `3972fc6`. The mandatory
workspace lookup remains first. `structurizrCmd.RunE` then passes the complete
private `pluginDiagramsExportDependencies` value to `runStructurizrDiagrams`,
and production selects `process.System()` explicitly.

After the unchanged exact `.structurizr/` deletion and before unchanged
dot-file discovery, the flow calls `_ = process.Execute` once with:

- executable `structurizr-cli`;
- ordered arguments `export`, `-w`, the caller-supplied workspace, `-format`,
  `dot`, `-output`, `.structurizr/`;
- empty `Dir`, nil stdin/stdout/stderr, and `Start: false`.

The synchronous result is still ignored and discovery always follows the
attempt. The later Graphviz `dot` and macOS `open` requests remain direct and
unchanged. Iteration, paths, printed text, errors, Cobra behavior, public API,
`pkg/structurizr`, the adapter, scanner, inventory, audit apparatus, and every
completed effect remain unchanged.

Four focused private contracts bring the suite to 364 tests. They prove the
complete production dependency selection, exact command and ordered arbitrary
workspace bytes, one request, zero-value process fields, exact ignored injected
error, continued discovery, safe zero dependency, rejection of an empty
recorded population, and absence of unrelated export, `dot`, or `open`
requests. They use temporary working directories and launch no external
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
Q3.9, and all 228 numeric debt leaves. Product move 55 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `3972fc6` reports:

- Measured commit tree: `395bbad053237f8769be804e211c7cb0f49550cc`.
- Structured scorecard SHA-256:
  `ced6611a0990d55285007f9b9d19a990bbbe2bef9ee2fbe18084420fe9956fa6`.
- Absolute L0: 8 of 8.
- 364 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 9 direct external sites outside five declared adapters of 31 production
  effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 73 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative audit output is `target/quality-audit/scorecard.json`. Exact
gate tools and retained logs are under `/private/tmp/ply-p355-gate.zBORR2` while
present.

## Next Objective

Complete one focused P3 production-effect move: route only the ignored macOS
`open` request after each successful Graphviz conversion in
`cmd/plugin_diagrams.go` through the existing process adapter.

Start red with private recording contracts that prove exact command `open`,
the exact single `outputPngFile` argument, empty directory and streams,
synchronous mode, one request per helper invocation, exact ignored injected
error, placement after successful direct `dot`, continued iteration, complete
production dependency selection, safe zero dependency, non-empty recorded
population, and no export, `dot`, or unrelated process request. Invoke no real
external tool.

Then add only the private dependency composition needed to select
`process.System()` in production and replace the single ignored
`structurizr.Run(exec.Command("open", outputPngFile))` request with equivalent
ignored `process.Execute`. Preserve the completed export move and all direct
Graphviz behavior.

Leave workspace lookup, deletion and discovery, iteration and output paths,
printed text, Graphviz execution/error return, Cobra registration and help,
public API, `pkg/structurizr`, process adapter, inventory, scanner, baseline,
audit apparatus, clock/server work, mutation harnesses, and P4-P8 unchanged.

## Verification Notes

- Focused red contracts failed only for the missing private export dependency
  boundary and helpers; focused green tests launch no real process.
- Focused cmd/process/structurizr and caller tests: pass.
- API/CLI and subprocess compatibility: pass with generated reports outside
  the worktree.
- Complete and uncached tests, race, vet, pinned lint, build, Make contracts,
  all production-script meta-contracts, and four host acceptance flows: pass.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass.
- The first launcher attempts with an explicitly nested temporary root hit the
  documented partial-raw-log signal timing flake. The standalone default-path
  contract and the complete Make run both passed all 62 controls.
- Focused audit: expected exit 1, four improved, two held, zero regressed, one
  non-comparable, zero dirty paths.
- Full clean audit: expected exit 1, 15 findings, zero comparable regressions,
  zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete plugin-diagrams and
structurizr code/tests/callers, the process adapter and every complete double,
the import-aware scanner, API/CLI contracts, and the T15 repair and baseline
reproduction README before editing. Confirm branch, HEAD, clean status,
reciprocal archive links, launcher `--check`, and exact product commit
`3972fc6`.

Make one focused plugin-diagrams open implementation commit, then the normal
separate continuity commit. Stop before the direct Graphviz request, another
plugin-diagrams effect, clock/server work, mutation harnesses, P4-P8, audit or
scanner changes, inventory changes, publication, or distribution. Do not push,
merge, stash, revert, launch a successor, or remove the worktree.
