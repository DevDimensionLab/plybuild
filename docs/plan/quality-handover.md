# Quality Upgrade Handover

Generated: 2026-08-24T11:32:08+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `312d16897e52`.
- Restart preparation base: `312d16897e52`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no implementation changes.
- No push, merge, release, publication, stash, revert, or worktree removal was
  performed.

P2B implementation commit:

```text
312d168 build: make distribution local-only
```

## Continuity Checkpoint

`codex-dev-start.sh` is in `NEXT` state for P3. P2B is complete; P3 is active
and P4-P8 remain queued in the machine-readable plan block. The active archive
is
`docs/plan/agent-sessions/2026-08-24T113208+0200-remove-process-exits.md`.
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
so normal handoff does not require a test edit. The ignored `.agent-task` path
is absent and is not task authority.

## P2B Result

P2B completed its single measured move:

1. `test/makefile_distribution_test.sh` began red on the active top-level
   `brews:` key. Its final controls reject Homebrew formula, Homebrew cask, and
   Snap publisher sections in the default config, require remote release to be
   disabled, and reject a surviving standalone brew config.
2. `.goreleaser.yml` now declares configuration version 2, pins its schema link
   to the official v2.17.1 tag, has `release.disable: true`, and contains no
   inactive package-manager publisher. `.goreleaser.brews.yml` was removed.
3. `make snapshot` records exactly
   `release --snapshot --clean --skip=publish`; `make release` delegates to the
   same path. Five publication credential variables are cleared. The recording
   executable rejects leaked credentials, and `make release-brew` fails before
   the recorder is called.
4. The distribution contract is part of both `make test` and `make preflight`.
   GoReleaser and every publisher remained unexecuted; actual snapshot artifact
   acceptance is still honestly unverified until P6.

The flags and schema were checked against official GoReleaser v2.17.1 sources.
That version's command definition states that `--snapshot` implies skipping
announce, publish, and validate; the explicit `--skip=publish`, disabled remote
release config, absent publishers, and cleared credentials provide independent
fail-closed layers.

## Measured Quality State

The clean full audit at `312d16897e52` reports:

- Absolute L0: 8 of 8.
- 37 test functions, zero skipped; 7 of 22 packages have tests.
- Q0.8: 0 of 7 production scripts lack a meta-test.
- Ratchets: two improved, five held, zero regressed.
- Process-exiting calls outside `main`: 127.
- Direct production effects outside the five declared adapters: 80 of 80.
- Declared seam swap tests: 0 of 8.
- Mutation harnesses: 0 of 8.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Q2.8 and Q2.9 remain honestly `UNMEASURABLE` pending criterion-bound manual
  evidence; the executable magnitude, bad-input, and read-only controls exist.
- Full audit exit: expected 1 for 16 documented non-passing criteria; no audit
  returned exit 2.
- Measurement identity: clean, with zero dirty paths.

Regenerate `target/quality-audit/scorecard.json`; it is ignored output, not
persistent evidence. Compatibility checks create ignored reports under
`target/compatibility`; remove those generated reports before the full audit or
the audit will correctly mark the measurement tree dirty.

## Decisions And Learned Facts

1. GoReleaser is not installed on the host and was not invoked. The official
   v2.17.1 tagged schema and command source were inspected before production
   flags were bound.
2. There is no production publication target. `make release` is deliberately a
   local snapshot alias, while the inactive Homebrew target explains its
   boundary and exits non-zero.
3. Passing `GOLANGCI_LINT` as a Make command-line variable propagates through
   `MAKEFLAGS` and defeats the lint meta-test's missing-binary mutant. Provide
   `APIDIFF` and `GOLANGCI_LINT` as environment variables for `make preflight`.
4. The sandbox denies the default Go and golangci-lint cache paths. Use isolated
   `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` directories under
   `/private/tmp`.
5. P2A API, Cobra, subprocess, and four host acceptance contracts remain green.
   The distribution move left the process-exit and direct-effect ratchets held.
6. Both `main.go` and `cmd/ply/main.go` remain active entry points. P3 must bind
   them to one error/exit boundary before introducing adapters.

## Next Objective

P3 is active. Its first measured move is the exit-only-main boundary: start
with a failing contract that proves both executable entry points delegate to
one error-returning command path, move process termination to the two `main`
boundaries, and reduce Q1.2 process-exiting calls outside `main` from 127 to
zero while preserving exported `cmd.Execute()` and `cmd.RootCmd` signatures and
all P2A CLI/API/acceptance behavior.

Do not begin subprocess, HTTP, or filesystem adapter work until the process
boundary is green and measured. Do not enter Docker, dependencies, cloud,
Spring, distribution, or mutation scope. If external callers demonstrably
require `cmd.Execute()` itself to terminate the process, stop for an explicit
migration decision rather than hiding the exit.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed at the P2B boundary:

- Focused distribution contract and its Homebrew/Snap/credential/target
  controls: PASS.
- API and CLI compatibility plus both compatibility meta-tests: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- `bash .quality/tools/test-quality-audit.sh`: PASS, 15 controls.
- Clean full audit at `312d16897e52`: expected exit 1, L0 8 of 8, two
  improved, five held, zero regressed ratchets, and zero dirty paths.

Tool paths used at P2B were
`/private/tmp/ply-p2b-api.4umBuM/bin/apidiff` and
`/private/tmp/ply-p2b-lint.SGWVGp/bin/golangci-lint`; probe before reuse because
temporary paths are not persistent dependencies. The successful gate used
`/private/tmp/ply-p2b-gate.gLKunN` for writable caches.

Environment on 2026-08-24: host Go 1.26.2 on Darwin arm64, module Go 1.18,
`/bin/bash` 3.2.57, PATH Bash 5.3.9, Docker daemon unavailable, GoReleaser
unavailable, and shellcheck unavailable.

## Stop Conditions

Stop and report rather than forcing progress when:

- the audit exits 2;
- a comparable ratchet regresses;
- public CLI or Go API compatibility cannot be established;
- both executable entry points cannot share one explicit error/exit boundary
  without an approved compatibility migration;
- P3 move 1 would require adapter, Docker, cloud, Spring, dependency,
  distribution, or mutation scope; or
- the first measured P3 move is complete.
