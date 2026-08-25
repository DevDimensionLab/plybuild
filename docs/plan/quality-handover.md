# Quality Upgrade Handover

Generated: 2026-08-25T05:51:00+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: d03e96d64992.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, and d03e96d
in roadmap order. The separate operational continuity implementation is
1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T055100+0200-migrate-file-ide-read-dir.md.
The recursive-grep predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 26 Preserved

Exported file.GrepRecursive(targetDir string, keyword string) (files []string,
err error) retains its signature and observable behavior. It is the production
wrapper around one private complete filesystem dependency, and production
selects filesystem.System(). The private helper delegates the exact root and
callback through filesystem.Walk. The unchanged callback ignores fi and errIn,
calls Grep(path, keyword) exactly once for each delivered path, returns its
error unchanged when non-nil, appends the exact path only for a hit, and
otherwise returns nil. Every accumulated path and exact final walk error is
preserved, and only a final error equal to io.EOF is normalized to nil.

Seven top-level recording contracts bring the suite to 216 tests. They prove
complete dependency delivery, production system selection, exact root and
callback delivery, callback order, ignored metadata and errors, exact
case-sensitive keyword decisions against immutable checked-in fixtures,
suppressed per-path read failures, ordered accumulation of all hits including
duplicate callback paths, nil no-match results, exact non-EOF errors with the
current partial result, exact-only EOF normalization, safe zero behavior
without developer-path access, and non-empty walk and callback populations.
The contracts perform no real filesystem mutation.

Only GrepRecursive's direct filepath.Walk moved. The adapter and every complete
filesystem recording double remain unchanged. Grep, OpenLines, Open,
FindFirst, FindAll, IDE cleanup, SuffixIn, Render, ClearDir, Move, DeleteAll,
DeleteSingleFile, Bitbucket, Wpost, supervisor, inventory, public API, callers,
and every P2A contract remain unchanged.

## Measured Quality State

The clean full audit at d03e96d64992 reports:

- Absolute L0: 8 of 8.
- 216 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 20 guarded safe-writer sites, 15 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 40 direct external sites outside five declared adapters of 56
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 40 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: tree c5381b37accb70091181684b5d99cabcbf66a3d2, status
  SHA-256 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d,
  and zero dirty paths.

The clean report is
/private/tmp/ply-grep-recursive-clean-audit-d03e96d/scorecard.json and the
focused report is /private/tmp/ply-grep-recursive-focused-audit. Compatibility,
tool caches, and all authoritative audit output stayed outside the measured
tree.

## Decisions And Learned Facts

1. GrepRecursive reuses the existing resolvable filesystem Walk dependency;
   no adapter operation or complete recording double changed.
2. GrepRecursive's callback body stayed unchanged, including deliberately
   ignoring callback metadata and errors and invoking the existing Grep once
   for every callback path.
3. Grep and OpenLines remain uninjected and unchanged. Their existing
   case-sensitive strings.Contains match and suppression of read failures are
   exercised with immutable checked-in fixtures.
4. Q1.3 moved from 41/57 to 40/56 because the direct GrepRecursive walk
   disappeared without adding a production adapter effect site.
5. Non-recursive IDE cleanup is the next isolated flow. Its direct
   ioutil.ReadDir can move behind one new filesystem ReadDir operation,
   nominally moving Q1.3 to 39/56.
6. To preserve behavior exactly, the system ReadDir operation should retain
   ioutil.ReadDir's []fs.FileInfo result and filename sorting. Update every
   complete filesystem recording double for the interface extension.
7. The non-recursive IDE loop has two ordered, independent substring checks:
   non-directory names containing .iml and directory names containing .idea.
   Preserve entry order, Path formatting, logAndDelete calls, first error,
   counts, and the exact report.
8. Recursive IDE cleanup, FindAll, logAndDelete, and DeleteAll remain outside
   the next move. Do not repair their existing behavior while moving ReadDir.
9. .quality/inventory is baseline-checksum-bound. Do not relabel seams or claim
   P5 mutation coverage.
10. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are available under /private/tmp/ply-quality-tools for the
    lifetime of that temporary directory.
11. The ignored .agent-task/current.md compatibility mirror is absent.
12. Supervisor success requires process, structured-stream, and committed
    repository evidence; final prose is observable only.

## Next Objective

Move only non-recursive removeIntellijFile(targetDir string, dryRun bool)'s
direct ioutil.ReadDir(targetDir) through one new ReadDir operation on the
existing zero-value-safe filesystem adapter.

Start red with adapter and file recording contracts for the complete
dependency, production system selection, exact target, ordered FileInfo
results and exact errors, system filename sorting, exact .iml file and .idea
directory substring decisions, ignored opposite types, ordered dry-run paths
and logs, counts and report, read-error short circuit, safe zero behavior
without developer-path access, and non-empty directory-read and entry
populations. Keep any fixture mutation confined to guarded temporary
directories.

Keep exported RemoveIntellijFiles and its recursive flag selection unchanged.
Production must select filesystem.System() only for the private non-recursive
composition. Extend the complete adapter and all complete recording doubles
with ReadDir returning []fs.FileInfo, and retain ioutil.ReadDir in the system
implementation for exact sorted behavior. Do not change the loop body,
recursive cleanup, FindAll, GrepRecursive, logAndDelete, DeleteAll, Path,
inventory, public API, or another operation, and do not enter HTTP/process,
config, clock/server, P4, P5, or later roadmap work.

Expected direction is nominally Q1.3 39/56 with Q1.2, Q1.4, and exact Q2.1
held. Regenerate the structured result and require zero comparable ratchet
regressions.

## Verification Notes

Completed from implementation commit d03e96d64992:

- Red evidence: after moving the Go cache to an isolated writable directory,
  focused compilation failed only on absent private recursive-grep dependency
  and helper symbols. The earlier host-cache permission failure was
  environmental and not accepted as red evidence.
- Focused filesystem, file, config caller, and adjacent complete-double package
  tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and pinned apidiff and golangci-lint tools outside the
  repository.
- API, CLI, and subprocess compatibility: PASS.
- make test and make test-install: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable.
- Clean full audit: exit 1, 16 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: d03e96d649925d54634d6625870d9d8591df695e
  (quality: route recursive grep walk through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, all file code/tests and RemoveIntellijFiles callers, the
adapter and complete doubles, and the named audits before editing. Confirm
branch, HEAD, status, reciprocal links, and ./codex-dev-start.sh --check. Begin
red for only the non-recursive IDE directory-read boundary, then finish with
one implementation commit and one separate handoff-only commit.

Stop before recursive IDE cleanup, FindAll, GrepRecursive, logAndDelete,
DeleteAll, another adapter operation or family, Q1.4 expansion, P4, mutation,
Docker, cloud, distribution, or publication. Stop on API/CLI change,
comparable ratchet regression, full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
