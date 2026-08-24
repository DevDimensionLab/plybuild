# Agent Session: Migrate Kibana POST HTTP

Status: ANSWERED - HISTORY
Session ID: `2026-08-24T192424+0200-migrate-kibana-post-http`
Created: `2026-08-24T19:24:24+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `fc1140ee2109e89d0c61779f00b7aaf664e7e4d65a86146fc3d72fe4f1834db1`
Previous: [2026-08-24T181207+0200-migrate-bitbucket-json-http.md](2026-08-24T181207+0200-migrate-bitbucket-json-http.md)
Next: [2026-08-24T200448+0200-migrate-wpost-http-filesystem.md](2026-08-24T200448+0200-migrate-wpost-http-filesystem.md)
Outcome: P3 move 10 completed at `07ac6ce`: Kibana POST request execution now uses the HTTP adapter, Q1.1 is 9 of 25, Q1.3 is 62 of 72, Q1.4 remains 7 of 8, exact Q2.1 remains 0 of 8, and the clean gate passed with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the next P3 move in one measured change: route the coherent,
command-reachable Kibana POST request execution through the existing HTTP
adapter. Keep Q1.2 at zero, reduce only migrated Q1.3 sites, hold Q1.1 at 9 of
25, Q1.4 at 7 of 8 and exact Q2.1 at 0 of 8, and preserve every P2A contract
with zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, all of `pkg/kibana` and every caller/test, relevant
`cmd/plugin_diagrams.go` follow-ups, `internal/adapter/httpclient` and every
anonymous/basic-auth/bearer/redirect contract, adjacent `pkg/http` JSON, Wget,
and unused Wpost flows without changing them, the Kibana response types, retry
and recursive interval-query callers, the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1
scanners, and the P2A compatibility, subprocess, and host acceptance contracts.
Regenerate ignored reports outside the measured tree or remove them before the
clean audit. Implementation commit `789ae23` has 112 tests across 16 of 25
packages, Q0.6 has 17 guarded safe-writer sites and zero unsafe direct test
writes, Q1.1 is 9 of 25, Q1.2 is 0, and Q1.3 is 64 violations of 74 production
effect sites with clock and server absent. Q1.4 is 7 of 8 and exact Q2.1 is 0
of 8 executable harnesses; the upstream filename-only denominator sees five
non-executable `mutate-*` paths. The clean gate passed, the full audit exited 1
for 16 documented findings and never 2, and comparable ratchets were five
improved, two held, and zero regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Kibana POST request-execution move. Do not push,
merge, publish, remove the worktree, stash inherited changes, revert user work,
or run destructive Git commands. Do not invoke a publisher or distribution
command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve exported `kibana.POST`,
`kibana.KibanaFetchRequest`, `kibana.KibanaResponse`, and
`kibana.ExecuteKibanaQuery` signatures and every caller, plus private
`internalPOST`. Preserve the exact request URL, first `"size"` rewrite to 500,
POST selection and body bytes, `accept-language`, `authorization`,
`content-type`, and `kbn-version` headers, standard fresh-client redirect
behavior, response body read and close behavior, propagated read errors,
absence of status rejection, newline split, JSON unmarshal targets and errors,
zero-hit retry selection and output, recursive interval queries, and later
command output-file follow-ups.

Migrate only request execution used by `internalPOST` and its coherent Kibana
callers. Extend `httpclient.Request` only as narrowly as needed for complete
POST method/body/header values. Use resolvable interfaces and complete
dependency values, not stored function dependencies. Preserve every existing
anonymous, basic-auth, and bearer JSON GET adapter caller. Do not duplicate
generic request execution or broaden into Wpost, filesystem, retry sleep/clock
injection, response format repair, parsing helpers, filtering, command writes,
XML, other JSON, Git, cloud, server, Docker, dependencies, distribution, or
formal mutation-harness scope. Do not create or relabel a seam driver or claim
a new seam or P5 harness.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`, all of
`pkg/kibana` and every caller/test, relevant `cmd/plugin_diagrams.go` command
and output follow-ups, `internal/adapter/httpclient` and every caller/test, the
adjacent HTTP JSON/Wget/Wpost contracts, Kibana retry and recursive query
behavior, the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 implementation in `.quality/tools`,
and the P2A compatibility and host acceptance contracts before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording HTTP adapter and Kibana POST contracts. Prove the
   complete URL, POST method, rewritten body and four headers, complete
   dependency delivery, dependency errors after recording, non-success body
   parsing, response read/close and read errors, header/result JSON content and
   unmarshal errors, standard redirects, safe defaults, caller selection, and
   rejection of an empty recorded population. Keep retry coverage deterministic
   without waiting 15 seconds or migrating time.

2. Extend `httpclient.Request` with the thinnest complete POST representation
   and route only `internalPOST` request execution through
   `httpclient.Dependencies`. Preserve its distinct response parsing lifecycle
   and all higher-level Kibana selection. Give the coherent Kibana caller flow
   only the private complete-dependency boundary needed by recording tests.
   Leave every GET path, Wpost, retry clock, parsing/filtering, filesystem, and
   command follow-up unchanged. Do not add or relabel a seam driver.

3. Run focused Kibana/httpclient and adjacent HTTP tests,
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 measurements, API/CLI compatibility, all four
   host acceptance flows, the full checkpoint gate, and empty-HOME count-2 from
   a clean commit. Expect Q1.1 to stay 9 of 25, Q1.2 to stay zero, Q1.4 to stay
   7 of 8, and exact Q2.1 to stay 0 of 8. Accept only the regenerated Q1.3
   value. The full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Kibana POST move
or record an exact resumable state. Rewrite the rolling handover, update the P3
measurements, answer this archive, create one linked NEXT archive, replace the
launcher's mutable regions, run the launcher contract, and make the separate
handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
