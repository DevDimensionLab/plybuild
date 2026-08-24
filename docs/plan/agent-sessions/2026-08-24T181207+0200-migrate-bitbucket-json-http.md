# Agent Session: Migrate Bitbucket JSON HTTP

Status: ANSWERED - HISTORY
Session ID: `2026-08-24T181207+0200-migrate-bitbucket-json-http`
Created: `2026-08-24T18:12:07+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `9f6c6738f76f82bcd1737c2f1f478edc573a3b5e55089b451df549452e121b0a`
Previous: [2026-08-24T173323+0200-migrate-spring-json-http.md](2026-08-24T173323+0200-migrate-spring-json-http.md)
Next: [2026-08-24T192424+0200-migrate-kibana-post-http.md](2026-08-24T192424+0200-migrate-kibana-post-http.md)
Outcome: P3 move 9 completed at 789ae23: Bitbucket bearer JSON GET now uses the HTTP adapter, Q1.1 is 9 of 25, Q1.3 is 64 of 74, Q1.4 remains 7 of 8, exact Q2.1 remains 0 of 8, and the clean gate passed with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the next P3 move in one measured change: route the coherent Bitbucket
bearer-token discovery JSON GET flow through the existing HTTP adapter. Keep
Q1.2 at zero, reduce only migrated Q1.3 sites, hold Q1.4 at 7 of 8 and exact
Q2.1 at 0 of 8, and preserve every P2A contract with zero comparable ratchet
regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, `pkg/bitbucket` request/response types and every caller/
test, `pkg/http/client.go` token JSON behavior and every caller/test,
`internal/adapter/httpclient` and all anonymous/basic-auth/redirect contracts,
the adjacent anonymous Spring JSON and Wget seams without changing them,
Bitbucket synchronization and clone/pull follow-ups, relevant command callers,
the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 scanners, and the P2A compatibility,
subprocess, and host acceptance contracts. Regenerate ignored reports outside
the measured tree or remove them before the clean audit. Implementation commit
`dee214c` has 101 tests across 15 of 25 packages, Q0.6 has 17 guarded
safe-writer sites and zero unsafe direct test writes, Q1.2 is 0, Q1.3 is 66
violations of 76 production effect sites with clock and server absent, and
Q1.4 is 7 of 8. Exact Q2.1 remains 0 of 8 executable harnesses; the upstream
filename-only denominator sees five non-executable `mutate-*` paths. The clean
gate passed, the full audit exited 1 for 16 documented findings and never 2,
and comparable ratchets were five improved, two held, and zero regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Bitbucket bearer-token JSON move. Do not push,
merge, publish, remove the worktree, stash inherited changes, revert user work,
or run destructive Git commands. Do not invoke a publisher or distribution
command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve exported
`http.GetJsonWithAccessToken`, `bitbucket.With`, and `bitbucket.QueryRepos`
signatures and every caller. Preserve private `Bitbucket.queryProjects`, exact
host-plus-path URLs and limits, standard fresh-client GET and redirect behavior,
the `Authorization: Bearer <token>` and `Content-Type: application/json`
headers, response body read and close behavior, ignored read errors, absence of
status rejection, debug logging, JSON unmarshal target and errors,
synchronization order and warnings, and later clone/pull follow-ups.

Migrate only request execution used by the token JSON helper and the coherent
Bitbucket query callers. Extend `httpclient.Request` only as narrowly as needed
for complete bearer JSON GET values. Use resolvable interfaces and complete
dependency values, not stored function dependencies. Preserve every existing
anonymous and basic-auth adapter caller. Do not duplicate generic request
execution or broaden into XML, anonymous JSON, Wget/Wpost, Kibana,
filesystem/clone/pull, Git, cloud, clock, server, Docker, dependencies,
distribution, or formal mutation-harness scope. Do not create a new P3 seam
driver or claim a new seam or P5 harness.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`, all of
`pkg/bitbucket` and every caller/test, `pkg/http/client.go` and every token JSON
caller/test, `internal/adapter/httpclient` and its tests, the completed Spring
JSON/download contracts and seam-driver pattern, relevant synchronization and
clone/pull follow-ups, the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 implementation in
`.quality/tools`, and the P2A compatibility and host acceptance contracts
before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording HTTP adapter, token JSON, and Bitbucket query
   contracts. Prove complete host-plus-path URLs, GET selection, bearer and
   content-type headers, parsed `ProjectList` and `ProjectRepos` content,
   dependency errors after complete request delivery, response read/close,
   ignored read errors, no status rejection, JSON unmarshal errors, standard
   redirects, safe defaults, complete dependency delivery, query selection,
   and rejection of an empty recorded population.

2. Extend `httpclient.Request` with the thinnest complete bearer JSON GET
   representation and route only `GetJsonWithAccessToken` request execution
   through `httpclient.Dependencies`. Preserve its distinct response lifecycle
   instead of adopting status rejection or read-error propagation. Give the
   coherent Bitbucket query flow the thinnest private complete-dependency
   boundary needed by recording tests. Leave anonymous/basic-auth behavior,
   Spring/Wget, clone/pull, synchronization follow-ups, and all excluded flows
   unchanged. Do not add or relabel a seam driver.

3. Run focused Bitbucket/HTTP/httpclient tests,
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 measurements, API/CLI compatibility, all four
   host acceptance flows, the full checkpoint gate, and empty-HOME count-2
   from a clean commit. Expect Q1.2 to stay zero, Q1.4 to stay 7 of 8, and
   exact Q2.1 to stay 0 of 8. Accept only regenerated Q1.1 and Q1.3 values. The
   full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Bitbucket token
JSON move or record an exact resumable state. Rewrite the rolling handover,
update the P3 measurements, answer this archive, create one linked NEXT archive,
replace the launcher's mutable regions, run the launcher contract, and make the
separate handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
