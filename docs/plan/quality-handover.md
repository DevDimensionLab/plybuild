# Quality Upgrade Handover

Generated: 2026-08-24T16:53:52+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `204e2223b5e3`.
- Restart preparation base: `204e2223b5e3`.
- Session head: use `git rev-parse --short=12 HEAD` after launch; the restart
  commit contains this handover and no implementation changes.
- No push, merge, release, publication, stash, revert, or worktree removal was
  performed.

P3 implementation commits:

```text
5c6f2fa quality: move process exits to main
03d6242 quality: move git processes behind adapter
f5ee37d quality: move Maven processes behind adapter
ee5e9ab quality: cover cloud clone seam
ffc4e77 quality: move Maven metadata HTTP behind adapter
204e222 quality: move template copy behind filesystem adapter
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main,
process-adapter, cloud-clone, Maven metadata HTTP, and template-copy filesystem
moves are complete. Spring download is the next coherent seam; later
filesystem moves and P4-P8 remain authorized in the machine-readable plan
block. The active archive is
`docs/plan/agent-sessions/2026-08-24T165352+0200-migrate-spring-download.md`.
Its predecessor is answered history, and the connected graph has exactly one
`NEXT` tail.

The launcher validates its attached branch/root, authorized queue, strict
session metadata, non-symlink planning inputs, reciprocal archive graph,
historical prompt digests, exact launcher/archive prompt bytes, and external
Codex executable. Mutable prompt data remains comment-encoded after the stable
execution boundary. A dirty worktree produces a static warning; no status
content enters the prompt.

`test/codex_dev_start_test.sh` derives active metadata from the launcher. Its
49 controls include premature `COMPLETE` rejection while authorized work
remains, so normal handoff requires no test edit. The ignored `.agent-task`
path is absent and is not task authority.

## P3 Move 6 Result

The sixth P3 move completed as one measured implementation change:

1. `internal/adapter/filesystem` now owns the four filesystem operations needed
   by a copied template file behind one resolvable `FileSystem` interface and
   complete `Dependencies` value. Its zero value returns `ErrNoFilesystem`
   before mutation; `Exists` treats the zero value as absent so the next safe
   operation returns that error. `System` is the only generic production
   implementation.
2. `pkg/file` retains the exported `CopyOrMerge` and `CopyFile` signatures.
   Private dependency-bearing helpers preserve destination selection, source
   read, destination-directory existence and legacy second probe, 0755
   creation, source mode lookup, source-before-destination logs, and final
   write. Existing-target merge internals, public generic helpers, and Git hook
   callers were not moved.
3. `pkg/template` uses a private complete-dependency interface whose production
   implementation delegates to the existing `file.CopyOrMerge`. Recording
   contracts prove the complete template source path stays before the complete
   resolved target path, dependency errors follow complete delivery, safe
   defaults do not write, and an empty recorded population fails.
4. File and adapter contracts prove missing-copy versus existing-merge
   selection, bytes and mode, missing-directory creation, exact operation and
   error order, source-before-destination logs, complete forwarding, zero-value
   safety, and real system behavior. The historical behavior in which a failed
   `MkdirAll` returns the preceding not-exist error is characterized exactly.
5. The inventory-bound `scripts/mutate-template` path is non-executable. Its
   Q0.8 meta-test binds the immutable `template-copy` label to all 14 focused Go
   contracts without declaring P5 mutations.

## Measured Quality State

The clean full audit at `204e2223b5e3` reports:

- Absolute L0: 8 of 8.
- 79 test functions, zero skipped; 14 of 25 packages have tests.
- Q0.6: 15 guarded safe-writer call sites and zero unsafe direct test writes.
- Q0.8: 0 of 11 production scripts lack a meta-test.
- Q1.1: 11 of 25 packages have no tests, the same debt count across the new
  adapter package.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 71 direct external sites outside the declared adapters of 80
  production effect sites. Filesystem is valid; clock and server remain absent,
  so Q1.3 is not comparable. The three direct copy sites left `pkg/file`, and
  the adapter contributes four in-boundary system sites.
- Q1.4: 6 of 8 declared seams covered: `git-process`, `maven-process`,
  `cloud-clone`, `maven-http`, `template-copy`, and `git-commit`.
- Q2.1: 0 of 8 subjects have a real executable harness. The upstream
  filename-only denominator sees four non-executable `mutate-*` paths; the
  exact-path validator records all four as non-executable. Q2.2/Q2.3 remain
  documented findings until P5.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Q2.8 and Q2.9 remain honestly `UNMEASURABLE` pending criterion-bound manual
  evidence; executable magnitude, bad-input, and read-only controls exist.
- Full audit exit: expected 1 for 16 documented findings; no accepted audit
  returned exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is not
  comparable while two declared adapter paths are absent.
- Measurement identity: clean at tree `f8da59562ca9`, with zero dirty paths.

Regenerate `target/quality-audit/scorecard.json`; it is ignored output, not
persistent evidence. Compatibility checks create ignored reports under
`target/compatibility`; direct them outside the worktree or remove regenerated
reports before the full audit.

## Decisions And Learned Facts

1. A function-valued effect dependency makes the type-aware Q1.3 scanner fail
   closed. Use resolvable interfaces and complete dependency values rather than
   stored functions.
2. Test recorders must write only through recognized in-memory concrete types
   or guarded `internal/testutil` helpers. Generic `io.Writer` writes and direct
   fixture mutations make Q0.6 fail closed.
3. `.quality/inventory` is baseline-checksum-bound. Seam labels therefore stay
   at their declared driver paths; changing those paths invalidates comparison.
4. Non-executable seam drivers keep exact Q2.1 at 0 of 8, while the pinned
   upstream tool counts their filenames. Q2.2/Q2.3 remain findings until P5
   builds real harnesses and T1-T10 controls.
5. Copy ordering is more specific than read/stat/write: the legacy path probes
   the target, reads the source, probes the destination directory twice when
   absent, creates it, stats the source mode, logs both paths, and writes. The
   failed-`MkdirAll` branch returns the prior not-exist error.
6. Only three old Q1.3 sites disappear for this move. Directory existence and
   creation use new private dependency-bearing helpers so the shared public
   `Exists` and `CreateDirectory` sites and their unrelated callers stay
   outside the migration.
7. An implicit directory in `testing/fstest.MapFS` can carry non-writable mode
   bits when copied literally. For a single writable fixture, use an existing
   `t.TempDir` plus `testutil.WriteFileOutsideWorkingTree`, or declare directory
   modes explicitly.
8. The Spring initializer uses `http.Wget`, not `Wpost`: it builds a GET URL as
   `<base>/starter.zip?<formData.Encode()>`, writes the response without a
   status-code check, then unzips and deletes the archive. Preserve that exact
   transport and lifecycle in the next move.
9. Spring archive naming uses the current working directory and
   `time.Now().Unix()`. Clock injection belongs to P4; a private helper can
   accept the already resolved archive path for deterministic download tests
   without migrating the clock early.
10. Passing `logger.StdOut()` in the complete Maven command exposes the returned
    `*os.File` as a filesystem capability in Q1.3. This is preserved stdout
    behavior, not a direct process effect.
11. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
    `make preflight`; Make command-line values propagate through `MAKEFLAGS` and
    defeat the lint meta-test's missing-binary mutant.
12. The sandbox denies default Go and golangci-lint cache paths. Use isolated
    `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` under `/private/tmp`.

## Next Objective

P3 remains active. Extend the existing HTTP and filesystem adapters through the
coherent Spring initializer download-to-archive flow and cover the declared
`spring-download` URL/archive-path swap. Begin red with recording contracts for
the complete encoded download URL before the complete resolved archive path,
response bytes reaching the archive, create/truncate and close behavior,
dependency errors and order, safe defaults, complete dependency delivery, and
rejection of an empty population.

Preserve exported `spring.DownloadInitializer` and `http.Wget` signatures and
every caller. Production must retain archive-path resolution and naming,
`formData.Encode()` query construction, the anonymous GET transport, lack of a
new status-code rejection, response and file close behavior, logging, error
order, unzip selection and errors, archive deletion and errors, and all later
build follow-ups. Use a private helper that accepts an already resolved archive
path if deterministic Spring recording tests need it; do not move the clock.

Route only the anonymous download request and archive creation/body-copy effects
through the existing adapters. Do not broaden into JSON/XML/authenticated HTTP,
`Wpost`, template/file merge internals, archive-path clock/current-directory
effects, unzip process internals, archive deletion, demo cleanup, Git, cloud,
Kibana, server, Docker, dependencies, distribution, or formal P5 harness work.
Bind only the immutable `spring-download` label through the declared
non-executable `scripts/mutate-spring` path and a Q0.8 meta-test. Expect Q1.2 to
stay zero, Q1.4 to rise from 6 of 8 to 7 of 8, exact Q2.1 to remain 0 of 8,
and zero comparable regressions; accept only regenerated Q1.3 values.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `204e2223b5e3`:

- Focused template/file/filesystem-adapter selection, path-order, bytes/mode,
  directory, logging, error-order, safe-default, system, and empty-population
  contracts: PASS (14 named tests).
- Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 focused audit: expected exit 1, zero unsafe test
  writes, 0 exits, 71 of 80 effects, 6 of 8 seams, and 0 of 8 executable
  harnesses.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 79 tests, 14 tested packages,
  five improved, two held, zero regressed, one not comparable, and zero dirty
  paths.

Tool paths used were `/private/tmp/ply-p2b-api.4umBuM/bin/apidiff` and
`/private/tmp/ply-p2b-lint.SGWVGp/bin/golangci-lint`; probe before reuse because
temporary paths are not persistent dependencies.

Environment on 2026-08-24: host Go 1.26.2 on Darwin arm64, module Go 1.18,
`/bin/bash` 3.2.57, and PATH Bash 5.3.9. Docker, public network, real cloud,
real Spring, and distribution execution were outside this move.

## Stop Conditions

Stop and report rather than forcing progress when:

- an audit exits 2 and the adapter/test shape cannot be corrected in scope;
- a comparable ratchet regresses;
- public CLI or Go API compatibility cannot be established;
- the Spring download seam cannot preserve exact URL/path order, bytes,
  create/truncate and close behavior, transport, logs, errors, unzip, deletion,
  and build follow-ups;
- the move requires JSON/XML/authenticated HTTP, `Wpost`, template/file merge
  internals, clock/current-directory migration, unzip process internals,
  archive deletion, demo cleanup, Git, cloud, Kibana, server, Docker,
  dependency, distribution, or formal mutation-harness scope; or
- the one coherent Spring initializer download move is complete.
