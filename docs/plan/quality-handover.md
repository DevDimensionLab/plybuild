# Quality Upgrade Handover

Generated: 2026-08-25T18:18:27+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: da7eebf9256c.
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
85f4c2b, 0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, and da7eebf in
roadmap order. The separate operational continuity implementation is 1b85711
and changes no Go quality denominator.

## Continuity Checkpoint

codex-dev-start.sh remains NEXT while P3 is active and P4-P8 are queued in the
machine-readable plan block. Its active archive is
docs/plan/agent-sessions/2026-08-25T181827+0200-migrate-shell-unzip-copy.md.
The shell Run predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484.
test/codex_dev_start_test.sh retains all 62 controls.

## P3 Move 49 Preserved

Public `pkg/shell.Run(string, ...string) Output` now selects one private
complete `runDependencies` composition containing `process.Dependencies`.
Only public Run selects `process.System()`. The private body constructs the
complete `process.Command` from the exact requested name and ordered variadic
arguments, logs the command description before stream wiring, assigns distinct
stdout and stderr buffers, and calls `process.Execute` once with empty Dir and
Stdin and Start false.

The public signature, `Output`, `String`, `FormatError`, stdout/stderr bytes,
synchronous completion, Maven analysis and both Git callers, and result
interpretation remain unchanged. The legacy dependency-error behavior is
intentionally preserved: after bytes delivered before an execution error,
Run returns early without assigning the error, so `Output.Err` remains nil.
The uncalled private `run(*exec.Cmd)` helper was removed. The process adapter,
its interface and implementation, Git-specific composition, Maven-specific
composition, structurizr, plugin diagrams, Unzip, and inventory did not change.

Four focused private contracts bring the suite to 336 tests. The shell Run
double records only the complete process request, copies ordered arguments,
and writes arbitrary in-memory stdout/stderr bytes. The contracts prove
complete production selection, exact name and argument bytes, one synchronous
request, empty Dir/Stdin, Start false, distinct writer identity and delivered
bytes, exact injected error observation through the legacy nil `Output.Err`,
safe-zero behavior, a non-empty recording population, and no unrelated process
request. No focused test launches a process, touches the network, changes the
working directory, or writes a file.

## Measured Quality State

The clean full audit from implementation commit da7eebf reports:

- Absolute L0: 8 of 8.
- 336 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 25 guarded safe-writer sites, 20 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside main.
- Q1.3: 11 direct external sites outside five declared adapters of 32
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 state-claim phrases across 63 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
/private/tmp/ply-shell-run-tools/reports/clean-full-da7eebf/scorecard.json and
the focused implementation report is
/private/tmp/ply-shell-run-tools/reports/focused-dirty/scorecard.json.
Compatibility reports, tool binaries, Go and linter caches, audit output, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. The existing Run helper discarded every `cmd.Run()` error rather than
   assigning it to `Output.Err`; the adapter move preserves that surprising
   public result exactly while retaining bytes written before failure.
2. One complete process request replaces the direct caller construction and
   execution. Q1.3 moved exactly from 12/33 to 11/32 without adding an adapter
   implementation site.
3. The command description remains name followed by ordered arguments, the log
   remains before buffer wiring, and execution remains after both writers are
   assigned.
4. Public Run callers remain Maven analysis, `GitDirty`, and `GitIsRepo`. Their
   requested arguments and result handling did not change.
5. The next narrow existing-adapter move is the direct
   `io.Copy(outFile, rc)` inside `pkg/shell.Unzip`. The existing complete unzip
   composition already carries `filesystem.Dependencies`, and the filesystem
   adapter already exposes `Copy(filesystem.File, io.Reader)`.
6. Legacy Unzip ignores both values returned by `io.Copy`, then closes the
   output file, closes the archive entry reader, and continues traversal.
   Preserve that ignored copy error and every later close error and precedence.
7. The focused unzip double can record exact destination identity and consume
   source bytes in memory without writing the destination. A pipe-backed
   `*os.File` can satisfy the unchanged OpenFile result without a repository
   write; any archive fixture remains guarded below `t.TempDir()`.
8. Directory entries, MkdirAll failure, OpenFile failure, and `f.Open` failure
   must produce no Copy request. The empty Copy recording population needs an
   explicit rejection contract.
