# Quality Upgrade Handover

Generated: 2026-08-25T20:09:05+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 0128cd0809c5.
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
85f4c2b, 0e10282, 9fcdfb5, d982f63, b9fe209, 2a684a0, da7eebf, 89918bd,
2f0a072, and 0128cd0 in roadmap order. The separate operational continuity
implementation is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT while P3 is active and P4-P8 are queued in
the machine-readable plan block. Its active archive will be
`docs/plan/agent-sessions/2026-08-25T200905+0200-migrate-shell-unzip-entry-close.md`.
The output-close predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 52 Preserved

Public `pkg/shell.Unzip(string, string) ([]string, error)` still selects the
complete private `unzipDependencies` composition with `filesystem.System()`.
Only the output-file close expression changed: direct
`err = outFile.Close()` is now
`err = filesystem.Close(dependencies.Files, outFile)`.

The existing filesystem adapter gained exactly one `Close(File) error`
operation, its zero-safe forwarder, and the exact `file.Close()` system
implementation. Every complete filesystem double received the mechanical
method addition. The dependency receives the exact opened output identity and
returns its exact error without fallback, retry, logging, cleanup, wrapping,
or another close attempt.

Archive open and deferred Close, traversal and entry order, joined paths,
zip-slip behavior, filename append order, directory short-circuit, parent
creation, output OpenFile path/flags/mode, entry open, Copy request and ignored
results, output-close assignment and following error branch, entry-reader
Close, later traversal, partial filenames, and all established error and close
precedence remain unchanged. No other filesystem operation, caller, public API,
or inventory entry changed.

Six focused contracts bring the suite to 352 tests. Adapter contracts prove
zero-safe behavior without closing the supplied file, exact file identity and
one dependency request, exact dependency and system errors, one exact system
close, and non-empty recording. Private Unzip contracts prove one exact opened
output identity per reached file, Copy-before-close order, exact close error
and partial filenames, output-close-error suppression of entry Close and later
entries, successful entry Close and later traversal, complete dependency
selection, non-empty recording, and no unrelated filesystem request. The
guarded archive fixture remains below `t.TempDir()` and injected pipes do not
target the repository.

## Measured Quality State

The clean full audit from implementation commit `0128cd0` reports:

- Absolute L0: 8 of 8.
- 352 test functions, zero skipped; 19 of 25 packages have tests.
- Q0.6: 26 guarded safe-writer sites, 21 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 6 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 10 direct external sites outside five declared adapters of 32
  production effect sites; missing clock and server make it non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.2: PASS because a real template contract references the README path.
- Q3.4: 0 state-claim phrases across 66 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
`/private/tmp/ply-quality-move52-clean-audit-final/scorecard.json`; the focused
implementation report is
`/private/tmp/ply-quality-move52-focused-audit/scorecard.json`. Compatibility
reports, pinned tool binaries, Go and linter caches, audit output, raw logs, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. The Close helper must receive the exact `filesystem.File` returned by
   OpenFile. Reopening, wrapping, fallback, retry, or direct helper-side Close
   would change identity or attempt count.
2. A zero filesystem dependency returns exact `filesystem.ErrNoFilesystem`
   without invoking a developer-supplied file's Close method.
3. A dependency-injected output-close error is returned after Copy with the
   reached filename retained. It prevents the entry-reader Close and every
   later entry; success still closes the entry reader and continues traversal.
4. Q1.3 holds exactly at 10/32. The direct output-file Close identity left, but
   the import-aware scanner retains filesystem provenance through the helper's
   opened-file argument and reports the replacement at the same caller line.
5. Custom ZIP method 93 plus a recording decompressor observes entry-reader
   Close without changing entry bytes. Raw entries require CRC32 and compressed
   and uncompressed sizes in the guarded fixture header.
6. Valid red evidence had the adapter tests fail to compile because Close was
   absent, while the private Unzip contracts rejected an empty output-close
   population and did not return the injected close error.
7. The pinned `apidiff` and `golangci-lint` binaries, all generated reports,
   and every build cache stayed under `/private/tmp`.
8. The next isolated Unzip effect is direct `err = rc.Close()` at
   `pkg/shell/command.go:115`. Keep the output-file Close and separate archive
   defer unchanged.
