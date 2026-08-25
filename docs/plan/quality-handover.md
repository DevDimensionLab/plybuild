# Quality Upgrade Handover

Generated: 2026-08-25T08:24:07+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 170b0aee0f1f.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, and 170b0ae in roadmap order. The
separate operational continuity implementation is 1b85711 and changes no Go
quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T082407+0200-migrate-template-markdown-write.md.
The filtered-walk predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 31 Preserved

Private filteredFilesFromTemplate(sourceDir string, filter []string) (files
[]string, err error) retains its signature and remains the production entry.
Only its private complete composition selects filesystem.System(). The exact
supplied sourceDir and unchanged callback pass through the existing
filesystem.Walk operation. No filesystem adapter operation or complete
recording double changed.

The callback body remains unchanged. It ignores its incoming error and first
calls info.IsDir(); nil info therefore retains the existing panic, while
non-nil info with a non-nil incoming error is processed normally. Directories
return nil without reading Name, evaluating filters, logging, or appending.
Non-directories retain the exact slash-based ReplaceAll root test, delivered
filter order, non-root pom.xml continue, substring matching, global .render
exception, debug logging, ordered append, partial results, nil results, and
exact final walk error. System lexical order, callback paths, metadata, and
incoming errors remain unchanged. Merge, getIgnores, copy, search, render,
config, file, Maven, and every caller remain unchanged.

Seven new top-level template contracts bring the suite to 250 tests. They
prove complete production selection and delivery, the exact supplied root,
non-empty walk-root and callback populations, delivered callback order, paths,
metadata and incoming errors, directory metadata observations, nil-info panic,
non-nil-info error processing, exact filtering and logging, ordered and nil
results, exact final errors, partial results, and safe defaults without
developer-path access. The production composition tests do not mutate the real
filesystem.

## Measured Quality State

The clean full audit at 170b0aee0f1f reports:

- Absolute L0: 8 of 8.
- 250 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 35 direct external sites outside five declared adapters of 52
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 45 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: tree f0ddd6353a067f45eaf675ff2b92070790addc54, status
  SHA-256 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d,
  and zero dirty paths.

The clean report is
/private/tmp/ply-filtered-clean-audit-170b0ae/scorecard.json and the focused
report is /private/tmp/ply-filtered-focused-audit/scorecard.json. All
authoritative audit output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. filteredFilesFromTemplate reuses the adapter's established Walk operation;
   this move does not extend the adapter or any complete filesystem recording
   double.
2. Q1.3 moved from 36/53 to 35/52 because the direct filtered walk disappeared
   while the already-counted adapter operation was reused.
3. The private production entry alone selects filesystem.System();
   templateCopyDependencies and every caller remain unchanged.
4. The exact sourceDir argument and callback body were not refactored. System
   lexical order, callback observations, slash-based root detection, ordered
   filters, logging, results, and errors are preserved.
5. The callback deliberately retains the nil-info incoming-error panic and
   processes a non-nil info despite a non-nil incoming error. Do not add an
   incoming-error branch or nil guard in later work.
6. The recording double rejects every unrelated adapter operation and records
   both non-empty walk-root and delivered-callback populations.
7. SaveTemplateListMarkdown is the next isolated flow. It can reuse WriteFile
   without extending the adapter or complete doubles, nominally moving Q1.3
   to 34/51.
8. SaveTemplateListMarkdown currently computes exactly
   `gitCfg.Implementation().Dir() + "/" + TemplatesDir + "/README.md"`, writes
   `[]byte(markdownDocument)` with mode 0644 once, and returns that full path
   with the exact write error on both success and failure.
9. The next move must not use filepath.Join, file.Path, clean or normalize the
   path, create directories, change ListAsMarkdown, or alter another template,
   config, file, Maven, merge, copy, search, or render operation.
10. .quality/inventory is baseline-checksum-bound. Do not relabel seams or
    claim P5 mutation coverage.
11. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are under /private/tmp/ply-quality-tools while that
    temporary directory remains available.
12. One complete make-test attempt hit the previously documented nested
    launcher signal-interruption partial-raw-log timing flake; its immediate
    full rerun passed. The complete preflight and independent launcher runs
    also passed all 62 controls.

## Next Objective

Move only template.SaveTemplateListMarkdown(gitCfg config.CloudConfig,
markdownDocument string) (string, error)'s direct
os.WriteFile(readmePath, []byte(markdownDocument), 0644) through the existing
filesystem WriteFile operation.

Start red with template recording contracts for the complete dependency,
production system selection, the exact receiver-derived README path,
non-empty recorded write populations, exact document bytes including empty and
non-ASCII content, exact 0644 mode, one write attempt, exact returned path on
success and failure, exact error identity, and safe zero behavior without
developer-path access. Production-composition tests must not mutate the real
filesystem.

Preserve the exported signature and every caller. Production must select
filesystem.System() only for a private complete markdown-write composition.
Reuse adapter WriteFile without extending the adapter,
templateCopyDependencies, or complete doubles. Do not change the exact path
expression, []byte conversion, mode, return statement, ListAsMarkdown,
createTemplateListRenderingModel, filteredFilesFromTemplate, merge sequencing,
getIgnores, copy/search/render behavior, another template/config/file method
or caller, inventory, public API, another production effect, or enter
clock/server, P4, P5, or later roadmap work.

Expected direction is nominally Q1.3 34/51 with Q0.6, Q1.2, Q1.4, and exact
Q2.1 held. Regenerate the structured result and require zero comparable
ratchet regressions.

## Verification Notes

Completed from implementation commit 170b0aee0f1f:

- Red evidence: focused compilation failed only on the absent private filtered
  walk dependency, system composition, and private composition function.
- Focused filesystem, template, config, file, Maven, and command-caller package
  tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test and make test-install: PASS; the complete make-test pass required
  one rerun after the known nested signal timing flake.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 35/52.
- Clean full audit: exit 1, 16 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: 170b0aee0f1f4d13c28a89d19cbbd1ce9c0fd7e8
  (quality: route filtered template walk through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete template code/tests and markdown-write callers,
the filesystem adapter and relevant doubles, and the named audits before
editing. Confirm branch, HEAD, status, reciprocal links, and
./codex-dev-start.sh --check. Begin red for only the template markdown write
boundary, then finish with one implementation commit and one separate
handoff-only commit.

Stop before list rendering, filtered walking, merge composition or sequencing,
copy/search/render behavior, getIgnores, another template/config/file or
filesystem operation, another adapter operation or family, Q1.4 expansion,
P4, mutation, Docker, cloud distribution, or publication. Stop on API/CLI
change, comparable ratchet regression, full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
