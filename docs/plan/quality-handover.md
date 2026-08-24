# Quality Upgrade Handover

Generated: 2026-08-24T20:04:48+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `07ac6ce809be`.
- Restart preparation base: `07ac6ce809be`.
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
dee214c quality: move Spring discovery behind HTTP adapter
789ae23 quality: move Bitbucket JSON behind HTTP adapter
07ac6ce quality: move Kibana POST behind HTTP adapter
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main,
process-adapter, cloud-clone, Maven metadata HTTP, template-copy filesystem,
Spring download, Spring discovery JSON, Bitbucket bearer JSON, and Kibana POST
moves are complete. The unused exported `Wpost` form-download lifecycle is the
next coherent adapter move, and P4-P8 remain authorized in the machine-readable
plan block. The active archive is
`docs/plan/agent-sessions/2026-08-24T200448+0200-migrate-wpost-http-filesystem.md`.
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

## P3 Move 10 Result

The tenth P3 move completed as one measured implementation change:

1. `httpclient.Request` now has one narrow `POST` value containing complete
   body bytes and headers. Anonymous requests retain standard `http.Get`, and
   basic-auth and bearer JSON retain their existing fresh-client GET paths.
   POST presence selects a fresh client, `http.MethodPost`, the supplied body,
   exact headers, and standard redirects.
2. Exported `kibana.POST`, `KibanaFetchRequest`, `KibanaResponse`, and
   `ExecuteKibanaQuery` signatures remain unchanged, as does private
   `internalPOST`. Only its request execution moved through
   `httpclient.Dependencies`; the exact URL, first `"size"` rewrite to 500,
   body, four headers, response read/close behavior, absence of status
   rejection, newline split, and JSON targets/errors remain owned by Kibana.
3. A private zero-value-safe query interface receives each complete
   `KibanaFetchRequest` through initial and recursive callers. Production
   delegates to the exported legacy `POST` boundary; zero-hit retry selection,
   output, and the 15-second sleep remain unchanged.
4. Eleven new top-level contracts raise the suite from 112 to 123 tests. They
   prove complete delivery, POST method/body/headers, redirects, non-success
   parsing, dependency/read/unmarshal errors, closure, safe defaults, recursive
   propagation, deterministic retry structure, and rejection of empty
   recorded populations.
5. No inventory entry, seam driver, or mutation label changed. Q1.4 therefore
   remains 7 of 8 and exact Q2.1 remains 0 of 8.

## Measured Quality State

The clean full audit at `07ac6ce809be` reports:

- Absolute L0: 8 of 8.
- 123 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 17 guarded safe-writer call sites and zero unsafe direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 62 direct external sites outside the declared adapters of 72
  production effect sites. Clock and server remain absent, so Q1.3 is not
  comparable. Both direct Kibana request-execution sites left the population.
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
- Measurement identity: clean at tree `9dc934ae917a`, with zero dirty paths.

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
5. `httpclient.POST` distinguishes POST presence even for an empty body and
   carries complete body and header values. It does not alter the anonymous,
   basic-auth, or bearer JSON GET selection paths.
6. Kibana `internalPOST` deliberately keeps its response lifecycle separate
   from the shared XML/anonymous helper: it parses non-success bodies, returns
   read errors, splits the first two newline records, and unmarshals distinct
   header/result targets.
7. The private Kibana query boundary uses idiomatic response-before-error
   ordering while the exported `POST` signature retains its legacy
   error-before-response order. Production adapts between them without storing
   a function dependency.
8. Retry coverage records syntax and selection without executing the legacy
   15-second sleep. Clock migration remains a separate P4 concern.
9. Moving `http.NewRequest` and fresh-client `Do` out of `internalPOST` removes
   exactly two violations and two production effect sites, yielding Q1.3 at 62
   of 72.
10. The remaining direct HTTP path is the unused exported `Wpost` flow. Its
    scanner population is one HTTP plus three filesystem effects; migrate that
    coherent request-to-file lifecycle without changing `Wget` or Kibana.
11. Passing `logger.StdOut()` in the complete Maven command exposes the
    returned `*os.File` as a filesystem capability in Q1.3. This is preserved
    stdout behavior, not a direct process effect.
12. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
    `make preflight`; Make command-line values propagate through `MAKEFLAGS`
    and defeat the lint meta-test's missing-binary mutant.
13. The sandbox denies default Go and golangci-lint cache paths. Use isolated
    `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` under `/private/tmp`.
14. Preflight regenerates ignored compatibility reports under `target`; remove
    only those generated files before the authoritative clean audit or direct
    report output outside the measured tree.

## Next Objective

P3 remains active. Route only the unused exported `http.Wpost` form-download
lifecycle through `internal/adapter/httpclient` and
`internal/adapter/filesystem`. Begin red with recording adapter and Wpost
contracts for the complete URL, exact encoded form body and content type,
standard `http.PostForm` redirect/client behavior, debug log, response and file
close order, request/create/copy error order, no status rejection, safe
defaults, complete dependency delivery, and rejection of empty populations.

Preserve the exported `Wpost` signature and every caller, the exact debug
message, form encoding, POST selection, content type, response-before-create
lifecycle, create/truncate behavior, file-before-response deferred close order,
ignored close errors, copy result, and lack of status rejection. Reuse the
existing narrow POST and filesystem values, extending request representation
only if characterization proves it is needed to preserve `http.PostForm`
behavior.

Do not change `Wget`, anonymous/basic-auth/bearer GET, Kibana, Spring, Maven,
Bitbucket, XML/JSON parsing, retry/clock, command writes, Git, cloud, server,
Docker, dependencies, distribution, or formal mutation-harness scope. Add no
seam driver or label. Q1.1 should remain 9 of 25, Q1.4 should remain 7 of 8,
and exact Q2.1 should remain 0 of 8; accept only the regenerated Q1.3 value.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `07ac6ce809be`:

- Red evidence: focused Kibana/HTTP-adapter tests failed only on the absent
  POST request value, dependency-aware internal POST, and query boundary.
- Focused Kibana/HTTP/httpclient plus adjacent HTTP, Spring, Maven, and
  Bitbucket contracts: PASS.
- Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 focused audit: expected exit 1, 17 guarded
  writes, zero unsafe test writes, 9 untested packages, 0 exits, 62 of 72
  effects, 7 of 8 seams, and 0 of 8 executable harnesses.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 123 tests, 16 tested packages,
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
- `Wpost` cannot preserve its exact URL, encoded form body, POST/content type,
  standard `http.PostForm` client and redirects, debug log, response/create/copy
  error order, no status rejection, create/truncate, close order, and ignored
  close errors;
- the move requires `Wget`, migrated GET or Kibana changes, clock/retry work,
  parsing repair, command writes, XML/JSON, Git, cloud, server, Docker,
  dependency, distribution, or formal mutation-harness scope; or
- the one coherent `Wpost` request-to-file move is complete.
