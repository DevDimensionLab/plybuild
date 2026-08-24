# Quality Upgrade Handover

Generated: 2026-08-24T12:22:34+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `5c6f2fa33719`.
- Restart preparation base: `5c6f2fa33719`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no implementation changes.
- No push, merge, release, publication, stash, revert, or worktree removal was
  performed.

P3 move 1 implementation commit:

```text
5c6f2fa quality: move process exits to main
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main move is
complete, P3 adapter moves remain active, and P4-P8 remain queued in the
machine-readable plan block. The active archive is
`docs/plan/agent-sessions/2026-08-24T122234+0200-introduce-process-adapter.md`.
Its predecessor is answered history, and the connected graph has exactly one
`NEXT` tail.

The launcher validates its attached branch/root, the authorized queue, strict
session metadata, non-symlink planning inputs, a reciprocal archive graph,
every historical prompt digest, exact launcher/archive prompt bytes, and the
external Codex executable. Mutable prompt data remains comment-encoded after
the stable execution boundary. A dirty worktree is allowed with a static
warning; no status content enters the prompt.

`test/codex_dev_start_test.sh` derives active metadata from the launcher. Its 49
controls include premature `COMPLETE` rejection while authorized work remains,
so normal handoff requires no test edit. The ignored `.agent-task` path remains
absent and is not task authority.

## P3 Move 1 Result

The first P3 move completed as one measured implementation change:

1. `cmd/entry_boundary_test.go` began red with both entry paths unbound, a
   missing error-returning command path, and all 127 process-terminating calls.
   Its syntax-aware scan now requires both executable files to delegate to the
   same `cmd.ExecuteE() error`, requires exit 1 selection inside each `main`,
   and rejects `log.Fatal*` or `os.Exit` everywhere else.
2. `cmd.ExecuteE()` owns Cobra execution and returns errors. The exported
   `cmd.Execute()` remains `func()` for API compatibility but no longer exits;
   `cmd.RootCmd` is unchanged. The compatible new symbol is recorded in the
   API allowlist. Root/status/upgrade/build help and unknown-command stdout,
   stderr, and exit 1 behavior remain green.
3. Cobra initializers, pre-run hooks, run hooks, and their unexported helpers
   now propagate errors. Both `main.go` and `cmd/ply/main.go` print returned
   errors and call `os.Exit(1)`. The two remaining package-level fatal sites
   were converted to non-terminating logging without changing exported
   signatures.
4. No process, HTTP, filesystem, clock, or server adapter was introduced. That
   ordered work begins with the next move.

## Measured Quality State

The clean full audit at `5c6f2fa33719` reports:

- Absolute L0: 8 of 8.
- 39 test functions, zero skipped; 9 of 22 packages have tests.
- Q0.8: 0 of 7 production scripts lack a meta-test.
- Process-exiting calls outside `main`: 0; Q1.2 is an absolute PASS and an
  improved ratchet from 127.
- Packages without tests: 13 of 22, improved from the stored 15-package debt.
- Direct production effects outside the five declared adapters: 80 of 80.
- Declared seam swap tests: 0 of 8.
- Mutation harnesses: 0 of 8.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Q2.8 and Q2.9 remain honestly `UNMEASURABLE` pending criterion-bound manual
  evidence; the executable magnitude, bad-input, and read-only controls exist.
- Full audit exit: expected 1 for documented project findings; no audit returned
  exit 2.
- Measurement identity: clean, with zero dirty paths.

The implementation audit exposed one inherited Q3.4 regression from a host
tool sentence added by restart commit `2034f63`. This handoff replaces that
sentence with a dated probe result. The final clean handoff audit must therefore
show four improved, three held, and zero regressed comparable ratchets.

Regenerate `target/quality-audit/scorecard.json`; it is ignored output, not
persistent evidence. Compatibility checks create ignored reports under
`target/compatibility`; direct them outside the worktree or remove those
generated reports before the full audit.

## Decisions And Learned Facts

1. Repository-wide local search found no Ply consumer of `cmd.Execute()` beyond
   the original root executable. The similarly named Trip functions belong to
   another Go module. No evidence requires `cmd.Execute()` to terminate.
2. `cmd.ExecuteE()` temporarily silences Cobra while it executes, restores the
   public `RootCmd` settings, reproduces unknown-command diagnostics, maps the
   private documentation-complete sentinel to success, and returns every other
   error to `main`.
3. Cobra's initializer callback cannot return an error. It now records the
   initialization error, and `InitGlobals` or the build boundary returns it
   before command work begins.
4. Returning `flag.ErrHelp` from `OkHelp` lets Cobra stop a command and retain
   successful help behavior without a process exit.
5. The clean gate improved Q1.1 because the root and `cmd` entry packages gained
   tests; this was a consequence of the required boundary contract, not adapter
   scope.
6. A host probe on 2026-08-24 found no GoReleaser executable, and no GoReleaser
   or publisher command was invoked. The retained local snapshot contract is
   outside P3 execution scope.
7. Passing `GOLANGCI_LINT` as a Make command-line variable propagates through
   `MAKEFLAGS` and defeats the lint meta-test's missing-binary mutant. Provide
   `APIDIFF` and `GOLANGCI_LINT` as environment variables for `make preflight`.
8. The sandbox denies the default Go and golangci-lint cache paths. Use isolated
   `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` directories under
   `/private/tmp`.

## Next Objective

P3 remains active. Introduce `internal/adapter/process` and migrate one coherent
git flow from `pkg/shell/git.go`. Start red with recording argument-order tests
for the declared `git-process` and `git-commit` seams: clone must keep URL before
target directory, and commit must keep target directory before message. The
recording double must fail on an empty call population and preserve the complete
dependency value passed to production code.

Reduce Q1.3 by exactly the migrated git call sites and increase Q1.4 only for
the swaps the new tests kill. Preserve the Q1.2 zero boundary, Go 1.18, all P2A
API/CLI/subprocess and host acceptance behavior, and zero comparable ratchet
regressions. Do not migrate Maven or cloud calls in the same move, and do not
enter HTTP, filesystem, clock, server, Docker, dependency, distribution,
Spring, or mutation-harness scope.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `5c6f2fa33719`:

- Focused entry-boundary tests for both executables: PASS.
- Both root and `cmd/ply` fresh binaries produced byte-identical help and
  unknown-command streams with exit 1 for bad input.
- API and CLI compatibility plus both compatibility meta-tests: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- `bash .quality/tools/test-quality-audit.sh`: PASS, 15 controls.
- Clean full audit: expected exit 1, L0 8 of 8, Q1.2 0, 39 tests, 9 tested
  packages, and zero dirty paths. The inherited handoff wording is corrected in
  the restart commit and rechecked there.

Tool paths used were `/private/tmp/ply-p2b-api.4umBuM/bin/apidiff` and
`/private/tmp/ply-p2b-lint.SGWVGp/bin/golangci-lint`; probe before reuse because
temporary paths are not persistent dependencies. The clean gate used
`/private/tmp/ply-p3-clean-gate.v5OyEp` for writable caches.

Environment on 2026-08-24: host Go 1.26.2 on Darwin arm64, module Go 1.18,
`/bin/bash` 3.2.57, and PATH Bash 5.3.9. Docker, public network, cloud, Spring,
and distribution execution were outside this move.

## Stop Conditions

Stop and report rather than forcing progress when:

- the audit exits 2;
- a comparable ratchet regresses;
- public CLI or Go API compatibility cannot be established;
- the git process adapter cannot preserve exact argument order and defaults;
- the move requires Maven, cloud, HTTP, filesystem, clock, server, Docker,
  dependency, distribution, Spring, or mutation-harness scope; or
- the one coherent git adapter move is complete.