9. Nominal Q1.3 direction is 10/31 when the direct unzip Copy caller site
   leaves and the existing system-adapter site remains singular. Q0.6, Q1.1,
   Q1.2, Q1.4, and exact Q2.1 should hold; regenerate every value.
10. Direct `zip.OpenReader`, `outFile.Close`, and `rc.Close` remain separate
    filesystem effects requiring later coherent authorization.
11. Clock and server adapter creation belongs to queued P4. Plugin-diagram,
    structurizr, Maven stdout, and remaining process sites need separate moves.
12. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams,
    make mutation harnesses executable, or claim P4/P5 coverage.
13. Full preflight uses APIDIFF and golangci-lint v2.12.2 from
    `/private/tmp/ply-open-browser-tools/bin`; keep reports and caches under an
    external move-specific root.

## Next Objective

Move only the file-entry payload transfer inside public `pkg/shell.Unzip` from
direct `io.Copy(outFile, rc)` to the existing
`filesystem.Copy(dependencies.Files, outFile, rc)` forwarding helper.

Start red by extending only the private unzip recording filesystem and focused
contracts. Prove exact opened destination identity, delivered arbitrary entry
bytes and order, one Copy request per reached file, the established
MkdirAll/OpenFile/entry-open/copy order, ignored copy count/error behavior,
continued close and later traversal, no Copy on directories or prior failures,
rejection of an empty Copy population, and no unrelated operation.

Replace only the direct copy call and ignore both adapter return values exactly
as before. Preserve public API, complete dependency selection, archive and
entry traversal, partial filenames, zip-slip behavior, close order and error
precedence. Change no adapter, archive open, direct close, other unzip branch,
Run, Git, process, HTTP, Spring, inventory, or other caller.

Expected nominal direction is Q1.3 10/31 with Q0.6 at 25 guarded sites, Q1.1
6/25, Q1.2 zero, Q1.4 7/8, and exact Q2.1 0/8. Require zero comparable ratchet
regressions and regenerate exact values.

## Verification Notes

Completed from implementation commit da7eebf9256ccb5d5b868290db027a605c34b47e:

- Valid red evidence: the focused shell Run contracts failed to compile only
  because `runDependencies`, `systemRunDependencies`, and
  `runWithDependencies` did not exist. No real process was invoked.
- Focused shell/process/Maven and relevant command, context, config,
  filesystem, HTTP, Spring, structurizr, profile, browser, tips, file,
  template, Bitbucket, Kibana, local-config, and caller package tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls on a complete
  rerun. One earlier attempt hit the documented nested signal/log timing flake.
- Make preflight meta-contracts, pinned linter with zero issues, complete
  preflight, complete `make test`, install, and all 15 audit meta-controls:
  PASS with external tools, reports, and caches. The first complete `make test`
  attempt hit the same signal/log timing flake; its complete rerun passed.
- API, CLI, and subprocess compatibility: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their meta-contracts:
  PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 11/32.
- Clean full audit with the pinned baseline: exit 1, 15 findings, L0 8/8, five
  improved, two held, zero regressed, one non-comparable, and zero dirty paths.
- Handoff launcher contract and `./codex-dev-start.sh --check`: PASS with the
  stable skeleton unchanged and all 62 controls green.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 64
  Markdown files and zero ratchet regressions.
- Implementation commit: da7eebf9256ccb5d5b868290db027a605c34b47e
  (refactor: route shell run through process adapter).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete shell and filesystem-adapter source/tests and
callers, all complete filesystem doubles, unzip and Spring contracts, and
relevant caller packages before editing. Confirm branch, HEAD, status,
reciprocal links, and `./codex-dev-start.sh --check`. Begin red only for the
unzip Copy request, then finish with one implementation commit and one separate
handoff-only commit.

Stop before changing archive open, either direct Close, another unzip function,
Run, Git-specific composition, structurizr, plugin diagrams, Maven process
composition, filesystem or HTTP adapter behavior, Spring, clock, server, public
API, inventory, Q1.4, P4, mutation, Docker, distribution, or publication. Stop
on API/CLI change, comparable ratchet regression, authoritative full-audit exit
2, or failure to isolate implementation from generated and handoff-only state.
