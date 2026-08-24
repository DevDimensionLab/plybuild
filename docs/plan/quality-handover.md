# Quality Upgrade Handover

Generated: 2026-08-24T16:18:00+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `ffc4e7796ad3`.
- Restart preparation base: `ffc4e7796ad3`.
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
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main,
process-adapter, cloud-clone, and Maven metadata HTTP moves are complete.
Filesystem and the later adapter moves remain active, and P4-P8 remain queued
in the machine-readable plan block. The active archive is
`docs/plan/agent-sessions/2026-08-24T161800+0200-migrate-template-copy-filesystem.md`.
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

## P3 Move 5 Result

The fifth P3 move completed as one measured implementation change:

1. `internal/adapter/httpclient` now owns complete anonymous or basic-auth GET
   execution behind a resolvable `Client` interface and `Dependencies` value.
   Its zero value returns `ErrNoClient` without issuing a request; `System`
   preserves `http.Get` for anonymous requests and the legacy
   `NewRequest`/`SetBasicAuth`/fresh-client path for authenticated requests.
2. `pkg/http` retains the exported `GetXml` and `GetAuthXml` signatures, XML
   decoding, formatted status errors, body reads, and ignored close errors. A
   shared response helper avoids duplicating that behavior. JSON/token,
   Wget/Wpost, Kibana, filesystem, clock, and server paths were untouched.
3. Maven uses a private complete-dependency value whose production client
   delegates to those existing XML helpers. Five recording contracts prove the
   full repository-derived URL, nil-auth anonymous selection, username before
   password, parsed XML, dependency errors, safe defaults, complete destination
   delivery, and rejection of an empty recorded population.
4. Adapter and HTTP tests cover complete request forwarding, real anonymous and
   basic-auth transport, successful parsing, status and read errors, and body
   closure. The inventory-bound non-executable Maven seam driver now binds both
   immutable Maven labels to eight contracts without declaring P5 mutations.

## Measured Quality State

The clean full audit at `ffc4e7796ad3` reports:

- Absolute L0: 8 of 8.
- 65 test functions, zero skipped; 13 of 24 packages have tests.
- Q0.6: 12 guarded safe-writer call sites and zero unsafe direct test writes.
- Q0.8: 0 of 10 production scripts lack a meta-test.
- Q1.1: 11 of 24 packages have no tests, improved from the baseline 15.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 74 direct external sites outside the declared adapters of 79 production
  effect sites. `internal/adapter/process` and `internal/adapter/httpclient` are
  valid; filesystem, clock, and server remain absent, so Q1.3 is not comparable.
  The two authenticated XML effects left; three HTTP effects now live inside
  the adapter, while the direct anonymous helper remains for out-of-scope JSON.
- Q1.4: 5 of 8 declared seams covered: `git-process`, `maven-process`,
  `cloud-clone`, `maven-http`, and `git-commit`.
- Q2.1: 0 of 8 subjects have a real executable harness. The upstream
  filename-only denominator sees three non-executable `mutate-*` paths, but the
  exact-path local validator records all three as non-executable. Q2.2/Q2.3
  remain documented findings until P5.
- Acceptance scripts: 4 of 4; Q2.5, Q2.6, Q2.7, and Q2.10 pass.
- Q2.8 and Q2.9 remain honestly `UNMEASURABLE` pending criterion-bound manual
  evidence; executable magnitude, bad-input, and read-only controls exist.
- Full audit exit: expected 1 for 16 documented findings; no accepted audit
  returned exit 2.
- Comparable ratchets: five improved, two held, zero regressed; Q1.3 is not
  comparable while three declared adapter paths are absent.
- Measurement identity: clean, with zero dirty paths.

Regenerate `target/quality-audit/scorecard.json`; it is ignored output, not
persistent evidence. Compatibility checks create ignored reports under
`target/compatibility`; direct them outside the worktree or remove those
generated reports before the full audit.

## Decisions And Learned Facts

1. A function-valued process dependency makes the type-aware Q1.3 scanner fail
   closed. Use resolvable interfaces and complete dependency values rather than
   stored functions.
