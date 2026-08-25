# Quality Upgrade Handover

Generated: 2026-08-25T19:34:22+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 2f0a07287de1.
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
and 2f0a072 in roadmap order. The separate operational continuity
implementation is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT while P3 is active and P4-P8 are queued in
the machine-readable plan block. Its active archive will be
`docs/plan/agent-sessions/2026-08-25T193422+0200-migrate-shell-unzip-output-close.md`.
The archive-open predecessor is answered history, and the reciprocal archive
graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 51 Preserved

Public `pkg/shell.Unzip(string, string) ([]string, error)` still selects the
complete private `unzipDependencies` composition with `filesystem.System()`.
Only the initial archive-open expression changed: direct `zip.OpenReader(src)`
is now `filesystem.OpenZipReader(dependencies.Files, src)`.

The existing filesystem adapter gained exactly one
`OpenZipReader(string) (*zip.ReadCloser, error)` operation, its zero-safe
forwarder, and the exact `zip.OpenReader(path)` system implementation. Every
complete filesystem double received the mechanical method addition. The
dependency returns its exact reader and error without wrapping, fallback,
retry, buffering, logging, cleanup, or path normalization.

The assignment, early archive-open return, deferred archive Close expression
and placement, returned archive identity, traversal and entry order, joined
paths, zip-slip behavior, filename append order, directory short-circuit,
parent creation, output OpenFile path/flags/mode, entry open, Copy request and
ignored results, output-file then entry-reader Close attempts, later traversal,
partial filenames, and all established error and close precedence remain
unchanged. No other filesystem operation, caller, public API, or inventory
entry changed.

Seven focused contracts bring the suite to 346 tests. Adapter contracts prove
zero-safe behavior, exact path/reader/error forwarding, non-empty recording,
and exact system entry order, names, and bytes. Private Unzip contracts use an
injected `*zip.ReadCloser` with a deliberately nonexistent supplied source path
to prove one exact request, injected-reader identity through traversal, ordered
filenames, the existing deferred archive Close, exact archive-open error with
nil partial filenames, complete dependency selection, zero-safe failure, and
no unrelated filesystem request. Guarded archive fixtures remain below
`t.TempDir()`.

## Measured Quality State

The clean full audit from implementation commit `2f0a072` reports:

- Absolute L0: 8 of 8.
- 346 test functions, zero skipped; 19 of 25 packages have tests.
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
- Q3.4: 0 state-claim phrases across 65 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
`/private/tmp/ply-quality-move51-clean-final/scorecard.json`; the final focused
implementation report is
`/private/tmp/ply-quality-move51-focused-final/scorecard.json`. Compatibility
reports, tool binaries, Go and linter caches, audit output, raw logs, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. The filesystem helper must return the exact concrete `*zip.ReadCloser` from
   the dependency. Rebuilding or wrapping it would change archive identity,
   traversal, or close behavior.
2. A deliberately nonexistent supplied source path plus a valid injected
   reader proves both that the source is forwarded unchanged and that private
   traversal uses the injected reader instead of reopening the path.
3. The existing deferred `func() { _ = r.Close() }()` remains immediately after
   the archive-open error check. A second Close returning `os.ErrClosed` proves
   it still runs without adding another defer or cleanup path.
4. Q1.3 improves exactly from 11/32 to 10/32. The direct public-package
   `zip.OpenReader` site left the violation set, while the exact system call is
   inside the declared filesystem adapter.
5. The focused system archive-open fixture adds one guarded write call site,
   moving Q0.6 from 25 to 26 safe-writer sites, 21 write and 5 copy, while
   unsafe direct test writes stay zero.
6. The first standalone launcher run hit the timing-sensitive synthetic signal
   log check; its immediate complete rerun passed all 62 controls. The full
   preflight rerun and `make test` launcher runs also passed all 62.
7. The first preflight attempt reached lint after compatibility, tests, and vet
   but the sandbox rejected the linter's default user cache. Setting
   `GOLANGCI_LINT_CACHE` under `/private/tmp` produced a complete green rerun
   with zero lint issues.
