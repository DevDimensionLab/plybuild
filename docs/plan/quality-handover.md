# Quality Upgrade Handover

Generated: 2026-08-25T08:59:31+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 8d493455a98f.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, and 8d49345 in roadmap order. The
separate operational continuity implementation is 1b85711 and changes no Go
quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T085931+0200-migrate-config-project-write.md.
The template-markdown-write predecessor is answered history, and the
reciprocal archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 32 Preserved

Exported SaveTemplateListMarkdown(gitCfg config.CloudConfig,
markdownDocument string) (string, error) retains its signature and remains the
production entry. Only its private complete markdown-write composition selects
filesystem.System(). The exact receiver-derived README path,
[]byte(markdownDocument), and mode 0644 pass through the existing
filesystem.WriteFile operation. No filesystem adapter operation,
templateCopyDependencies, or complete recording double changed.

The private composition retains the exact
gitCfg.Implementation().Dir() + "/" + TemplatesDir + "/README.md" evaluation.
It does not clean, join, normalize, preflight, create directories, retry,
inspect permissions, validate, or log. It makes one write attempt and returns
the complete readmePath with the exact dependency error on success and
failure. List rendering, the filtered walk, merge sequencing, getIgnores,
copy, search, render, config, file, Maven, and every caller remain unchanged.

Four new top-level template contracts bring the suite to 254 tests. They prove
complete production selection and delivery, exact receiver and path
evaluation, non-empty write populations, empty, non-ASCII, and arbitrary
document bytes, mode 0644, one write, full returned paths, exact error identity,
and safe defaults without developer-path access. The production composition
tests do not mutate the real filesystem.

## Measured Quality State

The clean full audit for the committed move-32 handoff reports:

- Absolute L0: 8 of 8.
- 254 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 22 guarded safe-writer sites, 17 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 34 direct external sites outside five declared adapters of 51
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 47 Markdown files.
- Acceptance: 4 of 4 host flows pass.
- Full audit: expected exit 1 for 15 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The clean report is
/private/tmp/ply-markdown-clean-audit-final/scorecard.json and the focused
report is /private/tmp/ply-markdown-focused-final.tBowLx/scorecard.json. All
authoritative audit output and tool caches stayed outside the measured tree.

## Decisions And Learned Facts

1. SaveTemplateListMarkdown reuses the adapter's established WriteFile
   operation; this move does not extend the adapter or any complete filesystem
   recording double.
2. Q1.3 moved from 35/52 to 34/51 because the direct markdown write disappeared
   while the already-counted adapter operation was reused.
3. The exported production entry alone selects filesystem.System() for its
   private complete composition; templateCopyDependencies and every caller
   remain unchanged.
4. The receiver-derived path expression, byte conversion, mode, single attempt,
   returned path, and exact error are preserved without normalization or
   additional effects.
5. Empty input becomes a non-nil empty byte slice at the dependency boundary;
   invalid UTF-8 and non-ASCII bytes represented by the Go string pass through
   unchanged.
6. The recording double rejects every unrelated adapter operation and rejects
   an empty write population.
7. ProjectConfiguration.WriteTo is the next isolated flow. It can reuse
   WriteFile without extending the adapter or complete doubles, nominally
   moving Q1.3 to 33/50.
8. ProjectConfiguration.WriteTo logs the exact target first, evaluates
   json.MarshalIndent(config, "", "    "), returns a marshal error before any
   write, then writes the exact target and serialized bytes with mode 0644.
9. The next move must preserve nil receivers and JSON bytes and must not change
   logging, serialization, ProjectConfig, project initialization,
   SortAndWritePom, cloud config, template, file, Maven, or command sequencing.
10. .quality/inventory is baseline-checksum-bound. Do not relabel seams or
    claim P5 mutation coverage.
11. Supply APIDIFF and GOLANGCI_LINT as environment variables for preflight,
    and keep GOLANGCI_LINT_CACHE and GOCACHE in writable external directories.
    The pinned tools are under /private/tmp/ply-quality-tools while that
    temporary directory remains available.
12. The complete preflight, make-test target, and independent launcher runs
    passed all 62 launcher controls without a retry.

## Next Objective

Move only (*config.ProjectConfiguration).WriteTo(targetFile string) error's
direct os.WriteFile(targetFile, data, 0644) through the existing filesystem
WriteFile operation.

Start red with config recording contracts for the complete dependency,
production system selection, exact caller target, non-empty write populations,
exact JSON bytes for nil and representative receivers, mode 0644, one attempt,
exact error identity, unchanged receiver state, and safe zero behavior without
developer-path access. Production-composition tests must not mutate the real
filesystem.

Preserve the exported method, ProjectConfig, and every caller. Production must
select filesystem.System() only for a private complete project-config-write
composition. Reuse adapter WriteFile without extending the adapter or complete
doubles. Do not change logging, json.MarshalIndent, marshal error handling,
target bytes, mode, ProjectConfigPath, projectConfigFile, SortAndWritePom,
project initialization, cloud config, template/file/Maven behavior, command
sequencing, inventory, public API, another production effect, or enter
clock/server, P4, P5, or later roadmap work.

Expected direction is nominally Q1.3 33/50 with Q0.6, Q1.2, Q1.4, and exact
Q2.1 held. Regenerate the structured result and require zero comparable
ratchet regressions.

## Verification Notes

Completed from implementation commit 8d493455a98f:

- Red evidence: focused compilation failed only on the absent private markdown
  write dependency, system composition, and private composition function.
- Focused filesystem, template, config, file, Maven, and command-caller package
  tests: PASS.
- Standalone /bin/bash launcher contract: PASS all 62 controls.
- Make preflight meta-contracts and full preflight: PASS, including all 15
  audit meta-controls and all 62 launcher controls.
- API, CLI, and subprocess compatibility: PASS.
- make test and make test-install: PASS without a retry.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME go test ./... -count=2: PASS.
- Focused audit: exit 1, four improved, two held, zero regressed, one
  non-comparable; Q1.3 is exactly 34/51.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Implementation commit: 8d493455a98ffd484705c68f007749a0dc93addd
  (quality: route template markdown write through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete config code/tests and project-write callers, the
filesystem adapter and relevant doubles, and the named audits before editing.
Confirm branch, HEAD, status, reciprocal links, and
./codex-dev-start.sh --check. Begin red for only the project-config write
boundary, then finish with one implementation commit and one separate
handoff-only commit.

Stop before JSON or logging refactors, SortAndWritePom, project initialization,
cloud config, template/file/Maven behavior, command composition, another
config or filesystem operation, another adapter operation or family, Q1.4
expansion, P4, mutation, Docker, cloud distribution, or publication. Stop on
API/CLI change, comparable ratchet regression, full-audit exit 2, or failure to
isolate implementation from generated and handoff-only state.
