# Quality Upgrade Handover

Generated: 2026-08-25T14:34:03+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 85f4c2be48b1.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438,
527a8b9, 886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc, and
85f4c2b in roadmap order. The separate operational continuity implementation
is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T143403+0200-migrate-profile-editor-process.md.
The tips-list directory-read predecessor is answered history, and the
reciprocal archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 43 Preserved

Exported `tips.List`, `tips.LocalDir`, `TipsDir`, their callers, and
`tipsShowCmd` remain unchanged. `List` still evaluates the exact
`file.Path("%s/%s", gitCfg.Implementation().Dir(), TipsDir)` result, makes one
directory-entry read, returns the exact error with a nil result before entry
inspection, and on success preserves delivered order and entry identity while
short-circuiting `IsDir` before `Name` and accepting only case-sensitive `.md`
non-directories. Empty and all-filtered results remain nil.

`filesystem.FileSystem` now has one distinct
`ReadDirEntries(string) ([]fs.DirEntry, error)` operation. Its zero value
returns exact `filesystem.ErrNoFilesystem`; its system implementation maps
directly to `os.ReadDir`. Established `ReadDir`, its `[]fs.FileInfo` result,
and its `ioutil.ReadDir` system mapping remain unchanged. All 35 pre-existing
complete filesystem test doubles add only the required method: the adapter
double records it and the other 34 reject it as unrelated. The new tips-list
double records only this operation, so there are now 36 complete doubles for
the next filesystem interface extension.

Eight focused contracts bring the suite to 304 tests and direct coverage in
`pkg/tips`. They prove exact arbitrary paths, ordered entry identity, nil,
empty, representative, partial-entry, and arbitrary-error results, safe zero
behavior, direct system sorting and failure semantics, complete production
selection, exact error and nil-on-error behavior, short-circuit filtering, and
no unrelated operation. Production-composition contracts neither read the real
filesystem nor run another command. No public API, CLI, caller, inventory,
seam, mutation harness, completed adapter operation, or completed effect changed.

## Measured Quality State

The clean full audit from implementation commit 85f4c2b reports:

- Absolute L0: 8 of 8.
- 304 test functions, zero skipped; 18 of 25 packages have tests.
- Q0.6: 24 guarded safe-writer sites, 19 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 7 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 21 direct external sites outside five declared adapters of 40
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 57 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
/private/tmp/ply-tips-list-clean-audit.D4FNIT/scorecard.json and the focused
implementation report is
/private/tmp/ply-tips-list-implementation-audit.nOapo6/scorecard.json. All
audit output, compatibility reports, tool binaries, caches, and empty-HOME
state stayed outside the measured tree.

## Decisions And Learned Facts

1. `ReadDirEntries` is intentionally distinct from established `ReadDir`:
   `fs.DirEntry` identity, lazy metadata, and `os.ReadDir` behavior are
   observable and cannot be reproduced through `fs.FileInfo` conversion.
2. Required `os.ReadDir` inside the exact system adapter remains in the
   scanner population, so the measured Q1.3 result is 21/40.
3. The filtering contract proves that an error wins over partial entries and
   that a directory's name is never read; `Info` and `Type` are never called.
4. The two guarded temporary fixture writes raise the Q0.6 safe-writer call
   population from 22 to 24 without changing its zero skipped/unsafe metric.
5. The next coherent measured process flow is the profile editor launch in
   `cmd/profile.go`: direct `exec.Command` and `cmd.Run` are two remaining Q1.3
   sites in one edit branch.
6. Existing `process.Command` carries name, args, directory, stdout, and stderr
   but not stdin. Exact editor behavior requires adding `Stdin io.Reader` and
   mapping it directly to `exec.Cmd.Stdin`; do not invent a second process
   operation or a function-valued dependency.
7. Preserve exact editor selection: one `os.Getenv("EDITOR")`, empty fallback
   to `vim`, exact single `ctx.LocalConfig.FilePath()` argument, exact
   `os.Stdin` and `os.Stdout`, empty directory, nil stderr, one run attempt, and
   exact run error. Do not route environment lookup or another profile effect.
8. A failed editor run returns before sync, reset, and print. A successful edit
   retains the existing later sync/reset/print selection and ordering.
9. The expected Q1.3 direction is nominally 21/40 to 19/38 when the two direct
   process sites leave; regenerate the exact scanner result. Q1.1 should hold
   at 7/25 because `cmd` already has tests.
10. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams,
    make mutation harnesses executable, or claim P5 coverage.
11. Full preflight needs external APIDIFF, GOCACHE, GOTMPDIR,
    GOLANGCI_LINT_CACHE, and golangci-lint v2.12.2. The first move-43 preflight
    attempt reached lint and failed only because its default cache was not
    writable; the external-cache rerun passed every repository control.

## Next Objective

Move only the profile edit branch's direct `exec.Command(editor,
ctx.LocalConfig.FilePath())` and `cmd.Run()` through the existing process
adapter.

Start red by extending the process adapter contracts for exact stdin identity
and direct system stdin delivery. Add focused profile-editor recording
contracts for complete production selection, exact EDITOR/fallback executable,
exact one-argument local-config path, exact stdin/stdout identity, empty
directory, nil stderr, one attempt, exact error, non-empty population, and no
real command. Then add only `Stdin io.Reader` to `process.Command`, map it
directly, and replace the direct editor construction/run through a private
complete profile-edit composition.

Preserve every command object and registration, `ConfigOpts`, global context,
EDITOR lookup and fallback, local-config path evaluation, all later edit/sync/
reset/print sequencing, terminal-width config behavior, process zero defaults,
Maven process behavior, every completed adapter and tips move, public API/CLI,
inventory, and all callers. Do not migrate another profile branch, command,
process, filesystem, environment, clock, server, P4, or P5 effect.

Expected direction is Q1.1 7/25 and nominal Q1.3 19/38 with Q0.6, Q1.2,
Q1.4, and exact Q2.1 held. Require zero comparable ratchet regressions and
regenerate exact values.

## Verification Notes

Completed from implementation commit 85f4c2be48b1f2d33dda0143542fb39a2621d2c4:

- Red evidence: focused compilation failed only on absent `ReadDirEntries`,
  `listDependencies`, `systemListDependencies`, and private `list` symbols.
- Focused tips/filesystem/process and relevant command, context, config, file,
  template, Maven, structurizr, Bitbucket, HTTP, Kibana, Spring, and shell
  package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls during the
  implementation gate and on the complete handoff rerun. The first handoff run
  hit only the signal probe's intermittent partial-log timing failure.
- Make preflight meta-contract, complete preflight, and all 15 audit
  meta-controls: PASS after the external linter-cache rerun.
- API, CLI, and subprocess compatibility: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their meta-contracts:
  PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS from dirty and clean states.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 21/40.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 58
  Markdown files and zero ratchet regressions.
- Implementation commit: 85f4c2be48b1f2d33dda0143542fb39a2621d2c4
  (refactor: route tips directory read through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete `cmd/profile.go` and relevant command/context/
config tests and callers, complete process adapter and tests, Maven process
contracts, all current tips/filesystem contracts, and named audits before
editing. Confirm branch, HEAD, status, reciprocal links, and
`./codex-dev-start.sh --check`. Begin red only for process stdin and the profile
editor boundary, then finish with one implementation commit and one separate
handoff-only commit.

Stop before changing EDITOR or fallback selection, local-config path behavior,
another profile branch or command, later sync/reset/print behavior, another
adapter operation or family, inventory, Q1.4, P4, mutation, Docker,
distribution, or publication. Stop on API/CLI change, comparable ratchet
regression, authoritative full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
