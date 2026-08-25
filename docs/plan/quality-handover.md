# Quality Upgrade Handover

Generated: 2026-08-25T22:01:07+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: /Users/perottochristensen/github/ply/upgrade-quality
- Branch: codex/upgrade-quality
- Base: master at 5635d50
- Measured implementation head and restart preparation base: 70bee0ed16af.
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
2f0a072, 0128cd0, and 70bee0e in roadmap order. The separate operational
continuity implementation is 1b85711 and changes no Go quality denominator.

## Continuity Checkpoint

`codex-dev-start.sh` remains NEXT while P3 is active and P4-P8 are queued in
the machine-readable plan block. Its active archive is
`docs/plan/agent-sessions/2026-08-25T220107+0200-resume-shell-unzip-archive-close.md`.
The blocked archive-close predecessor is answered history, and the reciprocal
archive graph has exactly one NEXT tail.

Normal launch remains a Bash 3.2-compatible non-interactive supervisor with
byte-exact archived prompts, unique external raw JSONL logs, structured stream
validation, post-turn repository revalidation, clean committed handoff
requirements, signal forwarding, and no implicit resume. The stable normalized
skeleton digest remains
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.
`test/codex_dev_start_test.sh` retains all 62 controls.

## P3 Move 53 Preserved

Public `pkg/shell.Unzip(string, string) ([]string, error)` still selects the
complete private `unzipDependencies` composition with `filesystem.System()`.
Only the entry-reader close expression changed: direct `err = rc.Close()` is
now `err = filesystem.CloseReader(dependencies.Files, rc)`.

The existing filesystem adapter gained exactly one
`CloseReader(io.ReadCloser) error` operation, its zero-safe forwarder, and the
exact `reader.Close()` system implementation. Every complete filesystem double
received the mechanical method addition. The dependency receives the exact
opened entry reader already delivered as the Copy source and returns its exact
error without fallback, retry, logging, cleanup, wrapping, or another close
attempt.

Archive open and deferred Close, traversal and entry order, joined paths,
zip-slip behavior, filename append order, directory short-circuit, parent
creation, output OpenFile path/flags/mode, entry open, Copy request and ignored
results, output-file Close, the entry-close assignment and following error
branch, later traversal, partial filenames, and all established error and close
precedence remain unchanged. No other filesystem operation, caller, public API,
or inventory entry changed.

Six focused contracts bring the suite to 358 tests. Adapter contracts prove
zero-safe behavior without closing a supplied reader, exact reader identity and
one dependency request, exact dependency and system errors, one exact system
close, and non-empty recording. Private Unzip contracts prove the exact opened
reader identity, one entry-close request per reached file after Copy and
successful output-file Close, Copy/output-Close/entry-Close order, exact close
error and partial filenames, short-circuit before later entries, suppression
after an output-close error, successful entry Close and later traversal,
complete dependency selection, non-empty recording, and no unrelated
filesystem request. The guarded archive fixture remains below `t.TempDir()`
and injected pipes do not target the repository.

## Measured Quality State

The clean full audit from implementation commit `70bee0e` reports:

- Absolute L0: 8 of 8.
- 358 test functions, zero skipped; 19 of 25 packages have tests.
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
- Q3.4: 0 state-claim phrases across 67 Markdown files at the clean
  implementation checkpoint.
- Acceptance: 4 of 4 host flows pass.
- Full audit: exit 1 for 15 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Clean identity: zero dirty paths.

The authoritative clean report is
`/private/tmp/ply-quality-move53-clean-audit/scorecard.json`; the focused
implementation report is
`/private/tmp/ply-quality-move53-focused-audit/scorecard.json`. Compatibility
reports, pinned tool binaries, Go and linter caches, audit output, raw logs, and
empty-HOME state stayed outside the measured tree.

## Decisions And Learned Facts

1. A ZIP entry's `f.Open()` returns a concrete checksum reader that wraps the
   custom decompressor reader. The exact CloseReader identity is therefore the
   Copy source, while the inner custom reader remains useful for observing the
   eventual close attempt.
2. A zero filesystem dependency returns exact `filesystem.ErrNoFilesystem`
   without invoking a developer-supplied reader's Close method.
3. A dependency-injected entry-close error is returned after Copy and
   output-file Close with the reached filename retained. It prevents every
   later entry; an output-close error prevents the entry-close request.
4. Q1.3 holds exactly at 10/32. The former interface-typed entry-reader Close
   was absent from its violation set, and the new CloseReader adapter request
   is also absent. The scanner and inventory were not changed.
5. Valid red evidence had the adapter tests fail to compile because CloseReader
   was absent, while the private Unzip contracts rejected an empty entry-close
   population and did not return the injected close error.
6. The pinned `apidiff` and `golangci-lint` binaries, all generated reports,
   and every build cache stayed under `/private/tmp`.
7. The initial preflight reached only a sandbox-denied default linter cache.
   Redirecting `GOLANGCI_LINT_CACHE` outside the tree produced the complete
   green preflight with zero linter issues.
8. One explicit `make test-agent-start` run hit the documented nested
   signal/partial-log timing flake; the immediate full rerun passed all 62
   controls, as did the earlier standalone and preflight runs.
9. A fresh archive-close turn proved that `*zip.ReadCloser` implements
   `io.Closer` but not the existing `io.ReadCloser` CloseReader parameter. It
   reverted its attempted red contracts, left the worktree clean, and made no
   commit. The user then explicitly authorized widening only that parameter to
   `io.Closer`, with mechanical interface/helper/system/double updates.
