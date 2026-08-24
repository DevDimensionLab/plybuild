# Quality Upgrade Handover

Generated: 2026-08-24T19:24:24+02:00

This is a rolling handover. Rewrite it at each checkpoint; do not append a
session diary.

## Repository

- Worktree: `/Users/perottochristensen/github/ply/upgrade-quality`
- Branch: `codex/upgrade-quality`
- Base: `master` at `5635d50`
- Measured implementation head: `789ae23af28b`.
- Restart preparation base: `789ae23af28b`.
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
```

## Continuity Checkpoint

`codex-dev-start.sh` remains in `NEXT` state for P3. The exit-only-main,
process-adapter, cloud-clone, Maven metadata HTTP, template-copy filesystem,
Spring download, Spring discovery JSON, and Bitbucket bearer JSON moves are
complete. The command-reachable Kibana POST flow is the next coherent HTTP
move; the unused `Wpost` lifecycle stays separate, and P4-P8 remain authorized
in the machine-readable plan block. The active archive is
`docs/plan/agent-sessions/2026-08-24T192424+0200-migrate-kibana-post-http.md`.
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

## P3 Move 9 Result

The ninth P3 move completed as one measured implementation change:

1. `httpclient.Request` now has one narrow `BearerJSON` access-token value.
   Anonymous requests still select standard `http.Get`; basic-auth requests
   retain their existing fresh-client path. Bearer JSON selects a fresh-client
   GET with exact authorization and JSON content-type headers and standard
   redirects.
2. Exported `http.GetJsonWithAccessToken` retains its signature and exact
   `json.Unmarshal(body, &response)` target. Only request execution moved
   through `httpclient.Dependencies`; no status rejection was introduced, read
   errors remain ignored, and debug logging plus body closure remain unchanged.
3. Exported `bitbucket.With` and `bitbucket.QueryRepos` plus private
   `Bitbucket.queryProjects` retain their signatures. A private zero-value-safe
   interface receives each complete request and decode destination. The exact
   project/repository URLs and limits, parsed types, lowercase repository
   selection, synchronization order, warning behavior, and clone/pull
   follow-ups are preserved.
4. Eleven new top-level contracts raise the suite from 101 to 112 tests and add
   `pkg/bitbucket` to the tested population. They prove complete delivery,
   GET/headers, redirects, non-success parsing, ignored read errors, closure,
   logging, errors, safe defaults, Bitbucket response content, synchronization
   selection and warnings, and rejection of an empty population.
5. No inventory entry, seam driver, or mutation label changed. Q1.4 therefore
   remains 7 of 8 and exact Q2.1 remains 0 of 8.

## Measured Quality State

The clean full audit at `789ae23af28b` reports:

- Absolute L0: 8 of 8.
- 112 test functions, zero skipped; 16 of 25 packages have tests.
- Q0.6: 17 guarded safe-writer call sites and zero unsafe direct test writes.
- Q0.8: 0 of 12 production scripts lack a meta-test.
- Q1.1: 9 of 25 packages have no tests.
- Q1.2: 0 process-exiting calls outside `main`.
- Q1.3: 64 direct external sites outside the declared adapters of 74
  production effect sites. Clock and server remain absent, so Q1.3 is not
  comparable. Both direct token request sites left the population.
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
- Measurement identity: clean at tree `25e44ff84b21`, with zero dirty paths.

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
5. A `BearerJSON` request value is sufficient to distinguish presence even for
   an empty access token and to preserve the paired authorization/content-type
   headers without broadening anonymous or basic-auth values.
6. `GetJsonWithAccessToken` deliberately does not use the shared XML/anonymous
   response helper: its legacy contract parses non-success bodies and ignores
   read errors before JSON decoding.
7. The Bitbucket query boundary receives whole `httpclient.Request` and decode
   destination values. `SynchronizeAllRepos` uses that same dependency for
   projects and lowercase repository queries before unchanged clone/pull work.
8. Moving `http.NewRequest` and fresh-client `Do` out of the token helper
   removes exactly two violations and two production effect sites, yielding
   Q1.3 at 64 of 74.
9. The remaining direct HTTP paths are the unused exported `Wpost` flow (one
   HTTP plus three filesystem effects) and command-reachable Kibana POST (two
   HTTP effects). Keep their request/response and file lifecycles separate.
10. Kibana `internalPOST` rewrites the first `"size"` value to 500, sends four
    exact headers, parses the first two newline-separated JSON records, and is
    retried by `POST` after a 15-second sleep when no hits arrive. Migrate only
    request execution next; clock/retry behavior is separate.
11. Passing `logger.StdOut()` in the complete Maven command exposes the
    returned `*os.File` as a filesystem capability in Q1.3. This is preserved
    stdout behavior, not a direct process effect.
12. Provide `APIDIFF` and `GOLANGCI_LINT` as environment variables for
    `make preflight`; Make command-line values propagate through `MAKEFLAGS`
    and defeat the lint meta-test's missing-binary mutant.
13. The sandbox denies default Go and golangci-lint cache paths. Use isolated
    `GOCACHE`, `GOTMPDIR`, and `GOLANGCI_LINT_CACHE` under `/private/tmp`.

## Next Objective

P3 remains active. Route only the command-reachable Kibana POST request
execution through `internal/adapter/httpclient`. Begin red with recording
adapter and Kibana contracts for exact POST URL, rewritten body, all four
headers, complete dependency delivery, dependency/read/JSON errors after
request recording, no status rejection, body closure, parsed header/result
content, standard fresh-client redirects, safe defaults, caller selection, and
rejection of an empty recorded population.

Preserve exported `kibana.POST`, `KibanaFetchRequest`, `KibanaResponse`, and
`ExecuteKibanaQuery` signatures and every caller, plus private `internalPOST`.
Preserve the first size rewrite to 500, fresh-client POST behavior, exact
headers, body bytes, response close/read/error behavior, no status rejection,
newline split and JSON targets/errors, zero-hit retry selection and output,
recursive interval queries, and command file/output follow-ups. Extend
`httpclient.Request` only as narrowly as needed for a complete POST value and
preserve every anonymous/basic-auth/bearer GET caller.

Do not broaden into `Wpost`, filesystem, retry sleep/clock injection, response
format repair, parsing helpers, result filtering, command output writes, XML,
other JSON, Git, cloud, server, Docker, dependencies, distribution, or formal
mutation-harness work. Add no seam driver or label. Q1.1 should remain 9 of 25,
Q1.4 should remain 7 of 8, and exact Q2.1 should remain 0 of 8; accept only the
regenerated Q1.3 value.

## Start

The active P3 task can be started from any directory:

```sh
/Users/perottochristensen/github/ply/upgrade-quality/codex-dev-start.sh
```

Do not start through `.agent-task/current.md`; it is not task authority. Every
session prepares its successor automatically before stopping.

## Verification Notes

Completed from clean implementation commit `789ae23af28b`:

- Red evidence: focused Bitbucket/HTTP/adapter tests failed only on the absent
  bearer request value, token dependency helper, and Bitbucket query boundary.
- Focused Bitbucket/HTTP/httpclient plus adjacent Spring/Maven contracts: PASS;
  the unchanged Spring driver still runs all 22 named contracts.
- Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 focused audit: expected exit 1, 17 guarded
  writes, zero unsafe test writes, 9 untested packages, 0 exits, 64 of 74
  effects, 7 of 8 seams, and 0 of 8 executable harnesses.
- API and CLI compatibility plus root/status/upgrade/build and unknown-command
  subprocess surfaces: PASS.
- Host install, status, upgrade, and build acceptance: PASS.
- `make preflight`, `make test`, `make test-install`, `make test-agent-start`,
  uncached tests, race tests, and `go vet ./...`: PASS.
- Empty-HOME `go test ./... -count=2` with isolated writable state: PASS.
- Audit meta-suite: PASS twice, 15 controls and all 228 baseline numeric leaves.
- Clean full audit: expected exit 1, L0 8 of 8, 112 tests, 16 tested packages,
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
- the Kibana POST flow cannot preserve exact URL, rewritten body, headers,
  fresh-client and redirect behavior, body read/close, no status rejection,
  newline/JSON behavior, retry selection, callers, and later command follow-ups;
- the move requires `Wpost`, filesystem, clock/retry migration, parsing repair,
  result filtering, command writes, XML, other JSON, Git, cloud, server, Docker,
  dependency, distribution, or formal mutation-harness scope; or
- the one coherent Kibana POST request-execution move is complete.
