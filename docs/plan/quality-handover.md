# Quality Upgrade Handover

Generated: 2026-08-25T05:22:43+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 04cfe44bf675.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, and 04cfe44 in roadmap
order. The separate operational continuity implementation is 1b85711 and
changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T052243+0200-migrate-file-grep-recursive.md.
The find-all predecessor is answered history, and the reciprocal archive graph
has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 25 Preserved

Exported file.FindAll(suffix string, excludes []string, dir string) (result
[]string, err error) retains its signature and observable behavior. It is the
production wrapper around one private complete filesystem dependency, and
production selects filesystem.System(). The private helper delegates the exact
root and callback through filesystem.Walk. The unchanged callback ignores fi
and errIn, matches only strings.HasSuffix(path, suffix) && !SuffixIn(path,
excludes), preserves SuffixIn's exact ordered strings.Contains behavior,
appends every exact match in callback order, and always returns nil. A no-match
result remains nil. Every accumulated path and exact final error is preserved,
and only a final walk error equal to io.EOF is normalized to nil.

Seven top-level file recording contracts bring the suite to 209 tests. They
prove complete dependency delivery, production system selection, exact root
and callback delivery, callback order, ignored metadata and errors,
case-sensitive suffix matching, ordered substring exclusions including an
empty exclusion, ordered all-match accumulation, nil no-match results, exact
non-EOF errors with the current partial result, exact-only EOF normalization,
safe zero behavior without developer-path access, and non-empty walk and
callback populations. The contracts perform no real filesystem mutation.

Only FindAll's direct filepath.Walk moved. The adapter and every complete
filesystem recording double remain unchanged. FindFirst, GrepRecursive, IDE
cleanup, SuffixIn, Render, ClearDir, Move, DeleteAll, DeleteSingleFile,
append-open, file-read, directory-create, file-create, file-overwrite,
file-existence, Bitbucket, Wpost, supervisor, inventory, public API, callers,
and every P2A contract remain unchanged.

## Measured Quality State

The clean full audit at 04cfe44bf675 reports:

- Absolute L0: 8 of 8.
- 209 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 20 guarded safe-writer sites, 15 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 41 direct external sites outside five declared adapters of 57
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 39 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: tree 8dd8b0b910efa481a555a8c07dcf2a879f7d258a, status
  SHA-256 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d,
  and zero dirty paths.

The clean report is /private/tmp/ply-find-all-clean-audit-04cfe44 and the
focused report is /private/tmp/ply-find-all-focused-audit. Generated ignored
compatibility reports were removed before the clean audit; all authoritative
audit output stayed outside the measured tree.

## Decisions And Learned Facts

1. FindAll reuses the existing resolvable filesystem Walk dependency; no
   adapter operation or complete recording double changed.
2. FindAll's callback body stayed unchanged, including deliberately ignoring
   callback metadata and errors and always returning nil.
3. Q1.3 moved from 42/58 to 41/57 because the direct FindAll walk disappeared
   without adding a production adapter effect site.
4. GrepRecursive is the next isolated flow. It can reuse the existing Walk
   operation, nominally moving Q1.3 to 40/56 without another adapter extension.
5. GrepRecursive deliberately calls Grep once for each callback path, appends
   every hit in callback order, propagates a Grep error, and normalizes only a
   final walk error exactly equal to io.EOF. Grep and OpenLines stay unchanged.
6. Existing immutable pkg/file/test fixtures can exercise recursive-grep
   callback decisions without real filesystem mutation; do not inject or move
   the non-recursive Grep/OpenLines/Open chain.
7. .quality/inventory is baseline-checksum-bound. Do not relabel seams or claim
   P5 mutation coverage.
8. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight, and
   keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
9. The ignored .agent-task/current.md compatibility mirror is absent.
10. Supervisor success requires process, structured-stream, and committed
    repository evidence; final prose is observable only.
11. The first handoff-only launcher run hit the previously observed
    partial-raw-log signal-fixture flake at control 26. The immediate complete
    rerun passed all 62 controls; never weaken or skip this contract.

## Next Objective

Move only exported file.GrepRecursive(targetDir string, keyword string) (files
[]string, err error)'s direct filepath.Walk(targetDir, callback) through the
existing zero-value-safe filesystem Walk operation.

Start red with file recording contracts for the complete dependency,
production system selection, exact root and callback delivery, callback order,
ignored metadata/errors, exact case-sensitive keyword decisions through the
unchanged per-path Grep calls, ordered all-hit accumulation, no-match nil slice
and nil error, exact non-EOF walk errors with the current partial result, final
io.EOF normalization, safe zero behavior without developer-path access, and
non-empty walk/callback populations. Use immutable checked-in fixtures and
perform no real filesystem mutation.

Keep GrepRecursive's signature and callback body unchanged. Production must
select filesystem.System(), and the private helper must carry one complete
filesystem dependency. Reuse adapter Walk without extending the adapter or
changing complete doubles. Do not move or inject Grep, OpenLines, Open,
FindFirst, FindAll, IDE cleanup, callers, or another operation, and do not enter
HTTP/process, config, clock/server, P4, P5, or later roadmap work.

Expected direction is nominally Q1.3 40/56 with Q1.2, Q1.4, and exact Q2.1
held. Regenerate the structured result and require zero comparable ratchet
regressions.

## Verification Notes

Completed from implementation commit 04cfe44bf675:

- Red evidence: focused compilation failed only on absent private find-all
  dependency and helper symbols.
- Focused filesystem, file, all caller, and adjacent complete-double package
  tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls before and after
  measured planning-note edits. The first handoff-only rerun hit the known
  control-26 signal-fixture flake; the immediate complete rerun passed all 62.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls.
- API, CLI, subprocess compatibility: PASS.
- make test and make test-install: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable.
- Clean full audit: exit 1, 16 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: 04cfe44bf6750d9309f774ffec21a6af576db0bd
  (quality: route find-all walk through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, all file code/tests and GrepRecursive callers, the
adapter and complete doubles, and the named audits before editing. Confirm
branch, HEAD, status, reciprocal links, and ./codex-dev-start.sh --check. Begin
red for only GrepRecursive's walk boundary, then finish with one implementation
commit and one separate handoff-only commit.

Stop before non-recursive Grep/OpenLines/Open, IDE cleanup, another adapter
operation or family, Q1.4 expansion, P4, mutation, Docker, cloud, distribution,
or publication. Stop on API/CLI change, comparable ratchet regression,
full-audit exit 2, or failure to isolate implementation from generated and
handoff-only state.
