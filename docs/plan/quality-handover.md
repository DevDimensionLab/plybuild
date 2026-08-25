# Quality Upgrade Handover

Generated: 2026-08-25T18:57:49+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 89918bd3e4ae.
- After launch, obtain the session head with `git rev-parse --short=12 HEAD`;
  the restart commit contains this handover and no product implementation
  changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits through this checkpoint are 5c6f2fa, 03d6242,
f5ee37d, ee5e9ab, ffc4e77, 204e222, 89d0f76, dee214c, 789ae23, 07ac6ce,
a7eb3ef, e13a036, f59a3f0, 61714a5, 60e5aac, e054082, 27c0d1a, c2f3597,
9e2d669, a4deb76, acda4e3, 838daa1, 224a691, 549685d, 04cfe44, d03e96d,
d6cb593, 82e631d, cb94f81, dfcfa75, 170b0ae, 8d49345, 7754575, 2566438,
527a8b9, 886dff0, 89f43aa, fd45ebf, d2e330b, b232dda, 6928a72, b2d37cc,
85f4c2b, 0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, and
89918bd in roadmap order. The separate operational continuity implementation
is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT while P3 is active and P4-P8 are queued in
the machine-readable plan block. Its active archive is
`docs/plan/agent-sessions/2026-08-25T185749+0200-migrate-shell-unzip-archive-open.md`.
The unzip Copy predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 50 Preserved

Public `pkg/shell.Unzip(string, string) ([]string, error)` still selects the
complete private `unzipDependencies` composition with `filesystem.System()`.
Only the file-entry payload expression changed: direct
`io.Copy(outFile, rc)` is now
`filesystem.Copy(dependencies.Files, outFile, rc)`. Both returned values remain
ignored exactly as before.

The joined paths, zip-slip check, filename append and archive order, directory
short-circuit, parent creation, OpenFile path, flags and mode, entry open, one
payload attempt per reached file, output-file Close then entry-reader Close,
later traversal, partial results, error identity, and every earlier error and
close error precedence remain unchanged. Archive open and close, both direct
entry closes, the filesystem interface and implementation, callers, public
API, and inventory did not change.

Three focused contracts bring the suite to 339 tests. The existing private
unzip filesystem double now records MkdirAll, OpenFile, and Copy only. It keeps
the exact opened destination identity, consumes exact ordered arbitrary source
bytes in memory, and returns arbitrary count/error values without writing the
destination. The contracts prove one Copy request per reached file entry,
ignored count/error behavior, parent/open/copy order, continued close and later
traversal, no Copy for directories or every earlier failure, non-empty Copy
recordings, complete dependency preservation, and no unrelated filesystem
request. Pipe-backed output files prove close without a repository write, and
the guarded zip fixture stays below `t.TempDir()`.

## Measured Quality State

The clean full audit from implementation commit `89918bd` reports:

- Absolute L0: 8 of 8.
- 339 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 25 guarded safe-writer sites, 20 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 11 direct external sites outside five declared adapters of 32
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 state-claim phrases across 64 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
`/private/tmp/ply-unzip-copy-tools/reports/clean-full-89918bd/scorecard.json`
and the final focused implementation report is
`/private/tmp/ply-unzip-copy-tools/reports/focused-dirty-commit/scorecard.json`.
Compatibility reports were removed before the clean audit; tool binaries, Go
and linter caches, audit output, raw logs, and empty-HOME state stayed outside
the measured tree.

## Decisions And Learned Facts

1. Legacy Unzip ignored both values from `io.Copy`; the adapter request must
   keep ignoring both values. An injected Copy error does not stop either close
   or later archive traversal.
2. The recording Copy operation must consume the archive reader to preserve the
   exact byte-delivery observation, but it must not write the destination.
   Pipe-backed `*os.File` values preserve the unchanged OpenFile type and exact
   identity.
3. A closed output pipe can be proved with `Stat` returning `os.ErrClosed`
   before reading the empty pipe. Writing to the pipe is not suitable because
   the Q0.6 scanner correctly treats direct `os.File.Write` in tests as an
   unsafe write site.
4. The direct `io.Copy` identity is gone, but measured Q1.3 holds at 11/32
   instead of the nominal 10/31. The import-aware scanner follows the opened
   destination and archive-reader arguments through the local adapter request
   and continues to classify that caller line as filesystem-effect provenance.
   This is an exact instrument result, not an adapter or inventory change.
