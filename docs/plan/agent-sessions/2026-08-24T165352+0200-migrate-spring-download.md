# Agent Session: Migrate Spring Download

Status: NEXT
Session ID: `2026-08-24T165352+0200-migrate-spring-download`
Created: `2026-08-24T16:53:52+02:00`
Source: `codex-dev-start.sh`
Prompt SHA-256: `eecaf214fa80300a4e6eaefb583f697471edafca0e62216801eb98d9a33c9180`
Previous: [2026-08-24T161800+0200-migrate-template-copy-filesystem.md](2026-08-24T161800+0200-migrate-template-copy-filesystem.md)
Next: none
Outcome: pending

The block below is the byte-exact Codex prompt argument, including its terminal LF.

<!-- CODEX_SESSION_PROMPT_BEGIN -->
# Mission

Complete the next P3 move in one measured change: extend the existing HTTP and
filesystem boundaries through the coherent Spring initializer download-to-
archive flow and kill the `spring-download` URL/archive-path argument swap.
Keep Q1.2 at zero, reduce only migrated Q1.3 sites, and preserve every P2A
contract with zero comparable ratchet regressions.

# Authorized Roadmap

The user has authorized the ordered roadmap through P8. P3 is active, P4-P8
are queued in the machine-readable block in docs/plan/quality-upgrade.md, and
the launcher must remain NEXT until every authorized checkpoint is complete.

# Measurements At Start

Before editing, inspect branch, HEAD, status, the rolling handover, the P3 plan,
`.quality/inventory`, `pkg/spring/io.go` and every caller/type/test,
`pkg/http/client.go` download helpers and every caller/test, the existing
`internal/adapter/httpclient` and `internal/adapter/filesystem` interfaces and
recording/system tests, `pkg/shell` unzip behavior, relevant build follow-ups,
the Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 scanners, and the P2A compatibility, subprocess,
and host acceptance contracts. Regenerate ignored reports outside the measured
tree or remove them before the clean audit. Implementation commit `204e222`
has 79 tests across 14 of 25 packages, Q0.6 has zero unsafe direct test writes,
Q1.2 is 0, Q1.3 is 71 violations of 80 production effect sites with clock and
server absent, and Q1.4 is 6 of 8. Exact Q2.1 remains 0 of 8 executable
harnesses; the upstream filename-only denominator sees four non-executable
`mutate-*` paths. The clean gate passed, the full audit exited 1 for 16
documented findings and never 2, and comparable ratchets were five improved,
two held, and zero regressed.

# Role And Boundaries

Work autonomously in this worktree on `codex/upgrade-quality`. Make one focused
implementation commit for the Spring initializer download move. Do not push,
merge, publish, remove the worktree, stash inherited changes, revert user work,
or run destructive Git commands. Do not invoke a publisher or distribution
command.

Keep Go 1.18. Preserve `cmd.Execute()`, `cmd.ExecuteE()`, `cmd.RootCmd`, the
normalized Cobra tree, root/status/upgrade/build and unknown-command streams
and exits, and all four host acceptance flows. Preserve exported
`spring.DownloadInitializer` and `http.Wget` signatures and every caller.
Production must retain archive-path error timing and `spring-<Unix>.zip`
naming, the exact `<base>/starter.zip?<formData.Encode()>` URL, anonymous GET
transport and redirect behavior, the absence of a new HTTP status rejection,
archive create/truncate and response/file close behavior, body-copy bytes,
logging, error order and semantics, unzip selection and errors, archive
deletion and errors, and later build follow-ups.

Migrate only the anonymous download request and archive creation/body-copy
effects required by this coherent flow. Use resolvable interfaces and complete
dependency values, not stored function dependencies. A private helper may
accept the already resolved archive path for deterministic recording tests; do
not migrate `os.Getwd` or `time.Now` before P4. Do not broaden into JSON/XML or
authenticated HTTP, `Wpost`, template/file merge internals, unzip process
internals, archive deletion, demo cleanup, Git, cloud, Kibana, clock, server,
Docker, dependencies, distribution, or formal mutation-harness scope.

# Required Reading

Read docs/plan/quality-handover.md, the P3 section and checkpoint gate in
docs/plan/quality-upgrade.md, docs/design/quality-lift.md,
docs/design/agent-session-continuity.md, `.quality/inventory`,
`pkg/spring/io.go`, Spring types and every caller/test, `pkg/http/client.go` and
the Wget/Wpost/GET tests and callers, both existing adapter packages and tests,
the shell unzip and file-delete implementations without moving them, the
Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 implementation in `.quality/tools`, and the P2A
compatibility and host acceptance contracts before editing.

# Three Moves

These are three ordered steps inside the single measured P3 move.

1. Start red with recording Spring download and adapter contracts. Prove the
   complete encoded download URL stays before the complete resolved archive
   path, response bytes reach that path, create/truncate and close behavior is
   preserved, dependency errors retain their order and behavior, safe defaults
   perform no request or mutation, the complete dependency reaches the flow,
   unzip and delete follow-ups retain their selection and errors, and an empty
   recorded population fails.

2. Extend the existing adapters only as needed and route the migrated `Wget`
   request/create/body-copy effects through them. Preserve the anonymous
   transport, response/file lifecycle, no-status-check behavior, logging, and
   error order without duplicating Spring or generic file behavior. Bind the
   immutable `spring-download` label through a new non-executable
   `scripts/mutate-spring` seam driver and its Q0.8 meta-test, without claiming
   the later P5 harness.

3. Run focused Spring/HTTP/filesystem/argument-swap tests,
   Q0.6/Q1.2/Q1.3/Q1.4/Q2.1 measurements, API/CLI compatibility, all four host
   acceptance flows, the full checkpoint gate, and empty-HOME count-2 from a
   clean commit. Expect Q1.2 to stay zero, `spring-download` to raise Q1.4 from
   6 of 8 to 7 of 8, and exact Q2.1 to remain 0 of 8. Accept only regenerated
   Q1.3 values. The full audit may exit 1 for documented findings but never 2.

# Automatic Handoff

Before this agent session ends, finish and commit the coherent Spring download
move or record an exact resumable state. Rewrite the rolling handover, update
the P3 measurements, answer this archive, create one linked NEXT archive,
replace the launcher's mutable regions, run the launcher contract, and make the
separate handoff-only commit `docs: prepare next agent session`.

Keep P3 active until its measured exit is documented. Do not launch the next
session. COMPLETE is valid only after every authorized checkpoint through P8
is complete.
<!-- CODEX_SESSION_PROMPT_END -->
