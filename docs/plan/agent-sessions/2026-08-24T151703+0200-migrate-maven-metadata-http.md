# Agent Session: Migrate Maven Metadata HTTP

Status: NEXT
Session ID: `2026-08-24T151703+0200-migrate-maven-metadata-http`
Created: `2026-08-24T15:17:03+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `fc99d15c250eec97129984279a00c650ffadbccd267840300ba3c04f02a8f926`
Previous: [2026-08-24T140112+0200-cover-cloud-clone-seam.md](2026-08-24T140112+0200-cover-cloud-clone-seam.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the next P3 move in one measured change: introduce the declared
`internal/adapter/httpclient` boundary for the coherent Maven metadata XML
request flow and kill the `maven-http` username/password argument swap. Keep
Q1.2 at zero, reduce only the migrated Q1.3 sites, and preserve every P2A
contract with zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, `pkg/maven/query.go` and its callers/tests, Maven
repository/auth types, `pkg/http/client.go` and every affected XML helper
caller, existing HTTP transport tests, the process and cloud seam patterns,
the Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 scanners, and the P2A compatibility, subprocess,
and host acceptance contracts. Regenerate ignored reports outside the measured
tree or remove them before the clean audit. Implementation commit `ee5e9ab`
has 55 tests across 11 of 23 packages, Q0.6 has zero unsafe direct test writes,
Q1.2 is 0, Q1.3 is 76 violations of 78 production effect sites with four
adapter paths absent, and Q1.4 is 4 of 8. Exact Q2.1 remains 0 of 8 executable
harnesses; the upstream filename-only denominator sees three non-executable
`mutate-*` paths. The clean gate passed, the full audit exited 1 for 16
documented findings and never 2, and comparable ratchets were five improved,
two held, and zero regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Maven metadata HTTP move. Do not push, merge,
publish, remove the worktree, stash inherited changes, revert user work, or run
destructive Git commands. Do not invoke a publisher or distribution command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve the exported
`Repository.GetMetaData`, `GetBannedModel`, `http.GetXml`, and
`http.GetAuthXml` signatures and every caller. Production must retain the
repository-derived metadata URL, anonymous GET and basic-auth behavior, XML
decoding, response status/body/close semantics, logging, and formatted errors.

Migrate only the network execution needed by the coherent anonymous/basic-auth
XML metadata flow. Use a resolvable interface and complete dependency value,
not stored function dependencies. Do not create a second generic HTTP
implementation or broaden into JSON/token requests, Wget/Wpost, Kibana,
filesystem, clock, server, Docker, dependency upgrades, Spring, distribution,
or formal mutation-harness scope.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`,
`pkg/maven/query.go`, `pkg/maven/query_test.go`, Maven repository/auth types and
all metadata callers, `pkg/http/client.go`, every affected XML helper caller,
the existing transport tests, the process/cloud interface and recording tests,
the Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 implementation in `.quality/tools`, and the P2A
compatibility and host acceptance contracts before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording anonymous and authenticated metadata contracts.
   Prove the complete repository-derived metadata URL, username before
   password, anonymous selection when auth is nil, parsed XML result,
   dependency errors, response behavior, safe defaults, complete dependency
   delivery, and rejection of an empty recorded population.

2. Introduce the thinnest zero-value-safe `internal/adapter/httpclient`
   interface boundary and route production through the existing HTTP/XML
   behavior without duplicating request construction or decoding. Preserve
   exact URL construction, GET/basic-auth selection, credential order, logging,
   status/body/close and error semantics. Bind the immutable `maven-http` label
   through the existing non-executable `scripts/mutate-maven-sorting` driver
   and its Q0.8 meta-test, without claiming the later P5 harness.

3. Run focused Maven/HTTP/adapter/argument-swap tests,
   Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 measurements, API/CLI compatibility, all four host
   acceptance flows, the full checkpoint gate, and empty-HOME count-2 from a
   clean commit. Expect Q1.2 to stay zero, `maven-http` to raise Q1.4 from 4 of
   8 to 5 of 8, and exact Q2.1 to remain 0 of 8. Accept only regenerated Q1.3
   values. The full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Maven metadata
HTTP move or record an exact resumable state. Rewrite the rolling handover,
update the P3 measurements, answer this archive, create one linked NEXT archive,
replace the launcher's mutable regions, run the launcher contract, and make the
separate handoff-only commit `docs: prepare next agent session`.

Keep P3 active for later adapter moves. Do not launch the next session.
COMPLETE is valid only after every authorized checkpoint through P8 is
complete.
<!-- CODEX_SESSION_PROMPT_END -->
