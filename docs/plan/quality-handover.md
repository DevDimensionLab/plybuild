# Quality Upgrade Handover

Generated: 2026-08-25T07:53:01+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: dfcfa75da0cd.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, and dfcfa75 in roadmap order. The separate
operational continuity implementation is 1b85711 and changes no Go quality
denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T075301+0200-migrate-template-filtered-walk.md.
The Templates predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 30 Preserved

Exported GitCloudConfig.Templates() (templates []CloudTemplate, err error)
retains its signature and remains the production entry. Only its private
complete composition selects filesystem.System(). The exact receiver-root
expression `root := file.Path("%s/templates",
gitCfg.Implementation().Dir())` is unchanged, and that exact root and the
existing callback pass through the existing filesystem.Walk operation. No
filesystem adapter operation or complete recording double changed.

The callback body remains unchanged. An incoming non-nil walk error is
suppressed with nil before info is dereferenced. With a nil incoming error,
the exact info.Name() comparisons select only ply.json and co-pilot.json.
Matched paths retain the exact splits, relative template name,
InitProjectFromDirectory(relPath[0]) call, logged and returned project-load
error, and CloudTemplate append. System lexical order, delivered callback
paths and metadata, ignored entries, partial ordered results, and the exact
final walk error remain unchanged. CloudConfig, Implementation, DirConfig.Dir,
file.Path, Templates callers, Examples, GitHookFiles, project loading, and
every other config operation remain unchanged.

Seven new top-level config contracts bring the suite to 243 tests. They prove
complete production selection and delivery, exact receiver-derived templates
root, non-empty walk-root and callback populations, delivered callback order,
paths, metadata and incoming errors, incoming-error suppression, exact current
and legacy filename matching, ignored entries, exact relative names and loaded
projects, exact project-load and walk errors, partial ordered results, nil
results, and safe defaults without developer-path access. The production
composition tests do not mutate the real filesystem.

## Measured Quality State

The clean full audit at dfcfa75da0cd reports:

- Absolute L0: 8 of 8.
- 243 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 36 direct external sites outside five declared adapters of 53
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 44 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: tree 77a24e64e0e4fe2797cdb5d380e9ae3d4231870d, status
  SHA-256 6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d,
  and zero dirty paths.

The clean report is
/private/tmp/ply-templates-clean-audit/scorecard.json and the focused report is
/private/tmp/ply-templates-focused-audit/scorecard.json. All authoritative
audit output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. Templates reuses the adapter's established Walk operation; this move does
   not extend the adapter or any complete filesystem recording double.
2. Q1.3 moved from 37/54 to 36/53 because the direct Templates walk
   disappeared while the already-counted adapter operation was reused.
3. The Templates production wrapper alone selects filesystem.System().
   GitHookFiles and Examples retain their completed compositions.
4. The templates-root expression and callback body were not refactored. Exact
   receiver composition, system lexical order, callback observations,
   matching, project loading, names, partial results, and errors are preserved.
5. The incoming-error guard retains short-circuit behavior: non-nil errors
   return nil without dereferencing info or propagating the incoming error.
6. The legacy and project-load-error fixtures exercise real unchanged project
   loading while the injected walk prevents production filesystem mutation.
7. The recording double rejects every unrelated adapter operation and records
   both non-empty walk-root and delivered-callback populations.
8. template.filteredFilesFromTemplate is the next isolated flow. It can reuse
   Walk without extending the adapter or complete doubles, nominally moving
   Q1.3 to 35/52.
9. filteredFilesFromTemplate passes sourceDir directly to filepath.Walk. Its
   callback ignores its incoming err argument and calls info.IsDir() first;
   system walk errors with nil info therefore retain the existing panic. With
   non-nil info, even a non-nil incoming error is processed by metadata and
   filter logic. Preserve both observations exactly.
10. The filtered callback's slash-based root detection, ordered filters,
    non-root pom.xml exception, substring matching, .render exception, logging,
    ordered results, partial results, and final walk error must remain exact.
11. .quality/inventory is baseline-checksum-bound. Do not relabel seams or
    claim P5 mutation coverage.
12. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are under /private/tmp/ply-quality-tools while that
    temporary directory remains available.
13. The first complete make-test attempt hit the previously documented
    launcher signal-interruption partial-raw-log timing flake; its immediate
    full rerun passed. The implementation audit initially exposed an unrelated
    Q3.2 heuristic pass caused only by an ignored fixture named README.md; that
    inert entry was renamed before the final commit and the exact clean audit
    restored the expected 16 findings.

## Next Objective

Move only template.filteredFilesFromTemplate(sourceDir string, filter
[]string)'s direct filepath.Walk(sourceDir, callback) through the existing
filesystem Walk operation.

Start red with template recording contracts for the complete dependency,
production system selection, exact sourceDir root, delivered callback order,
paths, metadata and incoming errors, exact directory handling and metadata
observations, nil-info incoming-error panic behavior, non-nil-info incoming
error processing, exact slash-based root detection, ordered ignore behavior,
the non-root pom.xml exception, substring matches, the .render exception,
exact logging, ordered and nil results, exact final walk errors, partial
results, safe zero behavior without developer-path access, and non-empty
walk-root and callback populations. Production-composition tests must not
mutate the real filesystem.

Keep filteredFilesFromTemplate private and preserve MergeTemplate, merge,
templateCopyDependencies, and every public signature. Production must select
filesystem.System() only for a private complete filtered-walk composition.
Reuse adapter Walk without extending the adapter or complete doubles. Do not
change the callback body, merge sequencing, getIgnores, copy/search/render
behavior, SaveTemplateListMarkdown, another template/config/file operation or
caller, inventory, public API, another production effect, or enter clock/server,
P4, P5, or later roadmap work.

Expected direction is nominally Q1.3 35/52 with Q0.6, Q1.2, Q1.4, and exact
Q2.1 held. Regenerate the structured result and require zero comparable
ratchet regressions.

## Verification Notes

Completed from implementation commit dfcfa75da0cd:

- Red evidence: focused compilation failed only on the absent private Templates
  dependency, system composition, and private composition method.
- Focused filesystem, config, and command-caller package tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test, make test-install, and make test-agent-start: PASS; the complete
  make-test pass required one rerun after the known signal timing flake.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 36/53.
- Clean full audit: exit 1, 16 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: dfcfa75da0cd373e69255e496742085479392c1e
  (quality: route Templates walk through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete template code/tests and filtered-walk callers,
the filesystem adapter and relevant doubles, and the named audits before
editing. Confirm branch, HEAD, status, reciprocal links, and
./codex-dev-start.sh --check. Begin red for only the filtered template walk
boundary, then finish with one implementation commit and one separate
handoff-only commit.

Stop before merge composition or sequencing, copy/search/render behavior,
getIgnores, SaveTemplateListMarkdown, another template/config/file or
filesystem operation, another adapter operation or family, Q1.4 expansion,
P4, mutation, Docker, cloud distribution, or publication. Stop on API/CLI
change, comparable ratchet regression, full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
