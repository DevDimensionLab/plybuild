# Quality Upgrade Handover

Generated: 2026-08-25T17:15:00+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: b9fe2091e134.
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
85f4c2b, 0e10282, 9fcdfb5, d982f63, and b9fe209 in roadmap order. The
separate operational continuity implementation is 1b85711 and changes no Go
quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T171500+0200-migrate-shell-unzip-file-open.md.
The directory-creation predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 47 Preserved

Public `pkg/shell.Unzip(string, string) ([]string, error)` now selects one
private complete composition containing `filesystem.Dependencies`, with
`filesystem.System()` chosen only by the production entry. The established
unzip body remains behind that composition in its original order. Only the two
former direct `os.MkdirAll` calls now use
`filesystem.MkdirAll(dependencies.Files, ..., os.ModePerm)`.

Directory entries still append the exact joined path before one recursive-
directory attempt and continue immediately after success. File entries still
append the exact joined path before one attempt on `filepath.Dir(fpath)` and
return before file open on failure. Exact errors and partial filenames remain
unwrapped. Archive open and deferred close, entry traversal, path join, debug
logging, zip-slip validation and text, direct file open, direct entry open,
copy, both file closes, final results, and all error precedence remain
unchanged. No filesystem interface or adapter implementation changed.

Six focused shell contracts bring the suite to 330 tests. They prove complete
production selection, both exact path forms, exact `os.ModePerm`, attempt
counts, filename and entry ordering, exact failure identity and partial
results, safe-zero behavior at the earliest reached directory boundary, a
non-empty recorded population, and no unrelated adapter operation. Their one
zip fixture write uses the central guarded writer under `t.TempDir()`.

## Measured Quality State

The clean full audit from implementation commit b9fe209 reports:

- Absolute L0: 8 of 8.
- 330 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 25 guarded safe-writer sites, 20 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 13 direct external sites outside five declared adapters of 34
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 61 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
/private/tmp/ply-shell-unzip-tools/reports/clean-full-audit/scorecard.json and
the focused implementation report is
/private/tmp/ply-shell-unzip-tools/reports/focused-audit-final/scorecard.json.
Compatibility reports, tool binaries, Go and linter caches, audit output, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. Both unzip directory creations share the existing `MkdirAll` adapter
   operation, but their path inputs and short-circuit behavior are distinct and
   are recorded separately.
2. The public production composition is selected before the unchanged direct
   archive open; `filesystem.System()` itself performs no filesystem operation.
3. The two former caller sites left Q1.3 without adding a system operation,
   producing the predicted 13/34 result.
4. The guarded zip fixture adds one write-safe call site, so Q0.6 improves from
   24 to 25 sites while unsafe direct test writes remain zero.
5. The narrow next existing-adapter move is the direct `os.OpenFile` in the
   file-entry branch of `pkg/shell.Unzip`; the existing `filesystem.OpenFile`
   operation already preserves `*os.File`, flags, mode, result, and error.
6. Reuse the private unzip composition. Do not add another production
   composition or select `filesystem.System()` again inside the private body.
7. The file-parent `MkdirAll` remains earlier. An injected file-open failure
   must occur only after that exact successful attempt, return the current
   partial filenames and exact open error, and prevent `f.Open`, copy, and
   close operations.
8. Directory entries must continue without any file-open attempt. The focused
   shell double may broaden only from recording `MkdirAll` to recording the
   existing `OpenFile` operation; pre-existing complete doubles stay unchanged.
9. Nominal Q1.3 direction is 12/33 when the direct file-open caller leaves and
   the existing adapter implementation remains singular. Q0.6, Q1.1, Q1.2,
   Q1.4, and exact Q2.1 should hold; regenerate every value.
10. Direct `zip.OpenReader`, `f.Open`, `io.Copy`, `outFile.Close`, and
    `rc.Close` need later, separately authorized boundaries. Do not combine
    them with file-open selection.
11. Clock and server adapter creation belongs to queued P4. The remaining
    process sites also require their own coherent P3 moves.
12. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams,
    make mutation harnesses executable, or claim P4/P5 coverage.
13. Full preflight uses APIDIFF and golangci-lint v2.12.2 from
    `/private/tmp/ply-open-browser-tools/bin`; keep reports and caches under an
    external move-specific root.

## Next Objective

Move only the direct `os.OpenFile` call in the file-entry branch of
`pkg/shell.Unzip` through the existing `filesystem.OpenFile` helper.

Start red by extending only the focused shell recording double and contracts
to prove the exact joined file path, exact
`os.O_WRONLY|os.O_CREATE|os.O_TRUNC` flags, exact zip entry mode, one parent
`MkdirAll` followed by one file-open attempt, exact open error identity, current
partial filenames, and no later adapter or legacy file-entry operation on open
failure. Assert that directory entries retain their continue short-circuit and
make no file-open attempt, and that an empty recorded file-open population is
rejected.

Reuse the existing private complete unzip composition and production
`filesystem.System()` selection. Replace only direct `os.OpenFile` with
`filesystem.OpenFile(dependencies.Files, ...)`. Add no filesystem interface
operation and change no adapter implementation. Preserve the public signature,
every caller, archive open and deferred close, traversal order, path join,
logging, traversal validation, filename append, both directory-creation
adapter attempts, directory continue, parent selection, exact flags and mode,
direct `f.Open`, copy, file close, entry close, final results, and every error
precedence.

Expected nominal direction is Q1.3 12/33 with Q0.6 at 25 guarded sites, Q1.1
6/25, Q1.2 zero, Q1.4 7/8, and exact Q2.1 0/8. Require zero comparable ratchet
regressions and regenerate exact values.

## Verification Notes

Completed from implementation commit b9fe2091e134b465d6c959881375cc53bdeec5c5:

- Valid red evidence: isolated focused compilation failed only on absent
  `unzipDependencies`, `systemUnzipDependencies`, and
  `unzipWithDependencies`; an earlier host-cache permission failure was
  rejected as red evidence.
- Focused shell/filesystem/Spring and relevant process, command, context,
  config, HTTP, Maven, structurizr, profile, browser, tips, file, template,
  Bitbucket, Kibana, local-config, and caller package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls after one
  documented signal-fixture timing flake and a complete rerun.
- Make preflight meta-contracts, pinned linter with zero issues, one complete
  preflight, complete `make test`, install, and all 15 audit meta-controls:
  PASS with external tools, reports, and caches. The first preflight and first
  two `make test` attempts hit the documented nested signal/log timing flake;
  complete reruns passed.
- API, CLI, and subprocess compatibility: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their meta-contracts:
  PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 13/34.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 62
  Markdown files and zero ratchet regressions.
- Implementation commit: b9fe2091e134b465d6c959881375cc53bdeec5c5
  (refactor: route shell unzip mkdir through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete shell implementation/tests/callers, complete
filesystem adapter/tests and every complete double, Spring archive/download
contracts, relevant file, HTTP, process, config, Maven, structurizr, context,
Bitbucket, Wpost, local-config, browser, profile, and audit contracts before
editing. Confirm branch, HEAD, status, reciprocal links, and
`./codex-dev-start.sh --check`. Begin red only for the direct unzip file-open
selection, then finish with one implementation commit and one separate
handoff-only commit.

Stop before changing archive open, directory-creation behavior, path traversal
validation, entry ordering, flags or entry mode, entry open, copy or close
behavior, process construction, another filesystem operation or caller,
Spring, clock, server, public API, inventory, Q1.4, P4, mutation, Docker,
distribution, or publication. Stop on API/CLI change, comparable ratchet
regression, authoritative full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
