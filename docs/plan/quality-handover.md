# Quality Upgrade Handover

Generated: 2026-08-25T00:31:30+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `e054082b5210`.
- Restart preparation base: `e054082b5210`.
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
```

The separate operational continuity implementation is:

```text
1b85711 quality: supervise non-interactive agent sessions
```

It changes no Go quality denominator and is separate from P3 move numbering.

## Continuity Checkpoint

`codex-dev-start.sh` remains `NEXT` while P3 is active and P4-P8 remain queued
in the machine-readable plan block. Its active archive is
`docs/plan/agent-sessions/2026-08-25T003130+0200-migrate-file-open.md`.
The directory-create predecessor is answered history, and the reciprocal archive
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
executable. Mutable header and prompt data remain inert after the stable
execution boundary; the pinned normalized skeleton digest is
`4755da4dd8645ac890df241d329319a130ac5d778c0bd127061c9667afb2d484`.

## P3 Move 16 Preserved

Exported `file.CreateDirectory(path string) error`, all callers, and every P2A
API/CLI/subprocess contract remain unchanged. The package owns a private
complete directory-create dependency containing `filesystem.Dependencies`.
Production selects `filesystem.System()`, and the private helper delegates the
probe and recursive creation to `filesystem.Stat` and `filesystem.MkdirAll`.

Only a missing-path stat error selects exactly one creation with mode `0755`.
Existing paths and other stat errors still return nil, and a failed creation
still returns the original missing-path stat error. The dependency's zero
value is safely non-mutating and returns nil. Five recording contracts prove
the complete path and dependency, both branches, mode, exact errors,
production system selection, safe defaults, and non-empty stat and creation
populations. The adapter contract, inventory, seam drivers, mutation labels,
public API, and callers did not change. `file.OpenFile` retains its direct
`os.OpenFile` behavior.

## Measured Quality State

The clean full audit at `e054082b5210` reports:

- Absolute L0: 8 of 8.
- 153 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 18 guarded safe-writer sites, 13 write and 5 copy, with zero unsafe
  direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 52 direct external sites outside declared adapters of 62 production
  effect sites. Clock and server are absent, making this ratchet
  non-comparable.
- Q1.4: 7 of 8 declared seams covered.
- Exact Q2.1: 0 of 8 subjects have an executable harness.
- Q3.4: 0 temporary-state claims.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Full audit: expected exit 1 for 16 documented findings, never 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is the
  single non-comparable ratchet.
- Measurement identity: clean at tree
  `228362c84985d7ea1c9bdd1a6d5588d247c4319b`, status SHA-256
  `6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d`
  and zero dirty paths.

The clean report is
`/private/tmp/ply-file-directory-clean.3NSITS/report`. Focused, preflight, and
compatibility reports were also kept under `/private/tmp`, so no ignored
report entered the measured tree.

## Decisions And Learned Facts

1. The private directory-create dependency requires the complete
   `filesystem.Dependencies` value. Production selects `filesystem.System()`;
   the helper remains independently recordable without a function-valued
   effect dependency.
2. `filesystem.Stat` and `filesystem.MkdirAll` already supply the complete
   operations required by `CreateDirectory`; the adapter contract did not
   change.
3. The zero dependency is safely non-mutating and returns nil because
   `filesystem.ErrNoFilesystem` is not a missing-path error, preserving the
   legacy non-missing-stat-error branch.
4. Exported `file.Open` is the next isolated filesystem flow. The existing
   adapter's `ReadFile` operation can replace its direct open/read/close flow
   while exact bytes, errors, signature, and callers stay fixed.
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
9. The ignored inherited `.agent-task/current.md` is not task authority. A
   clean local clone produced the implementation checkpoint identity without
   deleting or changing that user-owned file.
10. Supervisor success is the conjunction of process exit, structured
    terminal stream, and committed repository evidence. Final prose is
    observable only.

## Next Objective

Move only exported `file.Open` behind the existing filesystem adapter. Start
with recording contracts in `pkg/file`: require the complete path and
dependency; prove representative exact bytes, dependency-error propagation
with an empty result even when the dependency supplies partial bytes, safe
zero-value behavior, production selection of `filesystem.System()`, and a
non-empty recorded read population. Exercise no developer path and perform no
real filesystem mutation.

Keep `file.Open(filePath string) ([]byte, error)`, all its resource, config,
render, Kibana, and package-local callers, `file.OpenFile`, and every observable
result unchanged. Use a private complete dependency boundary, with the
exported function as the production wrapper. Delegate only the current
open/read/ignored-close flow to `filesystem.ReadFile`; do not move another file
operation, modify the adapter contract, or enter Bitbucket, HTTP/process,
config behavior, clock/server, P4, P5, or later roadmap work.

Expected direction is two fewer Q1.3 violations and two fewer production
effect sites, nominally 50 of 60, but regenerate the exact structured
measurement and accept it only with zero comparable ratchet regressions.
Q1.2, Q1.4, and exact Q2.1 should remain unchanged.

## Verification Notes

Completed from clean implementation commit `e054082b5210`:

- Red directory-create evidence: new contracts failed to compile because the
  private dependency constructor and helper did not exist.
- Focused file/filesystem and all relevant caller tests: PASS.
- `/bin/bash test/codex_dev_start_test.sh`: PASS, 62 controls.
- Make preflight meta-contracts and complete `make preflight`: PASS, including
  15 audit meta-controls.
- API/CLI and subprocess compatibility: PASS.
- `make test`, `make test-install`, uncached tests, race tests, and `go vet`:
  PASS.
- Host install, status, upgrade, and build acceptance: PASS, 4 of 4.
- Empty-HOME `go test ./... -count=2`: PASS.
- Focused seven-ratchet audit: expected exit 1 with Q1.3 at 52 of 62, four
  improved, two held, zero regressed, and one not comparable.
- Clean full audit: expected exit 1, 16 documented findings, L0 8 of 8, five
  improved, two held, zero regressed, and zero dirty paths.
- One repeated `make test` hit the documented partial-raw-log signal-fixture
  flake; standalone launcher, complete preflight, and the immediate full rerun
  each passed all 62 controls.

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

- exported `file.Open` behavior, its signature, or caller
  compatibility cannot be preserved;
- the filesystem dependency is incomplete, unresolvable, function-valued, or
  unsafe at its zero value;
- the effect cannot be isolated without moving another file operation or
  changing the existing adapter contract without proof;
- a comparable ratchet regresses, the audit exits 2, or the tree cannot be
  measured cleanly;
- the work requires P4-P8 implementation, publication, distribution, or a
  real successor launch; or
- the focused `file.Open` move and its separate automatic handoff
  are complete.
