# Quality Upgrade Handover

Generated: 2026-08-26T07:28:02+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`.
- Branch: `codex/upgrade-quality`.
- Base: `master` at `5635d50`.
- Focused P3.54 product implementation: `3bd07e99c3170fd0f1c0dda076252c6eb5833db7`.
- Truthful P3.54 apparatus/checkpoint commit:
  `ce736a233e1be1e17f290ad8cf3a327bf70ccfd0`.
- Product and test paths have no diff from `3bd07e9` through `ce736a2`.
- After launch, obtain the continuity head with `git rev-parse --short=12 HEAD`;
  that final commit changes only planning, archive, and launcher mutable state.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through move 54 are 5c6f2fa, 03d6242, f5ee37d,
ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce, a7eb3ef,
e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597, 9e2d669,
a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, d6cb593,
82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438, 527a8b9,
886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, 85f4c2b,
0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd, 2f0a072,
0128cd0, 70bee0e, and 3bd07e9 in roadmap order. The separate operational
continuity implementation is 1b85711. The focused T15 audit-apparatus repair is
ce736a2; neither operational commit changes a Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT because P3 is active and P4-P8 are queued. Its
active archive is
`docs/plan/agent-sessions/2026-08-26T072802+0200-migrate-plugin-diagrams-export.md`.
The P3.54 repair archive is answered history, reciprocal links are connected,
and the archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 54 And Checkpoint

Public `pkg/shell.Unzip(string, string) ([]string, error)` still selects the
complete private `unzipDependencies` composition with `filesystem.System()`.
Only its deferred archive close changed at product commit `3bd07e9`: the same
anonymous defer after successful archive open now calls
`filesystem.CloseReader(dependencies.Files, r)` and ignores the same result.

The exact `*zip.ReadCloser`, attempt count, defer ordering, traversal, entry and
output closes, bytes, filenames, partial results, zip-slip behavior, all
established errors, callers, public API, scanner, and inventory remain
unchanged. Two focused contracts bring the suite to 360 tests across 19 of 25
packages and prove successful, zip-slip, later-error, archive-open-error, and
ignored-close-error behavior without a real process or network request.

The complete P3.54 checkpoint passes truthfully at `ce736a2`. Compatibility,
build, tests, vet, pinned lint, all 62 launcher controls, Make contracts, all
production-script meta-tests, the 15-control audit suite, uncached and race
tests, empty-HOME count-2, and clean identity all pass. The full audit returns
the expected exit 1 for documented debt, never 2.

## T15 Apparatus Repair

The inherited T15 failure was reproduced before repair at exact predecessor
`70bee0e` and implementation `3bd07e9`: both passed T1-T14, then the old parser
rejected historical test-created ignored paths beyond the authorized inventory
overlay.

The focused red contract records every dirty path, type, mode, raw digest, and
checkout-root-independent digest in old/new clean-HOME execution replicas. It
verifies both replicas before audit execution and rejects a tracked-source drift
probe. T15 then uses separate equivalent execution and structured replicas:
historical tests may leave their truthful effects in the execution replicas,
while each parser measures a separately verified pristine exact `5635d50`
checkout plus the declared inventory overlay. No output path is excluded or
deleted, and every effect remains in the retained manifest.

That orchestration exposed a distinct historical environment precondition. An
empty HOME makes the untouched old tests pass, but stored raw baseline debt
records their failure. The directly related recipe and migration metadata now
pin an apparatus-owned directory-shaped `.co-pilot/profiles/.active_profile`
fixture. Its marker SHA-256 is
`cb95f24c35d3987f8aba51231aade19580ffe9242324804fcad2ccff350d1c9a`.
The parser, old and new scorecards, raw baseline debt, instrument identities,
exact source commit/tree, inventory, scanner, current product fixtures, and
product behavior did not change.

The apparatus repair is the four-path commit `ce736a2`:

- `.quality/tools/test-quality-audit.sh` owns replica orchestration, focused
  drift/effect contracts, and pinned HOME validation.
- `.quality/baseline/instrument-migration.json` binds the reproduction HOME.
- `.quality/baseline/reproduction-home/.../.fixture` is the pinned marker.
- `.quality/baseline/README.md` gives the two-replica reproduction recipe and
  explains the clean-HOME effect probe.

Complete repaired meta-suites pass over the implementation-content tree and
over exact predecessor `70bee0e` with only those apparatus paths overlaid. They
reproduce old scorecard
`d420887d73aabf496ff276fcc55d13ad379ac49c9322fad58808b5e28fdba7df`,
new scorecard
`5fb3226009cfbf0d29f63fa03592157cce4efcec6e38583b64a86f6288e89490`,
normalized raw body
`cf23c9dca987acd3a966693f933c4c7eca9f7d11c51f0fbabfe1d697f3d7497f`,
Q3.9, and all 228 numeric debt leaves.

## Measured Quality State

The authoritative clean full audit from `ce736a2` reports:

- Measured commit tree: `20266d8c4ff215679c981f9e4bf5d005fb435656`.
- Structured scorecard SHA-256:
  `c7ff3bd2cf32f8615c0fb1329c980bafba9b5e988cf4e785f4b8dc608d3f4243`.
- Absolute L0: 8 of 8.
- 360 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 10 direct external sites outside five declared adapters of 32
  production effect sites; missing clock and server keep it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 state-claim phrases across 72 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative audit output is `target/quality-audit/scorecard.json`. Gate
logs, exact tools, paired proof work, and preserved generated compatibility
reports are under `/private/tmp/ply-p354-checkpoint.z0VASE`,
`/private/tmp/ply-p354-meta-green2.plwJUi`, and
`/private/tmp/ply-p354-predecessor-proof.68Zacv` while retained.

## Next Objective

Complete one focused P3 production-effect move: migrate only the ignored
`structurizr-cli export` process request in `cmd/plugin_diagrams.go` through the
existing `internal/adapter/process` boundary.

Start red with private recording contracts that prove the complete command
name and ordered arguments, empty working directory and streams, synchronous
mode, one request, exact ignored process error, continued `FindAll` behavior,
complete system dependency selection, safe zero dependency, non-empty recorded
population, and no unrelated process request. Invoke no real external tool.

Then introduce only the private dependency seam required to select
`process.System()` in production and replace the one
`structurizr.Run(exec.Command("structurizr-cli", ...))` request with
`process.Execute` carrying the equivalent `process.Command`. Preserve ignored
error semantics and all later behavior.

Leave the Graphviz `dot` and macOS `open` requests direct, along with
`structurizr.RunWithOutputToFile`, file deletion/discovery, workspace and output
paths, Cobra registration/flags/help, initialization, public API, the process
adapter, inventory, scanner, audit apparatus, and all completed moves.

## Verification Notes

- Focused inherited failure: red after T1-T14 at the old parser dirty-tree
  identity check.
- Focused orchestration green exposed the empty-HOME/stored-debt mismatch and
  justified only the recipe, metadata, and apparatus fixture expansion.
- Final complete implementation-content meta-suite: 15 of 15 controls pass.
- Final exact-predecessor apparatus-overlay meta-suite: 15 of 15 controls pass.
- `make preflight`: pass after selecting exact apidiff and golangci-lint
  binaries and a writable external linter cache.
- `make test`, `make test-install`, `make test-agent-start`, uncached tests,
  race tests, vet, repeated audit meta-suite, and empty-HOME count-2: pass.
- Full audit: expected exit 1, 15 findings, zero regressions, zero dirty paths.
- The preflight-generated API/CLI reports were preserved outside the worktree
  before the authoritative clean audit; their hashes are recorded in the gate
  evidence directory.
- Product/test diff from `3bd07e9` through `ce736a2`: empty.

## Start And Stop

Read this handover, the linked NEXT archive, the P3 tail and checkpoint gate,
both design documents, `.quality/inventory`, complete plugin-diagrams and
structurizr code/tests/callers, the existing process adapter and every complete
double, and relevant compatibility/audit contracts before editing. Confirm
branch, HEAD, clean status, reciprocal archive links, launcher `--check`, the
apparatus commit, and the clean checkpoint identities above.

Make one focused plugin-diagrams export implementation commit, then the normal
separate continuity commit. Stop before the `dot` or `open` request, another
plugin-diagrams effect, clock/server work, mutation harnesses, P4-P8, audit or
scanner changes, inventory changes, publication, or distribution. Do not push,
merge, stash, revert, launch a successor, or remove the worktree.