5. The current move did not change the audit, filesystem interface, adapter
   implementation, archive open, or either close to force the expected metric.
   Comparable ratchets have zero regressions.
6. The next coherent Unzip effect is direct `zip.OpenReader(src)`. It needs one
   narrow archive-open operation on the existing filesystem adapter, exact
   system forwarding, complete-double updates, and focused request/error
   contracts. Keep the returned concrete `*zip.ReadCloser` so traversal and
   deferred close behavior remain unchanged.
7. The archive-open move must not absorb `filesystem.Copy`, `outFile.Close`,
   `rc.Close`, or any entry operation. Those remain separately authorized work.
8. Clock and server adapter creation belongs to queued P4. Plugin diagrams,
   structurizr process construction, Maven stdout provenance, and remaining
   process sites need separate P3 moves.
9. `.quality/inventory` remains baseline-checksum-bound. Do not relabel seams,
   make mutation harnesses executable, or claim P4/P5 coverage.
10. Full preflight uses APIDIFF and golangci-lint v2.12.2 from
    `/private/tmp/ply-open-browser-tools/bin`; keep reports and caches under an
    external move-specific root.

## Next Objective

Move only the direct `zip.OpenReader(src)` archive-open effect at the start of
private `pkg/shell.Unzip` traversal through one narrow operation added to the
existing filesystem adapter. Keep public `Unzip` and private
`unzipDependencies` complete, and keep `filesystem.System()` selection only at
the public entry.

Start red with focused filesystem adapter and private unzip archive-open
recording contracts. Prove exact source-path forwarding, complete system
selection, one open request, exact returned reader identity, exact injected
open error and nil partial filenames, safe-zero behavior, rejection of an
empty recording population, and no unrelated filesystem request. Keep the
existing deferred archive Close expression and placement unchanged. Any real
fixture open stays below `t.TempDir()`;
launch no process, touch no network, and change no working directory.

Add only an `OpenZipReader(string) (*zip.ReadCloser, error)` operation and its
zero-safe forwarding helper/system implementation to the existing filesystem
adapter, with the mechanical signature addition required by every complete
filesystem double. Replace only `zip.OpenReader(src)` with the helper call.
Preserve traversal, entry identity and order, zip-slip behavior, filename
append and partial results, MkdirAll/OpenFile/Copy requests, ignored Copy
results, output and entry close order/errors, deferred archive Close, and all
earlier precedence. Change no other adapter operation or caller.

Expected nominal direction is Q1.3 10/32, with Q0.6 at 25 guarded sites,
Q1.1 6/25, Q1.2 zero, Q1.4 7/8, and exact Q2.1 0/8. Regenerate exact values and
require zero comparable ratchet regressions; do not reshape the scanner to
obtain an expected number.

## Verification Notes

Completed from implementation commit `89918bd3e4ae877bf01364fbee741e8a7b93d47f`:

- Valid red evidence: the focused shell suite failed only because the new Copy
  recording population was empty while production still called direct
  `io.Copy`; no process or network request ran.
- Focused shell/filesystem/Spring and relevant process, Maven, command,
  context, config, HTTP, structurizr, profile, browser, tips, file, template,
  Bitbucket, Wpost, local-config, Kibana, and caller tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls.
- Make preflight meta-contracts, pinned linter with zero issues, complete
  preflight, complete `make test`, install, and all 15 audit meta-controls:
  PASS with external tools, reports, logs, and caches.
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
  normalized stable skeleton unchanged and all 62 controls green.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 65
  Markdown files and zero ratchet regressions.
- Implementation commit: `89918bd3e4ae877bf01364fbee741e8a7b93d47f`
  (`refactor: route unzip copy through filesystem adapter`).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete shell and filesystem-adapter source/tests and
callers, every complete filesystem double, completed unzip/Spring/process
contracts, and relevant caller packages before editing. Confirm branch, HEAD,
status, reciprocal links, and `./codex-dev-start.sh --check`. Begin red only for
the archive-open request, then finish with one implementation commit and one
separate handoff-only commit.

Stop before changing Copy, either direct entry Close, output-file open,
directory creation, another unzip function, Run, Git-specific composition,
structurizr, plugin diagrams, Maven process composition, another filesystem
operation or caller, HTTP, Spring, clock, server, public API, inventory, Q1.4,
P4, mutation, Docker, distribution, or publication. Stop on API/CLI change,
comparable ratchet regression, authoritative full-audit exit 2, or failure to
isolate implementation from generated and handoff-only state.
