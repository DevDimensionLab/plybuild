# Quality Upgrade Handover

Generated: 2026-08-25T01:04:57+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `27c0d1a34334`.
- Restart preparation base: `27c0d1a34334`.
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
```

The separate operational continuity implementation is:

```text
1b85711 quality: supervise non-interactive agent sessions
```

It changes no Go quality denominator and is separate from P3 move numbering.

## Continuity Checkpoint

`codex-dev-start.sh` stays `NEXT` while P3 is active and P4-P8 are queued in
the machine-readable plan block. Its active archive is
`docs/plan/agent-sessions/2026-08-25T010457+0200-migrate-file-open-file.md`.
The file-read predecessor is answered history, and the reciprocal archive graph
has exactly one `NEXT` tail.

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

## P3 Move 17 Preserved

Exported `file.Open(filePath string) ([]byte, error)`, all resource, config,
Kibana, render, and package-local callers, and every P2A API/CLI/subprocess
contract are unchanged. The package owns a private complete file-read
dependency containing `filesystem.Dependencies`. Production selects
`filesystem.System()`, and the private helper delegates the complete path to
the existing `filesystem.ReadFile` operation.

Successful reads return the dependency's exact complete bytes. Every dependency
error discards partial data and returns a non-nil empty `[]byte{}` with that
exact error. The zero dependency returns `filesystem.ErrNoFilesystem` and a
non-nil empty result without accessing the supplied path. Four recording
contracts prove the complete dependency and path, exact empty, multiline, and
binary bytes, partial-data error normalization, exact error identity,
production system selection, safe defaults, and a non-empty recorded
population.

Only `Open`'s direct `os.Open`, `io.ReadAll`, and ignored deferred close were
removed. The filesystem adapter contract, inventory, seam drivers, mutation
labels, public API, and callers did not change. `file.OpenFile` keeps its direct
`os.OpenFile` append-open operation and is the next isolated flow.

## Measured Quality State

The clean full audit at `27c0d1a34334` reports:

- Absolute L0: 8 of 8.
- 157 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 18 guarded safe-writer sites, 13 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 50 direct external sites outside declared adapters of 60 production
  effect sites. Clock and server are absent, making this ratchet
  non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Full audit: expected exit 1 for 16 documented findings.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Measurement identity: clean at tree
  `7efbc5d33c18c6c54fc593bd7f5007f720d7c278`, status SHA-256
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`,
  and zero dirty paths.

The clean report is
`/private/tmp/ply-file-open-clean.ujkFUr/report-final`. Focused, preflight, and
compatibility reports were also kept under `/private/tmp`, so no generated
report entered the measured commit.

## Decisions And Learned Facts

1. The private file-read dependency requires the complete
   `filesystem.Dependencies` value. Production selects `filesystem.System()`;
   the helper is independently recordable without a function-valued effect
   dependency.
2. `filesystem.ReadFile` already supplies the complete operation required by
   `Open`, so the adapter contract did not change.
3. The explicit error branch is necessary because a dependency can return
   partial bytes with an error; the legacy contract discards those bytes and
   returns a non-nil empty result with the exact error.
4. Exported `file.OpenFile` is the next isolated filesystem flow. Its existence
   and empty-create operations already have adapter support; only append-open
   needs a narrow existing-adapter extension. Moving one direct call into one
   adapter implementation should change Q1.3 nominally from 50/60 to 49/60.
5. A function-valued effect dependency makes Q1.3 fail closed. Production
   moves use resolvable interfaces and complete dependency values.
6. `.quality/inventory` is baseline-checksum-bound. Do not relabel seams or
   claim P5 mutation coverage.
7. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
   `make preflight`; Make command-line values propagate through `MAKEFLAGS`
   and defeat the missing-binary mutant.
8. Set `GOLANGCI_LINT_CACHE` and `GOCACHE` to writable external directories
   when the default caches reject writes. Keep every generated report outside
   the measured tree.
9. The tracked full-audit baseline is
   `.quality/baseline/scorecard.json`. One clean-clone invocation used the
   nonexistent `.quality/baseline.json` and failed closed before ratchet
   comparison; the corrected complete audit produced the authoritative exit 1
   report above.
