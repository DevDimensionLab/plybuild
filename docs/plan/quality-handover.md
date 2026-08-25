# Quality Upgrade Handover

Generated: 2026-08-25T03:03:37+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `acda4e335348`.
- Restart preparation base: `acda4e335348`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no product implementation changes.
- No push, merge, release, publication, stash, revert, successor launch, or
  worktree removal was performed.

P3 implementation commits are:

```text
5c6f2fa quality: move process exits to main
03d6242 quality: move git processes behind adapter
f5ee37d quality: move Maven processes behind adapter
ee5e9ab quality: cover cloud clone seam
ffc4e77 quality: move Maven metadata HTTP behind adapter
204e222 quality: move template copy behind filesystem adapter
89d0f76 quality: move Spring download behind adapters
dee214c quality: move Spring discovery behind HTTP adapter
789ae23 quality: move Bitbucket JSON behind HTTP adapter
07ac6ce quality: move Kibana POST behind HTTP adapter
a7eb3ef quality: move Wpost behind HTTP and filesystem adapters
e13a036 quality: move Bitbucket selection behind filesystem adapter
f59a3f0 quality: move file existence behind filesystem adapter
61714a5 quality: move file overwrite behind filesystem adapter
60e5aac quality: move file create behind filesystem adapter
e054082 quality: move directory create behind filesystem adapter
27c0d1a quality: route file open through filesystem adapter
c2f3597 quality: route file append open through filesystem adapter
9e2d669 quality: route single-file delete through filesystem adapter
a4deb76 quality: route recursive delete through filesystem adapter
acda4e3 quality: route file move through filesystem adapter
```

The separate operational continuity implementation is:

```text
1b85711 quality: supervise non-interactive agent sessions
```

It changes no Go quality denominator and is separate from P3 move numbering.

## Continuity Checkpoint

`codex-dev-start.sh` stays `NEXT` while P3 is active and P4-P8 are queued in
the machine-readable plan block. Its active archive is
`docs/plan/agent-sessions/2026-08-25T030337+0200-migrate-file-clear-dir.md`.
The file-move predecessor is answered history, and the reciprocal archive
graph has exactly one `NEXT` tail.

Normal launch is a Bash 3.2-compatible, non-interactive supervisor. Each
generation resolves an external Codex executable and invokes exact
`codex exec` arguments for normal service tier, workspace-write, the repository
working directory, JSONL, an external final-message path, and the byte-exact
validated prompt. It does not use interactive mode or resume a thread.

Every turn receives a unique physical log directory outside the worktree. Raw
JSONL is preserved byte-for-byte while concise progress is streamed. The
parser decodes JSON, treats messages only as data, requires the characterized
thread/turn lifecycle and one final `turn.completed`, and rejects empty,
malformed, truncated, contradictory, failed, error, post-terminal, and
non-zero-exit streams.

After a successful child exit, the parent re-reads the on-disk launcher and
validates branch/root, queue, full archive graph, prompt, and normalized stable
skeleton. A next turn requires a clean changed HEAD, an answered old archive,
a changed session ID, exactly one new committed `NEXT` archive, and reciprocal
history. A valid `COMPLETE` requires the whole authorized queue complete.
Signals are forwarded to the active child/parser and never start a successor.

`test/codex_dev_start_test.sh` has 62 controls. They characterize exact
argv/prompt bytes, logging/progress, malicious event data, two-generation
continuation, `COMPLETE`, every terminal/JSON/process/no-progress/dirty/handoff/
contract failure, and signal interruption without using the real Codex
executable. Mutable header and prompt data are inert after the stable execution
boundary; the pinned normalized skeleton digest is
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.

## P3 Move 21 Preserved

Exported `file.Move(source, destination string) error` retains its public
signature and observable behavior. It is the production wrapper around a
private dependency containing one complete `filesystem.Dependencies` value.
Production selects `filesystem.System()`; no function-valued effect dependency
or new adapter family was added.

The helper delegates both complete paths unchanged through the adapter and
returns the dependency's exact error. It does not probe, normalize, wrap,
retry, copy, remove, or add selection behavior. The safe zero value returns
`filesystem.ErrNoFilesystem` without developer path access.

The existing filesystem adapter gained only `Rename`. Its system
implementation delegates once to `os.Rename`, and all complete recording
doubles were updated mechanically. Recording contracts cover the complete
dependency, exact source and destination, nil success, exact dependency-error
identity, production system selection, safe zero behavior, and a non-empty
recorded population. No file-move recording contract performs real filesystem
mutation.

Only the direct `os.Rename(source, destination)` in exported `file.Move` moved.
`DeleteAll`, `DeleteSingleFile`, `ClearDir`, every other file operation,
inventory, public API, all callers, Bitbucket, Wpost, supervisor, and every P2A
contract remain unchanged.

## Measured Quality State

The clean full audit at `acda4e335348` reports:

