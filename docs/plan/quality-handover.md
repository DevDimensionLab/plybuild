# Quality Upgrade Handover

Generated: 2026-08-24T14:01:12+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `f5ee37d74665`.
- Restart preparation base: `f5ee37d74665`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no implementation changes.
- No push, merge, release, publication, stash, revert, or worktree removal was
  performed.

P3 implementation commits:

```text
5c6f2fa quality: move process exits to main
03d6242 quality: move git processes behind adapter
f5ee37d quality: move Maven processes behind adapter
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main, Git,
and Maven process-adapter moves are complete. Cloud and the later adapter moves
remain active, and P4-P8 remain queued in the machine-readable plan block. The
active archive is
`docs/plan/agent-sessions/2026-08-24T140112+0200-cover-cloud-clone-seam.md`.
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

## P3 Move 3 Result

The third P3 move completed as one measured implementation change:

1. `maven.RunOn` retains its exported signature and returned callback while a
   private callback receives the complete `process.Dependencies` value.
   Production explicitly selects `process.System()` and still runs the same
   Maven, dot, and ktlint executable/argv in the project directory.
2. `pkg/maven/command.go` no longer calls `exec.Command` or `cmd.Run` directly.
   Its complete `process.Command` preserves argument values and order, project
   `Dir`, inherited environment/stdin, nil stderr, the unchanged
   `logger.StdOut()` writer, and the dependency's returned error.
3. Recording tests prove executable-before-every-argument ordering, full
   dependency values, complete command forwarding, project path, exact stdout
   writer identity, error propagation, a safe no-process default, and rejection
   of an empty asserted call population.
4. The inventory-bound `scripts/mutate-maven-sorting` path is a non-executable
   seam-test driver, not the P5 harness. Its meta-test runs all three Maven
   contracts and rejects an empty population, missing label, executability, or
   mutation declarations. Q1.4 rises without claiming Q2.1 coverage.

## Measured Quality State

The clean full audit at `f5ee37d74665` reports:

- Absolute L0: 8 of 8.
- 50 test functions, zero skipped; 11 of 23 packages have tests.
- Q0.8: 0 of 9 production scripts lack a meta-test.
- Q1.1: 12 of 23 packages have no tests, improved from the baseline 15.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 76 direct external sites outside the declared adapters of 78 production
  effect sites. `internal/adapter/process` is valid; the other four adapters are
  absent, so Q1.3 remains not comparable. Both direct Maven process effects
  left; the preserved logger stdout is exposed as one filesystem capability at
  the adapter call.
- Q1.4: 3 of 8 declared seams covered: `git-process`, `maven-process`, and
  `git-commit`.
- Q2.1: 0 of 8 subjects have a real executable harness. The upstream
  filename-only denominator sees two `mutate-*` paths, but the exact-path local
  validator records both as non-executable. Q2.2/Q2.3 therefore remain
  documented findings until the real harnesses are built in P5.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Q2.8 and Q2.9 remain honestly `UNMEASURABLE` pending criterion-bound manual
  evidence; executable magnitude, bad-input, and read-only controls exist.
- Full audit exit: expected 1 for 16 documented findings; no accepted audit
  returned exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is not
  comparable while four declared adapter paths are absent.
- Measurement identity: clean, with zero dirty paths.

The scanner reports the same five-improved / two-held comparable ratchets as
the preceding move because Q1.4 was already improved and advances from 2 of 8
to 3 of 8. Q1.3 remains not comparable until all adapter paths exist.

Regenerate `target/quality-audit/scorecard.json`; it is ignored output, not
persistent evidence. Compatibility checks create ignored reports under
`target/compatibility`; direct them outside the worktree or remove those
generated reports before the full audit.

## Decisions And Learned Facts

1. A function-valued process dependency makes the type-aware Q1.3 scanner fail
   closed. Use the `process.Runner` interface; do not replace it with a stored
   function without extending and meta-testing the scanner in separate scope.
2. Test recorders must write only through resolved in-memory concrete types.
   Generic `io.Writer` writes made the Q0.6 scanner fail closed, while asserted
   `*bytes.Buffer` destinations are recognized as memory-only.
3. Resolve temporary-directory symlinks before comparing an executed process's
   physical working directory. On macOS, `/tmp` and `/private/tmp` can name the
   same directory.
4. `.quality/inventory` is baseline-checksum-bound. Seam labels therefore remain
   at their declared driver paths; changing those paths invalidates baseline
   comparison.
5. The non-executable seam drivers keep exact Q2.1 at 0 of 8, but the pinned
   upstream tool counts filenames for its mutation denominator. Q2.2/Q2.3 will
   report findings until P5 turns those paths into real executable harnesses and
   adds their T1-T10 meta-controls.
6. `shell.run` historically returns an `Output` whose `Err` field is not filled
   from `cmd.Run()`. The Git move preserves that behavior rather than combining
   a seam refactor with an unrelated compatibility change.
7. The adapter's system path inherits environment and stdin, attaches the
   caller's writers, and assigns `Command.Dir`; empty `Dir` preserves the Git
   helpers' original process working directory.
8. Passing `logger.StdOut()` in the complete adapter command exposes the
   returned `*os.File` as a filesystem capability in the type-aware scan. This
   is the preserved stdout contract, not a remaining direct process effect;
   Q1.3 therefore moves from 77/79 to 76/78 rather than 75/77.
9. Passing `GOLANGCI_LINT` as a Make command-line variable propagates through
   `MAKEFLAGS` and defeats the lint meta-test's missing-binary mutant. Provide
   `APIDIFF` and `GOLANGCI_LINT` as environment variables for `make preflight`.
10. The sandbox denies default Go and golangci-lint cache paths. Use isolated
   `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` directories under
   `/private/tmp`.

## Next Objective

P3 remains active. Cover the declared `cloud-clone` argument swap in
`GitCloudConfig.Refresh` as one measured move. Begin red with a recording
boundary that proves the configured cloud URL stays before the target
directory on clone, an existing repository selects pull, complete values and
errors survive, defaults are inert, and empty asserted populations fail.

Preserve the exported `CloudConfig` interface, `GitCloudConfig`,
`OpenGitCloudConfig`, and all Refresh callers. Keep cache-first behavior: pull
when `<target>/.git` exists, clone otherwise, with the same URL, target, error
formatting, logging, and filesystem probe. Reuse the existing Git/process path;
do not duplicate Git argv or broaden into generic shell, Maven, HTTP,
filesystem-adapter, Docker, dependency, Spring, distribution, or P5 harness
work.

Bind only the immutable `cloud-clone` label through the declared non-executable
`scripts/mutate-config-cloud` driver and its Q0.8 meta-test. Expect Q1.2 to stay
zero, Q1.4 to rise from 3 of 8 to 4 of 8, exact Q2.1 to remain 0 of 8, and zero
comparable regressions. Accept only regenerated scanner values.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `f5ee37d74665`:

- Focused process-adapter, Maven argument, safe-default, error, stdout,
  working-directory, and empty-population contracts: PASS.
- Q1.2/Q1.3/Q1.4 focused audit: expected exit 1, exact 0, 76 of 78, and 3 of 8;
  clean measurement identity and zero regressions.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 50 tests, 11 tested packages,
  five improved, two held, zero regressed, one not comparable, and zero dirty
  paths.

Tool paths used were `/private/tmp/ply-p2b-api.4umBuM/bin/apidiff` and
`/private/tmp/ply-p2b-lint.SGWVGp/bin/golangci-lint`; probe before reuse because
temporary paths are not persistent dependencies.

Environment on 2026-08-24: host Go 1.26.2 on Darwin arm64, module Go 1.18,
`/bin/bash` 3.2.57, and PATH Bash 5.3.9. Docker, public network, cloud, Spring,
and distribution execution were outside this move.

## Stop Conditions

Stop and report rather than forcing progress when:

- an audit exits 2 and the adapter/test shape cannot be corrected in scope;
- a comparable ratchet regresses;
- public CLI or Go API compatibility cannot be established;
- the cloud seam cannot preserve exact clone/pull selection, URL, target,
  cache-first behavior, default, and error behavior;
- the move requires HTTP, filesystem-adapter, clock, server, Docker, dependency,
  distribution, Spring, or formal mutation-harness scope; or
- the one coherent cloud-clone seam move is complete.
