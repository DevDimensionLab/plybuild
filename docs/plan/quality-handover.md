# Quality Upgrade Handover

Generated: 2026-08-25T17:48:00+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 2a684a09e33e.
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
85f4c2b, 0e10282, 9fcdfb5, d982f63, b9fe209, and 2a684a0 in roadmap order.
The separate operational continuity implementation is 1b85711 and changes no
Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh stays NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T174800+0200-migrate-shell-run-process.md.
The unzip file-open predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 48 Preserved

Public `pkg/shell.Unzip(string, string) ([]string, error)` retains the private
complete composition containing `filesystem.Dependencies`, with
`filesystem.System()` selected only by the production entry. The established
body remains in its original order. Only the former direct file-entry
`os.OpenFile` call uses
`filesystem.OpenFile(dependencies.Files, fpath,
os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())`.

Directory entries still append their exact joined path, attempt one
`MkdirAll`, and continue without an OpenFile attempt. File entries still
append their joined path, complete one parent `MkdirAll` on
`filepath.Dir(fpath)`, and then attempt OpenFile once with the exact flags and
entry mode. An OpenFile error returns the exact current partial filenames and
error before `f.Open`, copy, or either close. Archive open and deferred close,
traversal, debug logging, zip-slip validation and text, direct entry open,
copy, file close, entry close, final results, and every other error precedence
remain unchanged. No filesystem interface or adapter implementation changed.

Two focused top-level contracts bring the suite to 332 tests. The focused
shell double records the already-established MkdirAll operation and existing
OpenFile operation. The contracts prove exact path, flags, mode, operation and
entry order, partial results, exact open-error identity and precedence,
directory no-open short-circuit, a non-empty OpenFile population, and no
unrelated adapter operation. The guarded zip fixture stays under `t.TempDir()`.

## Measured Quality State

The clean full audit from implementation commit 2a684a0 reports:

- Absolute L0: 8 of 8.
- 332 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 25 guarded safe-writer sites, 20 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 12 direct external sites outside five declared adapters of 33
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 temporary-state claims across 62 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
/private/tmp/ply-shell-unzip-open-tools/reports/clean-full-2a684a0-correct/scorecard.json
and the focused implementation report is
/private/tmp/ply-shell-unzip-open-tools/reports/focused-audit/scorecard.json.
Compatibility reports, tool binaries, Go and linter caches, audit output, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. Reusing the existing OpenFile operation removed the single unzip caller
   site without adding a system-adapter site; Q1.3 moved exactly to 12/33.
2. The focused double records OpenFile only after successful parent MkdirAll.
   Directory entries and parent-creation failures record no open attempt.
3. An unsupported zip method after the injected open boundary proves the exact
   open error wins before the unchanged `f.Open` error, copy, or close work.
4. The fixture uses `zip.Writer.CreateRaw` only for that in-memory unsupported
   entry and makes one guarded archive write below `t.TempDir()`.
5. Direct `zip.OpenReader`, `io.Copy`, `outFile.Close`, and `rc.Close` remain
   separate filesystem effects requiring later coherent authorization.
6. The narrow next existing-adapter move is public `pkg/shell.Run`: its direct
   `exec.Command` construction is one Q1.3 process site, and the existing
   process adapter already accepts name, ordered args, streams, execution mode,
   and returns the dependency error.
7. Characterize legacy Run error behavior before implementation. The existing
   helper returns without assigning the `cmd.Run()` error into `Output.Err`;
   preservation is required even if it appears surprising.
8. Public Run callers are Maven analysis, `GitDirty`, and `GitIsRepo`. Keep
   their signatures, argument order, output capture, and result interpretation
   unchanged.
9. The next production composition selects `process.System()` only at public
   Run and passes a complete private dependency struct to the private body.
   Change no process interface or adapter implementation.
10. Nominal Q1.3 direction is 11/32 when the one direct Run construction site
    leaves and the existing adapter implementation remains singular. Q0.6,
    Q1.1, Q1.2, Q1.4, and exact Q2.1 should hold; regenerate every value.