8. The next coherent Unzip effect is the direct output-file
   `err = outFile.Close()` at `pkg/shell/command.go:110`. Keep the separate
   archive defer and entry-reader `rc.Close()` unchanged.
9. The scanner retained filesystem provenance through the earlier Copy helper's
   opened-file and reader arguments. A Close helper receiving the opened file
   may likewise remain a caller violation; regenerate the exact result and do
   not reshape the scanner or broaden the move.
10. Clock and server adapter creation belongs to queued P4. Plugin diagrams,
    Maven, entry-reader Close, and remaining process effects require separate
    moves. `.quality/inventory` remains checksum-bound and unchanged.

## Next Objective

Move only the direct output-file `err = outFile.Close()` effect in private
`pkg/shell.Unzip` through one narrow Close operation on the existing filesystem
adapter. Keep public `Unzip`, private `unzipDependencies`, and production
`filesystem.System()` selection complete.

Start red with focused filesystem-adapter and private Unzip output-close
recording contracts. Prove the exact opened destination identity, one close
request for every reached file after its one Copy request, exact injected close
error and partial filenames, Copy-before-close order, suppression of
entry-reader Close and later traversal after an output-close error, continued
entry close and later traversal after success, zero-safe adapter behavior,
rejection of an empty recording population, complete dependency preservation,
and absence of unrelated filesystem requests. Launch no process, touch no
network, and change no working directory.

Add only `Close(File) error` to the existing filesystem interface, its
zero-safe forwarding helper and exact `file.Close()` system implementation,
plus the mechanical method addition required by every complete filesystem
double. Replace only `err = outFile.Close()` with
`err = filesystem.Close(dependencies.Files, outFile)`. Keep its assignment,
following error branch, exact error identity, and every surrounding operation
unchanged.

Do not change archive open or deferred Close, MkdirAll, OpenFile, entry open,
Copy, `rc.Close()`, another unzip branch, shell Run, another adapter operation
or caller, public API, inventory, audit, mutation harnesses, or later roadmap
work. Expect Q0.6 to hold at 26 guarded sites, Q1.1 at 6/25, Q1.2 at zero,
Q1.4 at 7/8, and Q2.1 at 0/8. Because opened-file provenance may remain at the
helper call while the new in-boundary implementation adds one production site,
regenerate Q1.3 rather than forcing a nominal value; require zero comparable
ratchet regressions.

## Verification Notes

Completed from implementation commit `2f0a07287de1a148240c532f18198242039bba29`:

- Valid red evidence: the adapter suite failed to compile because
  `OpenZipReader` did not exist, and the two private shell tests returned direct
  path-open errors instead of using the injected reader/error. No process or
  network request ran.
- Focused shell/filesystem/Spring and relevant process, Maven, command,
  context, config, HTTP, structurizr, profile, browser, tips, file, template,
  Bitbucket, Wpost, local-config, Kibana, and caller tests: PASS.
- Standalone `/bin/bash` launcher contract complete rerun: PASS all 62 controls.
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
- Handoff launcher contract and `./codex-dev-start.sh --check`: PASS with the
  normalized stable skeleton unchanged and all 62 controls green.
- Handoff-state Q3.4 audit: exit 0 with zero state-claim phrases across 66
  Markdown files and zero ratchet regressions.
- Implementation commit: `2f0a07287de1a148240c532f18198242039bba29`
  (`refactor: route unzip archive open through filesystem adapter`).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete shell/filesystem-adapter source and tests, every
complete filesystem double, completed Unzip/Spring/process contracts, the
import-aware effect scanner, and relevant callers before editing. Confirm
branch, HEAD, status, reciprocal links, and `./codex-dev-start.sh --check`.
Begin red only for output-file Close, then finish with one implementation
commit and one separate handoff-only commit.

Stop before changing archive open or Close, Copy, entry-reader Close, entry
open, output OpenFile, directory creation, another unzip function, Run, Git,
structurizr, plugin diagrams, Maven, HTTP, Spring, clock, server, public API,
inventory, Q1.4, P4, mutation, Docker, distribution, or publication. Stop on
API/CLI change, comparable ratchet regression, authoritative full-audit exit 2,
or failure to isolate implementation from generated and handoff-only state.