- Absolute L0: 8 of 8.
- 180 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 18 guarded safe-writer sites, 13 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 46 direct external sites outside declared adapters of 60 production
  effect sites. Clock and server are absent, making this ratchet
  non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Measurement identity: clean at tree
  `8e744e35dfa795e1eb086d50b9de792282c8165f`, status SHA-256
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`,
  and zero dirty paths.

The clean report is
`/private/tmp/ply-file-move-clean-audit.YE9gz5/report`. The focused report is
`/private/tmp/ply-file-move-focused.1f2lbZ/report`. Preflight, compatibility,
Go-gate, acceptance, and hermetic outputs were kept under `/private/tmp`, so no
generated report entered the measured commit.

## Decisions And Learned Facts

1. The private file-move dependency carries the complete
   `filesystem.Dependencies` value and delegates exactly once; it does not
   reconstruct or store a function-valued effect dependency.
2. Adapter `Rename` is zero-value-safe and returns the exact dependency error;
   production selection and exact adapter delegation are independently covered.
3. The exported wrapper preserves both paths without probing, normalization,
   wrapping, retry, copying, removal, or new selection behavior.
4. `file.ClearDir(dirPath string, excludes []string) error` is the next isolated
   filesystem flow. It has one direct `filepath.Glob(filepath.Join(dirPath,
   "*"))` selection and one repeated `os.RemoveAll(file)` call. Moving both
   should change Q1.3 nominally from 46/60 to 44/60.
5. `ClearDir` must preserve the exact glob pattern, glob result order, existing
   ordered substring-exclusion checks, log behavior, first removal-error
   short-circuit, and exact glob or removal error.
6. A function-valued effect dependency makes Q1.3 fail closed. Production
   moves use resolvable interfaces and complete dependency values.
7. `.quality/inventory` is baseline-checksum-bound. Do not relabel seams or
   claim P5 mutation coverage.
8. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
   `make preflight`; Make command-line values propagate through `MAKEFLAGS`
   and defeat the missing-binary mutant.
9. Set `GOLANGCI_LINT_CACHE` and `GOCACHE` to writable external directories
   when the default caches reject writes. Keep every generated report outside
   the measured tree.
10. The ignored `.agent-task/current.md` compatibility mirror is absent.
    Tracked launcher, handover, archive, and plan state govern the next session.
11. Supervisor success is the conjunction of process exit, structured
    terminal stream, and committed repository evidence. Final prose is
    observable only.
12. The partial-raw-log signal fixture can still flake. An immediate complete
    rerun must pass all 62 controls; never weaken or skip the contract. It did
    not flake during move 21.

## Next Objective

Move only exported `file.ClearDir(dirPath string, excludes []string) error`'s
direct glob selection and repeated recursive removal behind a narrow extension
of the existing filesystem adapter. Start with recording contracts in
`pkg/file` and the adapter: require the complete dependency, exact glob pattern,
ordered matches, exact removal paths, existing exclusion decisions, exact
returned errors, first-error short-circuit, safe zero-value behavior with no
developer path access, production selection of `filesystem.System()`, and
non-empty recorded glob and removal populations. Perform no real filesystem
mutation in the file recording contracts.

Keep `ClearDir`'s public signature and every observable result unchanged. Use
one private complete dependency boundary, with the exported function as the
production wrapper. Add only the narrow `Glob` operation to
`internal/adapter/filesystem`, reuse its existing `RemoveAll`, and update
complete recording doubles mechanically. Do not move another file operation,
`pkg/shell`, or enter Bitbucket, HTTP/process, config behavior, clock/server,
P4, P5, or later roadmap work.

Expected direction is two fewer Q1.3 violations with the same production
effect population, nominally 44 of 60, but regenerate the exact structured
measurement and accept it only with zero comparable ratchet regressions.
Q1.2, Q1.4, and exact Q2.1 should remain unchanged.

## Verification Notes

Completed from implementation commit `acda4e335348`:

- Red file-move evidence: recording contracts first failed because the private
  dependency constructor and helper did not exist; adapter contracts failed
  because `Rename` did not exist.
- Focused filesystem, file, config caller, Bitbucket, and HTTP tests: PASS.
- `/bin/bash test/codex_dev_start_test.sh`: PASS, all 62 controls.
- Make preflight meta-contracts and complete `make preflight`: PASS, including
  all 15 audit meta-controls.
- API/CLI compatibility and subprocess contracts: PASS.
- `make test`, `make test-install`, full tests, race tests, and `go vet`: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-ratchet audit: expected exit 1; four improved, two held, zero
  regressed, and one not comparable.
- Clean full audit: expected exit 1, 16 documented findings, L0 8 of 8, five
  improved, two held, zero regressed, and one not comparable.
- `git diff --check`: PASS before the implementation commit.
- Implementation commit: `acda4e335348660d4bfd87b3e6005395ba67d42d`
  (`quality: route file move through filesystem adapter`).

## Start

1. Read this handover, the linked NEXT archive, the P3 section and checkpoint
   gate in `docs/plan/quality-upgrade.md`, both design documents, inventory,
   complete file implementation/tests, callers, adapter, complete recording
   doubles, and the named audit implementations before editing.
2. Confirm branch, HEAD, status, reciprocal archive links, and
   `./codex-dev-start.sh --check`.
3. Reproduce focused clean baselines as needed, then begin red with recording
   contracts for only `ClearDir`, adapter `Glob`, and the existing adapter
   `RemoveAll` composition.
4. Finish with one focused implementation commit and one separate handoff-only
   commit. Leave the launcher `NEXT`; do not launch a successor.

## Stop Conditions

- Stop before another file operation, another adapter family, Q1.4 expansion,
  P4, mutation harnesses, Docker, cloud, distribution, or publication.
- Stop if the public API or CLI contract would change.
- Stop if a comparable ratchet regresses.
- Stop if the full audit exits 2.
- Stop if tracked source changes are not isolated from generated reports and
  handoff-only state.
