# Quality Upgrade Handover

Generated: 2026-08-24T18:12:07+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `dee214cf6650`.
- Restart preparation base: `dee214cf6650`.
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
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main,
process-adapter, cloud-clone, Maven metadata HTTP, template-copy filesystem,
Spring download, and Spring discovery JSON moves are complete. The Bitbucket
bearer-token JSON GET flow is the next coherent HTTP move; later adapter moves
and P4-P8 remain authorized in the machine-readable plan block. The active
archive is
`docs/plan/agent-sessions/2026-08-24T181207+0200-migrate-bitbucket-json-http.md`.
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

## P3 Move 8 Result

The eighth P3 move completed as one measured implementation change:

1. `http.GetJson` retains its exported signature and exact
   `json.Unmarshal(body, &parsed)` target. Its anonymous request now passes
   through `httpclient.Dependencies` and the existing shared status, body-read,
   and close lifecycle used by XML.
2. Production still selects standard anonymous `http.Get`, including default
   transport and redirects. Exact status `>= 400` error text, dependency/read/
   unmarshal errors, response closure, and JSON content are recorded.
3. `spring.GetRoot`, `spring.GetDependencies`, and `spring.Validate` retain
   their exported signatures. A private zero-value-safe interface receives the
   complete `httpclient.Request` and decode destination, preserving the base
   root and `<base>/dependencies` URLs, response types, validation selection,
   user-dependency order, invalid-dependency text, and later build follow-ups.
4. Nine new top-level contracts raise the suite from 92 to 101 tests. They
   prove complete anonymous requests and destinations, parsed root/dependency
   content, status/read/unmarshal/dependency errors, body closure, redirects,
   safe defaults, validation behavior, and rejection of an empty population.
5. The inventory-bound `scripts/mutate-spring` path remains non-executable. Its
   Q0.8 meta-test binds the single immutable `spring-download` label to all 22
   focused download and discovery contracts without declaring P5 mutations or
   a new seam.

## Measured Quality State

The clean full audit at `dee214cf6650` reports:

- Absolute L0: 8 of 8.
- 101 test functions, zero skipped; 15 of 25 packages have tests.
- Q0.6: 17 guarded safe-writer call sites and zero unsafe direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 10 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 66 direct external sites outside the declared adapters of 76
  production effect sites. Clock and server remain absent, so Q1.3 is not
  comparable. The remaining direct anonymous JSON request left the population.
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
- Measurement identity: clean at tree `440987584701`, with zero dirty paths.

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
5. `http.GetJson` now uses the same `getHTTPResponse` and `responseBody`
   lifecycle as XML. Retaining `json.Unmarshal(body, &parsed)` is intentional
   compatibility, even though `parsed` normally already holds a pointer.
6. `httpclient.System` still selects standard anonymous `http.Get`, and the new
   system contract records a two-request redirect with anonymous GET on both
   hops and closure of both bodies.
7. The Spring discovery boundary accepts complete request and destination
   values; `Validate` passes the dependency value whole into dependency
   discovery and makes no request for an empty dependency list.
8. Validation still iterates user dependencies in order and iterates the valid
   dependency map when formatting the valid-key list. Tests use one valid key
   when asserting exact invalid-dependency text, avoiding a false ordering
   contract on Go map iteration.
9. The remaining direct HTTP sites are the bearer-token JSON helper (two),
   unused `Wpost` download path (one HTTP plus three filesystem), and Kibana
   POST (two). The Bitbucket token flow has real callers and is the next narrow
   HTTP move; keep `Wpost` and Kibana separate.
10. Passing `logger.StdOut()` in the complete Maven command exposes the
    returned `*os.File` as a filesystem capability in Q1.3. This is preserved
    stdout behavior, not a direct process effect.
11. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
    `make preflight`; Make command-line values propagate through `MAKEFLAGS`
    and defeat the lint meta-test's missing-binary mutant.
12. The sandbox denies default Go and golangci-lint cache paths. Use isolated
    `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` under `/private/tmp`.

## Next Objective

P3 remains active. Route the coherent Bitbucket bearer-token JSON GET flow
through `internal/adapter/httpclient`. Begin red with recording contracts for
the complete host-plus-path URLs, GET method, `Authorization: Bearer <token>`
and `Content-Type: application/json` headers, parsed `ProjectList` and
`ProjectRepos` content, dependency errors after complete request delivery,
response body read/close behavior, ignored read errors, lack of status
rejection, JSON unmarshal errors, redirects, safe defaults, complete dependency
delivery, query selection, and rejection of an empty recorded population.

Preserve exported `http.GetJsonWithAccessToken`, `bitbucket.With`, and
`bitbucket.QueryRepos` signatures and every caller. Preserve private
`Bitbucket.queryProjects`, exact project/repository paths and limits, fresh
standard client redirect behavior, logging, response read/close and ignored
read errors, no status rejection, JSON target/errors, synchronization order,
warning behavior, and later clone/pull follow-ups. Extend `httpclient.Request`
only as narrowly as needed for complete bearer JSON GET values; preserve
anonymous and basic-auth behavior and all existing adapter callers.

Migrate only request execution for the token JSON helper and the coherent
Bitbucket query callers. Do not broaden into XML, anonymous JSON, Wget/Wpost,
Kibana, filesystem/clone/pull, Git, cloud, clock, server, Docker,
dependencies, distribution, or formal mutation-harness work. Q1.4 should
remain 7 of 8 and exact Q2.1 should remain 0 of 8; accept only regenerated
Q1.1 and Q1.3 values.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `dee214cf6650`:

- Red evidence: focused HTTP/Spring tests failed only on the absent `getJson`
  and Spring discovery dependency boundary.
- Focused HTTP/Spring/httpclient and seam-driver contracts: PASS; the driver
  ran all 22 named download/discovery contracts.
- Q0.6/Q0.8/Q1.2/Q1.3/Q1.4/Q2.1 focused audit: expected exit 1, zero unsafe
  test writes, 12 of 12 scripts paired, 0 exits, 66 of 76 effects, 7 of 8
  seams, and 0 of 8 executable harnesses.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS twice, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 101 tests, 15 tested packages,
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
- the token JSON flow cannot preserve exact URLs, GET and bearer/content-type
  headers, standard redirect behavior, body read/close and ignored read errors,
  no status rejection, JSON behavior, callers, and later clone/pull follow-ups;
- the move requires XML, anonymous JSON, Wget/Wpost, Kibana, filesystem,
  clone/pull, Git, cloud, clock, server, Docker, dependency, distribution, or
  formal mutation-harness scope; or
- the one coherent Bitbucket bearer-token JSON move is complete.
