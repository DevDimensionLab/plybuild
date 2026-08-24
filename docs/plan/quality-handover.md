# Quality Upgrade Handover

Generated: 2026-08-24T13:02:35+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `03d6242de38d`.
- Restart preparation base: `03d6242de38d`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no implementation changes.
- No push, merge, release, publication, stash, revert, or worktree removal was
  performed.

P3 implementation commits:

```text
5c6f2fa quality: move process exits to main
03d6242 quality: move git processes behind adapter
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main and Git
process-adapter moves are complete. Maven, cloud, and the later adapter moves
remain active, and P4-P8 remain queued in the machine-readable plan block. The
active archive is
`docs/plan/agent-sessions/2026-08-24T130235+0200-migrate-maven-process-flow.md`.
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

## P3 Move 2 Result

The second P3 move completed as one measured implementation change:

1. `internal/adapter/process` now owns operating-system execution. A complete
   `Command` carries name, argument vector, working directory, stdout, and
   stderr. A `Runner` interface keeps the scanner resolvable, the zero-value
   `Dependencies` is an inert default, and production explicitly selects
   `System()`.
2. `pkg/shell/git.go` routes clone, pull, init, add, and commit through the
   adapter without changing exported signatures. Production retains the same
   `git` argv, empty process `Dir`, `-C` target-directory arguments, output
   buffers, and legacy error behavior. Git dirty/repository probes continue to
   use generic `shell.Run`; Maven and other process sites were not migrated.
3. Recording tests capture complete arguments and working directories, prove
   clone URL-before-target and commit target-before-message ordering, preserve
   the complete message value and shell output, prove the safe zero-value
   default performs no Git mutation, and reject an empty recorded population.
   The adapter also executes an explicit system-path contract.
4. The inventory-bound `scripts/mutate-file-shell` path is a non-executable
   seam-test driver, not the P5 mutation harness. Its meta-test runs all five
   Git recording contracts and rejects an empty test population, missing
   labels, executability, or mutation declarations. This raises Q1.4 without
   claiming Q2.1 coverage.

## Measured Quality State

The clean full audit at `03d6242de38d` reports:

- Absolute L0: 8 of 8.
- 47 test functions, zero skipped; 11 of 23 packages have tests.
- Q0.8: 0 of 8 production scripts lack a meta-test.
- Q1.1: 12 of 23 packages have no tests, improved from the baseline 15.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 77 direct external sites outside the declared adapters of 79 production
  effect sites. `internal/adapter/process` is valid; the other four adapters are
  absent, so Q1.3 remains not comparable. The three direct Git sites left the
  violation population and the adapter owns two process effects.
- Q1.4: 2 of 8 declared seams covered: `git-process` and `git-commit`.
- Q2.1: 0 of 8 subjects have a real executable harness. The upstream
  filename-only denominator sees one `mutate-*` path, but the exact-path local
  validator records it as non-executable. Q2.2/Q2.3 therefore remain documented
  findings until the real file-shell harness is built in P5.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Q2.8 and Q2.9 remain honestly `UNMEASURABLE` pending criterion-bound manual
  evidence; executable magnitude, bad-input, and read-only controls exist.
- Full audit exit: expected 1 for 16 documented findings; no accepted audit
  returned exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is not
  comparable while four declared adapter paths are absent.
- Measurement identity: clean, with zero dirty paths.

The earlier four-improved / three-held forecast described Q1.4 at 0 of 8.
Covering the two required Git seams necessarily moves Q1.4 from held to
improved, yielding the scanner's exact five-improved / two-held result.

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
4. `.quality/inventory` is baseline-checksum-bound. The Git seam labels therefore
   remain at their declared `scripts/mutate-file-shell` path; changing the path
   invalidates baseline comparison.
5. The non-executable seam driver keeps Q2.1 at 0 of 8, but the pinned upstream
   tool counts filenames for its mutation denominator. Q2.2/Q2.3 will report
   findings until P5 turns that path into the real executable harness and adds
   its T1-T10 meta-controls.
6. `shell.run` historically returns an `Output` whose `Err` field is not filled
   from `cmd.Run()`. The Git move preserves that behavior rather than combining
   a seam refactor with an unrelated compatibility change.
7. The adapter's system path inherits environment and stdin, attaches the
   caller's writers, and assigns `Command.Dir`; empty `Dir` preserves the Git
   helpers' original process working directory.
8. Passing `GOLANGCI_LINT` as a Make command-line variable propagates through
   `MAKEFLAGS` and defeats the lint meta-test's missing-binary mutant. Provide
   `APIDIFF` and `GOLANGCI_LINT` as environment variables for `make preflight`.
9. The sandbox denies default Go and golangci-lint cache paths. Use isolated
   `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` directories under
   `/private/tmp`.

## Next Objective

P3 remains active. Reuse `internal/adapter/process` for one coherent Maven
subprocess flow in `pkg/maven/command.go` and cover the declared
`maven-process` swap. Begin red with a recorder that proves the executable stays
before all arguments, the project path remains the process working directory,
the logger stdout writer reaches the process unchanged, the full dependency
value reaches the production path, defaults remain inert, and empty recorded
populations fail.

Preserve exported `maven.RunOn` and its returned callback signature, error
behavior, exact argv, working directory, and stdout wiring. Migrate only the two
direct process effects in `pkg/maven/command.go`; do not migrate cloud, generic
shell, diagrams, profile/editor, browser-opening, HTTP, filesystem, clock, or
server paths in the same move. Bind the immutable inventory label through the
declared `scripts/mutate-maven-sorting` path without claiming the later P5
mutation harness.

Expect Q1.2 to stay zero, the two Maven violations to leave Q1.3, and Q1.4 to
rise from 2 of 8 to 3 of 8. Accept only regenerated scanner values and zero
comparable regressions.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `03d6242de38d`:

- Focused process-adapter, Git argument, safe-default, output, working-directory,
  and empty-population contracts: PASS.
- Q1.2/Q1.3/Q1.4 focused audit: expected exit 1, exact 0, 77 of 79, and 2 of 8;
  clean measurement identity and zero regressions.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 47 tests, 11 tested packages,
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
- the Maven adapter move cannot preserve exact executable, argument, directory,
  stdout, default, and error behavior;
- the move requires cloud, HTTP, filesystem, clock, server, Docker, dependency,
  distribution, Spring, or formal mutation-harness scope; or
- the one coherent Maven process move is complete.
