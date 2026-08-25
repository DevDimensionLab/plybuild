# Quality Upgrade Handover

Generated: 2026-08-25T06:18:02+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: d6cb5939ccce.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d, and
d6cb593 in roadmap order. The separate operational continuity implementation
is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T061802+0200-migrate-config-git-hook-read-dir.md.
The non-recursive IDE predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 27 Preserved

The existing zero-value-safe filesystem adapter now has one additional
ReadDir(path string) ([]fs.FileInfo, error) operation. Its system implementation
delegates directly to ioutil.ReadDir(path), preserving exact FileInfo values,
errors, and filename sorting. A zero-value dependency returns
filesystem.ErrNoFilesystem without filesystem access. Every complete
filesystem recording double implements the extension.

Exported file.RemoveIntellijFiles(targetDir string, recursive bool, dryRun
bool) (string, error) retains its signature, recursive flag selection, and
observable behavior. Only the non-recursive private composition selects
filesystem.System() and passes the exact targetDir through filesystem.ReadDir.
The recursive branch, FindAll, GrepRecursive, logAndDelete, DeleteAll, Path,
the command caller, and every other file operation remain unchanged.

The non-recursive loop still iterates entries in dependency-delivered order and
performs two ordered, independent checks: a non-directory exact name containing
.iml, then a directory exact name containing .idea. Each selected entry still
uses Path("%s/%s", targetDir, f.Name()), calls logAndDelete, and stops at its
first error. Counts and the exact `Iml files: %d, .idea dirs: %d` report remain
unchanged. An initial read error returns an empty report and the exact error.

Eight new top-level adapter and IDE contracts bring the suite to 224 tests.
They prove complete production selection and delivery, exact target, delivered
entry order and metadata, system filename sorting, exact errors, classification
and ignored opposite types, case sensitivity, dry-run paths and logs, counts
and report, safe defaults without developer-path access, recursive selection,
and non-empty recorded read and entry populations. Fixture mutation is limited
to guarded temporary directories; production-composition tests do not mutate
the real filesystem.

## Measured Quality State

The clean full audit at d6cb5939ccce reports:

- Absolute L0: 8 of 8.
- 224 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 39 direct external sites outside five declared adapters of 56
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 41 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: tree b52698f0308dd83db251eea8edd980d340d7d907, status
  SHA-256 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d,
  and zero dirty paths.

The clean report is
/private/tmp/ply-ide-read-clean-audit-d6cb593/scorecard.json and the focused
report is /private/tmp/ply-ide-read-focused-audit/scorecard.json. Generated
ignored compatibility reports were removed before the clean audit; all
authoritative audit output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. The adapter ReadDir signature deliberately retains []fs.FileInfo because
   the migrated behavior depends on IsDir and ioutil.ReadDir's existing sorted
   result and metadata contract.
2. Q1.3 moved from 40/56 to 39/56: the direct IDE read disappeared and the new
   adapter implementation contributes the corresponding production effect, so
   the population remains 56.
3. The non-recursive IDE production wrapper alone selects filesystem.System().
   The recursive branch and its existing behavior were not refactored.
4. The loop body was not coalesced or repaired. Its exact independent checks,
   substring matching, ordered logAndDelete calls, paths, counts, report, and
   errors are preserved.
5. Every complete filesystem double needed an explicit ReadDir method; the new
   IDE double rejects every unrelated adapter operation.
6. The two new system-sorting fixture writes use the central guarded writer,
   moving Q0.6 safe-writer sites from 20 to 22 while unsafe writes stay zero.
7. GitCloudConfig.GitHookFiles is the next isolated flow. It can reuse ReadDir
   without extending the adapter or complete doubles, nominally moving Q1.3 to
   38/55.
8. GitHookFiles constructs root with
   file.Path("%s/%s", gitCfg.Implementation().Dir(), path), iterates delivered
   entries, excludes directories, and returns only names. Preserve nil results
   for empty/all-directory populations and exact initial read errors.
9. Examples has a separate direct ReadDir, Templates has a direct Walk, and
   every other config operation remains outside the next move.
10. .quality/inventory is baseline-checksum-bound. Do not relabel seams or
    claim P5 mutation coverage.
11. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are under /private/tmp/ply-quality-tools while that
    temporary directory remains available.
12. The standalone launcher and one later make-test run each encountered the
    known partial-raw-log signal-fixture timing flake. Immediate complete
    reruns passed all 62 controls, as did the full preflight. Never weaken or
    skip that contract.

## Next Objective

Move only GitCloudConfig.GitHookFiles(path string) ([]string, error)'s direct
ioutil.ReadDir(root) through the existing filesystem ReadDir operation.

Start red with config recording contracts for the complete dependency,
production system selection, exact receiver-derived root, delivered FileInfo
order and metadata, exact read errors, non-directory selection, ignored
directories, exact filename-only ordered results, nil results for empty and
all-directory populations, safe zero behavior without developer-path access,
and non-empty read and entry populations. Production-composition tests must not
mutate the real filesystem.

Keep GitHookFiles and CloudConfig public signatures unchanged. Production must
select filesystem.System() only for the private Git-hook-files composition.
Reuse adapter ReadDir without extending the adapter or complete doubles. Do not
change the root expression or loop body, Examples, Templates, another config
method or caller, inventory, public API, another production effect, or enter
HTTP/process, clock/server, P4, P5, or later roadmap work.

Expected direction is nominally Q1.3 38/55 with Q0.6, Q1.2, Q1.4, and exact
Q2.1 held. Regenerate the structured result and require zero comparable
ratchet regressions.

## Verification Notes

Completed from implementation commit d6cb5939ccce:

- Red evidence: focused compilation failed only on absent adapter ReadDir,
  private non-recursive IDE dependency/system composition, and the changed
  helper signature.
- Focused filesystem, file, command caller, and adjacent complete-double
  package tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls on the immediate
  complete rerun; the initial run hit the known signal-fixture timing flake.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test and make test-install: PASS on the immediate complete rerun; the
  first make-test run hit the same known signal-fixture timing flake.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 39/56.
- Clean full audit: exit 1, 16 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: d6cb5939cccefc3fd7ae8e80ef81c6bd5e922342
  (quality: route IDE directory read through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete config code/tests and GitHookFiles callers, the
CloudConfig interface and implementations/doubles, the filesystem adapter and
relevant doubles, and the named audits before editing. Confirm branch, HEAD,
status, reciprocal links, and ./codex-dev-start.sh --check. Begin red for only
the Git-hook directory-read boundary, then finish with one implementation
commit and one separate handoff-only commit.

Stop before Examples, Templates, another config or filesystem operation,
another adapter operation or family, Q1.4 expansion, P4, mutation, Docker,
cloud distribution, or publication. Stop on API/CLI change, comparable ratchet
regression, full-audit exit 2, or failure to isolate implementation from
generated and handoff-only state.