9. The present scanner does not report the interface-typed entry-reader Close.
   A CloseReader helper or in-boundary implementation may change or hold the
   production population; regenerate the exact result without reshaping it.
10. Clock and server adapter creation belongs to queued P4. Plugin diagrams,
    Maven provenance, the deferred archive Close, and remaining process effects
    require separate moves. `.quality/inventory` remains checksum-bound and
    unchanged.

## Next Objective

Move only the direct entry-reader `err = rc.Close()` effect in private
`pkg/shell.Unzip` through one narrow CloseReader operation on the existing
filesystem adapter. Keep public `Unzip`, private `unzipDependencies`, and
production `filesystem.System()` selection complete.

Start red with focused filesystem-adapter and private Unzip entry-close
recording contracts. Prove exact opened reader identity, one close request for
every reached file after Copy and successful output-file Close, exact injected
close error and partial filenames, Copy/output-Close/entry-Close order,
suppression of later traversal after an entry-close error, suppression of entry
Close after an output-close error, continued traversal after success, zero-safe
adapter behavior, rejection of an empty recording population, complete
dependency preservation, and absence of unrelated filesystem requests. Launch
no process, touch no network, and change no working directory.

Add only `CloseReader(io.ReadCloser) error` to the existing filesystem
interface, its zero-safe forwarding helper and exact `reader.Close()` system
implementation, plus the mechanical method addition required by every complete
filesystem double. Replace only `err = rc.Close()` with
`err = filesystem.CloseReader(dependencies.Files, rc)`. Keep its assignment,
following error branch, exact error identity, and every surrounding operation
unchanged.

Do not change archive open or deferred Close, output-file Close, Copy,
OpenFile, MkdirAll, entry open, another unzip branch, shell Run, another adapter
operation or caller, public API, inventory, audit, mutation harnesses, or later
roadmap work. Expect Q0.6 to hold at 26 guarded sites, Q1.1 at 6/25, Q1.2 at
zero, Q1.4 at 7/8, and Q2.1 at 0/8. Regenerate Q1.3 because the direct reader
Close is absent from its present violation set; require zero comparable
ratchet regressions.

## Verification Notes

Completed from implementation commit `0128cd0809c54043d35db906278f59edce1b85db`:

- Valid red evidence: adapter tests failed to compile because `Close` did not
  exist; private shell tests observed no adapter close request and no injected
  error. No process or network request ran.
- Focused shell/filesystem/Spring and relevant process, Maven, command,
  context, config, HTTP, structurizr, profile, browser, tips, file, template,
  Bitbucket, Wpost, local-config, Kibana, and caller tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls.
- Make meta-contracts, pinned linter with zero issues, complete preflight,
  complete `make test`, install, and all 15 audit meta-controls: PASS with
  external reports and caches.
- API, CLI, and subprocess compatibility: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their meta-contracts:
  PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 10/32.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- Handoff launcher contract: the first outer run and the second run's nested
  source-archive check hit the documented signal/partial-log timing flake; the
  third complete run and the final post-digest run passed all 62 controls, and
  `./codex-dev-start.sh --check` passed for the linked NEXT archive.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 67
  Markdown files and zero ratchet regressions.
- Implementation commit: `0128cd0809c54043d35db906278f59edce1b85db`
  (`refactor: route unzip output close through filesystem adapter`).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete shell/filesystem-adapter source and tests, every
complete filesystem double, completed Unzip/Spring/process contracts, the
import-aware effect scanner, and relevant callers before editing. Confirm
branch, HEAD, status, reciprocal links, and `./codex-dev-start.sh --check`.
Begin red only for entry-reader Close, then finish with one implementation
commit and one separate handoff-only commit.

Stop before changing archive open or deferred Close, output-file Close, Copy,
entry open, output OpenFile, directory creation, another unzip function, Run,
Git, structurizr, plugin diagrams, Maven, HTTP, Spring, clock, server, public
API, inventory, Q1.4, P4, mutation, Docker, distribution, or publication. Stop
on API/CLI change, comparable ratchet regression, authoritative full-audit exit
2, or failure to isolate implementation from generated and handoff-only state.