10. The deferred archive Close is absent from the current Q1.3 violation set.
    Regenerate the exact scanner result after moving it rather than assuming a
    denominator change.
11. Clock and server adapter creation belongs to queued P4. Plugin diagrams,
    Maven provenance, remaining HTTP closes, and remaining process effects
    require separate moves. `.quality/inventory` remains checksum-bound and
    unchanged.

## Next Objective

Move only the deferred archive-reader `_ = r.Close()` effect in private
`pkg/shell.Unzip` through
`filesystem.CloseReader(dependencies.Files, r)`. Widen the existing
CloseReader parameter from `io.ReadCloser` to `io.Closer` and update its
interface, zero-safe helper, exact system implementation, and complete doubles
mechanically. Keep public Unzip, private `unzipDependencies`, production
`filesystem.System()` selection, adapter behavior, and every other signature
unchanged.

Start red with focused private Unzip archive-close recording contracts. Prove
the exact `*zip.ReadCloser` identity, one deferred request after every
successful archive open, placement after all reached entry work, execution on
zip-slip and established later error returns, no request after archive-open
failure, ignored injected close error without result or error-precedence
change, complete dependency preservation, distinction from entry-reader-close
requests, non-empty recording, and absence of unrelated filesystem requests.
Launch no process, touch no public network, and change no working directory.

Replace only `defer func() { _ = r.Close() }()` with
`defer func() { _ = filesystem.CloseReader(dependencies.Files, r) }()`. Keep
the defer placement, anonymous function, ignored assignment, archive variable,
surrounding open branch, traversal, and every earlier and later operation
unchanged.

Do not change OpenZipReader, add another filesystem operation, or change
CloseReader beyond the authorized `io.Closer` parameter widening and its
mechanical system/helper/double updates. Do not change entry-reader Close,
output-file Close, Copy, entry open, output OpenFile,
directory creation, another unzip function, Run, Git, structurizr, plugin
diagrams, Maven, HTTP, Spring, clock, server, public API, inventory, audit,
mutation harnesses, or later-roadmap work. Expect Q0.6 to hold at 26 guarded
sites, Q1.1 at 6/25, Q1.2 at zero, Q1.4 at 7/8, and Q2.1 at 0/8. Regenerate
Q1.3 because the deferred archive Close is absent from its current violation
set; require zero comparable ratchet regressions.

## Verification Notes

Completed from implementation commit
`70bee0ed16af4b1e2dec179d6a0ba1ae4237c9e0`:

- Valid red evidence: adapter tests failed to compile because CloseReader did
  not exist; private shell tests observed no adapter reader-close request and
  no injected error. No process or network request ran.
- The first archive-close turn stopped before implementation because
  `*zip.ReadCloser` does not implement `io.ReadCloser`; it reverted its focused
  red contracts, passed the focused baseline tests, and made no commit. The
  successor prompt records the user's `io.Closer` authorization.
- Focused shell/filesystem/Spring and relevant process, Maven, command,
  context, config, HTTP, structurizr, profile, browser, tips, file, template,
  Bitbucket, Wpost, local-config, Kibana, and caller tests: PASS.
- Standalone `/bin/bash` launcher contract: PASS all 62 controls.
- Make meta-contracts, pinned linter with zero issues, complete preflight,
  complete `make test`, install, and all 15 audit meta-controls: PASS with
  external reports and caches.
- API, CLI, and subprocess compatibility: PASS.
- Uncached tests, race, and vet: PASS.
- Host install, status, upgrade, and build acceptance plus their
  meta-contracts: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-criterion audit: exit 1, four improved, two held, zero
  regressed, one non-comparable; Q1.3 is exactly 10/32.
- Clean full audit: exit 1, 15 findings, L0 8/8, five improved, two held, zero
  regressed, one non-comparable, and zero dirty paths.
- One explicit Make launcher run hit the documented nested signal/partial-log
  timing flake; its immediate complete rerun passed all 62 controls.
- The final NEXT archive/launcher contract passed all 62 controls, launcher
  `--check` passed, and the handoff-state Q3.4 audit exited 0 with zero phrases
  across 68 Markdown files and zero ratchet regressions.
- Implementation commit:
  `70bee0ed16af4b1e2dec179d6a0ba1ae4237c9e0`
  (`refactor: route unzip entry close through filesystem adapter`).

## Start And Stop

Read this handover, the linked NEXT archive, P3 and its gate, both design
documents, inventory, complete shell/filesystem-adapter source and tests, every
complete filesystem double, completed Unzip/Spring/process contracts, the
import-aware effect scanner, and relevant callers before editing. Confirm
branch, HEAD, status, reciprocal links, and `./codex-dev-start.sh --check`.
Begin red only for deferred archive Close, then finish with one implementation
commit and one separate handoff-only commit.

Stop before changing archive open, adding another filesystem operation, or
changing CloseReader beyond the authorized `io.Closer` parameter widening and
mechanical updates. Stop before changing entry-reader or output-file Close,
Copy, entry open, output OpenFile, directory creation,
another unzip function, Run, Git, structurizr, plugin diagrams, Maven, HTTP,
Spring, clock, server, public API, inventory, Q1.4, P4, mutation, Docker,
distribution, or publication. Stop on API/CLI change, comparable ratchet
regression, authoritative full-audit exit 2, or failure to isolate
implementation from generated and handoff-only state.
