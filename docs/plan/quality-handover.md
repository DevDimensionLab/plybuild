# Quality Upgrade Handover

Generated: 2026-08-26T08:40:00+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.56 product implementation:
  `6f4bc3737ba07a2a4614617d272ef44558d39b77`.
- Preceding P3.55 product implementation:
  `3972fc6500cdb68a3b48e12dda2c90c6c1e18783`.
- Truthful P3.54 audit-apparatus repair:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  the final commit changes only planning, archives, and launcher mutable state.
- No push, merge, release, publication, distribution, stash, revert, successor
  launch, or worktree removal was performed.

P3 implementation commits through move 56 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, 3bd07e9, 3972fc6, and 6f4bc37 in roadmap order. The
separate operational continuity implementation is 1b85711. The focused T15
audit-apparatus repair is ce736a2; neither operational commit changes a Go
quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T084000+0200-migrate-plugin-diagrams-dot.md`.
The P3.56 open archive is answered history and links reciprocally to that
archive. The graph has exactly one NEXT tail. The archived prompt SHA-256 is
`5c260abb2033cc27c5c7ee8ee57d6f485cc4e44ed753bf926b408d99d3a5ba05`.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 56

Only the ignored macOS `open` request after a successful Graphviz conversion in
the private plugin-diagrams flow changed at product commit `6f4bc37`.
`structurizrCmd.RunE` retains the mandatory workspace lookup and passes both
complete private dependency values to `runStructurizrDiagrams`. Production
selects `process.System()` explicitly for the new
`pluginDiagramsOpenDependencies` value.

After the unchanged direct `dot` conversion succeeds and before the range
advances, `openStructurizrDiagram` calls `_ = process.Execute` once with:

- executable `open`;
- exact single argument `outputPngFile`;
- empty `Dir`, nil stdin/stdout/stderr, and `Start: false`.

The synchronous result is still ignored and the next discovered file is always
attempted after the open attempt. The completed export request, direct Graphviz
request and returned error, discovery, iteration, output paths, printed text,
output capture/write, Cobra behavior, public API, `pkg/structurizr`, adapters,
scanner, inventory, and audit apparatus remain unchanged.

Five focused private contracts bring the suite to 369 tests. They prove the
complete production dependency selection, exact command and arbitrary output
path bytes, one request, zero-value directory/streams/Start, exact directly
ignored injected error, placement only after the direct Graphviz error gate,
continued iteration, safe zero dependency, rejection of an empty recorded
population, and absence of export, `dot`, or unrelated process requests. They
use a recorder and source AST, launch no external program, touch no network, and
write no repository fixture.

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
Q3.9, and all 228 numeric debt leaves. Product move 56 did not modify the
apparatus, parser, scanner, baseline, inventory, or reproduction recipe.

## Measured Quality State

The authoritative clean full audit from product commit `6f4bc37` reports:

- Measured commit tree: `cf7c1af5b9ecbbbb7b218cfa6951354d6bf94854`.
- Structured scorecard SHA-256:
  `efaaee797a34d182d03d6830d0af4c342b7e744180775942119a77c7e0f88090`.
- Absolute L0: 8 of 8.
- 369 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 8 direct external sites outside five declared adapters of 30 production
  effect sites; missing clock and server keep it non-comparable. The remaining
  plugin-diagrams violation is only the unchanged direct Graphviz request.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 74 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for the same 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative audit output is `target/quality-audit/scorecard.json`.
Focused audit output is under `/private/tmp/ply-p356-focused-audit`; pinned
verification tools and compatibility reports are under `/private/tmp`.

## Next Objective

Complete one focused P3 production-effect move: route only the direct Graphviz
`dot` conversion request in `cmd/plugin_diagrams.go` through the process adapter
while preserving its existing captured-output behavior.

Start red with private recording contracts for exact executable `dot`, ordered
arguments `file`, `-Tpng`, empty directory, nil stdin, distinct stdout/stderr
buffers, synchronous mode, one attempt, exact process error precedence, exact
successful stdout bytes/path/`0644` write, ignored write error, safe zero
dependencies, non-empty populations, following completed open placement, and no
unrelated request. Invoke no real external tool.

Then add only the private process/filesystem dependency composition needed to
preserve the current conversion behavior. Use `process.Execute` for the direct
Graphviz request and `filesystem.WriteFile` for its existing successful output
write. Keep process errors returned before any write/open, keep write errors
ignored, keep the completed open helper after success, and leave
`pkg/structurizr`, its public API, both adapters, and every other flow unchanged.

Leave workspace lookup, deletion, completed export and open requests, discovery,
iteration and output paths, printed text, Cobra registration/help, public API,
inventory, scanner, baseline, audit apparatus, clock/server work, mutation
harnesses, and P4-P8 unchanged.

## Verification Notes

- Focused red contracts failed only for the missing private open dependency and
  helper; focused green tests launch no real process.
- Focused cmd/process/structurizr and relevant caller tests: pass.
- API/CLI and subprocess compatibility: pass with generated reports outside the
  worktree.
- Complete and uncached tests, race, vet, pinned lint, build, Make contracts,
  all production-script meta-contracts, and four host acceptance flows: pass.
- All 62 launcher controls and the repaired 15-control audit suite: pass.
- Empty-HOME count-2: pass; its disposable HOME remains under `/private/tmp`.
- Focused audit: expected exit 1, five improved, two held, zero regressed, one
  non-comparable, zero dirty paths.
- Full clean audit: expected exit 1, 15 findings, zero comparable regressions,
  zero dirty paths.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete plugin-diagrams and
structurizr code/tests/callers, the process and filesystem adapters and relevant
complete doubles, the import-aware scanner, API/CLI contracts, and the T15
repair and baseline reproduction README before editing. Confirm branch, HEAD,
clean status, reciprocal archive links, launcher `--check`, and exact product
commit `6f4bc37`.

Make one focused plugin-diagrams Graphviz implementation commit, then the normal
separate continuity commit. Stop before another plugin-diagrams effect,
clock/server work, mutation harnesses, P4-P8, audit or scanner changes, inventory
changes, publication, or distribution. Do not push, merge, stash, revert, launch
a successor, or remove the worktree.