10. The ignored inherited `.agent-task/current.md` is only a compatibility
    mirror, not task authority. The clean local clone produced the
    implementation checkpoint identity without modifying that file.
11. Supervisor success is the conjunction of process exit, structured
    terminal stream, and committed repository evidence. Final prose is
    observable only.

## Next Objective

Move only exported `file.OpenFile`'s append-open operation behind a narrow
extension of the existing filesystem adapter. Start with recording contracts
in `pkg/file` and the adapter: require the complete dependency and path, exact
`os.O_APPEND|os.O_WRONLY` flags, mode `0644`, exact returned `*os.File` and
error, existing/missing/non-missing-stat-error sequencing, exact empty creation,
create-error short circuit, safe zero-value behavior, production selection of
`filesystem.System()`, and non-empty recorded populations. Perform no real
filesystem mutation.

Keep `file.OpenFile(fileName string) (*os.File, error)`, `Exists`, `CreateFile`,
`Open`, and every observable result unchanged. Use one private complete
dependency boundary, with the exported function as the production wrapper.
Drive existence, conditional empty write, and append-open from that same
dependency. Add only the narrow `OpenFile` operation to
`internal/adapter/filesystem` and update complete recording doubles
mechanically. Do not move `pkg/shell`'s direct `os.OpenFile` or another file
operation, and do not enter Bitbucket, HTTP/process, config behavior,
clock/server, P4, P5, or later roadmap work.

Expected direction is one fewer Q1.3 violation with the same production effect
population, nominally 49 of 60, but regenerate the exact structured
measurement and accept it only with zero comparable ratchet regressions.
Q1.2, Q1.4, and exact Q2.1 should remain unchanged.

## Verification Notes

Completed from clean implementation commit `27c0d1a34334`:

- Red file-read evidence: the new contracts failed to compile because the
  private dependency constructor and helper did not exist.
- Focused file/filesystem and all relevant caller tests: PASS.
- `/bin/bash test/codex_dev_start_test.sh`: PASS, 62 controls.
- Make preflight meta-contracts and a complete `make preflight` rerun: PASS,
  including 15 audit meta-controls.
- API/CLI and subprocess compatibility: PASS.
- `make test`, `make test-install`, uncached tests, race tests, and `go vet`:
  PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-ratchet audit: expected exit 1 with Q1.3 at 50 of 60, four
  improved, two held, zero regressed, and one not comparable.
- Clean full audit: expected exit 1, 16 documented findings, L0 8 of 8, five
  improved, two held, zero regressed, and zero dirty paths.
- The first complete `make preflight` run and the first handoff launcher run
  each hit the documented nested partial-raw-log signal-fixture flake; the
  standalone launcher and immediate complete reruns passed all 62 controls.
- One clean-audit command used the wrong baseline path and exited 2 before
  comparison; the corrected tracked-baseline run completed with the
  authoritative exit 1 result.

Tool paths used were `/private/tmp/ply-p2b-api.4umBuM/bin/apidiff` and
`/private/tmp/ply-p2b-lint.SGWVGp/bin/golangci-lint`; probe before reuse.

Environment: host Go 1.26.2 on Darwin arm64, module Go 1.18, `/bin/bash`
3.2.57, PATH Bash 5.3.9, and `golangci-lint` 2.12.2.

## Start

From any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

This starts the non-interactive supervisor and may run successive fresh
missions after valid handoffs. Do not invoke it while validating the handoff;
use `--check` or `--print-prompt`. `.agent-task/current.md` is not task
authority.

## Stop Conditions

Stop and report rather than forcing progress when:

- exported `file.OpenFile` behavior, its signature, or compatibility cannot be
  preserved;
- the filesystem dependency is incomplete, unresolvable, function-valued, or
  unsafe at its zero value;
- the effect cannot be isolated without moving another file operation or
  adding more than the narrow existing-adapter operation;
- a comparable ratchet regresses, the audit exits 2 after a valid invocation,
  or the tree cannot be measured cleanly;
- the work requires P4-P8 implementation, publication, distribution, or a
  real successor launch; or
- the focused `file.OpenFile` move and its separate automatic handoff are
  complete.