11. Clock and server adapter creation belongs to queued P4. Plugin-diagram,
    structurizr, Maven stdout, and remaining process sites need separate moves.
12. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams,
    make mutation harnesses executable, or claim P4/P5 coverage.
13. Full preflight uses APIDIFF and golangci-lint v2.12.2 from
    `/private/tmp/ply-open-browser-tools/bin`; keep reports and caches under an
    external move-specific root.

## Next Objective

Move only public `pkg/shell.Run` process construction and synchronous
execution through the existing `process.Execute` helper.

Start red with a focused private recording double and contracts proving exact
name and ordered variadic arguments, one complete synchronous request with
Start false and empty Dir/Stdin, exact stdout/stderr writer identity and
delivered bytes, the established dependency-error-to-Output behavior,
safe-zero behavior, rejection of an empty recorded population, and no real or
unrelated process request.

Add one private complete process dependency composition and select
`process.System()` only in public Run. Preserve public API, Output methods,
debug logging position, stream wiring, one execution attempt, caller behavior,
and legacy error handling. Change no adapter, Git-specific composition,
structurizr, plugin diagrams, Maven-specific process path, Unzip, filesystem,
HTTP, Spring, inventory, or other caller.

Expected nominal direction is Q1.3 11/32 with Q0.6 at 25 guarded sites, Q1.1
6/25, Q1.2 zero, Q1.4 7/8, and exact Q2.1 0/8. Require zero comparable ratchet
regressions and regenerate exact values.

## Verification Notes

Completed from implementation commit 2a684a09e33e8de2edaac9dd8f96bea4b22d1c18:

- Valid red evidence: the focused file-open contract reached the legacy branch
  and returned the host `os.OpenFile` PathError instead of the injected exact
  dependency error; the empty-population guard passed. A host-cache permission
  failure was rejected as red evidence.
- Focused shell/filesystem/Spring and relevant process, command, context,
  config, HTTP, Maven, structurizr, profile, browser, tips, file, template,
  Bitbucket, Kibana, local-config, and caller package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls.
- Make preflight meta-contracts, pinned linter with zero issues, complete
  preflight, complete `make test`, install, and all 15 audit meta-controls:
  PASS with external tools, reports, and caches. The first complete preflight
  attempt hit the documented nested signal/log timing flake; the complete
  rerun passed all 62 controls.
- API, CLI, and subprocess compatibility: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their meta-contracts:
  PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 12/33.
- Clean full audit with the pinned baseline: exit 1, 15 findings, L0 8/8, five
  improved, two held, zero regressed, one non-comparable, and zero dirty paths.
- One discarded clean-audit invocation used a nonexistent baseline filename;
  it failed closed before ratchet comparison. The complete rerun used
  `.quality/baseline/scorecard.json` and produced the authoritative report.
- Handoff launcher contract: PASS all 62 controls after the first handoff run
  hit the documented signal/log timing flake and a complete rerun passed.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 63
  Markdown files and zero ratchet regressions.
- Implementation commit: 2a684a09e33e8de2edaac9dd8f96bea4b22d1c18
  (refactor: route unzip open through filesystem adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete shell and process-adapter source/tests/callers,
all complete process doubles, Maven analyze, and relevant completed contracts
before editing. Confirm branch, HEAD, status, reciprocal links, and
`./codex-dev-start.sh --check`. Begin red only for public shell Run process
selection, then finish with one implementation commit and one separate
handoff-only commit.

Stop before changing another shell function, Git-specific composition,
structurizr, plugin diagrams, Maven process composition, any unzip effect,
filesystem or HTTP behavior, Spring, clock, server, public API, inventory,
Q1.4, P4, mutation, Docker, distribution, or publication. Stop on API/CLI
change, comparable ratchet regression, authoritative full-audit exit 2, or
failure to isolate implementation from generated and handoff-only state.