2. Test recorders must write only through resolved in-memory concrete types.
   Generic `io.Writer` writes make Q0.6 fail closed, while asserted
   `*bytes.Buffer` destinations are recognized as memory-only.
3. Resolve temporary-directory symlinks before comparing an executed process's
   physical working directory. On macOS, `/tmp` and `/private/tmp` can name the
   same directory.
4. `.quality/inventory` is baseline-checksum-bound. Seam labels therefore remain
   at their declared driver paths; changing those paths invalidates comparison.
5. Non-executable seam drivers keep exact Q2.1 at 0 of 8, but the pinned
   upstream tool counts filenames for its mutation denominator. Q2.2/Q2.3 stay
   findings until P5 builds real harnesses and their T1-T10 controls.
6. The adapter system path must preserve the legacy transport choices, not just
   produce equivalent XML: anonymous XML used `http.Get`, while authenticated
   XML built a GET request, set basic auth, and used a fresh `http.Client`.
7. The anonymous `get` helper is shared with JSON. Routing that helper wholesale
   through the adapter would broaden this move into JSON; the accepted split
   keeps JSON direct, moves only XML execution, and shares status/body/close
   processing to avoid behavior duplication.
8. Passing `logger.StdOut()` in the complete process command exposes the
   returned `*os.File` as a filesystem capability in Q1.3. This is preserved
   stdout behavior, not a direct process effect.
9. Passing `GOLANGCI_LINT` as a Make command-line variable propagates through
   `MAKEFLAGS` and defeats the lint meta-test's missing-binary mutant. Provide
   `APIDIFF` and `GOLANGCI_LINT` as environment variables for `make preflight`.
10. The sandbox denies default Go and golangci-lint cache paths. Use isolated
    `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` directories under
    `/private/tmp`.
11. Direct fixture writes are rejected even when the destination derives from
    `t.TempDir()`. Use guarded `internal/testutil` writers or in-memory values.

## Next Objective

P3 remains active. Introduce `internal/adapter/filesystem` through the coherent
non-existing-destination template copy flow and cover the declared
`template-copy` source/destination swap. Begin red with recording contracts for
the complete source path before the complete resolved target path, copy versus
merge selection, source bytes and mode, missing destination-directory creation,
dependency errors, safe defaults, complete dependency delivery, and rejection
of an empty population.

Preserve exported `template.MergeTemplate`, `file.CopyOrMerge`, and
`file.CopyFile` signatures and every caller. Production must retain the exact
template filtering and target-path derivation, existing-destination merge
branch, non-existing-destination copy behavior, source bytes/mode, directory
creation, logging, error order, search/replace, render, and Maven merge
follow-ups. Use resolvable interfaces and complete dependency values; route the
migrated OS calls through one thin adapter without duplicating generic file
behavior.

Bind only the immutable `template-copy` label through the declared
non-executable `scripts/mutate-template` path and a Q0.8 meta-test. Keep existing
destination merge implementations, unrelated write/delete/rename/grep/render
paths, Git hooks, cloud, Spring/download, HTTP, clock, server, Docker,
dependencies, distribution, and formal P5 harness work outside the move.
Expect Q1.2 to stay zero, Q1.4 to rise from 5 of 8 to 6 of 8, exact Q2.1 to
remain 0 of 8, and zero comparable regressions; accept only regenerated Q1.3
values.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `ffc4e7796ad3`:

- Focused Maven/HTTP/adapter, anonymous/basic-auth, argument-order,
  response-lifecycle, safe-default, error, and empty-population contracts: PASS.
- Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 focused audit: expected exit 1, zero unsafe test
  writes, 0 exits, 74 of 79 effects, 5 of 8 seams, and 0 of 8 executable
  harnesses.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 65 tests, 13 tested packages,
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
- the template-copy seam cannot preserve exact source/destination order,
  selection, bytes, mode, directory creation, logging, errors, and follow-ups;
- the move requires existing-destination merge internals, unrelated filesystem
  operations, Git hooks, cloud, Spring/download, HTTP, clock, server, Docker,
  dependency, distribution, or formal mutation-harness scope; or
- the one coherent template-copy filesystem move is complete.
