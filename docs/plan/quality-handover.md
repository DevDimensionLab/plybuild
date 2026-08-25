# Quality Upgrade Handover

Generated: 2026-08-25T13:58:54+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: b2d37cc43cde.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438,
527a8b9, 886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, and b2d37cc
in roadmap order. The separate operational continuity implementation is
1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T135854+0200-migrate-tips-list-read-dir.md.
The local-config directory-mkdir predecessor is answered history, and the
reciprocal archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 42 Preserved

Exported `LocalConfigDir.CheckOrCreateConfigDir`, every caller, and both later
local-config create and write lifecycles remain unchanged. The method still
evaluates exact `localCfgDir.Implementation().Path`, performs the completed
one-attempt adapter `Stat` with that exact arbitrary string, and applies exact
`os.IsNotExist(err)`. Nil stat, existing-directory stat, and arbitrary
non-not-exist errors return nil without mkdir. Only the not-exist branch
selects `filesystem.System()` for a private complete directory-create
composition and makes one exact single-level `Mkdir(dir, 0755)` attempt; the
exact mkdir error remains the result and success returns nil.

`filesystem.FileSystem` now has one distinct `Mkdir(string, fs.FileMode) error`
operation. Its zero value returns exact `filesystem.ErrNoFilesystem`; its
system implementation maps directly to `os.Mkdir`; established `MkdirAll`
remains unchanged. All 34 pre-existing complete filesystem test doubles add
only the required rejecting method, while the dedicated adapter and local-
config doubles record it. There are now 35 complete doubles for the next
interface extension.

Six focused contracts bring the suite to 296 tests. They prove exact paths and
modes, non-empty populations, one attempt, exact error identity, safe zero
behavior, system success, missing-parent and existing-directory failure, no
recursive creation, complete production selection, exact `0755`, nil and
arbitrary error results, and no unrelated operation. Production-composition
contracts neither mutate the real filesystem nor run another command. No
other adapter, stat/predicate behavior, public API, caller, inventory, seam,
mutation harness, or completed effect changed.

## Measured Quality State

The clean full audit from implementation commit b2d37cc reports:

- Absolute L0: 8 of 8.
- 296 test functions, zero skipped; 17 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 8 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 22 direct external sites outside five declared adapters of 40
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 56 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
/private/tmp/ply-mkdir-clean-audit.9gjsCn/scorecard.json and the focused report
is /private/tmp/ply-mkdir-focused-audit.hbihYW/scorecard.json. All audit output,
compatibility reports, tool binaries, and caches stayed outside the measured
tree.

## Decisions And Learned Facts

1. The new operation is intentionally distinct from `MkdirAll`: single-level
   failure behavior for missing parents and existing directories is observable.
2. Required `os.Mkdir` inside the exact system adapter remains in the scanner's
   production population, so Q1.3 is 22/40 rather than the nominal 22/39.
3. The local-config method preserves stat-before-predicate-before-mkdir order;
   the mkdir composition is selected only after exact `os.IsNotExist` succeeds.
4. The dedicated recording double rejects all unrelated adapter operations and
   an empty mkdir population. Safe zero behavior does not touch a developer path.
5. `tips.List` cannot reuse established `ReadDir`: its result is
   `[]os.DirEntry`/`[]fs.DirEntry`, while established `ReadDir` returns
   `[]fs.FileInfo` and maps to `ioutil.ReadDir`.
6. The next exact operation must therefore be
   `ReadDirEntries(string) ([]fs.DirEntry, error)` with direct `os.ReadDir`
   system semantics, leaving `ReadDir` untouched and updating all 35 current
   complete doubles only for interface compatibility.
7. Preserve `tips.List` filtering precisely: read error returns the named nil
   result and exact error without inspecting entries; success preserves
   delivered order and entry identity; `IsDir` short-circuits `Name`; only
   non-directories with case-sensitive `.md` suffix are retained; do not call
   `Info`.
8. Moving `tips.List` adds direct package contracts, so expected direction is
   Q1.1 from 8/25 to 7/25 and Q1.3 from 22/40 to 21/40. Regenerate both.
9. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams or
   claim P5 mutation coverage.
10. Full preflight needs APIDIFF and GOLANGCI_LINT discovered from an external
    PATH, with GOCACHE, GOTMPDIR, and GOLANGCI_LINT_CACHE in writable external
    directories. A command-line make override propagates into nested meta-tests
    and is not contract-correct.

## Next Objective

Move only direct `os.ReadDir(LocalDir(gitCfg))` in `tips.List` through one new
exact filesystem adapter `ReadDirEntries` operation.

Start red at the adapter boundary for exact arbitrary path, delivered entry
slice and error identity, non-empty population, one attempt, safe zero behavior,
and direct system sorting/error semantics. Then add direct tips-list contracts
for complete production selection, exact `LocalDir(gitCfg)`, read-error early
return, delivered order and identity, exact directory and case-sensitive `.md`
filtering, short-circuit behavior, nil/empty/all-filtered populations, and no
unrelated operation. Production-composition contracts must not inspect the real
filesystem or run another command.

Preserve `tips.List`'s public signature and callers, exact directory selection,
the existing `ReadDir` operation, tips-show read behavior, all filtering and
return semantics, and every completed adapter move. Add only the distinct
internal operation, its system mapping and contracts, the required complete-
double methods, and focused tips tests. Do not change another tips effect,
adapter operation, inventory, public API/CLI, or later roadmap work.

Expected direction is Q1.1 7/25 and Q1.3 21/40 with Q0.6, Q1.2, Q1.4, and
exact Q2.1 held. Require zero comparable ratchet regressions.

## Verification Notes

Completed from implementation commit b2d37cc43cdecbfe78a37caf11c65344115bedd4:

- Red evidence: focused compilation failed only on absent adapter `Mkdir` and
  private local-config directory-create dependencies/helpers.
- Focused config/filesystem/process and relevant command, context, file,
  template, Maven, tips, structurizr, Bitbucket, HTTP, Kibana, Spring, and shell
  package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls.
- Make preflight meta-contract, full preflight, and all 15 audit meta-controls:
  PASS after contract-correct external cache and tool setup.
- API, CLI, and subprocess compatibility: PASS.
- Full make test, install, and agent-start targets: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS from dirty and clean states.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 22/40.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases and zero
  ratchet regressions.
- Implementation commit: b2d37cc43cdecbfe78a37caf11c65344115bedd4
  (quality: route local config directory mkdir through filesystem).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete tips implementation/tests and callers, relevant
filesystem/file/template/config/Maven/structurizr behavior, complete filesystem
and process adapters, all 35 complete filesystem doubles, and named audits
before editing. Confirm branch, HEAD, status, reciprocal links, and
`./codex-dev-start.sh --check`. Begin red only for the distinct filesystem
`ReadDirEntries` contract and tips-list boundary, then finish with one
implementation commit and one separate handoff-only commit.

Stop before changing established `ReadDir`, tips-show read, filtering semantics,
another adapter operation or family, another tips/config effect, Q1.4 expansion,
P4, mutation, Docker, cloud distribution, or publication. Stop on API/CLI
change, comparable ratchet regression, authoritative full-audit exit 2, or
failure to isolate implementation from generated and handoff-only state.
