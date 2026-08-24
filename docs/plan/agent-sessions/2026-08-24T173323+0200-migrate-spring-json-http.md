# Agent Session: Migrate Spring JSON HTTP

Status: ANSWERED - HISTORY
Session ID: `2026-08-24T173323+0200-migrate-spring-json-http`
Created: `2026-08-24T17:33:23+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `5a898a58d4900dc83f8747b2723dfeb23358552db28de2d5508555c7d964f5bf`
Previous: [2026-08-24T165352+0200-migrate-spring-download.md](2026-08-24T165352+0200-migrate-spring-download.md)
Next: [2026-08-24T181207+0200-migrate-bitbucket-json-http.md](2026-08-24T181207+0200-migrate-bitbucket-json-http.md)
Outcome: P3 move 8 completed at dee214c: anonymous Spring discovery JSON now uses the HTTP adapter and shared response lifecycle, Q1.3 is 66 of 76, Q1.4 remains 7 of 8, exact Q2.1 remains 0 of 8, and the clean gate passed with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the next P3 move in one measured change: route the coherent anonymous
Spring initializer discovery JSON GET flow through the existing HTTP adapter.
Keep Q1.2 at zero, reduce only migrated Q1.3 sites, hold Q1.4 at 7 of 8 and
exact Q2.1 at 0 of 8, and preserve every P2A contract with zero comparable
ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, `pkg/spring/io.go`, Spring discovery/validation types and
every caller/test, `pkg/http/client.go` JSON and shared response helpers and
every caller/test, `internal/adapter/httpclient` and its recording/system
contracts, the adjacent Spring download seam without changing it, relevant
build follow-ups, the Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 scanners, and the P2A
compatibility, subprocess, and host acceptance contracts. Regenerate ignored
reports outside the measured tree or remove them before the clean audit.
Implementation commit `89d0f76` has 92 tests across 15 of 25 packages, Q0.6
has 17 guarded safe-writer sites and zero unsafe direct test writes, Q1.2 is 0,
Q1.3 is 67 violations of 77 production effect sites with clock and server
absent, and Q1.4 is 7 of 8. Exact Q2.1 remains 0 of 8 executable harnesses; the
upstream filename-only denominator sees five non-executable `mutate-*` paths.
The clean gate passed, the full audit exited 1 for 16 documented findings and
never 2, and comparable ratchets were five improved, two held, and zero
regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Spring initializer discovery JSON move. Do not
push, merge, publish, remove the worktree, stash inherited changes, revert user
work, or run destructive Git commands. Do not invoke a publisher or
distribution command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve exported
`spring.GetRoot`, `spring.GetDependencies`, `spring.Validate`, and
`http.GetJson` signatures and every caller. Production must retain the exact
base root and `<base>/dependencies` URLs, standard anonymous GET transport and
redirect behavior, status `>= 400` error text, response body read and close
behavior, JSON unmarshal target and errors, dependency-validation order,
invalid-dependency reporting, and later build follow-ups.

Migrate only the anonymous request execution used by `GetJson` and the coherent
Spring discovery callers. Use resolvable interfaces and complete dependency
values, not stored function dependencies. Reuse the existing HTTP adapter and
shared response lifecycle without duplicating generic request or JSON
behavior. Do not broaden into token/authenticated HTTP, XML, Wget/Wpost,
Kibana, filesystem/archive/delete/demo cleanup, unzip internals, Git, cloud,
clock, server, Docker, dependencies, distribution, or formal mutation-harness
scope.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`,
`pkg/spring/io.go`, Spring types and every discovery/validation caller/test,
`pkg/http/client.go`, every JSON/shared-response caller/test,
`internal/adapter/httpclient` and its tests, the completed Spring download
contracts and seam-driver pattern, the Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 implementation
in `.quality/tools`, and the P2A compatibility and host acceptance contracts
before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording Spring discovery and HTTP JSON contracts. Prove
   complete root and `/dependencies` URLs, anonymous request selection, parsed
   `IoRootResponse` and `IoDependenciesResponse` content, dependency errors
   after complete request delivery, status/read/unmarshal errors and response
   closure, safe defaults, complete dependency delivery, validation selection
   and errors, and rejection of an empty recorded population.

2. Route only `GetJson` request execution through `httpclient.Dependencies` and
   the shared response lifecycle, then give the coherent Spring discovery flow
   the thinnest private complete-dependency boundary needed by recording tests.
   Preserve JSON decoding and validation behavior without changing the
   completed Wget/download flow. Extend the existing non-executable
   `scripts/mutate-spring` driver and its Q0.8 meta-test only with the focused
   contracts; keep its immutable `spring-download` label single and do not
   claim a new seam or P5 harness.

3. Run focused Spring/HTTP/httpclient tests, Q0.6/Q1.2/Q1.3/Q1.4/Q2.1
   measurements, API/CLI compatibility, all four host acceptance flows, the
   full checkpoint gate, and empty-HOME count-2 from a clean commit. Expect
   Q1.2 to stay zero, Q1.4 to stay 7 of 8, and exact Q2.1 to stay 0 of 8.
   Accept only regenerated Q1.3 values. The full audit may exit 1 for
   documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Spring JSON move
or record an exact resumable state. Rewrite the rolling handover, update the P3
measurements, answer this archive, create one linked NEXT archive, replace the
launcher's mutable regions, run the launcher contract, and make the separate
handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
