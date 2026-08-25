# Quality Upgrade Handover

Generated: 2026-08-25T03:41:36+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `838daa153a29`.
- Restart preparation base: `838daa153a29`.
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
838daa1 quality: route clear directory through filesystem adapter
```

The separate operational continuity implementation is:

```text
1b85711 quality: supervise non-interactive agent sessions
```

It changes no Go quality denominator and is separate from P3 move numbering.

## Continuity Checkpoint

`codex-dev-start.sh` stays `NEXT` while P3 is active and P4-P8 are queued in
the machine-readable plan block. Its active archive is
`docs/plan/agent-sessions/2026-08-25T034048+0200-migrate-file-render-create.md`.
The clear-directory predecessor is answered history, and the reciprocal
archive graph has exactly one `NEXT` tail.

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

## P3 Move 22 Preserved

Exported `file.ClearDir(dirPath string, excludes []string) error` retains its
public signature and observable behavior. It is the production wrapper around
a private dependency containing one complete `filesystem.Dependencies` value.
Production selects `filesystem.System()`; no function-valued effect dependency
or new adapter family was added.

The helper delegates the exact `filepath.Join(dirPath, "*")` pattern, consumes
the returned matches without sorting, normalization, or deduplication, and
delegates each selected removal path unchanged through the existing adapter
`RemoveAll`. It retains the complete ordered `strings.Contains(file, exclude)`
loop, including multiple skip logs when several exclusions match. Removal logs
and removals retain glob order. The first removal error short-circuits, and
glob and removal errors retain exact identity without wrapping or retrying.

The existing filesystem adapter gained only zero-value-safe `Glob`; its system
implementation delegates once to `filepath.Glob`. All complete recording
doubles were updated mechanically, and the existing adapter `RemoveAll` was
reused. The safe zero value returns `filesystem.ErrNoFilesystem` before path
access. Recording contracts cover the complete dependency, production system
selection, exact pattern, deliberately unsorted results, ordered exclusions,
logs and removals, nil success, exact error identities, first-error
short-circuiting, safe zero behavior, and non-empty glob and removal
populations. No clear-directory recording contract mutates the real
filesystem.

Only the direct glob and repeated recursive removal in exported `ClearDir`
moved. `Move`, `DeleteAll`, `DeleteSingleFile`, append-open, file-read,
directory-create, file-create, file-overwrite, file-existence, Bitbucket,
Wpost, supervisor, inventory, public API, callers, and every P2A contract
remain unchanged.

## Measured Quality State

The clean full audit at `838daa153a29` reports:

- Absolute L0: 8 of 8.
- 187 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 18 guarded safe-writer sites, 13 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 44 direct external sites outside five declared adapters of 59
  production effect sites. Clock and server are absent, making this ratchet
  non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims across 36 Markdown files.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Full audit: expected exit 1 for 16 documented findings, never exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Measurement identity: clean at tree
  `2d94331a1459ccfd6551ee5140f0a2a2fe2e9f36`, status SHA-256
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`,
  and zero dirty paths.

The clean report is
`/private/tmp/ply-clear-dir-clean-audit.2ndObO/report`. The focused report is
`/private/tmp/ply-clear-dir-focused.2FjmNN/report`. Preflight, compatibility,
Go-gate, acceptance, and hermetic outputs were kept under `/private/tmp`, so no
generated report entered the measured commit.

## Decisions And Learned Facts

1. The private clear-directory dependency carries the complete
   `filesystem.Dependencies` value. It does not reconstruct or store a
   function-valued effect dependency.
2. Adapter `Glob` is zero-value-safe and returns ordered matches plus the exact
   dependency error; production selection and adapter delegation are covered
   independently.
3. The exclusion loop deliberately does not break after a match. A path that
   contains two exclusions produces two skip logs and no removal, preserving
   the original complete loop and logging behavior.
4. Q1.3 moved from 46/60 to 44/59 rather than the nominal 44/60. The two direct
   `ClearDir` sites disappeared, system `Glob` entered the adapter population,
   and reused system `RemoveAll` was already counted, contracting the total
   effect-site denominator by one.
5. Exported `file.Render` is the next isolated flow. It has one direct
   `os.Create(outputFilePath)` after its already-adapter-backed input read and
   template parse. Reusing adapter `Create` should move Q1.3 nominally from
   44/59 to 43/58; regenerate the exact result.
6. `Render` must preserve input-read ordering and errors, `template.Must` parse
   panic behavior, the exact output path, create/truncation behavior, emitted
   bytes, execution errors, nil success, and the existing absence of an added
   close.
7. A function-valued effect dependency makes Q1.3 fail closed. Production
   moves use resolvable interfaces and complete dependency values.
