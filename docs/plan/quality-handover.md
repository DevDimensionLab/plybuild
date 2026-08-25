# Quality Upgrade Handover

Generated: 2026-08-25T04:15:54+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 224a691e1e43.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, and 224a691 in roadmap order. The separate
operational continuity implementation is 1b85711 and changes no Go quality
denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T041554+0200-migrate-file-find-first.md.
The render-create predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 23 Preserved

Exported file.Render(inputFilePath string, outputFilePath string, r interface{})
error retains its signature and observable behavior. It is the production
wrapper around a private dependency containing one complete
filesystem.Dependencies value, and production selects filesystem.System().
No function-valued effect dependency, adapter operation, adapter family, or
public API was added.

The private helper delegates the exact input path to existing adapter ReadFile,
preserves its exact error, and performs the same template.Must parse before
output creation. Only after successful read and parse does it delegate the
exact output path to existing adapter Create. Template execution preserves the
exact rendered bytes, create/truncation behavior, exact writer error, nil
success, and existing lack of a close.

Eight top-level recording contracts replaced the old worktree-mutating render
test. They prove complete dependency delivery, production system selection,
exact paths and read-create-write order, input and parse short-circuiting,
exact create and writer errors, exact bytes, nil success, no close, a safe zero
without developer-path access, and non-empty read/create/write populations.
They perform no real filesystem mutation. Both template callers and every
complete filesystem recording double remain unchanged.

Only Render's direct os.Create(outputFilePath) moved. ClearDir, Move, DeleteAll,
DeleteSingleFile, append-open, file-read, directory-create, file-create,
file-overwrite, file-existence, Bitbucket, Wpost, supervisor, inventory,
public API, callers, and every P2A contract remain unchanged.

## Measured Quality State

The clean full audit at 224a691e1e43 reports:

- Absolute L0: 8 of 8.
- 194 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 18 guarded safe-writer sites, 13 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 43 direct external sites outside five declared adapters of 58
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 37 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: tree de654906116dd91700e1207c7d4f8a6fd1b8fd92, status
  SHA-256 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d,
  and zero dirty paths.

The clean report is /private/tmp/ply-render-clean-audit-224a691 and the focused
report is /private/tmp/ply-render-focused-20260825. Compatibility, preflight,
Go-gate, acceptance, and hermetic outputs also stayed under /private/tmp.

## Decisions And Learned Facts

1. Render carries one complete filesystem dependency. Reused ReadFile keeps the
   input boundary and makes the zero value safe before any create.
2. template.Must remains before Create. Invalid syntax still panics with the
   input path as template name and creates no output.
3. The adapter file intentionally remains open on success and execution error.
4. Q1.3 moved from 44/59 to 43/58 because direct render creation disappeared
   and already-counted system Create was reused.
5. file.FindFirst is the next isolated flow. Its one direct filepath.Walk site
   should become one adapter Walk site, nominally moving Q1.3 to 42/58.
6. FindFirst deliberately ignores callback fi and errIn, uses only
   strings.HasSuffix, assigns the exact first match, returns io.EOF to stop,
   and normalizes only a final walk error equal to io.EOF.
7. Function-valued effect dependencies fail Q1.3 closed. Use resolvable
   interfaces and complete dependency values.
8. .quality/inventory is baseline-checksum-bound. Do not relabel seams or claim
   P5 mutation coverage.
9. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight, and
   keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
10. The ignored .agent-task/current.md compatibility mirror is absent.
11. Supervisor success requires process, structured-stream, and committed
    repository evidence; final prose is observable only.
12. The partial-raw-log signal fixture can flake. The first complete make test
    run did; the standalone launcher, preflight, and immediate complete rerun
    all passed 62 controls. Never weaken or skip this contract.

## Next Objective

Move only exported file.FindFirst(fileSuffix string, dir string) (result string,
err error)'s direct filepath.Walk(dir, callback) behind one narrow,
zero-value-safe Walk operation in the existing filesystem adapter.

Start red with adapter and file recording contracts for the complete
dependency, production system selection, exact root and callback delivery,
callback order, ignored metadata/errors, suffix decisions, first-match result
and stop, no-match empty result, exact non-EOF walk errors and partial result,
final io.EOF normalization, safe zero behavior without developer path access,
and non-empty walk/callback populations. File contracts must not mutate the
real filesystem.

Keep FindFirst's signature and callback body unchanged. The system adapter
must delegate once to filepath.Walk, and its zero must return
filesystem.ErrNoFilesystem before root access. Update every complete
filesystem recording double mechanically. Do not move FindAll, GrepRecursive,
IDE cleanup, or another operation, and do not enter HTTP/process, config,
clock/server, P4, P5, or later roadmap work.

Expected direction is nominally Q1.3 42/58 with Q1.2, Q1.4, and exact Q2.1
held. Regenerate the structured result and require zero comparable ratchet
regressions.

## Verification Notes

Completed from implementation commit 224a691e1e43:

- Red evidence: focused file compilation failed on the absent private render
  dependency constructor and helper.
- Focused filesystem, file, template, Bitbucket, and HTTP tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls.
- API, CLI, subprocess compatibility: PASS.
- make test: first run hit the known signal flake; immediate full rerun PASS.
- make test-install, uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable.
- Clean full audit: exit 1, 16 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable.
- Implementation commit: 224a691e1e43112fd8df36a7b4c579d60246fb5c
  (quality: route render create through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, all file code/tests and FindFirst callers, the adapter
and complete doubles, and the named audits before editing. Confirm branch,
HEAD, status, reciprocal links, and ./codex-dev-start.sh --check. Begin red for
only FindFirst's walk boundary, then finish with one implementation commit and
one separate handoff-only commit.

Stop before FindAll, recursive grep, IDE cleanup, another adapter family,
Q1.4 expansion, P4, mutation, Docker, cloud, distribution, or publication.
Stop on API/CLI change, comparable ratchet regression, full-audit exit 2, or
failure to isolate implementation from generated and handoff-only state.
