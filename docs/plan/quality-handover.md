# Quality Upgrade Handover

Generated: 2026-08-25T07:12:58+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: cb94f8162176.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, and cb94f81 in roadmap order. The separate operational
continuity implementation is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T071258+0200-migrate-config-templates-walk.md.
The Examples predecessor is answered history, and the reciprocal archive graph
has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 29 Preserved

Exported GitCloudConfig.Examples() (templates []string, err error) retains its
signature and remains the production entry. Only its private complete
composition selects filesystem.System(). The exact receiver-root expression
`examplesDir := file.Path("%s/examples", gitCfg.Implementation().Dir())` is
unchanged, and that exact root passes through the existing filesystem.ReadDir
operation. No filesystem adapter operation or complete recording double
changed.

The loop body remains unchanged. It iterates entries in dependency-delivered
order, selects only `item.IsDir()`, and appends each exact `item.Name()`,
preserving directory-name-only ordered results. Production filename sorting
remains the adapter's ioutil.ReadDir behavior. An initial read error returns a
nil result and the exact error. Empty and all-file populations return a nil
result and nil error. CloudConfig, Implementation, DirConfig.Dir, file.Path,
Examples callers, GitHookFiles, Templates, and every other config operation
remain unchanged.

Six new top-level config contracts bring the suite to 236 tests. They prove
complete production selection and delivery, exact receiver-derived examples
root, delivered entry order and metadata, exact errors, directory
classification and ignored files, exact directory-name-only results, nil
results, safe defaults without developer-path access, and non-empty recorded
read and entry populations. The production-composition tests do not mutate the
real filesystem.

## Measured Quality State

The clean full audit at cb94f8162176 reports:

- Absolute L0: 8 of 8.
- 236 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 37 direct external sites outside five declared adapters of 54
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 43 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: tree 5d55e212d55578a937314fb38bfc52edc0246236, status
  SHA-256 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d,
  and zero dirty paths.

The clean report is
/private/tmp/ply-examples-clean-audit/scorecard.json and the focused report is
/private/tmp/ply-examples-focused-audit/scorecard.json. All authoritative audit
output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. Examples reuses the adapter's established ReadDir operation; this move does
   not extend the adapter or any complete filesystem recording double.
2. Q1.3 moved from 38/55 to 37/54 because the direct Examples read disappeared
   while the already-counted adapter operation was reused.
3. The Examples production wrapper alone selects filesystem.System().
   GitHookFiles retains its completed composition, while Templates retains its
   original direct walk.
4. The examples-root expression and loop body were not refactored. Exact
   receiver composition, directory classification, delivered order, names, nil
   results, and errors are preserved.
5. The recording double rejects every unrelated adapter operation and records
   both non-empty directory-read and delivered-entry populations.
6. The tests deliberately deliver unsorted metadata-rich entries. They prove
   dependency order while relying on the adapter's existing system-sorting
   contract for production order.
7. GitCloudConfig.Templates is the next isolated flow. It can reuse Walk
   without extending the adapter or complete doubles, nominally moving Q1.3 to
   36/53.
8. Templates constructs root with
   file.Path("%s/templates", gitCfg.Implementation().Dir()) and passes an inline
   filepath.WalkFunc. Preserve the callback body exactly, including suppression
   of incoming walk errors, exact current/legacy filename matching, path splits,
   project loading, relative names, partial ordered results, logged project
   errors, and final walk error.
9. InitProjectFromDirectory performs separate existing effects. Do not migrate,
   inject, or refactor them as part of the Templates walk move.
10. .quality/inventory is baseline-checksum-bound. Do not relabel seams or
    claim P5 mutation coverage.
11. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are under /private/tmp/ply-quality-tools while that
    temporary directory remains available.
12. This move's standalone launcher, full preflight, complete make-test,
    standalone audit-meta, and empty-HOME count-2 runs passed without a rerun.

## Next Objective

Move only GitCloudConfig.Templates() (templates []CloudTemplate, err error)'s
direct filepath.Walk(root, callback) through the existing filesystem Walk
operation.

Start red with config recording contracts for the complete dependency,
production system selection, exact receiver-derived templates root, delivered
callback order, paths, metadata, and incoming errors, exact suppression of
incoming walk errors, current and legacy project-file matching, ignored
entries, exact relative names and loaded projects, exact project-load and walk
errors, partial ordered results, safe zero behavior without developer-path
access, and non-empty walk-root and callback populations.
Production-composition tests must not mutate the real filesystem.

Keep Templates and CloudConfig public signatures unchanged. Production must
select filesystem.System() only for the private Templates composition. Reuse
adapter Walk without extending the adapter or complete doubles. Do not change
the templates-root expression or callback body, InitProjectFromDirectory,
Examples, GitHookFiles, another config method or caller, inventory, public API,
another production effect, or enter HTTP/process, clock/server, P4, P5, or
later roadmap work.

Expected direction is nominally Q1.3 36/53 with Q0.6, Q1.2, Q1.4, and exact
Q2.1 held. Regenerate the structured result and require zero comparable
ratchet regressions.

## Verification Notes

Completed from implementation commit cb94f8162176:

- Red evidence: focused compilation failed only on the absent private Examples
  dependency and system composition.
- Focused filesystem, config, and both command-caller package tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test, make test-install, and make test-agent-start: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 37/54.
- Clean full audit: exit 1, 16 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: cb94f816217601c0f7b2fdc7d0caef366a638afe
  (quality: route Examples read through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete config code/tests and Templates callers, the
CloudConfig interface and implementations/doubles, the filesystem adapter and
relevant doubles, and the named audits before editing. Confirm branch, HEAD,
status, reciprocal links, and ./codex-dev-start.sh --check. Begin red for only
the Templates walk boundary, then finish with one implementation commit and one
separate handoff-only commit.

Stop before InitProjectFromDirectory, another config or filesystem operation,
another adapter operation or family, Q1.4 expansion, P4, mutation, Docker,
cloud distribution, or publication. Stop on API/CLI change, comparable ratchet
regression, full-audit exit 2, or failure to isolate implementation from
generated and handoff-only state.