8. `.quality/inventory` is baseline-checksum-bound. Do not relabel seams or
   claim P5 mutation coverage.
9. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
   `make preflight`; Make command-line values propagate through `MAKEFLAGS`
   and defeat the missing-binary mutant.
10. Set `GOLANGCI_LINT_CACHE` and `GOCACHE` to writable external directories
    when the default caches reject writes. Keep every generated report outside
    the measured tree.
11. The ignored `.agent-task/current.md` compatibility mirror is absent.
    Tracked launcher, handover, archive, and plan state govern the next session.
12. Supervisor success is the conjunction of process exit, structured
    terminal stream, and committed repository evidence. Final prose is
    observable only.
13. The partial-raw-log signal fixture can still flake. An immediate complete
    rerun must pass all 62 controls; never weaken or skip the contract. The
    first standalone launcher and first `make test` runs flaked during move 22;
    both immediate complete reruns passed.

## Next Objective

Move only exported
`file.Render(inputFilePath string, outputFilePath string, r interface{}) error`'s
direct `os.Create(outputFilePath)` behind the existing filesystem adapter's
`Create`. Start with recording contracts in `pkg/file`: require one complete
filesystem dependency, exact input and output paths, input-read-before-create
ordering, exact create-error identity, exact rendered bytes, unchanged writer
execution errors, nil success, safe zero-value behavior without developer path
access, production selection of `filesystem.System()`, and non-empty read,
create, and write populations. Perform no real filesystem mutation in the
recording contracts.

Keep `Render`'s public signature and every observable result unchanged. Use one
private complete dependency boundary with the exported function as the
production wrapper. Reuse existing adapter `ReadFile` for the input and
existing adapter `Create` for the output through that complete dependency so
the private zero value is safe; do not add an adapter operation. Preserve
`template.Must` parse panic behavior, creation only after successful read and
parse, output truncation, emitted bytes, execution errors, the existing lack of
an added close, and both existing template callers. Do not move another file
operation, `pkg/shell`, or enter Bitbucket, HTTP/process, config behavior,
clock/server, P4, P5, or later roadmap work.

Expected direction is one fewer Q1.3 violation and one fewer production effect
site, nominally 43 of 58, but regenerate the exact structured measurement and
accept it only with zero comparable ratchet regressions. Q1.2, Q1.4, and exact
Q2.1 should remain unchanged.

## Verification Notes

Completed from implementation commit `838daa153a29`:

- Red clear-directory evidence: file contracts first failed because the
  private dependency constructor and helper did not exist; adapter contracts
  failed because `Glob` did not exist.
- Focused filesystem, file, Bitbucket, HTTP, and caller-adjacent tests: PASS.
- `/bin/bash test/codex_dev_start_test.sh`: the first run hit the known signal
  fixture flake; the immediate complete rerun passed all 62 controls.
- Make preflight meta-contracts and complete `make preflight`: PASS, including
  all 15 audit meta-controls.
- API/CLI compatibility and subprocess contracts: PASS.
- `make test`: the first run hit the same known signal fixture flake; the
  immediate complete rerun passed. `make test-install`: PASS.
- Full tests, race tests, and `go vet`: PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-ratchet audit: expected exit 1; four improved, two held, zero
  regressed, and one not comparable.
- Clean full audit: expected exit 1, 16 documented findings, L0 8 of 8, five
  improved, two held, zero regressed, and one not comparable.
- `git diff --check`: PASS before the implementation commit.
- Implementation commit: `838daa153a29d13603a9f3bc1d532321f6db64d9`
  (`quality: route clear directory through filesystem adapter`).

## Start

1. Read this handover, the linked NEXT archive, the P3 section and checkpoint
   gate in `docs/plan/quality-upgrade.md`, both design documents, inventory,
   complete file implementation/tests, Render callers, adapter, complete
   recording doubles, and the named audit implementations before editing.
2. Confirm branch, HEAD, status, reciprocal archive links, and
   `./codex-dev-start.sh --check`.
3. Reproduce focused clean baselines as needed, then begin red with recording
   contracts for only Render's existing ReadFile/Create composition.
4. Finish with one focused implementation commit and one separate handoff-only
   commit. Leave the launcher `NEXT`; do not launch a successor.

## Stop Conditions

- Stop before another file operation, an adapter interface extension, another
  adapter family, Q1.4 expansion, P4, mutation harnesses, Docker, cloud,
  distribution, or publication.
- Stop if the public API or CLI contract would change.
- Stop if a comparable ratchet regresses.
- Stop if the full audit exits 2.
- Stop if tracked source changes are not isolated from generated reports and
  handoff-only state.
