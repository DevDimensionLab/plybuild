# Quality Upgrade Handover

Generated: 2026-08-24T17:33:23+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `89d0f76cf5ba`.
- Restart preparation base: `89d0f76cf5ba`.
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
89d0f76 quality: move Spring download behind adapters
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main,
process-adapter, cloud-clone, Maven metadata HTTP, template-copy filesystem,
and Spring download moves are complete. Spring initializer discovery JSON is
the next coherent HTTP flow; later adapter moves and P4-P8 remain authorized
in the machine-readable plan block. The active archive is
`docs/plan/agent-sessions/2026-08-24T173323+0200-migrate-spring-json-http.md`.
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

## P3 Move 7 Result

The seventh P3 move completed as one measured implementation change:

1. `internal/adapter/filesystem` now also owns archive `Create` and body `Copy`
   operations behind its existing resolvable `FileSystem` interface and
   complete `Dependencies` value. Its zero value returns `ErrNoFilesystem`
   before create or copy; `System` retains `os.Create` create/truncate and
   `io.Copy` behavior.
2. `pkg/http.Wget` retains its exported signature. A private complete
   dependency value routes the anonymous request through `httpclient` and the
   archive create/copy through `filesystem`. Production retains `http.Get`
   redirect behavior, no status rejection, exact body bytes, response and file
   lifecycle, file-before-response deferred close order, and error precedence.
3. `spring.DownloadInitializer` retains its exported signature and resolves
   `os.Getwd` plus `time.Now().Unix()` before entering a private deterministic
   helper. One complete interface dependency records download, unzip, and
   delete selection without moving clock/current-directory effects, unzip
   internals, or archive deletion.
4. Thirteen adapter, HTTP, and Spring contracts prove complete encoded URL
   before complete archive path, response bytes, create/truncate, anonymous
   system transport, no status rejection, close order, dependency errors,
   logging, unzip/delete follow-ups, safe defaults, complete delivery,
   `spring-<Unix>.zip` naming, and rejection of an empty population.
5. The inventory-bound `scripts/mutate-spring` path is non-executable. Its Q0.8
   meta-test binds the immutable `spring-download` label to all 13 focused Go
   contracts without declaring P5 mutations.

## Measured Quality State

The clean full audit at `89d0f76cf5ba` reports:

- Absolute L0: 8 of 8.
- 92 test functions, zero skipped; 15 of 25 packages have tests.
- Q0.6: 17 guarded safe-writer call sites and zero unsafe direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 10 of 25 packages have no tests; `pkg/spring` joined the tested
  population.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 67 direct external sites outside the declared adapters of 77
  production effect sites. Clock and server remain absent, so Q1.3 is not
  comparable. Four direct Wget sites left and one system create site joined
  the filesystem adapter population.
- Q1.4: 7 of 8 declared seams covered: `git-process`, `maven-process`,
  `cloud-clone`, `maven-http`, `spring-download`, `template-copy`, and
  `git-commit`.
- Q2.1: 0 of 8 subjects have a real executable harness. The upstream
  filename-only denominator sees five non-executable `mutate-*` paths; the
  exact-path validator records all five as non-executable. Q2.2/Q2.3 remain
  documented findings until P5.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Q2.8 and Q2.9 remain honestly `UNMEASURABLE` pending criterion-bound manual
  evidence; executable magnitude, bad-input, and read-only controls exist.
- Full audit exit: expected 1 for 16 documented findings; no accepted audit
  returned exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is not
  comparable while two declared adapter paths are absent.
- Measurement identity: clean at tree `963b62a93fbf`, with zero dirty paths.

Reports were written under `/private/tmp`; no ignored quality or compatibility
artifact entered the measured tree. Keep future generated reports outside the
worktree or remove them before the clean audit.

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
5. The original Wget path contributed four outside-adapter sites: anonymous
   request, create, file close, and body copy. Moving all four and adding one
   `os.Create` system site produced the scanner's exact 67-of-77 result; do not
   substitute a forecast for regenerated values.
6. `httpclient.System` still selects the standard anonymous `http.Get` path,
   which preserves default transport and redirect behavior. Wget deliberately
   performs no response status check, including for 4xx/5xx bodies.
7. Wget registers the response close before the archive close, so Go's LIFO
   defers close the archive file before the response body. Close errors remain
   ignored exactly as before.
8. The private Spring helper accepts an already resolved archive path. This
   gives deterministic recording tests without migrating `os.Getwd` or
   `time.Now` before the clock phase.
9. The remaining anonymous `pkg/http.get` direct request serves only
   `GetJson`; its production callers are Spring `GetRoot` and
   `GetDependencies`, with `Validate` consuming the latter. That is the next
   narrow HTTP move.
10. Passing `logger.StdOut()` in the complete Maven command exposes the
    returned `*os.File` as a filesystem capability in Q1.3. This is preserved
    stdout behavior, not a direct process effect.
11. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
    `make preflight`; Make command-line values propagate through `MAKEFLAGS`
    and defeat the lint meta-test's missing-binary mutant.
12. The sandbox denies default Go and golangci-lint cache paths. Use isolated
    `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` under `/private/tmp`.

## Next Objective

P3 remains active. Route the coherent anonymous Spring initializer discovery
JSON flow through the existing HTTP adapter. Begin red with recording
contracts for the complete root and `/dependencies` URLs, anonymous request
selection, parsed `IoRootResponse` and `IoDependenciesResponse` content,
dependency and JSON/status/read errors, response closure, safe defaults,
complete dependency delivery, `Validate` selection, and rejection of an empty
population.

Preserve exported `spring.GetRoot`, `spring.GetDependencies`,
`spring.Validate`, and `http.GetJson` signatures and every caller. Production
must retain exact base URLs, standard anonymous GET transport and redirects,
status `>= 400` error text, response body read/close behavior, JSON unmarshal
semantics, validation order, invalid-dependency reporting, and build
follow-ups. Reuse `httpclient.Dependencies` and the shared response helper; do
not duplicate request or JSON behavior.

Migrate only the anonymous GET execution used by `GetJson` and the coherent
Spring discovery callers. Do not broaden into token/authenticated HTTP, XML,
Wget/Wpost, Kibana, filesystem/archive/delete/demo cleanup, unzip internals,
Git, cloud, clock, server, Docker, dependencies, distribution, or formal
mutation-harness work. Q1.4 should remain 7 of 8 and exact Q2.1 should remain
0 of 8; accept only regenerated Q1.3 values.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `89d0f76cf5ba`:

- Focused Spring/HTTP/filesystem-adapter URL/path, bytes, create/truncate,
  lifecycle, error-order, safe-default, follow-up, system, naming, and
  empty-population contracts: PASS (13 named tests).
- Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 focused audit: expected exit 1, zero unsafe test
  writes, 0 exits, 67 of 77 effects, 7 of 8 seams, and 0 of 8 executable
  harnesses.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 92 tests, 15 tested packages,
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
- the Spring JSON flow cannot preserve exact URLs, anonymous transport,
  redirects, status/body/close, JSON, validation, errors, and build follow-ups;
- the move requires token/authenticated HTTP, XML, Wget/Wpost, Kibana,
  filesystem/archive/delete/demo cleanup, unzip internals, Git, cloud, clock,
  server, Docker, dependency, distribution, or formal mutation-harness scope;
  or
- the one coherent Spring initializer JSON move is complete.
