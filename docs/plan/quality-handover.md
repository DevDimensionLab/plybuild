# Quality Upgrade Handover

Generated: 2026-08-25T16:33:21+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: d982f6334944.
- After launch, obtain the session head with git rev-parse --short=12 HEAD; the
  restart commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438,
527a8b9, 886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc,
85f4c2b, 0e10282, 9fcdfb5, and d982f63 in roadmap order. The separate
operational continuity implementation is 1b85711 and changes no Go quality
denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T163321+0200-migrate-shell-unzip-directory-creation.md.
The Spring archive-path predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 46 Preserved

The complete internal filesystem interface now includes only one additional
`WorkingDirectory() (string, error)` operation. Its forwarding helper returns
the exact safe-zero pair `"", filesystem.ErrNoFilesystem`, otherwise makes one
dependency attempt and returns the exact string and error. The system
implementation directly returns `os.Getwd()` without rewriting or wrapping.
Every earlier filesystem operation and caller remains unchanged, and all
pre-existing complete doubles reject the unrelated operation.

Private `pkg/spring.archivePath() (string, error)` remains the production entry
and selects `filesystem.System()` through one private complete composition. A
working-directory error returns the same empty path and exact error before the
direct `time.Now().Unix()` call. Success preserves arbitrary directory bytes
in `file.Path("%s/spring-%d.zip", curDir, now)`. Every Spring caller,
discovery, request, download, unzip, delete, and demo-file behavior remains
unchanged; no clock seam was added.

Eight focused contracts bring the suite to 324 tests. They prove exact adapter
value/error forwarding, one attempt, direct system behavior, exact safe zero,
complete Spring production selection, arbitrary directory composition, direct
Unix filename construction, exact error identity and empty error result,
rejection of empty recording populations, and no unrelated filesystem
operation.

## Measured Quality State

The clean full audit from implementation commit d982f63 reports:

- Absolute L0: 8 of 8.
- 324 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 24 guarded safe-writer sites, 19 write and 5 copy, and zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 15 direct external sites outside five declared adapters of 36
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 60 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
/private/tmp/ply-spring-archive-tools/reports/full-clean-d982f63/scorecard.json
and the focused implementation report is
/private/tmp/ply-spring-archive-tools/reports/focused-dirty/scorecard.json.
Compatibility reports, tool binaries, Go and linter caches, audit output, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. A working-directory read is a distinct filesystem operation. Exact caller
   bytes and errors pass through the adapter without cleaning, fallback, retry,
   logging, or wrapping.
2. The Spring production composition preserves Getwd-before-time ordering. An
   error exits before the unchanged direct clock selection and path formatting.
3. The former Spring caller site left Q1.3 while the system adapter remained in
   the production population, producing the predicted 15/36 result.
4. Clock and server adapter creation belongs to queued P4. Do not move the
   remaining Spring or Kibana clock call or either webservice server operation
   during P3.
5. The remaining direct process constructions are three sites in
   `cmd/plugin_diagrams.go` and one generic site in `pkg/shell/command.go`.
   Their exported structurizr command signatures and exact
   `exec.Cmd.String()` logging make them wider than an isolated filesystem
   mutation move.
6. The Q1.3 scanner also attributes `pkg/maven/command.go:20` to filesystem
   even though the source calls `process.Execute`; investigate the scanner
   receipt before treating that line as production debt.
7. The narrow next existing-adapter move is the two direct `os.MkdirAll` calls
   in `pkg/shell.Unzip`. They share exact `os.ModePerm`, established ordering,
   and an existing `filesystem.MkdirAll` operation.
8. Route only those two directory creations through a private complete
   filesystem composition. Leave direct `zip.OpenReader`, `os.OpenFile`,
   `io.Copy`, file closes, archive-entry reads, and every other shell effect for
   later coherent moves.
9. Nominal Q1.3 direction is 13/34 when the two caller sites leave and no new
   system operation is added. Q0.6, Q1.1, Q1.2, Q1.4, and exact Q2.1 should
   hold; regenerate every value.
10. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams,
    make mutation harnesses executable, or claim P4/P5 coverage.
11. Full preflight uses APIDIFF and golangci-lint v2.12.2 from
    `/private/tmp/ply-open-browser-tools/bin`; put new reports and caches under
    an external move-specific root.

## Next Objective

Move only the two direct `os.MkdirAll` calls in `pkg/shell.Unzip` through the
existing filesystem adapter operation.

Start red with focused private shell contracts for complete system dependency
selection, both exact path forms, exact `os.ModePerm`, one attempt per reached
directory creation, exact error identity and legacy partial filename results,
safe zero behavior at the directory-creation boundary, non-empty recorded
populations, and no unexpected filesystem operation. Preserve archive-open,
entry ordering, traversal validation, filename append ordering, directory
short-circuiting, parent selection, file flags and mode, copy and close
behavior, and every error precedence.

Add no filesystem interface operation. Keep public
`Unzip(string, string) ([]string, error)` and every caller unchanged. Select
`filesystem.System()` only through a private complete shell composition and
replace only the two direct `os.MkdirAll` calls with the established
`filesystem.MkdirAll` helper. Do not move archive open, file open, copy, close,
process execution, Spring, clock, server, P4, or P5 work.

Expected nominal direction is Q1.3 13/34 with Q0.6 at 24 guarded sites, Q1.1
6/25, Q1.2 zero, Q1.4 7/8, and exact Q2.1 0/8. Require zero comparable ratchet
regressions and regenerate exact values.

## Verification Notes

Completed from implementation commit d982f6334944a4e74c90c90823757bcda2b720ec:

- Valid red evidence: isolated focused compilation failed only on absent
  `WorkingDirectory`, `archivePathDependencies`,
  `systemArchivePathDependencies`, and `archivePathWithDependencies`. An
  earlier host-cache permission failure was rejected as red evidence.
- Focused Spring/filesystem and relevant process, command, context, config,
  HTTP, Maven, shell, structurizr, profile, browser, tips, file, template,
  Bitbucket, Kibana, local-config, and caller package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls after one
  documented signal-fixture timing flake and an immediate complete rerun.
- Make preflight meta-contracts, pinned linter with zero issues, one complete
  preflight, and all 15 audit meta-controls: PASS with external tools, reports,
  and caches. The first preflight attempt was rejected after the host denied
  the default linter cache; the isolated-cache rerun passed completely.
- API, CLI, and subprocess compatibility: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their meta-contracts:
  PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 15/36.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 61
  Markdown files and zero ratchet regressions.
- Implementation commit: d982f6334944a4e74c90c90823757bcda2b720ec
  (refactor: route spring archive cwd through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete shell implementation/tests/callers, complete
filesystem adapter/tests and every complete double, Spring archive/download
contracts, relevant file, HTTP, process, config, Maven, structurizr, context,
Bitbucket, Wpost, local-config, browser, profile, and audit contracts before
editing. Confirm branch, HEAD, status, reciprocal links, and
`./codex-dev-start.sh --check`. Begin red only for the private shell unzip
directory-creation composition, then finish with one implementation commit and
one separate handoff-only commit.

Stop before changing archive open, path traversal validation, entry ordering,
file open flags or mode, copy or close behavior, process construction, another
filesystem operation or caller, Spring, clock, server, public API, inventory,
Q1.4, P4, mutation, Docker, distribution, or publication. Stop on API/CLI
change, comparable ratchet regression, authoritative full-audit exit 2, or
failure to isolate implementation from generated and handoff-only state.
