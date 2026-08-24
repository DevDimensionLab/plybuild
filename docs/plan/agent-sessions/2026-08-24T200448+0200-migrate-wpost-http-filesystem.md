# Agent Session: Migrate Wpost Through HTTP And Filesystem Adapters

Status: ANSWERED - HISTORY
Session ID: `2026-08-24T200448+0200-migrate-wpost-http-filesystem`
Created: `2026-08-24T20:04:48+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `d4292c7572e7c114a1e02acdcfbc8f1c7ae604a40e1a34df8a7bf163939e778f`
Previous: [2026-08-24T192424+0200-migrate-kibana-post-http.md](2026-08-24T192424+0200-migrate-kibana-post-http.md)
Next: [2026-08-24T212408+0200-supervise-noninteractive-sessions.md](2026-08-24T212408+0200-supervise-noninteractive-sessions.md)
Outcome: P3 move 11 completed at `a7eb3ef`: Wpost now routes its exact form POST request, file creation, and body copy through the HTTP and filesystem adapters, Q1.1 is 9 of 25, Q1.3 is 58 of 68, Q1.4 remains 7 of 8, exact Q2.1 remains 0 of 8, and the clean gate passed with zero comparable ratchet regressions.

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the next P3 move in one measured change: route the unused exported
`pkg/http.Wpost` form-download request-to-file lifecycle through the existing
HTTP and filesystem adapters. Keep Q1.2 at zero, reduce only migrated Q1.3
sites, hold Q1.1 at 9 of 25, Q1.4 at 7 of 8 and exact Q2.1 at 0 of 8, and
preserve every P2A contract with zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, all of `pkg/http` and every caller/test, the current
`internal/adapter/httpclient` POST/GET implementation and every caller/test,
`internal/adapter/filesystem` and every caller/test, adjacent Wget and Spring
download contracts, the completed Kibana POST flow, Go's `http.PostForm`
semantics, the Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 scanners, and the P2A
compatibility, subprocess, and host acceptance contracts. Regenerate ignored
reports outside the measured tree or remove them before the clean audit.
Implementation commit `07ac6ce` has 123 tests across 16 of 25 packages, Q0.6
has 17 guarded safe-writer sites and zero unsafe direct test writes, Q1.1 is 9
of 25, Q1.2 is 0, and Q1.3 is 62 violations of 72 production effect sites with
clock and server absent. Q1.4 is 7 of 8 and exact Q2.1 is 0 of 8 executable
harnesses; the upstream filename-only denominator sees five non-executable
`mutate-*` paths. The clean gate passed, the full audit exited 1 for 16
documented findings and never 2, and comparable ratchets were five improved,
two held, and zero regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Wpost request-to-file move. Do not push, merge,
publish, remove the worktree, stash inherited changes, revert user work, or run
destructive Git commands. Do not invoke a publisher or distribution command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve exported
`http.Wpost(downloadUrl, filePath string, formData url.Values) error` and every
caller. Preserve the exact debug message and values, complete request URL,
`formData.Encode()` body, POST method, `application/x-www-form-urlencoded`
content type, standard `http.PostForm` client and redirect behavior,
response-before-file-create order, absence of status rejection, create/truncate
behavior, copy bytes and errors, file-before-response deferred close order, and
ignored close errors.

Migrate only Wpost's coherent request execution, file creation, and body copy
through complete `httpclient.Dependencies` and `filesystem.Dependencies`
values. Reuse the existing narrow POST and filesystem representations, making
only the smallest adapter representation change that characterization proves
necessary to preserve `http.PostForm` behavior. Use resolvable interfaces and
complete dependency values, not stored function dependencies. Preserve every
anonymous, basic-auth, bearer JSON, and Kibana POST adapter caller. Do not
duplicate generic request or filesystem execution or broaden into Wget,
Spring, retry/clock injection, parsing helpers, command writes, XML/JSON, Git,
cloud, server, Docker, dependencies, distribution, or formal mutation-harness
scope. Do not create or relabel a seam driver or claim a new seam or P5
harness.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`, all of
`pkg/http` and every caller/test, `internal/adapter/httpclient` and every
caller/test, `internal/adapter/filesystem` and every caller/test, adjacent Wget
and Spring download contracts, completed Kibana POST contracts, the standard
library `http.PostForm` behavior supported by the Go 1.18 API, the
Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 implementation in `.quality/tools`, and the P2A
compatibility and host acceptance contracts before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording HTTP/filesystem and Wpost contracts. Prove the
   complete URL, encoded form body, method, content type, debug log, standard
   client redirects, complete dependency delivery, request/create/copy errors
   after recording, no status rejection, response/file close order and ignored
   close errors, exact copied bytes, safe defaults, and rejection of empty
   recorded populations.

2. Give Wpost only the private complete-dependency boundary needed by recording
   tests and route its request, create, and copy effects through the existing
   adapters. Preserve the response and file lifecycle in Wpost. Extend the HTTP
   request value only if needed for exact PostForm semantics. Leave Wget, all
   GET paths, Kibana, retry/clock, parsing, commands, and other filesystem flows
   unchanged. Do not add or relabel a seam driver.

3. Run focused HTTP/httpclient/filesystem and adjacent Spring/Kibana tests,
   Q0.6/Q1.1/Q1.2/Q1.3/Q1.4/Q2.1 measurements, API/CLI compatibility, all four
   host acceptance flows, the full checkpoint gate, and empty-HOME count-2 from
   a clean commit. Expect Q1.1 to stay 9 of 25, Q1.2 to stay zero, Q1.4 to stay
   7 of 8, and exact Q2.1 to stay 0 of 8. Accept only the regenerated Q1.3
   value. The full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Wpost move or
record an exact resumable state. Rewrite the rolling handover, update the P3
measurements, answer this archive, create one linked NEXT archive, replace the
launcher's mutable regions, run the launcher contract, and make the separate
handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
